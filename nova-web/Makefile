GO ?= go
WIRE ?= wire
BIN_DIR ?= bin

SERVER_BIN := $(BIN_DIR)/nova-web

.DEFAULT_GOAL := help

.PHONY: help wire build run test fmt clean

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  wire   Generate Wire code"
	@echo "  build  Build cmd/server"
	@echo "  run    Run cmd/server"
	@echo "  test   Run tests"
	@echo "  fmt    Format Go code"
	@echo "  clean  Remove build artifacts"

wire:
	$(WIRE) ./internal/di

build:
	@mkdir -p "$(BIN_DIR)"
	$(GO) build -o "$(SERVER_BIN)" ./cmd/server

run:
	$(GO) run ./cmd/server

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -rf "$(BIN_DIR)"
