#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TEST_DIR=$(mktemp -d)
ORIGINAL_DIR=$(pwd)

# Function to print colored output
print_status() {
    echo -e "${BLUE}[TEST-INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[TEST-SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[TEST-WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[TEST-ERROR]${NC} $1"
}

print_test_header() {
    echo ""
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}TEST: $1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Cleanup function
cleanup() {
    cd "$ORIGINAL_DIR"
    rm -rf "$TEST_DIR"
}

trap cleanup EXIT

# Test 1: Validate coverage workflow files exist
test_workflow_files() {
    print_test_header "Coverage Workflow Files Validation"
    
    local errors=0
    
    # Check GitHub Actions workflow
    if [ -f ".github/workflows/coverage.yml" ]; then
        print_success "GitHub Actions coverage workflow exists"
    else
        print_error "GitHub Actions coverage workflow missing"
        errors=$((errors + 1))
    fi
    
    # Check Codecov configuration
    if [ -f "codecov.yml" ]; then
        print_success "Codecov configuration exists"
    else
        print_error "Codecov configuration missing"
        errors=$((errors + 1))
    fi
    
    # Check coverage configuration
    if [ -f ".coverage.yml" ]; then
        print_success "Coverage configuration exists"
    else
        print_error "Coverage configuration missing"
        errors=$((errors + 1))
    fi
    
    # Check coverage check script
    if [ -x "scripts/check-coverage.sh" ]; then
        print_success "Coverage check script exists and is executable"
    else
        print_error "Coverage check script missing or not executable"
        errors=$((errors + 1))
    fi
    
    return $errors
}

# Test 2: Validate current test coverage
test_current_coverage() {
    print_test_header "Current Test Coverage Validation"
    
    print_status "Running tests with coverage..."
    GO111MODULE=on go test -coverprofile=test-coverage.out ./...
    
    if [ -f "test-coverage.out" ]; then
        local coverage=$(go tool cover -func=test-coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        print_status "Current coverage: $coverage%"
        
        if (( $(echo "$coverage >= 80" | bc -l) )); then
            print_success "Coverage meets minimum threshold (80%)"
            return 0
        else
            print_warning "Coverage below minimum threshold: $coverage% < 80%"
            return 1
        fi
    else
        print_error "Coverage profile not generated"
        return 1
    fi
}

# Test 3: Validate coverage check script
test_coverage_script() {
    print_test_header "Coverage Check Script Validation"
    
    # Generate a test coverage file
    GO111MODULE=on go test -coverprofile=script-test-coverage.out ./...
    
    if [ ! -f "script-test-coverage.out" ]; then
        print_error "Failed to generate test coverage file"
        return 1
    fi
    
    # Test script with default threshold
    print_status "Testing coverage script with default threshold..."
    if ./scripts/check-coverage.sh -f script-test-coverage.out --quiet; then
        print_success "Coverage script passed with default threshold"
    else
        print_warning "Coverage script failed with default threshold"
    fi
    
    # Test script with very high threshold (should fail)
    print_status "Testing coverage script with high threshold (should fail)..."
    if ./scripts/check-coverage.sh -f script-test-coverage.out -t 99.0 --quiet; then
        print_error "Coverage script should have failed with 99% threshold"
        return 1
    else
        print_success "Coverage script correctly failed with high threshold"
    fi
    
    # Test HTML generation
    print_status "Testing HTML report generation..."
    if ./scripts/check-coverage.sh -f script-test-coverage.out --html --quiet; then
        if [ -f "coverage.html" ]; then
            print_success "HTML coverage report generated successfully"
            rm -f coverage.html
        else
            print_error "HTML coverage report not found"
            return 1
        fi
    else
        print_error "Failed to generate HTML coverage report"
        return 1
    fi
    
    # Test JSON generation
    print_status "Testing JSON report generation..."
    if ./scripts/check-coverage.sh -f script-test-coverage.out --json --quiet; then
        if [ -f "coverage.json" ]; then
            print_success "JSON coverage report generated successfully"
            # Validate JSON format
            if command -v jq >/dev/null 2>&1; then
                if jq . coverage.json >/dev/null 2>&1; then
                    print_success "JSON coverage report is valid JSON"
                else
                    print_error "JSON coverage report is not valid JSON"
                    return 1
                fi
            fi
            rm -f coverage.json
        else
            print_error "JSON coverage report not found"
            return 1
        fi
    else
        print_error "Failed to generate JSON coverage report"
        return 1
    fi
    
    # Test badge generation
    print_status "Testing badge generation..."
    if ./scripts/check-coverage.sh -f script-test-coverage.out --badge --quiet; then
        if [ -f "coverage-badge.svg" ]; then
            print_success "Coverage badge generated successfully"
            rm -f coverage-badge.svg
        else
            print_error "Coverage badge not found"
            return 1
        fi
    else
        print_error "Failed to generate coverage badge"
        return 1
    fi
    
    # Cleanup
    rm -f script-test-coverage.out
    
    return 0
}

# Test 4: Validate GitHub Actions workflow syntax
test_workflow_syntax() {
    print_test_header "GitHub Actions Workflow Syntax Validation"
    
    # Check if the workflow YAML is valid
    if command -v yamllint >/dev/null 2>&1; then
        if yamllint .github/workflows/coverage.yml; then
            print_success "GitHub Actions workflow YAML is valid"
        else
            print_error "GitHub Actions workflow YAML has syntax errors"
            return 1
        fi
    else
        print_warning "yamllint not available, skipping YAML syntax validation"
    fi
    
    # Check for required jobs
    local required_jobs=("test-coverage" "integration-tests" "performance-tests" "coverage-comparison")
    local missing_jobs=()
    
    for job in "${required_jobs[@]}"; do
        if grep -q "^  $job:" .github/workflows/coverage.yml; then
            print_success "Required job '$job' found in workflow"
        else
            print_error "Required job '$job' missing from workflow"
            missing_jobs+=("$job")
        fi
    done
    
    if [ ${#missing_jobs[@]} -eq 0 ]; then
        print_success "All required jobs found in workflow"
        return 0
    else
        print_error "Missing jobs: ${missing_jobs[*]}"
        return 1
    fi
}

# Test 5: Validate Codecov configuration
test_codecov_config() {
    print_test_header "Codecov Configuration Validation"
    
    # Check for required sections in codecov.yml
    local required_sections=("codecov" "coverage" "comment")
    local missing_sections=()
    
    for section in "${required_sections[@]}"; do
        if grep -q "^$section:" codecov.yml; then
            print_success "Required section '$section' found in Codecov config"
        else
            print_error "Required section '$section' missing from Codecov config"
            missing_sections+=("$section")
        fi
    done
    
    # Check for coverage thresholds
    if grep -q "target:" codecov.yml; then
        print_success "Coverage targets configured in Codecov"
    else
        print_error "No coverage targets found in Codecov config"
        missing_sections+=("targets")
    fi
    
    if [ ${#missing_sections[@]} -eq 0 ]; then
        print_success "Codecov configuration is complete"
        return 0
    else
        print_error "Codecov configuration issues: ${missing_sections[*]}"
        return 1
    fi
}

# Test 6: Package-specific coverage validation
test_package_coverage() {
    print_test_header "Package-Specific Coverage Validation"
    
    # Generate coverage for analysis
    GO111MODULE=on go test -coverprofile=package-coverage.out ./...
    
    if [ ! -f "package-coverage.out" ]; then
        print_error "Failed to generate package coverage file"
        return 1
    fi
    
    # Expected high-coverage packages
    local high_coverage_packages=("netbird" "utils")
    local coverage_issues=0
    
    for package in "${high_coverage_packages[@]}"; do
        # Get package coverage (this is a simplified check)
        local package_line=$(go tool cover -func=package-coverage.out | grep "pkg/$package/" | head -1)
        if [ -n "$package_line" ]; then
            print_success "Package '$package' has test coverage"
        else
            print_warning "Package '$package' may not have comprehensive coverage"
            coverage_issues=$((coverage_issues + 1))
        fi
    done
    
    # Check that main package has some coverage consideration
    local main_coverage=$(go tool cover -func=package-coverage.out | grep -c "main.go" || echo "0")
    if [ "$main_coverage" -gt 0 ]; then
        print_status "Main package coverage detected"
    else
        print_status "Main package coverage not detected (expected for main.go)"
    fi
    
    rm -f package-coverage.out
    
    if [ $coverage_issues -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Test 7: Performance test coverage validation
test_performance_coverage() {
    print_test_header "Performance Test Coverage Validation"
    
    # Check if performance tests exist
    if find . -name "*_test.go" -exec grep -l "Performance\|Benchmark\|StressTest" {} \; | grep -q .; then
        print_success "Performance tests found"
        
        # Try to run performance tests
        print_status "Running performance tests..."
        if GO111MODULE=on go test -run="Performance|StressTest" -timeout=30s ./... >/dev/null 2>&1; then
            print_success "Performance tests executed successfully"
            return 0
        else
            print_warning "Performance tests execution had issues (may be timeout-related)"
            return 1
        fi
    else
        print_error "No performance tests found"
        return 1
    fi
}

# Main test execution
main() {
    print_status "Starting Coverage Workflow Validation Tests"
    echo "Test directory: $TEST_DIR"
    echo "Working directory: $ORIGINAL_DIR"
    echo ""
    
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    
    # Array of test functions
    local tests=(
        "test_workflow_files"
        "test_current_coverage"
        "test_coverage_script"
        "test_workflow_syntax"
        "test_codecov_config"
        "test_package_coverage"
        "test_performance_coverage"
    )
    
    # Run all tests
    for test_func in "${tests[@]}"; do
        total_tests=$((total_tests + 1))
        if $test_func; then
            passed_tests=$((passed_tests + 1))
        else
            failed_tests=$((failed_tests + 1))
        fi
    done
    
    # Summary
    echo ""
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}TEST SUMMARY${NC}"
    echo -e "${BLUE}================================${NC}"
    echo "Total tests: $total_tests"
    echo -e "${GREEN}Passed: $passed_tests${NC}"
    echo -e "${RED}Failed: $failed_tests${NC}"
    
    if [ $failed_tests -eq 0 ]; then
        echo ""
        print_success "All coverage workflow validation tests passed!"
        return 0
    else
        echo ""
        print_error "Some coverage workflow validation tests failed!"
        return 1
    fi
}

# Run main function
main "$@"