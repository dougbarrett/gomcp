package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDomains(t *testing.T) {
	t.Run("returns empty list when no domains exist", func(t *testing.T) {
		registry, _ := testRegistry(t)

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}
		if len(result.Domains) != 0 {
			t.Errorf("expected 0 domains, got %d", len(result.Domains))
		}
	})

	t.Run("detects model files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create a model file
		modelsDir := filepath.Join(tmpDir, "internal", "models")
		if err := os.MkdirAll(modelsDir, 0755); err != nil {
			t.Fatalf("failed to create models dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(modelsDir, "product.go"), []byte("package models"), 0644); err != nil {
			t.Fatalf("failed to create model file: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Domains) != 1 {
			t.Errorf("expected 1 domain, got %d", len(result.Domains))
		}

		found := false
		for _, d := range result.Domains {
			if d.Name == "product" && d.HasModel {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find product domain with HasModel=true")
		}
	})

	t.Run("skips base.go in models", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create models directory with base.go
		modelsDir := filepath.Join(tmpDir, "internal", "models")
		if err := os.MkdirAll(modelsDir, 0755); err != nil {
			t.Fatalf("failed to create models dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(modelsDir, "base.go"), []byte("package models"), 0644); err != nil {
			t.Fatalf("failed to create base file: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// base.go should be skipped
		for _, d := range result.Domains {
			if d.Name == "base" {
				t.Errorf("base.go should be skipped")
			}
		}
	})

	t.Run("detects repository directories", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create a repository directory
		repoDir := filepath.Join(tmpDir, "internal", "repository", "product")
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			t.Fatalf("failed to create repo dir: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, d := range result.Domains {
			if d.Name == "product" && d.HasRepository {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find product domain with HasRepository=true")
		}
	})

	t.Run("detects service directories", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create a service directory
		serviceDir := filepath.Join(tmpDir, "internal", "services", "product")
		if err := os.MkdirAll(serviceDir, 0755); err != nil {
			t.Fatalf("failed to create service dir: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, d := range result.Domains {
			if d.Name == "product" && d.HasService {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find product domain with HasService=true")
		}
	})

	t.Run("detects controller directories", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create a web/controller directory
		webDir := filepath.Join(tmpDir, "internal", "web", "product")
		if err := os.MkdirAll(webDir, 0755); err != nil {
			t.Fatalf("failed to create web dir: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, d := range result.Domains {
			if d.Name == "product" && d.HasController {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find product domain with HasController=true")
		}
	})

	t.Run("skips middleware and layouts directories", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create middleware and layouts directories
		for _, dir := range []string{"middleware", "layouts", "components"} {
			webDir := filepath.Join(tmpDir, "internal", "web", dir)
			if err := os.MkdirAll(webDir, 0755); err != nil {
				t.Fatalf("failed to create %s dir: %v", dir, err)
			}
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// middleware, layouts, components should be skipped
		skipDirs := []string{"middleware", "layouts", "components"}
		for _, d := range result.Domains {
			for _, skip := range skipDirs {
				if d.Name == skip {
					t.Errorf("%s directory should be skipped", skip)
				}
			}
		}
	})

	t.Run("detects views in controller directories", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create a web directory with views
		viewsDir := filepath.Join(tmpDir, "internal", "web", "product", "views")
		if err := os.MkdirAll(viewsDir, 0755); err != nil {
			t.Fatalf("failed to create views dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(viewsDir, "list.templ"), []byte("package views"), 0644); err != nil {
			t.Fatalf("failed to create view file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(viewsDir, "show.templ"), []byte("package views"), 0644); err != nil {
			t.Fatalf("failed to create view file: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		found := false
		for _, d := range result.Domains {
			if d.Name == "product" && d.HasController {
				found = true
				if len(d.Views) != 2 {
					t.Errorf("expected 2 views, got %d", len(d.Views))
				}
				// Check view names
				hasListView := false
				hasShowView := false
				for _, v := range d.Views {
					if v == "list" {
						hasListView = true
					}
					if v == "show" {
						hasShowView = true
					}
				}
				if !hasListView {
					t.Errorf("expected list view")
				}
				if !hasShowView {
					t.Errorf("expected show view")
				}
				break
			}
		}
		if !found {
			t.Errorf("expected to find product domain with views")
		}
	})

	t.Run("combines information from multiple directories", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create a complete domain structure
		dirs := []string{
			"internal/models",
			"internal/repository/product",
			"internal/services/product",
			"internal/web/product",
		}
		for _, dir := range dirs {
			if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
				t.Fatalf("failed to create dir %s: %v", dir, err)
			}
		}

		// Create model file
		if err := os.WriteFile(filepath.Join(tmpDir, "internal/models/product.go"), []byte("package models"), 0644); err != nil {
			t.Fatalf("failed to create model file: %v", err)
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have exactly one domain with all flags set
		if len(result.Domains) != 1 {
			t.Errorf("expected 1 domain, got %d", len(result.Domains))
		}

		d := result.Domains[0]
		if d.Name != "product" {
			t.Errorf("expected product domain, got %s", d.Name)
		}
		if !d.HasModel {
			t.Errorf("expected HasModel=true")
		}
		if !d.HasRepository {
			t.Errorf("expected HasRepository=true")
		}
		if !d.HasService {
			t.Errorf("expected HasService=true")
		}
		if !d.HasController {
			t.Errorf("expected HasController=true")
		}
	})

	t.Run("handles multiple domains", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create multiple domain structures
		domains := []string{"product", "user", "order"}
		for _, domain := range domains {
			modelsDir := filepath.Join(tmpDir, "internal", "models")
			if err := os.MkdirAll(modelsDir, 0755); err != nil {
				t.Fatalf("failed to create models dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(modelsDir, domain+".go"), []byte("package models"), 0644); err != nil {
				t.Fatalf("failed to create model file: %v", err)
			}
		}

		result, err := listDomains(registry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Domains) != 3 {
			t.Errorf("expected 3 domains, got %d", len(result.Domains))
		}

		// Check all domains are present
		domainNames := make(map[string]bool)
		for _, d := range result.Domains {
			domainNames[d.Name] = true
		}
		for _, domain := range domains {
			if !domainNames[domain] {
				t.Errorf("expected to find domain %s", domain)
			}
		}
	})
}
