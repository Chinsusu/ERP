.PHONY: help build run test clean migrate proto

## help: Show this help message
help:
	@echo "Auth Service - Available Make Targets:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the service binary
build:
	@echo "Building auth-service..."
	go build -o bin/auth-service cmd/main.go

## run: Run the service locally
run:
	@echo "Running auth-service..."
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

## proto: Generate code from protobuf
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	find . -name "*.test" -delete

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t erp-auth-service:latest .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm --env-file .env -p 8081:8081 -p 9081:9081 erp-auth-service:latest

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
