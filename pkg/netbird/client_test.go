package netbird

import (
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		token    string
		expected struct {
			baseURL string
			token   string
		}
	}{
		{
			name:    "basic client creation",
			baseURL: "https://api.netbird.io",
			token:   "test-token",
			expected: struct {
				baseURL string
				token   string
			}{
				baseURL: "https://api.netbird.io",
				token:   "test-token",
			},
		},
		{
			name:    "client creation with trailing slash",
			baseURL: "https://api.netbird.io/",
			token:   "another-token",
			expected: struct {
				baseURL string
				token   string
			}{
				baseURL: "https://api.netbird.io",
				token:   "another-token",
			},
		},
		{
			name:    "client creation with multiple trailing slashes",
			baseURL: "https://api.netbird.io///",
			token:   "multi-slash-token",
			expected: struct {
				baseURL string
				token   string
			}{
				baseURL: "https://api.netbird.io//",
				token:   "multi-slash-token",
			},
		},
		{
			name:    "client creation with empty token",
			baseURL: "https://api.netbird.io",
			token:   "",
			expected: struct {
				baseURL string
				token   string
			}{
				baseURL: "https://api.netbird.io",
				token:   "",
			},
		},
		{
			name:    "client creation with localhost URL",
			baseURL: "http://localhost:8080",
			token:   "localhost-token",
			expected: struct {
				baseURL string
				token   string
			}{
				baseURL: "http://localhost:8080",
				token:   "localhost-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.token)

			if client == nil {
				t.Fatal("Expected client to be non-nil")
			}

			if client.baseURL != tt.expected.baseURL {
				t.Errorf("Expected baseURL to be %q, got %q", tt.expected.baseURL, client.baseURL)
			}

			if client.token != tt.expected.token {
				t.Errorf("Expected token to be %q, got %q", tt.expected.token, client.token)
			}

			if client.httpClient == nil {
				t.Error("Expected httpClient to be non-nil")
			}

			// Check that HTTP client has the expected timeout
			expectedTimeout := 30 * time.Second
			if client.httpClient.Timeout != expectedTimeout {
				t.Errorf("Expected HTTP client timeout to be %v, got %v", expectedTimeout, client.httpClient.Timeout)
			}
		})
	}
}

func TestClient_GetHTTPClient(t *testing.T) {
	client := NewClient("https://api.netbird.io", "test-token")

	httpClient := client.GetHTTPClient()

	if httpClient == nil {
		t.Fatal("Expected HTTP client to be non-nil")
	}

	if httpClient != client.httpClient {
		t.Error("Expected returned HTTP client to be the same as internal client")
	}

	// Verify it's a properly configured HTTP client
	if httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected HTTP client timeout to be 30s, got %v", httpClient.Timeout)
	}
}

func TestClient_GetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "URL without trailing slash",
			inputURL: "https://api.netbird.io",
			expected: "https://api.netbird.io",
		},
		{
			name:     "URL with trailing slash",
			inputURL: "https://api.netbird.io/",
			expected: "https://api.netbird.io",
		},
		{
			name:     "localhost URL",
			inputURL: "http://localhost:8080",
			expected: "http://localhost:8080",
		},
		{
			name:     "localhost URL with trailing slash",
			inputURL: "http://localhost:8080/",
			expected: "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.inputURL, "test-token")
			baseURL := client.GetBaseURL()

			if baseURL != tt.expected {
				t.Errorf("Expected base URL to be %q, got %q", tt.expected, baseURL)
			}
		})
	}
}

func TestClient_GetToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "normal token",
			token: "test-token-123",
		},
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "long token",
			token: strings.Repeat("a", 100),
		},
		{
			name:  "token with special characters",
			token: "token-with-!@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("https://api.netbird.io", tt.token)
			token := client.GetToken()

			if token != tt.token {
				t.Errorf("Expected token to be %q, got %q", tt.token, token)
			}
		})
	}
}

func TestClient_Immutability(t *testing.T) {
	originalURL := "https://api.netbird.io"
	originalToken := "original-token"
	client := NewClient(originalURL, originalToken)

	// Get the values
	baseURL := client.GetBaseURL()
	token := client.GetToken()
	httpClient := client.GetHTTPClient()

	// Verify original values
	if baseURL != originalURL {
		t.Errorf("Expected base URL to be %q, got %q", originalURL, baseURL)
	}
	if token != originalToken {
		t.Errorf("Expected token to be %q, got %q", originalToken, token)
	}
	if httpClient == nil {
		t.Fatal("Expected HTTP client to be non-nil")
	}

	// Modify the returned HTTP client timeout to ensure it doesn't affect the original
	httpClient.Timeout = 60 * time.Second

	// Verify that getting the client again still returns the modified client
	// (since we return the actual client reference, not a copy)
	secondHTTPClient := client.GetHTTPClient()
	if secondHTTPClient.Timeout != 60*time.Second {
		t.Error("Expected HTTP client modifications to persist (reference semantics)")
	}

	// But base URL and token should remain unchanged
	if client.GetBaseURL() != originalURL {
		t.Error("Base URL should not change")
	}
	if client.GetToken() != originalToken {
		t.Error("Token should not change")
	}
}

func TestClient_HTTPClientConfiguration(t *testing.T) {
	client := NewClient("https://api.netbird.io", "test-token")
	httpClient := client.GetHTTPClient()

	// Test that the HTTP client is properly configured
	if httpClient == nil {
		t.Fatal("Expected HTTP client to be non-nil")
	}

	// Check timeout
	expectedTimeout := 30 * time.Second
	if httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected HTTP client timeout to be %v, got %v", expectedTimeout, httpClient.Timeout)
	}

	// Verify it's an actual http.Client (it's already *http.Client type)
}

func TestClient_URLNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no trailing slash",
			input:    "https://api.netbird.io",
			expected: "https://api.netbird.io",
		},
		{
			name:     "single trailing slash",
			input:    "https://api.netbird.io/",
			expected: "https://api.netbird.io",
		},
		{
			name:     "multiple trailing slashes only removes one",
			input:    "https://api.netbird.io//",
			expected: "https://api.netbird.io/",
		},
		{
			name:     "path with trailing slash",
			input:    "https://api.netbird.io/v1/",
			expected: "https://api.netbird.io/v1",
		},
		{
			name:     "path without trailing slash",
			input:    "https://api.netbird.io/v1",
			expected: "https://api.netbird.io/v1",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just slash",
			input:    "/",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.input, "test-token")
			result := client.GetBaseURL()

			if result != tt.expected {
				t.Errorf("Expected normalized URL to be %q, got %q", tt.expected, result)
			}
		})
	}
}
