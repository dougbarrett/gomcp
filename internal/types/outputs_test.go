package types

import (
	"testing"
)

// TestNewSuccessResult tests creating a success result.
func TestNewSuccessResult(t *testing.T) {
	result := NewSuccessResult("Operation completed")

	if !result.Success {
		t.Error("Success should be true")
	}
	if result.Message != "Operation completed" {
		t.Errorf("Message = %q, want %q", result.Message, "Operation completed")
	}
	if len(result.FilesCreated) != 0 {
		t.Error("FilesCreated should be empty")
	}
	if len(result.FilesUpdated) != 0 {
		t.Error("FilesUpdated should be empty")
	}
	if len(result.NextSteps) != 0 {
		t.Error("NextSteps should be empty")
	}
}

// TestNewSuccessResultWithFiles tests creating a success result with files.
func TestNewSuccessResultWithFiles(t *testing.T) {
	created := []string{"file1.go", "file2.go"}
	updated := []string{"main.go"}

	result := NewSuccessResultWithFiles("Files generated", created, updated)

	if !result.Success {
		t.Error("Success should be true")
	}
	if result.Message != "Files generated" {
		t.Errorf("Message = %q, want %q", result.Message, "Files generated")
	}
	if len(result.FilesCreated) != 2 {
		t.Errorf("len(FilesCreated) = %d, want 2", len(result.FilesCreated))
	}
	if len(result.FilesUpdated) != 1 {
		t.Errorf("len(FilesUpdated) = %d, want 1", len(result.FilesUpdated))
	}
}

// TestNewSuccessResultWithFiles_Empty tests with empty slices.
func TestNewSuccessResultWithFiles_Empty(t *testing.T) {
	result := NewSuccessResultWithFiles("Done", nil, nil)

	if !result.Success {
		t.Error("Success should be true")
	}
	if result.FilesCreated != nil {
		t.Error("FilesCreated should be nil")
	}
	if result.FilesUpdated != nil {
		t.Error("FilesUpdated should be nil")
	}
}

// TestNewErrorResult tests creating an error result.
func TestNewErrorResult(t *testing.T) {
	result := NewErrorResult("Something went wrong")

	if result.Success {
		t.Error("Success should be false")
	}
	if result.Message != "Something went wrong" {
		t.Errorf("Message = %q, want %q", result.Message, "Something went wrong")
	}
}

// TestScaffoldResult_WithNextSteps tests adding next steps.
func TestScaffoldResult_WithNextSteps(t *testing.T) {
	result := NewSuccessResult("Done").WithNextSteps(
		"Run go mod tidy",
		"Run templ generate",
		"Start the server",
	)

	if len(result.NextSteps) != 3 {
		t.Errorf("len(NextSteps) = %d, want 3", len(result.NextSteps))
	}
	if result.NextSteps[0] != "Run go mod tidy" {
		t.Errorf("NextSteps[0] = %q, want %q", result.NextSteps[0], "Run go mod tidy")
	}
}

// TestScaffoldResult_WithNextSteps_Empty tests with no steps.
func TestScaffoldResult_WithNextSteps_Empty(t *testing.T) {
	result := NewSuccessResult("Done").WithNextSteps()

	if len(result.NextSteps) != 0 {
		t.Errorf("len(NextSteps) = %d, want 0", len(result.NextSteps))
	}
}

// TestScaffoldResult_WithFilesCreated tests adding created files.
func TestScaffoldResult_WithFilesCreated(t *testing.T) {
	result := NewSuccessResult("Done").
		WithFilesCreated("file1.go").
		WithFilesCreated("file2.go", "file3.go")

	if len(result.FilesCreated) != 3 {
		t.Errorf("len(FilesCreated) = %d, want 3", len(result.FilesCreated))
	}
}

// TestScaffoldResult_WithFilesUpdated tests adding updated files.
func TestScaffoldResult_WithFilesUpdated(t *testing.T) {
	result := NewSuccessResult("Done").
		WithFilesUpdated("main.go").
		WithFilesUpdated("config.go")

	if len(result.FilesUpdated) != 2 {
		t.Errorf("len(FilesUpdated) = %d, want 2", len(result.FilesUpdated))
	}
}

// TestScaffoldResult_Chaining tests method chaining.
func TestScaffoldResult_Chaining(t *testing.T) {
	result := NewSuccessResult("All done").
		WithFilesCreated("new.go").
		WithFilesUpdated("existing.go").
		WithNextSteps("Test", "Deploy")

	if !result.Success {
		t.Error("Success should be true")
	}
	if len(result.FilesCreated) != 1 {
		t.Error("Should have 1 created file")
	}
	if len(result.FilesUpdated) != 1 {
		t.Error("Should have 1 updated file")
	}
	if len(result.NextSteps) != 2 {
		t.Error("Should have 2 next steps")
	}
}

// TestDomainInfo tests DomainInfo struct.
func TestDomainInfo(t *testing.T) {
	info := DomainInfo{
		Name:          "product",
		HasModel:      true,
		HasRepository: true,
		HasService:    true,
		HasController: true,
		Views:         []string{"list.templ", "show.templ", "form.templ"},
	}

	if info.Name != "product" {
		t.Errorf("Name = %q, want %q", info.Name, "product")
	}
	if !info.HasModel {
		t.Error("HasModel should be true")
	}
	if !info.HasRepository {
		t.Error("HasRepository should be true")
	}
	if !info.HasService {
		t.Error("HasService should be true")
	}
	if !info.HasController {
		t.Error("HasController should be true")
	}
	if len(info.Views) != 3 {
		t.Errorf("len(Views) = %d, want 3", len(info.Views))
	}
}

// TestDomainInfo_Partial tests DomainInfo with partial implementation.
func TestDomainInfo_Partial(t *testing.T) {
	info := DomainInfo{
		Name:     "auth",
		HasModel: true,
		// Other fields default to false/nil
	}

	if info.Name != "auth" {
		t.Errorf("Name = %q, want %q", info.Name, "auth")
	}
	if !info.HasModel {
		t.Error("HasModel should be true")
	}
	if info.HasRepository {
		t.Error("HasRepository should be false")
	}
	if info.Views != nil {
		t.Error("Views should be nil")
	}
}

// TestNewListDomainsResult tests creating a list domains result.
func TestNewListDomainsResult(t *testing.T) {
	domains := []DomainInfo{
		{Name: "product", HasModel: true},
		{Name: "user", HasModel: true, HasRepository: true},
	}

	result := NewListDomainsResult(domains)

	if !result.Success {
		t.Error("Success should be true")
	}
	if len(result.Domains) != 2 {
		t.Errorf("len(Domains) = %d, want 2", len(result.Domains))
	}
	if result.Domains[0].Name != "product" {
		t.Errorf("Domains[0].Name = %q, want %q", result.Domains[0].Name, "product")
	}
}

// TestNewListDomainsResult_Empty tests with empty domains.
func TestNewListDomainsResult_Empty(t *testing.T) {
	result := NewListDomainsResult(nil)

	if !result.Success {
		t.Error("Success should be true")
	}
	if result.Domains != nil {
		t.Error("Domains should be nil")
	}
}

// TestNewListDomainsError tests creating an error list domains result.
func TestNewListDomainsError(t *testing.T) {
	result := NewListDomainsError("Failed to scan project")

	if result.Success {
		t.Error("Success should be false")
	}
	if result.Message != "Failed to scan project" {
		t.Errorf("Message = %q, want %q", result.Message, "Failed to scan project")
	}
	if result.Domains != nil {
		t.Error("Domains should be nil on error")
	}
}
