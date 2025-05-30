---
layout: default
title: Docker Compose
parent: Installation
nav_order: 1
---

# Docker Compose Installation
{: .no_toc }

The easiest way to get started with NetBird API Exporter for development and testing.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Overview

Docker Compose provides the simplest installation method, perfect for:
- **Development and testing environments**
- **Quick proof-of-concept setups**
- **Local monitoring setups**
- **Learning and experimentation**

## Prerequisites

Before starting, ensure you have:

- **Docker 20.10+** and **Docker Compose 2.0+** installed
- **NetBird API token** ([get one here](../getting-started/authentication))
- **Port 8080** available (or choose a different port)
- **Outbound internet access** to reach NetBird API

## Quick Start

### Step 1: Clone the Repository

```bash
git clone https://github.com/matanbaruch/netbird-api-exporter.git
cd netbird-api-exporter
```

### Step 2: Configure Environment

Create your environment file:

```bash
cp env.example .env
```

Edit the `.env` file with your settings:

```bash
# NetBird API Configuration
NETBIRD_API_TOKEN=nb_api_your_token_here
NETBIRD_API_URL=https://api.netbird.io

# Exporter Configuration
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=info
```

{: .important }
> **Required**: You must set your `NETBIRD_API_TOKEN`. Get your token from the [authentication guide](../getting-started/authentication).

### Step 3: Start the Exporter

```bash
docker-compose up -d
```

This will:
1. Pull the latest NetBird API Exporter image
2. Start the container in detached mode
3. Expose metrics on port 8080

### Step 4: Verify Installation

Check that the exporter is running:

```bash
# Health check
curl http://localhost:8080/health

# View metrics
curl http://localhost:8080/metrics
```

## Configuration Options

### Environment Variables

The Docker Compose setup supports all standard configuration options:

| Variable | Default | Description |
|----------|---------|-------------|
| `NETBIRD_API_TOKEN` | **Required** | Your NetBird API authentication token |
| `NETBIRD_API_URL` | `https://api.netbird.io` | NetBird API base URL |
| `LISTEN_ADDRESS` | `:8080` | Address and port for the exporter to listen on |
| `METRICS_PATH` | `/metrics` | Path where Prometheus metrics are exposed |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

### Custom Port

To use a different port, update your `.env` file:

```bash
LISTEN_ADDRESS=:9090
```

And modify the `docker-compose.yml` ports mapping:

```yaml
services:
  netbird-api-exporter:
    # ... other configuration
    ports:
      - "9090:9090"  # Change both sides to match LISTEN_ADDRESS
```

### Custom API URL

For self-hosted NetBird instances:

```bash
NETBIRD_API_URL=https://your-netbird-api.example.com
```

## Docker Compose Files

### Basic Configuration

The included `docker-compose.yml` provides a basic setup:

```yaml
version: '3.8'

services:
  netbird-api-exporter:
    image: ghcr.io/matanbaruch/netbird-api-exporter:latest
    container_name: netbird-api-exporter
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### Production Configuration

For production-like setups, create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  netbird-api-exporter:
    image: ghcr.io/matanbaruch/netbird-api-exporter:latest
    container_name: netbird-api-exporter
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: always
    
    # Resource limits
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: '0.2'
        reservations:
          memory: 64M
          cpus: '0.1'
    
    # Security settings
    read_only: true
    user: "65534:65534"  # nobody user
    cap_drop:
      - ALL
    security_opt:
      - no-new-privileges:true
    
    # Health checks
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    
    # Logging
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  # Optional: Prometheus for testing
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
    restart: unless-stopped
```

## Complete Setup with Prometheus

For a complete monitoring stack, use this enhanced configuration:

### Step 1: Create Prometheus Configuration

Create `prometheus.yml`:

```yaml
global:
  scrape_interval: 30s
  evaluation_interval: 30s

scrape_configs:
  - job_name: 'netbird-api-exporter'
    static_configs:
      - targets: ['netbird-api-exporter:8080']
    scrape_interval: 30s
    metrics_path: /metrics
    
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

### Step 2: Create Complete Docker Compose

Create `docker-compose.monitoring.yml`:

```yaml
version: '3.8'

services:
  netbird-api-exporter:
    image: ghcr.io/matanbaruch/netbird-api-exporter:latest
    container_name: netbird-api-exporter
    env_file:
      - .env
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - monitoring

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
    restart: unless-stopped
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    restart: unless-stopped
    networks:
      - monitoring

volumes:
  prometheus_data:
  grafana_data:

networks:
  monitoring:
    driver: bridge
```

### Step 3: Start the Complete Stack

```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

Access the services:
- **NetBird Exporter**: http://localhost:8080
- **Prometheus**: http://localhost:9090  
- **Grafana**: http://localhost:3000 (admin/admin)

## Management Commands

### View Logs

```bash
# All logs
docker-compose logs

# Follow logs
docker-compose logs -f

# Specific service logs
docker-compose logs netbird-api-exporter
```

### Restart Services

```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart netbird-api-exporter
```

### Update to Latest Version

```bash
# Pull latest images
docker-compose pull

# Restart with new images
docker-compose up -d
```

### Stop and Clean Up

```bash
# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Remove everything including images
docker-compose down -v --rmi all
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```
Error: bind: address already in use
```

**Solution**: Change the port in `.env` and `docker-compose.yml`:
```bash
LISTEN_ADDRESS=:9090
```

#### 2. Permission Denied
```
Error: permission denied while trying to connect to Docker daemon
```

**Solution**: Add your user to the docker group:
```bash
sudo usermod -aG docker $USER
# Log out and back in
```

#### 3. API Token Not Working
```
Error: 401 Unauthorized
```

**Solution**: Verify your token in `.env`:
- Check for extra spaces or newlines
- Ensure token has proper permissions
- Test token manually with curl

### Debug Mode

Enable debug logging:

```bash
# Add to .env
LOG_LEVEL=debug
```

Restart the container:
```bash
docker-compose restart netbird-api-exporter
```

View debug logs:
```bash
docker-compose logs -f netbird-api-exporter
```

### Health Checks

Monitor container health:

```bash
# Check container status
docker-compose ps

# View health check logs
docker inspect netbird-api-exporter | grep Health -A 10
```

## Security Considerations

### Environment File Security

Protect your `.env` file:

```bash
# Set restrictive permissions
chmod 600 .env

# Add to .gitignore
echo ".env" >> .gitignore
```

### Network Security

For production deployments:

1. **Use Docker networks** to isolate services
2. **Don't expose unnecessary ports** to the host
3. **Use reverse proxy** (nginx, traefik) for HTTPS
4. **Implement firewall rules** to restrict access

## Next Steps

Once your Docker Compose setup is running:

1. **[Configure Prometheus](../usage/prometheus-setup)** to scrape your exporter
2. **[Create Grafana dashboards](../usage/grafana-dashboards)** for visualization  
3. **[Explore available metrics](../reference/metrics)** and queries
4. **[Set up alerting](../usage/alerting)** for important events

For production deployments, consider:
- **[Helm installation](helm)** for Kubernetes environments
- **[systemd installation](systemd)** for traditional Linux servers 