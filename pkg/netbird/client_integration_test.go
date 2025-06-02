package netbird

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClient_Integration_RealAPI(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	baseURL := getenvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	client := NewClient(baseURL, token)

	// Test basic connectivity
	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		fmt.Sprintf("%s/api/users", client.GetBaseURL()),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", client.GetToken()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.GetHTTPClient().Do(req)
	if err != nil {
		t.Fatalf("Failed to make API request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}
}

func TestClient_Integration_AllEndpoints(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	baseURL := getenvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	client := NewClient(baseURL, token)

	endpoints := []string{
		"/api/peers",
		"/api/groups",
		"/api/users",
		"/api/dns/nameservers",
		"/api/dns/settings",
		"/api/networks",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req, err := http.NewRequestWithContext(
				context.Background(),
				"GET",
				fmt.Sprintf("%s%s", client.GetBaseURL(), endpoint),
				nil,
			)
			if err != nil {
				t.Fatalf("Failed to create request for %s: %v", endpoint, err)
			}

			req.Header.Set("Authorization", fmt.Sprintf("Token %s", client.GetToken()))
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.GetHTTPClient().Do(req)
			if err != nil {
				t.Fatalf("Failed to make API request to %s: %v", endpoint, err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Logf("Failed to close response body: %v", err)
				}
			}()

			// Accept both 200 and 429 (rate limited) as valid responses for integration tests
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Endpoint %s returned unexpected status %d", endpoint, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				t.Logf("Endpoint %s is rate limited - this is expected in integration tests", endpoint)
			}
		})
	}
}

func TestClient_Integration_Timeout(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	baseURL := getenvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	client := NewClient(baseURL, token)

	// Verify that the client has a proper timeout configured
	httpClient := client.GetHTTPClient()
	if httpClient.Timeout == 0 {
		t.Error("Expected HTTP client to have a timeout configured")
	}

	if httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected HTTP client timeout to be 30s, got %v", httpClient.Timeout)
	}
}

func TestClient_Integration_ErrorHandling(t *testing.T) {
	// Test with invalid token
	invalidClient := NewClient("https://api.netbird.io", "invalid-token")

	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		fmt.Sprintf("%s/api/users", invalidClient.GetBaseURL()),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", invalidClient.GetToken()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := invalidClient.GetHTTPClient().Do(req)
	if err != nil {
		t.Fatalf("Failed to make API request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	// Should get 401 Unauthorized for invalid token
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid token, got %d", resp.StatusCode)
	}
}

func TestClient_Integration_RateLimiting(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	baseURL := getenvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	client := NewClient(baseURL, token)

	// Make multiple rapid requests to test rate limiting behavior
	const numRequests = 3
	results := make([]int, numRequests)

	for i := 0; i < numRequests; i++ {
		req, err := http.NewRequestWithContext(
			context.Background(),
			"GET",
			fmt.Sprintf("%s/api/users", client.GetBaseURL()),
			nil,
		)
		if err != nil {
			t.Fatalf("Request %d: Failed to create request: %v", i, err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Token %s", client.GetToken()))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.GetHTTPClient().Do(req)
		if err != nil {
			t.Errorf("Request %d: Failed to make API request: %v", i, err)
			continue
		}
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}

		results[i] = resp.StatusCode

		// Small delay between requests
		time.Sleep(100 * time.Millisecond)
	}

	// Check that we got reasonable responses
	for i, statusCode := range results {
		if statusCode != http.StatusOK && statusCode != http.StatusTooManyRequests {
			t.Errorf("Request %d: Unexpected status code %d", i, statusCode)
		}
	}
}

func TestClient_Integration_CustomBaseURL(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	// Test with custom base URL (should still work with default)
	customURL := getenvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	client := NewClient(customURL, token)

	if client.GetBaseURL() != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.GetBaseURL())
	}

	// Test that client works with the custom URL
	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		fmt.Sprintf("%s/api/users", client.GetBaseURL()),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", client.GetToken()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.GetHTTPClient().Do(req)
	if err != nil {
		t.Fatalf("Failed to make API request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	// Accept both 200 and 429 as valid responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 200 or 429, got %d", resp.StatusCode)
	}
}

// Helper function for environment variable defaults
func getenvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
