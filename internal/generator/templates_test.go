package generator

import (
	"embed"
	"strings"
	"testing"
)

//go:embed testdata/*.tmpl
var templatesTestFS embed.FS

// TestParseTemplate tests template parsing with custom delimiters.
func TestParseTemplate(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		data     interface{}
		wantErr  bool
		contains string
	}{
		{
			name:     "simple variable",
			content:  "Hello, [[.Name]]!",
			data:     map[string]string{"Name": "World"},
			wantErr:  false,
			contains: "Hello, World!",
		},
		{
			name:     "multiple variables",
			content:  "[[.First]] [[.Last]]",
			data:     map[string]string{"First": "John", "Last": "Doe"},
			wantErr:  false,
			contains: "John Doe",
		},
		{
			name:     "with standard go braces",
			content:  "func main() { fmt.Println([[.Msg]]) }",
			data:     map[string]string{"Msg": `"hello"`},
			wantErr:  false,
			contains: `func main() { fmt.Println("hello") }`,
		},
		{
			name:     "empty template",
			content:  "",
			data:     nil,
			wantErr:  false,
			contains: "",
		},
		{
			name:     "no variables",
			content:  "static content {{ not a template }}",
			data:     nil,
			wantErr:  false,
			contains: "{{ not a template }}",
		},
		{
			name:    "syntax error",
			content: "[[.Name",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "unclosed action",
			content: "[[ if .X ]] no end",
			data:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := ParseTemplate(tt.name, tt.content)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseTemplate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseTemplate() error = %v", err)
			}

			if tmpl == nil {
				t.Fatal("ParseTemplate() returned nil template")
			}
		})
	}
}

// TestExecuteTemplateString tests template execution from strings.
func TestExecuteTemplateString(t *testing.T) {
	tests := []struct {
		name    string
		content string
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "simple variable",
			content: "Hello, [[.Name]]!",
			data:    map[string]string{"Name": "World"},
			want:    "Hello, World!",
		},
		{
			name:    "struct data",
			content: "[[.First]] [[.Last]]",
			data:    struct{ First, Last string }{"John", "Doe"},
			want:    "John Doe",
		},
		{
			name:    "with range",
			content: "[[range .Items]][[.]] [[end]]",
			data:    map[string][]string{"Items": {"a", "b", "c"}},
			want:    "a b c ",
		},
		{
			name:    "with if",
			content: "[[if .Show]]visible[[end]]",
			data:    map[string]bool{"Show": true},
			want:    "visible",
		},
		{
			name:    "with if false",
			content: "[[if .Show]]visible[[end]]",
			data:    map[string]bool{"Show": false},
			want:    "",
		},
		{
			name:    "preserve go braces",
			content: "type Config struct {\n\tName string `json:\"name\"`\n}",
			data:    nil,
			want:    "type Config struct {\n\tName string `json:\"name\"`\n}",
		},
		{
			name:    "helper function - toLower",
			content: "[[.Name | toLower]]",
			data:    map[string]string{"Name": "HELLO"},
			want:    "hello",
		},
		{
			name:    "helper function - toUpper",
			content: "[[.Name | toUpper]]",
			data:    map[string]string{"Name": "hello"},
			want:    "HELLO",
		},
		{
			name:    "helper function - toPascalCase",
			content: "[[.Name | toPascalCase]]",
			data:    map[string]string{"Name": "user_profile"},
			want:    "UserProfile",
		},
		{
			name:    "helper function - toCamelCase",
			content: "[[.Name | toCamelCase]]",
			data:    map[string]string{"Name": "user_profile"},
			want:    "userProfile",
		},
		{
			name:    "helper function - toSnakeCase",
			content: "[[.Name | toSnakeCase]]",
			data:    map[string]string{"Name": "UserProfile"},
			want:    "user_profile",
		},
		{
			name:    "helper function - pluralize",
			content: "[[.Name | pluralize]]",
			data:    map[string]string{"Name": "product"},
			want:    "products",
		},
		{
			name:    "missing field in map returns no value",
			content: "[[.Missing]]",
			data:    map[string]string{"Name": "test"},
			want:    "<no value>", // Go templates return "<no value>" for missing map keys
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecuteTemplateString(tt.name, tt.content, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Error("ExecuteTemplateString() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ExecuteTemplateString() error = %v", err)
			}

			if result != tt.want {
				t.Errorf("ExecuteTemplateString() = %q, want %q", result, tt.want)
			}
		})
	}
}

// TestLoadTemplate tests template loading from embedded FS.
func TestLoadTemplate(t *testing.T) {
	// Test loading existing template
	tmpl, err := LoadTemplate(templatesTestFS, "testdata/simple.tmpl")
	if err != nil {
		t.Fatalf("LoadTemplate() error = %v", err)
	}
	if tmpl == nil {
		t.Fatal("LoadTemplate() returned nil")
	}

	// Test loading non-existent template
	_, err = LoadTemplate(templatesTestFS, "testdata/nonexistent.tmpl")
	if err == nil {
		t.Error("LoadTemplate() should error for non-existent template")
	}
}

// TestExecuteTemplate tests template execution from embedded FS.
func TestExecuteTemplate(t *testing.T) {
	result, err := ExecuteTemplate(templatesTestFS, "testdata/simple.tmpl", map[string]string{"Name": "Test"})
	if err != nil {
		t.Fatalf("ExecuteTemplate() error = %v", err)
	}

	expected := "Hello, Test!\n"
	if result != expected {
		t.Errorf("ExecuteTemplate() = %q, want %q", result, expected)
	}
}

// TestMustParseTemplate tests must versions.
func TestMustParseTemplate(t *testing.T) {
	// Valid template should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustParseTemplate panicked unexpectedly: %v", r)
		}
	}()

	tmpl := MustParseTemplate("test", "Hello, [[.Name]]!")
	if tmpl == nil {
		t.Error("MustParseTemplate returned nil")
	}
}

// TestMustParseTemplate_Panic tests that invalid templates cause panic.
func TestMustParseTemplate_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseTemplate should panic on invalid template")
		}
	}()

	MustParseTemplate("invalid", "[[.Name")
}

// TestMustLoadTemplate tests must load versions.
func TestMustLoadTemplate(t *testing.T) {
	// Valid template should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustLoadTemplate panicked unexpectedly: %v", r)
		}
	}()

	tmpl := MustLoadTemplate(templatesTestFS, "testdata/simple.tmpl")
	if tmpl == nil {
		t.Error("MustLoadTemplate returned nil")
	}
}

// TestMustLoadTemplate_Panic tests that missing templates cause panic.
func TestMustLoadTemplate_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLoadTemplate should panic on missing template")
		}
	}()

	MustLoadTemplate(templatesTestFS, "testdata/nonexistent.tmpl")
}

// TestTemplateExists tests template existence check.
func TestTemplateExists(t *testing.T) {
	if !TemplateExists(templatesTestFS, "testdata/simple.tmpl") {
		t.Error("TemplateExists should return true for existing template")
	}

	if TemplateExists(templatesTestFS, "testdata/nonexistent.tmpl") {
		t.Error("TemplateExists should return false for non-existent template")
	}
}

// TestListTemplates tests template listing.
func TestListTemplates(t *testing.T) {
	templates, err := ListTemplates(templatesTestFS, "testdata")
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(templates) == 0 {
		t.Error("ListTemplates() should return at least one template")
	}

	// Check that simple.tmpl is in the list
	found := false
	for _, tmpl := range templates {
		if strings.HasSuffix(tmpl, "simple.tmpl") {
			found = true
			break
		}
	}
	if !found {
		t.Error("ListTemplates() should include simple.tmpl")
	}
}

// TestListTemplates_NonExistent tests listing from non-existent directory.
func TestListTemplates_NonExistent(t *testing.T) {
	_, err := ListTemplates(templatesTestFS, "nonexistent")
	if err == nil {
		t.Error("ListTemplates() should error for non-existent directory")
	}
}

// TestGetTemplateInfo tests getting template information.
func TestGetTemplateInfo(t *testing.T) {
	info, err := GetTemplateInfo(templatesTestFS, "testdata/simple.tmpl")
	if err != nil {
		t.Fatalf("GetTemplateInfo() error = %v", err)
	}

	if info.Name != "testdata/simple.tmpl" {
		t.Errorf("info.Name = %q, want %q", info.Name, "testdata/simple.tmpl")
	}

	if info.Size == 0 {
		t.Error("info.Size should not be 0")
	}

	if !strings.Contains(info.Content, "Hello") {
		t.Error("info.Content should contain template text")
	}
}

// TestGetTemplateInfo_NonExistent tests getting info for non-existent template.
func TestGetTemplateInfo_NonExistent(t *testing.T) {
	_, err := GetTemplateInfo(templatesTestFS, "testdata/nonexistent.tmpl")
	if err == nil {
		t.Error("GetTemplateInfo() should error for non-existent template")
	}
}

// TestDelimiterConstants tests delimiter constant values.
func TestDelimiterConstants(t *testing.T) {
	if LeftDelim != "[[" {
		t.Errorf("LeftDelim = %q, want %q", LeftDelim, "[[")
	}
	if RightDelim != "]]" {
		t.Errorf("RightDelim = %q, want %q", RightDelim, "]]")
	}
}

// TestComplexTemplate tests a more complex template scenario.
func TestComplexTemplate(t *testing.T) {
	content := `package [[.PackageName]]

import (
	"[[.ModulePath]]/internal/models"
)

type [[.ServiceName]] struct {
	repo [[.RepoType]]
}

func New[[.ServiceName]](repo [[.RepoType]]) *[[.ServiceName]] {
	return &[[.ServiceName]]{repo: repo}
}

[[range .Methods]]
func (s *[[$.ServiceName]]) [[.Name]]([[.Params]]) [[.Returns]] {
	// TODO: implement
}
[[end]]`

	data := map[string]interface{}{
		"PackageName": "product",
		"ModulePath":  "github.com/example/app",
		"ServiceName": "ProductService",
		"RepoType":    "*repository.ProductRepository",
		"Methods": []map[string]string{
			{"Name": "Create", "Params": "ctx context.Context, dto CreateDTO", "Returns": "(*models.Product, error)"},
			{"Name": "GetByID", "Params": "ctx context.Context, id uint", "Returns": "(*models.Product, error)"},
		},
	}

	result, err := ExecuteTemplateString("complex", content, data)
	if err != nil {
		t.Fatalf("ExecuteTemplateString() error = %v", err)
	}

	// Verify key elements
	checks := []string{
		"package product",
		"github.com/example/app/internal/models",
		"type ProductService struct",
		"func NewProductService",
		"func (s *ProductService) Create",
		"func (s *ProductService) GetByID",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Result should contain %q", check)
		}
	}
}
