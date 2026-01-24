.PHONY: help build run test clean migrate proto

## help: Show this help message
help:
	@echo "Supplier Service - Available Make Targets:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the service binary
build:
	@echo "Building supplier-service..."
	go build -o bin/supplier-service cmd/main.go

## run: Run the service locally
run:
	@echo "Running supplier-service..."
	go run cmd/main.go

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## migrate-up: Run database migrations
migrate-up:
	@echo "Running migrations..."
	@for file in migrations/*up.sql; do \
		echo "Applying $$file..."; \
		PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -f $$file; \
	done

## migrate-down: Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	@for file in $$(ls migrations/*down.sql | sort -r); do \
		echo "Rolling back $$file..."; \
		PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -f $$file; \
	done

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	find . -name "*.test" -delete

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t erp-supplier-service:latest .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm --env-file .env -p 8084:8084 -p 9084:9084 erp-supplier-service:latest

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
