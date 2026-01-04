package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterScaffoldConfig registers the scaffold_config tool.
func RegisterScaffoldConfig(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "scaffold_config",
		Description: "Generate TOML configuration files for pages, menus, app settings, or messages.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldConfigInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldConfig(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldConfig(registry *Registry, input types.ScaffoldConfigInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateConfigType(input.ConfigType); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	if input.Name == "" {
		return types.NewErrorResult("config name is required"), nil
	}

	locale := input.Locale
	if locale == "" {
		locale = "en"
	}

	if err := utils.ValidateLocale(locale); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	// Create generator
	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Determine config path based on type
	var configPath string
	switch input.ConfigType {
	case "page":
		configPath = filepath.Join("config", locale, "pages", input.Name+".toml")
	case "menu":
		configPath = filepath.Join("config", locale, "menu.toml")
	case "app":
		configPath = filepath.Join("config", locale, "app.toml")
	case "messages":
		configPath = filepath.Join("config", locale, "messages", input.Name+".toml")
	}

	// Generate basic TOML structure
	content := generateTOMLContent(input.ConfigType, input.Name, input.Content)

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := gen.EnsureDir(dir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Write config file
	if err := gen.GenerateFileFromString(configPath, content); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate config: %v", err)), nil
	}

	result := gen.Result()

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create %s config '%s'", input.ConfigType, input.Name),
			FilesCreated: result.FilesCreated,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created %s config '%s'", input.ConfigType, input.Name),
		FilesCreated: result.FilesCreated,
	}, nil
}

// generateTOMLContent generates TOML content based on config type.
func generateTOMLContent(configType, name string, content map[string]interface{}) string {
	switch configType {
	case "page":
		return fmt.Sprintf(`[meta]
title = "%s"
description = ""

[content]
heading = "%s"
empty_state = "No items found"
`, utils.ToLabel(name), utils.ToLabel(name))

	case "menu":
		return `[[main]]
label = "Dashboard"
url = "/dashboard"
icon = "home"
order = 1

[[main]]
label = "Settings"
url = "/settings"
icon = "settings"
order = 100
`

	case "app":
		return `[server]
address = ":8080"
debug = true

[database]
driver = "sqlite"
dsn = "data.db"
`

	case "messages":
		return fmt.Sprintf(`[%s]
# Add your messages here
`, name)

	default:
		return ""
	}
}
