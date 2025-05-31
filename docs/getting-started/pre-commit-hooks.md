# Pre-commit Hooks

Pre-commit hooks help maintain code quality by automatically running checks before each commit. The NetBird API Exporter project provides two options for setting up pre-commit hooks.

## What the Hooks Do

The pre-commit hooks automatically run:

- **`go fmt`** - Formats your Go code according to Go standards
- **`golangci-lint`** - Runs comprehensive linting checks
- **`go vet`** - Performs static analysis to find potential issues
- **`go test`** - Runs all unit tests to ensure code functionality

## Quick Setup

### Option 1: Interactive Setup (Recommended)

```bash
make setup-precommit
```

This will guide you through choosing between the simple git hook or the more advanced pre-commit framework.

### Option 2: Simple Git Hook

```bash
make install-hooks
```

This installs a basic bash script that runs fmt, lint, and test before each commit.

### Option 3: Pre-commit Framework

If you prefer the pre-commit framework with more features:

```bash
# Install pre-commit (choose one)
pip install pre-commit
# or
pipx install pre-commit
# or
brew install pre-commit  # macOS
# or
apt install pre-commit   # Ubuntu/Debian

# Install hooks
pre-commit install
```

## Pre-commit Framework vs Simple Git Hook

### Simple Git Hook

**Pros:**

- No external dependencies
- Simple to understand and modify
- Works out of the box

**Cons:**

- Runs checks on all files (slower)
- Less configurable
- Basic functionality only

### Pre-commit Framework

**Pros:**

- Only runs on changed files (faster)
- Highly configurable
- Includes additional checks (YAML, JSON, file sizes, etc.)
- Better integration with CI/CD
- Can run hooks on specific file types

**Cons:**

- Requires installing pre-commit tool
- More complex configuration

## Configuration Files

### Simple Git Hook

The hook script is located at `scripts/pre-commit` and can be customized by editing the bash script.

### Pre-commit Framework

Configuration is in `.pre-commit-config.yaml`. You can modify it to:

- Add or remove hooks
- Change which files are checked
- Configure hook behavior

Example configuration:

```yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: "go fmt"
        entry: make fmt
        language: system
        files: '\.go$'
        pass_filenames: false
```

## Usage

### Normal Workflow

Once installed, the hooks run automatically:

```bash
# Make changes to Go files
vim pkg/exporters/peer.go

# Stage changes
git add .

# Commit (hooks run automatically)
git commit -m "Add new feature"
```

### If Hooks Fail

When hooks fail, you'll see output like:

```
🔧 Running lint...
❌ Linting failed
Please fix the linting issues before committing.
```

Fix the issues and commit again:

```bash
# Fix the issues
make fmt
make lint

# Stage fixes
git add .

# Commit again
git commit -m "Add new feature"
```

### Bypassing Hooks (Not Recommended)

In emergency situations, you can bypass hooks:

```bash
git commit --no-verify -m "Emergency fix"
```

## Manual Testing

### Test All Hooks

```bash
# Pre-commit framework
pre-commit run --all-files

# Or run individual checks
make fmt
make lint
make test
```

### Test on Specific Files

```bash
# Pre-commit framework
pre-commit run --files pkg/exporters/peer.go

# Or use make targets
make fmt
```

## Troubleshooting

### Hook Not Running

1. Check if hook is installed:

   ```bash
   ls -la .git/hooks/pre-commit
   ```

2. Verify hook is executable:
   ```bash
   chmod +x .git/hooks/pre-commit
   ```

### Pre-commit Framework Issues

1. Check pre-commit installation:

   ```bash
   pre-commit --version
   ```

2. Update hooks:

   ```bash
   pre-commit autoupdate
   ```

3. Clear cache:
   ```bash
   pre-commit clean
   ```

### Lint Failures

1. Run lint manually to see full output:

   ```bash
   make lint
   ```

2. Auto-fix common issues:

   ```bash
   make fmt
   ```

3. Check golangci-lint configuration:
   ```bash
   golangci-lint run --help
   ```

## Development Workflow

### Before First Commit

```bash
# Set up development environment
make deps
make setup-precommit

# Verify everything works
make check
```

### Daily Development

```bash
# Regular development workflow
git add .
git commit -m "Your changes"  # Hooks run automatically

# If hooks fail, fix and retry
make fmt lint
git add .
git commit -m "Your changes"
```

## CI Integration

The same checks that run in pre-commit hooks should also run in CI:

```yaml
# Example GitHub Actions
- name: Run checks
  run: |
    make fmt
    make lint
    make test
```

This ensures consistency between local development and CI environments.

## Advanced Configuration

### Custom Hook Scripts

You can modify `scripts/pre-commit` to add custom checks:

```bash
# Add security scanning
echo "🔒 Running security scan..."
if ! make security; then
    echo "❌ Security scan failed"
    exit 1
fi
```

### Skip Specific Hooks

```bash
# Skip only tests
SKIP=go-test git commit -m "WIP: work in progress"

# Skip multiple hooks
SKIP=go-test,golangci-lint git commit -m "Quick fix"
```

### Hook Performance

For large repositories, consider:

- Using `files` patterns to limit scope
- Setting `pass_filenames: true` where possible
- Using `stages` to run expensive checks only on push

## Best Practices

1. **Always run hooks**: Don't bypass unless absolutely necessary
2. **Fix issues promptly**: Don't let lint errors accumulate
3. **Keep hooks fast**: Slow hooks discourage usage
4. **Test hook changes**: Verify hooks work before sharing
5. **Document custom hooks**: Make it easy for others to understand

## Getting Help

If you encounter issues:

1. Check this documentation
2. Run `make help` for available commands
3. Check hook output for specific error messages
4. Review the hook scripts in `scripts/` directory
