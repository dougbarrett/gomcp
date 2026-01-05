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

// RegisterScaffoldSeed registers the scaffold_seed tool.
func RegisterScaffoldSeed(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_seed",
		Description: `Generate database seeders for test data.

Features:
- with_faker: true for realistic fake data (names, emails, addresses, etc.)
- count: Number of records to seed (default: 10)
- dependencies: Other seeders to run first (e.g., ["user"] before "order")

Register the seeder in cmd/seed/main.go after generating.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldSeedInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldSeed(registry *Registry, input types.ScaffoldSeedInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateDomainName(input.Domain); err != nil {
		return types.NewErrorResult(err.Error()), nil
	}

	// Validate dependencies
	for _, dep := range input.Dependencies {
		if err := utils.ValidateDomainName(dep); err != nil {
			return types.NewErrorResult(fmt.Sprintf("dependency '%s': %v", dep, err)), nil
		}
	}

	// Validate fields if provided
	for _, field := range input.Fields {
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
	data := buildSeedData(input, modulePath)

	// Determine output path
	seedDir := filepath.Join("cmd", "seed", "seeders")
	outputPath := filepath.Join(seedDir, utils.ToSnakeCase(input.Domain)+"_seeder.go")

	// Ensure directory exists
	if err := gen.EnsureDir(seedDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate the seeder file
	if err := gen.GenerateFile("seed/seeder.go.tmpl", outputPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate seeder: %v", err)), nil
	}

	// Get result
	result := gen.Result()

	nextSteps := []string{
		"go mod tidy",
		fmt.Sprintf("Register the seeder in cmd/seed/main.go"),
	}

	if input.WithFaker {
		nextSteps = append([]string{"go get github.com/brianvoe/gofakeit/v6"}, nextSteps...)
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create seeder for domain '%s'", input.Domain),
			FilesCreated: result.FilesCreated,
			NextSteps:    nextSteps,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created seeder for domain '%s'", input.Domain),
		FilesCreated: result.FilesCreated,
		FilesUpdated: result.FilesUpdated,
		NextSteps:    nextSteps,
	}, nil
}

// buildSeedData creates SeedData from ScaffoldSeedInput.
func buildSeedData(input types.ScaffoldSeedInput, modulePath string) generator.SeedData {
	modelName := utils.ToModelName(input.Domain)
	tableName := utils.ToTableName(input.Domain)

	// Set default count
	count := input.Count
	if count <= 0 {
		count = 10
	}

	// Build fields
	fields := generator.NewFieldDataList(input.Fields)

	return generator.SeedData{
		ModulePath:   modulePath,
		DomainName:   input.Domain,
		ModelName:    modelName,
		TableName:    tableName,
		Fields:       fields,
		Count:        count,
		WithFaker:    input.WithFaker,
		Dependencies: input.Dependencies,
	}
}
