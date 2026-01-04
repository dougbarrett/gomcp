package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/modifier"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterUpdateDIWiring registers the update_di_wiring tool.
func RegisterUpdateDIWiring(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_di_wiring",
		Description: "Update the main.go dependency injection wiring to include specified domains. Uses marker comments to inject imports, repository/service/controller instantiations, and route registrations.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.UpdateDIWiringInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := updateDIWiring(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func updateDIWiring(registry *Registry, input types.UpdateDIWiringInput) (types.ScaffoldResult, error) {
	if len(input.Domains) == 0 {
		return types.NewErrorResult("at least one domain is required"), nil
	}

	// Validate domains
	for _, domain := range input.Domains {
		if err := utils.ValidateDomainName(domain); err != nil {
			return types.NewErrorResult(fmt.Sprintf("domain '%s': %v", domain, err)), nil
		}
	}

	// Get module path
	modulePath, err := utils.GetModulePath(registry.WorkingDir)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to get module path: %v", err)), nil
	}

	// Find main.go
	mainGoPath := filepath.Join(registry.WorkingDir, "cmd", "web", "main.go")
	if !utils.FileExists(mainGoPath) {
		return types.NewErrorResult("main.go not found at cmd/web/main.go"), nil
	}

	if input.DryRun {
		// Just validate that we could do it
		injector, err := modifier.NewInjector(mainGoPath)
		if err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to read main.go: %v", err)), nil
		}

		// Check for required markers
		requiredMarkers := []string{
			modifier.MarkerImportsStart,
			modifier.MarkerReposStart,
			modifier.MarkerServicesStart,
			modifier.MarkerControllersStart,
			modifier.MarkerRoutesStart,
		}

		var missingMarkers []string
		for _, marker := range requiredMarkers {
			if !injector.HasMarker(marker) {
				missingMarkers = append(missingMarkers, marker)
			}
		}

		if len(missingMarkers) > 0 {
			return types.NewErrorResult(fmt.Sprintf("main.go is missing required markers: %v", missingMarkers)), nil
		}

		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would update main.go with wiring for %d domain(s)", len(input.Domains)),
			FilesUpdated: []string{"cmd/web/main.go"},
		}, nil
	}

	// Load injector
	injector, err := modifier.NewInjector(mainGoPath)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to read main.go: %v", err)), nil
	}

	// Inject wiring for each domain
	for _, domain := range input.Domains {
		pkgName := utils.ToPackageName(domain)

		// Inject imports
		repoImport := fmt.Sprintf("%s/internal/repository/%s", modulePath, pkgName)
		if err := injector.InjectImport(repoImport); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject repo import for '%s': %v", domain, err)), nil
		}

		serviceImport := fmt.Sprintf("%s/internal/services/%s", modulePath, pkgName)
		if err := injector.InjectImport(serviceImport); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject service import for '%s': %v", domain, err)), nil
		}

		controllerImport := fmt.Sprintf("%s/internal/web/%s", modulePath, pkgName)
		if err := injector.InjectImport(controllerImport); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject controller import for '%s': %v", domain, err)), nil
		}

		// Inject repository
		if err := injector.InjectRepo(domain, modulePath); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject repo for '%s': %v", domain, err)), nil
		}

		// Inject service
		if err := injector.InjectService(domain); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject service for '%s': %v", domain, err)), nil
		}

		// Inject controller
		if err := injector.InjectController(domain); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject controller for '%s': %v", domain, err)), nil
		}

		// Inject route
		if err := injector.InjectRoute(domain); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject route for '%s': %v", domain, err)), nil
		}
	}

	// Save changes
	if err := injector.Save(); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to save main.go: %v", err)), nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully updated main.go with wiring for %d domain(s)", len(input.Domains)),
		FilesUpdated: []string{"cmd/web/main.go"},
	}, nil
}
