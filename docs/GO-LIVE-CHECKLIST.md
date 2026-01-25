# PRE GO-LIVE CHECKLIST

## 1 Week Before

### Infrastructure
- [ ] Production server provisioned (scripts/server-setup.sh)
- [ ] SSL certificates installed (scripts/setup-ssl.sh)
- [ ] Firewall configured (UFW enabled)
- [ ] Backup system tested (scripts/backup.sh)
- [ ] Monitoring semi-configured (logs to /var/log/erp)

### Application
- [ ] All services unit-tested (Phase 7.1-7.3)
- [ ] Integration tests logic verified (Phase 7.4)
- [ ] Performance k6 scripts ready (Phase 7.5)
- [ ] Security audit recommendations implemented (Phase 7.6)

### People
- [ ] Training completed for all users (docs/TRAINING-GUIDE.md)
- [ ] Support team briefed
- [ ] Escalation path defined

---

## 1 Day Before

### Final Checks
- [ ] Final environment health check
- [ ] Backup taken of any existing systems
- [ ] Rollback plan reviewed
- [ ] Communication sent to users
- [ ] Support channels (Slack/Internal) ready

### Go/No-Go Decision
- [ ] All checklist items completed
- [ ] Sign-off received
