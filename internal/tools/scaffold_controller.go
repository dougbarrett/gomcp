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
		Name:        "scaffold_controller",
		Description: "Create a standalone HTTP controller. Use this when you need custom endpoints without the full domain scaffolding.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldControllerInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldController(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldController(registry *Registry, input types.ScaffoldControllerInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateDomainName(input.DomainName); err != nil {
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
	pkgName := utils.ToPackageName(input.DomainName)
	urlPath := input.BasePath
	if urlPath == "" {
		urlPath = utils.ToURLPath(input.DomainName)
	}

	data := generator.DomainData{
		ModulePath:   modulePath,
		DomainName:   input.DomainName,
		ModelName:    utils.ToModelName(input.DomainName),
		PackageName:  pkgName,
		VariableName: utils.ToVariableName(input.DomainName),
		URLPath:      urlPath,
	}

	// Create directory
	controllerDir := filepath.Join("internal", "web", pkgName)
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
