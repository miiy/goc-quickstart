GO ?= go
NPM ?= npm
WIRE ?= wire
BIN_DIR ?= bin

SERVER_BIN := $(BIN_DIR)/nova-web

.DEFAULT_GOAL := help

.PHONY: help wire frontend-install frontend-build frontend-dev dev build run test fmt clean

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  wire              Generate Wire code"
	@echo "  frontend-install  Install frontend dependencies"
	@echo "  frontend-build    Build React assets"
	@echo "  frontend-dev      Run Vite dev server"
	@echo "  dev               Run Vite dev server and Go server"
	@echo "  build             Build React assets and cmd/server"
	@echo "  run               Run cmd/server"
	@echo "  test              Run tests"
	@echo "  fmt               Format Go code"
	@echo "  clean             Remove build artifacts"

wire:
	$(WIRE) ./internal/di

frontend-install:
	$(NPM) ci

frontend-build: frontend-install
	$(NPM) run build

frontend-dev:
	$(NPM) run dev

dev:
	@trap 'kill 0' INT TERM EXIT; \
	$(NPM) run dev & \
	$(GO) run ./cmd/server

build: frontend-build
	@mkdir -p "$(BIN_DIR)"
	$(GO) build -o "$(SERVER_BIN)" ./cmd/server

run: frontend-build
	$(GO) run ./cmd/server

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf "$(BIN_DIR)" dist
