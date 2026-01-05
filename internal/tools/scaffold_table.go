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

// RegisterScaffoldTable registers the scaffold_table tool.
func RegisterScaffoldTable(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_table",
		Description: `NEVER write table HTML manually. Use this tool instead.

Generates data tables with Tailwind CSS styling featuring:
- HTMX-powered sorting (click column headers)
- Pagination with page size controls
- Search functionality
- Row actions (view, edit, delete, custom)
- Bulk actions for batch operations
- Column formatting: text, currency, date, datetime, badge, link

Columns support sortable: true and custom badge_config for status fields.

Run 'templ generate' after creating tables.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldTableInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldTable(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldTable(registry *Registry, input types.ScaffoldTableInput) (types.ScaffoldResult, error) {
	// Validate input
	if input.TableName == "" {
		return types.NewErrorResult("table name is required"), nil
	}

	// Support nested paths like "admin/users"
	if err := utils.ValidateDomainPath(input.Domain); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if len(input.Columns) == 0 {
		return types.NewErrorResult("at least one column is required"), nil
	}

	// Validate columns
	for _, col := range input.Columns {
		if col.Key == "" {
			return types.NewErrorResult("column key is required"), nil
		}
	}

	// Validate row actions
	validActionTypes := map[string]bool{"view": true, "edit": true, "delete": true, "custom": true}
	for _, action := range input.RowActions {
		if !validActionTypes[action.Type] {
			return types.NewErrorResult(fmt.Sprintf("invalid row action type '%s'", action.Type)), nil
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

	// For nested paths like "admin/users", use base name for package/model
	baseDomain := utils.ParseDomainPath(input.Domain)
	domainDir := utils.DomainPathToDir(input.Domain)
	basePkgName := utils.ToPackageName(baseDomain)

	// Prepare template data
	data := buildTableData(input, modulePath, baseDomain)

	// Determine output path - use full path for directory structure
	viewDir := filepath.Join("internal", "web", domainDir, "views")
	outputPath := filepath.Join(viewDir, input.TableName+".templ")

	// Ensure directory exists
	if err := gen.EnsureDir(viewDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate the table file using the table template
	if err := gen.GenerateFile("views/table.templ.tmpl", outputPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate table: %v", err)), nil
	}

	// Get result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	nextSteps := []string{
		"templ generate",
		fmt.Sprintf("Import the table in internal/web/%s/%s.go", domainDir, basePkgName),
		fmt.Sprintf("Add list handler that uses the table in the %s controller", basePkgName),
	}

	suggestedTools := []types.ToolHint{
		{
			Tool:        "scaffold_form",
			Description: fmt.Sprintf("Create forms for creating/editing %s records", input.Domain),
			Priority:    "optional",
		},
		{
			Tool:        "scaffold_seed",
			Description: fmt.Sprintf("Create a seeder to populate %s with test data", input.Domain),
			Priority:    "optional",
		},
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:        true,
			Message:        fmt.Sprintf("Dry run: Would create table '%s' for domain '%s'", input.TableName, input.Domain),
			FilesCreated:   result.FilesCreated,
			NextSteps:      nextSteps,
			SuggestedTools: suggestedTools,
		}, nil
	}

	return types.ScaffoldResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully created table '%s' for domain '%s'", input.TableName, input.Domain),
		FilesCreated:   result.FilesCreated,
		FilesUpdated:   result.FilesUpdated,
		NextSteps:      nextSteps,
		SuggestedTools: suggestedTools,
	}, nil
}

// buildTableData creates TableData from ScaffoldTableInput.
// baseDomain is the base name extracted from nested paths (e.g., "users" from "admin/users").
func buildTableData(input types.ScaffoldTableInput, modulePath, baseDomain string) generator.TableData {
	modelName := utils.ToModelName(baseDomain)
	pkgName := utils.ToPackageName(baseDomain)
	varName := utils.ToVariableName(baseDomain)
	urlPath := utils.ToURLPath(input.Domain) // Keep full path for URL

	// Build columns
	columns := generator.NewColumnDataList(input.Columns)

	// Build row actions
	rowActions := make([]generator.RowActionData, len(input.RowActions))
	for i, action := range input.RowActions {
		rowActions[i] = generator.NewRowActionData(action)
	}

	// Get boolean defaults (nil means true for pagination, sorting, search)
	withPagination := true
	if input.WithPagination != nil {
		withPagination = *input.WithPagination
	}

	withSorting := true
	if input.WithSorting != nil {
		withSorting = *input.WithSorting
	}

	withSearch := true
	if input.WithSearch != nil {
		withSearch = *input.WithSearch
	}

	return generator.TableData{
		ModulePath:      modulePath,
		DomainName:      input.Domain,
		ModelName:       modelName,
		PackageName:     pkgName,
		VariableName:    varName,
		TableName:       input.TableName,
		URLPath:         urlPath,
		Columns:         columns,
		WithPagination:  withPagination,
		WithSorting:     withSorting,
		WithSearch:      withSearch,
		WithBulkActions: input.WithBulkActions,
		RowActions:      rowActions,
	}
}
