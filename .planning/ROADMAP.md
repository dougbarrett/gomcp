# Roadmap: go-mcp Wizard Scaffolding Improvement

## Overview

A focused improvement pass fixing known bugs and aligning wizard scaffolding templates with the quality standards of other scaffold_* tools. Starting with critical bug fixes (controller and service template issues), then systematically improving each template layer, adding test coverage, and validating the complete workflow.

## Domain Expertise

None

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

- [x] **Phase 1: Analysis** (Complete) - Understand wizard templates vs other scaffold_* patterns
- [x] **Phase 2: Controller Template Fixes** (Complete) - Fix Bug #a9479784 + align with patterns
- [x] **Phase 3: Service Template Fixes** (Complete) - Fix Bug #cb94adf6 + align with patterns
- [x] **Phase 4: View Template Improvements** (Complete) - Align view templates with other patterns
- [x] **Phase 5: Draft System Review** (Complete) - Review and improve draft save/resume functionality
- [x] **Phase 6: Generator Logic Review** (Complete) - Review wizard generator code for issues
- [x] **Phase 7: Test Coverage** (Complete) - Add/improve tests for wizard scaffolding
- [ ] **Phase 8: End-to-End Validation** - Generate test wizard, verify it compiles and runs
- [x] **Phase 9: Wizard Bug Fixes** (Complete) - Fix 4 bugs discovered during validation

## Phase Details

### Phase 1: Analysis
**Goal**: Understand current wizard template structure, compare to other scaffold_* tools, identify all discrepancies and issues beyond the reported bugs
**Depends on**: Nothing (first phase)
**Research**: Unlikely (internal codebase patterns)
**Plans**: TBD

Plans:
- [x] 01-01: Map wizard template files and generator structure
- [x] 01-02: Compare wizard patterns and document all issues (merged 01-03)

### Phase 2: Controller Template Fixes
**Goal**: Fix Bug #a9479784 - controller using non-existent web.Response methods (HXRedirect, Component, CSRFToken), double slashes in URLs, missing models import
**Depends on**: Phase 1
**Research**: Unlikely (following existing patterns)
**Plans**: 2

Plans:
- [x] 02-01: Add missing imports and render method (foundation fixes)
- [x] 02-02: Fix response method calls and URL construction

### Phase 3: Service Template Fixes
**Goal**: Fix Bug #cb94adf6 - wizard service missing repository import, ensure service template matches scaffold_service patterns
**Depends on**: Phase 2
**Research**: Unlikely (following existing patterns)
**Plans**: 1

Plans:
- [x] 03-01: Fix repository import and align with scaffold_service patterns (SVC-001, SVC-002)

### Phase 4: View Template Improvements
**Goal**: Align wizard view templates with scaffold_view and scaffold_form patterns for consistency
**Depends on**: Phase 3
**Research**: Unlikely (internal patterns)
**Plans**: TBD

Plans:
- [x] 04-01: Review wizard view templates against scaffold_view patterns
- [x] 04-02: Apply consistency fixes to wizard view templates

### Phase 5: Draft System Review
**Goal**: Review wizard draft persistence system for correctness and completeness
**Depends on**: Phase 4
**Research**: Unlikely (internal patterns)
**Plans**: TBD

Plans:
- [x] 05-01: Review draft model and repository templates
- [x] 05-02: Verify draft save/resume workflow is correct

### Phase 6: Generator Logic Review
**Goal**: Review wizard generator code (scaffold_wizard.go) for issues and alignment with other generators
**Depends on**: Phase 5
**Research**: Unlikely (internal patterns)
**Plans**: TBD

Plans:
- [x] 06-01: Review scaffold_wizard.go against scaffold_domain.go patterns
- [x] 06-02: Fix any generator logic issues found

### Phase 7: Test Coverage
**Goal**: Add/improve test coverage for wizard scaffolding
**Depends on**: Phase 6
**Research**: Unlikely (following existing test patterns)
**Plans**: TBD

Plans:
- [x] 07-01: Add tests for wizard DI wiring and metadata tracking
- [x] 07-02: Add tests for wizard generator
- [x] 07-03: Add tests for wizard templates

### Phase 8: End-to-End Validation
**Goal**: Generate a complete test wizard, verify it compiles, and run it to confirm fixes work
**Depends on**: Phase 7
**Research**: Unlikely (validation only)
**Plans**: TBD

Plans:
- [x] 08-01: Generate test wizard with all step types (created MCP test harness, 4 new bugs discovered)
- [x] 08-02: Verify generated code compiles and runs (Phase 9 fixes verified, 1 new minor bug logged)
- [ ] 08-03: Close bug reports if fixes verified

### Phase 9: Wizard Bug Fixes
**Goal**: Fix 4 bugs discovered during end-to-end validation that prevent generated wizard code from compiling
**Depends on**: Phase 8
**Research**: Unlikely (fixing known issues)
**Plans**: 3

**Bugs to fix:**
- b88f8cab: scaffold_wizard doesn't auto-generate wizard components
- 5ab684ea: Wizard controller has unused 'resp' variable
- ab2c40cb: Wizard controller references wrong DTO type name (CreateOrderDTO vs CreateOrderInput)
- 69d42b6e: belongs_to display assumes related model has Name field

Plans:
- [x] 09-01: Controller template fixes (bugs 5ab684ea, ab2c40cb)
- [x] 09-02: Wizard component auto-generation (bug b88f8cab)
- [x] 09-03: belongs_to display field fix (bug 69d42b6e)

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8 → 9

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Analysis | 2/2 | Complete | 2026-01-06 |
| 2. Controller Template Fixes | 2/2 | Complete | 2026-01-06 |
| 3. Service Template Fixes | 1/1 | Complete | 2026-01-06 |
| 4. View Template Improvements | 2/2 | Complete | 2026-01-07 |
| 5. Draft System Review | 2/2 | Complete | 2026-01-07 |
| 6. Generator Logic Review | 2/2 | Complete | 2026-01-07 |
| 7. Test Coverage | 3/3 | Complete | 2026-01-08 |
| 8. End-to-End Validation | 2/3 | In progress | - |
| 9. Wizard Bug Fixes | 3/3 | Complete | 2026-01-08 |
