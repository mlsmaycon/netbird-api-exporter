#!/bin/bash

# Script to help update the CHANGELOG.md file
# Usage: ./scripts/update-changelog.sh [type] [description]
# Types: breaking, feature, bugfix, security, deprecated, removed

set -e

CHANGELOG_FILE="CHANGELOG.md"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CHANGELOG_PATH="$PROJECT_ROOT/$CHANGELOG_FILE"

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

# Create a temporary file for processing
TEMP_FILE=$(mktemp)

# Process the changelog
{
    # Read until we find the [Unreleased] section
    while IFS= read -r line; do
        echo "$line"
        if [[ "$line" == "## [Unreleased]" ]]; then
            break
        fi
    done

    # Skip any existing content until we find the next section or end of unreleased
    FOUND_SECTION=false
    ADDED_ENTRY=false
    
    while IFS= read -r line; do
        # If we hit the next version section, we need to add our entry
        if [[ "$line" == "## ["*"]"* ]] && [[ "$line" != "## [Unreleased]" ]]; then
            if [ "$ADDED_ENTRY" = false ]; then
                # Add our section and entry before the next version
                echo ""
                echo "### $SECTION"
                echo "- $DESCRIPTION"
                echo ""
            fi
            echo "$line"
            break
        fi
        
        # Check if this line is our target section
        if [[ "$line" == "### $SECTION" ]]; then
            echo "$line"
            FOUND_SECTION=true
            # Read and copy existing entries in this section
            while IFS= read -r subline; do
                if [[ "$subline" == "### "* ]] || [[ "$subline" == "## ["*"]"* ]]; then
                    # Add our new entry before moving to the next section
                    echo "- $DESCRIPTION"
                    echo "$subline"
                    ADDED_ENTRY=true
                    break
                elif [[ "$subline" =~ ^[[:space:]]*$ ]] && [ "$FOUND_SECTION" = true ]; then
                    # Empty line after our section - add entry here
                    echo "- $DESCRIPTION"
                    echo "$subline"
                    ADDED_ENTRY=true
                    break
                else
                    echo "$subline"
                fi
            done
            break
        elif [[ "$line" == "### "* ]]; then
            # Different section found, we need to add our section before it
            if [ "$ADDED_ENTRY" = false ]; then
                echo ""
                echo "### $SECTION"
                echo "- $DESCRIPTION"
                echo ""
            fi
            echo "$line"
            ADDED_ENTRY=true
            break
        elif [[ "$line" =~ ^[[:space:]]*$ ]]; then
            # Empty line - could be end of unreleased section
            continue
        else
            echo "$line"
        fi
    done

    # Copy the rest of the file
    while IFS= read -r line; do
        echo "$line"
    done
} < "$CHANGELOG_PATH" > "$TEMP_FILE"

# Replace the original file
mv "$TEMP_FILE" "$CHANGELOG_PATH"

echo -e "${GREEN}Successfully added entry to CHANGELOG.md:${NC}"
echo -e "${YELLOW}Section:${NC} $SECTION"
echo -e "${YELLOW}Description:${NC} $DESCRIPTION"
echo ""
echo -e "${BLUE}Tip:${NC} Review the changes with: git diff CHANGELOG.md" 