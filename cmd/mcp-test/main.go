// Copyright 2025 Doug Barrett. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The mcp-test command is a reusable test harness that connects to an MCP server
// using Claude as the orchestrator. It validates that MCP tools work correctly
// by having Claude execute scaffolding commands and verifying the results.
//
// Usage: mcp-test [flags]
//
// Example:
//
//	export ANTHROPIC_API_KEY=sk-...
//	mcp-test --workdir=/tmp/wizard-test-project --task="scaffold a wizard project"
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	workdir     = flag.String("workdir", "/tmp/mcp-test", "Working directory for scaffolding")
	task        = flag.String("task", "", "Task description for Claude to execute")
	mcpBinary   = flag.String("mcp", "", "Path to MCP server binary (default: gomcp in PATH)")
	model       = flag.String("model", string(anthropic.ModelClaudeSonnet4_5_20250929), "Claude model to use")
	maxTurns    = flag.Int("max-turns", 20, "Maximum conversation turns")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
	verifyBuild = flag.Bool("verify-build", true, "Run 'go build' after scaffolding to verify")
)

func main() {
	flag.Parse()

	if *task == "" {
		fmt.Fprintln(os.Stderr, "Usage: mcp-test --task=\"<task description>\"")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, "  mcp-test --workdir=/tmp/wizard-test --task=\"scaffold a project with wizard\"")
		os.Exit(2)
	}

	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	// Setup working directory
	if err := os.MkdirAll(*workdir, 0755); err != nil {
		return fmt.Errorf("creating workdir: %w", err)
	}
	log.Printf("Working directory: %s", *workdir)

	// Find MCP binary
	mcpPath := *mcpBinary
	if mcpPath == "" {
		var err error
		mcpPath, err = exec.LookPath("gomcp")
		if err != nil {
			return fmt.Errorf("gomcp not found in PATH, use --mcp flag to specify: %w", err)
		}
	}
	log.Printf("MCP server: %s", mcpPath)

	// Start MCP server
	cmd := exec.Command(mcpPath)
	cmd.Dir = *workdir
	transport := &mcp.CommandTransport{Command: cmd}

	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-test",
		Version: "v1.0.0",
	}, nil)

	session, err := mcpClient.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("connecting to MCP server: %w", err)
	}
	defer session.Close()
	log.Printf("Connected to MCP server")

	// Get available tools
	tools, err := collectTools(ctx, session)
	if err != nil {
		return fmt.Errorf("getting tools: %w", err)
	}
	log.Printf("Found %d MCP tools", len(tools))

	if *verbose {
		for _, tool := range tools {
			log.Printf("  - %s: %s", tool.Name, truncate(tool.Description, 60))
		}
	}

	// Convert MCP tools to Anthropic tool format
	anthropicTools := convertToAnthropicTools(tools)

	// Create Claude client
	claudeClient := anthropic.NewClient()

	// Build system prompt
	systemPrompt := buildSystemPrompt()

	// Run conversation loop
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(*task)),
	}

	for turn := 0; turn < *maxTurns; turn++ {
		log.Printf("\n=== Turn %d ===", turn+1)

		// Call Claude
		response, err := claudeClient.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     anthropic.Model(*model),
			MaxTokens: 4096,
			System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
			Messages:  messages,
			Tools:     anthropicTools,
		})
		if err != nil {
			return fmt.Errorf("calling Claude: %w", err)
		}

		// Process response
		var toolUses []anthropic.ToolUseBlock
		var textContent strings.Builder

		for _, block := range response.Content {
			switch b := block.AsAny().(type) {
			case anthropic.TextBlock:
				textContent.WriteString(b.Text)
				if *verbose {
					log.Printf("[Claude]: %s", b.Text)
				}
			case anthropic.ToolUseBlock:
				toolUses = append(toolUses, b)
				log.Printf("[Tool Call]: %s", b.Name)
				if *verbose {
					inputJSON, _ := json.MarshalIndent(b.Input, "", "  ")
					log.Printf("  Input: %s", string(inputJSON))
				}
			}
		}

		// Add assistant response to messages
		messages = append(messages, response.ToParam())

		// Check if done
		if response.StopReason == "end_turn" && len(toolUses) == 0 {
			log.Printf("\n=== Claude completed task ===")
			if textContent.Len() > 0 {
				fmt.Println(textContent.String())
			}
			break
		}

		// Execute tool calls
		if len(toolUses) > 0 {
			var toolResults []anthropic.ContentBlockParamUnion

			for _, toolUse := range toolUses {
				result, err := executeToolCall(ctx, session, toolUse)
				if err != nil {
					log.Printf("[Tool Error]: %s: %v", toolUse.Name, err)
					toolResults = append(toolResults, anthropic.NewToolResultBlock(
						toolUse.ID,
						fmt.Sprintf("Error: %v", err),
						true, // is_error
					))
				} else {
					if *verbose {
						log.Printf("[Tool Result]: %s", truncate(result, 500))
					}
					toolResults = append(toolResults, anthropic.NewToolResultBlock(
						toolUse.ID,
						result,
						false,
					))
				}
			}

			messages = append(messages, anthropic.NewUserMessage(toolResults...))
		}
	}

	// Verify build if requested
	if *verifyBuild {
		if err := verifyGoBuild(ctx); err != nil {
			return fmt.Errorf("build verification failed: %w", err)
		}
	}

	log.Printf("\n=== Test completed successfully ===")
	return nil
}

func collectTools(ctx context.Context, session *mcp.ClientSession) ([]*mcp.Tool, error) {
	var tools []*mcp.Tool
	for tool, err := range session.Tools(ctx, nil) {
		if err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}
	return tools, nil
}

func convertToAnthropicTools(mcpTools []*mcp.Tool) []anthropic.ToolUnionParam {
	var tools []anthropic.ToolUnionParam
	for _, t := range mcpTools {
		// Convert MCP InputSchema to JSON
		schemaJSON, err := json.Marshal(t.InputSchema)
		if err != nil {
			log.Printf("Warning: skipping tool %s due to schema error: %v", t.Name, err)
			continue
		}

		// Parse as raw JSON for Anthropic
		var schema anthropic.ToolInputSchemaParam
		if err := json.Unmarshal(schemaJSON, &schema); err != nil {
			log.Printf("Warning: skipping tool %s due to schema parse error: %v", t.Name, err)
			continue
		}

		description := t.Description
		tool := anthropic.ToolParam{
			Name:        t.Name,
			Description: anthropic.String(description),
			InputSchema: schema,
		}
		tools = append(tools, anthropic.ToolUnionParam{OfTool: &tool})
	}
	return tools
}

func executeToolCall(ctx context.Context, session *mcp.ClientSession, toolUse anthropic.ToolUseBlock) (string, error) {
	// Parse input JSON to map
	var inputMap map[string]any
	if err := json.Unmarshal(toolUse.Input, &inputMap); err != nil {
		return "", fmt.Errorf("parsing tool input: %w", err)
	}

	// Call MCP tool
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolUse.Name,
		Arguments: inputMap,
	})
	if err != nil {
		return "", err
	}

	// Extract text content from result
	var output strings.Builder
	for _, content := range result.Content {
		switch c := content.(type) {
		case *mcp.TextContent:
			output.WriteString(c.Text)
			output.WriteString("\n")
		default:
			// Handle other content types if needed
			j, _ := json.Marshal(content)
			output.WriteString(string(j))
			output.WriteString("\n")
		}
	}

	if result.IsError {
		return output.String(), fmt.Errorf("tool returned error")
	}

	return output.String(), nil
}

func buildSystemPrompt() string {
	return `You are a test automation assistant validating MCP scaffolding tools.

Your task is to use the available MCP tools to scaffold a Go web application project.

IMPORTANT GUIDELINES:
1. Use scaffold_project first to create the project structure
2. Use scaffold_domain to create any required domains (models, repositories, services, controllers)
3. Use scaffold_wizard if creating wizard functionality
4. Always use dry_run: false to actually create files (this is a test environment)
5. Be methodical - create dependencies before dependent items
6. After scaffolding, report what was created

The working directory is already set up for scaffolding. You can use the tools directly.

When you're done scaffolding, summarize:
- What was created
- Any errors encountered
- Files that were generated`
}

func verifyGoBuild(ctx context.Context) error {
	// Check if go.mod exists
	goModPath := filepath.Join(*workdir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		log.Printf("No go.mod found, skipping build verification")
		return nil
	}

	log.Printf("\n=== Verifying Go build ===")

	// Run go mod tidy
	tidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	tidyCmd.Dir = *workdir
	tidyOutput, err := tidyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %w\n%s", err, string(tidyOutput))
	}
	log.Printf("go mod tidy: OK")

	// Run go build
	buildCmd := exec.CommandContext(ctx, "go", "build", "./...")
	buildCmd.Dir = *workdir
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed: %w\n%s", err, string(buildOutput))
	}
	log.Printf("go build ./...: OK")

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
