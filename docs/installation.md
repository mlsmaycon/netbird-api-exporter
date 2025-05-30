---
layout: default
title: Installation
nav_order: 3
has_children: true
---

# Installation
{: .no_toc }

Choose the installation method that best fits your environment and requirements.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Installation Methods

The NetBird API Exporter supports multiple deployment methods to fit different environments and use cases:

| Method | Best For | Complexity | Setup Time |
|--------|----------|------------|------------|
| [Docker Compose](installation/docker-compose) | Development & Testing | ⭐ Easy | 5 minutes |
| [Docker](installation/docker) | Container Environments | ⭐⭐ Medium | 10 minutes |
| [Helm](installation/helm) | Kubernetes Clusters | ⭐⭐ Medium | 15 minutes |
| [systemd](installation/systemd) | Linux Servers | ⭐⭐⭐ Advanced | 20 minutes |
| [Binary](installation/binary) | Custom Setups | ⭐⭐⭐ Advanced | 25 minutes |

## Quick Decision Guide

### Choose Docker Compose if:
- You're just getting started or testing the exporter
- You want the fastest setup with minimal configuration
- You're running on a development machine
- You prefer infrastructure-as-code with docker-compose.yml

### Choose Docker if:
- You're deploying to a single container host
- You need more control over container configuration
- You're integrating with existing Docker infrastructure
- You want a lightweight deployment without Compose

### Choose Helm if:
- You're running Kubernetes
- You want cloud-native deployment with proper scaling
- You need service discovery and load balancing
- You prefer GitOps workflows

### Choose systemd if:
- You're running on traditional Linux servers
- You want native OS integration and startup
- You need tight resource control and security
- You prefer traditional service management

### Choose Binary if:
- You have specific OS or architecture requirements
- You want maximum control over the deployment
- You need to embed it in other applications
- You're building custom automation around the exporter

## Prerequisites

Before proceeding with any installation method, ensure you have:

### Required
1. **NetBird API Token** - [Get your token here](getting-started/authentication)
2. **Network Access** - Outbound HTTPS to NetBird API (`api.netbird.io` or your custom API URL)
3. **Prometheus Server** - To scrape metrics from the exporter

### Method-Specific Requirements

#### Docker/Docker Compose
- Docker 20.10+ and Docker Compose 2.0+
- Port 8080 available (or custom port)

#### Helm/Kubernetes
- Kubernetes 1.19+
- Helm 3.0+
- kubectl configured to access your cluster

#### systemd/Linux
- Linux distribution with systemd
- sudo/root access for service installation
- Go 1.21+ (if building from source)

#### Binary
- Go 1.21+ development environment
- Git for cloning the repository

## Configuration Overview

All installation methods use the same environment variables for configuration:

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `NETBIRD_API_URL` | `https://api.netbird.io` | No | NetBird API base URL |
| `NETBIRD_API_TOKEN` | - | **Yes** | NetBird API authentication token |
| `LISTEN_ADDRESS` | `:8080` | No | Address and port to listen on |
| `METRICS_PATH` | `/metrics` | No | Path where metrics are exposed |
| `LOG_LEVEL` | `info` | No | Log level (debug, info, warn, error) |

{: .important }
> **Security Note**: Always store your `NETBIRD_API_TOKEN` securely using your platform's secret management system.

## Post-Installation Steps

After completing any installation:

1. **Verify the service is running**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **Check metrics are available**:
   ```bash
   curl http://localhost:8080/metrics
   ```

3. **Configure Prometheus** to scrape the exporter:
   ```yaml
   scrape_configs:
     - job_name: 'netbird-api-exporter'
       static_configs:
         - targets: ['localhost:8080']
   ```

4. **Set up monitoring dashboards** in Grafana

## Migration Between Methods

You can easily migrate between deployment methods:

### From Docker Compose to Kubernetes
1. Export your environment variables from `.env`
2. Create Kubernetes secrets with the same values
3. Deploy using Helm with the secrets

### From Binary to systemd
1. Copy your binary to `/usr/local/bin/`
2. Create systemd service file with your current environment
3. Enable and start the service

### From Any Method to Docker
1. Note your current configuration
2. Stop the existing deployment
3. Run Docker container with same environment variables

## Getting Help

If you encounter issues during installation:

1. **Check the logs** for error messages specific to your deployment method
2. **Review the [troubleshooting guide](reference/troubleshooting)**
3. **Verify your [API token configuration](getting-started/authentication)**
4. **Ask for help** in [GitHub Discussions](https://github.com/matanbaruch/netbird-api-exporter/discussions)

---

## Next Steps

Choose your preferred installation method and follow the detailed guide:

- **[Docker Compose Installation](installation/docker-compose)** - Quick and easy setup
- **[Docker Installation](installation/docker)** - Container deployment
- **[Helm Installation](installation/helm)** - Kubernetes deployment  
- **[systemd Installation](installation/systemd)** - Linux service
- **[Binary Installation](installation/binary)** - Build from source 