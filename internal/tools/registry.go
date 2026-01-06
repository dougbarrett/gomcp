// Package tools contains MCP tool implementations for scaffolding.
package tools

import (
	"os"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/templates"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registry holds references needed by tools.
type Registry struct {
	// WorkingDir is the base directory for scaffolding operations.
	WorkingDir string
}

// NewRegistry creates a new tool registry.
func NewRegistry(workingDir string) *Registry {
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	return &Registry{
		WorkingDir: workingDir,
	}
}

// NewGenerator creates a generator for the given project path.
func (r *Registry) NewGenerator(projectPath string) *generator.Generator {
	basePath := projectPath
	if projectPath == "" {
		basePath = r.WorkingDir
	}
	return generator.NewGenerator(templates.FS, basePath)
}

// CheckForConflicts checks if the generator has conflicts and returns a conflict result if so.
// Returns nil if there are no conflicts.
func CheckForConflicts(result generator.GeneratorResult) *types.ScaffoldResult {
	if !result.HasConflicts {
		return nil
	}

	conflicts := make([]types.FileConflict, len(result.Conflicts))
	for i, c := range result.Conflicts {
		conflicts[i] = types.FileConflict{
			Path:            c.Path,
			Description:     c.Description,
			ProposedContent: c.ProposedContent,
		}
	}

	conflictResult := types.NewConflictResult(conflicts)
	return &conflictResult
}

// RegisterAll registers all scaffolding tools with the server.
func (r *Registry) RegisterAll(server *mcp.Server) {
	// Phase 2: Project scaffolding
	RegisterScaffoldProject(server, r)

	// Phase 3: Domain layer tools
	RegisterScaffoldDomain(server, r)
	RegisterScaffoldRepository(server, r)
	RegisterScaffoldService(server, r)
	RegisterScaffoldServiceForRepo(server, r)
	RegisterScaffoldController(server, r)

	// Phase 4: View layer tools
	RegisterScaffoldView(server, r)
	RegisterScaffoldComponent(server, r)
	RegisterScaffoldForm(server, r)
	RegisterScaffoldTable(server, r)
	RegisterScaffoldModal(server, r)

	// Phase 5: Configuration and utility tools
	RegisterScaffoldPage(server, r)
	RegisterScaffoldConfig(server, r)
	RegisterScaffoldSeed(server, r)
	RegisterListDomains(server, r)
	RegisterAnalyzeDomain(server, r)
	RegisterUpdateDIWiring(server, r)

	// Phase 6: Extend tools for custom logic
	RegisterExtendRepository(server, r)
	RegisterExtendService(server, r)
	RegisterExtendController(server, r)

	// Utility tools
	RegisterReportBug(server, r)
}
