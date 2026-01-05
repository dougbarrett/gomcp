package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/modifier"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterScaffoldDomain registers the scaffold_domain tool.
func RegisterScaffoldDomain(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_domain",
		Description: `PRIMARY TOOL for adding new features. ALWAYS use this instead of manually creating models, repositories, services, or controllers.

Generates ALL layers at once following clean architecture:
- Model (internal/models/{domain}.go)
- Repository with CRUD operations (internal/repository/{domain}/)
- Service with DTOs (internal/services/{domain}/)
- Controller with HTTP handlers (internal/web/{domain}/)
- Optional CRUD views (with_crud_views: true, default)

Supports relationships: belongs_to, has_one, has_many, many_to_many

Supported field types:
- Scalars: string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool
- Time: time.Time, *time.Time (nullable)
- Pointers (nullable): *string, *int, *int64, *uint, *float64, *bool
- Slices: []byte, []string, []int, []uint
- Custom types: any valid Go identifier (e.g., Status, models.Role)

Automatically wires into main.go DI container. Run 'go mod tidy' and 'templ generate' after.

Use dry_run: true to preview all generated files first.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldDomainInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldDomain(registry *Registry, input types.ScaffoldDomainInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateDomainName(input.DomainName); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	// Check for reserved domain names that conflict with auth scaffolding
	if utils.ToPackageName(input.DomainName) == "user" {
		// Check if User model exists (created by scaffold_project with_auth)
		userModelPath := filepath.Join(registry.WorkingDir, "internal", "models", "user.go")
		if utils.FileExists(userModelPath) {
			return types.NewErrorResult(
				"Cannot scaffold 'user' domain: A User model already exists (likely from scaffold_project with_auth). " +
					"The auth system uses a specialized User model with PasswordHash, Active, and Role fields. " +
					"To add user-related features, extend the existing auth service instead.",
			), nil
		}
	}

	if len(input.Fields) == 0 {
		return types.NewErrorResult("at least one field is required"), nil
	}

	// Validate each field
	for _, field := range input.Fields {
		if err := utils.ValidateFieldName(field.Name); err != nil {
			return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
		}
		if err := utils.ValidateFieldType(field.Type); err != nil {
			return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
		}
		if field.FormType != "" {
			if err := utils.ValidateFormType(field.FormType); err != nil {
				return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
			}
		}
	}

	// Validate relationships
	for _, rel := range input.Relationships {
		if err := utils.ValidateRelationshipType(rel.Type); err != nil {
			return types.NewErrorResult(fmt.Sprintf("relationship to '%s': %v", rel.Model, err)), nil
		}
		if err := utils.ValidateRelationshipModel(rel.Model); err != nil {
			return types.NewErrorResult(fmt.Sprintf("relationship: %v", err)), nil
		}
		if rel.OnDelete != "" {
			if err := utils.ValidateOnDelete(rel.OnDelete); err != nil {
				return types.NewErrorResult(fmt.Sprintf("relationship to '%s': %v", rel.Model, err)), nil
			}
		}
	}

	// Get module path from go.mod
	modulePath, err := utils.GetModulePath(registry.WorkingDir)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to get module path: %v", err)), nil
	}

	// Create generator
	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Prepare template data
	data := generator.NewDomainData(input, modulePath)

	// Create directories
	pkgName := utils.ToPackageName(input.DomainName)
	directories := []string{
		filepath.Join("internal", "repository", pkgName),
		filepath.Join("internal", "services", pkgName),
		filepath.Join("internal", "web", pkgName),
	}

	if input.GetWithCrudViews() {
		directories = append(directories, filepath.Join("internal", "web", pkgName, "views"))
	}

	for _, dir := range directories {
		if err := gen.EnsureDir(dir); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to create directory %s: %v", dir, err)), nil
		}
	}

	// Generate model
	modelPath := filepath.Join("internal", "models", pkgName+".go")
	if err := gen.GenerateFile("domain/model.go.tmpl", modelPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate model: %v", err)), nil
	}

	// Generate repository
	repoPath := filepath.Join("internal", "repository", pkgName, pkgName+".go")
	if err := gen.GenerateFile("domain/repository.go.tmpl", repoPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate repository: %v", err)), nil
	}

	// Generate service
	servicePath := filepath.Join("internal", "services", pkgName, pkgName+".go")
	if err := gen.GenerateFile("domain/service.go.tmpl", servicePath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate service: %v", err)), nil
	}

	// Generate DTOs
	dtoPath := filepath.Join("internal", "services", pkgName, "dto.go")
	if err := gen.GenerateFile("domain/dto.go.tmpl", dtoPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate DTOs: %v", err)), nil
	}

	// Generate controller
	controllerPath := filepath.Join("internal", "web", pkgName, pkgName+".go")
	if err := gen.GenerateFile("domain/controller.go.tmpl", controllerPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate controller: %v", err)), nil
	}

	// Prepare result
	result := gen.Result()

	// Check for conflicts - if any files would be overwritten, return conflict result
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	// Inject into main.go if not dry run
	if !input.DryRun {
		mainGoPath := filepath.Join(registry.WorkingDir, "cmd", "web", "main.go")
		if utils.FileExists(mainGoPath) {
			if err := injectDomainWiring(mainGoPath, modulePath, pkgName, input.DomainName); err != nil {
				// Log warning but don't fail
				fmt.Printf("Warning: could not inject DI wiring: %v\n", err)
			} else {
				result.FilesUpdated = append(result.FilesUpdated, "cmd/web/main.go")
			}
		}
	}
	nextSteps := []string{
		"go mod tidy",
		"templ generate",
		fmt.Sprintf("Add business logic to internal/services/%s/%s.go", pkgName, pkgName),
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create domain '%s' with %d files", input.DomainName, len(result.FilesCreated)),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created domain '%s'", input.DomainName),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}

// injectDomainWiring injects the domain wiring into main.go.
func injectDomainWiring(mainGoPath, modulePath, pkgName, domainName string) error {
	injector, err := modifier.NewInjector(mainGoPath)
	if err != nil {
		return err
	}

	// Inject imports
	repoImport := fmt.Sprintf("%s/internal/repository/%s", modulePath, pkgName)
	if err := injector.InjectImport(repoImport); err != nil {
		return err
	}

	serviceImport := fmt.Sprintf("%s/internal/services/%s", modulePath, pkgName)
	if err := injector.InjectImport(serviceImport); err != nil {
		return err
	}

	controllerImport := fmt.Sprintf("%s/internal/web/%s", modulePath, pkgName)
	if err := injector.InjectImport(controllerImport); err != nil {
		return err
	}

	// Inject repository
	if err := injector.InjectRepo(domainName, modulePath); err != nil {
		return err
	}

	// Inject service
	if err := injector.InjectService(domainName); err != nil {
		return err
	}

	// Inject controller
	if err := injector.InjectController(domainName); err != nil {
		return err
	}

	// Inject route
	if err := injector.InjectRoute(domainName); err != nil {
		return err
	}

	return injector.Save()
}
