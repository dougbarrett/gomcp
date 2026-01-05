package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterScaffoldPage registers the scaffold_page tool.
func RegisterScaffoldPage(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_page",
		Description: `Create complete pages with layout and optional TOML configuration.

Layouts: default, dashboard, landing, blank

Set create_toml_config: true to generate a TOML file for i18n content.

Sections allow composing pages from: hero, content, table, cards, form.

Run 'templ generate' after creating pages.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldPageInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldPage(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldPage(registry *Registry, input types.ScaffoldPageInput) (types.ScaffoldResult, error) {
	// Validate input
	if input.PageName == "" {
		return types.NewErrorResult("page name is required"), nil
	}

	if input.Route == "" {
		return types.NewErrorResult("route is required"), nil
	}

	if err := utils.ValidateURLPath(input.Route); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if err := utils.ValidateLayoutType(input.Layout); err != nil {
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
	data := buildPageData(input, modulePath)

	// Determine page output path
	pageDir := filepath.Join("internal", "web", "pages")
	pagePath := filepath.Join(pageDir, utils.ToSnakeCase(input.PageName)+".templ")

	// Ensure directory exists
	if err := gen.EnsureDir(pageDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate the page file using list template as base for now
	// In a full implementation, we'd have a dedicated page.templ.tmpl
	if err := gen.GenerateFile("views/list.templ.tmpl", pagePath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate page: %v", err)), nil
	}

	// Generate TOML config if requested
	if input.CreateTomlConfig {
		configDir := filepath.Join("config", "en", "pages")
		configPath := filepath.Join(configDir, utils.ToSnakeCase(input.PageName)+".toml")

		if err := gen.EnsureDir(configDir); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to create config directory: %v", err)), nil
		}

		configData := buildPageConfigData(input)
		if err := gen.GenerateFile("config/page.toml.tmpl", configPath, configData); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate config: %v", err)), nil
		}
	}

	// Get result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	nextSteps := []string{
		"templ generate",
		fmt.Sprintf("Add route handler for '%s' in your router", input.Route),
	}

	suggestedTools := []types.ToolHint{
		{
			Tool:        "scaffold_component",
			Description: "Create reusable UI components for this page",
			Priority:    "optional",
		},
		{
			Tool:        "scaffold_modal",
			Description: "Add modal dialogs to this page",
			Priority:    "optional",
		},
		types.HintScaffoldDomain,
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:        true,
			Message:        fmt.Sprintf("Dry run: Would create page '%s' at route '%s'", input.PageName, input.Route),
			FilesCreated:   result.FilesCreated,
			NextSteps:      nextSteps,
			SuggestedTools: suggestedTools,
		}, nil
	}

	return types.ScaffoldResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully created page '%s' at route '%s'", input.PageName, input.Route),
		FilesCreated:   result.FilesCreated,
		FilesUpdated:   result.FilesUpdated,
		NextSteps:      nextSteps,
		SuggestedTools: suggestedTools,
	}, nil
}

// buildPageData creates PageData from ScaffoldPageInput.
func buildPageData(input types.ScaffoldPageInput, modulePath string) generator.PageData {
	layout := input.Layout
	if layout == "" {
		layout = "default"
	}

	// Build sections
	sections := make([]generator.SectionData, len(input.Sections))
	for i, section := range input.Sections {
		sections[i] = generator.SectionData{
			Type:   section.Type,
			Config: section.Config,
		}
	}

	// Generate title from page name
	title := utils.ToLabel(input.PageName)
	modelName := utils.ToPascalCase(input.PageName)
	pkgName := utils.ToPackageName(input.PageName)
	varName := utils.ToVariableName(input.PageName)

	urlPathSegment := strings.TrimPrefix(input.Route, "/")
	return generator.PageData{
		ModulePath:        modulePath,
		PageName:          input.PageName,
		ModelName:         modelName,
		PackageName:       pkgName,
		VariableName:      varName,
		URLPath:           input.Route,
		URLPathSegment:    urlPathSegment,
		Route:             input.Route,
		Layout:            layout,
		Sections:          sections,
		Title:             title,
		Description:       fmt.Sprintf("%s page", title),
		Fields:            []generator.FieldData{},
		WithPagination:    false,
		WithSearch:        false,
		EmptyStateMessage: "",
	}
}

// buildPageConfigData creates ConfigData for page TOML.
func buildPageConfigData(input types.ScaffoldPageInput) generator.ConfigData {
	title := utils.ToLabel(input.PageName)
	layout := input.Layout
	if layout == "" {
		layout = "default"
	}

	content := map[string]interface{}{
		"title":       title,
		"description": fmt.Sprintf("%s page", title),
		"heading":     title,
	}

	return generator.ConfigData{
		ConfigType:  "page",
		Name:        input.PageName,
		PageName:    input.PageName,
		Locale:      "en",
		Content:     content,
		Title:       title,
		Description: fmt.Sprintf("%s page", title),
		Heading:     title,
		Layout:      layout,
	}
}
