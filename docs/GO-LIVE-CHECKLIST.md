# PRE GO-LIVE CHECKLIST

**Version**: 1.0  
**Date**: 2026-01-25

---

## 1 TUẦN TRƯỚC GO-LIVE

### Infrastructure
- [ ] Production server provisioned (`scripts/server-setup.sh`)
- [ ] SSL certificates installed (`scripts/setup-ssl.sh`)
- [ ] Firewall configured (UFW rules: 80, 443, 22)
- [ ] Backup system tested (`scripts/backup.sh`)
- [ ] Monitoring stack running (Prometheus, Grafana, Loki)
- [ ] DNS configured pointing to production server

### Application
- [ ] All 13 services deployed to staging
- [ ] All unit tests passing (Phase 7.1-7.3)
- [ ] Integration tests verified (Phase 7.4)
- [ ] Load tests passed - p95 < 200ms (Phase 7.5)
- [ ] Security audit completed (Phase 7.6)
- [ ] UAT sign-off received from business

### Database
- [ ] All migrations tested on staging
- [ ] Seed data verified (roles, permissions, categories)
- [ ] Backup/restore tested successfully
- [ ] Database performance optimized (indexes checked)

### Data Migration (if migrating from old system)
- [ ] Migration scripts tested
- [ ] Data mapping document approved
- [ ] Opening balances verified
- [ ] Master data imported and validated
- [ ] Historical data migration plan confirmed

### People
- [ ] Training completed for all departments
- [ ] Support team briefed and ready
- [ ] On-call schedule confirmed
- [ ] Escalation path documented
- [ ] Communication plan ready

---

## 1 NGÀY TRƯỚC GO-LIVE

### Final Checks
- [ ] All services healthy on staging
- [ ] Final backup of current systems
- [ ] Rollback plan reviewed by team
- [ ] Communication sent to all users
- [ ] Support channels active (Slack/Teams/Hotline)

### Environment
- [ ] Production environment variables verified
- [ ] SSL certificate valid (not expiring soon)
- [ ] Domain resolving correctly
- [ ] Rate limits configured appropriately

### Team Readiness
- [ ] All team members confirmed available
- [ ] Contact numbers verified
- [ ] Access credentials distributed
- [ ] War room / virtual meeting set up

---

## GO/NO-GO DECISION

### Required Sign-offs

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Tech Lead | | | |
| QA Lead | | | |
| Business Owner | | | |
| IT Security | | | |

### Go Criteria
- ✅ All checklist items completed
- ✅ No P1/P2 bugs in staging
- ✅ Performance within targets
- ✅ All sign-offs received

### No-Go Triggers
- ❌ Critical infrastructure issues
- ❌ Data migration failures
- ❌ Security vulnerabilities unresolved
- ❌ Key personnel unavailable
