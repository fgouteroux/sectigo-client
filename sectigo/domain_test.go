package sectigo

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
		w.Write([]byte(`{"error":"Domain not found"}`))
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
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "Domain not found")
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
		w.Write([]byte(`{"error":"Invalid domain name"}`))
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
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "Invalid domain name")
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
		w.Write([]byte(`{"error":"Domain not found"}`))
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
	assert.Contains(t, err.Error(), "404")
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
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(domains.Domains))
	assert.Equal(t, 1, domains.TotalCount)
	assert.Equal(t, "example.com", domains.Domains[0].Name)
}

func TestListAllDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/domain/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
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
	domains, err := client.ListAllDomain(ctx, ListDomainParams{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(domains))
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
}

func TestCheckDomainValidationStatus(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(GetDomainValidationStatusResponse{
			Status: "validated",
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
	err := client.CheckDomainValidationStatus(ctx, "example.com", 3, 1*time.Millisecond)
	assert.NoError(t, err)
}

func TestCheckDomainValidationStatus_Retry(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	attempt := 0
	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		if attempt < 2 {
			_ = json.NewEncoder(w).Encode(GetDomainValidationStatusResponse{
				Status: "NOT_VALIDATED",
			})
			attempt++
		} else {
			_ = json.NewEncoder(w).Encode(GetDomainValidationStatusResponse{
				Status: "validated",
			})
		}
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
	err := client.CheckDomainValidationStatus(ctx, "example.com", 5, 1*time.Millisecond)
	assert.NoError(t, err)
}

func TestCheckDomainValidationStatus_MaxRetriesReached(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v2/validation/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(GetDomainValidationStatusResponse{
			Status: "NOT_VALIDATED",
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
	err := client.CheckDomainValidationStatus(ctx, "example.com", 3, 1*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max retries reached")
}

func TestListDomainValidation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]DomainValidation{
			{Domain: "example.com", DcvStatus: "validated"},
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
	validations, err := client.ListDomainValidation(ctx, ListDomainValidationParams{Size: 10})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(validations.Domains))
}

func TestListAllDomainValidation(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/dcv/v1/validation", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]DomainValidation{
			{Domain: "example.com", DcvStatus: "validated"},
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
	validations, err := client.ListAllDomainValidation(ctx, ListDomainValidationParams{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(validations))
}