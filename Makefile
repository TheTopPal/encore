# Encore - Makefile
# Usage: make <target>

## App

run: ## Run the server locally
	go run cmd/server/main.go

## Code Quality

lint: ## Run linters
	golangci-lint run ./...

format: ## Run formatters (gofumpt + gci)
	golangci-lint fmt ./...

check: format lint ## Format and lint in one command

## Help

help: ## Show available targets
	@echo   run        - Run the server locally
	@echo   lint       - Run linters
	@echo   format     - Run formatters (gofumpt + gci)
	@echo   check      - Format and lint in one command
	@echo   help       - Show available targets

.PHONY: run lint format check help