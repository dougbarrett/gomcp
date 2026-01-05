// Package utils provides utility functions for scaffolding operations.
package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// validProjectNameRegex matches valid project names: alphanumeric, hyphens, underscores.
var validProjectNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// validModulePathRegex matches valid Go module paths.
var validModulePathRegex = regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9_.]*(/[-a-zA-Z0-9_.]+)*$`)

// validIdentifierRegex matches valid Go identifiers.
var validIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// validDatabaseTypes are the supported database types.
var validDatabaseTypes = map[string]bool{
	"":         true, // empty defaults to sqlite
	"sqlite":   true,
	"postgres": true,
	"mysql":    true,
}

// validViewTypes are the supported view types.
var validViewTypes = map[string]bool{
	"list":   true,
	"show":   true,
	"form":   true,
	"card":   true,
	"table":  true,
	"custom": true,
}

// validFormTypes are the supported form field types.
var validFormTypes = map[string]bool{
	"":         true, // empty defaults to input
	"input":    true,
	"textarea": true,
	"select":   true,
	"checkbox": true,
	"switch":   true,
	"date":     true,
	"time":     true,
	"datetime": true,
	"email":    true,
	"password": true,
	"number":   true,
	"rating":   true,
	"tags":     true,
	"slider":   true,
}

// validGoTypes are commonly valid Go types for model fields.
var validGoTypes = map[string]bool{
	"string":     true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"float32":    true,
	"float64":    true,
	"bool":       true,
	"time.Time":  true,
	"*time.Time": true,
	"[]byte":     true,
	"[]string":   true,
}

// validModalTypes are the supported modal types.
var validModalTypes = map[string]bool{
	"dialog":  true,
	"sheet":   true,
	"confirm": true,
}

// validConfigTypes are the supported config types.
var validConfigTypes = map[string]bool{
	"page":     true,
	"menu":     true,
	"app":      true,
	"messages": true,
}

// validLayoutTypes are the supported layout types.
var validLayoutTypes = map[string]bool{
	"":          true, // empty defaults to "default"
	"default":   true,
	"dashboard": true,
	"landing":   true,
	"blank":     true,
}

// goReservedWords are Go language reserved keywords.
var goReservedWords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
}

// commonReservedNames are common reserved directory/package names.
var commonReservedNames = map[string]bool{
	"internal": true, "vendor": true, "testdata": true, "cmd": true,
	"pkg": true, "api": true, "web": true, "test": true, "tests": true,
}

const (
	maxProjectNameLength = 128
	maxModulePathLength  = 256
	maxDomainNameLength  = 64
	maxFieldNameLength   = 64
)

// ValidateProjectName validates a project name.
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}

	if len(name) > maxProjectNameLength {
		return fmt.Errorf("project name is too long (max %d characters)", maxProjectNameLength)
	}

	if !validProjectNameRegex.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only alphanumeric characters, hyphens, and underscores")
	}

	if goReservedWords[strings.ToLower(name)] {
		return fmt.Errorf("project name '%s' is a Go reserved word", name)
	}

	return nil
}

// ValidateModulePath validates a Go module path.
func ValidateModulePath(path string) error {
	if path == "" {
		return fmt.Errorf("module path is required")
	}

	if len(path) > maxModulePathLength {
		return fmt.Errorf("module path is too long (max %d characters)", maxModulePathLength)
	}

	// Check for spaces
	if strings.ContainsAny(path, " \t\n\r") {
		return fmt.Errorf("module path cannot contain whitespace")
	}

	// Basic format check
	if !validModulePathRegex.MatchString(path) {
		return fmt.Errorf("invalid module path format: %s", path)
	}

	// Check for double slashes
	if strings.Contains(path, "//") {
		return fmt.Errorf("module path cannot contain double slashes")
	}

	// Check for trailing slash
	if strings.HasSuffix(path, "/") {
		return fmt.Errorf("module path cannot end with a slash")
	}

	return nil
}

// ValidateDatabaseType validates a database type.
func ValidateDatabaseType(dbType string) error {
	if !validDatabaseTypes[dbType] {
		return fmt.Errorf("invalid database type '%s': must be sqlite, postgres, or mysql", dbType)
	}
	return nil
}

// ValidateDomainName validates a domain name.
func ValidateDomainName(name string) error {
	if name == "" {
		return fmt.Errorf("domain name is required")
	}

	if len(name) > maxDomainNameLength {
		return fmt.Errorf("domain name is too long (max %d characters)", maxDomainNameLength)
	}

	// Check first character
	if !unicode.IsLetter(rune(name[0])) && name[0] != '_' {
		return fmt.Errorf("domain name must start with a letter or underscore")
	}

	// Check all characters
	for i, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("domain name contains invalid character '%c' at position %d", r, i)
		}
	}

	// Check for reserved words
	lower := strings.ToLower(name)
	if goReservedWords[lower] {
		return fmt.Errorf("domain name '%s' is a Go reserved word", name)
	}

	if commonReservedNames[lower] {
		return fmt.Errorf("domain name '%s' is a reserved name", name)
	}

	return nil
}

// ValidateFieldName validates a field name.
func ValidateFieldName(name string) error {
	if name == "" {
		return fmt.Errorf("field name is required")
	}

	if len(name) > maxFieldNameLength {
		return fmt.Errorf("field name is too long (max %d characters)", maxFieldNameLength)
	}

	if !validIdentifierRegex.MatchString(name) {
		return fmt.Errorf("field name '%s' is not a valid Go identifier", name)
	}

	// Check first character is uppercase (exported field)
	if !unicode.IsUpper(rune(name[0])) {
		return fmt.Errorf("field name '%s' must start with an uppercase letter (exported)", name)
	}

	if goReservedWords[strings.ToLower(name)] {
		return fmt.Errorf("field name '%s' is a Go reserved word", name)
	}

	return nil
}

// ValidateFieldType validates a Go type for a field.
func ValidateFieldType(fieldType string) error {
	if fieldType == "" {
		return fmt.Errorf("field type is required")
	}

	// Check against known valid types
	if validGoTypes[fieldType] {
		return nil
	}

	// Allow pointer types to known types
	if strings.HasPrefix(fieldType, "*") {
		baseType := strings.TrimPrefix(fieldType, "*")
		if validGoTypes[baseType] {
			return nil
		}
	}

	// Allow slice types
	if strings.HasPrefix(fieldType, "[]") {
		baseType := strings.TrimPrefix(fieldType, "[]")
		if validGoTypes[baseType] || validGoTypes["[]"+baseType] {
			return nil
		}
	}

	// Allow custom types that look like valid identifiers (e.g., "CustomType", "models.Status")
	parts := strings.Split(fieldType, ".")
	for _, part := range parts {
		cleanPart := strings.TrimPrefix(part, "*")
		cleanPart = strings.TrimPrefix(cleanPart, "[]")
		if cleanPart != "" && !validIdentifierRegex.MatchString(cleanPart) {
			return fmt.Errorf("field type '%s' contains invalid identifier '%s'", fieldType, cleanPart)
		}
	}

	return nil
}

// ValidateFormType validates a form field type.
func ValidateFormType(formType string) error {
	if !validFormTypes[formType] {
		return fmt.Errorf("invalid form type '%s': must be one of input, textarea, select, checkbox, switch, date, time, datetime, email, password, number, rating, tags, slider", formType)
	}
	return nil
}

// ValidateViewType validates a view type.
func ValidateViewType(viewType string) error {
	if viewType == "" {
		return fmt.Errorf("view type is required")
	}
	if !validViewTypes[viewType] {
		return fmt.Errorf("invalid view type '%s': must be one of list, show, form, card, table, custom", viewType)
	}
	return nil
}

// ValidateModalType validates a modal type.
func ValidateModalType(modalType string) error {
	if modalType == "" {
		return fmt.Errorf("modal type is required")
	}
	if !validModalTypes[modalType] {
		return fmt.Errorf("invalid modal type '%s': must be one of dialog, sheet, confirm", modalType)
	}
	return nil
}

// ValidateConfigType validates a config type.
func ValidateConfigType(configType string) error {
	if configType == "" {
		return fmt.Errorf("config type is required")
	}
	if !validConfigTypes[configType] {
		return fmt.Errorf("invalid config type '%s': must be one of page, menu, app, messages", configType)
	}
	return nil
}

// ValidateLayoutType validates a layout type.
func ValidateLayoutType(layoutType string) error {
	if !validLayoutTypes[layoutType] {
		return fmt.Errorf("invalid layout type '%s': must be one of default, dashboard, landing, blank", layoutType)
	}
	return nil
}

// ValidateHTTPMethod validates an HTTP method.
func ValidateHTTPMethod(method string) error {
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true,
	}
	upper := strings.ToUpper(method)
	if !validMethods[upper] {
		return fmt.Errorf("invalid HTTP method '%s': must be GET, POST, PUT, PATCH, or DELETE", method)
	}
	return nil
}

// ValidateURLPath validates a URL path.
func ValidateURLPath(path string) error {
	if path == "" {
		return nil // empty is allowed, will use default
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("URL path must start with /")
	}
	// Check for invalid characters
	invalidChars := regexp.MustCompile(`[^a-zA-Z0-9/_\-{}:]`)
	if invalidChars.MatchString(path) {
		return fmt.Errorf("URL path contains invalid characters")
	}
	return nil
}

// ValidateLocale validates a locale code.
func ValidateLocale(locale string) error {
	if locale == "" {
		return nil // empty defaults to "en"
	}
	// Simple locale format: xx or xx-XX
	localeRegex := regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)
	if !localeRegex.MatchString(locale) {
		return fmt.Errorf("invalid locale '%s': must be in format 'xx' or 'xx-XX'", locale)
	}
	return nil
}

// ValidateComponentName validates a component name.
func ValidateComponentName(name string) error {
	if name == "" {
		return fmt.Errorf("component name is required")
	}
	if !validIdentifierRegex.MatchString(name) {
		return fmt.Errorf("component name '%s' is not a valid identifier", name)
	}
	return nil
}

// validRelationshipTypes are the supported relationship types.
var validRelationshipTypes = map[string]bool{
	"belongs_to":   true,
	"has_one":      true,
	"has_many":     true,
	"many_to_many": true,
}

// validOnDeleteActions are the supported ON DELETE actions.
var validOnDeleteActions = map[string]bool{
	"":          true, // empty defaults to CASCADE
	"CASCADE":   true,
	"SET NULL":  true,
	"RESTRICT":  true,
	"NO ACTION": true,
}

// ValidateRelationshipType validates a relationship type.
func ValidateRelationshipType(relType string) error {
	if relType == "" {
		return fmt.Errorf("relationship type is required")
	}
	if !validRelationshipTypes[relType] {
		return fmt.Errorf("invalid relationship type '%s': must be belongs_to, has_one, has_many, or many_to_many", relType)
	}
	return nil
}

// ValidateRelationshipModel validates a related model name.
func ValidateRelationshipModel(model string) error {
	if model == "" {
		return fmt.Errorf("related model name is required")
	}
	// Model should be PascalCase identifier
	if !validIdentifierRegex.MatchString(model) {
		return fmt.Errorf("related model '%s' is not a valid identifier", model)
	}
	// First character should be uppercase
	if !unicode.IsUpper(rune(model[0])) {
		return fmt.Errorf("related model '%s' must be in PascalCase", model)
	}
	return nil
}

// ValidateOnDelete validates an ON DELETE action.
func ValidateOnDelete(action string) error {
	upper := strings.ToUpper(action)
	if !validOnDeleteActions[upper] {
		return fmt.Errorf("invalid ON DELETE action '%s': must be CASCADE, SET NULL, RESTRICT, or NO ACTION", action)
	}
	return nil
}
