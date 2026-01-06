package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dbb1dev/go-mcp/internal/generator"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterScaffoldWizard registers the scaffold_wizard tool.
func RegisterScaffoldWizard(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "scaffold_wizard",
		Description: `Create multi-step wizard flows for complex entity creation.

Step types:
- form: Standard form fields for a step
- select: Selection from existing records (e.g., select client)
- has_many: Add multiple related items (select_existing or create_inline mode)
- summary: Review all selections before submit

Layout options (same as scaffold_domain):
- "dashboard" (default): Views wrapped in DashboardPage layout with sidebar
- "base": Views wrapped in BasePage layout without sidebar
- "auth": Views wrapped in AuthPage layout (centered card)
- "none": Views rendered without layout wrapper

Route group options (same as scaffold_domain):
- "public" (default): No authentication required
- "authenticated": Requires user login (RequireAuth middleware)
- "admin": Requires admin role (RequireAuth + RequireAdmin middleware)

Form style options:
- "page" (default for wizards): Each step is a full page
- "modal": Steps displayed in modal overlays (for simple wizards)

Draft persistence:
- Wizard progress is saved to database (with_drafts: true by default)
- Users can resume incomplete wizards via /domain/wizard/{draft_id}
- Old drafts are cleaned up automatically

Examples:

1. Simple public wizard (guest checkout):
   scaffold_wizard: {
     wizard_name: "checkout",
     domain: "order",
     steps: [
       {name: "Your Details", type: "form", fields: ["email", "phone"]},
       {name: "Review", type: "summary"}
     ]
   }

2. Admin order creation wizard:
   scaffold_wizard: {
     wizard_name: "create_order",
     domain: "order",
     route_group: "admin",
     layout: "dashboard",
     steps: [
       {name: "Select Client", type: "select", fields: ["client_id"], searchable: true},
       {name: "Add Products", type: "has_many", child_domain: "orderitem", has_many_mode: "select_existing"},
       {name: "Review & Confirm", type: "summary", fields: ["discount_id", "notes"]}
     ],
     success_redirect: "/admin/orders"
   }

3. Authenticated user wizard with inline creation:
   scaffold_wizard: {
     wizard_name: "create_project",
     domain: "project",
     route_group: "authenticated",
     steps: [
       {name: "Project Details", type: "form", fields: ["name", "description"]},
       {name: "Add Team Members", type: "has_many", child_domain: "projectmember", has_many_mode: "create_inline"},
       {name: "Review", type: "summary"}
     ]
   }

Use dry_run: true to preview all generated files first.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ScaffoldWizardInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := scaffoldWizard(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

func scaffoldWizard(registry *Registry, input types.ScaffoldWizardInput) (types.ScaffoldResult, error) {
	// Validate input
	if input.WizardName == "" {
		return types.NewErrorResult("wizard_name is required"), nil
	}
	if input.Domain == "" {
		return types.NewErrorResult("domain is required"), nil
	}
	if len(input.Steps) == 0 {
		return types.NewErrorResult("at least one step is required"), nil
	}

	// Validate wizard name
	if err := utils.ValidateComponentName(input.WizardName); err != nil {
		return types.NewErrorResult(fmt.Sprintf("invalid wizard_name: %v", err)), nil
	}

	// Validate domain name
	if err := utils.ValidateDomainName(input.Domain); err != nil {
		return types.NewErrorResult(fmt.Sprintf("invalid domain: %v", err)), nil
	}

	// Validate steps
	for i, step := range input.Steps {
		if step.Name == "" {
			return types.NewErrorResult(fmt.Sprintf("step %d: name is required", i+1)), nil
		}
		stepType := step.Type
		if stepType == "" {
			stepType = "form"
		}
		if stepType != "form" && stepType != "select" && stepType != "has_many" && stepType != "summary" {
			return types.NewErrorResult(fmt.Sprintf("step %d: invalid type '%s', must be form, select, has_many, or summary", i+1, stepType)), nil
		}
		if stepType == "has_many" && step.ChildDomain == "" {
			return types.NewErrorResult(fmt.Sprintf("step %d: child_domain is required for has_many steps", i+1)), nil
		}
		if step.HasManyMode != "" && step.HasManyMode != "select_existing" && step.HasManyMode != "create_inline" {
			return types.NewErrorResult(fmt.Sprintf("step %d: invalid has_many_mode '%s', must be select_existing or create_inline", i+1, step.HasManyMode)), nil
		}
	}

	// Get module path from go.mod
	modulePath, err := utils.GetModulePath(registry.WorkingDir)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to get module path: %v", err)), nil
	}

	// Create generator
	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Prepare template data
	data := generator.NewWizardData(input, modulePath)

	// Create directories
	pkgName := utils.ToPackageName(input.Domain)
	wizardName := utils.ToSnakeCase(input.WizardName)
	directories := []string{
		filepath.Join("internal", "web", pkgName),
		filepath.Join("internal", "web", pkgName, "views"),
	}

	// Add wizard draft directories if drafts are enabled
	if data.WithDrafts {
		directories = append(directories,
			filepath.Join("internal", "models"),
			filepath.Join("internal", "repository", "wizarddraft"),
			filepath.Join("internal", "services", "wizarddraft"),
		)
	}

	for _, dir := range directories {
		if err := gen.EnsureDir(dir); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to create directory %s: %v", dir, err)), nil
		}
	}

	// Generate wizard controller
	controllerPath := filepath.Join("internal", "web", pkgName, "wizard_"+wizardName+".go")
	if err := gen.GenerateFile("wizard/controller.go.tmpl", controllerPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate wizard controller: %v", err)), nil
	}

	// Generate main wizard view
	wizardViewPath := filepath.Join("internal", "web", pkgName, "views", "wizard_"+wizardName+".templ")
	if err := gen.GenerateFile("wizard/wizard_view.templ.tmpl", wizardViewPath, data); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to generate wizard view: %v", err)), nil
	}

	// Generate step views based on step types
	for i, step := range data.Steps {
		stepViewPath := filepath.Join("internal", "web", pkgName, "views", fmt.Sprintf("wizard_%s_step%d.templ", wizardName, i+1))
		var templatePath string
		switch step.Type {
		case "form":
			templatePath = "wizard/step_form.templ.tmpl"
		case "select":
			templatePath = "wizard/step_select.templ.tmpl"
		case "has_many":
			templatePath = "wizard/step_has_many.templ.tmpl"
		case "summary":
			templatePath = "wizard/step_summary.templ.tmpl"
		default:
			templatePath = "wizard/step_form.templ.tmpl"
		}

		// Create step-specific data
		stepData := struct {
			generator.WizardData
			Step generator.WizardStepData
		}{
			WizardData: data,
			Step:       step,
		}

		if err := gen.GenerateFile(templatePath, stepViewPath, stepData); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate step %d view: %v", i+1, err)), nil
		}
	}

	// Generate draft model, repo, and service if drafts are enabled
	if data.WithDrafts {
		draftData := generator.WizardDraftData{
			ModulePath: modulePath,
		}

		// Generate draft model
		draftModelPath := filepath.Join("internal", "models", "wizard_draft.go")
		if err := gen.GenerateFile("wizard/draft_model.go.tmpl", draftModelPath, draftData); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate draft model: %v", err)), nil
		}

		// Generate draft repository
		draftRepoPath := filepath.Join("internal", "repository", "wizarddraft", "wizarddraft.go")
		if err := gen.GenerateFile("wizard/draft_repository.go.tmpl", draftRepoPath, draftData); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate draft repository: %v", err)), nil
		}

		// Generate draft service
		draftServicePath := filepath.Join("internal", "services", "wizarddraft", "wizarddraft.go")
		if err := gen.GenerateFile("wizard/draft_service.go.tmpl", draftServicePath, draftData); err != nil {
			return types.NewErrorResult(fmt.Sprintf("failed to generate draft service: %v", err)), nil
		}
	}

	// Get result
	result := gen.Result()

	// Check for conflicts
	if conflictResult := CheckForConflicts(result); conflictResult != nil {
		return *conflictResult, nil
	}

	nextSteps := []string{
		"templ generate",
		"go mod tidy",
		fmt.Sprintf("Register wizard routes in cmd/web/main.go"),
		fmt.Sprintf("Add wizard link to domain views (e.g., a 'New with Wizard' button)"),
	}

	if data.WithDrafts {
		nextSteps = append(nextSteps, "Add WizardDraft to database AutoMigrate")
	}

	suggestedTools := []types.ToolHint{
		{
			Tool:        "scaffold_component",
			Description: "Create additional wizard UI components",
			Example:     `scaffold_component: { component_type: "wizard" }`,
			Priority:    "optional",
		},
	}

	if input.DryRun {
		return types.ScaffoldResult{
			Success:        true,
			Message:        fmt.Sprintf("Dry run: Would create wizard '%s' for domain '%s' with %d steps", input.WizardName, input.Domain, len(input.Steps)),
			FilesCreated:   result.FilesCreated,
			NextSteps:      nextSteps,
			SuggestedTools: suggestedTools,
		}, nil
	}

	return types.ScaffoldResult{
		Success:        true,
		Message:        fmt.Sprintf("Successfully created wizard '%s' for domain '%s' with %d steps", input.WizardName, input.Domain, len(input.Steps)),
		FilesCreated:   result.FilesCreated,
		FilesUpdated:   result.FilesUpdated,
		NextSteps:      nextSteps,
		SuggestedTools: suggestedTools,
	}, nil
}
