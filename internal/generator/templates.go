package generator

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

const (
	// LeftDelim is the left template delimiter.
	LeftDelim = "[["
	// RightDelim is the right template delimiter.
	RightDelim = "]]"
)

// LoadTemplate loads a template from the embedded filesystem.
// The template uses [[ ]] delimiters instead of {{ }}.
func LoadTemplate(fs embed.FS, name string) (*template.Template, error) {
	content, err := fs.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", name, err)
	}

	return ParseTemplate(name, string(content))
}

// ParseTemplate parses a template string with custom delimiters.
func ParseTemplate(name, content string) (*template.Template, error) {
	tmpl := template.New(name).
		Delims(LeftDelim, RightDelim).
		Funcs(TemplateFuncMap())

	parsed, err := tmpl.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return parsed, nil
}

// ExecuteTemplate loads and executes a template from the embedded filesystem.
func ExecuteTemplate(fs embed.FS, name string, data any) (string, error) {
	tmpl, err := LoadTemplate(fs, name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// ExecuteTemplateString parses and executes a template string.
func ExecuteTemplateString(name, content string, data any) (string, error) {
	tmpl, err := ParseTemplate(name, content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// MustLoadTemplate loads a template or panics.
func MustLoadTemplate(fs embed.FS, name string) *template.Template {
	tmpl, err := LoadTemplate(fs, name)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// MustParseTemplate parses a template or panics.
func MustParseTemplate(name, content string) *template.Template {
	tmpl, err := ParseTemplate(name, content)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// TemplateExists checks if a template exists in the embedded filesystem.
func TemplateExists(fs embed.FS, name string) bool {
	_, err := fs.ReadFile(name)
	return err == nil
}

// ListTemplates lists all templates matching a pattern in the embedded filesystem.
func ListTemplates(fs embed.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read template directory %s: %w", dir, err)
	}

	var templates []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			templates = append(templates, dir+"/"+entry.Name())
		}
	}

	return templates, nil
}

// TemplateInfo contains information about a template.
type TemplateInfo struct {
	// Name is the template name/path.
	Name string
	// Content is the raw template content.
	Content string
	// Size is the content size in bytes.
	Size int
}

// GetTemplateInfo returns information about a template.
func GetTemplateInfo(fs embed.FS, name string) (*TemplateInfo, error) {
	content, err := fs.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", name, err)
	}

	return &TemplateInfo{
		Name:    name,
		Content: string(content),
		Size:    len(content),
	}, nil
}
