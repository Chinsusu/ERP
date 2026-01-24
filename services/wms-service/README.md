# WMS Service (Warehouse Management System)

The most critical service for the cosmetics ERP system, handling all warehouse operations with FEFO (First Expired First Out) logic, lot traceability, and cold storage monitoring.

## Features

- **FEFO Logic**: First Expired First Out - automatically issues materials with earliest expiry dates first
- **Lot Traceability**: Full tracking from supplier → warehouse → production → customer
- **Quarantine Zone**: All incoming goods pass QC before storage
- **Cold Storage Monitoring**: Temperature logging for sensitive materials (2-8°C)
- **GRN Management**: Goods Receipt Notes with complete workflow (DRAFT → QC → COMPLETE)
- **Stock Movements**: Full audit trail of all inventory transactions
- **Stock Reservations**: Reserve materials for work orders and sales orders
- **Real-time Events**: NATS JetStream integration for event-driven architecture

## Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin (HTTP)
- **Database**: PostgreSQL 16
- **ORM**: GORM
- **Event Bus**: NATS JetStream
- **Logging**: Zap

## Ports

- **HTTP**: 8086
- **gRPC**: 9086
- **Database**: 5432 (wms_db)

## Quick Start

```bash
# Navigate to service directory
cd services/wms-service

# Download dependencies
make deps

# Run database migrations
make migrate-up

# Run the service
make run
```

## API Endpoints

### Warehouse & Location
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/warehouses` | List warehouses |
| GET | `/api/v1/warehouses/:id` | Get warehouse details |
| GET | `/api/v1/warehouses/:id/zones` | Get zones in warehouse |
| GET | `/api/v1/zones/:id/locations` | Get locations in zone |

### Stock
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/stock` | Query stock |
| GET | `/api/v1/stock/by-material/:id` | Stock by material with summary |
| GET | `/api/v1/stock/expiring?days=90` | Expiring stock (FEFO) |
| GET | `/api/v1/stock/low-stock?threshold=100` | Low stock alerts |

### Lots
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/lots` | List lots |
| GET | `/api/v1/lots/:id` | Get lot details |
| GET | `/api/v1/lots/:id/movements` | Lot movement history (traceability) |

### GRN (Goods Receipt Notes)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/grn` | Create GRN (from PO) |
| GET | `/api/v1/grn` | List GRNs |
| GET | `/api/v1/grn/:id` | Get GRN details |
| PATCH | `/api/v1/grn/:id/complete` | Complete GRN (after QC) |

### Health Checks
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Service health |
| GET | `/ready` | Readiness check |
| GET | `/live` | Liveness check |

## Database Schema (14 Tables)

1. `warehouses` - Warehouse master data
2. `zones` - Zones within warehouses (RECEIVING, QUARANTINE, STORAGE, COLD, PICKING, SHIPPING)
3. `locations` - Storage locations (Aisle-Rack-Shelf-Bin structure)
4. `lots` - Batch/Lot tracking with expiry
5. `stock` - Current stock by location and lot
6. `stock_movements` - Audit trail of all movements
7. `stock_reservations` - Stock reservations for orders
8. `grns` - Goods Receipt Notes
9. `grn_line_items` - GRN line items
10. `goods_issues` - Goods Issue documents
11. `gi_line_items` - GI line items
12. `stock_adjustments` - Stock adjustments
13. `inventory_counts` - Inventory count documents
14. `temperature_logs` - Cold storage temperature logs

## FEFO Logic (First Expired First Out)

The core algorithm for cosmetics warehouse management:

```go
// Issues stock from lots with earliest expiry dates first
func IssueStockFEFO(materialID uuid.UUID, quantity float64) ([]LotIssued, error) {
    // 1. Get available stock sorted by lot.expiry_date ASC
    // 2. Loop through lots, issuing from earliest expiry first
    // 3. Create movement records for each lot
    // 4. Return list of lots issued
}
```

**Example**: If you have 3 lots of ingredient "Vitamin E":
- LOT-202401-0001: Expires 2024-06-30, Qty: 50kg
- LOT-202401-0002: Expires 2024-08-15, Qty: 100kg  
- LOT-202401-0003: Expires 2024-12-01, Qty: 75kg

Requesting 80kg will issue:
1. 50kg from LOT-202401-0001 (expires first)
2. 30kg from LOT-202401-0002 (next to expire)

## Cosmetics-Specific Features

### Quarantine Workflow
```
Goods Arrive → Quarantine Zone → QC Inspection → Pass/Fail
                                        ↓ Pass
                              Move to Storage Zone
                                        ↓ Fail
                              Return/Dispose
```

### Lot Traceability
- Each lot has: Lot Number, Supplier Lot, Manufactured Date, Expiry Date
- Track movements: GRN → Stock → Work Order/Sales Order
- Support recall: Identify all products using a specific lot

### Cold Storage (2-8°C)
- Zones marked as COLD type
- Temperature logging at configurable intervals
- Alerts for out-of-range temperatures

## Events Published

- `wms.grn.created` - GRN created
- `wms.grn.completed` - GRN completed (triggers PO update)
- `wms.stock.received` - Stock received
- `wms.stock.issued` - Stock issued (with FEFO details)
- `wms.stock.reserved` - Stock reserved
- `wms.stock.low_stock_alert` - Low stock warning
- `wms.lot.expiring_soon` - Lot expiring (90/30/7 days)
- `wms.lot.expired` - Lot expired

## Events Subscribed

- `procurement.po.confirmed` - Prepare for receiving
- `manufacturing.wo.started` - Reserve materials
- `sales.order.confirmed` - Reserve products

## Environment Variables

```env
SERVICE_NAME=wms-service
PORT=8086
GRPC_PORT=9086

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=wms_db

NATS_URL=nats://localhost:4222

# WMS Specific
ENABLE_FEFO=true
EXPIRY_ALERT_DAYS=90,30,7
LOW_STOCK_CHECK_INTERVAL=1h
COLD_STORAGE_MIN_TEMP=2
COLD_STORAGE_MAX_TEMP=8
```

## Project Structure

```
wms-service/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── domain/
│   │   ├── entity/             # Domain entities
│   │   └── repository/         # Repository interfaces
│   ├── infrastructure/
│   │   ├── event/              # Event publisher
│   │   └── persistence/        # PostgreSQL repositories
│   ├── usecase/                # Business logic
│   │   ├── grn/
│   │   ├── lot/
│   │   ├── stock/
│   │   └── warehouse/
│   └── delivery/
│       └── http/               # HTTP handlers
│           ├── dto/
│           ├── handler/
│           └── router/
├── migrations/                 # SQL migrations
├── Makefile
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build Docker image
make docker-build

# Generate protobuf (for gRPC)
make proto
```

## Testing FEFO

```bash
# 1. Create GRN with multiple lots of different expiry dates
curl -X POST http://localhost:8086/api/v1/grn \
  -H "Content-Type: application/json" \
  -d '{
    "grn_date": "2024-01-15",
    "warehouse_id": "a1b2c3d4-1111-1111-1111-111111111111",
    "items": [
      {"material_id": "...", "received_qty": 100, "expiry_date": "2024-06-30", "unit_id": "..."},
      {"material_id": "...", "received_qty": 100, "expiry_date": "2024-12-31", "unit_id": "..."}
    ]
  }'

# 2. Complete GRN (QC passed)
curl -X PATCH http://localhost:8086/api/v1/grn/{id}/complete \
  -H "Content-Type: application/json" \
  -d '{"qc_status": "PASSED"}'

# 3. Check expiring stock
curl http://localhost:8086/api/v1/stock/expiring?days=180

# 4. Issue stock - verify FEFO order
# (Earlier expiry lots should be issued first)
```
