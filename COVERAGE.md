# Code Coverage

[![Coverage Status](https://img.shields.io/badge/coverage-84.2%25-brightgreen)](./coverage.html)

## Quick Start

```bash
# Run tests with coverage
make test-all

# Check coverage locally
./scripts/check-coverage.sh

# Generate HTML report
./scripts/check-coverage.sh --html

# Validate coverage workflow
./scripts/test-coverage-workflow.sh
```

## Coverage Status

| Package | Coverage | Target |
|---------|----------|--------|
| **Overall** | **84.2%** | 80%+ |
| pkg/exporters | 95.5% | 90%+ |
| pkg/netbird | 100.0% | 95%+ |
| pkg/utils | 100.0% | 95%+ |

## CI/CD Integration

- ✅ GitHub Actions workflow for automated testing
- ✅ PR coverage comments with detailed analysis
- ✅ Coverage gates preventing quality regression
- ✅ Codecov integration for trend analysis
- ✅ Performance and integration test coverage

## Quick Commands

```bash
# Development
make test-unit                    # Unit tests only
make test-integration            # Integration tests (requires NETBIRD_API_TOKEN)
make test-performance            # Performance tests
./scripts/check-coverage.sh -t 85 # Custom threshold

# Reports
./scripts/check-coverage.sh --html --badge   # Generate reports
./scripts/check-coverage.sh --json          # JSON output
./scripts/check-coverage.sh --verbose       # Detailed breakdown
```

## Files

- `.github/workflows/coverage.yml` - CI/CD workflow
- `codecov.yml` - Codecov configuration
- `.coverage.yml` - Local coverage settings
- `scripts/check-coverage.sh` - Coverage validation script
- `scripts/test-coverage-workflow.sh` - Workflow validation
- `docs/coverage-monitoring.md` - Complete documentation

[📖 **Full Documentation**](./docs/coverage-monitoring.md)