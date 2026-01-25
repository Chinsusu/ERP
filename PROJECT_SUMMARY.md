# ERP Cosmetics System - Project Summary

**Updated**: January 25, 2026  
**Repository**: https://github.com/Chinsusu/ERP  
**Status**: 10 Services + Frontend (Phase 1-5.1 Complete)

---

## Executive Summary

Complete ERP system for cosmetics manufacturing with Clean Architecture, microservices pattern, and event-driven design. Phase 5.1 Frontend now complete with Vue 3 + PrimeVue.

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
| **WMS Service** | 80 files | ✅ Running |
| **Manufacturing Service** | 54 files | ✅ Running |
| **Sales Service** | 54 files | ✅ Running |
| **Marketing Service** | 45 files | ✅ Running |
| **Frontend (Vue 3)** | 54 files | ✅ **NEW** |

**Total**: ~540+ files, ~48,000+ LOC

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

### 7. WMS Service (Port 8086) ✅ NEW - CRITICAL
Warehouse Management with FEFO logic for cosmetics industry.

| Component | Details |
|-----------|---------|
| **Tables** | 14 (warehouses, zones, locations, lots, stock, movements, reservations, GRN, goods issue, inventory counts, temperature logs) |
| **Endpoints** | 25+ (warehouses, stock, lots, GRN, goods issue, reservations, adjustments, transfers, inventory counts) |
| **FEFO Logic** | First Expired First Out - critical for cosmetics |
| **Lot Traceability** | Full tracking from supplier → warehouse → production |
| **QC Workflow** | Quarantine → QC Pass/Fail → Storage Zone |
| **Cold Storage** | Temperature monitoring (2-8°C) |
| **gRPC Methods** | CheckStockAvailability, ReserveStock, ReleaseReservation, IssueStock (FEFO), GetLotInfo, ReceiveStock |
| **Events Published** | grn.created, stock.received, stock.issued, lot.expiring_soon, stock.low |
| **Event Subscribers** | procurement.po.received, sales.order.confirmed/cancelled, manufacturing.wo.started |
| **Scheduler** | Daily expiry checks, hourly low stock alerts |
| **Unit Tests** | 24 tests (Lot, Stock, GRN, GI, Reservation workflows) |

---

## Phase 3: Operations ✅

### 8. Manufacturing Service (Port 8087) ✅
BOM, Work Orders, QC, NCR, and Traceability.

| Component | Details |
|-----------|---------|
| **Tables** | 11 (boms, bom_line_items, bom_versions, work_orders, wo_line_items, wo_material_issues, qc_checkpoints, qc_inspections, qc_inspection_items, ncrs, batch_traceability) |
| **Endpoints** | 25+ (BOM CRUD/approve, WO lifecycle, QC, NCR, Traceability) |
| **BOM Security** | AES-256-GCM encryption for formula_details |
| **WO Lifecycle** | PLANNED → RELEASED → IN_PROGRESS → QC_PENDING → COMPLETED |
| **QC Types** | IQC (Incoming), IPQC (In-Process), FQC (Final) |
| **Traceability** | Forward (material→products) and Backward (product→materials) |
| **Events** | bom.created/approved, wo.created/started/completed, qc.passed/failed, ncr.created |

---

## Phase 4: Commercial ✅ NEW

### 9. Sales Service (Port 8088) ✅ NEW
Customer management, quotations, sales orders, shipments.

| Component | Details |
|-----------|---------|
| **Tables** | 11 (customer_groups, customers, contacts, addresses, quotations, quotation_items, sales_orders, so_items, shipments, returns, return_items) |
| **Endpoints** | 35 (customers CRUD, quotations, sales orders, shipments) |
| **Credit Control** | Automatic credit limit check on order confirmation |
| **SO Lifecycle** | DRAFT → CONFIRMED → SHIPPED → DELIVERED |
| **Quotation Workflow** | Convert quotation to sales order with one click |
| **Events** | customer.created, order.confirmed/cancelled, shipment.shipped |

### 10. Marketing Service (Port 8089) ✅ NEW
KOL/Influencer database, campaigns, sample distribution.

| Component | Details |
|-----------|---------|
| **Tables** | 8 (kol_tiers, kols, campaigns, collaborations, sample_requests, sample_items, sample_shipments, kol_posts) |
| **Endpoints** | 24 (KOLs CRUD, campaigns, sample requests, approvals) |
| **KOL Tiers** | MEGA (>1M), MACRO (100K-1M), MICRO (10K-100K), NANO (<10K) |
| **Sample Workflow** | DRAFT → APPROVED → SHIPPED → DELIVERED → FEEDBACK |
| **Campaign ROI** | Budget, spend, impressions, engagement, conversions tracking |
| **Events** | campaign.created/launched, sample.approved/shipped, kol_post.recorded |

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
| **Supplier** | **8084** | **9084** | **supplier_db** | ✅ Running |
| **Procurement** | **8085** | **9085** | **procurement_db** | ✅ Running |
| **WMS** | **8086** | **9086** | **wms_db** | ✅ Running |
| **Manufacturing** | **8087** | **9087** | **manufacturing_db** | ✅ Running |
| **Sales** | **8088** | **9088** | **sales_db** | ✅ **NEW** |
| **Marketing** | **8089** | **9089** | **marketing_db** | ✅ **NEW** |
| Notification | 8090 | 9090 | notification_db | 📋 Planned |
| File | 8091 | 9091 | - | 📋 Planned |
| Reporting | 8092 | 9092 | - | 📋 Planned |

---

## Git History

| Version | Description |
|---------|-------------|
| **v0.11.0** | **Marketing Service - complete (Phase 4.2)** |
| **v0.10.0** | **Sales Service - complete (Phase 4.1)** |
| v0.9.0 | Manufacturing Service - complete |
| v0.8.0 | WMS Service - complete (CRITICAL) |
| v0.7.0 | Procurement Service - complete |
| v0.6.0 | Supplier Service - complete |
| v0.5.0 | API Gateway - complete |
| v0.4.0 | Master Data Service - complete |
| v0.3.0 | User Service - complete |
| v0.2.0 | Auth Service - complete |
| v0.1.0 | Infrastructure setup |

**Latest Commit**: `de49d0a`

---

## Completed Phases

| Phase | Services | Status |
|-------|----------|--------|
| **Phase 1: Core** | API Gateway, Auth, User, Master Data | ✅ Complete |
| **Phase 2: Supply Chain** | Supplier, Procurement, WMS | ✅ Complete |
| **Phase 3: Operations** | Manufacturing | ✅ Complete |
| **Phase 4: Commercial** | Sales, Marketing | ✅ Complete |
| **Phase 5.1: Frontend** | Vue 3 + PrimeVue | ✅ **Complete** |
| **Phase 5.2: CRUD Pages** | DataTable, Forms | 📋 Planned |

---

## Next Steps (Phase 5)

### Ready to Implement
- [ ] Notification Service (email, SMS, push notifications)
- [ ] File Service (document upload, MinIO integration)
- [ ] Reporting Service (dashboards, exports)

### Integration Points
- Notifications subscribe to all service events
- File service used by all services for document attachments

---

## Metrics

| Metric | Value |
|--------|-------|
| Total Files | 540+ |
| Lines of Code | ~48,000+ |
| Backend Services | 10 complete |
| Frontend | Vue 3 + PrimeVue |
| Database Tables | 68 active |
| API Endpoints | 150+ |
| NATS Events | 50+ defined |
| Unit Tests | 24+ |

---

**Repository**: https://github.com/Chinsusu/ERP  
**Updated**: 2026-01-25T07:35:00Z
