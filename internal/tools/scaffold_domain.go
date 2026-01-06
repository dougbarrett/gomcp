package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/metadata"
	"github.com/dbb1dev/go-mcp/internal/modifier"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ScaffolderVersion is the current version of the scaffolding tools.
// Used for tracking which version generated the code for future upgrades.
const ScaffolderVersion = "0.1.0"

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

Layout options (layout parameter):
- "dashboard" (default): Views wrapped in DashboardPage layout with sidebar
- "base": Views wrapped in BasePage layout without sidebar
- "none": Views rendered without layout wrapper

Route group options (route_group parameter):
- "public" (default): No authentication required
- "authenticated": Requires user login (RequireAuth middleware)
- "admin": Requires admin role (RequireAuth + RequireAdmin middleware)

Form style options (form_style parameter):
- "modal" (default): Forms displayed in popup modal overlays
- "page": Forms displayed as full page navigation (like user management)

Form type options (form_type parameter on fields):
- "input" (default): Standard text input
- "textarea": Multi-line text area
- "select": Dropdown select (use options: ["opt1", "opt2"] to specify choices)
- "checkbox": Boolean checkbox
- "number": Numeric input
- "email": Email input with validation
- "password": Password input (masked)
- "date": Date picker
- "datetime": Date and time picker

Examples:

1. Simple public domain (blog posts):
   scaffold_domain: {
     domain_name: "post",
     fields: [
       {name: "Title", type: "string"},
       {name: "Content", type: "string", form_type: "textarea"},
       {name: "PublishedAt", type: "*time.Time"}
     ]
   }

2. Authenticated user feature (orders):
   scaffold_domain: {
     domain_name: "order",
     route_group: "authenticated",
     fields: [
       {name: "Total", type: "float64"},
       {name: "Status", type: "string"}
     ],
     relationships: [
       {type: "belongs_to", model: "User"}
     ]
   }

3. Admin-only management (categories):
   scaffold_domain: {
     domain_name: "category",
     route_group: "admin",
     layout: "dashboard",
     fields: [
       {name: "Name", type: "string"},
       {name: "Slug", type: "string", gorm_tags: "uniqueIndex"}
     ]
   }

4. Domain with select dropdown (discounts with status):
   scaffold_domain: {
     domain_name: "discount",
     fields: [
       {name: "Name", type: "string"},
       {name: "Amount", type: "float64", form_type: "number"},
       {name: "DiscountType", type: "string", form_type: "select", options: ["percentage", "fixed"]},
       {name: "Status", type: "string", form_type: "select", options: ["draft", "active", "expired"]}
     ]
   }

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

	// Generate CRUD views if requested
	if input.GetWithCrudViews() {
		viewsDir := filepath.Join("internal", "web", pkgName, "views")

		// Generate list view
		listPath := filepath.Join(viewsDir, "list.templ")
		if err := gen.GenerateFile("views/list.templ.tmpl", listPath, data); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate list view: %v", err)), nil
		}

		// Generate show view
		showPath := filepath.Join(viewsDir, "show.templ")
		if err := gen.GenerateFile("views/show.templ.tmpl", showPath, data); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate show view: %v", err)), nil
		}

		// Generate form view
		formPath := filepath.Join(viewsDir, pkgName+"_form.templ")
		if err := gen.GenerateFile("views/form.templ.tmpl", formPath, data); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate form view: %v", err)), nil
		}

		// Generate partials (card, empty state, etc.)
		partialsPath := filepath.Join(viewsDir, "partials.templ")
		if err := gen.GenerateFile("views/partials.templ.tmpl", partialsPath, data); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate partials: %v", err)), nil
		}
	}

	// Prepare result
	result := gen.Result()

	// Check for conflicts - if any files would be overwritten, return conflict result
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	// Inject into main.go, database.go, and base_layout.templ if not dry run
	if !input.DryRun {
		mainGoPath := filepath.Join(registry.WorkingDir, "cmd", "web", "main.go")
		databaseGoPath := filepath.Join(registry.WorkingDir, "internal", "database", "database.go")
		layoutPath := filepath.Join(registry.WorkingDir, "internal", "web", "layouts", "base_layout.templ")
		if utils.FileExists(mainGoPath) {
			if err := injectDomainWiring(mainGoPath, databaseGoPath, modulePath, pkgName, input.DomainName, data.RouteGroup, input.Relationships, data.WithCrudViews); err != nil {
				// Log warning but don't fail
				fmt.Printf("Warning: could not inject DI wiring: %v\n", err)
			} else {
				result.FilesUpdated = append(result.FilesUpdated, "cmd/web/main.go")
				if utils.FileExists(databaseGoPath) {
					result.FilesUpdated = append(result.FilesUpdated, "internal/database/database.go")
				}
				// Track layout update for authenticated/admin routes
				if (data.RouteGroup == "authenticated" || data.RouteGroup == "admin") && utils.FileExists(layoutPath) {
					result.FilesUpdated = append(result.FilesUpdated, "internal/web/layouts/base_layout.templ")
				}
			}
		}

		// Inject inverse relationships into related models
		if len(input.Relationships) > 0 {
			injectInverseRelationships(registry.WorkingDir, input.DomainName, input.Relationships, &result.FilesUpdated)
		}
	}
	nextSteps := []string{
		"go mod tidy",
		"templ generate",
		fmt.Sprintf("Add business logic to internal/services/%s/%s.go", pkgName, pkgName),
	}

	// Suggest tools for extending the domain
	suggestedTools := []types.ToolHint{
		{
			Tool:        "extend_service",
			Description: fmt.Sprintf("Add custom business logic methods to the %s service", input.DomainName),
			Example:     fmt.Sprintf(`extend_service: { domain: "%s", methods: [{ name: "Archive", params: [{ name: "id", type: "uint" }], returns: "error" }] }`, input.DomainName),
			Priority:    "optional",
		},
		{
			Tool:        "extend_repository",
			Description: fmt.Sprintf("Add custom query methods to the %s repository", input.DomainName),
			Example:     fmt.Sprintf(`extend_repository: { domain: "%s", methods: [{ name: "FindByStatus", params: [{ name: "status", type: "string" }], returns: "[]models.%s, error" }] }`, input.DomainName, utils.ToPascalCase(input.DomainName)),
			Priority:    "optional",
		},
		{
			Tool:        "extend_controller",
			Description: fmt.Sprintf("Add custom HTTP endpoints to the %s controller", input.DomainName),
			Example:     fmt.Sprintf(`extend_controller: { domain: "%s", endpoints: [{ name: "Archive", method: "POST", path: "/{id}/archive" }] }`, input.DomainName),
			Priority:    "optional",
		},
		types.HintScaffoldSeed,
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:        true,
			Message:        fmt.Sprintf("Dry run: Would create domain '%s' with %d files", input.DomainName, len(result.FilesCreated)),
			FilesCreated:   result.FilesCreated,
			NextSteps:      nextSteps,
			SuggestedTools: suggestedTools,
		}, nil
	}

	// Save scaffold metadata for future sync/upgrade capabilities
	metaStore := metadata.NewStore(registry.WorkingDir)
	if err := metaStore.SaveDomain(input.DomainName, input, ScaffolderVersion); err != nil {
		// Log warning but don't fail - metadata is optional
		fmt.Printf("Warning: could not save scaffold metadata: %v\n", err)
	} else {
		result.FilesUpdated = append(result.FilesUpdated, ".mcp/scaffold-metadata.json")
	}

	return types.ScaffoldResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully created domain '%s'", input.DomainName),
		FilesCreated:   result.FilesCreated,
		FilesUpdated:   result.FilesUpdated,
		NextSteps:      nextSteps,
		SuggestedTools: suggestedTools,
	}, nil
}

// injectInverseRelationships injects inverse relationship fields into related models.
// For example, if Order has belongs_to: User, this will add Orders []Order to User model.
func injectInverseRelationships(workingDir string, domainName string, relationships []types.RelationshipDef, filesUpdated *[]string) {
	for _, rel := range relationships {
		var inverseFieldCode string
		var inverseModelPath string

		relatedPkgName := utils.ToPackageName(rel.Model)
		inverseModelPath = filepath.Join(workingDir, "internal", "models", relatedPkgName+".go")

		// Check if the related model file exists
		if !utils.FileExists(inverseModelPath) {
			continue
		}

		// Check if the model has the relationship markers
		injector, err := modifier.NewInjector(inverseModelPath)
		if err != nil {
			continue
		}

		if !injector.HasMarker(modifier.MarkerRelationshipsStart) {
			continue
		}

		// Determine the inverse relationship based on the relationship type
		modelName := utils.ToModelName(domainName)
		switch rel.Type {
		case "belongs_to":
			// belongs_to -> has_many (e.g., Order belongs_to User -> User has_many Orders)
			// The foreignKey is the FK column on the child model (Order.UserID), not the child model's name
			fieldName := utils.Pluralize(modelName)
			relatedModelName := utils.ToModelName(rel.Model)
			foreignKey := relatedModelName + "ID"
			inverseFieldCode = fmt.Sprintf(`%s []%s `+"`"+`gorm:"foreignKey:%s" json:"%s,omitempty"`+"`",
				fieldName, modelName, foreignKey, utils.ToSnakeCase(fieldName))

		case "has_one":
			// has_one -> belongs_to (e.g., User has_one Profile -> Profile belongs_to User)
			// The inverse side needs the foreign key pointing back
			inverseFieldCode = fmt.Sprintf(`%sID uint `+"`"+`json:"%s_id,omitempty"`+"`"+`
	%s *%s `+"`"+`gorm:"foreignKey:%sID" json:"%s,omitempty"`+"`",
				modelName, utils.ToSnakeCase(modelName),
				modelName, modelName, modelName, utils.ToSnakeCase(modelName))

		case "has_many":
			// has_many -> belongs_to (e.g., User has_many Posts -> Post belongs_to User)
			// The inverse side (the "many" side) needs the foreign key pointing back
			inverseFieldCode = fmt.Sprintf(`%sID uint `+"`"+`json:"%s_id,omitempty"`+"`"+`
	%s *%s `+"`"+`gorm:"foreignKey:%sID" json:"%s,omitempty"`+"`",
				modelName, utils.ToSnakeCase(modelName),
				modelName, modelName, modelName, utils.ToSnakeCase(modelName))

		case "many_to_many":
			// many_to_many -> many_to_many (bidirectional)
			fieldName := utils.Pluralize(modelName)
			joinTable := rel.JoinTable
			if joinTable == "" {
				// Default join table name (alphabetical order)
				names := []string{utils.ToSnakeCase(utils.Pluralize(domainName)), utils.ToSnakeCase(utils.Pluralize(rel.Model))}
				if names[0] > names[1] {
					names[0], names[1] = names[1], names[0]
				}
				joinTable = names[0] + "_" + names[1]
			}
			inverseFieldCode = fmt.Sprintf(`%s []%s `+"`"+`gorm:"many2many:%s" json:"%s,omitempty"`+"`",
				fieldName, modelName, joinTable, utils.ToSnakeCase(fieldName))

		default:
			continue
		}

		// Inject the inverse relationship
		if err := injector.InjectRelationship(inverseFieldCode); err != nil {
			continue
		}

		if err := injector.Save(); err != nil {
			continue
		}

		*filesUpdated = append(*filesUpdated, filepath.Join("internal", "models", relatedPkgName+".go"))
	}
}

// injectDomainWiring injects the domain wiring into main.go, database.go, and base_layout.templ.
func injectDomainWiring(mainGoPath, databaseGoPath, modulePath, pkgName, domainName, routeGroup string, relationships []types.RelationshipDef, withCrudViews bool) error {
	// Inject into main.go
	mainInjector, err := modifier.NewInjector(mainGoPath)
	if err != nil {
		return err
	}

	// Inject imports with aliases to avoid naming conflicts
	repoImport := fmt.Sprintf("%s/internal/repository/%s", modulePath, pkgName)
	repoAlias := utils.ToRepoImportAlias(domainName)
	if err := mainInjector.InjectImportWithAlias(repoImport, repoAlias); err != nil {
		return err
	}

	serviceImport := fmt.Sprintf("%s/internal/services/%s", modulePath, pkgName)
	serviceAlias := utils.ToServiceImportAlias(domainName)
	if err := mainInjector.InjectImportWithAlias(serviceImport, serviceAlias); err != nil {
		return err
	}

	controllerImport := fmt.Sprintf("%s/internal/web/%s", modulePath, pkgName)
	controllerAlias := utils.ToControllerImportAlias(domainName)
	if err := mainInjector.InjectImportWithAlias(controllerImport, controllerAlias); err != nil {
		return err
	}

	// Inject repository
	if err := mainInjector.InjectRepo(domainName, modulePath); err != nil {
		return err
	}

	// Inject service
	if err := mainInjector.InjectService(domainName); err != nil {
		return err
	}

	// Collect belongs_to relationships for controller injection
	var relatedDomains []string
	if withCrudViews {
		for _, rel := range relationships {
			if rel.Type == "belongs_to" {
				relatedDomains = append(relatedDomains, rel.Model)
			}
		}
	}

	// Inject controller with related services if needed
	if err := mainInjector.InjectControllerWithRelations(domainName, relatedDomains); err != nil {
		return err
	}

	// Inject route with route group
	if err := mainInjector.InjectRouteWithGroup(domainName, routeGroup); err != nil {
		return err
	}

	if err := mainInjector.Save(); err != nil {
		return err
	}

	// Inject model into database.go AutoMigrate
	if databaseGoPath != "" && utils.FileExists(databaseGoPath) {
		dbInjector, err := modifier.NewInjector(databaseGoPath)
		if err != nil {
			return err
		}

		modelName := utils.ToModelName(domainName)
		if err := dbInjector.InjectModel(modelName); err != nil {
			return err
		}

		if err := dbInjector.Save(); err != nil {
			return err
		}
	}

	// Inject navigation item into base_layout.templ for authenticated/admin routes
	if routeGroup == "authenticated" || routeGroup == "admin" {
		// Get directory of main.go to find base_layout.templ
		baseDir := filepath.Dir(filepath.Dir(mainGoPath)) // cmd/web -> cmd -> project root
		layoutPath := filepath.Join(baseDir, "internal", "web", "layouts", "base_layout.templ")

		if utils.FileExists(layoutPath) {
			navInjector, err := modifier.NewInjector(layoutPath)
			if err == nil {
				// Use default "folder" icon - domains can customize later
				if err := navInjector.InjectNavItem(domainName, routeGroup, "folder"); err == nil {
					_ = navInjector.Save()
				}
			}
		}
	}

	return nil
}
