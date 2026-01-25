#!/bin/bash
# cleanup-backups.sh - Cleanup old logs and backups

LOG_DIR="/var/log/erp"
BACKUP_DIR="/backups"

echo "Running maintenance cleanup at $(date)"

# Clean logs older than 30 days
find ${LOG_DIR} -name "*.log" -mtime +30 -exec rm {} \;

# Clean old Docker images (optional but recommended)
docker image prune -a --force --filter "until=168h"

echo "✓ Cleanup completed!"
