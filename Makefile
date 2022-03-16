.PHONY: all
all:
	@$(MAKE) --no-print-directory deps
	@$(MAKE) --no-print-directory protobuf
	@$(MAKE) --no-print-directory app

.PHONY: deps
deps:
	go mod tidy

.PHONY: protobuf
protobuf:
	### 1. grpc first
	protoc -I ./pb -I /usr/local/include \
	--go_out ./pb \
	--go_opt paths=source_relative \
	--go-grpc_out ./pb \
	--go-grpc_opt paths=source_relative \
	./pb/*.proto

	### 2. gateway of restful
	protoc -I ./pb -I /usr/local/include \
	--grpc-gateway_out ./pb \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out ./pb \
	--openapiv2_opt logtostderr=true \
	./pb/*.proto

.PHONY: app
app:
	go build

