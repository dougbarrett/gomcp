# Architecture

**Analysis Date:** 2026-01-06

## Pattern Overview

**Overall:** MCP Server with Template-Based Code Generation

**Key Characteristics:**
- Model Context Protocol (MCP) server using stdio transport
- Template-driven code scaffolding for Go web applications
- Clean architecture with layered, modular design
- Embedded template filesystem for portability
- Marker-based code injection for non-destructive updates

## Layers

**Transport Layer:**
- Purpose: MCP protocol communication via stdio
- Contains: Server initialization, protocol handling
- Location: `cmd/gomcp/main.go`
- Depends on: MCP SDK, Server layer
- Used by: MCP clients (Claude Code, etc.)

**Server Layer:**
- Purpose: MCP server configuration and setup
- Contains: Server name, version, tool instructions
- Location: `internal/server/server.go`
- Depends on: Tools layer (registry)
- Used by: Transport layer

**Tool Registry Layer:**
- Purpose: Central tool registration and coordination
- Contains: Phase-based tool organization, generator factory
- Location: `internal/tools/registry.go`
- Depends on: Generator layer, individual tool implementations
- Used by: Server layer

**Tool Implementation Layer:**
- Purpose: MCP tool command handlers
- Contains: 30+ scaffolding tools (scaffold_*, extend_*, list_*, analyze_*)
- Location: `internal/tools/scaffold_*.go`
- Depends on: Generator, Types, Utils layers
- Used by: Tool Registry

**Generator Layer:**
- Purpose: Core template processing engine
- Contains: Template execution, conflict detection, dry-run support
- Location: `internal/generator/`
- Depends on: Templates (embedded), Utils layer
- Used by: Tool implementations

**Template Layer:**
- Purpose: Embedded code templates
- Contains: Project, domain, view, component, auth, wizard templates
- Location: `internal/templates/` (embedded via `embed.go`)
- Depends on: Nothing (static assets)
- Used by: Generator layer

**Modifier Layer:**
- Purpose: Non-destructive code injection
- Contains: Marker-based injection for DI wiring
- Location: `internal/modifier/inject.go`
- Depends on: Utils layer
- Used by: Tool implementations (update_di_wiring)

**Utility Layers:**
- Purpose: Shared helpers and validation
- Contains: Filesystem, naming conventions, validation
- Location: `internal/utils/`, `internal/types/`
- Depends on: Standard library
- Used by: All layers

## Data Flow

**MCP Tool Execution:**

1. User invokes tool via MCP client
2. MCP SDK receives request via stdio (`cmd/gomcp/main.go`)
3. Server routes to registered tool handler (`internal/tools/registry.go`)
4. Tool validates input against schema (`internal/types/inputs.go`)
5. Tool transforms input to template data (`internal/generator/data.go`)
6. Generator loads templates from embedded FS (`internal/templates/embed.go`)
7. Generator executes templates with data (`internal/generator/generator.go`)
8. Content written to filesystem (or dry-run preview)
9. Result returned to MCP client (files created/updated/conflicted)

**State Management:**
- Stateless - Each tool invocation is independent
- Filesystem is the only persistent state
- Metadata tracked in `.scaffold-metadata.json` for domain analysis

## Key Abstractions

**Generator:**
- Purpose: Template processing with conflict tracking
- Examples: `internal/generator/generator.go`
- Pattern: Builder with methods like `GenerateFile()`, `GenerateFileIfNotExists()`
- Methods: `NewGenerator()`, `GenerateFile()`, `GetResult()`

**Registry:**
- Purpose: Tool lifecycle management
- Examples: `internal/tools/registry.go`
- Pattern: Registry with phase-based organization
- Methods: `RegisterAll()`, `NewGenerator()`

**Tool:**
- Purpose: Individual scaffolding operation
- Examples: `scaffold_domain.go`, `scaffold_form.go`, `scaffold_wizard.go`
- Pattern: MCP Tool interface with InputSchema and handler
- Workflow: Register → Validate → Generate → Report

**Data Transformers:**
- Purpose: Convert tool input to template data
- Examples: `ProjectData`, `DomainData`, `FieldData` in `internal/generator/data.go`
- Pattern: Struct with inference methods for defaults

**Modifier:**
- Purpose: Inject code between markers
- Examples: `internal/modifier/inject.go`
- Pattern: Regex-based marker detection and injection
- Markers: `MCP:IMPORTS:START`, `MCP:REPOS:START`, etc.

## Entry Points

**Primary - gomcp MCP Server:**
- Location: `cmd/gomcp/main.go`
- Triggers: MCP client connection via stdio
- Responsibilities: Initialize server, register tools, serve requests

**Secondary - mcp-mgr Bug Tracker:**
- Location: `cmd/mcp-mgr/main.go`
- Triggers: MCP client for bug management
- Responsibilities: Bug list, get, close operations via HTTP API

## Error Handling

**Strategy:** Return structured errors, never panic in tool handlers

**Patterns:**
- Error wrapping with `fmt.Errorf(..., %w, err)` for context
- Validation errors returned before generation starts
- File conflicts reported (not overwritten) unless force mode
- All errors bubble up to MCP response

**Types:**
- Validation errors: Invalid input, missing required fields
- Generation errors: Template execution failures
- Filesystem errors: Permission, path issues
- Conflict errors: Existing file would be overwritten

## Cross-Cutting Concerns

**Logging:**
- Minimal - errors returned via MCP protocol
- No structured logging framework

**Validation:**
- Comprehensive validation in `internal/utils/validation.go`
- Domain names, field types, paths all validated
- Reserved Go keywords checked

**Template Delimiters:**
- Uses `[[ ]]` instead of `{{ }}` to avoid conflicts with Go templates in generated code

**Dry-Run Support:**
- All generators support `dry_run: true`
- Preview changes without writing files
- Conflict detection without overwriting

---

*Architecture analysis: 2026-01-06*
*Update when major patterns change*
