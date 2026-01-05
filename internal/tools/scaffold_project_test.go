package tools

import (
	"os"
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldProject(t *testing.T) {
	t.Run("validates project name", func(t *testing.T) {
		tests := []struct {
			name        string
			projectName string
			wantErr     string
		}{
			{
				name:        "empty project name",
				projectName: "",
				wantErr:     "project name is required",
			},
			{
				name:        "invalid chars",
				projectName: "my@project",
				wantErr:     "project name must start with a letter",
			},
			{
				name:        "starts with number",
				projectName: "123project",
				wantErr:     "project name must start with a letter",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, _ := testRegistry(t)
				input := types.ScaffoldProjectInput{
					ProjectName: tt.projectName,
					ModulePath:  "github.com/test/project",
				}

				result, err := scaffoldProject(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Errorf("expected failure, got success")
				}
				if tt.wantErr != "" && !containsString(result.Message, tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, result.Message)
				}
			})
		}
	})

	t.Run("validates module path", func(t *testing.T) {
		tests := []struct {
			name       string
			modulePath string
			wantErr    string
		}{
			{
				name:       "empty module path",
				modulePath: "",
				wantErr:    "module path is required",
			},
			{
				name:       "invalid module path",
				modulePath: "not a valid module",
				wantErr:    "cannot contain whitespace",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, _ := testRegistry(t)
				input := types.ScaffoldProjectInput{
					ProjectName: "testproject",
					ModulePath:  tt.modulePath,
				}

				result, err := scaffoldProject(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if result.Success {
					t.Errorf("expected failure, got success")
				}
				if tt.wantErr != "" && !containsString(result.Message, tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, result.Message)
				}
			})
		}
	})

	t.Run("validates database type", func(t *testing.T) {
		tests := []struct {
			name    string
			dbType  string
			wantErr bool
		}{
			{name: "sqlite", dbType: "sqlite", wantErr: false},
			{name: "postgres", dbType: "postgres", wantErr: false},
			{name: "mysql", dbType: "mysql", wantErr: false},
			{name: "empty defaults to sqlite", dbType: "", wantErr: false},
			{name: "invalid db type", dbType: "mongodb", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, _ := testRegistry(t)
				input := types.ScaffoldProjectInput{
					ProjectName:  "testproject",
					ModulePath:   "github.com/test/project",
					DatabaseType: tt.dbType,
					DryRun:       true,
				}

				result, err := scaffoldProject(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr && result.Success {
					t.Errorf("expected failure for db type %q, got success", tt.dbType)
				}
				if !tt.wantErr && !result.Success {
					t.Errorf("expected success for db type %q, got: %s", tt.dbType, result.Message)
				}
			})
		}
	})

	t.Run("generates project files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName:  "myapp",
			ModulePath:   "github.com/test/myapp",
			DatabaseType: "sqlite",
			DryRun:       false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check expected files were created
		expectedFiles := []string{
			"go.mod",
			"cmd/web/main.go",
			"cmd/seed/main.go",
			"internal/config/config.go",
			"internal/database/database.go",
			"internal/models/base.go",
			"internal/web/router.go",
			"internal/web/middleware/middleware.go",
			"internal/web/response.go",
			"internal/web/layouts/base.templ",
			"internal/web/components/common.templ",
			"Taskfile.yml",
			".air.toml",
			"assets/css/input.css",
			"config/en/app.toml",
			"config/en/menu.toml",
			".gitignore",
		}

		projectDir := tmpDir + "/myapp"
		for _, f := range expectedFiles {
			path := projectDir + "/" + f
			if !fileExists(path) {
				t.Errorf("expected file %s to exist", f)
			}
		}
	})

	t.Run("creates directory structure", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/test/myapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check expected directories were created
		expectedDirs := []string{
			"cmd/web",
			"cmd/seed",
			"internal/config",
			"internal/database",
			"internal/models",
			"internal/repository",
			"internal/services",
			"internal/web/middleware",
			"internal/web/layouts",
			"internal/web/components",
			"config/en/pages",
			"assets/css",
			"assets/js",
			"components",
			"utils",
		}

		projectDir := tmpDir + "/myapp"
		for _, d := range expectedDirs {
			path := projectDir + "/" + d
			if !dirExists(path) {
				t.Errorf("expected directory %s to exist", d)
			}
		}
	})

	t.Run("go.mod contains correct module path", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/myorg/myapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		goModPath := tmpDir + "/myapp/go.mod"
		content := readFile(t, goModPath)
		if !containsString(content, "module github.com/myorg/myapp") {
			t.Errorf("go.mod should contain correct module path, got: %s", content)
		}
	})

	t.Run("main.go imports correct module", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/myorg/myapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		mainPath := tmpDir + "/myapp/cmd/web/main.go"
		content := readFile(t, mainPath)
		if !containsString(content, "github.com/myorg/myapp/internal/config") {
			t.Errorf("main.go should import config package, got: %s", content)
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "dryrunapp",
			ModulePath:  "github.com/test/dryrunapp",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Project directory should not exist
		projectDir := tmpDir + "/dryrunapp"
		if dirExists(projectDir) {
			t.Errorf("expected project directory to NOT exist in dry run")
		}

		// Should report what would be created
		if len(result.FilesCreated) == 0 {
			t.Errorf("expected FilesCreated to list files that would be created")
		}
	})

	t.Run("returns next steps", func(t *testing.T) {
		registry, _ := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/test/myapp",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		if len(result.NextSteps) == 0 {
			t.Errorf("expected NextSteps to be populated")
		}

		// Should include cd and go mod tidy
		foundCD := false
		foundGoModTidy := false
		for _, step := range result.NextSteps {
			if containsString(step, "cd myapp") {
				foundCD = true
			}
			if containsString(step, "go mod tidy") {
				foundGoModTidy = true
			}
		}
		if !foundCD {
			t.Errorf("expected NextSteps to include 'cd myapp'")
		}
		if !foundGoModTidy {
			t.Errorf("expected NextSteps to include 'go mod tidy'")
		}
	})

	t.Run("rejects existing project directory", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)

		// Create an existing project directory
		existingDir := tmpDir + "/existingapp"
		if err := os.MkdirAll(existingDir, 0755); err != nil {
			t.Fatalf("failed to create test directory: %v", err)
		}

		input := types.ScaffoldProjectInput{
			ProjectName: "existingapp",
			ModulePath:  "github.com/test/existingapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when project directory already exists")
		}
		if !containsString(result.Message, "already exists") {
			t.Errorf("expected error message about existing directory, got: %s", result.Message)
		}
	})

	t.Run("with_auth generates auth files", func(t *testing.T) {
		registry, _ := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "authapp",
			ModulePath:  "github.com/test/authapp",
			WithAuth:    true,
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check that auth files are included
		authFiles := []string{
			"internal/models/user.go",
			"internal/repository/user/user.go",
			"internal/services/auth/auth.go",
			"internal/services/auth/session.go",
			"internal/web/middleware/auth.go",
			"internal/web/auth/auth.go",
			"internal/web/auth/views/layout.templ",
			"internal/web/auth/views/login.templ",
			"internal/web/auth/views/register.templ",
		}

		for _, authFile := range authFiles {
			found := false
			for _, createdFile := range result.FilesCreated {
				if containsString(createdFile, authFile) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected auth file %q to be created, files: %v", authFile, result.FilesCreated)
			}
		}

		// Should have base files (19) + auth files (9) = 28 files
		expectedFileCount := 28
		if len(result.FilesCreated) != expectedFileCount {
			t.Errorf("expected %d files with auth, got %d", expectedFileCount, len(result.FilesCreated))
		}
	})

	t.Run("without_auth does not generate auth files", func(t *testing.T) {
		registry, _ := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "noauthapp",
			ModulePath:  "github.com/test/noauthapp",
			WithAuth:    false,
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check that auth files are NOT included
		for _, createdFile := range result.FilesCreated {
			if containsString(createdFile, "internal/models/user.go") ||
				containsString(createdFile, "internal/repository/user") ||
				containsString(createdFile, "internal/services/auth") ||
				containsString(createdFile, "internal/web/auth") {
				t.Errorf("unexpected auth file %q when WithAuth is false", createdFile)
			}
		}
	})

	t.Run("creates correct number of files", func(t *testing.T) {
		registry, _ := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "myapp",
			ModulePath:  "github.com/test/myapp",
			DryRun:      true,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Should have 19 files based on the template list (including tailwind.config.js and output.css)
		expectedFileCount := 19
		if len(result.FilesCreated) != expectedFileCount {
			t.Errorf("expected %d files, got %d: %v", expectedFileCount, len(result.FilesCreated), result.FilesCreated)
		}
	})

	t.Run("middleware includes CSRF protection", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "csrfapp",
			ModulePath:  "github.com/test/csrfapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check middleware.go includes CSRF middleware
		middlewarePath := tmpDir + "/csrfapp/internal/web/middleware/middleware.go"
		content := readFile(t, middlewarePath)

		if !containsString(content, "CSRF") {
			t.Errorf("middleware.go should include CSRF function")
		}
		if !containsString(content, "gorilla/csrf") {
			t.Errorf("middleware.go should import gorilla/csrf")
		}
		if !containsString(content, "GetCSRFToken") {
			t.Errorf("middleware.go should include GetCSRFToken function")
		}
	})

	t.Run("router includes CSRF middleware", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "csrfrouterapp",
			ModulePath:  "github.com/test/csrfrouterapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check router.go includes CSRF middleware usage
		routerPath := tmpDir + "/csrfrouterapp/internal/web/router.go"
		content := readFile(t, routerPath)

		if !containsString(content, "middleware.CSRF") {
			t.Errorf("router.go should use middleware.CSRF")
		}
		if !containsString(content, "middleware.InjectCSRFToken") {
			t.Errorf("router.go should use middleware.InjectCSRFToken")
		}
	})

	t.Run("config includes session config", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "sessioncfgapp",
			ModulePath:  "github.com/test/sessioncfgapp",
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check config.go includes session configuration
		configPath := tmpDir + "/sessioncfgapp/internal/config/config.go"
		content := readFile(t, configPath)

		if !containsString(content, "SessionConfig") {
			t.Errorf("config.go should include SessionConfig struct")
		}
		if !containsString(content, "SESSION_SECRET") {
			t.Errorf("config.go should read SESSION_SECRET from environment")
		}
	})

	t.Run("auth templates include CSRF token", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		input := types.ScaffoldProjectInput{
			ProjectName: "authcsrfapp",
			ModulePath:  "github.com/test/authcsrfapp",
			WithAuth:    true,
			DryRun:      false,
		}

		result, err := scaffoldProject(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check login.templ includes CSRF token
		loginPath := tmpDir + "/authcsrfapp/internal/web/auth/views/login.templ"
		loginContent := readFile(t, loginPath)

		if !containsString(loginContent, "CSRFToken") {
			t.Errorf("login.templ should include CSRFToken in props")
		}
		if !containsString(loginContent, "X-CSRF-Token") {
			t.Errorf("login.templ should include X-CSRF-Token header for HTMX")
		}

		// Check register.templ includes CSRF token
		registerPath := tmpDir + "/authcsrfapp/internal/web/auth/views/register.templ"
		registerContent := readFile(t, registerPath)

		if !containsString(registerContent, "CSRFToken") {
			t.Errorf("register.templ should include CSRFToken in props")
		}

		// Check auth controller passes CSRF token
		controllerPath := tmpDir + "/authcsrfapp/internal/web/auth/auth.go"
		controllerContent := readFile(t, controllerPath)

		if !containsString(controllerContent, "GetCSRFToken") {
			t.Errorf("auth controller should call GetCSRFToken")
		}
	})
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
