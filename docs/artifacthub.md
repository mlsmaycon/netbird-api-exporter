# Publishing to Artifact Hub

This document explains how the NetBird API Exporter is published to [Artifact Hub](https://artifacthub.io/) and how you can add it manually if needed.

## Automated Publishing

The NetBird API Exporter is automatically published to Artifact Hub through our CI/CD pipeline. Here's how it works:

### 1. OCI Package Publishing

When a new release is created, our GitHub Actions workflow:

1. **Packages the Helm chart** into an OCI artifact
2. **Pushes it to GitHub Container Registry** at `ghcr.io/matanbaruch/netbird-api-exporter/charts/netbird-api-exporter`
3. **Publishes repository metadata** using ORAS to enable Artifact Hub integration

### 2. Repository Metadata

The [`artifacthub-repo.yml`](../artifacthub-repo.yml) file in the repository root contains metadata that enables:

- Repository ownership verification
- Enhanced package information display
- Integration with Artifact Hub's discovery features

### 3. Chart Annotations

The [`Chart.yaml`](../charts/netbird-api-exporter/Chart.yaml) includes Artifact Hub specific annotations:

- `artifacthub.io/license`: License information (MIT)
- `artifacthub.io/links`: Useful links (repository, docs, changelog, issues)
- `artifacthub.io/maintainers`: Maintainer contact information
- `artifacthub.io/recommendations`: Related packages (Prometheus, Grafana)
- `artifacthub.io/changes`: Release notes and changes

## Manual Repository Addition

If you need to add the repository to Artifact Hub manually, follow these steps:

### Prerequisites

1. Sign in to [Artifact Hub](https://artifacthub.io/)
2. Ensure you have access to add repositories

### Steps

1. **Navigate to Control Panel**

   - Click on your profile in the top right
   - Select "Control Panel"

2. **Add Repository**

   - Go to the "Repositories" tab
   - Click "Add Repository"

3. **Configure Repository**

   - **Kind**: Helm charts
   - **Name**: `netbird-api-exporter`
   - **Display name**: `NetBird API Exporter`
   - **URL**: `oci://ghcr.io/matanbaruch/netbird-api-exporter/charts/netbird-api-exporter`
   - **Description**: A Prometheus exporter that collects metrics from the NetBird API

4. **Save and Wait**
   - Click "Add Repository"
   - Wait for Artifact Hub to index the repository (may take a few minutes)

## Repository URL Format

For OCI-based Helm repositories, use this URL format:

```
oci://ghcr.io/matanbaruch/netbird-api-exporter/charts/netbird-api-exporter
```

## Verification

After the repository is added, you can verify it's working by:

1. **Browse the Package**

   - Visit the [package page](https://artifacthub.io/packages/helm/netbird-api-exporter/netbird-api-exporter)
   - Check that all metadata is displayed correctly

2. **Test Installation**
   ```bash
   helm upgrade --install netbird-api-exporter \
     oci://ghcr.io/matanbaruch/netbird-api-exporter/charts/netbird-api-exporter \
     --set netbird.apiToken=your_token_here
   ```

## Verified Publisher Status

To enable the "Verified Publisher" badge:

1. **Get Repository ID**

   - After adding the repository, note the repository ID from the control panel

2. **Update Metadata**

   - Add the `repositoryID` field to [`artifacthub-repo.yml`](../artifacthub-repo.yml)
   - The next release will include this metadata

3. **Verification**
   - Artifact Hub will verify ownership during the next indexing cycle
   - The verified publisher badge will appear automatically

## Troubleshooting

### Repository Not Indexing

If the repository doesn't appear or isn't indexing:

1. **Check URL Format**: Ensure the OCI URL is correct
2. **Verify Access**: Make sure the repository is publicly accessible
3. **Check Logs**: Look for error messages in the Artifact Hub control panel
4. **Wait**: Initial indexing can take up to 30 minutes

### Missing Metadata

If package information is incomplete:

1. **Check Annotations**: Verify Chart.yaml has proper annotations
2. **Repository Metadata**: Ensure artifacthub-repo.yml is properly formatted
3. **Re-index**: Force a re-index by updating the repository

### Installation Issues

If users report installation problems:

1. **Test Locally**: Verify the chart installs correctly
2. **Check Versions**: Ensure version tags are semantic versions
3. **Documentation**: Update installation instructions if needed

## Related Links

- [Artifact Hub Documentation](https://artifacthub.io/docs/)
- [Helm OCI Support](https://helm.sh/blog/storing-charts-in-oci/)
- [ORAS CLI](https://oras.land/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
