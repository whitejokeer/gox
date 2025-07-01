# GOX Makefile
.PHONY: help build test lint dev clean install deps fmt vet check

# Variables
BINARY_NAME=gox
BUILD_DIR=build
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
help: ## Show this help message
	@echo "GOX Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Dependencies
deps: ## Download and tidy dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Build
build: deps ## Build the GOX CLI binary
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gox

# Install
install: ## Install the GOX CLI globally
	@echo "📦 Installing $(BINARY_NAME) globally..."
	go install $(LDFLAGS) ./cmd/gox

# Development
dev: ## Start development mode with hot reload
	@echo "🚀 Starting development server..."
	go run $(LDFLAGS) ./cmd/gox dev

# Testing
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage report
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Linting and formatting
fmt: ## Format Go code
	@echo "✨ Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

lint: fmt vet ## Run all linting tools
	@echo "📝 Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Quality checks
check: lint test ## Run all quality checks (lint + test)

# Cleanup
clean: ## Clean build artifacts
	@echo "🧹 Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean

# Release (for CI)
release: clean check build ## Build release binary
	@echo "🚀 Release build complete!"

# Docker (optional)
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Help for development
dev-setup: ## Set up development environment
	@echo "🛠️  Setting up development environment..."
	@echo "Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if ! command -v air >/dev/null 2>&1; then \
		go install github.com/air-verse/air@latest; \
	fi
	@echo "✅ Development environment ready!"

# Watch for changes (requires air)
watch: ## Watch for changes and rebuild (requires air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not installed. Install it with: go install github.com/air-verse/air@latest"; \
		echo "Or use 'make dev-setup' to install all development tools"; \
	fi