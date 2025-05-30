---
layout: default
title: Binary
parent: Installation
nav_order: 5
---

# Binary Installation
{: .no_toc }

Build and run NetBird API Exporter from source code for maximum customization and control.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Overview

Binary installation provides the most flexible deployment method, perfect for:
- **Custom builds and modifications**
- **Development and testing**
- **Specific OS or architecture requirements**
- **Embedding in other applications**
- **Maximum performance optimization**

## Prerequisites

Before starting, ensure you have:

- **Go 1.21+** development environment
- **Git** for cloning the repository
- **Make** (optional, for using Makefile)
- **NetBird API token** ([get one here](../getting-started/authentication))
- **Network access** to GitHub and Go module proxy

## Quick Start

### Step 1: Install Go

If you don't have Go installed:

```bash
# Linux
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# macOS with Homebrew
brew install go

# Windows with Chocolatey
choco install golang

# Verify installation
go version
```

### Step 2: Clone Repository

```bash
# Clone the repository
git clone https://github.com/matanbaruch/netbird-api-exporter.git
cd netbird-api-exporter

# Verify Go modules
go mod download
go mod verify
```

### Step 3: Build Binary

```bash
# Basic build
go build -o netbird-api-exporter .

# Or using Make
make build

# Verify binary
./netbird-api-exporter --version
```

### Step 4: Configure Environment

```bash
# Set required environment variables
export NETBIRD_API_TOKEN="nb_api_your_token_here"
export NETBIRD_API_URL="https://api.netbird.io"
export LISTEN_ADDRESS=":8080"
export METRICS_PATH="/metrics"
export LOG_LEVEL="info"
```

{: .important }
> **Required**: Replace `nb_api_your_token_here` with your actual NetBird API token from the [authentication guide](../getting-started/authentication).

### Step 5: Run the Exporter

```bash
# Run with environment variables
./netbird-api-exporter

# Or run with inline environment
NETBIRD_API_TOKEN="nb_api_your_token_here" ./netbird-api-exporter
```

### Step 6: Verify Installation

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test metrics endpoint
curl http://localhost:8080/metrics
```

## Build Options

### Production Build

Create an optimized production build:

```bash
# Production build with optimizations
go build -ldflags="-s -w" -o netbird-api-exporter .

# Or using Make
make build-prod

# Build with version information
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD)

go build \
  -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
  -o netbird-api-exporter .
```

### Cross-Platform Builds

Build for different platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o netbird-api-exporter-linux-amd64 .

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o netbird-api-exporter-linux-arm64 .

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o netbird-api-exporter-darwin-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o netbird-api-exporter-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o netbird-api-exporter-windows-amd64.exe .

# Build all platforms using Make
make build-all
```

### Static Binary

Create a static binary with no external dependencies:

```bash
# Static build for Linux
CGO_ENABLED=0 GOOS=linux go build \
  -a -installsuffix cgo \
  -ldflags="-s -w -extldflags '-static'" \
  -o netbird-api-exporter-static .

# Verify static linking
ldd netbird-api-exporter-static
# Should output: "not a dynamic executable"
```

## Development Setup

### Development Environment

Set up a complete development environment:

```bash
# Clone repository
git clone https://github.com/matanbaruch/netbird-api-exporter.git
cd netbird-api-exporter

# Install development dependencies
go mod download

# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Set up git hooks (optional)
cp scripts/pre-commit .git/hooks/
chmod +x .git/hooks/pre-commit
```

### Make Commands

The project includes a Makefile with common tasks:

```bash
# Show available commands
make help

# Build binary
make build

# Run tests
make test

# Run linting
make lint

# Format code
make fmt

# Run all checks (test + lint)
make check

# Clean build artifacts
make clean

# Build for all platforms
make build-all
```

### Configuration File

Create a configuration file for development:

```bash
# Create config file
cat > config.env << 'EOF'
NETBIRD_API_TOKEN=nb_api_your_token_here
NETBIRD_API_URL=https://api.netbird.io
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=debug
EOF

# Load and run
set -a; source config.env; set +a
./netbird-api-exporter
```

## Advanced Usage

### Custom Build Tags

Use build tags for conditional compilation:

```bash
# Build with debug features
go build -tags debug -o netbird-api-exporter .

# Build for production
go build -tags production -o netbird-api-exporter .

# Build with specific features
go build -tags "metrics,health,pprof" -o netbird-api-exporter .
```

### Memory Optimization

Build with memory optimizations:

```bash
# Build with garbage collector optimizations
GOGC=off go build -ldflags="-s -w" -o netbird-api-exporter .

# Build with memory allocator optimizations
go build -ldflags="-s -w" -gcflags="all=-N -l" -o netbird-api-exporter .
```

### Profiling Build

Build with profiling support:

```bash
# Build with race detection (development only)
go build -race -o netbird-api-exporter .

# Build with profiling
go build -ldflags="-s -w" -o netbird-api-exporter .

# Run with profiling
./netbird-api-exporter --enable-pprof
```

## Testing

### Unit Tests

Run comprehensive tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific tests
go test -run TestPeerExporter ./pkg/exporters/

# Run tests with race detection
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

### Integration Tests

```bash
# Run integration tests (requires API token)
export NETBIRD_API_TOKEN="nb_api_your_token_here"
go test -tags integration ./tests/integration/

# Run with custom API URL
export NETBIRD_API_URL="https://your-api.example.com"
go test -tags integration ./tests/integration/
```

### Load Testing

```bash
# Build optimized binary for load testing
go build -ldflags="-s -w" -o netbird-api-exporter .

# Run load test with hey
hey -n 1000 -c 10 http://localhost:8080/metrics

# Run with custom duration
hey -z 60s -c 5 http://localhost:8080/health
```

## Deployment

### Production Deployment

Prepare for production deployment:

```bash
# Build production binary
make build-prod

# Create deployment directory
sudo mkdir -p /opt/netbird-api-exporter
sudo cp netbird-api-exporter /opt/netbird-api-exporter/

# Create configuration
sudo tee /opt/netbird-api-exporter/config.env > /dev/null << 'EOF'
NETBIRD_API_TOKEN=nb_api_your_token_here
NETBIRD_API_URL=https://api.netbird.io
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=info
EOF

# Set permissions
sudo chmod 755 /opt/netbird-api-exporter/netbird-api-exporter
sudo chmod 600 /opt/netbird-api-exporter/config.env
```

### Startup Script

Create a startup script:

```bash
sudo tee /opt/netbird-api-exporter/start.sh > /dev/null << 'EOF'
#!/bin/bash

# Load configuration
set -a
source /opt/netbird-api-exporter/config.env
set +a

# Start exporter
exec /opt/netbird-api-exporter/netbird-api-exporter
EOF

sudo chmod +x /opt/netbird-api-exporter/start.sh
```

### Process Management

Use a process manager for production:

```bash
# Using supervisor
sudo tee /etc/supervisor/conf.d/netbird-api-exporter.conf > /dev/null << 'EOF'
[program:netbird-api-exporter]
command=/opt/netbird-api-exporter/start.sh
directory=/opt/netbird-api-exporter
user=nobody
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/netbird-api-exporter.log
EOF

# Reload supervisor
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start netbird-api-exporter
```

## Troubleshooting

### Build Issues

#### 1. Go Version Mismatch

```bash
# Check Go version
go version

# Update Go if needed
sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
```

#### 2. Module Download Issues

```bash
# Clear module cache
go clean -modcache

# Retry download
go mod download

# Use direct module download
GOPROXY=direct go mod download
```

#### 3. Build Failures

```bash
# Clean build cache
go clean -cache

# Rebuild all dependencies
go build -a .

# Check for missing dependencies
go mod tidy
```

### Runtime Issues

#### 1. Binary Won't Start

```bash
# Check binary permissions
ls -la netbird-api-exporter

# Make executable
chmod +x netbird-api-exporter

# Check for missing libraries
ldd netbird-api-exporter
```

#### 2. Environment Variables

```bash
# List current environment
env | grep NETBIRD

# Test environment loading
set -a; source config.env; set +a; env | grep NETBIRD
```

#### 3. Network Issues

```bash
# Test API connectivity
curl -H "Authorization: Bearer $NETBIRD_API_TOKEN" \
     https://api.netbird.io/api/health

# Check DNS resolution
nslookup api.netbird.io

# Test port binding
netstat -tlnp | grep :8080
```

### Performance Issues

#### 1. Memory Usage

```bash
# Monitor memory usage
top -p $(pgrep netbird-api-exporter)

# Enable memory profiling
./netbird-api-exporter --enable-pprof &
go tool pprof http://localhost:8080/debug/pprof/heap
```

#### 2. CPU Usage

```bash
# Monitor CPU usage
htop -p $(pgrep netbird-api-exporter)

# Enable CPU profiling
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
```

## Development Workflow

### Code Quality

Maintain code quality with these tools:

```bash
# Format code
go fmt ./...
goimports -w .

# Lint code
golangci-lint run

# Run security checks
gosec ./...

# Check for vulnerabilities
go list -json -m all | nancy sleuth

# Generate documentation
godoc -http=:6060
```

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/new-exporter

# Make changes and test
make test
make lint

# Commit changes
git add .
git commit -m "Add new exporter functionality"

# Push and create pull request
git push origin feature/new-exporter
```

### Release Process

```bash
# Tag release
git tag -a v0.2.0 -m "Release version 0.2.0"
git push origin v0.2.0

# Build release binaries
make build-all

# Create release artifacts
tar -czf netbird-api-exporter-v0.2.0-linux-amd64.tar.gz netbird-api-exporter-linux-amd64
```

## Next Steps

Once your binary is built and running:

1. **[Configure Prometheus](../usage/prometheus-setup)** to scrape metrics
2. **[Set up monitoring](../usage/grafana-dashboards)** with Grafana dashboards
3. **[Explore metrics](../reference/metrics)** and create custom queries
4. **[Contribute to development](../contributing)** by submitting improvements

For production deployments, consider:
- **[systemd installation](systemd)** for service management
- **[Docker](docker)** for containerization
- **[Helm](helm)** for Kubernetes deployment 