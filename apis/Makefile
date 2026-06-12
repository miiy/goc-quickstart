BUF ?= buf
ROOT_DIR := $(abspath ..)

GO_GEN_DIR := gen/go
OPENAPI_GEN_DIR := gen/openapiv2

AUTH_DST := $(ROOT_DIR)/auth-service/gen/go/blog/auth/v1
POST_DST := $(ROOT_DIR)/post-service/gen/go/blog/post/v1
USER_DST := $(ROOT_DIR)/user-service/gen/go/blog/user/v1
UPLOAD_DST := $(ROOT_DIR)/upload-service/gen/go/blog/upload/v1
GATEWAY_DST := $(ROOT_DIR)/api-gateway/gen/go/blog
APIDOC_DST := $(ROOT_DIR)/apidoc-server/gen/openapiv2

.DEFAULT_GOAL := help

.PHONY: help proto deps generate copy copy-go copy-openapi copy-auth copy-post copy-user copy-upload \
	copy-gateway copy-apidoc lint format clean clean-copies clean-all

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  proto          Update deps, clean, generate, and copy generated files"
	@echo "  deps           Update buf dependencies"
	@echo "  generate       Generate Go and OpenAPI files under gen/"
	@echo "  copy           Copy generated files to services"
	@echo "  clean          Remove apis/gen"
	@echo "  clean-copies   Remove copied generated files from services"
	@echo "  clean-all      Remove apis/gen and copied generated files"
	@echo "  lint           Run buf lint"
	@echo "  format         Format proto files"

proto: deps clean generate copy

deps:
	$(BUF) dep update

generate:
	$(BUF) generate

copy: copy-go copy-openapi

copy-go: copy-auth copy-post copy-user copy-upload copy-gateway

copy-openapi: copy-apidoc

copy-auth:
	@rm -rf "$(AUTH_DST)"
	@mkdir -p "$(AUTH_DST)"
	cp -R "$(GO_GEN_DIR)/blog/auth/v1/." "$(AUTH_DST)/"

copy-post:
	@rm -rf "$(POST_DST)"
	@mkdir -p "$(POST_DST)"
	cp -R "$(GO_GEN_DIR)/blog/post/v1/." "$(POST_DST)/"

copy-user:
	@rm -rf "$(USER_DST)"
	@mkdir -p "$(USER_DST)"
	cp -R "$(GO_GEN_DIR)/blog/user/v1/." "$(USER_DST)/"

copy-upload:
	@rm -rf "$(UPLOAD_DST)"
	@mkdir -p "$(UPLOAD_DST)"
	cp -R "$(GO_GEN_DIR)/blog/upload/v1/." "$(UPLOAD_DST)/"

copy-gateway:
	@rm -rf "$(GATEWAY_DST)"
	@mkdir -p "$(GATEWAY_DST)"
	cp -R "$(GO_GEN_DIR)/blog/." "$(GATEWAY_DST)/"

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
	rm -rf "$(AUTH_DST)" "$(POST_DST)" "$(USER_DST)" "$(GATEWAY_DST)" "$(APIDOC_DST)"

clean-all: clean clean-copies
