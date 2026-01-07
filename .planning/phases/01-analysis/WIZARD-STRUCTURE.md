# Wizard Template Structure Analysis

This document maps the wizard scaffolding system: template files, data requirements, generated output, and generator logic.

---

## Template Inventory

### 1. Controller Template

**File:** `internal/templates/wizard/controller.go.tmpl`
**Purpose:** Generates HTTP handlers for the wizard flow
**Output:** `internal/web/{domain}/{wizard_name}_wizard.go`

**Responsibilities:**
- Route registration (`/wizard/{wizard_name}/...`)
- Start handler (creates draft if enabled, redirects to step 1)
- Resume handler (loads existing draft)
- Per-step GET handlers (render step views)
- Per-step POST handlers (process form submissions)
- Submit handler (creates final entity, deletes draft)

**Key Patterns:**
- Uses `web.NewResponse(w, r)` for responses
- Calls `resp.HXRedirect()`, `resp.Component()`, `resp.CSRFToken()`, `resp.Error()`
- Injects domain service and optional draft service

---

### 2. Draft Service Template

**File:** `internal/templates/wizard/draft_service.go.tmpl`
**Purpose:** Business logic for wizard draft persistence
**Output:** `internal/services/wizarddraft/wizarddraft.go`

**Responsibilities:**
- CreateOrUpdate: Create new or update existing draft
- GetByID: Retrieve draft by ID
- GetStepData: Decode JSON step data from draft
- Delete: Remove completed draft
- FindExisting: Find draft by wizard name + user
- CleanupExpired: Remove old drafts

**Key Patterns:**
- Uses `Repository` interface for data access
- Stores step data as JSON-encoded string
- Tracks current step number for resume

---

### 3. Draft Repository Template

**File:** `internal/templates/wizard/draft_repository.go.tmpl`
**Purpose:** GORM-based data access for wizard drafts
**Output:** `internal/repository/wizarddraft/wizarddraft.go`

**Responsibilities:**
- CRUD operations for WizardDraft model
- FindByWizardAndUser query
- DeleteExpired cleanup query

**Key Patterns:**
- Uses `gorm.DB` for database operations
- Supports soft delete via GORM
- Handles nullable UserID for anonymous drafts

---

### 4. Draft Model Template

**File:** `internal/templates/wizard/draft_model.go.tmpl`
**Purpose:** GORM model for wizard draft table
**Output:** `internal/models/wizard_draft.go`

**Fields:**
- ID (uint, primary key)
- WizardName (string, indexed)
- Domain (string, indexed)
- UserID (*uint, nullable, indexed)
- StepData (string, text/JSON)
- CurrentStep (int)
- CreatedAt, UpdatedAt, DeletedAt (timestamps)

---

### 5. Wizard View Template

**File:** `internal/templates/wizard/wizard_view.templ.tmpl`
**Purpose:** Common wizard layout and step navigation
**Output:** `internal/web/{domain}/views/wizard.templ`

**Exports:**
- `{WizardName}WizardProps` struct
- `{WizardName}Steps()` function (returns step configuration)
- `{WizardName}WizardLayout` templ component
- `{WizardName}WizardCard` templ component

**Key Patterns:**
- Wraps step content with WizardSteps progress indicator
- Uses components: Wizard, WizardSteps, WizardStepContent, WizardCard
- Tracks step status (completed, active, pending)

---

### 6. Step Form Template

**File:** `internal/templates/wizard/step_form.templ.tmpl`
**Purpose:** Form step for collecting user input fields
**Output:** `internal/web/{domain}/views/step_{N}.templ`

**Generates:**
- Props struct with CurrentStep, TotalSteps, CSRFToken, Errors, (DraftID, StepData if drafts)
- Step templ component rendering form fields

**Key Patterns:**
- Uses HTMX for form submission (`hx-post`, `hx-target="#main-content"`)
- Iterates over `.Step.FieldNames` to render inputs
- Uses components: Label, Input, FormError, WizardNav
- Pre-fills values from StepData when drafts enabled

---

### 7. Step Select Template

**File:** `internal/templates/wizard/step_select.templ.tmpl`
**Purpose:** Selection step for choosing from a list of options
**Output:** `internal/web/{domain}/views/step_{N}.templ`

**Generates:**
- Props struct with Options (model slice), SearchQuery (if searchable)
- Radio button list for selection
- Optional search input with HTMX filtering

**Key Patterns:**
- References `models.{ModelName}` for Options type
- Uses HTMX search with `hx-trigger="keyup changed delay:300ms"`
- Stores selection in `selected_id` field

---

### 8. Step Has Many Template

**File:** `internal/templates/wizard/step_has_many.templ.tmpl`
**Purpose:** Has-many step for selecting/creating multiple related items
**Output:** `internal/web/{domain}/views/step_{N}.templ`

**Modes:**
- `select_existing`: Select items from list with quantities
- `create_inline`: Add new items inline (template placeholder)

**Generates:**
- Props struct with AvailableItems, SelectedItems (custom struct)
- Two-column layout (available | selected)
- Uses Hyperscript for client-side item management

**Key Patterns:**
- References `models.{ChildModelName}` for AvailableItems
- Uses `{WizardName}SelectedItem` struct for selected items with quantity
- Serializes selections to hidden input `selected_items`

---

### 9. Step Summary Template

**File:** `internal/templates/wizard/step_summary.templ.tmpl`
**Purpose:** Review step before final submission
**Output:** `internal/web/{domain}/views/step_{N}.templ`

**Generates:**
- Props struct with SummaryItems (components.SummaryItem slice)
- Summary display with edit links back to previous steps
- Optional additional fields for final input

**Key Patterns:**
- Uses `components.WizardSummary` for display
- Generates `buildSummaryFromStepData()` helper when drafts enabled
- Posts to `/wizard/{wizard_name}/submit` endpoint

---

## Template Data Requirements

### Common Data (WizardData)

All templates receive `WizardData` struct containing:

| Field | Type | Description |
|-------|------|-------------|
| `.WizardName` | string | Lowercase wizard name (e.g., "create_order") |
| `.WizardNamePascal` | string | PascalCase (e.g., "CreateOrder") |
| `.PackageName` | string | Go package name for domain |
| `.ModulePath` | string | Go module path (e.g., "github.com/user/app") |
| `.Domain` | string | Domain name for the wizard target |
| `.ModelName` | string | Target model name (PascalCase) |
| `.URLPath` | string | Base URL path (e.g., "/orders") |
| `.Steps` | []Step | Array of step configurations |
| `.TotalSteps` | int | Count of steps |
| `.WithDrafts` | bool | Enable draft persistence feature |
| `.SuccessRedirect` | string | Redirect URL after successful completion |

### Step Data

Each step in `.Steps` contains:

| Field | Type | Description |
|-------|------|-------------|
| `.Name` | string | Human-readable step name |
| `.Type` | string | Step type: "form", "select", "has_many", "summary" |
| `.Number` | int | Step number (1-indexed) |
| `.IsFirst` | bool | True if first step |
| `.IsLast` | bool | True if last step |
| `.FieldNames` | []string | Form field names for this step |
| `.Searchable` | bool | Enable search for select/has_many steps |
| `.HasManyMode` | string | "select_existing" or "create_inline" |
| `.ChildModelName` | string | Model name for has_many child items |

### Conditional Sections

Templates use Go template conditionals for optional features:

```go
[[- if .WithDrafts]]
// Draft-related code
[[- end]]

[[- range $i, $step := .Steps]]
// Per-step code generation
[[- end]]

[[- if eq .Type "select"]]
// Select-specific code
[[- end]]
```

---

## Generated Output Structure

For a wizard named "create_order" in the "order" domain with drafts enabled:

```
internal/
├── models/
│   └── wizard_draft.go           # WizardDraft model (shared)
│
├── repository/
│   └── wizarddraft/
│       └── wizarddraft.go        # Draft repository (shared)
│
├── services/
│   └── wizarddraft/
│       └── wizarddraft.go        # Draft service (shared)
│
└── web/
    └── order/
        ├── create_order_wizard.go     # Controller
        └── views/
            ├── wizard.templ           # Layout + step navigation
            ├── step_1.templ           # First step view
            ├── step_2.templ           # Second step view
            └── step_N.templ           # Additional step views
```

**Template -> Output Mapping:**

| Template | Output Path | When Generated |
|----------|-------------|----------------|
| controller.go.tmpl | `internal/web/{domain}/{wizard}_wizard.go` | Always |
| wizard_view.templ.tmpl | `internal/web/{domain}/views/wizard.templ` | Always |
| step_form.templ.tmpl | `internal/web/{domain}/views/step_{N}.templ` | When step.Type == "form" |
| step_select.templ.tmpl | `internal/web/{domain}/views/step_{N}.templ` | When step.Type == "select" |
| step_has_many.templ.tmpl | `internal/web/{domain}/views/step_{N}.templ` | When step.Type == "has_many" |
| step_summary.templ.tmpl | `internal/web/{domain}/views/step_{N}.templ` | When step.Type == "summary" |
| draft_model.go.tmpl | `internal/models/wizard_draft.go` | When WithDrafts == true |
| draft_repository.go.tmpl | `internal/repository/wizarddraft/wizarddraft.go` | When WithDrafts == true |
| draft_service.go.tmpl | `internal/services/wizarddraft/wizarddraft.go` | When WithDrafts == true |

---

## Observations and Issues Found

### Issue 1: Missing `models` Import in Controller

**Location:** `controller.go.tmpl` lines 151-155 (inside Step{N}Submit with drafts)
**Problem:** Uses `models.WizardDraft` type but the `models` package is not imported
**Impact:** Generated code will not compile

```go
[[- if $.WithDrafts]]
// Get or create draft
draftID := r.FormValue("draft_id")
var draft *models.WizardDraft  // <-- models not imported!
```

### Issue 2: web.Response Method Calls

**Location:** `controller.go.tmpl` throughout
**Problem:** Uses methods that may not exist on `web.Response`:
- `resp.HXRedirect(url)` - May need to be `resp.Redirect(url)` or similar
- `resp.Component(view)` - May need to be `resp.Render(view)` or similar
- `resp.CSRFToken()` - Needs verification against actual web.Response implementation

**Impact:** This is Bug #a9479784 - non-existent web.Response methods

### Issue 3: Potential URL Double Slashes

**Location:** `controller.go.tmpl` various URL constructions
**Examples:**
```go
resp.HXRedirect(fmt.Sprintf("[[.URLPath]]/wizard/[[.WizardName]]/%d/step/1", draft.ID))
```
If `URLPath` is `/orders`, this generates `/orders/wizard/create_order/...` which is correct.
But if `URLPath` is `/orders/`, this would generate `/orders//wizard/...` with double slash.

**Impact:** Part of Bug #a9479784 - URL construction issues

### Issue 4: Inconsistent Template Delimiters

All templates use `[[` and `]]` as delimiters (non-standard), which is intentional to avoid conflicts with Go template syntax in generated `.templ` files.

---

## Template Function Usage

Templates use these template functions:

| Function | Usage | Description |
|----------|-------|-------------|
| `add` | `[[add $i 1]]` | Adds numbers (used for 1-indexed step numbers) |
| `sub` | `[[sub .Step.Number 1]]` | Subtracts numbers (used for previous step) |
| `toLower` | `[[.ModelName \| toLower]]` | Lowercase string |
| `toLabel` | `[[. \| toLabel]]` | Convert field name to human label |
| `empty` | `[[- if empty .Step.FieldNames]]` | Check if slice is empty |

---

## Dependencies Between Templates

```
                    ┌─────────────────────┐
                    │  controller.go.tmpl │
                    └──────────┬──────────┘
                               │ imports
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ▼                      ▼                      ▼
┌───────────────┐    ┌─────────────────┐    ┌────────────────┐
│ domain/service│    │ wizarddraft/    │    │ views/         │
│ (existing)    │    │ service         │    │ step_*.templ   │
└───────────────┘    └────────┬────────┘    └────────┬───────┘
                              │                      │
                              │ depends on           │ uses
                              ▼                      ▼
                     ┌─────────────────┐    ┌────────────────┐
                     │ wizarddraft/    │    │ wizard.templ   │
                     │ repository      │    │ (layout)       │
                     └────────┬────────┘    └────────────────┘
                              │
                              │ depends on
                              ▼
                     ┌─────────────────┐
                     │ models/         │
                     │ wizard_draft.go │
                     └─────────────────┘
```

---

## Generator Logic

### Source Files

- **Generator:** `internal/tools/scaffold_wizard.go`
- **Data Types:** `internal/generator/data.go`
- **Input Types:** `internal/types/inputs.go`

### Input Structure

**ScaffoldWizardInput:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `wizard_name` | string | (required) | Wizard identifier (e.g., "create_order") |
| `domain` | string | (required) | Target domain (e.g., "order") |
| `steps` | []WizardStepDef | (required) | Step definitions |
| `layout` | string | "dashboard" | View layout: dashboard, base, auth, none |
| `route_group` | string | "public" | Middleware context: public, authenticated, admin |
| `form_style` | string | "page" | How steps are displayed: page, modal |
| `success_redirect` | string | "/{domain}" | Redirect URL after completion |
| `with_drafts` | *bool | true | Enable draft persistence |
| `dry_run` | bool | false | Preview without writing files |

**WizardStepDef:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | (required) | Step display name |
| `type` | string | "form" | Step type: form, select, has_many, summary |
| `fields` | []string | [] | Field names for this step |
| `child_domain` | string | "" | Related domain for has_many steps |
| `has_many_mode` | string | "select_existing" | How items are added: select_existing, create_inline |
| `searchable` | bool | false | Enable search for select/has_many steps |
| `validation_rules` | map | {} | Per-field validation rules |

### Input Validation

The generator validates input in this order:

1. **Required Fields:**
   - `wizard_name` must be non-empty
   - `domain` must be non-empty
   - `steps` must have at least one step

2. **Name Validation:**
   - `wizard_name` validated via `utils.ValidateComponentName()`
   - `domain` validated via `utils.ValidateDomainName()`

3. **Step Validation:**
   - Each step must have a `name`
   - Step `type` must be: form, select, has_many, or summary (defaults to form)
   - `has_many` steps require `child_domain`
   - `has_many_mode` must be: select_existing or create_inline

4. **Module Path:**
   - Retrieved from go.mod via `utils.GetModulePath()`

### Data Transformation

**Input → WizardData transformation** (`generator.NewWizardData`):

```
ScaffoldWizardInput                 WizardData
─────────────────────────────────────────────────────
wizard_name        →  WizardName, WizardNamePascal
domain             →  Domain, ModelName, PackageName, VariableName
                      URLPath (/domain), URLPathSegment (domain)
steps[]            →  Steps[] (with Number, IsFirst, IsLast flags)
layout             →  Layout (default: "dashboard")
route_group        →  RouteGroup (default: "public")
form_style         →  FormStyle (default: "page")
success_redirect   →  SuccessRedirect (default: /{domain})
with_drafts        →  WithDrafts (default: true)
```

**Step Processing:**
- Each step gets `Number` (1-indexed), `IsFirst`, `IsLast` flags
- Step type defaults to "form" if empty
- HasManyMode defaults to "select_existing" if empty
- ChildModelName computed from ChildDomain via `utils.ToModelName()`

**Feature Flags (derived from steps):**
- `HasSelectSteps`: true if any step is type "select"
- `HasHasManySteps`: true if any step is type "has_many"
- `HasSummaryStep`: true if any step is type "summary"
- `HasFormSteps`: true if any step is type "form"

### Generation Flow

1. **Create Directories:**
   ```
   internal/web/{pkgName}/
   internal/web/{pkgName}/views/
   ```
   Plus if WithDrafts:
   ```
   internal/models/
   internal/repository/wizarddraft/
   internal/services/wizarddraft/
   ```

2. **Generate Controller:**
   ```
   wizard/controller.go.tmpl → internal/web/{pkgName}/wizard_{wizardName}.go
   ```

3. **Generate Wizard View:**
   ```
   wizard/wizard_view.templ.tmpl → internal/web/{pkgName}/views/wizard_{wizardName}.templ
   ```

4. **Generate Step Views (per step):**
   ```
   For each step:
     - Determine template by step.Type (form/select/has_many/summary)
     - Create stepData = { WizardData, Step }
     - Generate: internal/web/{pkgName}/views/wizard_{wizardName}_step{N}.templ
   ```

5. **Generate Draft Infrastructure (if WithDrafts):**
   ```
   wizard/draft_model.go.tmpl      → internal/models/wizard_draft.go
   wizard/draft_repository.go.tmpl → internal/repository/wizarddraft/wizarddraft.go
   wizard/draft_service.go.tmpl    → internal/services/wizarddraft/wizarddraft.go
   ```

6. **Check for Conflicts:**
   - Calls `CheckForConflicts(result)` to detect overwrites

7. **Return Result:**
   - Success message with file list
   - Next steps for manual wiring

### Post-Generation (Manual Steps)

The generator outputs these next steps:

1. `templ generate` - Compile templ files
2. `go mod tidy` - Update dependencies
3. Register wizard routes in `cmd/web/main.go`
4. Add wizard link to domain views (e.g., "New with Wizard" button)
5. (If drafts) Add `WizardDraft` to database AutoMigrate

**Not Automated:**
- Route registration in main.go (no DI wiring like scaffold_domain)
- Link from existing domain views to wizard
- Database migration for WizardDraft model

### Comparison to scaffold_domain

| Aspect | scaffold_wizard | scaffold_domain |
|--------|-----------------|-----------------|
| DI Wiring | Manual | Automatic via update_di_wiring |
| Controller | Wizard-specific handlers | CRUD handlers |
| Views | Step-based views | List/Show/Form views |
| Routes | `/wizard/{name}/...` | `/{domain}/...` |
| Draft Support | Built-in | N/A |
| Service | Uses existing domain service | Generates new service |

---

## Summary: Key Issues for Phase 2+

Based on this analysis, the following issues need attention in subsequent phases:

### Bug #a9479784 (Controller Template)

1. **Missing `models` import** - Line 151 uses `models.WizardDraft` without import
2. **Non-existent `web.Response` methods:**
   - `resp.HXRedirect()` - should be `resp.Redirect()` or similar
   - `resp.Component()` - should be `resp.Render()` or similar
   - `resp.CSRFToken()` - needs verification
3. **Potential URL double slashes** - URLPath concatenation

### Bug #cb94adf6 (Service Template)

Need to verify during Phase 3 analysis - draft_service.go.tmpl appears correct, but the bug report mentions missing repository import.

### Patterns to Align With

- Check `scaffold_controller` templates for correct `web.Response` method usage
- Check `scaffold_domain` for proper import patterns
- Verify URL path construction patterns

---

*Document created: Phase 1, Plan 1 - Wizard Structure Analysis*
*Updated with Generator Logic: Task 2*
