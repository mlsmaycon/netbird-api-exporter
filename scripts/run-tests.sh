#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
PACKAGE_PATH="./..."
INTEGRATION_PACKAGE="./pkg/..."
UNIT_TIMEOUT="30s"
INTEGRATION_TIMEOUT="5m"
COVER_PROFILE="coverage.out"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
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

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    
    if go test -v -timeout "$UNIT_TIMEOUT" -race -coverprofile="$COVER_PROFILE" -skip="Integration_" "$PACKAGE_PATH"; then
        print_success "Unit tests passed!"
        
        # Generate coverage report
        if [ -f "$COVER_PROFILE" ]; then
            coverage=$(go tool cover -func="$COVER_PROFILE" | grep total | awk '{print $3}')
            print_status "Test coverage: $coverage"
            
            # Generate HTML coverage report
            go tool cover -html="$COVER_PROFILE" -o coverage.html
            print_status "Coverage report generated: coverage.html"
        fi
    else
        print_error "Unit tests failed!"
        return 1
    fi
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    
    if [ -z "$NETBIRD_API_TOKEN" ]; then
        print_warning "NETBIRD_API_TOKEN not set - integration tests will be skipped"
        print_status "Set NETBIRD_API_TOKEN environment variable to run integration tests"
        return 0
    fi
    
    print_status "NETBIRD_API_TOKEN is set - running integration tests"
    
    if go test -v -timeout "$INTEGRATION_TIMEOUT" -run="Integration_" "$INTEGRATION_PACKAGE"; then
        print_success "Integration tests passed!"
    else
        print_error "Integration tests failed!"
        return 1
    fi
}

# Function to run benchmark tests
run_benchmark_tests() {
    print_status "Running benchmark tests..."
    
    if go test -v -bench=. -benchmem -run=^$ "$PACKAGE_PATH"; then
        print_success "Benchmark tests completed!"
    else
        print_error "Benchmark tests failed!"
        return 1
    fi
}

# Function to run performance tests
run_performance_tests() {
    print_status "Running performance tests..."
    
    if go test -v -timeout="5m" -run="Performance|StressTest" "$PACKAGE_PATH"; then
        print_success "Performance tests passed!"
    else
        print_error "Performance tests failed!"
        return 1
    fi
}

# Function to clean up test artifacts
cleanup() {
    print_status "Cleaning up test artifacts..."
    rm -f "$COVER_PROFILE" coverage.html
    print_success "Cleanup completed!"
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS] [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  unit                Run unit tests only"
    echo "  integration         Run integration tests only (requires NETBIRD_API_TOKEN)"
    echo "  benchmark           Run benchmark tests only"
    echo "  performance         Run performance tests only"
    echo "  all                 Run all tests (default)"
    echo "  clean               Clean up test artifacts"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -v, --verbose       Enable verbose output"
    echo "  --no-cache          Disable Go test cache"
    echo "  --timeout DURATION  Set timeout for tests (default: $UNIT_TIMEOUT for unit, $INTEGRATION_TIMEOUT for integration)"
    echo ""
    echo "Environment Variables:"
    echo "  NETBIRD_API_TOKEN   API token for NetBird Cloud (required for integration tests)"
    echo "  NETBIRD_API_URL     Custom NetBird API URL (default: https://api.netbird.io)"
    echo ""
    echo "Examples:"
    echo "  $0 unit                          Run only unit tests"
    echo "  $0 integration                   Run only integration tests"
    echo "  NETBIRD_API_TOKEN=xxx $0 all     Run all tests with API token"
    echo "  $0 --timeout 1m unit             Run unit tests with 1 minute timeout"
}

# Parse command line arguments
VERBOSE=false
NO_CACHE=""
COMMAND="all"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        --no-cache)
            NO_CACHE="-count=1"
            shift
            ;;
        --timeout)
            UNIT_TIMEOUT="$2"
            INTEGRATION_TIMEOUT="$2"
            shift 2
            ;;
        unit|integration|benchmark|performance|all|clean)
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

# Add no-cache flag to go test commands if specified
if [ -n "$NO_CACHE" ]; then
    print_status "Test cache disabled"
fi

# Set verbose flag for go test if specified
VERBOSE_FLAG=""
if [ "$VERBOSE" = true ]; then
    VERBOSE_FLAG="-v"
    print_status "Verbose mode enabled"
fi

# Main execution
print_status "Starting NetBird API Exporter test suite..."
print_status "Command: $COMMAND"

case $COMMAND in
    unit)
        run_unit_tests
        ;;
    integration)
        run_integration_tests
        ;;
    benchmark)
        run_benchmark_tests
        ;;
    performance)
        run_performance_tests
        ;;
    all)
        run_unit_tests
        echo ""
        run_integration_tests
        echo ""
        run_performance_tests
        ;;
    clean)
        cleanup
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac

print_success "Test suite completed successfully!" 