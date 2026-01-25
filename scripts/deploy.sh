#!/bin/bash
# deploy.sh - Rolling deployment script for ERP Cosmetics
# Usage: ./deploy.sh <version>

set -e

VERSION=${1:-latest}
COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="/backups/pre-deploy-$(date +%Y%m%d_%H%M%S)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $1"; }
warn() { echo -e "${YELLOW}[$(date +'%H:%M:%S')] WARNING:${NC} $1"; }
error() { echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"; exit 1; }

# Check version argument
if [ -z "$1" ]; then
    warn "No version specified, using 'latest'"
fi

log "Starting deployment of version: $VERSION"

# Pre-deploy backup
log "Creating pre-deploy backup..."
mkdir -p "$BACKUP_DIR"
./scripts/backup-db.sh "$BACKUP_DIR" || warn "Backup failed, continuing anyway"

# Pull new images
log "Pulling new images..."
export VERSION=$VERSION
docker-compose -f $COMPOSE_FILE pull

# Deploy infrastructure first (if needed)
log "Ensuring infrastructure is running..."
docker-compose -f $COMPOSE_FILE up -d postgres redis nats minio
sleep 5

# Rolling update services
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

for service in "${SERVICES[@]}"; do
    log "Deploying $service..."
    docker-compose -f $COMPOSE_FILE up -d --no-deps "$service"
    sleep 3
    
    # Quick health check
    if docker-compose -f $COMPOSE_FILE ps "$service" | grep -q "Up"; then
        log "$service is running"
    else
        error "$service failed to start! Check logs: docker-compose logs $service"
    fi
done

# Deploy frontend
log "Deploying frontend..."
docker-compose -f $COMPOSE_FILE up -d --no-deps frontend
sleep 5

# Full health check
log "Running health check..."
./scripts/health-check.sh || warn "Some services may not be fully healthy"

# Summary
log "========================================="
log "Deployment completed!"
log "Version: $VERSION"
log "Backup: $BACKUP_DIR"
log "========================================="
log "Next steps:"
log "  1. Verify: curl http://localhost:8080/health"
log "  2. Test login"
log "  3. Monitor logs: docker-compose logs -f"
