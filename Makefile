.PHONY: help build test test-libs test-services run-api lint tidy sync

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build all services"
	@echo "  make test           - Run all tests"
	@echo "  make test-libs      - Run tests for all libraries"
	@echo "  make test-services  - Run tests for all services"
	@echo "  make run-api        - Run the api-gateway service"
	@echo "  make lint           - Run linter on all modules"
	@echo "  make tidy           - Run go mod tidy on all modules"
	@echo "  make sync           - Sync go.work dependencies"

# Build all services
build:
	@echo "Building all services..."
	@for svc in services/*/; do \
		echo "Building $$svc..."; \
		cd $$svc && go build ./cmd/... && cd - > /dev/null; \
	done
	@echo "Build complete!"

# Run all tests
test: test-libs test-services

# Run tests for libraries
test-libs:
	@echo "Testing libraries..."
	@for lib in libs/*/; do \
		echo "Testing $$lib..."; \
		cd $$lib && go test ./... -v && cd - > /dev/null; \
	done

# Run tests for services
test-services:
	@echo "Testing services..."
	@for svc in services/*/; do \
		echo "Testing $$svc..."; \
		cd $$svc && go test ./... -v && cd - > /dev/null; \
	done

# Run api-gateway service
run-api:
	@echo "Starting api-gateway service..."
	cd services/api-gateway && go run ./cmd/

# Run linter on all modules
lint:
	@echo "Running linter..."
	@for d in libs/* services/*; do \
		echo "Linting $$d..."; \
		cd $$d && golangci-lint run ./... && cd - > /dev/null; \
	done

# Run go mod tidy on all modules
tidy:
	@echo "Running go mod tidy on all modules..."
	@for d in libs/* services/*; do \
		echo "Tidying $$d..."; \
		cd $$d && go mod tidy && cd - > /dev/null; \
	done
	@echo "Tidy complete!"

# Sync go.work dependencies
sync:
	@echo "Syncing go.work dependencies..."
	go work sync
	@echo "Sync complete!"

# Build specific service
build-%:
	@echo "Building services/$*..."
	cd services/$* && go build ./cmd/...

# Test specific service
test-%:
	@echo "Testing services/$*..."
	cd services/$* && go test ./... -v

# Add a new library
new-lib:
	@read -p "Enter library name: " name; \
	mkdir -p libs/$$name; \
	cd libs/$$name && go mod init monorepo/libs/$$name; \
	echo "Created libs/$$name"; \
	echo "Don't forget to add it to go.work!"

# Add a new service
new-service:
	@read -p "Enter service name: " name; \
	mkdir -p services/$$name/cmd; \
	mkdir -p services/$$name/internal; \
	cd services/$$name && go mod init monorepo/services/$$name; \
	echo "Created services/$$name"; \
	echo "Don't forget to add it to go.work!"
