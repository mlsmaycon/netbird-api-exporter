# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.35] - 2025-06-02

## [0.1.34] - 2025-06-02

## [0.1.33] - 2025-06-02

## [0.1.32] - 2025-06-01

## [0.1.31] - 2025-06-01

### Bugfix
- Fix update-changelog.sh to properly update Helm chart artifacthub.io/changes with correct format and handle empty descriptions
Files modified in this change:
- Added: .github/dependabot.yml
- Modified: scripts/update-changelog.sh

## [0.1.30] - 2025-06-01

## [0.1.29] - 2025-06-01


### Bugfix
- Fix ArtifactHub.io annotations by properly quoting strings and correcting alternative name
Files modified in this change:
- Modified: charts/netbird-api-exporter/Chart.yaml

## [0.1.28] - 2025-06-01
### Features

- Refactor GitHub issue templates with improved structure and NetBird-specific context
  Files modified in this change:
- Deleted: .github/ISSUE_TEMPLATE.md
- Modified: .github/ISSUE_TEMPLATE/bug_report.md
- Deleted: .github/ISSUE_TEMPLATE/custom.md
- Modified: .github/ISSUE_TEMPLATE/feature_request.md
- New: .github/ISSUE_TEMPLATE/config.yml
- New: .github/ISSUE_TEMPLATE/documentation.md

## [0.1.27] - 2025-06-01

### Features

- Enhanced update-changelog.sh script to also update Helm chart artifacthub.io/changes annotation and include uncommitted files summary
  Files modified in this change:
- Modified: .github/workflows/release.yml
- Modified: ARCHITECTURE.md
- Modified: CHANGELOG.md
- Modified: README.md
- Modified: SECURITY.md
- Modified: charts/netbird-api-exporter/Chart.yaml
- Modified: scripts/update-changelog.sh

### Bugfix

- Fix Helm chart name to resolve packaging filename mismatch
  Files modified in this change:
- Modified: charts/netbird-api-exporter/Chart.yaml

### Security

- Add build provenance attestations for Docker images and Go binaries

## [0.1.26] - 2025-06-01

### Features

- Add comprehensive GitHub badges including build status, Go version, Prometheus, Docker, Kubernetes, license, and distribution metrics
  Files modified in this change:
- Modified: README.md

## [0.1.25] - 2025-06-01

### Features

- Add icon to helm chart for better visual identification in Artifact Hub
  Files modified in this change:
- Modified: charts/netbird-api-exporter/Chart.yaml

### Bugfix

- Fix Helm chart name to resolve packaging filename mismatch
  Files modified in this change:
- Modified: charts/netbird-api-exporter/Chart.yaml

## [0.1.24] - 2025-06-01

### Security

- Add build provenance attestations for Docker images and Go binaries

## [0.1.23] - 2025-06-01

### Bugfix

- Fix GitHub Actions workflow by installing ORAS before using it in release pipeline
- Fix GitHub Actions release workflow Helm packaging error with 'latest' version

## [0.1.22] - 2025-06-01

### Features

- Only commit version bump after successful Docker and Helm builds
- Only push git tags after successful Docker and Helm builds

## [0.1.21] - 2025-06-01

### Features

- Enhanced release automation and build process improvements

## [0.1.20] - 2025-06-01

### Features

- Build process and CI/CD improvements

## [0.1.19] - 2025-06-01

### Features

- Release pipeline enhancements

## [0.1.18] - 2025-06-01

### Features

- Add publishing to artifacthub.io for the OCI package of the helm chart

## [0.1.17] - 2025-05-31

### Features

- Add comprehensive Grafana dashboard with documentation
- Add comprehensive SECURITY.md with vulnerability reporting procedures and security guidelines

## [0.1.16] - 2025-05-31

### Bugfix

- Fix bundler version compatibility for Dependabot (update from v1.17.2 to v2.5.0)

## [0.1.15] - 2025-05-31

### Features

- Enhanced development environment and tooling improvements

## [0.1.14] - 2025-05-31

### Features

- Add comprehensive test coverage for NetBird client and types
- Migrate from legacy .cursorrules to modern .cursor/rules MDC format

### Bugfix

- Fix changelog update script breaking the changelog structure and creating duplicate sections

## [0.1.13] - 2025-05-31

### Features

- Development workflow improvements

## [0.1.12] - 2025-05-31

### Features

- Improve development workflow with automatic environment setup

## [0.1.11] - 2025-05-31

### Features

- Enhanced release process with automated changelog generation
- Add changelog file for better release tracking

## [0.1.10] - 2025-05-31

### Features

- Enhanced release process with automated changelog generation
- Add changelog file for better release tracking

<!--
When adding entries to the changelog, use the following guidelines:

### BREAKING CHANGES
- List any breaking changes that require user action
- Include migration instructions if needed

### Features
- New functionality and enhancements
- Use present tense (e.g., "Add support for...")

### Bugfix
- Bug fixes and corrections
- Reference issue numbers if applicable (e.g., "Fix memory leak in exporter (#123)")

### Security
- Security-related changes
- Vulnerability fixes

### Deprecated
- Features that will be removed in future versions

### Removed
- Features that have been removed
-->

## [0.1.9] - 2024-12-16

### Bugfix

- Refactor logging in main.go and networks.go for consistency
- Adjusted formatting in debug logging middleware
- Updated log fields in NetworksExporter to enhance readability and maintain uniformity in log output

## [0.1.8] - 2024-12-16

### Features

- Implement debug logging middleware for HTTP requests and responses
- Enhanced logging in exporters (DNS, Groups, Networks, Peers, Users) to capture API fetch details and metrics updates

## [0.1.7] - 2024-12-16

### Bugfix

- Update installation instructions in README and helm.md to reflect OCI registry usage
- Removed the need for adding a Helm repository
- Added examples for installing the NetBird API Exporter directly from the OCI registry
- Added configurations for ExternalSecret and production setups

## [0.1.6] - 2024-12-16

### Features

- Enhanced Docker security and permissions
- Added wget to Dockerfile dependencies
- Created app directory with proper permissions for the nobody user
- Adjusted binary copy command to set ownership and executable permissions

### Bugfix

- Update README.md documentation

## [0.1.5] - 2024-12-16

### Bugfix

- Add Jekyll build artifacts to gitignore
- Fix GitHub Pages deployment with custom blue color scheme for Just the Docs theme
- Fix Liquid template syntax errors in docker.md by escaping Docker template syntax
- Remove conflicting just-the-docs gem from Gemfile to avoid conflicts with github-pages
- Add jekyll-remote-theme plugin for remote theme support
- Comment out missing logo reference
- Add Gemfile.lock for reproducible builds

## [0.1.4] - 2024-12-16

### Features

- Enhanced GitHub Actions workflows by adding paths-ignore for documentation files
- Transition to Just the Docs theme for documentation
- Improved navigation and refined installation guides
- Introduced new sections for better organization and accessibility
- Removed deprecated navigation data file

## [0.1.3] - 2024-12-16

### Features

- Transition to ReadTheDocs theme for documentation
- Enhanced navigation and improved installation guides
- Added detailed sections for Docker, Docker Compose, and security considerations
- Streamlined content organization and updated links for better accessibility

## [0.1.2] - 2024-12-16

### Features

- Refactored documentation structure and enhanced installation guides
- Added comprehensive installation methods including Docker, Docker Compose, Helm, and systemd
- Introduced new configuration files and workflows for GitHub Pages deployment
- Improved overall clarity and accessibility of the documentation

## [0.1.1] - 2024-12-16

### Bugfix

- Updated README.md to remove contributing section
- Enhanced GitHub Actions workflow by adding permissions for write access
- Adjusted checkout settings for better version control

<<<<<<< HEAD
[0.1.10]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.9...v0.1.10
[0.1.11]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.10...v0.1.11
[0.1.12]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.11...v0.1.12
[0.1.13]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.12...v0.1.13
[0.1.14]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.13...v0.1.14
[0.1.15]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.14...v0.1.15
[0.1.16]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.15...v0.1.16
[0.1.17]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.16...v0.1.17
[0.1.18]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.17...v0.1.18
[0.1.19]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.18...v0.1.19
[0.1.20]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.19...v0.1.20
[0.1.21]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.20...v0.1.21
[0.1.22]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.21...v0.1.22
[0.1.23]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.22...v0.1.23
[0.1.24]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.23...v0.1.24
[0.1.25]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.24...v0.1.25
[0.1.26]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.25...v0.1.26
[0.1.27]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.26...v0.1.27
[0.1.28]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.27...v0.1.28
[0.1.29]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.28...v0.1.29
[0.1.30]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.29...v0.1.30
[0.1.31]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.30...v0.1.31
[0.1.32]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.31...v0.1.32
[0.1.33]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.32...v0.1.33
[0.1.34]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.33...v0.1.34
[0.1.35]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.34...v0.1.35
[Unreleased]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.35...HEAD
=======
[0.1.29]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.28...v0.1.29
[0.1.30]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.29...v0.1.30
[0.1.31]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.30...v0.1.31
[0.1.32]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.31...v0.1.32
[0.1.33]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.32...v0.1.33
[0.1.34]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.33...v0.1.34
[0.1.35]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.34...v0.1.35
[Unreleased]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.35...HEAD
[0.1.27]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.26...v0.1.27
[0.1.26]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.25...v0.1.26
[0.1.25]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.24...v0.1.25
[0.1.24]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.23...v0.1.24
[0.1.23]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.22...v0.1.23
[0.1.22]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.21...v0.1.22
[0.1.21]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.20...v0.1.21
[0.1.20]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.19...v0.1.20
[0.1.19]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.18...v0.1.19
[0.1.18]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.17...v0.1.18
[0.1.17]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.16...v0.1.17
[0.1.16]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.15...v0.1.16
[0.1.15]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.14...v0.1.15
[0.1.14]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.13...v0.1.14
[0.1.13]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.12...v0.1.13
[0.1.12]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.11...v0.1.12
[0.1.11]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.10...v0.1.11
[0.1.10]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.9...v0.1.10
>>>>>>> 3d4287c (Refactor GitHub issue templates and update CHANGELOG)
[0.1.9]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.8...v0.1.9
[0.1.8]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.7...v0.1.8
[0.1.7]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.6...v0.1.7
[0.1.6]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/matanbaruch/netbird-api-exporter/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/matanbaruch/netbird-api-exporter/releases/tag/v0.1.1
