PROTO_DIR := proto
GO_MODULE := k8s-sandbox
BUILD_DIR := build
API_IMAGE := $(GO_MODULE)-api:local
WORKER_IMAGE := $(GO_MODULE)-worker:local
CLIENT_IMAGE := $(GO_MODULE)-client:local

all: build test build docker-images kube-deploy

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

build: tidy fmt vet api-build worker-build client-build

clean-build:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	@echo "Build directory cleaned."

docker-images: docker-api-image docker-worker-image docker-client-image

docker-api-image: clean-api-image proto-gen
	@echo "Building $(API_IMAGE)..."
	docker build -t $(API_IMAGE) -f Dockerfile.api .
	@echo "Finsihed $(API_IMAGE) build."

docker-worker-image: clean-worker-image proto-gen
	@echo "Building $(WORKER_IMAGE)..."
	docker build -t $(WORKER_IMAGE) -f Dockerfile.worker .
	@echo "Finsihed $(WORKER_IMAGE) build."

docker-client-image: clean-client-image proto-gen
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

clean-docker-images: clean-api-image clean-worker-image clean-client-image

start-db:
	docker compose -f ./deployments/db/docker-compose.yml up -d

stop-db:
	docker compose -f ./deployments/db/docker-compose.yml down

start-redis:
	docker compose -f ./deployments/redis/docker-compose.yml up -d

stop-redis:
	docker compose -f ./deployments/redis/docker-compose.yml down

K8S_DIR := deployments

kube-deploy-api: docker-api-image
	kubectl apply -f $(K8S_DIR)/api

kube-deploy-worker: docker-worker-image
	kubectl apply -f $(K8S_DIR)/worker

kube-deploy-client: docker-client-image
	kubectl apply -f $(K8S_DIR)/client

kube-deploy: start-db start-redis kube-deploy-api kube-deploy-worker kube-deploy-client

kube-clean-api:
	kubectl delete --ignore-not-found=true -f $(K8S_DIR)/api

kube-clean-worker:
	kubectl delete --ignore-not-found=true -f $(K8S_DIR)/worker

kube-clean-client:
	kubectl delete --ignore-not-found=true -f $(K8S_DIR)/client

clean-kube: kube-clean-client kube-clean-api kube-clean-worker

clean: clean-kube stop-redis stop-db clean-docker-images clean-build clean-proto

