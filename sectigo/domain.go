package sectigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DomainRequest represents the structure of the JSON payload for the domain creation request.
type DomainRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Active      bool                `json:"active"`
	Delegations []DelegationRequest `json:"delegations"`
}

type DelegationRequest struct {
	OrgId     int      `json:"orgId"`
	CertTypes []string `json:"certTypes"`
}

// ApproveDelegationRequest represents the structure of the JSON payload for approving a delegation.
type ApproveDelegationRequest struct {
	OrgId int `json:"orgId"`
}

// DelegateDomainRequest represents the structure of the JSON payload for delegating a domain.
type DelegateDomainRequest struct {
	DomainIds []int    `json:"domainIds"`
	OrgId     int      `json:"orgId"`
	CertTypes []string `json:"certTypes"`
}

// ListDomainParams represents the parameters for listing domains.
type ListDomainParams struct {
	Size     int
	Position int
	Name     string
	State    string
	Status   string
	OrgId    int
}

// Domain represents a domain in the response.
type Domain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DomainDetails represents the detailed information of a domain.
type DomainDetails struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	DelegationStatus string `json:"delegationStatus"`
	State            string `json:"state"`
	ValidationStatus string `json:"validationStatus"`
	ValidationMethod string `json:"validationMethod"`
	DcvValidation    string `json:"dcvValidation"`
	DcvExpiration    string `json:"dcvExpiration"`
	CtLogMonitoring  struct {
		Enabled           bool   `json:"enabled"`
		IncludeSubdomains bool   `json:"includeSubdomains"`
		BucketId          string `json:"bucketId"`
	} `json:"ctLogMonitoring"`
	Delegations []struct {
		OrgId     int      `json:"orgId"`
		CertTypes []string `json:"certTypes"`
		Status    string   `json:"status"`
	} `json:"delegations"`
}

// ListDomainResponse represents the response structure for listing domains.
type ListDomainResponse struct {
	Domains    []Domain `json:"domains"`
	TotalCount int      `json:"total_count"`
}

// StartDomainCNameValidationRequest represents the structure of the JSON payload for starting CNAME validation.
type StartDomainCNameValidationRequest struct {
	Domain string `json:"domain"`
}

// StartDomainCNameValidationResponse represents the response structure for starting CNAME validation.
type StartDomainCNameValidationResponse struct {
	Host  string `json:"host"`
	Point string `json:"point"`
}

// SubmitDomainCNameValidationRequest represents the structure of the JSON payload for submitting CNAME validation.
type SubmitDomainCNameValidationRequest struct {
	Domain string `json:"domain"`
}

// SubmitDomainCNameValidationResponse represents the response structure for submitting CNAME validation.
type SubmitDomainCNameValidationResponse struct {
	OrderStatus string `json:"orderStatus"`
	Message     string `json:"message"`
	Status      string `json:"status"`
}

// GetDomainValidationStatusRequest represents the structure of the JSON payload for getting domain validation status.
type GetDomainValidationStatusRequest struct {
	Domain string `json:"domain"`
}

// GetDomainValidationStatusResponse represents the response structure for getting domain validation status.
type GetDomainValidationStatusResponse struct {
	Status         string `json:"status"`
	OrderStatus    string `json:"orderStatus"`
	ExpirationDate string `json:"expirationDate"`
}

// ListDomainValidationParams represents the parameters for searching domain validation statuses.
type ListDomainValidationParams struct {
	Size        int
	Position    int
	Domain      string
	Org         int
	Department  int
	DcvStatus   string
	OrderStatus string
	ExpiresIn   int
}

// DomainValidation represents a domain validation entry in the response.
type DomainValidation struct {
	Domain         string `json:"domain"`
	DcvStatus      string `json:"dcvStatus"`
	DcvOrderStatus string `json:"dcvOrderStatus"`
	DcvMethod      string `json:"dcvMethod"`
}

// ListDomainValidationResponse represents the response structure for listing domain validations.
type ListDomainValidationResponse struct {
	Domains    []DomainValidation `json:"domains"`
	TotalCount int                `json:"total_count"`
}

// GetDomainDetails sends a request to get detailed information about a specific domain via the Sectigo API.
func (c *Client) GetDomainDetails(ctx context.Context, domainID int) (*DomainDetails, error) {
	url := fmt.Sprintf("%s/api/domain/v1/%d", c.BaseURL, domainID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var domainDetails DomainDetails
	err = json.Unmarshal(body, &domainDetails)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &domainDetails, nil
}

// CreateDomain sends a request to create a new domain via the Sectigo API.
func (c *Client) CreateDomain(ctx context.Context, domainRequest DomainRequest) error {
	url := fmt.Sprintf("%s/api/domain/v1", c.BaseURL)
	jsonPayload, err := json.Marshal(domainRequest)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	_, _, err = c.sendRequest(ctx, req, http.StatusCreated)
	return err
}

// DeleteDomain sends a request to delete a domain via the Sectigo API.
func (c *Client) DeleteDomain(ctx context.Context, domainID int) error {
	url := fmt.Sprintf("%s/api/domain/v1/%d", c.BaseURL, domainID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	_, _, err = c.sendRequest(ctx, req, http.StatusNoContent)
	return err
}

// ApproveDelegation sends a request to approve a delegation via the Sectigo API.
func (c *Client) ApproveDelegation(ctx context.Context, domainID int, approveRequest ApproveDelegationRequest) error {
	url := fmt.Sprintf("%s/api/domain/v1/%d/delegation/approve", c.BaseURL, domainID)
	jsonPayload, err := json.Marshal(approveRequest)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	_, _, err = c.sendRequest(ctx, req, http.StatusOK)
	return err
}

// DelegateDomain sends a request to delegate a domain via the Sectigo API.
func (c *Client) DelegateDomain(ctx context.Context, delegateRequest DelegateDomainRequest) error {
	url := fmt.Sprintf("%s/api/domain/v1/delegation", c.BaseURL)
	jsonPayload, err := json.Marshal(delegateRequest)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	_, _, err = c.sendRequest(ctx, req, http.StatusOK)
	return err
}

// ListDomain sends a request to list domains via the Sectigo API with query parameters and parses the response.
func (c *Client) ListDomain(ctx context.Context, params ListDomainParams) (*ListDomainResponse, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/domain/v1", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("size", fmt.Sprintf("%d", params.Size))
	queryParams.Add("position", fmt.Sprintf("%d", params.Position))
	if params.Name != "" {
		queryParams.Add("name", params.Name)
	}
	if params.State != "" {
		queryParams.Add("state", params.State)
	}
	if params.Status != "" {
		queryParams.Add("status", params.Status)
	}
	if params.OrgId > 0 {
		queryParams.Add("orgId", fmt.Sprintf("%d", params.OrgId))
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var domains []Domain
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	listDomainResponse := ListDomainResponse{Domains: domains}
	totalCountHeader := resp.Header.Get("X-Total-Count")
	if totalCountHeader != "" {
		listDomainResponse.TotalCount, _ = strconv.Atoi(totalCountHeader)
	}

	return &listDomainResponse, nil
}

// ListAllDomain sends requests to list all domains by iterating through the results using the X-Total-Count header.
func (c *Client) ListAllDomain(ctx context.Context, params ListDomainParams) ([]Domain, error) {
	var allDomains []Domain
	position := 0
	size := 200

	for {
		params.Position = position
		params.Size = size
		listDomainResponse, err := c.ListDomain(ctx, params)
		if err != nil {
			return nil, err
		}

		allDomains = append(allDomains, listDomainResponse.Domains...)

		if len(listDomainResponse.Domains) < params.Size || position+params.Size >= listDomainResponse.TotalCount {
			break
		}

		position += params.Size
	}

	return allDomains, nil
}

// StartDomainCNameValidation sends a request to start CNAME validation for a domain via the Sectigo API.
func (c *Client) StartDomainCNameValidation(ctx context.Context, request StartDomainCNameValidationRequest) (*StartDomainCNameValidationResponse, error) {
	url := fmt.Sprintf("%s/api/dcv/v1/validation/start/domain/cname", c.BaseURL)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var validationResponse StartDomainCNameValidationResponse
	err = json.Unmarshal(body, &validationResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &validationResponse, nil
}

// SubmitDomainCNameValidation sends a request to submit CNAME validation for a domain via the Sectigo API.
func (c *Client) SubmitDomainCNameValidation(ctx context.Context, request SubmitDomainCNameValidationRequest) (*SubmitDomainCNameValidationResponse, error) {
	url := fmt.Sprintf("%s/api/dcv/v1/validation/submit/domain/cname", c.BaseURL)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var validationResponse SubmitDomainCNameValidationResponse
	err = json.Unmarshal(body, &validationResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &validationResponse, nil
}

// GetDomainValidationStatus sends a request to get the validation status for a domain via the Sectigo API.
func (c *Client) GetDomainValidationStatus(ctx context.Context, request GetDomainValidationStatusRequest) (*GetDomainValidationStatusResponse, error) {
	url := fmt.Sprintf("%s/api/dcv/v2/validation/status", c.BaseURL)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var validationResponse GetDomainValidationStatusResponse
	err = json.Unmarshal(body, &validationResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &validationResponse, nil
}

// CheckDomainValidationStatus checks the domain validation status with retries.
func (c *Client) CheckDomainValidationStatus(ctx context.Context, domain string, maxRetries int, retryInterval time.Duration) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err := c.GetDomainValidationStatus(ctx, GetDomainValidationStatusRequest{Domain: domain})
		if err != nil {
			return err
		}

		if response.Status == "NOT_VALIDATED" {
			log.Println("Domain is not validated, retrying...")
			time.Sleep(retryInterval)
			continue
		}

		return nil
	}

	return fmt.Errorf("max retries reached, domain is still not validated")
}

// ListDomainValidation sends a request to search for domain validation statuses via the Sectigo API.
func (c *Client) ListDomainValidation(ctx context.Context, params ListDomainValidationParams) (*ListDomainValidationResponse, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/dcv/v1/validation", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("size", fmt.Sprintf("%d", params.Size))
	queryParams.Add("position", fmt.Sprintf("%d", params.Position))
	if params.Domain != "" {
		queryParams.Add("domain", params.Domain)
	}
	if params.Org > 0 {
		queryParams.Add("org", fmt.Sprintf("%d", params.Org))
	}
	if params.Department > 0 {
		queryParams.Add("department", fmt.Sprintf("%d", params.Department))
	}
	if params.DcvStatus != "" {
		queryParams.Add("dcvStatus", params.DcvStatus)
	}
	if params.OrderStatus != "" {
		queryParams.Add("orderStatus", params.OrderStatus)
	}
	if params.ExpiresIn > 0 {
		queryParams.Add("expiresIn", fmt.Sprintf("%d", params.ExpiresIn))
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var domainsValidation []DomainValidation
	err = json.Unmarshal(body, &domainsValidation)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	listDomainValidationResponse := ListDomainValidationResponse{Domains: domainsValidation}
	totalCountHeader := resp.Header.Get("X-Total-Count")
	if totalCountHeader != "" {
		listDomainValidationResponse.TotalCount, _ = strconv.Atoi(totalCountHeader)
	}

	return &listDomainValidationResponse, nil
}

// ListAllDomainValidation sends requests to list all domains by iterating through the results using the X-Total-Count header.
func (c *Client) ListAllDomainValidation(ctx context.Context, params ListDomainValidationParams) ([]DomainValidation, error) {
	var allDomainsValidation []DomainValidation
	position := 0
	size := 200

	for {
		params.Position = position
		params.Size = size
		listDomainValidationResponse, err := c.ListDomainValidation(ctx, params)
		if err != nil {
			return nil, err
		}

		allDomainsValidation = append(allDomainsValidation, listDomainValidationResponse.Domains...)

		if len(listDomainValidationResponse.Domains) < params.Size || position+params.Size >= listDomainValidationResponse.TotalCount {
			break
		}

		position += params.Size
	}

	return allDomainsValidation, nil
}