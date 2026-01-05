package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldTable(t *testing.T) {
	t.Run("validates table name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldTableInput{
			TableName: "",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name"},
			},
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty table name")
		}
		if !strings.Contains(strings.ToLower(result.Message), "table name") {
			t.Errorf("expected error about table name, got %q", result.Message)
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				input := types.ScaffoldTableInput{
					TableName: "test_table",
					Domain:    tt.domainName,
					Columns: []types.ColumnDef{
						{Key: "name", Label: "Name"},
					},
				}
				result, err := scaffoldTable(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Error("expected failure for invalid domain name")
				}
			})
		}
	})

	t.Run("requires at least one column", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldTableInput{
			TableName: "test_table",
			Domain:    "product",
			Columns:   []types.ColumnDef{},
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty columns")
		}
		if !strings.Contains(strings.ToLower(result.Message), "column") {
			t.Errorf("expected error about columns, got %q", result.Message)
		}
	})

	t.Run("validates column key is required", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldTableInput{
			TableName: "test_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "", Label: "Name"},
			},
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty column key")
		}
		if !strings.Contains(strings.ToLower(result.Message), "column key") {
			t.Errorf("expected error about column key, got %q", result.Message)
		}
	})

	t.Run("validates row action types", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldTableInput{
			TableName: "test_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name"},
			},
			RowActions: []types.RowActionDef{
				{Type: "invalid_action"},
			},
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid row action type")
		}
		if !strings.Contains(strings.ToLower(result.Message), "row action type") {
			t.Errorf("expected error about row action type, got %q", result.Message)
		}
	})

	t.Run("accepts valid row action types", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		validTypes := []string{"view", "edit", "delete", "custom"}

		for _, actionType := range validTypes {
			t.Run(actionType, func(t *testing.T) {
				input := types.ScaffoldTableInput{
					TableName: "test_table_" + actionType,
					Domain:    "product",
					Columns: []types.ColumnDef{
						{Key: "name", Label: "Name"},
					},
					RowActions: []types.RowActionDef{
						{Type: actionType, Label: "Action"},
					},
				}
				result, err := scaffoldTable(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !result.Success {
					t.Errorf("expected success for action type %q, got: %s", actionType, result.Message)
				}
			})
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldTableInput{
			TableName: "product_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name"},
			},
		}
		result, err := scaffoldTable(registry, input)
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

	t.Run("generates table with all options", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		withPagination := true
		withSorting := true
		withSearch := true

		input := types.ScaffoldTableInput{
			TableName: "product_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name", Sortable: true},
				{Key: "price", Label: "Price", Sortable: true, Format: "currency"},
				{Key: "created_at", Label: "Created", Sortable: true, Format: "date"},
			},
			WithPagination:  &withPagination,
			WithSorting:     &withSorting,
			WithSearch:      &withSearch,
			WithBulkActions: true,
			RowActions: []types.RowActionDef{
				{Type: "view", Label: "View"},
				{Type: "edit", Label: "Edit"},
				{Type: "delete", Label: "Delete", Confirm: true, ConfirmMessage: "Are you sure?"},
			},
		}
		result, err := scaffoldTable(registry, input)
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

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package views") {
			t.Error("expected file to contain 'package views'")
		}
	})

	t.Run("defaults for pagination, sorting, search", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		// When nil, defaults should be true
		input := types.ScaffoldTableInput{
			TableName: "default_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name"},
			},
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldTableInput{
			TableName: "product_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name"},
			},
			DryRun: true,
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "product", "views", "product_table.templ")
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

		input := types.ScaffoldTableInput{
			TableName: "product_table",
			Domain:    "product",
			Columns: []types.ColumnDef{
				{Key: "name", Label: "Name"},
			},
		}
		result, err := scaffoldTable(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}

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

func TestBuildTableData(t *testing.T) {
	withPagination := true
	withSorting := false
	withSearch := true

	input := types.ScaffoldTableInput{
		TableName: "product_list",
		Domain:    "product_category",
		Columns: []types.ColumnDef{
			{Key: "name", Label: "Name", Sortable: true},
			{Key: "price", Label: "Price"},
		},
		WithPagination:  &withPagination,
		WithSorting:     &withSorting,
		WithSearch:      &withSearch,
		WithBulkActions: true,
		RowActions: []types.RowActionDef{
			{Type: "edit", Label: "Edit"},
		},
	}

	data := buildTableData(input, "github.com/example/app", "product_category")

	if data.ModulePath != "github.com/example/app" {
		t.Errorf("expected ModulePath to be github.com/example/app, got %s", data.ModulePath)
	}
	if data.DomainName != "product_category" {
		t.Errorf("expected DomainName to be product_category, got %s", data.DomainName)
	}
	if data.ModelName != "ProductCategory" {
		t.Errorf("expected ModelName to be ProductCategory, got %s", data.ModelName)
	}
	if data.TableName != "product_list" {
		t.Errorf("expected TableName to be product_list, got %s", data.TableName)
	}
	if !data.WithPagination {
		t.Error("expected WithPagination to be true")
	}
	if data.WithSorting {
		t.Error("expected WithSorting to be false")
	}
	if !data.WithSearch {
		t.Error("expected WithSearch to be true")
	}
	if !data.WithBulkActions {
		t.Error("expected WithBulkActions to be true")
	}
	if len(data.Columns) != 2 {
		t.Errorf("expected 2 columns, got %d", len(data.Columns))
	}
	if len(data.RowActions) != 1 {
		t.Errorf("expected 1 row action, got %d", len(data.RowActions))
	}
}
