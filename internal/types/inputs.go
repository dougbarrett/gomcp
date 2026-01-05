// Package types contains input and output types for MCP tools.
package types

// ScaffoldProjectInput is the input for the scaffold_project tool.
type ScaffoldProjectInput struct {
	// ProjectName is the name of the project (used for display and defaults).
	ProjectName string `json:"project_name"`
	// ModulePath is the Go module path (e.g., github.com/user/project).
	ModulePath string `json:"module_path"`
	// DatabaseType is the database driver: sqlite, postgres, or mysql.
	DatabaseType string `json:"database_type,omitempty"`
	// WithAuth enables authentication scaffolding.
	WithAuth bool `json:"with_auth,omitempty"`
	// InCurrentDir generates files in the current directory instead of a subdirectory.
	InCurrentDir bool `json:"in_current_dir,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// FieldDef defines a model field for scaffolding.
type FieldDef struct {
	// Name is the field name in PascalCase (e.g., "FirstName").
	Name string `json:"name"`
	// Type is the Go type (e.g., "string", "int", "time.Time").
	Type string `json:"type"`
	// GORMTags are optional GORM struct tags (e.g., "size:255;not null").
	GORMTags string `json:"gorm_tags,omitempty"`
	// JSONTag is the JSON field name (defaults to snake_case of Name).
	JSONTag string `json:"json_tag,omitempty"`
	// FormType is the form input type: input, textarea, select, checkbox, date, email, password, number, rating, tags, slider.
	FormType string `json:"form_type,omitempty"`
	// Required indicates if the field is required in forms.
	Required bool `json:"required,omitempty"`
	// Label is the display label for forms (defaults to Name with spaces).
	Label string `json:"label,omitempty"`
}

// RelationshipDef defines a model relationship.
type RelationshipDef struct {
	// Type is the relationship type: belongs_to, has_one, has_many, many_to_many.
	Type string `json:"type"`
	// Model is the related model name in PascalCase (e.g., "User", "OrderItem").
	Model string `json:"model"`
	// ForeignKey is the foreign key field name. Defaults to {Model}ID for belongs_to.
	ForeignKey string `json:"foreign_key,omitempty"`
	// References is the referenced field. Defaults to "ID".
	References string `json:"references,omitempty"`
	// JoinTable is the join table name (for many_to_many).
	JoinTable string `json:"join_table,omitempty"`
	// OnDelete is the delete behavior: CASCADE, SET NULL, RESTRICT. Defaults to CASCADE.
	OnDelete string `json:"on_delete,omitempty"`
	// Preload indicates if the relationship should be preloaded by default.
	Preload bool `json:"preload,omitempty"`
}

// ScaffoldDomainInput is the input for the scaffold_domain tool.
type ScaffoldDomainInput struct {
	// DomainName is the domain name in singular form (e.g., "product").
	DomainName string `json:"domain_name"`
	// Fields is the list of model fields.
	Fields []FieldDef `json:"fields"`
	// Relationships is the list of model relationships.
	Relationships []RelationshipDef `json:"relationships,omitempty"`
	// WithCrudViews generates CRUD templ views. Defaults to true.
	WithCrudViews *bool `json:"with_crud_views,omitempty"`
	// WithSoftDelete includes soft delete support. Defaults to true.
	WithSoftDelete *bool `json:"with_soft_delete,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// GetWithCrudViews returns the WithCrudViews value with default true.
func (s ScaffoldDomainInput) GetWithCrudViews() bool {
	if s.WithCrudViews == nil {
		return true
	}
	return *s.WithCrudViews
}

// GetWithSoftDelete returns the WithSoftDelete value with default true.
func (s ScaffoldDomainInput) GetWithSoftDelete() bool {
	if s.WithSoftDelete == nil {
		return true
	}
	return *s.WithSoftDelete
}

// MethodDef defines a service or repository method.
type MethodDef struct {
	// Name is the method name in PascalCase.
	Name string `json:"name"`
	// Description describes what the method does.
	Description string `json:"description,omitempty"`
	// Params is the list of parameter definitions.
	Params []ParamDef `json:"params,omitempty"`
	// Returns is the return type (e.g., "*models.Product", "error").
	Returns string `json:"returns,omitempty"`
}

// ParamDef defines a method parameter.
type ParamDef struct {
	// Name is the parameter name.
	Name string `json:"name"`
	// Type is the parameter type.
	Type string `json:"type"`
}

// ScaffoldRepositoryInput is the input for the scaffold_repository tool.
type ScaffoldRepositoryInput struct {
	// DomainName is the domain name (e.g., "product").
	DomainName string `json:"domain_name"`
	// ModelName is the model struct name (e.g., "Product"). Defaults to PascalCase of DomainName.
	ModelName string `json:"model_name,omitempty"`
	// Methods is the list of custom methods beyond CRUD.
	Methods []MethodDef `json:"methods,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ScaffoldServiceInput is the input for the scaffold_service tool.
type ScaffoldServiceInput struct {
	// DomainName is the domain name (e.g., "product").
	DomainName string `json:"domain_name"`
	// Methods is the list of service methods.
	Methods []MethodDef `json:"methods,omitempty"`
	// Dependencies is the list of other services/repos to inject.
	Dependencies []string `json:"dependencies,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ActionDef defines a controller action.
type ActionDef struct {
	// Name is the action name (e.g., "List", "Create").
	Name string `json:"name"`
	// Method is the HTTP method (GET, POST, PUT, DELETE).
	Method string `json:"method"`
	// Path is the URL path suffix (e.g., "/", "/{id}").
	Path string `json:"path"`
	// WithView generates a corresponding templ view.
	WithView bool `json:"with_view,omitempty"`
}

// ScaffoldControllerInput is the input for the scaffold_controller tool.
type ScaffoldControllerInput struct {
	// DomainName is the domain/feature name.
	DomainName string `json:"domain_name"`
	// Actions is the list of controller actions.
	Actions []ActionDef `json:"actions"`
	// BasePath is the URL base path. Defaults to /{domain_name}.
	BasePath string `json:"base_path,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ViewConfig contains view-specific configuration.
type ViewConfig struct {
	// ItemsVariable is the template variable name for items (for list views).
	ItemsVariable string `json:"items_variable,omitempty"`
	// ItemComponent is the component type: card, table_row, or custom.
	ItemComponent string `json:"item_component,omitempty"`
	// WithPagination enables pagination.
	WithPagination bool `json:"with_pagination,omitempty"`
	// WithSearch enables search functionality.
	WithSearch bool `json:"with_search,omitempty"`
	// WithFilters enables filter controls.
	WithFilters bool `json:"with_filters,omitempty"`
	// WithSorting enables column sorting (for table views).
	WithSorting bool `json:"with_sorting,omitempty"`
	// WithBulkActions enables bulk actions (for table views).
	WithBulkActions bool `json:"with_bulk_actions,omitempty"`
	// WithSoftDelete indicates if soft delete is enabled.
	WithSoftDelete bool `json:"with_soft_delete,omitempty"`
	// EmptyStateMessage is shown when no items exist.
	EmptyStateMessage string `json:"empty_state_message,omitempty"`
	// Columns is the list of table columns (for table views).
	Columns []ColumnDef `json:"columns,omitempty"`
	// Fields is the list of form fields (for form views).
	Fields []FieldDef `json:"fields,omitempty"`
	// RowActions is the list of row actions (for table views).
	RowActions []RowActionDef `json:"row_actions,omitempty"`
	// SubmitURL is the HTMX post URL (for form views).
	SubmitURL string `json:"submit_url,omitempty"`
	// Method is POST or PUT (for form views).
	Method string `json:"method,omitempty"`
	// SuccessRedirect is the HX-Redirect target (for form views).
	SuccessRedirect string `json:"success_redirect,omitempty"`
}

// ColumnDef defines a table column.
type ColumnDef struct {
	// Key is the field key.
	Key string `json:"key"`
	// Label is the display label.
	Label string `json:"label"`
	// Sortable enables sorting on this column.
	Sortable bool `json:"sortable,omitempty"`
	// Format is the display format: text, currency, date, datetime, badge, link.
	Format string `json:"format,omitempty"`
	// Width is the optional CSS width.
	Width string `json:"width,omitempty"`
	// BadgeConfig contains badge configuration (for badge format).
	BadgeConfig map[string]string `json:"badge_config,omitempty"`
}

// ScaffoldViewInput is the input for the scaffold_view tool.
type ScaffoldViewInput struct {
	// DomainName is the domain name.
	DomainName string `json:"domain_name"`
	// ViewType is list, show, form, card, table, or custom.
	ViewType string `json:"view_type"`
	// ViewName is the view file name (without extension).
	ViewName string `json:"view_name"`
	// Config contains view-specific configuration.
	Config ViewConfig `json:"config,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ScaffoldFormInput is the input for the scaffold_form tool.
type ScaffoldFormInput struct {
	// FormName is the form component name.
	FormName string `json:"form_name"`
	// Domain is the domain name.
	Domain string `json:"domain"`
	// Action is create or edit.
	Action string `json:"action"`
	// Fields is the list of form fields.
	Fields []FieldDef `json:"fields"`
	// SubmitEndpoint is the form submission URL.
	SubmitEndpoint string `json:"submit_endpoint,omitempty"`
	// ValidationRules contains field validation rules.
	ValidationRules map[string]string `json:"validation_rules,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// RowActionDef defines a table row action.
type RowActionDef struct {
	// Type is view, edit, delete, or custom.
	Type string `json:"type"`
	// Label is the action button label.
	Label string `json:"label,omitempty"`
	// URL is the action URL (for custom actions).
	URL string `json:"url,omitempty"`
	// Confirm shows a confirmation dialog.
	Confirm bool `json:"confirm,omitempty"`
	// ConfirmMessage is the confirmation message.
	ConfirmMessage string `json:"confirm_message,omitempty"`
}

// ScaffoldTableInput is the input for the scaffold_table tool.
type ScaffoldTableInput struct {
	// TableName is the table component name.
	TableName string `json:"table_name"`
	// Domain is the domain name.
	Domain string `json:"domain"`
	// Columns is the list of table columns.
	Columns []ColumnDef `json:"columns"`
	// WithPagination enables pagination. Defaults to true.
	WithPagination *bool `json:"with_pagination,omitempty"`
	// WithSorting enables column sorting. Defaults to true.
	WithSorting *bool `json:"with_sorting,omitempty"`
	// WithSearch enables search functionality. Defaults to true.
	WithSearch *bool `json:"with_search,omitempty"`
	// WithBulkActions enables bulk actions. Defaults to false.
	WithBulkActions bool `json:"with_bulk_actions,omitempty"`
	// RowActions is the list of row actions.
	RowActions []RowActionDef `json:"row_actions,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// GetWithPagination returns the WithPagination value with default true.
func (s ScaffoldTableInput) GetWithPagination() bool {
	if s.WithPagination == nil {
		return true
	}
	return *s.WithPagination
}

// GetWithSorting returns the WithSorting value with default true.
func (s ScaffoldTableInput) GetWithSorting() bool {
	if s.WithSorting == nil {
		return true
	}
	return *s.WithSorting
}

// GetWithSearch returns the WithSearch value with default true.
func (s ScaffoldTableInput) GetWithSearch() bool {
	if s.WithSearch == nil {
		return true
	}
	return *s.WithSearch
}

// TriggerConfig defines how a modal is triggered.
type TriggerConfig struct {
	// ButtonText is the trigger button text.
	ButtonText string `json:"button_text,omitempty"`
	// ButtonVariant is the button variant: default, outline, destructive.
	ButtonVariant string `json:"button_variant,omitempty"`
	// HTMXURL is the URL to load content from.
	HTMXURL string `json:"htmx_url,omitempty"`
}

// ScaffoldModalInput is the input for the scaffold_modal tool.
type ScaffoldModalInput struct {
	// ModalName is the modal component name.
	ModalName string `json:"modal_name"`
	// ModalType is dialog, sheet, or confirm.
	ModalType string `json:"modal_type"`
	// Title is the modal title.
	Title string `json:"title,omitempty"`
	// ContentType is form, info, or confirm.
	ContentType string `json:"content_type,omitempty"`
	// TriggerConfig defines how the modal is triggered.
	TriggerConfig TriggerConfig `json:"trigger_config,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// PropDef defines a component property.
type PropDef struct {
	// Name is the property name.
	Name string `json:"name"`
	// Type is the Go type.
	Type string `json:"type"`
	// Default is the default value.
	Default string `json:"default,omitempty"`
	// Required indicates if the prop is required.
	Required bool `json:"required,omitempty"`
}

// ScaffoldComponentInput is the input for the scaffold_component tool.
type ScaffoldComponentInput struct {
	// ComponentName is the component name.
	ComponentName string `json:"component_name"`
	// ComponentType is card, modal, dropdown, form_field, table, or custom.
	ComponentType string `json:"component_type,omitempty"`
	// Props is the list of component properties.
	Props []PropDef `json:"props,omitempty"`
	// WithHTMX includes HTMX attributes.
	WithHTMX bool `json:"with_htmx,omitempty"`
	// AlpineState contains Alpine.js state if needed.
	AlpineState map[string]interface{} `json:"alpine_state,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// SectionDef defines a page section.
type SectionDef struct {
	// Type is the section type: hero, content, table, cards, form.
	Type string `json:"type"`
	// Config contains section-specific configuration.
	Config map[string]interface{} `json:"config,omitempty"`
}

// ScaffoldPageInput is the input for the scaffold_page tool.
type ScaffoldPageInput struct {
	// PageName is the page name.
	PageName string `json:"page_name"`
	// Route is the URL route.
	Route string `json:"route"`
	// Layout is the layout type: default, dashboard, landing, blank.
	Layout string `json:"layout,omitempty"`
	// Sections is the list of page sections.
	Sections []SectionDef `json:"sections,omitempty"`
	// CreateTomlConfig generates a TOML config file.
	CreateTomlConfig bool `json:"create_toml_config,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ScaffoldConfigInput is the input for the scaffold_config tool.
type ScaffoldConfigInput struct {
	// ConfigType is page, menu, app, or messages.
	ConfigType string `json:"config_type"`
	// Name is the config file name (without extension).
	Name string `json:"name"`
	// Locale is the locale code. Defaults to "en".
	Locale string `json:"locale,omitempty"`
	// Content is the configuration content.
	Content map[string]interface{} `json:"content,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ScaffoldSeedInput is the input for the scaffold_seed tool.
type ScaffoldSeedInput struct {
	// Domain is the domain name.
	Domain string `json:"domain"`
	// Fields is the list of fields to seed.
	Fields []FieldDef `json:"fields,omitempty"`
	// Count is the number of records to seed.
	Count int `json:"count,omitempty"`
	// WithFaker uses faker for realistic data.
	WithFaker bool `json:"with_faker,omitempty"`
	// Dependencies is the list of seeders to run first.
	Dependencies []string `json:"dependencies,omitempty"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ListDomainsInput is the input for the list_domains tool.
type ListDomainsInput struct {
	// No input required - scans the project structure.
}

// UpdateDIWiringInput is the input for the update_di_wiring tool.
type UpdateDIWiringInput struct {
	// Domains is the list of domains to wire.
	Domains []string `json:"domains"`
	// DryRun previews changes without writing files.
	DryRun bool `json:"dry_run,omitempty"`
}

// ReportBugInput is the input for the report_bug tool.
type ReportBugInput struct {
	// Title is a short summary of the bug.
	Title string `json:"title"`
	// Description is a detailed description of the bug.
	Description string `json:"description"`
}
