PROTO_DIR := proto

proto-build:
	@echo "Generating Go code from proto files..."
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=k8s-sandbox \
		--go-grpc_opt=module=k8s-sandbox \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/*.proto
	@echo "Done generating proto code."

clean-proto:
	rm -rf internal/gen/proto

proto-gen: clean-proto proto-build

clean: clean-proto
