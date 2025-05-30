---
layout: default
title: Docker
parent: Installation
nav_order: 2
---

# Docker Installation
{: .no_toc }

Deploy NetBird API Exporter using Docker containers for flexible containerized environments.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Overview

Docker installation provides a containerized deployment method, perfect for:
- **Container-based infrastructure**
- **CI/CD pipelines**
- **Cloud platforms**
- **Development environments**
- **Microservices architectures**

## Prerequisites

Before starting, ensure you have:

- **Docker 20.10+** installed
- **NetBird API token** ([get one here](../getting-started/authentication))
- **Port 8080** available (or choose a different port)
- **Outbound internet access** to reach NetBird API

## Quick Start

### Step 1: Pull the Image

```bash
# Pull the latest image from GitHub Container Registry
docker pull ghcr.io/matanbaruch/netbird-api-exporter:latest
```

### Step 2: Run the Container

```bash
# Basic run command
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  -e NETBIRD_API_TOKEN="nb_api_your_token_here" \
  -e NETBIRD_API_URL="https://api.netbird.io" \
  --restart unless-stopped \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

{: .important }
> **Required**: Replace `nb_api_your_token_here` with your actual NetBird API token from the [authentication guide](../getting-started/authentication).

### Step 3: Verify Installation

```bash
# Check container status
docker ps

# Test health endpoint
curl http://localhost:8080/health

# Test metrics endpoint
curl http://localhost:8080/metrics
```

## Configuration Options

### Environment Variables

Configure the exporter using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `NETBIRD_API_TOKEN` | **Required** | Your NetBird API authentication token |
| `NETBIRD_API_URL` | `https://api.netbird.io` | NetBird API base URL |
| `LISTEN_ADDRESS` | `:8080` | Address and port for the exporter to listen on |
| `METRICS_PATH` | `/metrics` | Path where Prometheus metrics are exposed |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

### Complete Example

```bash
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  -e NETBIRD_API_TOKEN="nb_api_your_token_here" \
  -e NETBIRD_API_URL="https://api.netbird.io" \
  -e LISTEN_ADDRESS=":8080" \
  -e METRICS_PATH="/metrics" \
  -e LOG_LEVEL="info" \
  --restart unless-stopped \
  --memory="128m" \
  --cpus="0.2" \
  --read-only \
  --user 65534:65534 \
  --cap-drop ALL \
  --security-opt no-new-privileges:true \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

## Production Deployment

### Using Environment File

Create an environment file for better configuration management:

```bash
# Create environment file
cat > netbird-exporter.env << 'EOF'
NETBIRD_API_TOKEN=nb_api_your_token_here
NETBIRD_API_URL=https://api.netbird.io
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=info
EOF

# Set secure permissions
chmod 600 netbird-exporter.env

# Run with environment file
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  --env-file netbird-exporter.env \
  --restart unless-stopped \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

### Security Hardened Deployment

```bash
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  --env-file netbird-exporter.env \
  --restart unless-stopped \
  \
  # Resource limits
  --memory="128m" \
  --memory-swap="128m" \
  --cpus="0.2" \
  --pids-limit 100 \
  \
  # Security settings
  --read-only \
  --tmpfs /tmp:noexec,nosuid,size=10m \
  --user 65534:65534 \
  --cap-drop ALL \
  --security-opt no-new-privileges:true \
  --security-opt seccomp=unconfined \
  \
  # Network security
  --network bridge \
  \
  # Health check
  --health-cmd="curl -f http://localhost:8080/health || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  --health-start-period=40s \
  \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

### Custom Network

For better network isolation:

```bash
# Create custom network
docker network create netbird-monitoring

# Run container in custom network
docker run -d \
  --name netbird-api-exporter \
  --network netbird-monitoring \
  -p 8080:8080 \
  --env-file netbird-exporter.env \
  --restart unless-stopped \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

## Docker with Prometheus

### Complete Monitoring Stack

Create a simple monitoring setup with Prometheus:

```bash
# Create network
docker network create monitoring

# Create Prometheus configuration
mkdir -p prometheus-config
cat > prometheus-config/prometheus.yml << 'EOF'
global:
  scrape_interval: 30s
  evaluation_interval: 30s

scrape_configs:
  - job_name: 'netbird-api-exporter'
    static_configs:
      - targets: ['netbird-api-exporter:8080']
    scrape_interval: 30s
    metrics_path: /metrics
EOF

# Run NetBird exporter
docker run -d \
  --name netbird-api-exporter \
  --network monitoring \
  --env-file netbird-exporter.env \
  --restart unless-stopped \
  ghcr.io/matanbaruch/netbird-api-exporter:latest

# Run Prometheus
docker run -d \
  --name prometheus \
  --network monitoring \
  -p 9090:9090 \
  -v $(pwd)/prometheus-config:/etc/prometheus:ro \
  --restart unless-stopped \
  prom/prometheus:latest \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/prometheus \
  --web.console.libraries=/etc/prometheus/console_libraries \
  --web.console.templates=/etc/prometheus/consoles

# Run Grafana
docker run -d \
  --name grafana \
  --network monitoring \
  -p 3000:3000 \
  -e GF_SECURITY_ADMIN_PASSWORD=admin \
  --restart unless-stopped \
  grafana/grafana:latest
```

Access the services:
- **NetBird Exporter**: http://localhost:8080
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

## Management Commands

### Container Management

```bash
# Start container
docker start netbird-api-exporter

# Stop container
docker stop netbird-api-exporter

# Restart container
docker restart netbird-api-exporter

# Remove container
docker rm netbird-api-exporter

# View container status
docker ps -f name=netbird-api-exporter

# Inspect container configuration
docker inspect netbird-api-exporter
```

### Logs and Monitoring

```bash
# View logs
docker logs netbird-api-exporter

# Follow logs
docker logs -f netbird-api-exporter

# View recent logs
docker logs --tail 50 netbird-api-exporter

# View logs with timestamps
docker logs -t netbird-api-exporter

# Check container stats
docker stats netbird-api-exporter
```

### Updates

```bash
# Pull latest image
docker pull ghcr.io/matanbaruch/netbird-api-exporter:latest

# Stop and remove old container
docker stop netbird-api-exporter
docker rm netbird-api-exporter

# Run new container
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  --env-file netbird-exporter.env \
  --restart unless-stopped \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

## Building Custom Images

### Build from Source

```bash
# Clone repository
git clone https://github.com/matanbaruch/netbird-api-exporter.git
cd netbird-api-exporter

# Build image
docker build -t netbird-api-exporter:local .

# Run locally built image
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  --env-file netbird-exporter.env \
  netbird-api-exporter:local
```

### Custom Dockerfile

Create a custom Dockerfile for specific needs:

```dockerfile
FROM ghcr.io/matanbaruch/netbird-api-exporter:latest

# Add custom CA certificates
COPY custom-ca.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

# Set custom user
USER 65534:65534

# Custom health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

## Troubleshooting

### Common Issues

#### 1. Container Won't Start

```bash
# Check container logs
docker logs netbird-api-exporter

# Check container configuration
docker inspect netbird-api-exporter

# Common causes:
# - Invalid API token
# - Port already in use
# - Resource constraints
```

#### 2. Permission Denied Errors

```bash
# Check if running as non-root user
docker exec netbird-api-exporter id

# Verify file permissions in container
docker exec netbird-api-exporter ls -la /app

# Use proper user flag
docker run --user 65534:65534 ...
```

#### 3. Network Connectivity Issues

```bash
# Test DNS resolution in container
docker exec netbird-api-exporter nslookup api.netbird.io

# Test API connectivity
docker exec netbird-api-exporter curl -v https://api.netbird.io/api/health

# Check container network
docker network ls
docker network inspect bridge
```

### Debug Mode

Enable debug logging:

```bash
# Run with debug logging
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  -e NETBIRD_API_TOKEN="nb_api_your_token_here" \
  -e LOG_LEVEL="debug" \
  ghcr.io/matanbaruch/netbird-api-exporter:latest

# View debug logs
docker logs -f netbird-api-exporter
```

### Health Checks

Check container health:

```bash
# View health status
docker inspect --format='{{.State.Health.Status}}' netbird-api-exporter

# View health check logs
docker inspect --format='{{range .State.Health.Log}}{{.Output}}{{end}}' netbird-api-exporter

# Manual health check
docker exec netbird-api-exporter curl -f http://localhost:8080/health
```

## Security Considerations

### Image Security

- **Use official images** from trusted registries
- **Scan images** for vulnerabilities regularly
- **Keep images updated** with latest security patches
- **Use specific tags** instead of `latest` in production

### Runtime Security

```bash
# Run with security best practices
docker run -d \
  --name netbird-api-exporter \
  -p 8080:8080 \
  --env-file netbird-exporter.env \
  \
  # Read-only filesystem
  --read-only \
  --tmpfs /tmp:noexec,nosuid,size=10m \
  \
  # Drop all capabilities
  --cap-drop ALL \
  \
  # Run as non-root user
  --user 65534:65534 \
  \
  # Security options
  --security-opt no-new-privileges:true \
  --security-opt seccomp=unconfined \
  \
  # Resource limits
  --memory="128m" \
  --cpus="0.2" \
  --pids-limit 100 \
  \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

### Environment Security

```bash
# Secure environment file
chmod 600 netbird-exporter.env
chown root:root netbird-exporter.env

# Use Docker secrets (in Swarm mode)
echo "nb_api_your_token_here" | docker secret create netbird_token -

# Use in container
docker service create \
  --name netbird-api-exporter \
  --secret netbird_token \
  -e NETBIRD_API_TOKEN_FILE=/run/secrets/netbird_token \
  ghcr.io/matanbaruch/netbird-api-exporter:latest
```

## Next Steps

Once your Docker deployment is running:

1. **[Configure Prometheus](../usage/prometheus-setup)** to scrape your exporter
2. **[Create Grafana dashboards](../usage/grafana-dashboards)** for visualization
3. **[Set up alerting](../usage/alerting)** for important metrics
4. **[Explore metrics](../reference/metrics)** and custom queries

For production environments, consider:
- **[Helm installation](helm)** for Kubernetes orchestration
- **[systemd installation](systemd)** for traditional server deployment 