package sectigo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
