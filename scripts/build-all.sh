#!/bin/bash
set -e

echo "🚀 Building all ERP services..."
echo "================================"

# Version tag
VERSION=$(date +%Y%m%d)-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
echo "Version: $VERSION"
echo ""

# List of all services
SERVICES=(
    "api-gateway"
    "auth-service"
    "user-service"
    "master-data-service"
    "supplier-service"
    "procurement-service"
    "wms-service"
    "manufacturing-service"
    "sales-service"
    "marketing-service"
    "notification-service"
    "file-service"
    "reporting-service"
)

# Build each service
for svc in "${SERVICES[@]}"; do
    echo "📦 Building $svc..."
    docker build \
        -f Dockerfile.service \
        --build-arg SERVICE=$svc \
        -t erp/$svc:$VERSION \
        -t erp/$svc:latest \
        . || {
            echo "❌ Failed to build $svc"
            exit 1
        }
    echo "✅ $svc built successfully"
    echo ""
done

echo "🎨 Building frontend..."
cd frontend
docker build -t erp/frontend:$VERSION -t erp/frontend:latest . || {
    echo "❌ Failed to build frontend"
    exit 1
}
cd ..
echo "✅ Frontend built successfully"
echo ""

echo "🎉 All images built successfully!"
echo ""
echo "Built images:"
docker images | grep "erp/" | grep -E "latest|$VERSION"
