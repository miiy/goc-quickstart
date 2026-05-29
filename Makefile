.PHONY: proto build test clean docker-up docker-down help

# Variables
SERVICES := auth-service user-service post-service api-gateway web
GO := go
BUF := buf
WIRE := wire

# Default target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  proto       Generate proto files"
	@echo "  wire        Generate wire dependencies"
	@echo "  build       Build all services"
	@echo "  test        Run all tests"
	@echo "  clean       Clean build artifacts"
	@echo "  docker-up   Start services with docker-compose"
	@echo "  docker-down Stop docker-compose services"
	@echo "  lint        Run linter"
	@echo "  fmt         Format code"

# Generate proto files
proto:
	@echo ">>> Generating proto files..."
	cd apis && $(BUF) dep update && $(BUF) generate proto
	@echo ">>> Copying generated files to services..."
	cp -r apis/gen/go/blog/auth/v1/* auth-service/gen/go/blog/auth/v1/
	cp -r apis/gen/go/blog/post/v1/* post-service/gen/go/blog/post/v1/
	cp -r apis/gen/go/blog/user/v1/* user-service/gen/go/blog/user/v1/
	cp -r apis/gen/go/blog/* api-gateway/gen/go/blog/
	@echo ">>> Proto generation complete"

# Generate wire dependencies
wire:
	@echo ">>> Generating wire dependencies..."
	for service in auth-service post-service user-service; do \
		echo "  -> $$service"; \
		cd $$service && $(WIRE) ./internal/di/ && cd ..; \
	done
	@echo ">>> Wire generation complete"

# Build all services
build:
	@echo ">>> Building all services..."
	for service in $(SERVICES); do \
		echo "  -> $$service"; \
		cd $$service && $(GO) build -o bin/server ./cmd/... && cd ..; \
	done
	@echo ">>> Build complete"

# Run tests
test:
	@echo ">>> Running tests..."
	for service in $(SERVICES); do \
		echo "  -> $$service"; \
		cd $$service && $(GO) test -v ./... && cd ..; \
	done

# Clean build artifacts
clean:
	@echo ">>> Cleaning build artifacts..."
	for service in $(SERVICES); do \
		rm -rf $$service/bin/; \
	done
	rm -rf apis/gen/
	@echo ">>> Clean complete"

# Docker commands
docker-up:
	@echo ">>> Starting services with docker-compose..."
	docker-compose up -d

docker-down:
	@echo ">>> Stopping docker-compose services..."
	docker-compose down

docker-build:
	@echo ">>> Building docker images..."
	docker-compose build

# Lint
lint:
	@echo ">>> Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo ">>> Formatting code..."
	for service in $(SERVICES); do \
		cd $$service && $(GO) fmt ./... && cd ..; \
	done
	@echo ">>> Format complete"

# Development
dev-auth:
	cd auth-service && $(GO) run cmd/server/main.go

dev-user:
	cd user-service && $(GO) run cmd/server/main.go

dev-post:
	cd post-service && $(GO) run cmd/server/main.go

dev-gateway:
	cd api-gateway && $(GO) run cmd/main.go

dev-web:
	cd web && $(GO) run cmd/server/main.go
