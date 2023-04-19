.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: swagger
swagger: ## generate swagger file
	swagger generate spec -m -o swagger.yaml

.PHONY: swagger-install
swagger-install: ## install swagger for linux
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest

.PHONY: build
build: ## build binary file
	CGO_ENABLED=0 go build -o server cmd/mmr_boost_server/main.go

.PHONY: proto
proto:
	protoc -I./pkg/proto/server \
		-I./pkg/proto \
		--go_out ./pkg/proto/server \
		--go_opt=paths=source_relative \
		--go-grpc_out ./pkg/proto/server \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out ./pkg/proto/server \
		--grpc-gateway_opt=paths=source_relative \
		--openapiv2_out ./ \
		--openapiv2_opt allow_merge=true,merge_file_name=api,output_format=yaml \
		./pkg/proto/server/*.proto
