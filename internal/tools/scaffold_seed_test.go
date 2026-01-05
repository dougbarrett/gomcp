package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldSeed(t *testing.T) {
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
				input := types.ScaffoldSeedInput{
					Domain: tt.domainName,
				}
				result, err := scaffoldSeed(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Error("expected failure for invalid domain name")
				}
			})
		}
	})

	t.Run("validates dependencies", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldSeedInput{
			Domain:       "product",
			Dependencies: []string{"valid_dep", "invalid@dep"},
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid dependency name")
		}
	})

	t.Run("validates field names", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldSeedInput{
			Domain: "product",
			Fields: []types.FieldDef{
				{Name: "invalid-name", Type: "string"},
			},
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid field name")
		}
	})

	t.Run("validates field types", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldSeedInput{
			Domain: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "invalid_type"},
			},
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid field type")
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldSeedInput{
			Domain: "product",
		}
		result, err := scaffoldSeed(registry, input)
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

	t.Run("generates seeder", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldSeedInput{
			Domain: "product",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
				{Name: "Price", Type: "float64"},
			},
			Count: 50,
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "cmd", "seed", "seeders", "product_seeder.go")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package seeders") {
			t.Error("expected file to contain 'package seeders'")
		}
		if !strings.Contains(content, "ProductSeeder") {
			t.Error("expected file to contain 'ProductSeeder'")
		}
	})

	t.Run("generates seeder with faker", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldSeedInput{
			Domain: "user",
			Fields: []types.FieldDef{
				{Name: "Name", Type: "string"},
				{Name: "Email", Type: "string"},
			},
			WithFaker: true,
			Count:     100,
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "cmd", "seed", "seeders", "user_seeder.go")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "gofakeit") {
			t.Error("expected file to contain gofakeit import when WithFaker is true")
		}
	})

	t.Run("generates seeder with dependencies", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldSeedInput{
			Domain: "order",
			Fields: []types.FieldDef{
				{Name: "Total", Type: "float64"},
			},
			Dependencies: []string{"user", "product"},
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "cmd", "seed", "seeders", "order_seeder.go")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "Dependencies") {
			t.Error("expected file to contain Dependencies function")
		}
	})

	t.Run("defaults count to 10", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldSeedInput{
			Domain: "product",
			// Count not specified, should default to 10
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "cmd", "seed", "seeders", "product_seeder.go")
		content := readFile(t, expectedPath)
		if !strings.Contains(content, "10") {
			t.Error("expected default count of 10 in generated file")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldSeedInput{
			Domain: "product",
			DryRun: true,
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "cmd", "seed", "seeders", "product_seeder.go")
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

		input := types.ScaffoldSeedInput{
			Domain: "product",
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}
	})

	t.Run("returns faker step when WithFaker", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldSeedInput{
			Domain:    "product",
			WithFaker: true,
		}
		result, err := scaffoldSeed(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hasFakerStep := false
		for _, step := range result.NextSteps {
			if strings.Contains(step, "gofakeit") {
				hasFakerStep = true
				break
			}
		}
		if !hasFakerStep {
			t.Error("expected NextSteps to include gofakeit installation step")
		}
	})
}

func TestBuildSeedData(t *testing.T) {
	input := types.ScaffoldSeedInput{
		Domain: "product_category",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string"},
			{Name: "Active", Type: "bool"},
		},
		Count:        25,
		WithFaker:    true,
		Dependencies: []string{"product"},
	}

	data := buildSeedData(input, "github.com/example/app")

	if data.ModulePath != "github.com/example/app" {
		t.Errorf("expected ModulePath to be github.com/example/app, got %s", data.ModulePath)
	}
	if data.DomainName != "product_category" {
		t.Errorf("expected DomainName to be product_category, got %s", data.DomainName)
	}
	if data.ModelName != "ProductCategory" {
		t.Errorf("expected ModelName to be ProductCategory, got %s", data.ModelName)
	}
	if data.TableName != "product_categories" {
		t.Errorf("expected TableName to be product_categories, got %s", data.TableName)
	}
	if data.Count != 25 {
		t.Errorf("expected Count to be 25, got %d", data.Count)
	}
	if !data.WithFaker {
		t.Error("expected WithFaker to be true")
	}
	if len(data.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(data.Fields))
	}
	if len(data.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(data.Dependencies))
	}
}

func TestBuildSeedData_Relationships(t *testing.T) {
	input := types.ScaffoldSeedInput{
		Domain: "order",
		Fields: []types.FieldDef{
			{Name: "Total", Type: "float64"},
		},
		Count:     20,
		WithFaker: true,
		Relationships: []types.SeedRelationshipDef{
			{Field: "UserID", Model: "User", Strategy: "random"},
			{Field: "ProductID", Model: "Product", Strategy: "distribute"},
		},
	}

	data := buildSeedData(input, "github.com/example/app")

	if !data.HasRelationships {
		t.Error("expected HasRelationships to be true")
	}
	if len(data.Relationships) != 2 {
		t.Errorf("expected 2 relationships, got %d", len(data.Relationships))
	}

	// Check first relationship
	rel := data.Relationships[0]
	if rel.Field != "UserID" {
		t.Errorf("expected Field to be UserID, got %s", rel.Field)
	}
	if rel.Model != "User" {
		t.Errorf("expected Model to be User, got %s", rel.Model)
	}
	if rel.ModelVar != "user" {
		t.Errorf("expected ModelVar to be user, got %s", rel.ModelVar)
	}
	if rel.Strategy != "random" {
		t.Errorf("expected Strategy to be random, got %s", rel.Strategy)
	}

	// Check second relationship
	rel2 := data.Relationships[1]
	if rel2.Strategy != "distribute" {
		t.Errorf("expected Strategy to be distribute, got %s", rel2.Strategy)
	}
}

func TestBuildSeedData_Distributions(t *testing.T) {
	input := types.ScaffoldSeedInput{
		Domain: "user",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string"},
			{Name: "Role", Type: "string"},
		},
		Count:     10,
		WithFaker: true,
		Distributions: []types.SeedDistributionDef{
			{
				Field: "Role",
				Values: []types.SeedValueDef{
					{Value: `"admin"`, Count: 2},
					{Value: `"user"`, Count: 8},
				},
			},
		},
	}

	data := buildSeedData(input, "github.com/example/app")

	if !data.HasDistributions {
		t.Error("expected HasDistributions to be true")
	}
	if len(data.Distributions) != 1 {
		t.Errorf("expected 1 distribution, got %d", len(data.Distributions))
	}

	dist := data.Distributions[0]
	if dist.Field != "Role" {
		t.Errorf("expected Field to be Role, got %s", dist.Field)
	}
	if dist.TotalCount != 10 {
		t.Errorf("expected TotalCount to be 10, got %d", dist.TotalCount)
	}
	if len(dist.Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(dist.Values))
	}
}

func TestBuildSeedData_DefaultStrategy(t *testing.T) {
	input := types.ScaffoldSeedInput{
		Domain: "order",
		Relationships: []types.SeedRelationshipDef{
			{Field: "UserID", Model: "User"}, // No strategy specified
		},
	}

	data := buildSeedData(input, "github.com/example/app")

	if data.Relationships[0].Strategy != "random" {
		t.Errorf("expected default Strategy to be random, got %s", data.Relationships[0].Strategy)
	}
}
