package tools

import (
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldService(t *testing.T) {
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

				input := types.ScaffoldServiceInput{
					DomainName: tt.domainName,
				}

				result, err := scaffoldService(registry, input)
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

	t.Run("requires go.mod", func(t *testing.T) {
		registry, _ := testRegistry(t)
		// Don't setup go.mod

		input := types.ScaffoldServiceInput{
			DomainName: "product",
		}

		result, err := scaffoldService(registry, input)
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

	t.Run("generates service and dto files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldServiceInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldService(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check service file was created
		servicePath := tmpDir + "/internal/services/product/product.go"
		if !fileExists(servicePath) {
			t.Errorf("expected service file to exist at %s", servicePath)
		}

		// Check DTO file was created
		dtoPath := tmpDir + "/internal/services/product/dto.go"
		if !fileExists(dtoPath) {
			t.Errorf("expected dto file to exist at %s", dtoPath)
		}

		// Check service content
		content := readFile(t, servicePath)
		if !containsString(content, "package product") {
			t.Errorf("expected package product, got: %s", content)
		}
		if !containsString(content, "Service") {
			t.Errorf("expected Service interface")
		}
	})

	t.Run("handles snake_case domain name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldServiceInput{
			DomainName: "order_item",
			DryRun:     false,
		}

		result, err := scaffoldService(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created in correct directory
		servicePath := tmpDir + "/internal/services/orderitem/orderitem.go"
		if !fileExists(servicePath) {
			t.Errorf("expected service file to exist at %s", servicePath)
		}

		// Check content has correct model name
		content := readFile(t, servicePath)
		if !containsString(content, "OrderItem") {
			t.Errorf("expected OrderItem model name in service")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldServiceInput{
			DomainName: "product",
			DryRun:     true,
		}

		result, err := scaffoldService(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check files were NOT created
		servicePath := tmpDir + "/internal/services/product/product.go"
		if fileExists(servicePath) {
			t.Errorf("expected service file to NOT exist in dry run")
		}

		dtoPath := tmpDir + "/internal/services/product/dto.go"
		if fileExists(dtoPath) {
			t.Errorf("expected dto file to NOT exist in dry run")
		}

		// Should report what would be created (2 files)
		if len(result.FilesCreated) != 2 {
			t.Errorf("expected 2 files in FilesCreated, got %d", len(result.FilesCreated))
		}
	})

	t.Run("service has business logic methods", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldServiceInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldService(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check service methods
		servicePath := tmpDir + "/internal/services/product/product.go"
		content := readFile(t, servicePath)

		methods := []string{"Create", "Get", "Update", "Delete", "List"}
		for _, method := range methods {
			if !containsString(content, method) {
				t.Errorf("expected service to have %s method", method)
			}
		}
	})

	t.Run("dto has request and response types", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldServiceInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldService(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check DTO content
		dtoPath := tmpDir + "/internal/services/product/dto.go"
		content := readFile(t, dtoPath)

		// Check for input/response types
		if !containsString(content, "Input") {
			t.Errorf("expected DTO to have Input types")
		}
		if !containsString(content, "Response") {
			t.Errorf("expected DTO to have Response types")
		}
	})

	t.Run("includes module path in imports", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/myorg/myproject")

		input := types.ScaffoldServiceInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldService(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check imports use correct module path
		servicePath := tmpDir + "/internal/services/product/product.go"
		content := readFile(t, servicePath)
		if !containsString(content, "github.com/myorg/myproject") {
			t.Errorf("expected import with correct module path")
		}
	})
}
