.PHONY: build run test clean migrate-up migrate-down

# Build the application
build:
	@echo "Building Manufacturing Service..."
	@go build -o bin/manufacturing-service ./cmd/main.go

# Run the application
run:
	@go run ./cmd/main.go

# Run tests
test:
	@echo "Running all tests..."
	@go test ./... -v

# Run Manufacturing tests
test-mfg:
	@echo "Running Manufacturing Service tests..."
	@go test -v ./internal/domain/entity/... ./internal/usecase/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep "internal/domain/entity/bom.go"
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@rm -rf bin/

# Run database migrations up
migrate-up:
	@migrate -path ./migrations -database "postgres://postgres:postgres123@localhost:5438/manufacturing_db?sslmode=disable" up

# Run database migrations down
migrate-down:
	@migrate -path ./migrations -database "postgres://postgres:postgres123@localhost:5438/manufacturing_db?sslmode=disable" down

# Download dependencies
deps:
	@go mod download
	@go mod tidy

# Generate mocks (if needed)
mocks:
	@mockgen -source=internal/domain/repository/repository.go -destination=internal/mocks/repository_mock.go -package=mocks

# Docker build
docker-build:
	@docker build -t manufacturing-service .

# Docker run
docker-run:
	@docker run -p 8087:8087 -p 9087:9087 manufacturing-service
