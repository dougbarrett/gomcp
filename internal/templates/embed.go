// Package templates provides embedded template files for code generation.
package templates

import "embed"

// FS is the embedded filesystem containing all template files.
// Templates use [[ ]] delimiters instead of {{ }} to avoid conflicts with Go templates.
//
//go:embed project/*.tmpl domain/*.tmpl views/*.tmpl components/*.tmpl config/*.tmpl seed/*.tmpl auth/*.tmpl
var FS embed.FS

// Template directories:
// - project/    : Project scaffolding templates (go.mod, main.go, config, etc.)
// - domain/     : Domain layer templates (model, repository, service, controller, dto)
// - views/      : View templates (list, show, form, table, partials)
// - components/ : Component templates (card, modal, form_field)
// - config/     : Configuration templates (page.toml)
// - seed/       : Seeder templates (seeder.go)
// - auth/       : Authentication templates (user_model, middleware, service, controller, views)
