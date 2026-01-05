package tools

import (
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldController(t *testing.T) {
	t.Run("validates domain name", func(t *testing.T) {
		tests := []struct {
			name       string
			domainName string
			wantErr    string
		}{
			{
				name:       "empty domain",
				domainName: "",
				wantErr:    "domain name is required",
			},
			{
				name:       "invalid chars",
				domainName: "my@domain",
				wantErr:    "invalid character",
			},
			{
				name:       "starts with number",
				domainName: "123domain",
				wantErr:    "must start with a letter",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, tmpDir := testRegistry(t)
				setupGoMod(t, tmpDir, "github.com/test/project")

				input := types.ScaffoldControllerInput{
					DomainName: tt.domainName,
					Actions: []types.ActionDef{
						{Name: "list", Method: "GET", Path: "/"},
					},
				}

				result, err := scaffoldController(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Errorf("expected failure, got success")
				}
				if !containsString(result.Message, tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, result.Message)
				}
			})
		}
	})

	t.Run("requires at least one action", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions:    []types.ActionDef{},
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when no actions provided")
		}
		if !containsString(result.Message, "at least one action") {
			t.Errorf("expected error about actions, got: %s", result.Message)
		}
	})

	t.Run("validates action names", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions: []types.ActionDef{
				{Name: "", Method: "GET", Path: "/"},
			},
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when action name is empty")
		}
		if !containsString(result.Message, "action name is required") {
			t.Errorf("expected error about action name, got: %s", result.Message)
		}
	})

	t.Run("validates HTTP methods", func(t *testing.T) {
		tests := []struct {
			name    string
			method  string
			wantErr bool
		}{
			{name: "GET", method: "GET", wantErr: false},
			{name: "POST", method: "POST", wantErr: false},
			{name: "PUT", method: "PUT", wantErr: false},
			{name: "DELETE", method: "DELETE", wantErr: false},
			{name: "PATCH", method: "PATCH", wantErr: false},
			{name: "invalid", method: "INVALID", wantErr: true},
			{name: "empty", method: "", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, tmpDir := testRegistry(t)
				setupGoMod(t, tmpDir, "github.com/test/project")

				input := types.ScaffoldControllerInput{
					DomainName: "product",
					Actions: []types.ActionDef{
						{Name: "test", Method: tt.method, Path: "/test"},
					},
					DryRun: true,
				}

				result, err := scaffoldController(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr && result.Success {
					t.Errorf("expected failure for method %q, got success", tt.method)
				}
				if !tt.wantErr && !result.Success {
					t.Errorf("expected success for method %q, got: %s", tt.method, result.Message)
				}
			})
		}
	})

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)
		// Don't setup go.mod

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
			},
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when go.mod is missing")
		}
		if !containsString(result.Message, "module path") {
			t.Errorf("expected error about module path, got: %s", result.Message)
		}
	})

	t.Run("generates controller file", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
				{Name: "create", Method: "POST", Path: "/"},
			},
			DryRun: false,
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created
		controllerPath := tmpDir + "/internal/web/product/product.go"
		if !fileExists(controllerPath) {
			t.Errorf("expected controller file to exist at %s", controllerPath)
		}

		// Check content
		content := readFile(t, controllerPath)
		if !containsString(content, "package product") {
			t.Errorf("expected package product, got: %s", content)
		}
		if !containsString(content, "Controller") {
			t.Errorf("expected Controller struct")
		}
	})

	t.Run("handles snake_case domain name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "order_item",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
			},
			DryRun: false,
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created in correct directory
		controllerPath := tmpDir + "/internal/web/orderitem/orderitem.go"
		if !fileExists(controllerPath) {
			t.Errorf("expected controller file to exist at %s", controllerPath)
		}

		// Check content has correct model name
		content := readFile(t, controllerPath)
		if !containsString(content, "OrderItem") {
			t.Errorf("expected OrderItem model name in controller")
		}
	})

	t.Run("uses custom base path", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			BasePath:   "/api/v1/products",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
			},
			DryRun: false,
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check content uses custom base path
		controllerPath := tmpDir + "/internal/web/product/product.go"
		content := readFile(t, controllerPath)
		if !containsString(content, "/api/v1/products") {
			t.Errorf("expected custom base path in controller")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
			},
			DryRun: true,
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was NOT created
		controllerPath := tmpDir + "/internal/web/product/product.go"
		if fileExists(controllerPath) {
			t.Errorf("expected controller file to NOT exist in dry run")
		}

		// Should report what would be created
		if len(result.FilesCreated) == 0 {
			t.Errorf("expected FilesCreated to list files that would be created")
		}
	})

	t.Run("controller has HTTP handler methods", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
			},
			DryRun: false,
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check controller has handler-related content
		controllerPath := tmpDir + "/internal/web/product/product.go"
		content := readFile(t, controllerPath)

		if !containsString(content, "http.ResponseWriter") {
			t.Errorf("expected http.ResponseWriter in handler signature")
		}
		if !containsString(content, "*http.Request") {
			t.Errorf("expected *http.Request in handler signature")
		}
	})

	t.Run("includes module path in imports", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/myorg/myproject")

		input := types.ScaffoldControllerInput{
			DomainName: "product",
			Actions: []types.ActionDef{
				{Name: "list", Method: "GET", Path: "/"},
			},
			DryRun: false,
		}

		result, err := scaffoldController(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check imports use correct module path
		controllerPath := tmpDir + "/internal/web/product/product.go"
		content := readFile(t, controllerPath)
		if !containsString(content, "github.com/myorg/myproject") {
			t.Errorf("expected import with correct module path")
		}
	})
}
