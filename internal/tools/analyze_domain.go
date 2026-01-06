package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/metadata"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// RegisterAnalyzeDomain registers the analyze_domain tool.
func RegisterAnalyzeDomain(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "analyze_domain",
		Description: `Analyze scaffolded domains to detect differences from current templates.

This tool compares existing scaffolded code against what the current templates would generate,
helping identify opportunities to upgrade code to use new template features.

Requirements:
- Domain must have been scaffolded with metadata tracking (scaffold_domain saves this automatically)
- For domains without metadata, use list_domains to see what exists

Use cases:
- Check if a domain can benefit from new template features (e.g., form_style, belongs_to display)
- See what changes would be made before running sync_domain
- Audit scaffolded code for drift from templates

Output includes:
- List of files with differences
- Unified diff showing what would change
- Summary of added/removed lines

Examples:
1. Analyze a specific domain:
   analyze_domain: { domain: "order" }

2. Analyze all domains:
   analyze_domain: {}

3. Analyze only views:
   analyze_domain: { domain: "order", layers: ["views"] }

4. Include unchanged files:
   analyze_domain: { domain: "order", show_unchanged: true }`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.AnalyzeDomainInput) (*mcp.CallToolResult, types.AnalyzeDomainResult, error) {
		result, err := ExecuteAnalyzeDomain(ctx, registry, input)
		if err != nil {
			return nil, types.AnalyzeDomainResult{Success: false, Message: err.Error()}, nil
		}
		return nil, result, nil
	})
}

// ExecuteAnalyzeDomain executes the analyze_domain tool.
func ExecuteAnalyzeDomain(ctx context.Context, registry *Registry, input types.AnalyzeDomainInput) (types.AnalyzeDomainResult, error) {
	metaStore := metadata.NewStore(registry.WorkingDir)

	var domainsToAnalyze []string

	if input.Domain != "" {
		// Analyze specific domain
		exists, err := metaStore.Exists(input.Domain)
		if err != nil {
			return types.AnalyzeDomainResult{
				Success: false,
				Message: fmt.Sprintf("Failed to check domain metadata: %v", err),
			}, nil
		}
		if !exists {
			return types.AnalyzeDomainResult{
				Success: false,
				Message: fmt.Sprintf("No metadata found for domain '%s'. Domain may have been scaffolded before metadata tracking was added, or may not exist.", input.Domain),
			}, nil
		}
		domainsToAnalyze = []string{input.Domain}
	} else {
		// Analyze all domains with metadata
		domains, err := metaStore.ListDomains()
		if err != nil {
			return types.AnalyzeDomainResult{
				Success: false,
				Message: fmt.Sprintf("Failed to list domains: %v", err),
			}, nil
		}
		if len(domains) == 0 {
			return types.AnalyzeDomainResult{
				Success: false,
				Message: "No domains with metadata found. Scaffold a domain first, or check .mcp/scaffold-metadata.json exists.",
			}, nil
		}
		domainsToAnalyze = domains
	}

	var analyses []types.DomainAnalysis
	totalChanges := 0

	for _, domainName := range domainsToAnalyze {
		analysis, err := analyzeSingleDomain(registry, metaStore, domainName, input.Layers, input.ShowUnchanged)
		if err != nil {
			// Include error in analysis but continue
			analysis = types.DomainAnalysis{
				Domain:     domainName,
				HasChanges: false,
				Files: []types.FileAnalysis{{
					Path:   "",
					Status: "error",
					Diff:   err.Error(),
				}},
			}
		}
		if analysis.HasChanges {
			totalChanges++
		}
		analyses = append(analyses, analysis)
	}

	message := fmt.Sprintf("Analyzed %d domain(s)", len(analyses))
	if totalChanges > 0 {
		message += fmt.Sprintf(", %d with changes available", totalChanges)
	} else {
		message += ", all up to date"
	}

	result := types.AnalyzeDomainResult{
		Success: true,
		Message: message,
		Domains: analyses,
	}

	// Suggest sync_domain if there are changes
	if totalChanges > 0 {
		result.SuggestedTool = &types.ToolHint{
			Tool:        "sync_domain",
			Description: "Apply the detected changes to update scaffolded code",
			Priority:    "optional",
		}
	}

	return result, nil
}

// analyzeSingleDomain analyzes a single domain and returns the analysis.
func analyzeSingleDomain(registry *Registry, metaStore *metadata.Store, domainName string, layers []string, showUnchanged bool) (types.DomainAnalysis, error) {
	domainMeta, exists, err := metaStore.GetDomain(domainName)
	if err != nil {
		return types.DomainAnalysis{}, err
	}
	if !exists {
		return types.DomainAnalysis{}, fmt.Errorf("domain metadata not found")
	}

	analysis := types.DomainAnalysis{
		Domain:            domainName,
		ScaffoldedAt:      domainMeta.ScaffoldedAt.Format("2006-01-02 15:04:05"),
		ScaffolderVersion: domainMeta.ScaffolderVersion,
		CurrentVersion:    ScaffolderVersion,
		HasChanges:        false,
	}

	// Get module path
	modulePath, err := utils.GetModulePath(registry.WorkingDir)
	if err != nil {
		return types.DomainAnalysis{}, fmt.Errorf("failed to get module path: %w", err)
	}

	// Create generator in dry run mode with content storage
	gen := registry.NewGenerator("")
	gen.SetDryRun(true)
	gen.SetStoreContent(true)
	gen.SetForceOverwrite(true) // Allow "overwriting" to capture all files

	// Prepare template data using the stored input
	domainInput := domainMeta.Input
	data := generator.NewDomainData(domainInput, modulePath)

	// Generate all domain files (same logic as scaffold_domain)
	pkgName := utils.ToPackageName(domainInput.DomainName)

	// Generate model
	modelPath := filepath.Join("internal", "models", pkgName+".go")
	if err := gen.GenerateFile("domain/model.go.tmpl", modelPath, data); err != nil {
		return types.DomainAnalysis{}, fmt.Errorf("failed to generate model: %w", err)
	}

	// Generate repository
	repoPath := filepath.Join("internal", "repository", pkgName, pkgName+".go")
	if err := gen.GenerateFile("domain/repository.go.tmpl", repoPath, data); err != nil {
		return types.DomainAnalysis{}, fmt.Errorf("failed to generate repository: %w", err)
	}

	// Generate service
	servicePath := filepath.Join("internal", "services", pkgName, pkgName+".go")
	if err := gen.GenerateFile("domain/service.go.tmpl", servicePath, data); err != nil {
		return types.DomainAnalysis{}, fmt.Errorf("failed to generate service: %w", err)
	}

	// Generate DTOs
	dtoPath := filepath.Join("internal", "services", pkgName, "dto.go")
	if err := gen.GenerateFile("domain/dto.go.tmpl", dtoPath, data); err != nil {
		return types.DomainAnalysis{}, fmt.Errorf("failed to generate DTOs: %w", err)
	}

	// Generate controller
	controllerPath := filepath.Join("internal", "web", pkgName, pkgName+".go")
	if err := gen.GenerateFile("domain/controller.go.tmpl", controllerPath, data); err != nil {
		return types.DomainAnalysis{}, fmt.Errorf("failed to generate controller: %w", err)
	}

	// Generate CRUD views if requested
	if domainInput.GetWithCrudViews() {
		viewsDir := filepath.Join("internal", "web", pkgName, "views")

		// Generate list view
		listPath := filepath.Join(viewsDir, "list.templ")
		if err := gen.GenerateFile("views/list.templ.tmpl", listPath, data); err != nil {
			return types.DomainAnalysis{}, fmt.Errorf("failed to generate list view: %w", err)
		}

		// Generate show view
		showPath := filepath.Join(viewsDir, "show.templ")
		if err := gen.GenerateFile("views/show.templ.tmpl", showPath, data); err != nil {
			return types.DomainAnalysis{}, fmt.Errorf("failed to generate show view: %w", err)
		}

		// Generate form view
		formPath := filepath.Join(viewsDir, pkgName+"_form.templ")
		if err := gen.GenerateFile("views/form.templ.tmpl", formPath, data); err != nil {
			return types.DomainAnalysis{}, fmt.Errorf("failed to generate form view: %w", err)
		}
	}

	// Build layer filter
	layerFilter := make(map[string]bool)
	if len(layers) > 0 {
		for _, l := range layers {
			layerFilter[strings.ToLower(l)] = true
		}
	}

	// Get result and compare files
	result := gen.Result()
	dmp := diffmatchpatch.New()

	// Combine created and updated files for comparison
	allFiles := append(result.FilesCreated, result.FilesUpdated...)

	for _, filePath := range allFiles {
		// Check layer filter
		if len(layerFilter) > 0 && !matchesLayer(filePath, layerFilter) {
			continue
		}

		fullPath := filepath.Join(registry.WorkingDir, filePath)
		existingContent, err := os.ReadFile(fullPath)

		var fileAnalysis types.FileAnalysis
		fileAnalysis.Path = filePath

		if err != nil {
			if os.IsNotExist(err) {
				fileAnalysis.Status = "missing"
				fileAnalysis.Diff = "File does not exist (may have been deleted)"
				analysis.HasChanges = true
			} else {
				fileAnalysis.Status = "error"
				fileAnalysis.Diff = err.Error()
			}
		} else {
			// Get generated content
			generatedContent := gen.GetFileContent(filePath)
			if generatedContent == "" {
				continue // Skip if we can't get generated content
			}

			existingStr := string(existingContent)

			if existingStr == generatedContent {
				if showUnchanged {
					fileAnalysis.Status = "unchanged"
				} else {
					continue // Skip unchanged files
				}
			} else {
				fileAnalysis.Status = "modified"
				analysis.HasChanges = true

				// Generate unified diff
				diffs := dmp.DiffMain(existingStr, generatedContent, true)
				fileAnalysis.Diff = generateUnifiedDiff(filePath, diffs)

				// Count changes
				for _, d := range diffs {
					lines := strings.Count(d.Text, "\n")
					if lines == 0 && len(d.Text) > 0 {
						lines = 1
					}
					switch d.Type {
					case diffmatchpatch.DiffInsert:
						fileAnalysis.LinesAdded += lines
					case diffmatchpatch.DiffDelete:
						fileAnalysis.LinesRemoved += lines
					}
				}
			}
		}

		analysis.Files = append(analysis.Files, fileAnalysis)
	}

	return analysis, nil
}

// matchesLayer checks if a file path matches the layer filter.
func matchesLayer(filePath string, layerFilter map[string]bool) bool {
	if strings.Contains(filePath, "/models/") && layerFilter["model"] {
		return true
	}
	if strings.Contains(filePath, "/repository/") && layerFilter["repository"] {
		return true
	}
	if strings.Contains(filePath, "/services/") && layerFilter["service"] {
		return true
	}
	if strings.Contains(filePath, "/web/") && !strings.Contains(filePath, "/views/") && layerFilter["controller"] {
		return true
	}
	if strings.Contains(filePath, "/views/") && layerFilter["views"] {
		return true
	}
	return false
}

// generateUnifiedDiff creates a human-readable unified diff.
func generateUnifiedDiff(filePath string, diffs []diffmatchpatch.Diff) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- a/%s (existing)\n", filePath))
	sb.WriteString(fmt.Sprintf("+++ b/%s (new template)\n", filePath))

	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		for i, line := range lines {
			// Skip empty last line from split
			if i == len(lines)-1 && line == "" {
				continue
			}

			switch diff.Type {
			case diffmatchpatch.DiffEqual:
				sb.WriteString(fmt.Sprintf(" %s\n", line))
			case diffmatchpatch.DiffInsert:
				sb.WriteString(fmt.Sprintf("+%s\n", line))
			case diffmatchpatch.DiffDelete:
				sb.WriteString(fmt.Sprintf("-%s\n", line))
			}
		}
	}

	return sb.String()
}
