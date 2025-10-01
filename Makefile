# Makefile for PolicyReporter-FSM project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Project parameters
BINARY_NAME=policyreporter-fsm
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./main.go
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-v $(LDFLAGS)

# Test flags
TEST_FLAGS=-v -race -coverprofile=$(COVERAGE_FILE)
BENCH_FLAGS=-v -bench=. -benchmem

.PHONY: all build clean test coverage lint fmt vet deps help

# Default target
all: clean deps fmt vet lint test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Build for Linux
build-linux:
	@echo "Building $(BINARY_UNIX)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_UNIX) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	$(GOTEST) -v ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) $(BENCH_FLAGS) ./...

# Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Show coverage in terminal
coverage-text: test
	@echo "Coverage summary:"
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Verify dependencies
deps-verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

# Run the application
run:
	@echo "Running application..."
	$(GOCMD) run $(MAIN_PATH)

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	$(GOCMD) list -json -m all | nancy sleuth

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./...

# Run all checks (CI pipeline)
ci: clean deps fmt vet lint test coverage-text

# Development workflow
dev: fmt vet test

# Release build
release: clean deps fmt vet lint test build-linux
	@echo "Release build completed"

# Docker build (if Dockerfile exists)
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Docker run (if Dockerfile exists)
docker-run:
	@echo "Running Docker container..."
	docker run --rm $(BINARY_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Run full build pipeline (clean, deps, fmt, vet, lint, test, build)"
	@echo "  build        - Build the binary"
	@echo "  build-linux  - Build binary for Linux"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests with coverage"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  bench        - Run benchmarks"
	@echo "  coverage     - Generate HTML coverage report"
	@echo "  coverage-text- Show coverage summary in terminal"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code with gofmt"
	@echo "  vet          - Run go vet"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo "  deps-verify  - Verify dependencies"
	@echo "  run          - Run the application"
	@echo "  install-tools- Install development tools"
	@echo "  security     - Check for security vulnerabilities"
	@echo "  docs         - Generate documentation"
	@echo "  ci           - Run CI pipeline"
	@echo "  dev          - Development workflow (fmt, vet, test)"
	@echo "  release      - Release build"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help message"
