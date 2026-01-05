package tools

import (
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldRepository(t *testing.T) {
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

				input := types.ScaffoldRepositoryInput{
					DomainName: tt.domainName,
				}

				result, err := scaffoldRepository(registry, input)
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

		input := types.ScaffoldRepositoryInput{
			DomainName: "product",
		}

		result, err := scaffoldRepository(registry, input)
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

	t.Run("generates repository file", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldRepositoryInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldRepository(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created
		repoPath := tmpDir + "/internal/repository/product/product.go"
		if !fileExists(repoPath) {
			t.Errorf("expected repository file to exist at %s", repoPath)
		}

		// Check content
		content := readFile(t, repoPath)
		if !containsString(content, "package product") {
			t.Errorf("expected package product, got: %s", content)
		}
		if !containsString(content, "Repository") {
			t.Errorf("expected Repository interface")
		}
		if !containsString(content, "NewRepository") {
			t.Errorf("expected NewRepository function")
		}
	})

	t.Run("uses custom model name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldRepositoryInput{
			DomainName: "product",
			ModelName:  "Item",
			DryRun:     false,
		}

		result, err := scaffoldRepository(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check content uses custom model name
		repoPath := tmpDir + "/internal/repository/product/product.go"
		content := readFile(t, repoPath)
		if !containsString(content, "Item") {
			t.Errorf("expected custom model name 'Item' in repository")
		}
	})

	t.Run("handles snake_case domain name", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldRepositoryInput{
			DomainName: "order_item",
			DryRun:     false,
		}

		result, err := scaffoldRepository(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created in correct directory
		repoPath := tmpDir + "/internal/repository/orderitem/orderitem.go"
		if !fileExists(repoPath) {
			t.Errorf("expected repository file to exist at %s", repoPath)
		}

		// Check content has correct model name
		content := readFile(t, repoPath)
		if !containsString(content, "OrderItem") {
			t.Errorf("expected OrderItem model name in repository")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldRepositoryInput{
			DomainName: "product",
			DryRun:     true,
		}

		result, err := scaffoldRepository(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was NOT created
		repoPath := tmpDir + "/internal/repository/product/product.go"
		if fileExists(repoPath) {
			t.Errorf("expected repository file to NOT exist in dry run")
		}

		// Should report what would be created
		if len(result.FilesCreated) == 0 {
			t.Errorf("expected FilesCreated to list files that would be created")
		}
	})

	t.Run("repository has CRUD methods", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldRepositoryInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldRepository(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check CRUD methods
		repoPath := tmpDir + "/internal/repository/product/product.go"
		content := readFile(t, repoPath)

		methods := []string{"Create", "FindByID", "Update", "Delete", "FindAll"}
		for _, method := range methods {
			if !containsString(content, method) {
				t.Errorf("expected repository to have %s method", method)
			}
		}
	})

	t.Run("includes module path in imports", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/myorg/myproject")

		input := types.ScaffoldRepositoryInput{
			DomainName: "product",
			DryRun:     false,
		}

		result, err := scaffoldRepository(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check imports use correct module path
		repoPath := tmpDir + "/internal/repository/product/product.go"
		content := readFile(t, repoPath)
		if !containsString(content, "github.com/myorg/myproject/internal/models") {
			t.Errorf("expected import with correct module path")
		}
	})
}
