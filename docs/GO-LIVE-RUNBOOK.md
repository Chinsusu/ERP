# GO-LIVE RUNBOOK

**Version**: 1.0  
**Date**: 2026-01-25

---

## TỔNG QUAN

**Ngày Go-Live**: Saturday (giảm thiểu ảnh hưởng)  
**Khung giờ**: 6:00 AM - 12:00 PM  
**Hypercare**: 24h đầu tiên

---

## LIÊN HỆ KHẨN CẤP

| Role | Name | Phone | Backup |
|------|------|-------|--------|
| Tech Lead | | | |
| DBA | | | |
| DevOps | | | |
| QA Lead | | | |
| Business Owner | | | |

---

## TIMELINE CHI TIẾT

### Phase 1: Chuẩn bị (6:00 - 6:30)

| Time | Task | Owner | Status |
|------|------|-------|--------|
| 6:00 | Team standup call - xác nhận ready | Tech Lead | ☐ |
| 6:10 | Verify tất cả thành viên online | Tech Lead | ☐ |
| 6:15 | Khóa hệ thống cũ (read-only) | DBA | ☐ |
| 6:20 | Backup cuối cùng hệ thống cũ | DBA | ☐ |
| 6:30 | **CHECKPOINT 1**: Confirm proceed | All | ☐ |

### Phase 2: Data Migration (6:30 - 8:00)

| Time | Task | Owner | Status |
|------|------|-------|--------|
| 6:30 | Export data từ hệ thống cũ | DBA | ☐ |
| 7:00 | Chạy migration scripts | Dev Lead | ☐ |
| 7:30 | Validate record counts | QA | ☐ |
| 7:45 | Validate critical business data | Business | ☐ |
| 8:00 | **CHECKPOINT 2**: Data migration OK | All | ☐ |

### Phase 3: Application Deployment (8:00 - 9:00)

| Time | Task | Owner | Status |
|------|------|-------|--------|
| 8:00 | Pull latest images | DevOps | ☐ |
| 8:10 | Deploy infrastructure (DB, Redis, NATS) | DevOps | ☐ |
| 8:20 | Run database migrations (`make migrate`) | Dev Lead | ☐ |
| 8:30 | Deploy all 13 services | DevOps | ☐ |
| 8:45 | Health check all services | DevOps | ☐ |
| 9:00 | **CHECKPOINT 3**: Deployment complete | All | ☐ |

### Phase 4: Verification (9:00 - 10:30)

| Time | Task | Owner | Status |
|------|------|-------|--------|
| 9:00 | Smoke test - Login | QA | ☐ |
| 9:10 | Smoke test - Master Data | QA | ☐ |
| 9:20 | Smoke test - WMS (GRN, FEFO) | QA | ☐ |
| 9:35 | Smoke test - Manufacturing | QA | ☐ |
| 9:50 | Smoke test - Sales | QA | ☐ |
| 10:00 | End-to-end test: Order → Stock | QA | ☐ |
| 10:15 | Performance spot check | DevOps | ☐ |
| 10:30 | **CHECKPOINT 4**: Verification OK | QA Lead | ☐ |

### Phase 5: Go-Live (10:30 - 11:00)

| Time | Task | Owner | Status |
|------|------|-------|--------|
| 10:30 | Enable user access | Admin | ☐ |
| 10:35 | Monitor first user logins | DevOps | ☐ |
| 10:45 | Monitor system metrics | DevOps | ☐ |
| 11:00 | **GO-LIVE ANNOUNCEMENT** | PM | ☐ |

### Phase 6: Hypercare (11:00 - 12:00+)

| Time | Task | Owner | Status |
|------|------|-------|--------|
| 11:00 | Support channels active | Support | ☐ |
| 11:30 | First hour report | Tech Lead | ☐ |
| 12:00 | Handover to support team | Tech Lead | ☐ |
| +24h | Hypercare monitoring | Team | ☐ |

---

## ROLLBACK PROCEDURE

### Khi nào Rollback?

| Trigger | Condition |
|---------|-----------|
| 🔴 Critical | Lỗi ảnh hưởng > 50% users |
| 🔴 Data | Phát hiện data corruption |
| 🔴 Performance | < 50% target (p95 > 400ms) |
| 🔴 Security | Phát hiện breach |

### Decision Matrix

| Issue | Severity | Rollback? | Approver |
|-------|----------|-----------|----------|
| Minor UI bug | P3 | No, fix forward | Dev Lead |
| Single feature broken | P2 | No, disable feature | Tech Lead |
| Core flow broken | P1 | **YES** | Tech Lead + Business |
| Data corruption | P1 | **YES** | Tech Lead + DBA |

### Rollback Steps

```bash
# 1. Announce maintenance
echo "Đang bảo trì hệ thống..."

# 2. Disable user access
# (Toggle maintenance mode)

# 3. Stop new services
docker-compose -f docker-compose.prod.yml down

# 4. Restore database backup
./scripts/restore-db.sh /backups/pre-golive.sql.gz

# 5. Restart old system (if applicable)

# 6. Verify old system working

# 7. Communicate to users
```

### Rollback Deadline
- ⏱️ Quyết định rollback phải được đưa ra **trước 11:00 AM**
- Sau 11:00 → chỉ fix forward, không rollback

---

## POST GO-LIVE MONITORING

### Metrics to Watch (First 24h)

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Error rate | < 0.1% | > 1% |
| Response p95 | < 200ms | > 500ms |
| CPU usage | < 70% | > 85% |
| Memory usage | < 70% | > 85% |
| Active users | Baseline | -50% or spike |

### Checklist - End of Day 1

- [ ] No P1/P2 incidents
- [ ] All critical flows working
- [ ] User feedback collected
- [ ] Metrics within targets
- [ ] Next day plan confirmed

---

## APPENDIX: Commands Reference

```bash
# Deploy
./scripts/deploy.sh v1.0.0

# Health check
./scripts/health-check.sh

# View logs
docker-compose logs -f --tail=100 api-gateway

# Backup
./scripts/backup-db.sh

# Restore
./scripts/restore-db.sh /path/to/backup.sql.gz
```
