---
layout: default
title: systemd
parent: Installation
nav_order: 4
---

# systemd Installation
{: .no_toc }

Deploy NetBird API Exporter as a native Linux service using systemd.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Overview

The systemd installation provides native OS integration, perfect for:
- **Traditional Linux servers**
- **Production environments requiring OS-level service management**
- **Environments with strict security requirements**
- **Infrastructure requiring fine-grained resource control**
- **Servers without container orchestration**

## Prerequisites

Before starting, ensure you have:

- **Linux distribution with systemd** (CentOS/RHEL 7+, Ubuntu 16.04+, Debian 8+)
- **sudo/root access** for service installation
- **NetBird API token** ([get one here](../getting-started/authentication))
- **Go 1.21+** (if building from source)
- **curl or wget** for downloading pre-built binaries

## Quick Start

### Step 1: Download Binary

Choose your method for getting the binary:

#### Option A: Download Pre-built Binary

```bash
# Download latest release (replace VERSION with actual version)
VERSION="0.1.0"
ARCH="amd64"  # or arm64

curl -L "https://github.com/matanbaruch/netbird-api-exporter/releases/download/v${VERSION}/netbird-api-exporter-linux-${ARCH}" \
  -o /tmp/netbird-api-exporter

# Make executable
chmod +x /tmp/netbird-api-exporter
```

#### Option B: Build from Source

```bash
# Clone repository
git clone https://github.com/matanbaruch/netbird-api-exporter.git
cd netbird-api-exporter

# Build binary
go build -o /tmp/netbird-api-exporter .
```

### Step 2: Install Binary

```bash
# Copy to system location
sudo cp /tmp/netbird-api-exporter /usr/local/bin/
sudo chmod +x /usr/local/bin/netbird-api-exporter

# Verify installation
/usr/local/bin/netbird-api-exporter --version
```

### Step 3: Create User and Directories

```bash
# Create dedicated user (no login shell for security)
sudo useradd --system --no-create-home --shell /bin/false netbird-api-exporter

# Create configuration directory
sudo mkdir -p /etc/netbird-api-exporter
sudo chown netbird-api-exporter:netbird-api-exporter /etc/netbird-api-exporter
sudo chmod 750 /etc/netbird-api-exporter

# Create log directory (optional)
sudo mkdir -p /var/log/netbird-api-exporter
sudo chown netbird-api-exporter:netbird-api-exporter /var/log/netbird-api-exporter
sudo chmod 755 /var/log/netbird-api-exporter
```

### Step 4: Create Configuration

Create environment configuration file:

```bash
sudo tee /etc/netbird-api-exporter/config > /dev/null << 'EOF'
# NetBird API Configuration
NETBIRD_API_TOKEN=nb_api_your_token_here
NETBIRD_API_URL=https://api.netbird.io

# Exporter Configuration
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=info
EOF

# Secure the configuration file
sudo chmod 640 /etc/netbird-api-exporter/config
sudo chown root:netbird-api-exporter /etc/netbird-api-exporter/config
```

{: .important }
> **Security**: Replace `nb_api_your_token_here` with your actual NetBird API token from the [authentication guide](../getting-started/authentication).

### Step 5: Create systemd Service

Create the service file using the provided template:

```bash
# Download service file from repository
sudo curl -L "https://raw.githubusercontent.com/matanbaruch/netbird-api-exporter/main/netbird-exporter.service" \
  -o /etc/systemd/system/netbird-api-exporter.service

# Or create manually
sudo tee /etc/systemd/system/netbird-api-exporter.service > /dev/null << 'EOF'
[Unit]
Description=NetBird API Exporter
Documentation=https://github.com/matanbaruch/netbird-api-exporter
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=netbird-api-exporter
Group=netbird-api-exporter
ExecStart=/usr/local/bin/netbird-api-exporter
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=netbird-api-exporter

# Security settings
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
RestrictRealtime=yes
RestrictNamespaces=yes

# Load environment from file
EnvironmentFile=/etc/netbird-api-exporter/config

[Install]
WantedBy=multi-user.target
EOF
```

### Step 6: Start and Enable Service

```bash
# Reload systemd to read the new service
sudo systemctl daemon-reload

# Enable service to start at boot
sudo systemctl enable netbird-api-exporter

# Start the service
sudo systemctl start netbird-api-exporter

# Check status
sudo systemctl status netbird-api-exporter
```

### Step 7: Verify Installation

```bash
# Check service status
sudo systemctl is-active netbird-api-exporter

# Test health endpoint
curl http://localhost:8080/health

# Test metrics endpoint
curl http://localhost:8080/metrics
```

## Configuration

### Environment File

The main configuration file is `/etc/netbird-api-exporter/config`:

```bash
# NetBird API Configuration
NETBIRD_API_TOKEN=nb_api_your_actual_token_here
NETBIRD_API_URL=https://api.netbird.io

# Exporter Configuration
LISTEN_ADDRESS=:8080
METRICS_PATH=/metrics
LOG_LEVEL=info

# Optional: Custom timeout settings
# HTTP_TIMEOUT=30s
# API_RETRY_ATTEMPTS=3
```

### Service Configuration

The systemd service file provides comprehensive security and operational settings:

```ini
[Unit]
Description=NetBird API Exporter
Documentation=https://github.com/matanbaruch/netbird-api-exporter
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=netbird-api-exporter
Group=netbird-api-exporter
ExecStart=/usr/local/bin/netbird-api-exporter
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=netbird-api-exporter

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
PrivateTmp=yes
PrivateDevices=yes
RestrictRealtime=yes
RestrictNamespaces=yes
LockPersonality=yes
MemoryDenyWriteExecute=yes
RestrictAddressFamilies=AF_UNIX AF_INET AF_INET6

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

# Load environment from file
EnvironmentFile=/etc/netbird-api-exporter/config

[Install]
WantedBy=multi-user.target
```

## Management Commands

### Service Control

```bash
# Start service
sudo systemctl start netbird-api-exporter

# Stop service
sudo systemctl stop netbird-api-exporter

# Restart service
sudo systemctl restart netbird-api-exporter

# Reload configuration (restart required for env changes)
sudo systemctl reload-or-restart netbird-api-exporter

# Enable auto-start at boot
sudo systemctl enable netbird-api-exporter

# Disable auto-start at boot
sudo systemctl disable netbird-api-exporter
```

### Status and Monitoring

```bash
# Check service status
sudo systemctl status netbird-api-exporter

# Check if service is running
sudo systemctl is-active netbird-api-exporter

# Check if service is enabled
sudo systemctl is-enabled netbird-api-exporter

# View recent logs
sudo journalctl -u netbird-api-exporter -n 50

# Follow logs in real-time
sudo journalctl -u netbird-api-exporter -f

# View logs for specific time range
sudo journalctl -u netbird-api-exporter --since "1 hour ago"
```

### Configuration Updates

```bash
# Edit configuration
sudo nano /etc/netbird-api-exporter/config

# Validate configuration syntax (check logs after restart)
sudo systemctl restart netbird-api-exporter
sudo systemctl status netbird-api-exporter

# Reload systemd if service file changed
sudo systemctl daemon-reload
sudo systemctl restart netbird-api-exporter
```

## Security Hardening

### File Permissions

```bash
# Binary permissions
sudo chown root:root /usr/local/bin/netbird-api-exporter
sudo chmod 755 /usr/local/bin/netbird-api-exporter

# Configuration permissions
sudo chown root:netbird-api-exporter /etc/netbird-api-exporter/config
sudo chmod 640 /etc/netbird-api-exporter/config

# Service file permissions
sudo chown root:root /etc/systemd/system/netbird-api-exporter.service
sudo chmod 644 /etc/systemd/system/netbird-api-exporter.service
```

### Enhanced Service Security

For maximum security, create an enhanced service file:

```bash
sudo tee /etc/systemd/system/netbird-api-exporter.service > /dev/null << 'EOF'
[Unit]
Description=NetBird API Exporter
Documentation=https://github.com/matanbaruch/netbird-api-exporter
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=netbird-api-exporter
Group=netbird-api-exporter
ExecStart=/usr/local/bin/netbird-api-exporter
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=netbird-api-exporter

# Enhanced security settings
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectKernelLogs=yes
ProtectControlGroups=yes
PrivateTmp=yes
PrivateDevices=yes
PrivateUsers=yes
ProtectHostname=yes
ProtectClock=yes
RestrictRealtime=yes
RestrictNamespaces=yes
RestrictSUIDSGID=yes
LockPersonality=yes
MemoryDenyWriteExecute=yes
RemoveIPC=yes

# Network restrictions
RestrictAddressFamilies=AF_UNIX AF_INET AF_INET6
IPAddressDeny=any
IPAddressAllow=localhost
IPAddressAllow=10.0.0.0/8
IPAddressAllow=172.16.0.0/12
IPAddressAllow=192.168.0.0/16

# Capability restrictions
CapabilityBoundingSet=
AmbientCapabilities=

# System call filtering
SystemCallFilter=@system-service
SystemCallFilter=~@privileged @resources @mount

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096
LimitCORE=0
LimitAS=1G
TasksMax=4096

# Environment
EnvironmentFile=/etc/netbird-api-exporter/config
Environment=HOME=/var/lib/netbird-api-exporter

[Install]
WantedBy=multi-user.target
EOF
```

### Firewall Configuration

Configure firewall to allow only necessary access:

```bash
# For UFW (Ubuntu/Debian)
sudo ufw allow from 10.0.0.0/8 to any port 8080
sudo ufw allow from 172.16.0.0/12 to any port 8080
sudo ufw allow from 192.168.0.0/16 to any port 8080

# For firewalld (CentOS/RHEL/Fedora)
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="10.0.0.0/8" port protocol="tcp" port="8080" accept'
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="172.16.0.0/12" port protocol="tcp" port="8080" accept'
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.0.0/16" port protocol="tcp" port="8080" accept'
sudo firewall-cmd --reload

# For iptables
sudo iptables -A INPUT -s 10.0.0.0/8 -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -s 172.16.0.0/12 -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -s 192.168.0.0/16 -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8080 -j DROP
```

## Monitoring and Logging

### Log Management

Configure log rotation to prevent disk space issues:

```bash
# Create logrotate configuration
sudo tee /etc/logrotate.d/netbird-api-exporter > /dev/null << 'EOF'
/var/log/netbird-api-exporter/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    copytruncate
    notifempty
    create 644 netbird-api-exporter netbird-api-exporter
}
EOF
```

### SystemD Journal Configuration

Configure journal limits for the service:

```bash
# Create journal configuration directory
sudo mkdir -p /etc/systemd/system/netbird-api-exporter.service.d

# Create journal limits
sudo tee /etc/systemd/system/netbird-api-exporter.service.d/journal.conf > /dev/null << 'EOF'
[Service]
# Limit journal size for this service
StandardOutput=journal
StandardError=journal
SyslogIdentifier=netbird-api-exporter
LogRateLimitIntervalSec=30
LogRateLimitBurst=1000
EOF

# Reload systemd
sudo systemctl daemon-reload
sudo systemctl restart netbird-api-exporter
```

### Health Monitoring

Create a simple health check script:

```bash
sudo tee /usr/local/bin/netbird-exporter-healthcheck > /dev/null << 'EOF'
#!/bin/bash

ENDPOINT="http://localhost:8080/health"
TIMEOUT=10

if curl -s --max-time $TIMEOUT "$ENDPOINT" > /dev/null; then
    echo "$(date): NetBird API Exporter is healthy"
    exit 0
else
    echo "$(date): NetBird API Exporter health check failed"
    exit 1
fi
EOF

sudo chmod +x /usr/local/bin/netbird-exporter-healthcheck

# Test the health check
/usr/local/bin/netbird-exporter-healthcheck
```

## Troubleshooting

### Common Issues

#### 1. Service Won't Start

```bash
# Check service status and logs
sudo systemctl status netbird-api-exporter
sudo journalctl -u netbird-api-exporter -n 50

# Common causes:
# - Binary not found or not executable
# - Configuration file syntax error
# - Port already in use
# - Permission issues
```

#### 2. Permission Denied Errors

```bash
# Check file permissions
ls -la /usr/local/bin/netbird-api-exporter
ls -la /etc/netbird-api-exporter/config

# Fix permissions
sudo chown root:root /usr/local/bin/netbird-api-exporter
sudo chmod 755 /usr/local/bin/netbird-api-exporter
sudo chown root:netbird-api-exporter /etc/netbird-api-exporter/config
sudo chmod 640 /etc/netbird-api-exporter/config
```

#### 3. API Authentication Issues

```bash
# Test API token manually
curl -H "Authorization: Bearer $(grep NETBIRD_API_TOKEN /etc/netbird-api-exporter/config | cut -d= -f2)" \
     https://api.netbird.io/api/peers

# Check configuration file
sudo cat /etc/netbird-api-exporter/config
```

### Debug Mode

Enable debug logging temporarily:

```bash
# Edit configuration
sudo sed -i 's/LOG_LEVEL=info/LOG_LEVEL=debug/' /etc/netbird-api-exporter/config

# Restart service
sudo systemctl restart netbird-api-exporter

# View debug logs
sudo journalctl -u netbird-api-exporter -f

# Revert to info level when done
sudo sed -i 's/LOG_LEVEL=debug/LOG_LEVEL=info/' /etc/netbird-api-exporter/config
sudo systemctl restart netbird-api-exporter
```

## Updates

### Update Binary

```bash
# Download new version
VERSION="0.2.0"
curl -L "https://github.com/matanbaruch/netbird-api-exporter/releases/download/v${VERSION}/netbird-api-exporter-linux-amd64" \
  -o /tmp/netbird-api-exporter-new

# Stop service
sudo systemctl stop netbird-api-exporter

# Backup current binary
sudo cp /usr/local/bin/netbird-api-exporter /usr/local/bin/netbird-api-exporter.backup

# Install new binary
sudo cp /tmp/netbird-api-exporter-new /usr/local/bin/netbird-api-exporter
sudo chmod +x /usr/local/bin/netbird-api-exporter

# Start service
sudo systemctl start netbird-api-exporter

# Verify update
sudo systemctl status netbird-api-exporter
```

### Update Configuration

```bash
# Edit configuration
sudo nano /etc/netbird-api-exporter/config

# Restart service to apply changes
sudo systemctl restart netbird-api-exporter

# Verify configuration
sudo systemctl status netbird-api-exporter
```

## Uninstallation

To completely remove the NetBird API Exporter:

```bash
# Stop and disable service
sudo systemctl stop netbird-api-exporter
sudo systemctl disable netbird-api-exporter

# Remove service file
sudo rm /etc/systemd/system/netbird-api-exporter.service
sudo systemctl daemon-reload

# Remove binary
sudo rm /usr/local/bin/netbird-api-exporter

# Remove configuration
sudo rm -rf /etc/netbird-api-exporter

# Remove user
sudo userdel netbird-api-exporter

# Remove logs (optional)
sudo rm -rf /var/log/netbird-api-exporter

# Remove logrotate configuration
sudo rm /etc/logrotate.d/netbird-api-exporter
```

## Next Steps

Once your systemd service is running:

1. **[Configure Prometheus](../usage/prometheus-setup)** to scrape metrics from the exporter
2. **[Set up monitoring dashboards](../usage/grafana-dashboards)** in Grafana
3. **[Configure alerting](../usage/alerting)** for important metrics
4. **[Review security settings](../reference/security)** and hardening options

For other deployment methods:
- **[Docker Compose](docker-compose)** for containerized development
- **[Helm](helm)** for Kubernetes clusters 