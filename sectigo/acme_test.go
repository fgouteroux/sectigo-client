package sectigo

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAcmeAccount(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccount{
			{
				ID:             1,
				Name:           "Account 1",
				OrganizationID: 1,
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
		Size:           10,
		OrganizationId: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(accounts.Accounts))
	assert.Equal(t, "Account 1", accounts.Accounts[0].Name)
}

func TestListAcmeAccount_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Server error"}`))
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
	_, err := client.ListAcmeAccount(ctx, ListAcmeAccountParams{OrganizationId: 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestListAllAcmeAccount(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "2")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccount{
			{ID: 1, Name: "Account 1", OrganizationID: 1},
			{ID: 2, Name: "Account 2", OrganizationID: 1},
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
	accounts, err := client.ListAllAcmeAccount(ctx, ListAcmeAccountParams{OrganizationId: 1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(accounts))
	assert.Equal(t, "Account 1", accounts[0].Name)
	assert.Equal(t, "Account 2", accounts[1].Name)
}

func TestListAcmeAccountDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccountDomain{
			{
				Name:       "example.com",
				ValidUntil: "2023-12-31",
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
	domains, err := client.ListAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1, Size: 10})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(domains.Domains))
	assert.Equal(t, "example.com", domains.Domains[0].Name)
}

func TestListAcmeAccountDomain_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Server error"}`))
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
	_, err := client.ListAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1, Size: 10})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestListAllAcmeAccountDomain(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "2")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]AcmeAccountDomain{
			{Name: "example.com", ValidUntil: "2023-12-31"},
			{Name: "example2.com", ValidUntil: "2023-12-31"},
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
	domains, err := client.ListAllAcmeAccountDomain(ctx, ListAcmeAccountDomainParams{AccountID: 1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(domains))
	assert.Equal(t, "example.com", domains[0].Name)
	assert.Equal(t, "example2.com", domains[1].Name)
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
	err := client.AddAcmeAccountDomains(ctx, AcmeAccountDomainParams{
		AccountID: 1,
		Domains:   []string{"example1.com", "example2.com"},
	})
	assert.NoError(t, err)
}

func TestAddAcmeAccountDomains_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/acme/v2/account/1/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code":-993,"description":"Certificate orders currently restricted"}`))
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
	err := client.AddAcmeAccountDomains(ctx, AcmeAccountDomainParams{
		AccountID: 1,
		Domains:   []string{"example1.com"},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "Certificate orders currently restricted")
}