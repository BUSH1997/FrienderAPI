.PHONY: build
build:
	go build -o ./build/main_server/main ./cmd/main_server

ROOT_PATH = $(PWD)
BIN_OUTPUT_PATH = $(ROOT_PATH)/bin

.PHONY: run_public_sync
run_public_sync: build_public_sync
	$(BIN_OUTPUT_PATH)/public_sync

.PHONY: build_public_sync
build_public_sync:
	go build -o $(BIN_OUTPUT_PATH)/public_sync ./cmd/public_sync

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. --proto_path=$(ROOT_PATH)/grpc $(ROOT_PATH)/grpc/api.proto
