package exporters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"netbird-api-exporter/pkg/netbird"
)

func TestNewUsersExporter(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewUsersExporter(client)

	if exporter == nil {
		t.Fatal("Expected exporter to be non-nil")
	}

	if exporter.client != client {
		t.Error("Expected client to be set correctly")
	}

	// Check that all metrics are initialized
	metrics := map[string]interface{}{
		"usersTotal":           exporter.usersTotal,
		"usersByRole":          exporter.usersByRole,
		"usersByStatus":        exporter.usersByStatus,
		"usersServiceUsers":    exporter.usersServiceUsers,
		"usersBlocked":         exporter.usersBlocked,
		"usersByIssued":        exporter.usersByIssued,
		"usersLastLogin":       exporter.usersLastLogin,
		"usersAutoGroupsCount": exporter.usersAutoGroupsCount,
		"usersRestricted":      exporter.usersRestricted,
		"usersPermissions":     exporter.usersPermissions,
		"scrapeErrorsTotal":    exporter.scrapeErrorsTotal,
		"scrapeDuration":       exporter.scrapeDuration,
	}

	for name, metric := range metrics {
		if metric == nil {
			t.Errorf("Expected %s metric to be non-nil", name)
		}
	}
}

func TestUsersExporter_Describe(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewUsersExporter(client)

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

func TestUsersExporter_Collect_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users" {
			http.NotFound(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Token ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		users := []netbird.User{
			{
				ID:            "user1",
				Email:         "admin@example.com",
				Name:          "Admin User",
				Role:          "admin",
				Status:        "active",
				LastLogin:     time.Now(),
				AutoGroups:    []string{"group1", "group2"},
				IsCurrent:     false,
				IsServiceUser: false,
				IsBlocked:     false,
				Issued:        "api",
				Permissions: netbird.UserPermissions{
					IsRestricted: false,
					Modules: map[string]map[string]bool{
						"peers": {
							"read":   true,
							"write":  true,
							"delete": true,
						},
					},
				},
			},
			{
				ID:            "user2",
				Email:         "service@example.com",
				Name:          "Service User",
				Role:          "user",
				Status:        "active",
				LastLogin:     time.Now().Add(-time.Hour),
				AutoGroups:    []string{"group1"},
				IsCurrent:     false,
				IsServiceUser: true,
				IsBlocked:     false,
				Issued:        "integration",
				Permissions: netbird.UserPermissions{
					IsRestricted: true,
					Modules: map[string]map[string]bool{
						"peers": {
							"read":   true,
							"write":  false,
							"delete": false,
						},
					},
				},
			},
			{
				ID:            "user3",
				Email:         "blocked@example.com",
				Name:          "Blocked User",
				Role:          "user",
				Status:        "inactive",
				LastLogin:     time.Now().Add(-24 * time.Hour),
				AutoGroups:    []string{},
				IsCurrent:     false,
				IsServiceUser: false,
				IsBlocked:     true,
				Issued:        "api",
				Permissions: netbird.UserPermissions{
					IsRestricted: false,
					Modules:      map[string]map[string]bool{},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			t.Errorf("Failed to encode users: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewUsersExporter(client)

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

	// Check users total
	totalFound := false
	for _, family := range families {
		if family.GetName() == "netbird_users_total" {
			totalFound = true
			if len(family.GetMetric()) > 0 {
				value := family.GetMetric()[0].GetGauge().GetValue()
				if value != 3 {
					t.Errorf("Expected users total to be 3, got %f", value)
				}
			}
			break
		}
	}

	if !totalFound {
		t.Error("Expected to find users total metric")
	}
}

func TestUsersExporter_Collect_APIError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewUsersExporter(client)

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

	// Even with API errors, some metrics might still be collected (like error counters)
	// This test ensures the collection doesn't panic or hang
}

func TestUsersExporter_Collect_EmptyResponse(t *testing.T) {
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
	exporter := NewUsersExporter(client)

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

func TestUsersExporter_UpdateMetrics(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewUsersExporter(client)

	// Test data
	users := []netbird.User{
		{
			ID:            "user1",
			Email:         "admin@example.com",
			Name:          "Admin User",
			Role:          "admin",
			Status:        "active",
			LastLogin:     time.Now(),
			AutoGroups:    []string{"group1", "group2"},
			IsServiceUser: false,
			IsBlocked:     false,
			Issued:        "api",
			Permissions: netbird.UserPermissions{
				IsRestricted: false,
				Modules: map[string]map[string]bool{
					"peers": {
						"read":  true,
						"write": true,
					},
				},
			},
		},
		{
			ID:            "user2",
			Email:         "user@example.com",
			Name:          "Regular User",
			Role:          "user",
			Status:        "active",
			IsServiceUser: true,
			IsBlocked:     false,
			Issued:        "integration",
			AutoGroups:    []string{"group1"},
			Permissions: netbird.UserPermissions{
				IsRestricted: true,
				Modules:      map[string]map[string]bool{},
			},
		},
	}

	// Call updateMetrics directly
	exporter.updateMetrics(users)

	// Check metric values using a registry
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Verify some key metrics
	expectedMetrics := map[string]float64{
		"netbird_users_total": 2,
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

func TestUsersExporter_FetchUsers_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/api/users" {
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

		users := []netbird.User{
			{
				ID:     "user1",
				Email:  "test@example.com",
				Name:   "Test User",
				Role:   "user",
				Status: "active",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			t.Errorf("Failed to encode users: %v", err)
		}
	}))
	defer server.Close()

	client := netbird.NewClient(server.URL, "test-token")
	exporter := NewUsersExporter(client)

	users, err := exporter.fetchUsers()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if users[0].Email != "test@example.com" {
		t.Errorf("Expected user email to be test@example.com, got %s", users[0].Email)
	}
}

func TestUsersExporter_FetchUsers_InvalidJSON(t *testing.T) {
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
	exporter := NewUsersExporter(client)

	_, err := exporter.fetchUsers()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestUsersExporter_MetricsReset(t *testing.T) {
	client := netbird.NewClient("https://api.netbird.io", "test-token")
	exporter := NewUsersExporter(client)

	// Set some initial values
	exporter.usersTotal.WithLabelValues().Set(10)
	exporter.usersByRole.WithLabelValues("admin").Set(5)

	// Create empty users to test reset
	users := []netbird.User{}
	exporter.updateMetrics(users)

	// Collect and verify metrics are reset
	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	for _, family := range families {
		if family.GetName() == "netbird_users_total" {
			if len(family.GetMetric()) > 0 {
				value := family.GetMetric()[0].GetGauge().GetValue()
				if value != 0 {
					t.Errorf("Expected users total to be reset to 0, got %f", value)
				}
			}
		}
	}
}
