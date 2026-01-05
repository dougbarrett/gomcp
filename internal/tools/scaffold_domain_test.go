package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldDomain(t *testing.T) {
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
				input := types.ScaffoldDomainInput{
					DomainName: tt.domainName,
					Fields: []types.FieldDef{
						{Name: "Name", Type: "string"},
					},
				}
				result, err := scaffoldDomain(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Error("expected failure for invalid domain name")
				}
			})
		}
	})

	t.Run("requires at least one field", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields:     []types.FieldDef{},
		}
		result, err := scaffoldDomain(registry, input)
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

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "invalid-name", Type: "string"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid field name")
		}
	})

	t.Run("validates field types", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "invalid_type"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid field type")
		}
	})

	t.Run("validates form types", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string", FormType: "invalid_form_type"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid form type")
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
		}
		result, err := scaffoldDomain(registry, input)
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

	t.Run("generates domain files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string", Required: true},
				{Name: "Price", Type: "float64"},
				{Name: "Active", Type: "bool"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check model file
		modelPath := filepath.Join(tmpDir, "internal", "models", "product.go")
		if !fileExists(modelPath) {
			t.Errorf("expected model file to be created at %s", modelPath)
		}

		// Check repository file
		repoPath := filepath.Join(tmpDir, "internal", "repository", "product", "product.go")
		if !fileExists(repoPath) {
			t.Errorf("expected repository file to be created at %s", repoPath)
		}

		// Check service file
		servicePath := filepath.Join(tmpDir, "internal", "services", "product", "product.go")
		if !fileExists(servicePath) {
			t.Errorf("expected service file to be created at %s", servicePath)
		}

		// Check DTO file
		dtoPath := filepath.Join(tmpDir, "internal", "services", "product", "dto.go")
		if !fileExists(dtoPath) {
			t.Errorf("expected DTO file to be created at %s", dtoPath)
		}

		// Check controller file
		controllerPath := filepath.Join(tmpDir, "internal", "web", "product", "product.go")
		if !fileExists(controllerPath) {
			t.Errorf("expected controller file to be created at %s", controllerPath)
		}
	})

	t.Run("generates domain with snake_case name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldDomainInput{
			DomainName: "product_category",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check that files use correct package name
		modelPath := filepath.Join(tmpDir, "internal", "models", "productcategory.go")
		if !fileExists(modelPath) {
			t.Errorf("expected model file to be created at %s", modelPath)
		}

		content := readFile(t, modelPath)
		if !strings.Contains(content, "ProductCategory") {
			t.Error("expected model to contain 'ProductCategory' struct")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
			DryRun: true,
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		modelPath := filepath.Join(tmpDir, "internal", "models", "product.go")
		if fileExists(modelPath) {
			t.Error("expected model file NOT to be created in dry run mode")
		}

		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report files that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldDomainInput{
			DomainName: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}

		hasGoModTidy := false
		hasTemplGenerate := false
		for _, step := range result.NextSteps {
			if strings.Contains(step, "go mod tidy") {
				hasGoModTidy = true
			}
			if strings.Contains(step, "templ generate") {
				hasTemplGenerate = true
			}
		}
		if !hasGoModTidy {
			t.Error("expected NextSteps to include 'go mod tidy'")
		}
		if !hasTemplGenerate {
			t.Error("expected NextSteps to include 'templ generate'")
		}
	})

	t.Run("creates multiple files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldDomainInput{
			DomainName: "user",
			Fields: []types.FieldDef{
				{Name: "Email", Type: "string", Required: true},
				{Name: "Name", Type: "string"},
				{Name: "Age", Type: "int"},
			},
		}
		result, err := scaffoldDomain(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Should create at least 5 files (model, repo, service, dto, controller)
		if len(result.FilesCreated) < 5 {
			t.Errorf("expected at least 5 files created, got %d", len(result.FilesCreated))
		}
	})
}
