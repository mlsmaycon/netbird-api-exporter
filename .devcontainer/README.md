# NetBird API Exporter Development Container

This directory contains the development container configuration for the NetBird API Exporter project.

## Features

- **Go 1.23**: Latest Go version matching the project requirements
- **Development Tools**: Pre-installed Go development tools including:
  - `golangci-lint`: For code linting
  - `goimports`: For import formatting
  - `air`: For hot reloading during development
  - `gopls`: Go language server
  - `pre-commit`: For automated code quality checks
- **VS Code Extensions**: Automatically installs relevant extensions for Go development
- **Docker-in-Docker**: Ability to build and run Docker containers from within the devcontainer
- **Port Forwarding**: Automatic forwarding of port 8080 for the exporter
- **Prometheus Integration**: Optional Prometheus setup for testing metrics

## Getting Started

1. **Open in VS Code**: Open the project in VS Code and click "Reopen in Container" when prompted, or use the Command Palette (`Ctrl+Shift+P`) and select "Dev Containers: Reopen in Container"

   > **Important**: If you're upgrading from a previous version, you **must** rebuild the container to get Go 1.23 and all development tools. Use "Dev Containers: Rebuild Container" from the Command Palette.

2. **Environment Variables**: The devcontainer sets default environment variables. You'll need to provide your NetBird API token:

   ```bash
   export NETBIRD_API_TOKEN="your-token-here"
   ```

3. **Pre-commit Hooks**: Pre-commit hooks are automatically installed during container creation and will run formatting, linting, and tests on each commit. You can also run them manually:

   ```bash
   # Run all pre-commit hooks on all files
   pre-commit run --all-files

   # Or use make targets
   make setup-precommit  # Interactive setup
   make check           # Run all quality checks
   ```

4. **Run the Application**:

   ```bash
   # Development mode with hot reloading
   air

   # Or build and run normally
   go run main.go
   ```

5. **Access the Application**:
   - Application: http://localhost:8080
   - Metrics: http://localhost:8080/metrics
   - Health Check: http://localhost:8080/health
   - Prometheus (if using docker-compose): http://localhost:9090

## Development Workflow

### Hot Reloading

The container includes `air` for hot reloading. The `.air.toml` configuration file in the root directory will be used automatically.

### Pre-commit Hooks

Pre-commit hooks are automatically installed and configured. They run on every commit to ensure code quality:

```bash
# Run all hooks manually
pre-commit run --all-files

# Run specific hook
pre-commit run go-fmt

# Skip hooks for a commit (not recommended)
git commit --no-verify -m "message"
```

### Linting and Formatting

- Code is automatically formatted on save using `goimports`
- Linting is performed using `golangci-lint`
- Run manually: `make lint` or `golangci-lint run`
- Pre-commit hooks automatically run these checks before commits

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Building

```bash
# Build the application
go build -o netbird-api-exporter

# Or use the Makefile
make build
```

## Docker Compose Setup

The included `docker-compose.yml` provides:

- The development container
- Prometheus instance for testing metrics collection

To use with docker-compose:

```bash
# Update devcontainer.json to use docker-compose
# Change "image" to "dockerComposeFile": "docker-compose.yml"
# Add "service": "app"
```

## Environment Variables

The following environment variables are pre-configured in the devcontainer:

- `NETBIRD_API_URL`: https://api.netbird.io
- `LISTEN_ADDRESS`: :8080
- `METRICS_PATH`: /metrics
- `LOG_LEVEL`: debug

You can override these by setting them in your shell or updating the `containerEnv` section in `devcontainer.json`.

## Customization

### Adding Extensions

Edit the `extensions` array in `devcontainer.json` to add VS Code extensions.

### Modifying Settings

Update the `settings` section in `devcontainer.json` to customize VS Code behavior.

### Additional Tools

Add tools in the `postCreateCommand` or modify the Dockerfile for persistent installations.

## Troubleshooting

### Port Conflicts

If port 8080 is already in use, update the `forwardPorts` section in `devcontainer.json`.

### Permission Issues

The devcontainer runs as the `vscode` user. If you encounter permission issues, check file ownership and permissions.

### Docker Socket Access

Docker-in-Docker is enabled via socket mounting. Ensure Docker is running on your host system.

### Go Version Mismatch

If you see errors like "go.mod requires go >= 1.23 (running go 1.21.13)", rebuild the devcontainer:

1. Use `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
2. Select "Dev Containers: Rebuild Container"
3. Wait for the container to rebuild with Go 1.23

### Pre-commit Not Found

If pre-commit is not found after container rebuild:

1. Check if it's installed: `pipx list | grep pre-commit`
2. Verify PATH includes user bin: `echo $PATH | grep local/bin`
3. Reinstall if needed: `pipx install pre-commit`
4. Install hooks: `pre-commit install`
