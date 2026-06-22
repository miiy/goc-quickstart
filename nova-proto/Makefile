BUF ?= buf
ROOT_DIR := $(abspath ..)

GO_GEN_DIR := gen/go
OPENAPI_GEN_DIR := gen/openapiv2

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
APIDOC_DST := $(ROOT_DIR)/nova-apidoc/gen/openapiv2

.DEFAULT_GOAL := help

.PHONY: help proto deps generate copy copy-go copy-openapi copy-auth copy-post copy-user copy-file \
	copy-gateway copy-apidoc lint format clean clean-copies clean-all

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  proto          Update deps, clean, generate, and copy generated files"
	@echo "  deps           Update buf dependencies"
	@echo "  generate       Generate Go and OpenAPI files under gen/"
	@echo "  copy           Copy generated files to services"
	@echo "  clean          Remove nova-proto/gen"
	@echo "  clean-copies   Remove copied generated files from services"
	@echo "  clean-all      Remove nova-proto/gen and copied generated files"
	@echo "  lint           Run buf lint"
	@echo "  format         Format proto files"

proto: deps clean generate copy

deps:
	$(BUF) dep update

generate:
	$(BUF) generate

copy: copy-go copy-openapi

copy-go: copy-auth copy-post copy-user copy-file copy-gateway

copy-openapi: copy-apidoc

copy-auth:
	@rm -rf "$(AUTH_BASE)"
	@mkdir -p "$(AUTH_DST)"
	cp -R "$(GO_GEN_DIR)/nova/auth/v1/." "$(AUTH_DST)/"

copy-post:
	@rm -rf "$(POST_BASE)"
	@mkdir -p "$(POST_DST)"
	cp -R "$(GO_GEN_DIR)/nova/post/v1/." "$(POST_DST)/"

copy-user:
	@rm -rf "$(USER_BASE)"
	@mkdir -p "$(USER_DST)"
	cp -R "$(GO_GEN_DIR)/nova/user/v1/." "$(USER_DST)/"

copy-file:
	@rm -rf "$(FILE_BASE)"
	@mkdir -p "$(FILE_DST)"
	cp -R "$(GO_GEN_DIR)/nova/file/v1/." "$(FILE_DST)/"

copy-gateway:
	@rm -rf "$(GATEWAY_DST)"
	@mkdir -p "$(GATEWAY_DST)"
	cp -R "$(GO_GEN_DIR)/nova/." "$(GATEWAY_DST)/"

copy-apidoc:
	@rm -rf "$(APIDOC_DST)"
	@mkdir -p "$(APIDOC_DST)"
	cp -R "$(OPENAPI_GEN_DIR)/." "$(APIDOC_DST)/"

lint:
	$(BUF) lint

format:
	$(BUF) format -w

clean:
	rm -rf "$(GO_GEN_DIR)" "$(OPENAPI_GEN_DIR)"

clean-copies:
	rm -rf "$(AUTH_BASE)" "$(POST_BASE)" "$(USER_BASE)" "$(FILE_BASE)" "$(GATEWAY_DST)" "$(APIDOC_DST)"

clean-all: clean clean-copies
