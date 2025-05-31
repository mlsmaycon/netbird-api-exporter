#!/bin/bash

# Verification script for NetBird API Exporter development environment
# This script checks that all required tools are properly installed

set -e

echo "🔍 NetBird API Exporter Environment Verification"
echo "==============================================="
echo

# Check Go version
echo "📋 Checking Go version..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    echo "✅ Go: $GO_VERSION"
    
    # Check if Go version meets minimum requirement
    if [[ $(echo "$GO_VERSION 1.23" | tr " " "\n" | sort -V | head -1) != "1.23" ]]; then
        echo "⚠️  Go version $GO_VERSION is older than required 1.23"
    fi
else
    echo "❌ Go: Not found"
    exit 1
fi

# Check golangci-lint
echo
echo "📋 Checking golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    LINT_VERSION=$(golangci-lint version --format=short 2>/dev/null || echo "unknown")
    echo "✅ golangci-lint: $LINT_VERSION"
else
    echo "❌ golangci-lint: Not found"
fi

# Check goimports
echo
echo "📋 Checking goimports..."
if command -v goimports &> /dev/null; then
    echo "✅ goimports: Available"
else
    echo "❌ goimports: Not found"
fi

# Check air
echo
echo "📋 Checking air (live reload)..."
if command -v air &> /dev/null; then
    echo "✅ air: Available"
else
    echo "❌ air: Not found"
fi

# Check pre-commit
echo
echo "📋 Checking pre-commit..."
# Ensure PATH includes common locations for pipx and pip user installs
export PATH="$HOME/.local/bin:/usr/local/bin:$PATH"

if command -v pre-commit &> /dev/null; then
    PRECOMMIT_VERSION=$(pre-commit --version | cut -d' ' -f2)
    echo "✅ pre-commit: $PRECOMMIT_VERSION"
    
    # Check if hooks are installed
    if [ -f ".git/hooks/pre-commit" ]; then
        echo "✅ pre-commit hooks: Installed"
    else
        echo "⚠️  pre-commit hooks: Not installed (run 'pre-commit install')"
    fi
else
    echo "❌ pre-commit: Not found"
    echo "   Check these locations:"
    ls -la ~/.local/bin/pre-commit 2>/dev/null && echo "   Found in ~/.local/bin/" || true
    which pipx &>/dev/null && pipx list | grep pre-commit || true
fi

# Check make targets
echo
echo "📋 Checking make targets..."
if command -v make &> /dev/null; then
    echo "✅ make: Available"
    echo "📁 Available targets:"
    make help | grep -E "(fmt|lint|test|check)" | head -4
else
    echo "❌ make: Not found"
fi

# Check project dependencies
echo
echo "📋 Checking project dependencies..."
if [ -f "go.mod" ]; then
    echo "✅ go.mod: Found"
    if go mod verify &> /dev/null; then
        echo "✅ Dependencies: Verified"
    else
        echo "⚠️  Dependencies: Issues found (run 'go mod tidy')"
    fi
else
    echo "❌ go.mod: Not found"
fi

# Test basic functionality
echo
echo "📋 Testing basic functionality..."

# Test go fmt
if go fmt ./... &> /dev/null; then
    echo "✅ go fmt: Working"
else
    echo "❌ go fmt: Issues found"
fi

# Test go vet (only if no linting issues)
if go vet ./... &> /dev/null; then
    echo "✅ go vet: Working"
else
    echo "⚠️  go vet: Issues found (check with 'make lint')"
fi

# Test go build
if go build -o /tmp/netbird-api-exporter . &> /dev/null; then
    echo "✅ go build: Working"
    rm -f /tmp/netbird-api-exporter
else
    echo "❌ go build: Failed"
fi

echo
echo "🎉 Environment verification complete!"
echo
echo "Next steps:"
echo "  • Set your NetBird API token: export NETBIRD_API_TOKEN=your_token_here"
echo "  • Run the application: make run"
echo "  • Start development mode: make dev"
echo "  • Run all checks: make check" 