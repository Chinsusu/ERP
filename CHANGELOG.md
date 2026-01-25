# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added - Phase 11: Monitoring & Maintenance - 2026-01-25

**Complete Monitoring & Alerting Setup (~1,100 LOC)**

**Monitoring & Alerting**:
- **Alert Rules** (261 lines): erp-alerts.yml
  - 10+ technical alerts (ServiceDown, HighErrorRate, DatabaseConnectionHigh, Memory/Disk)
  - 10+ business alerts (LowStock, LotsExpiring, CertificateExpiring, FailedWorkOrders)
  - Severity-based routing (critical vs warning)
- **AlertManager Config** (127 lines): Email notifications with routing
  - Critical alerts → IT Lead + DevOps (1h repeat)
  - Warning alerts → IT Team (4h repeat)
  - Business alerts → Operations (business hours only)
  - Inhibition rules to prevent alert storms
- **Grafana Dashboard** (326 lines): ERP Overview dashboard
  - 10 panels: Service Health, Request Rate, Response Time, Error Rate
  - Database Connections, CPU/Memory Usage by Service
  - Business metrics: Low Stock Alerts, Expiring Lots
- **Maintenance Guide** (343 lines): Comprehensive operations manual
  - Daily checklists (AM/PM), Weekly/Monthly tasks
  - Incident Management (P1/P2/P3 severity levels)
  - Troubleshooting guide (common issues + solutions)
  - Backup verification procedures

**Total Lines**: 1,057 lines of monitoring configuration and documentation

---

### Added - Phase 9-10: Documentation & Deployment - 2026-01-25

**Complete Documentation & Deployment Package (~1,700 LOC)**

**Phase 9: Documentation & Training**
- **User Manual (Vietnamese)**: Complete manual covering all 12 modules (274 lines)
  - Login, Materials (INCI/CAS), Suppliers (GMP), Procurement (PR/PO)
  - WMS (FEFO), Manufacturing (BOM Security, Traceability)
  - Sales, Marketing, Reports, Admin, FAQ
- **Training Guide**: 6 training sessions with hands-on labs (215 lines)
  - WMS FEFO exercises, Manufacturing BOM, Procurement workflows
  - Assessment questions and answers
- **OpenAPI Specification**: Full API documentation (743 lines)
  - 40+ endpoints covering Auth, Users, Materials, Suppliers
  - Procurement, WMS, Manufacturing, Sales, Reports
  - Complete schemas and permission requirements

**Phase 10: Deployment & Go-Live**
- **Go-Live Checklist**: Pre-deployment verification (92 lines)
  - Infrastructure, Application, Data, People sections
  - Go/No-Go decision matrix
- **Go-Live Runbook**: Detailed timeline and procedures (181 lines)
  - 6-phase timeline (6AM-12PM Saturday)
  - Rollback procedure and decision matrix
  - Post go-live monitoring checklist
- **deploy.sh**: Rolling deployment automation (91 lines)
  - Pre-deploy backup, health checks
  - Service-by-service deployment
- **setup-ssl.sh**: SSL certificate automation (100 lines)
  - Let's Encrypt integration
  - Auto-renewal configuration

**Total Lines**: 1,696 lines of documentation and scripts

---

### Added - File Service (Phase 6.3) - 2026-01-25

**Complete File Service Implementation (~20 files, ~1,200 LOC)**

**Database Layer (6 migrations)**
- `files` - File metadata and storage references
- `file_categories` - Category definitions with validation rules
- Seed: 10 file categories (DOCUMENT, IMAGE, CERTIFICATE, etc.)

**Domain Layer**
- `File` - File entity with metadata and helper methods
- `FileCategory` - Category with validation (extensions, size limits)

**Infrastructure Layer**
- `MinIOClient` - S3-compatible storage (upload, download, presigned URLs)
- PostgreSQL repositories for files and categories

**API Endpoints (8 total)**
- `POST /upload` - Single file upload
- `POST /upload/multiple` - Multiple file upload
- `GET /:id/download` - Download file
- `GET /:id/url` - Get presigned URL
- `GET /entity/:type/:id` - Files by entity
- `GET /categories` - List categories

**Port**: 8091 (HTTP)

---

### Added - Reporting Service (Phase 6.2) - 2026-01-25

**Complete Reporting Service Implementation (~35 files, ~2,000 LOC)**

**Database Layer (10 migrations)**
- `report_definitions` - Report templates with SQL queries and parameters
- `report_executions` - Execution tracking with status and results
- `dashboards` - User dashboards with layout configuration
- `widgets` - Dashboard widgets (KPI, Charts, Tables)
- Seed: 10 default reports, 1 main dashboard with 8 widgets

**Domain Layer (4 entities + 4 repository interfaces)**
- `ReportDefinition` - Report templates with columns
- `ReportExecution` - Execution tracking
- `Dashboard` - Dashboard configuration
- `Widget` - Widget configuration

**Infrastructure Layer**
- `csv_exporter.go` - CSV export
- `excel_exporter.go` - Excel (XLSX) export with formatting
- `stats_aggregator.go` - KPI aggregation from services

**Use Cases (3)**
- `DashboardUseCase` - Dashboard + widget CRUD
- `ReportUseCase` - Execute and export reports
- `StatsUseCase` - Real-time KPIs

**API Endpoints (20 total)**
- Dashboards: CRUD + widgets management
- Reports: List, execute, download
- Stats: Inventory/Sales/Production/Procurement KPIs

**Pre-built Reports**: STOCK_SUMMARY, EXPIRY_REPORT, STOCK_MOVEMENT, LOW_STOCK_ITEMS, PO_SUMMARY, SUPPLIER_PERFORMANCE, PRODUCTION_OUTPUT, QC_SUMMARY, SALES_SUMMARY, TOP_PRODUCTS

**Port**: 8092 (HTTP)

---

### Added - Notification Service (Phase 6.1) - 2026-01-25

**Complete Notification Service Implementation (~40 files, ~2,500 LOC)**

**Database Layer (12 migrations)**
- `notification_templates` - Email/in-app notification templates with variables
- `notifications` - Outbound notification queue (email, SMS, push)
- `user_notifications` - In-app notifications with read/dismiss tracking
- `alert_rules` - Configurable automated alert rules
- `email_logs` - Email delivery tracking
- Seed data: 9 default templates, 6 default alert rules

**Domain Layer (5 entities + 4 repository interfaces)**
- `Notification` - With status tracking and retry logic
- `NotificationTemplate` - With template rendering
- `UserNotification` - In-app notifications
- `AlertRule` - Configurable alerts
- `EmailLog` - Email tracking

**Infrastructure Layer**
- `smtp_sender.go` - SMTP email sender with TLS support + mock sender
- `publisher.go` - NATS event publisher
- `subscriber.go` - Event subscriptions from 9 service events

**Use Case Layer (4 use cases)**
- `SendNotification` - Send email/in-app/both with template rendering
- `TemplateUseCase` - Template CRUD
- `UserNotificationUseCase` - In-app notification management
- `AlertRuleUseCase` - Alert rule management

**HTTP Handlers (4 handlers + router)**
- `NotificationHandler` - Send, template CRUD
- `UserNotificationHandler` - List, read, delete
- `AlertRuleHandler` - CRUD + activate/deactivate
- `HealthHandler` - Health/ready checks

**API Endpoints (17 total)**
- `POST /api/v1/notifications/send`
- Templates: GET/POST/PUT/DELETE `/api/v1/notifications/templates`
- In-App: GET/POST/PATCH/DELETE `/api/v1/notifications/in-app`
- Alert Rules: GET/POST/PUT/DELETE `/api/v1/alert-rules`

**Event Subscriptions**
- `wms.stock.low_stock_alert` → Low stock notification
- `wms.lot.expiring_soon` → Lot expiry notification
- `supplier.certification.expiring` → Certificate expiry notification
- `procurement.pr.submitted` → PR approval notification
- `procurement.po.created` → PO notification
- `manufacturing.qc.failed` → QC failed notification
- `sales.order.confirmed` → Order confirmation

**Default Alert Rules**
- Stock below reorder point
- Lot expiring in 30/7 days
- Certificate expiring in 90/30 days
- Approval pending > 24 hours

**Configuration**
- Port: 8090 (HTTP)
- SMTP configuration via environment variables
- Background email processor with retry logic

---

### Added - Frontend Business Modules (Phase 5.2) - 2026-01-25

**Business Module Pages (~15 files, ~2,500 LOC)**
- Materials: List, Detail (tabs with suppliers), Form (create/edit)
- Suppliers: List (with ratings), Detail (contacts, certifications, evaluations), Form
- Procurement: Purchase Requisitions list, Purchase Orders list
- WMS: Stock overview (dashboard with stats), GRN list
- Manufacturing: BOM list, Work Orders (with progress bars)

**API Layer (~5 files, ~400 LOC)**
- `material.api.ts` - Materials CRUD with filters
- `supplier.api.ts` - Suppliers with contacts, certifications, evaluations
- `procurement.api.ts` - PR and PO operations
- `wms.api.ts` - Stock, lots, GRN, movements
- `manufacturing.api.ts` - BOM, work orders, QC, traceability

**TypeScript Types**
- `business.types.ts` - 25+ interfaces (Material, Supplier, PR, PO, Stock, BOM, WorkOrder, etc.)

**Composables**
- `useApi.ts` - TanStack Query wrappers (useApiQuery, usePaginatedQuery, useApiMutation)

**Mock Data System**
- `src/mocks/data.ts` - 50+ mock records for testing
- `src/mocks/handlers.ts` - Axios interceptors for mock API
- Auto-enabled in DEV mode (set VITE_USE_MOCKS=false to disable)

**Common Features**
- Server-side pagination with DataTable
- Filter bars with search, dropdowns
- Status tags with color severity
- Action buttons (view, edit)
- Responsive card layouts

**Build Output**
- Production build: 5.05s, 30+ assets
- Total frontend files: 55+

---

### Added - Frontend Setup (Phase 5.1) - VUE 3 + PRIMEVUE

**Implementation Complete (~35 files, ~3,000 LOC)**
- Vue 3.4+ with Composition API and TypeScript
- PrimeVue 4 with Aura theme (dark mode support)
- Pinia for state management
- Vue Router 4 with auth guards
- Axios with JWT token management and auto-refresh
- TanStack Query ready for data fetching

**Project Structure**
- `src/api/` - Axios instance with interceptors, auth API, user API
- `src/composables/` - useAuth, usePagination
- `src/stores/` - auth.store (RBAC permissions), app.store (dark mode)
- `src/router/` - Routes with lazy loading and permission guards
- `src/layouts/` - DefaultLayout, AuthLayout
- `src/components/layout/` - AppHeader, AppSidebar with gradient design

**Core Features**
- JWT authentication with token refresh
- RBAC permission checking (service:resource:action format)
- Dark mode toggle with localStorage persistence
- Collapsible sidebar navigation
- Permission-based menu filtering
- Responsive design

**Pages**
- LoginPage with gradient branding
- ForgotPasswordPage with success state
- DashboardPage with stats cards
- Placeholder pages for all modules

**Build Output**
- Production build: 29 assets (~445KB gzip main bundle)
- Dev server: http://localhost:5173

---

### Added - Marketing Service (Phase 4.2) - COMMERCIAL

**Implementation Complete (~40 files, ~4,000 LOC)**
- KOL/Influencer database with tier classification (MEGA/MACRO/MICRO/NANO)
- Campaign management with budget and ROI tracking
- Sample request workflow (DRAFT → APPROVED → SHIPPED → DELIVERED)
- KOL post tracking with engagement metrics
- Sample shipment tracking

**Database Layer (18 migration files, 8 tables)**
- `kol_tiers` - Tier classification
- `kols` - KOL master with social media handles
- `campaigns` - Marketing campaigns with performance metrics
- `kol_collaborations` - KOL-Campaign relationships
- `sample_requests`, `sample_items`, `sample_shipments` - Sample workflow
- `kol_posts` - Content tracking

**API Endpoints (24 total)**
- KOL Tiers: GET /kol-tiers
- KOLs: GET/POST/PUT/DELETE /kols, GET /:id/posts
- Campaigns: GET/POST/PUT /campaigns, PATCH /:id/launch, GET /:id/performance
- Samples: GET/POST /samples/requests, PATCH /:id/approve, PATCH /:id/ship

**Events Published**
- `marketing.campaign.created`, `marketing.campaign.launched`
- `marketing.sample_request.created`, `marketing.sample_request.approved`
- `marketing.sample.shipped`, `marketing.kol_post.recorded`

---

### Added - Sales Service (Phase 4.1) - COMMERCIAL


**Implementation Complete (~50 files, ~4,500 LOC)**
- First service in Phase 4: Commercial Operations
- Customer management with credit limit control
- Quotation workflow with convert-to-order
- Sales Order lifecycle (DRAFT → CONFIRMED → SHIPPED → DELIVERED)
- Shipment tracking with carrier integration
- Return processing workflow

**Database Layer (22 migration files, 11 tables)**
- `customer_groups` - Pricing tiers (VIP, Gold, Silver, etc.)
- `customers` - Customer master with credit_limit, current_balance
- `customer_contacts`, `customer_addresses` - CRM data
- `quotations`, `quotation_line_items` - Sales quotations
- `sales_orders`, `so_line_items` - Orders with shipped_quantity tracking
- `shipments` - Delivery tracking
- `returns`, `return_line_items` - Sales returns

**Sales Order Lifecycle**
```
DRAFT → CONFIRMED → PROCESSING → PARTIALLY_SHIPPED → SHIPPED → DELIVERED
                                                           ↓
                                                     CANCELLED
```

**Credit Limit Control**
```go
// Automatic credit check on order confirmation
if !customer.CanPlaceOrder(orderAmount) {
    return ErrInsufficientCredit
}
// current_balance updated on confirm, released on cancel
```

**API Endpoints (25+ total)**
- Customers: POST/GET/PUT/DELETE /customers, GET /:id/credit-check
- Quotations: POST/GET /quotations, PATCH /:id/send, POST /:id/convert-to-order
- Sales Orders: POST/GET /sales-orders, PATCH /:id/confirm|ship|deliver|cancel
- Shipments: POST/GET /shipments, PATCH /:id/ship|deliver

**Events Published**
- `sales.customer.created`
- `sales.quotation.sent`
- `sales.order.created`, `sales.order.confirmed` (→ WMS reserves stock)
- `sales.order.shipped`, `sales.order.delivered`
- `sales.order.cancelled` (→ WMS releases reservations)
- `sales.return.created`

**Integration Points**
- WMS subscribes to `sales.order.confirmed` for stock reservation
- WMS subscribes to `sales.order.cancelled` to release reservations
- Manufacturing triggered via `sales.order.confirmed` for MTO items

---
### Added - WMS Extended Features
- **Goods Issue with FEFO**: Issue stock automatically using First Expired First Out logic
- **Stock Reservations**: Reserve stock for sales orders, work orders
- **Stock Adjustments**: Cycle count, damage, expiry corrections
- **Stock Transfers**: Move stock between locations
- **Scheduler**: Daily expiry checks, hourly low stock alerts
- **gRPC Proto**: Ready for inter-service communication
- **Inventory Count**: Full physical inventory workflow (Draft → In Progress → Completed)
- **gRPC Server**: Complete gRPC server implementation, starts on port 9091
- **Unit Tests**: 24 tests covering Lot FEFO, Stock operations, GRN/GI workflows
- **Event Subscribers**: React to events from other services automatically

**Event Handlers**:
- `procurement.po.received` → Auto-create GRN
- `sales.order.confirmed` → Reserve stock for sales orders
- `sales.order.cancelled` → Release reservations
- `manufacturing.wo.started` → Reserve materials for production

**New Endpoints**:
- POST /api/v1/goods-issue - Create goods issue (FEFO)
- POST /api/v1/reservations - Reserve stock
- DELETE /api/v1/reservations/:id - Release reservation
- POST /api/v1/adjustments - Create stock adjustment
- POST /api/v1/transfers - Transfer stock between locations
- GET /api/v1/stock/availability/:material_id - Check availability
- POST /api/v1/inventory-counts - Create inventory count
- GET /api/v1/inventory-counts - List inventory counts
- GET /api/v1/inventory-counts/:id - Get inventory count
- PATCH /api/v1/inventory-counts/:id/start - Start counting
- POST /api/v1/inventory-counts/:id/record - Record count
- PATCH /api/v1/inventory-counts/:id/complete - Complete and apply variance

**gRPC Methods**:
- CheckStockAvailability, ReserveStock, ReleaseReservation
- IssueStock (FEFO), GetLotInfo, GetLotsByMaterial, ReceiveStock

---

### Added - Manufacturing Service (Phase 3.1) - OPERATIONS

**Implementation Complete (~54 files, ~4,700 LOC)**
- First service in Phase 3: Operations
- BOM management with AES-256-GCM formula encryption
- Work Order lifecycle management
- QC (IQC/IPQC/FQC) and NCR tracking
- Full lot traceability (forward/backward)

**Database Layer (24 migration files, 11 tables)**
- `boms`, `bom_line_items`, `bom_versions` - BOM with encrypted formula
- `work_orders`, `wo_line_items`, `wo_material_issues` - Production
- `qc_checkpoints`, `qc_inspections`, `qc_inspection_items` - Quality Control
- `ncrs` - Non-Conformance Reports
- `batch_traceability` - Material → Product lot mapping

**BOM Security (Critical Feature)**
```go
// AES-256-GCM encryption for formula_details
formula_details BYTEA  -- encrypted JSON
confidentiality_level: PUBLIC | INTERNAL | CONFIDENTIAL | RESTRICTED
```

**Work Order Lifecycle**
```
PLANNED → RELEASED → IN_PROGRESS → QC_PENDING → COMPLETED
                                        ↓
                                    CANCELLED
```

**API Endpoints (25+ total)**
- BOM: POST/GET /boms, GET /boms/:id, POST /boms/:id/approve
- Work Orders: POST/GET /work-orders, PATCH /:id/release|start|complete
- QC: GET /qc-checkpoints, POST/GET /qc-inspections, PATCH /:id/approve
- NCR: POST/GET /ncrs, PATCH /:id/close
- Traceability: GET /traceability/backward/:lot_id, /forward/:lot_id

**Events Published**
- `manufacturing.bom.created`, `manufacturing.bom.approved`
- `manufacturing.wo.created`, `manufacturing.wo.released`, `manufacturing.wo.started`, `manufacturing.wo.completed`
- `manufacturing.qc.passed`, `manufacturing.qc.failed`
- `manufacturing.ncr.created`

**Cosmetics-Specific Features**
- Formula confidentiality (trade secret protection)
- GMP-compliant QC checkpoints (IQC, IPQC, FQC)
- Batch/Lot traceability for regulatory compliance

---

## [0.9.0] - 2026-01-24

### Added - WMS Service Complete (Phase 2.3) - CRITICAL SERVICE

**Implementation Complete (~65 files, ~4,500 LOC)**
- The most critical service for cosmetics ERP
- FEFO (First Expired First Out) logic for cosmetics industry
- Full lot traceability from supplier → warehouse → production
- Cold storage monitoring support (2-8°C)

**Database Layer (30 migration files)**
- 14 tables: warehouses, zones, locations, lots, stock, stock_movements, stock_reservations, grns, grn_line_items, goods_issues, gi_line_items, stock_adjustments, inventory_counts, temperature_logs
- Comprehensive seed data: 3 warehouses, 9 zones, 15 locations
- Indexes optimized for FEFO queries (expiry_date ASC)

**Domain Layer (9 entities)**
- Lot entity with expiry tracking: DaysUntilExpiry(), IsExpired(), IsExpiringSoon()
- Stock entity with availability: GetAvailableQuantity(), Reserve(), Issue()
- FEFO support: LotIssued struct tracks which lots were used
- GRN workflow: DRAFT → QC → COMPLETED with quarantine zone
- Zone types: RECEIVING, QUARANTINE, STORAGE, COLD, PICKING, SHIPPING

**FEFO Logic (Core Feature)**
```go
// Issues stock from earliest expiring lots first
func GetAvailableStockFEFO(materialID) → sorted by lots.expiry_date ASC
func IssueStockFEFO(materialID, quantity) → []LotIssued
```

**API Endpoints (20+ total)**
- GET /api/v1/warehouses - List warehouses with zones
- GET /api/v1/warehouses/:id/zones - Get zones in warehouse
- GET /api/v1/zones/:id/locations - Get locations in zone
- GET /api/v1/stock - Query stock with filters
- GET /api/v1/stock/by-material/:id - Stock by material with summary
- GET /api/v1/stock/expiring?days=90 - **Expiring stock (FEFO)**
- GET /api/v1/stock/low-stock - Low stock alerts
- GET /api/v1/lots - List lots with status/QC filters
- GET /api/v1/lots/:id - Get lot details
- GET /api/v1/lots/:id/movements - **Lot traceability**
- POST /api/v1/grn - Create GRN (from PO)
- GET /api/v1/grn - List GRNs
- GET /api/v1/grn/:id - Get GRN with line items
- PATCH /api/v1/grn/:id/complete - Complete GRN after QC

**Events Published**
- wms.grn.created, wms.grn.completed
- wms.stock.received, wms.stock.issued
- wms.stock.reserved, wms.stock.low_stock_alert
- wms.lot.expiring_soon, wms.lot.expired

**Cosmetics-Specific Features**
- FEFO: First Expired First Out (not FIFO)
- Lot Traceability: Track materials from supplier to production
- Quarantine Zone: All goods pass QC before storage
- Cold Storage: Temperature monitoring for sensitive materials
- Expiry Alerts: 90, 30, 7 days configurable thresholds

### Ports
- HTTP: 8086
- gRPC: 9086 (planned)
- Database: wms_db

---

## [0.8.0] - 2026-01-24

### Added - Procurement Service Complete (Phase 2.2)

**Implementation Complete (~35 files, ~2,500 LOC)**
- Full microservice for Purchase Requisitions and Purchase Orders
- Multi-level approval workflow based on amount thresholds
- PR to PO conversion with line item transfer

**Database Layer (14 migration files)**
- 7 tables: purchase_requisitions, pr_line_items, pr_approvals, purchase_orders, po_line_items, po_amendments, po_receipts
- Indexes for performance optimization
- Foreign key relationships for data integrity

**Domain Layer**
- PurchaseRequisition entity with workflow methods (Submit, Approve, Reject)
- PurchaseOrder entity with workflow methods (Confirm, Cancel, Close)
- PRLineItem, PRApproval, POLineItem, POAmendment, POReceipt entities
- Multi-level approval logic (Auto < 10M, Dept Manager < 50M, Procurement < 200M, CFO > 200M VND)

**API Endpoints (13 total)**
- POST/GET /api/v1/purchase-requisitions - Create and list PRs
- GET /api/v1/purchase-requisitions/:id - Get PR details
- POST /api/v1/purchase-requisitions/:id/submit - Submit for approval
- POST /api/v1/purchase-requisitions/:id/approve - Approve PR
- POST /api/v1/purchase-requisitions/:id/reject - Reject PR
- POST /api/v1/purchase-requisitions/:id/convert-to-po - Convert to PO
- GET /api/v1/purchase-orders - List POs
- GET /api/v1/purchase-orders/:id - Get PO details
- POST /api/v1/purchase-orders/:id/confirm - Confirm PO
- POST /api/v1/purchase-orders/:id/cancel - Cancel PO
- POST /api/v1/purchase-orders/:id/close - Close PO
- GET /api/v1/purchase-orders/:id/receipts - Get PO receipts

**Events Published**
- procurement.pr.created, procurement.pr.submitted
- procurement.pr.approved, procurement.pr.rejected
- procurement.po.created, procurement.po.confirmed
- procurement.po.received, procurement.po.closed

### Ports
- HTTP: 8085
- gRPC: 9085 (planned)

---

## [0.7.0] - 2026-01-24

### Added - Supplier Service Complete (Phase 2.1)

**Implementation Complete (~40 files, ~2,800 LOC)**
- Full microservice for supplier management
- Cosmetics industry-specific certification tracking
- Approved Supplier List (ASL) management

**Database Layer (16 migration files)**
- 7 tables: suppliers, supplier_addresses, supplier_contacts, supplier_certifications, supplier_evaluations, approved_supplier_list, supplier_price_lists
- Seed data: 4 sample suppliers with addresses, contacts, certifications
- GMP, ISO, ORGANIC, ECOCERT, HALAL, COSMOS certification types

**Domain Layer**
- Supplier entity with status workflow (PENDING → APPROVED → BLOCKED)
- Certification with expiry tracking (VALID, EXPIRING_SOON, EXPIRED)
- Evaluation with 5 categories and weighted rating calculation
- HasValidGMP() business method for compliance checking

**API Endpoints (16 total)**
- POST/GET /api/v1/suppliers - Create and list suppliers
- GET /api/v1/suppliers/:id - Get supplier details
- PUT /api/v1/suppliers/:id - Update supplier
- POST /api/v1/suppliers/:id/approve - Approve supplier
- POST /api/v1/suppliers/:id/block - Block supplier
- GET/POST /api/v1/suppliers/:id/addresses - Manage addresses
- GET/POST /api/v1/suppliers/:id/contacts - Manage contacts
- GET/POST /api/v1/suppliers/:id/certifications - Manage certifications
- GET /api/v1/certifications/expiring - Get expiring certs
- GET/POST /api/v1/suppliers/:id/evaluations - Manage evaluations

**Events Published**
- supplier.created, supplier.approved, supplier.blocked
- supplier.certification.added, supplier.certification.expiring
- supplier.certification.expired, supplier.evaluation.completed

### Ports
- HTTP: 8084
- gRPC: 9084 (planned)

---

## [0.6.0] - 2026-01-24

### Added - API Gateway Complete

**Implementation Complete (17 files)**
- Single entry point for all ERP services on port 8080
- Dynamic routing to 15 backend services

**Middleware Chain**
- RequestID: UUID generation and propagation
- CORS: Cross-origin request handling
- Logger: Structured request/response logging
- Recovery: Panic handling
- RateLimiter: Redis-based sliding window (100/min user, 30/min IP)
- Auth: JWT validation with blacklist check
- CircuitBreaker: Fault tolerance (5 failures → 30s open)

**Proxy Layer**
- Reverse proxy with header enrichment
- Service registry with health checking
- X-User-ID, X-Request-ID header injection

**Routes Configured**
- Auth, User, Master Data, Supplier, Procurement
- WMS, Manufacturing, Sales, Marketing
- Notifications, Files, Reports

**Health Endpoints**
- GET /health - Aggregate service health
- GET /ready, /live - Kubernetes probes
- GET /health/:service - Individual service status

### Ports
- HTTP: 8080

## [0.5.0] - 2026-01-24

### Added - Master Data Service Complete

**Implementation Complete (47 files, ~3,000 LOC)**
- Full microservice for materials, products, categories, and units of measure
- Complete cosmetics industry-specific features

**Database Layer (16 migration files)**
- 7 tables: categories, units_of_measure, unit_conversions, materials, material_specifications, products, product_images
- Hierarchical categories using materialized path pattern
- Unit conversion support (bidirectional)
- Seed data: 12 units (KG, G, MG, L, ML, PCS, etc.)
- Seed data: 8 unit conversions
- Seed data: 27 categories (material categories + product categories)

**Domain Layer (4 entity files, 4 repository interfaces)**
- Category entity with tree structure support
- UnitOfMeasure and UnitConversion entities
- Material entity with cosmetics-specific fields
- Product entity with license tracking
- MaterialSpecification for extended attributes
- ProductImage for multi-image support

**Cosmetics Industry Features**
- INCI Names: International Nomenclature of Cosmetic Ingredients
- CAS Numbers: Chemical Abstracts Service registry
- Allergen Tracking: is_allergen, allergen_info fields
- Storage Conditions: AMBIENT, COLD (2-8°C), FROZEN
- Cosmetic License: Number and expiry date tracking
- Auto-generated Codes: RM-0001, PKG-0001, FG-SERUM-0001

**Infrastructure Layer (5 files)**
- CategoryRepository with tree loading
- UnitRepository with conversion logic
- MaterialRepository with full-text search
- ProductRepository with image management
- Event publisher for NATS

**Use Case Layer (4 use case packages)**
- Category: CRUD, GetTree
- Unit: CRUD, Convert
- Material: CRUD, Search, AddSpecification
- Product: CRUD, Search, AddImage

**HTTP API (20+ endpoints)**
- Categories: CRUD + /tree endpoint
- Units: CRUD + /convert endpoint
- Materials: CRUD + /search + /specifications
- Products: CRUD + /search + /by-category + /images

**Configuration & Deployment**
- Makefile with build, run, migrate targets
- Dockerfile for production builds
- Comprehensive README.md with examples

### Ports
- HTTP: 8083
- gRPC: 9083 (placeholder)
- Database: master_data_db

## [0.2.0] - 2026-01-24

### Added - Auth Service Template (Complete & Tested)

**Core Implementation (50 files, ~3,500 LOC)**
- Complete authentication microservice serving as template for all 15 services
- Clean Architecture: Domain → UseCase → Delivery → Infrastructure layers

**Database Layer (18 files)**
- 9 database migrations with up/down SQL scripts
- 8 tables: roles, permissions, role_permissions, user_credentials, user_roles, refresh_tokens, sessions, password_reset_tokens
- Seed data: 5 default roles (Super Admin, Admin, Manager, Staff, Viewer)
- Seed data: 42 default permissions covering 8 services
- Default admin user: admin@company.vn (bcrypt hashed)

**Security Features**
- JWT tokens: Access (15min) + Refresh (7 days) with token rotation
- Password security: Bcrypt hashing (cost 12)
- Account protection: Lockout after 5 failed attempts (30 minutes)
- RBAC: Permission format service:resource:action with wildcard support

**Testing Results**
- ✅ Health check, Login, Refresh token, Logout - All PASS
- ✅ Database: 8 tables migrated, seed data loaded
- ✅ Infrastructure: PostgreSQL, Redis, NATS all connected

### Fixed
- Viper configuration: Added explicit environment variable binding
- Logger usage: Fixed zap.Field usage in main.go
- Database config: Fixed NowFunc field name in GORM config


### Infrastructure Implementation - 2026-01-23

**Added:**
- **Docker Compose Infrastructure**
  - docker-compose.yml - Production environment with all services
  - docker-compose.dev.yml - Development environment with hot reload
  - PostgreSQL 16 with 13 service-specific databases
  - Redis 7 for caching and sessions
  - NATS with JetStream for event-driven architecture
  - MinIO for S3-compatible object storage
  - Monitoring stack: Prometheus, Grafana, Loki, Jaeger
  - Nginx reverse proxy with rate limiting

- **Shared Go Libraries** (`shared/pkg/`)
  - `config` - Viper-based configuration management with environment variables
  - `database` - GORM PostgreSQL connection with pooling and transactions
  - `logger` - Zap structured logging (JSON/console formats)
  - `jwt` - JWT token generation and verification (access/refresh tokens)
  - `middleware` - HTTP middlewares (CORS, auth, logging, recovery, request ID)
  - `errors` - Custom error types with HTTP status code mapping
  - `validator` - Struct validation with user-friendly error messages
  - `response` - Standard HTTP response helpers with pagination
  - `nats` - NATS JetStream client for pub/sub messaging
  - `grpc` - gRPC server/client with logging and recovery interceptors

- **Protobuf Definitions** (`shared/proto/`)
  - common.proto - Shared message types for cross-service communication

- **Deployment Configurations** (`deploy/`)
  - Nginx: Reverse proxy with rate limiting and SSL/TLS ready
  - Prometheus: Metrics scraping for all 15 microservices
  - Loki: Log aggregation configuration
  - Grafana: Auto-provisioned datasources (Prometheus, Loki, Jaeger)

- **Utility Scripts** (`scripts/`)
  - init-databases.sql - PostgreSQL database initialization for all services
  - health-check.sh - Comprehensive health monitoring for all infrastructure
  - backup.sh - Automated backup for PostgreSQL, Redis, MinIO, and configs

- **Developer Tools**
  - Root Makefile with 20+ automation targets
  - .env.example with 93 documented environment variables
  - shared/README.md with library usage documentation

### Documentation
- Created comprehensive walkthrough for infrastructure setup
- Documented shared library APIs with code examples
- SSL certificate setup instructions

### Added
- Initial project documentation structure
- README.md with project overview and quick start guide
- 01-ARCHITECTURE.md with comprehensive system architecture documentation
  - Microservices architecture overview
  - Service communication patterns (REST, gRPC, Event-driven)
  - Database per service strategy
  - API Gateway design
  - Security architecture (JWT + RBAC)
  - Docker Compose deployment architecture
  - Monitoring and observability strategy
  - Offline capability design
  - Disaster recovery plan

### Documentation
- Created comprehensive architecture diagrams using Mermaid
- Documented all 15 microservices with port allocations
- Defined security patterns and encryption strategies
- Specified resource allocation for Docker containers

---

## [1.0.0] - 2026-01-23

### Initial Release
- Project initialization
- Documentation framework setup

## [Unreleased] - 2026-01-23

### Added
- 02-SERVICE-SPECIFICATIONS.md - Comprehensive specifications for all 15 microservices
  - API Gateway routing and middleware configuration
  - Auth Service with JWT and RBAC implementation
  - User Service with department management
  - Master Data Service for materials and products (with INCI, CAS numbers)
  - Supplier Service with certification tracking (GMP, ISO, Organic)
  - Procurement Service with PR/PO workflows
  - WMS Service with lot management and FEFO logic
  - Manufacturing Service with BOM versioning, work orders, and QC
  - Sales Service with customer and order management
  - Summary of remaining services (Marketing, Finance, Reporting, Notification, AI, File)

### Documentation
- Detailed database schemas for each service
- REST API endpoints with HTTP methods
- gRPC internal communication methods
- Event publishing/subscription patterns
- Permission requirements per service
- Service dependency mapping with Mermaid diagram

### Phase 2 Completed - 2026-01-23

**Added:**
- 03-AUTH-SERVICE.md - Complete authentication & authorization documentation
  - JWT access & refresh token flows with sequence diagrams
  - RBAC permission system (format: service:resource:action)
  - User-role and role-permission management
  - Account security (bcrypt, lockout, token rotation)
  - Redis caching strategy for permissions
  - Comprehensive API endpoints and gRPC methods
  
- 04-USER-SERVICE.md - User & department management
  - Employee information management
  - Department hierarchy (tree structure)
  - User preferences and avatar handling
  - Employment status tracking
  - Emergency contact management
  
- 05-MASTER-DATA-SERVICE.md - Materials & products master data
  - Material management with INCI names & CAS numbers
  - Product master data with cosmetic licenses
  - Category hierarchy for materials and products
  - Units of measure with conversion logic
  - Material-supplier approved lists
  - Allergen tracking and storage requirements
  
- 13-API-GATEWAY.md - API Gateway implementation
  - Complete middleware chain (CORS, logging, rate limiting, auth)
  - Routing configuration for all 15 services
  - Circuit breaker pattern for fault tolerance
  - Rate limiting strategy (100 req/min per user)
  - Permission-based route protection
  - Error handling and monitoring metrics

### Phase 3 Completed - 2026-01-23

**Added:**
- 06-SUPPLIER-SERVICE.md - Supplier management documentation
  - GMP, ISO22716, Organic, Ecocert certification tracking
  - Certificate expiry monitoring (90, 30, 7 days alerts)
  - Approved Supplier List (ASL) with approval workflow
  - Supplier evaluation & rating system
  - Quarterly performance reviews
  - Multi-address and multi-contact support
  
- 07-PROCUREMENT-SERVICE.md - Procurement workflow documentation
  - Purchase Requisition (PR) with multi-level approval
  - Purchase Order (PO) creation from PR
  - RFQ (Request for Quotation) process
  - Supplier comparison and selection
  - PO tracking with partial receipt support
  - Budget control and approval rules
  
- 08-WMS-SERVICE.md - Warehouse management system
  - FEFO (First Expired First Out) logic for cosmetics
  - Lot/batch traceability from supplier to customer
  - Cold storage (2-8°C) temperature monitoring
  - Warehouse/Zone/Location hierarchy
  - Stock reservation for production and sales
  - GRN (Goods Receipt Note) with QC workflow
  - Quarantine → Storage zone management
  - Expiry alerts and low stock monitoring

### Phase 4 Completed - 2026-01-23

**Added:**
- 09-MANUFACTURING-SERVICE.md - Manufacturing & production management
  - BOM (Bill of Materials) with AES-256-GCM encryption for formula protection
  - BOM versioning and costing analysis
  - Work Order lifecycle management
  - Material issue with FEFO from WMS
  - QC checkpoints (IQC, IPQC, FQC)
  - NCR (Non-Conformance Report) workflow
  - Batch/Lot traceability (forward & backward)
  - GMP compliance documentation
  
- 10-SALES-SERVICE.md - Sales order management
  - Customer management (retail, wholesale, distributors)
  - Quotation workflow
  - Sales Order with credit limit checks
  - Stock reservation integration with WMS
  - Order fulfillment tracking
  - Customer-specific pricing
  
- 11-MARKETING-SERVICE.md - Marketing & KOL management
  - KOL/Influencer database (tier classification: MEGA, MACRO, MICRO, NANO)
  - Sample request & approval workflow
  - Sample distribution tracking
  - Campaign management and ROI tracking
  - KOL post performance monitoring
  - Marketing budget control

### Phase 5 & 6 Completed - 2026-01-23

**Added:**
- 12-NOTIFICATION-SERVICE.md - Notification management system
  - Email notifications via SMTP
  - In-app notifications
  - Alert rules configuration
  - Notification templates
  - Event subscriptions from all services
  
- 14-EVENT-CATALOG.md - Event-driven architecture documentation
  - Complete event catalog for all 15 services
  - Event naming conventions
  - Publisher/subscriber mapping
  - Event flow diagrams
  - Event schema examples
  
- 15-DATABASE-SCHEMAS.md - Database design documentation
  - Database-per-service strategy (13 databases)
  - Common patterns (UUID, timestamps, soft delete)
  - Index strategies and performance optimization
  - Cosmetics-specific schemas (INCI, CAS, lot tracking, BOM encryption)
  - Backup and migration strategies
  
- 16-DEPLOYMENT.md - Docker Compose deployment guide
  - System requirements (hardware, software)
  - Complete docker-compose.yml configuration
  - Environment variables setup
  - Nginx reverse proxy configuration
  - Database initialization scripts
  - Monitoring stack setup (Prometheus, Grafana)
  - Backup and scaling strategies
  
- 17-IMPLEMENTATION-ROADMAP.md - 9-month implementation plan
  - Phase 1: Infrastructure & Core (Month 1-2)
  - Phase 2: Supply Chain (Month 3-4)
  - Phase 3: Production & Sales (Month 5-6)
  - Phase 4: Advanced Features (Month 7-8)
  - Phase 5: UAT & Go-Live (Month 9)
  - Team composition and resource allocation
  - Risk mitigation strategies
  - Success criteria
  
- 18-GLOSSARY.md - Comprehensive terminology reference
  - General ERP terms
  - Cosmetics-specific terms (INCI, CAS, GMP, FEFO)
  - Business process terms (BOM, PR, PO, GRN, WO, SO)
  - Quality control terms (IQC, IPQC, FQC, NCR)
  - Warehouse terms
  - Technical terms
  - Marketing terms (KOL, engagement rate, ROI)
  - Status values

---

## PROJECT COMPLETION SUMMARY

**Total Documentation**: 19 files (README + CHANGELOG + .gitignore + 16 docs)
**Total Lines**: ~13,000+ lines
**Coverage**: 
- ✅ All 15 microservices documented
- ✅ Architecture & design patterns
- ✅ API specifications & event catalog
- ✅ Database schemas
- ✅ Deployment guide
- ✅ Implementation roadmap
- ✅ Complete glossary

**Tech Stack Documented**:
- Backend: Go 1.22+ (Gin, gRPC, GORM)
- Frontend: Vue 3 + PrimeVue + TypeScript
- Database: PostgreSQL (Database-per-service)
- Cache: Redis
- Message Queue: NATS
- Storage: MinIO
- Deployment: Docker Compose
- Monitoring: Prometheus + Grafana + Loki + Jaeger

**Industry-Specific Features**:
- ✅ INCI name & CAS number tracking
- ✅ GMP & ISO 22716 compliance
- ✅ FEFO (First Expired First Out) logic
- ✅ Lot/batch traceability
- ✅ Cold storage monitoring (2-8°C)
- ✅ BOM formula encryption (AES-256)
- ✅ Certificate expiry monitoring
- ✅ KOL/Influencer management
- ✅ Sample distribution tracking
- ✅ QC checkpoints (IQC, IPQC, FQC)

**Ready for**: Development team kickoff, implementation, and deployment.

## [0.3.0] - 2026-01-24

### Added - User Service Foundation

**Database Layer (9 files)**
- 4 database migrations (8 files total with up/down)
- 3 tables: departments (hierarchical with materialized path), users, user_profiles
- Seed data: 8 default departments (EXEC, IT, HR, FIN, OPS, MFG, WH, QC)
- Seed data: Admin user matching Auth Service (user_id: 09c95223-32b6-4b50-87e7-4ea1333ae072)

**Domain Layer (3 files)**
- User entity: Email, employee code auto-generation, department/manager relationships
- Department entity: Hierarchical structure with parent_id, materialized path, level computation
- UserProfile entity: Extended user information (DOB, address, emergency contact, join date)

**Project Setup**
- go.mod with dependencies
- Clean Architecture structure prepared

### In Progress
- Repository implementations
- Use cases for user and department management
- HTTP REST API endpoints
- gRPC service implementation
- Integration with Auth Service


## [0.4.0] - 2026-01-24

### Added - User Service Complete & Tested

**Implementation Complete (39 files, ~2,500 LOC)**
- Full microservice for user and department management
- Following Auth Service template architecture

**Database Layer**
- 3 tables: departments (hierarchical), users, user_profiles
- Materialized path pattern for efficient tree queries
- Employee code auto-generation: EMP{YYYYMMDD}{sequence}
- Seed data: 8 departments (EXEC, IT, HR, FIN, OPS, MFG, WH, QC)
- Seed data: Admin user matching Auth Service

**Domain Layer**
- User entity with validation and business logic
- Department entity with path computation
- UserProfile entity for extended information
- 3 repository interfaces

**Infrastructure Layer**
- PostgreSQL repositories with tree building
- Auth Service gRPC client (placeholder)
- NATS event publisher

**Use Case Layer**
- User: CreateUser, GetUser, ListUsers
- Department: CreateDepartment, GetDepartmentTree

**HTTP API (6 endpoints)**
- GET /api/v1/users - List with pagination
- POST /api/v1/users - Create user
- GET /api/v1/users/:id - Get user
- GET /api/v1/departments - Get tree
- POST /api/v1/departments - Create
- GET /health, /ready - Health checks

**Testing Results**
- ✅ Health check: PASS
- ✅ Get departments: PASS (8 departments returned)
- ✅ Create user: PASS (employee code auto-generated)
- ✅ List users: PASS (2 users displayed)

### Fixed
- Removed unused imports in user_repo.go and create_user.go
- Fixed compile errors for successful deployment

### Infrastructure
- Restarted PostgreSQL, Redis, NATS containers
- Created user_db database
- Ran all migrations successfully

