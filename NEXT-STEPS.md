# 📋 NEXT STEPS AFTER BUILD COMPLETES

## ✅ Current Status
- Cloudflare Tunnel: **CONNECTED** ✅
- Infrastructure: **RUNNING** ✅  
- Build Process: **IN PROGRESS** 🔄

---

## 🚀 STEP 1: Verify Build Completion

Sau khi build script hoàn tất, kiểm tra:

```bash
# Check build status
docker images | grep erp/

# Should see 14 images (13 services + frontend):
# erp/api-gateway:latest
# erp/api-gateway:20260125-bbae72a
# erp/auth-service:latest
# ... (tất cả 13 services)
# erp/frontend:latest
```

---

## 🚀 STEP 2: Start All Services

```bash
cd /opt/ERP

# Start tất cả services
docker-compose up -d

# Hoặc start từng nhóm:
# 1. Backend services
docker-compose up -d auth-service user-service api-gateway

# 2. Business services
docker-compose up -d master-data-service supplier-service procurement-service \
                     wms-service manufacturing-service sales-service \
                     marketing-service notification-service file-service \
                     reporting-service

# 3. Frontend
docker-compose up -d frontend

# 4. Nginx (đã có cloudflared rồi, có thể bỏ qua nginx nếu muốn)
# docker-compose up -d nginx
```

---

## 🚀 STEP 3: Check Service Health

```bash
# View all running containers
docker-compose ps

# Check logs của từng service
docker-compose logs -f auth-service
docker-compose logs -f api-gateway

# Check health của tất cả services
for port in 8080 8081 8082 8083 8084 8085 8086 8087 8088 8089 8090 8091 8092; do
    echo "Checking port $port..."
    curl -s http://localhost:$port/health || echo "FAILED"
done
```

---

## 🚀 STEP 4: Run Database Migrations

Nếu services cần migrations:

```bash
# Auth service
docker-compose exec auth-service ./migrate up

# User service  
docker-compose exec user-service ./migrate up

# Repeat for all services...
```

Hoặc nếu có script:

```bash
./scripts/run-migrations.sh
```

---

## 🚀 STEP 5: Create Initial Admin User

```bash
# Tạo admin user đầu tiên
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.vn",
    "password": "Admin@123456",
    "full_name": "System Administrator",
    "role": "admin"
  }'
```

---

## 🚀 STEP 6: Test Application

### 6.1 Test qua Cloudflare Tunnel

```bash
# Test health endpoint
curl https://erp.xelu.top/health

# Test API
curl https://erp.xelu.top/api/v1/health

# Test login
curl -X POST https://erp.xelu.top/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.vn",
    "password": "Admin@123456"
  }'
```

### 6.2 Test trong Browser

1. Mở: **https://erp.xelu.top**
2. Kiểm tra login page hiển thị
3. Login với:
   - Email: `admin@company.vn`
   - Password: `Admin@123456`
4. Kiểm tra dashboard

### 6.3 Test Grafana

1. Mở: **https://grafana.erp.xelu.top**
2. Login:
   - Username: `admin`
   - Password: (value trong `.env` - `GRAFANA_ADMIN_PASSWORD`)
3. Check dashboards

---

## 🚀 STEP 7: Configure Cloudflare Security (Optional)

### 7.1 SSL/TLS Settings
1. Vào https://dash.cloudflare.com
2. Chọn domain `xelu.top`
3. **SSL/TLS** → **Overview** → Set to **"Full"** hoặc **"Full (strict)"**

### 7.2 WAF (Web Application Firewall)
1. **Security** → **WAF**
2. Enable **Cloudflare Managed Ruleset**
3. Enable **OWASP Core Ruleset**

### 7.3 Rate Limiting
1. **Security** → **WAF** → **Rate limiting rules**
2. Add rule:
   - Name: "API Rate Limit"
   - If: `URI Path contains /api/`
   - Then: Rate limit
   - Requests: 100 per minute

### 7.4 Bot Protection
1. **Security** → **Bots**
2. Enable **Bot Fight Mode**

---

## 🚀 STEP 8: Setup Monitoring Alerts

### 8.1 Configure Alertmanager

Edit `/opt/ERP/deploy/monitoring/alertmanager.yml`:

```yaml
route:
  receiver: 'email'
  
receivers:
  - name: 'email'
    email_configs:
      - to: 'admin@company.vn'
        from: 'erp-alerts@company.vn'
        smarthost: 'smtp.gmail.com:587'
        auth_username: 'erp-alerts@company.vn'
        auth_password: 'your-app-password'
```

Restart alertmanager:
```bash
docker-compose restart alertmanager
```

---

## 🚀 STEP 9: Setup Backups

### 9.1 Database Backup Script

File `/opt/ERP/scripts/backup-db.sh` đã có, setup cron:

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /opt/ERP/scripts/backup-db.sh >> /var/log/erp-backup.log 2>&1

# Add weekly cleanup on Sunday
0 4 * * 0 /opt/ERP/scripts/cleanup-backups.sh >> /var/log/erp-backup.log 2>&1
```

### 9.2 Test Backup

```bash
# Run manual backup
./scripts/backup-db.sh

# Verify
ls -lh /opt/ERP/backups/
```

---

## 🚀 STEP 10: Final Verification Checklist

```markdown
## Infrastructure
- [ ] All containers running (`docker ps`)
- [ ] No containers restarting
- [ ] Disk space OK (`df -h`)
- [ ] Memory usage OK (`free -h`)

## Application
- [ ] Login works
- [ ] API responses correct  
- [ ] Frontend loads
- [ ] Can create/read data

## Monitoring
- [ ] Grafana accessible
- [ ] Dashboards showing data
- [ ] Prometheus targets UP
- [ ] Logs flowing to Loki

## Security
- [ ] Cloudflare Tunnel: HEALTHY
- [ ] SSL certificate valid
- [ ] WAF enabled
- [ ] Strong passwords set
- [ ] Firewall: only SSH open

## Backup
- [ ] Backup script works
- [ ] Backup files created
- [ ] Cron jobs configured
```

---

## 📞 QUICK COMMANDS

```bash
# View all logs
docker-compose logs -f

# Restart a service
docker-compose restart auth-service

# Stop all
docker-compose down

# Start all
docker-compose up -d

# Rebuild a service
docker-compose up -d --build auth-service

# Check Cloudflare Tunnel
docker logs erp-cloudflared

# Check nginx
docker logs erp-nginx
```

---

## 🚨 TROUBLESHOOTING

### Service won't start
```bash
# Check logs
docker-compose logs service-name

# Check if port is in use
netstat -tulpn | grep PORT

# Restart service
docker-compose restart service-name
```

### 502 Bad Gateway
```bash
# Check if backend is running
docker ps | grep erp-

# Check nginx can reach backend
docker exec erp-nginx curl http://api-gateway:8080/health

# Check Cloudflare Tunnel
docker logs erp-cloudflared
```

### Database connection errors
```bash
# Check postgres is running
docker-compose ps postgres

# Check connection
docker-compose exec postgres psql -U postgres -c "SELECT 1"

# Check service can connect
docker-compose exec auth-service env | grep DB_
```

---

## 🎉 SUCCESS CRITERIA

Deployment thành công khi:

1. ✅ `https://erp.xelu.top` accessible
2. ✅ Login works
3. ✅ API returns data
4. ✅ Grafana shows metrics
5. ✅ All services healthy
6. ✅ No errors in logs
7. ✅ Cloudflare Tunnel: HEALTHY

---

**Sau khi build xong, bắt đầu từ STEP 2!** 🚀
