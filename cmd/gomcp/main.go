// Package main is the entry point for the MCP scaffolding server.
package main

import (
	"context"
	"log"
	"os"

	"github.com/dbb1dev/go-mcp/internal/server"
	"github.com/dbb1dev/go-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Get working directory from environment or use current directory
	workingDir := os.Getenv("MCP_SCAFFOLD_WORKDIR")
	if workingDir == "" {
		var err error
		workingDir, err = os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
	}

	// Create server configuration
	cfg := &server.Config{
		WorkingDir: workingDir,
	}

	// Create MCP server
	srv := server.New(cfg)

	// Create tool registry and register all tools
	registry := tools.NewRegistry(workingDir)
	registry.RegisterAll(srv)

	// Run the server with stdio transport
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
