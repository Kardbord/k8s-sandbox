PROTO_DIR := proto
GO_MODULE := k8s-sandbox
BUILD_DIR := build
API_IMAGE := $(GO_MODULE)-api:local
WORKER_IMAGE := $(GO_MODULE)-worker:local
CLIENT_IMAGE := $(GO_MODULE)-client:local

all: tidy fmt vet test api-build worker-build client-build docker

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

docker: docker-api docker-worker docker-client

docker-api: clean-api-image proto-gen
	@echo "Building $(API_IMAGE)..."
	docker build -t $(API_IMAGE) -f Dockerfile.api .
	@echo "Finsihed $(API_IMAGE) build."

docker-worker: clean-worker-image proto-gen
	@echo "Building $(WORKER_IMAGE)..."
	docker build -t $(WORKER_IMAGE) -f Dockerfile.worker .
	@echo "Finsihed $(WORKER_IMAGE) build."

docker-client: clean-client-image proto-gen
	@echo "Building $(CLIENT_IMAGE)..."
	docker build -t $(CLIENT_IMAGE) -f Dockerfile.client .
	@echo "Finished $(CLIENT_IMAGE) build."

clean-api-image:
	@echo "Cleaning API Docker image..."
	docker rmi $(API_IMAGE) || true
	@echo "Cleaned API Docker image..."

clean-worker-image:
	@echo "Cleaning worker Docker image..."
	docker rmi $(WORKER_IMAGE) || true
	@echo "Cleaned worker Docker image..."

clean-client-image:
	@echo "Cleaning client Docker image..."
	docker rmi $(CLIENT_IMAGE) || true
	@echo "Cleaned client Docker image..."

clean-docker: clean-api-image clean-worker-image clean-client-image

clean: clean-docker clean-build clean-proto

