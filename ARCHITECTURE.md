# NetBird API Exporter Architecture

## Overview

The NetBird API Exporter has been refactored into a modular architecture that makes it easy to add support for new NetBird APIs while keeping the code organized and maintainable.

## Directory Structure

```bash
netbird-api-exporter/
├── main.go                     # Clean application entry point
├── pkg/
│   ├── netbird/               # NetBird API client and types
│   │   ├── client.go          # Base HTTP client for NetBird API
│   │   └── types.go           # Shared data structures (Peer, Group, etc.)
│   ├── exporters/             # Prometheus exporters for different APIs
│   │   ├── exporter.go        # Main composite exporter
│   │   ├── peers.go           # Peers API exporter
│   │   └── groups.go          # Groups API exporter (placeholder)
│   └── utils/                 # Utility functions
│       └── config.go          # Configuration helpers
├── go.mod
├── go.sum
└── README.md
```

## Key Components

### 1. NetBird Client (`pkg/netbird/`)

- **`client.go`**: Provides a reusable HTTP client for all NetBird API calls
- **`types.go`**: Contains shared data structures used across different APIs

### 2. Exporters (`pkg/exporters/`)

- **`exporter.go`**: Main composite exporter that combines all individual API exporters
- **`peers.go`**: Handles all peers-related metrics collection
- **`groups.go`**: Placeholder for groups API metrics (ready for implementation)

### 3. Main Application (`main.go`)

Clean and focused on:

- Configuration management
- HTTP server setup
- Graceful shutdown
- Health checks

## Adding New NetBird APIs

To add support for a new NetBird API (e.g., Users, Policies, Routes), follow these steps:

### Step 1: Add Types (if needed)

If the new API introduces new data structures, add them to `pkg/netbird/types.go`:

```go
// User represents a NetBird user from the API
type User struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    // ... other fields
}
```

### Step 2: Create API Exporter

Create a new file `pkg/exporters/{api_name}.go`:

```go
package exporters

import (
    // ... imports
    "netbird-api-exporter/pkg/netbird"
)

// UsersExporter handles users-specific metrics collection
type UsersExporter struct {
    client *netbird.Client
    
    // Prometheus metrics
    usersTotal *prometheus.GaugeVec
    // ... other metrics
}

// NewUsersExporter creates a new users exporter
func NewUsersExporter(client *netbird.Client) *UsersExporter {
    return &UsersExporter{
        client: client,
        // Initialize metrics...
    }
}

// Implement prometheus.Collector interface
func (e *UsersExporter) Describe(ch chan<- *prometheus.Desc) { /* ... */ }
func (e *UsersExporter) Collect(ch chan<- prometheus.Metric) { /* ... */ }

// API-specific methods
func (e *UsersExporter) fetchUsers() ([]netbird.User, error) { /* ... */ }
func (e *UsersExporter) updateMetrics(users []netbird.User) { /* ... */ }
```

### Step 3: Update Main Exporter

Add the new exporter to `pkg/exporters/exporter.go`:

```go
type NetBirdExporter struct {
    client        *netbird.Client
    peersExporter *PeersExporter
    usersExporter *UsersExporter  // Add this
    // ... other exporters
}

func NewNetBirdExporter(baseURL, token string) *NetBirdExporter {
    client := netbird.NewClient(baseURL, token)
    
    return &NetBirdExporter{
        client:        client,
        peersExporter: NewPeersExporter(client),
        usersExporter: NewUsersExporter(client),  // Add this
        // ... initialize other exporters
    }
}

func (e *NetBirdExporter) Describe(ch chan<- *prometheus.Desc) {
    e.peersExporter.Describe(ch)
    e.usersExporter.Describe(ch)  // Add this
    // ... other exporters
}

func (e *NetBirdExporter) Collect(ch chan<- prometheus.Metric) {
    // ... existing code ...
    
    // Add collection for new exporter
    func() {
        defer func() {
            if r := recover(); r != nil {
                logrus.WithField("panic", r).Error("Panic during users collection")
                e.scrapeErrors.Inc()
            }
        }()
        e.usersExporter.Collect(ch)
    }()
}
```

## Benefits of This Architecture

1. **Separation of Concerns**: Each API has its own dedicated exporter
2. **Reusability**: Common client and types are shared
3. **Extensibility**: Easy to add new APIs without modifying existing code
4. **Maintainability**: Clean structure makes debugging and updates easier
5. **Testability**: Each component can be tested independently
6. **Performance**: Individual APIs can be optimized independently

## Example: Current Peers API

The peers API implementation demonstrates the pattern:

- **Data Fetching**: `fetchPeers()` method handles HTTP request
- **Metrics Processing**: `updateMetrics()` method processes data and updates Prometheus metrics
- **Error Handling**: Graceful error handling with logging
- **Reset Logic**: Metrics are reset before each collection cycle

## Configuration

The application uses environment variables for configuration:

- `NETBIRD_API_URL`: NetBird API base URL (default: <https://api.netbird.io>)
- `NETBIRD_API_TOKEN`: NetBird API authentication token (required)
- `LISTEN_ADDRESS`: HTTP server listen address (default: :8080)
- `METRICS_PATH`: Prometheus metrics endpoint path (default: /metrics)
- `LOG_LEVEL`: Logging level (default: info)

## Future Enhancements

- Add configuration options to enable/disable specific APIs
- Implement API response caching for better performance
- Add custom scrape intervals per API
- Include API health status metrics
- Add rate limiting and retry logic
