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

// RegisterScaffoldService registers the scaffold_service tool.
func RegisterScaffoldService(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_service",
		Description: `Create a standalone service layer. Use scaffold_domain instead for complete features.

Only use this tool when you need JUST a service layer. Common cases:
- Adding business logic for an existing repository
- Creating cross-domain services

Generates:
- internal/services/{domain}/{domain}.go (service with repository dependency)
- internal/services/{domain}/dto.go (Create/Update DTOs)

Prefer scaffold_domain for new features - it generates all layers consistently.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldServiceInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldService(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldService(registry *Registry, input types.ScaffoldServiceInput) (types.ScaffoldResult, error) {
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
	pkgName := utils.ToPackageName(input.DomainName)

	data := generator.DomainData{
		ModulePath:   modulePath,
		DomainName:   input.DomainName,
		ModelName:    utils.ToModelName(input.DomainName),
		PackageName:  pkgName,
		VariableName: utils.ToVariableName(input.DomainName),
		TableName:    utils.ToTableName(input.DomainName),
		URLPath:      utils.ToURLPath(input.DomainName),
	}

	// Create directory
	serviceDir := filepath.Join("internal", "services", pkgName)
	if err := gen.EnsureDir(serviceDir); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to create directory: %v", err)), nil
	}

	// Generate service
	servicePath := filepath.Join(serviceDir, pkgName+".go")
	if err := gen.GenerateFile("domain/service.go.tmpl", servicePath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate service: %v", err)), nil
	}

	// Generate DTOs
	dtoPath := filepath.Join(serviceDir, "dto.go")
	if err := gen.GenerateFile("domain/dto.go.tmpl", dtoPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate DTOs: %v", err)), nil
	}

	// Prepare result
	result := gen.Result()

	if input.DryRun {
		return types.ScaffoldResult{
			Success:      true,
			Message:      fmt.Sprintf("Dry run: Would create service for '%s'", input.DomainName),
			FilesCreated: result.FilesCreated,
		}, nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully created service for '%s'", input.DomainName),
		FilesCreated: result.FilesCreated,
	}, nil
}
