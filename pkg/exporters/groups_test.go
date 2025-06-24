package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/matanbaruch/netbird-api-exporter/pkg/netbird"
)

func TestNewGroupsExporter(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewGroupsExporter(client)

	if exporter == nil {
		t.Fatal("Expected exporter to be non-nil")
	}

	if exporter.client != client {
		t.Error("Expected client to be set correctly")
	}

	// Check that all metrics are initialized
	if exporter.groupsTotal == nil {
		t.Error("Expected groupsTotal metric to be non-nil")
	}
	if exporter.groupPeersCount == nil {
		t.Error("Expected groupPeersCount metric to be non-nil")
	}
	if exporter.groupResourcesCount == nil {
		t.Error("Expected groupResourcesCount metric to be non-nil")
	}
	if exporter.groupInfo == nil {
		t.Error("Expected groupInfo metric to be non-nil")
	}
	if exporter.groupResourcesByType == nil {
		t.Error("Expected groupResourcesByType metric to be non-nil")
	}
}

func TestGroupsExporter_Describe(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewGroupsExporter(client)

	ch := make(chan *prometheus.Desc, 10)
	go func() {
		exporter.Describe(ch)
		close(ch)
	}()

	count := 0
	for desc := range ch {
		if desc == nil {
			t.Error("Expected metric description to be non-nil")
		}
		count++
	}

	if count == 0 {
		t.Error("Expected at least one metric description")
	}
}

func TestGroupsExporter_Collect_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/groups" {
			http.NotFound(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Token ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		groups := []netbird.Group{
			{
				ID:             "group1",
				Name:           "admin-group",
				PeersCount:     5,
				ResourcesCount: 2,
				Issued:         "api",
				Peers: []netbird.GroupPeer{
					{ID: "peer1", Name: "peer-1"},
					{ID: "peer2", Name: "peer-2"},
				},
				Resources: []netbird.GroupResource{
					{ID: "resource1", Type: "host"},
				},
			},
			{
				ID:             "group2",
				Name:           "user-group",
				PeersCount:     10,
				ResourcesCount: 0,
				Issued:         "integration",
				Peers: []netbird.GroupPeer{
					{ID: "peer3", Name: "peer-3"},
					{ID: "peer4", Name: "peer-4"},
					{ID: "peer5", Name: "peer-5"},
				},
				Resources: []netbird.GroupResource{},
			},
			{
				ID:             "group3",
				Name:           "service-group",
				PeersCount:     2,
				ResourcesCount: 3,
				Issued:         "api",
				Peers:          []netbird.GroupPeer{},
				Resources: []netbird.GroupResource{
					{ID: "resource2", Type: "domain"},
					{ID: "resource3", Type: "subnet"},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(groups); err != nil {
			t.Errorf("Failed to encode groups: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewGroupsExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 50)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	metrics := make([]prometheus.Metric, 0)
	for metric := range ch {
		if metric != nil {
			metrics = append(metrics, metric)
		}
	}

	if len(metrics) == 0 {
		t.Error("Expected at least one metric to be collected")
	}

	// Test specific metric values
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Check groups total
	totalFound := false
	for _, family := range families {
		if family.GetName() == "netbird_groups" {
			totalFound = true
			if len(family.GetMetric()) > 0 {
				value := family.GetMetric()[0].GetGauge().GetValue()
				if value != 3 {
					t.Errorf("Expected groups total to be 3, got %f", value)
				}
			}
			break
		}
	}

	if !totalFound {
		t.Error("Expected to find groups total metric")
	}
}

func TestGroupsExporter_Collect_APIError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewGroupsExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 50)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	// Should still complete without error (though may not collect useful data)
	metricCount := 0
	for range ch {
		metricCount++
	}

	// Even with API errors, some metrics might still be collected
	// This test ensures the collection doesn't panic or hang
}

func TestGroupsExporter_Collect_EmptyResponse(t *testing.T) {
	// Create mock server that returns empty array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`[]`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewGroupsExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 50)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	metricCount := 0
	for range ch {
		metricCount++
	}

	// Should still collect metrics (zeros)
	if metricCount == 0 {
		t.Error("Expected at least one metric (even if zero)")
	}
}

func TestGroupsExporter_Collect_InvalidJSON(t *testing.T) {
	// Create mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`invalid json`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewGroupsExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 50)
	go func() {
		exporter.Collect(ch)
		close(ch)
	}()

	// Should complete without panic
	for range ch {
		// Drain channel
	}
}

func TestGroupsExporter_UpdateMetrics(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewGroupsExporter(client)

	// Test data
	groups := []netbird.Group{
		{
			ID:             "group1",
			Name:           "test-group-1",
			PeersCount:     5,
			ResourcesCount: 2,
			Issued:         "api",
		},
		{
			ID:             "group2",
			Name:           "test-group-2",
			PeersCount:     3,
			ResourcesCount: 1,
			Issued:         "integration",
		},
	}

	// Call updateMetrics directly
	exporter.updateMetrics(groups)

	// Check metric values using a registry
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Verify some key metrics
	expectedMetrics := map[string]float64{
		"netbird_groups": 2,
	}

	for _, family := range families {
		if expectedValue, exists := expectedMetrics[family.GetName()]; exists {
			if len(family.GetMetric()) > 0 {
				actualValue := family.GetMetric()[0].GetGauge().GetValue()
				if actualValue != expectedValue {
					t.Errorf("Expected %s to be %f, got %f", family.GetName(), expectedValue, actualValue)
				}
			}
		}
	}
}

func TestGroupsExporter_FetchGroups_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/api/groups" {
			http.NotFound(w, r)
			return
		}

		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			http.Error(w, "Invalid Accept header", http.StatusBadRequest)
			return
		}

		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Token ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		groups := []netbird.Group{
			{
				ID:             "group1",
				Name:           "test-group",
				PeersCount:     1,
				ResourcesCount: 0,
				Issued:         "api",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(groups); err != nil {
			t.Errorf("Failed to encode groups: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewGroupsExporter(client)

	groups, err := exporter.fetchGroups()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}

	if groups[0].Name != "test-group" {
		t.Errorf("Expected group name to be test-group, got %s", groups[0].Name)
	}
}

func TestGroupsExporter_FetchGroups_InvalidJSON(t *testing.T) {
	// Create mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`invalid json`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewGroupsExporter(client)

	_, err := exporter.fetchGroups()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestGroupsExporter_MetricsReset(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewGroupsExporter(client)

	// Set some initial values
	exporter.groupsTotal.WithLabelValues().Set(10)
	exporter.groupPeersCount.WithLabelValues("group1", "test-group", "api").Set(5)

	// Create empty groups to test reset
	groups := []netbird.Group{}
	exporter.updateMetrics(groups)

	// Collect and verify metrics are reset
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	for _, family := range families {
		if family.GetName() == "netbird_groups" {
			if len(family.GetMetric()) > 0 {
				value := family.GetMetric()[0].GetGauge().GetValue()
				if value != 0 {
					t.Errorf("Expected groups total to be reset to 0, got %f", value)
				}
			}
		}
	}
}

func TestGroupsExporter_MetricLabels(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewGroupsExporter(client)

	// Test data with specific values for label verification
	groups := []netbird.Group{
		{
			ID:             "group1",
			Name:           "test-group",
			PeersCount:     5,
			ResourcesCount: 2,
			Issued:         "api",
		},
	}

	// Call updateMetrics directly instead of Collect to avoid API calls
	exporter.updateMetrics(groups)

	// Collect metrics to verify labels
	ch := make(chan prometheus.Metric, 50)
	go func() {
		// Only collect the metrics we've set, not trigger API calls
		exporter.groupsTotal.Collect(ch)
		exporter.groupPeersCount.Collect(ch)
		exporter.groupResourcesCount.Collect(ch)
		exporter.groupInfo.Collect(ch)
		close(ch)
	}()

	labelFound := false
	for metric := range ch {
		// Check if metric has expected labels - this is a basic check
		// In a real scenario, you'd want to inspect the metric's labels more thoroughly
		if metric != nil {
			labelFound = true
		}
	}

	if !labelFound {
		t.Error("Expected to find metrics with labels")
	}
}
