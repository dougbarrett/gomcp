// Package templates provides embedded template files for code generation.
package templates

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

// FS is the embedded filesystem containing all template files.
// Templates use [[ ]] delimiters instead of {{ }} to avoid conflicts with Go templates.
//
//go:embed project/*.tmpl domain/*.tmpl views/*.tmpl components/*.tmpl config/*.tmpl seed/*.tmpl auth/*.tmpl usermgmt/*.tmpl usermgmt/views/*.tmpl wizard/*.tmpl
var FS embed.FS

// Template directories:
// - project/    : Project scaffolding templates (go.mod, main.go, config, etc.)
// - domain/     : Domain layer templates (model, repository, service, controller, dto)
// - views/      : View templates (list, show, form, table, partials)
// - components/ : Component templates (card, modal, form_field, wizard)
// - config/     : Configuration templates (page.toml)
// - seed/       : Seeder templates (seeder.go)
// - auth/       : Authentication templates (user_model, middleware, service, controller, views)
// - usermgmt/   : User management templates (service, controller, views)
// - wizard/     : Wizard templates (controller, views, draft model/repo/service)

// Categories of templates available.
var Categories = []string{
	"project",
	"domain",
	"views",
	"components",
	"config",
	"seed",
	"auth",
	"usermgmt",
	"wizard",
}

// ReadTemplate reads a template file by path and returns its contents.
// The path should be relative to the templates directory (e.g., "project/main.go.tmpl").
func ReadTemplate(path string) ([]byte, error) {
	return FS.ReadFile(path)
}

// TemplateExists checks if a template exists at the given path.
func TemplateExists(path string) bool {
	_, err := FS.ReadFile(path)
	return err == nil
}

// ListTemplates returns all template paths in the embedded filesystem.
func ListTemplates() ([]string, error) {
	var templates []string
	err := fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			templates = append(templates, path)
		}
		return nil
	})
	return templates, err
}

// ListTemplatesInCategory returns all template paths in a specific category.
func ListTemplatesInCategory(category string) ([]string, error) {
	var templates []string
	entries, err := FS.ReadDir(category)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			templates = append(templates, filepath.Join(category, entry.Name()))
		}
	}
	return templates, nil
}

// CountTemplates returns the total number of templates.
func CountTemplates() (int, error) {
	templates, err := ListTemplates()
	if err != nil {
		return 0, err
	}
	return len(templates), nil
}
