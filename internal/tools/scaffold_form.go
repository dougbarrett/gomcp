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

// RegisterScaffoldForm registers the scaffold_form tool.
func RegisterScaffoldForm(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_form",
		Description: `NEVER write form HTML manually. Use this tool instead.

Generates HTMX-powered forms with templui components that:
- Auto-submit via HTMX with loading states
- Display validation errors inline
- Support all field types: input, textarea, select, checkbox, date, email, password, number, rating, tags, slider

IMPORTANT: Each form file handles BOTH create and edit operations using the IsEdit prop.
Only scaffold ONE form per domain - it generates FormCreate() and FormEdit() helper functions.
Do NOT scaffold separate forms for create and edit actions.

Specify action: 'create' or 'edit' to set the default form behavior.
Set submit_endpoint for the HTMX post URL.

Run 'templ generate' after creating forms.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldFormInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldForm(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldForm(registry *Registry, input types.ScaffoldFormInput) (types.ScaffoldResult, error) {
	// Validate input
	if input.FormName == "" {
		return types.NewErrorResult("form name is required"), nil
	}

	if err := utils.ValidateDomainName(input.Domain); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if input.Action != "create" && input.Action != "edit" {
		return types.NewErrorResult("action must be 'create' or 'edit'"), nil
	}

	if len(input.Fields) == 0 {
		return types.NewErrorResult("at least one field is required"), nil
	}

	// Validate each field
	for _, field := range input.Fields {
		if err := utils.ValidateFieldName(field.Name); err != nil {
			return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
		}
		if field.FormType != "" {
			if err := utils.ValidateFormType(field.FormType); err != nil {
				return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
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
	data := generator.NewFormData(input, modulePath)

	// Determine output path
	pkgName := utils.ToPackageName(input.Domain)
	viewDir := filepath.Join("internal", "web", pkgName, "views")
	outputPath := filepath.Join(viewDir, input.FormName+".templ")

	// Ensure directory exists
	if err := gen.EnsureDir(viewDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate the form file using the form template
	if err := gen.GenerateFile("views/form.templ.tmpl", outputPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate form: %v", err)), nil
	}

	// Get result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	nextSteps := []string{
		"templ generate",
		fmt.Sprintf("Import the form in internal/web/%s/%s.go", pkgName, pkgName),
		fmt.Sprintf("Add form handler in the %s controller", pkgName),
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create %s form '%s' for domain '%s'", input.Action, input.FormName, input.Domain),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created %s form '%s' for domain '%s'", input.Action, input.FormName, input.Domain),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}
