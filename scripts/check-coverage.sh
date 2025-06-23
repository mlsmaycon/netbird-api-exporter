#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_THRESHOLD=80.0
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
COVERAGE_JSON="coverage.json"

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

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -t, --threshold FLOAT   Set coverage threshold (default: $COVERAGE_THRESHOLD)"
    echo "  -f, --file FILE         Coverage profile file (default: $COVERAGE_FILE)"
    echo "  -o, --output DIR        Output directory for reports (default: current dir)"
    echo "  -q, --quiet             Quiet mode - only show final result"
    echo "  -v, --verbose           Verbose mode - show detailed package coverage"
    echo "  -j, --json              Generate JSON coverage report"
    echo "  --html                  Generate HTML coverage report"
    echo "  --badge                 Generate coverage badge"
    echo "  --compare FILE          Compare with previous coverage file"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Run with defaults"
    echo "  $0 -t 85.0 --html --badge           # Custom threshold with reports"
    echo "  $0 --compare previous-coverage.out   # Compare with previous run"
}

# Parse command line arguments
QUIET=false
VERBOSE=false
GENERATE_JSON=false
GENERATE_HTML=false
GENERATE_BADGE=false
OUTPUT_DIR="."
COMPARE_FILE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--threshold)
            COVERAGE_THRESHOLD="$2"
            shift 2
            ;;
        -f|--file)
            COVERAGE_FILE="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -j|--json)
            GENERATE_JSON=true
            shift
            ;;
        --html)
            GENERATE_HTML=true
            shift
            ;;
        --badge)
            GENERATE_BADGE=true
            shift
            ;;
        --compare)
            COMPARE_FILE="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Check if coverage file exists
if [ ! -f "$COVERAGE_FILE" ]; then
    print_error "Coverage file '$COVERAGE_FILE' not found"
    print_status "Run 'go test -coverprofile=$COVERAGE_FILE ./...' first"
    exit 1
fi

# Extract overall coverage
OVERALL_COVERAGE=$(GO111MODULE=on go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

if [ -z "$OVERALL_COVERAGE" ]; then
    print_error "Could not extract coverage percentage from $COVERAGE_FILE"
    exit 1
fi

[ "$QUIET" = false ] && print_status "Overall coverage: $OVERALL_COVERAGE%"

# Check threshold
if (( $(echo "$OVERALL_COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
    THRESHOLD_STATUS="FAIL"
    THRESHOLD_COLOR="$RED"
else
    THRESHOLD_STATUS="PASS"
    THRESHOLD_COLOR="$GREEN"
fi

# Verbose mode - show per-package coverage
if [ "$VERBOSE" = true ] && [ "$QUIET" = false ]; then
    print_status "Per-package coverage breakdown:"
    GO111MODULE=on go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | while read line; do
        echo "  $line"
    done
    echo ""
fi

# Generate HTML report
if [ "$GENERATE_HTML" = true ]; then
    HTML_FILE="$OUTPUT_DIR/$COVERAGE_HTML"
    GO111MODULE=on go tool cover -html="$COVERAGE_FILE" -o "$HTML_FILE"
    [ "$QUIET" = false ] && print_success "HTML coverage report generated: $HTML_FILE"
fi

# Generate JSON report
if [ "$GENERATE_JSON" = true ]; then
    JSON_FILE="$OUTPUT_DIR/$COVERAGE_JSON"
    
    # Create JSON coverage report
    cat > "$JSON_FILE" << EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "overall_coverage": $OVERALL_COVERAGE,
  "threshold": $COVERAGE_THRESHOLD,
  "status": "$THRESHOLD_STATUS",
  "packages": [
EOF
    
    # Add package-level coverage
    FIRST=true
    GO111MODULE=on go tool cover -func="$COVERAGE_FILE" | grep -v "total:" | while read file func coverage; do
        package=$(dirname "$file")
        if [ "$FIRST" = true ]; then
            FIRST=false
        else
            echo "," >> "$JSON_FILE"
        fi
        echo "    {\"package\": \"$package\", \"file\": \"$file\", \"function\": \"$func\", \"coverage\": \"$coverage\"}" >> "$JSON_FILE"
    done
    
    cat >> "$JSON_FILE" << EOF
  ]
}
EOF
    
    [ "$QUIET" = false ] && print_success "JSON coverage report generated: $JSON_FILE"
fi

# Generate coverage badge
if [ "$GENERATE_BADGE" = true ]; then
    # Determine badge color based on coverage
    if (( $(echo "$OVERALL_COVERAGE >= 90" | bc -l) )); then
        BADGE_COLOR="brightgreen"
    elif (( $(echo "$OVERALL_COVERAGE >= 80" | bc -l) )); then
        BADGE_COLOR="green"
    elif (( $(echo "$OVERALL_COVERAGE >= 70" | bc -l) )); then
        BADGE_COLOR="yellow"
    elif (( $(echo "$OVERALL_COVERAGE >= 60" | bc -l) )); then
        BADGE_COLOR="orange"
    else
        BADGE_COLOR="red"
    fi
    
    BADGE_FILE="$OUTPUT_DIR/coverage-badge.svg"
    COVERAGE_ENCODED=$(echo "$OVERALL_COVERAGE" | sed 's/\./%2E/g')
    curl -s "https://img.shields.io/badge/coverage-${COVERAGE_ENCODED}%25-${BADGE_COLOR}" > "$BADGE_FILE"
    [ "$QUIET" = false ] && print_success "Coverage badge generated: $BADGE_FILE"
fi

# Compare with previous coverage if provided
if [ -n "$COMPARE_FILE" ] && [ -f "$COMPARE_FILE" ]; then
    PREVIOUS_COVERAGE=$(GO111MODULE=on go tool cover -func="$COMPARE_FILE" | grep total | awk '{print $3}' | sed 's/%//')
    COVERAGE_DIFF=$(echo "$OVERALL_COVERAGE - $PREVIOUS_COVERAGE" | bc -l)
    
    if [ "$QUIET" = false ]; then
        print_status "Coverage comparison:"
        echo "  Previous: $PREVIOUS_COVERAGE%"
        echo "  Current:  $OVERALL_COVERAGE%"
        
        if (( $(echo "$COVERAGE_DIFF > 0" | bc -l) )); then
            print_success "Coverage improved by $COVERAGE_DIFF%"
        elif (( $(echo "$COVERAGE_DIFF < 0" | bc -l) )); then
            print_warning "Coverage decreased by ${COVERAGE_DIFF#-}%"
        else
            print_status "Coverage unchanged"
        fi
    fi
fi

# Final result
echo ""
echo -e "${THRESHOLD_COLOR}======================================${NC}"
echo -e "${THRESHOLD_COLOR}Coverage: $OVERALL_COVERAGE% | Threshold: $COVERAGE_THRESHOLD% | Status: $THRESHOLD_STATUS${NC}"
echo -e "${THRESHOLD_COLOR}======================================${NC}"

# Exit with appropriate code
if [ "$THRESHOLD_STATUS" = "FAIL" ]; then
    exit 1
else
    exit 0
fi