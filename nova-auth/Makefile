GO ?= go
WIRE ?= wire
BIN_DIR ?= bin

SERVER_BIN := $(BIN_DIR)/nova-auth
CLIENT_BIN := $(BIN_DIR)/client

.DEFAULT_GOAL := help

.PHONY: help wire build build-server build-client run run-client test fmt clean

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  wire         Generate Wire code"
	@echo "  build        Build server and client"
	@echo "  build-server Build cmd/server"
	@echo "  build-client Build cmd/client"
	@echo "  run          Run cmd/server"
	@echo "  run-client   Run cmd/client"
	@echo "  test         Run tests"
	@echo "  fmt          Format Go code"
	@echo "  clean        Remove build artifacts"

wire:
	$(WIRE) ./internal/di

build: build-server build-client

build-server:
	@mkdir -p "$(BIN_DIR)"
	$(GO) build -o "$(SERVER_BIN)" ./cmd/server

build-client:
	@mkdir -p "$(BIN_DIR)"
	$(GO) build -o "$(CLIENT_BIN)" ./cmd/client

run:
	$(GO) run ./cmd/server

run-client:
	$(GO) run ./cmd/client

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf "$(BIN_DIR)"
