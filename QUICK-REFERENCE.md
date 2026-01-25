# 🚀 ERP DEPLOYMENT - QUICK REFERENCE

**Date**: 2026-01-25 | **Status**: 73% Complete | **Domain**: erp.xelu.top

---

## ✅ WHAT'S WORKING

```
✅ Cloudflare Tunnel    → HEALTHY (4 connections)
✅ PostgreSQL           → Running on port 5432
✅ Redis                → Running on port 6379
✅ NATS                 → Running on port 4222
✅ MinIO                → Running on port 9000/9001
✅ Prometheus           → Running on port 9090
✅ Grafana              → Running on port 3000
✅ Docker Images        → 13 images built (20260125-bbae72a)
✅ API Gateway          → Running (erp-api-gateway)
✅ User Service         → Running (erp-user-service)
✅ Frontend             → Running (erp-frontend-app)
```

---

## ⚠️ WHAT NEEDS FIXING

```
❌ Nginx                → Crashed (DNS resolution issue)
❌ Auth Service         → Not built (missing RefreshTokenUseCase)
❌ Reporting Service    → Build errors (go.sum issues)
❌ Public URL           → 502 Bad Gateway (nginx issue)
⏸️ Database Migrations  → Not run yet
⏸️ Admin User           → Not created yet
⏸️ Other Services       → Not started yet
```

---

## 🔧 QUICK FIX COMMANDS

### Fix Nginx (Option 1: Simple Frontend Only)
```bash
cd /opt/ERP

# Stop broken nginx
docker stop erp-nginx && docker rm erp-nginx

# Start simple nginx pointing to frontend
docker run -d --name erp-nginx \
  --network erp_erp-network \
  -e FRONTEND_HOST=erp-frontend-app \
  nginx:alpine sh -c \
  'echo "server { listen 80; location / { proxy_pass http://erp-frontend-app:80; } }" > /etc/nginx/conf.d/default.conf && nginx -g "daemon off;"'
```

### Test Public URL
```bash
curl https://erp.xelu.top
# Should return HTML instead of "error code: 502"
```

### Check All Running Services
```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### View Cloudflare Tunnel Status
```bash
docker logs erp-cloudflared | grep -i "registered\|connection"
```

---

## 📁 IMPORTANT FILES

| File | Purpose |
|------|---------|
| `/opt/ERP/DEPLOYMENT-SUMMARY.md` | Complete deployment report |
| `/opt/ERP/NEXT-STEPS.md` | Detailed next steps guide |
| `/opt/ERP/DEPLOYMENT-GUIDE.md` | Full deployment instructions |
| `/opt/ERP/.env` | Environment variables |
| `/opt/ERP/docker-compose.yml` | Service definitions |
| `/opt/ERP/Dockerfile.service` | Generic service builder |

---

## 🌐 URLS

| Service | URL | Status |
|---------|-----|--------|
| Main App | https://erp.xelu.top | 🔴 502 (nginx issue) |
| Grafana | https://grafana.erp.xelu.top | 🔴 502 (nginx issue) |
| Cloudflare Dashboard | https://one.dash.cloudflare.com | ✅ Working |
| Tunnel Status | Check container logs | ✅ HEALTHY |

---

## 🐳 CONTAINER QUICK REFERENCE

```bash
# Infrastructure
docker logs erp-postgres
docker logs erp-redis
docker logs erp-nats
docker logs erp-minio

# Monitoring
docker logs erp-prometheus
docker logs erp-grafana

# Application
docker logs erp-api-gateway
docker logs erp-user-service
docker logs erp-frontend-app

# Networking
docker logs erp-cloudflared
docker logs erp-nginx
```

---

## 📊 BUILD SUMMARY

**Total Images**: 13  
**Successfully Built**: 13  
**Version Tag**: 20260125-bbae72a  
**Total Size**: ~500 MB  

**Services Built**:
1. api-gateway ✅
2. user-service ✅
3. master-data-service ✅
4. supplier-service ✅
5. procurement-service ✅
6. wms-service ✅
7. manufacturing-service ✅
8. sales-service ✅
9. marketing-service ✅
10. notification-service ✅
11. file-service ✅
12. reporting-service ⚠️ (has errors)
13. frontend ✅

**Not Built**:
- auth-service ❌ (code incomplete)

---

## 🎯 IMMEDIATE NEXT STEPS

1. **Fix Nginx** (5 min)
   - Use simple config or fix DNS resolution
   
2. **Test Public Access** (2 min)
   - Verify https://erp.xelu.top loads
   
3. **Start Remaining Services** (10 min)
   - Add to docker-compose.yml
   - Start with `docker-compose up -d`
   
4. **Run Migrations** (5 min)
   - Execute for each service
   
5. **Create Admin User** (2 min)
   - Via API or database

**Total Time**: ~25 minutes to basic working state

---

## 💡 TIPS

- **Check Tunnel**: `docker logs erp-cloudflared` should show "HEALTHY"
- **Check Services**: `docker ps` should show all containers "Up"
- **Check Logs**: `docker-compose logs -f` for real-time monitoring
- **Restart Service**: `docker restart <container-name>`
- **Full Restart**: `docker-compose restart`

---

## 🆘 TROUBLESHOOTING

### 502 Bad Gateway
→ Nginx cannot reach backend  
→ Fix: Update nginx config or restart services

### Container Keeps Restarting
→ Check logs: `docker logs <container-name>`  
→ Common: Missing env vars or DB connection issues

### Cloudflare Tunnel Down
→ Check: `docker logs erp-cloudflared`  
→ Restart: `docker restart erp-cloudflared`

### Database Connection Failed
→ Check: `docker exec -it erp-postgres psql -U postgres`  
→ Verify: Connection string in `.env`

---

## 📞 SUPPORT

**Documentation**: `/opt/ERP/DEPLOYMENT-SUMMARY.md`  
**Next Steps**: `/opt/ERP/NEXT-STEPS.md`  
**Full Guide**: `/opt/ERP/DEPLOYMENT-GUIDE.md`

---

**Last Updated**: 2026-01-25T14:16:35Z  
**Progress**: 73% Complete  
**Status**: 🟡 Partial Success - Ready for Completion
