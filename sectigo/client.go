package sectigo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

// Config represents the configuration for the Sectigo client.
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
// Modified to include response body in error messages for better debugging.
// expectedStatus can be a specific status code (200, 201, 204, etc.) or 0 to accept any 2xx status code.
func (c *Client) sendRequest(ctx context.Context, req *http.Request, expectedStatus int) (*http.Response, []byte, error) {
	req = req.WithContext(ctx)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making request: %w", err)
	}

	// Always read the body first
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return resp, nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check if status code matches expectation
	statusOK := false
	if expectedStatus == 0 {
		// Accept any 2xx status code
		statusOK = resp.StatusCode >= 200 && resp.StatusCode < 300
	} else {
		// Check for specific status code
		statusOK = resp.StatusCode == expectedStatus
	}

	if !statusOK {
		// Include response body in error message
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "... (truncated)"
		}
		return resp, body, fmt.Errorf("failed request, status code: %d, response: %s", resp.StatusCode, bodyStr)
	}

	return resp, body, nil
}