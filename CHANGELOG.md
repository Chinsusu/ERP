# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

