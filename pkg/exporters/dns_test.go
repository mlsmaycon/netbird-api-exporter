package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"netbird-api-exporter/pkg/netbird"
)

func TestNewDNSExporter(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewDNSExporter(client)

	if exporter == nil {
		t.Fatal("Expected exporter to be non-nil")
	}

	if exporter.client != client {
		t.Error("Expected client to be set correctly")
	}

	// Check that all metrics are initialized
	if exporter.nameserverGroupsTotal == nil {
		t.Error("Expected nameserverGroupsTotal metric to be non-nil")
	}
	if exporter.nameserverGroupsEnabled == nil {
		t.Error("Expected nameserverGroupsEnabled metric to be non-nil")
	}
	if exporter.nameserverGroupsPrimary == nil {
		t.Error("Expected nameserverGroupsPrimary metric to be non-nil")
	}
	if exporter.nameserverGroupDomains == nil {
		t.Error("Expected nameserverGroupDomains metric to be non-nil")
	}
	if exporter.nameserversTotal == nil {
		t.Error("Expected nameserversTotal metric to be non-nil")
	}
	if exporter.nameserversByType == nil {
		t.Error("Expected nameserversByType metric to be non-nil")
	}
	if exporter.nameserversByPort == nil {
		t.Error("Expected nameserversByPort metric to be non-nil")
	}
	if exporter.dnsManagementDisabled == nil {
		t.Error("Expected dnsManagementDisabled metric to be non-nil")
	}
}

func TestDNSExporter_Describe(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewDNSExporter(client)

	ch := make(chan *prometheus.Desc, 20)
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

func TestDNSExporter_Collect_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Token ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		switch r.URL.Path {
		case "/api/dns/nameservers":
			nameserverGroups := []netbird.NameserverGroup{
				{
					ID:          "ns1",
					Name:        "primary-dns",
					Description: "Primary DNS servers",
					Nameservers: []netbird.Nameserver{
						{IP: "8.8.8.8", NSType: "udp", Port: 53},
						{IP: "8.8.4.4", NSType: "udp", Port: 53},
					},
					Enabled:              true,
					Groups:               []string{"group1", "group2"},
					Primary:              true,
					Domains:              []string{"example.com", "test.com"},
					SearchDomainsEnabled: true,
				},
				{
					ID:          "ns2",
					Name:        "secondary-dns",
					Description: "Secondary DNS servers",
					Nameservers: []netbird.Nameserver{
						{IP: "1.1.1.1", NSType: "udp", Port: 53},
					},
					Enabled:              false,
					Groups:               []string{"group3"},
					Primary:              false,
					Domains:              []string{"internal.com"},
					SearchDomainsEnabled: false,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(nameserverGroups); err != nil {
				t.Errorf("Failed to encode nameserver groups: %v", err)
			}

		case "/api/dns/settings":
			settings := netbird.DNSSettings{
				Items: netbird.DNSSettingsItems{
					DisabledManagementGroups: []string{"group4", "group5"},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(settings); err != nil {
				t.Errorf("Failed to encode settings: %v", err)
			}

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewDNSExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 100)
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

	// Check nameserver groups total
	totalFound := false
	for _, family := range families {
		if family.GetName() == "netbird_dns_nameserver_groups" {
			totalFound = true
			if len(family.GetMetric()) > 0 {
				value := family.GetMetric()[0].GetGauge().GetValue()
				if value != 2 {
					t.Errorf("Expected nameserver groups total to be 2, got %f", value)
				}
			}
			break
		}
	}

	if !totalFound {
		t.Error("Expected to find nameserver groups total metric")
	}
}

func TestDNSExporter_Collect_APIError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewDNSExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 100)
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

func TestDNSExporter_Collect_EmptyResponse(t *testing.T) {
	// Create mock server that returns empty responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch r.URL.Path {
		case "/api/dns/nameservers":
			if _, err := w.Write([]byte(`[]`)); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		case "/api/dns/settings":
			if _, err := w.Write([]byte(`{"items": {"disabled_management_groups": []}}`)); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewDNSExporter(client)

	// Collect metrics
	ch := make(chan prometheus.Metric, 100)
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

func TestDNSExporter_FetchNameserverGroups_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/api/dns/nameservers" {
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

		nameserverGroups := []netbird.NameserverGroup{
			{
				ID:          "ns1",
				Name:        "test-dns",
				Description: "Test DNS servers",
				Nameservers: []netbird.Nameserver{
					{IP: "8.8.8.8", NSType: "udp", Port: 53},
				},
				Enabled: true,
				Primary: true,
				Domains: []string{"example.com"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(nameserverGroups); err != nil {
			t.Errorf("Failed to encode nameserver groups: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewDNSExporter(client)

	nameserverGroups, err := exporter.fetchNameserverGroups()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(nameserverGroups) != 1 {
		t.Errorf("Expected 1 nameserver group, got %d", len(nameserverGroups))
	}

	if nameserverGroups[0].Name != "test-dns" {
		t.Errorf("Expected nameserver group name to be test-dns, got %s", nameserverGroups[0].Name)
	}
}

func TestDNSExporter_FetchDNSSettings_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/api/dns/settings" {
			http.NotFound(w, r)
			return
		}

		settings := netbird.DNSSettings{
			Items: netbird.DNSSettingsItems{
				DisabledManagementGroups: []string{"group1", "group2"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(settings); err != nil {
			t.Errorf("Failed to encode settings: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewDNSExporter(client)

	settings, err := exporter.fetchDNSSettings()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(settings.Items.DisabledManagementGroups) != 2 {
		t.Errorf("Expected 2 disabled groups, got %d", len(settings.Items.DisabledManagementGroups))
	}
}

func TestDNSExporter_FetchNameserverGroups_InvalidJSON(t *testing.T) {
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
	exporter := NewDNSExporter(client)

	_, err := exporter.fetchNameserverGroups()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestDNSExporter_UpdateNameserverMetrics(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewDNSExporter(client)

	// Test data
	nameserverGroups := []netbird.NameserverGroup{
		{
			ID:          "ns1",
			Name:        "primary-dns",
			Description: "Primary DNS",
			Nameservers: []netbird.Nameserver{
				{IP: "8.8.8.8", NSType: "udp", Port: 53},
				{IP: "8.8.4.4", NSType: "udp", Port: 53},
			},
			Enabled:              true,
			Primary:              true,
			Domains:              []string{"example.com", "test.com"},
			SearchDomainsEnabled: true,
		},
		{
			ID:          "ns2",
			Name:        "secondary-dns",
			Description: "Secondary DNS",
			Nameservers: []netbird.Nameserver{
				{IP: "1.1.1.1", NSType: "udp", Port: 53},
			},
			Enabled:              false,
			Primary:              false,
			Domains:              []string{"internal.com"},
			SearchDomainsEnabled: false,
		},
	}

	// Call updateNameserverMetrics directly
	exporter.updateNameserverMetrics(nameserverGroups)

	// Check metric values using a registry
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Verify some key metrics
	expectedMetrics := map[string]float64{
		"netbird_dns_nameserver_groups": 2,
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
