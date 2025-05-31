---
layout: default
title: Getting Started
nav_order: 2
has_children: true
---

# Getting Started
{: .no_toc }

This guide will help you get the NetBird API Exporter up and running in your environment.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Prerequisites

Before installing the NetBird API Exporter, ensure you have:

1. **NetBird Account**: An active NetBird account with access to the dashboard
2. **API Token**: A valid NetBird API token with appropriate permissions
3. **Prometheus Server**: A Prometheus server or compatible metrics collection system
4. **Deployment Environment**: One of the following:
   - Docker and Docker Compose
   - Kubernetes cluster (for Helm deployment)
   - Linux server (for systemd deployment)
   - Go development environment (for binary deployment)

## System Requirements

### Minimum Requirements
- **CPU**: 100m (0.1 CPU core)
- **Memory**: 64MB RAM
- **Disk**: 10MB storage
- **Network**: Outbound HTTPS access to NetBird API

### Recommended Requirements
- **CPU**: 200m (0.2 CPU core)
- **Memory**: 128MB RAM
- **Disk**: 50MB storage
- **Network**: Stable internet connection

## Supported Platforms

The NetBird API Exporter supports the following platforms:

### Container Platforms
- Docker 20.10+
- Docker Compose 2.0+
- Kubernetes 1.19+
- Podman 3.0+

### Operating Systems
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Helm Versions
- Helm 3.0+

## Next Steps

1. **[Get your API token](getting-started/authentication)** - Learn how to create and configure your NetBird API token
2. **Choose your installation method**:
   - [Docker Compose](installation/docker-compose) - Easiest for development and testing
   - [Docker](installation/docker) - Best for container environments
   - [Helm](installation/helm) - Ideal for Kubernetes clusters
   - [systemd](installation/systemd) - Perfect for Linux servers
   - [Binary](installation/binary) - For custom setups

3. **[Configure Prometheus](usage/prometheus-setup)** - Set up Prometheus to scrape metrics
4. **[Create dashboards](usage/grafana-dashboards)** - Build monitoring dashboards in Grafana

## Quick Verification

Once you have the exporter running, you can quickly verify it's working:

### Health Check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "0.1.0"
}
```

### Metrics Endpoint
```bash
curl http://localhost:8080/metrics
```

You should see Prometheus metrics output starting with:
```
# HELP netbird_peers_total Total number of NetBird peers
# TYPE netbird_peers_total gauge
netbird_peers_total 42
```

## Troubleshooting

If you encounter issues during setup:

1. **Check the logs** for error messages
2. **Verify your API token** is correct and has proper permissions
3. **Ensure network connectivity** to NetBird API
4. **Review the [troubleshooting guide](reference/troubleshooting)**

## Getting Help

- 📖 **Documentation**: Browse the complete guide in the sidebar
- 🐛 **Issues**: Report problems on [GitHub Issues](https://github.com/matanbaruch/netbird-api-exporter/issues)
- 💡 **Discussions**: Ask questions in [GitHub Discussions](https://github.com/matanbaruch/netbird-api-exporter/discussions) 
