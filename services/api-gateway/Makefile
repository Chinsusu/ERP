.PHONY: help build run test clean

## help: Show this help message
help:
	@echo "API Gateway - Available Make Targets:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the service binary
build:
	@echo "Building api-gateway..."
	go build -o bin/api-gateway cmd/main.go

## run: Run the service locally
run:
	@echo "Running api-gateway..."
	go run cmd/main.go

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t erp-api-gateway:latest .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm --env-file .env -p 8080:8080 erp-api-gateway:latest

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
