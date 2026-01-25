.PHONY: help build run test clean migrate

help:
	@echo "File Service - Available Make Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'

## build: Build the service binary
build:
	go build -o bin/file-service cmd/main.go

## run: Run the service locally
run:
	go run cmd/main.go

## test: Run all tests
test:
	go test -v ./...

## migrate-up: Run database migrations
migrate-up:
	@for file in migrations/*up.sql; do \
		echo "Applying $$file..."; \
		PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -f $$file; \
	done

## migrate-down: Rollback database migrations
migrate-down:
	@for file in $$(ls migrations/*down.sql | sort -r); do \
		echo "Rolling back $$file..."; \
		PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -f $$file; \
	done

## clean: Remove build artifacts
clean:
	rm -rf bin/

## docker-build: Build Docker image
docker-build:
	docker build -t erp-file-service:latest .

## deps: Download dependencies
deps:
	go mod download
	go mod tidy
