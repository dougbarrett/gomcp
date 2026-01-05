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

// RegisterScaffoldController registers the scaffold_controller tool.
func RegisterScaffoldController(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_controller",
		Description: `Create a standalone HTTP controller. Use scaffold_domain instead for complete features.

Only use this tool when you need JUST HTTP handlers. Common cases:
- Adding endpoints for an existing service
- Creating utility/health endpoints

Generates: internal/web/{domain}/{domain}.go with HTTP handlers and routing.

Prefer scaffold_domain for new features - it generates all layers consistently.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldControllerInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldController(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldController(registry *Registry, input types.ScaffoldControllerInput) (types.ScaffoldResult, error) {
	// Validate input - support nested paths like "admin/users"
	if err := utils.ValidateDomainPath(input.DomainName); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if len(input.Actions) == 0 {
		return types.NewErrorResult("at least one action is required"), nil
	}

	// Validate actions
	for _, action := range input.Actions {
		if action.Name == "" {
			return types.NewErrorResult("action name is required"), nil
		}
		if err := utils.ValidateHTTPMethod(action.Method); err != nil {
			return types.NewErrorResult(fmt.Sprintf("action '%s': %v", action.Name, err)), nil
		}
	}

	// Validate base path if provided
	if input.BasePath != "" {
		if err := utils.ValidateURLPath(input.BasePath); err != nil {
			return types.NewErrorResult(err.Error()), nil
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
	// For nested paths like "admin/users", use base name for package/model
	baseDomain := utils.ParseDomainPath(input.DomainName)
	domainDir := utils.DomainPathToDir(input.DomainName)
	pkgName := utils.ToPackageName(baseDomain)
	urlPath := input.BasePath
	if urlPath == "" {
		urlPath = utils.ToURLPath(input.DomainName) // Keep full path for URL
	}

	// Default layout to "dashboard"
	layout := input.Layout
	if layout == "" {
		layout = "dashboard"
	}

	// Default route group to "public"
	routeGroup := input.RouteGroup
	if routeGroup == "" {
		routeGroup = "public"
	}

	data := generator.DomainData{
		ModulePath:   modulePath,
		DomainName:   input.DomainName,
		ModelName:    utils.ToModelName(baseDomain),
		PackageName:  pkgName,
		VariableName: utils.ToVariableName(baseDomain),
		URLPath:      urlPath,
		Layout:       layout,
		RouteGroup:   routeGroup,
	}

	// Create directory - use full path for nested domains
	controllerDir := filepath.Join("internal", "web", domainDir)
	if err := gen.EnsureDir(controllerDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate controller
	controllerPath := filepath.Join(controllerDir, pkgName+".go")
	if err := gen.GenerateFile("domain/controller.go.tmpl", controllerPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate controller: %v", err)), nil
	}

	// Prepare result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create controller for '%s'", input.DomainName),
			FilesCreated: result.FilesCreated,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created controller for '%s'", input.DomainName),
		FilesCreated: result.FilesCreated,
	}, nil
}
