gen:
	protoc --proto_path=./gen/proto --go_out=./gen --go_opt=paths=source_relative \
		--go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
		--descriptor_set_out=./gen/e-commerce.protoset \
		./gen/proto/*.proto	
	@echo "Generated Go code from proto files."

gen-common:
	protoc --proto_path=./gen/proto/common --go_out=./gen/proto/common --go_opt=paths=source_relative \
		--go-grpc_out=./gen/proto/common --go-grpc_opt=paths=source_relative \
		--descriptor_set_out=./gen/e-commerce.protoset \
		./gen/proto/common/*.proto	
	@echo "Generated Go code from proto files."

.PHONY: gen gen-common
.DEFAULT_GOAL := gen
