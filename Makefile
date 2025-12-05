PROTO_DIR := proto

proto-gen:
	@echo "Generating Go code from proto files..."
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=k8s-sandbox \
		--go-grpc_opt=module=k8s-sandbox \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/*.proto
	@echo "Done generating proto code."

