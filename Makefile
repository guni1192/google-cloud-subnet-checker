# Makefile for building, testing, and linting a Go project

# Variables
GO=go
LINTER=golangci-lint

# Default target
all: build

# Build the project
build:
	@echo "Building the project..."
	mkdir -p bin/
	$(GO) build -o bin/google-cloud-subnet-checker cmd/google-cloud-subnet-checker/main.go

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	$(LINTER) run

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/

# Phony targets
.PHONY: all build test lint clean