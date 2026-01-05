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
		Name: "update_di_wiring",
		Description: `IMPORTANT: Run this after scaffold_domain to wire up new domains.

Automatically updates cmd/web/main.go to:
- Add import statements for repository, service, controller
- Instantiate repository, service, and controller
- Register routes

Uses marker comments (MCP:IMPORTS:START, MCP:REPOS:START, etc.) in main.go.
scaffold_domain calls this automatically, but use this tool to re-wire or add missing domains.

Use dry_run: true to verify markers exist without making changes.`,
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

		// Special handling for "auth" domain - it uses user repository, not auth repository
		// and is typically scaffolded with scaffold_project --with_auth, not scaffold_domain
		if pkgName == "auth" {
			// Auth is special - skip standard wiring as it uses user repo and has custom initialization
			// The auth service/controller are wired directly in main.go by scaffold_project with_auth
			continue
		}

		// Inject imports with aliases to avoid package name conflicts
		repoImport := fmt.Sprintf("%s/internal/repository/%s", modulePath, pkgName)
		repoAlias := utils.ToRepoImportAlias(domain)
		if err := injector.InjectImportWithAlias(repoImport, repoAlias); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject repo import for '%s': %v", domain, err)), nil
		}

		serviceImport := fmt.Sprintf("%s/internal/services/%s", modulePath, pkgName)
		serviceAlias := utils.ToServiceImportAlias(domain)
		if err := injector.InjectImportWithAlias(serviceImport, serviceAlias); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to inject service import for '%s': %v", domain, err)), nil
		}

		controllerImport := fmt.Sprintf("%s/internal/web/%s", modulePath, pkgName)
		controllerAlias := utils.ToControllerImportAlias(domain)
		if err := injector.InjectImportWithAlias(controllerImport, controllerAlias); err != nil {
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

	suggestedTools := []types.ToolHint{
		types.HintScaffoldDomain,
		types.HintScaffoldSeed,
	}

	return types.ScaffoldResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully updated main.go with wiring for %d domain(s)", len(input.Domains)),
		FilesUpdated:   []string{"cmd/web/main.go"},
		SuggestedTools: suggestedTools,
	}, nil
}
