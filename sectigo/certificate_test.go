package sectigo

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListSSL(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
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
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sslCertificates.SSLCertificates))
	assert.Equal(t, "example.com", sslCertificates.SSLCertificates[0].CommonName)
}

func TestListSSL_Error(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Server error"}`)) //nolint:errcheck
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
	assert.Contains(t, err.Error(), "500")
}

func TestListAllSSL(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/api/ssl/v1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Header().Set("X-Total-Count", "1")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]SSLCertificate{
			{
				SSLId:                   1,
				CommonName:              "example.com",
				SubjectAlternativeNames: []string{"example.com"},
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
	sslCertificates, err := client.ListAllSSL(ctx, ListSSLParams{Size: 10})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sslCertificates))
	assert.Equal(t, "example.com", sslCertificates[0].CommonName)
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
		w.Write([]byte(`{"error":"Failed to revoke"}`)) //nolint:errcheck
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
	assert.Contains(t, err.Error(), "500")
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

	// Test empty reason
	err := client.RevokeSSLById(ctx, 1, "")
	assert.Error(t, err)
	assert.Equal(t, "reason must be between 1 and 512 characters", err.Error())

	// Test too long reason
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Certificate not found"}`)) //nolint:errcheck
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
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "Certificate not found")
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

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SSLDetails{
			CommonName: "ccmqa.com",
			SSLId:      1740,
			Status:     "Requested",
			Comments:   "some comments",
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
		SSLId:            1740,
		Term:             365,
		RequesterAdminId: 1,
		ApproverAdminId:  -1,
		Comments:         "some comments",
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Invalid request"}`)) //nolint:errcheck
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
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "Invalid request")
}

func TestAutoRenewDetails_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		details  AutoRenewDetails
		expected string
	}{
		{
			name:     "empty details returns null",
			details:  AutoRenewDetails{},
			expected: "null",
		},
		{
			name:     "with state only",
			details:  AutoRenewDetails{State: "Scheduled"},
			expected: `{"state":"Scheduled","daysBeforeExpiration":0}`,
		},
		{
			name:     "with days only",
			details:  AutoRenewDetails{DaysBeforeExpiration: 30},
			expected: `{"state":"","daysBeforeExpiration":30}`,
		},
		{
			name:     "with both fields",
			details:  AutoRenewDetails{State: "Scheduled", DaysBeforeExpiration: 30},
			expected: `{"state":"Scheduled","daysBeforeExpiration":30}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.details)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestValidateUpdateSSLDetailsRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateSSLDetailsRequest
		expectedErr string
	}{
		{
			name:        "valid request",
			request:     UpdateSSLDetailsRequest{SSLId: 1740},
			expectedErr: "",
		},
		{
			name:        "invalid SSLId",
			request:     UpdateSSLDetailsRequest{SSLId: 0},
			expectedErr: "sslId must be at least 1",
		},
		{
			name:        "invalid term",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, Term: -1},
			expectedErr: "term must be at least 1",
		},
		{
			name:        "invalid certTypeId",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, CertTypeId: -1},
			expectedErr: "certTypeId must be at least 1",
		},
		{
			name:        "invalid orgId",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, OrgId: -1},
			expectedErr: "orgId must be at least 1",
		},
		{
			name:        "invalid CSR",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, CSR: "invalid_csr!"},
			expectedErr: "csr must match the regular expression",
		},
		{
			name:        "invalid comments length",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, Comments: strings.Repeat("a", 1025)},
			expectedErr: "comments maximum length is 1024 characters",
		},
		{
			name: "invalid custom field name",
			request: UpdateSSLDetailsRequest{
				SSLId:        1740,
				CustomFields: []CustomField{{Name: "", Value: "value"}},
			},
			expectedErr: "custom field name must not be null",
		},
		{
			name: "invalid autoRenewDetails state",
			request: UpdateSSLDetailsRequest{
				SSLId:            1740,
				AutoRenewDetails: &AutoRenewDetails{State: "Invalid"},
			},
			expectedErr: "autoRenewDetails.state allowed values",
		},
		{
			name:        "invalid requesterAdminId",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, RequesterAdminId: -1},
			expectedErr: "requesterAdminId must be at least 1",
		},
		{
			name:        "invalid approverAdminId",
			request:     UpdateSSLDetailsRequest{SSLId: 1740, ApproverAdminId: -2},
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
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}
