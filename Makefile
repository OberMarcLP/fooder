.PHONY: help all backend frontend db db-stop clean install

# Load environment variables from .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/Makefile://' | sort | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

all: db ## Start database (run 'make backend' and 'make frontend' in separate terminals)
	@echo "Database is running. In separate terminals, run:"
	@echo "  make backend"
	@echo "  make frontend"

db: ## Start PostgreSQL database
	@echo "Starting PostgreSQL..."
	@docker compose up -d db
	@echo "Waiting for database to be ready..."
	@until docker exec fooder-db pg_isready -U fooder -d fooder > /dev/null 2>&1; do sleep 1; done
	@echo "Database ready at postgres://fooder:fooder_secret@localhost:5432/fooder"

db-stop: ## Stop PostgreSQL database
	@docker compose stop db

backend: ## Run backend with Go
	@cd backend && go run ./cmd/server

frontend: install ## Run frontend with Vite dev server
	@cd frontend && npm run dev

install: ## Install frontend dependencies
	@cd frontend && npm install

clean: ## Stop and remove database container and volume
	@docker compose down -v
