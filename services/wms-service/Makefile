.PHONY: build run test clean migrate-up migrate-down docker-build docker-run deps proto

# Variables
APP_NAME=wms-service
BUILD_DIR=./bin
MAIN_FILE=./cmd/main.go

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

# Run with hot reload (requires air)
dev:
	@air

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run WMS service tests
test-wms:
	@echo "Running WMS service tests..."
	@go test -v ./internal/domain/entity/... ./internal/usecase/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep "internal/domain/entity/stock_fefo.go"
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	@migrate -path ./migrations -database "postgres://postgres:postgres123@localhost:5432/wms_db?sslmode=disable" up

# Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	@migrate -path ./migrations -database "postgres://postgres:postgres123@localhost:5432/wms_db?sslmode=disable" down

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8086:8086 --env-file .env $(APP_NAME):latest

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go-grpc_out=. api/proto/*.proto

# Lint the code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  migrate-up    - Run database migrations up"
	@echo "  migrate-down  - Run database migrations down"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  deps          - Download dependencies"
	@echo "  proto         - Generate protobuf files"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
