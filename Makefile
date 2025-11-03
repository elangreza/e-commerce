gen:
	protoc --proto_path=./gen/proto --go_out=./gen --go_opt=paths=source_relative \
		--go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
		--descriptor_set_out=./gen/e-commerce.protoset \
		./gen/proto/*.proto	
	@echo "Generated Go code from proto files."

.PHONY: gen
.DEFAULT_GOAL := gen
