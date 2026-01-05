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

// RegisterScaffoldRepository registers the scaffold_repository tool.
func RegisterScaffoldRepository(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_repository",
		Description: `Create a standalone GORM repository. Use scaffold_domain instead for complete features.

Only use this tool when you need JUST a repository layer without service/controller. Common cases:
- Adding a repository for an existing model
- Creating a shared data access layer

Generates: internal/repository/{domain}/{domain}.go with standard CRUD operations.

Prefer scaffold_domain for new features - it generates all layers consistently.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldRepositoryInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldRepository(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldRepository(registry *Registry, input types.ScaffoldRepositoryInput) (types.ScaffoldResult, error) {
	// Validate input
	if err := utils.ValidateDomainName(input.DomainName); err != nil {
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
	modelName := input.ModelName
	if modelName == "" {
		modelName = utils.ToModelName(input.DomainName)
	}

	pkgName := utils.ToPackageName(input.DomainName)

	data := generator.DomainData{
		ModulePath:   modulePath,
		DomainName:   input.DomainName,
		ModelName:    modelName,
		PackageName:  pkgName,
		VariableName: utils.ToVariableName(input.DomainName),
		TableName:    utils.ToTableName(input.DomainName),
		URLPath:      utils.ToURLPath(input.DomainName),
	}

	// Create directory
	repoDir := filepath.Join("internal", "repository", pkgName)
	if err := gen.EnsureDir(repoDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate repository
	repoPath := filepath.Join(repoDir, pkgName+".go")
	if err := gen.GenerateFile("domain/repository.go.tmpl", repoPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate repository: %v", err)), nil
	}

	// Prepare result
	result := gen.Result()

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create repository for '%s'", input.DomainName),
			FilesCreated: result.FilesCreated,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created repository for '%s'", input.DomainName),
		FilesCreated: result.FilesCreated,
	}, nil
}
