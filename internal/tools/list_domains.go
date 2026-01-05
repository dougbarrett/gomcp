package tools

import (
	"context"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterListDomains registers the list_domains tool.
func RegisterListDomains(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_domains",
		Description: `Discover all scaffolded domains in the project.

Scans internal/models, internal/repository, internal/services, and internal/web directories.

Returns for each domain:
- Which layers exist (model, repository, service, controller)
- List of views in internal/web/{domain}/views

Use this to understand project structure before adding new domains.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ListDomainsInput) (*mcp.CallToolResult, types.ListDomainsResult, error) {
		result, err := listDomains(registry)
		if err != nil {
			return nil, types.NewListDomainsError(err.Error()), nil
		}
		return nil, result, nil
	})
}

func listDomains(registry *Registry) (types.ListDomainsResult, error) {
	// Scan for domains in various directories
	domains := make(map[string]*types.DomainInfo)

	// Scan models directory
	modelsDir := filepath.Join(registry.WorkingDir, "internal", "models")
	modelFiles, _ := utils.ListFiles(modelsDir, "*.go")
	for _, file := range modelFiles {
		name := filepath.Base(file)
		name = name[:len(name)-3] // remove .go extension

		// Skip base.go and other non-domain files
		if name == "base" || name == "models" {
			continue
		}

		if _, ok := domains[name]; !ok {
			domains[name] = &types.DomainInfo{Name: name}
		}
		domains[name].HasModel = true
	}

	// Scan repository directories
	repoDir := filepath.Join(registry.WorkingDir, "internal", "repository")
	repoDirs, _ := utils.ListDirs(repoDir)
	for _, dir := range repoDirs {
		if _, ok := domains[dir]; !ok {
			domains[dir] = &types.DomainInfo{Name: dir}
		}
		domains[dir].HasRepository = true
	}

	// Scan services directories
	servicesDir := filepath.Join(registry.WorkingDir, "internal", "services")
	serviceDirs, _ := utils.ListDirs(servicesDir)
	for _, dir := range serviceDirs {
		if _, ok := domains[dir]; !ok {
			domains[dir] = &types.DomainInfo{Name: dir}
		}
		domains[dir].HasService = true
	}

	// Scan web directories (controllers)
	webDir := filepath.Join(registry.WorkingDir, "internal", "web")
	webDirs, _ := utils.ListDirs(webDir)
	for _, dir := range webDirs {
		// Skip common directories
		if dir == "middleware" || dir == "layouts" || dir == "components" {
			continue
		}

		if _, ok := domains[dir]; !ok {
			domains[dir] = &types.DomainInfo{Name: dir}
		}
		domains[dir].HasController = true

		// Check for views
		viewsDir := filepath.Join(webDir, dir, "views")
		viewFiles, _ := utils.ListFiles(viewsDir, "*.templ")
		for _, viewFile := range viewFiles {
			viewName := filepath.Base(viewFile)
			viewName = viewName[:len(viewName)-6] // remove .templ extension
			domains[dir].Views = append(domains[dir].Views, viewName)
		}
	}

	// Convert map to slice
	result := make([]types.DomainInfo, 0, len(domains))
	for _, info := range domains {
		result = append(result, *info)
	}

	return types.NewListDomainsResult(result), nil
}
