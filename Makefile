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

## Docker

up: ## Start PostgreSQL and app in Docker
	docker compose up -d

down: ## Stop all containers
	docker compose down

db: ## Start only PostgreSQL
	docker compose up -d db

migrate-up: ## Run migrations
	docker compose up migrate

migrate-down: ## Rollback last migration
	docker compose run --rm migrate -path /migrations -database "postgres://encore:encore@db:5432/encore?sslmode=disable" down 1

## Help

help: ## Show available targets
	@echo   run          - Run the server locally
	@echo   lint         - Run linters
	@echo   format       - Run formatters (gofumpt + gci)
	@echo   check        - Format and lint in one command
	@echo   up           - Start PostgreSQL and app in Docker
	@echo   down         - Stop all containers
	@echo   db           - Start only PostgreSQL
	@echo   migrate-up   - Run migrations
	@echo   migrate-down - Rollback last migration
	@echo   help         - Show available targets

.PHONY: run lint format check up down db migrate-up migrate-down help