package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldForm(t *testing.T) {
	t.Run("validates form name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldFormInput{
			FormName: "",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty form name")
		}
		if !strings.Contains(strings.ToLower(result.Message), "form name") {
			t.Errorf("expected error about form name, got %q", result.Message)
		}
	})

	t.Run("validates domain name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		tests := []struct {
			name       string
			domainName string
		}{
			{"empty domain", ""},
			{"invalid chars", "user@profile"},
			{"starts with number", "123user"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				input := types.ScaffoldFormInput{
					FormName: "test_form",
					Domain:   tt.domainName,
					Action:   "create",
					Fields: []types.FieldDef{
						{Name: "Name", Type: "string"},
					},
				}
				result, err := scaffoldForm(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Error("expected failure for invalid domain name")
				}
			})
		}
	})

	t.Run("validates action", func(t *testing.T) {
		registry, _ := testRegistry(t)

		tests := []struct {
			action  string
			wantErr bool
		}{
			{"create", false},
			{"edit", false},
			{"delete", true},
			{"invalid", true},
			{"", true},
		}

		for _, tt := range tests {
			t.Run(tt.action, func(t *testing.T) {
				input := types.ScaffoldFormInput{
					FormName: "test_form",
					Domain:   "product",
					Action:   tt.action,
					Fields: []types.FieldDef{
						{Name: "Name", Type: "string"},
					},
				}
				result, err := scaffoldForm(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr && result.Success {
					t.Errorf("expected failure for action %q", tt.action)
				}
				if !tt.wantErr && !result.Success {
					// This will fail because go.mod is missing, which is fine
					// We're just testing action validation at this point
					if strings.Contains(result.Message, "action") {
						t.Errorf("unexpected action error for %q: %s", tt.action, result.Message)
					}
				}
			})
		}
	})

	t.Run("requires at least one field", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldFormInput{
			FormName: "test_form",
			Domain:   "product",
			Action:   "create",
			Fields:   []types.FieldDef{},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty fields")
		}
		if !strings.Contains(strings.ToLower(result.Message), "field") {
			t.Errorf("expected error about fields, got %q", result.Message)
		}
	})

	t.Run("validates field names", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldFormInput{
			FormName: "test_form",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "invalid-name", Type: "string"},
			},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid field name")
		}
	})

	t.Run("validates form type", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldFormInput{
			FormName: "test_form",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string", FormType: "invalid_type"},
			},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid form type")
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldFormInput{
			FormName: "product_form",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
		}
		result, err := scaffoldForm(registry, input)
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

	t.Run("generates create form", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldFormInput{
			FormName: "product_create_form",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string", Required: true},
				{Name: "Price", Type: "float64"},
				{Name: "Active", Type: "bool"},
			},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_create_form.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package views") {
			t.Error("expected file to contain 'package views'")
		}
	})

	t.Run("generates edit form", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldFormInput{
			FormName: "product_edit_form",
			Domain:   "product",
			Action:   "edit",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
				{Name: "Price", Type: "float64"},
			},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_edit_form.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldFormInput{
			FormName: "product_form",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
			DryRun: true,
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_form.templ")
		if fileExists(expectedPath) {
			t.Error("expected file NOT to be created in dry run mode")
		}

		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report the file that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldFormInput{
			FormName: "product_form",
			Domain:   "product",
			Action:   "create",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
		}
		result, err := scaffoldForm(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}
	})
}
