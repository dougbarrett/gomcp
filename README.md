# Go MCP Scaffolding Tool

A Model Context Protocol (MCP) server that provides scaffolding tools for Go web applications. Generate complete project structures, domains, views, and configurations following clean architecture principles.

## Features

- **Project Scaffolding**: Initialize complete Go web application projects
- **Domain-Driven Design**: Generate models, repositories, services, and controllers
- **View Generation**: Create templ-based views with reusable components
- **Configuration Management**: Multi-locale TOML configuration files
- **Dependency Injection**: Automatic DI wiring with marker-based code injection

## Installation

```bash
go install github.com/dbb1dev/go-mcp/cmd/gomcp@latest
```

## Usage

The tool runs as an MCP server and is designed to be used with MCP-compatible clients.

### Claude Code

Add the MCP server using the Claude Code CLI:

```bash
claude mcp add gomcp -- gomcp
```

Or for a specific project directory:

```bash
claude mcp add gomcp -e MCP_SCAFFOLD_WORKDIR=/path/to/your/project -- gomcp
```

### VS Code (GitHub Copilot)

Add to your VS Code settings (`.vscode/settings.json`):

```json
{
  "github.copilot.chat.mcpServers": {
    "gomcp": {
      "command": "gomcp",
      "env": {
        "MCP_SCAFFOLD_WORKDIR": "${workspaceFolder}"
      }
    }
  }
}
```

### VS Code (Continue)

Add to your Continue config (`~/.continue/config.json`):

```json
{
  "experimental": {
    "modelContextProtocolServers": [
      {
        "transport": {
          "type": "stdio",
          "command": "gomcp",
          "env": {
            "MCP_SCAFFOLD_WORKDIR": "/path/to/your/project"
          }
        }
      }
    ]
  }
}
```

### Manual

Run the MCP server directly:

```bash
gomcp
```

## Current Capabilities

### Project Scaffolding (`scaffold_project`)

Creates a complete Go web application structure:

| Component            | Description                             |
| -------------------- | --------------------------------------- |
| `cmd/web/main.go`    | Application entry point with DI wiring  |
| `cmd/seed/main.go`   | Database seeding entry point            |
| `internal/config/`   | Configuration management                |
| `internal/database/` | GORM database setup                     |
| `internal/models/`   | Base model with timestamps              |
| `internal/web/`      | Router, middleware, layouts, components |
| `config/`            | TOML configuration files                |
| `Taskfile.yml`       | Task runner configuration               |
| `.air.toml`          | Hot reload configuration                |

**Supported databases**: SQLite, PostgreSQL, MySQL

**Authentication scaffolding** (with `with_auth: true`):

When enabled, generates a complete authentication system:

| Component                           | Description                                        |
| ----------------------------------- | -------------------------------------------------- |
| `internal/models/user.go`           | User model with password hashing (bcrypt)          |
| `internal/repository/user/user.go`  | User data access layer                             |
| `internal/services/auth/auth.go`    | Login, register, logout, password change           |
| `internal/services/auth/session.go` | Session management with gorilla/sessions           |
| `internal/web/middleware/auth.go`   | RequireAuth, RequireAdmin, OptionalAuth middleware |
| `internal/web/auth/auth.go`         | Auth HTTP handlers                                 |
| `internal/web/auth/views/*.templ`   | Login and registration pages                       |

Features:

- Secure password hashing with bcrypt
- Cookie-based session management
- Remember me functionality (extended session)
- Role-based access control (user/admin)
- HTMX-compatible redirects
- Flash messages for errors/success
- Last login tracking

### Domain Scaffolding (`scaffold_domain`)

Generates a complete domain with all layers:

```
internal/
├── models/{domain}.go           # GORM model
├── repository/{domain}/         # Data access layer
│   └── {domain}.go
├── services/{domain}/           # Business logic layer
│   ├── {domain}.go
│   └── dto.go
└── web/{domain}/                # HTTP handlers
    ├── {domain}.go
    └── views/                   # templ views (optional)
```

**Supported field types**:

- `string`, `int`, `int64`, `uint`, `float32`, `float64`
- `bool`, `time.Time`
- Pointer types (`*string`, `*int`, etc.)

**Relationship support**:

Define model associations with the `relationships` field:

```json
{
  "domain_name": "order",
  "fields": [{ "name": "Total", "type": "float64" }],
  "relationships": [
    { "type": "belongs_to", "model": "User", "preload": true },
    { "type": "has_many", "model": "OrderItem" },
    { "type": "many_to_many", "model": "Tag", "join_table": "order_tags" }
  ]
}
```

| Relationship Type | Description                  | Generated                 |
| ----------------- | ---------------------------- | ------------------------- |
| `belongs_to`      | Foreign key on this model    | FK field + relation field |
| `has_one`         | Foreign key on related model | Relation field            |
| `has_many`        | One-to-many relationship     | Slice relation field      |
| `many_to_many`    | Join table relationship      | Slice relation field      |

Relationship options:

- `preload`: Auto-load relationship in queries
- `foreign_key`: Custom foreign key field name
- `references`: Referenced field (default: "ID")
- `join_table`: Join table name (for many_to_many)
- `on_delete`: DELETE constraint (CASCADE, SET NULL, RESTRICT, NO ACTION)

### Standalone Layer Tools

| Tool                  | Description                               |
| --------------------- | ----------------------------------------- |
| `scaffold_repository` | Generate only the repository layer        |
| `scaffold_service`    | Generate only the service layer with DTOs |
| `scaffold_controller` | Generate only the HTTP controller         |

### View Tools

| Tool                 | Description                                          |
| -------------------- | ---------------------------------------------------- |
| `scaffold_view`      | Generate templ views (list, show, form, table, card) |
| `scaffold_form`      | Generate HTMX-powered forms                          |
| `scaffold_table`     | Generate data tables with pagination/sorting         |
| `scaffold_modal`     | Generate modal dialogs                               |
| `scaffold_component` | Generate reusable templ components                   |
| `scaffold_page`      | Generate page templates with TOML config             |

### Configuration Tools

| Tool               | Description                                            |
| ------------------ | ------------------------------------------------------ |
| `scaffold_config`  | Generate TOML config files (page, menu, app, messages) |
| `scaffold_seed`    | Generate database seeder with optional faker support   |
| `list_domains`     | List all scaffolded domains in the project             |
| `update_di_wiring` | Update main.go with DI wiring for domains              |

### Code Injection

The tool uses marker comments to inject code into existing files:

```go
// MCP:IMPORTS:START
// MCP:IMPORTS:END

// MCP:REPOS:START
// MCP:REPOS:END

// MCP:SERVICES:START
// MCP:SERVICES:END

// MCP:CONTROLLERS:START
// MCP:CONTROLLERS:END

// MCP:ROUTES:START
// MCP:ROUTES:END
```

## Technology Stack

Generated projects use:

- **[Go](https://go.dev/)** - Programming language
- **[GORM](https://gorm.io/)** - ORM
- **[templ](https://templ.guide/)** - Type-safe HTML templating
- **[HTMX](https://htmx.org/)** - Frontend interactivity
- **[Tailwind CSS](https://tailwindcss.com/)** - Styling
- **[Task](https://taskfile.dev/)** - Task runner
- **[Air](https://github.com/cosmtrek/air)** - Hot reload

---

## Roadmap

### Planned Features

#### API-Only Mode

RESTful JSON API generation:

- JSON response handlers (no templ views)
- OpenAPI/Swagger specification
- API versioning (`/api/v1/...`)
- Rate limiting middleware
- CORS configuration

#### Migration Support

Database migration management:

- `scaffold_migration` tool
- Up/down migration files
- Migration runner integration
- Seed data tied to migrations

#### Enhanced Validation

Field-level validation rules:

```json
{ "name": "Email", "type": "string", "validation": "required,email,unique" }
```

- Service layer validation
- Form validation in views
- Custom validation messages
- Cross-field validation

#### Test Generation

Automated test scaffolding:

- Repository tests with mock database
- Service tests with mock repository
- Controller/handler tests
- Integration test templates

#### Search & Filtering

Enhanced list views:

- Full-text search across fields
- Field-specific filters
- Date range filters
- Sort by multiple columns

#### Audit Trail

Automatic tracking:

- `created_by`, `updated_by` fields
- Automatic population from auth context
- Audit log table option

#### Background Jobs

Async job scaffolding:

- Job handler templates
- Queue integration
- Retry logic
- Job status tracking

#### Admin Dashboard

Auto-generated admin interface:

- CRUD for all domains
- Dashboard with statistics
- User management
- Activity logs

---

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/tools/... -v
```

### Test Coverage

| Package              | Coverage |
| -------------------- | -------- |
| `internal/types`     | 100%     |
| `internal/modifier`  | 97.2%    |
| `internal/utils`     | 90.4%    |
| `internal/generator` | 81.1%    |
| `internal/tools`     | 81.1%    |

### Project Structure

```
go-mcp/
├── cmd/gomcp/            # MCP server entry point
├── internal/
│   ├── generator/        # Template generation engine
│   ├── modifier/         # Code injection system
│   ├── server/           # MCP server setup
│   ├── templates/        # Embedded Go templates
│   │   ├── project/      # Project scaffolding templates
│   │   ├── domain/       # Domain layer templates
│   │   ├── views/        # View templates
│   │   ├── components/   # Component templates
│   │   ├── config/       # Config templates
│   │   └── seed/         # Seeder templates
│   ├── tools/            # MCP tool implementations
│   ├── types/            # Input/output types
│   └── utils/            # Validation and naming utilities
└── README.md
```

## License

MIT
