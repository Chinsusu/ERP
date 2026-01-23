.PHONY: help docker-up docker-down docker-dev docker-logs health-check backup migrate proto test build clean

# Default target
.DEFAULT_GOAL := help

# Load environment variables from .env if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

## help: Show this help message
help:
	@echo "ERP Cosmetics - Available Make Targets:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""

## docker-up: Start all services in production mode
docker-up:
	@echo "Starting ERP Cosmetics services..."
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	@./scripts/health-check.sh

## docker-dev: Start all services in development mode with hot reload
docker-dev:
	@echo "Starting ERP Cosmetics in development mode..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
	@echo "Services started in development mode"
	@./scripts/health-check.sh

## docker-down: Stop all services
docker-down:
	@echo "Stopping all services..."
	docker-compose down

## docker-down-volumes: Stop all services and remove volumes (WARNING: data loss!)
docker-down-volumes:
	@echo "WARNING: This will delete all data! Press Ctrl+C to cancel..."
	@sleep 5
	docker-compose down -v

## docker-logs: Follow logs from all services
docker-logs:
	docker-compose logs -f

## docker-logs-service: Follow logs from a specific service (usage: make docker-logs-service SERVICE=auth-service)
docker-logs-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make docker-logs-service SERVICE=<service-name>"; \
		exit 1; \
	fi
	docker-compose logs -f $(SERVICE)

## health-check: Check health of all services
health-check:
	@./scripts/health-check.sh

## backup: Backup all databases and configurations
backup:
	@./scripts/backup.sh

## migrate: Run database migrations for all services
migrate:
	@echo "Running database migrations..."
	@for service in services/*-service; do \
		if [ -f "$$service/Makefile" ]; then \
			echo "Migrating $$service..."; \
			$(MAKE) -C $$service migrate || true; \
		fi \
	done

## migrate-service: Run migration for a specific service (usage: make migrate-service SERVICE=auth-service)
migrate-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make migrate-service SERVICE=<service-name>"; \
		exit 1; \
	fi
	$(MAKE) -C services/$(SERVICE) migrate

## proto: Generate protobuf code from .proto files
proto:
	@echo "Generating protobuf code..."
	@if [ -d "shared/proto" ]; then \
		cd shared/proto && \
		for file in *.proto; do \
			protoc --go_out=. --go_opt=paths=source_relative \
				   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
				   $$file; \
		done; \
		echo "Protobuf generation complete"; \
	else \
		echo "No proto directory found"; \
	fi

## test: Run tests for all services
test:
	@echo "Running tests..."
	@for service in services/*-service; do \
		if [ -f "$$service/Makefile" ]; then \
			echo "Testing $$service..."; \
			$(MAKE) -C $$service test || true; \
		fi \
	done

## test-service: Run tests for a specific service (usage: make test-service SERVICE=auth-service)
test-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make test-service SERVICE=<service-name>"; \
		exit 1; \
	fi
	$(MAKE) -C services/$(SERVICE) test

## build: Build all services
build:
	@echo "Building all services..."
	docker-compose build

## build-service: Build a specific service (usage: make build-service SERVICE=auth-service)
build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make build-service SERVICE=<service-name>"; \
		exit 1; \
	fi
	docker-compose build $(SERVICE)

## clean: Remove all built binaries and caches
clean:
	@echo "Cleaning up..."
	@for service in services/*-service; do \
		if [ -f "$$service/Makefile" ]; then \
			$(MAKE) -C $$service clean || true; \
		fi \
	done
	@find . -name "*.test" -type f -delete
	@find . -name "*.out" -type f -delete
	@echo "Cleanup complete"

## init: Initialize the project (create .env, setup directories)
init:
	@echo "Initializing ERP Cosmetics project..."
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "⚠️  Please edit .env and update passwords and secrets!"; \
	else \
		echo ".env already exists, skipping..."; \
	fi
	@mkdir -p deploy/nginx/ssl
	@mkdir -p deploy/grafana/provisioning/datasources
	@mkdir -p deploy/grafana/provisioning/dashboards
	@mkdir -p backups
	@chmod +x scripts/*.sh
	@echo "✓ Initialization complete"

## ps: Show running services
ps:
	docker-compose ps

## restart: Restart all services
restart:
	@echo "Restarting all services..."
	docker-compose restart

## restart-service: Restart a specific service (usage: make restart-service SERVICE=auth-service)
restart-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make restart-service SERVICE=<service-name>"; \
		exit 1; \
	fi
	docker-compose restart $(SERVICE)

## shell: Open a shell in a service container (usage: make shell SERVICE=auth-service)
shell:
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make shell SERVICE=<service-name>"; \
		exit 1; \
	fi
	docker-compose exec $(SERVICE) sh

## db-shell: Open PostgreSQL shell
db-shell:
	docker-compose exec postgres psql -U $(POSTGRES_USER:-postgres)

## redis-shell: Open Redis CLI
redis-shell:
	docker-compose exec redis redis-cli -a $(REDIS_PASSWORD)
