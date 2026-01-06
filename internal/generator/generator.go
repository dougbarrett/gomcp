// Package generator provides template-based code generation capabilities.
package generator

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/utils"
)

// FileConflict represents a file that would be overwritten.
type FileConflict struct {
	// Path is the relative file path.
	Path string
	// Description explains the purpose of this file.
	Description string
	// ProposedContent is the content that would be written.
	ProposedContent string
}

// Generator handles template-based file generation.
type Generator struct {
	// fs is the embedded filesystem containing templates.
	fs embed.FS
	// basePath is the base directory for generated files.
	basePath string
	// dryRun if true, no files are written.
	dryRun bool
	// forceOverwrite if true, allows overwriting existing files.
	forceOverwrite bool
	// storeContent if true, stores generated content for later retrieval (used for analysis).
	storeContent bool
	// filesCreated tracks created files.
	filesCreated []string
	// filesUpdated tracks updated files.
	filesUpdated []string
	// conflicts tracks files that would be overwritten.
	conflicts []FileConflict
	// generatedContent stores generated file content when storeContent is true.
	generatedContent map[string]string
}

// GeneratorResult contains the results of generation.
type GeneratorResult struct {
	// FilesCreated is the list of created files.
	FilesCreated []string
	// FilesUpdated is the list of updated files.
	FilesUpdated []string
	// Conflicts is the list of files that would be overwritten.
	Conflicts []FileConflict
	// HasConflicts is true if there are any conflicts.
	HasConflicts bool
}

// NewGenerator creates a new Generator.
func NewGenerator(fs embed.FS, basePath string) *Generator {
	return &Generator{
		fs:           fs,
		basePath:     basePath,
		filesCreated: make([]string, 0),
		filesUpdated: make([]string, 0),
		conflicts:    make([]FileConflict, 0),
	}
}

// SetForceOverwrite sets whether to allow overwriting existing files.
func (g *Generator) SetForceOverwrite(force bool) {
	g.forceOverwrite = force
}

// SetDryRun sets the dry run mode.
func (g *Generator) SetDryRun(dryRun bool) {
	g.dryRun = dryRun
}

// IsDryRun returns whether dry run mode is enabled.
func (g *Generator) IsDryRun() bool {
	return g.dryRun
}

// SetStoreContent sets whether to store generated content for later retrieval.
// This is useful for analysis/comparison without writing files.
func (g *Generator) SetStoreContent(store bool) {
	g.storeContent = store
	if store && g.generatedContent == nil {
		g.generatedContent = make(map[string]string)
	}
}

// GetFileContent retrieves the generated content for a file path.
// Only works if SetStoreContent(true) was called before generation.
func (g *Generator) GetFileContent(outputPath string) string {
	if g.generatedContent == nil {
		return ""
	}
	return g.generatedContent[outputPath]
}

// BasePath returns the base path for generation.
func (g *Generator) BasePath() string {
	return g.basePath
}

// EnsureDir creates a directory relative to the base path.
func (g *Generator) EnsureDir(relPath string) error {
	fullPath := filepath.Join(g.basePath, relPath)

	if g.dryRun {
		return nil
	}

	return utils.EnsureDir(fullPath)
}

// GenerateFile generates a file from a template.
// templatePath is the path within the embedded FS.
// outputPath is the path relative to basePath.
// data is the template data.
func (g *Generator) GenerateFile(templatePath, outputPath string, data any) error {
	return g.GenerateFileWithDescription(templatePath, outputPath, data, "")
}

// GenerateFileWithDescription generates a file from a template with a description for conflict reporting.
func (g *Generator) GenerateFileWithDescription(templatePath, outputPath string, data any, description string) error {
	fullOutputPath := filepath.Join(g.basePath, outputPath)

	// Check if file exists
	fileExists := utils.FileExists(fullOutputPath)

	// Load and execute template
	content, err := ExecuteTemplate(g.fs, templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	// Store content for later retrieval if enabled
	if g.storeContent {
		if g.generatedContent == nil {
			g.generatedContent = make(map[string]string)
		}
		g.generatedContent[outputPath] = content
	}

	// If file exists and we're not forcing overwrite, record as conflict
	if fileExists && !g.forceOverwrite {
		desc := description
		if desc == "" {
			desc = inferFileDescription(outputPath)
		}
		g.conflicts = append(g.conflicts, FileConflict{
			Path:            outputPath,
			Description:     desc,
			ProposedContent: content,
		})
		return nil
	}

	if g.dryRun {
		if fileExists {
			g.filesUpdated = append(g.filesUpdated, outputPath)
		} else {
			g.filesCreated = append(g.filesCreated, outputPath)
		}
		return nil
	}

	// Ensure parent directory exists
	dir := filepath.Dir(fullOutputPath)
	if err := utils.EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := utils.WriteFileString(fullOutputPath, content, true); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullOutputPath, err)
	}

	if fileExists {
		g.filesUpdated = append(g.filesUpdated, outputPath)
	} else {
		g.filesCreated = append(g.filesCreated, outputPath)
	}

	return nil
}

// inferFileDescription infers a description based on the file path.
func inferFileDescription(path string) string {
	if strings.Contains(path, "/models/") {
		return "Model definition with struct fields and GORM tags"
	}
	if strings.Contains(path, "/repository/") {
		return "Repository layer with database CRUD operations"
	}
	if strings.Contains(path, "/services/") {
		if strings.HasSuffix(path, "dto.go") {
			return "Data Transfer Objects for service layer"
		}
		return "Service layer with business logic"
	}
	if strings.Contains(path, "/web/") && strings.Contains(path, "/views/") {
		return "Templ view template for UI rendering"
	}
	if strings.Contains(path, "/web/") && !strings.Contains(path, "/views/") {
		return "HTTP controller with route handlers"
	}
	if strings.Contains(path, "/middleware/") {
		return "HTTP middleware"
	}
	if strings.Contains(path, "/config/") {
		return "Configuration file"
	}
	if strings.HasSuffix(path, "_seeder.go") {
		return "Database seeder for test data"
	}
	return "Generated source file"
}

// GenerateFileIfNotExists generates a file only if it doesn't exist.
func (g *Generator) GenerateFileIfNotExists(templatePath, outputPath string, data any) error {
	fullOutputPath := filepath.Join(g.basePath, outputPath)

	if utils.FileExists(fullOutputPath) {
		return nil
	}

	return g.GenerateFile(templatePath, outputPath, data)
}

// GenerateFileFromString generates a file from a string content.
func (g *Generator) GenerateFileFromString(outputPath, content string) error {
	return g.GenerateFileFromStringWithDescription(outputPath, content, "")
}

// GenerateFileFromStringWithDescription generates a file from string content with a description.
func (g *Generator) GenerateFileFromStringWithDescription(outputPath, content, description string) error {
	fullOutputPath := filepath.Join(g.basePath, outputPath)

	// Check if file exists
	fileExists := utils.FileExists(fullOutputPath)

	// If file exists and we're not forcing overwrite, record as conflict
	if fileExists && !g.forceOverwrite {
		desc := description
		if desc == "" {
			desc = inferFileDescription(outputPath)
		}
		g.conflicts = append(g.conflicts, FileConflict{
			Path:            outputPath,
			Description:     desc,
			ProposedContent: content,
		})
		return nil
	}

	if g.dryRun {
		if fileExists {
			g.filesUpdated = append(g.filesUpdated, outputPath)
		} else {
			g.filesCreated = append(g.filesCreated, outputPath)
		}
		return nil
	}

	// Ensure parent directory exists
	dir := filepath.Dir(fullOutputPath)
	if err := utils.EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := utils.WriteFileString(fullOutputPath, content, true); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullOutputPath, err)
	}

	if fileExists {
		g.filesUpdated = append(g.filesUpdated, outputPath)
	} else {
		g.filesCreated = append(g.filesCreated, outputPath)
	}

	return nil
}

// Result returns the generation result.
func (g *Generator) Result() GeneratorResult {
	return GeneratorResult{
		FilesCreated: g.filesCreated,
		FilesUpdated: g.filesUpdated,
		Conflicts:    g.conflicts,
		HasConflicts: len(g.conflicts) > 0,
	}
}

// HasConflicts returns true if there are any file conflicts.
func (g *Generator) HasConflicts() bool {
	return len(g.conflicts) > 0
}

// Conflicts returns the list of file conflicts.
func (g *Generator) Conflicts() []FileConflict {
	return g.conflicts
}

// Reset clears the tracked files.
func (g *Generator) Reset() {
	g.filesCreated = make([]string, 0)
	g.filesUpdated = make([]string, 0)
	g.conflicts = make([]FileConflict, 0)
}

// FullPath returns the full path for a relative path.
func (g *Generator) FullPath(relPath string) string {
	return filepath.Join(g.basePath, relPath)
}

// FileExists checks if a file exists relative to base path.
func (g *Generator) FileExists(relPath string) bool {
	return utils.FileExists(g.FullPath(relPath))
}

// DirExists checks if a directory exists relative to base path.
func (g *Generator) DirExists(relPath string) bool {
	return utils.DirExists(g.FullPath(relPath))
}

// ReadFile reads a file relative to base path.
func (g *Generator) ReadFile(relPath string) (string, error) {
	return utils.ReadFileString(g.FullPath(relPath))
}

// WriteFile writes a file relative to base path.
func (g *Generator) WriteFile(relPath, content string) error {
	fullPath := g.FullPath(relPath)
	fileExists := utils.FileExists(fullPath)

	// If file exists and we're not forcing overwrite, record as conflict
	if fileExists && !g.forceOverwrite {
		g.conflicts = append(g.conflicts, FileConflict{
			Path:            relPath,
			Description:     inferFileDescription(relPath),
			ProposedContent: content,
		})
		return nil
	}

	if g.dryRun {
		if fileExists {
			g.filesUpdated = append(g.filesUpdated, relPath)
		} else {
			g.filesCreated = append(g.filesCreated, relPath)
		}
		return nil
	}

	if err := utils.WriteFileString(fullPath, content, true); err != nil {
		return err
	}

	if fileExists {
		g.filesUpdated = append(g.filesUpdated, relPath)
	} else {
		g.filesCreated = append(g.filesCreated, relPath)
	}

	return nil
}

// ListGeneratedFiles returns all files that were created or updated.
func (g *Generator) ListGeneratedFiles() []string {
	all := make([]string, 0, len(g.filesCreated)+len(g.filesUpdated))
	all = append(all, g.filesCreated...)
	all = append(all, g.filesUpdated...)
	return all
}

// Summary returns a summary of the generation.
func (g *Generator) Summary() string {
	var sb strings.Builder
	if len(g.filesCreated) > 0 {
		sb.WriteString(fmt.Sprintf("Created %d file(s):\n", len(g.filesCreated)))
		for _, f := range g.filesCreated {
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}
	if len(g.filesUpdated) > 0 {
		sb.WriteString(fmt.Sprintf("Updated %d file(s):\n", len(g.filesUpdated)))
		for _, f := range g.filesUpdated {
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}
	if len(g.filesCreated) == 0 && len(g.filesUpdated) == 0 {
		sb.WriteString("No files generated.\n")
	}
	return sb.String()
}
