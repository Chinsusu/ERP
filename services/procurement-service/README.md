# Procurement Service

Procurement Service manages Purchase Requisitions (PR) and Purchase Orders (PO) for the ERP Cosmetics System.

## Features

- **Purchase Requisitions (PR)**
  - Create, update, delete PRs
  - Multi-level approval workflow based on amount thresholds
  - PR status tracking (DRAFT → SUBMITTED → PENDING_APPROVAL → APPROVED/REJECTED)
  - Convert approved PRs to POs

- **Purchase Orders (PO)**
  - Create POs from approved PRs
  - PO confirmation and cancellation
  - Track partial receipts from WMS
  - PO amendment history
  - PO status tracking (DRAFT → CONFIRMED → PARTIALLY_RECEIVED → FULLY_RECEIVED → CLOSED)

- **Cosmetics-Specific Features**
  - Integration with Approved Supplier List (ASL)
  - Material specification tracking
  - Certification requirements per material

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP), gRPC
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: NATS
- **ORM**: GORM

## Ports

| Type | Port |
|------|------|
| HTTP | 8085 |
| gRPC | 9085 |

## Quick Start

```bash
# Create database
createdb procurement_db

# Copy environment file
cp .env.example .env

# Run migrations
make migrate-up

# Run service
make run
```

## API Endpoints

### Purchase Requisitions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/purchase-requisitions | Create new PR |
| GET | /api/v1/purchase-requisitions | List PRs with filters |
| GET | /api/v1/purchase-requisitions/:id | Get PR by ID |
| POST | /api/v1/purchase-requisitions/:id/submit | Submit PR for approval |
| POST | /api/v1/purchase-requisitions/:id/approve | Approve PR |
| POST | /api/v1/purchase-requisitions/:id/reject | Reject PR |
| POST | /api/v1/purchase-requisitions/:id/convert-to-po | Convert PR to PO |

### Purchase Orders

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/purchase-orders | List POs with filters |
| GET | /api/v1/purchase-orders/:id | Get PO by ID |
| POST | /api/v1/purchase-orders/:id/confirm | Confirm PO |
| POST | /api/v1/purchase-orders/:id/cancel | Cancel PO |
| POST | /api/v1/purchase-orders/:id/close | Close fully received PO |
| GET | /api/v1/purchase-orders/:id/receipts | Get PO receipts |

## Database Schema

### Tables

1. **purchase_requisitions** - PR header information
2. **pr_line_items** - PR line items with materials
3. **pr_approvals** - PR approval history
4. **purchase_orders** - PO header information
5. **po_line_items** - PO line items with received quantities
6. **po_amendments** - PO amendment history
7. **po_receipts** - Receipt records from WMS

## Events Published

| Event | Description |
|-------|-------------|
| procurement.pr.created | PR created |
| procurement.pr.submitted | PR submitted for approval |
| procurement.pr.approved | PR approved |
| procurement.pr.rejected | PR rejected |
| procurement.po.created | PO created |
| procurement.po.confirmed | PO confirmed (WMS subscribes) |
| procurement.po.received | PO goods received |
| procurement.po.closed | PO closed |

## Events Subscribed

| Event | Source | Action |
|-------|--------|--------|
| wms.grn.completed | WMS Service | Update PO received quantities |
| supplier.blocked | Supplier Service | Alert on blocked supplier POs |

## Approval Levels

| Amount (VND) | Required Approval |
|--------------|-------------------|
| < 10,000,000 | Department Manager |
| 10,000,000 - 50,000,000 | Procurement Manager |
| 50,000,000 - 200,000,000 | Finance Director |
| > 200,000,000 | General Director |

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP port | 8085 |
| GRPC_PORT | gRPC port | 9085 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_NAME | Database name | procurement_db |
| NATS_URL | NATS server URL | nats://localhost:4222 |

## Project Structure

```
procurement-service/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── domain/
│   │   ├── entity/             # Domain entities
│   │   └── repository/         # Repository interfaces
│   ├── infrastructure/
│   │   ├── persistence/        # PostgreSQL implementations
│   │   └── event/              # NATS event publisher
│   ├── usecase/
│   │   ├── pr/                 # PR use cases
│   │   └── po/                 # PO use cases
│   └── delivery/
│       └── http/
│           ├── dto/            # Request/Response DTOs
│           ├── handler/        # HTTP handlers
│           └── router/         # Gin router
├── migrations/                 # Database migrations
├── api/
│   └── proto/                  # Protobuf definitions
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

## Dependencies

- **Supplier Service** - For Approved Supplier List
- **Master Data Service** - For material information
- **User Service** - For approver information
- **WMS Service** - For goods receipt integration
