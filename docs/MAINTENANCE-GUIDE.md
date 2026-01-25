# HƯỚNG DẪN BẢO TRÌ HỆ THỐNG ERP
## Maintenance Guide for ERP Cosmetics

**Version**: 1.0  
**Updated**: 2026-01-25

---

## 1. KIỂM TRA HÀNG NGÀY (Daily Operations)

### Morning Check (8:00 AM)

#### System Health
- [ ] **All services UP**: Check Grafana dashboard → Service Health panel
  - Expected: 13/13 services healthy
  - If down: Check `docker ps` and logs
- [ ] **No critical alerts**: Check AlertManager → Active Alerts
  - Investigate any P1/P2 alerts from overnight
- [ ] **Response time normal**: p95 < 200ms
  - If high: Check slow query log and resource usage
- [ ] **Error rate < 0.1%**: Check Grafana → Error Rate panel
  - If spiking: Review application logs for stack traces
- [ ] **Disk usage < 80%**: `df -h` on production server
  - If high: Run `./scripts/cleanup-backups.sh`

#### Application Health
- [ ] **Login working**: Test manual login at https://erp.company.vn
- [ ] **API Gateway responding**: `curl https://erp.company.vn/api/v1/health`
- [ ] **Database connections normal**: Check Grafana → DB Connections < 50
- [ ] **Redis cache healthy**: `docker exec erp-redis redis-cli PING`
- [ ] **NATS queue healthy**: Check http://localhost:8222/varz

#### Business Alerts
- [ ] **Check low stock alerts**: WMS → Stock → Low Stock Report
  - Create PRs for items below reorder point
- [ ] **Check expiring lots (30 days)**: WMS → Lots Expiring
  - Prioritize FEFO usage
- [ ] **Check pending approvals**: Procurement → Pending Approvals
  - Escalate if > 7 days old

### Evening Check (5:00 PM)
- [ ] **Review day's incidents**: Check incident log, update status
- [ ] **Verify backup completed**: Check `/backups/` for today's backup
  - Expected: `erp_prod_YYYYMMDD.sql.gz`
- [ ] **Check alert queue**: Any unresolved warnings
- [ ] **Update operations log**: File: `docs/ops-log.md`

---

## 2. BẢO TRÌ ĐỊNH KỲ (Periodic Maintenance)

### Hàng tuần (Every Monday, 9:00 AM)

#### Performance Review
- [ ] **Review weekly metrics**: Grafana → Last 7 days view
  - Response time trends
  - Error rate patterns
  - Traffic patterns
- [ ] **Identify slow queries**: 
  ```sql
  SELECT query, mean_exec_time, calls
  FROM pg_stat_statements
  ORDER BY mean_exec_time DESC
  LIMIT 10;
  ```
- [ ] **Check database growth**: 
  ```sql
  SELECT pg_size_pretty(pg_database_size('wms_db'));
  ```
- [ ] **Review error logs**: 
  ```bash
  docker logs erp-api-gateway --since 7d | grep ERROR
  ```

#### Security
- [ ] **Review failed login attempts**: 
  ```sql
  SELECT email, COUNT(*) 
  FROM auth.login_attempts 
  WHERE success = false 
    AND created_at > NOW() - INTERVAL '7 days'
  GROUP BY email 
  HAVING COUNT(*) > 10;
  ```
- [ ] **Check for suspicious activities**: Review audit logs
- [ ] **Review audit logs**: Look for unusual access patterns
- [ ] **Update blocked IPs**: Add to firewall rules if needed

#### Housekeeping
- [ ] **Clean old logs (>30 days)**: 
  ```bash
  ./scripts/cleanup-backups.sh
  ```
- [ ] **Vacuum PostgreSQL**: 
  ```sql
  VACUUM ANALYZE;
  ```
- [ ] **Clear Redis cache (if needed)**: 
  ```bash
  docker exec erp-redis redis-cli FLUSHDB
  ```
- [ ] **Archive old data**: Move completed orders > 1 year to archive DB

---

### Hàng tháng (First Saturday, 2:00 AM - 6:00 AM)

#### Planned Downtime Tasks
- [ ] **Apply security patches**: 
  ```bash
  apt update && apt upgrade
  ```
- [ ] **Update Docker images**: 
  ```bash
  docker-compose pull
  ./scripts/deploy.sh latest
  ```
- [ ] **Database maintenance**:
  - Full VACUUM
  - Reindex large tables
  - Update statistics
- [ ] **SSL certificate check**: 
  ```bash
  openssl x509 -in /etc/letsencrypt/live/erp.company.vn/cert.pem -noout -dates
  ```
- [ ] **Backup verification test**: 
  ```bash
  ./scripts/test-restore.sh /backups/erp_prod_latest.sql.gz
  ```

#### Capacity Planning
- [ ] **Review resource utilization**: CPU, RAM, Disk trends
- [ ] **Forecast growth**: Project next 3-6 months
- [ ] **Plan scaling**: Add resources if >70% utilized

#### Documentation
- [ ] **Update runbooks**: Add new issues encountered
- [ ] **Review procedures**: Update outdated steps
- [ ] **Update contact lists**: Verify on-call schedule

---

## 3. QUY TRÌNH KHÔI PHỤC (Recovery Procedures)

### Backup Verification Procedure

#### Monthly Backup Test

**Purpose**: Verify backups can be restored successfully

**Steps**:
1. **Prepare test environment**:
   ```bash
   docker run -d --name pg-restore-test \
     -e POSTGRES_PASSWORD=test \
     postgres:16-alpine
   ```

2. **Restore latest backup**:
   ```bash
   LATEST_BACKUP=$(ls -t /backups/db/*.sql.gz | head -1)
   gunzip -c $LATEST_BACKUP | docker exec -i pg-restore-test psql -U postgres
   ```

3. **Verify data**:
   ```sql
   SELECT 'users', COUNT(*) FROM users
   UNION ALL
   SELECT 'materials', COUNT(*) FROM materials
   UNION ALL
   SELECT 'stock', COUNT(*) FROM stock;
   ```

4. **Document results**:
   - Backup file: ________________
   - Backup date: ________________
   - Restore time: _____ minutes
   - Data verification: PASS / FAIL
   - Issues found: ________________

5. **Cleanup**:
   ```bash
   docker rm -f pg-restore-test
   ```

---

## 4. QUẢN LÝ SỰ CỐ (Incident Management)

### Severity Levels

| Level | Description | Response Time | Examples |
|-------|-------------|---------------|----------|
| **P1** | System down | 15 min | All services down, data loss |
| **P2** | Major impact | 30 min | Key function unavailable (WMS, Manufacturing) |
| **P3** | Minor impact | 4 hours | Non-critical bug, UI glitch |
| **P4** | Low impact | 24 hours | Cosmetic issues, feature requests |

### Escalation Path

```
P1: On-call → Tech Lead → CTO → CEO
P2: On-call → Tech Lead → CTO
P3: On-call → Tech Lead
P4: Support Team → Backlog
```

### Incident Response Steps

1. **Detection**: Alert received OR user reported
2. **Triage**: Assess severity, notify team
3. **Investigation**: Gather logs, identify cause
4. **Resolution**: Implement fix, test
5. **Communication**: Update stakeholders
6. **Documentation**: Complete incident report

---

## 5. TROUBLESHOOTING GUIDE

### Service Won't Start

**Symptoms**:
- Container exits immediately
- Health check failing

**Steps**:
1. Check logs: `docker logs <service>`
2. Check config: Environment variables
3. Check dependencies: Database, Redis, NATS
4. Check ports: Port conflicts
5. Check resources: Memory, disk

### Database Connection Issues

**Symptoms**:
- "connection refused" errors
- Timeout errors

**Steps**:
1. Check PostgreSQL running: `docker ps | grep postgres`
2. Check connection count: 
   ```sql
   SELECT count(*) FROM pg_stat_activity;
   ```
3. Check pg_hba.conf if auth issues
4. Restart connection pool

### High Memory Usage

**Symptoms**:
- OOM killer triggered
- Service restarts frequently

**Steps**:
1. Check memory: `docker stats`
2. Check for memory leaks
3. Review recent deployments
4. Increase memory limit temporarily
5. Fix root cause

### Slow Queries

**Symptoms**:
- High response time
- Database CPU high

**Steps**:
1. Find slow queries: `pg_stat_statements`
2. Check missing indexes
3. Analyze query plan: `EXPLAIN ANALYZE`
4. Add index or optimize query

### FEFO Not Working

**Symptoms**:
- Wrong lot selected for issue
- Complaint about lot selection

**Steps**:
1. Check lot expiry dates in DB
2. Verify FEFO logic in `wms-service/usecase/issue`
3. Check reservations aren't blocking lots
4. Verify no manual lot selection overrides

---

## 6. COMMON COMMANDS

### Docker
```bash
# View all services
docker-compose ps

# View logs
docker-compose logs -f --tail=100 api-gateway

# Restart service
docker-compose restart wms-service

# Execute command in container
docker exec -it erp-postgres psql -U postgres
```

### Database
```bash
# Connect to database
docker exec -it erp-postgres psql -U postgres -d wms_db

# Backup
./scripts/backup-db.sh

# Restore
./scripts/restore-db.sh /path/to/backup.sql.gz
```

### Monitoring
```bash
# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check active alerts
curl http://localhost:9093/api/v1/alerts

# View metrics
curl http://localhost:8086/metrics
```

---

## 7. CONTACTS

| Role | Name | Phone | Email | On-Call |
|------|------|-------|-------|---------|
| Tech Lead | | | | Mon-Fri |
| DevOps | | | | 24/7 |
| DBA | | | | Business hours |
| Support | | | | Business hours |

---

**Next Review**: 2026-02-25  
**Document Owner**: IT Operations Team
