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
- [ ] **Phase 3: Service Template Fixes** - Fix Bug #cb94adf6 + align with patterns
- [ ] **Phase 4: View Template Improvements** - Align view templates with other patterns
- [ ] **Phase 5: Draft System Review** - Review and improve draft save/resume functionality
- [ ] **Phase 6: Generator Logic Review** - Review wizard generator code for issues
- [ ] **Phase 7: Test Coverage** - Add/improve tests for wizard scaffolding
- [ ] **Phase 8: End-to-End Validation** - Generate test wizard, verify it compiles and runs

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
**Plans**: TBD

Plans:
- [ ] 03-01: Fix repository import in wizarddraft service template
- [ ] 03-02: Align service template with scaffold_service patterns

### Phase 4: View Template Improvements
**Goal**: Align wizard view templates with scaffold_view and scaffold_form patterns for consistency
**Depends on**: Phase 3
**Research**: Unlikely (internal patterns)
**Plans**: TBD

Plans:
- [ ] 04-01: Review wizard view templates against scaffold_view patterns
- [ ] 04-02: Apply consistency fixes to wizard view templates

### Phase 5: Draft System Review
**Goal**: Review wizard draft persistence system for correctness and completeness
**Depends on**: Phase 4
**Research**: Unlikely (internal patterns)
**Plans**: TBD

Plans:
- [ ] 05-01: Review draft model and repository templates
- [ ] 05-02: Verify draft save/resume workflow is correct

### Phase 6: Generator Logic Review
**Goal**: Review wizard generator code (scaffold_wizard.go) for issues and alignment with other generators
**Depends on**: Phase 5
**Research**: Unlikely (internal patterns)
**Plans**: TBD

Plans:
- [ ] 06-01: Review scaffold_wizard.go against scaffold_domain.go patterns
- [ ] 06-02: Fix any generator logic issues found

### Phase 7: Test Coverage
**Goal**: Add/improve test coverage for wizard scaffolding
**Depends on**: Phase 6
**Research**: Unlikely (following existing test patterns)
**Plans**: TBD

Plans:
- [ ] 07-01: Review existing test patterns in codebase
- [ ] 07-02: Add tests for wizard generator
- [ ] 07-03: Add tests for wizard templates

### Phase 8: End-to-End Validation
**Goal**: Generate a complete test wizard, verify it compiles, and run it to confirm fixes work
**Depends on**: Phase 7
**Research**: Unlikely (validation only)
**Plans**: TBD

Plans:
- [ ] 08-01: Generate test wizard with all step types
- [ ] 08-02: Verify generated code compiles and runs
- [ ] 08-03: Close bug reports if fixes verified

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Analysis | 2/2 | Complete | 2026-01-06 |
| 2. Controller Template Fixes | 2/2 | Complete | 2026-01-06 |
| 3. Service Template Fixes | 0/2 | Not started | - |
| 4. View Template Improvements | 0/2 | Not started | - |
| 5. Draft System Review | 0/2 | Not started | - |
| 6. Generator Logic Review | 0/2 | Not started | - |
| 7. Test Coverage | 0/3 | Not started | - |
| 8. End-to-End Validation | 0/3 | Not started | - |
