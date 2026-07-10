SERVICES := nova-auth nova-user nova-post nova-file nova-gateway nova-web nova-apidoc
WIRE_SERVICES := nova-auth nova-user nova-post nova-file nova-web
GOLANGCI_LINT ?= golangci-lint
DOCKER_COMPOSE ?= docker-compose

.DEFAULT_GOAL := help

.PHONY: help proto proto-deps proto-generate proto-copy proto-clean openapi openapi-deps \
	openapi-validate openapi-generate openapi-generate-go-gin-server openapi-generate-ts-client \
	openapi-generate-swagger-json openapi-copy-apidoc openapi-copy-web-client openapi-clean \
	wire build test lint fmt clean docker-up docker-down docker-build dev dev-auth \
	dev-auth-client dev-user dev-post dev-file dev-gateway dev-web dev-apidoc

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  proto          Update buf deps, clean, generate, and copy generated files"
	@echo "  proto-deps     Update buf dependencies"
	@echo "  proto-generate Generate proto files under nova-contracts/gen/go/rpc"
	@echo "  proto-copy     Copy generated files from nova-contracts/gen to services"
	@echo "  proto-clean    Remove generated proto files"
	@echo "  openapi        Validate and generate OpenAPI outputs"
	@echo "  openapi-deps   Install OpenAPI generator npm dependencies"
	@echo "  openapi-validate Validate nova-contracts/openapi/openapi.yaml"
	@echo "  openapi-generate Generate OpenAPI Go server, TypeScript client, and swagger.json"
	@echo "  openapi-generate-ts-client Generate TypeScript frontend client and copy it to nova-web"
	@echo "  openapi-copy-apidoc Copy swagger.json to nova-apidoc"
	@echo "  openapi-copy-web-client Copy TypeScript frontend client to nova-web"
	@echo "  openapi-clean  Remove generated OpenAPI outputs"
	@echo "  wire           Generate Wire code for DI-based services"
	@echo "  build          Build all Go projects"
	@echo "  test           Run tests for all Go projects"
	@echo "  lint           Run contract lint/validation and golangci-lint"
	@echo "  fmt            Format all Go projects"
	@echo "  clean          Clean build artifacts"
	@echo "  docker-up      Start docker-compose services"
	@echo "  docker-down    Stop docker-compose services"
	@echo "  docker-build   Build docker-compose images"
	@echo "  dev            Start all services under the nova-launcher supervisor (Ctrl+C stops all)"
	@echo "  dev-*          Run a single project locally"

proto:
	$(MAKE) -C nova-contracts proto

proto-deps:
	$(MAKE) -C nova-contracts deps

proto-generate:
	$(MAKE) -C nova-contracts generate

proto-copy:
	$(MAKE) -C nova-contracts copy

proto-clean:
	$(MAKE) -C nova-contracts clean-all

openapi:
	$(MAKE) -C nova-contracts openapi

openapi-deps:
	$(MAKE) -C nova-contracts openapi-deps

openapi-validate:
	$(MAKE) -C nova-contracts openapi-validate

openapi-generate:
	$(MAKE) -C nova-contracts openapi-generate

openapi-generate-go-gin-server:
	$(MAKE) -C nova-contracts openapi-generate-go-gin-server

openapi-generate-ts-client:
	$(MAKE) -C nova-contracts openapi-generate-ts-client

openapi-generate-swagger-json:
	$(MAKE) -C nova-contracts openapi-generate-swagger-json

openapi-copy-apidoc:
	$(MAKE) -C nova-contracts openapi-copy-apidoc

openapi-copy-web-client:
	$(MAKE) -C nova-contracts openapi-copy-web-client

openapi-clean:
	$(MAKE) -C nova-contracts openapi-clean

wire:
	@set -e; for service in $(WIRE_SERVICES); do \
		echo ">>> wire: $$service"; \
		$(MAKE) -C $$service wire; \
	done

build:
	@set -e; for service in $(SERVICES); do \
		echo ">>> build: $$service"; \
		$(MAKE) -C $$service build; \
	done

test:
	@set -e; for service in $(SERVICES); do \
		echo ">>> test: $$service"; \
		$(MAKE) -C $$service test; \
	done

lint:
	$(MAKE) -C nova-contracts lint
	@set -e; for service in $(SERVICES); do \
		echo ">>> lint: $$service"; \
		(cd $$service && $(GOLANGCI_LINT) run ./...); \
	done

fmt:
	@set -e; for service in $(SERVICES); do \
		echo ">>> fmt: $$service"; \
		$(MAKE) -C $$service fmt; \
	done

clean:
	@set -e; for service in $(SERVICES); do \
		echo ">>> clean: $$service"; \
		$(MAKE) -C $$service clean; \
	done
	$(MAKE) -C nova-contracts clean

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-build:
	$(DOCKER_COMPOSE) build

# Start every service under a single supervisor (nova-launcher). Ctrl+C tears
# them all down atomically. Use ONLY=svc1,svc2 to launch a subset.
dev:
	cd nova-launcher && go run . $(if $(ONLY),-only $(ONLY))

dev-auth:
	$(MAKE) -C nova-auth run

dev-auth-client:
	$(MAKE) -C nova-auth run-client

dev-user:
	$(MAKE) -C nova-user run

dev-post:
	$(MAKE) -C nova-post run

dev-file:
	$(MAKE) -C nova-file run

dev-gateway:
	$(MAKE) -C nova-gateway run

dev-web:
	$(MAKE) -C nova-web run

dev-apidoc:
	$(MAKE) -C nova-apidoc run
