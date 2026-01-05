package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldComponent(t *testing.T) {
	t.Run("validates component name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		tests := []struct {
			name          string
			componentName string
			wantErr       bool
		}{
			{"empty name", "", true},
			{"valid name", "UserCard", false},
			{"valid snake_case", "user_card", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				input := types.ScaffoldComponentInput{
					ComponentName: tt.componentName,
					ComponentType: "card",
				}
				result, err := scaffoldComponent(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr && result.Success {
					t.Error("expected failure for invalid component name")
				}
			})
		}
	})

	t.Run("generates card component", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "ProductCard",
			ComponentType: "card",
			Props: []types.PropDef{
				{Name: "Title", Type: "string", Required: true},
				{Name: "Description", Type: "string"},
			},
		}
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "product_card.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package components") {
			t.Error("expected file to contain 'package components'")
		}
	})

	t.Run("generates modal component", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "ConfirmModal",
			ComponentType: "modal",
		}
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "confirm_modal.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("generates form_field component", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "CustomField",
			ComponentType: "form_field",
		}
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "custom_field.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("generates custom component defaults to card", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "CustomWidget",
			ComponentType: "custom",
		}
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "custom_widget.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("with HTMX enabled", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "DynamicCard",
			ComponentType: "card",
			WithHTMX:      true,
		}
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "dynamic_card.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "TestComponent",
			ComponentType: "card",
			DryRun:        true,
		}
		result, err := scaffoldComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "test_component.templ")
		if fileExists(expectedPath) {
			t.Error("expected file NOT to be created in dry run mode")
		}

		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report the file that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldComponentInput{
			ComponentName: "TestComponent",
			ComponentType: "card",
		}
		result, err := scaffoldComponent(registry, input)
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

func TestGetComponentTemplatePath(t *testing.T) {
	tests := []struct {
		componentType string
		expected      string
	}{
		{"card", "components/card.templ.tmpl"},
		{"modal", "components/card.templ.tmpl"}, // modal falls back to card; use scaffold_modal for full modal support
		{"form_field", "components/form_field.templ.tmpl"},
		{"custom", "components/card.templ.tmpl"},
		{"unknown", "components/card.templ.tmpl"},
		{"", "components/card.templ.tmpl"},
	}

	for _, tt := range tests {
		t.Run(tt.componentType, func(t *testing.T) {
			got := getComponentTemplatePath(tt.componentType)
			if got != tt.expected {
				t.Errorf("getComponentTemplatePath(%q) = %q, want %q", tt.componentType, got, tt.expected)
			}
		})
	}
}

func TestBuildComponentData(t *testing.T) {
	input := types.ScaffoldComponentInput{
		ComponentName: "UserCard",
		ComponentType: "card",
		Props: []types.PropDef{
			{Name: "Name", Type: "string", Required: true},
			{Name: "Age", Type: "int", Default: "0"},
		},
		WithHTMX: true,
		AlpineState: map[string]interface{}{
			"open": false,
		},
	}

	data := buildComponentData(input)

	if data.ComponentName != "UserCard" {
		t.Errorf("expected ComponentName to be UserCard, got %s", data.ComponentName)
	}
	if data.ComponentType != "card" {
		t.Errorf("expected ComponentType to be card, got %s", data.ComponentType)
	}
	if len(data.Props) != 2 {
		t.Errorf("expected 2 props, got %d", len(data.Props))
	}
	if !data.WithHTMX {
		t.Error("expected WithHTMX to be true")
	}
	if data.AlpineState == nil {
		t.Error("expected AlpineState to be set")
	}
}
