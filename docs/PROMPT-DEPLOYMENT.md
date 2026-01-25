# 🚀 PROMPT: DEPLOY ERP COSMETICS TO PRODUCTION

## CONTEXT

Hệ thống ERP mỹ phẩm đã hoàn thành development và testing. Cần deploy lên production server.

**Project Info:**
- 13 microservices (Go)
- Frontend: Vue 3 + TypeScript
- Database: PostgreSQL 16
- Cache: Redis 7
- Message Queue: NATS
- Object Storage: MinIO
- Monitoring: Prometheus + Grafana + Loki

**Server Requirements:**
- OS: Ubuntu 22.04 LTS
- CPU: 16 cores (minimum 8)
- RAM: 64GB (minimum 32GB)
- Disk: 500GB SSD
- Domain: erp.company.vn (thay bằng domain thực)

---

## PHASE 1: SERVER SETUP (30 phút)

### 1.1 SSH vào server và chạy script setup

```bash
# Upload scripts lên server
scp -r scripts/ root@YOUR_SERVER_IP:/opt/erp/

# SSH vào server
ssh root@YOUR_SERVER_IP

# Chạy server setup
cd /opt/erp
chmod +x scripts/*.sh
./scripts/server-setup.sh
```

### 1.2 Nội dung script server-setup.sh cần có:

```bash
#!/bin/bash
set -e

echo "=== ERP Server Setup ==="

# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
systemctl enable docker
systemctl start docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Install utilities
apt install -y \
  git \
  htop \
  ncdu \
  certbot \
  python3-certbot-nginx \
  nginx \
  fail2ban \
  ufw

# Configure firewall
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw --force enable

# Configure fail2ban
cat > /etc/fail2ban/jail.local << 'EOF'
[sshd]
enabled = true
maxretry = 3
bantime = 3600

[nginx-http-auth]
enabled = true
EOF

systemctl enable fail2ban
systemctl restart fail2ban

# Create app user
useradd -m -s /bin/bash erp || true
usermod -aG docker erp

# Create directories
mkdir -p /opt/erp/{config,logs,backups,ssl}
mkdir -p /data/{postgres,redis,minio,prometheus,grafana,loki}
chown -R erp:erp /opt/erp
chown -R erp:erp /data

echo "✅ Server setup completed!"
```

---

## PHASE 2: SSL CERTIFICATE (10 phút)

### 2.1 Setup SSL với Let's Encrypt

```bash
# Thay YOUR_DOMAIN và YOUR_EMAIL
export DOMAIN="erp.company.vn"
export EMAIL="admin@company.vn"

# Stop nginx temporarily
systemctl stop nginx || true

# Get certificate
certbot certonly --standalone \
  -d $DOMAIN \
  --non-interactive \
  --agree-tos \
  --email $EMAIL

# Setup auto-renewal
echo "0 0 1 * * certbot renew --quiet --post-hook 'systemctl reload nginx'" | crontab -

echo "✅ SSL certificate installed!"
```

---

## PHASE 3: CONFIGURATION (20 phút)

### 3.1 Tạo file .env cho production

```bash
cat > /opt/erp/.env << 'EOF'
# ===================
# ENVIRONMENT
# ===================
ENVIRONMENT=production
LOG_LEVEL=info

# ===================
# DOMAIN
# ===================
DOMAIN=erp.company.vn
CORS_ALLOWED_ORIGINS=https://erp.company.vn

# ===================
# DATABASE
# ===================
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=erp_admin
POSTGRES_PASSWORD=<GENERATE_STRONG_PASSWORD_32_CHARS>
POSTGRES_DB=erp_main

# Service databases
AUTH_DB_NAME=erp_auth
USER_DB_NAME=erp_user
MASTER_DATA_DB_NAME=erp_master_data
SUPPLIER_DB_NAME=erp_supplier
PROCUREMENT_DB_NAME=erp_procurement
WMS_DB_NAME=erp_wms
MANUFACTURING_DB_NAME=erp_manufacturing
SALES_DB_NAME=erp_sales
MARKETING_DB_NAME=erp_marketing
NOTIFICATION_DB_NAME=erp_notification
FILE_DB_NAME=erp_file
REPORTING_DB_NAME=erp_reporting

# ===================
# REDIS
# ===================
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=<GENERATE_STRONG_PASSWORD_32_CHARS>

# ===================
# NATS
# ===================
NATS_URL=nats://nats:4222

# ===================
# MINIO
# ===================
MINIO_ENDPOINT=minio:9000
MINIO_ROOT_USER=erp_minio_admin
MINIO_ROOT_PASSWORD=<GENERATE_STRONG_PASSWORD_32_CHARS>
MINIO_USE_SSL=false

# ===================
# JWT
# ===================
JWT_SECRET=<GENERATE_STRONG_SECRET_64_CHARS>
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h

# ===================
# BOM ENCRYPTION KEY (32 bytes for AES-256)
# ===================
BOM_ENCRYPTION_KEY=<GENERATE_32_BYTE_KEY>

# ===================
# SMTP (Email)
# ===================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=erp-notifications@company.vn
SMTP_PASSWORD=<APP_PASSWORD>
SMTP_FROM=ERP System <erp-notifications@company.vn>

# ===================
# SERVICE PORTS
# ===================
API_GATEWAY_PORT=8080
AUTH_SERVICE_PORT=8081
AUTH_GRPC_PORT=9081
USER_SERVICE_PORT=8082
USER_GRPC_PORT=9082
MASTER_DATA_SERVICE_PORT=8083
SUPPLIER_SERVICE_PORT=8084
PROCUREMENT_SERVICE_PORT=8085
WMS_SERVICE_PORT=8086
WMS_GRPC_PORT=9086
MANUFACTURING_SERVICE_PORT=8087
SALES_SERVICE_PORT=8088
MARKETING_SERVICE_PORT=8089
NOTIFICATION_SERVICE_PORT=8090
FILE_SERVICE_PORT=8091
REPORTING_SERVICE_PORT=8092

# ===================
# MONITORING
# ===================
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=<GENERATE_STRONG_PASSWORD>
PROMETHEUS_RETENTION_TIME=30d
EOF

# Generate passwords
echo ""
echo "🔐 Generate passwords với command sau:"
echo "openssl rand -base64 32"
```

### 3.2 Tạo file docker-compose.prod.yml

```yaml
version: '3.9'

networks:
  erp-network:
    driver: bridge

volumes:
  postgres-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/postgres
  redis-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/redis
  minio-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/minio
  prometheus-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/prometheus
  grafana-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/grafana
  loki-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /data/loki

services:
  # ==================
  # INFRASTRUCTURE
  # ==================
  postgres:
    image: postgres:16-alpine
    container_name: erp-postgres
    restart: always
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./scripts/init-databases.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - erp-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 8G

  redis:
    image: redis:7-alpine
    container_name: erp-redis
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD} --maxmemory 2gb --maxmemory-policy allkeys-lru
    volumes:
      - redis-data:/data
    networks:
      - erp-network
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G

  nats:
    image: nats:2.10-alpine
    container_name: erp-nats
    restart: always
    command: "-js -m 8222 --store_dir /data"
    volumes:
      - /data/nats:/data
    networks:
      - erp-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 3s
      retries: 5

  minio:
    image: minio/minio:latest
    container_name: erp-minio
    restart: always
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    volumes:
      - minio-data:/data
    networks:
      - erp-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 10s
      retries: 3

  # ==================
  # MONITORING
  # ==================
  prometheus:
    image: prom/prometheus:latest
    container_name: erp-prometheus
    restart: always
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=${PROMETHEUS_RETENTION_TIME}'
    volumes:
      - ./deploy/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./deploy/monitoring/rules:/etc/prometheus/rules
      - prometheus-data:/prometheus
    networks:
      - erp-network

  alertmanager:
    image: prom/alertmanager:latest
    container_name: erp-alertmanager
    restart: always
    volumes:
      - ./deploy/monitoring/alertmanager.yml:/etc/alertmanager/alertmanager.yml
    networks:
      - erp-network

  grafana:
    image: grafana/grafana:latest
    container_name: erp-grafana
    restart: always
    environment:
      - GF_SECURITY_ADMIN_USER=${GRAFANA_ADMIN_USER}
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_ROOT_URL=https://${DOMAIN}/grafana
    volumes:
      - grafana-data:/var/lib/grafana
      - ./deploy/grafana/provisioning:/etc/grafana/provisioning
    networks:
      - erp-network

  loki:
    image: grafana/loki:latest
    container_name: erp-loki
    restart: always
    command: -config.file=/etc/loki/loki-config.yml
    volumes:
      - ./deploy/loki/loki-config.yml:/etc/loki/loki-config.yml
      - loki-data:/loki
    networks:
      - erp-network

  # ==================
  # API GATEWAY
  # ==================
  api-gateway:
    image: erp/api-gateway:${VERSION:-latest}
    container_name: erp-api-gateway
    restart: always
    environment:
      - PORT=${API_GATEWAY_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - AUTH_GRPC_URL=auth-service:${AUTH_GRPC_PORT}
      - LOG_LEVEL=${LOG_LEVEL}
      - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS}
    networks:
      - erp-network
    depends_on:
      redis:
        condition: service_healthy
      auth-service:
        condition: service_healthy
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "5"

  # ==================
  # AUTH SERVICE
  # ==================
  auth-service:
    image: erp/auth-service:${VERSION:-latest}
    container_name: erp-auth-service
    restart: always
    environment:
      - PORT=${AUTH_SERVICE_PORT}
      - GRPC_PORT=${AUTH_GRPC_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${AUTH_DB_NAME}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - NATS_URL=${NATS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - JWT_ACCESS_TOKEN_EXPIRE=${JWT_ACCESS_TOKEN_EXPIRE}
      - JWT_REFRESH_TOKEN_EXPIRE=${JWT_REFRESH_TOKEN_EXPIRE}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      nats:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:${AUTH_SERVICE_PORT}/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "5"

  # ==================
  # USER SERVICE
  # ==================
  user-service:
    image: erp/user-service:${VERSION:-latest}
    container_name: erp-user-service
    restart: always
    environment:
      - PORT=${USER_SERVICE_PORT}
      - GRPC_PORT=${USER_GRPC_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${USER_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:${USER_SERVICE_PORT}/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

  # ==================
  # MASTER DATA SERVICE
  # ==================
  master-data-service:
    image: erp/master-data-service:${VERSION:-latest}
    container_name: erp-master-data-service
    restart: always
    environment:
      - PORT=${MASTER_DATA_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${MASTER_DATA_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:${MASTER_DATA_SERVICE_PORT}/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # ==================
  # SUPPLIER SERVICE
  # ==================
  supplier-service:
    image: erp/supplier-service:${VERSION:-latest}
    container_name: erp-supplier-service
    restart: always
    environment:
      - PORT=${SUPPLIER_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${SUPPLIER_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # PROCUREMENT SERVICE
  # ==================
  procurement-service:
    image: erp/procurement-service:${VERSION:-latest}
    container_name: erp-procurement-service
    restart: always
    environment:
      - PORT=${PROCUREMENT_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${PROCUREMENT_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # WMS SERVICE (CRITICAL)
  # ==================
  wms-service:
    image: erp/wms-service:${VERSION:-latest}
    container_name: erp-wms-service
    restart: always
    environment:
      - PORT=${WMS_SERVICE_PORT}
      - GRPC_PORT=${WMS_GRPC_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${WMS_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:${WMS_SERVICE_PORT}/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G

  # ==================
  # MANUFACTURING SERVICE (CRITICAL)
  # ==================
  manufacturing-service:
    image: erp/manufacturing-service:${VERSION:-latest}
    container_name: erp-manufacturing-service
    restart: always
    environment:
      - PORT=${MANUFACTURING_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${MANUFACTURING_DB_NAME}
      - NATS_URL=${NATS_URL}
      - BOM_ENCRYPTION_KEY=${BOM_ENCRYPTION_KEY}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      wms-service:
        condition: service_healthy
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G

  # ==================
  # SALES SERVICE
  # ==================
  sales-service:
    image: erp/sales-service:${VERSION:-latest}
    container_name: erp-sales-service
    restart: always
    environment:
      - PORT=${SALES_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${SALES_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # MARKETING SERVICE
  # ==================
  marketing-service:
    image: erp/marketing-service:${VERSION:-latest}
    container_name: erp-marketing-service
    restart: always
    environment:
      - PORT=${MARKETING_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${MARKETING_DB_NAME}
      - NATS_URL=${NATS_URL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # NOTIFICATION SERVICE
  # ==================
  notification-service:
    image: erp/notification-service:${VERSION:-latest}
    container_name: erp-notification-service
    restart: always
    environment:
      - PORT=${NOTIFICATION_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${NOTIFICATION_DB_NAME}
      - NATS_URL=${NATS_URL}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASSWORD=${SMTP_PASSWORD}
      - SMTP_FROM=${SMTP_FROM}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_healthy

  # ==================
  # FILE SERVICE
  # ==================
  file-service:
    image: erp/file-service:${VERSION:-latest}
    container_name: erp-file-service
    restart: always
    environment:
      - PORT=${FILE_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${FILE_DB_NAME}
      - MINIO_ENDPOINT=${MINIO_ENDPOINT}
      - MINIO_ACCESS_KEY_ID=${MINIO_ROOT_USER}
      - MINIO_SECRET_ACCESS_KEY=${MINIO_ROOT_PASSWORD}
      - MINIO_USE_SSL=${MINIO_USE_SSL}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_healthy

  # ==================
  # REPORTING SERVICE
  # ==================
  reporting-service:
    image: erp/reporting-service:${VERSION:-latest}
    container_name: erp-reporting-service
    restart: always
    environment:
      - PORT=${REPORTING_SERVICE_PORT}
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=${POSTGRES_HOST}
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${REPORTING_DB_NAME}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # FRONTEND
  # ==================
  frontend:
    image: erp/frontend:${VERSION:-latest}
    container_name: erp-frontend
    restart: always
    networks:
      - erp-network
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 256M

  # ==================
  # NGINX (Reverse Proxy)
  # ==================
  nginx:
    image: nginx:alpine
    container_name: erp-nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./deploy/nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./deploy/nginx/erp.conf:/etc/nginx/conf.d/erp.conf
      - /etc/letsencrypt:/etc/letsencrypt:ro
    networks:
      - erp-network
    depends_on:
      - api-gateway
      - frontend
      - grafana
```

---

## PHASE 4: BUILD & PUSH IMAGES (30 phút)

### 4.1 Build tất cả services

```bash
# Trên máy development
cd /path/to/ERP-main

# Set version
export VERSION=$(date +%Y%m%d)-$(git rev-parse --short HEAD)
echo "Building version: $VERSION"

# Build all services
for svc in api-gateway auth-service user-service master-data-service \
           supplier-service procurement-service wms-service \
           manufacturing-service sales-service marketing-service \
           notification-service file-service reporting-service; do
    echo "Building $svc..."
    docker build -t erp/$svc:$VERSION -t erp/$svc:latest ./services/$svc
done

# Build frontend
echo "Building frontend..."
cd frontend
npm ci
npm run build
docker build -t erp/frontend:$VERSION -t erp/frontend:latest .
cd ..

echo "✅ All images built!"
```

### 4.2 Push to Registry (nếu dùng private registry)

```bash
# Login to registry
docker login registry.company.vn

# Tag and push
for svc in api-gateway auth-service user-service master-data-service \
           supplier-service procurement-service wms-service \
           manufacturing-service sales-service marketing-service \
           notification-service file-service reporting-service frontend; do
    docker tag erp/$svc:$VERSION registry.company.vn/erp/$svc:$VERSION
    docker tag erp/$svc:latest registry.company.vn/erp/$svc:latest
    docker push registry.company.vn/erp/$svc:$VERSION
    docker push registry.company.vn/erp/$svc:latest
done
```

### 4.3 Hoặc save images và copy lên server

```bash
# Save all images
docker save erp/api-gateway erp/auth-service erp/user-service \
    erp/master-data-service erp/supplier-service erp/procurement-service \
    erp/wms-service erp/manufacturing-service erp/sales-service \
    erp/marketing-service erp/notification-service erp/file-service \
    erp/reporting-service erp/frontend | gzip > erp-images.tar.gz

# Copy to server
scp erp-images.tar.gz root@YOUR_SERVER_IP:/opt/erp/

# On server: load images
ssh root@YOUR_SERVER_IP
cd /opt/erp
gunzip -c erp-images.tar.gz | docker load
```

---

## PHASE 5: DATABASE MIGRATION (15 phút)

### 5.1 Start database first

```bash
cd /opt/erp
docker-compose -f docker-compose.prod.yml up -d postgres
sleep 10

# Check postgres is ready
docker-compose -f docker-compose.prod.yml exec postgres pg_isready
```

### 5.2 Run migrations cho từng service

```bash
# Auth service migrations
docker-compose -f docker-compose.prod.yml run --rm auth-service \
    /app/migrate -path /app/migrations -database \
    "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${AUTH_DB_NAME}?sslmode=disable" up

# User service migrations
docker-compose -f docker-compose.prod.yml run --rm user-service \
    /app/migrate -path /app/migrations -database \
    "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${USER_DB_NAME}?sslmode=disable" up

# Repeat for all services...
# master-data-service, supplier-service, procurement-service,
# wms-service, manufacturing-service, sales-service,
# marketing-service, notification-service, file-service, reporting-service
```

---

## PHASE 6: DEPLOY (15 phút)

### 6.1 Deploy tất cả services

```bash
cd /opt/erp

# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Watch logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 6.2 Health check

```bash
#!/bin/bash
# health-check.sh

SERVICES=(
    "api-gateway:8080"
    "auth-service:8081"
    "user-service:8082"
    "master-data-service:8083"
    "supplier-service:8084"
    "procurement-service:8085"
    "wms-service:8086"
    "manufacturing-service:8087"
    "sales-service:8088"
    "marketing-service:8089"
    "notification-service:8090"
    "file-service:8091"
    "reporting-service:8092"
)

echo "=== Health Check ==="
for svc in "${SERVICES[@]}"; do
    name="${svc%%:*}"
    port="${svc##*:}"
    status=$(docker exec erp-$name wget -q -O- http://localhost:$port/health 2>/dev/null || echo "FAILED")
    if [[ "$status" == *"ok"* ]] || [[ "$status" == *"healthy"* ]]; then
        echo "✅ $name: OK"
    else
        echo "❌ $name: FAILED"
    fi
done
```

---

## PHASE 7: VERIFY (15 phút)

### 7.1 Test endpoints

```bash
# Test login
curl -X POST https://erp.company.vn/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@company.vn","password":"Admin@123"}'

# Test với token
TOKEN="<token_from_login>"
curl https://erp.company.vn/api/v1/users \
    -H "Authorization: Bearer $TOKEN"
```

### 7.2 Check monitoring

```bash
# Grafana: https://erp.company.vn/grafana
# Default: admin / <GRAFANA_ADMIN_PASSWORD>

# Prometheus: Internal only
# http://localhost:9090 (qua SSH tunnel)
```

---

## PHASE 8: BACKUP SETUP (10 phút)

### 8.1 Setup cron jobs

```bash
# Edit crontab
crontab -e

# Add backup jobs
# Database backup - daily at 2 AM
0 2 * * * /opt/erp/scripts/backup-db.sh >> /var/log/erp/backup.log 2>&1

# Cleanup old backups - weekly on Sunday
0 4 * * 0 /opt/erp/scripts/cleanup-backups.sh >> /var/log/erp/backup.log 2>&1
```

### 8.2 Test backup

```bash
# Run manual backup
/opt/erp/scripts/backup-db.sh

# Verify backup
ls -la /opt/erp/backups/
```

---

## 📋 POST-DEPLOYMENT CHECKLIST

```markdown
## Infrastructure
- [ ] All services running (docker ps)
- [ ] Health checks passing
- [ ] SSL certificate valid
- [ ] Firewall configured

## Application
- [ ] Login works
- [ ] API responses correct
- [ ] Frontend loads
- [ ] File upload works

## Monitoring
- [ ] Grafana accessible
- [ ] Dashboards showing data
- [ ] Alerts configured
- [ ] Log aggregation working

## Security
- [ ] Strong passwords set
- [ ] JWT secret rotated
- [ ] BOM encryption key secure
- [ ] fail2ban active

## Backup
- [ ] Backup script works
- [ ] Backup files created
- [ ] Restore tested

## Documentation
- [ ] Access credentials documented (secure location)
- [ ] Runbook available
- [ ] Team notified
```

---

## 🚨 ROLLBACK PROCEDURE

```bash
# Nếu có vấn đề, rollback về version trước:

# 1. Stop current version
docker-compose -f docker-compose.prod.yml down

# 2. Restore database backup
gunzip -c /opt/erp/backups/db_backup_YYYYMMDD.sql.gz | \
    docker exec -i erp-postgres psql -U ${POSTGRES_USER}

# 3. Deploy previous version
export VERSION=<previous_version>
docker-compose -f docker-compose.prod.yml up -d

# 4. Verify
./scripts/health-check.sh
```

---

## 📞 SUPPORT CONTACTS

| Role | Name | Contact |
|------|------|---------|
| Tech Lead | [Name] | [Phone/Email] |
| DBA | [Name] | [Phone/Email] |
| DevOps | [Name] | [Phone/Email] |

---

**Created**: 2026-01-25  
**Version**: 1.0
