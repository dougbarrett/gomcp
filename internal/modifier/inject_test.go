package modifier

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewInjectorFromContent tests creating injector from content string.
func TestNewInjectorFromContent(t *testing.T) {
	content := "package main\n\nfunc main() {}\n"
	injector := NewInjectorFromContent(content)

	if injector == nil {
		t.Fatal("NewInjectorFromContent returned nil")
	}
	if injector.Content() != content {
		t.Errorf("Content() = %q, want %q", injector.Content(), content)
	}
}

// TestNewInjector tests creating injector from file.
func TestNewInjector(t *testing.T) {
	// Create a temp file
	tmpDir, err := os.MkdirTemp("", "inject-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	content := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	injector, err := NewInjector(testFile)
	if err != nil {
		t.Fatalf("NewInjector() error = %v", err)
	}
	if injector.Content() != content {
		t.Errorf("Content() = %q, want %q", injector.Content(), content)
	}
}

// TestNewInjector_FileNotFound tests error on missing file.
func TestNewInjector_FileNotFound(t *testing.T) {
	_, err := NewInjector("/nonexistent/path/file.go")
	if err == nil {
		t.Error("NewInjector() should error for non-existent file")
	}
}

// TestInjector_HasMarker tests marker detection.
func TestInjector_HasMarker(t *testing.T) {
	content := `package main

// MCP:REPOS:START
// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	tests := []struct {
		marker string
		want   bool
	}{
		{MarkerReposStart, true},
		{MarkerReposEnd, true},
		{"MCP:NONEXISTENT", false},
		{MarkerImportsStart, false},
	}

	for _, tt := range tests {
		t.Run(tt.marker, func(t *testing.T) {
			got := injector.HasMarker(tt.marker)
			if got != tt.want {
				t.Errorf("HasMarker(%q) = %v, want %v", tt.marker, got, tt.want)
			}
		})
	}
}

// TestInjector_InjectBetweenMarkers tests basic injection.
func TestInjector_InjectBetweenMarkers(t *testing.T) {
	content := `package main

	// MCP:REPOS:START
	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	code := `userRepo := user.NewRepository(db)`
	err := injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, code)
	if err != nil {
		t.Fatalf("InjectBetweenMarkers() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, "userRepo := user.NewRepository(db)") {
		t.Error("Injected code not found in result")
	}
	if !strings.Contains(result, "// MCP:REPOS:START") {
		t.Error("Start marker should be preserved")
	}
	if !strings.Contains(result, "// MCP:REPOS:END") {
		t.Error("End marker should be preserved")
	}
}

// TestInjector_InjectBetweenMarkers_Duplicate tests that duplicate code is skipped.
func TestInjector_InjectBetweenMarkers_Duplicate(t *testing.T) {
	content := `package main

	// MCP:REPOS:START
	userRepo := user.NewRepository(db)
	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	code := `userRepo := user.NewRepository(db)`
	err := injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, code)
	if err != nil {
		t.Fatalf("InjectBetweenMarkers() error = %v", err)
	}

	// Count occurrences
	result := injector.Content()
	count := strings.Count(result, "userRepo := user.NewRepository(db)")
	if count != 1 {
		t.Errorf("Code should appear once, found %d times", count)
	}
}

// TestInjector_InjectBetweenMarkers_Multiple tests multiple injections.
func TestInjector_InjectBetweenMarkers_Multiple(t *testing.T) {
	content := `package main

	// MCP:REPOS:START
	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, `userRepo := user.NewRepository(db)`)
	injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, `productRepo := product.NewRepository(db)`)
	injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, `orderRepo := order.NewRepository(db)`)

	result := injector.Content()
	if !strings.Contains(result, "userRepo") {
		t.Error("userRepo not found")
	}
	if !strings.Contains(result, "productRepo") {
		t.Error("productRepo not found")
	}
	if !strings.Contains(result, "orderRepo") {
		t.Error("orderRepo not found")
	}
}

// TestInjector_InjectBetweenMarkers_MissingStart tests error on missing start marker.
func TestInjector_InjectBetweenMarkers_MissingStart(t *testing.T) {
	content := `package main

	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, "code")
	if err == nil {
		t.Error("Should error on missing start marker")
	}
	if !strings.Contains(err.Error(), "start marker not found") {
		t.Errorf("Error should mention start marker: %v", err)
	}
}

// TestInjector_InjectBetweenMarkers_MissingEnd tests error on missing end marker.
func TestInjector_InjectBetweenMarkers_MissingEnd(t *testing.T) {
	content := `package main

	// MCP:REPOS:START

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, "code")
	if err == nil {
		t.Error("Should error on missing end marker")
	}
	if !strings.Contains(err.Error(), "end marker not found") {
		t.Errorf("Error should mention end marker: %v", err)
	}
}

// TestInjector_InjectBetweenMarkers_WrongOrder tests error when markers are in wrong order.
func TestInjector_InjectBetweenMarkers_WrongOrder(t *testing.T) {
	content := `package main

	// MCP:REPOS:END
	// MCP:REPOS:START

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, "code")
	if err == nil {
		t.Error("Should error when end comes before start")
	}
	if !strings.Contains(err.Error(), "before end marker") {
		t.Errorf("Error should mention order: %v", err)
	}
}

// TestInjector_InjectAfterMarker tests injection after a marker.
func TestInjector_InjectAfterMarker(t *testing.T) {
	content := `package main

	// MCP:MARKER

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectAfterMarker("MCP:MARKER", "// Injected after")
	if err != nil {
		t.Fatalf("InjectAfterMarker() error = %v", err)
	}

	result := injector.Content()
	markerIdx := strings.Index(result, "// MCP:MARKER")
	injectedIdx := strings.Index(result, "// Injected after")

	if injectedIdx <= markerIdx {
		t.Error("Injected code should appear after marker")
	}
}

// TestInjector_InjectAfterMarker_NotFound tests error on missing marker.
func TestInjector_InjectAfterMarker_NotFound(t *testing.T) {
	content := `package main
func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectAfterMarker("MCP:NONEXISTENT", "code")
	if err == nil {
		t.Error("Should error on missing marker")
	}
}

// TestInjector_InjectBeforeMarker tests injection before a marker.
func TestInjector_InjectBeforeMarker(t *testing.T) {
	content := `package main

	// MCP:MARKER

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectBeforeMarker("MCP:MARKER", "// Injected before")
	if err != nil {
		t.Fatalf("InjectBeforeMarker() error = %v", err)
	}

	result := injector.Content()
	markerIdx := strings.Index(result, "// MCP:MARKER")
	injectedIdx := strings.Index(result, "// Injected before")

	if injectedIdx >= markerIdx {
		t.Error("Injected code should appear before marker")
	}
}

// TestInjector_InjectBeforeMarker_NotFound tests error on missing marker.
func TestInjector_InjectBeforeMarker_NotFound(t *testing.T) {
	content := `package main
func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectBeforeMarker("MCP:NONEXISTENT", "code")
	if err == nil {
		t.Error("Should error on missing marker")
	}
}

// TestInjector_ReplaceMarkerContent tests content replacement.
func TestInjector_ReplaceMarkerContent(t *testing.T) {
	content := `package main

	// MCP:REPOS:START
	oldCode := "should be removed"
	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.ReplaceMarkerContent(MarkerReposStart, MarkerReposEnd, `newCode := "replacement"`)
	if err != nil {
		t.Fatalf("ReplaceMarkerContent() error = %v", err)
	}

	result := injector.Content()
	if strings.Contains(result, "oldCode") {
		t.Error("Old content should be removed")
	}
	if !strings.Contains(result, "newCode") {
		t.Error("New content should be present")
	}
}

// TestInjector_ReplaceMarkerContent_Empty tests replacing with empty content.
func TestInjector_ReplaceMarkerContent_Empty(t *testing.T) {
	content := `package main

	// MCP:REPOS:START
	codeToRemove := true
	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.ReplaceMarkerContent(MarkerReposStart, MarkerReposEnd, "")
	if err != nil {
		t.Fatalf("ReplaceMarkerContent() error = %v", err)
	}

	result := injector.Content()
	if strings.Contains(result, "codeToRemove") {
		t.Error("Content should be removed")
	}
	if !strings.Contains(result, "// MCP:REPOS:START") {
		t.Error("Start marker should be preserved")
	}
	if !strings.Contains(result, "// MCP:REPOS:END") {
		t.Error("End marker should be preserved")
	}
}

// TestInjector_InjectImport tests import injection.
func TestInjector_InjectImport(t *testing.T) {
	content := `package main

import (
	"fmt"
)

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectImport("github.com/example/pkg")
	if err != nil {
		t.Fatalf("InjectImport() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, `"github.com/example/pkg"`) {
		t.Error("Import should be added")
	}
	if !strings.Contains(result, `"fmt"`) {
		t.Error("Existing imports should be preserved")
	}
}

// TestInjector_InjectImport_Duplicate tests that duplicate imports are skipped.
func TestInjector_InjectImport_Duplicate(t *testing.T) {
	content := `package main

import (
	"fmt"
	"github.com/example/pkg"
)

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectImport("github.com/example/pkg")
	if err != nil {
		t.Fatalf("InjectImport() error = %v", err)
	}

	result := injector.Content()
	count := strings.Count(result, `"github.com/example/pkg"`)
	if count != 1 {
		t.Errorf("Import should appear once, found %d times", count)
	}
}

// TestInjector_InjectImport_WithMarkers tests import injection with markers.
func TestInjector_InjectImport_WithMarkers(t *testing.T) {
	content := `package main

import (
	// MCP:IMPORTS:START
	// MCP:IMPORTS:END
	"fmt"
)

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectImport("github.com/example/pkg")
	if err != nil {
		t.Fatalf("InjectImport() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, `"github.com/example/pkg"`) {
		t.Error("Import should be added")
	}
}

// TestInjector_InjectImport_NoImportBlock tests error when no import block.
func TestInjector_InjectImport_NoImportBlock(t *testing.T) {
	content := `package main

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectImport("github.com/example/pkg")
	if err == nil {
		t.Error("Should error when no import block found")
	}
}

// TestInjector_InjectModel tests model injection.
func TestInjector_InjectModel(t *testing.T) {
	content := `package main

	// MCP:MODELS:START
	// MCP:MODELS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectModel("Product")
	if err != nil {
		t.Fatalf("InjectModel() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, "&models.Product{},") {
		t.Error("Model should be injected")
	}
}

// TestInjector_InjectRepo tests repository injection.
func TestInjector_InjectRepo(t *testing.T) {
	content := `package main

	// MCP:REPOS:START
	// MCP:REPOS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectRepo("product", "github.com/example/app")
	if err != nil {
		t.Fatalf("InjectRepo() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, "productRepo := productrepo.NewRepository(db)") {
		t.Error("Repo should be injected")
	}
}

// TestInjector_InjectService tests service injection.
func TestInjector_InjectService(t *testing.T) {
	content := `package main

	// MCP:SERVICES:START
	// MCP:SERVICES:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectService("product")
	if err != nil {
		t.Fatalf("InjectService() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, "productService := productsvc.NewService(productRepo)") {
		t.Error("Service should be injected")
	}
}

// TestInjector_InjectController tests controller injection.
func TestInjector_InjectController(t *testing.T) {
	content := `package main

	// MCP:CONTROLLERS:START
	// MCP:CONTROLLERS:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectController("product")
	if err != nil {
		t.Fatalf("InjectController() error = %v", err)
	}

	result := injector.Content()
	if !strings.Contains(result, "productController := productctrl.NewController(productService)") {
		t.Error("Controller should be injected")
	}
}

// TestInjector_InjectRoute tests route injection.
func TestInjector_InjectRoute(t *testing.T) {
	content := `package main

	// MCP:ROUTES:START
	// MCP:ROUTES:END

func main() {}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectRoute("product")
	if err != nil {
		t.Fatalf("InjectRoute() error = %v", err)
	}

	result := injector.Content()
	expectedRoute := "productController.RegisterRoutes(router)"
	if !strings.Contains(result, expectedRoute) {
		t.Errorf("Route should be injected.\nExpected to contain: %s\nActual content:\n%s", expectedRoute, result)
	}
}

// TestInjector_Save tests saving to file.
func TestInjector_Save(t *testing.T) {
	// Create a temp file
	tmpDir, err := os.MkdirTemp("", "inject-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

	// MCP:REPOS:START
	// MCP:REPOS:END

func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	injector, err := NewInjector(testFile)
	if err != nil {
		t.Fatalf("NewInjector() error = %v", err)
	}

	injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, `userRepo := user.NewRepository(db)`)

	err = injector.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Read back and verify
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(data), "userRepo := user.NewRepository(db)") {
		t.Error("Saved file should contain injected code")
	}
}

// TestInjector_Save_NoFilePath tests error when no file path.
func TestInjector_Save_NoFilePath(t *testing.T) {
	injector := NewInjectorFromContent("content")

	err := injector.Save()
	if err == nil {
		t.Error("Save() should error when no file path set")
	}
}

// TestInjector_SaveTo tests saving to a different file.
func TestInjector_SaveTo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "inject-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	injector := NewInjectorFromContent("test content")

	outputFile := filepath.Join(tmpDir, "output.txt")
	err = injector.SaveTo(outputFile)
	if err != nil {
		t.Fatalf("SaveTo() error = %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(data) != "test content" {
		t.Errorf("SaveTo() wrote %q, want %q", string(data), "test content")
	}
}

// TestIndentCode tests the indentCode helper.
func TestIndentCode(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		indent string
		want   string
	}{
		{
			name:   "single line",
			code:   "x := 1",
			indent: "\t",
			want:   "\tx := 1",
		},
		{
			name:   "multiple lines",
			code:   "x := 1\ny := 2\nz := 3",
			indent: "\t\t",
			want:   "\t\tx := 1\n\t\ty := 2\n\t\tz := 3",
		},
		{
			name:   "with spaces",
			code:   "line1\nline2",
			indent: "    ",
			want:   "    line1\n    line2",
		},
		{
			name:   "empty indent",
			code:   "line1\nline2",
			indent: "",
			want:   "line1\nline2",
		},
		{
			name:   "with empty lines",
			code:   "line1\n\nline2",
			indent: "\t",
			want:   "\tline1\n\n\tline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := indentCode(tt.code, tt.indent)
			if got != tt.want {
				t.Errorf("indentCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestMarkerConstants tests that marker constants are properly defined.
func TestMarkerConstants(t *testing.T) {
	// Verify markers have expected format
	markers := []string{
		MarkerModelsStart, MarkerModelsEnd,
		MarkerReposStart, MarkerReposEnd,
		MarkerServicesStart, MarkerServicesEnd,
		MarkerControllersStart, MarkerControllersEnd,
		MarkerRoutesStart, MarkerRoutesEnd,
		MarkerImportsStart, MarkerImportsEnd,
	}

	for _, marker := range markers {
		if !strings.HasPrefix(marker, "MCP:") {
			t.Errorf("Marker %q should start with 'MCP:'", marker)
		}
		if !strings.HasSuffix(marker, ":START") && !strings.HasSuffix(marker, ":END") {
			t.Errorf("Marker %q should end with :START or :END", marker)
		}
	}
}

// TestInjector_FullDomainInjection tests injecting a complete domain.
func TestInjector_FullDomainInjection(t *testing.T) {
	content := `package main

import (
	// MCP:IMPORTS:START
	// MCP:IMPORTS:END
	"fmt"
)

func main() {
	// MCP:MODELS:START
	// MCP:MODELS:END

	// MCP:REPOS:START
	// MCP:REPOS:END

	// MCP:SERVICES:START
	// MCP:SERVICES:END

	// MCP:CONTROLLERS:START
	// MCP:CONTROLLERS:END

	// MCP:ROUTES:START
	// MCP:ROUTES:END
}
`
	injector := NewInjectorFromContent(content)

	// Inject complete domain
	domain := "product"
	injector.InjectImport("github.com/example/app/internal/repository/product")
	injector.InjectImport("github.com/example/app/internal/services/product")
	injector.InjectImport("github.com/example/app/internal/web/product")
	injector.InjectModel("Product")
	injector.InjectRepo(domain, "github.com/example/app")
	injector.InjectService(domain)
	injector.InjectController(domain)
	injector.InjectRoute(domain)

	result := injector.Content()

	// Verify all injections
	checks := []string{
		`"github.com/example/app/internal/repository/product"`,
		`"github.com/example/app/internal/services/product"`,
		`"github.com/example/app/internal/web/product"`,
		"&models.Product{},",
		"productRepo := productrepo.NewRepository(db)",
		"productService := productsvc.NewService(productRepo)",
		"productController := productctrl.NewController(productService)",
		"productController.RegisterRoutes(router)",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Result should contain %q", check)
		}
	}
}

// TestInjector_PreservesIndentation tests that indentation is preserved.
func TestInjector_PreservesIndentation(t *testing.T) {
	content := `package main

func main() {
		// MCP:REPOS:START
		// MCP:REPOS:END
}
`
	injector := NewInjectorFromContent(content)

	err := injector.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, `userRepo := user.NewRepository(db)`)
	if err != nil {
		t.Fatalf("InjectBetweenMarkers() error = %v", err)
	}

	result := injector.Content()
	// The injected code should have the same indentation as the end marker (2 tabs)
	if !strings.Contains(result, "\t\tuserRepo := user.NewRepository(db)") {
		t.Error("Injected code should preserve indentation")
	}
}
