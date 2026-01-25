.PHONY: help build-all build-services build-frontend deploy clean

VERSION ?= $(shell date +%Y%m%d)-$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

help:
	@echo "ERP Build & Deploy Commands"
	@echo "============================"
	@echo "make build-all        - Build all services and frontend"
	@echo "make build-services   - Build all backend services"
	@echo "make build-frontend   - Build frontend"
	@echo "make deploy           - Deploy to production"
	@echo "make clean            - Clean build artifacts"
	@echo ""
	@echo "Current VERSION: $(VERSION)"

build-all: build-services build-frontend

build-services:
	@echo "Building all services with version: $(VERSION)"
	@for svc in api-gateway auth-service user-service master-data-service \
	           supplier-service procurement-service wms-service \
	           manufacturing-service sales-service marketing-service \
	           notification-service file-service reporting-service; do \
		echo "Building $$svc..."; \
		docker build \
			--build-arg VERSION=$(VERSION) \
			-t erp/$$svc:$(VERSION) \
			-t erp/$$svc:latest \
			-f services/$$svc/Dockerfile \
			--build-context shared=./shared \
			. || exit 1; \
	done
	@echo "✅ All services built successfully!"

build-frontend:
	@echo "Building frontend..."
	@cd frontend && docker build -t erp/frontend:$(VERSION) -t erp/frontend:latest .
	@echo "✅ Frontend built successfully!"

deploy:
	@echo "Deploying ERP system..."
	@docker-compose up -d
	@echo "✅ Deployment complete!"

clean:
	@echo "Cleaning up..."
	@docker system prune -f
	@echo "✅ Cleanup complete!"
