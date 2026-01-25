# 🚀 HƯỚNG DẪN DEPLOY LÊN PRODUCTION

## ✅ Checklist Trước Khi Deploy

- [x] Docker đã cài: ✅ (version 28.2.2)
- [x] Docker Compose đã cài: ✅ (version 1.29.2)  
- [x] Cloudflared đã cài: ✅ (version 2026.1.1)
- [x] Tunnel ID: `5e9dfecb-38c1-4c8a-a2b1-127c45ce1092`
- [x] Domain: `erp.xelu.top`
- [ ] Tunnel connected qua Cloudflare Dashboard
- [ ] Build Docker images
- [ ] Start services

---

## 🎯 BƯỚC 1: SETUP CLOUDFLARE TUNNEL QUA DASHBOARD (5 phút)

### 1.1 Truy cập Cloudflare Zero Trust

1. Vào https://one.dash.cloudflare.com
2. Login với Cloudflare account
3. Nếu chưa có Zero Trust, click **"Get started для"** (miễn phí)

### 1.2 Tạo/Quản lý Tunnel

1. Menu bên trái: **Networks** → **Tunnels**
2. Tìm tunnel `erp-tunnel` (ID: `5e9dfecb-38c1-4c8a-a2b1-127c45ce1092`)
   - Nếu chưa có, click **Create a tunnel**
   - Nếu đã có, click vào tunnel name để edit

### 1.3 Configure Public Hostnames

Click tab **"Public Hostname"**, rồi add 2 hostnames:

#### Hostname 1: Main ERP App
- Click **"Add a public hostname"**
- **Subdomain**: `erp`
- **Domain**: `xelu.top`
- **Path**: (để trống)
- **Service**:
  - Type: `HTTP`
  - URL: `localhost:80`
- **Additional settings**: 
  - No TLS Verify: ON
- Click **Save hostname**

#### Hostname 2: Grafana Monitoring
- Click **"Add a public hostname"** again
- **Subdomain**: `grafana.erp`
- **Domain**: `xelu.top`
- **Path**: (để trống)
- **Service**:
  - Type: `HTTP`
  - URL: `localhost:3000`
- **Additional settings**:
  - No TLS Verify: ON
- Click **Save hostname**

### 1.4 Install Connector (Lấy Token)

1. Click tab **"Install and run a connector"**
2. Chọn **Docker**
3. Copy lệnh docker run, nó sẽ có dạng:

```bash
docker run cloudflare/cloudflared:latest tunnel --no-autoupdate run --token eyJhIjoiNzg5Mzg3NDMy...VERY_LONG_TOKEN
```

4. **Copy phần TOKEN** (phần sau `--token`)

5. Tạo file lưu token:

```bash
cd /opt/ERP
echo "TUNNEL_TOKEN=eyJhIjoiNzg5Mzg3NDMy...YOUR_LONG_TOKEN" >> .env
```

**Thay `YOUR_LONG_TOKEN` bằng token thực từ dashboard**

---

## 🎯 BƯỚC 2: BUILD DOCKER IMAGES (10-20 phút)

### 2.1 Check các services có Dockerfile

```bash
cd /opt/ERP

# Liệt kê services
ls -la services/

# Check Dockerfile
find services/ -name Dockerfile
```

### 2.2 Build từng service

```bash
# API Gateway
docker build -t erp/api-gateway:latest services/api-gateway/

# Auth Service
docker build -t erp/auth-service:latest services/auth-service/

# User Service
docker build -t erp/user-service:latest services/user-service/

# Master Data Service (nếu có)
docker build -t erp/master-data-service:latest services/master-data-service/

# WMS Service
docker build -t erp/wms-service:latest services/wms-service/

# Manufacturing Service
docker build -t erp/manufacturing-service:latest services/manufacturing-service/

# Sales Service
docker build -t erp/sales-service:latest services/sales-service/

# Marketing Service
docker build -t erp/marketing-service:latest services/marketing-service/

# Notification Service
docker build -t erp/notification-service:latest services/notification-service/

# Reporting Service
docker build -t erp/reporting-service:latest services/reporting-service/

# File Service
docker build -t erp/file-service:latest services/file-service/
```

**Hoặc build tất cả một lúc:**

```bash
docker-compose build
```

### 2.3 Verify images

```bash
docker images | grep erp
```

---

## 🎯 BƯỚC 3: CẬP NHẬT .ENV CHO PRODUCTION

### 3.1 Update domain

```bash
cd /opt/ERP
nano .env
```

Sửa dòng:
```bash
DOMAIN=erp.xelu.top
CORS_ALLOWED_ORIGINS=https://erp.xelu.top,https://grafana.erp.xelu.top
```

### 3.2 Update passwords (QUAN TRỌNG!)

```bash
# Generate strong passwords
openssl rand -base64 32

# Update trong .env:
POSTGRES_PASSWORD=<strong_password_1>
REDIS_PASSWORD=<strong_password_2>
MINIO_ROOT_PASSWORD=<strong_password_3>
GRAFANA_ADMIN_PASSWORD=<strong_password_4>
JWT_SECRET=<strong_64_char_string>
```

### 3.3 Add TUNNEL_TOKEN

```bash
# Thêm vào cuối file .env
TUNNEL_TOKEN=<token_from_cloudflare_dashboard>
```

---

## 🎯 BƯỚC 4: TẠO DOCKER-COMPOSE.PROD.YML

File `/opt/ERP/docker-compose.prod.yml` đã được tạo tại:
- `/opt/erp/cloudflared/config.yml` ✅
- `/opt/erp/deploy/nginx/nginx-cloudflare.conf` ✅

Cần thêm cloudflared service vào docker-compose:

```bash
cd /opt/ERP
nano docker-compose.yml
```

Thêm service sau vào đầu (sau phần volumes):

```yaml
  # ====================
  # Cloudflare Tunnel
  # ====================
  cloudflared:
    image: cloudflare/cloudflared:latest
    container_name: erp-cloudflared
    restart: always
    command: tunnel --no-autoupdate run --token ${TUNNEL_TOKEN}
    networks:
      - erp-network
    depends_on:
      - nginx
```

Và sửa nginx service để không expose port ra ngoài:

```yaml
  nginx:
    image: nginx:alpine
    container_name: erp-nginx
    volumes:
      - ./deploy/nginx/nginx-cloudflare.conf:/etc/nginx/nginx.conf  # ← Đổi file config
    # KHÔNG cần expose ports - cloudflared sẽ handle
    # ports:
    #   - "80:80"
    #   - "443:443"
    networks:
      - erp-network
    restart: unless-stopped
    depends_on:
      - api-gateway
```

---

## 🎯 BƯỚC 5: KHỞI ĐỘNG HỆ THỐNG (5 phút)

### 5.1 Start infrastructure trước

```bash
cd /opt/ERP

# Start database, cache, message queue
docker-compose up -d postgres redis nats minio
```

### 5.2 Chờ PostgreSQL ready

```bash
# Check health
docker-compose ps

# Chờ postgres healthy (màu xanh)
watch -n 2 'docker-compose ps postgres'
```

### 5.3 Start monitoring

```bash
docker-compose up -d prometheus grafana loki
```

### 5.4 Start services

```bash
# Start backend services
docker-compose up -d auth-service user-service api-gateway

# Nếu có frontend
docker-compose up -d frontend

# Start nginx
docker-compose up -d nginx

# Start Cloudflare Tunnel
docker-compose up -d cloudflared
```

### 5.5 Check logs

```bash
# Check cloudflared connection
docker logs erp-cloudflared

# Should see: "Connection ... registered"

# Check nginx
docker logs erp-nginx

# Check API gateway
docker logs erp-api-gateway
```

---

## 🎯 BƯỚC 6: VERIFY DEPLOYMENT

### 6.1 Check Cloudflare Dashboard

1. Vào https://one.dash.cloudflare.com
2. **Networks** → **Tunnels** → `erp-tunnel`
3. Status phải là **HEALTHY** (màu xanh)

### 6.2 Test URLs

```bash
# Test main app
curl https://erp.xelu.top/health

# Test API
curl https://erp.xelu.top/api/v1/health

# Test Grafana
curl https://grafana.erp.xelu.top
```

### 6.3 Test trong browser

1. Mở https://erp.xelu.top
2. Kiểm tra login page hiển thị
3. Test login với:
   - Email: `admin@company.vn`
   - Password: `Admin@123` (hoặc password mặc định của bạn)

4. Kiểm tra Grafana: https://grafana.erp.xelu.top
   - Username: `admin`
   - Password: (value trong .env `GRAFANA_ADMIN_PASSWORD`)

---

## 🎯 BƯỚC 7: CLOUDFLARE SECURITY SETTINGS

### 7.1 SSL/TLS Settings

1. Vào https://dash.cloudflare.com
2. Chọn domain `xelu.top`
3. **SSL/TLS** → **Overview**
4. Set to: **Full (strict)** hoặc **Full**

### 7.2 Security Settings

1. **Security** → **Settings**:
   - Security Level: **Medium**
   - Browser Integrity Check: **ON**

2. **Security** → **WAF** (Web Application Firewall):
   - Enable **Cloudflare Managed Ruleset**
   - Enable **OWASP Core Ruleset**

3. **Security** → **Bots**:
   - Bot Fight Mode: **ON**

### 7.3 Caching (Optional)

1. **Caching** → **Configuration**:
   - Caching Level: **Standard**

2. **Rules** → **Page Rules**:
   - Add rule: `erp.xelu.top/api/*`
   - Setting: **Cache Level** = **Bypass**

---

## 🔥 FIREWALL (QUAN TRỌNG!)

Với Cloudflare Tunnel, **CHỈ CẦN MỞ SSH**:

```bash
# Setup UFW
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh

# Enable
sudo ufw enable

# Check
sudo ufw status
```

**KHÔNG cần mở port 80, 443!** ✅

---

## ✅ POST-DEPLOYMENT CHECKLIST

- [ ] All services running: `docker-compose ps`
- [ ] Cloudflare tunnel: **HEALTHY**
- [ ] Can access: https://erp.xelu.top
- [ ] Can login to ERP
- [ ] Can access: https://grafana.erp.xelu.top
- [ ] SSL certificate: Valid (Green lock in browser)
- [ ] Firewall: Only SSH open
- [ ] Monitoring: Grafana dashboards working
- [ ] Logs: No critical errors in `docker logs`

---

## 🚨 TROUBLESHOOTING

### Tunnel shows "UNHEALTHY"

```bash
# Check cloudflared logs
docker logs -f erp-cloudflared

# Restart tunnel
docker-compose restart cloudflared
```

### 502 Bad Gateway

```bash
# Check nginx can reach services
docker exec erp-nginx curl http://api-gateway:8080/health

# Check service logs
docker logs erp-api-gateway
docker logs erp-auth-service
```

### SSL Certificate Error

- Cloudflare Dashboard → **SSL/TLS** → Change to **"Full"** (not Full strict)

---

## 📞 QUICK COMMANDS

```bash
# View all logs
docker-compose logs -f

# Restart all
docker-compose restart

# Stop all
docker-compose down

# Start all
docker-compose up -d

# Rebuild and restart
docker-compose up -d --build

# Check status
docker-compose ps
```

---

**🎉 Deploy thành công!** 

Access your ERP at: https://erp.xelu.top
