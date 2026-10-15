.DEFAULT_GOAL := help
.PHONY: help install dev build lint format api

help: ## List the web commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-10s %s\n", $$1, $$2}'

install: ## Install the npm dependencies
	npm install

dev: ## Run the development server
	npm run dev -- --host 0.0.0.0

build: ## Build the static bundle into dist/
	npm run build

lint: ## ESLint and Prettier check, including architecture boundaries
	npx eslint .
	npx prettier --check .

format: ## Format the code with Prettier
	npx prettier --write .

api: ## Regenerate the API client from ../back/api/openapi.yaml
	npx orval
