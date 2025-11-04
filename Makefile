.PHONY: help build run test clean docker-build docker-up docker-down docker-logs template

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@go build -o server ./cmd/server
	@echo "Build complete! Binary: ./server"

run: ## Run the application locally
	@echo "Running application..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f server
	@rm -f coverage.out coverage.html
	@rm -f generated/*.xlsx
	@echo "Clean complete!"

template: ## Generate Excel template
	@echo "Generating Excel template..."
	@go run scripts/simple_template.go
	@echo "Template generated!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build
	@echo "Docker build complete!"

docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	@docker-compose up -d
	@echo "Services started!"
	@echo "Acts Service: http://localhost:8080"
	@echo "Mongo Express: http://localhost:8081"

docker-down: ## Stop Docker Compose services
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped!"

docker-logs: ## Show Docker Compose logs
	@docker-compose logs -f

docker-restart: docker-down docker-up ## Restart Docker Compose services

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

all: clean install-deps template build ## Clean, install deps, generate template, and build

.DEFAULT_GOAL := help

