// Package modifier provides code injection capabilities using marker comments.
package modifier

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dbb1dev/go-mcp/internal/utils"
)

// MarkerPrefix is the prefix used for marker comments.
const MarkerPrefix = "// MCP:"

// Common marker names
const (
	MarkerModelsStart         = "MCP:MODELS:START"
	MarkerModelsEnd           = "MCP:MODELS:END"
	MarkerReposStart          = "MCP:REPOS:START"
	MarkerReposEnd            = "MCP:REPOS:END"
	MarkerServicesStart       = "MCP:SERVICES:START"
	MarkerServicesEnd         = "MCP:SERVICES:END"
	MarkerControllersStart    = "MCP:CONTROLLERS:START"
	MarkerControllersEnd      = "MCP:CONTROLLERS:END"
	MarkerRoutesStart         = "MCP:ROUTES:START"
	MarkerRoutesEnd           = "MCP:ROUTES:END"
	MarkerImportsStart        = "MCP:IMPORTS:START"
	MarkerImportsEnd          = "MCP:IMPORTS:END"
	MarkerRelationshipsStart  = "MCP:RELATIONSHIPS:START"
	MarkerRelationshipsEnd    = "MCP:RELATIONSHIPS:END"
	// Route group markers
	MarkerRoutesPublicStart        = "MCP:ROUTES:PUBLIC:START"
	MarkerRoutesPublicEnd          = "MCP:ROUTES:PUBLIC:END"
	MarkerRoutesAuthenticatedStart = "MCP:ROUTES:AUTHENTICATED:START"
	MarkerRoutesAuthenticatedEnd   = "MCP:ROUTES:AUTHENTICATED:END"
	MarkerRoutesAdminStart         = "MCP:ROUTES:ADMIN:START"
	MarkerRoutesAdminEnd           = "MCP:ROUTES:ADMIN:END"
	// Navigation item markers (in base_layout.templ)
	MarkerNavItemsStart      = "MCP:NAV_ITEMS:START"
	MarkerNavItemsEnd        = "MCP:NAV_ITEMS:END"
	MarkerNavItemsAdminStart = "MCP:NAV_ITEMS_ADMIN:START"
	MarkerNavItemsAdminEnd   = "MCP:NAV_ITEMS_ADMIN:END"
)

// Injector handles code injection into files using marker comments.
type Injector struct {
	filePath string
	content  string
}

// NewInjector creates a new injector for the given file.
func NewInjector(filePath string) (*Injector, error) {
	content, err := utils.ReadFileString(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return &Injector{
		filePath: filePath,
		content:  content,
	}, nil
}

// NewInjectorFromContent creates a new injector from content string.
func NewInjectorFromContent(content string) *Injector {
	return &Injector{
		content: content,
	}
}

// InjectBetweenMarkers injects code between START and END markers.
// It adds the code before the END marker, preserving existing content.
func (i *Injector) InjectBetweenMarkers(startMarker, endMarker, code string) error {
	// Find the markers
	startPattern := regexp.MustCompile(`(?m)^(\s*)//\s*` + regexp.QuoteMeta(startMarker) + `\s*$`)
	endPattern := regexp.MustCompile(`(?m)^(\s*)//\s*` + regexp.QuoteMeta(endMarker) + `\s*$`)

	startMatch := startPattern.FindStringSubmatchIndex(i.content)
	endMatch := endPattern.FindStringSubmatchIndex(i.content)

	if startMatch == nil {
		return fmt.Errorf("start marker not found: %s", startMarker)
	}
	if endMatch == nil {
		return fmt.Errorf("end marker not found: %s", endMarker)
	}
	if startMatch[0] >= endMatch[0] {
		return fmt.Errorf("start marker must come before end marker")
	}

	// Get the indentation from the end marker
	indent := ""
	if endMatch[2] != -1 && endMatch[3] != -1 {
		indent = i.content[endMatch[2]:endMatch[3]]
	}

	// Check if code already exists between markers
	existingContent := i.content[startMatch[1]:endMatch[0]]
	if strings.Contains(existingContent, strings.TrimSpace(code)) {
		// Code already exists, skip injection
		return nil
	}

	// Prepare the code with proper indentation
	indentedCode := indentCode(code, indent)

	// Insert before the end marker
	insertPos := endMatch[0]
	i.content = i.content[:insertPos] + indentedCode + "\n" + i.content[insertPos:]

	return nil
}

// InjectAfterMarker injects code after a marker.
func (i *Injector) InjectAfterMarker(marker, code string) error {
	pattern := regexp.MustCompile(`(?m)^(\s*)//\s*` + regexp.QuoteMeta(marker) + `\s*$`)
	match := pattern.FindStringSubmatchIndex(i.content)

	if match == nil {
		return fmt.Errorf("marker not found: %s", marker)
	}

	// Get the indentation
	indent := ""
	if match[2] != -1 && match[3] != -1 {
		indent = i.content[match[2]:match[3]]
	}

	// Prepare the code with proper indentation
	indentedCode := indentCode(code, indent)

	// Insert after the marker line
	insertPos := match[1]
	i.content = i.content[:insertPos] + "\n" + indentedCode + i.content[insertPos:]

	return nil
}

// InjectBeforeMarker injects code before a marker.
func (i *Injector) InjectBeforeMarker(marker, code string) error {
	pattern := regexp.MustCompile(`(?m)^(\s*)//\s*` + regexp.QuoteMeta(marker) + `\s*$`)
	match := pattern.FindStringSubmatchIndex(i.content)

	if match == nil {
		return fmt.Errorf("marker not found: %s", marker)
	}

	// Get the indentation
	indent := ""
	if match[2] != -1 && match[3] != -1 {
		indent = i.content[match[2]:match[3]]
	}

	// Prepare the code with proper indentation
	indentedCode := indentCode(code, indent)

	// Insert before the marker line
	insertPos := match[0]
	i.content = i.content[:insertPos] + indentedCode + "\n" + i.content[insertPos:]

	return nil
}

// ReplaceMarkerContent replaces all content between markers with new content.
func (i *Injector) ReplaceMarkerContent(startMarker, endMarker, code string) error {
	startPattern := regexp.MustCompile(`(?m)^(\s*)//\s*` + regexp.QuoteMeta(startMarker) + `\s*$`)
	endPattern := regexp.MustCompile(`(?m)^(\s*)//\s*` + regexp.QuoteMeta(endMarker) + `\s*$`)

	startMatch := startPattern.FindStringSubmatchIndex(i.content)
	endMatch := endPattern.FindStringSubmatchIndex(i.content)

	if startMatch == nil {
		return fmt.Errorf("start marker not found: %s", startMarker)
	}
	if endMatch == nil {
		return fmt.Errorf("end marker not found: %s", endMarker)
	}
	if startMatch[0] >= endMatch[0] {
		return fmt.Errorf("start marker must come before end marker")
	}

	// Get the indentation from the end marker
	indent := ""
	if endMatch[2] != -1 && endMatch[3] != -1 {
		indent = i.content[endMatch[2]:endMatch[3]]
	}

	// Prepare the code with proper indentation
	indentedCode := ""
	if code != "" {
		indentedCode = indentCode(code, indent) + "\n"
	}

	// Replace content between markers
	i.content = i.content[:startMatch[1]+1] + indentedCode + i.content[endMatch[0]:]

	return nil
}

// HasMarker checks if a marker exists in the content.
func (i *Injector) HasMarker(marker string) bool {
	pattern := regexp.MustCompile(`(?m)//\s*` + regexp.QuoteMeta(marker))
	return pattern.MatchString(i.content)
}

// Content returns the current content.
func (i *Injector) Content() string {
	return i.content
}

// Save writes the content back to the file.
func (i *Injector) Save() error {
	if i.filePath == "" {
		return fmt.Errorf("no file path set")
	}
	return utils.WriteFileString(i.filePath, i.content, true)
}

// SaveTo writes the content to the specified file.
func (i *Injector) SaveTo(filePath string) error {
	return utils.WriteFileString(filePath, i.content, true)
}

// indentCode adds indentation to each line of code.
func indentCode(code, indent string) string {
	lines := strings.Split(strings.TrimSpace(code), "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// InjectImport adds an import statement to the imports section.
func (i *Injector) InjectImport(importPath string) error {
	return i.InjectImportWithAlias(importPath, "")
}

// InjectImportWithAlias adds an import statement with an optional alias to the imports section.
func (i *Injector) InjectImportWithAlias(importPath, alias string) error {
	// Check if import already exists
	importPattern := regexp.MustCompile(`"` + regexp.QuoteMeta(importPath) + `"`)
	if importPattern.MatchString(i.content) {
		return nil // Import already exists
	}

	// Format the import statement
	importStmt := `"` + importPath + `"`
	if alias != "" {
		importStmt = alias + ` "` + importPath + `"`
	}

	// Try to inject using markers first
	if i.HasMarker(MarkerImportsStart) && i.HasMarker(MarkerImportsEnd) {
		return i.InjectBetweenMarkers(MarkerImportsStart, MarkerImportsEnd, importStmt)
	}

	// Fall back to finding import block
	importBlockPattern := regexp.MustCompile(`(?m)^import \(\n((?:\s+.*\n)*)\)`)
	match := importBlockPattern.FindStringSubmatchIndex(i.content)
	if match == nil {
		return fmt.Errorf("no import block found")
	}

	// Insert before the closing paren
	insertPos := match[1] - 1
	i.content = i.content[:insertPos] + "\t" + importStmt + "\n" + i.content[insertPos:]

	return nil
}

// InjectModel adds a model to the AutoMigrate call.
func (i *Injector) InjectModel(modelName string) error {
	modelCode := "&models." + modelName + "{},"
	return i.InjectBetweenMarkers(MarkerModelsStart, MarkerModelsEnd, modelCode)
}

// InjectRepo adds a repository instantiation.
func (i *Injector) InjectRepo(domainName, modulePath string) error {
	varName := utils.ToRepoVariableName(domainName)
	pkgAlias := utils.ToRepoImportAlias(domainName)
	code := fmt.Sprintf(`%s := %s.NewRepository(db)`, varName, pkgAlias)
	return i.InjectBetweenMarkers(MarkerReposStart, MarkerReposEnd, code)
}

// InjectService adds a service instantiation.
func (i *Injector) InjectService(domainName string) error {
	varName := utils.ToServiceVariableName(domainName)
	repoVarName := utils.ToRepoVariableName(domainName)
	pkgAlias := utils.ToServiceImportAlias(domainName)
	code := fmt.Sprintf(`%s := %s.NewService(%s)`, varName, pkgAlias, repoVarName)
	return i.InjectBetweenMarkers(MarkerServicesStart, MarkerServicesEnd, code)
}

// InjectController adds a controller instantiation.
func (i *Injector) InjectController(domainName string) error {
	return i.InjectControllerWithRelations(domainName, nil)
}

// InjectControllerWithRelations adds a controller instantiation with related services.
// relatedDomains is a list of domain names for belongs_to relationships that need their
// services injected into the controller.
func (i *Injector) InjectControllerWithRelations(domainName string, relatedDomains []string) error {
	varName := utils.ToControllerVariableName(domainName)
	serviceVarName := utils.ToServiceVariableName(domainName)
	pkgAlias := utils.ToControllerImportAlias(domainName)

	// Build constructor arguments
	args := serviceVarName
	for _, relDomain := range relatedDomains {
		relServiceVarName := utils.ToServiceVariableName(relDomain)
		args += ", " + relServiceVarName
	}

	code := fmt.Sprintf(`%s := %s.NewController(%s)`, varName, pkgAlias, args)
	return i.InjectBetweenMarkers(MarkerControllersStart, MarkerControllersEnd, code)
}

// InjectRoute adds a route registration to the default (public) route group.
// Routes are mounted at the default URL path (e.g., /products for "product" domain).
// Users can later modify the path in main.go to mount under custom prefixes like /admin/products.
func (i *Injector) InjectRoute(domainName string) error {
	return i.InjectRouteWithGroup(domainName, "public")
}

// InjectRouteWithGroup adds a route registration to the specified route group.
// Valid groups: "public" (no auth), "authenticated" (requires login), "admin" (requires admin role).
// Falls back to the general MCP:ROUTES markers if group-specific markers are not found.
func (i *Injector) InjectRouteWithGroup(domainName, routeGroup string) error {
	varName := utils.ToControllerVariableName(domainName)
	urlPath := utils.ToURLPath(domainName)

	// For authenticated routes, the chi.Router variable is 'r' inside the group
	routerVar := "router"
	if routeGroup == "authenticated" || routeGroup == "admin" {
		routerVar = "r"
	}

	code := fmt.Sprintf(`%s.Route("%s", %s.RegisterRoutes)`, routerVar, urlPath, varName)

	// Determine which markers to use based on route group
	var startMarker, endMarker string
	switch routeGroup {
	case "authenticated":
		startMarker = MarkerRoutesAuthenticatedStart
		endMarker = MarkerRoutesAuthenticatedEnd
	case "admin":
		startMarker = MarkerRoutesAdminStart
		endMarker = MarkerRoutesAdminEnd
	default: // "public" or empty
		startMarker = MarkerRoutesPublicStart
		endMarker = MarkerRoutesPublicEnd
	}

	// Try group-specific markers first
	if i.HasMarker(startMarker) && i.HasMarker(endMarker) {
		return i.InjectBetweenMarkers(startMarker, endMarker, code)
	}

	// Fall back to general routes markers with router variable
	code = fmt.Sprintf(`router.Route("%s", %s.RegisterRoutes)`, urlPath, varName)
	return i.InjectBetweenMarkers(MarkerRoutesStart, MarkerRoutesEnd, code)
}

// InjectRelationship adds a relationship field to a model struct.
// This is used to inject inverse relationships when scaffolding related domains.
func (i *Injector) InjectRelationship(fieldCode string) error {
	return i.InjectBetweenMarkers(MarkerRelationshipsStart, MarkerRelationshipsEnd, fieldCode)
}

// InjectNavItem adds a navigation item to the sidebar in base_layout.templ.
// routeGroup determines which nav section to add to: "authenticated" or "admin".
// icon should be a valid icon name (e.g., "folder", "list", "file-text").
func (i *Injector) InjectNavItem(domainName, routeGroup, icon string) error {
	urlPath := utils.ToURLPath(domainName)
	label := utils.ToLabel(domainName)
	// Pluralize the label for nav items (e.g., "Product" -> "Products")
	label = utils.Pluralize(label)

	// Default icon if not provided
	if icon == "" {
		icon = "folder"
	}

	// Generate the navItem templ call
	code := fmt.Sprintf(`@navItem("%s", "%s", "%s", false)`, urlPath, icon, label)

	// Determine which markers to use based on route group
	var startMarker, endMarker string
	switch routeGroup {
	case "admin":
		startMarker = MarkerNavItemsAdminStart
		endMarker = MarkerNavItemsAdminEnd
	default: // "authenticated" or empty
		startMarker = MarkerNavItemsStart
		endMarker = MarkerNavItemsEnd
	}

	// Check if markers exist
	if !i.HasMarker(startMarker) || !i.HasMarker(endMarker) {
		return fmt.Errorf("navigation markers not found: %s, %s", startMarker, endMarker)
	}

	return i.InjectBetweenMarkers(startMarker, endMarker, code)
}
