# Grafana Query Examples

This document provides additional Prometheus query examples for creating custom Grafana panels beyond the pre-built dashboard.

## Advanced Peer Queries

### Peer Health Score

Calculate a composite health score based on multiple factors:

```promql
# Peer health score (0-100)
(
  (netbird_peers_connected{connected="true"} * 40) +
  (netbird_peers_ssh_enabled{ssh_enabled="true"} * 20) +
  ((netbird_peers_login_expired{login_expired="false"} * 20)) +
  ((netbird_peers_approval_required{approval_required="false"} * 20))
) / netbird_peers_total * 100
```

### Peer Activity Timeline

Show peer activity over time:

```promql
# Recently active peers (last 1 hour)
count(time() - netbird_peer_last_seen_timestamp < 3600)

# Peer activity buckets
count(time() - netbird_peer_last_seen_timestamp < 300) # 5 minutes
count(time() - netbird_peer_last_seen_timestamp < 3600) # 1 hour
count(time() - netbird_peer_last_seen_timestamp < 86400) # 1 day
```

### Geographic Distribution Analysis

```promql
# Top 10 countries by peer count
topk(10, sum by (country_code) (netbird_peers_by_country))

# Cities with most peers
topk(10, sum by (city_name) (netbird_peers_by_country))

# Geographic diversity score
count(count by (country_code) (netbird_peers_by_country))
```

## User Analytics

### User Engagement Metrics

```promql
# Active users (logged in within last 7 days)
count(time() - netbird_user_last_login_timestamp < 604800)

# User login frequency
histogram_quantile(0.5,
  rate(time() - netbird_user_last_login_timestamp[7d])
)

# Service user ratio
netbird_users_service_users{is_service_user="true"} / netbird_users_total
```

### Permission Analysis

```promql
# Users with admin permissions
count by (module) (netbird_user_permissions{permission="admin", value="1"})

# Permission distribution
sum by (permission) (netbird_user_permissions{value="1"})

# Restricted users percentage
netbird_users_restricted{is_restricted="true"} / netbird_users_total * 100
```

## Group Efficiency Metrics

### Group Utilization

```promql
# Average group utilization
avg(netbird_group_peers_count / (netbird_group_peers_count + netbird_group_resources_count))

# Groups with optimal size (5-20 peers)
count(netbird_group_peers_count >= 5 and netbird_group_peers_count <= 20)

# Oversized groups (>50 peers)
count(netbird_group_peers_count > 50)
```

### Resource Efficiency

```promql
# Resource-to-peer ratio
netbird_group_resources_count / netbird_group_peers_count

# Groups with unused resources
count(netbird_group_resources_count > 0 and netbird_group_peers_count == 0)

# Resource type diversity per group
count by (group_id) (netbird_group_resources_by_type > 0)
```

## DNS Performance Metrics

### DNS Configuration Health

```promql
# DNS group coverage
(netbird_dns_nameserver_groups_total -
 netbird_dns_management_disabled_groups_count) /
 netbird_dns_nameserver_groups_total * 100

# Average domains per active group
avg(netbird_dns_nameserver_group_domains_count > 0)

# Nameserver redundancy ratio
avg(netbird_dns_nameservers_total)
```

### DNS Load Distribution

```promql
# Port usage distribution
sum by (port) (netbird_dns_nameservers_by_port) /
sum(netbird_dns_nameservers_by_port) * 100

# Protocol preference
netbird_dns_nameservers_by_type{ns_type="UDP"} /
sum(netbird_dns_nameservers_by_type) * 100
```

## Network Topology Analysis

### Network Complexity

```promql
# Network complexity score
(netbird_network_routers_count * 0.3 +
 netbird_network_resources_count * 0.4 +
 netbird_network_policies_count * 0.3)

# Hub networks (high routing peer count)
topk(5, netbird_network_routing_peers_count)

# Network efficiency ratio
netbird_network_resources_count / netbird_network_routers_count
```

### Routing Analysis

```promql
# Total routing capacity
sum(netbird_network_routing_peers_count)

# Networks without routing
count(netbird_network_routing_peers_count == 0)

# Average policies per network
avg(netbird_network_policies_count)
```

## Performance Monitoring

### API Performance Trends

```promql
# 95th percentile response time trend
histogram_quantile(0.95,
  rate(netbird_users_scrape_duration_seconds_bucket[5m])
)

# Error rate by component
rate(netbird_users_scrape_errors_total[5m]) * 100
rate(netbird_groups_scrape_errors_total[5m]) * 100
rate(netbird_networks_scrape_errors_total[5m]) * 100
```

### System Health Score

```promql
# Overall system health (0-100)
(
  (1 - rate(netbird_users_scrape_errors_total[5m])) * 25 +
  (1 - rate(netbird_groups_scrape_errors_total[5m])) * 25 +
  (1 - rate(netbird_networks_scrape_errors_total[5m])) * 25 +
  (netbird_peers_connected{connected="true"} / netbird_peers_total) * 25
) * 100
```

## Custom Panel Types

### Heatmap Panels

Create heatmaps for:

- Peer activity by hour of day
- User login patterns
- API response times by endpoint

### Graph Panels with Annotations

Add annotations for:

- System maintenance windows
- Configuration changes
- Incident markers

### Table Panels with Conditional Formatting

Color-code tables based on:

- Connection status (green=connected, red=disconnected)
- Performance thresholds
- Resource utilization levels

### Gauge Panels

Create gauges for:

- System health scores (0-100)
- Resource utilization percentages
- SLA compliance metrics

## Alerting Rules

### Critical Alerts

```yaml
# High disconnected peer rate
alert: HighPeerDisconnectionRate
expr: (netbird_peers_connected{connected="false"} / netbird_peers_total) > 0.2
for: 5m
labels:
  severity: critical
annotations:
  summary: "{{ $value }}% of peers are disconnected"

# API errors
alert: APIErrorRate
expr: rate(netbird_users_scrape_errors_total[5m]) > 0.1
for: 2m
labels:
  severity: warning
annotations:
  summary: "High API error rate: {{ $value }} errors/sec"
```

### Warning Alerts

```yaml
# No recent peer activity
alert: StalePolatePeers
expr: count(time() - netbird_peer_last_seen_timestamp > 86400) > 0
for: 1h
labels:
  severity: warning
annotations:
  summary: "{{ $value }} peers haven't been seen in 24+ hours"

# Large groups
alert: OversizedGroups
expr: netbird_group_peers_count > 100
for: 15m
labels:
  severity: warning
annotations:
  summary: "Group {{ $labels.group_name }} has {{ $value }} peers"
```

## Tips for Custom Dashboards

1. **Use Variables**: Create dashboard variables for filtering by group, network, or time range
2. **Set Thresholds**: Use appropriate color coding for metrics (green/yellow/red)
3. **Add Links**: Link panels to detailed views or external documentation
4. **Optimize Queries**: Use recording rules for complex calculations
5. **Consider Mobile**: Ensure dashboards work well on mobile devices
6. **Add Context**: Include helpful text panels with metric explanations

For more query examples and Grafana configuration tips, see the [official Prometheus documentation](https://prometheus.io/docs/prometheus/latest/querying/basics/) and [Grafana documentation](https://grafana.com/docs/grafana/latest/panels/).
