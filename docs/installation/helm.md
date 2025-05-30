---
layout: default
title: Helm
parent: Installation
nav_order: 3
---

# Helm Installation
{: .no_toc }

Deploy NetBird API Exporter to Kubernetes using Helm charts.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Overview

The Helm chart provides a cloud-native deployment method, perfect for:
- **Kubernetes clusters**
- **Production environments**
- **GitOps workflows**
- **Multi-environment deployments**
- **Scalable monitoring setups**

## Prerequisites

Before starting, ensure you have:

- **Kubernetes 1.19+** cluster with kubectl access
- **Helm 3.0+** installed
- **NetBird API token** ([get one here](../getting-started/authentication))
- **Sufficient RBAC permissions** to create deployments and services

## Quick Start

### Step 1: Add Helm Repository

```bash
# Add the chart repository
helm repo add netbird-api-exporter oci://ghcr.io/matanbaruch/netbird-api-exporter/charts
helm repo update
```

### Step 2: Create Secrets

Create a Kubernetes secret with your NetBird API token:

```bash
kubectl create secret generic netbird-api-token \
  --from-literal=token="nb_api_your_token_here"
```

### Step 3: Install the Chart

```bash
helm upgrade --install netbird-api-exporter \
  netbird-api-exporter/netbird-api-exporter \
  --set netbird.existingSecret="netbird-api-token"
```

### Step 4: Verify Installation

```bash
# Check pod status
kubectl get pods -l app.kubernetes.io/name=netbird-api-exporter

# Check service
kubectl get svc netbird-api-exporter

# Port forward to test
kubectl port-forward svc/netbird-api-exporter 8080:8080

# Test health endpoint
curl http://localhost:8080/health
```

## Installation Methods

### Method 1: Using Values File (Recommended)

Create a `values.yaml` file:

```yaml
# NetBird Configuration
netbird:
  apiToken: "nb_api_your_token_here"
  apiUrl: "https://api.netbird.io"

# Service Configuration
service:
  type: ClusterIP
  port: 8080
  annotations: {}

# Deployment Configuration
replicaCount: 1

image:
  repository: ghcr.io/matanbaruch/netbird-api-exporter
  pullPolicy: IfNotPresent
  tag: "latest"

# Resource Management
resources:
  limits:
    cpu: 200m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi

# Monitoring
serviceMonitor:
  enabled: true
  interval: 30s
  path: /metrics
  labels:
    release: prometheus
```

Install with values file:

```bash
helm upgrade --install netbird-api-exporter \
  netbird-api-exporter/netbird-api-exporter \
  -f values.yaml
```

### Method 2: Using Existing Secret

For production environments, use Kubernetes secrets:

```bash
# Create secret
kubectl create secret generic netbird-api-secret \
  --from-literal=token="nb_api_your_token_here"

# Install chart referencing secret
helm upgrade --install netbird-api-exporter \
  netbird-api-exporter/netbird-api-exporter \
  --set netbird.existingSecret="netbird-api-secret"
```

### Method 3: Command Line Values

For quick testing:

```bash
helm upgrade --install netbird-api-exporter \
  netbird-api-exporter/netbird-api-exporter \
  --set netbird.apiToken="nb_api_your_token_here" \
  --set netbird.apiUrl="https://api.netbird.io" \
  --set service.type="LoadBalancer"
```

## Configuration

### Complete Values File

Here's a comprehensive `values.yaml` example:

```yaml
# Default values for netbird-api-exporter

# NetBird API Configuration
netbird:
  # API token for authentication (required)
  apiToken: ""
  
  # Alternative: use existing secret containing the token
  existingSecret: ""
  secretKey: "token"
  
  # NetBird API URL
  apiUrl: "https://api.netbird.io"

# Container image configuration
image:
  repository: ghcr.io/matanbaruch/netbird-api-exporter
  pullPolicy: IfNotPresent
  tag: "latest"

# Deployment configuration
replicaCount: 1
nameOverride: ""
fullnameOverride: ""

# Security context
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 65534
  fsGroup: 65534

securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 65534

# Service configuration
service:
  type: ClusterIP
  port: 8080
  targetPort: 8080
  annotations: {}

# Ingress configuration
ingress:
  enabled: false
  className: ""
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: netbird-exporter.local
      paths:
        - path: /
          pathType: Prefix
  tls: []

# Resource limits and requests
resources:
  limits:
    cpu: 200m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi

# Horizontal Pod Autoscaler
autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 5
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80

# Node selection
nodeSelector: {}
tolerations: []
affinity: {}

# ServiceMonitor for Prometheus Operator
serviceMonitor:
  enabled: false
  namespace: ""
  interval: 30s
  path: /metrics
  labels: {}
  annotations: {}

# Prometheus rules
prometheusRule:
  enabled: false
  namespace: ""
  labels: {}
  rules: []

# Pod disruption budget
podDisruptionBudget:
  enabled: false
  minAvailable: 1

# Liveness and readiness probes
livenessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 30
  periodSeconds: 30
  timeoutSeconds: 10
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health
    port: http
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

# Additional environment variables
extraEnv: []
  # - name: LOG_LEVEL
  #   value: "debug"

# Annotations for pods
podAnnotations: {}

# Additional labels for all resources
commonLabels: {}
```

### Environment Variables

All configuration is done through values.yaml, but you can add extra environment variables:

```yaml
extraEnv:
  - name: LOG_LEVEL
    value: "debug"
  - name: LISTEN_ADDRESS
    value: ":8080"
  - name: METRICS_PATH
    value: "/metrics"
```

## Production Configuration

### High Availability Setup

```yaml
# values-production.yaml
replicaCount: 3

resources:
  limits:
    cpu: 500m
    memory: 256Mi
  requests:
    cpu: 200m
    memory: 128Mi

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app.kubernetes.io/name
            operator: In
            values:
            - netbird-api-exporter
        topologyKey: kubernetes.io/hostname

podDisruptionBudget:
  enabled: true
  minAvailable: 2

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80
```

### Security Hardening

```yaml
# values-security.yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 65534
  fsGroup: 65534
  seccompProfile:
    type: RuntimeDefault

securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 65534

# Network policies (if supported)
networkPolicy:
  enabled: true
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            name: monitoring
      ports:
      - protocol: TCP
        port: 8080
  egress:
    - to: []
      ports:
      - protocol: TCP
        port: 443  # HTTPS to NetBird API
```

## Monitoring Integration

### Prometheus Operator

If you're using the Prometheus Operator:

```yaml
serviceMonitor:
  enabled: true
  namespace: monitoring
  interval: 30s
  path: /metrics
  labels:
    release: prometheus-operator
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
```

### Traditional Prometheus

For traditional Prometheus setups, add scrape configuration:

```yaml
# Add to your Prometheus config
scrape_configs:
  - job_name: 'netbird-api-exporter'
    kubernetes_sd_configs:
      - role: endpoints
        namespaces:
          names:
          - default  # Update with your namespace
    relabel_configs:
      - source_labels: [__meta_kubernetes_service_name]
        action: keep
        regex: netbird-api-exporter
```

### Grafana Integration

Deploy Grafana with the exporter:

```yaml
# Add to values.yaml
grafana:
  enabled: true
  adminPassword: "admin"
  service:
    type: LoadBalancer
  datasources:
    datasources.yaml:
      apiVersion: 1
      datasources:
      - name: Prometheus
        type: prometheus
        url: http://prometheus:9090
        isDefault: true
```

## Ingress Setup

### Nginx Ingress

```yaml
ingress:
  enabled: true
  className: "nginx"
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: netbird-exporter.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: netbird-exporter-tls
      hosts:
        - netbird-exporter.example.com
```

### Traefik Ingress

```yaml
ingress:
  enabled: true
  className: "traefik"
  annotations:
    traefik.ingress.kubernetes.io/router.tls: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
  hosts:
    - host: netbird-exporter.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: netbird-exporter-tls
      hosts:
        - netbird-exporter.example.com
```

## Management Commands

### Install/Upgrade

```bash
# Install
helm install netbird-api-exporter netbird-api-exporter/netbird-api-exporter -f values.yaml

# Upgrade
helm upgrade netbird-api-exporter netbird-api-exporter/netbird-api-exporter -f values.yaml

# Install or upgrade
helm upgrade --install netbird-api-exporter netbird-api-exporter/netbird-api-exporter -f values.yaml
```

### Status and Information

```bash
# Check status
helm status netbird-api-exporter

# Get values
helm get values netbird-api-exporter

# Get manifest
helm get manifest netbird-api-exporter

# List releases
helm list
```

### Rollback

```bash
# See revision history
helm history netbird-api-exporter

# Rollback to previous version
helm rollback netbird-api-exporter

# Rollback to specific revision
helm rollback netbird-api-exporter 2
```

### Uninstall

```bash
# Uninstall release
helm uninstall netbird-api-exporter

# Uninstall and delete secrets
kubectl delete secret netbird-api-token
```

## Troubleshooting

### Common Issues

#### 1. Pod Not Starting

```bash
# Check pod status
kubectl get pods -l app.kubernetes.io/name=netbird-api-exporter

# Describe pod
kubectl describe pod <pod-name>

# Check logs
kubectl logs <pod-name>
```

#### 2. Secret Not Found

```bash
# Verify secret exists
kubectl get secrets

# Check secret content
kubectl describe secret netbird-api-token
```

#### 3. Service Not Accessible

```bash
# Check service
kubectl get svc netbird-api-exporter

# Check endpoints
kubectl get endpoints netbird-api-exporter

# Port forward for testing
kubectl port-forward svc/netbird-api-exporter 8080:8080
```

### Debug Mode

Enable debug logging:

```yaml
extraEnv:
  - name: LOG_LEVEL
    value: "debug"
```

Upgrade the deployment:
```bash
helm upgrade netbird-api-exporter netbird-api-exporter/netbird-api-exporter -f values.yaml
```

### Resource Issues

Check resource usage:

```bash
# Pod resource usage
kubectl top pods -l app.kubernetes.io/name=netbird-api-exporter

# Node resources
kubectl top nodes

# Describe node for resource pressure
kubectl describe node <node-name>
```

## Multi-Environment Deployment

### Development Environment

```yaml
# values-dev.yaml
replicaCount: 1

resources:
  requests:
    cpu: 50m
    memory: 32Mi
  limits:
    cpu: 100m
    memory: 64Mi

serviceMonitor:
  enabled: false
```

### Staging Environment

```yaml
# values-staging.yaml
replicaCount: 2

resources:
  requests:
    cpu: 100m
    memory: 64Mi
  limits:
    cpu: 200m
    memory: 128Mi

serviceMonitor:
  enabled: true
  labels:
    release: prometheus-staging
```

### Production Environment

```yaml
# values-prod.yaml
replicaCount: 3

resources:
  requests:
    cpu: 200m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10

podDisruptionBudget:
  enabled: true
  minAvailable: 2

serviceMonitor:
  enabled: true
  labels:
    release: prometheus-production
```

Deploy to each environment:

```bash
# Development
helm upgrade --install netbird-api-exporter-dev \
  netbird-api-exporter/netbird-api-exporter \
  -f values-dev.yaml \
  --namespace development

# Staging
helm upgrade --install netbird-api-exporter-staging \
  netbird-api-exporter/netbird-api-exporter \
  -f values-staging.yaml \
  --namespace staging

# Production
helm upgrade --install netbird-api-exporter-prod \
  netbird-api-exporter/netbird-api-exporter \
  -f values-prod.yaml \
  --namespace production
```

## Next Steps

Once your Helm deployment is running:

1. **[Configure monitoring](../usage/prometheus-setup)** with Prometheus and Grafana
2. **[Set up alerting](../usage/alerting)** for important metrics
3. **[Explore metrics](../reference/metrics)** and create custom dashboards
4. **[Implement GitOps](../technical/gitops)** for automated deployments

For other deployment methods:
- **[Docker Compose](docker-compose)** for local development
- **[systemd](systemd)** for traditional servers 