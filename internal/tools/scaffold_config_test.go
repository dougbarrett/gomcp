package tools

import (
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

func TestScaffoldConfig(t *testing.T) {
	t.Run("validates config type", func(t *testing.T) {
		tests := []struct {
			name       string
			configType string
			wantErr    bool
		}{
			{name: "page type", configType: "page", wantErr: false},
			{name: "menu type", configType: "menu", wantErr: false},
			{name: "app type", configType: "app", wantErr: false},
			{name: "messages type", configType: "messages", wantErr: false},
			{name: "invalid type", configType: "invalid", wantErr: true},
			{name: "empty type", configType: "", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, _ := testRegistry(t)
				setupGoMod(t, registry.WorkingDir, "github.com/test/project")

				input := types.ScaffoldConfigInput{
					ConfigType: tt.configType,
					Name:       "test",
					DryRun:     true,
				}

				result, err := scaffoldConfig(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr && result.Success {
					t.Errorf("expected failure for config type %q, got success", tt.configType)
				}
				if !tt.wantErr && !result.Success {
					t.Errorf("expected success for config type %q, got: %s", tt.configType, result.Message)
				}
			})
		}
	})

	t.Run("requires config name", func(t *testing.T) {
		registry, _ := testRegistry(t)
		setupGoMod(t, registry.WorkingDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "page",
			Name:       "",
			DryRun:     true,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Success {
			t.Errorf("expected failure when name is empty")
		}
		if !containsString(result.Message, "name is required") {
			t.Errorf("expected error about name, got: %s", result.Message)
		}
	})

	t.Run("validates locale", func(t *testing.T) {
		tests := []struct {
			name    string
			locale  string
			wantErr bool
		}{
			{name: "en locale", locale: "en", wantErr: false},
			{name: "empty defaults to en", locale: "", wantErr: false},
			{name: "es locale", locale: "es", wantErr: false},
			{name: "fr locale", locale: "fr", wantErr: false},
			{name: "invalid locale", locale: "invalid-locale", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				registry, _ := testRegistry(t)
				setupGoMod(t, registry.WorkingDir, "github.com/test/project")

				input := types.ScaffoldConfigInput{
					ConfigType: "page",
					Name:       "test",
					Locale:     tt.locale,
					DryRun:     true,
				}

				result, err := scaffoldConfig(registry, input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tt.wantErr && result.Success {
					t.Errorf("expected failure for locale %q, got success", tt.locale)
				}
				if !tt.wantErr && !result.Success {
					t.Errorf("expected success for locale %q, got: %s", tt.locale, result.Message)
				}
			})
		}
	})

	t.Run("generates page config", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "page",
			Name:       "products",
			Locale:     "en",
			DryRun:     false,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created
		configPath := tmpDir + "/config/en/pages/products.toml"
		if !fileExists(configPath) {
			t.Errorf("expected config file to exist at %s", configPath)
		}

		// Check content
		content := readFile(t, configPath)
		if !containsString(content, "[meta]") {
			t.Errorf("expected [meta] section in page config")
		}
		if !containsString(content, "title") {
			t.Errorf("expected title in page config")
		}
	})

	t.Run("generates menu config", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "menu",
			Name:       "main",
			Locale:     "en",
			DryRun:     false,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created
		configPath := tmpDir + "/config/en/menu.toml"
		if !fileExists(configPath) {
			t.Errorf("expected config file to exist at %s", configPath)
		}

		// Check content
		content := readFile(t, configPath)
		if !containsString(content, "[[main]]") {
			t.Errorf("expected [[main]] section in menu config")
		}
		if !containsString(content, "Dashboard") {
			t.Errorf("expected Dashboard in menu config")
		}
	})

	t.Run("generates app config", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "app",
			Name:       "settings",
			Locale:     "en",
			DryRun:     false,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created
		configPath := tmpDir + "/config/en/app.toml"
		if !fileExists(configPath) {
			t.Errorf("expected config file to exist at %s", configPath)
		}

		// Check content
		content := readFile(t, configPath)
		if !containsString(content, "[server]") {
			t.Errorf("expected [server] section in app config")
		}
		if !containsString(content, "[database]") {
			t.Errorf("expected [database] section in app config")
		}
	})

	t.Run("generates messages config", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "messages",
			Name:       "errors",
			Locale:     "en",
			DryRun:     false,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created
		configPath := tmpDir + "/config/en/messages/errors.toml"
		if !fileExists(configPath) {
			t.Errorf("expected config file to exist at %s", configPath)
		}

		// Check content
		content := readFile(t, configPath)
		if !containsString(content, "[errors]") {
			t.Errorf("expected [errors] section in messages config")
		}
	})

	t.Run("dry run does not create files", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "page",
			Name:       "dryrun",
			Locale:     "en",
			DryRun:     true,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was NOT created
		configPath := tmpDir + "/config/en/pages/dryrun.toml"
		if fileExists(configPath) {
			t.Errorf("expected config file to NOT exist in dry run")
		}

		// Should report what would be created
		if len(result.FilesCreated) == 0 {
			t.Errorf("expected FilesCreated to list files that would be created")
		}
	})

	t.Run("respects locale for path", func(t *testing.T) {
		registry, tmpDir := testRegistry(t)
		setupGoMod(t, tmpDir, "github.com/test/project")

		input := types.ScaffoldConfigInput{
			ConfigType: "page",
			Name:       "products",
			Locale:     "es",
			DryRun:     false,
		}

		result, err := scaffoldConfig(registry, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Success {
			t.Fatalf("expected success, got: %s", result.Message)
		}

		// Check file was created in correct locale directory
		configPath := tmpDir + "/config/es/pages/products.toml"
		if !fileExists(configPath) {
			t.Errorf("expected config file to exist at %s", configPath)
		}
	})
}
