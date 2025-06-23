# Code Coverage Monitoring

This document describes the comprehensive code coverage monitoring system implemented for the NetBird API Exporter project.

## Overview

The coverage monitoring system provides:
- Automated coverage tracking on every pull request and push
- Coverage gates that prevent code quality degradation
- Detailed coverage reports and badges
- Integration with Codecov for trending analysis
- Local development tools for coverage validation

## System Components

### 1. GitHub Actions Workflow (`.github/workflows/coverage.yml`)

**Main Jobs:**
- **test-coverage**: Runs unit tests with coverage analysis
- **integration-tests**: Runs integration tests when API tokens are available
- **performance-tests**: Runs performance and benchmark tests
- **coverage-comparison**: Compares PR coverage against main branch

**Features:**
- Coverage threshold enforcement (80% minimum)
- Automatic PR comments with coverage details
- HTML and badge generation
- Codecov integration
- Artifact uploading for 30-day retention

### 2. Coverage Configuration Files

#### `.coverage.yml` - Local Coverage Configuration
```yaml
coverage:
  global_threshold: 80.0
  packages:
    "netbird-api-exporter/pkg/exporters": 90.0
    "netbird-api-exporter/pkg/netbird": 95.0
    "netbird-api-exporter/pkg/utils": 95.0
```

#### `codecov.yml` - Codecov Integration
- Project-level coverage targeting 80%
- Patch coverage targeting 80%
- Proper ignore patterns for test files
- PR comment configuration

### 3. Local Development Scripts

#### `scripts/check-coverage.sh` - Coverage Validation Script
```bash
# Basic usage
./scripts/check-coverage.sh

# With custom threshold
./scripts/check-coverage.sh -t 85.0

# Generate all reports
./scripts/check-coverage.sh --html --json --badge

# Compare with previous run
./scripts/check-coverage.sh --compare previous-coverage.out
```

#### `scripts/test-coverage-workflow.sh` - Workflow Validation
```bash
# Validate entire coverage system
./scripts/test-coverage-workflow.sh
```

## Current Coverage Status

**Baseline Coverage: 84.2%**

| Package | Coverage | Status |
|---------|----------|--------|
| pkg/exporters | 95.5% | ✅ Excellent |
| pkg/netbird | 100.0% | ✅ Perfect |
| pkg/utils | 100.0% | ✅ Perfect |
| main | 0.0% | ⚠️ Expected (integration testing) |

## Usage Guide

### For Developers

#### Running Tests Locally
```bash
# Run all tests with coverage
make test-all

# Run specific test types
make test-unit
make test-integration  # Requires NETBIRD_API_TOKEN
make test-performance
```

#### Checking Coverage Before Commit
```bash
# Quick coverage check
./scripts/check-coverage.sh

# Detailed coverage with reports
./scripts/check-coverage.sh --html --verbose
```

#### Improving Coverage
1. **Identify uncovered code:**
   ```bash
   ./scripts/check-coverage.sh --html
   # Open coverage.html in browser
   ```

2. **Focus on critical functions:**
   - API endpoints in `pkg/exporters/`
   - Business logic in `pkg/netbird/`
   - Utility functions in `pkg/utils/`

3. **Add meaningful tests:**
   - Unit tests for individual functions
   - Integration tests for API interactions
   - Error handling scenarios

### For CI/CD

#### Workflow Triggers
- **Push to main/develop**: Full coverage analysis
- **Pull Requests**: Coverage analysis with comparison
- **Integration tests**: Only when `NETBIRD_API_TOKEN` is available

#### Required Secrets
```yaml
# GitHub Repository Secrets
CODECOV_TOKEN: <codecov-upload-token>
NETBIRD_API_TOKEN: <optional-for-integration-tests>
```

#### Workflow Outputs
- Coverage percentage and status
- HTML coverage reports
- Coverage badges
- PR comments with detailed analysis
- Codecov integration for trending

## Coverage Gates and Quality Control

### Automatic Failing Conditions
1. **Overall coverage below 80%**
2. **Coverage decrease > 1% from main branch**
3. **Missing critical test coverage in new code**

### Warning Conditions
1. **Coverage between 78-80%**
2. **Significant coverage changes without explanation**

### Manual Review Required
1. **Integration test failures**
2. **Performance test degradation**
3. **Coverage configuration changes**

## Troubleshooting

### Common Issues

#### "Coverage file not found"
```bash
# Generate coverage file first
go test -coverprofile=coverage.out ./...
```

#### "Module path issues"
```bash
# Ensure GO111MODULE is enabled
GO111MODULE=on go test -coverprofile=coverage.out ./...
```

#### "Integration tests skipped"
```bash
# Set API token for integration tests
export NETBIRD_API_TOKEN="your-token-here"
make test-integration
```

#### "Codecov upload fails"
- Verify `CODECOV_TOKEN` is set in repository secrets
- Check Codecov project configuration
- Review network connectivity in CI environment

### Debug Commands

```bash
# Check coverage file contents
go tool cover -func=coverage.out

# Generate HTML report manually
go tool cover -html=coverage.out -o coverage.html

# Test coverage script
./scripts/check-coverage.sh --verbose

# Validate workflow configuration
./scripts/test-coverage-workflow.sh
```

## Best Practices

### Writing Tests for Coverage
1. **Test all public functions**
2. **Cover error paths and edge cases**
3. **Use table-driven tests for multiple scenarios**
4. **Mock external dependencies**
5. **Test concurrent behavior where applicable**

### Coverage Strategy
1. **Aim for 90%+ on critical business logic**
2. **80%+ overall project coverage**
3. **100% on utility and helper functions**
4. **Focus on quality over quantity**

### Maintenance
1. **Review coverage trends monthly**
2. **Update thresholds as code matures**
3. **Add tests for new features immediately**
4. **Remove obsolete test exclusions**

## Integration with Development Workflow

### Pre-commit Hooks
```bash
# Install coverage check as pre-commit hook
make install-hooks
```

### IDE Integration
Most IDEs support Go coverage visualization:
- VSCode: Go extension with coverage highlighting
- GoLand: Built-in coverage analysis
- Vim/Neovim: Various coverage plugins

### Continuous Monitoring
- Codecov dashboard for trend analysis
- GitHub Actions for immediate feedback
- Weekly coverage reports (can be automated)

## Advanced Configuration

### Custom Coverage Thresholds
Edit `.coverage.yml` to set package-specific thresholds:
```yaml
packages:
  "your/critical/package": 95.0
  "your/standard/package": 80.0
```

### Excluding Files from Coverage
```yaml
exclude:
  - "**/*_test.go"
  - "**/testdata/**"
  - "**/generated/**"
```

### Notification Integration
Configure Slack/email notifications for coverage changes:
```yaml
notifications:
  slack:
    enabled: true
    webhook_url: "your-webhook-url"
```

## Metrics and Reporting

### Available Reports
1. **HTML Coverage Report** (`coverage.html`)
2. **JSON Coverage Data** (`coverage.json`)
3. **Coverage Badge** (`coverage-badge.svg`)
4. **Text Summary** (console output)

### Key Metrics
- Overall coverage percentage
- Per-package coverage breakdown
- Coverage trend over time
- Test execution time
- Code quality indicators

## Support and Maintenance

### Regular Tasks
- [ ] Review coverage reports weekly
- [ ] Update coverage thresholds quarterly
- [ ] Verify CI/CD pipeline monthly
- [ ] Update dependencies as needed

### Getting Help
1. Check this documentation first
2. Run validation scripts for diagnostics
3. Review GitHub Actions logs
4. Check Codecov dashboard for insights

---

For additional questions or issues, please:
- Open a GitHub issue with the `coverage` label
- Include relevant logs and configuration
- Describe expected vs actual behavior