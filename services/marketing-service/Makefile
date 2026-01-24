.PHONY: build run test clean migrate-up migrate-down docker-build docker-run deps proto lint fmt

# Build the service
build:
	go build -o bin/marketing-service ./cmd/main.go

# Run the service
run: build
	./bin/marketing-service

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/marketing_db?sslmode=disable" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/marketing_db?sslmode=disable" down

# Build Docker image
docker-build:
	docker build -t marketing-service:latest .

# Run Docker container
docker-run:
	docker run -p 8089:8089 -p 9089:9089 --env-file .env marketing-service:latest

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate protobuf
proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
