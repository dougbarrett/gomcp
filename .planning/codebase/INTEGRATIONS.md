# External Integrations

**Analysis Date:** 2026-01-06

## APIs & External Services

**Model Context Protocol (MCP):**
- MCP Server Implementation - Provides scaffolding tools via stdio transport
  - SDK/Client: modelcontextprotocol/go-sdk v1.2.0 (`go.mod`)
  - Transport: stdio (standard input/output)
  - Entry point: `cmd/gomcp/main.go`

**Bug Tracking API:**
- Internal API for bug management - `cmd/mcp-mgr/client/client.go`
  - Integration method: HTTP REST API
  - Auth: API Key header (`x-api-key`)
  - Configuration: `~/.config/mcp-mgr/config.json`
  - Endpoints: GET /bugs, GET /bugs/{id}, POST /bugs/{id}/close

## Data Storage

**Databases (gomcp itself):**
- None - CLI scaffolding tool with no persistent storage

**Databases (generated projects support):**
- SQLite - Default for local development (`internal/templates/project/go.mod.tmpl`)
- PostgreSQL - Production option
- MySQL - Alternative production option
- Connection: via DB_DSN environment variable

**Bug Tracker Backend:**
- DynamoDB - NoSQL database for bug tracking (`cmd/mcp-mgr/cdk/lib/bug-tracker-stack.ts`)
  - Table: `gomcp-bugs` with GSI on status/created_at
  - Billing: Pay-per-request mode
  - SDK: AWS SDK v2 (dynamodb v1.38.1)

**File Storage:**
- Local filesystem only - Template output written to working directory
- Embedded filesystem for templates (`internal/templates/embed.go`)

**Caching:**
- None currently

## Authentication & Identity

**Auth Provider (generated projects):**
- Custom session-based auth - gorilla/sessions
- Password hashing via golang.org/x/crypto/bcrypt
- CSRF protection via gorilla/csrf (`internal/templates/project/middleware.go.tmpl`)
- Role-Based Access Control - User and admin roles (`internal/templates/auth/`)

**MCP-MGR API Auth:**
- API Key authentication - `x-api-key` header (`cmd/mcp-mgr/cdk/lib/bug-tracker-stack.ts`)

## Monitoring & Observability

**Error Tracking:**
- None - Errors returned via MCP protocol responses

**Analytics:**
- None

**Logs:**
- Standard output/error for MCP server
- AWS CloudWatch for Lambda functions (mcp-mgr)

## CI/CD & Deployment

**Hosting (gomcp):**
- Distributed as compiled Go binary
- Runs locally via MCP client

**Hosting (mcp-mgr):**
- AWS Lambda - Serverless compute (`cmd/mcp-mgr/cdk/lib/bug-tracker-stack.ts`)
  - Memory: 128MB
  - Timeout: 30 seconds
  - Runtime: Custom Go bootstrap

- API Gateway - REST API endpoint
  - CORS enabled
  - API Key protection for management endpoints

**Infrastructure:**
- AWS CDK v2.170.0 - Infrastructure as Code
  - Stack: `cmd/mcp-mgr/cdk/lib/bug-tracker-stack.ts`
  - Deploys: Lambda, DynamoDB, API Gateway

## Environment Configuration

**Development:**
- No environment variables required for gomcp
- Templates embed all resources via `//go:embed`
- Build with `go build -o /tmp/gomcp ./cmd/gomcp`

**Production (mcp-mgr):**
- AWS credentials configured via CDK deployment
- API Gateway URL and API Key stored in deployment outputs
- Configuration: `~/.config/mcp-mgr/config.json`

## Webhooks & Callbacks

**Incoming:**
- None - gomcp is a request-response tool

**Outgoing:**
- None

## Third-Party Code Generation

**Template Engine:**
- Go text/template with custom delimiters `[[ ]]`
- Embedded templates via `//go:embed` directive
- Template locations: `internal/templates/*`

**Generated Technology Stack:**
- HTMX - Frontend interactivity without JavaScript
- Tailwind CSS - Utility-first styling
- templ - Type-safe Go HTML templates
- GORM - Database ORM

---

*Integration audit: 2026-01-06*
*Update when adding/removing external services*
