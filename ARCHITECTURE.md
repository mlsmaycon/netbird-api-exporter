# NetBird API Exporter Architecture

## Overview

The NetBird API Exporter is built with a modular architecture that makes it easy to add support for new NetBird APIs while keeping the code organized and maintainable. The project follows Go best practices and provides comprehensive Prometheus metrics for NetBird deployments.

## Directory Structure

```bash
netbird-api-exporter/
├── main.go                     # Clean application entry point
├── pkg/                        # Core application packages
│   ├── netbird/               # NetBird API client and types
│   │   ├── client.go          # Base HTTP client for NetBird API
│   │   └── types.go           # Shared data structures (Peer, Group, User, etc.)
│   ├── exporters/             # Prometheus exporters for different APIs
│   │   ├── exporter.go        # Main composite exporter
│   │   ├── peers.go           # Peers API exporter
│   │   ├── groups.go          # Groups API exporter
│   │   ├── users.go           # Users API exporter
│   │   ├── networks.go        # Networks API exporter
│   │   ├── dns.go             # DNS API exporter
│   │   └── *_test.go          # Comprehensive test suite for each exporter
│   └── utils/                 # Utility functions
│       └── config.go          # Configuration helpers
├── charts/                     # Kubernetes deployment
│   └── netbird-api-exporter/  # Helm chart for K8s deployment
│       └── templates/         # K8s resource templates
│           └── tests/         # Helm chart tests
├── docs/                       # GitHub Pages documentation
│   ├── _config.yml            # Jekyll configuration
│   ├── _sass/                 # Styling for documentation
│   ├── getting-started/       # Getting started guides
│   ├── installation/          # Installation documentation
│   ├── index.md               # Documentation homepage
│   └── *.md                   # Various documentation files
├── tmp/                        # Temporary files (gitignored)
├── docker-compose.yml          # Local development with Docker Compose
├── Dockerfile                  # Container image definition
├── env.example                 # Environment variables template
├── Makefile                    # Build and development automation
├── netbird-exporter.service    # Systemd service file
├── prometheus.yml.example      # Example Prometheus configuration
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── LICENSE                     # Project license
├── README.md                   # Main project documentation
├── CONTRIBUTING.md             # Contribution guidelines
└── ARCHITECTURE.md             # This file
```

## Core Components

### 1. NetBird Client (`pkg/netbird/`)

The NetBird client package provides a unified interface for interacting with the NetBird API.

#### `client.go`

- Provides a reusable HTTP client for all NetBird API calls
- Handles authentication with API tokens
- Implements common error handling and logging
- Manages HTTP timeouts and retries

#### `types.go`

Contains shared data structures used across different APIs:

- `Peer`: NetBird peer information
- `Group`: NetBird group details
- `User`: User account information
- `Network`: Network configuration
- `DNSNameServerGroup`: DNS configuration
- Common utility types and constants

### 2. Exporters (`pkg/exporters/`)

Each exporter implements the `prometheus.Collector` interface and focuses on a specific NetBird API endpoint.

#### `exporter.go` - Main Composite Exporter

- Orchestrates all individual API exporters
- Provides unified metrics collection
- Implements graceful error handling and recovery
- Manages scrape timing and error metrics

#### `peers.go` - Peers API Exporter

Provides comprehensive peer metrics:

- Total peer count
- Connection status
- Operating system distribution
- Geographic distribution
- Group memberships
- SSH configuration
- Login status and approval requirements

#### `groups.go` - Groups API Exporter

Handles group-related metrics:

- Total group count
- Peer counts per group
- Resource assignments
- Group types and issued status

#### `users.go` - Users API Exporter

Manages user account metrics:

- User counts by role and status
- Service user identification
- Login timestamps
- Permission mappings
- Auto-group assignments

#### `networks.go` - Networks API Exporter

Tracks network configuration:

- Network counts and information
- Router assignments
- Resource associations
- Policy applications
- Routing peer configurations

#### `dns.go` - DNS API Exporter

Monitors DNS configuration:

- Nameserver group management
- Domain configurations
- DNS management status
- Nameserver types and ports

### 3. Main Application (`main.go`)

Clean and focused application entry point:

- Configuration management via environment variables
- HTTP server setup with configurable endpoints
- Graceful shutdown handling
- Health check endpoints
- Structured logging setup

### 4. Configuration (`pkg/utils/config.go`)

Centralized configuration management:

- Environment variable parsing
- Default value handling
- Configuration validation
- Type-safe configuration access

## Testing Architecture

### Test Structure

Each exporter has corresponding test files following Go conventions:

- `*_test.go` files alongside implementation
- Table-driven tests for comprehensive coverage
- Mock NetBird API responses
- Error condition testing
- Metrics validation

### Test Patterns

```go
func TestExporter_CollectMetrics(t *testing.T) {
    tests := []struct {
        name           string
        apiResponse    string
        expectedMetrics map[string]float64
        expectError    bool
    }{
        // Test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Deployment Architecture

### Container Deployment

- **Dockerfile**: Multi-stage build for optimal image size
- **docker-compose.yml**: Local development environment
- **Base Image**: Minimal Alpine Linux for security

### Kubernetes Deployment

- **Helm Chart**: Production-ready Kubernetes deployment
- **ConfigMaps**: Environment-based configuration
- **ServiceMonitor**: Prometheus operator integration
- **RBAC**: Minimal required permissions

### Systemd Service

- **netbird-exporter.service**: Native Linux service deployment
- Automatic restart and logging integration
- Environment file support

## Configuration Management

### Environment Variables

```env
# NetBird API Configuration
NETBIRD_API_URL=https://api.netbird.io
NETBIRD_API_TOKEN=your_token_here

# Server Configuration
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics

# Logging
LOG_LEVEL=info
```

### Configuration Hierarchy

1. Environment variables (highest priority)
2. Configuration files
3. Default values (lowest priority)

## Documentation Architecture

### GitHub Pages Integration

- **Jekyll-based**: Automatic documentation generation
- **Responsive Design**: Mobile-friendly documentation
- **Search Integration**: Full-text documentation search

### Documentation Structure

- **Getting Started**: Quick setup guides
- **Installation**: Detailed deployment instructions
- **API Reference**: Metrics documentation
- **Examples**: Real-world configuration examples

## Adding New NetBird APIs

The modular architecture makes adding new APIs straightforward:

### Step 1: Define Types

Add new data structures to `pkg/netbird/types.go`:

```go
// Policy represents a NetBird policy from the API
type Policy struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Enabled     bool     `json:"enabled"`
    Rules       []Rule   `json:"rules"`
}
```

### Step 2: Create Exporter

Create `pkg/exporters/policies.go`:

```go
package exporters

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/matanbaruch/netbird-api-exporter/pkg/netbird"
)

type PoliciesExporter struct {
    client *netbird.Client

    // Prometheus metrics
    policiesTotal    *prometheus.GaugeVec
    policiesEnabled  *prometheus.GaugeVec
    scrapeErrors     prometheus.Counter
    scrapeDuration   prometheus.Histogram
}

func NewPoliciesExporter(client *netbird.Client) *PoliciesExporter {
    return &PoliciesExporter{
        client: client,
        policiesTotal: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "netbird_policies",
                Help: "Total number of NetBird policies",
            },
            []string{},
        ),
        // Initialize other metrics...
    }
}

func (e *PoliciesExporter) Describe(ch chan<- *prometheus.Desc) {
    e.policiesTotal.Describe(ch)
    // Describe other metrics...
}

func (e *PoliciesExporter) Collect(ch chan<- prometheus.Metric) {
    // Implementation...
}
```

### Step 3: Add to Main Exporter

Update `pkg/exporters/exporter.go`:

```go
type NetBirdExporter struct {
    client           *netbird.Client
    peersExporter    *PeersExporter
    groupsExporter   *GroupsExporter
    usersExporter    *UsersExporter
    networksExporter *NetworksExporter
    dnsExporter      *DNSExporter
    policiesExporter *PoliciesExporter  // Add new exporter
}
```

### Step 4: Add Tests

Create comprehensive tests in `pkg/exporters/policies_test.go`:

```go
func TestPoliciesExporter_Collect(t *testing.T) {
    // Test implementation
}
```

### Step 5: Update Documentation

- Add metrics to README.md
- Update architecture documentation
- Add configuration examples

## Performance Considerations

### API Rate Limiting

- Respect NetBird API rate limits
- Implement exponential backoff
- Cache responses when appropriate

### Memory Management

- Efficient metric reset patterns
- Garbage collection optimization
- Resource cleanup on shutdown

### Monitoring

- Self-monitoring with exporter metrics
- Scrape duration tracking
- Error rate monitoring

## Security Architecture

### API Token Management

- Secure token storage
- Token rotation support
- Environment-based configuration

### Container Security

- Non-root user execution
- Minimal base images
- Regular security updates

### Network Security

- TLS encryption for API calls
- Configurable timeouts
- Certificate validation

### Build Security & Supply Chain

- **Artifact Attestations**: All releases include signed build provenance attestations
- **Sigstore Integration**: Container images and binaries are signed using Sigstore
- **GitHub Actions**: Builds occur in secure, auditable GitHub Actions environment
- **Reproducible Builds**: Pinned dependencies and deterministic build process
- **SLSA Compliance**: Following SLSA (Supply-chain Levels for Software Artifacts) guidelines

#### Verification Process

Users can verify the authenticity of artifacts:

```bash
# Verify Docker image attestation
gh attestation verify oci://ghcr.io/matanbaruch/netbird-api-exporter:latest --owner matanbaruch

# Verify binary attestations
gh attestation verify netbird-api-exporter-linux-amd64 --owner matanbaruch
```

See [docs/security/artifact-attestations.md](docs/security/artifact-attestations.md) for complete verification instructions.

## Future Enhancements

### Planned Features

- Configuration-driven API selection
- Custom scrape intervals per API
- API response caching
- Rate limiting with retry logic
- Custom metric filtering
- Multi-tenant support

### Monitoring Improvements

- API health status metrics
- Response time percentiles
- Error categorization
- Performance profiling

### Deployment Enhancements

- Operator pattern for Kubernetes
- Auto-scaling based on load
- Multi-region deployment support
- Blue-green deployment patterns

## Contributing to the Architecture

When contributing to the project architecture:

1. **Follow Patterns**: Use existing patterns for consistency
2. **Test Coverage**: Ensure comprehensive test coverage
3. **Documentation**: Update architectural documentation
4. **Backwards Compatibility**: Consider existing deployments
5. **Performance**: Profile and benchmark changes
6. **Security**: Consider security implications

## Metrics Standards

### Naming Conventions

- Use `netbird_` prefix for all metrics
- Use `snake_case` for metric names
- Include units in metric names when appropriate
- Use consistent label naming

### Label Guidelines

- Keep label cardinality reasonable
- Use meaningful label names
- Avoid high-cardinality labels
- Consider label stability over time

### Metric Types

- **Gauges**: For current state values
- **Counters**: For monotonically increasing values
- **Histograms**: For timing and size distributions
- **Summaries**: For quantile calculations

---

This architecture provides a solid foundation for monitoring NetBird deployments while maintaining flexibility for future enhancements and easy maintenance.
