package generator

import (
	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/dbb1dev/go-mcp/internal/utils"
)

// ProjectData is the template data for project scaffolding.
type ProjectData struct {
	// ProjectName is the project directory name.
	ProjectName string
	// ModulePath is the Go module path.
	ModulePath string
	// DatabaseType is sqlite, postgres, or mysql.
	DatabaseType string
	// WithAuth enables authentication scaffolding.
	WithAuth bool
}

// NewProjectData creates ProjectData from ScaffoldProjectInput.
func NewProjectData(input types.ScaffoldProjectInput) ProjectData {
	dbType := input.DatabaseType
	if dbType == "" {
		dbType = "sqlite"
	}
	return ProjectData{
		ProjectName:  input.ProjectName,
		ModulePath:   input.ModulePath,
		DatabaseType: dbType,
		WithAuth:     input.WithAuth,
	}
}

// FieldData is the template data for a model field.
type FieldData struct {
	// Name is the field name in PascalCase.
	Name string
	// Type is the Go type.
	Type string
	// GORMTags are the GORM struct tags.
	GORMTags string
	// JSONName is the JSON field name.
	JSONName string
	// Omitempty adds omitempty to JSON tag.
	Omitempty bool
	// FormType is the form input type.
	FormType string
	// Required indicates if the field is required.
	Required bool
	// Label is the display label.
	Label string
}

// NewFieldData creates FieldData from a FieldDef.
func NewFieldData(field types.FieldDef) FieldData {
	jsonTag := field.JSONTag
	if jsonTag == "" {
		jsonTag = utils.ToJSONTag(field.Name)
	}

	label := field.Label
	if label == "" {
		label = utils.ToLabel(field.Name)
	}

	formType := field.FormType
	if formType == "" {
		formType = inferFormType(field.Type)
	}

	return FieldData{
		Name:      field.Name,
		Type:      field.Type,
		GORMTags:  field.GORMTags,
		JSONName:  jsonTag,
		Omitempty: !field.Required,
		FormType:  formType,
		Required:  field.Required,
		Label:     label,
	}
}

// NewFieldDataList creates a list of FieldData from FieldDefs.
func NewFieldDataList(fields []types.FieldDef) []FieldData {
	result := make([]FieldData, len(fields))
	for i, field := range fields {
		result[i] = NewFieldData(field)
	}
	return result
}

// inferFormType infers the form type from a Go type.
func inferFormType(goType string) string {
	switch goType {
	case "string":
		return "input"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "float32", "float64":
		return "number"
	case "bool":
		return "checkbox"
	case "time.Time", "*time.Time":
		return "datetime"
	default:
		return "input"
	}
}

// DomainData is the template data for domain scaffolding.
type DomainData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// DomainName is the domain name (e.g., "product").
	DomainName string
	// ModelName is the model struct name (e.g., "Product").
	ModelName string
	// PackageName is the package name (e.g., "product").
	PackageName string
	// VariableName is the variable name (e.g., "product").
	VariableName string
	// TableName is the database table name (e.g., "products").
	TableName string
	// URLPath is the URL path (e.g., "/products").
	URLPath string
	// Fields is the list of fields.
	Fields []FieldData
	// WithSoftDelete enables soft delete.
	WithSoftDelete bool
	// WithCrudViews generates CRUD views.
	WithCrudViews bool
}

// NewDomainData creates DomainData from ScaffoldDomainInput and module path.
func NewDomainData(input types.ScaffoldDomainInput, modulePath string) DomainData {
	return DomainData{
		ModulePath:     modulePath,
		DomainName:     input.DomainName,
		ModelName:      utils.ToModelName(input.DomainName),
		PackageName:    utils.ToPackageName(input.DomainName),
		VariableName:   utils.ToVariableName(input.DomainName),
		TableName:      utils.ToTableName(input.DomainName),
		URLPath:        utils.ToURLPath(input.DomainName),
		Fields:         NewFieldDataList(input.Fields),
		WithSoftDelete: input.GetWithSoftDelete(),
		WithCrudViews:  input.GetWithCrudViews(),
	}
}

// ColumnData is the template data for a table column.
type ColumnData struct {
	// Key is the field key.
	Key string
	// Label is the display label.
	Label string
	// Sortable enables sorting.
	Sortable bool
	// Format is the display format.
	Format string
	// Width is the CSS width.
	Width string
}

// NewColumnData creates ColumnData from a ColumnDef.
func NewColumnData(col types.ColumnDef) ColumnData {
	label := col.Label
	if label == "" {
		label = utils.ToLabel(col.Key)
	}
	return ColumnData{
		Key:      col.Key,
		Label:    label,
		Sortable: col.Sortable,
		Format:   col.Format,
		Width:    col.Width,
	}
}

// NewColumnDataList creates a list of ColumnData from ColumnDefs.
func NewColumnDataList(cols []types.ColumnDef) []ColumnData {
	result := make([]ColumnData, len(cols))
	for i, col := range cols {
		result[i] = NewColumnData(col)
	}
	return result
}

// ViewData is the template data for view scaffolding.
type ViewData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// DomainName is the domain name.
	DomainName string
	// ModelName is the model struct name.
	ModelName string
	// PackageName is the package name.
	PackageName string
	// VariableName is the variable name.
	VariableName string
	// URLPath is the URL path.
	URLPath string
	// ViewType is list, show, form, card, table, or custom.
	ViewType string
	// ViewName is the view file name.
	ViewName string
	// Fields is the list of fields.
	Fields []FieldData
	// Columns is the list of columns (for table views).
	Columns []ColumnData
	// WithPagination enables pagination.
	WithPagination bool
	// WithSearch enables search.
	WithSearch bool
	// WithFilters enables filters.
	WithFilters bool
	// EmptyStateMessage is shown when empty.
	EmptyStateMessage string
	// SubmitURL is the form submission URL.
	SubmitURL string
	// Method is POST or PUT.
	Method string
	// SuccessRedirect is the redirect target.
	SuccessRedirect string
}

// FormData is the template data for form scaffolding.
type FormData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// DomainName is the domain name.
	DomainName string
	// ModelName is the model struct name.
	ModelName string
	// FormName is the form component name.
	FormName string
	// Action is create or edit.
	Action string
	// Fields is the list of form fields.
	Fields []FieldData
	// SubmitEndpoint is the submission URL.
	SubmitEndpoint string
	// Method is POST or PUT.
	Method string
	// IsCreate is true for create forms.
	IsCreate bool
	// IsEdit is true for edit forms.
	IsEdit bool
}

// NewFormData creates FormData from ScaffoldFormInput.
func NewFormData(input types.ScaffoldFormInput, modulePath string) FormData {
	method := "POST"
	if input.Action == "edit" {
		method = "PUT"
	}
	return FormData{
		ModulePath:     modulePath,
		DomainName:     input.Domain,
		ModelName:      utils.ToModelName(input.Domain),
		FormName:       input.FormName,
		Action:         input.Action,
		Fields:         NewFieldDataList(input.Fields),
		SubmitEndpoint: input.SubmitEndpoint,
		Method:         method,
		IsCreate:       input.Action == "create",
		IsEdit:         input.Action == "edit",
	}
}

// TableData is the template data for table scaffolding.
type TableData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// DomainName is the domain name.
	DomainName string
	// ModelName is the model struct name.
	ModelName string
	// TableName is the table component name.
	TableName string
	// URLPath is the URL path.
	URLPath string
	// Columns is the list of columns.
	Columns []ColumnData
	// WithPagination enables pagination.
	WithPagination bool
	// WithSorting enables sorting.
	WithSorting bool
	// WithSearch enables search.
	WithSearch bool
	// WithBulkActions enables bulk actions.
	WithBulkActions bool
	// RowActions is the list of row actions.
	RowActions []RowActionData
}

// RowActionData is the template data for a row action.
type RowActionData struct {
	// Type is view, edit, delete, or custom.
	Type string
	// Label is the button label.
	Label string
	// URL is the action URL.
	URL string
	// Confirm shows confirmation dialog.
	Confirm bool
	// ConfirmMessage is the confirmation message.
	ConfirmMessage string
	// IsView is true for view actions.
	IsView bool
	// IsEdit is true for edit actions.
	IsEdit bool
	// IsDelete is true for delete actions.
	IsDelete bool
}

// NewRowActionData creates RowActionData from a RowActionDef.
func NewRowActionData(action types.RowActionDef) RowActionData {
	return RowActionData{
		Type:           action.Type,
		Label:          action.Label,
		URL:            action.URL,
		Confirm:        action.Confirm,
		ConfirmMessage: action.ConfirmMessage,
		IsView:         action.Type == "view",
		IsEdit:         action.Type == "edit",
		IsDelete:       action.Type == "delete",
	}
}

// ModalData is the template data for modal scaffolding.
type ModalData struct {
	// ModalName is the modal component name.
	ModalName string
	// ModalType is dialog, sheet, or confirm.
	ModalType string
	// Title is the modal title.
	Title string
	// ContentType is form, info, or confirm.
	ContentType string
	// TriggerButton is the trigger button text.
	TriggerButton string
	// TriggerVariant is the button variant.
	TriggerVariant string
	// HTMXURL is the content load URL.
	HTMXURL string
	// IsDialog is true for dialog modals.
	IsDialog bool
	// IsSheet is true for sheet modals.
	IsSheet bool
	// IsConfirm is true for confirm modals.
	IsConfirm bool
}

// NewModalData creates ModalData from ScaffoldModalInput.
func NewModalData(input types.ScaffoldModalInput) ModalData {
	return ModalData{
		ModalName:      input.ModalName,
		ModalType:      input.ModalType,
		Title:          input.Title,
		ContentType:    input.ContentType,
		TriggerButton:  input.TriggerConfig.ButtonText,
		TriggerVariant: input.TriggerConfig.ButtonVariant,
		HTMXURL:        input.TriggerConfig.HTMXURL,
		IsDialog:       input.ModalType == "dialog",
		IsSheet:        input.ModalType == "sheet",
		IsConfirm:      input.ModalType == "confirm",
	}
}

// ComponentData is the template data for component scaffolding.
type ComponentData struct {
	// ComponentName is the component name.
	ComponentName string
	// ComponentType is the component type.
	ComponentType string
	// Props is the list of properties.
	Props []PropData
	// WithHTMX includes HTMX attributes.
	WithHTMX bool
	// AlpineState contains Alpine.js state.
	AlpineState map[string]interface{}
}

// PropData is the template data for a component prop.
type PropData struct {
	// Name is the prop name.
	Name string
	// Type is the Go type.
	Type string
	// Default is the default value.
	Default string
	// Required is true if required.
	Required bool
}

// NewPropData creates PropData from a PropDef.
func NewPropData(prop types.PropDef) PropData {
	return PropData{
		Name:     prop.Name,
		Type:     prop.Type,
		Default:  prop.Default,
		Required: prop.Required,
	}
}

// PageData is the template data for page scaffolding.
type PageData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// PageName is the page name.
	PageName string
	// Route is the URL route.
	Route string
	// Layout is the layout type.
	Layout string
	// Sections is the list of sections.
	Sections []SectionData
	// Title is the page title.
	Title string
	// Description is the page description.
	Description string
}

// SectionData is the template data for a page section.
type SectionData struct {
	// Type is the section type.
	Type string
	// Config contains section configuration.
	Config map[string]interface{}
}

// ConfigData is the template data for config scaffolding.
type ConfigData struct {
	// ConfigType is the config type.
	ConfigType string
	// Name is the config file name.
	Name string
	// Locale is the locale code.
	Locale string
	// Content is the configuration content.
	Content map[string]interface{}
}

// SeedData is the template data for seeder scaffolding.
type SeedData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// DomainName is the domain name.
	DomainName string
	// ModelName is the model struct name.
	ModelName string
	// TableName is the database table name.
	TableName string
	// Fields is the list of fields.
	Fields []FieldData
	// Count is the number of records.
	Count int
	// WithFaker uses faker for data.
	WithFaker bool
	// Dependencies is the list of dependencies.
	Dependencies []string
}

// AuthData is the template data for auth scaffolding.
type AuthData struct {
	// ModulePath is the Go module path.
	ModulePath string
	// ProjectName is the project name.
	ProjectName string
	// SessionType is cookie or jwt.
	SessionType string
}

// NewAuthData creates AuthData.
func NewAuthData(modulePath, projectName string) AuthData {
	return AuthData{
		ModulePath:  modulePath,
		ProjectName: projectName,
		SessionType: "cookie", // default to cookie-based sessions
	}
}
