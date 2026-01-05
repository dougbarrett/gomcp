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

// RegisterScaffoldComponent registers the scaffold_component tool.
func RegisterScaffoldComponent(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "scaffold_component",
		Description: "Create reusable templ components with templui patterns. Supports card, modal, dropdown, form_field, table, and custom component types.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldComponentInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldComponent(registry *Registry, input types.ScaffoldComponentInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateComponentName(input.ComponentName); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	// Create generator
	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Prepare template data
	data := buildComponentData(input)

	// Determine output path
	outputDir := filepath.Join("internal", "web", "components")
	outputPath := filepath.Join(outputDir, utils.ToSnakeCase(input.ComponentName)+".templ")

	// Ensure directory exists
	if err := gen.EnsureDir(outputDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Determine template based on component type
	templatePath := getComponentTemplatePath(input.ComponentType)

	// Generate the component file
	if err := gen.GenerateFile(templatePath, outputPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate component: %v", err)), nil
	}

	// Get result
	result := gen.Result()

	nextSteps := []string{
		"templ generate",
		"Import the component where needed",
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create %s component '%s'", input.ComponentType, input.ComponentName),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created %s component '%s'", input.ComponentType, input.ComponentName),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}

// getComponentTemplatePath returns the template path for a component type.
func getComponentTemplatePath(componentType string) string {
	switch componentType {
	case "card":
		return "components/card.templ.tmpl"
	case "form_field":
		return "components/form_field.templ.tmpl"
	default:
		// Default to card for modal, custom, and unknown types
		// Note: For full modal support, use scaffold_modal tool
		return "components/card.templ.tmpl"
	}
}

// buildComponentData creates ComponentData from ScaffoldComponentInput.
func buildComponentData(input types.ScaffoldComponentInput) generator.ComponentData {
	// Build props
	props := make([]generator.PropData, len(input.Props))
	for i, prop := range input.Props {
		props[i] = generator.NewPropData(prop)
	}

	componentType := input.ComponentType
	if componentType == "" {
		componentType = "custom"
	}

	return generator.ComponentData{
		ComponentName: input.ComponentName,
		ComponentType: componentType,
		Props:         props,
		WithHTMX:      input.WithHTMX,
		AlpineState:   input.AlpineState,
	}
}
