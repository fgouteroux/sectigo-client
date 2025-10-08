package sectigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

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

// ListAcmeAccountParams represents query parameters for listing ACME accounts.
type ListAcmeAccountParams struct {
	Position           int
	Size               int
	OrganizationId     int
	Name               string
	AcmeServer         string
	CertValidationType string
	Status             string
}

// ListAcmeAccountDomainParams represents query parameters for listing ACME account domains.
type ListAcmeAccountDomainParams struct {
	AccountID                   int
	Position                    int
	Size                        int
	Name                        string
	ExpiresWithinNextDays       int
	StickyExpiresWithinNextDays int
}

// ListAcmeAccountDomainResponse represents the response structure for listing ACME account domains.
type ListAcmeAccountDomainResponse struct {
	Domains    []AcmeAccountDomain `json:"domains"`
	TotalCount int                 `json:"total_count"`
}

// AcmeAccountDomain represents an ACME account domain.
type AcmeAccountDomain struct {
	Name                string `json:"name"`
	ValidUntil          string `json:"validUntil"`
	StickyUntil         string `json:"stickyUntil"`
	OvAnchorOrderNumber int    `json:"ovAnchorOrderNumber"`
	OvAnchorID          string `json:"ovAnchorID"`
	EvAnchorOrderNumber int    `json:"evAnchorOrderNumber"`
	EvAnchorID          string `json:"evAnchorID"`
}

// AcmeAccountDomainParams represents parameters for adding domains to an ACME account.
type AcmeAccountDomainParams struct {
	AccountID int
	Domains   []string
}

// AcmeAccountDomainName represents a domain name in ACME account operations.
type AcmeAccountDomainName struct {
	Name string `json:"name"`
}

// AcmeAccountDomainRequest represents the request structure for adding domains to an ACME account.
type AcmeAccountDomainRequest struct {
	Domains []AcmeAccountDomainName `json:"domains"`
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

	resp, body, err := c.sendRequest(ctx, req, http.StatusOK)
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

	resp, body, err := c.sendRequest(ctx, req, http.StatusOK)
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

	_, _, err = c.sendRequest(ctx, req, http.StatusOK)
	return err
}