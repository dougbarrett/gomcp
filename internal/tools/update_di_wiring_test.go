package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

// Sample main.go content with markers for testing
const mainGoWithMarkers = `package main

import (
	"net/http"
	// MCP:IMPORTS:START
	// MCP:IMPORTS:END
)

func main() {
	db := setupDB()

	// MCP:REPOS:START
	// MCP:REPOS:END

	// MCP:SERVICES:START
	// MCP:SERVICES:END

	// MCP:CONTROLLERS:START
	// MCP:CONTROLLERS:END

	router := http.NewServeMux()

	// MCP:ROUTES:START
	// MCP:ROUTES:END

	http.ListenAndServe(":8080", router)
}
`

const mainGoWithoutMarkers = `package main

import (
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", nil)
}
`

func TestUpdateDIWiring(t *testing.T) {
	t.Run("requires at least one domain", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{},
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when no domains provided")
		}
		if !containsString(result.Message, "at least one domain") {
			t.Errorf("expected error about domains, got: %s", result.Message)
		}
	})

	t.Run("validates domain names", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"valid", "123invalid"},
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure for invalid domain name")
		}
		if !containsString(result.Message, "123invalid") {
			t.Errorf("expected error about invalid domain, got: %s", result.Message)
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		// Don't setup go.mod
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when go.mod is missing")
		}
		if !containsString(result.Message, "module path") {
			t.Errorf("expected error about module path, got: %s", result.Message)
		}
	})

	t.Run("requires main.go at expected path", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		// Don't create main.go

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when main.go is missing")
		}
		if !containsString(result.Message, "main.go not found") {
			t.Errorf("expected error about main.go, got: %s", result.Message)
		}
	})

	t.Run("dry run checks for markers", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithoutMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  true,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when markers are missing")
		}
		if !containsString(result.Message, "missing required markers") {
			t.Errorf("expected error about missing markers, got: %s", result.Message)
		}
	})

	t.Run("dry run succeeds with markers", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  true,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}
		if len(result.FilesUpdated) == 0 {
			t.Errorf("expected FilesUpdated to be populated")
		}
	})

	t.Run("injects imports for domain", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check main.go was updated
		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		// Check imports
		if !containsString(content, "github.com/test/project/internal/repository/product") {
			t.Errorf("expected repository import")
		}
		if !containsString(content, "github.com/test/project/internal/services/product") {
			t.Errorf("expected services import")
		}
		if !containsString(content, "github.com/test/project/internal/web/product") {
			t.Errorf("expected web/controller import")
		}
	})

	t.Run("injects repository instantiation", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		if !containsString(content, "productRepo") {
			t.Errorf("expected productRepo variable")
		}
		if !containsString(content, "product.NewRepository") {
			t.Errorf("expected product.NewRepository call")
		}
	})

	t.Run("injects service instantiation", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		if !containsString(content, "productService") {
			t.Errorf("expected productService variable")
		}
		if !containsString(content, "product.NewService") {
			t.Errorf("expected product.NewService call")
		}
	})

	t.Run("injects controller instantiation", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		if !containsString(content, "productController") {
			t.Errorf("expected productController variable")
		}
		if !containsString(content, "product.NewController") {
			t.Errorf("expected product.NewController call")
		}
	})

	t.Run("injects route registration", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		// Check for full route registration with correct router variable name
		if !containsString(content, "productController.RegisterRoutes(router)") {
			t.Errorf("expected route registration with router variable, got:\n%s", content)
		}
	})

	t.Run("handles multiple domains", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"product", "user", "order"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		// Check all domains have their wiring
		for _, domain := range []string{"product", "user", "order"} {
			if !containsString(content, domain+"Repo") {
				t.Errorf("expected %sRepo variable", domain)
			}
			if !containsString(content, domain+"Service") {
				t.Errorf("expected %sService variable", domain)
			}
			if !containsString(content, domain+"Controller") {
				t.Errorf("expected %sController variable", domain)
			}
		}
	})

	t.Run("handles snake_case domain names", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")
		setupMainGo(t, tmpDir, mainGoWithMarkers)

		input := types.UpdateDIWiringInput{
			Domains: []string{"order_item"},
			DryRun:  false,
		}

		result, err := updateDIWiring(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainGoPath := filepath.Join(tmpDir, "cmd", "web", "main.go")
		content := readFile(t, mainGoPath)

		// Check package name is lowercase without underscore
		if !containsString(content, "orderitem.NewRepository") {
			t.Errorf("expected orderitem package name")
		}
	})
}

// setupMainGo creates main.go with the given content
func setupMainGo(t *testing.T, tmpDir, content string) {
	t.Helper()
	mainDir := filepath.Join(tmpDir, "cmd", "web")
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		t.Fatalf("failed to create cmd/web dir: %v", err)
	}
	mainPath := filepath.Join(mainDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}
}
