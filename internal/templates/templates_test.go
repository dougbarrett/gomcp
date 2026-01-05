// Package templates provides embedded template files for code generation.
package templates

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/dbb1dev/go-mcp/internal/generator"
)

const (
	leftDelim  = "[["
	rightDelim = "]]"
)

// parseTemplate parses a template with the correct delimiters and functions.
func parseTemplate(name, content string) (*template.Template, error) {
	return template.New(name).
		Delims(leftDelim, rightDelim).
		Funcs(generator.TemplateFuncMap()).
		Parse(content)
}

// TestAllTemplatesParseCorrectly tests that all embedded templates can be parsed.
func TestAllTemplatesParseCorrectly(t *testing.T) {
	dirs := []string{
		"project",
		"domain",
		"views",
		"components",
		"config",
		"seed",
		"auth",
	}

	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			entries, err := FS.ReadDir(dir)
			if err != nil {
				t.Fatalf("Failed to read directory %s: %v", dir, err)
			}

			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmpl") {
					continue
				}

				templatePath := dir + "/" + entry.Name()
				t.Run(entry.Name(), func(t *testing.T) {
					content, err := FS.ReadFile(templatePath)
					if err != nil {
						t.Fatalf("Failed to read template %s: %v", templatePath, err)
					}

					_, err = parseTemplate(entry.Name(), string(content))
					if err != nil {
						t.Errorf("Failed to parse template %s: %v", templatePath, err)
					}
				})
			}
		})
	}
}

// TestProjectTemplatesExecute tests that project templates execute with valid data.
func TestProjectTemplatesExecute(t *testing.T) {
	projectData := struct {
		ProjectName  string
		ModulePath   string
		DatabaseType string
		WithAuth     bool
	}{
		ProjectName:  "testproject",
		ModulePath:   "github.com/test/testproject",
		DatabaseType: "sqlite",
		WithAuth:     true,
	}

	templates := []string{
		"project/go.mod.tmpl",
		"project/main.go.tmpl",
		"project/database.go.tmpl",
		"project/base_model.go.tmpl",
		"project/response.go.tmpl",
		"project/base_layout.templ.tmpl",
		"project/air.toml.tmpl",
		"project/tailwind_input.css.tmpl",
		"project/app.toml.tmpl",
		"project/gitignore.tmpl",
		"project/taskfile.yml.tmpl",
		"project/config.go.tmpl",
		"project/middleware.go.tmpl",
		"project/router.go.tmpl",
		"project/common_components.templ.tmpl",
		"project/menu.toml.tmpl",
		"project/seed_main.go.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, projectData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}
		})
	}
}

// TestDomainTemplatesExecute tests that domain templates execute with valid data.
func TestDomainTemplatesExecute(t *testing.T) {
	domainData := struct {
		ModulePath           string
		DomainName           string
		ModelName            string
		PackageName          string
		VariableName         string
		TableName            string
		URLPath              string
		Fields               []generator.FieldData
		Relationships        []generator.RelationshipData
		HasRelationships     bool
		PreloadRelationships []generator.RelationshipData
		WithSoftDelete       bool
		WithCrudViews        bool
	}{
		ModulePath:   "github.com/test/testproject",
		DomainName:   "product",
		ModelName:    "Product",
		PackageName:  "product",
		VariableName: "product",
		TableName:    "products",
		URLPath:      "/products",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", GORMTags: "size:255", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
			{Name: "Price", Type: "float64", JSONName: "price", Required: true, Label: "Price", FormType: "number"},
			{Name: "Description", Type: "string", JSONName: "description", Label: "Description", FormType: "textarea"},
			{Name: "Active", Type: "bool", JSONName: "active", Label: "Active", FormType: "checkbox"},
		},
		Relationships:        []generator.RelationshipData{},
		HasRelationships:     false,
		PreloadRelationships: []generator.RelationshipData{},
		WithSoftDelete:       true,
		WithCrudViews:        true,
	}

	templates := []string{
		"domain/model.go.tmpl",
		"domain/repository.go.tmpl",
		"domain/service.go.tmpl",
		"domain/controller.go.tmpl",
		"domain/dto.go.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, domainData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}
		})
	}
}

// TestViewTemplatesExecute tests that view templates execute with valid data.
func TestViewTemplatesExecute(t *testing.T) {
	viewData := struct {
		ModulePath        string
		DomainName        string
		ModelName         string
		PackageName       string
		VariableName      string
		URLPath           string
		ViewType          string
		ViewName          string
		Fields            []generator.FieldData
		Columns           []generator.ColumnData
		WithPagination    bool
		WithSearch        bool
		WithFilters       bool
		WithSorting       bool
		WithBulkActions   bool
		WithSoftDelete    bool
		RowActions        []generator.RowActionData
		EmptyStateMessage string
		SubmitURL         string
		Method            string
		SuccessRedirect   string
	}{
		ModulePath:   "github.com/test/testproject",
		DomainName:   "product",
		ModelName:    "Product",
		PackageName:  "product",
		VariableName: "product",
		URLPath:      "/products",
		ViewType:     "list",
		ViewName:     "ProductList",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
			{Name: "Price", Type: "float64", JSONName: "price", Required: true, Label: "Price", FormType: "number"},
		},
		Columns: []generator.ColumnData{
			{Key: "name", Label: "Name", Sortable: true},
			{Key: "price", Label: "Price", Sortable: true, Format: "currency"},
		},
		WithPagination:  true,
		WithSearch:      true,
		WithFilters:     false,
		WithSorting:     true,
		WithBulkActions: false,
		WithSoftDelete:  false,
		RowActions: []generator.RowActionData{
			{Type: "view", Label: "View", IsView: true},
			{Type: "edit", Label: "Edit", IsEdit: true},
			{Type: "delete", Label: "Delete", IsDelete: true, Confirm: true, ConfirmMessage: "Are you sure?"},
		},
		EmptyStateMessage: "No products found",
		SubmitURL:         "/products",
		Method:            "POST",
		SuccessRedirect:   "/products",
	}

	templates := []string{
		"views/list.templ.tmpl",
		"views/show.templ.tmpl",
		"views/form.templ.tmpl",
		"views/table.templ.tmpl",
		"views/partials.templ.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, viewData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}
		})
	}
}

// TestAuthTemplatesExecute tests that auth templates execute with valid data.
func TestAuthTemplatesExecute(t *testing.T) {
	authData := struct {
		ModulePath  string
		ProjectName string
		SessionType string
	}{
		ModulePath:  "github.com/test/testproject",
		ProjectName: "testproject",
		SessionType: "cookie",
	}

	templates := []string{
		"auth/user_model.go.tmpl",
		"auth/auth_middleware.go.tmpl",
		"auth/session.go.tmpl",
		"auth/user_repository.go.tmpl",
		"auth/auth_service.go.tmpl",
		"auth/auth_controller.go.tmpl",
		"auth/login.templ.tmpl",
		"auth/register.templ.tmpl",
		"auth/auth_layout.templ.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, authData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}
		})
	}
}

// TestComponentTemplatesExecute tests that component templates execute with valid data.
func TestComponentTemplatesExecute(t *testing.T) {
	// Card and form_field use ComponentData
	componentData := struct {
		ModulePath    string
		ComponentName string
		ComponentType string
		Props         []generator.PropData
		WithHTMX      bool
		AlpineState   map[string]interface{}
	}{
		ModulePath:    "github.com/test/testproject",
		ComponentName: "UserCard",
		ComponentType: "card",
		Props: []generator.PropData{
			{Name: "Title", Type: "string", Required: true},
			{Name: "Description", Type: "string", Default: ""},
		},
		WithHTMX: true,
		AlpineState: map[string]interface{}{
			"open": false,
		},
	}

	// Modal uses ModalData
	modalData := struct {
		ModulePath     string
		ModalName      string
		ModalType      string
		Title          string
		ContentType    string
		TriggerButton  string
		TriggerVariant string
		HTMXURL        string
		IsDialog       bool
		IsSheet        bool
		IsConfirm      bool
	}{
		ModulePath:     "github.com/test/testproject",
		ModalName:      "ConfirmDelete",
		ModalType:      "confirm",
		Title:          "Confirm Delete",
		ContentType:    "confirm",
		TriggerButton:  "Delete",
		TriggerVariant: "destructive",
		HTMXURL:        "/api/delete",
		IsDialog:       false,
		IsSheet:        false,
		IsConfirm:      true,
	}

	t.Run("components/card.templ.tmpl", func(t *testing.T) {
		content, err := FS.ReadFile("components/card.templ.tmpl")
		if err != nil {
			t.Fatalf("Failed to read template: %v", err)
		}

		tmpl, err := parseTemplate("card.templ.tmpl", string(content))
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, componentData)
		if err != nil {
			t.Errorf("Failed to execute template: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("Template produced empty output")
		}
	})

	t.Run("components/modal.templ.tmpl", func(t *testing.T) {
		content, err := FS.ReadFile("components/modal.templ.tmpl")
		if err != nil {
			t.Fatalf("Failed to read template: %v", err)
		}

		tmpl, err := parseTemplate("modal.templ.tmpl", string(content))
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, modalData)
		if err != nil {
			t.Errorf("Failed to execute template: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("Template produced empty output")
		}
	})

	t.Run("components/form_field.templ.tmpl", func(t *testing.T) {
		content, err := FS.ReadFile("components/form_field.templ.tmpl")
		if err != nil {
			t.Fatalf("Failed to read template: %v", err)
		}

		tmpl, err := parseTemplate("form_field.templ.tmpl", string(content))
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, componentData)
		if err != nil {
			t.Errorf("Failed to execute template: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("Template produced empty output")
		}
	})
}

// TestSeedTemplatesExecute tests that seed templates execute with valid data.
func TestSeedTemplatesExecute(t *testing.T) {
	seedData := struct {
		ModulePath       string
		DomainName       string
		ModelName        string
		TableName        string
		Fields           []generator.FieldData
		Count            int
		WithFaker        bool
		Dependencies     []string
		Relationships    []generator.SeedRelationshipData
		Distributions    []generator.SeedDistributionData
		HasRelationships bool
		HasDistributions bool
	}{
		ModulePath: "github.com/test/testproject",
		DomainName: "product",
		ModelName:  "Product",
		TableName:  "products",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name"},
			{Name: "Price", Type: "float64", JSONName: "price"},
		},
		Count:            10,
		WithFaker:        true,
		Dependencies:     []string{"user"},
		Relationships:    []generator.SeedRelationshipData{},
		Distributions:    []generator.SeedDistributionData{},
		HasRelationships: false,
		HasDistributions: false,
	}

	templates := []string{
		"seed/seeder.go.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, seedData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}
		})
	}
}

// TestConfigTemplatesExecute tests that config templates execute with valid data.
func TestConfigTemplatesExecute(t *testing.T) {
	type breadcrumbData struct {
		Key   string
		Label string
		Href  string
	}
	type columnData struct {
		Key      string
		Label    string
		Sortable bool
		Format   string
		Width    string
	}
	type filterData struct {
		Key   string
		Label string
		Type  string
	}
	type actionData struct {
		Key     string
		Label   string
		Variant string
		Icon    string
		Href    string
	}

	configData := struct {
		ConfigType        string
		Name              string
		PageName          string
		Locale            string
		Content           map[string]interface{}
		Title             string
		Description       string
		Heading           string
		Layout            string
		EmptyStateMessage string
		WithBreadcrumbs   bool
		WithTable         bool
		WithFilters       bool
		WithActions       bool
		Sidebar           bool
		Breadcrumbs       []breadcrumbData
		Columns           []columnData
		Filters           []filterData
		Actions           []actionData
	}{
		ConfigType:        "page",
		Name:              "dashboard",
		PageName:          "Dashboard",
		Locale:            "en",
		Content:           map[string]interface{}{"key": "value"},
		Title:             "Dashboard",
		Description:       "Main dashboard page",
		Heading:           "Dashboard",
		Layout:            "dashboard",
		EmptyStateMessage: "No data available",
		WithBreadcrumbs:   true,
		WithTable:         true,
		WithFilters:       false,
		WithActions:       true,
		Sidebar:           true,
		Breadcrumbs: []breadcrumbData{
			{Key: "home", Label: "Home", Href: "/"},
			{Key: "dashboard", Label: "Dashboard", Href: "/dashboard"},
		},
		Columns: []columnData{
			{Key: "name", Label: "Name", Sortable: true},
			{Key: "status", Label: "Status", Sortable: false, Format: "badge"},
		},
		Filters: []filterData{},
		Actions: []actionData{
			{Key: "create", Label: "Create", Variant: "primary", Icon: "plus", Href: "/dashboard/new"},
		},
	}

	templates := []string{
		"config/page.toml.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, configData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}
		})
	}
}

// TestDomainTemplatesWithRelationships tests domain templates with relationship data.
func TestDomainTemplatesWithRelationships(t *testing.T) {
	fkField := generator.FieldData{Name: "UserID", Type: "uint", JSONName: "user_id"}
	domainData := struct {
		ModulePath           string
		DomainName           string
		ModelName            string
		PackageName          string
		VariableName         string
		TableName            string
		URLPath              string
		Fields               []generator.FieldData
		Relationships        []generator.RelationshipData
		HasRelationships     bool
		PreloadRelationships []generator.RelationshipData
		WithSoftDelete       bool
		WithCrudViews        bool
	}{
		ModulePath:   "github.com/test/testproject",
		DomainName:   "order",
		ModelName:    "Order",
		PackageName:  "order",
		VariableName: "order",
		TableName:    "orders",
		URLPath:      "/orders",
		Fields: []generator.FieldData{
			{Name: "Total", Type: "float64", JSONName: "total", Required: true, Label: "Total"},
			{Name: "Status", Type: "string", JSONName: "status", Required: true, Label: "Status"},
		},
		Relationships: []generator.RelationshipData{
			{
				Type:            "belongs_to",
				Model:           "User",
				FieldName:       "User",
				ForeignKey:      "UserID",
				References:      "ID",
				IsBelongsTo:     true,
				GORMTag:         "foreignKey:UserID;references:ID",
				Preload:         true,
				ForeignKeyField: &fkField,
			},
			{
				Type:       "has_many",
				Model:      "OrderItem",
				FieldName:  "OrderItems",
				ForeignKey: "OrderID",
				References: "ID",
				IsHasMany:  true,
				GORMTag:    "foreignKey:OrderID;references:ID",
				Preload:    true,
			},
		},
		HasRelationships: true,
		PreloadRelationships: []generator.RelationshipData{
			{
				Type:        "belongs_to",
				Model:       "User",
				FieldName:   "User",
				ForeignKey:  "UserID",
				IsBelongsTo: true,
				Preload:     true,
			},
		},
		WithSoftDelete: true,
		WithCrudViews:  true,
	}

	templates := []string{
		"domain/model.go.tmpl",
		"domain/repository.go.tmpl",
		"domain/service.go.tmpl",
		"domain/controller.go.tmpl",
		"domain/dto.go.tmpl",
	}

	for _, tmplPath := range templates {
		t.Run(tmplPath, func(t *testing.T) {
			content, err := FS.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template: %v", err)
			}

			tmpl, err := parseTemplate(tmplPath, string(content))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, domainData)
			if err != nil {
				t.Errorf("Failed to execute template: %v", err)
			}

			if buf.Len() == 0 {
				t.Error("Template produced empty output")
			}

			// Verify relationship code is present
			output := buf.String()
			if tmplPath == "domain/model.go.tmpl" {
				if !strings.Contains(output, "User") {
					t.Error("Model template should contain User relationship")
				}
			}
		})
	}
}

// TestTemplateDelimiters verifies templates use correct delimiters.
func TestTemplateDelimiters(t *testing.T) {
	dirs := []string{
		"project",
		"domain",
		"views",
		"components",
		"config",
		"seed",
		"auth",
	}

	// Files that use {{ }} for non-Go-template purposes (e.g., taskfile variables)
	skipFiles := map[string]bool{
		"project/taskfile.yml.tmpl": true,
	}

	for _, dir := range dirs {
		entries, err := FS.ReadDir(dir)
		if err != nil {
			t.Fatalf("Failed to read directory %s: %v", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tmpl") {
				continue
			}

			templatePath := dir + "/" + entry.Name()

			// Skip files that legitimately use {{ }} for other purposes
			if skipFiles[templatePath] {
				continue
			}

			t.Run(templatePath, func(t *testing.T) {
				content, err := FS.ReadFile(templatePath)
				if err != nil {
					t.Fatalf("Failed to read template: %v", err)
				}

				contentStr := string(content)

				// Check for incorrect {{ }} delimiters in template actions
				// (they should use [[ ]] instead)
				// Note: {{ }} is valid inside templ components, so we check for template actions
				if strings.Contains(contentStr, "{{.") || strings.Contains(contentStr, "{{ .") {
					t.Errorf("Template uses incorrect {{ }} delimiters for template variables, should use [[ ]]")
				}
				if strings.Contains(contentStr, "{{if") || strings.Contains(contentStr, "{{ if") {
					t.Errorf("Template uses incorrect {{ }} delimiters for if statements, should use [[ ]]")
				}
				if strings.Contains(contentStr, "{{range") || strings.Contains(contentStr, "{{ range") {
					t.Errorf("Template uses incorrect {{ }} delimiters for range statements, should use [[ ]]")
				}
				if strings.Contains(contentStr, "{{end") || strings.Contains(contentStr, "{{ end") {
					t.Errorf("Template uses incorrect {{ }} delimiters for end statements, should use [[ ]]")
				}
			})
		}
	}
}
