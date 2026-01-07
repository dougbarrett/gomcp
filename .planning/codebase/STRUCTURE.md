# Codebase Structure

**Analysis Date:** 2026-01-06

## Directory Layout

```
go-mcp/
├── cmd/                    # Command entry points
│   ├── gomcp/             # Primary MCP server
│   │   └── main.go        # Scaffolding server entry
│   └── mcp-mgr/           # Bug tracking MCP server
│       ├── main.go        # Bug tracker entry
│       ├── client/        # HTTP API client
│       ├── tools/         # Bug management tools
│       ├── types/         # Bug type definitions
│       └── cdk/           # AWS CDK infrastructure
│
├── internal/              # Private packages
│   ├── server/            # MCP server configuration
│   ├── tools/             # Tool implementations (30+ files)
│   ├── generator/         # Template processing engine
│   ├── templates/         # Embedded template files
│   ├── modifier/          # Code injection system
│   ├── metadata/          # Scaffolding metadata tracking
│   ├── types/             # Shared type definitions
│   └── utils/             # Utilities
│
├── go.mod                 # Module definition
├── go.sum                 # Dependency lock file
├── README.md              # Project documentation
├── CLAUDE.md              # Project instructions
└── coverage.out           # Test coverage report
```

## Directory Purposes

**cmd/gomcp/**
- Purpose: Primary MCP server entry point
- Contains: `main.go` - Server initialization
- Key files: `main.go` initializes MCP SDK, registers tools
- Subdirectories: None

**cmd/mcp-mgr/**
- Purpose: Secondary MCP server for bug tracking
- Contains: Server, client, tools, types, infrastructure
- Key files: `main.go`, `client/client.go`, `types/bug.go`
- Subdirectories:
  - `client/` - HTTP client for backend API
  - `tools/` - list_bugs, get_bug, close_bug implementations
  - `types/` - Bug type definitions
  - `cdk/` - AWS CDK stack (Lambda, DynamoDB, API Gateway)

**internal/server/**
- Purpose: MCP server factory and configuration
- Contains: `server.go`, `server_test.go`
- Key files: `server.go` - Server name, version, instructions

**internal/tools/**
- Purpose: All scaffolding tool implementations
- Contains: 30+ Go files, comprehensive tests
- Key files:
  - `registry.go` - Central tool registration
  - `scaffold_project.go` - Project initialization
  - `scaffold_domain.go` - Full domain generation
  - `scaffold_wizard.go` - Multi-step wizard generation
  - `update_di_wiring.go` - DI container updates
  - `analyze_domain.go` - Domain comparison

**internal/generator/**
- Purpose: Core template processing engine
- Contains: Generator, data transformers, helpers
- Key files:
  - `generator.go` - Template execution, conflict detection
  - `data.go` - Input to template data transformation
  - `helpers.go` - Template helper functions
  - `templates.go` - Template loading utilities

**internal/templates/**
- Purpose: Embedded template asset library
- Contains: All code generation templates
- Key files: `embed.go` - Embedded filesystem definition
- Subdirectories:
  - `project/` - Project scaffolding (go.mod, main.go, taskfile, etc.)
  - `domain/` - Domain layer (model, repository, service, controller)
  - `views/` - templ views (list, show, form, table)
  - `components/` - UI components (card, modal, form_field)
  - `config/` - TOML configuration templates
  - `seed/` - Database seeder templates
  - `auth/` - Authentication system templates
  - `usermgmt/` - User management templates
  - `wizard/` - Multi-step wizard templates

**internal/modifier/**
- Purpose: Non-destructive code injection
- Contains: Marker-based code injector
- Key files: `inject.go` - Injection logic with marker detection

**internal/metadata/**
- Purpose: Track scaffolded domains for analysis
- Contains: Metadata storage and retrieval
- Key files: `metadata.go` - JSON-based metadata tracking

**internal/types/**
- Purpose: Shared type definitions
- Contains: Input and output type structures
- Key files:
  - `inputs.go` - Tool input types with validation
  - `outputs.go` - Tool output types

**internal/utils/**
- Purpose: Shared utilities
- Contains: Filesystem, naming, validation helpers
- Key files:
  - `filesystem.go` - File operations (copy, create, etc.)
  - `naming.go` - Naming conventions (PascalCase, camelCase, etc.)
  - `validation.go` - Input validation (500+ lines)

## Key File Locations

**Entry Points:**
- `cmd/gomcp/main.go` - Primary MCP server
- `cmd/mcp-mgr/main.go` - Bug tracking server

**Configuration:**
- `go.mod` - Module definition and dependencies
- `CLAUDE.md` - Project instructions for AI assistants

**Core Logic:**
- `internal/tools/registry.go` - Tool registration
- `internal/generator/generator.go` - Template processing
- `internal/modifier/inject.go` - Code injection

**Testing:**
- `internal/tools/*_test.go` - Tool tests
- `internal/generator/*_test.go` - Generator tests
- `internal/tools/integration_test.go` - E2E tests
- `coverage.out` - Coverage report

**Documentation:**
- `README.md` - User-facing documentation

## Naming Conventions

**Files:**
- `snake_case.go` - All Go source files
- `*_test.go` - Test files (co-located)
- `*.tmpl` - Template files
- `*.templ.tmpl` - templ component templates

**Directories:**
- `lowercase` - All directories
- Singular names for packages: `generator`, `modifier`
- Plural for collections: `tools`, `templates`

**Special Patterns:**
- `scaffold_*.go` - Scaffolding tool implementations
- `extend_*.go` - Extension tool implementations
- `embed.go` - Embedded filesystem declarations

## Where to Add New Code

**New Scaffolding Tool:**
- Primary code: `internal/tools/scaffold_{name}.go`
- Register in: `internal/tools/registry.go`
- Tests: `internal/tools/scaffold_{name}_test.go`

**New Template:**
- Implementation: `internal/templates/{category}/{name}.tmpl`
- Update embed: `internal/templates/embed.go` (if new directory)

**New Generator Feature:**
- Implementation: `internal/generator/{feature}.go`
- Tests: `internal/generator/{feature}_test.go`

**New Type Definitions:**
- Implementation: `internal/types/inputs.go` or `outputs.go`
- Tests: `internal/types/*_test.go`

**Utilities:**
- Shared helpers: `internal/utils/{purpose}.go`
- Tests: `internal/utils/{purpose}_test.go`

## Special Directories

**internal/templates/**
- Purpose: Embedded template assets
- Source: Go `//go:embed` directive
- Committed: Yes (source of truth)

**cmd/mcp-mgr/cdk/**
- Purpose: AWS infrastructure as code
- Source: TypeScript CDK stack
- Committed: Yes (infrastructure definition)

---

*Structure analysis: 2026-01-06*
*Update when directory structure changes*
