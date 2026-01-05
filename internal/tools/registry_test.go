package tools

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegistry_RegisterAll(t *testing.T) {
	t.Run("registers all tools without panic", func(t *testing.T) {
		registry := NewRegistry("/tmp/test")
		server := mcp.NewServer(&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		}, nil)

		// This should not panic
		registry.RegisterAll(server)
	})

	t.Run("can be called multiple times safely", func(t *testing.T) {
		registry := NewRegistry("/tmp/test")
		server := mcp.NewServer(&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		}, nil)

		// Register twice - should not panic
		registry.RegisterAll(server)
		registry.RegisterAll(server)
	})
}

func TestRegistry_WorkingDirEdgeCases(t *testing.T) {
	t.Run("preserves path with spaces", func(t *testing.T) {
		path := "/path/with spaces/project"
		registry := NewRegistry(path)

		if registry.WorkingDir != path {
			t.Errorf("expected path with spaces to be preserved")
		}
	})

	t.Run("handles path with special characters", func(t *testing.T) {
		path := "/path/with-dashes_and_underscores/project"
		registry := NewRegistry(path)

		if registry.WorkingDir != path {
			t.Errorf("expected path with special chars to be preserved")
		}
	})

	t.Run("handles very long path", func(t *testing.T) {
		longPath := "/very/long/path"
		for i := 0; i < 20; i++ {
			longPath += "/segment"
		}
		registry := NewRegistry(longPath)

		if registry.WorkingDir != longPath {
			t.Errorf("expected long path to be preserved")
		}
	})
}

func TestRegistry_GeneratorIntegration(t *testing.T) {
	t.Run("generator can generate files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		gen := registry.NewGenerator("")
		gen.SetDryRun(true)

		// Ensure a directory can be created
		err := gen.EnsureDir("internal/models")
		if err != nil {
			t.Errorf("expected EnsureDir to succeed: %v", err)
		}
	})

	t.Run("generator uses correct base path", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		gen := registry.NewGenerator("")
		if gen.BasePath() != tmpDir {
			t.Errorf("expected generator BasePath to match registry WorkingDir")
		}
	})
}
