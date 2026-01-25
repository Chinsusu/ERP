#!/bin/bash
# backup-db.sh - Automated PostgreSQL backup

BACKUP_DIR="/backups/db"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/erp_prod_${TIMESTAMP}.sql.gz"

mkdir -p ${BACKUP_DIR}

echo "Starting database backup at $(date)"

# Perform backup using pg_dump
# Assuming running inside docker or having access to DB
docker exec erp-postgres pg_dump -U user -d erp_prod | gzip > ${BACKUP_FILE}

if [ $? -eq 0 ]; then
    echo "✓ Backup successful: ${BACKUP_FILE}"
    # Rotate backups (keep last 30 days)
    find ${BACKUP_DIR} -name "erp_prod_*.sql.gz" -mtime +30 -delete
else
    echo "✗ Backup failed!"
    exit 1
fi
