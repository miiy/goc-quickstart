TARGET=server

.PHONY: proto
proto:
	@buf dep update
	@buf generate proto

.PHONY: help
help:
	@echo "make proto: proto file"
