.DEFAULT_GOAL := help
.PHONY: help install run build lint format api

help: ## List the back commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

install: ## Download the Go modules
	go mod download

run: ## Run the server with live reload
	air

build: ## Build the server binary into bin/
	go build -o bin/server ./cmd/server

lint: ## Lint, check formatting and architecture boundaries
	golangci-lint run ./...

format: ## Format the code
	golangci-lint fmt ./...

api: ## Regenerate the server code from api/openapi.yaml
	go tool oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
