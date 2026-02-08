.PHONY: dev dev-build dev-down dev-logs test test-api test-web build build-api build-web migrate migrate-up migrate-down migrate-create clean help

# Development
dev: ## Start all services in development mode
	docker compose up

dev-build: ## Build and start all services in development mode
	docker compose up --build

dev-down: ## Stop all services
	docker compose down

dev-logs: ## View logs from all services
	docker compose logs -f

# Testing
test: test-api test-web ## Run all tests

test-api: ## Run API tests
	cd api && go test ./...

test-web: ## Run web tests
	cd web && npm test

# Build
build: build-api build-web ## Build all production artifacts

build-api: ## Build API binary
	cd api && go build -ldflags="-w -s" -o bin/server ./cmd/server

build-web: ## Build web assets
	cd web && npm run build

# Database migrations
migrate: migrate-up ## Alias for migrate-up

migrate-up: ## Run database migrations up
	cd api && go run ./cmd/migrate up

migrate-down: ## Rollback last migration
	cd api && go run ./cmd/migrate down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	cd api && go run ./cmd/migrate create $(name)

# Docker production
prod: ## Start all services in production mode
	docker compose -f docker-compose.prod.yml up -d

prod-build: ## Build production images
	docker compose -f docker-compose.prod.yml build

prod-down: ## Stop production services
	docker compose -f docker-compose.prod.yml down

# Utilities
clean: ## Clean build artifacts and docker volumes
	rm -rf api/bin api/tmp
	rm -rf web/dist web/node_modules/.cache
	docker compose down -v

db-shell: ## Open PostgreSQL shell
	docker compose exec db psql -U memwright -d memwright

api-shell: ## Open shell in API container
	docker compose exec api sh

web-shell: ## Open shell in web container
	docker compose exec web sh

# Help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
