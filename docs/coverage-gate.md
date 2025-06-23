# Coverage Gate Implementation

This document describes the comprehensive coverage gate system implemented for the NetBird API Exporter project.

## Overview

The coverage gate ensures that pull requests maintain or improve code coverage, preventing regressions in test coverage. When coverage decreases beyond an acceptable threshold, the PR is automatically blocked from merging.

## Components

### 1. GitHub Actions Workflow (`.github/workflows/test.yml`)

#### Enhanced Unit Tests Job
- Extracts coverage percentage from `go tool cover`
- Outputs coverage metrics for downstream jobs
- Generates coverage badge color based on thresholds

#### Coverage Gate Job
- **Trigger**: Only runs on pull requests
- **Dependencies**: Requires unit tests to pass first
- **Permissions**: Can write to PRs, statuses, and checks

**Process Flow**:
1. Checkout PR branch with full git history
2. Run unit tests and extract PR coverage
3. Checkout base branch and extract base coverage
4. Calculate coverage difference using `bc` arithmetic
5. Generate detailed coverage report in Markdown
6. Update/create PR comment with coverage analysis
7. Set GitHub status check (`coverage/gate`)
8. Fail job if coverage decreases beyond threshold

#### Test Summary Integration
- Includes coverage gate results in overall test summary
- Shows coverage percentage and badge information
- Treats coverage gate failure as critical test failure

### 2. Configuration File (`.github/coverage-gate.yml`)

Comprehensive configuration supporting:

- **Thresholds**: Minimum coverage, max decrease, target coverage
- **Scope**: Include/exclude file patterns
- **Gate Settings**: Failure conditions and exemptions
- **Reporting**: PR comments, package details, badges
- **Status Checks**: Context names and templates
- **Integrations**: GitHub labels and status checks
- **Advanced**: Precision, timeouts, debug options

### 3. Local Testing Script (`scripts/coverage-gate.sh`)

Features:
- Test coverage gate logic locally before pushing
- Compare any two branches
- Configurable thresholds
- Detailed reporting with package breakdown
- Clean git state management (stashing/unstashing)

**Usage Examples**:
```bash
# Test current branch against main
./scripts/coverage-gate.sh

# Test against different base branch
./scripts/coverage-gate.sh -b develop

# Use custom threshold
./scripts/coverage-gate.sh -t 1.0

# Generate detailed report
./scripts/coverage-gate.sh -r
```

## Coverage Gate Logic

### Threshold Calculation
```bash
# Extract coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Calculate difference
DIFF=$(echo "$PR_COVERAGE - $BASE_COVERAGE" | bc -l)

# Check threshold (default: -0.5%)
DECREASED=$(echo "$DIFF < -0.5" | bc -l)
```

### Failure Conditions
The coverage gate fails when:
1. Coverage decreases by more than 0.5% (configurable)
2. Base branch tests fail to run
3. PR branch tests fail to run
4. Coverage calculation fails

### Success Conditions
The coverage gate passes when:
1. Coverage improves (positive difference)
2. Coverage stays the same (zero difference)
3. Coverage decreases but within acceptable threshold

## Reporting Features

### PR Comments
Automatically generated coverage reports include:
- Base vs PR coverage comparison
- Percentage difference calculation
- Pass/fail status with clear visual indicators
- Package-level coverage breakdown
- Actionable recommendations for failed gates

### Status Checks
- GitHub status check with context `coverage/gate`
- Required status check for branch protection
- Descriptive messages showing coverage changes
- Links to detailed check results

### Test Summary Integration
- Coverage information in workflow summaries
- Badge generation based on coverage thresholds
- Integration with overall test results

## Configuration Options

### Thresholds (in `coverage-gate.yml`)
```yaml
thresholds:
  minimum_coverage: 60.0      # Minimum acceptable coverage
  max_decrease: 0.5           # Maximum allowed decrease
  target_coverage: 80.0       # Project target
  warning_threshold: 70.0     # Warning level
```

### File Exemptions
```yaml
coverage_scope:
  exclude:
    - "**/*_test.go"           # Test files
    - "**/mock_*.go"           # Mock files
    - "cmd/main.go"            # Entry points
```

### Gate Behavior
```yaml
gate_settings:
  fail_on_decrease: true       # Fail on any decrease
  fail_on_low_coverage: false  # Don't fail on absolute coverage
```

## Branch Protection Setup

To enforce coverage gates, configure branch protection rules:

1. **Required Status Checks**: Add `coverage/gate`
2. **Require branches to be up to date**: Enabled
3. **Restrictions**: Configure as needed

## Monitoring and Maintenance

### Adjusting Thresholds
- Monitor coverage trends over time
- Adjust `max_decrease` based on team workflow
- Consider project complexity when setting targets

### Exemption Management
- Regular review of exempted files
- Add exemptions for legitimate cases (CLI entry points, generated code)
- Document exemption reasons

### Performance Considerations
- Coverage calculation runs twice per PR (base + current)
- Uses git checkout which requires clean working directory
- Caching enabled to speed up subsequent runs

## Troubleshooting

### Common Issues

1. **"Tests failed for branch"**
   - Check if tests pass locally on the branch
   - Ensure all dependencies are available in CI

2. **"Coverage Gate Failed" with small changes**
   - Check if the threshold is too strict
   - Consider if new code needs more tests

3. **Git checkout failures**
   - Ensure full git history is available (`fetch-depth: 0`)
   - Check for uncommitted changes

### Debug Mode
Enable debug logging in `coverage-gate.yml`:
```yaml
advanced:
  debug: true
```

## Benefits

1. **Quality Assurance**: Prevents coverage regressions
2. **Automated Enforcement**: No manual coverage reviews needed
3. **Developer Feedback**: Clear guidance on coverage requirements
4. **Flexibility**: Configurable thresholds and exemptions
5. **Local Testing**: Test coverage gates before pushing
6. **Integration**: Works with existing CI/CD pipeline

## Future Enhancements

- Coverage trending over time
- Integration with external coverage services
- File-level coverage requirements
- Coverage improvement goals and tracking
- Automated exemption suggestions based on file patterns