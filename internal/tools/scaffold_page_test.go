package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldPage(t *testing.T) {
	t.Run("validates page name", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldPageInput{
			PageName: "",
			Route:    "/test",
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty page name")
		}
		if !strings.Contains(strings.ToLower(result.Message), "page name") {
			t.Errorf("expected error about page name, got %q", result.Message)
		}
	})

	t.Run("validates route", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldPageInput{
			PageName: "TestPage",
			Route:    "",
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for empty route")
		}
		if !strings.Contains(strings.ToLower(result.Message), "route") {
			t.Errorf("expected error about route, got %q", result.Message)
		}
	})

	t.Run("validates route format", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldPageInput{
			PageName: "TestPage",
			Route:    "invalid-route", // missing leading slash
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid route format")
		}
	})

	t.Run("validates layout type", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldPageInput{
			PageName: "TestPage",
			Route:    "/test",
			Layout:   "invalid_layout",
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Error("expected failure for invalid layout type")
		}
	})

	t.Run("accepts valid layout types", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		validLayouts := []string{"default", "dashboard", "landing", "blank", ""}

		for _, layout := range validLayouts {
			t.Run(layout, func(t *testing.T) {
				pageName := "test_page_" + layout
				if layout == "" {
					pageName = "test_page_empty"
				}
				input := types.ScaffoldPageInput{
					PageName: pageName,
					Route:    "/" + pageName,
					Layout:   layout,
				}
				result, err := scaffoldPage(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !result.Success {
					t.Errorf("expected success for layout %q, got: %s", layout, result.Message)
				}
			})
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.ScaffoldPageInput{
			PageName: "TestPage",
			Route:    "/test",
		}
		result, err := scaffoldPage(registry, input)
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

	t.Run("generates page", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldPageInput{
			PageName: "Dashboard",
			Route:    "/dashboard",
			Layout:   "dashboard",
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		expectedPath := filepath.Join(tmpDir, "internal", "web", "pages", "dashboard.templ")
		if !fileExists(expectedPath) {
			t.Errorf("expected file to be created at %s", expectedPath)
		}
	})

	t.Run("generates page with TOML config", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldPageInput{
			PageName:         "Settings",
			Route:            "/settings",
			Layout:           "default",
			CreateTomlConfig: true,
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check page file
		pagePath := filepath.Join(tmpDir, "internal", "web", "pages", "settings.templ")
		if !fileExists(pagePath) {
			t.Errorf("expected page file to be created at %s", pagePath)
		}

		// Check config file
		configPath := filepath.Join(tmpDir, "config", "en", "pages", "settings.toml")
		if !fileExists(configPath) {
			t.Errorf("expected config file to be created at %s", configPath)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldPageInput{
			PageName:         "TestPage",
			Route:            "/test",
			CreateTomlConfig: true,
			DryRun:           true,
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		pagePath := filepath.Join(tmpDir, "internal", "web", "pages", "test_page.templ")
		if fileExists(pagePath) {
			t.Error("expected page file NOT to be created in dry run mode")
		}

		configPath := filepath.Join(tmpDir, "config", "en", "pages", "test_page.toml")
		if fileExists(configPath) {
			t.Error("expected config file NOT to be created in dry run mode")
		}

		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report files that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/example/testapp")

		input := types.ScaffoldPageInput{
			PageName: "TestPage",
			Route:    "/test",
		}
		result, err := scaffoldPage(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}
	})
}
