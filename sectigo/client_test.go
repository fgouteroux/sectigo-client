package sectigo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestNewClient(t *testing.T) {
	client := NewClient(Config{
		URL:      "https://cert-manager.com",
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    false,
	})

	assert.NotNil(t, client)
	assert.Equal(t, "https://cert-manager.com", client.BaseURL)
	assert.False(t, client.Debug)
}

func TestNewClient_WithDebug(t *testing.T) {
	client := NewClient(Config{
		URL:      "https://cert-manager.com",
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    true,
	})

	assert.NotNil(t, client)
	assert.True(t, client.Debug)
}

func TestRoundTrip(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "test", r.Header.Get("login"))
		assert.Equal(t, "test", r.Header.Get("customerUri"))
		assert.Equal(t, "test", r.Header.Get("password"))
		assert.Equal(t, "application/json;charset=utf-8", r.Header.Get("Content-Type"))

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

	defer resp.Body.Close() //nolint:errcheck
	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	_ = json.Unmarshal(body, &result)
	assert.Equal(t, "value", result["key"])
}

func TestRoundTrip_WithDebug(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"value"}`)) //nolint:errcheck
	})

	client := NewClient(Config{
		URL:      mockClient.Server.URL,
		Username: "test",
		Customer: "test",
		Password: "test",
		Debug:    true,
	})

	transport := client.Client.Transport.(*authTransport)
	req, _ := http.NewRequest("GET", mockClient.Server.URL+"/test", io.NopCloser(strings.NewReader(`{"test":"data"}`)))
	resp, err := transport.RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendRequest_Success(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`)) //nolint:errcheck
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
	req, _ := http.NewRequestWithContext(ctx, "GET", mockClient.Server.URL+"/test", nil)
	resp, body, err := client.sendRequest(ctx, req, http.StatusOK)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "success")
}

func TestSendRequest_ErrorWithBody(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code":-993,"description":"Certificate orders currently restricted"}`)) //nolint:errcheck
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
	req, _ := http.NewRequestWithContext(ctx, "GET", mockClient.Server.URL+"/test", nil)
	_, _, err := client.sendRequest(ctx, req, http.StatusOK)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status code: 400")
	assert.Contains(t, err.Error(), "Certificate orders currently restricted")
}

func TestSendRequest_Accept2xx(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"created"}`)) //nolint:errcheck
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
	req, _ := http.NewRequestWithContext(ctx, "GET", mockClient.Server.URL+"/test", nil)
	resp, body, err := client.sendRequest(ctx, req, 0) // 0 = accept any 2xx

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, string(body), "created")
}

func TestSendRequest_LongBodyTruncated(t *testing.T) {
	mockClient := NewMockClient()
	defer mockClient.Close()

	longBody := strings.Repeat("a", 600)
	mockClient.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(longBody)) //nolint:errcheck
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
	req, _ := http.NewRequestWithContext(ctx, "GET", mockClient.Server.URL+"/test", nil)
	_, _, err := client.sendRequest(ctx, req, http.StatusOK)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "... (truncated)")
}
