package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

// Edge case tests verify behavior with unusual inputs and boundary conditions.

func TestEdgeCases_DomainNames(t *testing.T) {
	t.Run("single character domain name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "a",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
			DryRun: true,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Single char should be valid
		if !result.Success {
			t.Errorf("single character domain should be valid, got: %s", result.Message)
		}
	})

	t.Run("very long domain name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		longName := strings.Repeat("a", 100)
		input := types.ScaffoldDomainInput{
			DomainName: longName,
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
			DryRun: true,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Very long names should still work (filesystem limits may apply later)
		_ = result
	})

	t.Run("domain name with underscores", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "order_line_item",
			Fields: []types.FieldDef{
				{Name: "Quantity", Type: "int"},
			},
			DryRun: false,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check that model name is properly converted to PascalCase
		// Package name removes underscores: order_line_item -> orderlineitem
		modelPath := filepath.Join(tmpDir, "internal/models/orderlineitem.go")
		if !fileExists(modelPath) {
			t.Errorf("expected model file to exist at %s", modelPath)
		} else {
			content := readFile(t, modelPath)
			if !strings.Contains(content, "OrderLineItem") {
				t.Errorf("expected OrderLineItem struct in model")
			}
		}
	})

	t.Run("domain name that is Go keyword", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		// "type" is a Go keyword
		input := types.ScaffoldDomainInput{
			DomainName: "type",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
			DryRun: true,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should fail because "type" is a reserved keyword
		if result.Success {
			t.Errorf("Go keyword 'type' should be rejected as domain name")
		}
	})
}

func TestEdgeCases_FieldTypes(t *testing.T) {
	t.Run("all supported field types", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "allfieldtypes",
			Fields: []types.FieldDef{
				{Name: "StringField", Type: "string"},
				{Name: "IntField", Type: "int"},
				{Name: "Int64Field", Type: "int64"},
				{Name: "Float32Field", Type: "float32"},
				{Name: "Float64Field", Type: "float64"},
				{Name: "BoolField", Type: "bool"},
				{Name: "TimeField", Type: "time.Time"},
				{Name: "UintField", Type: "uint"},
			},
			DryRun: false,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success with all field types, got: %s", result.Message)
		}

		// Verify model contains all fields
		modelPath := filepath.Join(tmpDir, "internal/models/allfieldtypes.go")
		content := readFile(t, modelPath)

		expectedFields := []string{"StringField", "IntField", "Int64Field", "Float32Field", "Float64Field", "BoolField", "TimeField", "UintField"}
		for _, field := range expectedFields {
			if !strings.Contains(content, field) {
				t.Errorf("expected model to contain field %s", field)
			}
		}
	})

	t.Run("pointer types", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "pointertypes",
			Fields: []types.FieldDef{
				{Name: "OptionalString", Type: "*string"},
				{Name: "OptionalInt", Type: "*int"},
			},
			DryRun: false,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success with pointer types, got: %s", result.Message)
		}
	})
}

func TestEdgeCases_ProjectNames(t *testing.T) {
	t.Run("project name with hyphens", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldProjectInput{
			ProjectName: "my-awesome-project",
			ModulePath:  "github.com/test/my-awesome-project",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("project name with hyphens should be valid, got: %s", result.Message)
		}
	})

	t.Run("project name with underscores", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldProjectInput{
			ProjectName: "my_awesome_project",
			ModulePath:  "github.com/test/my_awesome_project",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("project name with underscores should be valid, got: %s", result.Message)
		}
	})

	t.Run("project name starting with lowercase", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldProjectInput{
			ProjectName: "myproject",
			ModulePath:  "github.com/test/myproject",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("lowercase project name should be valid, got: %s", result.Message)
		}
	})
}

func TestEdgeCases_ModulePaths(t *testing.T) {
	t.Run("simple module path", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "myapp",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("simple module path should be valid, got: %s", result.Message)
		}
	})

	t.Run("deeply nested module path", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/myorg/platform/services/api/myapp",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("deeply nested module path should be valid, got: %s", result.Message)
		}
	})

	t.Run("module path with version", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/myorg/myapp/v2",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("module path with version should be valid, got: %s", result.Message)
		}
	})
}

func TestEdgeCases_EmptyFields(t *testing.T) {
	t.Run("domain with many fields", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		// Create 50 fields
		fields := make([]types.FieldDef, 50)
		for i := 0; i < 50; i++ {
			fields[i] = types.FieldDef{
				Name: strings.Repeat("Field", 1) + string(rune('A'+i%26)) + string(rune('0'+i/26)),
				Type: "string",
			}
		}

		input := types.ScaffoldDomainInput{
			DomainName: "manyfields",
			Fields:     fields,
			DryRun:     true,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("domain with many fields should work, got: %s", result.Message)
		}
	})
}

func TestEdgeCases_FileOverwrite(t *testing.T) {
	t.Run("scaffolding same domain twice", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
			DryRun: false,
		}

		// First scaffold
		result1, err := scaffoldDomain(registry, input)
		if err != nil || !result1.Success {
			t.Fatalf("first scaffold failed: %v / %s", err, result1.Message)
		}

		// Second scaffold with same name (should work - files get overwritten)
		result2, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("second scaffold error: %v", err)
		}
		// This may succeed or fail depending on implementation
		_ = result2
	})
}

func TestEdgeCases_SpecialCharacters(t *testing.T) {
	t.Run("field name with numbers", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "item",
			Fields: []types.FieldDef{
				{Name: "Field1", Type: "string"},
				{Name: "Field2Value", Type: "string"},
			},
			DryRun: true,
		}

		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("field names with numbers should be valid, got: %s", result.Message)
		}
	})
}

func TestEdgeCases_ConfigLocales(t *testing.T) {
	t.Run("various locale formats", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		locales := []string{"en", "es", "fr", "de", "pt"}
		for _, locale := range locales {
			input := types.ScaffoldConfigInput{
				ConfigType: "page",
				Name:       "test",
				Locale:     locale,
				DryRun:     false,
			}

			result, err := scaffoldConfig(registry, input)
			if err != nil {
				t.Fatalf("unexpected error for locale %s: %v", locale, err)
			}
			if !result.Success {
				t.Errorf("locale %s should be valid, got: %s", locale, result.Message)
			}

			// Verify file was created in correct locale directory
			configPath := filepath.Join(tmpDir, "config", locale, "pages", "test.toml")
			if !fileExists(configPath) {
				t.Errorf("expected config file for locale %s to exist", locale)
			}
		}
	})
}

func TestEdgeCases_DryRunConsistency(t *testing.T) {
	t.Run("dry run returns same files as actual run", func(t *testing.T) {
		registry1, tmpDir1 := testRegistry(t)
		setupGoMod(t, tmpDir1, "github.com/test/project")

		registry2, tmpDir2 := testRegistry(t)
		setupGoMod(t, tmpDir2, "github.com/test/project")

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
				{Name: "Price", Type: "float64"},
			},
		}

		// Dry run
		input.DryRun = true
		dryResult, err := scaffoldDomain(registry1, input)
		if err != nil || !dryResult.Success {
			t.Fatalf("dry run failed: %v / %s", err, dryResult.Message)
		}

		// Actual run
		input.DryRun = false
		actualResult, err := scaffoldDomain(registry2, input)
		if err != nil || !actualResult.Success {
			t.Fatalf("actual run failed: %v / %s", err, actualResult.Message)
		}

		// File counts should match
		if len(dryResult.FilesCreated) != len(actualResult.FilesCreated) {
			t.Errorf("dry run reported %d files, actual created %d files",
				len(dryResult.FilesCreated), len(actualResult.FilesCreated))
		}

		// Dry run should not have created files
		for _, f := range dryResult.FilesCreated {
			fullPath := filepath.Join(tmpDir1, f)
			if fileExists(fullPath) {
				t.Errorf("dry run should not create file %s", f)
			}
		}

		// Actual run should have created files
		for _, f := range actualResult.FilesCreated {
			fullPath := filepath.Join(tmpDir2, f)
			if !fileExists(fullPath) {
				t.Errorf("actual run should create file %s", f)
			}
		}
	})
}

func TestEdgeCases_ConcurrentScaffolding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping concurrent test in short mode")
	}

	t.Run("concurrent domain scaffolding", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		domains := []string{"user", "product", "order", "category", "tag"}
		errors := make(chan error, len(domains))

		for _, domain := range domains {
			go func(d string) {
				input := types.ScaffoldDomainInput{
					DomainName: d,
					Fields: []types.FieldDef{
						{Name: "Name", Type: "string"},
					},
					DryRun: false,
				}
				result, err := scaffoldDomain(registry, input)
				if err != nil {
					errors <- err
				} else if !result.Success {
					errors <- os.ErrInvalid // Use a generic error
				} else {
					errors <- nil
				}
			}(domain)
		}

		// Collect results
		for i := 0; i < len(domains); i++ {
			if err := <-errors; err != nil {
				t.Errorf("concurrent scaffolding failed: %v", err)
			}
		}

		// Verify all domains were created
		for _, domain := range domains {
			modelPath := filepath.Join(tmpDir, "internal/models", domain+".go")
			if !fileExists(modelPath) {
				t.Errorf("expected model for %s to exist after concurrent scaffolding", domain)
			}
		}
	})
}
