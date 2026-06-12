SERVICES := nova-auth nova-user nova-post nova-file nova-gateway nova-web nova-apidoc
WIRE_SERVICES := nova-auth nova-user nova-post nova-file nova-web
GOLANGCI_LINT ?= golangci-lint
DOCKER_COMPOSE ?= docker-compose

.DEFAULT_GOAL := help

.PHONY: help proto proto-deps proto-generate proto-copy proto-clean wire build test lint fmt clean \
	docker-up docker-down docker-build dev-auth dev-auth-client dev-user dev-post dev-file dev-gateway \
	dev-web dev-apidoc

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  proto          Update buf deps, clean, generate, and copy generated files"
	@echo "  proto-deps     Update buf dependencies"
	@echo "  proto-generate Generate proto and OpenAPI files under nova-proto/gen"
	@echo "  proto-copy     Copy generated files from nova-proto/gen to services"
	@echo "  proto-clean    Remove generated proto/OpenAPI files"
	@echo "  wire           Generate Wire code for DI-based services"
	@echo "  build          Build all Go projects"
	@echo "  test           Run tests for all Go projects"
	@echo "  lint           Run buf lint and golangci-lint"
	@echo "  fmt            Format all Go projects"
	@echo "  clean          Clean build artifacts"
	@echo "  docker-up      Start docker-compose services"
	@echo "  docker-down    Stop docker-compose services"
	@echo "  docker-build   Build docker-compose images"
	@echo "  dev-*          Run a single project locally"

proto:
	$(MAKE) -C apis proto

proto-deps:
	$(MAKE) -C apis deps

proto-generate:
	$(MAKE) -C apis generate

proto-copy:
	$(MAKE) -C apis copy

proto-clean:
	$(MAKE) -C apis clean-all

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
	$(MAKE) -C apis lint
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
	$(MAKE) -C apis clean

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-build:
	$(DOCKER_COMPOSE) build

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
