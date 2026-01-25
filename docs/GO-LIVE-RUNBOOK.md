# GO-LIVE RUNBOOK

## Timeline: Saturday 6:00 AM - 12:00 PM

### Phase 1: Preparation (6:00 - 6:30)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 6:00 | Team standup call | Tech Lead | |
| 6:10 | Verify all team members online | Tech Lead | |
| 6:15 | Lock old system (read-only) | DBA | |
| 6:20 | Take final backup of old system | DBA | |
| 6:30 | Confirm ready to proceed | All | |

### Phase 2: Data Migration (6:30 - 8:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 6:30 | Export data from old system | DBA | |
| 7:00 | Run migration scripts | Dev Lead | |
| 7:30 | Validate data counts | QA | |
| 8:00 | Data migration sign-off | All | |

### Phase 3: Application Deployment (8:00 - 9:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 8:00 | Deploy all services (scripts/deploy.sh) | DevOps | |
| 8:15 | Run database migrations (make migrate) | Dev Lead | |
| 8:45 | Health check all services | DevOps | |
| 9:00 | Deployment complete | Tech Lead | |

### Phase 4: Verification & Smoke Test (9:00 - 10:30)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 9:00 | Smoke test - Login | QA | |
| 9:10 | Smoke test - WMS FEFO Flow | QA | |
| 9:30 | Smoke test - Mfg BOM Access | QA | |
| 10:30 | Verification sign-off | QA Lead | |

### Phase 5: Go-Live (10:30 - 11:00)
| Time | Task | Owner | Status |
|------|------|-------|--------|
| 10:30 | Enable user access | Admin | |
| 11:00 | Official go-live announcement | PM | |

---

## Rollback Procedure

### Trigger Conditions
- Lỗi nghiêm trọng ảnh hưởng > 50% người dùng.
- Phát hiện sai lệch dữ liệu (Data corruption).
- Hiệu năng hệ thống < 50% mục tiêu.

### Các bước Rollback
1. Thông báo bảo hệ thống.
2. Ngắt truy cập người dùng.
3. Dừng tất cả dịch vụ mới.
4. Restore bản backup dữ liệu cuối cùng của hệ thống cũ.
5. Kiểm tra và mở lại hệ thống cũ.
