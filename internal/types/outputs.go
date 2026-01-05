package types

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
	// NextSteps is the list of suggested next actions.
	NextSteps []string `json:"next_steps,omitempty"`
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
