# 16 - DEPLOYMENT GUIDE

## TỔNG QUAN

Hướng dẫn deploy hệ thống ERP sử dụng **Docker Compose** cho môi trường on-premise.

---

## SYSTEM REQUIREMENTS

### Hardware

**Minimum**:
- CPU: 4 cores
- RAM: 16GB
- Disk: 500GB SSD

**Recommended (production ~100 users)**:
- CPU: 8-16 cores
- RAM: 32-64GB
- Disk: 1TB SSD (RAID 10 preferred)

### Software

- **OS**: Ubuntu 22.04 LTS / CentOS 8 / RHEL 8
- **Docker**: 24.0+
- **Docker Compose**: 2.20+
- **Git**: 2.x

---

## QUICK START

### 1. Clone Repository

```bash
git clone git@github.com:Chinsusu/ERP.git
cd ERP
```

### 2. Configure Environment

```bash
cp .env.example .env
nano .env  # Edit configuration
```

### 3. Deploy

```bash
docker-compose up -d
```

### 4. Initialize Database

```bash
./scripts/init-databases.sh
```

### 5. Access System

- **Frontend**: http://localhost (or your domain)
- **API Gateway**: http://localhost:8080

---

## DOCKER COMPOSE ARCHITECTURE

```yaml
version: '3.9'

services:
  # Infrastructure
  nginx:
  postgres:
  redis:
  nats:
  minio:
  
  # Monitoring
  prometheus:
  grafana:
  loki:
  jaeger:
  
  # Services (15 microservices)
  api-gateway:
  auth-service:
  user-service:
  master-data-service:
  supplier-service:
  procurement-service:
  wms-service:
  manufacturing-service:
  sales-service:
  marketing-service:
  finance-service:
  reporting-service:
  notification-service:
  ai-service:
  file-service:
  
  # Frontend
  frontend:
```

---

## ENVIRONMENT VARIABLES

### .env File

```bash
# Environment
ENVIRONMENT=production
DOMAIN=erp.company.com

# Security
JWT_SECRET=<generate-64-char-random-string>
BOM_ENCRYPTION_KEY=<generate-32-byte-hex>

# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=<strong-password>

# Redis
REDIS_PASSWORD=<strong-password>

# MinIO
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=<strong-password>

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=noreply@company.com
SMTP_PASSWORD=<app-password>

# Monitoring
GRAFANA_ADMIN_PASSWORD=<strong-password>
```

---

## FULL DOCKER COMPOSE

```yaml
version: '3.9'

networks:
  erp-network:
    driver: bridge

volumes:
  postgres-data:
  redis-data:
  minio-data:
  grafana-data:
  prometheus-data:

services:
  # ====================
  # Nginx Reverse Proxy
  # ====================
  nginx:
    image: nginx:alpine
    container_name: erp-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - api-gateway
      - frontend
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # PostgreSQL
  # ====================
  postgres:
    image: postgres:16-alpine
    container_name: erp-postgres
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - erp-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ====================
  # Redis
  # ====================
  redis:
    image: redis:7-alpine
    container_name: erp-redis
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # NATS
  # ====================
  nats:
    image: nats:latest
    container_name: erp-nats
    command: "-js -m 8222"
    ports:
      - "4222:4222"
      - "8222:8222"
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # MinIO
  # ====================
  minio:
    image: minio/minio:latest
    container_name: erp-minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    volumes:
      - minio-data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # API Gateway
  # ====================
  api-gateway:
    build: ./services/api-gateway
    container_name: erp-api-gateway
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - REDIS_HOST=redis
      - AUTH_SERVICE_URL=http://auth-service:8081
    depends_on:
      - redis
      - auth-service
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # Auth Service
  # ====================
  auth-service:
    build: ./services/auth-service
    environment:
      - AUTH_SERVICE_PORT=8081
      - AUTH_DB_HOST=postgres
      - AUTH_DB_NAME=auth_db
      - REDIS_HOST=redis
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # (Other 13 services...)
  # ====================
  
  # ====================
  # Frontend
  # ====================
  frontend:
    build: ./frontend
    container_name: erp-frontend
    environment:
      - VITE_API_URL=http://localhost:8080
    ports:
      - "3000:80"
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # Monitoring: Prometheus
  # ====================
  prometheus:
    image: prom/prometheus:latest
    container_name: erp-prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    ports:
      - "9090:9090"
    networks:
      - erp-network
    restart: unless-stopped

  # ====================
  # Monitoring: Grafana
  # ====================
  grafana:
    image: grafana/grafana:latest
    container_name: erp-grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards
    ports:
      - "3001:3000"
    networks:
      - erp-network
    restart: unless-stopped
```

---

## DATABASE INITIALIZATION

```sql
-- scripts/init-db.sql
CREATE DATABASE auth_db;
CREATE DATABASE user_db;
CREATE DATABASE master_data_db;
CREATE DATABASE supplier_db;
CREATE DATABASE procurement_db;
CREATE DATABASE wms_db;
CREATE DATABASE manufacturing_db;
CREATE DATABASE sales_db;
CREATE DATABASE marketing_db;
CREATE DATABASE finance_db;
CREATE DATABASE reporting_db;
CREATE DATABASE notification_db;
CREATE DATABASE ai_db;

-- Create users
CREATE USER auth_service WITH PASSWORD 'auth_password';
GRANT ALL PRIVILEGES ON DATABASE auth_db TO auth_service;

-- Repeat for other services...
```

---

## NGINX CONFIGURATION

```nginx
upstream api_gateway {
    server api-gateway:8080;
}

server {
    listen 80;
    server_name erp.company.com;
    
    # Frontend
    location / {
        proxy_pass http://frontend:80;
    }
    
    # API
    location /api/ {
        proxy_pass http://api_gateway/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## MONITORING DASHBOARDS

Access:
- **Grafana**: http://localhost:3001 (admin / password from .env)
- **Prometheus**: http://localhost:9090

---

## BACKUP SCRIPT

```bash
#!/bin/bash
# scripts/backup.sh

BACKUP_DIR="/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR

# Backup databases
docker exec erp-postgres pg_dumpall -U postgres | gzip > $BACKUP_DIR/postgres.sql.gz

# Backup MinIO
docker exec erp-minio mc mirror /data $BACKUP_DIR/minio

# Backup configs
cp -r .env docker-compose.yml nginx/ $BACKUP_DIR/

echo "Backup completed: $BACKUP_DIR"
```

---

## SCALING

### Horizontal Scaling (Multiple Replicas)

```yaml
services:
  api-gateway:
    deploy:
      replicas: 3
```

### Load Balancing

Nginx handles load balancing:

```nginx
upstream api_gateway {
    server api-gateway-1:8080;
    server api-gateway-2:8080;
    server api-gateway-3:8080;
}
```

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
