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
