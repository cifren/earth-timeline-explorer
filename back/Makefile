.DEFAULT_GOAL := help
.PHONY: help install run build test test-integration lint format api zones

help: ## List the back commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

install: ## Download the Go modules
	go mod download

run: ## Run the server with live reload
	air

build: ## Build the server binary into bin/
	go build -o bin/server ./cmd/server

test: ## Run the unit tests with the race detector
	go test -race ./...

test-integration: ## Run the integration tests (needs the Redis of the stack)
	go test -race -tags integration ./...

lint: ## Lint, check formatting and architecture boundaries
	golangci-lint run ./...

format: ## Format the code
	golangci-lint fmt ./...

api: ## Regenerate the server code from api/openapi.yaml
	go tool oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml

zones: ## Regenerate the zone outlines and country points from Natural Earth
	go run ./cmd/zones
