package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

// Template validation tests verify that generated code is syntactically valid.
// These tests require the Go toolchain to be available.

func TestTemplateValidation_ProjectTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping template validation in short mode")
	}

	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go command not available")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "validationtest",
		ModulePath:   "github.com/test/validationtest",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	result, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !result.Success {
		t.Fatalf("scaffoldProject failed: %s", result.Message)
	}

	projectDir := filepath.Join(tmpDir, "validationtest")

	// Validate all .go files have valid Go syntax
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Use go fmt to check syntax (it will fail on syntax errors)
		cmd := exec.Command("gofmt", "-e", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("syntax error in %s: %s", path, string(output))
		}

		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk project directory: %v", err)
	}
}

func TestTemplateValidation_DomainTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping template validation in short mode")
	}

	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go command not available")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project first
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "domainvalidation",
		ModulePath:   "github.com/test/domainvalidation",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil || !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %v / %s", err, projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "domainvalidation")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold a domain with various field types
	domainInput := types.ScaffoldDomainInput{
		DomainName: "article",
		Fields: []types.FieldDef{
			{Name: "Title", Type: "string", Required: true},
			{Name: "Content", Type: "string"},
			{Name: "Views", Type: "int"},
			{Name: "Rating", Type: "float64"},
			{Name: "Published", Type: "bool"},
			{Name: "PublishedAt", Type: "time.Time"},
		},
		DryRun: false,
	}

	domainResult, err := scaffoldDomain(projectRegistry, domainInput)
	if err != nil || !domainResult.Success {
		t.Fatalf("scaffoldDomain failed: %v / %s", err, domainResult.Message)
	}

	// Validate all generated .go files
	goFiles := []string{
		"internal/models/article.go",
		"internal/repository/article/article.go",
		"internal/services/article/article.go",
		"internal/services/article/dto.go",
		"internal/web/article/article.go",
	}

	for _, file := range goFiles {
		path := filepath.Join(projectDir, file)
		if !fileExists(path) {
			t.Errorf("expected file %s to exist", file)
			continue
		}

		cmd := exec.Command("gofmt", "-e", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("syntax error in %s: %s", file, string(output))
		}
	}
}

func TestTemplateValidation_ViewTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping template validation in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "viewvalidation",
		ModulePath:   "github.com/test/viewvalidation",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil || !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %v / %s", err, projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "viewvalidation")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold views
	viewTypes := []string{"list", "show", "form", "table"}
	for _, viewType := range viewTypes {
		viewInput := types.ScaffoldViewInput{
			DomainName: "post",
			ViewType:   viewType,
			ViewName:   viewType,
			DryRun:     false,
		}

		viewResult, err := scaffoldView(projectRegistry, viewInput)
		if err != nil || !viewResult.Success {
			t.Errorf("scaffoldView(%s) failed: %v / %s", viewType, err, viewResult.Message)
			continue
		}

		// Check the file exists and has valid templ syntax (basic check)
		viewPath := filepath.Join(projectDir, "internal/web/post/views", viewType+".templ")
		if !fileExists(viewPath) {
			t.Errorf("expected view file %s.templ to exist", viewType)
			continue
		}

		// Read the file and check for basic templ structure
		content := readFile(t, viewPath)
		if !strings.Contains(content, "package views") {
			t.Errorf("view %s.templ should have package declaration", viewType)
		}
		if !strings.Contains(content, "templ ") {
			t.Errorf("view %s.templ should have templ function", viewType)
		}
	}
}

func TestTemplateValidation_SeedTemplates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping template validation in short mode")
	}

	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go command not available")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "seedvalidation",
		ModulePath:   "github.com/test/seedvalidation",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil || !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %v / %s", err, projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "seedvalidation")
	projectRegistry := NewRegistry(projectDir)

	// First create a domain
	domainInput := types.ScaffoldDomainInput{
		DomainName: "user",
		Fields: []types.FieldDef{
			{Name: "Email", Type: "string", Required: true},
			{Name: "Name", Type: "string"},
		},
		DryRun: false,
	}

	domainResult, err := scaffoldDomain(projectRegistry, domainInput)
	if err != nil || !domainResult.Success {
		t.Fatalf("scaffoldDomain failed: %v / %s", err, domainResult.Message)
	}

	// Scaffold a seeder
	seedInput := types.ScaffoldSeedInput{
		Domain: "user",
		Count:  10,
		Fields: []types.FieldDef{
			{Name: "Email", Type: "string"},
			{Name: "Name", Type: "string"},
		},
		WithFaker: true,
		DryRun:    false,
	}

	seedResult, err := scaffoldSeed(projectRegistry, seedInput)
	if err != nil || !seedResult.Success {
		t.Fatalf("scaffoldSeed failed: %v / %s", err, seedResult.Message)
	}

	// Validate seeder file syntax
	seederPath := filepath.Join(projectDir, "cmd/seed/seeders/user_seeder.go")
	if !fileExists(seederPath) {
		t.Fatalf("expected seeder file to exist")
	}

	cmd := exec.Command("gofmt", "-e", seederPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("syntax error in seeder: %s", string(output))
	}
}

func TestTemplateValidation_ConfigTemplates(t *testing.T) {
	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "configvalidation",
		ModulePath:   "github.com/test/configvalidation",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil || !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %v / %s", err, projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "configvalidation")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold configs
	// Note: "menu" and "app" configs are already created by scaffoldProject,
	// so we only scaffold "page" and "messages" here to avoid conflicts.
	configTypes := []string{"page", "messages"}
	for _, configType := range configTypes {
		configInput := types.ScaffoldConfigInput{
			ConfigType: configType,
			Name:       "test",
			Locale:     "en",
			DryRun:     false,
		}

		configResult, err := scaffoldConfig(projectRegistry, configInput)
		if err != nil || !configResult.Success {
			t.Errorf("scaffoldConfig(%s) failed: %v / %s", configType, err, configResult.Message)
			continue
		}
	}

	// Validate TOML files have valid syntax (basic check)
	tomlFiles := []string{
		"config/en/pages/test.toml",
		"config/en/menu.toml",
		"config/en/app.toml",
		"config/en/messages/test.toml",
	}

	for _, file := range tomlFiles {
		path := filepath.Join(projectDir, file)
		if !fileExists(path) {
			t.Errorf("expected config file %s to exist", file)
			continue
		}

		// Basic TOML syntax check - should have at least one section or key
		content := readFile(t, path)
		if len(content) == 0 {
			t.Errorf("config file %s is empty", file)
		}
		if !strings.Contains(content, "[") && !strings.Contains(content, "=") {
			t.Errorf("config file %s doesn't look like valid TOML", file)
		}
	}
}

func TestTemplateValidation_AuthAndUserManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping template validation in short mode")
	}

	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go command not available")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project with auth and user management
	projectInput := types.ScaffoldProjectInput{
		ProjectName:        "authvalidation",
		ModulePath:         "github.com/test/authvalidation",
		DatabaseType:       "sqlite",
		WithAuth:           true,
		WithUserManagement: true,
		DryRun:             false,
	}

	result, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !result.Success {
		t.Fatalf("scaffoldProject failed: %s", result.Message)
	}

	projectDir := filepath.Join(tmpDir, "authvalidation")

	// Verify auth files exist
	authFiles := []string{
		"internal/models/user.go",
		"internal/models/role.go",
		"internal/repository/user/user.go",
		"internal/services/auth/auth.go",
		"internal/services/auth/session.go",
		"internal/web/auth/auth.go",
		"internal/web/middleware/auth.go",
	}

	for _, file := range authFiles {
		path := filepath.Join(projectDir, file)
		if !fileExists(path) {
			t.Errorf("expected auth file %s to exist", file)
		}
	}

	// Verify user management files exist
	userMgmtFiles := []string{
		"internal/services/user/user.go",
		"internal/web/users/users.go",
		"internal/web/users/views/list.templ",
		"internal/web/users/views/form.templ",
		"internal/web/users/views/show.templ",
		"internal/web/users/views/password.templ",
	}

	for _, file := range userMgmtFiles {
		path := filepath.Join(projectDir, file)
		if !fileExists(path) {
			t.Errorf("expected user management file %s to exist", file)
		}
	}

	// Validate all .go files have valid syntax
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		cmd := exec.Command("gofmt", "-e", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("syntax error in %s: %s", path, string(output))
		}

		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk project directory: %v", err)
	}

	// Verify base_layout.templ has admin section
	baseLayoutPath := filepath.Join(projectDir, "internal/web/layouts/base.templ")
	if fileExists(baseLayoutPath) {
		content := readFile(t, baseLayoutPath)
		if !strings.Contains(content, "middleware.IsAdmin") {
			t.Error("base layout should contain admin visibility check")
		}
		if !strings.Contains(content, "MCP:NAV_ITEMS_ADMIN") {
			t.Error("base layout should contain admin nav markers")
		}
		if !strings.Contains(content, "/admin/users") {
			t.Error("base layout should contain Users nav item for user management")
		}
	}
}
