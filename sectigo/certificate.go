package sectigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// ListSSLParams represents the parameters for listing SSL certificates.
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

// SSLCertificate represents an SSL certificate.
type SSLCertificate struct {
	SSLId                   int      `json:"sslId"`
	CommonName              string   `json:"commonName"`
	SubjectAlternativeNames []string `json:"subjectAlternativeNames"`
	SerialNumber            string   `json:"serialNumber"`
}

// ListSSLResponse represents the response structure for listing SSL certificates.
type ListSSLResponse struct {
	SSLCertificates []SSLCertificate
	TotalCount      int
}

// RevokeSSLParams represents the parameters for revoking an SSL certificate.
type RevokeSSLParams struct {
	SSLId  int    `json:"sslId"`
	Reason string `json:"reason"`
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

// UpdateSSLDetailsRequest represents the request body for updating SSL certificate details
type UpdateSSLDetailsRequest struct {
	SSLId                   int               `json:"sslId"`
	Term                    int               `json:"term,omitempty"`
	CertTypeId              int               `json:"certTypeId,omitempty"`
	OrgId                   int               `json:"orgId,omitempty"`
	CommonName              string            `json:"commonName,omitempty"`
	CSR                     string            `json:"csr,omitempty"`
	ExternalRequester       string            `json:"externalRequester,omitempty"`
	Comments                string            `json:"comments,omitempty"`
	SubjectAlternativeNames []string          `json:"subjectAlternativeNames,omitempty"`
	CustomFields            []CustomField     `json:"customFields,omitempty"`
	AutoRenewDetails        *AutoRenewDetails `json:"autoRenewDetails,omitempty"`
	SuspendNotifications    bool              `json:"suspendNotifications,omitempty"`
	Requester               string            `json:"requester,omitempty"`
	RequesterAdminId        int               `json:"requesterAdminId,omitempty"`
	ApproverAdminId         int               `json:"approverAdminId,omitempty"`
}

// ListSSL sends a request to list SSL certificates via the Sectigo API.
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

	resp, body, err := c.sendRequest(ctx, req, http.StatusOK)
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

// ListAllSSL sends requests to list all SSL certificates by iterating through the results using the X-Total-Count header.
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

// RevokeSSLById sends a request to revoke an SSL certificate by ID via the Sectigo API.
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

	_, _, err = c.sendRequest(ctx, req, http.StatusNoContent)
	return err
}

// GetSSLDetails retrieves detailed information about an SSL certificate
func (c *Client) GetSSLDetails(ctx context.Context, sslId int) (*SSLDetails, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/api/ssl/v1/%d", c.BaseURL, sslId))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var sslDetails SSLDetails
	err = json.Unmarshal(body, &sslDetails)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &sslDetails, nil
}

// validateUpdateSSLDetailsRequest validates the request parameters
func validateUpdateSSLDetailsRequest(request UpdateSSLDetailsRequest) error {
	if request.SSLId < 1 {
		return fmt.Errorf("sslId must be at least 1")
	}

	if request.Term != 0 && request.Term < 1 {
		return fmt.Errorf("term must be at least 1")
	}

	if request.CertTypeId != 0 && request.CertTypeId < 1 {
		return fmt.Errorf("certTypeId must be at least 1")
	}

	if request.OrgId != 0 && request.OrgId < 1 {
		return fmt.Errorf("orgId must be at least 1")
	}

	if request.CSR != "" {
		csrRegex := regexp.MustCompile(`^[a-zA-Z0-9-+=\/\s]+$`)
		if !csrRegex.MatchString(request.CSR) {
			return fmt.Errorf("csr must match the regular expression [a-zA-Z0-9-+=\\/\\s]+")
		}
		if len(request.CSR) > 32767 {
			return fmt.Errorf("csr size must be between 1 and 32767 inclusive")
		}
	}

	if len(request.Comments) > 1024 {
		return fmt.Errorf("comments maximum length is 1024 characters or can be empty")
	}

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

	if request.RequesterAdminId != 0 && request.RequesterAdminId < 1 {
		return fmt.Errorf("requesterAdminId must be at least 1")
	}

	if request.ApproverAdminId != 0 && request.ApproverAdminId < -1 {
		return fmt.Errorf("approverAdminId must be at least -1")
	}

	return nil
}

// UpdateSSLDetails updates the details of an SSL certificate
func (c *Client) UpdateSSLDetails(ctx context.Context, request UpdateSSLDetailsRequest) (*SSLDetails, error) {
	if err := validateUpdateSSLDetailsRequest(request); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(fmt.Sprintf("%s/api/ssl/v1", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error parsing base URL: %v", err)
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", baseURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	_, body, err := c.sendRequest(ctx, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var sslDetails SSLDetails
	err = json.Unmarshal(body, &sslDetails)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &sslDetails, nil
}
