---
phase: 02-controller-template-fixes
plan: 01
subsystem: templates
tags: [go-templates, wizard, controller, imports, templ]

# Dependency graph
requires:
  - phase: 01-analysis
    provides: Issue catalog identifying CTRL-004, CTRL-005, CTRL-006, CTRL-007
provides:
  - Missing middleware import added to wizard controller
  - Missing models import (conditional) added to wizard controller
  - Missing templ import added to wizard controller
  - Local render method added to wizard controller
affects: [02-02-PLAN, controller-template-fixes]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Local render method pattern for templ component rendering"
    - "Conditional imports with [[ if ]] for optional features"

key-files:
  created: []
  modified:
    - internal/templates/wizard/controller.go.tmpl

key-decisions:
  - "Follow domain controller import ordering: stdlib → internal → external"
  - "Models import is conditional on WithDrafts to avoid unused import"

patterns-established:
  - "Wizard controller follows same render pattern as domain controller"

issues-created: []

# Metrics
duration: 3min
completed: 2026-01-06
---

# Phase 2 Plan 1: Add Missing Imports and Render Method Summary

**Added middleware, models, and templ imports plus local render method to wizard controller template, fixing 4 CTRL issues as foundation for handler method fixes**

## Performance

- **Duration:** 3 min
- **Started:** 2026-01-06T~14:00:00Z
- **Completed:** 2026-01-06T~14:03:00Z
- **Tasks:** 2
- **Files modified:** 1

## Accomplishments

- Added middleware import for CSRF token access (CTRL-004)
- Added conditional models import for WizardDraft type (CTRL-005)
- Added templ import for Component type in render signature (CTRL-006)
- Added local render method matching domain controller pattern (CTRL-007)

## Task Commits

Each task was committed atomically:

1. **Task 1: Add missing imports** - `e044002` (fix)
2. **Task 2: Add local render method** - `417ee40` (feat)

**Plan metadata:** (to be committed with this file)

## Files Created/Modified

- `internal/templates/wizard/controller.go.tmpl` - Added 3 imports and render method

## Decisions Made

- Followed domain controller import ordering (stdlib, internal, external)
- Made models import conditional on `[[if .WithDrafts]]` to avoid unused import errors

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - both tasks completed successfully.

## Next Phase Readiness

- Foundation fixes complete (imports + render method)
- Ready for Plan 02-02: Fix response method calls (HXRedirect, CSRFToken, Component)
- Handler methods can now use the local render method

---
*Phase: 02-controller-template-fixes*
*Completed: 2026-01-06*
