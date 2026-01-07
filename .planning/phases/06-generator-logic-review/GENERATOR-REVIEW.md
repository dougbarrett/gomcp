# Generator Logic Review: scaffold_wizard.go vs scaffold_domain.go

**Review Date:** 2026-01-07
**Files Analyzed:**
- `internal/tools/scaffold_wizard.go` (287 lines)
- `internal/tools/scaffold_domain.go` (558 lines)
- `internal/modifier/inject.go` (405 lines)
- `internal/metadata/metadata.go` (186 lines)

## Executive Summary

scaffold_wizard.go is missing 6 key features that scaffold_domain.go has, resulting in:
- **Manual DI wiring required** - Users must manually register wizard routes
- **Manual AutoMigrate required** - Users must manually add WizardDraft model
- **No metadata tracking** - Cannot sync/upgrade wizards in future

## Feature Comparison Table

| Feature | scaffold_domain.go | scaffold_wizard.go | Status |
|---------|-------------------|-------------------|--------|
| modifier import | ✅ Line 10 | ❌ | Missing |
| metadata import | ✅ Line 9 | ❌ | Missing |
| ScaffolderVersion constant | ✅ Lines 16-18 | ❌ | Missing |
| DI wiring injection | ✅ Lines 293-317 | ❌ | Manual step |
| database.go AutoMigrate | ✅ Lines 522-537 | ❌ | Manual step |
| Metadata tracking | ✅ Lines 357-364 | ❌ | Missing |
| Input validation | ✅ Lines 142-192 | ✅ Lines 96-134 | OK |
| Dry run support | ✅ Lines 347-354 | ✅ Lines 268-276 | OK |
| Conflict checking | ✅ Lines 287-290 | ✅ Lines 244-246 | OK |

---

## Issue Catalog

### GEN-001: Missing modifier import

**Severity:** Critical
**Location:** `internal/tools/scaffold_wizard.go` lines 4-12

**Current imports (scaffold_wizard.go):**
```go
import (
    "context"
    "fmt"
    "path/filepath"

    "github.com/dbb1dev/go-mcp/internal/generator"
    "github.com/dbb1dev/go-mcp/internal/types"
    "github.com/dbb1dev/go-mcp/internal/utils"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)
```

**Expected imports (from scaffold_domain.go lines 4-14):**
```go
import (
    "context"
    "fmt"
    "path/filepath"

    "github.com/dbb1dev/go-mcp/internal/generator"
    "github.com/dbb1dev/go-mcp/internal/metadata"
    "github.com/dbb1dev/go-mcp/internal/modifier"
    "github.com/dbb1dev/go-mcp/internal/types"
    "github.com/dbb1dev/go-mcp/internal/utils"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)
```

**Fix Approach:**
Add missing imports to scaffold_wizard.go:
```go
"github.com/dbb1dev/go-mcp/internal/modifier"
```

---

### GEN-002: Missing metadata import

**Severity:** Medium
**Location:** `internal/tools/scaffold_wizard.go` lines 4-12

**Issue:** No metadata package import

**Fix Approach:**
Add missing import:
```go
"github.com/dbb1dev/go-mcp/internal/metadata"
```

---

### GEN-003: No version tracking constant

**Severity:** Low
**Location:** `internal/tools/scaffold_wizard.go` - missing entirely

**scaffold_domain.go (lines 16-18):**
```go
// ScaffolderVersion is the current version of the scaffolding tools.
// Used for tracking which version generated the code for future upgrades.
const ScaffolderVersion = "0.1.0"
```

**scaffold_wizard.go:** No version constant

**Fix Approach:**
Option A: Reuse ScaffolderVersion from scaffold_domain.go (both in same package, already defined)
Option B: Add separate WizardScaffolderVersion constant

**Recommendation:** Option A - reuse existing constant since both tools are in same package.

---

### GEN-004: No DI wiring injection

**Severity:** Critical
**Location:** `internal/tools/scaffold_wizard.go` lines 248-256

**Current (scaffold_wizard.go lines 248-256):**
```go
nextSteps := []string{
    "templ generate",
    "go mod tidy",
    fmt.Sprintf("Register wizard routes in cmd/web/main.go"),
    fmt.Sprintf("Add wizard link to domain views (e.g., a 'New with Wizard' button)"),
}

if data.WithDrafts {
    nextSteps = append(nextSteps, "Add WizardDraft to database AutoMigrate")
}
```

**Expected (from scaffold_domain.go lines 293-317):**
```go
// Inject into main.go, database.go, and base_layout.templ if not dry run
if !input.DryRun {
    mainGoPath := filepath.Join(registry.WorkingDir, "cmd", "web", "main.go")
    databaseGoPath := filepath.Join(registry.WorkingDir, "internal", "database", "database.go")
    layoutPath := filepath.Join(registry.WorkingDir, "internal", "web", "layouts", "base_layout.templ")
    if utils.FileExists(mainGoPath) {
        if err := injectDomainWiring(mainGoPath, databaseGoPath, modulePath, pkgName, input.DomainName, data.RouteGroup, input.Relationships, data.WithCrudViews); err != nil {
            // Log warning but don't fail
            fmt.Printf("Warning: could not inject DI wiring: %v\n", err)
        } else {
            result.FilesUpdated = append(result.FilesUpdated, "cmd/web/main.go")
            if utils.FileExists(databaseGoPath) {
                result.FilesUpdated = append(result.FilesUpdated, "internal/database/database.go")
            }
            // ...
        }
    }
}
```

**Fix Approach:**
Create new function `injectWizardWiring()` that:
1. Injects wizard controller import into main.go
2. Injects wizard draft repo/service imports (if WithDrafts)
3. Injects wizard draft repo instantiation (if WithDrafts)
4. Injects wizard draft service instantiation (if WithDrafts)
5. Injects wizard controller instantiation
6. Injects wizard route registration

**Required modifier markers:**
- The wizard controller needs to be wired similarly to domain controllers
- Routes should use existing `MCP:ROUTES:*` markers based on route_group

---

### GEN-005: No database.go AutoMigrate injection

**Severity:** Critical
**Location:** `internal/tools/scaffold_wizard.go` - missing entirely

**scaffold_domain.go (lines 522-537):**
```go
// Inject model into database.go AutoMigrate
if databaseGoPath != "" && utils.FileExists(databaseGoPath) {
    dbInjector, err := modifier.NewInjector(databaseGoPath)
    if err != nil {
        return err
    }

    modelName := utils.ToModelName(domainName)
    if err := dbInjector.InjectModel(modelName); err != nil {
        return err
    }

    if err := dbInjector.Save(); err != nil {
        return err
    }
}
```

**scaffold_wizard.go:** Only tells user to manually add WizardDraft

**Fix Approach:**
When `WithDrafts=true`, inject WizardDraft model into database.go:
```go
if data.WithDrafts && !input.DryRun {
    databaseGoPath := filepath.Join(registry.WorkingDir, "internal", "database", "database.go")
    if utils.FileExists(databaseGoPath) {
        dbInjector, err := modifier.NewInjector(databaseGoPath)
        if err == nil {
            if err := dbInjector.InjectModel("WizardDraft"); err == nil {
                _ = dbInjector.Save()
                result.FilesUpdated = append(result.FilesUpdated, "internal/database/database.go")
            }
        }
    }
}
```

---

### GEN-006: No metadata tracking

**Severity:** Medium
**Location:** `internal/tools/scaffold_wizard.go` - missing entirely

**scaffold_domain.go (lines 357-364):**
```go
// Save scaffold metadata for future sync/upgrade capabilities
metaStore := metadata.NewStore(registry.WorkingDir)
if err := metaStore.SaveDomain(input.DomainName, input, ScaffolderVersion); err != nil {
    // Log warning but don't fail - metadata is optional
    fmt.Printf("Warning: could not save scaffold metadata: %v\n", err)
} else {
    result.FilesUpdated = append(result.FilesUpdated, ".mcp/scaffold-metadata.json")
}
```

**scaffold_wizard.go:** No metadata tracking

**Fix Approach:**
Option A: Extend metadata.Store to support wizard metadata (new SaveWizard method)
Option B: Create separate wizard metadata tracking

**Recommendation:** Option A - extend existing metadata system for consistency.

**Required changes to metadata/metadata.go:**
1. Add `WizardMetadata` struct similar to `DomainMetadata`
2. Add `Wizards map[string]WizardMetadata` to `ProjectMetadata`
3. Add `SaveWizard()` method to `Store`

**Then in scaffold_wizard.go:**
```go
// Save wizard metadata for future sync/upgrade capabilities
if !input.DryRun {
    metaStore := metadata.NewStore(registry.WorkingDir)
    if err := metaStore.SaveWizard(input.WizardName, input, ScaffolderVersion); err != nil {
        fmt.Printf("Warning: could not save wizard metadata: %v\n", err)
    } else {
        result.FilesUpdated = append(result.FilesUpdated, ".mcp/scaffold-metadata.json")
    }
}
```

---

## Implementation Priority

| Priority | Issue | Effort | Impact |
|----------|-------|--------|--------|
| 1 | GEN-004: DI wiring injection | High | Critical - eliminates manual step |
| 2 | GEN-005: AutoMigrate injection | Medium | Critical - eliminates manual step |
| 3 | GEN-001: modifier import | Low | Required for GEN-004/005 |
| 4 | GEN-002: metadata import | Low | Required for GEN-006 |
| 5 | GEN-006: metadata tracking | Medium | Enables sync/upgrade |
| 6 | GEN-003: version constant | Low | Already available in package |

## Implementation Order

**Phase 6 Plan 02 should:**
1. Add missing imports (GEN-001, GEN-002)
2. Create `injectWizardWiring()` function for DI wiring (GEN-004)
3. Add database.go AutoMigrate injection for WizardDraft (GEN-005)
4. Extend metadata system for wizards (GEN-006)
5. Reuse existing ScaffolderVersion (GEN-003 - no change needed)

---

## Notes

- The validation patterns in scaffold_wizard.go (lines 96-134) are adequate
- Dry run support is properly implemented
- Conflict checking is properly implemented
- The core generation logic is correct - only the post-generation wiring is missing
