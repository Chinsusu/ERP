# ERP Cosmetics System - Project Summary

**Updated**: January 24, 2026  
**Repository**: https://github.com/Chinsusu/ERP  
**Status**: 6 Services Implemented (Phase 1 Core + Phase 2 Supply Chain)

---

## Executive Summary

Complete ERP system for cosmetics manufacturing with Clean Architecture, microservices pattern, and event-driven design. Phase 2 (Supply Chain) now complete with Supplier and Procurement services.

| Component | Files | Status |
|-----------|-------|--------|
| **Infrastructure** | 30+ | ✅ Production-ready |
| **Shared Libraries** | 10 packages | ✅ Production-ready |
| **API Gateway** | 17 files | ✅ Running |
| **Auth Service** | 50 files | ✅ Tested |
| **User Service** | 39 files | ✅ Tested |
| **Master Data Service** | 47 files | ✅ Built |
| **Supplier Service** | 40 files | ✅ Running |
| **Procurement Service** | 35 files | ✅ Running |

**Total**: ~260+ files, ~22,000+ LOC

---

## Architecture Overview

```
                    ┌─────────────────────────────────────┐
                    │           Nginx (SSL/TLS)           │
                    └─────────────────────────────────────┘
                                     │
                    ┌─────────────────────────────────────┐
                    │        API Gateway (:8080)          │
                    │  Rate Limit • Auth • Circuit Break  │
                    └─────────────────────────────────────┘
                                     │
    ┌───────────┬───────────┬────────┴────────┬────────────┬─────────────┐
    │           │           │                 │            │             │
┌───▼───┐  ┌────▼───┐  ┌────▼─────┐   ┌───────▼───┐  ┌─────▼─────┐       │
│ Auth  │  │ User   │  │ Master   │   │ Supplier  │  │Procurement│  ... 15
│ :8081 │  │ :8082  │  │  Data    │   │  :8084    │  │  :8085    │ services
└───┬───┘  └────┬───┘  │  :8083   │   └─────┬─────┘  └─────┬─────┘
    │           │      └────┬─────┘         │              │
┌───▼───┐  ┌────▼───┐  ┌────▼─────┐   ┌─────▼─────┐  ┌─────▼─────┐
│auth_db│  │user_db │  │master_   │   │supplier_  │  │procure_   │
│       │  │        │  │data_db   │   │   db      │  │ment_db    │
└───────┘  └────────┘  └──────────┘   └───────────┘  └───────────┘
                                            │              │
                                       ┌────┴──────────────┴────┐
                                       │    NATS JetStream      │
                                       │    (Event Bus)         │
                                       └────────────────────────┘
```

---

## Phase 1: Core Services ✅

### 1. API Gateway (Port 8080)
Single entry point for all requests.

| Feature | Implementation |
|---------|---------------|
| **Routing** | 15 services configured |
| **Rate Limiting** | Redis (100/min user, 30/min IP) |
| **Auth** | JWT validation + blacklist |
| **Circuit Breaker** | 5 failures → 30s open |
| **Health** | `/health`, `/ready`, `/live` |

### 2. Auth Service (Port 8081)
Authentication & authorization with JWT + RBAC.

| Component | Details |
|-----------|---------|
| **Tables** | 8 (roles, permissions, credentials, sessions) |
| **Endpoints** | login, logout, refresh |
| **Security** | Bcrypt, 5-attempt lockout, token rotation |
| **Defaults** | 5 roles, 42 permissions, admin user |

### 3. User Service (Port 8082)
User and department management.

| Component | Details |
|-----------|---------|
| **Tables** | 3 (departments, users, profiles) |
| **Endpoints** | users CRUD, departments tree |
| **Features** | Auto employee code (EMP20260124xxx) |
| **Defaults** | 8 departments (EXEC, IT, HR, etc.) |

### 4. Master Data Service (Port 8083)
Materials, products, categories, units.

| Component | Details |
|-----------|---------|
| **Tables** | 7 (categories, units, materials, products, etc.) |
| **Endpoints** | 20+ (CRUD + search + specifications) |
| **Cosmetics** | INCI names, CAS numbers, allergens |
| **Defaults** | 12 units, 27 categories, 8 conversions |

---

## Phase 2: Supply Chain ✅

### 5. Supplier Service (Port 8084) ✅ NEW
Supplier management with cosmetics-specific certifications.

| Component | Details |
|-----------|---------|
| **Tables** | 7 (suppliers, addresses, contacts, certifications, evaluations, ASL, price lists) |
| **Endpoints** | 16 (suppliers CRUD, certifications, evaluations) |
| **Certifications** | GMP, ISO9001, ISO22716, ORGANIC, ECOCERT, HALAL, COSMOS |
| **Features** | Expiry tracking, rating calculation, approval workflow |
| **Events** | supplier.created, approved, blocked, certification.expiring |

### 6. Procurement Service (Port 8085) ✅ NEW
Purchase Requisitions and Purchase Orders management.

| Component | Details |
|-----------|---------|
| **Tables** | 7 (PRs, PR items, approvals, POs, PO items, amendments, receipts) |
| **Endpoints** | 13 (PR workflow, PO workflow, receipts) |
| **PR Workflow** | DRAFT → SUBMITTED → APPROVED/REJECTED → CONVERTED_TO_PO |
| **PO Workflow** | DRAFT → CONFIRMED → PARTIALLY_RECEIVED → FULLY_RECEIVED → CLOSED |
| **Approval Levels** | Auto (<10M), Dept Manager (<50M), Procurement (<200M), CFO (>200M VND) |
| **Events** | pr.created/submitted/approved, po.created/confirmed/received/closed |
| **gRPC** | GetPO, GetPOsBySupplier, UpdatePOReceivedQty (WMS integration) |
| **Subscribers** | wms.grn.completed, supplier.blocked |

---

## Infrastructure Stack

```yaml
Services:
  - PostgreSQL 16     # 13 databases
  - Redis 7           # Caching & sessions
  - NATS JetStream    # Event streaming
  - MinIO             # Object storage
  - Prometheus        # Metrics
  - Grafana           # Dashboards
  - Loki              # Logging
  - Jaeger            # Tracing
  - Nginx             # Reverse proxy
```

---

## Shared Libraries

```
shared/pkg/
├── config/       # Viper configuration
├── database/     # GORM PostgreSQL
├── logger/       # Zap structured logging
├── middleware/   # HTTP middlewares
├── errors/       # Custom error types
├── validator/    # Input validation
├── jwt/          # JWT utilities
├── grpc/         # gRPC helpers
├── nats/         # Event pub/sub
└── response/     # HTTP responses
```

---

## Quick Start

### 1. Start Infrastructure
```bash
cd /opt/ERP
docker start erp-postgres erp-redis erp-nats
```

### 2. Start Services
```bash
# API Gateway
cd services/api-gateway && make run
# http://localhost:8080

# Auth Service
cd services/auth-service && make run
# http://localhost:8081

# User Service
cd services/user-service && make run
# http://localhost:8082

# Master Data Service
cd services/master-data-service && make run
# http://localhost:8083

# Supplier Service
cd services/supplier-service && make run
# http://localhost:8084

# Procurement Service
cd services/procurement-service && make run
# http://localhost:8085
```

### 3. Test Endpoints
```bash
# Gateway health
curl http://localhost:8080/health

# Supplier list
curl http://localhost:8084/api/v1/suppliers

# Create Purchase Requisition
curl -X POST http://localhost:8085/api/v1/purchase-requisitions \
  -H "Content-Type: application/json" \
  -d '{"required_date":"2026-02-15","items":[{"material_id":"...","quantity":100}]}'
```

---

## Port Allocation

| Service | HTTP | gRPC | Database | Status |
|---------|------|------|----------|--------|
| API Gateway | 8080 | - | - | ✅ Running |
| Auth | 8081 | 9081 | auth_db | ✅ Complete |
| User | 8082 | 9082 | user_db | ✅ Complete |
| Master Data | 8083 | 9083 | master_data_db | ✅ Complete |
| **Supplier** | **8084** | **9084** | **supplier_db** | ✅ **NEW** |
| **Procurement** | **8085** | **9085** | **procurement_db** | ✅ **NEW** |
| WMS | 8086 | 9086 | wms_db | 📋 Planned |
| Manufacturing | 8087 | 9087 | manufacturing_db | 📋 Planned |
| Sales | 8088 | 9088 | sales_db | 📋 Planned |
| Marketing | 8089 | 9089 | marketing_db | 📋 Planned |
| Notification | 8090 | 9090 | notification_db | 📋 Planned |
| File | 8091 | 9091 | - | 📋 Planned |
| Reporting | 8092 | 9092 | - | 📋 Planned |

---

## Git History

| Version | Description |
|---------|-------------|
| **v0.8.0** | **Procurement Service - complete** |
| **v0.7.0** | **Supplier Service - complete** |
| v0.6.0 | API Gateway - complete |
| v0.5.0 | Master Data Service - complete |
| v0.4.0 | User Service - complete |
| v0.2.0 | Auth Service - complete |
| v0.1.0 | Infrastructure setup |

**Latest Commit**: `47e9d87`

---

## Completed Phases

| Phase | Services | Status |
|-------|----------|--------|
| **Phase 1: Core** | API Gateway, Auth, User, Master Data | ✅ Complete |
| **Phase 2: Supply Chain** | Supplier, Procurement | ✅ Complete |
| **Phase 3: Operations** | WMS, Manufacturing | 📋 Next |
| **Phase 4: Commercial** | Sales, Marketing | 📋 Planned |
| **Phase 5: Support** | Notifications, Files, Reporting | 📋 Planned |

---

## Next Steps (Phase 3)

### Ready to Implement
- [ ] WMS Service (warehouse, inventory, lots, FEFO, GRN)
- [ ] Manufacturing Service (BOM, work orders, production, QC)

### Integration Points
- WMS subscribes to `procurement.po.confirmed`
- WMS publishes `wms.grn.completed` → Procurement updates received qty
- Manufacturing subscribes to `wms.materials.available`

---

## Metrics

| Metric | Value |
|--------|-------|
| Total Files | 260+ |
| Lines of Code | ~22,000+ |
| Services | 6 complete, 9 planned |
| Database Tables | 32 active |
| API Endpoints | 65+ |
| NATS Events | 20+ defined |

---

**Repository**: https://github.com/Chinsusu/ERP  
**Updated**: 2026-01-24T06:25:00Z
