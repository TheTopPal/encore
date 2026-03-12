# Encore

Concert ticket platform.

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [golangci-lint v2](https://golangci-lint.run/welcome/install/)
- [Docker & Docker Compose](https://docs.docker.com/get-started/get-docker/)

## Stack

- **Language:** Go 1.26
- **Router:** Chi
- **Database:** PostgreSQL 17 (pgx)
- **Migrations:** golang-migrate

## Getting Started

```bash
# Clone
git clone https://github.com/TheTopPal/encore
cd encore

# Start PostgreSQL
make db

# Run migrations
make migrate-up

# Run the server
make run
```

## Make Targets

```
make run          - Run the server locally
make lint         - Run linters
make format       - Run formatters (gofumpt + gci)
make check        - Format and lint in one command
make up           - Start PostgreSQL and app in Docker
make down         - Stop all containers
make db           - Start only PostgreSQL
make migrate-up   - Run migrations
make migrate-down - Rollback last migration
```
