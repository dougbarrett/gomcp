// Package generator provides template-based code generation capabilities.
package generator

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/utils"
)

// Generator handles template-based file generation.
type Generator struct {
	// fs is the embedded filesystem containing templates.
	fs embed.FS
	// basePath is the base directory for generated files.
	basePath string
	// dryRun if true, no files are written.
	dryRun bool
	// filesCreated tracks created files.
	filesCreated []string
	// filesUpdated tracks updated files.
	filesUpdated []string
}

// GeneratorResult contains the results of generation.
type GeneratorResult struct {
	// FilesCreated is the list of created files.
	FilesCreated []string
	// FilesUpdated is the list of updated files.
	FilesUpdated []string
}

// NewGenerator creates a new Generator.
func NewGenerator(fs embed.FS, basePath string) *Generator {
	return &Generator{
		fs:           fs,
		basePath:     basePath,
		filesCreated: make([]string, 0),
		filesUpdated: make([]string, 0),
	}
}

// SetDryRun sets the dry run mode.
func (g *Generator) SetDryRun(dryRun bool) {
	g.dryRun = dryRun
}

// IsDryRun returns whether dry run mode is enabled.
func (g *Generator) IsDryRun() bool {
	return g.dryRun
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
	fullOutputPath := filepath.Join(g.basePath, outputPath)

	// Check if file exists
	fileExists := utils.FileExists(fullOutputPath)

	// Load and execute template
	content, err := ExecuteTemplate(g.fs, templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
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
	fullOutputPath := filepath.Join(g.basePath, outputPath)

	// Check if file exists
	fileExists := utils.FileExists(fullOutputPath)

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
	}
}

// Reset clears the tracked files.
func (g *Generator) Reset() {
	g.filesCreated = make([]string, 0)
	g.filesUpdated = make([]string, 0)
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
