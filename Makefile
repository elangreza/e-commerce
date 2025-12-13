gen:
	protoc --proto_path=./gen/proto --go_out=./gen --go_opt=paths=source_relative \
		--go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
		--descriptor_set_out=./gen/e-commerce.protoset \
		./gen/proto/*.proto	
	@echo "Generated Go code from proto files."

.PHONY: build-builder build-runtime

build-builder:
	docker build -t e-commerce/cgo-builder:latest -f images/cgo/Dockerfile .

build-runtime:
	docker build -t e-commerce/runtime-base:latest -f images/runtime-base/Dockerfile .

build: build-builder build-runtime
	cp ./api/env.example ./api/api.env
	cp ./order/env.example ./order/order.env
	cp ./product/env.example ./product/product.env
	cp ./warehouse/env.example ./warehouse/warehouse.env
	cp ./shop/env.example ./shop/shop.env
	cp ./payment/env.example ./payment/payment.env
	docker compose up --build

copy-to-local:
	@echo "Copying Docker data to local development..."
	@chmod +x ./copy-docker-data.sh && ./copy-docker-data.sh to-local

copy-to-docker:
	@echo "Copying local development data to Docker..."
	@chmod +x ./copy-docker-data.sh && ./copy-docker-data.sh to-docker

.PHONY: gen build copy-to-local copy-to-docker
.DEFAULT_GOAL := gen
