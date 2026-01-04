package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAddTempluiComponent registers the add_templui_component tool.
func RegisterAddTempluiComponent(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_templui_component",
		Description: "Install templui components into the project. Available components include: button, card, input, select, checkbox, dialog, dropdown, table, tabs, toast, and more.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.AddTempluiComponentInput) (*mcp.CallToolResult, types.AddComponentResult, error) {
		result, err := addTempluiComponent(registry, input)
		if err != nil {
			return nil, types.NewAddComponentError(err.Error(), nil), nil
		}
		return nil, result, nil
	})
}

func addTempluiComponent(registry *Registry, input types.AddTempluiComponentInput) (types.AddComponentResult, error) {
	if len(input.Components) == 0 {
		return types.NewAddComponentError("at least one component is required", nil), nil
	}

	var added []string
	var skipped []string
	var errors []string

	for _, component := range input.Components {
		// Run templui add command
		args := []string{"add", component}
		if input.Force {
			args = append(args, "--force")
		}

		cmd := exec.Command("templui", args...)
		cmd.Dir = registry.WorkingDir

		output, err := cmd.CombinedOutput()
		if err != nil {
			// Check if it's because component already exists
			if strings.Contains(string(output), "already exists") {
				skipped = append(skipped, component)
			} else {
				errors = append(errors, fmt.Sprintf("%s: %s", component, strings.TrimSpace(string(output))))
			}
			continue
		}

		added = append(added, component)
	}

	if len(errors) > 0 && len(added) == 0 {
		return types.NewAddComponentError("failed to add components", errors), nil
	}

	result := types.NewAddComponentResult(added, skipped)
	if len(errors) > 0 {
		result.Errors = errors
	}

	if len(added) > 0 {
		result.Message = fmt.Sprintf("Added %d component(s)", len(added))
		if len(skipped) > 0 {
			result.Message += fmt.Sprintf(", skipped %d (already exist)", len(skipped))
		}
	} else if len(skipped) > 0 {
		result.Message = fmt.Sprintf("All %d component(s) already exist", len(skipped))
	}

	return result, nil
}
