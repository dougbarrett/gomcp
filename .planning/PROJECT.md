# go-mcp Wizard Scaffolding Improvement

## What This Is

A focused improvement pass on the wizard scaffolding feature of go-mcp, fixing known bugs and bringing the wizard templates up to the quality standard of other scaffold_* tools. The goal is a production-ready wizard scaffolding tool that generates correct, consistent code.

## Core Value

~~The 2 reported bugs must be fixed first. Everything else is valuable, but working code generation is non-negotiable.~~

**COMPLETE:** All bugs fixed and verified. Wizard scaffolding generates correct, compilable code.

## Requirements

### Validated

- Wizard scaffolding generates controller, views, and draft persistence — existing
- Multi-step wizard flow with form, select, has_many, and summary step types — existing
- Draft save/resume functionality via database — existing
- HTMX-powered step navigation — existing
- Integration with existing domain scaffolding patterns — existing

### Completed

- [x] Fix wizard controller using non-existent web.Response methods (Bug #a9479784) — Phase 2
- [x] Fix wizard service missing repository import (Bug #cb94adf6) — Phase 3
- [x] Improve template quality and consistency with other scaffold_* tools — Phases 4-6
- [x] Add/improve test coverage for wizard scaffolding — Phase 7
- [x] Ensure feature completeness for real-world wizard use cases — Phase 8
- [x] Fix additional bugs discovered during validation (Bugs #b88f8cab, #69d42b6e) — Phase 9

### Out of Scope

- UI redesign — Focus on code correctness, not visual changes
- Breaking changes — Existing wizard scaffolds should continue to work
- New step types beyond existing (form, select, has_many, summary) — Stabilize what exists first

## Context

**Existing Codebase:**
- MCP server using Go 1.24.3 with modelcontextprotocol/go-sdk
- Template-driven code generation with embedded templates in `internal/templates/`
- 30+ scaffolding tools following consistent patterns
- Generator uses `[[ ]]` delimiters to avoid conflicts with Go templates

**Bug Details:**

Bug #a9479784 - Wizard controller uses non-existent web.Response methods:
- `resp.HXRedirect(url)` should be `resp.Redirect(url)`
- `resp.Component(component)` should be `resp.Render(component)`
- `resp.CSRFToken()` doesn't exist — CSRF token comes from middleware context
- URLs have double slashes: `//orders/wizard/...`
- Missing models import in wizard controller

Bug #cb94adf6 - Wizard service missing repository import:
- Generated wizarddraft service uses `Repository` type without importing the repository package
- Needs qualified import: `wizarddraftrepo "module/internal/repository/wizarddraft"`

**Known Patterns:**
- Other scaffold_* tools use `internal/tools/scaffold_*.go` structure
- Templates in `internal/templates/wizard/`
- Tests follow `*_test.go` convention with table-driven tests

## Constraints

- **Consistency**: Must match patterns used in other scaffold_* tools (domain, form, table, etc.)
- **Backwards Compatible**: No changes that would break existing generated wizard code

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Fix bugs first, then improve | Core value — working code before polish | Done (Phases 2-3) |
| Match existing scaffold patterns | Consistency constraint from user | Done (Phases 4-6) |
| URL double slash was generator bug, not template | Root cause analysis in Phase 2 | Fixed in data.go |
| Use GenerateFileIfNotExists for wizard components | Preserve user customizations | Implemented in Phase 9 |
| Use DisplayField for belongs_to display | Match form/show template patterns | Implemented in Phase 9 |
| Phase 9 inserted for validation bugs | Fix bugs before declaring complete | All 4 bugs fixed |

## Project Completion

All 9 phases complete:
- Phase 1: Analysis — mapped templates, cataloged issues
- Phase 2: Controller template fixes — Bug #a9479784
- Phase 3: Service template fixes — Bug #cb94adf6
- Phase 4: View template improvements — consistency review
- Phase 5: Draft system review — verified persistence
- Phase 6: Generator logic review — metadata tracking
- Phase 7: Test coverage — added wizard tests
- Phase 8: End-to-end validation — verified compilation
- Phase 9: Wizard bug fixes — fixed 4 validation bugs

---
*Last updated: 2026-01-10 — PROJECT COMPLETE*
