PROTO_DIR := proto
GO_MODULE := k8s-sandbox
BUILD_DIR := build

all: tidy fmt vet test api-build worker-build client-build

tidy: proto-gen
	go mod tidy

fmt: proto-gen
	go fmt ./...

test: proto-gen
	go test ./...

vet: proto-gen
	go vet ./...

proto-build:
	@echo "Generating Go code from proto files..."
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=$(GO_MODULE) \
		--go-grpc_opt=module=$(GO_MODULE) \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/*.proto
	@echo "Done generating proto code."

clean-proto:
	@echo "Cleaning generated protobuf files..."
	rm -rf internal/gen/proto
	@echo "Generated protobuf files cleaned."

proto-gen: clean-proto proto-build

api-build: proto-gen
	@echo "Building API server..."
	go build -o $(BUILD_DIR)/ ./cmd/api
	@echo "API server built."

worker-build: proto-gen
	@echo "Building worker..."
	go build -o $(BUILD_DIR)/ ./cmd/worker
	@echo "Worker built."

client-build: proto-gen
	@echo "Building client..."
	go build -o $(BUILD_DIR)/ ./cmd/client
	@echo "Client built."

clean-build:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	@echo "Build directory cleaned."

clean: clean-build clean-proto
