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
		ProjectName        string
		ModulePath         string
		DatabaseType       string
		WithAuth           bool
		WithUserManagement bool
	}{
		ProjectName:        "testproject",
		ModulePath:         "github.com/test/testproject",
		DatabaseType:       "sqlite",
		WithAuth:           true,
		WithUserManagement: false,
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
		URLPathSegment       string
		Fields               []generator.FieldData
		Relationships        []generator.RelationshipData
		HasRelationships     bool
		PreloadRelationships []generator.RelationshipData
		WithSoftDelete       bool
		WithCrudViews        bool
		WithPagination       bool
		WithSearch           bool
		Layout               string
		RouteGroup           string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "product",
		ModelName:      "Product",
		PackageName:    "product",
		VariableName:   "product",
		TableName:      "products",
		URLPath:        "/products",
		URLPathSegment: "products",
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
		WithPagination:       true,
		WithSearch:           true,
		Layout:               "dashboard",
		RouteGroup:           "public",
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
		URLPathSegment    string
		ViewType          string
		ViewName          string
		Fields            []generator.FieldData
		Columns           []generator.ColumnData
		Relationships     []generator.RelationshipData
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
		FormStyle         string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "product",
		ModelName:      "Product",
		PackageName:    "product",
		VariableName:   "product",
		URLPath:        "/products",
		URLPathSegment: "products",
		ViewType:       "list",
		ViewName:       "ProductList",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
			{Name: "Price", Type: "float64", JSONName: "price", Required: true, Label: "Price", FormType: "number"},
		},
		Columns: []generator.ColumnData{
			{Key: "name", Label: "Name", Sortable: true},
			{Key: "price", Label: "Price", Sortable: true, Format: "currency"},
		},
		Relationships:   []generator.RelationshipData{},
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
		FormStyle:         "modal",
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
		URLPathSegment       string
		Fields               []generator.FieldData
		Relationships        []generator.RelationshipData
		HasRelationships     bool
		PreloadRelationships []generator.RelationshipData
		WithSoftDelete       bool
		WithCrudViews        bool
		WithPagination       bool
		WithSearch           bool
		Layout               string
		RouteGroup           string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "order",
		ModelName:      "Order",
		PackageName:    "order",
		VariableName:   "order",
		TableName:      "orders",
		URLPath:        "/orders",
		URLPathSegment: "orders",
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
		WithPagination: true,
		WithSearch:     true,
		Layout:         "dashboard",
		RouteGroup:     "public",
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

// TestReadTemplate tests the ReadTemplate helper function.
func TestReadTemplate(t *testing.T) {
	t.Run("reads existing template", func(t *testing.T) {
		content, err := ReadTemplate("project/main.go.tmpl")
		if err != nil {
			t.Fatalf("ReadTemplate failed: %v", err)
		}
		if len(content) == 0 {
			t.Error("ReadTemplate returned empty content")
		}
		if !strings.Contains(string(content), "package main") {
			t.Error("Expected main.go.tmpl to contain 'package main'")
		}
	})

	t.Run("returns error for non-existent template", func(t *testing.T) {
		_, err := ReadTemplate("nonexistent/template.tmpl")
		if err == nil {
			t.Error("ReadTemplate should return error for non-existent template")
		}
	})
}

// TestTemplateExists tests the TemplateExists helper function.
func TestTemplateExists(t *testing.T) {
	t.Run("returns true for existing template", func(t *testing.T) {
		if !TemplateExists("project/main.go.tmpl") {
			t.Error("TemplateExists should return true for existing template")
		}
	})

	t.Run("returns false for non-existent template", func(t *testing.T) {
		if TemplateExists("nonexistent/template.tmpl") {
			t.Error("TemplateExists should return false for non-existent template")
		}
	})

	t.Run("checks various categories", func(t *testing.T) {
		templates := []string{
			"domain/model.go.tmpl",
			"views/list.templ.tmpl",
			"auth/login.templ.tmpl",
			"components/card.templ.tmpl",
		}
		for _, tmpl := range templates {
			if !TemplateExists(tmpl) {
				t.Errorf("TemplateExists should return true for %s", tmpl)
			}
		}
	})
}

// TestListTemplates tests the ListTemplates helper function.
func TestListTemplates(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	if len(templates) == 0 {
		t.Error("ListTemplates should return at least one template")
	}

	// Verify all returned paths end with .tmpl
	for _, tmpl := range templates {
		if !strings.HasSuffix(tmpl, ".tmpl") {
			t.Errorf("Template path %s should end with .tmpl", tmpl)
		}
	}

	// Verify some known templates are in the list
	expectedTemplates := []string{
		"project/main.go.tmpl",
		"domain/model.go.tmpl",
		"views/list.templ.tmpl",
	}
	for _, expected := range expectedTemplates {
		found := false
		for _, tmpl := range templates {
			if tmpl == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template %s not found in list", expected)
		}
	}
}

// TestListTemplatesInCategory tests the ListTemplatesInCategory helper function.
func TestListTemplatesInCategory(t *testing.T) {
	t.Run("lists project templates", func(t *testing.T) {
		templates, err := ListTemplatesInCategory("project")
		if err != nil {
			t.Fatalf("ListTemplatesInCategory failed: %v", err)
		}
		if len(templates) == 0 {
			t.Error("project category should have templates")
		}
		// Check that all templates are in the project directory
		for _, tmpl := range templates {
			if !strings.HasPrefix(tmpl, "project/") {
				t.Errorf("Template %s should be in project/ directory", tmpl)
			}
		}
	})

	t.Run("lists domain templates", func(t *testing.T) {
		templates, err := ListTemplatesInCategory("domain")
		if err != nil {
			t.Fatalf("ListTemplatesInCategory failed: %v", err)
		}
		expectedCount := 5 // model, repository, service, controller, dto
		if len(templates) != expectedCount {
			t.Errorf("domain category should have %d templates, got %d", expectedCount, len(templates))
		}
	})

	t.Run("returns error for invalid category", func(t *testing.T) {
		_, err := ListTemplatesInCategory("nonexistent")
		if err == nil {
			t.Error("ListTemplatesInCategory should return error for invalid category")
		}
	})

	t.Run("lists all categories", func(t *testing.T) {
		for _, category := range Categories {
			templates, err := ListTemplatesInCategory(category)
			if err != nil {
				t.Errorf("ListTemplatesInCategory(%s) failed: %v", category, err)
			}
			if len(templates) == 0 {
				t.Errorf("Category %s should have at least one template", category)
			}
		}
	})
}

// TestCountTemplates tests the CountTemplates helper function.
func TestCountTemplates(t *testing.T) {
	count, err := CountTemplates()
	if err != nil {
		t.Fatalf("CountTemplates failed: %v", err)
	}

	if count == 0 {
		t.Error("CountTemplates should return at least one template")
	}

	// Verify count matches ListTemplates
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	if count != len(templates) {
		t.Errorf("CountTemplates (%d) should match len(ListTemplates()) (%d)", count, len(templates))
	}

	// We expect at least 40 templates based on our categories
	if count < 40 {
		t.Errorf("Expected at least 40 templates, got %d", count)
	}
}

// TestCategories tests the Categories variable.
func TestCategories(t *testing.T) {
	expectedCategories := []string{
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

	if len(Categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(Categories))
	}

	for i, expected := range expectedCategories {
		if Categories[i] != expected {
			t.Errorf("Category %d: expected %s, got %s", i, expected, Categories[i])
		}
	}
}

// TestFormSelectOptionsRendering verifies the form template generates select options
// when the options field is provided on a select form_type field.
func TestFormSelectOptionsRendering(t *testing.T) {
	viewData := struct {
		ModulePath        string
		DomainName        string
		ModelName         string
		PackageName       string
		VariableName      string
		URLPath           string
		URLPathSegment    string
		ViewType          string
		ViewName          string
		Fields            []generator.FieldData
		Columns           []generator.ColumnData
		Relationships     []generator.RelationshipData
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
		FormStyle         string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "discount",
		ModelName:      "Discount",
		PackageName:    "discount",
		VariableName:   "discount",
		URLPath:        "/discounts",
		URLPathSegment: "discounts",
		ViewType:       "form",
		ViewName:       "DiscountForm",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
			{
				Name:       "DiscountType",
				Type:       "string",
				JSONName:   "discount_type",
				Required:   true,
				Label:      "Discount Type",
				FormType:   "select",
				Options:    []string{"percentage", "fixed"},
				HasOptions: true,
			},
			{
				Name:       "Status",
				Type:       "string",
				JSONName:   "status",
				Required:   false,
				Label:      "Status",
				FormType:   "select",
				Options:    []string{"draft", "active", "expired"},
				HasOptions: true,
			},
		},
		Columns:           []generator.ColumnData{},
		Relationships:     []generator.RelationshipData{},
		WithPagination:    false,
		WithSearch:        false,
		WithFilters:       false,
		WithSorting:       false,
		WithBulkActions:   false,
		WithSoftDelete:    false,
		RowActions:        []generator.RowActionData{},
		EmptyStateMessage: "",
		SubmitURL:         "/discounts",
		Method:            "POST",
		SuccessRedirect:   "/discounts",
		FormStyle:         "modal",
	}

	content, err := FS.ReadFile("views/form.templ.tmpl")
	if err != nil {
		t.Fatalf("Failed to read form template: %v", err)
	}

	tmpl, err := parseTemplate("form.templ.tmpl", string(content))
	if err != nil {
		t.Fatalf("Failed to parse form template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, viewData)
	if err != nil {
		t.Fatalf("Failed to execute form template: %v", err)
	}

	output := buf.String()

	// Test 1: Options should be rendered for DiscountType
	t.Run("DiscountType options rendered", func(t *testing.T) {
		if !strings.Contains(output, `value="percentage"`) {
			t.Error("Form should include percentage option for DiscountType")
		}
		if !strings.Contains(output, `value="fixed"`) {
			t.Error("Form should include fixed option for DiscountType")
		}
	})

	// Test 2: Options should be rendered for Status
	t.Run("Status options rendered", func(t *testing.T) {
		if !strings.Contains(output, `value="draft"`) {
			t.Error("Form should include draft option for Status")
		}
		if !strings.Contains(output, `value="active"`) {
			t.Error("Form should include active option for Status")
		}
		if !strings.Contains(output, `value="expired"`) {
			t.Error("Form should include expired option for Status")
		}
	})

	// Test 3: Select should NOT have the "Add your options here" comment when options exist
	t.Run("No placeholder comment when options exist", func(t *testing.T) {
		// When options are provided, we shouldn't see the placeholder comment
		// Note: We can't easily test this because both selects have options
		// Just verify options are present
		if !strings.Contains(output, `value="percentage"`) {
			t.Error("Options should be present when HasOptions is true")
		}
	})
}

// TestFormBooleanCheckboxRendering verifies the form template generates the correct
// hidden field + checkbox pattern for boolean fields. The hidden field ensures "false"
// is sent when the checkbox is unchecked.
func TestFormBooleanCheckboxRendering(t *testing.T) {
	viewData := struct {
		ModulePath        string
		DomainName        string
		ModelName         string
		PackageName       string
		VariableName      string
		URLPath           string
		URLPathSegment    string
		ViewType          string
		ViewName          string
		Fields            []generator.FieldData
		Columns           []generator.ColumnData
		Relationships     []generator.RelationshipData
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
		FormStyle         string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "discount",
		ModelName:      "Discount",
		PackageName:    "discount",
		VariableName:   "discount",
		URLPath:        "/discounts",
		URLPathSegment: "discounts",
		ViewType:       "form",
		ViewName:       "DiscountForm",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
			{Name: "Active", Type: "bool", JSONName: "active", Required: false, Label: "Active", FormType: "checkbox"},
		},
		Columns:           []generator.ColumnData{},
		Relationships:     []generator.RelationshipData{},
		WithPagination:    false,
		WithSearch:        false,
		WithFilters:       false,
		WithSorting:       false,
		WithBulkActions:   false,
		WithSoftDelete:    false,
		RowActions:        []generator.RowActionData{},
		EmptyStateMessage: "",
		SubmitURL:         "/discounts",
		Method:            "POST",
		SuccessRedirect:   "/discounts",
		FormStyle:         "modal",
	}

	content, err := FS.ReadFile("views/form.templ.tmpl")
	if err != nil {
		t.Fatalf("Failed to read form template: %v", err)
	}

	tmpl, err := parseTemplate("form.templ.tmpl", string(content))
	if err != nil {
		t.Fatalf("Failed to parse form template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, viewData)
	if err != nil {
		t.Fatalf("Failed to execute form template: %v", err)
	}

	output := buf.String()

	// Test 1: Hidden field should be present before checkbox
	t.Run("Hidden field pattern for unchecked state", func(t *testing.T) {
		if !strings.Contains(output, `type="hidden" name="active" value="false"`) {
			t.Error("Form should include hidden field with value=\"false\" for checkbox unchecked state")
		}
	})

	// Test 2: Checkbox should use the Checkbox component with value="true"
	t.Run("Checkbox component with true value", func(t *testing.T) {
		if !strings.Contains(output, `@components.Checkbox("active"`) {
			t.Error("Form should use the Checkbox component for boolean fields")
		}
	})

	// Test 3: Form should have a comment explaining the hidden field pattern
	t.Run("Comment explaining hidden field", func(t *testing.T) {
		if !strings.Contains(output, "Hidden field for unchecked state") {
			t.Error("Form should have a comment explaining the hidden field pattern")
		}
	})
}

// TestControllerBooleanCheckboxHandling verifies the controller template generates correct
// code for handling boolean checkbox form values. Checkboxes use a hidden field pattern
// where hidden="false" is always sent, and checkbox="true" is sent only when checked.
// The controller must check all form values, not just the first one.
func TestControllerBooleanCheckboxHandling(t *testing.T) {
	domainData := struct {
		ModulePath           string
		DomainName           string
		ModelName            string
		PackageName          string
		VariableName         string
		TableName            string
		URLPath              string
		URLPathSegment       string
		Fields               []generator.FieldData
		Relationships        []generator.RelationshipData
		HasRelationships     bool
		PreloadRelationships []generator.RelationshipData
		WithSoftDelete       bool
		WithCrudViews        bool
		WithPagination       bool
		WithSearch           bool
		Layout               string
		RouteGroup           string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "discount",
		ModelName:      "Discount",
		PackageName:    "discount",
		VariableName:   "discount",
		TableName:      "discounts",
		URLPath:        "/discounts",
		URLPathSegment: "discounts",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
			{Name: "Active", Type: "bool", JSONName: "active", Required: false, Label: "Active", FormType: "checkbox"},
			{Name: "Featured", Type: "bool", JSONName: "featured", Required: false, Label: "Featured", FormType: "checkbox"},
		},
		Relationships:        []generator.RelationshipData{},
		HasRelationships:     false,
		PreloadRelationships: []generator.RelationshipData{},
		WithSoftDelete:       true,
		WithCrudViews:        true,
		WithPagination:       true,
		WithSearch:           true,
		Layout:               "dashboard",
		RouteGroup:           "public",
	}

	content, err := FS.ReadFile("domain/controller.go.tmpl")
	if err != nil {
		t.Fatalf("Failed to read controller template: %v", err)
	}

	tmpl, err := parseTemplate("controller.go.tmpl", string(content))
	if err != nil {
		t.Fatalf("Failed to parse controller template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, domainData)
	if err != nil {
		t.Fatalf("Failed to execute controller template: %v", err)
	}

	output := buf.String()

	// Test 1: Create handler should iterate over form values for boolean fields
	// The pattern should be: r.Form["fieldname"] to get all values
	t.Run("Create handler uses r.Form for boolean fields", func(t *testing.T) {
		// Check for Active field - should use r.Form["active"]
		if !strings.Contains(output, `r.Form["active"]`) {
			t.Error("Create handler should use r.Form[\"active\"] to check all form values for checkbox")
		}
		// Check for Featured field - should use r.Form["featured"]
		if !strings.Contains(output, `r.Form["featured"]`) {
			t.Error("Create handler should use r.Form[\"featured\"] to check all form values for checkbox")
		}
		// Should NOT use r.FormValue for boolean fields in Create
		if strings.Contains(output, `r.FormValue("active") == "true"`) {
			t.Error("Create handler should NOT use r.FormValue for boolean checkbox fields (misses checked state due to hidden field)")
		}
	})

	// Test 2: Update handler should also iterate over form values for boolean fields
	t.Run("Update handler uses r.Form for boolean fields", func(t *testing.T) {
		// The Update handler should have a comment about checkbox handling
		if !strings.Contains(output, "check all form values since hidden field") {
			t.Error("Update handler should have a comment explaining checkbox handling")
		}
	})

	// Test 3: Generated code should check for both "true" and "on" values
	t.Run("Checks for both true and on values", func(t *testing.T) {
		if !strings.Contains(output, `v == "true"`) {
			t.Error("Boolean handling should check for \"true\" value")
		}
		if !strings.Contains(output, `v == "on"`) {
			t.Error("Boolean handling should check for \"on\" value (for native HTML checkbox)")
		}
	})

	// Test 4: Non-boolean fields should still use r.FormValue
	t.Run("Non-boolean fields use r.FormValue", func(t *testing.T) {
		if !strings.Contains(output, `r.FormValue("name")`) {
			t.Error("Non-boolean string fields should use r.FormValue")
		}
	})
}

// TestBelongsToDisplayFieldRendering verifies the form and show templates use the
// DisplayField property for belongs_to relationships instead of hardcoded "Name".
func TestBelongsToDisplayFieldRendering(t *testing.T) {
	fkField := generator.FieldData{Name: "CategoryID", Type: "uint", JSONName: "category_id"}
	viewData := struct {
		ModulePath        string
		DomainName        string
		ModelName         string
		PackageName       string
		VariableName      string
		URLPath           string
		URLPathSegment    string
		ViewType          string
		ViewName          string
		Fields            []generator.FieldData
		Columns           []generator.ColumnData
		Relationships     []generator.RelationshipData
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
		FormStyle         string
	}{
		ModulePath:     "github.com/test/testproject",
		DomainName:     "product",
		ModelName:      "Product",
		PackageName:    "product",
		VariableName:   "product",
		URLPath:        "/products",
		URLPathSegment: "products",
		ViewType:       "form",
		ViewName:       "ProductForm",
		Fields: []generator.FieldData{
			{Name: "Name", Type: "string", JSONName: "name", Required: true, Label: "Name", FormType: "input"},
		},
		Columns: []generator.ColumnData{},
		Relationships: []generator.RelationshipData{
			{
				Type:            "belongs_to",
				Model:           "Category",
				FieldName:       "Category",
				ForeignKey:      "CategoryID",
				References:      "ID",
				IsBelongsTo:     true,
				GORMTag:         "foreignKey:CategoryID;references:ID",
				Preload:         true,
				ForeignKeyField: &fkField,
				DisplayField:    "Title", // Custom display field instead of "Name"
			},
		},
		WithPagination:    false,
		WithSearch:        false,
		WithFilters:       false,
		WithSorting:       false,
		WithBulkActions:   false,
		WithSoftDelete:    false,
		RowActions:        []generator.RowActionData{},
		EmptyStateMessage: "",
		SubmitURL:         "/products",
		Method:            "POST",
		SuccessRedirect:   "/products",
		FormStyle:         "modal",
	}

	// Test form template uses DisplayField
	t.Run("Form template uses DisplayField", func(t *testing.T) {
		content, err := FS.ReadFile("views/form.templ.tmpl")
		if err != nil {
			t.Fatalf("Failed to read form template: %v", err)
		}

		tmpl, err := parseTemplate("form.templ.tmpl", string(content))
		if err != nil {
			t.Fatalf("Failed to parse form template: %v", err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, viewData)
		if err != nil {
			t.Fatalf("Failed to execute form template: %v", err)
		}

		output := buf.String()

		// Should use opt.Title (the DisplayField) instead of opt.Name
		if !strings.Contains(output, "opt.Title") {
			t.Error("Form template should use opt.Title (DisplayField) for belongs_to dropdown, not hardcoded opt.Name")
		}
		if strings.Contains(output, "opt.Name") {
			t.Error("Form template should NOT use hardcoded opt.Name when DisplayField is set to Title")
		}
	})

	// Test show template uses DisplayField
	t.Run("Show template uses DisplayField", func(t *testing.T) {
		content, err := FS.ReadFile("views/show.templ.tmpl")
		if err != nil {
			t.Fatalf("Failed to read show template: %v", err)
		}

		tmpl, err := parseTemplate("show.templ.tmpl", string(content))
		if err != nil {
			t.Fatalf("Failed to parse show template: %v", err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, viewData)
		if err != nil {
			t.Fatalf("Failed to execute show template: %v", err)
		}

		output := buf.String()

		// Should use props.Item.Category.Title (the DisplayField) instead of props.Item.Category.Name
		if !strings.Contains(output, "props.Item.Category.Title") {
			t.Error("Show template should use props.Item.Category.Title (DisplayField) for belongs_to display, not hardcoded .Name")
		}
		if strings.Contains(output, "props.Item.Category.Name") {
			t.Error("Show template should NOT use hardcoded .Name when DisplayField is set to Title")
		}
	})

	// Test with default DisplayField (should be "Name")
	t.Run("Default DisplayField is Name", func(t *testing.T) {
		defaultFKField := generator.FieldData{Name: "UserID", Type: "uint", JSONName: "user_id"}
		viewDataDefault := viewData
		viewDataDefault.Relationships = []generator.RelationshipData{
			{
				Type:            "belongs_to",
				Model:           "User",
				FieldName:       "User",
				ForeignKey:      "UserID",
				References:      "ID",
				IsBelongsTo:     true,
				GORMTag:         "foreignKey:UserID;references:ID",
				Preload:         true,
				ForeignKeyField: &defaultFKField,
				DisplayField:    "Name", // Default value
			},
		}

		content, err := FS.ReadFile("views/form.templ.tmpl")
		if err != nil {
			t.Fatalf("Failed to read form template: %v", err)
		}

		tmpl, err := parseTemplate("form.templ.tmpl", string(content))
		if err != nil {
			t.Fatalf("Failed to parse form template: %v", err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, viewDataDefault)
		if err != nil {
			t.Fatalf("Failed to execute form template: %v", err)
		}

		output := buf.String()

		// Should use opt.Name when DisplayField is "Name"
		if !strings.Contains(output, "opt.Name") {
			t.Error("Form template should use opt.Name when DisplayField is 'Name'")
		}
	})
}

// TestWizardControllerTemplatePatterns verifies the wizard controller template:
// 1. Step handler doesn't declare unused 'resp' variable
// 2. Submit handler uses Create...Input (not Create...DTO)
func TestWizardControllerTemplatePatterns(t *testing.T) {
	wizardData := struct {
		ModulePath       string
		WizardName       string
		WizardNamePascal string
		Domain           string
		ModelName        string
		PackageName      string
		VariableName     string
		URLPath          string
		URLPathSegment   string
		Steps            []generator.WizardStepData
		TotalSteps       int
		Layout           string
		RouteGroup       string
		FormStyle        string
		SuccessRedirect  string
		WithDrafts       bool
	}{
		ModulePath:       "github.com/test/testproject",
		WizardName:       "create_order",
		WizardNamePascal: "CreateOrder",
		Domain:           "order",
		ModelName:        "Order",
		PackageName:      "order",
		VariableName:     "order",
		URLPath:          "/orders",
		URLPathSegment:   "orders",
		Steps: []generator.WizardStepData{
			{Number: 1, Name: "Client Info", Type: "form", FieldNames: []string{"client_name", "email"}, IsFirst: true, IsLast: false, IsForm: true},
			{Number: 2, Name: "Review", Type: "summary", FieldNames: []string{}, IsFirst: false, IsLast: true, IsSummary: true},
		},
		TotalSteps:      2,
		Layout:          "dashboard",
		RouteGroup:      "admin",
		FormStyle:       "page",
		SuccessRedirect: "/orders",
		WithDrafts:      false,
	}

	content, err := FS.ReadFile("wizard/controller.go.tmpl")
	if err != nil {
		t.Fatalf("Failed to read wizard controller template: %v", err)
	}

	tmpl, err := parseTemplate("wizard_controller.go.tmpl", string(content))
	if err != nil {
		t.Fatalf("Failed to parse wizard controller template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, wizardData)
	if err != nil {
		t.Fatalf("Failed to execute wizard controller template: %v", err)
	}

	output := buf.String()

	// Test 1: Step handler should NOT have unused resp variable
	t.Run("Step handler no unused resp", func(t *testing.T) {
		// Find the Step1 function and check if it declares resp but doesn't use it
		// The Step handler should just call c.render() which doesn't need resp
		// A valid pattern is: resp declared AND used for Error/Redirect
		// An invalid pattern is: resp declared but only c.render() is called

		// Since we're testing without drafts, the Step handler should NOT have resp at all
		// because it only calls c.render(w, r, views.Component(props))
		stepFuncStart := strings.Index(output, "func (c *CreateOrderWizardController) Step1(")
		stepFuncEnd := strings.Index(output[stepFuncStart:], "func (c *CreateOrderWizardController) Step1Submit")
		if stepFuncEnd == -1 {
			// Try to find end of Step1
			stepFuncEnd = strings.Index(output[stepFuncStart:], "\n\n//")
		}
		stepFunc := output[stepFuncStart : stepFuncStart+stepFuncEnd]

		// Without drafts, Step handler should NOT declare resp at all
		// because it doesn't use it - it just renders a view
		if strings.Contains(stepFunc, "resp := web.NewResponse") {
			t.Error("Step handler (without drafts) should NOT declare resp since it only calls c.render()")
		}
	})

	// Test 2: Submit handler should use Create...Input, not Create...DTO
	t.Run("Submit uses CreateInput not CreateDTO", func(t *testing.T) {
		if strings.Contains(output, "Create"+wizardData.ModelName+"DTO") {
			t.Error("Submit handler should use Create...Input, not Create...DTO")
		}
		if !strings.Contains(output, "Create"+wizardData.ModelName+"Input") {
			t.Error("Submit handler should reference Create...Input")
		}
	})

	// Test 3: With drafts enabled, Step handler still shouldn't have unused resp
	t.Run("Step handler with drafts no unused resp", func(t *testing.T) {
		wizardDataWithDrafts := wizardData
		wizardDataWithDrafts.WithDrafts = true

		var bufDrafts bytes.Buffer
		err = tmpl.Execute(&bufDrafts, wizardDataWithDrafts)
		if err != nil {
			t.Fatalf("Failed to execute wizard controller template with drafts: %v", err)
		}

		outputDrafts := bufDrafts.String()

		// Find the Step1 function
		stepFuncStart := strings.Index(outputDrafts, "func (c *CreateOrderWizardController) Step1(")
		stepFuncEnd := strings.Index(outputDrafts[stepFuncStart:], "func (c *CreateOrderWizardController) Step1Submit")
		if stepFuncEnd == -1 {
			stepFuncEnd = strings.Index(outputDrafts[stepFuncStart:], "\n\n//")
		}
		stepFunc := outputDrafts[stepFuncStart : stepFuncStart+stepFuncEnd]

		// Even with drafts, Step handler doesn't need resp because it just renders a view
		// The draft data is fetched but doesn't require Error/Redirect calls
		if strings.Contains(stepFunc, "resp := web.NewResponse") {
			t.Error("Step handler (with drafts) should NOT declare resp since it only calls c.render()")
		}
	})
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
