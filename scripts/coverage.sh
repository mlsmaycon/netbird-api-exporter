#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Coverage configuration
COVERAGE_DIR="coverage"
UNIT_PROFILE="$COVERAGE_DIR/unit.out"
INTEGRATION_PROFILE="$COVERAGE_DIR/integration.out"
MERGED_PROFILE="$COVERAGE_DIR/coverage.out"
MINIMUM_COVERAGE=80.0
PACKAGE_PATH="./..."

# Coverage thresholds by package (simplified for POSIX compatibility)
PKG_EXPORTERS_THRESHOLD=85.0
PKG_NETBIRD_THRESHOLD=80.0
PKG_UTILS_THRESHOLD=90.0

# Files/patterns to exclude from coverage
EXCLUDE_PATTERNS=(
    "*_test.go"
    "*/testdata/*"
    "*/vendor/*"
    "main.go"
)

# Function to print colored output
print_status() {
    echo -e "${BLUE}[COVERAGE]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to create coverage directory
setup_coverage_dir() {
    print_status "Setting up coverage directory..."
    mkdir -p "$COVERAGE_DIR"
    rm -f "$COVERAGE_DIR"/*.out "$COVERAGE_DIR"/*.html "$COVERAGE_DIR"/*.lcov "$COVERAGE_DIR"/*.json "$COVERAGE_DIR"/*.xml
}

# Function to run unit tests with coverage
run_unit_coverage() {
    print_status "Running unit tests with coverage..."
    go test -v -race -timeout=30s -coverprofile="$UNIT_PROFILE" -covermode=atomic -skip="Integration_" "$PACKAGE_PATH"
    
    if [ -f "$UNIT_PROFILE" ]; then
        print_success "Unit test coverage profile generated: $UNIT_PROFILE"
    else
        print_error "Failed to generate unit test coverage profile"
        return 1
    fi
}

# Function to run integration tests with coverage (if API token available)
run_integration_coverage() {
    if [ -z "$NETBIRD_API_TOKEN" ]; then
        print_warning "NETBIRD_API_TOKEN not set - skipping integration test coverage"
        return 0
    fi
    
    print_status "Running integration tests with coverage..."
    go test -v -timeout=5m -coverprofile="$INTEGRATION_PROFILE" -covermode=atomic -run="Integration_" "$PACKAGE_PATH"
    
    if [ -f "$INTEGRATION_PROFILE" ]; then
        print_success "Integration test coverage profile generated: $INTEGRATION_PROFILE"
    else
        print_warning "No integration test coverage profile generated"
    fi
}

# Function to merge coverage profiles
merge_coverage_profiles() {
    print_status "Merging coverage profiles..."
    
    local profiles=()
    [ -f "$UNIT_PROFILE" ] && profiles+=("$UNIT_PROFILE")
    [ -f "$INTEGRATION_PROFILE" ] && profiles+=("$INTEGRATION_PROFILE")
    
    if [ ${#profiles[@]} -eq 0 ]; then
        print_error "No coverage profiles found to merge"
        return 1
    elif [ ${#profiles[@]} -eq 1 ]; then
        print_status "Only one profile found, copying to merged profile..."
        cp "${profiles[0]}" "$MERGED_PROFILE"
    else
        print_status "Merging ${#profiles[@]} coverage profiles..."
        # Use gocovmerge if available, otherwise use go tool cover
        if command -v gocovmerge >/dev/null 2>&1; then
            gocovmerge "${profiles[@]}" > "$MERGED_PROFILE"
        else
            print_warning "gocovmerge not found, using manual merge (less accurate)"
            cat "${profiles[@]}" | grep -v "mode:" > "$MERGED_PROFILE.tmp"
            echo "mode: atomic" > "$MERGED_PROFILE"
            cat "$MERGED_PROFILE.tmp" >> "$MERGED_PROFILE"
            rm -f "$MERGED_PROFILE.tmp"
        fi
    fi
    
    print_success "Coverage profiles merged: $MERGED_PROFILE"
}

# Function to generate coverage reports in multiple formats
generate_coverage_reports() {
    if [ ! -f "$MERGED_PROFILE" ]; then
        print_error "Merged coverage profile not found: $MERGED_PROFILE"
        return 1
    fi
    
    print_status "Generating coverage reports..."
    
    # HTML report
    go tool cover -html="$MERGED_PROFILE" -o "$COVERAGE_DIR/coverage.html"
    print_success "HTML coverage report: $COVERAGE_DIR/coverage.html"
    
    # Text report with function-level coverage
    go tool cover -func="$MERGED_PROFILE" > "$COVERAGE_DIR/coverage.txt"
    print_success "Text coverage report: $COVERAGE_DIR/coverage.txt"
    
    # LCOV format (for CI/CD integration)
    if command -v gcov2lcov >/dev/null 2>&1; then
        gcov2lcov -infile="$MERGED_PROFILE" -outfile="$COVERAGE_DIR/coverage.lcov"
        print_success "LCOV coverage report: $COVERAGE_DIR/coverage.lcov"
    else
        print_warning "gcov2lcov not found, skipping LCOV report generation"
        print_status "Install with: go install github.com/jandelgado/gcov2lcov@latest"
    fi
    
    # JSON format (for programmatic processing)
    generate_json_report
    
    # XML format (for some CI systems)
    generate_xml_report
}

# Function to generate JSON coverage report
generate_json_report() {
    print_status "Generating JSON coverage report..."
    
    local total_coverage
    total_coverage=$(go tool cover -func="$MERGED_PROFILE" | grep total | awk '{print $3}' | sed 's/%//')
    
    local timestamp
    timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    cat > "$COVERAGE_DIR/coverage.json" << EOF
{
  "timestamp": "$timestamp",
  "total_coverage": $total_coverage,
  "minimum_threshold": $MINIMUM_COVERAGE,
  "coverage_file": "$MERGED_PROFILE",
  "reports": {
    "html": "$COVERAGE_DIR/coverage.html",
    "text": "$COVERAGE_DIR/coverage.txt",
    "lcov": "$COVERAGE_DIR/coverage.lcov",
    "xml": "$COVERAGE_DIR/coverage.xml"
  },
  "package_coverage": $(generate_package_coverage_json)
}
EOF
    
    print_success "JSON coverage report: $COVERAGE_DIR/coverage.json"
}

# Function to generate package-level coverage JSON
generate_package_coverage_json() {
    local json="{"
    local first=true
    
    while IFS= read -r line; do
        if [[ $line == *"/"* ]] && [[ $line != *"total:"* ]]; then
            local package=$(echo "$line" | awk '{print $1}')
            local coverage=$(echo "$line" | awk '{print $3}' | sed 's/%//')
            
            if [ "$first" = true ]; then
                first=false
            else
                json+=","
            fi
            json+="\"$package\":$coverage"
        fi
    done < <(go tool cover -func="$MERGED_PROFILE" | grep -E "^[^[:space:]]" | head -n -1)
    
    json+="}"
    echo "$json"
}

# Function to generate XML coverage report (Cobertura format)
generate_xml_report() {
    print_status "Generating XML coverage report..."
    
    if command -v gocover-cobertura >/dev/null 2>&1; then
        gocover-cobertura < "$MERGED_PROFILE" > "$COVERAGE_DIR/coverage.xml"
        print_success "XML coverage report: $COVERAGE_DIR/coverage.xml"
    else
        print_warning "gocover-cobertura not found, skipping XML report generation"
        print_status "Install with: go install github.com/boumenot/gocover-cobertura@latest"
    fi
}

# Function to check coverage thresholds
check_coverage_thresholds() {
    print_status "Checking coverage thresholds..."
    
    local overall_coverage
    overall_coverage=$(go tool cover -func="$MERGED_PROFILE" | grep total | awk '{print $3}' | sed 's/%//')
    
    print_status "Overall coverage: ${overall_coverage}%"
    print_status "Minimum threshold: ${MINIMUM_COVERAGE}%"
    
    # Check overall threshold
    if (( $(echo "$overall_coverage < $MINIMUM_COVERAGE" | bc -l) )); then
        print_error "Overall coverage ${overall_coverage}% is below minimum threshold ${MINIMUM_COVERAGE}%"
        return 1
    else
        print_success "Overall coverage ${overall_coverage}% meets minimum threshold ${MINIMUM_COVERAGE}%"
    fi
    
    # Check package-specific thresholds
    local threshold_failures=0
    
    # Check pkg/exporters threshold
    local exporters_coverage
    exporters_coverage=$(go tool cover -func="$MERGED_PROFILE" | grep "pkg/exporters" | awk '{sum+=$3; count++} END {if(count>0) print sum/count; else print 0}' | sed 's/%//')
    if [ -n "$exporters_coverage" ] && [ "$exporters_coverage" != "0" ]; then
        print_status "Package ./pkg/exporters coverage: ${exporters_coverage}% (threshold: ${PKG_EXPORTERS_THRESHOLD}%)"
        if (( $(echo "$exporters_coverage < $PKG_EXPORTERS_THRESHOLD" | bc -l) )); then
            print_error "Package ./pkg/exporters coverage ${exporters_coverage}% is below threshold ${PKG_EXPORTERS_THRESHOLD}%"
            threshold_failures=$((threshold_failures + 1))
        fi
    fi
    
    # Check pkg/netbird threshold
    local netbird_coverage
    netbird_coverage=$(go tool cover -func="$MERGED_PROFILE" | grep "pkg/netbird" | awk '{sum+=$3; count++} END {if(count>0) print sum/count; else print 0}' | sed 's/%//')
    if [ -n "$netbird_coverage" ] && [ "$netbird_coverage" != "0" ]; then
        print_status "Package ./pkg/netbird coverage: ${netbird_coverage}% (threshold: ${PKG_NETBIRD_THRESHOLD}%)"
        if (( $(echo "$netbird_coverage < $PKG_NETBIRD_THRESHOLD" | bc -l) )); then
            print_error "Package ./pkg/netbird coverage ${netbird_coverage}% is below threshold ${PKG_NETBIRD_THRESHOLD}%"
            threshold_failures=$((threshold_failures + 1))
        fi
    fi
    
    # Check pkg/utils threshold
    local utils_coverage
    utils_coverage=$(go tool cover -func="$MERGED_PROFILE" | grep "pkg/utils" | awk '{sum+=$3; count++} END {if(count>0) print sum/count; else print 0}' | sed 's/%//')
    if [ -n "$utils_coverage" ] && [ "$utils_coverage" != "0" ]; then
        print_status "Package ./pkg/utils coverage: ${utils_coverage}% (threshold: ${PKG_UTILS_THRESHOLD}%)"
        if (( $(echo "$utils_coverage < $PKG_UTILS_THRESHOLD" | bc -l) )); then
            print_error "Package ./pkg/utils coverage ${utils_coverage}% is below threshold ${PKG_UTILS_THRESHOLD}%"
            threshold_failures=$((threshold_failures + 1))
        fi
    fi
    
    if [ $threshold_failures -gt 0 ]; then
        print_error "$threshold_failures package(s) failed coverage thresholds"
        return 1
    else
        print_success "All package coverage thresholds met"
    fi
}

# Function to display coverage summary
display_coverage_summary() {
    print_status "Coverage Summary:"
    echo "=================="
    
    if [ -f "$COVERAGE_DIR/coverage.txt" ]; then
        tail -n 1 "$COVERAGE_DIR/coverage.txt"
    fi
    
    echo ""
    print_status "Reports generated in $COVERAGE_DIR/:"
    ls -la "$COVERAGE_DIR"/ | grep -E '\.(html|txt|lcov|json|xml)$' || true
    
    echo ""
    print_status "View HTML report: open $COVERAGE_DIR/coverage.html"
}

# Function to clean up coverage files
cleanup_coverage() {
    print_status "Cleaning up coverage files..."
    rm -rf "$COVERAGE_DIR"
    print_success "Coverage cleanup completed"
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  generate    Generate coverage reports (default)"
    echo "  check       Check coverage thresholds only"
    echo "  clean       Clean up coverage files"
    echo ""
    echo "Options:"
    echo "  -t, --threshold PERCENT     Set minimum coverage threshold (default: $MINIMUM_COVERAGE)"
    echo "  --unit-only                 Run unit tests only"
    echo "  --integration-only          Run integration tests only (requires NETBIRD_API_TOKEN)"
    echo "  --no-merge                  Don't merge coverage profiles"
    echo "  --no-threshold-check        Skip coverage threshold checks"
    echo "  -h, --help                  Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  NETBIRD_API_TOKEN          API token for integration tests"
    echo "  MINIMUM_COVERAGE           Minimum coverage percentage override"
    echo ""
    echo "Examples:"
    echo "  $0                         Generate full coverage report"
    echo "  $0 --unit-only             Generate coverage for unit tests only"
    echo "  $0 --threshold 85          Set minimum coverage to 85%"
    echo "  $0 check                   Check thresholds against existing coverage"
}

# Parse command line arguments
COMMAND="generate"
UNIT_ONLY=false
INTEGRATION_ONLY=false
NO_MERGE=false
NO_THRESHOLD_CHECK=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -t|--threshold)
            MINIMUM_COVERAGE="$2"
            shift 2
            ;;
        --unit-only)
            UNIT_ONLY=true
            shift
            ;;
        --integration-only)
            INTEGRATION_ONLY=true
            shift
            ;;
        --no-merge)
            NO_MERGE=true
            shift
            ;;
        --no-threshold-check)
            NO_THRESHOLD_CHECK=true
            shift
            ;;
        generate|check|clean)
            COMMAND="$1"
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Override minimum coverage from environment if set
if [ -n "$MINIMUM_COVERAGE_ENV" ]; then
    MINIMUM_COVERAGE="$MINIMUM_COVERAGE_ENV"
fi

# Main execution
print_status "NetBird API Exporter Coverage Tool"
print_status "Command: $COMMAND"
print_status "Minimum coverage threshold: ${MINIMUM_COVERAGE}%"

case $COMMAND in
    generate)
        setup_coverage_dir
        
        if [ "$INTEGRATION_ONLY" = false ]; then
            run_unit_coverage
        fi
        
        if [ "$UNIT_ONLY" = false ]; then
            run_integration_coverage
        fi
        
        if [ "$NO_MERGE" = false ]; then
            merge_coverage_profiles
            generate_coverage_reports
        fi
        
        if [ "$NO_THRESHOLD_CHECK" = false ]; then
            check_coverage_thresholds
        fi
        
        display_coverage_summary
        ;;
    check)
        if [ ! -f "$MERGED_PROFILE" ]; then
            print_error "Coverage profile not found: $MERGED_PROFILE"
            print_status "Run '$0 generate' first to create coverage data"
            exit 1
        fi
        check_coverage_thresholds
        ;;
    clean)
        cleanup_coverage
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac

print_success "Coverage operation completed successfully!"