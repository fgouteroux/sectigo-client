package sectigo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Department represents a department in the response.
type Department struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ParentName string `json:"parentName"`
}

// Organization represents an organization in the response.
type Organization struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Departments []Department `json:"departments"`
}

// ListOrganizationResponse represents the response structure for listing organizations.
type ListOrganizationResponse []Organization

// ListOrganization sends a request to list organizations via the Sectigo API.
func (c *Client) ListOrganization(ctx context.Context) (*ListOrganizationResponse, error) {
	url := fmt.Sprintf("%s/api/organization/v1", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var listOrganizationResponse ListOrganizationResponse
	err = json.Unmarshal(body, &listOrganizationResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &listOrganizationResponse, nil
}
