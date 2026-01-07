# go-mcp Wizard Scaffolding Improvement

## What This Is

A focused improvement pass on the wizard scaffolding feature of go-mcp, fixing known bugs and bringing the wizard templates up to the quality standard of other scaffold_* tools. The goal is a production-ready wizard scaffolding tool that generates correct, consistent code.

## Core Value

The 2 reported bugs must be fixed first. Everything else is valuable, but working code generation is non-negotiable.

## Requirements

### Validated

- Wizard scaffolding generates controller, views, and draft persistence — existing
- Multi-step wizard flow with form, select, has_many, and summary step types — existing
- Draft save/resume functionality via database — existing
- HTMX-powered step navigation — existing
- Integration with existing domain scaffolding patterns — existing

### Active

- [ ] Fix wizard controller using non-existent web.Response methods (Bug #a9479784)
- [ ] Fix wizard service missing repository import (Bug #cb94adf6)
- [ ] Improve template quality and consistency with other scaffold_* tools
- [ ] Add/improve test coverage for wizard scaffolding
- [ ] Ensure feature completeness for real-world wizard use cases

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
| Fix bugs first, then improve | Core value — working code before polish | — Pending |
| Match existing scaffold patterns | Consistency constraint from user | — Pending |

---
*Last updated: 2026-01-06 after initialization*
