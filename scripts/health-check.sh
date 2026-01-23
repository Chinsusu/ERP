#!/bin/bash

# Health check script for ERP services
# Usage: ./scripts/health-check.sh

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "ERP Cosmetics - Health Check"
echo "========================================"
echo ""

# Check PostgreSQL
echo -n "PostgreSQL... "
if docker exec erp-postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${RED}✗ Unhealthy${NC}"
fi

# Check Redis
echo -n "Redis... "
if docker exec erp-redis redis-cli -a ${REDIS_PASSWORD:-redis} --no-auth-warning ping > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${RED}✗ Unhealthy${NC}"
fi

# Check NATS
echo -n "NATS... "
if curl -f http://localhost:8222/healthz > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${RED}✗ Unhealthy${NC}"
fi

# Check MinIO
echo -n "MinIO... "
if curl -f http://localhost:9000/minio/health/live > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${RED}✗ Unhealthy${NC}"
fi

# Check API Gateway
echo -n "API Gateway... "
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${YELLOW}⚠ Not running or unhealthy${NC}"
fi

# Check Prometheus
echo -n "Prometheus... "
if curl -f http://localhost:9090/-/healthy > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${YELLOW}⚠ Not running or unhealthy${NC}"
fi

# Check Grafana
echo -n "Grafana... "
if curl -f http://localhost:3001/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${YELLOW}⚠ Not running or unhealthy${NC}"
fi

echo ""
echo "========================================"
echo "Health check complete"
echo "========================================"
