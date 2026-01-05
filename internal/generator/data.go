package generator

import (
	"strings"

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
	// WithUserManagement enables admin user management.
	WithUserManagement bool
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
	// Handle pointer types by stripping the * prefix
	baseType := strings.TrimPrefix(goType, "*")

	switch baseType {
	case "string":
		return "input"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "float32", "float64":
		return "number"
	case "bool":
		return "checkbox"
	case "time.Time":
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
	// URLPathSegment is the URL path without leading slash (e.g., "products").
	URLPathSegment string
	// Fields is the list of fields.
	Fields []FieldData
	// Relationships is the list of model relationships.
	Relationships []RelationshipData
	// HasRelationships is true if any relationships are defined.
	HasRelationships bool
	// PreloadRelationships is the list of relationships to preload.
	PreloadRelationships []RelationshipData
	// WithSoftDelete enables soft delete.
	WithSoftDelete bool
	// WithCrudViews generates CRUD views.
	WithCrudViews bool
	// WithPagination enables pagination in views.
	WithPagination bool
	// WithSearch enables search in views.
	WithSearch bool
	// Layout specifies the view layout: dashboard, base, auth, none. Defaults to "dashboard".
	Layout string
	// RouteGroup specifies the middleware context: public, authenticated, admin. Defaults to "public".
	RouteGroup string
}

// NewDomainData creates DomainData from ScaffoldDomainInput and module path.
func NewDomainData(input types.ScaffoldDomainInput, modulePath string) DomainData {
	relationships := NewRelationshipDataList(input.Relationships, input.DomainName)

	// Filter relationships that should be preloaded
	var preloadRels []RelationshipData
	for _, rel := range relationships {
		if rel.Preload {
			preloadRels = append(preloadRels, rel)
		}
	}

	withCrudViews := input.GetWithCrudViews()

	// Default layout to "dashboard"
	layout := input.Layout
	if layout == "" {
		layout = "dashboard"
	}

	// Default route group to "public"
	routeGroup := input.RouteGroup
	if routeGroup == "" {
		routeGroup = "public"
	}

	urlPath := utils.ToURLPath(input.DomainName)
	return DomainData{
		ModulePath:           modulePath,
		DomainName:           input.DomainName,
		ModelName:            utils.ToModelName(input.DomainName),
		PackageName:          utils.ToPackageName(input.DomainName),
		VariableName:         utils.ToVariableName(input.DomainName),
		TableName:            utils.ToTableName(input.DomainName),
		URLPath:              urlPath,
		URLPathSegment:       strings.TrimPrefix(urlPath, "/"),
		Fields:               NewFieldDataList(input.Fields),
		Relationships:        relationships,
		HasRelationships:     len(relationships) > 0,
		PreloadRelationships: preloadRels,
		WithSoftDelete:       input.GetWithSoftDelete(),
		WithCrudViews:        withCrudViews,
		WithPagination:       withCrudViews, // Enable pagination when CRUD views are generated
		WithSearch:           withCrudViews, // Enable search when CRUD views are generated
		Layout:               layout,
		RouteGroup:           routeGroup,
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
	// URLPathSegment is the URL path without leading slash.
	URLPathSegment string
	// ViewType is list, show, form, card, table, or custom.
	ViewType string
	// ViewName is the view file name.
	ViewName string
	// Fields is the list of fields.
	Fields []FieldData
	// Columns is the list of columns (for table views).
	Columns []ColumnData
	// Relationships is the list of model relationships.
	Relationships []RelationshipData
	// WithPagination enables pagination.
	WithPagination bool
	// WithSearch enables search.
	WithSearch bool
	// WithFilters enables filters.
	WithFilters bool
	// WithSorting enables sorting (for table views).
	WithSorting bool
	// WithBulkActions enables bulk actions (for table views).
	WithBulkActions bool
	// WithSoftDelete indicates if soft delete is enabled.
	WithSoftDelete bool
	// RowActions is the list of row actions (for table views).
	RowActions []RowActionData
	// EmptyStateMessage is shown when empty.
	EmptyStateMessage string
	// SubmitURL is the form submission URL.
	SubmitURL string
	// Method is POST or PUT.
	Method string
	// SuccessRedirect is the redirect target.
	SuccessRedirect string
	// Layout specifies the view layout: dashboard, base, auth, none. Defaults to "dashboard".
	Layout string
}

// FormData is the template data for form scaffolding.
type FormData struct {
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
	// URLPathSegment is the URL path without leading slash.
	URLPathSegment string
	// FormName is the form component name.
	FormName string
	// Action is create or edit.
	Action string
	// Fields is the list of form fields.
	Fields []FieldData
	// Relationships is the list of model relationships.
	Relationships []RelationshipData
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
	urlPath := utils.ToURLPath(input.Domain)
	return FormData{
		ModulePath:     modulePath,
		DomainName:     input.Domain,
		ModelName:      utils.ToModelName(input.Domain),
		PackageName:    utils.ToPackageName(input.Domain),
		VariableName:   utils.ToVariableName(input.Domain),
		URLPath:        urlPath,
		URLPathSegment: strings.TrimPrefix(urlPath, "/"),
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
	// PackageName is the package name.
	PackageName string
	// VariableName is the variable name.
	VariableName string
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
	// ModulePath is the Go module path.
	ModulePath string
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
func NewModalData(modulePath string, input types.ScaffoldModalInput) ModalData {
	return ModalData{
		ModulePath:     modulePath,
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
	// ModulePath is the Go module path.
	ModulePath string
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
	// ModelName is an alias for PageName (for template compatibility).
	ModelName string
	// PackageName is the package name.
	PackageName string
	// VariableName is the variable name.
	VariableName string
	// URLPath is the URL path (alias for Route).
	URLPath string
	// URLPathSegment is the URL path without leading slash.
	URLPathSegment string
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
	// Fields is an empty list for template compatibility.
	Fields []FieldData
	// WithPagination for template compatibility.
	WithPagination bool
	// WithSearch for template compatibility.
	WithSearch bool
	// EmptyStateMessage for template compatibility.
	EmptyStateMessage string
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
	// PageName is an alias for Name.
	PageName string
	// Locale is the locale code.
	Locale string
	// Content is the configuration content.
	Content map[string]interface{}
	// Title is the page title.
	Title string
	// Description is the page description.
	Description string
	// Heading is the page heading.
	Heading string
	// Layout is the layout type.
	Layout string
	// EmptyStateMessage is shown when empty.
	EmptyStateMessage string
	// WithBreadcrumbs enables breadcrumbs.
	WithBreadcrumbs bool
	// WithTable enables table config.
	WithTable bool
	// WithFilters enables filters config.
	WithFilters bool
	// WithActions enables actions config.
	WithActions bool
	// Sidebar enables sidebar.
	Sidebar bool
}

// SeedRelationshipData represents a relationship for seeding.
type SeedRelationshipData struct {
	// Field is the foreign key field name (e.g., "UserID").
	Field string
	// Model is the related model name (e.g., "User").
	Model string
	// ModelVar is the variable name (e.g., "user").
	ModelVar string
	// Strategy is how to assign: random, each, distribute.
	Strategy string
}

// SeedDistributionData represents value distribution for seeding.
type SeedDistributionData struct {
	// Field is the field name.
	Field string
	// Values is the list of value distributions.
	Values []SeedValueData
	// TotalCount is the sum of all value counts.
	TotalCount int
}

// SeedValueData represents a value and its count.
type SeedValueData struct {
	// Value is the field value.
	Value string
	// Count is how many records should have this value.
	Count int
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
	// Relationships is the list of relationships to seed.
	Relationships []SeedRelationshipData
	// Distributions is the list of value distributions.
	Distributions []SeedDistributionData
	// HasRelationships is true if there are relationships.
	HasRelationships bool
	// HasDistributions is true if there are distributions.
	HasDistributions bool
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

// RelationshipData is the template data for a model relationship.
type RelationshipData struct {
	// Type is the relationship type: belongs_to, has_one, has_many, many_to_many.
	Type string
	// Model is the related model name in PascalCase (e.g., "User").
	Model string
	// FieldName is the struct field name (e.g., "User" for belongs_to, "Orders" for has_many).
	FieldName string
	// ForeignKey is the foreign key field name.
	ForeignKey string
	// References is the referenced field (usually "ID").
	References string
	// JoinTable is the join table name (for many_to_many).
	JoinTable string
	// OnDelete is the delete behavior.
	OnDelete string
	// Preload indicates if the relationship should be preloaded by default.
	Preload bool
	// IsBelongsTo is true for belongs_to relationships.
	IsBelongsTo bool
	// IsHasOne is true for has_one relationships.
	IsHasOne bool
	// IsHasMany is true for has_many relationships.
	IsHasMany bool
	// IsManyToMany is true for many_to_many relationships.
	IsManyToMany bool
	// GORMTag is the complete GORM struct tag for the relationship.
	GORMTag string
	// ForeignKeyField is the FK field definition (for belongs_to).
	ForeignKeyField *FieldData
}

// NewRelationshipData creates RelationshipData from a RelationshipDef.
func NewRelationshipData(rel types.RelationshipDef, domainName string) RelationshipData {
	fieldName := rel.Model
	foreignKey := rel.ForeignKey
	references := rel.References
	onDelete := rel.OnDelete

	// Defaults
	if references == "" {
		references = "ID"
	}
	if onDelete == "" {
		onDelete = "CASCADE"
	}

	// Determine field name based on relationship type
	switch rel.Type {
	case "belongs_to":
		// belongs_to: field name is singular (e.g., "User")
		if foreignKey == "" {
			foreignKey = rel.Model + "ID"
		}
	case "has_one":
		// has_one: field name is singular (e.g., "Profile")
		if foreignKey == "" {
			foreignKey = utils.ToModelName(domainName) + "ID"
		}
	case "has_many":
		// has_many: field name is plural (e.g., "Orders")
		fieldName = utils.Pluralize(rel.Model)
		if foreignKey == "" {
			foreignKey = utils.ToModelName(domainName) + "ID"
		}
	case "many_to_many":
		// many_to_many: field name is plural (e.g., "Tags")
		fieldName = utils.Pluralize(rel.Model)
	}

	// Build GORM tag
	gormTag := buildGORMTag(rel.Type, foreignKey, references, rel.JoinTable, onDelete)

	// Create FK field for belongs_to relationships
	var fkField *FieldData
	if rel.Type == "belongs_to" {
		fkField = &FieldData{
			Name:      foreignKey,
			Type:      "uint",
			GORMTags:  "",
			JSONName:  utils.ToJSONTag(foreignKey),
			Omitempty: true,
		}
	}

	return RelationshipData{
		Type:            rel.Type,
		Model:           rel.Model,
		FieldName:       fieldName,
		ForeignKey:      foreignKey,
		References:      references,
		JoinTable:       rel.JoinTable,
		OnDelete:        onDelete,
		Preload:         rel.Preload,
		IsBelongsTo:     rel.Type == "belongs_to",
		IsHasOne:        rel.Type == "has_one",
		IsHasMany:       rel.Type == "has_many",
		IsManyToMany:    rel.Type == "many_to_many",
		GORMTag:         gormTag,
		ForeignKeyField: fkField,
	}
}

// buildGORMTag builds the GORM struct tag for a relationship.
func buildGORMTag(relType, foreignKey, references, joinTable, onDelete string) string {
	var parts []string

	switch relType {
	case "belongs_to":
		parts = append(parts, "foreignKey:"+foreignKey)
		parts = append(parts, "references:"+references)
	case "has_one", "has_many":
		parts = append(parts, "foreignKey:"+foreignKey)
		parts = append(parts, "references:"+references)
		if onDelete != "" && onDelete != "CASCADE" {
			parts = append(parts, "constraint:OnDelete:"+onDelete)
		}
	case "many_to_many":
		if joinTable != "" {
			parts = append(parts, "many2many:"+joinTable)
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ";")
}

// NewRelationshipDataList creates a list of RelationshipData from RelationshipDefs.
func NewRelationshipDataList(rels []types.RelationshipDef, domainName string) []RelationshipData {
	result := make([]RelationshipData, len(rels))
	for i, rel := range rels {
		result[i] = NewRelationshipData(rel, domainName)
	}
	return result
}
