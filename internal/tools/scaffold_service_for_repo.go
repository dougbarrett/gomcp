package tools

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterScaffoldServiceForRepo registers the scaffold_service_for_repo tool.
func RegisterScaffoldServiceForRepo(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_service_for_repo",
		Description: `Generate a service that wraps an existing custom repository.

Use this when you have a repository with custom methods (like the auth user repository)
and want to create a service layer that wraps those methods.

The tool:
1. Reads the existing repository interface from internal/repository/{domain}/{domain}.go
2. Generates a service with methods that wrap repository calls
3. Creates a separate package to avoid import cycles

Example: Create a user management service from the auth user repository:
  scaffold_service_for_repo: {
    service_name: "usermgmt",
    repository_domain: "user",
    exclude_methods: ["UpdateLastLogin", "UpdatePasswordHash"]
  }

This generates internal/services/usermgmt/ with a service that delegates to userrepo.Repository.

IMPORTANT: Unlike scaffold_domain, this tool does NOT automatically wire the service into main.go.
You must manually:
1. Import the service package in cmd/web/main.go
2. Instantiate the service with the repository
3. Wire it to any controllers that need it`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldServiceForRepoInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldServiceForRepo(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

// RepoMethod represents a method parsed from a repository interface.
type RepoMethod struct {
	Name       string
	Params     []RepoParam
	Returns    []string
	HasContext bool
}

// RepoParam represents a parameter of a repository method.
type RepoParam struct {
	Name string
	Type string
}

// ParsedRepository contains parsed repository information.
type ParsedRepository struct {
	Methods  []RepoMethod
	IsStruct bool // true if repository is a concrete struct, false if interface
}

func scaffoldServiceForRepo(registry *Registry, input types.ScaffoldServiceForRepoInput) (types.ScaffoldResult, error) {
	// Validate input
	if input.ServiceName == "" {
		return types.NewErrorResult("service_name is required"), nil
	}
	if input.RepositoryDomain == "" {
		return types.NewErrorResult("repository_domain is required"), nil
	}

	if err := utils.ValidateDomainName(input.ServiceName); err != nil {
		return types.NewErrorResult(fmt.Sprintf("service_name: %v", err)), nil
	}
	if err := utils.ValidateDomainName(input.RepositoryDomain); err != nil {
		return types.NewErrorResult(fmt.Sprintf("repository_domain: %v", err)), nil
	}

	// Get module path
	modulePath, err := utils.GetModulePath(registry.WorkingDir)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to get module path: %v", err)), nil
	}

	// Find repository file
	repoPkgName := utils.ToPackageName(input.RepositoryDomain)
	repoFilePath := filepath.Join(registry.WorkingDir, "internal", "repository", repoPkgName, repoPkgName+".go")
	if !utils.FileExists(repoFilePath) {
		return types.NewErrorResult(fmt.Sprintf("repository file not found: internal/repository/%s/%s.go", repoPkgName, repoPkgName)), nil
	}

	// Parse repository to extract methods
	parsed, err := parseRepositoryMethods(repoFilePath)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to parse repository: %v", err)), nil
	}

	if len(parsed.Methods) == 0 {
		return types.NewErrorResult("no methods found in repository"), nil
	}

	// Filter methods
	methods := filterMethods(parsed.Methods, input.IncludeMethods, input.ExcludeMethods)
	if len(methods) == 0 {
		return types.NewErrorResult("no methods remaining after filtering"), nil
	}

	// Create generator
	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Create service directory
	svcPkgName := utils.ToPackageName(input.ServiceName)
	svcDir := filepath.Join("internal", "services", svcPkgName)
	if err := gen.EnsureDir(svcDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate service code
	repoModelName := utils.ToModelName(input.RepositoryDomain)
	svcModelName := utils.ToModelName(input.ServiceName)
	repoImportAlias := utils.ToRepoImportAlias(input.RepositoryDomain)

	serviceCode := generateServiceCode(
		svcPkgName,
		svcModelName,
		modulePath,
		repoPkgName,
		repoImportAlias,
		repoModelName,
		methods,
		parsed.IsStruct,
	)

	// Write service file
	svcFilePath := filepath.Join(svcDir, svcPkgName+".go")
	if err := gen.WriteFile(svcFilePath, serviceCode); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to write service file: %v", err)), nil
	}

	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	nextSteps := []string{
		"go mod tidy",
		fmt.Sprintf("Wire %s service in cmd/web/main.go", svcPkgName),
		fmt.Sprintf("Add business logic to internal/services/%s/%s.go", svcPkgName, svcPkgName),
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create service '%s' wrapping '%s' repository with %d methods", input.ServiceName, input.RepositoryDomain, len(methods)),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created service '%s' wrapping '%s' repository with %d methods", input.ServiceName, input.RepositoryDomain, len(methods)),
		FilesCreated: result.FilesCreated,
		NextSteps:    nextSteps,
	}, nil
}

// parseRepositoryMethods parses a repository file and extracts methods.
// It supports both interface-based repositories and concrete struct repositories.
func parseRepositoryMethods(filePath string) (*ParsedRepository, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	result := &ParsedRepository{}
	var hasRepositoryStruct bool

	// First, check for Repository interface
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if typeSpec.Name.Name != "Repository" {
				continue
			}

			// Check if it's an interface
			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				methods := extractInterfaceMethods(interfaceType)
				if len(methods) > 0 {
					result.Methods = methods
					result.IsStruct = false
					return result, nil
				}
			}

			// Check if it's a struct (for concrete repository pattern)
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				hasRepositoryStruct = true
			}
		}
	}

	// If we found a Repository struct but no interface, extract methods from function declarations
	if hasRepositoryStruct {
		result.Methods = extractStructMethods(node)
		result.IsStruct = true
	}

	return result, nil
}

// extractInterfaceMethods extracts methods from an interface type.
func extractInterfaceMethods(interfaceType *ast.InterfaceType) []RepoMethod {
	var methods []RepoMethod

	for _, method := range interfaceType.Methods.List {
		if len(method.Names) == 0 {
			continue
		}

		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		repoMethod := extractMethodFromFuncType(method.Names[0].Name, funcType)
		methods = append(methods, repoMethod)
	}

	return methods
}

// extractStructMethods extracts methods from function declarations with *Repository receiver.
func extractStructMethods(node *ast.File) []RepoMethod {
	var methods []RepoMethod

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
			continue
		}

		// Check if this method is on *Repository
		for _, recv := range funcDecl.Recv.List {
			recvType := exprToString(recv.Type)
			if recvType == "*Repository" {
				// Skip unexported methods
				if funcDecl.Name.Name == "" || !ast.IsExported(funcDecl.Name.Name) {
					continue
				}

				repoMethod := extractMethodFromFuncType(funcDecl.Name.Name, funcDecl.Type)
				methods = append(methods, repoMethod)
			}
		}
	}

	return methods
}

// extractMethodFromFuncType extracts method details from a function type.
func extractMethodFromFuncType(name string, funcType *ast.FuncType) RepoMethod {
	repoMethod := RepoMethod{
		Name: name,
	}

	// Parse parameters
	if funcType.Params != nil {
		for _, param := range funcType.Params.List {
			paramType := exprToString(param.Type)
			if paramType == "context.Context" {
				repoMethod.HasContext = true
				continue // Don't add context to params list
			}

			// Handle multiple names for same type
			if len(param.Names) == 0 {
				repoMethod.Params = append(repoMethod.Params, RepoParam{
					Name: "",
					Type: paramType,
				})
			} else {
				for _, pname := range param.Names {
					repoMethod.Params = append(repoMethod.Params, RepoParam{
						Name: pname.Name,
						Type: paramType,
					})
				}
			}
		}
	}

	// Parse returns
	if funcType.Results != nil {
		for _, result := range funcType.Results.List {
			repoMethod.Returns = append(repoMethod.Returns, exprToString(result.Type))
		}
	}

	return repoMethod
}

// exprToString converts an AST expression to a string representation.
func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + exprToString(e.Elt)
		}
		return fmt.Sprintf("[%s]%s", exprToString(e.Len), exprToString(e.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprToString(e.Key), exprToString(e.Value))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.Ellipsis:
		return "..." + exprToString(e.Elt)
	default:
		return "unknown"
	}
}

// filterMethods filters methods based on include/exclude lists.
func filterMethods(methods []RepoMethod, include, exclude []string) []RepoMethod {
	// If include list is provided, only include those methods
	if len(include) > 0 {
		includeMap := make(map[string]bool)
		for _, m := range include {
			includeMap[m] = true
		}

		var filtered []RepoMethod
		for _, method := range methods {
			if includeMap[method.Name] {
				filtered = append(filtered, method)
			}
		}
		methods = filtered
	}

	// Apply exclude list
	if len(exclude) > 0 {
		excludeMap := make(map[string]bool)
		for _, m := range exclude {
			excludeMap[m] = true
		}

		var filtered []RepoMethod
		for _, method := range methods {
			if !excludeMap[method.Name] {
				filtered = append(filtered, method)
			}
		}
		methods = filtered
	}

	return methods
}

// generateServiceCode generates the service Go code.
func generateServiceCode(svcPkgName, svcModelName, modulePath, repoPkgName, repoImportAlias, repoModelName string, methods []RepoMethod, isStructRepo bool) string {
	var sb strings.Builder

	// Package declaration
	sb.WriteString(fmt.Sprintf("package %s\n\n", svcPkgName))

	// Imports
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString(fmt.Sprintf("\t%s \"%s/internal/repository/%s\"\n", repoImportAlias, modulePath, repoPkgName))
	sb.WriteString(")\n\n")

	// Service interface
	sb.WriteString(fmt.Sprintf("// Service defines the %s service interface.\n", svcModelName))
	sb.WriteString("type Service interface {\n")
	for _, method := range methods {
		sb.WriteString(fmt.Sprintf("\t%s\n", generateMethodSignature(method)))
	}
	sb.WriteString("}\n\n")

	// Determine repo type based on whether it's a struct or interface
	repoType := fmt.Sprintf("%s.Repository", repoImportAlias)
	if isStructRepo {
		repoType = fmt.Sprintf("*%s.Repository", repoImportAlias)
	}

	// Service implementation struct
	sb.WriteString(fmt.Sprintf("// service implements the %s Service interface.\n", svcModelName))
	sb.WriteString("type service struct {\n")
	sb.WriteString(fmt.Sprintf("\trepo %s\n", repoType))
	sb.WriteString("}\n\n")

	// Constructor
	sb.WriteString(fmt.Sprintf("// NewService creates a new %s service.\n", svcModelName))
	sb.WriteString(fmt.Sprintf("func NewService(repo %s) Service {\n", repoType))
	sb.WriteString("\treturn &service{repo: repo}\n")
	sb.WriteString("}\n")

	// Method implementations
	for _, method := range methods {
		sb.WriteString("\n")
		sb.WriteString(generateMethodImplementation(method))
	}

	return sb.String()
}

// generateMethodSignature generates a method signature for the interface.
func generateMethodSignature(method RepoMethod) string {
	var params []string
	params = append(params, "ctx context.Context")
	for _, p := range method.Params {
		if p.Name != "" {
			params = append(params, fmt.Sprintf("%s %s", p.Name, p.Type))
		} else {
			params = append(params, p.Type)
		}
	}

	returns := strings.Join(method.Returns, ", ")
	if len(method.Returns) > 1 {
		returns = "(" + returns + ")"
	}

	return fmt.Sprintf("%s(%s) %s", method.Name, strings.Join(params, ", "), returns)
}

// generateMethodImplementation generates a method implementation.
func generateMethodImplementation(method RepoMethod) string {
	var sb strings.Builder

	// Method signature
	var params []string
	params = append(params, "ctx context.Context")
	for _, p := range method.Params {
		if p.Name != "" {
			params = append(params, fmt.Sprintf("%s %s", p.Name, p.Type))
		} else {
			// Generate a name for unnamed parameters
			params = append(params, fmt.Sprintf("arg%d %s", len(params), p.Type))
		}
	}

	returns := strings.Join(method.Returns, ", ")
	if len(method.Returns) > 1 {
		returns = "(" + returns + ")"
	}

	sb.WriteString(fmt.Sprintf("// %s wraps the repository %s method.\n", method.Name, method.Name))
	sb.WriteString(fmt.Sprintf("func (s *service) %s(%s) %s {\n", method.Name, strings.Join(params, ", "), returns))

	// Method body - delegate to repo
	var args []string
	if method.HasContext {
		args = append(args, "ctx")
	}
	for i, p := range method.Params {
		if p.Name != "" {
			args = append(args, p.Name)
		} else {
			args = append(args, fmt.Sprintf("arg%d", i+1))
		}
	}

	sb.WriteString(fmt.Sprintf("\treturn s.repo.%s(%s)\n", method.Name, strings.Join(args, ", ")))
	sb.WriteString("}\n")

	return sb.String()
}
