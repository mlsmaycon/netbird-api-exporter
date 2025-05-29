.PHONY: build run clean test docker-build docker-run docker-compose-up docker-compose-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=netbird-api-exporter
DOCKER_IMAGE=netbird-api-exporter

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run with Docker
docker-run: docker-build
	docker run -d \
		--name $(DOCKER_IMAGE) \
		-p 8080:8080 \
		-e NETBIRD_API_TOKEN=${NETBIRD_API_TOKEN} \
		$(DOCKER_IMAGE)

# Stop and remove Docker container
docker-stop:
	docker stop $(DOCKER_IMAGE) || true
	docker rm $(DOCKER_IMAGE) || true

# Start with Docker Compose
docker-compose-up:
	docker-compose up -d

# Stop Docker Compose
docker-compose-down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f netbird-exporter

# Development mode with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Running golangci-lint..."
	GO111MODULE=on golangci-lint run
	@echo "Running go vet..."
	GO111MODULE=on go vet ./...
	@echo "Running go fmt check..."
	@test -z "$$(GO111MODULE=on gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt' to fix." && exit 1)
	@echo "All linting checks passed!"

# Security scan (requires gosec)
security:
	gosec ./...

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-windows-amd64.exe

# Run all checks (tests, linting, formatting)
check: test lint
	@echo "All checks passed!"

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the binary"
	@echo "  run             - Build and run the application"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  deps            - Download and tidy dependencies"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run with Docker"
	@echo "  docker-stop     - Stop Docker container"
	@echo "  docker-compose-up   - Start with Docker Compose"
	@echo "  docker-compose-down - Stop Docker Compose"
	@echo "  logs            - View container logs"
	@echo "  dev             - Development mode with live reload"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code (golangci-lint, go vet, format check)"
	@echo "  check           - Run all checks (tests + linting)"
	@echo "  security        - Run security scan"
	@echo "  build-all       - Build for multiple platforms"
	@echo "  help            - Show this help" 