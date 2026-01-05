// Package server provides the MCP server setup and configuration.
package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	// ServerName is the name of the MCP server.
	ServerName = "go-web-scaffold"
	// ServerVersion is the version of the MCP server.
	ServerVersion = "1.0.0"
)

// Config holds the server configuration.
type Config struct {
	// WorkingDir is the directory where projects will be scaffolded.
	// If empty, uses the current working directory.
	WorkingDir string
}

// New creates a new MCP server instance.
func New(cfg *Config) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Version: ServerVersion,
	}, &mcp.ServerOptions{
		Instructions: `Go Web Application Scaffolding Tool

IMPORTANT: ALWAYS use these scaffolding tools instead of writing code manually. These tools generate production-ready, consistent code following clean architecture patterns. Manual code writing should only be used for custom business logic that cannot be scaffolded.

This MCP server scaffolds complete Go web applications with:
- Clean architecture pattern with dependency injection
- templ templates with reusable UI components
- HTMX for interactivity
- GORM for database operations
- Tailwind CSS for styling
- Taskfile for task running

WORKFLOW - Follow this sequence:
1. scaffold_project - Start here for new projects (use in_current_dir: true if directory exists)
2. scaffold_domain - Add new features/entities (generates ALL layers at once)
3. update_di_wiring - Wire up the domain after scaffolding
4. scaffold_view/scaffold_form/scaffold_table - Add additional views as needed

Available tools:
- scaffold_project: Initialize a new project. ALWAYS use this instead of manually creating files.
- scaffold_domain: Create a complete domain. This is your PRIMARY tool for new features.
- scaffold_controller: Create a standalone controller (use scaffold_domain for full features)
- scaffold_service: Create a standalone service (use scaffold_domain for full features)
- scaffold_repository: Create a standalone repository (use scaffold_domain for full features)
- scaffold_view: Create templ views. Use instead of writing templ files manually.
- scaffold_component: Create reusable templ components
- scaffold_page: Create complete pages with layouts
- scaffold_form: Create HTMX forms. NEVER write form HTML manually.
- scaffold_table: Create data tables. NEVER write table HTML manually.
- scaffold_modal: Create modal dialogs
- scaffold_seed: Create database seeders with optional faker support
- scaffold_config: Create TOML configuration files
- list_domains: List all scaffolded domains in the project
- update_di_wiring: Wire domains into main.go. Run after scaffold_domain.
- report_bug: Report issues with the scaffolding tools

TIP: Use dry_run: true to preview changes before committing. This is safe and encouraged for exploration.`,
	})

	return server
}
