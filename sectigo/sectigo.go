package sectigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
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

// ListDomainResponse represents the response structure for listing domains.
type ListDomainValidationResponse struct {
	Domains    []DomainValidation `json:"domains"`
	TotalCount int                `json:"total_count"`
}

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

// AcmeAccount represents the acme account structure.
type AcmeAccount struct {
	ID                 int      `json:"id"`
	AccountID          string   `json:"accountId"`
	MacID              string   `json:"macId"`
	MacKey             string   `json:"macKey"`
	AcmeServer         string   `json:"acmeServer"`
	Name               string   `json:"name"`
	OrganizationID     int      `json:"organizationId"`
	CertValidationType string   `json:"certValidationType"`
	Status             string   `json:"status"`
	OvOrderNumber      int      `json:"ovOrderNumber"`
	OvAnchorID         string   `json:"ovAnchorID"`
	EvDetails          struct{} `json:"evDetails"`
	Contacts           string   `json:"contacts"`
}

// ListAcmeAccountResponse represents the response structure for listing acme accounts.
type ListAcmeAccountResponse struct {
	Accounts   []AcmeAccount `json:"accounts"`
	TotalCount int           `json:"total_count"`
}

// Define a struct for query parameters
type ListAcmeAccountParams struct {
	Position           int
	Size               int
	OrganizationId     int
	Name               string
	AcmeServer         string
	CertValidationType string
	Status             string
}

type ListAcmeAccountDomainParams struct {
	AccountID                   int
	Position                    int
	Size                        int
	Name                        string
	ExpiresWithinNextDays       int
	StickyExpiresWithinNextDays int
}

type ListAcmeAccountDomainResponse struct {
	Domains    []AcmeAccountDomain `json:"domains"`
	TotalCount int                 `json:"total_count"`
}

type AcmeAccountDomain struct {
	Name                string `json:"name"`
	ValidUntil          string `json:"validUntil"`
	StickyUntil         string `json:"stickyUntil"`
	OvAnchorOrderNumber int    `json:"ovAnchorOrderNumber"`
	OvAnchorID          string `json:"ovAnchorID"`
	EvAnchorOrderNumber int    `json:"evAnchorOrderNumber"`
	EvAnchorID          string `json:"evAnchorID"`
}

type AcmeAccountDomainParams struct {
	AccountID int
	Domains   []string
}

type AcmeAccountDomainName struct {
	Name string `json:"name"`
}

type AcmeAccountDomainRequest struct {
	Domains []AcmeAccountDomainName `json:"domains"`
}

// Client is a struct that holds the necessary information to make requests to the Sectigo API.
type Client struct {
	BaseURL string
	Client  *http.Client
	Debug   bool
}

// authTransport is a custom RoundTripper that adds authentication headers to each request.
type authTransport struct {
	login       string
	customerUri string
	password    string
	transport   http.RoundTripper
	debug       bool
}

type Config struct {
	URL      string
	Username string
	Customer string
	Password string
	Debug    bool
}

// RoundTrip implements the RoundTripper interface.
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("login", t.login)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("customerUri", t.customerUri)
	req.Header.Set("password", t.password)

	if t.debug {
		log.Printf("Request: %s %s\n", req.Method, req.URL.String())
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			log.Printf("Request Body: %s\n", string(body))
		}
	}

	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if t.debug {
		log.Printf("Response Status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		log.Printf("Response Body: %s\n", string(body))
	}

	return resp, nil
}

// NewClient initializes a new Sectigo API client with custom headers and optional debug mode.
func NewClient(config Config) *Client {
	// Create a new http.Client with the custom RoundTripper
	client := &http.Client{
		Transport: &authTransport{
			login:       config.Username,
			customerUri: config.Customer,
			password:    config.Password,
			transport:   http.DefaultTransport,
			debug:       config.Debug,
		},
	}

	return &Client{
		BaseURL: config.URL,
		Client:  client,
		Debug:   config.Debug,
	}
}

// sendRequest sends an HTTP request and returns the response body.
func (c *Client) sendRequest(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	req = req.WithContext(ctx)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, nil, fmt.Errorf("failed request, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	return resp, body, nil
}

// GetDomainDetails sends a request to get detailed information about a specific domain via the Sectigo API.
func (c *Client) GetDomainDetails(ctx context.Context, domainID int) (*DomainDetails, error) {
	url := fmt.Sprintf("%s/api/domain/v1/%d", c.BaseURL, domainID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req)
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

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create domain, status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteDomain sends a request to delete a domain via the Sectigo API.
func (c *Client) DeleteDomain(ctx context.Context, domainID int) error {
	url := fmt.Sprintf("%s/api/domain/v1/%d", c.BaseURL, domainID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete domain, status code: %d", resp.StatusCode)
	}

	return nil
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

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to approve delegation, status code: %d", resp.StatusCode)
	}

	return nil
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

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delegate domain, status code: %d", resp.StatusCode)
	}

	return nil
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

	resp, body, err := c.sendRequest(ctx, req)
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

	_, body, err := c.sendRequest(ctx, req)
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

	_, body, err := c.sendRequest(ctx, req)
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

	_, body, err := c.sendRequest(ctx, req)
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

	resp, body, err := c.sendRequest(ctx, req)
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

// ListOrganization sends a request to list organizations via the Sectigo API.
func (c *Client) ListOrganization(ctx context.Context) (*ListOrganizationResponse, error) {
	url := fmt.Sprintf("%s/api/organization/v1", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	_, body, err := c.sendRequest(ctx, req)
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

// ListAcmeAccount sends a request to list ACME accounts via the Sectigo API.
func (c *Client) ListAcmeAccount(ctx context.Context, params ListAcmeAccountParams) (*ListAcmeAccountResponse, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/acme/v2/account", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("size", fmt.Sprintf("%d", params.Size))
	queryParams.Add("position", fmt.Sprintf("%d", params.Position))
	queryParams.Add("organizationId", fmt.Sprintf("%d", params.OrganizationId))

	if params.Name != "" {
		queryParams.Add("name", params.Name)
	}
	if params.AcmeServer != "" {
		queryParams.Add("acmeServer", params.AcmeServer)
	}
	if params.CertValidationType != "" {
		queryParams.Add("certValidationType", params.CertValidationType)
	}
	if params.Status != "" {
		queryParams.Add("status", params.Status)
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, body, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var accounts []AcmeAccount
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	listAcmeAccountResponse := ListAcmeAccountResponse{Accounts: accounts}
	totalCountHeader := resp.Header.Get("X-Total-Count")
	if totalCountHeader != "" {
		listAcmeAccountResponse.TotalCount, _ = strconv.Atoi(totalCountHeader)
	}

	return &listAcmeAccountResponse, nil
}

// ListAllAcmeAccount sends requests to list all ACME accounts by iterating through the results using the X-Total-Count header.
func (c *Client) ListAllAcmeAccount(ctx context.Context, params ListAcmeAccountParams) ([]AcmeAccount, error) {
	var allAcmeAccounts []AcmeAccount
	position := 0
	size := 200

	for {
		params.Position = position
		params.Size = size
		listAcmeAccountResponse, err := c.ListAcmeAccount(ctx, params)
		if err != nil {
			return nil, err
		}

		allAcmeAccounts = append(allAcmeAccounts, listAcmeAccountResponse.Accounts...)

		if len(listAcmeAccountResponse.Accounts) < params.Size || position+params.Size >= listAcmeAccountResponse.TotalCount {
			break
		}

		position += params.Size
	}

	return allAcmeAccounts, nil
}

// ListAcmeAccountDomain sends a request to list ACME account domains via the Sectigo API.
func (c *Client) ListAcmeAccountDomain(ctx context.Context, params ListAcmeAccountDomainParams) (*ListAcmeAccountDomainResponse, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/acme/v2/account/%d/domain", c.BaseURL, params.AccountID))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("size", fmt.Sprintf("%d", params.Size))
	queryParams.Add("position", fmt.Sprintf("%d", params.Position))

	if params.Name != "" {
		queryParams.Add("name", params.Name)
	}
	if params.ExpiresWithinNextDays > 0 {
		queryParams.Add("expiresWithinNextDays", fmt.Sprintf("%d", params.ExpiresWithinNextDays))
	}
	if params.StickyExpiresWithinNextDays != 0 {
		queryParams.Add("stickyExpiresWithinNextDays", fmt.Sprintf("%d", params.StickyExpiresWithinNextDays))
	}

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, body, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var domains []AcmeAccountDomain
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	listAcmeAccountDomainResponse := ListAcmeAccountDomainResponse{Domains: domains}
	totalCountHeader := resp.Header.Get("X-Total-Count")
	if totalCountHeader != "" {
		listAcmeAccountDomainResponse.TotalCount, _ = strconv.Atoi(totalCountHeader)
	}

	return &listAcmeAccountDomainResponse, nil
}

// ListAllAcmeAccountDomain sends requests to list all ACME account domains by iterating through the results using the X-Total-Count header.
func (c *Client) ListAllAcmeAccountDomain(ctx context.Context, params ListAcmeAccountDomainParams) ([]AcmeAccountDomain, error) {
	var allAcmeAccountDomains []AcmeAccountDomain
	position := 0
	size := 200

	for {
		params.Position = position
		params.Size = size
		listAcmeAccountDomainResponse, err := c.ListAcmeAccountDomain(ctx, params)
		if err != nil {
			return nil, err
		}

		allAcmeAccountDomains = append(allAcmeAccountDomains, listAcmeAccountDomainResponse.Domains...)

		if len(listAcmeAccountDomainResponse.Domains) < params.Size || position+params.Size >= listAcmeAccountDomainResponse.TotalCount {
			break
		}

		position += params.Size
	}

	return allAcmeAccountDomains, nil
}

// AddAcmeAccountDomains sends a request to add domains to an ACME account via the Sectigo API.
func (c *Client) AddAcmeAccountDomains(ctx context.Context, params AcmeAccountDomainParams) error {
	url := fmt.Sprintf("%s/api/acme/v2/account/%d/domain", c.BaseURL, params.AccountID)

	var domains AcmeAccountDomainRequest
	for _, domain := range params.Domains {
		domains.Domains = append(domains.Domains, AcmeAccountDomainName{Name: domain})
	}

	jsonPayload, err := json.Marshal(domains)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add domain to acme account, status code: %d", resp.StatusCode)
	}

	return nil
}

type ListSSLParams struct {
	Size                   int
	Position               int
	CommonName             string
	SubjectAlternativeName string
	Status                 string
	SSLTypeId              int
	DiscoveryStatus        string
	Vendor                 string
	OrgId                  int
	InstallStatus          string
	RenewalStatus          string
	Issuer                 string
	SerialNumber           string
	Requester              string
	ExternalRequester      string
	SignatureAlgorithm     string
	KeyAlgorithm           string
	KeySize                int
	KeyParam               string
	Sha1Hash               string
	Md5Hash                string
	KeyUsage               string
	ExtendedKeyUsage       string
	RequestedVia           string
}

type SSLCertificate struct {
	SSLId                   int      `json:"sslId"`
	CommonName              string   `json:"commonName"`
	SubjectAlternativeNames []string `json:"subjectAlternativeNames"`
	SerialNumber            string   `json:"serialNumber"`
}

type ListSSLResponse struct {
	SSLCertificates []SSLCertificate
	TotalCount      int
}

func (c *Client) ListSSL(ctx context.Context, params ListSSLParams) (*ListSSLResponse, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/ssl/v1", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("size", fmt.Sprintf("%d", params.Size))
	queryParams.Add("position", fmt.Sprintf("%d", params.Position))
	if params.CommonName != "" {
		queryParams.Add("commonName", params.CommonName)
	}
	if params.SubjectAlternativeName != "" {
		queryParams.Add("subjectAlternativeName", params.SubjectAlternativeName)
	}
	if params.Status != "" {
		queryParams.Add("status", params.Status)
	}
	if params.SSLTypeId > 0 {
		queryParams.Add("sslTypeId", fmt.Sprintf("%d", params.SSLTypeId))
	}
	if params.DiscoveryStatus != "" {
		queryParams.Add("discoveryStatus", params.DiscoveryStatus)
	}
	if params.Vendor != "" {
		queryParams.Add("vendor", params.Vendor)
	}
	if params.OrgId > 0 {
		queryParams.Add("orgId", fmt.Sprintf("%d", params.OrgId))
	}
	if params.InstallStatus != "" {
		queryParams.Add("installStatus", params.InstallStatus)
	}
	if params.RenewalStatus != "" {
		queryParams.Add("renewalStatus", params.RenewalStatus)
	}
	if params.Issuer != "" {
		queryParams.Add("issuer", params.Issuer)
	}
	if params.SerialNumber != "" {
		queryParams.Add("serialNumber", params.SerialNumber)
	}
	if params.Requester != "" {
		queryParams.Add("requester", params.Requester)
	}
	if params.ExternalRequester != "" {
		queryParams.Add("externalRequester", params.ExternalRequester)
	}
	if params.SignatureAlgorithm != "" {
		queryParams.Add("signatureAlgorithm", params.SignatureAlgorithm)
	}
	if params.KeyAlgorithm != "" {
		queryParams.Add("keyAlgorithm", params.KeyAlgorithm)
	}
	if params.KeySize > 0 {
		queryParams.Add("keySize", fmt.Sprintf("%d", params.KeySize))
	}
	if params.KeyParam != "" {
		queryParams.Add("keyParam", params.KeyParam)
	}
	if params.Sha1Hash != "" {
		queryParams.Add("sha1Hash", params.Sha1Hash)
	}
	if params.Md5Hash != "" {
		queryParams.Add("md5Hash", params.Md5Hash)
	}
	if params.KeyUsage != "" {
		queryParams.Add("keyUsage", params.KeyUsage)
	}
	if params.ExtendedKeyUsage != "" {
		queryParams.Add("extendedKeyUsage", params.ExtendedKeyUsage)
	}
	if params.RequestedVia != "" {
		queryParams.Add("requestedVia", params.RequestedVia)
	}
	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, body, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var sslCertificates []SSLCertificate
	err = json.Unmarshal(body, &sslCertificates)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	listSSLResponse := ListSSLResponse{SSLCertificates: sslCertificates}
	totalCountHeader := resp.Header.Get("X-Total-Count")
	if totalCountHeader != "" {
		listSSLResponse.TotalCount, _ = strconv.Atoi(totalCountHeader)
	}

	return &listSSLResponse, nil
}

func (c *Client) ListAllSSL(ctx context.Context, params ListSSLParams) ([]SSLCertificate, error) {
	var allSSLCertificates []SSLCertificate
	position := 0
	size := 200

	for {
		params.Position = position
		params.Size = size
		listSSLResponse, err := c.ListSSL(ctx, params)
		if err != nil {
			return nil, err
		}

		allSSLCertificates = append(allSSLCertificates, listSSLResponse.SSLCertificates...)

		if len(listSSLResponse.SSLCertificates) < params.Size || position+params.Size >= listSSLResponse.TotalCount {
			break
		}

		position += params.Size
	}

	return allSSLCertificates, nil
}

type RevokeSSLParams struct {
	SSLId  int    `json:"sslId"`
	Reason string `json:"reason"`
}

func (c *Client) RevokeSSLById(ctx context.Context, sslId int, reason string) error {
	if reason == "" || len(reason) > 512 {
		return fmt.Errorf("reason must be between 1 and 512 characters")
	}

	url := fmt.Sprintf("%s/api/ssl/v1/revoke/%d", c.BaseURL, sslId)
	reqBody := map[string]string{"reason": reason}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to revoke SSL certificate: %s", resp.Status)
	}

	return nil
}

// SSLDetails represents the detailed information about an SSL certificate
type SSLDetails struct {
	CommonName              string             `json:"commonName"`
	SSLId                   int                `json:"sslId"`
	Id                      int                `json:"id"`
	OrgId                   int                `json:"orgId"`
	Status                  string             `json:"status"`
	OrderNumber             int                `json:"orderNumber"`
	BackendCertId           string             `json:"backendCertId"`
	Vendor                  string             `json:"vendor"`
	CertType                CertType           `json:"certType"`
	SubType                 string             `json:"subType"`
	ValidationType          string             `json:"validationType"`
	Term                    int                `json:"term"`
	Owner                   string             `json:"owner"`
	OwnerId                 int                `json:"ownerId"`
	Requester               string             `json:"requester"`
	RequesterId             int                `json:"requesterId"`
	RequestedVia            string             `json:"requestedVia"`
	ExternalRequester       string             `json:"externalRequester"`
	Comments                string             `json:"comments"`
	Requested               string             `json:"requested"`
	Approved                string             `json:"approved"`
	Issued                  string             `json:"issued"`
	Declined                string             `json:"declined"`
	Expires                 string             `json:"expires"`
	Replaced                string             `json:"replaced"`
	Revoked                 string             `json:"revoked"`
	ReasonCode              int                `json:"reasonCode"`
	Renewed                 bool               `json:"renewed"`
	RenewedDate             string             `json:"renewedDate"`
	SerialNumber            string             `json:"serialNumber"`
	SignatureAlg            string             `json:"signatureAlg"`
	KeyAlgorithm            string             `json:"keyAlgorithm"`
	KeySize                 int                `json:"keySize"`
	KeyType                 string             `json:"keyType"`
	KeyUsages               []string           `json:"keyUsages"`
	ExtendedKeyUsages       []string           `json:"extendedKeyUsages"`
	SubjectAlternativeNames []string           `json:"subjectAlternativeNames"`
	CustomFields            []CustomField      `json:"customFields"`
	CertificateDetails      CertificateDetails `json:"certificateDetails"`
	AutoInstallDetails      AutoInstallDetails `json:"autoInstallDetails"`
	AutoRenewDetails        AutoRenewDetails   `json:"autoRenewDetails"`
	SuspendNotifications    bool               `json:"suspendNotifications"`
}

// CertType represents information about the Certificate Profile
type CertType struct {
	Id                  int      `json:"id"`
	UseSecondaryOrgName bool     `json:"useSecondaryOrgName"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Terms               []int    `json:"terms"`
	KeyTypes            KeyTypes `json:"keyTypes"`
}

// KeyTypes represents key types available for the Certificate Profile
type KeyTypes struct {
	RSA []string `json:"rsa"`
	EC  []string `json:"ec"`
}

// CustomField represents a custom field
type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CertificateDetails represents certificate details
type CertificateDetails struct {
	Issuer          string `json:"issuer"`
	Subject         string `json:"subject"`
	SubjectAltNames string `json:"subjectAltNames"`
	Md5Hash         string `json:"md5Hash"`
	Sha1Hash        string `json:"sha1Hash"`
}

// AutoInstallDetails represents auto-installation information
type AutoInstallDetails struct {
	State string     `json:"state"`
	Nodes []NodeInfo `json:"nodes"`
}

// NodeInfo represents information about an auto-installation node
type NodeInfo struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// AutoRenewDetails represents auto-renewal information
type AutoRenewDetails struct {
	State                string `json:"state"`
	DaysBeforeExpiration int    `json:"daysBeforeExpiration"`
}

// MarshalJSON implements json.Marshaler for AutoRenewDetails
func (a AutoRenewDetails) MarshalJSON() ([]byte, error) {
	// If both State and DaysBeforeExpiration are empty, return null
	if a.State == "" && a.DaysBeforeExpiration == 0 {
		return []byte("null"), nil
	}

	// Otherwise, use the default marshaling
	type Alias AutoRenewDetails
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&a),
	})
}

// GetSSLDetails retrieves detailed information about an SSL certificate
func (c *Client) GetSSLDetails(ctx context.Context, sslId int) (*SSLDetails, error) {
	// Construct the URL
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/ssl/v1/%d", c.BaseURL, sslId))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// Send the request
	resp, body, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get SSL details: %s", resp.Status)
	}

	// Parse the response
	var sslDetails SSLDetails
	err = json.Unmarshal(body, &sslDetails)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &sslDetails, nil
}

// UpdateSSLDetailsRequest represents the request body for updating SSL certificate details
type UpdateSSLDetailsRequest struct {
	SSLId                   int              `json:"sslId"`
	Term                    int              `json:"term,omitempty"`
	CertTypeId              int              `json:"certTypeId,omitempty"`
	OrgId                   int              `json:"orgId,omitempty"`
	CommonName              string           `json:"commonName,omitempty"`
	CSR                     string           `json:"csr,omitempty"`
	ExternalRequester       string           `json:"externalRequester,omitempty"`
	Comments                string           `json:"comments,omitempty"`
	SubjectAlternativeNames []string         `json:"subjectAlternativeNames,omitempty"`
	CustomFields            []CustomField    `json:"customFields,omitempty"`
	AutoRenewDetails        *AutoRenewDetails `json:"autoRenewDetails,omitempty"`
	SuspendNotifications    bool             `json:"suspendNotifications,omitempty"`
	Requester               string           `json:"requester,omitempty"`
	RequesterAdminId        int              `json:"requesterAdminId,omitempty"`
	ApproverAdminId         int              `json:"approverAdminId,omitempty"`
}

// validateUpdateSSLDetailsRequest validates the request parameters
func validateUpdateSSLDetailsRequest(request UpdateSSLDetailsRequest) error {
	// Validate required fields
	if request.SSLId < 1 {
		return fmt.Errorf("sslId must be at least 1")
	}

	// Validate term if provided
	if request.Term != 0 && request.Term < 1 {
		return fmt.Errorf("term must be at least 1")
	}

	// Validate certTypeId if provided
	if request.CertTypeId != 0 && request.CertTypeId < 1 {
		return fmt.Errorf("certTypeId must be at least 1")
	}

	// Validate orgId if provided
	if request.OrgId != 0 && request.OrgId < 1 {
		return fmt.Errorf("orgId must be at least 1")
	}

	// Validate CSR if provided
	if request.CSR != "" {
		csrRegex := regexp.MustCompile(`^[a-zA-Z0-9-+=\/\s]+$`)
		if !csrRegex.MatchString(request.CSR) {
			return fmt.Errorf("csr must match the regular expression [a-zA-Z0-9-+=\\/\\s]+")
		}
		if len(request.CSR) > 32767 {
			return fmt.Errorf("csr size must be between 1 and 32767 inclusive")
		}
	}

	// Validate comments if provided
	if len(request.Comments) > 1024 {
		return fmt.Errorf("comments maximum length is 1024 characters or can be empty")
	}

	// Validate custom fields if provided
	for _, field := range request.CustomFields {
		if field.Name == "" {
			return fmt.Errorf("custom field name must not be null")
		}
		if len(field.Name) < 1 || len(field.Name) > 256 {
			return fmt.Errorf("custom field name size must be between 1 and 256 inclusive")
		}
		if len(field.Value) > 256 {
			return fmt.Errorf("custom field value maximum length is 256 characters or can be empty")
		}
	}

	// Validate autoRenewDetails if provided
	if request.AutoRenewDetails != nil {
		if request.AutoRenewDetails.State != "" {
			if request.AutoRenewDetails.State != "Not scheduled" && request.AutoRenewDetails.State != "Scheduled" {
				return fmt.Errorf("autoRenewDetails.state allowed values are 'Not scheduled' and 'Scheduled'")
			}
		}
		if request.AutoRenewDetails.DaysBeforeExpiration != 0 && request.AutoRenewDetails.DaysBeforeExpiration < 1 {
			return fmt.Errorf("autoRenewDetails.daysBeforeExpiration must be at least 1")
		}
	}

	// Validate requesterAdminId if provided
	if request.RequesterAdminId != 0 && request.RequesterAdminId < 1 {
		return fmt.Errorf("requesterAdminId must be at least 1")
	}

	// Validate approverAdminId if provided
	if request.ApproverAdminId != 0 && request.ApproverAdminId < -1 {
		return fmt.Errorf("approverAdminId must be at least -1")
	}

	return nil
}

// UpdateSSLDetails updates the details of an SSL certificate
func (c *Client) UpdateSSLDetails(ctx context.Context, request UpdateSSLDetailsRequest) (*SSLDetails, error) {
	// First validate the request
	if err := validateUpdateSSLDetailsRequest(request); err != nil {
		return nil, err
	}

	// Construct the URL
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/ssl/v1", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	// Marshal the request body
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %v", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "PUT", baseURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, body, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed request, status code: %d", resp.StatusCode)
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update SSL details: %s", resp.Status)
	}

	// Parse the response
	var sslDetails SSLDetails
	err = json.Unmarshal(body, &sslDetails)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &sslDetails, nil
}
