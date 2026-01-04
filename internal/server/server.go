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

This MCP server provides tools to scaffold complete Go web applications with:
- Clean architecture pattern with dependency injection
- templ templates with templui components
- HTMX for interactivity
- GORM for database operations
- Tailwind CSS for styling
- Taskfile for task running

Available tools:
- scaffold_project: Initialize a new project with complete directory structure
- scaffold_domain: Create a complete domain (model, repository, service, controller, views)
- scaffold_controller: Create a standalone controller
- scaffold_service: Create a standalone service
- scaffold_repository: Create a standalone repository
- scaffold_view: Create templ views (list, show, form, table, etc.)
- scaffold_component: Create reusable templ components
- scaffold_page: Create complete pages with layouts
- scaffold_form: Create HTMX-powered forms
- scaffold_table: Create data tables with sorting/pagination
- scaffold_modal: Create modal dialogs
- scaffold_seed: Create database seeders
- scaffold_config: Create TOML configuration files
- add_templui_component: Install templui components
- list_domains: List all scaffolded domains
- update_di_wiring: Update main.go dependency injection

All tools support a dry_run parameter to preview changes without writing files.`,
	})

	return server
}
