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

// RegisterScaffoldView registers the scaffold_view tool.
func RegisterScaffoldView(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "scaffold_view",
		Description: "Create templ views using templui components. Supports list, show, form, card, table, and custom view types.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldViewInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldView(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldView(registry *Registry, input types.ScaffoldViewInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateDomainName(input.DomainName); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if err := utils.ValidateViewType(input.ViewType); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if input.ViewName == "" {
		return types.NewErrorResult("view name is required"), nil
	}

	// Validate fields if provided
	for _, field := range input.Config.Fields {
		if err := utils.ValidateFieldName(field.Name); err != nil {
			return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
		}
		if err := utils.ValidateFieldType(field.Type); err != nil {
			return types.NewErrorResult(fmt.Sprintf("field '%s': %v", field.Name, err)), nil
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
	data := buildViewData(input, modulePath)

	// Determine output path
	pkgName := utils.ToPackageName(input.DomainName)
	viewDir := filepath.Join("internal", "web", pkgName, "views")
	outputPath := filepath.Join(viewDir, input.ViewName+".templ")

	// Ensure directory exists
	if err := gen.EnsureDir(viewDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Determine template based on view type
	templatePath := getViewTemplatePath(input.ViewType)

	// Generate the view file
	if err := gen.GenerateFile(templatePath, outputPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate view: %v", err)), nil
	}

	// Get result
	result := gen.Result()

	nextSteps := []string{
		"templ generate",
		fmt.Sprintf("Import the view in internal/web/%s/%s.go", pkgName, pkgName),
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create %s view '%s' for domain '%s'", input.ViewType, input.ViewName, input.DomainName),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created %s view '%s' for domain '%s'", input.ViewType, input.ViewName, input.DomainName),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}

// getViewTemplatePath returns the template path for a view type.
func getViewTemplatePath(viewType string) string {
	switch viewType {
	case "list":
		return "views/list.templ.tmpl"
	case "show":
		return "views/show.templ.tmpl"
	case "form":
		return "views/form.templ.tmpl"
	case "table":
		return "views/table.templ.tmpl"
	case "card":
		// Card views use partials template
		return "views/partials.templ.tmpl"
	default:
		// Custom views use list as base
		return "views/list.templ.tmpl"
	}
}

// buildViewData creates ViewData from ScaffoldViewInput.
func buildViewData(input types.ScaffoldViewInput, modulePath string) generator.ViewData {
	modelName := utils.ToModelName(input.DomainName)
	pkgName := utils.ToPackageName(input.DomainName)
	varName := utils.ToVariableName(input.DomainName)
	urlPath := utils.ToURLPath(input.DomainName)

	// Build fields
	fields := generator.NewFieldDataList(input.Config.Fields)

	// Build columns
	columns := generator.NewColumnDataList(input.Config.Columns)

	// Set defaults
	emptyStateMsg := input.Config.EmptyStateMessage
	if emptyStateMsg == "" {
		emptyStateMsg = fmt.Sprintf("No %s found", utils.Pluralize(input.DomainName))
	}

	method := input.Config.Method
	if method == "" {
		method = "POST"
	}

	submitURL := input.Config.SubmitURL
	if submitURL == "" {
		submitURL = urlPath
	}

	successRedirect := input.Config.SuccessRedirect
	if successRedirect == "" {
		successRedirect = urlPath
	}

	// Build row actions
	rowActions := make([]generator.RowActionData, len(input.Config.RowActions))
	for i, action := range input.Config.RowActions {
		rowActions[i] = generator.NewRowActionData(action)
	}

	return generator.ViewData{
		ModulePath:        modulePath,
		DomainName:        input.DomainName,
		ModelName:         modelName,
		PackageName:       pkgName,
		VariableName:      varName,
		URLPath:           urlPath,
		ViewType:          input.ViewType,
		ViewName:          input.ViewName,
		Fields:            fields,
		Columns:           columns,
		WithPagination:    input.Config.WithPagination,
		WithSearch:        input.Config.WithSearch,
		WithFilters:       input.Config.WithFilters,
		WithSorting:       input.Config.WithSorting,
		WithBulkActions:   input.Config.WithBulkActions,
		WithSoftDelete:    input.Config.WithSoftDelete,
		RowActions:        rowActions,
		EmptyStateMessage: emptyStateMsg,
		SubmitURL:         submitURL,
		Method:            method,
		SuccessRedirect:   successRedirect,
	}
}
