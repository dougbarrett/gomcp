# Deferred Issues

Issues discovered during execution but not blocking current work.

## Open

### stepData unused variable in wizard submit handler
- **Discovered:** 08-02-PLAN.md
- **Location:** `templates/controller_wizard.tmpl` (wizard submit handler)
- **Issue:** The wizard submit handler retrieves `stepData` but doesn't map it to the DTO
- **Workaround:** Add `_ = stepData` statement (applied in 08-02 validation test)
- **Priority:** Low - template bug, easy fix
- **Suggested fix:** Either use stepData to populate DTO fields or remove the retrieval

## Closed

None yet.
