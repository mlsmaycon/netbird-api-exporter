package netbird

import "time"

// Peer represents a NetBird peer from the API
type Peer struct {
	ID                          string    `json:"id"`
	Name                        string    `json:"name"`
	IP                          string    `json:"ip"`
	ConnectionIP                string    `json:"connection_ip"`
	Connected                   bool      `json:"connected"`
	LastSeen                    time.Time `json:"last_seen"`
	OS                          string    `json:"os"`
	KernelVersion               string    `json:"kernel_version"`
	GeonnameID                  int       `json:"geoname_id"`
	Version                     string    `json:"version"`
	Groups                      []Group   `json:"groups"`
	SSHEnabled                  bool      `json:"ssh_enabled"`
	UserID                      string    `json:"user_id"`
	Hostname                    string    `json:"hostname"`
	UIVersion                   string    `json:"ui_version"`
	DNSLabel                    string    `json:"dns_label"`
	LoginExpirationEnabled      bool      `json:"login_expiration_enabled"`
	LoginExpired                bool      `json:"login_expired"`
	LastLogin                   time.Time `json:"last_login"`
	InactivityExpirationEnabled bool      `json:"inactivity_expiration_enabled"`
	ApprovalRequired            bool      `json:"approval_required"`
	CountryCode                 string    `json:"country_code"`
	CityName                    string    `json:"city_name"`
	SerialNumber                string    `json:"serial_number"`
	ExtraDNSLabels              []string  `json:"extra_dns_labels"`
	AccessiblePeersCount        int       `json:"accessible_peers_count"`
}

// GroupPeer represents a simplified peer reference within a group
type GroupPeer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GroupResource represents a resource within a group
type GroupResource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Group represents a NetBird group from the API
type Group struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	PeersCount     int             `json:"peers_count"`
	ResourcesCount int             `json:"resources_count"`
	Issued         string          `json:"issued"`
	Peers          []GroupPeer     `json:"peers,omitempty"`
	Resources      []GroupResource `json:"resources,omitempty"`
}

// User represents a NetBird user from the API
type User struct {
	ID            string          `json:"id"`
	Email         string          `json:"email"`
	Name          string          `json:"name"`
	Role          string          `json:"role"`
	Status        string          `json:"status"`
	LastLogin     time.Time       `json:"last_login"`
	AutoGroups    []string        `json:"auto_groups"`
	IsCurrent     bool            `json:"is_current"`
	IsServiceUser bool            `json:"is_service_user"`
	IsBlocked     bool            `json:"is_blocked"`
	Issued        string          `json:"issued"`
	Permissions   UserPermissions `json:"permissions"`
}

// UserPermissions represents user permissions from the API
type UserPermissions struct {
	IsRestricted bool                       `json:"is_restricted"`
	Modules      map[string]map[string]bool `json:"modules"`
}

// Nameserver represents a DNS nameserver configuration
type Nameserver struct {
	IP     string `json:"ip"`
	NSType string `json:"ns_type"`
	Port   int    `json:"port"`
}

// NameserverGroup represents a NetBird nameserver group from the API
type NameserverGroup struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Nameservers          []Nameserver `json:"nameservers"`
	Enabled              bool         `json:"enabled"`
	Groups               []string     `json:"groups"`
	Primary              bool         `json:"primary"`
	Domains              []string     `json:"domains"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

// DNSSettings represents NetBird DNS settings from the API
type DNSSettings struct {
	Items DNSSettingsItems `json:"items"`
}

// DNSSettingsItems represents the nested DNS settings data
type DNSSettingsItems struct {
	DisabledManagementGroups []string `json:"disabled_management_groups"`
}

// Network represents a NetBird network from the API
type Network struct {
	ID                string   `json:"id"`
	Routers           []string `json:"routers"`
	RoutingPeersCount int      `json:"routing_peers_count"`
	Resources         []string `json:"resources"`
	Policies          []string `json:"policies"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
}

// NetworkRouter represents a network router from the API
type NetworkRouter struct {
	ID         string   `json:"id"`
	Peer       string   `json:"peer"`
	PeerGroups []string `json:"peer_groups"`
	Metric     int      `json:"metric"`
	Masquerade bool     `json:"masquerade"`
	Enabled    bool     `json:"enabled"`
}

// NetworkResource represents a network resource from the API
type NetworkResource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Type        string `json:"type"`
}
