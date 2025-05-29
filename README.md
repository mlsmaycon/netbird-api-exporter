# NetBird API Exporter

A Prometheus exporter for NetBird API that provides comprehensive metrics about your NetBird network peers, groups, users, networks, and DNS configuration. This exporter fetches data from the [NetBird REST API](https://docs.netbird.io/ipa/resources/peers), [Groups API](https://docs.netbird.io/ipa/resources/groups), [Users API](https://docs.netbird.io/ipa/resources/users), [Networks API](https://docs.netbird.io/ipa/resources/networks), and [DNS API](https://docs.netbird.io/ipa/resources/dns) and exposes it in Prometheus format.

## Features

The exporter provides the following metrics:

### Peer Metrics
- **Total Peers**: Total number of peers in your NetBird network
- **Connection Status**: Number of connected vs disconnected peers  
- **Peer Distribution**: Breakdown by operating system, country, and groups
- **Last Seen Timestamps**: When each peer was last seen
- **SSH Status**: Number of peers with SSH enabled/disabled
- **Login Status**: Number of peers with expired vs valid logins
- **Approval Status**: Number of peers requiring approval
- **Accessible Peers**: Count of accessible peers for each peer

### Group Metrics
- **Total Groups**: Total number of groups in your NetBird network
- **Group Peer Counts**: Number of peers in each group
- **Group Resource Counts**: Number of resources in each group
- **Group Information**: Metadata about each group (name, ID, issued status)
- **Resource Types**: Breakdown of resources by type within each group
- **Group Performance**: Scrape duration and error metrics for groups API

### User Metrics
- **Total Users**: Total number of users in your NetBird network
- **User Distribution**: Breakdown by role, status, and type (service vs regular users)
- **User Status**: Number of blocked vs active users
- **User Permissions**: Detailed permissions by module and action
- **Auto Groups**: Number of auto-assigned groups per user
- **Last Login Timestamps**: When each user last logged in
- **User Performance**: Scrape duration and error metrics for users API

### DNS Metrics
- **Nameserver Groups**: Total number of nameserver groups and their status
- **Nameserver Configuration**: Breakdown by nameserver type (UDP/TCP) and port
- **DNS Domains**: Number of domains configured per nameserver group
- **DNS Management**: Groups with DNS management disabled
- **Primary Nameservers**: Number of primary vs secondary nameserver groups

### Network Metrics
- **Total Networks**: Total number of networks in your NetBird deployment
- **Network Routers**: Number of routers configured in each network
- **Network Resources**: Number of resources associated with each network
- **Network Policies**: Number of policies applied to each network
- **Routing Peers**: Number of routing peers in each network
- **Network Information**: Metadata about each network (name, description)

### Exporter Metrics
- **Scrape Duration**: Time taken to collect metrics from APIs
- **Error Counts**: Number of scrape errors encountered

## Metrics Exported

### Peer Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `netbird_peers_total` | Gauge | Total number of NetBird peers | - |
| `netbird_peers_connected` | Gauge | Number of connected/disconnected peers | `connected` |
| `netbird_peer_last_seen_timestamp` | Gauge | Last seen timestamp for each peer | `peer_id`, `peer_name`, `hostname` |
| `netbird_peers_by_os` | Gauge | Number of peers by operating system | `os` |
| `netbird_peers_by_country` | Gauge | Number of peers by country/city | `country_code`, `city_name` |
| `netbird_peers_by_group` | Gauge | Number of peers by group | `group_id`, `group_name` |
| `netbird_peers_ssh_enabled` | Gauge | Number of peers with SSH enabled/disabled | `ssh_enabled` |
| `netbird_peers_login_expired` | Gauge | Number of peers with expired/valid login | `login_expired` |
| `netbird_peers_approval_required` | Gauge | Number of peers requiring/not requiring approval | `approval_required` |
| `netbird_peer_accessible_peers_count` | Gauge | Number of accessible peers for each peer | `peer_id`, `peer_name` |

### Group Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `netbird_groups_total` | Gauge | Total number of NetBird groups | - |
| `netbird_group_peers_count` | Gauge | Number of peers in each NetBird group | `group_id`, `group_name`, `issued` |
| `netbird_group_resources_count` | Gauge | Number of resources in each NetBird group | `group_id`, `group_name`, `issued` |
| `netbird_group_info` | Gauge | Information about NetBird groups (always 1) | `group_id`, `group_name`, `issued` |
| `netbird_group_resources_by_type` | Gauge | Number of resources in each group by resource type | `group_id`, `group_name`, `resource_type` |
| `netbird_groups_scrape_errors_total` | Counter | Total number of errors encountered while scraping groups | `error_type` |
| `netbird_groups_scrape_duration_seconds` | Histogram | Time spent scraping groups from the NetBird API | - |

### User Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `netbird_users_total` | Gauge | Total number of NetBird users | - |
| `netbird_users_by_role` | Gauge | Number of users by role | `role` |
| `netbird_users_by_status` | Gauge | Number of users by status | `status` |
| `netbird_users_service_users` | Gauge | Number of service users vs regular users | `is_service_user` |
| `netbird_users_blocked` | Gauge | Number of blocked vs unblocked users | `is_blocked` |
| `netbird_users_by_issued` | Gauge | Number of users by issuance type | `issued` |
| `netbird_users_restricted` | Gauge | Number of users with restricted permissions | `is_restricted` |
| `netbird_user_last_login_timestamp` | Gauge | Last login timestamp for each user | `user_id`, `user_email`, `user_name` |
| `netbird_user_auto_groups_count` | Gauge | Number of auto groups assigned to each user | `user_id`, `user_email`, `user_name` |
| `netbird_user_permissions` | Gauge | User permissions by module and action | `user_id`, `user_email`, `module`, `permission`, `value` |
| `netbird_users_scrape_errors_total` | Counter | Total number of errors encountered while scraping users | `error_type` |
| `netbird_users_scrape_duration_seconds` | Histogram | Time spent scraping users from the NetBird API | - |

### DNS Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `netbird_dns_nameserver_groups_total` | Gauge | Total number of NetBird nameserver groups | - |
| `netbird_dns_nameserver_groups_enabled` | Gauge | Number of enabled/disabled nameserver groups | `enabled` |
| `netbird_dns_nameserver_groups_primary` | Gauge | Number of primary/secondary nameserver groups | `primary` |
| `netbird_dns_nameserver_group_domains_count` | Gauge | Number of domains configured in each nameserver group | `group_id`, `group_name` |
| `netbird_dns_nameservers_total` | Gauge | Total number of nameservers in each group | `group_id`, `group_name` |
| `netbird_dns_nameservers_by_type` | Gauge | Number of nameservers by type (UDP/TCP) | `ns_type` |
| `netbird_dns_nameservers_by_port` | Gauge | Number of nameservers by port | `port` |
| `netbird_dns_management_disabled_groups_count` | Gauge | Number of groups with DNS management disabled | - |

### Network Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `netbird_networks_total` | Gauge | Total number of networks in your NetBird deployment | - |
| `netbird_network_routers_count` | Gauge | Number of routers configured in each network | `network_id`, `network_name` |
| `netbird_network_resources_count` | Gauge | Number of resources associated with each network | `network_id`, `network_name` |
| `netbird_network_policies_count` | Gauge | Number of policies applied to each network | `network_id`, `network_name` |
| `netbird_network_routing_peers_count` | Gauge | Number of routing peers in each network | `network_id`, `network_name` |
| `netbird_network_info` | Gauge | Information about networks (always 1) | `network_id`, `network_name`, `description` |
| `netbird_networks_scrape_errors_total` | Counter | Total number of errors encountered while scraping networks | `error_type` |
| `netbird_networks_scrape_duration_seconds` | Histogram | Time spent scraping networks from the NetBird API | - |

### Exporter Metrics
| Metric Name | Type | Description | Labels |
|-------------|------|-------------|--------|
| `netbird_exporter_scrape_duration_seconds` | Histogram | Time spent scraping NetBird API | - |
| `netbird_exporter_scrape_errors_total` | Counter | Total number of scrape errors | - |

## Configuration

The exporter is configured via environment variables:

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `NETBIRD_API_URL` | `https://api.netbird.io` | No | NetBird API base URL |
| `NETBIRD_API_TOKEN` | - | **Yes** | NetBird API authentication token |
| `LISTEN_ADDRESS` | `:8080` | No | Address and port to listen on |
| `METRICS_PATH` | `/metrics` | No | Path where metrics are exposed |
| `LOG_LEVEL` | `info` | No | Log level (debug, info, warn, error) |

## Getting Your NetBird API Token

1. Log into your NetBird dashboard
2. Go to **Settings** → **API Keys**  
3. Create a new API key with appropriate permissions
4. Copy the token and use it as `NETBIRD_API_TOKEN`

## Installation & Usage

### Option 1: Docker Compose (Recommended)

1. Clone this repository:
```bash
git clone <repository-url>
cd netbird-api-exporter
```

2. Create environment file:
```bash
cp env.example .env
# Edit .env with your NetBird API token
```

3. Start the exporter:
```bash
docker-compose up -d
```

### Option 2: Docker

```bash
docker build -t netbird-exporter .
docker run -d \
  -p 8080:8080 \
  -e NETBIRD_API_TOKEN=your_token_here \
  --name netbird-exporter \
  netbird-exporter
```

### Option 3: Go Binary

1. Install dependencies:
```bash
go mod download
```

2. Build and run:
```bash
export NETBIRD_API_TOKEN=your_token_here
go build -o netbird-exporter
./netbird-exporter
```

## Endpoints

- **`/metrics`** - Prometheus metrics endpoint
- **`/health`** - Health check endpoint (returns JSON)
- **`/`** - Information page with links

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'netbird-exporter'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 30s
    metrics_path: /metrics
```

## Example Queries

Here are some useful Prometheus queries:

### Peer Queries
```promql
# Total number of peers
netbird_peers_total

# Percentage of connected peers
(netbird_peers_connected{connected="true"} / netbird_peers_total) * 100

# Peers by operating system
sum by (os) (netbird_peers_by_os)

# Peers that haven't been seen in over 1 hour
(time() - netbird_peer_last_seen_timestamp) > 3600

# Number of peers requiring approval
netbird_peers_approval_required{approval_required="true"}

# Average accessible peers per peer
avg(netbird_peer_accessible_peers_count)
```

### Group Queries
```promql
# Total number of groups
netbird_groups_total

# Groups with the most peers
topk(5, netbird_group_peers_count)

# Groups with the most resources
topk(5, netbird_group_resources_count)

# Average peers per group
avg(netbird_group_peers_count)

# Groups by issued method (API vs manual)
count by (issued) (netbird_group_info)

# Resource distribution by type across all groups
sum by (resource_type) (netbird_group_resources_by_type)

# Groups with no peers
netbird_group_peers_count == 0

# Groups with no resources
netbird_group_resources_count == 0

# Groups scrape error rate
rate(netbird_groups_scrape_errors_total[5m])
```

### User Queries
```promql
# Total number of users
netbird_users_total

# Users by role
sum by (role) (netbird_users_by_role)

# Users by status
sum by (status) (netbird_users_by_status)

# Service users vs regular users
netbird_users_service_users

# Blocked users
netbird_users_blocked

# Users by issuance type
sum by (issued) (netbird_users_by_issued)

# Users with restricted permissions
netbird_users_restricted

# Last login timestamp for each user
netbird_user_last_login_timestamp

# Auto groups assigned to each user
netbird_user_auto_groups_count

# User permissions by module and action
sum by (module, permission) (netbird_user_permissions)
```

### DNS Queries
```promql
# Total number of nameserver groups
netbird_dns_nameserver_groups_total

# Enabled vs disabled nameserver groups
netbird_dns_nameserver_groups_enabled

# Primary vs secondary nameserver groups
netbird_dns_nameserver_groups_primary

# Nameserver groups with the most domains
topk(5, netbird_dns_nameserver_group_domains_count)

# Nameserver groups with the most nameservers
topk(5, netbird_dns_nameservers_total)

# Nameserver distribution by type
sum by (ns_type) (netbird_dns_nameservers_by_type)

# Nameserver distribution by port
sum by (port) (netbird_dns_nameservers_by_port)

# Groups with DNS management disabled
netbird_dns_management_disabled_groups_count

# Average domains per nameserver group
avg(netbird_dns_nameserver_group_domains_count)

# Nameserver groups with no domains configured
netbird_dns_nameserver_group_domains_count == 0

# Total nameservers across all groups
sum(netbird_dns_nameservers_total)
```

### Network Queries
```promql
# Total number of networks
netbird_networks_total

# Networks with the most routers
topk(5, netbird_network_routers_count)

# Networks with the most resources
topk(5, netbird_network_resources_count)

# Networks with the most policies
topk(5, netbird_network_policies_count)

# Networks with the most routing peers
topk(5, netbird_network_routing_peers_count)

# Average routers per network
avg(netbird_network_routers_count)

# Average resources per network
avg(netbird_network_resources_count)

# Networks with no routers
netbird_network_routers_count == 0

# Networks with no resources
netbird_network_resources_count == 0

# Networks with no policies
netbird_network_policies_count == 0

# Total routers across all networks
sum(netbird_network_routers_count)

# Total resources across all networks
sum(netbird_network_resources_count)

# Networks scrape error rate
rate(netbird_networks_scrape_errors_total[5m])
```

## Grafana Dashboard

You can create a Grafana dashboard using these metrics. Example panels:

### Peer Panels
1. **Total Peers** - Single stat panel with `netbird_peers_total`
2. **Connection Status** - Pie chart with `netbird_peers_connected`
3. **OS Distribution** - Bar chart with `sum by (os) (netbird_peers_by_os)`
4. **Geographic Distribution** - World map with `sum by (country_code) (netbird_peers_by_country)`
5. **Last Seen Timeline** - Time series of peer activity

### Group Panels
1. **Total Groups** - Single stat panel with `netbird_groups_total`
2. **Group Sizes** - Bar chart with `netbird_group_peers_count` 
3. **Resource Distribution** - Pie chart with `sum by (resource_type) (netbird_group_resources_by_type)`
4. **Group Resource Counts** - Table with `netbird_group_resources_count`
5. **Groups by Creation Method** - Pie chart with `count by (issued) (netbird_group_info)`
6. **Empty Groups** - Table showing groups with `netbird_group_peers_count == 0`
7. **Group Performance** - Time series with `rate(netbird_groups_scrape_duration_seconds_sum[5m])`

### User Panels
1. **Total Users** - Single stat panel with `netbird_users_total`
2. **User Distribution** - Pie chart with `sum by (role) (netbird_users_by_role)`
3. **User Status** - Pie chart with `sum by (status) (netbird_users_by_status)`
4. **Service Users** - Pie chart with `netbird_users_service_users`
5. **Blocked Users** - Pie chart with `netbird_users_blocked`
6. **Users by Issuance Type** - Pie chart with `sum by (issued) (netbird_users_by_issued)`
7. **Users with Restricted Permissions** - Pie chart with `netbird_users_restricted`
8. **Last Login Timeline** - Time series of user activity
9. **Auto Groups** - Pie chart with `netbird_user_auto_groups_count`
10. **User Permissions** - Table with `netbird_user_permissions`

### DNS Panels
1. **Total Nameserver Groups** - Single stat panel with `netbird_dns_nameserver_groups_total`
2. **Nameserver Group Status** - Pie chart with `netbird_dns_nameserver_groups_enabled`
3. **Primary vs Secondary Groups** - Pie chart with `netbird_dns_nameserver_groups_primary`
4. **Domains per Group** - Bar chart with `netbird_dns_nameserver_group_domains_count`
5. **Nameservers per Group** - Bar chart with `netbird_dns_nameservers_total`
6. **Nameserver Types** - Pie chart with `sum by (ns_type) (netbird_dns_nameservers_by_type)`
7. **Nameserver Ports** - Pie chart with `sum by (port) (netbird_dns_nameservers_by_port)`
8. **DNS Management Disabled** - Single stat panel with `netbird_dns_management_disabled_groups_count`
9. **Groups without Domains** - Table showing groups with `netbird_dns_nameserver_group_domains_count == 0`

### Network Panels
1. **Total Networks** - Single stat panel with `netbird_networks_total`
2. **Network Routers** - Bar chart with `netbird_network_routers_count`
3. **Network Resources** - Bar chart with `netbird_network_resources_count`
4. **Network Policies** - Bar chart with `netbird_network_policies_count`
5. **Routing Peers** - Bar chart with `netbird_network_routing_peers_count`
6. **Network Information** - Table with `netbird_network_info`
7. **Networks without Routers** - Table showing networks with `netbird_network_routers_count == 0`
8. **Networks without Resources** - Table showing networks with `netbird_network_resources_count == 0`
9. **Network Performance** - Time series with `rate(netbird_networks_scrape_duration_seconds_sum[5m])`

### Performance Panels
1. **API Scrape Duration** - Time series with both peer and group scrape durations
2. **Error Rates** - Time series with error counters for both APIs

## Troubleshooting

### Common Issues

1. **Authentication errors**: Verify your `NETBIRD_API_TOKEN` is correct and has appropriate permissions
2. **Connection errors**: Check if the NetBird API URL is accessible from your network
3. **Missing metrics**: Ensure your NetBird account has peers registered

### Logs

Check logs for debugging:
```bash
# Docker Compose
docker-compose logs netbird-exporter

# Docker
docker logs netbird-exporter

# Binary
# Logs are output to stdout
```

### Enable Debug Logging

Set `LOG_LEVEL=debug` for more verbose output.

## Security Considerations

- Store your NetBird API token securely (use Docker secrets, Kubernetes secrets, etc.)
- Consider running the exporter in a private network
- Implement proper firewall rules to restrict access to the metrics endpoint
- Regularly rotate your API tokens

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Building from Source

```bash
go mod download
go build -o netbird-exporter
```

### Code Quality

This project uses several tools to maintain code quality:

#### Linting

Run linting checks:
```bash
make lint
```

This runs:
- `golangci-lint` - Comprehensive Go linting
- `go vet` - Go's built-in static analysis
- `gofmt` - Code formatting check

#### Formatting

Format code:
```bash
make fmt
```

#### Running All Checks

Run tests and linting together:
```bash
make check
```

#### Available Make Targets

```bash
make help
```

Shows all available targets including:
- `build` - Build the binary
- `test` - Run tests
- `lint` - Run linting checks
- `fmt` - Format code
- `check` - Run all checks (tests + linting)

### Continuous Integration

The project includes GitHub Actions workflows that automatically:
- Run linting checks on all pull requests
- Verify code formatting
- Run tests
- Check for security issues

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check` to ensure all tests and linting pass
5. Add tests if applicable
6. Submit a pull request

All pull requests must pass the CI checks including linting and tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## References

- [NetBird Peers API Documentation](https://docs.netbird.io/ipa/resources/peers)
- [NetBird Groups API Documentation](https://docs.netbird.io/ipa/resources/groups)
- [NetBird Users API Documentation](https://docs.netbird.io/ipa/resources/users)
- [NetBird Networks API Documentation](https://docs.netbird.io/ipa/resources/networks)
- [NetBird DNS API Documentation](https://docs.netbird.io/ipa/resources/dns)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [NetBird Official Documentation](https://docs.netbird.io/)