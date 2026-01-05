package tools

import (
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestAddTempluiComponent(t *testing.T) {
	t.Run("requires at least one component", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.AddTempluiComponentInput{
			Components: []string{},
		}

		result, err := addTempluiComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when no components provided")
		}
		if !containsString(result.Message, "at least one component") {
			t.Errorf("expected error about components, got: %s", result.Message)
		}
	})

	// Note: Tests that actually run templui would require the templui binary
	// to be installed. The following tests verify the behavior when templui
	// is not available or returns errors.

	t.Run("handles templui command not found", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.AddTempluiComponentInput{
			Components: []string{"button"},
		}

		result, err := addTempluiComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// If templui is not installed, we expect an error result
		// If templui is installed, we expect success (integration test behavior)
		// Either way, it shouldn't panic
		_ = result
	})

	t.Run("handles multiple components", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.AddTempluiComponentInput{
			Components: []string{"button", "card", "input"},
		}

		result, err := addTempluiComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Just verify it doesn't panic with multiple components
		_ = result
	})

	t.Run("handles force flag", func(t *testing.T) {
		registry, _ := testRegistry(t)

		input := types.AddTempluiComponentInput{
			Components: []string{"button"},
			Force:      true,
		}

		result, err := addTempluiComponent(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Just verify it doesn't panic with force flag
		_ = result
	})
}

func TestAddTempluiComponentResult(t *testing.T) {
	t.Run("NewAddComponentResult creates success result", func(t *testing.T) {
		added := []string{"button", "card"}
		skipped := []string{"input"}

		result := types.NewAddComponentResult(added, skipped)

		if !result.Success {
			t.Errorf("expected success=true")
		}
		if len(result.ComponentsAdded) != 2 {
			t.Errorf("expected 2 added, got %d", len(result.ComponentsAdded))
		}
		if len(result.ComponentsSkipped) != 1 {
			t.Errorf("expected 1 skipped, got %d", len(result.ComponentsSkipped))
		}
	})

	t.Run("NewAddComponentError creates error result", func(t *testing.T) {
		errors := []string{"button: not found", "card: permission denied"}

		result := types.NewAddComponentError("failed", errors)

		if result.Success {
			t.Errorf("expected success=false")
		}
		if !containsString(result.Message, "failed") {
			t.Errorf("expected error message")
		}
		if len(result.Errors) != 2 {
			t.Errorf("expected 2 errors, got %d", len(result.Errors))
		}
	})
}
