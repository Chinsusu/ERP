# 🚀 DEPLOY ERP VỚI CLOUDFLARE TUNNEL

## Ưu điểm so với cách truyền thống

| Traditional | Cloudflare Tunnel |
|-------------|-------------------|
| Cần public IP | ❌ Không cần |
| Mở port 80/443 | ❌ Không cần |
| Tự quản lý SSL | ✅ Cloudflare lo |
| DDoS protection | ✅ Miễn phí |
| Lộ IP server | ✅ Ẩn hoàn toàn |

---

## PREREQUISITES

1. **Cloudflare Account** (Free plan OK)
2. **Domain** đã add vào Cloudflare (DNS managed by Cloudflare)
3. **Server** (có thể ở nhà, VPS, hoặc behind NAT)

---

## PHASE 1: SETUP CLOUDFLARE TUNNEL (15 phút)

### 1.1 Cài đặt cloudflared trên server

```bash
# Ubuntu/Debian
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb -o cloudflared.deb
sudo dpkg -i cloudflared.deb

# Verify
cloudflared --version
```

### 1.2 Login Cloudflare

```bash
cloudflared tunnel login
# Browser sẽ mở ra, chọn domain của bạn
# File cert.pem sẽ được tạo tại ~/.cloudflared/
```

### 1.3 Tạo Tunnel

```bash
# Tạo tunnel tên "erp-tunnel"
cloudflared tunnel create erp-tunnel

# Output sẽ có dạng:
# Tunnel credentials written to /root/.cloudflared/<TUNNEL_ID>.json
# Created tunnel erp-tunnel with id <TUNNEL_ID>

# Lưu TUNNEL_ID
export TUNNEL_ID=<your-tunnel-id>
echo "TUNNEL_ID=$TUNNEL_ID" >> ~/.bashrc
```

### 1.4 Tạo DNS Record

```bash
# Tạo CNAME record trỏ domain về tunnel
cloudflared tunnel route dns erp-tunnel erp.company.vn

# Nếu cần subdomain cho Grafana
cloudflared tunnel route dns erp-tunnel grafana.erp.company.vn
```

---

## PHASE 2: CONFIGURE TUNNEL (10 phút)

### 2.1 Tạo config file

```bash
mkdir -p /opt/erp/cloudflared

cat > /opt/erp/cloudflared/config.yml << 'EOF'
tunnel: <TUNNEL_ID>
credentials-file: /etc/cloudflared/<TUNNEL_ID>.json

# Ingress rules - route traffic to services
ingress:
  # Main ERP application
  - hostname: erp.company.vn
    service: http://localhost:80
    originRequest:
      noTLSVerify: true
  
  # Grafana monitoring (optional subdomain)
  - hostname: grafana.erp.company.vn
    service: http://localhost:3001
    originRequest:
      noTLSVerify: true
  
  # Catch-all (required)
  - service: http_status:404
EOF

# Thay <TUNNEL_ID> bằng ID thực
sed -i "s/<TUNNEL_ID>/$TUNNEL_ID/g" /opt/erp/cloudflared/config.yml
```

### 2.2 Copy credentials

```bash
# Copy tunnel credentials
sudo mkdir -p /etc/cloudflared
sudo cp ~/.cloudflared/${TUNNEL_ID}.json /etc/cloudflared/
sudo cp ~/.cloudflared/cert.pem /etc/cloudflared/
```

---

## PHASE 3: DOCKER COMPOSE CHO CLOUDFLARE TUNNEL

### 3.1 docker-compose.prod.yml (simplified - no nginx SSL)

```yaml
version: '3.9'

networks:
  erp-network:
    driver: bridge

volumes:
  postgres-data:
  redis-data:
  minio-data:
  prometheus-data:
  grafana-data:
  loki-data:

services:
  # ==================
  # CLOUDFLARE TUNNEL
  # ==================
  cloudflared:
    image: cloudflare/cloudflared:latest
    container_name: erp-cloudflared
    restart: always
    command: tunnel --config /etc/cloudflared/config.yml run
    volumes:
      - /opt/erp/cloudflared/config.yml:/etc/cloudflared/config.yml:ro
      - /etc/cloudflared:/etc/cloudflared:ro
    networks:
      - erp-network
    depends_on:
      - nginx

  # ==================
  # NGINX (Internal reverse proxy)
  # ==================
  nginx:
    image: nginx:alpine
    container_name: erp-nginx
    restart: always
    # KHÔNG expose ports ra ngoài!
    # ports:
    #   - "80:80"
    volumes:
      - ./deploy/nginx/nginx-cloudflare.conf:/etc/nginx/nginx.conf:ro
    networks:
      - erp-network
    depends_on:
      - api-gateway
      - frontend

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

  redis:
    image: redis:7-alpine
    container_name: erp-redis
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    networks:
      - erp-network
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  nats:
    image: nats:2.10-alpine
    container_name: erp-nats
    restart: always
    command: "-js -m 8222"
    networks:
      - erp-network

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

  # ==================
  # MONITORING
  # ==================
  prometheus:
    image: prom/prometheus:latest
    container_name: erp-prometheus
    restart: always
    volumes:
      - ./deploy/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./deploy/monitoring/rules:/etc/prometheus/rules
      - prometheus-data:/prometheus
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
      - GF_SERVER_ROOT_URL=https://grafana.erp.company.vn
    volumes:
      - grafana-data:/var/lib/grafana
      - ./deploy/grafana/provisioning:/etc/grafana/provisioning
    ports:
      - "127.0.0.1:3001:3000"  # Only localhost, cloudflared will proxy
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
      - PORT=8080
      - ENVIRONMENT=${ENVIRONMENT}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - AUTH_GRPC_URL=auth-service:9081
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      - redis
      - auth-service

  # ==================
  # AUTH SERVICE
  # ==================
  auth-service:
    image: erp/auth-service:${VERSION:-latest}
    container_name: erp-auth-service
    restart: always
    environment:
      - PORT=8081
      - GRPC_PORT=9081
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${AUTH_DB_NAME}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - NATS_URL=nats://nats:4222
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

  # ==================
  # USER SERVICE
  # ==================
  user-service:
    image: erp/user-service:${VERSION:-latest}
    container_name: erp-user-service
    restart: always
    environment:
      - PORT=8082
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${USER_DB_NAME}
      - NATS_URL=nats://nats:4222
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # MASTER DATA SERVICE
  # ==================
  master-data-service:
    image: erp/master-data-service:${VERSION:-latest}
    container_name: erp-master-data-service
    restart: always
    environment:
      - PORT=8083
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${MASTER_DATA_DB_NAME}
      - NATS_URL=nats://nats:4222
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # SUPPLIER SERVICE
  # ==================
  supplier-service:
    image: erp/supplier-service:${VERSION:-latest}
    container_name: erp-supplier-service
    restart: always
    environment:
      - PORT=8084
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${SUPPLIER_DB_NAME}
      - NATS_URL=nats://nats:4222
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
      - PORT=8085
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${PROCUREMENT_DB_NAME}
      - NATS_URL=nats://nats:4222
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # WMS SERVICE
  # ==================
  wms-service:
    image: erp/wms-service:${VERSION:-latest}
    container_name: erp-wms-service
    restart: always
    environment:
      - PORT=8086
      - GRPC_PORT=9086
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${WMS_DB_NAME}
      - NATS_URL=nats://nats:4222
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # MANUFACTURING SERVICE
  # ==================
  manufacturing-service:
    image: erp/manufacturing-service:${VERSION:-latest}
    container_name: erp-manufacturing-service
    restart: always
    environment:
      - PORT=8087
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${MANUFACTURING_DB_NAME}
      - NATS_URL=nats://nats:4222
      - BOM_ENCRYPTION_KEY=${BOM_ENCRYPTION_KEY}
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy

  # ==================
  # SALES SERVICE
  # ==================
  sales-service:
    image: erp/sales-service:${VERSION:-latest}
    container_name: erp-sales-service
    restart: always
    environment:
      - PORT=8088
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${SALES_DB_NAME}
      - NATS_URL=nats://nats:4222
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
      - PORT=8089
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${MARKETING_DB_NAME}
      - NATS_URL=nats://nats:4222
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
      - PORT=8090
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${NOTIFICATION_DB_NAME}
      - NATS_URL=nats://nats:4222
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

  # ==================
  # FILE SERVICE
  # ==================
  file-service:
    image: erp/file-service:${VERSION:-latest}
    container_name: erp-file-service
    restart: always
    environment:
      - PORT=8091
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${FILE_DB_NAME}
      - MINIO_ENDPOINT=minio:9000
      - MINIO_ACCESS_KEY_ID=${MINIO_ROOT_USER}
      - MINIO_SECRET_ACCESS_KEY=${MINIO_ROOT_PASSWORD}
      - MINIO_USE_SSL=false
      - LOG_LEVEL=${LOG_LEVEL}
    networks:
      - erp-network
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_started

  # ==================
  # REPORTING SERVICE
  # ==================
  reporting-service:
    image: erp/reporting-service:${VERSION:-latest}
    container_name: erp-reporting-service
    restart: always
    environment:
      - PORT=8092
      - ENVIRONMENT=${ENVIRONMENT}
      - DB_HOST=postgres
      - DB_PORT=5432
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
```

### 3.2 Nginx config cho Cloudflare (không cần SSL)

```bash
cat > /opt/erp/deploy/nginx/nginx-cloudflare.conf << 'EOF'
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log;

    # Gzip
    gzip on;
    gzip_types text/plain application/json application/javascript text/css;

    # Upstream services
    upstream api_gateway {
        server api-gateway:8080;
        keepalive 32;
    }

    upstream frontend {
        server frontend:80;
    }

    server {
        listen 80;
        server_name _;

        # Real IP from Cloudflare
        set_real_ip_from 173.245.48.0/20;
        set_real_ip_from 103.21.244.0/22;
        set_real_ip_from 103.22.200.0/22;
        set_real_ip_from 103.31.4.0/22;
        set_real_ip_from 141.101.64.0/18;
        set_real_ip_from 108.162.192.0/18;
        set_real_ip_from 190.93.240.0/20;
        set_real_ip_from 188.114.96.0/20;
        set_real_ip_from 197.234.240.0/22;
        set_real_ip_from 198.41.128.0/17;
        set_real_ip_from 162.158.0.0/15;
        set_real_ip_from 104.16.0.0/13;
        set_real_ip_from 104.24.0.0/14;
        set_real_ip_from 172.64.0.0/13;
        set_real_ip_from 131.0.72.0/22;
        real_ip_header CF-Connecting-IP;

        # Health check
        location /health {
            return 200 'OK';
            add_header Content-Type text/plain;
        }

        # API routes
        location /api/ {
            proxy_pass http://api_gateway;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Frontend
        location / {
            proxy_pass http://frontend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;
        }
    }
}
EOF
```

---

## PHASE 4: CLOUDFLARE SETTINGS (5 phút)

### 4.1 Login Cloudflare Dashboard

Vào https://dash.cloudflare.com → Chọn domain

### 4.2 SSL/TLS Settings

```
SSL/TLS → Overview → Full (strict)
```

### 4.3 Security Settings (khuyến nghị)

```
Security → Settings:
- Security Level: Medium
- Challenge Passage: 30 minutes
- Browser Integrity Check: ON

Security → WAF:
- Enable Cloudflare Managed Rules (Free plan có basic)

Security → Bots:
- Bot Fight Mode: ON
```

### 4.4 Caching (optional)

```
Caching → Configuration:
- Caching Level: Standard
- Browser Cache TTL: 4 hours

Rules → Page Rules (nếu cần):
- erp.company.vn/api/* → Cache Level: Bypass
```

---

## PHASE 5: DEPLOY (10 phút)

### 5.1 Start services

```bash
cd /opt/erp

# Start all
docker-compose -f docker-compose.prod.yml up -d

# Check logs
docker-compose -f docker-compose.prod.yml logs -f cloudflared
```

### 5.2 Verify tunnel

```bash
# Check tunnel status
cloudflared tunnel info erp-tunnel

# Check từ Cloudflare Dashboard
# Zero Trust → Access → Tunnels → erp-tunnel → Status: HEALTHY
```

### 5.3 Test access

```bash
# Test từ internet
curl https://erp.company.vn/health
curl https://erp.company.vn/api/v1/health

# Test login
curl -X POST https://erp.company.vn/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@company.vn","password":"Admin@123"}'
```

---

## PHASE 6: SETUP TUNNEL AS SERVICE (5 phút)

### 6.1 Install as systemd service

```bash
# Install service
sudo cloudflared service install

# Hoặc manually create service
sudo cat > /etc/systemd/system/cloudflared.service << 'EOF'
[Unit]
Description=Cloudflare Tunnel
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/bin/cloudflared tunnel --config /opt/erp/cloudflared/config.yml run
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable cloudflared
sudo systemctl start cloudflared

# Check status
sudo systemctl status cloudflared
```

**Lưu ý**: Nếu chạy cloudflared trong Docker thì KHÔNG cần systemd service.

---

## FIREWALL SETUP (Strict Mode)

```bash
# Với Cloudflare Tunnel, chỉ cần mở SSH
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh

# KHÔNG cần mở 80, 443!
sudo ufw enable

# Verify
sudo ufw status
```

---

## 📋 CHECKLIST

```markdown
## Cloudflare Setup
- [ ] Domain added to Cloudflare
- [ ] cloudflared installed
- [ ] Tunnel created
- [ ] DNS record created (CNAME)
- [ ] config.yml configured

## Server
- [ ] Docker installed
- [ ] .env configured
- [ ] Images built/pulled
- [ ] docker-compose up

## Verify
- [ ] Tunnel status: HEALTHY
- [ ] https://erp.company.vn accessible
- [ ] Login works
- [ ] API responses correct

## Security
- [ ] Firewall: only SSH open
- [ ] Cloudflare SSL: Full (strict)
- [ ] WAF enabled
```

---

## TROUBLESHOOTING

### Tunnel không connect

```bash
# Check logs
docker logs erp-cloudflared
# hoặc
journalctl -u cloudflared -f

# Common issues:
# 1. credentials file không đúng path
# 2. tunnel ID sai trong config.yml
# 3. Network issue
```

### 502 Bad Gateway

```bash
# Check nginx có kết nối được service không
docker exec erp-nginx curl http://api-gateway:8080/health

# Check service logs
docker logs erp-api-gateway
```

### SSL Certificate Error

```bash
# Cloudflare Dashboard → SSL/TLS → Overview
# Đổi sang "Full" (không phải "Full strict") nếu có issue
```

---

## MULTIPLE SERVICES (Optional)

Nếu muốn expose thêm services:

```yaml
# /opt/erp/cloudflared/config.yml
ingress:
  - hostname: erp.company.vn
    service: http://nginx:80
  
  - hostname: grafana.erp.company.vn
    service: http://grafana:3000
  
  - hostname: minio.erp.company.vn
    service: http://minio:9001
  
  - service: http_status:404
```

Rồi tạo DNS records:
```bash
cloudflared tunnel route dns erp-tunnel grafana.erp.company.vn
cloudflared tunnel route dns erp-tunnel minio.erp.company.vn
```

---

## SO SÁNH VỚI TRADITIONAL DEPLOYMENT

| Aspect | Traditional | Cloudflare Tunnel |
|--------|-------------|-------------------|
| Setup time | 2-3 hours | 1 hour |
| SSL management | Manual (Let's Encrypt) | Automatic |
| DDoS protection | None/Paid | Free (basic) |
| IP exposure | Exposed | Hidden |
| Firewall rules | Complex | Simple (SSH only) |
| Cost | Free | Free |
| Latency | Direct | +5-20ms |

---

**Tổng thời gian deploy với Cloudflare Tunnel: ~1 giờ** 🚀
