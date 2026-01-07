# Testing Patterns

**Analysis Date:** 2026-01-06

## Test Framework

**Runner:**
- Go standard `testing` package
- No external test framework

**Assertion Library:**
- Built-in `t.Errorf()`, `t.Fatalf()`
- No assertion library (chai, testify, etc.)

**Run Commands:**
```bash
go test ./...                           # Run all tests
go test ./... -cover                    # Show coverage summary
go test -short ./...                    # Skip integration tests
go test ./internal/tools/...            # Run specific package
go test -run TestScaffoldDomain ./...   # Run specific test
go test ./... -coverprofile=coverage.out  # Generate coverage file
```

## Test File Organization

**Location:**
- Co-located with source: `*.go` + `*_test.go` in same directory
- No separate `tests/` directory

**Naming:**
- Unit tests: `{module}_test.go`
- Integration tests: `integration_test.go`
- Edge cases: `edge_cases_test.go`

**Structure:**
```
internal/
├── generator/
│   ├── generator.go
│   ├── generator_test.go
│   ├── data.go
│   ├── data_test.go
│   └── helpers_test.go
├── tools/
│   ├── scaffold_domain.go
│   ├── scaffold_domain_test.go
│   ├── integration_test.go
│   └── edge_cases_test.go
└── utils/
    ├── naming.go
    └── naming_test.go
```

## Test Structure

**Suite Organization:**
```go
func TestToPascalCase(t *testing.T) {
    tests := []struct {
        input string
        want  string
    }{
        {"", ""},
        {"user", "User"},
        {"user_profile", "UserProfile"},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got := ToPascalCase(tt.input)
            if got != tt.want {
                t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
            }
        })
    }
}
```

**Patterns:**
- Table-driven tests with `tests := []struct{...}`
- Subtests with `t.Run()` for organization
- `t.Helper()` for helper functions
- `t.Cleanup()` for resource cleanup

## Mocking

**Framework:**
- No mocking framework
- Interface-based testing where needed

**Patterns:**
```go
// Temporary directory for file operations
tmpDir, err := os.MkdirTemp("", "test-*")
if err != nil {
    t.Fatalf("failed to create temp dir: %v", err)
}
t.Cleanup(func() {
    os.RemoveAll(tmpDir)
})
```

**What to Mock:**
- Filesystem (using temp directories)
- Working directory isolation

**What NOT to Mock:**
- Template engine (test actual behavior)
- Validation logic (test real validation)

## Fixtures and Factories

**Test Data:**
```go
// Factory pattern in tests
func testRegistry(t *testing.T, workingDir string) *Registry {
    t.Helper()
    return NewRegistry(workingDir)
}

// Setup helpers
func setupGoMod(t *testing.T, dir string) {
    t.Helper()
    content := `module testmodule\n\ngo 1.21`
    err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644)
    if err != nil {
        t.Fatalf("failed to create go.mod: %v", err)
    }
}
```

**Location:**
- Factory functions in test file near usage
- Embedded test templates: `internal/generator/testdata/`
- Inline test data for simple cases

## Coverage

**Requirements:**
- No enforced minimum coverage
- Coverage tracked for awareness

**By Package (from README.md):**
- `internal/types`: 100%
- `internal/modifier`: 97.2%
- `internal/utils`: 90.4%
- `internal/generator`: 81.1%
- `internal/tools`: 81.1%

**Configuration:**
- Built-in Go coverage via `-cover` flag
- Coverage file: `coverage.out` (1,361 lines)

**View Coverage:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Types

**Unit Tests:**
- Scope: Single function in isolation
- Speed: Fast (<100ms per test)
- Examples: `internal/utils/naming_test.go`, `internal/generator/helpers_test.go`

**Integration Tests:**
- Scope: Multiple modules together
- Setup: Real file operations in temp directories
- Skip in short mode: `if testing.Short() { t.Skip("...") }`
- Examples: `internal/tools/integration_test.go`

**Edge Case Tests:**
- Scope: Boundary conditions and error cases
- Examples: `internal/tools/edge_cases_test.go`
- Tests: Invalid inputs, empty values, special characters

**Template Validation Tests:**
- Scope: Template loading and consistency
- Examples: `internal/tools/template_validation_test.go`
- Ensures all templates load without errors

## Common Patterns

**Async Testing:**
```go
// Not applicable - synchronous operations only
```

**Error Testing:**
```go
func TestValidateDomainName(t *testing.T) {
    tests := []struct {
        name       string
        domainName string
        wantErr    bool
    }{
        {"empty domain", "", true},
        {"valid domain", "user", false},
        {"invalid chars", "user@profile", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateDomainName(tt.domainName)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateDomainName(%q) error = %v, wantErr %v",
                    tt.domainName, err, tt.wantErr)
            }
        })
    }
}
```

**File System Testing:**
```go
func TestGenerateFile(t *testing.T) {
    tmpDir, err := os.MkdirTemp("", "generator-test-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    t.Cleanup(func() {
        os.RemoveAll(tmpDir)
    })

    g := NewGenerator(tmpDir)
    // ... test code
}
```

**Snapshot Testing:**
- Not used in this codebase
- Prefer explicit assertions for clarity

## Test Statistics

**Files:** 35 test files
**Lines:** ~15,382 lines of test code
**Source:** ~9,842 lines of non-test code
**Ratio:** 1.56:1 (test to source)

---

*Testing analysis: 2026-01-06*
*Update when test patterns change*
