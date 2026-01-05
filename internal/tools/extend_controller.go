package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/modifier"
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterExtendController registers the extend_controller tool.
func RegisterExtendController(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "extend_controller",
		Description: `Add custom endpoints to an existing controller.

Use this to add HTTP endpoints beyond standard CRUD. The controller must have been created
with scaffold_domain (which includes injection markers).

Example endpoints:
- POST /{id}/cancel: Trigger an action on a resource
- GET /{id}/stats: Get computed/aggregated data
- POST /{id}/duplicate: Create a copy
- PATCH /{id}/status: Partial update

Each endpoint adds a route registration and handler method. Use the body parameter
to provide the handler implementation, or leave empty for a TODO placeholder.

The handler has access to:
- c.service: The domain service
- web.NewResponse(w, r): Helper for JSON/HTMX responses
- chi.URLParam(r, "id"): URL parameters

Template variables available in body:
- [[.ModelName]]: The model name in PascalCase (e.g., "Order")
- [[.VariableName]]: The variable name in camelCase (e.g., "order")
- [[.PackageName]]: The package name (e.g., "order")
- [[.URLPath]]: The base URL path (e.g., "/orders")`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ExtendControllerInput) (*mcp.CallToolResult, types.ScaffoldResult, error) {
		result, err := extendController(registry, input)
		if err != nil {
			return nil, types.NewErrorResult(err.Error()), nil
		}
		return nil, result, nil
	})
}

// extendController adds endpoints to an existing controller.
func extendController(registry *Registry, input types.ExtendControllerInput) (types.ScaffoldResult, error) {
	if input.Domain == "" {
		return types.NewErrorResult("domain is required"), nil
	}
	if len(input.Endpoints) == 0 {
		return types.NewErrorResult("at least one endpoint is required"), nil
	}

	gen := registry.NewGenerator("")
	gen.SetDryRun(input.DryRun)

	// Derive names from domain
	packageName := utils.ToPackageName(input.Domain)
	modelName := utils.ToPascalCase(input.Domain)
	variableName := utils.ToCamelCase(input.Domain)
	urlPath := "/" + utils.ToKebabCase(utils.Pluralize(input.Domain))

	// Check if controller file exists
	controllerPath := filepath.Join("internal", "web", packageName, packageName+".go")
	if !gen.FileExists(controllerPath) {
		return types.NewErrorResult(fmt.Sprintf("controller file not found: %s. Use scaffold_domain first.", controllerPath)), nil
	}

	// Read the existing controller file
	content, err := gen.ReadFile(controllerPath)
	if err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to read controller file: %v", err)), nil
	}

	// Generate route registrations and handler implementations
	var routes []string
	var handlers []string

	for _, endpoint := range input.Endpoints {
		// Normalize HTTP method
		method := strings.ToUpper(endpoint.Method)
		chiMethod := strings.Title(strings.ToLower(method))

		// Route registration
		route := fmt.Sprintf("\t\tr.%s(\"%s\", c.%s)", chiMethod, endpoint.Path, endpoint.Name)
		routes = append(routes, route)

		// Handler implementation
		description := endpoint.Description
		if description == "" {
			description = fmt.Sprintf("handles %s %s%s", method, urlPath, endpoint.Path)
		}

		// Process body - replace template variables
		body := endpoint.Body
		if body == "" {
			body = `	res := web.NewResponse(w, r)
	// TODO: Implement this handler
	res.JSON(http.StatusOK, map[string]string{"status": "ok"})`
		}
		body = strings.ReplaceAll(body, "[[.ModelName]]", modelName)
		body = strings.ReplaceAll(body, "[[.VariableName]]", variableName)
		body = strings.ReplaceAll(body, "[[.PackageName]]", packageName)
		body = strings.ReplaceAll(body, "[[.URLPath]]", urlPath)

		handler := fmt.Sprintf(`
// %s %s
func (c *Controller) %s(w http.ResponseWriter, r *http.Request) {
%s
}`, endpoint.Name, description, endpoint.Name, body)
		handlers = append(handlers, handler)
	}

	// Use injector to add content between markers
	injector := modifier.NewInjectorFromContent(content)

	// Inject routes
	routeContent := strings.Join(routes, "\n")
	if err := injector.InjectBetweenMarkers("MCP:ROUTES:START", "MCP:ROUTES:END", routeContent); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to inject routes: %v. Make sure the controller was scaffolded with markers.", err)), nil
	}

	// Inject handlers
	handlerContent := strings.Join(handlers, "\n")
	if err := injector.InjectBetweenMarkers("MCP:HANDLERS:START", "MCP:HANDLERS:END", handlerContent); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to inject handlers: %v. Make sure the controller was scaffolded with markers.", err)), nil
	}

	// Write the modified file
	if err := gen.WriteFile(controllerPath, injector.Content()); err != nil {
		return types.NewErrorResult(fmt.Sprintf("failed to write controller file: %v", err)), nil
	}

	return types.ScaffoldResult{
		Success:      true,
		Message:      fmt.Sprintf("Added %d endpoint(s) to %s controller", len(input.Endpoints), modelName),
		FilesUpdated: []string{controllerPath},
	}, nil
}
