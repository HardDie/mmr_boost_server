PACKAGES := $(shell go list ./... | grep -v /vendor/)

.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build binary file
	CGO_ENABLED=0 go build -o server cmd/mmr_boost_server/main.go

.PHONY: linter
linter: ## run linters
	golangci-lint run --timeout 3m

.PHONY: test
test: ## run unit tests
	go test ./... -v -race -timeout 60s -count=1

.PHONY: coverage
coverage: ## show coverage of unit tests
	go test $(PACKAGES) -race -timeout 60s -count=1 -coverprofile cover.out || exit 1
	go tool cover -func cover.out
	rm cover.out

.PHONY: mocks
mocks: ## generate mocks
	mockery --all --dir internal --output ./internal/mocks --outpkg mocks

.PHONY: dependency
dependency: ## install dev dependency
	# gRPC generator
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# Mock generator
	go install github.com/vektra/mockery/v2@latest
	# Linter
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: proto
proto: ## generate go files from *.proto
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
