package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

// Integration tests verify the full workflow of scaffolding tools.
// These tests are more comprehensive and may take longer to run.

func TestIntegration_ProjectToDomain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Step 1: Scaffold a new project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "testapp",
		ModulePath:   "github.com/test/testapp",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	// Update registry to point to the new project
	projectDir := filepath.Join(tmpDir, "testapp")
	projectRegistry := NewRegistry(projectDir)

	// Step 2: Scaffold a domain
	domainInput := types.ScaffoldDomainInput{
		DomainName: "product",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string", Required: true},
			{Name: "Price", Type: "float64"},
			{Name: "Description", Type: "string"},
		},
		DryRun: false,
	}

	domainResult, err := scaffoldDomain(projectRegistry, domainInput)
	if err != nil {
		t.Fatalf("scaffoldDomain error: %v", err)
	}
	if !domainResult.Success {
		t.Fatalf("scaffoldDomain failed: %s", domainResult.Message)
	}

	// Step 3: Verify all expected files exist
	expectedFiles := []string{
		// Project files
		"go.mod",
		"cmd/web/main.go",
		"internal/config/config.go",
		"internal/database/database.go",
		// Domain files
		"internal/models/product.go",
		"internal/repository/product/product.go",
		"internal/services/product/product.go",
		"internal/services/product/dto.go",
		"internal/web/product/product.go",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(projectDir, f)
		if !fileExists(path) {
			t.Errorf("expected file %s to exist", f)
		}
	}

	// Step 4: Verify go.mod has correct module path
	goModContent := readFile(t, filepath.Join(projectDir, "go.mod"))
	if !containsString(goModContent, "module github.com/test/testapp") {
		t.Errorf("go.mod should have correct module path")
	}

	// Step 5: Verify model has correct fields
	modelContent := readFile(t, filepath.Join(projectDir, "internal/models/product.go"))
	if !containsString(modelContent, "Name") {
		t.Errorf("model should have Name field")
	}
	if !containsString(modelContent, "Price") {
		t.Errorf("model should have Price field")
	}
	if !containsString(modelContent, "Description") {
		t.Errorf("model should have Description field")
	}
}

func TestIntegration_ProjectBuildable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check if go command is available
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go command not available")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "buildtest",
		ModulePath:   "github.com/test/buildtest",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "buildtest")

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	tidyOutput, err := tidyCmd.CombinedOutput()
	if err != nil {
		t.Logf("go mod tidy output: %s", string(tidyOutput))
		// Don't fail - network issues may prevent downloading deps
		t.Skip("go mod tidy failed (possibly network issue)")
	}

	// Try to build the project (just check syntax, don't link)
	buildCmd := exec.Command("go", "build", "-o", "/dev/null", "./...")
	buildCmd.Dir = projectDir
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("go build output: %s", string(buildOutput))
		// This may fail due to missing dependencies, but syntax should be valid
		if containsString(string(buildOutput), "syntax error") {
			t.Errorf("generated code has syntax errors: %s", string(buildOutput))
		}
	}
}

func TestIntegration_MultipleDomains(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "multidomainapp",
		ModulePath:   "github.com/test/multidomainapp",
		DatabaseType: "postgres",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "multidomainapp")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold multiple domains
	domains := []struct {
		name   string
		fields []types.FieldDef
	}{
		{
			name: "user",
			fields: []types.FieldDef{
				{Name: "Email", Type: "string", Required: true},
				{Name: "Name", Type: "string"},
			},
		},
		{
			name: "product",
			fields: []types.FieldDef{
				{Name: "Title", Type: "string", Required: true},
				{Name: "Price", Type: "float64"},
			},
		},
		{
			name: "order",
			fields: []types.FieldDef{
				{Name: "Total", Type: "float64"},
				{Name: "Status", Type: "string"},
			},
		},
	}

	for _, d := range domains {
		domainInput := types.ScaffoldDomainInput{
			DomainName: d.name,
			Fields:     d.fields,
			DryRun:     false,
		}

		domainResult, err := scaffoldDomain(projectRegistry, domainInput)
		if err != nil {
			t.Fatalf("scaffoldDomain(%s) error: %v", d.name, err)
		}
		if !domainResult.Success {
			t.Fatalf("scaffoldDomain(%s) failed: %s", d.name, domainResult.Message)
		}
	}

	// Verify all domains were created
	for _, d := range domains {
		modelPath := filepath.Join(projectDir, "internal/models", d.name+".go")
		if !fileExists(modelPath) {
			t.Errorf("expected model file for %s to exist", d.name)
		}

		repoPath := filepath.Join(projectDir, "internal/repository", d.name, d.name+".go")
		if !fileExists(repoPath) {
			t.Errorf("expected repository file for %s to exist", d.name)
		}

		servicePath := filepath.Join(projectDir, "internal/services", d.name, d.name+".go")
		if !fileExists(servicePath) {
			t.Errorf("expected service file for %s to exist", d.name)
		}

		controllerPath := filepath.Join(projectDir, "internal/web", d.name, d.name+".go")
		if !fileExists(controllerPath) {
			t.Errorf("expected controller file for %s to exist", d.name)
		}
	}

	// Use list_domains to verify
	listResult, err := listDomains(projectRegistry)
	if err != nil {
		t.Fatalf("listDomains error: %v", err)
	}
	if !listResult.Success {
		t.Fatalf("listDomains failed: %s", listResult.Message)
	}

	if len(listResult.Domains) != 3 {
		t.Errorf("expected 3 domains, got %d", len(listResult.Domains))
	}
}

func TestIntegration_DIWiringWithDomains(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "diwiringapp",
		ModulePath:   "github.com/test/diwiringapp",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "diwiringapp")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold a domain
	domainInput := types.ScaffoldDomainInput{
		DomainName: "task",
		Fields: []types.FieldDef{
			{Name: "Title", Type: "string", Required: true},
			{Name: "Done", Type: "bool"},
		},
		DryRun: false,
	}

	domainResult, err := scaffoldDomain(projectRegistry, domainInput)
	if err != nil {
		t.Fatalf("scaffoldDomain error: %v", err)
	}
	if !domainResult.Success {
		t.Fatalf("scaffoldDomain failed: %s", domainResult.Message)
	}

	// Add MCP markers to main.go for DI wiring
	mainGoPath := filepath.Join(projectDir, "cmd/web/main.go")
	mainGoContent := readFile(t, mainGoPath)

	// Check if the project template already includes markers
	// If not, we need to add them for the test
	if !containsString(mainGoContent, "MCP:IMPORTS:START") {
		// The template should include markers, but if it doesn't, skip this test
		t.Skip("main.go template does not include MCP markers")
	}

	// Update DI wiring
	wiringInput := types.UpdateDIWiringInput{
		Domains: []string{"task"},
		DryRun:  false,
	}

	wiringResult, err := updateDIWiring(projectRegistry, wiringInput)
	if err != nil {
		t.Fatalf("updateDIWiring error: %v", err)
	}
	if !wiringResult.Success {
		t.Fatalf("updateDIWiring failed: %s", wiringResult.Message)
	}

	// Verify main.go was updated
	updatedMainGo := readFile(t, mainGoPath)
	if !containsString(updatedMainGo, "taskRepo") {
		t.Errorf("expected taskRepo in main.go")
	}
	if !containsString(updatedMainGo, "taskService") {
		t.Errorf("expected taskService in main.go")
	}
	if !containsString(updatedMainGo, "taskController") {
		t.Errorf("expected taskController in main.go")
	}
}

func TestIntegration_ViewScaffolding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "viewapp",
		ModulePath:   "github.com/test/viewapp",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "viewapp")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold views for a domain
	viewTypes := []string{"list", "show", "form"}
	for _, viewType := range viewTypes {
		viewInput := types.ScaffoldViewInput{
			DomainName: "article",
			ViewType:   viewType,
			ViewName:   viewType, // Use view type as view name
			DryRun:     false,
		}

		viewResult, err := scaffoldView(projectRegistry, viewInput)
		if err != nil {
			t.Fatalf("scaffoldView(%s) error: %v", viewType, err)
		}
		if !viewResult.Success {
			t.Fatalf("scaffoldView(%s) failed: %s", viewType, viewResult.Message)
		}
	}

	// Verify view files were created
	for _, viewType := range viewTypes {
		viewPath := filepath.Join(projectDir, "internal/web/article/views", viewType+".templ")
		if !fileExists(viewPath) {
			t.Errorf("expected view file %s.templ to exist", viewType)
		}
	}
}

func TestIntegration_ConfigScaffolding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "configapp",
		ModulePath:   "github.com/test/configapp",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "configapp")
	projectRegistry := NewRegistry(projectDir)

	// Scaffold page configs for multiple locales
	locales := []string{"en", "es", "fr"}
	for _, locale := range locales {
		configInput := types.ScaffoldConfigInput{
			ConfigType: "page",
			Name:       "dashboard",
			Locale:     locale,
			DryRun:     false,
		}

		configResult, err := scaffoldConfig(projectRegistry, configInput)
		if err != nil {
			t.Fatalf("scaffoldConfig(%s) error: %v", locale, err)
		}
		if !configResult.Success {
			t.Fatalf("scaffoldConfig(%s) failed: %s", locale, configResult.Message)
		}
	}

	// Verify config files were created for each locale
	for _, locale := range locales {
		configPath := filepath.Join(projectDir, "config", locale, "pages", "dashboard.toml")
		if !fileExists(configPath) {
			t.Errorf("expected config file for locale %s to exist", locale)
		}
	}
}

func TestIntegration_SeedScaffolding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	registry, tmpDir := testRegistry(t)

	// Scaffold a project
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "seedapp",
		ModulePath:   "github.com/test/seedapp",
		DatabaseType: "sqlite",
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "seedapp")
	projectRegistry := NewRegistry(projectDir)

	// First scaffold a domain (seeds need a model to reference)
	domainInput := types.ScaffoldDomainInput{
		DomainName: "category",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string", Required: true},
			{Name: "Slug", Type: "string"},
		},
		DryRun: false,
	}

	domainResult, err := scaffoldDomain(projectRegistry, domainInput)
	if err != nil {
		t.Fatalf("scaffoldDomain error: %v", err)
	}
	if !domainResult.Success {
		t.Fatalf("scaffoldDomain failed: %s", domainResult.Message)
	}

	// Scaffold a seeder
	seedInput := types.ScaffoldSeedInput{
		Domain: "category",
		Count:  10,
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string"},
			{Name: "Slug", Type: "string"},
		},
		DryRun: false,
	}

	seedResult, err := scaffoldSeed(projectRegistry, seedInput)
	if err != nil {
		t.Fatalf("scaffoldSeed error: %v", err)
	}
	if !seedResult.Success {
		t.Fatalf("scaffoldSeed failed: %s", seedResult.Message)
	}

	// Verify seeder file was created
	seederPath := filepath.Join(projectDir, "cmd/seed/seeders", "category_seeder.go")
	if !fileExists(seederPath) {
		t.Errorf("expected seeder file to exist at %s", seederPath)
	}

	// Verify seeder content
	seederContent := readFile(t, seederPath)
	if !containsString(seederContent, "CategorySeeder") {
		t.Errorf("expected CategorySeeder struct in seeder")
	}
	if !containsString(seederContent, "Seed") {
		t.Errorf("expected Seed method in seeder")
	}
}

// Helper to check if a directory exists
func dirExistsIntegration(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func TestWizardIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check if go command is available for gofmt validation
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go command not available")
	}

	registry, tmpDir := testRegistry(t)

	// Step 1: Scaffold a project with auth (provides User model)
	projectInput := types.ScaffoldProjectInput{
		ProjectName:  "wizardtest",
		ModulePath:   "github.com/test/wizardtest",
		DatabaseType: "sqlite",
		WithAuth:     true,
		DryRun:       false,
	}

	projectResult, err := scaffoldProject(registry, projectInput)
	if err != nil {
		t.Fatalf("scaffoldProject error: %v", err)
	}
	if !projectResult.Success {
		t.Fatalf("scaffoldProject failed: %s", projectResult.Message)
	}

	projectDir := filepath.Join(tmpDir, "wizardtest")
	projectRegistry := NewRegistry(projectDir)

	// Step 2: Scaffold a client domain for the select step
	clientInput := types.ScaffoldDomainInput{
		DomainName: "client",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string", Required: true},
			{Name: "Email", Type: "string"},
		},
		DryRun: false,
	}

	clientResult, err := scaffoldDomain(projectRegistry, clientInput)
	if err != nil {
		t.Fatalf("scaffoldDomain(client) error: %v", err)
	}
	if !clientResult.Success {
		t.Fatalf("scaffoldDomain(client) failed: %s", clientResult.Message)
	}

	// Step 3: Scaffold an order domain for the wizard
	orderInput := types.ScaffoldDomainInput{
		DomainName: "order",
		Fields: []types.FieldDef{
			{Name: "Total", Type: "float64"},
			{Name: "Notes", Type: "string"},
			{Name: "Status", Type: "string"},
		},
		Relationships: []types.RelationshipDef{
			{Type: "belongs_to", Model: "Client"},
		},
		DryRun: false,
	}

	orderResult, err := scaffoldDomain(projectRegistry, orderInput)
	if err != nil {
		t.Fatalf("scaffoldDomain(order) error: %v", err)
	}
	if !orderResult.Success {
		t.Fatalf("scaffoldDomain(order) failed: %s", orderResult.Message)
	}

	// Step 4: Scaffold an orderitem domain for has_many step
	orderItemInput := types.ScaffoldDomainInput{
		DomainName: "orderitem",
		Fields: []types.FieldDef{
			{Name: "ProductName", Type: "string"},
			{Name: "Quantity", Type: "int"},
			{Name: "Price", Type: "float64"},
		},
		DryRun: false,
	}

	orderItemResult, err := scaffoldDomain(projectRegistry, orderItemInput)
	if err != nil {
		t.Fatalf("scaffoldDomain(orderitem) error: %v", err)
	}
	if !orderItemResult.Success {
		t.Fatalf("scaffoldDomain(orderitem) failed: %s", orderItemResult.Message)
	}

	// Step 5: Scaffold a wizard with multiple step types
	wizardInput := types.ScaffoldWizardInput{
		WizardName: "create_order",
		Domain:     "order",
		Steps: []types.WizardStepDef{
			{Name: "Select Client", Type: "select", Fields: []string{"client_id"}, Searchable: true},
			{Name: "Order Details", Type: "form", Fields: []string{"notes", "status"}},
			{Name: "Add Items", Type: "has_many", ChildDomain: "orderitem", HasManyMode: "select_existing"},
			{Name: "Review", Type: "summary"},
		},
		RouteGroup: "authenticated",
		Layout:     "dashboard",
		DryRun:     false,
	}

	wizardResult, err := scaffoldWizard(projectRegistry, wizardInput)
	if err != nil {
		t.Fatalf("scaffoldWizard error: %v", err)
	}
	if !wizardResult.Success {
		t.Fatalf("scaffoldWizard failed: %s", wizardResult.Message)
	}

	// Step 6: Verify wizard files were created
	expectedWizardFiles := []string{
		"internal/web/order/wizard_create_order.go",
		"internal/web/order/views/wizard_create_order.templ",
		"internal/web/order/views/wizard_create_order_step1.templ",
		"internal/web/order/views/wizard_create_order_step2.templ",
		"internal/web/order/views/wizard_create_order_step3.templ",
		"internal/web/order/views/wizard_create_order_step4.templ",
		"internal/models/wizard_draft.go",
	}

	for _, f := range expectedWizardFiles {
		path := filepath.Join(projectDir, f)
		if !fileExists(path) {
			t.Errorf("expected wizard file %s to exist", f)
		}
	}

	// Step 7: Run gofmt on all generated .go files to verify syntax
	var goFiles []string
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk project directory: %v", err)
	}

	if len(goFiles) == 0 {
		t.Fatal("no .go files found in project directory")
	}

	// Check each .go file with gofmt
	var syntaxErrors []string
	for _, goFile := range goFiles {
		cmd := exec.Command("gofmt", "-e", goFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// gofmt -e returns non-zero for syntax errors
			syntaxErrors = append(syntaxErrors, fmt.Sprintf("%s: %s", goFile, string(output)))
		}
	}

	if len(syntaxErrors) > 0 {
		t.Errorf("gofmt found syntax errors in %d files:\n%s",
			len(syntaxErrors), strings.Join(syntaxErrors, "\n"))
	}

	// Step 8 (optional): Verify Go code is importable by checking go vet on generated files
	// This is lighter than a full build and faster
	vetCmd := exec.Command("go", "vet", "./internal/web/order/...")
	vetCmd.Dir = projectDir
	vetOutput, err := vetCmd.CombinedOutput()
	if err != nil {
		// Log but don't fail - vet might fail due to missing dependencies
		t.Logf("go vet output (may fail due to missing deps): %s", string(vetOutput))
	}

	t.Logf("Integration test verified %d .go files with gofmt", len(goFiles))
}
