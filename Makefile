# Encore - Makefile
# Usage: make <target>

run: ## Run the server locally
	go run cmd/server/main.go

lint: ## Run linter
	golangci-lint run ./...