package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dbb1dev/go-mcp/internal/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BugAPIURL is the endpoint for reporting bugs.
const BugAPIURL = "https://y3d5o56xre.execute-api.us-west-2.amazonaws.com/prod/bugs"

// RegisterReportBug registers the report_bug tool.
func RegisterReportBug(server *mcp.Server, registry *Registry) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "report_bug",
		Description: `Report issues or bugs with the scaffolding tools.

Use this when you encounter:
- Tool errors or unexpected behavior
- Missing functionality
- Incorrect code generation
- Documentation issues

Bugs are tracked and prioritized for fixes.`,
	}, func(ctx context.Context, req *mcp.CallToolRequest, input types.ReportBugInput) (*mcp.CallToolResult, types.ReportBugResult, error) {
		result, err := reportBug(input)
		if err != nil {
			return nil, types.NewReportBugError(err.Error()), nil
		}
		return nil, result, nil
	})
}

// bugPayload is the request body for creating a bug.
type bugPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// bugResponse is the response from the bug API.
type bugResponse struct {
	ID string `json:"id"`
}

func reportBug(input types.ReportBugInput) (types.ReportBugResult, error) {
	// Validate input
	if input.Title == "" {
		return types.NewReportBugError("title is required"), nil
	}
	if input.Description == "" {
		return types.NewReportBugError("description is required"), nil
	}

	// Create payload
	payload := bugPayload{
		Title:       input.Title,
		Description: input.Description,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return types.NewReportBugError(fmt.Sprintf("failed to marshal request: %v", err)), nil
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", BugAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return types.NewReportBugError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return types.NewReportBugError(fmt.Sprintf("failed to send request: %v", err)), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return types.NewReportBugError(fmt.Sprintf("API returned status %d", resp.StatusCode)), nil
	}

	// Parse response
	var bugResp bugResponse
	if err := json.NewDecoder(resp.Body).Decode(&bugResp); err != nil {
		return types.NewReportBugError(fmt.Sprintf("failed to parse response: %v", err)), nil
	}

	return types.NewReportBugResult(bugResp.ID), nil
}
