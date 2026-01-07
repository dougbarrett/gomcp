# Coding Conventions

**Analysis Date:** 2026-01-06

## Naming Patterns

**Files:**
- `snake_case.go` for all Go source files
- `*_test.go` co-located with source
- `scaffold_{domain}.go` for scaffolding tools
- `extend_{layer}.go` for extension tools

**Functions:**
- camelCase for all functions (Go standard)
- No special prefix for async functions
- `To{Format}` for conversion functions: `ToPascalCase`, `ToCamelCase` (`internal/utils/naming.go`)
- `Must{Action}` for functions that panic on error: `MustLoadTemplate`, `MustParseTemplate`

**Variables:**
- camelCase for variables: `workingDir`, `tmpDir`, `domainName`
- Domain suffixes: `userRepo`, `productService`, `orderController` (`internal/utils/naming.go:110-126`)
- Single-letter receivers: `(g *Generator)`, `(r *Registry)`

**Types:**
- PascalCase for structs: `FileConflict`, `Generator`, `DomainData`
- No `I` prefix for interfaces
- Input types: `ScaffoldProjectInput`, `ScaffoldDomainInput`
- Output types with `Result` suffix: `GeneratorResult`

## Code Style

**Formatting:**
- Standard `gofmt` formatting
- No custom `.editorconfig` at project root
- Tabs for indentation (Go standard)
- No line length limit enforced

**Linting:**
- No `.golangci.yml` configuration
- Relies on Go standard tooling
- Run: `go vet ./...`

## Import Organization

**Order:**
1. Standard library packages
2. Internal packages (`github.com/dougbarrett/go-mcp/internal/...`)
3. External packages

**Grouping:**
- Blank line between groups
- No explicit sorting within groups

**Path Aliases:**
- None used - full import paths

## Error Handling

**Patterns:**
- Return errors, catch at tool handler level
- Error wrapping with context: `fmt.Errorf("failed to X: %w", err)`
- 41 occurrences of `%w` error wrapping across codebase

**Error Types:**
- Validation errors returned early before generation
- File conflicts reported as structured data
- Template errors wrapped with template path context

**Example:**
```go
// From internal/generator/generator.go
err := g.executeTemplate(templatePath, data, content)
if err != nil {
    return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
}
```

## Logging

**Framework:**
- None - errors returned via MCP protocol
- No structured logging library

**Patterns:**
- Errors bubble up to MCP response
- No console logging in production code

## Comments

**When to Comment:**
- Package documentation at top of file
- Exported functions and types documented
- Complex logic explained inline

**Documentation Style:**
```go
// NewGenerator creates a new Generator instance for the given working directory.
// It initializes the template engine and sets up conflict tracking.
func NewGenerator(workingDir string) *Generator
```

**Struct Field Documentation:**
```go
// FileConflict represents a file that would be overwritten.
type FileConflict struct {
    // Path is the relative file path.
    Path string
    // Description explains the purpose of this file.
    Description string
}
```

**TODO Comments:**
- Format: `// TODO: description`
- Minimal TODOs in codebase (well-maintained)

## Function Design

**Size:**
- Most functions under 50 lines
- Complex functions in generator are well-encapsulated

**Parameters:**
- Context-based structs for complex inputs
- Options patterns for optional configuration

**Return Values:**
- Explicit returns
- Multiple returns with error last: `(result, error)`
- Early return for validation failures

## Module Design

**Exports:**
- Named exports only (no default exports in Go)
- Public API through package-level functions
- Internal helpers unexported (lowercase)

**Package Organization:**
- Single responsibility per package
- Internal packages for private code
- Minimal cross-package dependencies

## Template Conventions

**Delimiters:**
- Uses `[[ ]]` instead of `{{ }}` (`internal/generator/generator.go`)
- Avoids conflicts with Go templates in generated code

**Template Variables:**
- PascalCase: `[[.ModelName]]`, `[[.PackageName]]`
- camelCase for local: `[[.variableName]]`

**Markers for Injection:**
- Format: `// MCP:{SECTION}:START` / `// MCP:{SECTION}:END`
- Examples: `MCP:IMPORTS:START`, `MCP:REPOS:START`, `MCP:ROUTES:START`

---

*Convention analysis: 2026-01-06*
*Update when patterns change*
