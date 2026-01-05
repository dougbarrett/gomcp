package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseRepositoryMethods_Interface(t *testing.T) {
	// Create a test repository file with an interface
	tmpDir := t.TempDir()
	repoFile := filepath.Join(tmpDir, "repo.go")

	interfaceRepo := `package repo

import "context"

type Repository interface {
	Create(ctx context.Context, entity *Entity) error
	FindByID(ctx context.Context, id uint) (*Entity, error)
	FindAll(ctx context.Context) ([]Entity, error)
	Update(ctx context.Context, entity *Entity) error
	Delete(ctx context.Context, id uint) error
}
`
	if err := os.WriteFile(repoFile, []byte(interfaceRepo), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	parsed, err := parseRepositoryMethods(repoFile)
	if err != nil {
		t.Fatalf("parseRepositoryMethods error: %v", err)
	}

	if parsed.IsStruct {
		t.Error("expected IsStruct to be false for interface repository")
	}

	if len(parsed.Methods) != 5 {
		t.Errorf("expected 5 methods, got %d", len(parsed.Methods))
	}

	// Verify method names
	expectedMethods := []string{"Create", "FindByID", "FindAll", "Update", "Delete"}
	for i, expected := range expectedMethods {
		if parsed.Methods[i].Name != expected {
			t.Errorf("expected method %d to be %s, got %s", i, expected, parsed.Methods[i].Name)
		}
	}

	// Verify HasContext is set
	for _, method := range parsed.Methods {
		if !method.HasContext {
			t.Errorf("expected method %s to have context", method.Name)
		}
	}
}

func TestParseRepositoryMethods_Struct(t *testing.T) {
	// Create a test repository file with a struct (like auth user repo)
	tmpDir := t.TempDir()
	repoFile := filepath.Join(tmpDir, "repo.go")

	structRepo := `package repo

import (
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, entity *Entity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *Repository) FindByID(ctx context.Context, id uint) (*Entity, error) {
	var entity Entity
	err := r.db.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*Entity, error) {
	var entity Entity
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&entity).Error
	return &entity, err
}

func (r *Repository) Update(ctx context.Context, entity *Entity) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Entity{}, id).Error
}

// unexportedMethod should not be included
func (r *Repository) unexportedMethod() {}
`
	if err := os.WriteFile(repoFile, []byte(structRepo), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	parsed, err := parseRepositoryMethods(repoFile)
	if err != nil {
		t.Fatalf("parseRepositoryMethods error: %v", err)
	}

	if !parsed.IsStruct {
		t.Error("expected IsStruct to be true for struct repository")
	}

	// Should have 5 methods (excluding NewRepository which isn't a method, and unexportedMethod)
	if len(parsed.Methods) != 5 {
		t.Errorf("expected 5 methods, got %d: %v", len(parsed.Methods), methodNames(parsed.Methods))
	}

	// Verify method names
	expectedMethods := map[string]bool{
		"Create":      true,
		"FindByID":    true,
		"FindByEmail": true,
		"Update":      true,
		"Delete":      true,
	}

	for _, method := range parsed.Methods {
		if !expectedMethods[method.Name] {
			t.Errorf("unexpected method: %s", method.Name)
		}
		delete(expectedMethods, method.Name)
	}

	if len(expectedMethods) > 0 {
		t.Errorf("missing methods: %v", expectedMethods)
	}
}

func TestParseRepositoryMethods_NoRepository(t *testing.T) {
	tmpDir := t.TempDir()
	repoFile := filepath.Join(tmpDir, "repo.go")

	noRepo := `package repo

type SomeOtherType struct {
	field string
}

func (s *SomeOtherType) DoSomething() {}
`
	if err := os.WriteFile(repoFile, []byte(noRepo), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	parsed, err := parseRepositoryMethods(repoFile)
	if err != nil {
		t.Fatalf("parseRepositoryMethods error: %v", err)
	}

	if len(parsed.Methods) != 0 {
		t.Errorf("expected 0 methods for non-Repository type, got %d", len(parsed.Methods))
	}
}

func TestFilterMethods(t *testing.T) {
	methods := []RepoMethod{
		{Name: "Create"},
		{Name: "FindByID"},
		{Name: "FindByEmail"},
		{Name: "Update"},
		{Name: "Delete"},
	}

	t.Run("include filter", func(t *testing.T) {
		filtered := filterMethods(methods, []string{"Create", "FindByID"}, nil)
		if len(filtered) != 2 {
			t.Errorf("expected 2 methods, got %d", len(filtered))
		}
	})

	t.Run("exclude filter", func(t *testing.T) {
		filtered := filterMethods(methods, nil, []string{"Delete", "Update"})
		if len(filtered) != 3 {
			t.Errorf("expected 3 methods, got %d", len(filtered))
		}
	})

	t.Run("include and exclude", func(t *testing.T) {
		filtered := filterMethods(methods, []string{"Create", "FindByID", "Update"}, []string{"Update"})
		if len(filtered) != 2 {
			t.Errorf("expected 2 methods, got %d", len(filtered))
		}
	})
}

func methodNames(methods []RepoMethod) []string {
	names := make([]string, len(methods))
	for i, m := range methods {
		names[i] = m.Name
	}
	return names
}
