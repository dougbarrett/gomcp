package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldModal(t *testing.T) {
	t.Run("validates modal name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldModalInput{
			ModalName: "",
			ModalType: "dialog",
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty modal name")
		}
		if !strings.Contains(strings.ToLower(result.Message), "modal name") {
			t.Errorf("expected error about modal name, got %q", result.Message)
		}
	})

	t.Run("validates modal type", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldModalInput{
			ModalName: "test_modal",
			ModalType: "invalid_type",
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid modal type")
		}
		if !strings.Contains(strings.ToLower(result.Message), "modal type") {
			t.Errorf("expected error about modal type, got %q", result.Message)
		}
	})

	t.Run("accepts valid modal types", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		validTypes := []string{"dialog", "sheet", "confirm"}

		for _, modalType := range validTypes {
			t.Run(modalType, func(t *testing.T) {
				input := types.ScaffoldModalInput{
					ModalName: "test_" + modalType,
					ModalType: modalType,
					Title:     "Test Modal",
				}
				result, err := scaffoldModal(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !result.Success {
					t.Errorf("expected success for modal type %q, got: %s", modalType, result.Message)
				}

				expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "test_"+modalType+".templ")
				if !fileExists(expectedPath) {
					t.Errorf("expected file to be created at %s", expectedPath)
				}
			})
		}
	})

	t.Run("generates dialog modal", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldModalInput{
			ModalName:   "product_dialog",
			ModalType:   "dialog",
			Title:       "Product Details",
			ContentType: "info",
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "product_dialog.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}

		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package components") {
			t.Error("expected file to contain 'package components'")
		}
	})

	t.Run("generates sheet modal", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldModalInput{
			ModalName:   "filter_sheet",
			ModalType:   "sheet",
			Title:       "Filters",
			ContentType: "form",
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "filter_sheet.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("generates confirm modal", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldModalInput{
			ModalName:   "delete_confirm",
			ModalType:   "confirm",
			Title:       "Confirm Delete",
			ContentType: "confirm",
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "delete_confirm.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("includes HTMX URL in next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldModalInput{
			ModalName: "dynamic_modal",
			ModalType: "dialog",
			TriggerConfig: types.TriggerConfig{
				HTMXURL: "/api/content",
			},
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		hasHTMXStep := false
		for _, step := range result.NextSteps {
			if strings.Contains(step, "/api/content") {
				hasHTMXStep = true
				break
			}
		}
		if !hasHTMXStep {
			t.Error("expected NextSteps to include HTMX URL endpoint")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldModalInput{
			ModalName: "test_modal",
			ModalType: "dialog",
			DryRun:    true,
		}
		result, err := scaffoldModal(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "components", "test_modal.templ")
		if fileExists(expectedPath) {
			t.Error("expected file NOT to be created in dry run mode")
		}

		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report the file that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldModalInput{
			ModalName: "test_modal",
			ModalType: "dialog",
		}
		result, err := scaffoldModal(registry, input)
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
