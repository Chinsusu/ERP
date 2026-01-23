# 08 - WMS SERVICE (Warehouse Management)

## TỔNG QUAN

WMS Service là hệ thống quản lý kho quan trọng nhất, xử lý toàn bộ hoạt động nhập-xuất-tồn kho, lot/batch tracking, FEFO logic, và cold storage monitoring.

### Responsibilities

✅ Warehouse, zone, location hierarchy  
✅ Lot/Batch management với expiry tracking  
✅ Stock on hand, available, reserved  
✅ Goods Receipt Note (GRN) từ PO  
✅ Goods Issue với FEFO logic  
✅ Stock transfers & adjustments  
✅ Inventory counting/cycle count  
✅ Cold storage temperature monitoring  
✅ Stock reservation cho production/sales

### Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP) + gRPC
- **Database**: PostgreSQL

### Ports

- HTTP: `8086`
- gRPC: `9086`

---

## WHY WMS IS CRITICAL FOR COSMETICS

### 1. FEFO (First Expired First Out)

Khác với FIFO, cosmetics phải xuất hàng sắp HẾT HẠN trước:
- Serum: 2-3 years shelf life
- Creams: 1-2 years
- Raw materials: varies

### 2. Lot Traceability

Theo dõi từng lot từ nhà cung cấp → production → customer:
- Supplier lot number
- Internal lot number
- Manufacturing date
- Expiry date

### 3. Cold Storage (2-8°C)

Một số nguyên liệu yêu cầu bảo quản lạnh:
- Peptides
- Certain vitamins
- Plant extracts
- Temperature monitoring 24/7

### 4. Quarantine Zone

Hàng nhập về phải qua IQC trước khi vào kho chính:
- QUARANTINE zone
- QC passed → move to STORAGE zone

---

## WAREHOUSE HIERARCHY

```
Warehouse (main, cold_storage, quarantine, finished_goods)
  └── Zone (receiving, storage, picking, shipping, cold)
      └── Location (aisle, rack, shelf, bin)
          └── Stock (material + lot)
```

Example:
```
WH-MAIN (Main Warehouse)
  ├── ZONE-RECV (Receiving)
  ├── ZONE-QUAR (Quarantine)
  ├── ZONE-STOR (Storage)
  │   └── LOC-A01-R02-S03 (Aisle A01, Rack 02, Shelf 03)
  │       └── Stock: Material RM-001, Lot LOT-2024-001, Qty: 100kg
  └── ZONE-PICK (Picking)

WH-COLD (Cold Storage 2-8°C)
  └── ZONE-COLD-01
      └── LOC-C01-R01-S01
```

---

## DATABASE SCHEMA (Key Tables)

Chi tiết trong `02-SERVICE-SPECIFICATIONS.md`:

- `warehouses`, `zones`, `locations`
- `lots` (batch tracking)
- `stock` (quantity by location + lot)
- `stock_movements` (transaction log)
- `grns` + `grn_line_items`
- `stock_reservations`

### Key Constraint

```sql
-- Stock cannot be negative
CHECK (quantity >= 0)
CHECK (reserved_quantity >= 0)
CHECK (reserved_quantity <= quantity)

-- Available = quantity - reserved
available_quantity GENERATED ALWAYS AS (quantity - reserved_quantity) STORED
```

---

## API ENDPOINTS (Key)

### Stock Inquiry

```
GET /api/v1/warehouse/stock
GET /api/v1/warehouse/stock/material/:material_id
GET /api/v1/warehouse/stock/expiring
GET /api/v1/warehouse/stock/low-stock
```

#### GET /api/v1/warehouse/stock

**Permission**: `wms:stock:read`

**Query Parameters**:
- `warehouse_id`: uuid
- `zone_id`: uuid
- `material_id`: uuid
- `lot_id`: uuid
- `has_stock`: boolean (only locations with stock)
- `expiring_within_days`: int

**Response 200**:
```json
{
  "data": [
    {
      "location": {
        "id": "loc-uuid",
        "code": "A01-R02-S03",
        "warehouse": "Main Warehouse",
        "zone": "Storage"
      },
      "material": {
        "id": "mat-uuid",
        "material_code": "RM-001",
        "name": "Vitamin C"
      },
      "lot": {
        "id": "lot-uuid",
        "lot_number": "LOT-2024-001",
        "manufactured_date": "2024-01-15",
        "expiry_date": "2026-01-14",
        "days_until_expiry": 720
      },
      "quantity": 100,
      "reserved_quantity": 15,
      "available_quantity": 85,
      "uom": "KG"
    }
  ],
  "summary": {
    "total_quantity": 500,
    "total_reserved": 75,
    "total_available": 425
  }
}
```

---

#### GET /api/v1/warehouse/stock/expiring

Hàng sắp hết hạn (FEFO logic).

**Query Parameters**:
- `days`: int (default: 90)
- `material_id`: uuid (optional)

**Response 200**:
```json
{
  "data": [
    {
      "material_code": "RM-050",
      "material_name": "Hyaluronic Acid",
      "lot_number": "LOT-2023-150",
      "expiry_date": "2024-03-31",
      "days_until_expiry": 45,
      "quantity": 20,
      "location_code": "A05-R01-S02",
      "warehouse": "Main Warehouse"
    }
  ],
  "total": 15
}
```

---

### Goods Receipt (GRN)

```
GET  /api/v1/warehouse/grn
POST /api/v1/warehouse/grn
GET  /api/v1/warehouse/grn/:id
PUT  /api/v1/warehouse/grn/:id
PATCH /api/v1/warehouse/grn/:id/complete
```

#### POST /api/v1/warehouse/grn

Create GRN from PO.

**Permission**: `wms:grn:create`

**Request**:
```json
{
  "grn_date": "2024-01-23",
  "po_id": "po-uuid",
  "warehouse_id": "wh-main-uuid",
  "delivery_note_number": "DN-SUP-001",
  "vehicle_number": "79A-12345",
  "items": [
    {
      "po_line_item_id": "po-line-uuid",
      "material_id": "material-uuid",
      "received_quantity": 50,
      "uom_id": "kg-uuid",
      "supplier_lot_number": "SUP-LOT-123",
      "manufactured_date": "2024-01-10",
      "expiry_date": "2026-01-09",
      "location_id": "quarantine-loc-uuid"
    }
  ],
  "notes": "Goods received in good condition"
}
```

**Response 201**:
```json
{
  "id": "grn-uuid",
  "grn_number": "GRN-2024-001",
  "status": "DRAFT",
  "po_number": "PO-2024-001"
}
```

**Process**:
1. Create GRN (status=DRAFT)
2. Create lot for each item
3. Put stock in QUARANTINE zone
4. Await QC inspection
5. If QC pass → Complete GRN → Move to STORAGE
6. Publish event → Procurement updates PO

---

#### PATCH /api/v1/warehouse/grn/:id/complete

**Permission**: `wms:grn:complete`

**Request**:
```json
{
  "qc_status": "PASSED",
  "qc_notes": "All tests passed. Quality acceptable."
}
```

**Actions**:
1. Update GRN status = COMPLETED
2. Move stock from QUARANTINE → STORAGE
3. Publish `wms.grn.completed` event
4. Procurement Service updates PO

---

### Goods Issue

```
POST /api/v1/warehouse/issue
```

#### POST /api/v1/warehouse/issue

Issue materials (FEFO logic).

**Permission**: `wms:issue:create`

**Request**:
```json
{
  "issue_date": "2024-01-23",
  "issue_type": "PRODUCTION",
  "reference_id": "work-order-uuid",
  "reference_number": "WO-2024-001",
  "items": [
    {
      "material_id": "material-uuid",
      "quantity": 25,
      "uom_id": "kg-uuid"
    }
  ],
  "notes": "Materials for work order WO-2024-001"
}
```

**FEFO Logic**:
1. Find all lots của material có stock
2. Sort by expiry_date ASC (sớm nhất trước)
3. Issue từ lot sắp hết hạn trước
4. Nếu lot không đủ, lấy từ lot tiếp theo

**Response 200**:
```json
{
  "movement_number": "MOV-OUT-2024-001",
  "issued_from_lots": [
    {
      "lot_number": "LOT-2023-100",
      "quantity": 20,
      "expiry_date": "2025-06-30"
    },
    {
      "lot_number": "LOT-2024-001",
      "quantity": 5,
      "expiry_date": "2026-01-14"
    }
  ]
}
```

---

### Stock Reservation

```
POST /api/v1/warehouse/reserve
POST /api/v1/warehouse/reserve/:id/release
```

#### POST /api/v1/warehouse/reserve

Reserve stock cho sales order hoặc work order.

**Permission**: `wms:stock:reserve`

**Request**:
```json
{
  "material_id": "material-uuid",
  "quantity": 50,
  "uom_id": "kg-uuid",
  "reservation_type": "SALES_ORDER",
  "reference_id": "so-uuid",
  "reference_number": "SO-2024-001",
  "expires_at": "2024-02-23T00:00:00Z"
}
```

**Response 200**:
```json
{
  "reservation_id": "res-uuid",
  "reserved_quantity": 50,
  "reserved_from_locations": [
    {
      "location_code": "A01-R02-S03",
      "quantity": 50
    }
  ]
}
```

---

### Stock Adjustment

```
POST /api/v1/warehouse/adjustment
```

#### POST /api/v1/warehouse/adjustment

Adjust stock (cycle count corrections,損耗).

**Permission**: `wms:adjustment:create` (create), `wms:adjustment:approve` (approve)

**Request**:
```json
{
  "adjustment_date": "2024-01-23",
  "adjustment_type": "CYCLE_COUNT",
  "location_id": "loc-uuid",
  "material_id": "material-uuid",
  "lot_id": "lot-uuid",
  "system_quantity": 100,
  "actual_quantity": 98,
  "adjustment_quantity": -2,
  "uom_id": "kg-uuid",
  "reason": "Physical count discrepancy",
  "notes": "Minor evaporation loss"
}
```

---

### Stock Transfer

```
POST /api/v1/warehouse/transfer
```

Transfer stock giữa locations.

**Request**:
```json
{
  "transfer_date": "2024-01-23",
  "material_id": "material-uuid",
  "lot_id": "lot-uuid",
  "from_location_id": "loc-a-uuid",
  "to_location_id": "loc-b-uuid",
  "quantity": 10,
  "uom_id": "kg-uuid",
  "reason": "Reorganization"
}
```

---

## gRPC METHODS

### CheckStockAvailability

```protobuf
message CheckStockAvailabilityRequest {
  string material_id = 1;
  double quantity = 2;
  string uom_id = 3;
}

message CheckStockAvailabilityResponse {
  bool available = 1;
  double available_quantity = 2;
  double shortage = 3;
}
```

### ReserveStock

```protobuf
message ReserveStockRequest {
  string material_id = 1;
  double quantity = 2;
  string reservation_type = 3;
  string reference_id = 4;
}

message ReserveStockResponse {
  string reservation_id = 1;
  bool success = 2;
}
```

### IssueStock

```protobuf
message IssueStockRequest {
  string material_id = 1;
  double quantity = 2;
  bool use_fefo = 3;
  string reference_id = 4;
}

message IssueStockResponse {
  string movement_number = 1;
  repeated LotIssued lots_issued = 2;
}
```

---

## EVENTS

### Events Published

```yaml
wms.grn.created:
  payload:
    grn_id: uuid
    grn_number: string
    po_id: uuid

wms.grn.completed:
  payload:
    grn_id: uuid
    grn_number: string
    po_id: uuid
    items: array

wms.stock.received:
  payload:
    material_id: uuid
    lot_id: uuid
    quantity: decimal
    location_id: uuid

wms.stock.issued:
  payload:
    material_id: uuid
    quantity: decimal
    lots_used: array
    reference_type: string
    reference_id: uuid

wms.stock.low_stock_alert:
  payload:
    material_id: uuid
    material_code: string
    current_quantity: decimal
    reorder_point: decimal

wms.lot.expiring_soon:
  payload:
    lot_id: uuid
    lot_number: string
    material_id: uuid
    expiry_date: date
    days_until_expiry: int
    quantity: decimal
```

### Events Subscribed

```yaml
procurement.po.confirmed:
  action: Prepare for receiving
  handler: PrepareForGRN

manufacturing.workorder.started:
  action: Reserve materials
  handler: ReserveMaterialsForWO

sales.order.confirmed:
  action: Reserve products
  handler: ReserveProductsForSO
```

---

## BUSINESS LOGIC

### FEFO Implementation

```go
func IssueStockFEFO(materialID string, quantity float64) []LotIssued {
    // Get all lots with available stock
    lots := GetLotsWithStock(materialID)
    
    // Sort by expiry_date ASC (earliest first)
    sort.Slice(lots, func(i, j int) bool {
        return lots[i].ExpiryDate.Before(lots[j].ExpiryDate)
    })
    
    remaining := quantity
    issued := []LotIssued{}
    
    for _, lot := range lots {
        if remaining <= 0 {
            break
        }
        
        issueQty := min(lot.AvailableQuantity, remaining)
        
        // Create stock movement
        CreateStockMovement(lot, issueQty, "GOODS_ISSUE")
        
        issued = append(issued, LotIssued{
            LotNumber: lot.LotNumber,
            Quantity: issueQty,
            ExpiryDate: lot.ExpiryDate,
        })
        
        remaining -= issueQty
    }
    
    return issued
}
```

### Cold Storage Monitoring

**Daily Job**:
- Check temperature logs (from IoT sensors)
- Alert if temperature out of range (2-8°C)
- Track cold chain compliance

### Expiry Alert Schedule

- 90 days before: Email notification
- 30 days before: Daily email
- 7 days before: Critical alert
- On expiry: Auto-quarantine lot

---

## CONFIGURATION

```bash
# Server
WMS_SERVICE_PORT=8086
WMS_GRPC_PORT=9086

# Database
WMS_DB_HOST=postgres
WMS_DB_PORT=5437
WMS_DB_NAME=wms_db

# Business Rules
ENABLE_FEFO=true
DEFAULT_ZONE_QUARANTINE=ZONE-QUAR
DEFAULT_ZONE_STORAGE=ZONE-STOR

# Alerts
EXPIRY_ALERT_DAYS=90,30,7
LOW_STOCK_CHECK_INTERVAL=1h

# Cold Storage
COLD_STORAGE_MIN_TEMP=2
COLD_STORAGE_MAX_TEMP=8
TEMP_CHECK_INTERVAL=15m
```

---

## MONITORING METRICS

```
wms_stock_total_quantity{material_id, warehouse}
wms_stock_reserved_quantity{material_id}
wms_lots_expiring_count{days_range}
wms_grn_total{status}
wms_stock_movements_total{type}
wms_cold_storage_temperature{zone}
```

---

## DEPENDENCIES

- **Master Data Service** (gRPC): Get material info
- **Procurement Service** (Subscribe): PO events
- **Manufacturing Service** (Subscribe): Work order events
- **Sales Service** (Subscribe): Sales order events
- **NATS**: Events

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
