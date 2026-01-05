package generator

import (
	"testing"

	"github.com/dbb1dev/go-mcp/internal/types"
)

// TestNewProjectData tests ProjectData creation.
func TestNewProjectData(t *testing.T) {
	tests := []struct {
		name     string
		input    types.ScaffoldProjectInput
		wantDB   string
		wantAuth bool
	}{
		{
			name: "default database",
			input: types.ScaffoldProjectInput{
				ProjectName: "myapp",
				ModulePath:  "github.com/user/myapp",
			},
			wantDB:   "sqlite",
			wantAuth: false,
		},
		{
			name: "postgres with auth",
			input: types.ScaffoldProjectInput{
				ProjectName:  "myapp",
				ModulePath:   "github.com/user/myapp",
				DatabaseType: "postgres",
				WithAuth:     true,
			},
			wantDB:   "postgres",
			wantAuth: true,
		},
		{
			name: "mysql",
			input: types.ScaffoldProjectInput{
				ProjectName:  "myapp",
				ModulePath:   "github.com/user/myapp",
				DatabaseType: "mysql",
			},
			wantDB:   "mysql",
			wantAuth: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewProjectData(tt.input)

			if data.ProjectName != tt.input.ProjectName {
				t.Errorf("ProjectName = %q, want %q", data.ProjectName, tt.input.ProjectName)
			}
			if data.ModulePath != tt.input.ModulePath {
				t.Errorf("ModulePath = %q, want %q", data.ModulePath, tt.input.ModulePath)
			}
			if data.DatabaseType != tt.wantDB {
				t.Errorf("DatabaseType = %q, want %q", data.DatabaseType, tt.wantDB)
			}
			if data.WithAuth != tt.wantAuth {
				t.Errorf("WithAuth = %v, want %v", data.WithAuth, tt.wantAuth)
			}
		})
	}
}

// TestNewFieldData tests FieldData creation.
func TestNewFieldData(t *testing.T) {
	tests := []struct {
		name         string
		input        types.FieldDef
		wantJSONTag  string
		wantLabel    string
		wantFormType string
	}{
		{
			name: "basic field",
			input: types.FieldDef{
				Name: "UserName",
				Type: "string",
			},
			wantJSONTag:  "user_name",
			wantLabel:    "User Name",
			wantFormType: "input",
		},
		{
			name: "with custom json tag",
			input: types.FieldDef{
				Name:    "UserName",
				Type:    "string",
				JSONTag: "username",
			},
			wantJSONTag:  "username",
			wantLabel:    "User Name",
			wantFormType: "input",
		},
		{
			name: "with custom label",
			input: types.FieldDef{
				Name:  "UserName",
				Type:  "string",
				Label: "Full Name",
			},
			wantJSONTag:  "user_name",
			wantLabel:    "Full Name",
			wantFormType: "input",
		},
		{
			name: "with explicit form type",
			input: types.FieldDef{
				Name:     "Description",
				Type:     "string",
				FormType: "textarea",
			},
			wantJSONTag:  "description",
			wantLabel:    "Description",
			wantFormType: "textarea",
		},
		{
			name: "integer field",
			input: types.FieldDef{
				Name: "Count",
				Type: "int",
			},
			wantJSONTag:  "count",
			wantLabel:    "Count",
			wantFormType: "number",
		},
		{
			name: "boolean field",
			input: types.FieldDef{
				Name: "IsActive",
				Type: "bool",
			},
			wantJSONTag:  "is_active",
			wantLabel:    "Is Active",
			wantFormType: "checkbox",
		},
		{
			name: "time field",
			input: types.FieldDef{
				Name: "CreatedAt",
				Type: "time.Time",
			},
			wantJSONTag:  "created_at",
			wantLabel:    "Created At",
			wantFormType: "datetime",
		},
		{
			name: "required field",
			input: types.FieldDef{
				Name:     "Email",
				Type:     "string",
				Required: true,
			},
			wantJSONTag:  "email",
			wantLabel:    "Email",
			wantFormType: "input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewFieldData(tt.input)

			if data.Name != tt.input.Name {
				t.Errorf("Name = %q, want %q", data.Name, tt.input.Name)
			}
			if data.Type != tt.input.Type {
				t.Errorf("Type = %q, want %q", data.Type, tt.input.Type)
			}
			if data.JSONName != tt.wantJSONTag {
				t.Errorf("JSONName = %q, want %q", data.JSONName, tt.wantJSONTag)
			}
			if data.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", data.Label, tt.wantLabel)
			}
			if data.FormType != tt.wantFormType {
				t.Errorf("FormType = %q, want %q", data.FormType, tt.wantFormType)
			}
			if data.Required != tt.input.Required {
				t.Errorf("Required = %v, want %v", data.Required, tt.input.Required)
			}
			// Omitempty should be opposite of Required
			if data.Omitempty == tt.input.Required {
				t.Errorf("Omitempty = %v, should be opposite of Required (%v)", data.Omitempty, tt.input.Required)
			}
		})
	}
}

// TestNewFieldDataList tests creating list of FieldData.
func TestNewFieldDataList(t *testing.T) {
	fields := []types.FieldDef{
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
		{Name: "Active", Type: "bool"},
	}

	result := NewFieldDataList(fields)

	if len(result) != 3 {
		t.Fatalf("len(result) = %d, want 3", len(result))
	}

	if result[0].Name != "Name" {
		t.Errorf("result[0].Name = %q, want %q", result[0].Name, "Name")
	}
	if result[1].Name != "Age" {
		t.Errorf("result[1].Name = %q, want %q", result[1].Name, "Age")
	}
	if result[2].Name != "Active" {
		t.Errorf("result[2].Name = %q, want %q", result[2].Name, "Active")
	}
}

// TestNewFieldDataList_Empty tests empty field list.
func TestNewFieldDataList_Empty(t *testing.T) {
	result := NewFieldDataList(nil)
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}

	result = NewFieldDataList([]types.FieldDef{})
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0 for empty slice", len(result))
	}
}

// TestInferFormType tests form type inference.
func TestInferFormType(t *testing.T) {
	tests := []struct {
		goType   string
		wantForm string
	}{
		{"string", "input"},
		{"int", "number"},
		{"int8", "number"},
		{"int16", "number"},
		{"int32", "number"},
		{"int64", "number"},
		{"uint", "number"},
		{"uint8", "number"},
		{"uint16", "number"},
		{"uint32", "number"},
		{"uint64", "number"},
		{"float32", "number"},
		{"float64", "number"},
		{"bool", "checkbox"},
		{"time.Time", "datetime"},
		{"*time.Time", "datetime"},
		{"CustomType", "input"},
		{"[]byte", "input"},
	}

	for _, tt := range tests {
		t.Run(tt.goType, func(t *testing.T) {
			result := inferFormType(tt.goType)
			if result != tt.wantForm {
				t.Errorf("inferFormType(%q) = %q, want %q", tt.goType, result, tt.wantForm)
			}
		})
	}
}

// boolPtr is a helper to create *bool values for tests.
func boolPtr(b bool) *bool {
	return &b
}

// TestNewDomainData tests DomainData creation.
func TestNewDomainData(t *testing.T) {
	input := types.ScaffoldDomainInput{
		DomainName: "user_profile",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string", Required: true},
			{Name: "Email", Type: "string", Required: true},
			{Name: "Age", Type: "int"},
		},
		WithSoftDelete: boolPtr(true),
		WithCrudViews:  boolPtr(true),
	}

	data := NewDomainData(input, "github.com/user/app")

	// Check basic naming
	if data.DomainName != "user_profile" {
		t.Errorf("DomainName = %q, want %q", data.DomainName, "user_profile")
	}
	if data.ModelName != "UserProfile" {
		t.Errorf("ModelName = %q, want %q", data.ModelName, "UserProfile")
	}
	if data.PackageName != "userprofile" {
		t.Errorf("PackageName = %q, want %q", data.PackageName, "userprofile")
	}
	if data.VariableName != "userProfile" {
		t.Errorf("VariableName = %q, want %q", data.VariableName, "userProfile")
	}
	if data.TableName != "user_profiles" {
		t.Errorf("TableName = %q, want %q", data.TableName, "user_profiles")
	}
	if data.URLPath != "/user-profiles" {
		t.Errorf("URLPath = %q, want %q", data.URLPath, "/user-profiles")
	}

	// Check module path
	if data.ModulePath != "github.com/user/app" {
		t.Errorf("ModulePath = %q, want %q", data.ModulePath, "github.com/user/app")
	}

	// Check fields
	if len(data.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(data.Fields))
	}

	// Check options
	if !data.WithSoftDelete {
		t.Error("WithSoftDelete should be true")
	}
	if !data.WithCrudViews {
		t.Error("WithCrudViews should be true")
	}
}

// TestNewDomainData_Defaults tests DomainData with default values.
func TestNewDomainData_Defaults(t *testing.T) {
	input := types.ScaffoldDomainInput{
		DomainName: "product",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string"},
		},
		// WithSoftDelete and WithCrudViews are nil (not set)
		// The GetWithSoftDelete() and GetWithCrudViews() methods return true by default
	}

	data := NewDomainData(input, "github.com/example/app")

	// When nil, the default is true (as per the getter methods)
	if !data.WithSoftDelete {
		t.Error("WithSoftDelete should be true by default (nil -> true)")
	}
	if !data.WithCrudViews {
		t.Error("WithCrudViews should be true by default (nil -> true)")
	}
}

// TestNewColumnData tests ColumnData creation.
func TestNewColumnData(t *testing.T) {
	tests := []struct {
		name      string
		input     types.ColumnDef
		wantLabel string
	}{
		{
			name: "with custom label",
			input: types.ColumnDef{
				Key:      "user_name",
				Label:    "Full Name",
				Sortable: true,
			},
			wantLabel: "Full Name",
		},
		{
			name: "auto-generate label",
			input: types.ColumnDef{
				Key:      "user_name",
				Sortable: false,
			},
			wantLabel: "User Name",
		},
		{
			name: "with format and width",
			input: types.ColumnDef{
				Key:    "created_at",
				Format: "datetime",
				Width:  "200px",
			},
			wantLabel: "Created At",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewColumnData(tt.input)

			if data.Key != tt.input.Key {
				t.Errorf("Key = %q, want %q", data.Key, tt.input.Key)
			}
			if data.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", data.Label, tt.wantLabel)
			}
			if data.Sortable != tt.input.Sortable {
				t.Errorf("Sortable = %v, want %v", data.Sortable, tt.input.Sortable)
			}
			if data.Format != tt.input.Format {
				t.Errorf("Format = %q, want %q", data.Format, tt.input.Format)
			}
			if data.Width != tt.input.Width {
				t.Errorf("Width = %q, want %q", data.Width, tt.input.Width)
			}
		})
	}
}

// TestNewColumnDataList tests creating list of ColumnData.
func TestNewColumnDataList(t *testing.T) {
	cols := []types.ColumnDef{
		{Key: "name", Label: "Name"},
		{Key: "email", Label: "Email"},
		{Key: "created_at"},
	}

	result := NewColumnDataList(cols)

	if len(result) != 3 {
		t.Fatalf("len(result) = %d, want 3", len(result))
	}

	if result[2].Label != "Created At" {
		t.Errorf("result[2].Label = %q, want %q", result[2].Label, "Created At")
	}
}

// TestNewFormData tests FormData creation.
func TestNewFormData(t *testing.T) {
	tests := []struct {
		name       string
		input      types.ScaffoldFormInput
		wantMethod string
		wantCreate bool
		wantEdit   bool
	}{
		{
			name: "create form",
			input: types.ScaffoldFormInput{
				FormName:       "ProductForm",
				Domain:         "product",
				Action:         "create",
				SubmitEndpoint: "/products",
				Fields: []types.FieldDef{
					{Name: "Name", Type: "string"},
				},
			},
			wantMethod: "POST",
			wantCreate: true,
			wantEdit:   false,
		},
		{
			name: "edit form",
			input: types.ScaffoldFormInput{
				FormName:       "ProductForm",
				Domain:         "product",
				Action:         "edit",
				SubmitEndpoint: "/products/1",
				Fields: []types.FieldDef{
					{Name: "Name", Type: "string"},
				},
			},
			wantMethod: "PUT",
			wantCreate: false,
			wantEdit:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewFormData(tt.input, "github.com/example/app")

			if data.FormName != tt.input.FormName {
				t.Errorf("FormName = %q, want %q", data.FormName, tt.input.FormName)
			}
			if data.Method != tt.wantMethod {
				t.Errorf("Method = %q, want %q", data.Method, tt.wantMethod)
			}
			if data.IsCreate != tt.wantCreate {
				t.Errorf("IsCreate = %v, want %v", data.IsCreate, tt.wantCreate)
			}
			if data.IsEdit != tt.wantEdit {
				t.Errorf("IsEdit = %v, want %v", data.IsEdit, tt.wantEdit)
			}
			if data.ModelName != "Product" {
				t.Errorf("ModelName = %q, want %q", data.ModelName, "Product")
			}
		})
	}
}

// TestNewRowActionData tests RowActionData creation.
func TestNewRowActionData(t *testing.T) {
	tests := []struct {
		name       string
		input      types.RowActionDef
		wantView   bool
		wantEdit   bool
		wantDelete bool
	}{
		{
			name: "view action",
			input: types.RowActionDef{
				Type:  "view",
				Label: "View",
				URL:   "/products/{id}",
			},
			wantView:   true,
			wantEdit:   false,
			wantDelete: false,
		},
		{
			name: "edit action",
			input: types.RowActionDef{
				Type:  "edit",
				Label: "Edit",
				URL:   "/products/{id}/edit",
			},
			wantView:   false,
			wantEdit:   true,
			wantDelete: false,
		},
		{
			name: "delete action with confirm",
			input: types.RowActionDef{
				Type:           "delete",
				Label:          "Delete",
				URL:            "/products/{id}",
				Confirm:        true,
				ConfirmMessage: "Are you sure?",
			},
			wantView:   false,
			wantEdit:   false,
			wantDelete: true,
		},
		{
			name: "custom action",
			input: types.RowActionDef{
				Type:  "custom",
				Label: "Archive",
				URL:   "/products/{id}/archive",
			},
			wantView:   false,
			wantEdit:   false,
			wantDelete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewRowActionData(tt.input)

			if data.Type != tt.input.Type {
				t.Errorf("Type = %q, want %q", data.Type, tt.input.Type)
			}
			if data.Label != tt.input.Label {
				t.Errorf("Label = %q, want %q", data.Label, tt.input.Label)
			}
			if data.URL != tt.input.URL {
				t.Errorf("URL = %q, want %q", data.URL, tt.input.URL)
			}
			if data.Confirm != tt.input.Confirm {
				t.Errorf("Confirm = %v, want %v", data.Confirm, tt.input.Confirm)
			}
			if data.ConfirmMessage != tt.input.ConfirmMessage {
				t.Errorf("ConfirmMessage = %q, want %q", data.ConfirmMessage, tt.input.ConfirmMessage)
			}
			if data.IsView != tt.wantView {
				t.Errorf("IsView = %v, want %v", data.IsView, tt.wantView)
			}
			if data.IsEdit != tt.wantEdit {
				t.Errorf("IsEdit = %v, want %v", data.IsEdit, tt.wantEdit)
			}
			if data.IsDelete != tt.wantDelete {
				t.Errorf("IsDelete = %v, want %v", data.IsDelete, tt.wantDelete)
			}
		})
	}
}

// TestNewModalData tests ModalData creation.
func TestNewModalData(t *testing.T) {
	tests := []struct {
		name        string
		input       types.ScaffoldModalInput
		wantDialog  bool
		wantSheet   bool
		wantConfirm bool
	}{
		{
			name: "dialog modal",
			input: types.ScaffoldModalInput{
				ModalName:   "ProductModal",
				ModalType:   "dialog",
				Title:       "Add Product",
				ContentType: "form",
				TriggerConfig: types.TriggerConfig{
					ButtonText:    "Add",
					ButtonVariant: "primary",
				},
			},
			wantDialog:  true,
			wantSheet:   false,
			wantConfirm: false,
		},
		{
			name: "sheet modal",
			input: types.ScaffoldModalInput{
				ModalName:   "ProductDetails",
				ModalType:   "sheet",
				Title:       "Product Details",
				ContentType: "info",
			},
			wantDialog:  false,
			wantSheet:   true,
			wantConfirm: false,
		},
		{
			name: "confirm modal",
			input: types.ScaffoldModalInput{
				ModalName:   "DeleteConfirm",
				ModalType:   "confirm",
				Title:       "Confirm Delete",
				ContentType: "confirm",
			},
			wantDialog:  false,
			wantSheet:   false,
			wantConfirm: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewModalData("github.com/test/myapp", tt.input)

			if data.ModalName != tt.input.ModalName {
				t.Errorf("ModalName = %q, want %q", data.ModalName, tt.input.ModalName)
			}
			if data.ModalType != tt.input.ModalType {
				t.Errorf("ModalType = %q, want %q", data.ModalType, tt.input.ModalType)
			}
			if data.Title != tt.input.Title {
				t.Errorf("Title = %q, want %q", data.Title, tt.input.Title)
			}
			if data.IsDialog != tt.wantDialog {
				t.Errorf("IsDialog = %v, want %v", data.IsDialog, tt.wantDialog)
			}
			if data.IsSheet != tt.wantSheet {
				t.Errorf("IsSheet = %v, want %v", data.IsSheet, tt.wantSheet)
			}
			if data.IsConfirm != tt.wantConfirm {
				t.Errorf("IsConfirm = %v, want %v", data.IsConfirm, tt.wantConfirm)
			}
		})
	}
}

// TestNewPropData tests PropData creation.
func TestNewPropData(t *testing.T) {
	input := types.PropDef{
		Name:     "title",
		Type:     "string",
		Default:  "Untitled",
		Required: true,
	}

	data := NewPropData(input)

	if data.Name != input.Name {
		t.Errorf("Name = %q, want %q", data.Name, input.Name)
	}
	if data.Type != input.Type {
		t.Errorf("Type = %q, want %q", data.Type, input.Type)
	}
	if data.Default != input.Default {
		t.Errorf("Default = %q, want %q", data.Default, input.Default)
	}
	if data.Required != input.Required {
		t.Errorf("Required = %v, want %v", data.Required, input.Required)
	}
}

// TestNewAuthData tests AuthData creation.
func TestNewAuthData(t *testing.T) {
	data := NewAuthData("github.com/user/app", "myapp")

	if data.ModulePath != "github.com/user/app" {
		t.Errorf("ModulePath = %q, want %q", data.ModulePath, "github.com/user/app")
	}
	if data.ProjectName != "myapp" {
		t.Errorf("ProjectName = %q, want %q", data.ProjectName, "myapp")
	}
	if data.SessionType != "cookie" {
		t.Errorf("SessionType = %q, want %q (default)", data.SessionType, "cookie")
	}
}

// TestNewRelationshipData tests RelationshipData creation for all relationship types.
func TestNewRelationshipData(t *testing.T) {
	tests := []struct {
		name             string
		input            types.RelationshipDef
		domainName       string
		wantFieldName    string
		wantForeignKey   string
		wantReferences   string
		wantOnDelete     string
		wantIsBelongsTo  bool
		wantIsHasOne     bool
		wantIsHasMany    bool
		wantIsManyToMany bool
		wantHasFKField   bool
	}{
		{
			name: "belongs_to with defaults",
			input: types.RelationshipDef{
				Type:  "belongs_to",
				Model: "User",
			},
			domainName:      "order",
			wantFieldName:   "User",
			wantForeignKey:  "UserID",
			wantReferences:  "ID",
			wantOnDelete:    "CASCADE",
			wantIsBelongsTo: true,
			wantHasFKField:  true,
		},
		{
			name: "belongs_to with custom foreign key",
			input: types.RelationshipDef{
				Type:       "belongs_to",
				Model:      "User",
				ForeignKey: "AuthorID",
			},
			domainName:      "post",
			wantFieldName:   "User",
			wantForeignKey:  "AuthorID",
			wantReferences:  "ID",
			wantOnDelete:    "CASCADE",
			wantIsBelongsTo: true,
			wantHasFKField:  true,
		},
		{
			name: "has_one with defaults",
			input: types.RelationshipDef{
				Type:  "has_one",
				Model: "Profile",
			},
			domainName:     "user",
			wantFieldName:  "Profile",
			wantForeignKey: "UserID",
			wantReferences: "ID",
			wantOnDelete:   "CASCADE",
			wantIsHasOne:   true,
			wantHasFKField: false,
		},
		{
			name: "has_many with defaults",
			input: types.RelationshipDef{
				Type:  "has_many",
				Model: "Order",
			},
			domainName:     "user",
			wantFieldName:  "Orders",
			wantForeignKey: "UserID",
			wantReferences: "ID",
			wantOnDelete:   "CASCADE",
			wantIsHasMany:  true,
			wantHasFKField: false,
		},
		{
			name: "has_many with SET NULL on delete",
			input: types.RelationshipDef{
				Type:     "has_many",
				Model:    "Comment",
				OnDelete: "SET NULL",
			},
			domainName:     "post",
			wantFieldName:  "Comments",
			wantForeignKey: "PostID",
			wantReferences: "ID",
			wantOnDelete:   "SET NULL",
			wantIsHasMany:  true,
			wantHasFKField: false,
		},
		{
			name: "many_to_many with join table",
			input: types.RelationshipDef{
				Type:      "many_to_many",
				Model:     "Tag",
				JoinTable: "post_tags",
			},
			domainName:       "post",
			wantFieldName:    "Tags",
			wantForeignKey:   "",
			wantReferences:   "ID",
			wantOnDelete:     "CASCADE",
			wantIsManyToMany: true,
			wantHasFKField:   false,
		},
		{
			name: "many_to_many without join table",
			input: types.RelationshipDef{
				Type:  "many_to_many",
				Model: "Category",
			},
			domainName:       "product",
			wantFieldName:    "Categories",
			wantForeignKey:   "",
			wantReferences:   "ID",
			wantOnDelete:     "CASCADE",
			wantIsManyToMany: true,
			wantHasFKField:   false,
		},
		{
			name: "belongs_to with preload",
			input: types.RelationshipDef{
				Type:    "belongs_to",
				Model:   "Category",
				Preload: true,
			},
			domainName:      "product",
			wantFieldName:   "Category",
			wantForeignKey:  "CategoryID",
			wantReferences:  "ID",
			wantOnDelete:    "CASCADE",
			wantIsBelongsTo: true,
			wantHasFKField:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewRelationshipData(tt.input, tt.domainName)

			if data.Type != tt.input.Type {
				t.Errorf("Type = %q, want %q", data.Type, tt.input.Type)
			}
			if data.Model != tt.input.Model {
				t.Errorf("Model = %q, want %q", data.Model, tt.input.Model)
			}
			if data.FieldName != tt.wantFieldName {
				t.Errorf("FieldName = %q, want %q", data.FieldName, tt.wantFieldName)
			}
			if data.ForeignKey != tt.wantForeignKey {
				t.Errorf("ForeignKey = %q, want %q", data.ForeignKey, tt.wantForeignKey)
			}
			if data.References != tt.wantReferences {
				t.Errorf("References = %q, want %q", data.References, tt.wantReferences)
			}
			if data.OnDelete != tt.wantOnDelete {
				t.Errorf("OnDelete = %q, want %q", data.OnDelete, tt.wantOnDelete)
			}
			if data.IsBelongsTo != tt.wantIsBelongsTo {
				t.Errorf("IsBelongsTo = %v, want %v", data.IsBelongsTo, tt.wantIsBelongsTo)
			}
			if data.IsHasOne != tt.wantIsHasOne {
				t.Errorf("IsHasOne = %v, want %v", data.IsHasOne, tt.wantIsHasOne)
			}
			if data.IsHasMany != tt.wantIsHasMany {
				t.Errorf("IsHasMany = %v, want %v", data.IsHasMany, tt.wantIsHasMany)
			}
			if data.IsManyToMany != tt.wantIsManyToMany {
				t.Errorf("IsManyToMany = %v, want %v", data.IsManyToMany, tt.wantIsManyToMany)
			}
			if data.Preload != tt.input.Preload {
				t.Errorf("Preload = %v, want %v", data.Preload, tt.input.Preload)
			}

			// Check FK field
			hasFKField := data.ForeignKeyField != nil
			if hasFKField != tt.wantHasFKField {
				t.Errorf("has ForeignKeyField = %v, want %v", hasFKField, tt.wantHasFKField)
			}

			// If belongs_to, verify FK field details
			if tt.wantIsBelongsTo && data.ForeignKeyField != nil {
				if data.ForeignKeyField.Name != tt.wantForeignKey {
					t.Errorf("ForeignKeyField.Name = %q, want %q", data.ForeignKeyField.Name, tt.wantForeignKey)
				}
				if data.ForeignKeyField.Type != "uint" {
					t.Errorf("ForeignKeyField.Type = %q, want %q", data.ForeignKeyField.Type, "uint")
				}
			}
		})
	}
}

// TestNewRelationshipDataList tests creating list of RelationshipData.
func TestNewRelationshipDataList(t *testing.T) {
	rels := []types.RelationshipDef{
		{Type: "belongs_to", Model: "User"},
		{Type: "has_many", Model: "OrderItem"},
		{Type: "many_to_many", Model: "Tag", JoinTable: "order_tags"},
	}

	result := NewRelationshipDataList(rels, "order")

	if len(result) != 3 {
		t.Fatalf("len(result) = %d, want 3", len(result))
	}

	if result[0].FieldName != "User" {
		t.Errorf("result[0].FieldName = %q, want %q", result[0].FieldName, "User")
	}
	if result[1].FieldName != "OrderItems" {
		t.Errorf("result[1].FieldName = %q, want %q", result[1].FieldName, "OrderItems")
	}
	if result[2].FieldName != "Tags" {
		t.Errorf("result[2].FieldName = %q, want %q", result[2].FieldName, "Tags")
	}
}

// TestNewRelationshipDataList_Empty tests empty relationship list.
func TestNewRelationshipDataList_Empty(t *testing.T) {
	result := NewRelationshipDataList(nil, "order")
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}

	result = NewRelationshipDataList([]types.RelationshipDef{}, "order")
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0 for empty slice", len(result))
	}
}

// TestBuildGORMTag tests GORM tag generation for relationships.
func TestBuildGORMTag(t *testing.T) {
	tests := []struct {
		name       string
		relType    string
		foreignKey string
		references string
		joinTable  string
		onDelete   string
		want       string
	}{
		{
			name:       "belongs_to",
			relType:    "belongs_to",
			foreignKey: "UserID",
			references: "ID",
			want:       "foreignKey:UserID;references:ID",
		},
		{
			name:       "has_one",
			relType:    "has_one",
			foreignKey: "UserID",
			references: "ID",
			onDelete:   "CASCADE",
			want:       "foreignKey:UserID;references:ID",
		},
		{
			name:       "has_many with SET NULL",
			relType:    "has_many",
			foreignKey: "PostID",
			references: "ID",
			onDelete:   "SET NULL",
			want:       "foreignKey:PostID;references:ID;constraint:OnDelete:SET NULL",
		},
		{
			name:      "many_to_many with join table",
			relType:   "many_to_many",
			joinTable: "post_tags",
			want:      "many2many:post_tags",
		},
		{
			name:    "many_to_many without join table",
			relType: "many_to_many",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildGORMTag(tt.relType, tt.foreignKey, tt.references, tt.joinTable, tt.onDelete)
			if result != tt.want {
				t.Errorf("buildGORMTag() = %q, want %q", result, tt.want)
			}
		})
	}
}

// TestNewDomainData_WithRelationships tests DomainData creation with relationships.
func TestNewDomainData_WithRelationships(t *testing.T) {
	input := types.ScaffoldDomainInput{
		DomainName: "order",
		Fields: []types.FieldDef{
			{Name: "Total", Type: "float64"},
			{Name: "Status", Type: "string"},
		},
		Relationships: []types.RelationshipDef{
			{Type: "belongs_to", Model: "User", Preload: true},
			{Type: "has_many", Model: "OrderItem"},
			{Type: "many_to_many", Model: "Tag", JoinTable: "order_tags", Preload: true},
		},
	}

	data := NewDomainData(input, "github.com/user/app")

	// Check relationships
	if !data.HasRelationships {
		t.Error("HasRelationships should be true")
	}
	if len(data.Relationships) != 3 {
		t.Fatalf("len(Relationships) = %d, want 3", len(data.Relationships))
	}

	// Check preload relationships
	if len(data.PreloadRelationships) != 2 {
		t.Fatalf("len(PreloadRelationships) = %d, want 2", len(data.PreloadRelationships))
	}

	// Verify the preloaded ones are User and Tags
	preloadNames := make(map[string]bool)
	for _, rel := range data.PreloadRelationships {
		preloadNames[rel.FieldName] = true
	}
	if !preloadNames["User"] {
		t.Error("User should be in PreloadRelationships")
	}
	if !preloadNames["Tags"] {
		t.Error("Tags should be in PreloadRelationships")
	}
}

// TestNewDomainData_WithoutRelationships tests DomainData with no relationships.
func TestNewDomainData_WithoutRelationships(t *testing.T) {
	input := types.ScaffoldDomainInput{
		DomainName: "product",
		Fields: []types.FieldDef{
			{Name: "Name", Type: "string"},
		},
	}

	data := NewDomainData(input, "github.com/user/app")

	if data.HasRelationships {
		t.Error("HasRelationships should be false")
	}
	if len(data.Relationships) != 0 {
		t.Errorf("len(Relationships) = %d, want 0", len(data.Relationships))
	}
	if len(data.PreloadRelationships) != 0 {
		t.Errorf("len(PreloadRelationships) = %d, want 0", len(data.PreloadRelationships))
	}
}
