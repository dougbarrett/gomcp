package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterScaffoldProject registers the scaffold_project tool.
func RegisterScaffoldProject(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_project",
		Description: `ALWAYS use this tool to initialize new Go web projects. Never manually create go.mod, main.go, or project structure files.

Generates a complete, production-ready project with:
- Clean architecture (models, repositories, services, controllers)
- templ + HTMX for interactive UIs with Tailwind CSS styling
- Reusable UI components (buttons, cards, forms, tables, modals)
- GORM database setup (sqlite, postgres, or mysql)
- Taskfile for development commands
- Hot reload with Air

Directory behavior:
- Auto-detects if current directory name matches project_name and scaffolds in place
- Otherwise creates a new subdirectory with the project name
- Use in_current_dir: true to force scaffolding in current directory regardless of name

Options:
- in_current_dir: true to force scaffold in current directory
- with_auth: true to include full authentication system (login, register, sessions, middleware)
- dry_run: true to preview files without writing

After running: Execute 'go mod tidy' then 'task dev' to start.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldProjectInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldProject(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldProject(registry *Registry, input types.ScaffoldProjectInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateProjectName(input.ProjectName); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}
	if err := utils.ValidateModulePath(input.ModulePath); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}
	if err := utils.ValidateDatabaseType(input.DatabaseType); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	// Set defaults
	dbType := input.DatabaseType
	if dbType == "" {
		dbType = "sqlite"
	}

	// Auto-detect if we should scaffold in current directory:
	// If the current directory name matches the project name, use current dir
	currentDirName := filepath.Base(registry.WorkingDir)
	useCurrentDir := input.InCurrentDir || currentDirName == input.ProjectName

	// Create project path
	var projectPath string
	if useCurrentDir {
		projectPath = registry.WorkingDir
	} else {
		projectPath = filepath.Join(registry.WorkingDir, input.ProjectName)
		// Check if project already exists (only when creating subdirectory)
		if utils.DirExists(projectPath) && !input.DryRun {
			return types.NewErrorResult(fmt.Sprintf("project directory already exists: %s", projectPath)), nil
		}
	}

	// Create generator
	gen := registry.NewGenerator(projectPath)
	gen.SetDryRun(input.DryRun)

	// Prepare template data
	data := generator.ProjectData{
		ProjectName:  input.ProjectName,
		ModulePath:   input.ModulePath,
		DatabaseType: dbType,
		WithAuth:     input.WithAuth,
	}

	// Create directory structure
	directories := []string{
		"cmd/web",
		"cmd/seed",
		"internal/config",
		"internal/database",
		"internal/models",
		"internal/repository",
		"internal/services",
		"internal/web/middleware",
		"internal/web/layouts",
		"internal/web/components",
		"config/en/pages",
		"assets/css",
		"assets/js",
		"components",
		"utils",
	}

	// Add auth directories if WithAuth is enabled
	if input.WithAuth {
		directories = append(directories,
			"internal/repository/user",
			"internal/services/auth",
			"internal/web/auth",
			"internal/web/auth/views",
			"internal/web/dashboard",
			"internal/web/dashboard/views",
		)
	}

	for _, dir := range directories {
		if err := gen.EnsureDir(dir); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to create directory %s: %v", dir, err)), nil
		}
	}

	// Generate files
	files := []struct {
		template string
		output   string
	}{
		{"project/go.mod.tmpl", "go.mod"},
		{"project/main.go.tmpl", "cmd/web/main.go"},
		{"project/seed_main.go.tmpl", "cmd/seed/main.go"},
		{"project/config.go.tmpl", "internal/config/config.go"},
		{"project/database.go.tmpl", "internal/database/database.go"},
		{"project/base_model.go.tmpl", "internal/models/base.go"},
		{"project/router.go.tmpl", "internal/web/router.go"},
		{"project/middleware.go.tmpl", "internal/web/middleware/middleware.go"},
		{"project/response.go.tmpl", "internal/web/response.go"},
		{"project/base_layout.templ.tmpl", "internal/web/layouts/base.templ"},
		{"project/common_components.templ.tmpl", "internal/web/components/common.templ"},
		{"project/taskfile.yml.tmpl", "Taskfile.yml"},
		{"project/air.toml.tmpl", ".air.toml"},
		{"project/tailwind_input.css.tmpl", "assets/css/input.css"},
		{"project/tailwind_output.css.tmpl", "assets/css/output.css"},
		{"project/tailwind.config.js.tmpl", "tailwind.config.js"},
		{"project/app.toml.tmpl", "config/en/app.toml"},
		{"project/menu.toml.tmpl", "config/en/menu.toml"},
		{"project/gitignore.tmpl", ".gitignore"},
	}

	for _, f := range files {
		if err := gen.GenerateFile(f.template, f.output, data); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate %s: %v", f.output, err)), nil
		}
	}

	// Generate auth files if WithAuth is enabled
	if input.WithAuth {
		authData := generator.NewAuthData(input.ModulePath, input.ProjectName)
		authFiles := []struct {
			template string
			output   string
		}{
			// Role model (must be before User model since User references Role)
			{"auth/role_model.go.tmpl", "internal/models/role.go"},
			// User model and repository
			{"auth/user_model.go.tmpl", "internal/models/user.go"},
			{"auth/user_repository.go.tmpl", "internal/repository/user/user.go"},
			// Auth service
			{"auth/auth_service.go.tmpl", "internal/services/auth/auth.go"},
			{"auth/session.go.tmpl", "internal/services/auth/session.go"},
			// Auth middleware and controller
			{"auth/auth_middleware.go.tmpl", "internal/web/middleware/auth.go"},
			{"auth/auth_controller.go.tmpl", "internal/web/auth/auth.go"},
			// Auth views
			{"auth/auth_layout.templ.tmpl", "internal/web/auth/views/layout.templ"},
			{"auth/login.templ.tmpl", "internal/web/auth/views/login.templ"},
			{"auth/register.templ.tmpl", "internal/web/auth/views/register.templ"},
			// Dashboard
			{"auth/dashboard_controller.go.tmpl", "internal/web/dashboard/dashboard.go"},
			{"auth/dashboard.templ.tmpl", "internal/web/dashboard/views/dashboard.templ"},
		}

		for _, f := range authFiles {
			if err := gen.GenerateFile(f.template, f.output, authData); err != nil {
				return types.NewErrorResult(fmt.Sprintf("failed to generate auth file %s: %v", f.output, err)), nil
			}
		}
	}

	// Add MCP marker instructions to CLAUDE.md and AGENTS.md
	if !input.DryRun {
		if err := addMCPMarkerInstructions(projectPath); err != nil {
			// Non-fatal error - just log and continue
			// The project is still functional without these instructions
		}
	}

	// Prepare result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	var nextSteps []string
	if useCurrentDir {
		nextSteps = []string{
			"go mod tidy",
			"task dev  # Start development server",
		}
	} else {
		nextSteps = []string{
			fmt.Sprintf("cd %s", input.ProjectName),
			"go mod tidy",
			"task dev  # Start development server",
		}
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create project '%s' with %d files", input.ProjectName, len(result.FilesCreated)),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created project '%s'", input.ProjectName),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}

// mcpMarkerInstructions contains the warning about MCP marker comments.
const mcpMarkerInstructions = `## MCP Scaffolding Markers - DO NOT MODIFY

This project uses MCP (Model Context Protocol) scaffolding tools that rely on special marker comments for code injection. These markers enable the scaffolding tools to add new code (imports, dependencies, routes, etc.) to existing files without overwriting your changes.

### Critical: Never modify or remove these marker comments:

` + "```" + `
// MCP:IMPORTS:START
// MCP:IMPORTS:END

// MCP:MODELS:START
// MCP:MODELS:END

// MCP:REPOS:START
// MCP:REPOS:END

// MCP:SERVICES:START
// MCP:SERVICES:END

// MCP:CONTROLLERS:START
// MCP:CONTROLLERS:END

// MCP:ROUTES:START
// MCP:ROUTES:END

// MCP:REPO_INTERFACE:START
// MCP:REPO_INTERFACE:END

// MCP:REPO_METHODS:START
// MCP:REPO_METHODS:END

// MCP:SERVICE_INTERFACE:START
// MCP:SERVICE_INTERFACE:END

// MCP:SERVICE_METHODS:START
// MCP:SERVICE_METHODS:END

// MCP:HANDLERS:START
// MCP:HANDLERS:END
` + "```" + `

### Why this matters:
- The scaffolding tools inject code **between** these markers
- Removing or modifying the markers will break the ability to scaffold new domains
- You can safely add your own code **outside** of these marker pairs
- Code **inside** the markers may be regenerated - add custom code elsewhere

### Files that contain these markers:
- ` + "`cmd/web/main.go`" + ` - DI wiring (imports, repos, services, controllers, routes)
- ` + "`internal/repository/{domain}/{domain}.go`" + ` - Repository interface and methods
- ` + "`internal/services/{domain}/{domain}.go`" + ` - Service interface and methods
- ` + "`internal/web/{domain}/{domain}.go`" + ` - Controller routes and handlers
`

// addMCPMarkerInstructions adds the MCP marker warning to CLAUDE.md and AGENTS.md if they exist.
func addMCPMarkerInstructions(projectPath string) error {
	marker := "## MCP Scaffolding Markers"

	// Try to add to CLAUDE.md
	claudePath := filepath.Join(projectPath, "CLAUDE.md")
	if err := utils.AppendToFileIfNotContains(claudePath, marker, mcpMarkerInstructions); err != nil {
		// Try parent directory (in case scaffolding in subdirectory)
		parentClaudePath := filepath.Join(filepath.Dir(projectPath), "CLAUDE.md")
		utils.AppendToFileIfNotContains(parentClaudePath, marker, mcpMarkerInstructions)
	}

	// Try to add to AGENTS.md
	agentsPath := filepath.Join(projectPath, "AGENTS.md")
	if err := utils.AppendToFileIfNotContains(agentsPath, marker, mcpMarkerInstructions); err != nil {
		// Try parent directory
		parentAgentsPath := filepath.Join(filepath.Dir(projectPath), "AGENTS.md")
		utils.AppendToFileIfNotContains(parentAgentsPath, marker, mcpMarkerInstructions)
	}

	return nil
}
