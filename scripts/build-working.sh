#!/bin/bash
set -e

echo "🚀 Building working ERP services (skip broken ones)..."
echo "================================"

VERSION=$(date +%Y%m%d)-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
echo "Version: $VERSION"
echo ""

# List of services to build (excluding auth-service for now)
SERVICES=(
    "api-gateway"
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

SUCCESS=()
FAILED=()

# Build each service
for svc in "${SERVICES[@]}"; do
    echo "📦 Building $svc..."
    if docker build \
        -f Dockerfile.service \
        --build-arg SERVICE=$svc \
        -t erp/$svc:$VERSION \
        -t erp/$svc:latest \
        . 2>&1 | tee /tmp/build-$svc.log; then
        echo "✅ $svc built successfully"
        SUCCESS+=("$svc")
    else
        echo "❌ Failed to build $svc"
        FAILED+=("$svc")
    fi
    echo ""
done

echo "🎨 Building frontend..."
if cd frontend && docker build -t erp/frontend:$VERSION -t erp/frontend:latest .; then
    cd ..
    echo "✅ Frontend built successfully"
    SUCCESS+=("frontend")
else
    cd ..
    echo "❌ Failed to build frontend"
    FAILED+=("frontend")
fi
echo ""

echo "================================"
echo "📊 BUILD SUMMARY"
echo "================================"
echo "✅ Success (${#SUCCESS[@]}): ${SUCCESS[*]}"
echo "❌ Failed (${#FAILED[@]}): ${FAILED[*]}"
echo ""

if [ ${#FAILED[@]} -eq 0 ]; then
    echo "🎉 All services built successfully!"
    exit 0
else
    echo "⚠️  Some services failed to build"
    exit 1
fi
