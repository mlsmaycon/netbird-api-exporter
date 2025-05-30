# NetBird API Exporter Helm Chart

This Helm chart deploys the NetBird API Exporter, a Prometheus exporter that provides comprehensive metrics about your NetBird network including peers, groups, users, networks, and DNS configuration.

## Features

- Comprehensive NetBird metrics collection
- Secure API token handling via Kubernetes secrets
- Optional Prometheus Operator integration with ServiceMonitor
- Configurable resource limits and requests
- Health checks and readiness probes
- Optional ingress for external access
- Horizontal Pod Autoscaler support
- Pod Disruption Budget for high availability
- Security hardening with non-root containers and read-only filesystem

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- A valid NetBird API token

## Installing the Chart

### Install from local directory

```bash
# Clone the repository
git clone https://github.com/matanbaruch/netbird-api-exporter.git
cd netbird-api-exporter

# Install the chart
helm install my-netbird-api-exporter ./charts/netbird-api-exporter \
  --set netbird.apiToken="your-netbird-api-token"
```

### Install with custom values

```bash
helm install my-netbird-api-exporter ./charts/netbird-api-exporter \
  --values my-values.yaml
```

## Configuration

The following table lists the configurable parameters and their default values.

### Basic Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Image repository | `netbird-api-exporter` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `image.tag` | Image tag | `"latest"` |

### NetBird Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `netbird.apiUrl` | NetBird API URL | `"https://api.netbird.io"` |
| `netbird.apiToken` | NetBird API token (stored in secret) | `""` |

### Application Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `config.listenAddress` | Listen address and port | `":8080"` |
| `config.metricsPath` | Metrics endpoint path | `"/metrics"` |
| `config.logLevel` | Log level | `"info"` |

## Usage Examples

### Basic Installation

```bash
helm install netbird-api-exporter ./charts/netbird-api-exporter \
  --set netbird.apiToken="nb_token_xxx"
```

### With Prometheus Operator

```bash
helm install netbird-api-exporter ./charts/netbird-api-exporter \
  --set netbird.apiToken="nb_token_xxx" \
  --set serviceMonitor.enabled=true \
  --set serviceMonitor.additionalLabels.release=prometheus
```

### Production Configuration

```bash
helm install netbird-api-exporter ./charts/netbird-api-exporter \
  --set netbird.apiToken="nb_token_xxx" \
  --values ./charts/netbird-api-exporter/values-production.yaml
```

## Getting Your NetBird API Token

1. Log into your NetBird dashboard
2. Go to **Settings** → **API Keys**  
3. Create a new API key with appropriate permissions
4. Copy the token and use it as `netbird.apiToken`

## Monitoring

Once deployed, the exporter provides metrics at the `/metrics` endpoint. Key metrics include:

- `netbird_peers_total` - Total number of peers
- `netbird_peers_connected` - Connected/disconnected peers
- `netbird_groups_total` - Total number of groups
- `netbird_users_total` - Total number of users
- `netbird_networks_total` - Total number of networks
- `netbird_dns_nameserver_groups_total` - Total DNS nameserver groups

## Testing

Run the included test:

```bash
helm test netbird-api-exporter
```

## Troubleshooting

### Check pod logs

```bash
kubectl logs -l app.kubernetes.io/name=netbird-api-exporter
```

### Test connectivity

```bash
kubectl port-forward svc/netbird-api-exporter 8080:8080
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```
