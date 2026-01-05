package tools

import (
	"os"
	"path/filepath"
	"testing"
)

// testRegistry creates a registry with a temporary directory for testing.
func testRegistry(t *testing.T) (*Registry, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "tools-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	return NewRegistry(tmpDir), tmpDir
}

// setupGoMod creates a go.mod file in the temp directory.
func setupGoMod(t *testing.T, tmpDir, modulePath string) {
	t.Helper()
	goModContent := "module " + modulePath + "\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// readFile reads a file and returns its content.
func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}

// TestNewRegistry tests the registry constructor.
func TestNewRegistry(t *testing.T) {
	t.Run("with working dir", func(t *testing.T) {
		r := NewRegistry("/tmp/test")
		if r.WorkingDir != "/tmp/test" {
			t.Errorf("expected WorkingDir to be /tmp/test, got %s", r.WorkingDir)
		}
	})

	t.Run("empty working dir uses cwd", func(t *testing.T) {
		r := NewRegistry("")
		cwd, _ := os.Getwd()
		if r.WorkingDir != cwd {
			t.Errorf("expected WorkingDir to be %s, got %s", cwd, r.WorkingDir)
		}
	})
}

// TestNewGenerator tests the generator factory.
func TestNewGenerator(t *testing.T) {
	registry, tmpDir := testRegistry(t)

	t.Run("empty project path uses working dir", func(t *testing.T) {
		gen := registry.NewGenerator("")
		if gen.BasePath() != tmpDir {
			t.Errorf("expected BasePath to be %s, got %s", tmpDir, gen.BasePath())
		}
	})

	t.Run("custom project path", func(t *testing.T) {
		gen := registry.NewGenerator("/custom/path")
		if gen.BasePath() != "/custom/path" {
			t.Errorf("expected BasePath to be /custom/path, got %s", gen.BasePath())
		}
	})
}
