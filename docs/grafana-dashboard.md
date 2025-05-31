# Grafana Dashboard

This document describes how to use the NetBird API Exporter Grafana dashboard for monitoring your NetBird infrastructure.

## Overview

The NetBird API Exporter Grafana dashboard provides comprehensive visualizations of your NetBird deployment, including:

- **Overview**: High-level statistics for peers, users, groups, and networks
- **Peers**: Connection status, operating system distribution, geographic distribution
- **Users**: Role breakdown, status overview, service vs regular users
- **Groups**: Peer and resource counts per group
- **DNS**: Nameserver group configurations and types
- **Networks**: Network information and resource counts
- **Performance**: API scrape duration and error rates

## Installation

### Method 1: Import Dashboard JSON

1. **Download the dashboard file**: Get the `grafana-dashboard.json` file from the repository
2. **Open Grafana**: Navigate to your Grafana instance
3. **Import Dashboard**:
   - Go to `Dashboards` → `Import` (or press `+` → `Import`)
   - Upload the `grafana-dashboard.json` file or paste its contents
   - Click `Import`

### Method 2: Manual Setup

If you prefer to create the dashboard manually, you can use the dashboard JSON as a reference for the queries and panel configurations.

## Prerequisites

Before using the dashboard, ensure you have:

1. **Prometheus**: A running Prometheus instance configured to scrape the NetBird API Exporter
2. **Grafana**: A Grafana instance with Prometheus configured as a data source
3. **NetBird API Exporter**: The exporter running and exposing metrics

### Prometheus Configuration

Add the NetBird API Exporter to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: "netbird-api-exporter"
    static_configs:
      - targets: ["localhost:8080"] # Adjust to your exporter's address
    scrape_interval: 30s
    metrics_path: /metrics
```

### Grafana Data Source

Ensure your Prometheus data source is configured in Grafana:

1. Go to `Configuration` → `Data Sources`
2. Add or verify your Prometheus data source
3. The dashboard uses a variable `${DS_PROMETHEUS}` which should automatically detect your Prometheus data source

## Dashboard Sections

### Overview

- **Total Peers**: Number of NetBird peers
- **Total Users**: Number of NetBird users
- **Total Groups**: Number of NetBird groups
- **Total Networks**: Number of NetBird networks

### Peers

- **Connection Status**: Pie chart showing connected vs disconnected peers
- **Operating System Distribution**: Breakdown of peers by OS
- **Geographic Distribution**: Table showing peers by country and city

### Users

- **Role Distribution**: Users by role (admin, user, etc.)
- **Status Overview**: Users by status (active, invited, etc.)
- **User Types**: Service users vs regular users

### Groups

- **Peer Counts**: Table showing number of peers per group
- **Resource Counts**: Table showing number of resources per group

### DNS

- **Nameserver Groups**: Total count and enabled/disabled status
- **Nameserver Types**: Distribution by protocol type (UDP/TCP)

### Networks

- **Network Overview**: Table with network information and descriptions

### Performance & Errors

- **API Scrape Duration**: 95th percentile response times for API calls
- **Error Rates**: Rate of scraping errors by component

## Customization

### Time Range

The dashboard defaults to showing the last 1 hour of data. You can adjust this using the time picker in the top-right corner.

### Refresh Rate

The dashboard automatically refreshes every 30 seconds. You can change this in the refresh dropdown.

### Panel Modifications

You can customize any panel by:

1. Clicking the panel title
2. Selecting `Edit`
3. Modifying queries, visualizations, or thresholds as needed

### Additional Metrics

The NetBird API Exporter provides additional metrics that aren't included in the default dashboard. You can create custom panels using these metrics:

- `netbird_peer_last_seen_timestamp`: Individual peer last seen times
- `netbird_user_last_login_timestamp`: User login timestamps
- `netbird_peer_accessible_peers_count`: Accessible peer counts
- `netbird_dns_nameserver_group_domains_count`: Domain counts per nameserver group
- And more...

## Troubleshooting

### No Data Showing

1. **Check Prometheus**: Verify the exporter is being scraped successfully
2. **Check Data Source**: Ensure Grafana can connect to Prometheus
3. **Check Time Range**: Make sure the time range includes when the exporter was running
4. **Check Metrics**: Verify metrics are being exposed at the `/metrics` endpoint

### Missing Panels

1. **Check Metric Names**: Ensure the NetBird API Exporter is running the latest version
2. **Check Labels**: Some panels may not show data if certain labels are missing

### Performance Issues

1. **Reduce Time Range**: Shorter time ranges load faster
2. **Increase Scrape Interval**: Reduce the frequency of Prometheus scrapes if needed
3. **Optimize Queries**: Some complex queries may need optimization for large deployments

## Alerting

You can create Grafana alerts based on the metrics shown in this dashboard. Common alerting scenarios include:

- **High Error Rates**: Alert when API scrape errors exceed a threshold
- **Peer Connection Issues**: Alert when too many peers are disconnected
- **Performance Degradation**: Alert when API response times are too high

To create alerts:

1. Edit a panel that contains the metric you want to alert on
2. Go to the `Alert` tab
3. Configure alert conditions and notification channels

## Contributing

If you have suggestions for improving the dashboard or want to add new panels, please:

1. Fork the repository
2. Modify the `grafana-dashboard.json` file
3. Test your changes
4. Submit a pull request

For more information about the available metrics, see the [NetBird API Exporter documentation](../README.md).

For advanced query examples and custom panel ideas, see the [Grafana Query Examples](../examples/grafana-queries.md).
