package pkg

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"netbird-api-exporter/pkg/exporters"
	"netbird-api-exporter/pkg/netbird"
	"netbird-api-exporter/pkg/utils"
)

// Integration tests that require a real NetBird API token
// These tests are skipped if NETBIRD_API_TOKEN is not set

func TestIntegration_SkipIfNoToken(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration tests: NETBIRD_API_TOKEN environment variable not set")
	}
}

func getTestClient(t *testing.T) *netbird.Client {
	t.Helper()

	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	baseURL := utils.GetEnvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	client := netbird.NewClient(baseURL, token)

	if client == nil {
		t.Fatal("Failed to create NetBird client")
	}

	return client
}

func TestIntegration_NetBirdClient_HTTPStatus(t *testing.T) {
	client := getTestClient(t)

	// Test API connectivity by making a simple request
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

	httpClient := client.GetHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make API request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify proper content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}
}

func TestIntegration_NetBirdExporter_RealAPI(t *testing.T) {
	token := os.Getenv("NETBIRD_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: NETBIRD_API_TOKEN environment variable not set")
	}

	baseURL := utils.GetEnvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	exporter := exporters.NewNetBirdExporter(baseURL, token)

	// Set up Prometheus registry for testing
	registry := prometheus.NewRegistry()
	if err := registry.Register(exporter); err != nil {
		t.Fatalf("Failed to register exporter: %v", err)
	}

	// Test that Describe works
	ch := make(chan *prometheus.Desc, 100)
	go func() {
		exporter.Describe(ch)
		close(ch)
	}()

	descCount := 0
	for desc := range ch {
		if desc == nil {
			t.Error("Received nil metric description")
		}
		descCount++
	}

	if descCount == 0 {
		t.Error("Expected at least one metric description")
	}

	// Test that Collect works with real API
	metricsCh := make(chan prometheus.Metric, 1000)
	go func() {
		defer close(metricsCh)

		// Add timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			defer close(done)
			exporter.Collect(metricsCh)
		}()

		select {
		case <-ctx.Done():
			t.Error("Collect operation timed out")
			return
		case <-done:
			// Collection completed successfully
		}
	}()

	metricsCount := 0
	var metrics []prometheus.Metric
	for metric := range metricsCh {
		if metric == nil {
			t.Error("Received nil metric")
			continue
		}
		metrics = append(metrics, metric)
		metricsCount++
	}

	if metricsCount == 0 {
		t.Error("Expected at least one metric from real API")
	}

	t.Logf("Successfully collected %d metrics from NetBird API", metricsCount)

	// Verify we get expected metric families
	expectedMetricPrefixes := []string{
		"netbird_peers_",
		"netbird_groups_",
		"netbird_users_",
		"netbird_dns_",
		"netbird_networks_",
		"netbird_exporter_scrape_",
	}

	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		desc := metric.Desc()
		if desc != nil {
			// Extract metric name from descriptor
			descStr := desc.String()
			// Parse the fqName from the descriptor string
			// Format: Desc{fqName: "metric_name", help: "...", ...}
			start := strings.Index(descStr, `fqName: "`) + 9
			if start > 8 {
				end := strings.Index(descStr[start:], `"`)
				if end > 0 {
					metricName := descStr[start : start+end]
					metricNames[metricName] = true
				}
			}
		}
	}

	for _, prefix := range expectedMetricPrefixes {
		found := false
		for name := range metricNames {
			if strings.HasPrefix(name, prefix) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("No metrics found with prefix %s", prefix)
		}
	}
}

func TestIntegration_PeersExporter_RealData(t *testing.T) {
	client := getTestClient(t)

	peersExporter := exporters.NewPeersExporter(client)
	if peersExporter == nil {
		t.Fatal("Failed to create peers exporter")
	}

	// Test metrics collection
	metricsCh := make(chan prometheus.Metric, 100)
	go func() {
		defer close(metricsCh)
		peersExporter.Collect(metricsCh)
	}()

	var metrics []prometheus.Metric
	for metric := range metricsCh {
		metrics = append(metrics, metric)
	}

	// Should have at least the basic metrics even if no peers exist
	if len(metrics) == 0 {
		t.Error("Expected at least some peer metrics")
	}

	t.Logf("Collected %d peer metrics", len(metrics))
}

func TestIntegration_GroupsExporter_RealData(t *testing.T) {
	client := getTestClient(t)

	groupsExporter := exporters.NewGroupsExporter(client)
	if groupsExporter == nil {
		t.Fatal("Failed to create groups exporter")
	}

	metricsCh := make(chan prometheus.Metric, 100)
	go func() {
		defer close(metricsCh)
		groupsExporter.Collect(metricsCh)
	}()

	var metrics []prometheus.Metric
	for metric := range metricsCh {
		metrics = append(metrics, metric)
	}

	if len(metrics) == 0 {
		t.Error("Expected at least some group metrics")
	}

	t.Logf("Collected %d group metrics", len(metrics))
}

func TestIntegration_UsersExporter_RealData(t *testing.T) {
	client := getTestClient(t)

	usersExporter := exporters.NewUsersExporter(client)
	if usersExporter == nil {
		t.Fatal("Failed to create users exporter")
	}

	metricsCh := make(chan prometheus.Metric, 100)
	go func() {
		defer close(metricsCh)
		usersExporter.Collect(metricsCh)
	}()

	var metrics []prometheus.Metric
	for metric := range metricsCh {
		metrics = append(metrics, metric)
	}

	if len(metrics) == 0 {
		t.Error("Expected at least some user metrics")
	}

	t.Logf("Collected %d user metrics", len(metrics))
}

func TestIntegration_DNSExporter_RealData(t *testing.T) {
	client := getTestClient(t)

	dnsExporter := exporters.NewDNSExporter(client)
	if dnsExporter == nil {
		t.Fatal("Failed to create DNS exporter")
	}

	metricsCh := make(chan prometheus.Metric, 100)
	go func() {
		defer close(metricsCh)
		dnsExporter.Collect(metricsCh)
	}()

	var metrics []prometheus.Metric
	for metric := range metricsCh {
		metrics = append(metrics, metric)
	}

	if len(metrics) == 0 {
		t.Error("Expected at least some DNS metrics")
	}

	t.Logf("Collected %d DNS metrics", len(metrics))
}

func TestIntegration_NetworksExporter_RealData(t *testing.T) {
	client := getTestClient(t)

	networksExporter := exporters.NewNetworksExporter(client)
	if networksExporter == nil {
		t.Fatal("Failed to create networks exporter")
	}

	metricsCh := make(chan prometheus.Metric, 100)
	go func() {
		defer close(metricsCh)
		networksExporter.Collect(metricsCh)
	}()

	var metrics []prometheus.Metric
	for metric := range metricsCh {
		metrics = append(metrics, metric)
	}

	if len(metrics) == 0 {
		t.Error("Expected at least some network metrics")
	}

	t.Logf("Collected %d network metrics", len(metrics))
}

func TestIntegration_APIRateLimiting(t *testing.T) {
	client := getTestClient(t)

	// Test multiple rapid requests to ensure rate limiting doesn't break
	const numRequests = 5

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
			t.Errorf("Request %d: Failed to close response body: %v", i, err)
		}

		// Accept both 200 and 429 (rate limited) as valid responses
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Request %d: Unexpected status code %d", i, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			t.Logf("Request %d: Rate limited (status 429) - this is expected behavior", i)
		}

		// Small delay between requests
		time.Sleep(100 * time.Millisecond)
	}
}

func TestIntegration_APIErrorHandling(t *testing.T) {
	// Test with invalid token to ensure error handling works
	invalidClient := netbird.NewClient("https://api.netbird.io", "invalid-token")

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
			t.Errorf("Failed to close response body: %v", err)
		}
	}()

	// Should get 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid token, got %d", resp.StatusCode)
	}
}

func TestIntegration_MetricsAccuracy(t *testing.T) {
	client := getTestClient(t)

	// Test that metrics are consistent between multiple collections
	exporter := exporters.NewNetBirdExporter(client.GetBaseURL(), client.GetToken())

	// Collect metrics twice
	var firstCollection, secondCollection []prometheus.Metric

	// First collection
	metricsCh1 := make(chan prometheus.Metric, 1000)
	go func() {
		defer close(metricsCh1)
		exporter.Collect(metricsCh1)
	}()

	for metric := range metricsCh1 {
		firstCollection = append(firstCollection, metric)
	}

	// Wait a bit between collections
	time.Sleep(1 * time.Second)

	// Second collection
	metricsCh2 := make(chan prometheus.Metric, 1000)
	go func() {
		defer close(metricsCh2)
		exporter.Collect(metricsCh2)
	}()

	for metric := range metricsCh2 {
		secondCollection = append(secondCollection, metric)
	}

	// Both collections should have metrics
	if len(firstCollection) == 0 {
		t.Error("First collection returned no metrics")
	}
	if len(secondCollection) == 0 {
		t.Error("Second collection returned no metrics")
	}

	// The number of metrics should be similar (allowing for some variation due to timing)
	ratio := float64(len(secondCollection)) / float64(len(firstCollection))
	if ratio < 0.8 || ratio > 1.2 {
		t.Errorf("Metric counts vary too much between collections: %d vs %d (ratio: %.2f)",
			len(firstCollection), len(secondCollection), ratio)
	}

	t.Logf("First collection: %d metrics, Second collection: %d metrics",
		len(firstCollection), len(secondCollection))
}

func TestIntegration_LoggingConfiguration(t *testing.T) {
	// Test that logging is properly configured for integration tests
	client := getTestClient(t)

	// Temporarily set log level to debug to test logging
	originalLevel := logrus.GetLevel()
	logrus.SetLevel(logrus.DebugLevel)
	defer logrus.SetLevel(originalLevel)

	// Make a request that should generate debug logs
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
			t.Errorf("Failed to close response body: %v", err)
		}
	}()

	// Just verify the request completed successfully
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
