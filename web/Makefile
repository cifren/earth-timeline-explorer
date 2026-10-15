.DEFAULT_GOAL := help
.PHONY: help install dev build test test-integration lint format api

help: ## List the web commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

install: ## Install the npm dependencies
	npm install

dev: ## Run the development server
	npm run dev -- --host 0.0.0.0

build: ## Type-check and build the static bundle into dist/
	npm run build

test: ## Run the unit tests
	npm test

test-integration: ## Run the integration tests (pages driven like a user, API faked with MSW)
	npm run test:integration

lint: ## ESLint and Prettier check, including architecture boundaries
	npm run lint

format: ## Format the code with Prettier
	npm run format

api: ## Regenerate the API client from ../back/api/openapi.yaml
	npm run api
