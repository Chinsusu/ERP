#!/bin/bash

# Backup script for ERP Cosmetics
# Usage: ./scripts/backup.sh

set -e

BACKUP_DIR="/backups/erp-$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "========================================"
echo "ERP Cosmetics - Backup Starting"
echo "Backup directory: $BACKUP_DIR"
echo "========================================"
echo ""

# Backup PostgreSQL databases
echo "Backing up PostgreSQL databases..."
docker exec erp-postgres pg_dumpall -U postgres | gzip > "$BACKUP_DIR/postgres-all.sql.gz"
echo "✓ PostgreSQL backup complete"

# Backup individual databases
for DB in auth_db user_db master_data_db supplier_db procurement_db wms_db manufacturing_db sales_db marketing_db finance_db reporting_db notification_db ai_db; do
    echo "  - Backing up $DB..."
    docker exec erp-postgres pg_dump -U postgres "$DB" | gzip > "$BACKUP_DIR/${DB}.sql.gz"
done

# Backup Redis
echo ""
echo "Backing up Redis..."
docker exec erp-redis redis-cli -a ${REDIS_PASSWORD:-redis} --no-auth-warning --rdb /data/dump.rdb SAVE
docker cp erp-redis:/data/dump.rdb "$BACKUP_DIR/redis-dump.rdb"
echo "✓ Redis backup complete"

# Backup MinIO data
echo ""
echo "Backing up MinIO data..."
mkdir -p "$BACKUP_DIR/minio"
docker exec erp-minio mc mirror /data "$BACKUP_DIR/minio" 2>/dev/null || echo "MinIO backup requires mc client configuration"

# Backup configuration files
echo ""
echo "Backing up configuration files..."
cp .env "$BACKUP_DIR/.env.backup" 2>/dev/null || echo ".env not found"
cp docker-compose.yml "$BACKUP_DIR/" 2>/dev/null || echo "docker-compose.yml not found"
cp -r deploy "$BACKUP_DIR/" 2>/dev/null || echo "deploy directory not found"

# Create archive
echo ""
echo "Creating compressed archive..."
tar -czf "${BACKUP_DIR}.tar.gz" -C "$(dirname "$BACKUP_DIR")" "$(basename "$BACKUP_DIR")"
rm -rf "$BACKUP_DIR"

echo ""
echo "========================================"
echo "Backup complete!"
echo "Archive: ${BACKUP_DIR}.tar.gz"
echo "========================================"

# Cleanup old backups (keep last 30 days)
echo ""
echo "Cleaning up old backups (keeping last 30 days)..."
find /backups -name "erp-*.tar.gz" -mtime +30 -delete 2>/dev/null || true
echo "✓ Cleanup complete"
