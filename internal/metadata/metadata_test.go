package metadata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestStore_SaveAndLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Test saving a domain
	input := types.ScaffoldDomainInput{
		DomainName: "order",
		Fields: []types.FieldDef{
			{Name: "Title", Type: "string"},
			{Name: "Amount", Type: "float64"},
		},
		FormStyle: "modal",
	}

	err = store.SaveDomain("order", input, "0.1.0")
	if err != nil {
		t.Fatalf("SaveDomain() error = %v", err)
	}

	// Verify file was created
	metaPath := filepath.Join(tmpDir, MetadataDir, MetadataFile)
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Errorf("Metadata file was not created at %s", metaPath)
	}

	// Test loading
	meta, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if meta.Version != CurrentVersion {
		t.Errorf("Version = %q, want %q", meta.Version, CurrentVersion)
	}

	if len(meta.Domains) != 1 {
		t.Errorf("len(Domains) = %d, want 1", len(meta.Domains))
	}

	domain, exists := meta.Domains["order"]
	if !exists {
		t.Fatal("Domain 'order' not found in metadata")
	}

	if domain.ScaffolderVersion != "0.1.0" {
		t.Errorf("ScaffolderVersion = %q, want %q", domain.ScaffolderVersion, "0.1.0")
	}

	if domain.Input.DomainName != "order" {
		t.Errorf("Input.DomainName = %q, want %q", domain.Input.DomainName, "order")
	}

	if domain.Input.FormStyle != "modal" {
		t.Errorf("Input.FormStyle = %q, want %q", domain.Input.FormStyle, "modal")
	}

	if len(domain.Input.Fields) != 2 {
		t.Errorf("len(Input.Fields) = %d, want 2", len(domain.Input.Fields))
	}
}

func TestStore_GetDomain(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Test getting non-existent domain
	_, exists, err := store.GetDomain("nonexistent")
	if err != nil {
		t.Fatalf("GetDomain() error = %v", err)
	}
	if exists {
		t.Error("GetDomain() returned exists=true for non-existent domain")
	}

	// Save a domain and retrieve it
	input := types.ScaffoldDomainInput{
		DomainName: "product",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string"},
		},
	}
	err = store.SaveDomain("product", input, "0.1.0")
	if err != nil {
		t.Fatalf("SaveDomain() error = %v", err)
	}

	domain, exists, err := store.GetDomain("product")
	if err != nil {
		t.Fatalf("GetDomain() error = %v", err)
	}
	if !exists {
		t.Error("GetDomain() returned exists=false for existing domain")
	}
	if domain.Input.DomainName != "product" {
		t.Errorf("domain.Input.DomainName = %q, want %q", domain.Input.DomainName, "product")
	}
}

func TestStore_ListDomains(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Test listing empty
	domains, err := store.ListDomains()
	if err != nil {
		t.Fatalf("ListDomains() error = %v", err)
	}
	if len(domains) != 0 {
		t.Errorf("len(domains) = %d, want 0", len(domains))
	}

	// Add domains
	store.SaveDomain("order", types.ScaffoldDomainInput{DomainName: "order"}, "0.1.0")
	store.SaveDomain("product", types.ScaffoldDomainInput{DomainName: "product"}, "0.1.0")

	domains, err = store.ListDomains()
	if err != nil {
		t.Fatalf("ListDomains() error = %v", err)
	}
	if len(domains) != 2 {
		t.Errorf("len(domains) = %d, want 2", len(domains))
	}
}

func TestStore_RemoveDomain(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Add and remove domain
	store.SaveDomain("order", types.ScaffoldDomainInput{DomainName: "order"}, "0.1.0")

	exists, _ := store.Exists("order")
	if !exists {
		t.Error("Domain should exist after saving")
	}

	err = store.RemoveDomain("order")
	if err != nil {
		t.Fatalf("RemoveDomain() error = %v", err)
	}

	exists, _ = store.Exists("order")
	if exists {
		t.Error("Domain should not exist after removal")
	}
}

func TestStore_UpdateDomain(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Save initial
	input1 := types.ScaffoldDomainInput{
		DomainName: "order",
		FormStyle:  "modal",
	}
	store.SaveDomain("order", input1, "0.1.0")

	domain1, _, _ := store.GetDomain("order")
	originalTime := domain1.ScaffoldedAt

	// Update with new input
	input2 := types.ScaffoldDomainInput{
		DomainName: "order",
		FormStyle:  "page",
	}
	store.SaveDomain("order", input2, "0.2.0")

	domain2, _, _ := store.GetDomain("order")

	// Original scaffold time should be preserved
	if domain2.ScaffoldedAt != originalTime {
		t.Error("ScaffoldedAt should be preserved on update")
	}

	// Updated time should be set
	if domain2.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set on update")
	}

	// New values should be saved
	if domain2.ScaffolderVersion != "0.2.0" {
		t.Errorf("ScaffolderVersion = %q, want %q", domain2.ScaffolderVersion, "0.2.0")
	}

	if domain2.Input.FormStyle != "page" {
		t.Errorf("Input.FormStyle = %q, want %q", domain2.Input.FormStyle, "page")
	}
}

func TestStore_Exists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	exists, err := store.Exists("order")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() should return false for non-existent domain")
	}

	store.SaveDomain("order", types.ScaffoldDomainInput{DomainName: "order"}, "0.1.0")

	exists, err = store.Exists("order")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() should return true for existing domain")
	}
}

// Wizard metadata tests

func TestStore_SaveWizard(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Create wizard input with all fields
	withDrafts := true
	input := types.ScaffoldWizardInput{
		WizardName: "create_order",
		Domain:     "order",
		Steps: []types.WizardStepDef{
			{Name: "Select Client", Type: "select", Fields: []string{"client_id"}},
			{Name: "Add Items", Type: "has_many", ChildDomain: "orderitem"},
			{Name: "Review", Type: "summary"},
		},
		Layout:     "dashboard",
		RouteGroup: "admin",
		FormStyle:  "page",
		WithDrafts: &withDrafts,
	}

	err = store.SaveWizard("create_order", "order", input, "0.1.0")
	if err != nil {
		t.Fatalf("SaveWizard() error = %v", err)
	}

	// Load and verify
	meta, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify composite key "domain:wizardName"
	key := "order:create_order"
	wizard, exists := meta.Wizards[key]
	if !exists {
		t.Fatalf("Wizard with key %q not found in metadata", key)
	}

	// Verify wizard metadata fields
	if wizard.ScaffolderVersion != "0.1.0" {
		t.Errorf("ScaffolderVersion = %q, want %q", wizard.ScaffolderVersion, "0.1.0")
	}

	if wizard.Domain != "order" {
		t.Errorf("Domain = %q, want %q", wizard.Domain, "order")
	}

	if wizard.Input.WizardName != "create_order" {
		t.Errorf("Input.WizardName = %q, want %q", wizard.Input.WizardName, "create_order")
	}

	if wizard.Input.Layout != "dashboard" {
		t.Errorf("Input.Layout = %q, want %q", wizard.Input.Layout, "dashboard")
	}

	if wizard.Input.RouteGroup != "admin" {
		t.Errorf("Input.RouteGroup = %q, want %q", wizard.Input.RouteGroup, "admin")
	}

	if wizard.Input.FormStyle != "page" {
		t.Errorf("Input.FormStyle = %q, want %q", wizard.Input.FormStyle, "page")
	}

	if len(wizard.Input.Steps) != 3 {
		t.Errorf("len(Input.Steps) = %d, want 3", len(wizard.Input.Steps))
	}

	// Verify step details
	if wizard.Input.Steps[0].Type != "select" {
		t.Errorf("Steps[0].Type = %q, want %q", wizard.Input.Steps[0].Type, "select")
	}

	if wizard.Input.Steps[1].ChildDomain != "orderitem" {
		t.Errorf("Steps[1].ChildDomain = %q, want %q", wizard.Input.Steps[1].ChildDomain, "orderitem")
	}
}

func TestStore_WizardCompositeKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Create two wizards with same domain but different names
	input1 := types.ScaffoldWizardInput{
		WizardName: "create",
		Domain:     "order",
		Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "form"}},
	}
	input2 := types.ScaffoldWizardInput{
		WizardName: "checkout",
		Domain:     "order",
		Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "form"}},
	}

	store.SaveWizard("create", "order", input1, "0.1.0")
	store.SaveWizard("checkout", "order", input2, "0.1.0")

	// Load and verify both exist
	meta, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Should have 2 separate entries
	if len(meta.Wizards) != 2 {
		t.Errorf("len(Wizards) = %d, want 2", len(meta.Wizards))
	}

	// Both keys should exist
	key1 := "order:create"
	key2 := "order:checkout"

	if _, exists := meta.Wizards[key1]; !exists {
		t.Errorf("Wizard with key %q not found", key1)
	}

	if _, exists := meta.Wizards[key2]; !exists {
		t.Errorf("Wizard with key %q not found", key2)
	}
}

func TestStore_WizardDifferentDomains(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Create wizards with same name but different domains
	input1 := types.ScaffoldWizardInput{
		WizardName: "create",
		Domain:     "order",
		Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "form"}},
	}
	input2 := types.ScaffoldWizardInput{
		WizardName: "create",
		Domain:     "product",
		Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "form"}},
	}

	store.SaveWizard("create", "order", input1, "0.1.0")
	store.SaveWizard("create", "product", input2, "0.1.0")

	// Load and verify both exist with different keys
	meta, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Should have 2 separate entries
	if len(meta.Wizards) != 2 {
		t.Errorf("len(Wizards) = %d, want 2", len(meta.Wizards))
	}

	// Both keys should exist
	key1 := "order:create"
	key2 := "product:create"

	if _, exists := meta.Wizards[key1]; !exists {
		t.Errorf("Wizard with key %q not found", key1)
	}

	if _, exists := meta.Wizards[key2]; !exists {
		t.Errorf("Wizard with key %q not found", key2)
	}

	// Verify each wizard has correct domain
	if meta.Wizards[key1].Domain != "order" {
		t.Errorf("Wizard[%q].Domain = %q, want %q", key1, meta.Wizards[key1].Domain, "order")
	}

	if meta.Wizards[key2].Domain != "product" {
		t.Errorf("Wizard[%q].Domain = %q, want %q", key2, meta.Wizards[key2].Domain, "product")
	}
}

func TestStore_UpdateWizard(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Save initial wizard
	input1 := types.ScaffoldWizardInput{
		WizardName: "create",
		Domain:     "order",
		Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "form"}},
		Layout:     "dashboard",
	}
	store.SaveWizard("create", "order", input1, "0.1.0")

	// Get original timestamp
	meta1, _ := store.Load()
	originalTime := meta1.Wizards["order:create"].ScaffoldedAt

	// Update wizard
	input2 := types.ScaffoldWizardInput{
		WizardName: "create",
		Domain:     "order",
		Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "form"}, {Name: "Step 2", Type: "summary"}},
		Layout:     "base",
	}
	store.SaveWizard("create", "order", input2, "0.2.0")

	// Load and verify update
	meta2, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	wizard := meta2.Wizards["order:create"]

	// Original scaffold time should be preserved
	if wizard.ScaffoldedAt != originalTime {
		t.Error("ScaffoldedAt should be preserved on update")
	}

	// Updated time should be set
	if wizard.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set on update")
	}

	// New values should be saved
	if wizard.ScaffolderVersion != "0.2.0" {
		t.Errorf("ScaffolderVersion = %q, want %q", wizard.ScaffolderVersion, "0.2.0")
	}

	if wizard.Input.Layout != "base" {
		t.Errorf("Input.Layout = %q, want %q", wizard.Input.Layout, "base")
	}

	if len(wizard.Input.Steps) != 2 {
		t.Errorf("len(Input.Steps) = %d, want 2", len(wizard.Input.Steps))
	}
}
