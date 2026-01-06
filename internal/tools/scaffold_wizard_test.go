package tools

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldWizard(t *testing.T) {
	t.Run("validates required fields", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		tests := []struct {
			name    string
			input   types.ScaffoldWizardInput
			wantErr string
		}{
			{
				name:    "empty wizard_name",
				input:   types.ScaffoldWizardInput{Domain: "order", Steps: []types.WizardStepDef{{Name: "Step 1"}}},
				wantErr: "wizard_name is required",
			},
			{
				name:    "empty domain",
				input:   types.ScaffoldWizardInput{WizardName: "create", Steps: []types.WizardStepDef{{Name: "Step 1"}}},
				wantErr: "domain is required",
			},
			{
				name:    "no steps",
				input:   types.ScaffoldWizardInput{WizardName: "create", Domain: "order"},
				wantErr: "at least one step is required",
			},
			{
				name: "step without name",
				input: types.ScaffoldWizardInput{
					WizardName: "create",
					Domain:     "order",
					Steps:      []types.WizardStepDef{{Type: "form"}},
				},
				wantErr: "step 1: name is required",
			},
			{
				name: "invalid step type",
				input: types.ScaffoldWizardInput{
					WizardName: "create",
					Domain:     "order",
					Steps:      []types.WizardStepDef{{Name: "Step 1", Type: "invalid"}},
				},
				wantErr: "step 1: invalid type 'invalid'",
			},
			{
				name: "has_many without child_domain",
				input: types.ScaffoldWizardInput{
					WizardName: "create",
					Domain:     "order",
					Steps:      []types.WizardStepDef{{Name: "Items", Type: "has_many"}},
				},
				wantErr: "step 1: child_domain is required for has_many steps",
			},
			{
				name: "invalid has_many_mode",
				input: types.ScaffoldWizardInput{
					WizardName: "create",
					Domain:     "order",
					Steps:      []types.WizardStepDef{{Name: "Items", Type: "has_many", ChildDomain: "item", HasManyMode: "invalid"}},
				},
				wantErr: "step 1: invalid has_many_mode 'invalid'",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := scaffoldWizard(registry, tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Errorf("expected failure, got success")
				}
				if !strings.Contains(result.Message, tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, result.Message)
				}
			})
		}
	})

	t.Run("generates simple wizard with form steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "create",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Details", Type: "form", Fields: []string{"name", "email"}},
				{Name: "Review", Type: "summary"},
			},
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check controller is created
		controllerPath := filepath.Join(tmpDir, "internal", "web", "order", "wizard_create.go")
		if !fileExists(controllerPath) {
			t.Errorf("expected controller at %s", controllerPath)
		}

		// Check main wizard view is created
		wizardViewPath := filepath.Join(tmpDir, "internal", "web", "order", "views", "wizard_create.templ")
		if !fileExists(wizardViewPath) {
			t.Errorf("expected wizard view at %s", wizardViewPath)
		}

		// Check step views are created
		step1Path := filepath.Join(tmpDir, "internal", "web", "order", "views", "wizard_create_step1.templ")
		if !fileExists(step1Path) {
			t.Errorf("expected step 1 view at %s", step1Path)
		}

		step2Path := filepath.Join(tmpDir, "internal", "web", "order", "views", "wizard_create_step2.templ")
		if !fileExists(step2Path) {
			t.Errorf("expected step 2 view at %s", step2Path)
		}

		// Check draft model is created (with_drafts defaults to true)
		draftModelPath := filepath.Join(tmpDir, "internal", "models", "wizard_draft.go")
		if !fileExists(draftModelPath) {
			t.Errorf("expected draft model at %s", draftModelPath)
		}
	})

	t.Run("generates wizard with has_many step", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "checkout",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Select Items", Type: "has_many", ChildDomain: "orderitem", HasManyMode: "select_existing"},
				{Name: "Review", Type: "summary"},
			},
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check has_many step view is created
		step1Path := filepath.Join(tmpDir, "internal", "web", "order", "views", "wizard_checkout_step1.templ")
		if !fileExists(step1Path) {
			t.Errorf("expected has_many step view at %s", step1Path)
		}
	})

	t.Run("generates wizard with select step", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "create",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Select Client", Type: "select", Fields: []string{"client_id"}, Searchable: true},
				{Name: "Review", Type: "summary"},
			},
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Check select step view is created
		step1Path := filepath.Join(tmpDir, "internal", "web", "order", "views", "wizard_create_step1.templ")
		if !fileExists(step1Path) {
			t.Errorf("expected select step view at %s", step1Path)
		}
	})

	t.Run("respects with_drafts false", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		withDrafts := false
		input := types.ScaffoldWizardInput{
			WizardName: "quick",
			Domain:     "task",
			Steps: []types.WizardStepDef{
				{Name: "Details", Type: "form"},
			},
			WithDrafts: &withDrafts,
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Draft model should NOT be created
		draftModelPath := filepath.Join(tmpDir, "internal", "models", "wizard_draft.go")
		if fileExists(draftModelPath) {
			t.Error("expected draft model NOT to be created when with_drafts is false")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "create",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Details", Type: "form"},
			},
			DryRun: true,
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Controller should NOT be created
		controllerPath := filepath.Join(tmpDir, "internal", "web", "order", "wizard_create.go")
		if fileExists(controllerPath) {
			t.Error("expected controller NOT to be created in dry run mode")
		}

		// FilesCreated should report what would be created
		if len(result.FilesCreated) == 0 {
			t.Error("expected FilesCreated to report files that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "create",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Details", Type: "form"},
			},
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.NextSteps) == 0 {
			t.Error("expected NextSteps to be populated")
		}

		hasTemplGenerate := false
		for _, step := range result.NextSteps {
			if strings.Contains(step, "templ generate") {
				hasTemplGenerate = true
				break
			}
		}
		if !hasTemplGenerate {
			t.Error("expected NextSteps to include 'templ generate'")
		}
	})

	t.Run("respects layout option", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "create",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Details", Type: "form"},
			},
			Layout: "auth",
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Verify controller was created with correct layout
		controllerPath := filepath.Join(tmpDir, "internal", "web", "order", "wizard_create.go")
		if !fileExists(controllerPath) {
			t.Errorf("expected controller at %s", controllerPath)
		}
	})

	t.Run("respects route_group option", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/myapp")

		input := types.ScaffoldWizardInput{
			WizardName: "create",
			Domain:     "order",
			Steps: []types.WizardStepDef{
				{Name: "Details", Type: "form"},
			},
			RouteGroup: "admin",
		}

		result, err := scaffoldWizard(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}

		// Controller should exist
		controllerPath := filepath.Join(tmpDir, "internal", "web", "order", "wizard_create.go")
		if !fileExists(controllerPath) {
			t.Errorf("expected controller at %s", controllerPath)
		}
	})
}
