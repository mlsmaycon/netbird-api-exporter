#!/bin/bash

# Script to help update the CHANGELOG.md file and Helm chart artifacthub.io/changes
# Usage: ./scripts/update-changelog.sh [type] [description]
# Types: breaking, feature, bugfix, security, deprecated, removed

set -e

CHANGELOG_FILE="CHANGELOG.md"
CHART_FILE="charts/netbird-api-exporter/Chart.yaml"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CHANGELOG_PATH="$PROJECT_ROOT/$CHANGELOG_FILE"
CHART_PATH="$PROJECT_ROOT/$CHART_FILE"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_usage() {
    echo -e "${BLUE}Usage: $0 [type] [description]${NC}"
    echo ""
    echo -e "${YELLOW}Types:${NC}"
    echo "  breaking   - Breaking changes that require user action"
    echo "  feature    - New functionality and enhancements"
    echo "  bugfix     - Bug fixes and corrections"
    echo "  security   - Security-related changes"
    echo "  deprecated - Features that will be removed in future versions"
    echo "  removed    - Features that have been removed"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 feature \"Add support for custom metrics endpoint\""
    echo "  $0 bugfix \"Fix memory leak in DNS exporter (#123)\""
    echo "  $0 breaking \"Remove deprecated --old-flag parameter\""
    echo ""
    echo -e "${YELLOW}Note:${NC} This script will also update the Helm chart's artifacthub.io/changes annotation"
    echo "and include a summary of uncommitted files in the changelog entry."
}

# Function to get uncommitted files summary
get_uncommitted_summary() {
    local git_status
    git_status=$(git status --porcelain 2>/dev/null || echo "")
    
    if [ -z "$git_status" ]; then
        return 0
    fi
    
    echo ""
    echo "Files modified in this change:"
    
    # Parse git status output
    while IFS= read -r line; do
        if [ -z "$line" ]; then
            continue
        fi
        
        local status="${line:0:2}"
        local file="${line:3}"
        
        case "$status" in
            " M"|"M "|"MM")
                echo "- Modified: $file"
                ;;
            " A"|"A "|"AM")
                echo "- Added: $file"
                ;;
            " D"|"D "|"AD")
                echo "- Deleted: $file"
                ;;
            "??")
                echo "- New: $file"
                ;;
            "R ")
                echo "- Renamed: $file"
                ;;
            *)
                echo "- Changed: $file"
                ;;
        esac
    done <<< "$git_status"
}

# Function to update Helm chart artifacthub.io/changes
update_chart_changes() {
    local description="$1"
    
    if [ ! -f "$CHART_PATH" ]; then
        echo -e "${YELLOW}Warning: Chart.yaml not found at $CHART_PATH, skipping Helm chart update${NC}"
        return 0
    fi
    
    # Create a temporary file for processing Chart.yaml
    local temp_chart=$(mktemp)
    
    # Process the Chart.yaml file
    {
        local in_changes=false
        local changes_updated=false
        
        while IFS= read -r line; do
            if [[ "$line" =~ ^[[:space:]]*artifacthub\.io/changes:[[:space:]]*\|[[:space:]]*$ ]]; then
                echo "$line"
                echo "    - $description"
                in_changes=true
                changes_updated=true
                # Skip existing changes lines
                while IFS= read -r next_line; do
                    if [[ "$next_line" =~ ^[[:space:]]*-[[:space:]] ]] && [ "$in_changes" = true ]; then
                        continue  # Skip existing change entries
                    else
                        echo "$next_line"
                        in_changes=false
                        break
                    fi
                done
            else
                echo "$line"
            fi
        done
        
        # If we didn't find the changes section, this should not happen with our Chart.yaml
        if [ "$changes_updated" = false ]; then
            echo -e "${YELLOW}Warning: Could not find artifacthub.io/changes section in Chart.yaml${NC}" >&2
        fi
    } < "$CHART_PATH" > "$temp_chart"
    
    # Replace the original file
    mv "$temp_chart" "$CHART_PATH"
    
    if [ "$changes_updated" = true ]; then
        echo -e "${GREEN}Successfully updated Helm chart artifacthub.io/changes${NC}"
    fi
}

if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    print_usage
    exit 0
fi

if [ $# -ne 2 ]; then
    echo -e "${RED}Error: Both type and description are required${NC}"
    print_usage
    exit 1
fi

TYPE="$1"
DESCRIPTION="$2"

# Validate type
case "$TYPE" in
    breaking|feature|bugfix|security|deprecated|removed)
        ;;
    *)
        echo -e "${RED}Error: Invalid type '$TYPE'${NC}"
        print_usage
        exit 1
        ;;
esac

# Check if changelog file exists
if [ ! -f "$CHANGELOG_PATH" ]; then
    echo -e "${RED}Error: CHANGELOG.md not found at $CHANGELOG_PATH${NC}"
    exit 1
fi

# Map type to section header
case "$TYPE" in
    breaking)
        SECTION="BREAKING CHANGES"
        ;;
    feature)
        SECTION="Features"
        ;;
    bugfix)
        SECTION="Bugfix"
        ;;
    security)
        SECTION="Security"
        ;;
    deprecated)
        SECTION="Deprecated"
        ;;
    removed)
        SECTION="Removed"
        ;;
esac

# Get uncommitted files summary
UNCOMMITTED_SUMMARY=$(get_uncommitted_summary)

# Prepare the full description with uncommitted summary if available
FULL_DESCRIPTION="$DESCRIPTION"
if [ -n "$UNCOMMITTED_SUMMARY" ]; then
    FULL_DESCRIPTION="$DESCRIPTION$UNCOMMITTED_SUMMARY"
fi

# Create a temporary file for processing
TEMP_FILE=$(mktemp)

# Process the changelog with a simpler, more robust approach
{
    # Read and copy everything until we find the [Unreleased] section
    while IFS= read -r line; do
        echo "$line"
        if [[ "$line" == "## [Unreleased]" ]]; then
            break
        fi
    done
    
    # Now we're after the [Unreleased] line
    # Look for existing content under Unreleased section
    found_section=false
    added_entry=false
    in_unreleased=true
    buffer=""
    
    while IFS= read -r line && [ "$in_unreleased" = true ]; do
        # Check if we've reached the next version section
        if [[ "$line" =~ ^##\ \[.*\]\ -\ [0-9] ]]; then
            # We've reached the first versioned release
            # If we haven't added our entry yet, add it now
            if [ "$added_entry" = false ]; then
                if [ "$found_section" = false ]; then
                    echo ""
                    echo "### $SECTION"
                    echo "- $FULL_DESCRIPTION"
                else
                    # We found the section but haven't added our entry yet
                    # Add it to the buffer and then flush
                    echo "- $FULL_DESCRIPTION"
                fi
                echo ""
            fi
            echo "$line"
            in_unreleased=false
            break
        fi
        
        # Check if this is our target section
        if [[ "$line" == "### $SECTION" ]]; then
            echo "$line"
            found_section=true
            # Read the next line to see if we should add our entry immediately
            continue
        fi
        
        # Check if this is a different section
        if [[ "$line" =~ ^###\  ]] && [[ "$line" != "### $SECTION" ]]; then
            # Different section found
            if [ "$found_section" = false ]; then
                # Add our section before this one
                echo ""
                echo "### $SECTION"
                echo "- $FULL_DESCRIPTION"
                echo ""
                added_entry=true
            fi
            echo "$line"
            continue
        fi
        
        # If we're in our target section, add our entry before any existing entries
        if [ "$found_section" = true ] && [ "$added_entry" = false ]; then
            echo "- $FULL_DESCRIPTION"
            added_entry=true
        fi
        
        echo "$line"
    done
    
    # If we never found any sections in unreleased, add our section now
    if [ "$in_unreleased" = true ] && [ "$added_entry" = false ]; then
        echo ""
        echo "### $SECTION"
        echo "- $FULL_DESCRIPTION"
        echo ""
    fi
    
    # Copy the rest of the file
    while IFS= read -r line; do
        echo "$line"
    done
} < "$CHANGELOG_PATH" > "$TEMP_FILE"

# Replace the original file
mv "$TEMP_FILE" "$CHANGELOG_PATH"

# Update Helm chart changes
update_chart_changes "$DESCRIPTION"

echo -e "${GREEN}Successfully added entry to CHANGELOG.md:${NC}"
echo -e "${YELLOW}Section:${NC} $SECTION"
echo -e "${YELLOW}Description:${NC} $DESCRIPTION"
if [ -n "$UNCOMMITTED_SUMMARY" ]; then
    echo -e "${YELLOW}Uncommitted files summary added to changelog entry${NC}"
fi
echo ""
echo -e "${BLUE}Tip:${NC} Review the changes with:"
echo "  git diff CHANGELOG.md"
echo "  git diff charts/netbird-api-exporter/Chart.yaml" 