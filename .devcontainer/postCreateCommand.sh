#!/bin/bash

# Post-create command for NetBird API Exporter devcontainer
# This script sets up the development environment

set -e

echo "🚀 Setting up NetBird API Exporter development environment..."

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download

# Install Go development tools
echo "🔧 Installing Go development tools..."

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
echo "📍 Go version: $GO_VERSION"

# Install golangci-lint with version compatibility
echo "📦 Installing golangci-lint..."
if ! go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; then
    echo "⚠️  Latest golangci-lint failed, trying compatible version..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
fi

echo "📦 Installing goimports..."
go install golang.org/x/tools/cmd/goimports@latest

echo "📦 Installing air (live reload)..."
go install github.com/air-verse/air@latest

# Install pre-commit
echo "🪝 Installing pre-commit..."
if ! command -v pip &> /dev/null && ! command -v pip3 &> /dev/null; then
    echo "📦 Installing Python pip..."
    sudo apt update -qq
    sudo apt install -y python3-pip pipx
fi

# Install pre-commit using pipx (preferred) or pip
if command -v pipx &> /dev/null; then
    echo "📦 Installing pre-commit with pipx..."
    pipx install pre-commit
elif command -v pip3 &> /dev/null; then
    echo "📦 Installing pre-commit with pip3..."
    pip3 install --user pre-commit
    # Add pip user bin to PATH if not already there
    export PATH="$HOME/.local/bin:$PATH"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
else
    echo "❌ Could not install pre-commit: no pip or pipx found"
fi

# Install pre-commit hooks automatically
echo "🎯 Setting up pre-commit hooks..."
if [ -f ".pre-commit-config.yaml" ]; then
    pre-commit migrate-config 2>/dev/null || true  # Fix deprecated config if needed
    pre-commit install
    echo "✅ Pre-commit hooks installed successfully!"
else
    echo "⚠️  .pre-commit-config.yaml not found, skipping pre-commit hook installation"
fi

# Make scripts executable
echo "🔑 Making scripts executable..."
chmod +x scripts/*.sh

# Run initial checks to ensure everything is working
echo "🧪 Running initial checks..."
if make fmt lint test; then
    echo "✅ All checks passed!"
else
    echo "⚠️  Some checks failed. You may need to fix issues before committing."
fi

echo "🎉 Development environment setup complete!"
echo ""

# Run verification
echo "🔍 Running environment verification..."
if ./scripts/verify-setup.sh; then
    echo ""
    echo "✅ All systems ready for development!"
else
    echo ""
    echo "⚠️  Some issues detected. Check the output above."
fi

echo ""
echo "Available commands:"
echo "  make help            - Show all available make targets"
echo "  make setup-precommit - Interactive pre-commit setup"
echo "  make dev             - Start development mode with live reload"
echo "  make check           - Run all quality checks"
echo "  make verify          - Verify environment setup"
echo ""
echo "Pre-commit hooks are ready! They will run automatically on commits."
echo "To bypass hooks (not recommended): git commit --no-verify" 