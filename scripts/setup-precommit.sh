#!/bin/bash

# Setup script for pre-commit hooks in netbird-api-exporter
# Provides two options: simple git hook or pre-commit framework

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔧 NetBird API Exporter Pre-commit Setup"
echo "========================================"
echo
echo "This script will help you set up pre-commit hooks that run:"
echo "  • go fmt (code formatting)"
echo "  • go lint (code linting)"
echo "  • go test (unit tests)"
echo

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    echo "Please run this script from within the git repository."
    exit 1
fi

echo "Choose your preferred pre-commit setup:"
echo
echo "1) Simple Git Hook (bash script)"
echo "   - Basic setup using a bash script"
echo "   - No external dependencies"
echo "   - Runs fmt, lint, and test on all Go files"
echo
echo "2) Pre-commit Framework (recommended)"
echo "   - More sophisticated hook management"
echo "   - Requires installing pre-commit (pip install pre-commit)"
echo "   - Additional checks for YAML, JSON, file sizes, etc."
echo "   - Better performance (only runs on changed files when possible)"
echo
echo "3) Both (install both options)"
echo

read -p "Enter your choice (1/2/3): " choice

install_simple_hook() {
    echo
    echo "📋 Installing simple git hook..."
    
    # Copy the pre-commit script
    cp "$SCRIPT_DIR/pre-commit" "$PROJECT_ROOT/.git/hooks/pre-commit"
    chmod +x "$PROJECT_ROOT/.git/hooks/pre-commit"
    
    echo "✅ Simple git hook installed successfully!"
    echo "   Hook location: .git/hooks/pre-commit"
}

install_precommit_framework() {
    echo
    echo "📋 Installing pre-commit framework..."
    
    # Check if pre-commit is installed
    if ! command -v pre-commit &> /dev/null; then
        echo "⚠️  pre-commit is not installed."
        echo
        echo "To install pre-commit, run one of these commands:"
        echo "  • pip install pre-commit"
        echo "  • pipx install pre-commit"
        echo "  • brew install pre-commit (on macOS)"
        echo "  • apt install pre-commit (on Ubuntu/Debian)"
        echo
        read -p "Would you like to try installing with pip? (y/N): " install_pip
        
        if [[ "$install_pip" =~ ^[Yy]$ ]]; then
            echo "Installing pre-commit with pip..."
            pip install pre-commit || {
                echo "❌ Failed to install pre-commit with pip"
                echo "Please install pre-commit manually and run this script again."
                return 1
            }
        else
            echo "Please install pre-commit manually and run this script again."
            return 1
        fi
    fi
    
    # Install the pre-commit hooks
    cd "$PROJECT_ROOT"
    pre-commit install
    
    echo "✅ Pre-commit framework installed successfully!"
    echo "   Configuration: .pre-commit-config.yaml"
    echo "   You can run 'pre-commit run --all-files' to test all hooks"
}

case $choice in
    1)
        install_simple_hook
        ;;
    2)
        install_precommit_framework
        ;;
    3)
        install_simple_hook
        install_precommit_framework
        echo
        echo "⚠️  Note: Both hooks are installed. The pre-commit framework will take precedence."
        ;;
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo
echo "🎉 Setup complete!"
echo
echo "What happens now:"
echo "• When you commit, the hooks will automatically run"
echo "• fmt, lint, and test will be executed"
echo "• Commit will be blocked if any checks fail"
echo
echo "To bypass hooks (not recommended): git commit --no-verify"
echo
echo "Test your setup:"
echo "• Make a small change to a Go file"
echo "• Run: git add . && git commit -m 'test commit'"
echo "• The hooks should run automatically" 