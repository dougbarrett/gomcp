package server

import (
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	t.Run("ServerName is set", func(t *testing.T) {
		if ServerName == "" {
			t.Error("ServerName should not be empty")
		}
		if ServerName != "go-web-scaffold" {
			t.Errorf("ServerName = %q, want %q", ServerName, "go-web-scaffold")
		}
	})

	t.Run("ServerVersion is set", func(t *testing.T) {
		if ServerVersion == "" {
			t.Error("ServerVersion should not be empty")
		}
		// Version should follow semver pattern
		parts := strings.Split(ServerVersion, ".")
		if len(parts) != 3 {
			t.Errorf("ServerVersion %q should follow semver (x.y.z)", ServerVersion)
		}
	})
}

func TestConfig(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		cfg := &Config{}
		if cfg.WorkingDir != "" {
			t.Errorf("WorkingDir should be empty by default, got %q", cfg.WorkingDir)
		}
	})

	t.Run("with working dir", func(t *testing.T) {
		cfg := &Config{WorkingDir: "/tmp/test"}
		if cfg.WorkingDir != "/tmp/test" {
			t.Errorf("WorkingDir = %q, want %q", cfg.WorkingDir, "/tmp/test")
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("creates server with nil config", func(t *testing.T) {
		server := New(nil)
		if server == nil {
			t.Error("New(nil) should return a non-nil server")
		}
	})

	t.Run("creates server with empty config", func(t *testing.T) {
		cfg := &Config{}
		server := New(cfg)
		if server == nil {
			t.Error("New(&Config{}) should return a non-nil server")
		}
	})

	t.Run("creates server with working dir", func(t *testing.T) {
		cfg := &Config{WorkingDir: "/tmp/test"}
		server := New(cfg)
		if server == nil {
			t.Error("New should return a non-nil server")
		}
	})
}
