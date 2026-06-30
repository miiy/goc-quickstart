BUF ?= buf
NPM ?= npm
ROOT_DIR := $(abspath ..)

GO_RPC_GEN_DIR := gen/go/rpc
GO_HTTP_GEN_DIR := gen/go/http
APIDOC_OPENAPI_DIR := $(ROOT_DIR)/nova-apidoc/gen/openapi

# Each backend service only needs its OWN proto package; the gateway needs all.
# BASE is the whole `<svc>/gen/go/nova` tree — wiping it on every copy makes
# `make copy` self-cleaning, so stale off-target proto packages can't linger
# across services (e.g. nova-auth carrying user/post/file).
AUTH_BASE    := $(ROOT_DIR)/nova-auth/gen/go/nova
POST_BASE    := $(ROOT_DIR)/nova-post/gen/go/nova
USER_BASE    := $(ROOT_DIR)/nova-user/gen/go/nova
FILE_BASE    := $(ROOT_DIR)/nova-file/gen/go/nova
GATEWAY_BASE := $(ROOT_DIR)/nova-gateway/gen/go/nova

AUTH_DST := $(AUTH_BASE)/auth/v1
POST_DST := $(POST_BASE)/post/v1
USER_DST := $(USER_BASE)/user/v1
FILE_DST := $(FILE_BASE)/file/v1
GATEWAY_DST := $(GATEWAY_BASE)

.DEFAULT_GOAL := help

.PHONY: help proto deps generate copy copy-go copy-auth copy-post copy-user copy-file \
	copy-gateway openapi openapi-deps openapi-validate openapi-generate \
	openapi-generate-go-gin-server openapi-generate-go-client openapi-generate-swagger-json \
	openapi-copy-apidoc openapi-fmt openapi-clean lint format clean clean-copies clean-all

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  proto          Update deps, clean, generate, and copy generated files"
	@echo "  deps           Update buf dependencies"
	@echo "  generate       Generate Go protobuf/gRPC files under gen/go/rpc"
	@echo "  copy           Copy generated Go files to services"
	@echo "  openapi        Validate and generate OpenAPI outputs"
	@echo "  openapi-deps   Install OpenAPI generator npm dependencies"
	@echo "  openapi-validate Validate openapi/openapi.yaml"
	@echo "  openapi-generate Generate OpenAPI Go server/client and swagger.json"
	@echo "  openapi-copy-apidoc Copy swagger.json to nova-apidoc"
	@echo "  openapi-clean  Remove gen/go/http and copied nova-apidoc OpenAPI files"
	@echo "  clean          Remove nova-contracts/gen/go/rpc"
	@echo "  clean-copies   Remove copied generated files from services"
	@echo "  clean-all      Remove nova-contracts/gen/go/rpc and copied generated files"
	@echo "  lint           Run buf lint and OpenAPI validation"
	@echo "  format         Format proto files"

proto: deps clean generate copy

deps:
	$(BUF) dep update

generate:
	$(BUF) generate

copy: copy-go

copy-go: copy-auth copy-post copy-user copy-file copy-gateway

copy-auth:
	@rm -rf "$(AUTH_BASE)"
	@mkdir -p "$(AUTH_DST)"
	cp -R "$(GO_RPC_GEN_DIR)/nova/auth/v1/." "$(AUTH_DST)/"

copy-post:
	@rm -rf "$(POST_BASE)"
	@mkdir -p "$(POST_DST)"
	cp -R "$(GO_RPC_GEN_DIR)/nova/post/v1/." "$(POST_DST)/"

copy-user:
	@rm -rf "$(USER_BASE)"
	@mkdir -p "$(USER_DST)"
	cp -R "$(GO_RPC_GEN_DIR)/nova/user/v1/." "$(USER_DST)/"

copy-file:
	@rm -rf "$(FILE_BASE)"
	@mkdir -p "$(FILE_DST)"
	cp -R "$(GO_RPC_GEN_DIR)/nova/file/v1/." "$(FILE_DST)/"

copy-gateway:
	@rm -rf "$(GATEWAY_DST)"
	@mkdir -p "$(GATEWAY_DST)"
	cp -R "$(GO_RPC_GEN_DIR)/nova/." "$(GATEWAY_DST)/"

openapi: openapi-validate openapi-generate

openapi-deps:
	$(NPM) install

openapi-validate:
	$(NPM) run openapi:validate

openapi-generate:
	$(NPM) run openapi:generate

openapi-generate-go-gin-server:
	$(NPM) run openapi:generate:go-gin-server
	$(MAKE) openapi-fmt

openapi-generate-go-client:
	$(NPM) run openapi:generate:go-client
	$(MAKE) openapi-fmt

openapi-generate-swagger-json:
	$(NPM) run openapi:generate:swagger-json

openapi-copy-apidoc:
	$(NPM) run openapi:copy:apidoc

openapi-fmt:
	$(NPM) run openapi:fmt:go

openapi-clean:
	rm -rf "$(GO_HTTP_GEN_DIR)" "$(APIDOC_OPENAPI_DIR)"

lint:
	$(BUF) lint
	$(NPM) run openapi:validate

format:
	$(BUF) format -w

clean:
	rm -rf "$(GO_RPC_GEN_DIR)"

clean-copies:
	rm -rf "$(AUTH_BASE)" "$(POST_BASE)" "$(USER_BASE)" "$(FILE_BASE)" "$(GATEWAY_DST)"

clean-all: clean clean-copies
