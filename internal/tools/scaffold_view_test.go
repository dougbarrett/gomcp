package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldView(t *testing.T) {
	t.Run("validates domain name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		tests := []struct {
			name       string
			domainName string
			wantErr    string
		}{
			{"empty domain", "", "domain name"},
			{"invalid chars", "user@profile", "domain name"},
			{"starts with number", "123user", "domain name"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				input := types.ScaffoldViewInput{
					DomainName: tt.domainName,
					ViewType:   "list",
					ViewName:   "test",
				}
				result, err := scaffoldView(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Error("expected failure for invalid domain name")
				}
				if !strings.Contains(strings.ToLower(result.Message), tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, result.Message)
				}
			})
		}
	})

	t.Run("validates view type", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "invalid_type",
			ViewName:   "test",
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid view type")
		}
		if !strings.Contains(strings.ToLower(result.Message), "view type") {
			t.Errorf("expected error about view type, got %q", result.Message)
		}
	})

	t.Run("requires view name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "list",
			ViewName:   "",
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty view name")
		}
		if !strings.Contains(strings.ToLower(result.Message), "view name") {
			t.Errorf("expected error about view name, got %q", result.Message)
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "list",
			ViewName:   "product_list",
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure when go.mod is missing")
		}
		if !strings.Contains(strings.ToLower(result.Message), "module path") {
			t.Errorf("expected error about module path, got %q", result.Message)
		}
	})

	t.Run("generates list view", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "list",
			ViewName:   "product_list",
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check file was created
		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_list.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		// Check file content
		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package views") {
			t.Error("expected file to contain 'package views'")
		}
		if !strings.Contains(content, "Product") {
			t.Error("expected file to contain model name 'Product'")
		}
	})

	t.Run("generates show view", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "show",
			ViewName:   "product_show",
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_show.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("generates form view", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "form",
			ViewName:   "product_form",
			Config: types.ViewConfig{
				Fields: []types.FieldDef{
					{Name: "Name", Type: "string"},
					{Name: "Price", Type: "float64"},
				},
			},
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_form.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("generates table view", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "table",
			ViewName:   "product_table",
			Config: types.ViewConfig{
				Columns: []types.ColumnDef{
					{Key: "name", Label: "Name", Sortable: true},
					{Key: "price", Label: "Price", Sortable: true},
				},
				WithPagination: true,
				WithSearch:     true,
			},
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_table.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "list",
			ViewName:   "product_list",
			DryRun:     true,
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// File should NOT be created
		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_list.templ")
		if fileExists(expectedPath) {
			t.Error("expected file NOT to be created in dry run mode")
		}

		// But should report what would be created
		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report the file that would be created")
		}
	})

	t.Run("validates field definitions", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "form",
			ViewName:   "product_form",
			Config: types.ViewConfig{
				Fields: []types.FieldDef{
					{Name: "invalid-name", Type: "string"},
				},
			},
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid field name")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldViewInput{
			DomainName: "product",
			ViewType:   "list",
			ViewName:   "product_list",
		}
		result, err := scaffoldView(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}

		// Should include templ generate
		hasTemplGenerate := false
		for _, step := range result.NextSteps {
			if strings.Contains(step, "templ generate") {
				hasTemplGenerate = true
				break
			}
		}
		if !hasTemplGenerate {
			t.Error("expected NextSteps to include 'templ generate'")
		}
	})
}

func TestGetViewTemplatePath(t *testing.T) {
	tests := []struct {
		viewType string
		expected string
	}{
		{"list", "views/list.templ.tmpl"},
		{"show", "views/show.templ.tmpl"},
		{"form", "views/form.templ.tmpl"},
		{"table", "views/table.templ.tmpl"},
		{"card", "views/partials.templ.tmpl"},
		{"custom", "views/list.templ.tmpl"},
		{"unknown", "views/list.templ.tmpl"},
	}

	for _, tt := range tests {
		t.Run(tt.viewType, func(t *testing.T) {
			got := getViewTemplatePath(tt.viewType)
			if got != tt.expected {
				t.Errorf("getViewTemplatePath(%q) = %q, want %q", tt.viewType, got, tt.expected)
			}
		})
	}
}

func TestBuildViewData(t *testing.T) {
	input := types.ScaffoldViewInput{
		DomainName: "product_category",
		ViewType:   "list",
		ViewName:   "category_list",
		Config: types.ViewConfig{
			WithPagination:    true,
			WithSearch:        true,
			EmptyStateMessage: "No categories yet",
		},
	}

	data := buildViewData(input, "github.com/example/app")

	if data.ModulePath != "github.com/example/app" {
		t.Errorf("expected ModulePath to be github.com/example/app, got %s", data.ModulePath)
	}
	if data.DomainName != "product_category" {
		t.Errorf("expected DomainName to be product_category, got %s", data.DomainName)
	}
	if data.ModelName != "ProductCategory" {
		t.Errorf("expected ModelName to be ProductCategory, got %s", data.ModelName)
	}
	if data.PackageName != "productcategory" {
		t.Errorf("expected PackageName to be productcategory, got %s", data.PackageName)
	}
	if data.URLPath != "/product-categories" {
		t.Errorf("expected URLPath to be /product-categories, got %s", data.URLPath)
	}
	if !data.WithPagination {
		t.Error("expected WithPagination to be true")
	}
	if !data.WithSearch {
		t.Error("expected WithSearch to be true")
	}
	if data.EmptyStateMessage != "No categories yet" {
		t.Errorf("expected EmptyStateMessage to be 'No categories yet', got %s", data.EmptyStateMessage)
	}
}
