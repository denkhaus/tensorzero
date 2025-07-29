# TensorZero Go Client Makefile

.PHONY: help test test-unit test-openai docker-up docker-down docker-logs docker-status clean build deps

# Default target
help: ## Show this help message
	@echo "TensorZero Go Client - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Environment variables
OPENROUTER_API_KEY ?= $(shell echo $$OPENROUTER_API_KEY)
DOCKER_COMPOSE_FILE = docker/docker-compose.yaml
TENSORZERO_URL = http://localhost:3000

# Dependency management
deps: ## Install Go dependencies
	go mod download
	go mod tidy

# Build
build: deps ## Build the Go package
	go build ./...

# Testing targets
test: test-unit ## Run all tests (unit tests only by default)

test-unit: ## Run unit tests (no external dependencies)
	@echo "Running unit tests..."
	go test -tags=unit ./tests

test-openai: docker-up ## Run OpenAI compatibility tests against TensorZero instance
	@echo "Waiting for TensorZero to be ready..."
	@timeout 60 bash -c 'until curl -s $(TENSORZERO_URL)/health > /dev/null 2>&1; do echo "Waiting for TensorZero..."; sleep 2; done' || (echo "TensorZero failed to start" && exit 1)
	@echo "Running OpenAI compatibility tests..."
	go test -tags=integration ./tests/openai_test.go -v
	@echo "OpenAI compatibility tests completed"

test-all: docker-up ## Run all tests including OpenAI compatibility tests
	@echo "Running all tests..."
	$(MAKE) test-unit
	go test -tags=integration ./tests

# Docker management
docker-up: ## Start TensorZero services with Docker Compose
	@echo "Starting TensorZero services..."
	@if [ ! -f docker/.env ]; then \
		echo "Error: docker/.env file not found"; \
		echo "Please create docker/.env with OPENROUTER_API_KEY=your_openrouter_api_key"; \
		exit 1; \
	fi
	@if ! grep -q "OPENROUTER_API_KEY=" docker/.env; then \
		echo "Error: OPENROUTER_API_KEY not found in docker/.env"; \
		echo "Please add OPENROUTER_API_KEY=your_openrouter_api_key to docker/.env"; \
		exit 1; \
	fi
	cd docker && docker-compose -f docker-compose.yaml up -d
	@echo "Waiting for services to be healthy..."
	@timeout 120 bash -c 'until docker-compose -f $(DOCKER_COMPOSE_FILE) ps | grep -q "healthy"; do echo "Waiting for services..."; sleep 3; done' || (echo "Services failed to become healthy" && $(MAKE) docker-logs && exit 1)
	@echo "TensorZero services are ready!"
	@echo "Gateway: $(TENSORZERO_URL)"
	@echo "UI: http://localhost:4000"

docker-down: ## Stop TensorZero services
	@echo "Stopping TensorZero services..."
	cd docker && docker-compose -f docker-compose.yaml down

docker-restart: docker-down docker-up ## Restart TensorZero services

docker-logs: ## Show logs from TensorZero services
	cd docker && docker-compose -f docker-compose.yaml logs -f

docker-status: ## Show status of TensorZero services
	cd docker && docker-compose -f docker-compose.yaml ps

docker-clean: docker-down ## Clean up Docker resources
	@echo "Cleaning up Docker resources..."
	cd docker && docker-compose -f docker-compose.yaml down -v --remove-orphans
	docker system prune -f

# Development helpers
dev-setup: deps docker-up ## Set up development environment
	@echo "Development environment ready!"
	@echo "Run 'make test-all' to run all tests"

check-env: ## Check if required environment variables are set
	@echo "Checking environment..."
	@if [ ! -f docker/.env ]; then \
		echo "❌ docker/.env file not found"; \
		echo "   Please create docker/.env with OPENROUTER_API_KEY=your_openrouter_api_key"; \
		exit 1; \
	fi
	@if ! grep -q "OPENROUTER_API_KEY=" docker/.env; then \
		echo "❌ OPENROUTER_API_KEY not found in docker/.env"; \
		echo "   Please add OPENROUTER_API_KEY=your_openrouter_api_key to docker/.env"; \
		exit 1; \
	else \
		echo "✅ OPENROUTER_API_KEY is set in docker/.env"; \
	fi
	@if grep "OPENROUTER_API_KEY=" docker/.env | grep -q "dummy\|test"; then \
		echo "⚠️  OPENROUTER_API_KEY appears to be a dummy key"; \
		echo "   For full test functionality, use a real OpenRouter API key"; \
	fi
	@echo "✅ Environment check passed"

health-check: ## Check if TensorZero services are healthy
	@echo "Checking TensorZero health..."
	@curl -s $(TENSORZERO_URL)/health > /dev/null && echo "✅ TensorZero Gateway is healthy" || echo "❌ TensorZero Gateway is not responding"
	@curl -s http://localhost:4000 > /dev/null && echo "✅ TensorZero UI is healthy" || echo "❌ TensorZero UI is not responding"

# Cleanup
clean: docker-clean ## Clean up everything
	go clean -cache
	go clean -modcache
	@echo "Cleanup completed"

# Quick commands
up: docker-up ## Alias for docker-up
down: docker-down ## Alias for docker-down
logs: docker-logs ## Alias for docker-logs

# CI/CD helpers
ci-test: ## Run tests suitable for CI (with timeout)
	@echo "Running CI tests..."
	timeout 300 $(MAKE) test-all || (echo "Tests timed out or failed" && $(MAKE) docker-logs && exit 1)

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	go doc -all . > docs.txt
	@echo "Documentation generated in docs.txt"

# Version info
version: ## Show version information
	@echo "TensorZero Go Client"
	@echo "Go version: $(shell go version)"
	@echo "Docker version: $(shell docker --version 2>/dev/null || echo 'Docker not available')"
	@echo "Docker Compose version: $(shell docker-compose --version 2>/dev/null || echo 'Docker Compose not available')"

.DEFAULT_GOAL := help
