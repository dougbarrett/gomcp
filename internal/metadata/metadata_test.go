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
