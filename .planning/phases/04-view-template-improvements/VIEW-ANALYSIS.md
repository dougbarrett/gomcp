# Wizard View Template Pattern Analysis

## Overview

This analysis compares wizard view templates against domain view templates (scaffold_view) to identify consistency opportunities. The goal is to ensure wizard templates follow established patterns while respecting their distinct purpose (multi-step create flow vs CRUD operations).

## Component Usage Comparison

### Domain View Templates

| Component | form.templ.tmpl | list.templ.tmpl |
|-----------|----------------|-----------------|
| `components.Button` | ✓ | ✓ |
| `components.Card` | ✓ | ✓ |
| `components.CardHeader` | ✓ | ✓ |
| `components.CardContent` | ✓ | ✓ |
| `components.CardFooter` | - | ✓ |
| `components.Icon` | ✓ | ✓ |
| `components.Input` | ✓ | - |
| `components.Textarea` | ✓ | - |
| `components.Select` | ✓ | - |
| `components.Checkbox` | ✓ | - |
| `components.Label` | ✓ | - |
| `components.FormError` | ✓ | - |
| `components.Pagination` | - | ✓ |
| `components.ModalContainer` | - | ✓ |

### Wizard View Templates

| Component | wizard_view | step_form | step_select | step_has_many | step_summary |
|-----------|-------------|-----------|-------------|---------------|--------------|
| `components.Wizard` | ✓ | - | - | - | - |
| `components.WizardCard` | ✓ | - | - | - | - |
| `components.WizardSteps` | ✓ | - | - | - | - |
| `components.WizardStepContent` | ✓ | - | - | - | - |
| `components.WizardNav` | - | ✓ | ✓ | ✓ | ✓ |
| `components.WizardEmpty` | - | - | ✓ | - | - |
| `components.WizardSummary` | - | - | - | - | ✓ |
| `components.Label` | - | ✓ | - | - | ✓ |
| `components.Input` | - | ✓ | ✓ | ✓ | ✓ |
| `components.FormError` | - | ✓ | ✓ | ✓ | ✓ |
| `components.Icon` | - | - | - | ✓ | - |
| `components.Card` | - | - | - | - | ✓ |
| `components.CardContent` | - | - | - | - | ✓ |
| `components.SummaryItem` | - | - | - | - | ✓ (type) |

### Observation

Both template sets use appropriate components for their use cases. Wizard templates correctly use wizard-specific components (`WizardNav`, `WizardSteps`, `WizardEmpty`, `WizardSummary`) while sharing common form components (`Input`, `Label`, `FormError`).

## Props Pattern Comparison

### Domain Views (form.templ.tmpl)

```go
type ModelNameFormProps struct {
    Item       *models.ModelName
    Errors     map[string]string
    IsEdit     bool
    CSRFToken  string
    BasePath   string  // URL flexibility
    // Relationship options
}

func (p ModelNameFormProps) getBasePath() string {
    if p.BasePath != "" {
        return p.BasePath
    }
    return "/default-path"
}
```

**Key patterns:**
- `BasePath` field with `getBasePath()` helper for URL flexibility
- Supports both create and edit modes via `IsEdit`
- Relationship options passed for belongs_to dropdowns

### Wizard Views (step templates)

```go
type WizardNameStep1Props struct {
    CurrentStep int
    TotalSteps  int
    CSRFToken   string
    Errors      map[string]string
    DraftID     string              // Draft support
    StepData    map[string]interface{}  // Pre-filled values
}
```

**Key patterns:**
- No `BasePath`/`getBasePath()` - URLs constructed directly with template variables
- Step-specific props (CurrentStep, TotalSteps)
- Draft persistence support (DraftID, StepData)

### Assessment

The difference is **intentional and appropriate**:
- Domain views need URL flexibility for admin/public path variations
- Wizards have fixed URL patterns (`/wizard/name/step/N`) that don't need runtime configuration
- Adding `getBasePath()` to wizards would add complexity with no benefit

**Status: No change needed**

## HTMX Pattern Comparison

### Domain Views

```html
hx-post={ props.getBasePath() }
hx-put={ fmt.Sprintf("%s/%d", props.getBasePath(), props.Item.ID) }
hx-target="#main-content"
hx-swap="innerHTML"
```

### Wizard Views

```html
hx-post={ fmt.Sprintf("[[.WizardData.URLPath]]/wizard/[[.WizardData.WizardName]]/step/[[.Step.Number]]") }
hx-target="#main-content"
hx-swap="innerHTML"
```

**Assessment:** Both use consistent HTMX patterns (`hx-target="#main-content"`, `hx-swap="innerHTML"`). URL construction differs due to different URL patterns but follows the same approach.

**Status: No change needed**

## Empty State Pattern Comparison

### Domain Views

```go
// Dedicated empty state component
templ ModelNameEmptyState(basePath string) {
    <div class="text-center py-12">
        // Icon, message, action button
    </div>
}
```

### Wizard Views

| Template | Empty State Pattern |
|----------|---------------------|
| step_select.templ.tmpl | ✓ Uses `@components.WizardEmpty("message", "description")` |
| step_has_many.templ.tmpl | ⚠ Uses inline `<p class="text-sm text-gray-500...">No items available.</p>` |

### Inconsistency Found

`step_has_many.templ.tmpl` uses an inline empty message (line 86-89) instead of `@components.WizardEmpty`:

```go
// Current (inline)
<p class="text-sm text-gray-500 dark:text-gray-400 p-4 text-center">
    No items available.
</p>

// Could use (component - matches step_select.templ.tmpl)
@components.WizardEmpty("No items available", "Add items first.")
```

**Assessment:** Minor inconsistency. Both approaches work correctly. The inline approach is simpler and the styling is consistent with the component. This is **optional to fix**.

**Status: Optional improvement (low priority)**

## TODO Comment Analysis

### Current TODOs in Templates

| Template | Line | TODO Content |
|----------|------|--------------|
| step_select.templ.tmpl | 21-22 | `// TODO: Change this to the appropriate model type.` |
| step_select.templ.tmpl | 104-105 | `<!-- TODO: Customize display based on your model fields -->` |
| step_has_many.templ.tmpl | 97-98 | `<!-- TODO: Customize display based on your model fields -->` |
| step_has_many.templ.tmpl | 184-186 | `<!-- TODO: Add form fields for the child model -->` |
| step_summary.templ.tmpl | 131-132 | `// TODO: Customize this based on your wizard's fields` |

### Assessment

TODOs are clear and correctly placed. They indicate customization points for:
1. Model type changes (related model selection)
2. Display field customization
3. Form field addition for inline creation

**Status: No change needed**

## Identified Improvements

### Worth Fixing

None identified. All templates are functional and follow appropriate patterns for their use cases.

### Optional (Low Priority)

1. **step_has_many.templ.tmpl empty state**: Could use `@components.WizardEmpty` for consistency with step_select.templ.tmpl
   - Location: `step_has_many.templ.tmpl` lines 86-89
   - Current: Inline `<p>` element
   - Proposed: `@components.WizardEmpty("No items available", "Please add some items first.")`
   - Impact: Cosmetic consistency only, no functional difference

### Skip (Intentional Differences)

| Pattern | Domain Views | Wizard Views | Reason to Skip |
|---------|--------------|--------------|----------------|
| BasePath/getBasePath() | Yes | No | Wizards use fixed URL patterns, no flexibility needed |
| IsEdit mode | Yes | No | Wizards are create-only flows |
| Modal vs page form style | Yes | No | Wizards always use page-style multi-step layout |

## Conclusion

**Overall Assessment: Wizard view templates are well-designed and consistent with scaffold_view patterns where appropriate.**

The wizard templates correctly:
- Use wizard-specific components (WizardNav, WizardSteps, WizardSummary)
- Share common form components (Input, Label, FormError) with domain views
- Follow the same HTMX patterns (hx-target, hx-swap)
- Handle empty states appropriately
- Include clear TODO comments for customization points

The differences from domain views are intentional:
- Multi-step flow vs single-page CRUD
- Fixed wizard URL patterns vs flexible base paths
- Draft persistence support vs immediate save
- Create-only vs create/edit modes

**Recommendation: No template changes required.** The minor empty state inconsistency in step_has_many.templ.tmpl is cosmetic and doesn't affect functionality. Making that change would introduce risk with minimal benefit.

---

*Analysis completed: 2026-01-07*
*Phase: 04-view-template-improvements*
*Plan: 01*
