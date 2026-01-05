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

// RegisterScaffoldModal registers the scaffold_modal tool.
func RegisterScaffoldModal(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_modal",
		Description: `Create modal dialogs with HTMX content loading.

Modal types:
- dialog: Standard centered modal
- sheet: Slide-in panel from edge
- confirm: Confirmation dialog with actions

Features:
- HTMX loading for dynamic content (set htmx_url in trigger_config)
- Customizable trigger button (text, variant)
- Content types: form, info, confirm

Run 'templ generate' after creating modals.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldModalInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldModal(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldModal(registry *Registry, input types.ScaffoldModalInput) (types.ScaffoldResult, error) {
	// Validate input
	if input.ModalName == "" {
		return types.NewErrorResult("modal name is required"), nil
	}

	if err := utils.ValidateModalType(input.ModalType); err != nil {
		return types.NewErrorResult(err.Error()), nil
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
	data := generator.NewModalData(modulePath, input)

	// Determine output path - modals go to shared components directory
	outputDir := filepath.Join("internal", "web", "components")
	outputPath := filepath.Join(outputDir, input.ModalName+".templ")

	// Ensure directory exists
	if err := gen.EnsureDir(outputDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate the modal file using the modal template
	if err := gen.GenerateFile("components/modal.templ.tmpl", outputPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate modal: %v", err)), nil
	}

	// Get result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	nextSteps := []string{
		"templ generate",
		"Import the modal component where needed",
	}

	if data.HTMXURL != "" {
		nextSteps = append(nextSteps, fmt.Sprintf("Add handler for HTMX endpoint: %s", data.HTMXURL))
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create %s modal '%s'", input.ModalType, input.ModalName),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created %s modal '%s'", input.ModalType, input.ModalName),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}
