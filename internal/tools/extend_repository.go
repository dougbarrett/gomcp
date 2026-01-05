package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/modifier"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterExtendRepository registers the extend_repository tool.
func RegisterExtendRepository(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "extend_repository",
		Description: `Add custom methods to an existing repository.

Use this to add data access methods beyond standard CRUD. The repository must have been created
with scaffold_domain (which includes injection markers).

Example methods:
- FindByEmail: Look up by a specific field
- FindByStatus: Filter by status
- CountByCategory: Aggregate queries
- SoftDelete: Custom delete behavior

Each method is added to both the interface and the implementation. Use the body parameter
to provide the implementation, or leave empty for a TODO placeholder.

Template variables available in body:
- [[.ModelName]]: The model name in PascalCase (e.g., "Order")
- [[.VariableName]]: The variable name in camelCase (e.g., "order")
- [[.PackageName]]: The package name (e.g., "order")`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ExtendRepositoryInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := extendRepository(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

// extendRepository adds methods to an existing repository.
func extendRepository(registry *Registry, input types.ExtendRepositoryInput) (types.ScaffoldResult, error) {
	if input.Domain == "" {
		return types.NewErrorResult("domain is required"), nil
	}
	if len(input.Methods) == 0 {
		return types.NewErrorResult("at least one method is required"), nil
	}

	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Derive names from domain
	packageName := utils.ToPackageName(input.Domain)
	modelName := utils.ToPascalCase(input.Domain)
	variableName := utils.ToCamelCase(input.Domain)

	// Check if repository file exists
	repoPath := filepath.Join("internal", "repository", packageName, packageName+".go")
	if !gen.FileExists(repoPath) {
		return types.NewErrorResult(fmt.Sprintf("repository file not found: %s. Use scaffold_domain first.", repoPath)), nil
	}

	// Read the existing repository file
	content, err := gen.ReadFile(repoPath)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to read repository file: %v", err)), nil
	}

	// Generate interface signatures and method implementations
	var interfaceMethods []string
	var methodImpls []string

	for _, method := range input.Methods {
		// Build parameter list (always include ctx context.Context first)
		params := []string{"ctx context.Context"}
		for _, p := range method.Params {
			params = append(params, fmt.Sprintf("%s %s", p.Name, p.Type))
		}
		paramStr := strings.Join(params, ", ")

		// Build return type
		returnType := method.Returns
		if returnType == "" {
			returnType = "error"
		}
		// Wrap multiple returns in parentheses if not already
		if strings.Contains(returnType, ",") && !strings.HasPrefix(returnType, "(") {
			returnType = "(" + returnType + ")"
		}

		// Interface method signature
		interfaceSig := fmt.Sprintf("\t%s(%s) %s", method.Name, paramStr, returnType)
		interfaceMethods = append(interfaceMethods, interfaceSig)

		// Method implementation
		description := method.Description
		if description == "" {
			description = fmt.Sprintf("%s performs a custom repository operation.", method.Name)
		}

		// Process body - replace template variables
		body := method.Body
		if body == "" {
			body = "\t// TODO: Implement this method\n\treturn nil"
		}
		body = strings.ReplaceAll(body, "[[.ModelName]]", modelName)
		body = strings.ReplaceAll(body, "[[.VariableName]]", variableName)
		body = strings.ReplaceAll(body, "[[.PackageName]]", packageName)

		impl := fmt.Sprintf(`
// %s %s
func (r *repository) %s(%s) %s {
%s
}`, method.Name, description, method.Name, paramStr, returnType, body)
		methodImpls = append(methodImpls, impl)
	}

	// Use injector to add content between markers
	injector := modifier.NewInjectorFromContent(content)

	// Inject interface methods
	interfaceContent := strings.Join(interfaceMethods, "\n")
	if err := injector.InjectBetweenMarkers("MCP:REPO_INTERFACE:START", "MCP:REPO_INTERFACE:END", interfaceContent); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to inject interface methods: %v. Make sure the repository was scaffolded with markers.", err)), nil
	}

	// Inject method implementations
	implContent := strings.Join(methodImpls, "\n")
	if err := injector.InjectBetweenMarkers("MCP:REPO_METHODS:START", "MCP:REPO_METHODS:END", implContent); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to inject method implementations: %v. Make sure the repository was scaffolded with markers.", err)), nil
	}

	// Write the modified file
	if err := gen.WriteFile(repoPath, injector.Content()); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to write repository file: %v", err)), nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Added %d method(s) to %s repository", len(input.Methods), modelName),
		FilesUpdated: []string{repoPath},
	}, nil
}
