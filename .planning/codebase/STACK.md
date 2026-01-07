# Technology Stack

**Analysis Date:** 2026-01-06

## Languages

**Primary:**
- Go 1.24.3 - All application code (`go.mod`)

**Secondary:**
- TypeScript - AWS CDK infrastructure for bug tracker (`cmd/mcp-mgr/cdk/`)
- HTML/CSS - Generated via templ templates and Tailwind CSS

## Runtime

**Environment:**
- Go 1.24.3 - CLI tool and MCP server
- Node.js - AWS CDK tooling (`cmd/mcp-mgr/cdk/package.json`)
- Lambda Runtime: Go 1.24 with custom bootstrap

**Package Manager:**
- Go Modules - `go.mod` and `go.sum`
- npm - For CDK infrastructure (`cmd/mcp-mgr/cdk/package.json`)
- Lockfile: `go.sum` present

## Frameworks

**Core (gomcp scaffolding tool):**
- modelcontextprotocol/go-sdk v1.2.0 - MCP protocol implementation (`go.mod`)
- text/template - Go standard library template engine

**Generated Projects Use:**
- go-chi/chi v5.1.0 - HTTP router (`internal/templates/project/go.mod.tmpl`)
- GORM v1.25.12 - ORM for database operations (`internal/templates/domain/model.go.tmpl`)
- templ v0.3.857 - Type-safe HTML templating (`internal/templates/project/go.mod.tmpl`)
- HTMX v2.0.0 - Frontend interactivity (`internal/templates/project/router.go.tmpl`)
- Tailwind CSS - Utility-first CSS (`internal/templates/project/tailwind.config.js.tmpl`)
- gorilla/csrf v1.7.2 - CSRF protection (`internal/templates/project/middleware.go.tmpl`)
- gorilla/sessions v1.2.2 - Session management

**Testing:**
- Go standard `testing` package - All tests
- No external test framework

**Build/Dev:**
- Taskfile - Task automation (`internal/templates/project/taskfile.yml.tmpl`)
- Air - Hot reload for development (`internal/templates/project/air.toml.tmpl`)
- AWS CDK v2.170.0 - Infrastructure as Code (`cmd/mcp-mgr/cdk/package.json`)

## Key Dependencies

**Critical (gomcp):**
- modelcontextprotocol/go-sdk v1.2.0 - MCP protocol (`go.mod`)
- jinzhu/inflection v1.0.0 - String inflection utilities (`go.mod`)
- aws/aws-cdk-go v2.233.0 - AWS CDK Go bindings (`go.mod`)

**Infrastructure (generated projects):**
- GORM drivers: SQLite, PostgreSQL, MySQL (`internal/templates/project/go.mod.tmpl`)
- BurntSushi/toml v1.3.2 - TOML configuration parsing
- golang.org/x/crypto - Password hashing

## Configuration

**Environment:**
- No .env files for gomcp itself (CLI tool)
- Generated projects use environment variables: DB_DRIVER, DB_DSN, SESSION_SECRET
- TOML configuration files for i18n content (`config/{locale}/`)

**Build:**
- `go.mod` - Module definition
- Taskfile templates for development workflows
- Air configuration for hot reload

## Platform Requirements

**Development:**
- macOS/Linux/Windows (any platform with Go 1.24+)
- No external dependencies for gomcp
- Node.js required for mcp-mgr CDK deployment

**Production:**
- gomcp: Runs as MCP server via stdio transport
- mcp-mgr: AWS Lambda with DynamoDB, API Gateway

---

*Stack analysis: 2026-01-06*
*Update after major dependency changes*
