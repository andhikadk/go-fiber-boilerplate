.PHONY: help build run test clean migrate seed docker-up docker-down docker-logs docker-reset docker-dev docker-dev-logs docker-dev-down docker-dev-reset install-deps swagger swagger-install

# Variables
APP_NAME=go-fiber-boilerplate
MAIN_PATH=cmd/api/main.go

# Detect OS for binary extension
ifeq ($(OS),Windows_NT)
	BIN_EXT=.exe
else
	BIN_EXT=
endif

BINARY_NAME=./bin/$(APP_NAME)$(BIN_EXT)
AIR_BIN=./bin/app$(BIN_EXT)

help: ## Display this help screen
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-20s %s\n", $$1, $$2}'

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

build: swagger ## Build the application (generates Swagger docs first)
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/api
	@echo "Build complete: $(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/api

dev: ## Run in development mode with hot reload (requires air)
	@echo "Running in development mode..."
	@air --build.bin "$(AIR_BIN)" --build.cmd "go build -o $(AIR_BIN) ./cmd/api" || echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"

test: ## Run unit tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

swagger-install: ## Install swag CLI tool for generating Swagger documentation
	@echo "Installing swag CLI..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "swag installed successfully"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go
	@echo "Swagger documentation generated in ./docs"

swagger-fmt: ## Format Swagger comments
	@echo "Formatting Swagger comments..."
	@swag fmt

migrate: ## Run database migrations (AutoMigrate for dev)
	@echo "Running migrations..."
	@go run ./cmd/api -migrate=auto

migrate-sql: ## Run SQL migrations from files
	@echo "Running SQL migrations..."
	@go run ./cmd/api -migrate=sql

migrate-status: ## Show migration status
	@echo "Migration status..."
	@go run ./cmd/api -status

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@go run ./cmd/api -seed

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-up: ## Start Docker containers (docker compose)
	@echo "Starting Docker containers..."
	@docker compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker compose down

docker-logs: ## View Docker logs
	@docker compose logs -f

docker-reset: ## Reset Docker containers and volumes (removes all data)
	@echo "Resetting Docker containers and volumes..."
	@docker compose down -v
	@echo "Containers and volumes removed. Restart with: make docker-up"

docker-dev: ## Start Docker containers with hot reload (development mode)
	@echo "Starting Docker containers with hot reload..."
	@docker compose -f docker-compose.dev.yml up -d

docker-dev-logs: ## View Docker development logs
	@docker compose -f docker-compose.dev.yml logs -f

docker-dev-down: ## Stop Docker development containers
	@echo "Stopping Docker development containers..."
	@docker compose -f docker-compose.dev.yml down

docker-dev-reset: ## Reset Docker development containers and volumes
	@echo "Resetting Docker development containers and volumes..."
	@docker compose -f docker-compose.dev.yml down -v
	@echo "Development containers and volumes removed. Restart with: make docker-dev"

all: clean install-deps build test ## Clean, install, build and test

.DEFAULT_GOAL := help
