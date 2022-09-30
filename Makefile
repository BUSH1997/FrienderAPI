ROOT_PATH = $(PWD)
BIN_OUTPUT_PATH = $(ROOT_PATH)/bin

.PHONY: api
api: build_api
	$(BIN_OUTPUT_PATH)/api

build_api:
	go build -o $(BIN_OUTPUT_PATH)/api ./cmd/api

.PHONY: internalapi
internalapi: build_internalapi
	$(BIN_OUTPUT_PATH)/internalapi

build_internalapi:
	go build -o $(BIN_OUTPUT_PATH)/internalapi ./cmd/internalapi

.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. --proto_path=$(ROOT_PATH)/grpc $(ROOT_PATH)/grpc/api.proto

