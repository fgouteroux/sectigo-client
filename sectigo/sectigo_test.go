package sectigo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockClient is a mock HTTP client that returns predefined responses.
type MockClient struct {
	Client  *http.Client
	Mux     *http.ServeMux
	Server  *httptest.Server
	Handler http.Handler
}

// NewMockClient creates a new mock HTTP client.
func NewMockClient() *MockClient {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return &MockClient{
		Client:  server.Client(),
		Mux:     mux,
		Server:  server,
		Handler: mux,
	}
}

// Close shuts down the mock HTTP server.
func (m *MockClient) Close() {
	m.Server.Close()
}

func TestGetDomainDetails(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DomainDetails{
			ID:               1,
			Name:             "example.com",
			DelegationStatus: "delegated",
			State:            "active",
			ValidationStatus: "validated",
			ValidationMethod: "http",
			DcvValidation:    "valid",
			DcvExpiration:    "2023-12-31",
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	domainDetails, err := client.GetDomainDetails(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, domainDetails.ID)
	assert.Equal(t, "example.com", domainDetails.Name)
}

func TestGetDomainDetails_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusNotFound)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.GetDomainDetails(ctx, 1)
	assert.Error(t, err)
}

func TestCreateDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusCreated)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.CreateDomain(ctx, DomainRequest{
		Name:        "example.com",
		Description: "Test domain",
		Active:      true,
	})
	assert.NoError(t, err)
}

func TestCreateDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.CreateDomain(ctx, DomainRequest{
		Name:        "example.com",
		Description: "Test domain",
		Active:      true,
	})
	assert.Error(t, err)
}

func TestDeleteDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.DeleteDomain(ctx, 1)
	assert.NoError(t, err)
}

func TestDeleteDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNotFound)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.DeleteDomain(ctx, 1)
	assert.Error(t, err)
}

func TestApproveDelegation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/1/delegation/approve", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.ApproveDelegation(ctx, 1, ApproveDelegationRequest{OrgId: 1})
	assert.NoError(t, err)
}

func TestApproveDelegation_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/1/delegation/approve", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.ApproveDelegation(ctx, 1, ApproveDelegationRequest{OrgId: 1})
	assert.Error(t, err)
}

func TestDelegateDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/delegation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.DelegateDomain(ctx, DelegateDomainRequest{
		DomainIds: []int{1},
		OrgId:     1,
		CertTypes: []string{"type1"},
	})
	assert.NoError(t, err)
}

func TestDelegateDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1/delegation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.DelegateDomain(ctx, DelegateDomainRequest{
		DomainIds: []int{1},
		OrgId:     1,
		CertTypes: []string{"type1"},
	})
	assert.Error(t, err)
}

func TestListDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]Domain{
			{ID: 1, Name: "example.com"},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	domains, err := client.ListDomain(ctx, ListDomainParams{
		Size:     10,
		Position: 0,
		Name:     "example.com",
		State:    "active",
		Status:   "validated",
		OrgId:    1,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(domains.Domains))
	assert.Equal(t, 1, domains.TotalCount)
	assert.Equal(t, "example.com", domains.Domains[0].Name)
}

func TestListDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListDomain(ctx, ListDomainParams{Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestListAllDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]Domain{
			{ID: 1, Name: "example.com"},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	domains, err := client.ListAllDomain(ctx, ListDomainParams{Size: 10, Position: 0})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(domains))
	assert.Equal(t, "example.com", domains[0].Name)
}

func TestListAllDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAllDomain(ctx, ListDomainParams{Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestStartDomainCNameValidation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation/start/domain/cname", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(StartDomainCNameValidationResponse{
			Host:  "host",
			Point: "point",
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	response, err := client.StartDomainCNameValidation(ctx, StartDomainCNameValidationRequest{Domain: "example.com"})
	assert.NoError(t, err)
	assert.Equal(t, "host", response.Host)
	assert.Equal(t, "point", response.Point)
}

func TestStartDomainCNameValidation_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation/start/domain/cname", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.StartDomainCNameValidation(ctx, StartDomainCNameValidationRequest{Domain: "example.com"})
	assert.Error(t, err)
}

func TestSubmitDomainCNameValidation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation/submit/domain/cname", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SubmitDomainCNameValidationResponse{
			OrderStatus: "valid",
			Message:     "success",
			Status:      "success",
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	response, err := client.SubmitDomainCNameValidation(ctx, SubmitDomainCNameValidationRequest{Domain: "example.com"})
	assert.NoError(t, err)
	assert.Equal(t, "valid", response.OrderStatus)
	assert.Equal(t, "success", response.Message)
	assert.Equal(t, "success", response.Status)
}

func TestSubmitDomainCNameValidation_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation/submit/domain/cname", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.SubmitDomainCNameValidation(ctx, SubmitDomainCNameValidationRequest{Domain: "example.com"})
	assert.Error(t, err)
}

func TestGetDomainValidationStatus(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(GetDomainValidationStatusResponse{
			Status:         "validated",
			OrderStatus:    "valid",
			ExpirationDate: "2023-12-31",
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	response, err := client.GetDomainValidationStatus(ctx, GetDomainValidationStatusRequest{Domain: "example.com"})
	assert.NoError(t, err)
	assert.Equal(t, "validated", response.Status)
	assert.Equal(t, "valid", response.OrderStatus)
	assert.Equal(t, "2023-12-31", response.ExpirationDate)
}

func TestGetDomainValidationStatus_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.GetDomainValidationStatus(ctx, GetDomainValidationStatusRequest{Domain: "example.com"})
	assert.Error(t, err)
}

func TestCheckDomainValidationStatus(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(GetDomainValidationStatusResponse{
			Status:         "validated",
			OrderStatus:    "valid",
			ExpirationDate: "2023-12-31",
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.CheckDomainValidationStatus(ctx, "example.com", 3, 1*time.Second)
	assert.NoError(t, err)
}

func TestCheckDomainValidationStatus_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.CheckDomainValidationStatus(ctx, "example.com", 3, 1*time.Second)
	assert.Error(t, err)
}

func TestListDomainValidation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]DomainValidation{
			{
				Domain:         "example.com",
				DcvStatus:      "validated",
				DcvOrderStatus: "valid",
				DcvMethod:      "http",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	validations, err := client.ListDomainValidation(ctx, ListDomainValidationParams{
		Size:        10,
		Position:    0,
		Domain:      "example.com",
		Org:         1,
		Department:  1,
		DcvStatus:   "validated",
		OrderStatus: "valid",
		ExpiresIn:   90,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(validations.Domains))
	assert.Equal(t, 1, validations.TotalCount)
	assert.Equal(t, "example.com", validations.Domains[0].Domain)
}

func TestListDomainValidation_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListDomainValidation(ctx, ListDomainValidationParams{Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestListAllDomainValidation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]DomainValidation{
			{
				Domain:         "example.com",
				DcvStatus:      "validated",
				DcvOrderStatus: "valid",
				DcvMethod:      "http",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	validations, err := client.ListAllDomainValidation(ctx, ListDomainValidationParams{Size: 10, Position: 0})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(validations))
	assert.Equal(t, "example.com", validations[0].Domain)
}

func TestListAllDomainValidation_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAllDomainValidation(ctx, ListDomainValidationParams{Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestListOrganization(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/organization/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ListOrganizationResponse{
			{
				ID:   1,
				Name: "Test Organization",
				Departments: []Department{
					{ID: 1, Name: "Department 1", ParentName: "Parent 1"},
				},
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	organizations, err := client.ListOrganization(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*organizations))
	assert.Equal(t, "Test Organization", (*organizations)[0].Name)
}

func TestListOrganization_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/organization/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListOrganization(ctx)
	assert.Error(t, err)
}

func TestListAcmeAccount(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccount{
			{
				ID:                 1,
				AccountID:          "account1",
				MacID:              "mac1",
				MacKey:             "mackey1",
				AcmeServer:         "server1",
				Name:               "Account 1",
				OrganizationID:     1,
				CertValidationType: "type1",
				Status:             "active",
				OvOrderNumber:      1,
				OvAnchorID:         "anchor1",
				EvDetails:          struct{}{},
				Contacts:           "contact1",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	accounts, err := client.ListAcmeAccount(ctx, ListAcmeAccountParams{
		Size:               10,
		Position:           0,
		OrganizationId:     1,
		Name:               "Account 1",
		AcmeServer:         "server1",
		CertValidationType: "type1",
		Status:             "active",
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(accounts.Accounts))
	assert.Equal(t, 1, accounts.TotalCount)
	assert.Equal(t, "Account 1", accounts.Accounts[0].Name)
}

func TestListAcmeAccount_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAcmeAccount(ctx, ListAcmeAccountParams{Size: 10, Position: 0, OrganizationId: 1})
	assert.Error(t, err)
}

func TestListAllAcmeAccount(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccount{
			{
				ID:                 1,
				AccountID:          "account1",
				MacID:              "mac1",
				MacKey:             "mackey1",
				AcmeServer:         "server1",
				Name:               "Account 1",
				OrganizationID:     1,
				CertValidationType: "type1",
				Status:             "active",
				OvOrderNumber:      1,
				OvAnchorID:         "anchor1",
				EvDetails:          struct{}{},
				Contacts:           "contact1",
			},
			{
				ID:                 2,
				AccountID:          "account2",
				MacID:              "mac1",
				MacKey:             "mackey1",
				AcmeServer:         "server1",
				Name:               "Account 2",
				OrganizationID:     1,
				CertValidationType: "type1",
				Status:             "active",
				OvOrderNumber:      1,
				OvAnchorID:         "anchor1",
				EvDetails:          struct{}{},
				Contacts:           "contact1",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	accounts, err := client.ListAllAcmeAccount(ctx, ListAcmeAccountParams{Size: 10, Position: 0, OrganizationId: 1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(accounts))
	assert.Equal(t, "Account 1", accounts[0].Name)
	assert.Equal(t, "Account 2", accounts[1].Name)
}

func TestListAllAcmeAccount_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAllAcmeAccount(ctx, ListAcmeAccountParams{Size: 10, Position: 0, OrganizationId: 1})
	assert.Error(t, err)
}

func TestListAcmeAccountDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccountDomain{
			{
				Name:                "example.com",
				ValidUntil:          "2023-12-31",
				StickyUntil:         "2023-12-31",
				OvAnchorOrderNumber: 1,
				OvAnchorID:          "anchor1",
				EvAnchorOrderNumber: 1,
				EvAnchorID:          "anchor1",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	domains, err := client.ListAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1, Size: 10, Position: 0})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(domains.Domains))
	assert.Equal(t, 1, domains.TotalCount)
	assert.Equal(t, "example.com", domains.Domains[0].Name)
}

func TestListAcmeAccountDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1, Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestListAllAcmeAccountDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "2")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccountDomain{
			{
				Name:                "example.com",
				ValidUntil:          "2023-12-31",
				StickyUntil:         "2023-12-31",
				OvAnchorOrderNumber: 1,
				OvAnchorID:          "anchor1",
				EvAnchorOrderNumber: 1,
				EvAnchorID:          "anchor1",
			},
			{
				Name:                "example2.com",
				ValidUntil:          "2023-12-31",
				StickyUntil:         "2023-12-31",
				OvAnchorOrderNumber: 1,
				OvAnchorID:          "anchor1",
				EvAnchorOrderNumber: 1,
				EvAnchorID:          "anchor1",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	domains, err := client.ListAllAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1, Size: 10, Position: 0})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(domains))
	assert.Equal(t, "example.com", domains[0].Name)
	assert.Equal(t, "example2.com", domains[1].Name)
}

func TestListAllAcmeAccountDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAllAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1, Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestAddAcmeAccountDomains(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.AddAcmeAccountDomains(ctx, AcmeAccountDomainParams{AccountID: 1, Domains: []string{"example1.com", "example2.com"}})
	assert.NoError(t, err)
}

func TestAddAcmeAccountDomains_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.AddAcmeAccountDomains(ctx, AcmeAccountDomainParams{AccountID: 1, Domains: []string{"example1.com", "example2.com"}})
	assert.Error(t, err)
}

func TestRoundTrip(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"key": "value"})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})

	// Use the authTransport directly
	transport := client.Client.Transport.(*authTransport)
	req, _ := http.NewRequest("GET", mockClient.Server.URL+"/test", nil)
	resp, err := transport.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	_ = json.Unmarshal(body, &result)
	assert.Equal(t, "value", result["key"])
}

func TestListSSL(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]SSLCertificate{
			{
				SSLId:                   1,
				CommonName:              "example.com",
				SubjectAlternativeNames: []string{"example.com", "www.example.com"},
				SerialNumber:            "1234567890",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	sslCertificates, err := client.ListSSL(ctx, ListSSLParams{
		Size:       10,
		Position:   0,
		CommonName: "example.com",
		Status:     "validated",
		OrgId:      1,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sslCertificates.SSLCertificates))
	assert.Equal(t, 1, sslCertificates.TotalCount)
	assert.Equal(t, "example.com", sslCertificates.SSLCertificates[0].CommonName)
}

func TestListSSL_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListSSL(ctx, ListSSLParams{Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestListAllSSL(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]SSLCertificate{
			{
				SSLId:                   1,
				CommonName:              "example.com",
				SubjectAlternativeNames: []string{"example.com", "www.example.com"},
				SerialNumber:            "1234567890",
			},
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	sslCertificates, err := client.ListAllSSL(ctx, ListSSLParams{Size: 10, Position: 0})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sslCertificates))
	assert.Equal(t, "example.com", sslCertificates[0].CommonName)
}

func TestListAllSSL_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.ListAllSSL(ctx, ListSSLParams{Size: 10, Position: 0})
	assert.Error(t, err)
}

func TestRevokeSSLById(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1/revoke/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, "test reason", reqBody["reason"])

		w.WriteHeader(http.StatusNoContent)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.RevokeSSLById(ctx, 1, "test reason")
	assert.NoError(t, err)
}

func TestRevokeSSLById_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1/revoke/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	err := client.RevokeSSLById(ctx, 1, "test reason")
	assert.Error(t, err)
}

func TestRevokeSSLById_InvalidReason(t *testing.T) {
	client := NewClient(Config{
		URL:      "http://example.com",
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})

	ctx := context.Background()
	err := client.RevokeSSLById(ctx, 1, "")
	assert.Error(t, err)
	assert.Equal(t, "reason must be between 1 and 512 characters", err.Error())

	longReason := strings.Repeat("a", 513)
	err = client.RevokeSSLById(ctx, 1, longReason)
	assert.Error(t, err)
	assert.Equal(t, "reason must be between 1 and 512 characters", err.Error())
}

func TestGetSSLDetails(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1/1638", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SSLDetails{
			CommonName:   "example.com",
			SSLId:        1638,
			Id:           1638,
			Status:       "Issued",
			Vendor:       "Sectigo",
			RequestedVia: "ACME",
			SerialNumber: "1234567890",
			// Add other fields as needed
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	sslDetails, err := client.GetSSLDetails(ctx, 1638)
	assert.NoError(t, err)
	assert.NotNil(t, sslDetails)
	assert.Equal(t, "example.com", sslDetails.CommonName)
	assert.Equal(t, 1638, sslDetails.SSLId)
	assert.Equal(t, "Issued", sslDetails.Status)
}

func TestGetSSLDetails_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1/1638", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.GetSSLDetails(ctx, 1638)
	assert.Error(t, err)
}

func TestUpdateSSLDetails(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "application/json;charset=UTF-8", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		var request UpdateSSLDetailsRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)
		assert.Equal(t, 1740, request.SSLId)
		assert.Equal(t, 365, request.Term)
		assert.Equal(t, 5108, request.CertTypeId)
		assert.Equal(t, 10548, request.OrgId)
		assert.Equal(t, "ccmqa.com", request.CommonName)
		assert.Equal(t, "some comments", request.Comments)
		assert.Equal(t, []string{"ccmqa.com"}, request.SubjectAlternativeNames)
		assert.Equal(t, "Not scheduled", request.AutoRenewDetails.State)
		assert.Equal(t, 30, request.AutoRenewDetails.DaysBeforeExpiration)

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SSLDetails{
			CommonName:   "ccmqa.com",
			SSLId:        1740,
			Id:           1740,
			Status:       "Requested",
			Vendor:       "Vendor",
			RequestedVia: "Enrollment Form",
			Comments:     "some comments",
			// Add other fields as needed
		})
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	sslDetails, err := client.UpdateSSLDetails(ctx, UpdateSSLDetailsRequest{
		SSLId:                   1740,
		Term:                    365,
		RequesterAdminId:        1,
		ApproverAdminId:         -1,
		CertTypeId:              5108,
		OrgId:                   10548,
		CommonName:              "ccmqa.com",
		Comments:                "some comments",
		SubjectAlternativeNames: []string{"ccmqa.com"},
		AutoRenewDetails: AutoRenewDetails{
			State:                "Not scheduled",
			DaysBeforeExpiration: 30,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, sslDetails)
	assert.Equal(t, "ccmqa.com", sslDetails.CommonName)
	assert.Equal(t, 1740, sslDetails.SSLId)
	assert.Equal(t, "Requested", sslDetails.Status)
}

func TestUpdateSSLDetails_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})
	client.Client = mockClient.Client

	ctx := context.Background()
	_, err := client.UpdateSSLDetails(ctx, UpdateSSLDetailsRequest{
		SSLId: 1740,
	})
	assert.Error(t, err)
}

func TestValidateUpdateSSLDetailsRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateSSLDetailsRequest
		expectedErr string
	}{
		{
			name: "valid request",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
			},
			expectedErr: "",
		},
		{
			name: "invalid SSLId",
			request: UpdateSSLDetailsRequest{
				SSLId:            0,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
			},
			expectedErr: "sslId must be at least 1",
		},
		{
			name: "invalid term",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             0,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
			},
			expectedErr: "term must be at least 1",
		},
		{
			name: "invalid certTypeId",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       0,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
			},
			expectedErr: "certTypeId must be at least 1",
		},
		{
			name: "invalid orgId",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            0,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
			},
			expectedErr: "orgId must be at least 1",
		},
		{
			name: "invalid CSR",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
				CSR:              "invalid_csr!",
			},
			expectedErr: "csr must match the regular expression [a-zA-Z0-9-+=\\/\\s]+",
		},
		{
			name: "invalid comments",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
				Comments:         strings.Repeat("a", 1025),
			},
			expectedErr: "comments maximum length is 1024 characters or can be empty",
		},
		{
			name: "invalid custom field name",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
				CustomFields: []CustomField{
					{
						Name:  "",
						Value: "value",
					},
				},
			},
			expectedErr: "custom field name must not be null",
		},
		{
			name: "invalid autoRenewDetails state",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -1,
				AutoRenewDetails: AutoRenewDetails{
					State: "Invalid",
				},
			},
			expectedErr: "autoRenewDetails.state allowed values are 'Not scheduled' and 'Scheduled'",
		},
		{
			name: "invalid requesterAdminId",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 0,
				ApproverAdminId:  -1,
			},
			expectedErr: "requesterAdminId must be at least 1",
		},
		{
			name: "invalid approverAdminId",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				OrgId:            1,
				Term:             1,
				CertTypeId:       1,
				RequesterAdminId: 1,
				ApproverAdminId:  -2,
			},
			expectedErr: "approverAdminId must be at least -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateSSLDetailsRequest(tt.request)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			}
		})
	}
}
