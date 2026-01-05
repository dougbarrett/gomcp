package types

import (
	"encoding/json"
	"testing"
)

// boolPtr is a helper to create *bool values for tests.
func boolPtr(b bool) *bool {
	return &b
}

// TestScaffoldProjectInput_JSONUnmarshal tests JSON unmarshaling.
func TestScaffoldProjectInput_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		wantProject string
		wantModule  string
		wantDB      string
		wantAuth    bool
		wantDryRun  bool
	}{
		{
			name:        "full input",
			jsonInput:   `{"project_name":"myapp","module_path":"github.com/user/myapp","database_type":"postgres","with_auth":true,"dry_run":true}`,
			wantProject: "myapp",
			wantModule:  "github.com/user/myapp",
			wantDB:      "postgres",
			wantAuth:    true,
			wantDryRun:  true,
		},
		{
			name:        "minimal input",
			jsonInput:   `{"project_name":"test","module_path":"github.com/test/test"}`,
			wantProject: "test",
			wantModule:  "github.com/test/test",
			wantDB:      "",
			wantAuth:    false,
			wantDryRun:  false,
		},
		{
			name:        "sqlite database",
			jsonInput:   `{"project_name":"app","module_path":"mymodule","database_type":"sqlite"}`,
			wantProject: "app",
			wantModule:  "mymodule",
			wantDB:      "sqlite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input ScaffoldProjectInput
			err := json.Unmarshal([]byte(tt.jsonInput), &input)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if input.ProjectName != tt.wantProject {
				t.Errorf("ProjectName = %q, want %q", input.ProjectName, tt.wantProject)
			}
			if input.ModulePath != tt.wantModule {
				t.Errorf("ModulePath = %q, want %q", input.ModulePath, tt.wantModule)
			}
			if input.DatabaseType != tt.wantDB {
				t.Errorf("DatabaseType = %q, want %q", input.DatabaseType, tt.wantDB)
			}
			if input.WithAuth != tt.wantAuth {
				t.Errorf("WithAuth = %v, want %v", input.WithAuth, tt.wantAuth)
			}
			if input.DryRun != tt.wantDryRun {
				t.Errorf("DryRun = %v, want %v", input.DryRun, tt.wantDryRun)
			}
		})
	}
}

// TestFieldDef_JSONUnmarshal tests FieldDef JSON unmarshaling.
func TestFieldDef_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		jsonInput    string
		wantName     string
		wantType     string
		wantGORMTags string
		wantJSONTag  string
		wantFormType string
		wantRequired bool
		wantLabel    string
	}{
		{
			name:      "minimal field",
			jsonInput: `{"name":"ID","type":"uint"}`,
			wantName:  "ID",
			wantType:  "uint",
		},
		{
			name:         "full field",
			jsonInput:    `{"name":"Email","type":"string","gorm_tags":"unique;not null","json_tag":"email","form_type":"email","required":true,"label":"Email Address"}`,
			wantName:     "Email",
			wantType:     "string",
			wantGORMTags: "unique;not null",
			wantJSONTag:  "email",
			wantFormType: "email",
			wantRequired: true,
			wantLabel:    "Email Address",
		},
		{
			name:         "with gorm tags",
			jsonInput:    `{"name":"Price","type":"float64","gorm_tags":"precision:2"}`,
			wantName:     "Price",
			wantType:     "float64",
			wantGORMTags: "precision:2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var field FieldDef
			err := json.Unmarshal([]byte(tt.jsonInput), &field)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if field.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", field.Name, tt.wantName)
			}
			if field.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", field.Type, tt.wantType)
			}
			if field.GORMTags != tt.wantGORMTags {
				t.Errorf("GORMTags = %q, want %q", field.GORMTags, tt.wantGORMTags)
			}
			if field.JSONTag != tt.wantJSONTag {
				t.Errorf("JSONTag = %q, want %q", field.JSONTag, tt.wantJSONTag)
			}
			if field.FormType != tt.wantFormType {
				t.Errorf("FormType = %q, want %q", field.FormType, tt.wantFormType)
			}
			if field.Required != tt.wantRequired {
				t.Errorf("Required = %v, want %v", field.Required, tt.wantRequired)
			}
			if field.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", field.Label, tt.wantLabel)
			}
		})
	}
}

// TestScaffoldDomainInput_GetWithCrudViews tests the getter method.
func TestScaffoldDomainInput_GetWithCrudViews(t *testing.T) {
	tests := []struct {
		name  string
		input ScaffoldDomainInput
		want  bool
	}{
		{
			name:  "nil defaults to true",
			input: ScaffoldDomainInput{DomainName: "test"},
			want:  true,
		},
		{
			name:  "explicitly true",
			input: ScaffoldDomainInput{DomainName: "test", WithCrudViews: boolPtr(true)},
			want:  true,
		},
		{
			name:  "explicitly false",
			input: ScaffoldDomainInput{DomainName: "test", WithCrudViews: boolPtr(false)},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetWithCrudViews()
			if got != tt.want {
				t.Errorf("GetWithCrudViews() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestScaffoldDomainInput_GetWithSoftDelete tests the getter method.
func TestScaffoldDomainInput_GetWithSoftDelete(t *testing.T) {
	tests := []struct {
		name  string
		input ScaffoldDomainInput
		want  bool
	}{
		{
			name:  "nil defaults to true",
			input: ScaffoldDomainInput{DomainName: "test"},
			want:  true,
		},
		{
			name:  "explicitly true",
			input: ScaffoldDomainInput{DomainName: "test", WithSoftDelete: boolPtr(true)},
			want:  true,
		},
		{
			name:  "explicitly false",
			input: ScaffoldDomainInput{DomainName: "test", WithSoftDelete: boolPtr(false)},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetWithSoftDelete()
			if got != tt.want {
				t.Errorf("GetWithSoftDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestScaffoldDomainInput_JSONUnmarshal tests JSON unmarshaling with pointer bools.
func TestScaffoldDomainInput_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name           string
		jsonInput      string
		wantCrudViews  bool
		wantSoftDelete bool
	}{
		{
			name:           "defaults when not specified",
			jsonInput:      `{"domain_name":"product","fields":[{"name":"Name","type":"string"}]}`,
			wantCrudViews:  true,
			wantSoftDelete: true,
		},
		{
			name:           "explicitly false",
			jsonInput:      `{"domain_name":"product","fields":[{"name":"Name","type":"string"}],"with_crud_views":false,"with_soft_delete":false}`,
			wantCrudViews:  false,
			wantSoftDelete: false,
		},
		{
			name:           "explicitly true",
			jsonInput:      `{"domain_name":"product","fields":[{"name":"Name","type":"string"}],"with_crud_views":true,"with_soft_delete":true}`,
			wantCrudViews:  true,
			wantSoftDelete: true,
		},
		{
			name:           "mixed values",
			jsonInput:      `{"domain_name":"product","fields":[{"name":"Name","type":"string"}],"with_crud_views":true,"with_soft_delete":false}`,
			wantCrudViews:  true,
			wantSoftDelete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input ScaffoldDomainInput
			err := json.Unmarshal([]byte(tt.jsonInput), &input)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if input.GetWithCrudViews() != tt.wantCrudViews {
				t.Errorf("GetWithCrudViews() = %v, want %v", input.GetWithCrudViews(), tt.wantCrudViews)
			}
			if input.GetWithSoftDelete() != tt.wantSoftDelete {
				t.Errorf("GetWithSoftDelete() = %v, want %v", input.GetWithSoftDelete(), tt.wantSoftDelete)
			}
		})
	}
}

// TestScaffoldTableInput_GetWithPagination tests the getter method.
func TestScaffoldTableInput_GetWithPagination(t *testing.T) {
	tests := []struct {
		name  string
		input ScaffoldTableInput
		want  bool
	}{
		{
			name:  "nil defaults to true",
			input: ScaffoldTableInput{TableName: "test"},
			want:  true,
		},
		{
			name:  "explicitly true",
			input: ScaffoldTableInput{TableName: "test", WithPagination: boolPtr(true)},
			want:  true,
		},
		{
			name:  "explicitly false",
			input: ScaffoldTableInput{TableName: "test", WithPagination: boolPtr(false)},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetWithPagination()
			if got != tt.want {
				t.Errorf("GetWithPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestScaffoldTableInput_GetWithSorting tests the getter method.
func TestScaffoldTableInput_GetWithSorting(t *testing.T) {
	tests := []struct {
		name  string
		input ScaffoldTableInput
		want  bool
	}{
		{
			name:  "nil defaults to true",
			input: ScaffoldTableInput{TableName: "test"},
			want:  true,
		},
		{
			name:  "explicitly true",
			input: ScaffoldTableInput{TableName: "test", WithSorting: boolPtr(true)},
			want:  true,
		},
		{
			name:  "explicitly false",
			input: ScaffoldTableInput{TableName: "test", WithSorting: boolPtr(false)},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetWithSorting()
			if got != tt.want {
				t.Errorf("GetWithSorting() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestScaffoldTableInput_GetWithSearch tests the getter method.
func TestScaffoldTableInput_GetWithSearch(t *testing.T) {
	tests := []struct {
		name  string
		input ScaffoldTableInput
		want  bool
	}{
		{
			name:  "nil defaults to true",
			input: ScaffoldTableInput{TableName: "test"},
			want:  true,
		},
		{
			name:  "explicitly true",
			input: ScaffoldTableInput{TableName: "test", WithSearch: boolPtr(true)},
			want:  true,
		},
		{
			name:  "explicitly false",
			input: ScaffoldTableInput{TableName: "test", WithSearch: boolPtr(false)},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetWithSearch()
			if got != tt.want {
				t.Errorf("GetWithSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMethodDef_JSONUnmarshal tests MethodDef JSON unmarshaling.
func TestMethodDef_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"name":"FindByEmail","description":"Find user by email","params":[{"name":"email","type":"string"}],"returns":"*User"}`

	var method MethodDef
	err := json.Unmarshal([]byte(jsonInput), &method)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if method.Name != "FindByEmail" {
		t.Errorf("Name = %q, want %q", method.Name, "FindByEmail")
	}
	if method.Description != "Find user by email" {
		t.Errorf("Description = %q, want %q", method.Description, "Find user by email")
	}
	if len(method.Params) != 1 {
		t.Fatalf("len(Params) = %d, want 1", len(method.Params))
	}
	if method.Params[0].Name != "email" {
		t.Errorf("Params[0].Name = %q, want %q", method.Params[0].Name, "email")
	}
	if method.Returns != "*User" {
		t.Errorf("Returns = %q, want %q", method.Returns, "*User")
	}
}

// TestActionDef_JSONUnmarshal tests ActionDef JSON unmarshaling.
func TestActionDef_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name       string
		jsonInput  string
		wantName   string
		wantMethod string
		wantPath   string
		wantView   bool
	}{
		{
			name:       "list action",
			jsonInput:  `{"name":"List","method":"GET","path":"/","with_view":true}`,
			wantName:   "List",
			wantMethod: "GET",
			wantPath:   "/",
			wantView:   true,
		},
		{
			name:       "create action",
			jsonInput:  `{"name":"Create","method":"POST","path":"/"}`,
			wantName:   "Create",
			wantMethod: "POST",
			wantPath:   "/",
			wantView:   false,
		},
		{
			name:       "show action",
			jsonInput:  `{"name":"Show","method":"GET","path":"/{id}","with_view":true}`,
			wantName:   "Show",
			wantMethod: "GET",
			wantPath:   "/{id}",
			wantView:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action ActionDef
			err := json.Unmarshal([]byte(tt.jsonInput), &action)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if action.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", action.Name, tt.wantName)
			}
			if action.Method != tt.wantMethod {
				t.Errorf("Method = %q, want %q", action.Method, tt.wantMethod)
			}
			if action.Path != tt.wantPath {
				t.Errorf("Path = %q, want %q", action.Path, tt.wantPath)
			}
			if action.WithView != tt.wantView {
				t.Errorf("WithView = %v, want %v", action.WithView, tt.wantView)
			}
		})
	}
}

// TestColumnDef_JSONUnmarshal tests ColumnDef JSON unmarshaling.
func TestColumnDef_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"key":"created_at","label":"Created","sortable":true,"format":"datetime","width":"150px"}`

	var col ColumnDef
	err := json.Unmarshal([]byte(jsonInput), &col)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if col.Key != "created_at" {
		t.Errorf("Key = %q, want %q", col.Key, "created_at")
	}
	if col.Label != "Created" {
		t.Errorf("Label = %q, want %q", col.Label, "Created")
	}
	if !col.Sortable {
		t.Error("Sortable should be true")
	}
	if col.Format != "datetime" {
		t.Errorf("Format = %q, want %q", col.Format, "datetime")
	}
	if col.Width != "150px" {
		t.Errorf("Width = %q, want %q", col.Width, "150px")
	}
}

// TestRowActionDef_JSONUnmarshal tests RowActionDef JSON unmarshaling.
func TestRowActionDef_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		wantType    string
		wantLabel   string
		wantConfirm bool
	}{
		{
			name:      "view action",
			jsonInput: `{"type":"view","label":"View Details"}`,
			wantType:  "view",
			wantLabel: "View Details",
		},
		{
			name:        "delete with confirm",
			jsonInput:   `{"type":"delete","label":"Delete","confirm":true,"confirm_message":"Are you sure?"}`,
			wantType:    "delete",
			wantLabel:   "Delete",
			wantConfirm: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action RowActionDef
			err := json.Unmarshal([]byte(tt.jsonInput), &action)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if action.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", action.Type, tt.wantType)
			}
			if action.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", action.Label, tt.wantLabel)
			}
			if action.Confirm != tt.wantConfirm {
				t.Errorf("Confirm = %v, want %v", action.Confirm, tt.wantConfirm)
			}
		})
	}
}

// TestTriggerConfig_JSONUnmarshal tests TriggerConfig JSON unmarshaling.
func TestTriggerConfig_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"button_text":"Open Modal","button_variant":"primary","htmx_url":"/modal/content"}`

	var config TriggerConfig
	err := json.Unmarshal([]byte(jsonInput), &config)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if config.ButtonText != "Open Modal" {
		t.Errorf("ButtonText = %q, want %q", config.ButtonText, "Open Modal")
	}
	if config.ButtonVariant != "primary" {
		t.Errorf("ButtonVariant = %q, want %q", config.ButtonVariant, "primary")
	}
	if config.HTMXURL != "/modal/content" {
		t.Errorf("HTMXURL = %q, want %q", config.HTMXURL, "/modal/content")
	}
}

// TestPropDef_JSONUnmarshal tests PropDef JSON unmarshaling.
func TestPropDef_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"name":"title","type":"string","default":"Untitled","required":true}`

	var prop PropDef
	err := json.Unmarshal([]byte(jsonInput), &prop)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if prop.Name != "title" {
		t.Errorf("Name = %q, want %q", prop.Name, "title")
	}
	if prop.Type != "string" {
		t.Errorf("Type = %q, want %q", prop.Type, "string")
	}
	if prop.Default != "Untitled" {
		t.Errorf("Default = %q, want %q", prop.Default, "Untitled")
	}
	if !prop.Required {
		t.Error("Required should be true")
	}
}

// TestSectionDef_JSONUnmarshal tests SectionDef JSON unmarshaling.
func TestSectionDef_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"type":"hero","config":{"title":"Welcome","subtitle":"Get started"}}`

	var section SectionDef
	err := json.Unmarshal([]byte(jsonInput), &section)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if section.Type != "hero" {
		t.Errorf("Type = %q, want %q", section.Type, "hero")
	}
	if section.Config == nil {
		t.Fatal("Config should not be nil")
	}
	if section.Config["title"] != "Welcome" {
		t.Errorf("Config[title] = %v, want %q", section.Config["title"], "Welcome")
	}
}

// TestUpdateDIWiringInput_JSONUnmarshal tests UpdateDIWiringInput JSON unmarshaling.
func TestUpdateDIWiringInput_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"domains":["user","product","order"],"dry_run":true}`

	var input UpdateDIWiringInput
	err := json.Unmarshal([]byte(jsonInput), &input)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(input.Domains) != 3 {
		t.Errorf("len(Domains) = %d, want 3", len(input.Domains))
	}
	if !input.DryRun {
		t.Error("DryRun should be true")
	}
}

// TestScaffoldSeedInput_JSONUnmarshal tests ScaffoldSeedInput JSON unmarshaling.
func TestScaffoldSeedInput_JSONUnmarshal(t *testing.T) {
	jsonInput := `{"domain":"product","count":100,"with_faker":true,"dependencies":["category"]}`

	var input ScaffoldSeedInput
	err := json.Unmarshal([]byte(jsonInput), &input)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if input.Domain != "product" {
		t.Errorf("Domain = %q, want %q", input.Domain, "product")
	}
	if input.Count != 100 {
		t.Errorf("Count = %d, want %d", input.Count, 100)
	}
	if !input.WithFaker {
		t.Error("WithFaker should be true")
	}
	if len(input.Dependencies) != 1 {
		t.Errorf("len(Dependencies) = %d, want 1", len(input.Dependencies))
	}
}
