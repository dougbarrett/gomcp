# Analysis Findings: Wizard Template Issues

This document catalogs all issues identified by comparing wizard templates against the correct domain template patterns, organized by severity and target fix phase.

---

## Template Comparison Analysis

### Controller Template Comparison

**Files Compared:**
- Wizard: `internal/templates/wizard/controller.go.tmpl`
- Domain: `internal/templates/domain/controller.go.tmpl`

| Aspect | Wizard Controller | Domain Controller | Status |
|--------|-------------------|-------------------|--------|
| Response object | `resp` | `res` | Minor inconsistency |
| Redirect method | `resp.HXRedirect(url)` | `res.Redirect(url)` | **BUG: Method doesn't exist** |
| Render method | `resp.Component(view)` | `c.render(w, r, component)` | **BUG: Method doesn't exist** |
| CSRF token | `resp.CSRFToken()` | `middleware.GetCSRFToken(r.Context())` | **BUG: Method doesn't exist** |
| Middleware import | Missing | `"[[.ModulePath]]/internal/web/middleware"` | **BUG: Missing import** |
| Models import | Missing (when WithDrafts) | Imports as needed | **BUG: Missing import** |
| Service import | `[[.PackageName]]` | `[[.PackageName]]svc` with alias | Inconsistent but works |
| templ import | Missing | `"github.com/a-h/templ"` | Needed for render method |
| Local render method | Missing | `func (c *Controller) render(...)` | **BUG: Missing method** |

### Service Template Comparison

**Files Compared:**
- Wizard: `internal/templates/wizard/draft_service.go.tmpl`
- Domain: `internal/templates/domain/service.go.tmpl`

| Aspect | Wizard Draft Service | Domain Service | Status |
|--------|----------------------|----------------|--------|
| Repository import | Missing | `[[.PackageName]]repo "[[.ModulePath]]/internal/repository/[[.PackageName]]"` | **BUG: Missing import** |
| Repository type | `Repository` (undefined) | `[[.PackageName]]repo.Repository` | **BUG: Unqualified type** |
| Error definitions | None | `var Err[[.ModelName]]NotFound = errors.New(...)` | Enhancement |
| MCP injection markers | None | Has `// MCP:SERVICE_INTERFACE:START` etc. | Enhancement |

### Repository Template Comparison

**Files Compared:**
- Wizard: `internal/templates/wizard/draft_repository.go.tmpl`
- Domain: `internal/templates/domain/repository.go.tmpl`

| Aspect | Wizard Draft Repository | Domain Repository | Status |
|--------|-------------------------|-------------------|--------|
| QueryOption pattern | Missing | Full implementation with pagination/search/order | Enhancement |
| MCP injection markers | None | Has `// MCP:REPO_INTERFACE:START` etc. | Enhancement |
| Relationship preloading | N/A | Full preload support | N/A |
| FindAll method | Missing | Returns items + count | Enhancement |

---

## Issue Catalog

### Phase 2: Controller Template Fixes

#### Issue CTRL-001: Non-existent HXRedirect method

| Field | Value |
|-------|-------|
| **ID** | CTRL-001 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 63, 66, 89, 136, 181, 183, 188, 190, 248 |
| **Current Code** | `resp.HXRedirect(url)` |
| **Expected Code** | `resp.Redirect(url)` |
| **Fix Approach** | Replace all `HXRedirect` calls with `Redirect` |

#### Issue CTRL-002: Non-existent Component method

| Field | Value |
|-------|-------|
| **ID** | CTRL-002 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 136 |
| **Current Code** | `resp.Component(views.[[$.WizardNamePascal]]Step[[add $i 1]](props))` |
| **Expected Code** | `c.render(w, r, views.[[$.WizardNamePascal]]Step[[add $i 1]](props))` |
| **Fix Approach** | Add local `render` method to controller, use `c.render(w, r, component)` |

#### Issue CTRL-003: Non-existent CSRFToken method

| Field | Value |
|-------|-------|
| **ID** | CTRL-003 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 115 |
| **Current Code** | `CSRFToken: resp.CSRFToken(),` |
| **Expected Code** | `CSRFToken: middleware.GetCSRFToken(r.Context()),` |
| **Fix Approach** | Import middleware package, use `middleware.GetCSRFToken(r.Context())` |

#### Issue CTRL-004: Missing middleware import

| Field | Value |
|-------|-------|
| **ID** | CTRL-004 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 1-15 (import block) |
| **Current Code** | No middleware import |
| **Expected Code** | `"[[.ModulePath]]/internal/web/middleware"` |
| **Fix Approach** | Add middleware import to import block |

#### Issue CTRL-005: Missing models import (with drafts)

| Field | Value |
|-------|-------|
| **ID** | CTRL-005 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation when WithDrafts=true) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 1-15 (import block), 151 (usage) |
| **Current Code** | `var draft *models.WizardDraft` without import |
| **Expected Code** | Add `"[[.ModulePath]]/internal/models"` to imports |
| **Fix Approach** | Add conditional models import when WithDrafts is true |

#### Issue CTRL-006: Missing templ import

| Field | Value |
|-------|-------|
| **ID** | CTRL-006 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 1-15 (import block) |
| **Current Code** | No templ import |
| **Expected Code** | `"github.com/a-h/templ"` |
| **Fix Approach** | Add templ import for render method signature |

#### Issue CTRL-007: Missing local render method

| Field | Value |
|-------|-------|
| **ID** | CTRL-007 |
| **Bug Reference** | #a9479784 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | After RegisterRoutes (around line 48) |
| **Current Code** | No render method |
| **Expected Code** | `func (c *[[.WizardNamePascal]]WizardController) render(w http.ResponseWriter, r *http.Request, component templ.Component) { w.Header().Set("Content-Type", "text/html; charset=utf-8"); component.Render(r.Context(), w) }` |
| **Fix Approach** | Add render method matching domain controller pattern |

#### Issue CTRL-008: Potential URL double slashes

| Field | Value |
|-------|-------|
| **ID** | CTRL-008 |
| **Bug Reference** | #a9479784 |
| **Severity** | Medium (incorrect behavior if URLPath has trailing slash) |
| **Location** | `internal/templates/wizard/controller.go.tmpl` |
| **Lines** | 63, 66, 89, 181, 183, 188, 190, 248 |
| **Current Code** | `fmt.Sprintf("[[.URLPath]]/wizard/...")` |
| **Expected Code** | Use `strings.TrimRight(url, "/")` or ensure URLPath never has trailing slash |
| **Fix Approach** | Add URL path normalization or verify generator always produces clean paths |

---

### Phase 3: Service Template Fixes

#### Issue SVC-001: Missing repository import

| Field | Value |
|-------|-------|
| **ID** | SVC-001 |
| **Bug Reference** | #cb94adf6 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/draft_service.go.tmpl` |
| **Lines** | 1-8 (import block), 22 (usage) |
| **Current Code** | `repo Repository` without import |
| **Expected Code** | Add `wizarddraftrepo "[[.ModulePath]]/internal/repository/wizarddraft"` to imports, use `wizarddraftrepo.Repository` |
| **Fix Approach** | Add repository import with alias, qualify Repository type |

#### Issue SVC-002: Unqualified Repository type

| Field | Value |
|-------|-------|
| **ID** | SVC-002 |
| **Bug Reference** | #cb94adf6 |
| **Severity** | Critical (blocks compilation) |
| **Location** | `internal/templates/wizard/draft_service.go.tmpl` |
| **Lines** | 22, 26 |
| **Current Code** | `repo Repository` and `NewService(repo Repository)` |
| **Expected Code** | `repo wizarddraftrepo.Repository` and `NewService(repo wizarddraftrepo.Repository)` |
| **Fix Approach** | Qualify Repository type with package alias |

---

### Phase 4: View Template Improvements

#### Issue VIEW-001: Inconsistent component usage patterns

| Field | Value |
|-------|-------|
| **ID** | VIEW-001 |
| **Severity** | Low (inconsistency) |
| **Location** | `internal/templates/wizard/step_*.templ.tmpl` |
| **Current Code** | Various wizard-specific component patterns |
| **Expected Code** | Align with scaffold_view patterns |
| **Fix Approach** | Review and align component usage with scaffold_view templates |

---

### Phase 5: Draft System Review

#### Issue DRAFT-001: No MCP injection markers in draft templates

| Field | Value |
|-------|-------|
| **ID** | DRAFT-001 |
| **Severity** | Low (enhancement) |
| **Location** | `draft_service.go.tmpl`, `draft_repository.go.tmpl` |
| **Current Code** | No injection markers |
| **Expected Code** | Add `// MCP:*:START` and `// MCP:*:END` markers |
| **Fix Approach** | Add injection markers for extensibility via extend_* tools |

---

### Phase 6: Generator Logic Review

No issues identified in generator logic from template analysis. Review needed during Phase 6 execution.

---

## Issue Summary by Phase

| Phase | Issues | Critical | High | Medium | Low |
|-------|--------|----------|------|--------|-----|
| Phase 2: Controller | 8 | 7 | 0 | 1 | 0 |
| Phase 3: Service | 2 | 2 | 0 | 0 | 0 |
| Phase 4: Views | 1 | 0 | 0 | 0 | 1 |
| Phase 5: Drafts | 1 | 0 | 0 | 0 | 1 |
| Phase 6: Generator | 0 | 0 | 0 | 0 | 0 |
| **Total** | **12** | **9** | **0** | **1** | **2** |

---

## Bug Report Verification

### Bug #a9479784: Controller uses non-existent web.Response methods

**Status:** VERIFIED

**Details:**
- `resp.HXRedirect()` - Does not exist, should be `resp.Redirect()`
- `resp.Component()` - Does not exist, should use `c.render(w, r, component)`
- `resp.CSRFToken()` - Does not exist, should be `middleware.GetCSRFToken(r.Context())`

**Additional Issues Found:**
- Missing `middleware` import
- Missing `models` import (when WithDrafts=true)
- Missing `templ` import
- Missing local `render` method on controller
- Potential URL double slash issues

### Bug #cb94adf6: Service missing repository import

**Status:** VERIFIED

**Details:**
- `draft_service.go.tmpl` uses `Repository` type without importing the repository package
- Service is in `internal/services/wizarddraft` package
- Repository is in `internal/repository/wizarddraft` package
- These are separate packages, so import is required

**Fix Required:**
- Import `wizarddraftrepo "[[.ModulePath]]/internal/repository/wizarddraft"`
- Use qualified type `wizarddraftrepo.Repository`

---

## Correct Pattern Reference

### Domain Controller Pattern (controller.go.tmpl)

```go
package [[.PackageName]]

import (
    "net/http"
    "strconv"

    [[.PackageName]]svc "[[.ModulePath]]/internal/services/[[.PackageName]]"
    "[[.ModulePath]]/internal/web"
    "[[.ModulePath]]/internal/web/[[.PackageName]]/views"
    "[[.ModulePath]]/internal/web/middleware"  // CSRF token
    "github.com/a-h/templ"                     // For render method
    "github.com/go-chi/chi/v5"
)

type Controller struct {
    service [[.PackageName]]svc.Service
}

// render renders a templ component to the response.
func (c *Controller) render(w http.ResponseWriter, r *http.Request, component templ.Component) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    component.Render(r.Context(), w)
}

// Example handler using correct patterns:
func (c *Controller) Show(w http.ResponseWriter, r *http.Request) {
    res := web.NewResponse(w, r)

    // CSRF token via middleware
    csrfToken := middleware.GetCSRFToken(r.Context())

    // Redirect via Response
    res.Redirect("/some/path")

    // Render via controller method
    c.render(w, r, views.SomeView(props))
}
```

### Domain Service Pattern (service.go.tmpl)

```go
package [[.PackageName]]

import (
    "context"
    "errors"

    "[[.ModulePath]]/internal/models"
    [[.PackageName]]repo "[[.ModulePath]]/internal/repository/[[.PackageName]]"  // Aliased import
)

type service struct {
    repo [[.PackageName]]repo.Repository  // Qualified type
}

func NewService(repo [[.PackageName]]repo.Repository) Service {
    return &service{repo: repo}
}
```

---

*Document created: Phase 1, Plan 2 - Compare and Document Issues*
*Issues verified against actual template files and web.Response implementation*
