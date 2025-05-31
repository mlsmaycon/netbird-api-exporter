package netbird

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPeer_JSONMarshaling(t *testing.T) {
	// Create a sample peer
	lastSeen := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	lastLogin := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)

	peer := Peer{
		ID:                          "peer-123",
		Name:                        "test-peer",
		IP:                          "100.64.0.1",
		ConnectionIP:                "192.168.1.10",
		Connected:                   true,
		LastSeen:                    lastSeen,
		OS:                          "linux",
		KernelVersion:               "5.4.0",
		GeonnameID:                  5128581,
		Version:                     "0.21.0",
		Groups:                      []Group{{ID: "group1", Name: "test-group"}},
		SSHEnabled:                  true,
		UserID:                      "user-123",
		Hostname:                    "test-hostname",
		UIVersion:                   "0.21.0",
		DNSLabel:                    "test-dns-label",
		LoginExpirationEnabled:      true,
		LoginExpired:                false,
		LastLogin:                   lastLogin,
		InactivityExpirationEnabled: false,
		ApprovalRequired:            true,
		CountryCode:                 "US",
		CityName:                    "New York",
		SerialNumber:                "SN123456",
		ExtraDNSLabels:              []string{"extra1", "extra2"},
		AccessiblePeersCount:        5,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(peer)
	if err != nil {
		t.Fatalf("Failed to marshal peer to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledPeer Peer
	err = json.Unmarshal(jsonData, &unmarshaledPeer)
	if err != nil {
		t.Fatalf("Failed to unmarshal peer from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledPeer.ID != peer.ID {
		t.Errorf("Expected ID %q, got %q", peer.ID, unmarshaledPeer.ID)
	}
	if unmarshaledPeer.Name != peer.Name {
		t.Errorf("Expected Name %q, got %q", peer.Name, unmarshaledPeer.Name)
	}
	if unmarshaledPeer.Connected != peer.Connected {
		t.Errorf("Expected Connected %v, got %v", peer.Connected, unmarshaledPeer.Connected)
	}
	if !unmarshaledPeer.LastSeen.Equal(peer.LastSeen) {
		t.Errorf("Expected LastSeen %v, got %v", peer.LastSeen, unmarshaledPeer.LastSeen)
	}
	if len(unmarshaledPeer.Groups) != len(peer.Groups) {
		t.Errorf("Expected %d groups, got %d", len(peer.Groups), len(unmarshaledPeer.Groups))
	}
	if len(unmarshaledPeer.ExtraDNSLabels) != len(peer.ExtraDNSLabels) {
		t.Errorf("Expected %d extra DNS labels, got %d", len(peer.ExtraDNSLabels), len(unmarshaledPeer.ExtraDNSLabels))
	}
}

func TestGroup_JSONMarshaling(t *testing.T) {
	group := Group{
		ID:             "group-123",
		Name:           "test-group",
		PeersCount:     5,
		ResourcesCount: 2,
		Issued:         "api",
		Peers: []GroupPeer{
			{ID: "peer1", Name: "peer-1"},
			{ID: "peer2", Name: "peer-2"},
		},
		Resources: []GroupResource{
			{ID: "resource1", Type: "host"},
			{ID: "resource2", Type: "subnet"},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("Failed to marshal group to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledGroup Group
	err = json.Unmarshal(jsonData, &unmarshaledGroup)
	if err != nil {
		t.Fatalf("Failed to unmarshal group from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledGroup.ID != group.ID {
		t.Errorf("Expected ID %q, got %q", group.ID, unmarshaledGroup.ID)
	}
	if unmarshaledGroup.PeersCount != group.PeersCount {
		t.Errorf("Expected PeersCount %d, got %d", group.PeersCount, unmarshaledGroup.PeersCount)
	}
	if len(unmarshaledGroup.Peers) != len(group.Peers) {
		t.Errorf("Expected %d peers, got %d", len(group.Peers), len(unmarshaledGroup.Peers))
	}
	if len(unmarshaledGroup.Resources) != len(group.Resources) {
		t.Errorf("Expected %d resources, got %d", len(group.Resources), len(unmarshaledGroup.Resources))
	}
}

func TestUser_JSONMarshaling(t *testing.T) {
	lastLogin := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)

	user := User{
		ID:            "user-123",
		Email:         "test@example.com",
		Name:          "Test User",
		Role:          "admin",
		Status:        "active",
		LastLogin:     lastLogin,
		AutoGroups:    []string{"group1", "group2"},
		IsCurrent:     true,
		IsServiceUser: false,
		IsBlocked:     false,
		Issued:        "api",
		Permissions: UserPermissions{
			IsRestricted: false,
			Modules: map[string]map[string]bool{
				"peers": {"read": true, "write": true},
				"users": {"read": true, "write": false},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledUser.ID != user.ID {
		t.Errorf("Expected ID %q, got %q", user.ID, unmarshaledUser.ID)
	}
	if unmarshaledUser.Email != user.Email {
		t.Errorf("Expected Email %q, got %q", user.Email, unmarshaledUser.Email)
	}
	if unmarshaledUser.IsCurrent != user.IsCurrent {
		t.Errorf("Expected IsCurrent %v, got %v", user.IsCurrent, unmarshaledUser.IsCurrent)
	}
	if len(unmarshaledUser.AutoGroups) != len(user.AutoGroups) {
		t.Errorf("Expected %d auto groups, got %d", len(user.AutoGroups), len(unmarshaledUser.AutoGroups))
	}
}

func TestNameserverGroup_JSONMarshaling(t *testing.T) {
	nsGroup := NameserverGroup{
		ID:          "ns-123",
		Name:        "test-nameserver",
		Description: "Test nameserver group",
		Nameservers: []Nameserver{
			{IP: "8.8.8.8", NSType: "udp", Port: 53},
			{IP: "1.1.1.1", NSType: "tcp", Port: 53},
		},
		Enabled:              true,
		Groups:               []string{"group1", "group2"},
		Primary:              true,
		Domains:              []string{"example.com", "test.com"},
		SearchDomainsEnabled: true,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(nsGroup)
	if err != nil {
		t.Fatalf("Failed to marshal nameserver group to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledNSGroup NameserverGroup
	err = json.Unmarshal(jsonData, &unmarshaledNSGroup)
	if err != nil {
		t.Fatalf("Failed to unmarshal nameserver group from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledNSGroup.ID != nsGroup.ID {
		t.Errorf("Expected ID %q, got %q", nsGroup.ID, unmarshaledNSGroup.ID)
	}
	if unmarshaledNSGroup.Enabled != nsGroup.Enabled {
		t.Errorf("Expected Enabled %v, got %v", nsGroup.Enabled, unmarshaledNSGroup.Enabled)
	}
	if len(unmarshaledNSGroup.Nameservers) != len(nsGroup.Nameservers) {
		t.Errorf("Expected %d nameservers, got %d", len(nsGroup.Nameservers), len(unmarshaledNSGroup.Nameservers))
	}
	if len(unmarshaledNSGroup.Domains) != len(nsGroup.Domains) {
		t.Errorf("Expected %d domains, got %d", len(nsGroup.Domains), len(unmarshaledNSGroup.Domains))
	}
}

func TestDNSSettings_JSONMarshaling(t *testing.T) {
	dnsSettings := DNSSettings{
		Items: DNSSettingsItems{
			DisabledManagementGroups: []string{"group1", "group2"},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(dnsSettings)
	if err != nil {
		t.Fatalf("Failed to marshal DNS settings to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledDNSSettings DNSSettings
	err = json.Unmarshal(jsonData, &unmarshaledDNSSettings)
	if err != nil {
		t.Fatalf("Failed to unmarshal DNS settings from JSON: %v", err)
	}

	// Verify fields
	if len(unmarshaledDNSSettings.Items.DisabledManagementGroups) != len(dnsSettings.Items.DisabledManagementGroups) {
		t.Errorf("Expected %d disabled management groups, got %d",
			len(dnsSettings.Items.DisabledManagementGroups),
			len(unmarshaledDNSSettings.Items.DisabledManagementGroups))
	}
}

func TestNetwork_JSONMarshaling(t *testing.T) {
	network := Network{
		ID:                "net-123",
		Routers:           []string{"router1", "router2"},
		RoutingPeersCount: 2,
		Resources:         []string{"resource1", "resource2"},
		Policies:          []string{"policy1"},
		Name:              "test-network",
		Description:       "Test network description",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal network to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledNetwork Network
	err = json.Unmarshal(jsonData, &unmarshaledNetwork)
	if err != nil {
		t.Fatalf("Failed to unmarshal network from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledNetwork.ID != network.ID {
		t.Errorf("Expected ID %q, got %q", network.ID, unmarshaledNetwork.ID)
	}
	if unmarshaledNetwork.RoutingPeersCount != network.RoutingPeersCount {
		t.Errorf("Expected RoutingPeersCount %d, got %d", network.RoutingPeersCount, unmarshaledNetwork.RoutingPeersCount)
	}
	if len(unmarshaledNetwork.Routers) != len(network.Routers) {
		t.Errorf("Expected %d routers, got %d", len(network.Routers), len(unmarshaledNetwork.Routers))
	}
}

func TestNetworkRouter_JSONMarshaling(t *testing.T) {
	router := NetworkRouter{
		ID:         "router-123",
		Peer:       "peer-123",
		PeerGroups: []string{"group1", "group2"},
		Metric:     100,
		Masquerade: true,
		Enabled:    true,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(router)
	if err != nil {
		t.Fatalf("Failed to marshal network router to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledRouter NetworkRouter
	err = json.Unmarshal(jsonData, &unmarshaledRouter)
	if err != nil {
		t.Fatalf("Failed to unmarshal network router from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledRouter.ID != router.ID {
		t.Errorf("Expected ID %q, got %q", router.ID, unmarshaledRouter.ID)
	}
	if unmarshaledRouter.Metric != router.Metric {
		t.Errorf("Expected Metric %d, got %d", router.Metric, unmarshaledRouter.Metric)
	}
	if unmarshaledRouter.Masquerade != router.Masquerade {
		t.Errorf("Expected Masquerade %v, got %v", router.Masquerade, unmarshaledRouter.Masquerade)
	}
	if len(unmarshaledRouter.PeerGroups) != len(router.PeerGroups) {
		t.Errorf("Expected %d peer groups, got %d", len(router.PeerGroups), len(unmarshaledRouter.PeerGroups))
	}
}

func TestNetworkResource_JSONMarshaling(t *testing.T) {
	resource := NetworkResource{
		ID:          "resource-123",
		Name:        "test-resource",
		Description: "Test resource description",
		Address:     "192.168.1.0/24",
		Type:        "subnet",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal network resource to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledResource NetworkResource
	err = json.Unmarshal(jsonData, &unmarshaledResource)
	if err != nil {
		t.Fatalf("Failed to unmarshal network resource from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledResource.ID != resource.ID {
		t.Errorf("Expected ID %q, got %q", resource.ID, unmarshaledResource.ID)
	}
	if unmarshaledResource.Name != resource.Name {
		t.Errorf("Expected Name %q, got %q", resource.Name, unmarshaledResource.Name)
	}
	if unmarshaledResource.Address != resource.Address {
		t.Errorf("Expected Address %q, got %q", resource.Address, unmarshaledResource.Address)
	}
	if unmarshaledResource.Type != resource.Type {
		t.Errorf("Expected Type %q, got %q", resource.Type, unmarshaledResource.Type)
	}
}

func TestEmptyStructs(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
	}{
		{"Peer", &Peer{}},
		{"Group", &Group{}},
		{"User", &User{}},
		{"NameserverGroup", &NameserverGroup{}},
		{"DNSSettings", &DNSSettings{}},
		{"Network", &Network{}},
		{"NetworkRouter", &NetworkRouter{}},
		{"NetworkResource", &NetworkResource{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal empty struct
			jsonData, err := json.Marshal(tt.obj)
			if err != nil {
				t.Fatalf("Failed to marshal empty %s to JSON: %v", tt.name, err)
			}

			// Verify it's valid JSON
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s JSON: %v", tt.name, err)
			}
		})
	}
}

func TestNameserver_JSONMarshaling(t *testing.T) {
	nameserver := Nameserver{
		IP:     "8.8.8.8",
		NSType: "udp",
		Port:   53,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(nameserver)
	if err != nil {
		t.Fatalf("Failed to marshal nameserver to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledNameserver Nameserver
	err = json.Unmarshal(jsonData, &unmarshaledNameserver)
	if err != nil {
		t.Fatalf("Failed to unmarshal nameserver from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledNameserver.IP != nameserver.IP {
		t.Errorf("Expected IP %q, got %q", nameserver.IP, unmarshaledNameserver.IP)
	}
	if unmarshaledNameserver.NSType != nameserver.NSType {
		t.Errorf("Expected NSType %q, got %q", nameserver.NSType, unmarshaledNameserver.NSType)
	}
	if unmarshaledNameserver.Port != nameserver.Port {
		t.Errorf("Expected Port %d, got %d", nameserver.Port, unmarshaledNameserver.Port)
	}
}

func TestUserPermissions_JSONMarshaling(t *testing.T) {
	permissions := UserPermissions{
		IsRestricted: true,
		Modules: map[string]map[string]bool{
			"peers": {
				"read":  true,
				"write": false,
			},
			"groups": {
				"read":   true,
				"write":  true,
				"delete": false,
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(permissions)
	if err != nil {
		t.Fatalf("Failed to marshal user permissions to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaledPermissions UserPermissions
	err = json.Unmarshal(jsonData, &unmarshaledPermissions)
	if err != nil {
		t.Fatalf("Failed to unmarshal user permissions from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledPermissions.IsRestricted != permissions.IsRestricted {
		t.Errorf("Expected IsRestricted %v, got %v", permissions.IsRestricted, unmarshaledPermissions.IsRestricted)
	}
	if len(unmarshaledPermissions.Modules) != len(permissions.Modules) {
		t.Errorf("Expected %d modules, got %d", len(permissions.Modules), len(unmarshaledPermissions.Modules))
	}

	// Check specific module permissions
	if peersModule, exists := unmarshaledPermissions.Modules["peers"]; exists {
		if peersModule["read"] != true {
			t.Error("Expected peers read permission to be true")
		}
		if peersModule["write"] != false {
			t.Error("Expected peers write permission to be false")
		}
	} else {
		t.Error("Expected peers module to exist in permissions")
	}
}

func TestJSONFieldTags(t *testing.T) {
	// Test that the JSON field tags are working correctly
	peer := Peer{
		ID:        "test-id",
		Name:      "test-name",
		Connected: true,
	}

	jsonData, err := json.Marshal(peer)
	if err != nil {
		t.Fatalf("Failed to marshal peer: %v", err)
	}

	// Convert to string and check for expected field names
	jsonStr := string(jsonData)

	expectedFields := []string{
		`"id":"test-id"`,
		`"name":"test-name"`,
		`"connected":true`,
	}

	for _, expected := range expectedFields {
		if !contains(jsonStr, expected) {
			t.Errorf("Expected JSON to contain %q, but got: %s", expected, jsonStr)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			indexOfSubstring(s, substr) >= 0)))
}

// Helper function to find index of substring
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
