package types

import (
	"fmt"
	"strings"
)

// FileConflict represents a file that would be overwritten.
type FileConflict struct {
	// Path is the relative file path.
	Path string `json:"path"`
	// Description explains the purpose of this file.
	Description string `json:"description"`
	// ProposedContent is the content that would be written.
	ProposedContent string `json:"proposed_content"`
}

// ToolHint suggests a tool that could be called next.
type ToolHint struct {
	// Tool is the tool name (e.g., "scaffold_domain").
	Tool string `json:"tool"`
	// Description explains when/why to use this tool.
	Description string `json:"description"`
	// Example shows a sample invocation (optional).
	Example string `json:"example,omitempty"`
	// Priority indicates importance: "recommended" or "optional".
	Priority string `json:"priority"`
}

// ScaffoldResult is the result returned by scaffolding tools.
type ScaffoldResult struct {
	// Success indicates if the operation succeeded.
	Success bool `json:"success"`
	// Message describes the result.
	Message string `json:"message"`
	// FilesCreated is the list of files that were created.
	FilesCreated []string `json:"files_created,omitempty"`
	// FilesUpdated is the list of files that were updated.
	FilesUpdated []string `json:"files_updated,omitempty"`
	// NextSteps is the list of suggested next actions (shell commands).
	NextSteps []string `json:"next_steps,omitempty"`
	// SuggestedTools hints at which MCP tools to call next.
	SuggestedTools []ToolHint `json:"suggested_tools,omitempty"`
	// Conflicts is the list of files that would be overwritten.
	// When conflicts exist, Success is false and no files are written.
	Conflicts []FileConflict `json:"conflicts,omitempty"`
	// ConflictsXML is a structured XML representation of conflicts for LLM consumption.
	ConflictsXML string `json:"conflicts_xml,omitempty"`
}

// NewConflictResult creates a result indicating file conflicts that would overwrite existing files.
func NewConflictResult(conflicts []FileConflict) ScaffoldResult {
	return ScaffoldResult{
		Success:      false,
		Message:      fmt.Sprintf("Cannot proceed: %d file(s) already exist and would be overwritten. Review the proposed changes below.", len(conflicts)),
		Conflicts:    conflicts,
		ConflictsXML: GenerateConflictsXML(conflicts),
	}
}

// GenerateConflictsXML creates a structured XML representation of file conflicts.
// This format is designed for easy LLM parsing and decision-making.
func GenerateConflictsXML(conflicts []FileConflict) string {
	if len(conflicts) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("<file_conflicts>\n")
	sb.WriteString("  <summary>\n")
	sb.WriteString(fmt.Sprintf("    <total_conflicts>%d</total_conflicts>\n", len(conflicts)))
	sb.WriteString("    <action_required>Review each file and decide whether to apply the proposed changes manually or skip them.</action_required>\n")
	sb.WriteString("  </summary>\n")
	sb.WriteString("  <conflicts>\n")

	for i, conflict := range conflicts {
		sb.WriteString(fmt.Sprintf("    <file index=\"%d\">\n", i+1))
		sb.WriteString(fmt.Sprintf("      <path>%s</path>\n", escapeXML(conflict.Path)))
		sb.WriteString(fmt.Sprintf("      <description>%s</description>\n", escapeXML(conflict.Description)))
		sb.WriteString("      <proposed_content>\n")
		sb.WriteString("<![CDATA[\n")
		sb.WriteString(conflict.ProposedContent)
		if !strings.HasSuffix(conflict.ProposedContent, "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString("]]>\n")
		sb.WriteString("      </proposed_content>\n")
		sb.WriteString("      <suggested_actions>\n")
		sb.WriteString("        <action type=\"manual_merge\">Compare with existing file and merge changes manually</action>\n")
		sb.WriteString("        <action type=\"skip\">Keep the existing file unchanged</action>\n")
		sb.WriteString("        <action type=\"backup_and_replace\">Backup existing file, then apply proposed content</action>\n")
		sb.WriteString("      </suggested_actions>\n")
		sb.WriteString("    </file>\n")
	}

	sb.WriteString("  </conflicts>\n")
	sb.WriteString("</file_conflicts>")

	return sb.String()
}

// escapeXML escapes special characters for XML content.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// NewSuccessResult creates a success result with a message.
func NewSuccessResult(message string) ScaffoldResult {
	return ScaffoldResult{
		Success: true,
		Message: message,
	}
}

// NewSuccessResultWithFiles creates a success result with file lists.
func NewSuccessResultWithFiles(message string, created, updated []string) ScaffoldResult {
	return ScaffoldResult{
		Success:      true,
		Message:      message,
		FilesCreated: created,
		FilesUpdated: updated,
	}
}

// NewErrorResult creates an error result with a message.
func NewErrorResult(message string) ScaffoldResult {
	return ScaffoldResult{
		Success: false,
		Message: message,
	}
}

// WithNextSteps adds next steps to the result.
func (r ScaffoldResult) WithNextSteps(steps ...string) ScaffoldResult {
	r.NextSteps = steps
	return r
}

// WithFilesCreated adds created files to the result.
func (r ScaffoldResult) WithFilesCreated(files ...string) ScaffoldResult {
	r.FilesCreated = append(r.FilesCreated, files...)
	return r
}

// WithFilesUpdated adds updated files to the result.
func (r ScaffoldResult) WithFilesUpdated(files ...string) ScaffoldResult {
	r.FilesUpdated = append(r.FilesUpdated, files...)
	return r
}

// WithSuggestedTools adds tool hints to the result.
func (r ScaffoldResult) WithSuggestedTools(tools ...ToolHint) ScaffoldResult {
	r.SuggestedTools = append(r.SuggestedTools, tools...)
	return r
}

// Common tool hints for reuse across scaffolding tools.
var (
	HintScaffoldDomain = ToolHint{
		Tool:        "scaffold_domain",
		Description: "Add a new feature with model, repository, service, and controller",
		Example:     `scaffold_domain: { domain_name: "product", fields: [{ name: "Name", type: "string" }] }`,
		Priority:    "recommended",
	}
	HintUpdateDIWiring = ToolHint{
		Tool:        "update_di_wiring",
		Description: "Wire the new domain into main.go dependency injection",
		Priority:    "recommended",
	}
	HintScaffoldForm = ToolHint{
		Tool:        "scaffold_form",
		Description: "Create an HTMX-powered form for create/edit operations",
		Priority:    "optional",
	}
	HintScaffoldTable = ToolHint{
		Tool:        "scaffold_table",
		Description: "Create a data table with sorting, pagination, and actions",
		Priority:    "optional",
	}
	HintScaffoldView = ToolHint{
		Tool:        "scaffold_view",
		Description: "Create additional views (list, show, card)",
		Priority:    "optional",
	}
	HintExtendService = ToolHint{
		Tool:        "extend_service",
		Description: "Add custom business logic methods to the service",
		Priority:    "optional",
	}
	HintExtendRepository = ToolHint{
		Tool:        "extend_repository",
		Description: "Add custom data access methods to the repository",
		Priority:    "optional",
	}
	HintExtendController = ToolHint{
		Tool:        "extend_controller",
		Description: "Add custom HTTP endpoints to the controller",
		Priority:    "optional",
	}
	HintScaffoldSeed = ToolHint{
		Tool:        "scaffold_seed",
		Description: "Create a database seeder for test data",
		Priority:    "optional",
	}
	HintScaffoldPage = ToolHint{
		Tool:        "scaffold_page",
		Description: "Create a standalone page with layout",
		Priority:    "optional",
	}
)

// DomainInfo describes a scaffolded domain.
type DomainInfo struct {
	// Name is the domain name (e.g., "product").
	Name string `json:"name"`
	// HasModel indicates if a model file exists.
	HasModel bool `json:"has_model"`
	// HasRepository indicates if a repository exists.
	HasRepository bool `json:"has_repository"`
	// HasService indicates if a service exists.
	HasService bool `json:"has_service"`
	// HasController indicates if a controller exists.
	HasController bool `json:"has_controller"`
	// Views is the list of view files.
	Views []string `json:"views,omitempty"`
}

// ListDomainsResult is the result of the list_domains tool.
type ListDomainsResult struct {
	// Success indicates if the operation succeeded.
	Success bool `json:"success"`
	// Message describes the result.
	Message string `json:"message,omitempty"`
	// Domains is the list of domain info.
	Domains []DomainInfo `json:"domains,omitempty"`
}

// NewListDomainsResult creates a successful list domains result.
func NewListDomainsResult(domains []DomainInfo) ListDomainsResult {
	return ListDomainsResult{
		Success: true,
		Domains: domains,
	}
}

// NewListDomainsError creates an error list domains result.
func NewListDomainsError(message string) ListDomainsResult {
	return ListDomainsResult{
		Success: false,
		Message: message,
	}
}

// ReportBugResult is the result of the report_bug tool.
type ReportBugResult struct {
	// Success indicates if the bug was reported successfully.
	Success bool `json:"success"`
	// Message describes the result.
	Message string `json:"message"`
	// BugID is the ID of the created bug.
	BugID string `json:"bug_id,omitempty"`
}

// NewReportBugResult creates a successful report bug result.
func NewReportBugResult(bugID string) ReportBugResult {
	return ReportBugResult{
		Success: true,
		Message: "Bug reported successfully",
		BugID:   bugID,
	}
}

// NewReportBugError creates an error report bug result.
func NewReportBugError(message string) ReportBugResult {
	return ReportBugResult{
		Success: false,
		Message: message,
	}
}
