# 09 - MANUFACTURING SERVICE

## TỔNG QUAN

Manufacturing Service quản lý toàn bộ quy trình sản xuất mỹ phẩm, từ BOM (Bill of Materials) với công thức bảo mật, work order, material issue, production tracking, đến QC checkpoints và traceability.

### Responsibilities

✅ BOM (Bill of Materials) với versioning & encryption  
✅ BOM costing & material cost analysis  
✅ Work Order creation & tracking  
✅ Material issue to production (FEFO từ WMS)  
✅ Production completion & yield tracking  
✅ QC checkpoints (IQC, IPQC, FQC)  
✅ NCR (Non-Conformance Report)  
✅ Batch/Lot traceability (forward & backward)

### Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP) + gRPC
- **Database**: PostgreSQL
- **Encryption**: AES-256-GCM (BOM formulas)

### Ports

- HTTP: `8087`
- gRPC: `9087`

---

## WHY MANUFACTURING IS CRITICAL FOR COSMETICS

### 1. Formula Confidentiality

BOM chứa công thức sản phẩm - **trade secret** của công ty:
- Percentage của từng ingredient
- Processing instructions
- Critical parameters

→ **Must encrypt** `formula_details` field
→ **Strict RBAC**: Only R&D Manager, Production Manager can view full BOM

### 2. GMP Compliance

Good Manufacturing Practice yêu cầu:
- Document everything
- Batch records
- QC at multiple stages
- Traceability

### 3. Quality Control Checkpoints

**IQC (Incoming QC)**: Nguyên liệu nhập về  
**IPQC (In-Process QC)**: Quá trình sản xuất  
**FQC (Final QC)**: Thành phẩm trước khi xuất kho

### 4. Traceability

- **Backward trace**: Product lot → material lots used
- **Forward trace**: Material lot → which product lots

Example: Nếu Vitamin C có vấn đề → tìm tất cả sản phẩm đã dùng lot Vitamin C đó.

---

## DATABASE SCHEMA (Key Tables)

Chi tiết trong `02-SERVICE-SPECIFICATIONS.md`:

- `boms` + `bom_line_items` (with encrypted formula)
- `work_orders` + `wo_material_issues`
- `qc_inspections`
- `ncrs` (Non-Conformance Reports)
- `batch_traceability`

### BOM Encryption

```sql
CREATE TABLE boms (
    ...
    formula_details JSONB, -- ENCRYPTED with AES-256-GCM
    confidentiality_level VARCHAR(50) DEFAULT 'RESTRICTED',
    ...
);

-- Encrypted content example (before encryption):
{
  "processing_steps": [
    "Heat Phase A to 75°C",
    "Add Phase B slowly while stirring",
    "Cool to 40°C, add actives"
  ],
  "critical_parameters": {
    "mixing_speed": "1200 rpm",
    "mixing_time": "20 minutes",
    "cooling_rate": "5°C per 10 minutes"
  },
  "notes": "Vitamin C must be added last to prevent oxidation"
}
```

---

## API ENDPOINTS (Key)

### BOM Management

```
GET    /api/v1/manufacturing/boms
POST   /api/v1/manufacturing/boms
GET    /api/v1/manufacturing/boms/:id
PUT    /api/v1/manufacturing/boms/:id
POST   /api/v1/manufacturing/boms/:id/copy        # Create new version
PATCH  /api/v1/manufacturing/boms/:id/approve
POST   /api/v1/manufacturing/boms/:id/cost        # Calculate cost
```

#### POST /api/v1/manufacturing/boms

**Permission**: `manufacturing:bom:create`

**Request**:
```json
{
  "bom_number": "BOM-SERUM-001",
  "product_id": "product-uuid",
  "version": 1,
  "name": "Vitamin C Serum 30ml - Ver 1.0",
  "description": "Initial formula for Vitamin C brightening serum",
  "batch_size": 10,
  "batch_uom_id": "liter-uuid",
  "items": [
    {
      "line_number": 1,
      "material_id": "water-uuid",
      "quantity": 6.5,
      "uom_id": "liter-uuid",
      "item_type": "MATERIAL",
      "is_critical": false
    },
    {
      "line_number": 2,
      "material_id": "vitamin-c-uuid",
      "quantity": 2.0,
      "uom_id": "kg-uuid",
      "quantity_min": 1.95,
      "quantity_max": 2.05,
      "item_type": "MATERIAL",
      "is_critical": true,
      "notes": "Must be pure L-Ascorbic Acid"
    },
    {
      "line_number": 3,
      "material_id": "bottle-30ml-uuid",
      "quantity": 1000,
      "uom_id": "piece-uuid",
      "item_type": "PACKAGING"
    }
  ],
  "formula_details": {
    "processing_steps": [
      "Prepare Phase A: Mix water + preservative at room temp",
      "Prepare Phase B: Dissolve Vitamin C slowly",
      "Combine A + B, mix at 800rpm for 15min",
      "Add fragrance, mix 5min",
      "Fill into bottles"
    ],
    "critical_parameters": {
      "temperature_range": "20-25°C",
      "mixing_speed": "800 rpm",
      "ph_target": "3.5-4.0"
    }
  },
  "confidentiality_level": "CONFIDENTIAL"
}
```

**Response 201**:
```json
{
  "id": "bom-uuid",
  "bom_number": "BOM-SERUM-001",
  "product_name": "Vitamin C Serum 30ml",
  "version": 1,
  "status": "DRAFT",
  "material_cost": 125000,
  "total_cost": 145000
}
```

**Note**: `formula_details` được encrypt trước khi lưu vào DB.

---

#### GET /api/v1/manufacturing/boms/:id

**Permission**: `manufacturing:bom:read`

**Response** (for users with full access):
```json
{
  "id": "bom-uuid",
  "bom_number": "BOM-SERUM-001",
  "product": {
    "id": "product-uuid",
    "product_code": "FG-SERUM-001",
    "name": "Vitamin C Serum 30ml"
  },
  "version": 1,
  "batch_size": 10,
  "batch_uom": "Liter",
  "items": [
    {
      "line_number": 1,
      "material_code": "RM-001",
      "material_name": "Purified Water",
      "quantity": 6.5,
      "uom": "L",
      "unit_cost": 5000,
      "total_cost": 32500
    },
    {
      "line_number": 2,
      "material_code": "RM-050",
      "material_name": "Vitamin C (Ascorbic Acid)",
      "quantity": 2.0,
      "uom": "KG",
      "is_critical": true,
      "unit_cost": 2500000,
      "total_cost": 5000000
    }
  ],
  "formula_details": {
    "processing_steps": [...],
    "critical_parameters": {...}
  },
  "material_cost": 5125000,
  "labor_cost": 150000,
  "overhead_cost": 80000,
  "total_cost": 5355000,
  "status": "APPROVED",
  "effective_from": "2024-02-01"
}
```

**Response** (for users without full BOM access):
```json
{
  "id": "bom-uuid",
  "bom_number": "BOM-SERUM-001",
  "product_name": "Vitamin C Serum 30ml",
  "version": 1,
  "batch_size": 10,
  "status": "APPROVED",
  "message": "Full BOM details restricted. Contact R&D Manager."
}
```

---

### Work Orders

```
GET    /api/v1/manufacturing/work-orders
POST   /api/v1/manufacturing/work-orders
GET    /api/v1/manufacturing/work-orders/:id
PATCH  /api/v1/manufacturing/work-orders/:id/release
PATCH  /api/v1/manufacturing/work-orders/:id/start
PATCH  /api/v1/manufacturing/work-orders/:id/complete
POST   /api/v1/manufacturing/work-orders/:id/issue-material
```

#### POST /api/v1/manufacturing/work-orders

**Permission**: `manufacturing:wo:create`

**Request**:
```json
{
  "wo_date": "2024-01-23",
  "product_id": "product-uuid",
  "bom_id": "bom-uuid",
  "planned_quantity": 1000,
  "uom_id": "bottle-uuid",
  "planned_start_date": "2024-01-25",
  "planned_end_date": "2024-01-26",
  "batch_number": "BATCH-2024-001",
  "sales_order_id": "so-uuid",
  "production_line": "Line 1",
  "shift": "Morning",
  "priority": "HIGH"
}
```

**Response 201**:
```json
{
  "id": "wo-uuid",
  "wo_number": "WO-2024-001",
  "product_name": "Vitamin C Serum 30ml",
  "planned_quantity": 1000,
  "status": "PLANNED"
}
```

---

#### PATCH /api/v1/manufacturing/work-orders/:id/start

**Permission**: `manufacturing:wo:execute`

**Request**:
```json
{
  "actual_start_date": "2024-01-25T08:00:00Z",
  "supervisor_id": "user-uuid"
}
```

**Actions**:
1. Update status = IN_PROGRESS
2. Reserve materials từ WMS
3. Publish event `manufacturing.wo.started`

---

#### POST /api/v1/manufacturing/work-orders/:id/issue-material

Issue materials cho work order.

**Request**:
```json
{
  "bom_line_item_id": "bom-line-uuid",
  "quantity": 2.0,
  "uom_id": "kg-uuid"
}
```

**Response**:
```json
{
  "movement_number": "MOV-2024-001",
  "lots_used": [
    {
      "lot_number": "LOT-2023-150",
      "quantity": 2.0,
      "expiry_date": "2025-06-30"
    }
  ]
}
```

**Note**: WMS service handles FEFO selection.

---

#### PATCH /api/v1/manufacturing/work-orders/:id/complete

**Permission**: `manufacturing:wo:complete`

**Request**:
```json
{
  "actual_end_date": "2024-01-26T16:00:00Z",
  "actual_quantity": 980,
  "good_quantity": 975,
  "rejected_quantity": 5,
  "notes": "Minor issues with 5 bottles (seal defects)"
}
```

**Actions**:
1. Update WO status = COMPLETED
2. Create product lot in WMS
3. Create traceability records
4. Publish event `manufacturing.wo.completed`

---

### Quality Control

```
GET    /api/v1/manufacturing/qc
POST   /api/v1/manufacturing/qc
GET    /api/v1/manufacturing/qc/:id
PATCH  /api/v1/manufacturing/qc/:id/approve
```

#### POST /api/v1/manufacturing/qc

**Permission**: `manufacturing:qc:perform`

**Request**:
```json
{
  "inspection_date": "2024-01-26T14:00:00Z",
  "inspection_type": "FQC",
  "reference_type": "WORK_ORDER",
  "reference_id": "wo-uuid",
  "product_id": "product-uuid",
  "lot_id": "lot-uuid",
  "inspected_quantity": 980,
  "accepted_quantity": 975,
  "rejected_quantity": 5,
  "result": "PASSED",
  "test_results": {
    "ph": 3.8,
    "viscosity": "Good",
    "color": "Light yellow",
    "fragrance": "Pleasant citrus",
    "packaging": "5 defects found"
  },
  "notes": "Overall quality acceptable"
}
```

---

### Non-Conformance Reports

```
GET    /api/v1/manufacturing/ncr
POST   /api/v1/manufacturing/ncr
PATCH  /api/v1/manufacturing/ncr/:id/close
```

#### POST /api/v1/manufacturing/ncr

**Permission**: `manufacturing:ncr:create`

**Request**:
```json
{
  "ncr_date": "2024-01-26",
  "nc_type": "PROCESS",
  "reference_type": "WORK_ORDER",
  "reference_id": "wo-uuid",
  "description": "Mixing time was 15min instead of 20min due to equipment issue",
  "severity": "MEDIUM",
  "root_cause": "Mixer timer malfunction",
  "corrective_action": "Calibrated mixer timer. Reprocessed batch.",
  "preventive_action": "Add weekly timer calibration to PM schedule"
}
```

---

### Traceability

```
GET /api/v1/manufacturing/traceability/forward/:lot_id
GET /api/v1/manufacturing/traceability/backward/:lot_id
```

#### GET /api/v1/manufacturing/traceability/backward/:lot_id

Backward trace: Product lot → material lots used.

**Response**:
```json
{
  "finished_lot": {
    "lot_number": "BATCH-2024-001",
    "product_code": "FG-SERUM-001",
    "product_name": "Vitamin C Serum 30ml",
    "quantity": 975,
    "manufactured_date": "2024-01-26"
  },
  "work_order": {
    "wo_number": "WO-2024-001",
    "supervisor": "Nguyễn Văn A"
  },
  "materials_used": [
    {
      "material_code": "RM-001",
      "material_name": "Purified Water",
      "lot_number": "LOT-2024-010",
      "quantity": 6.5,
      "supplier": "ABC Water Co."
    },
    {
      "material_code": "RM-050",
      "material_name": "Vitamin C",
      "lot_number": "LOT-2023-150",
      "quantity": 2.0,
      "supplier": "XYZ Chemicals",
      "supplier_lot": "SUP-VC-2023-150"
    }
  ]
}
```

#### GET /api/v1/manufacturing/traceability/forward/:lot_id

Forward trace: Material lot → product lots.

**Response**:
```json
{
  "material_lot": {
    "lot_number": "LOT-2023-150",
    "material_code": "RM-050",
    "material_name": "Vitamin C",
    "supplier": "XYZ Chemicals"
  },
  "used_in_products": [
    {
      "product_lot": "BATCH-2024-001",
      "product_name": "Vitamin C Serum 30ml",
      "wo_number": "WO-2024-001",
      "quantity_used": 2.0,
      "production_date": "2024-01-26"
    },
    {
      "product_lot": "BATCH-2024-005",
      "product_name": "Vitamin C Cream 50g",
      "wo_number": "WO-2024-005",
      "quantity_used": 1.5,
      "production_date": "2024-01-28"
    }
  ],
  "total_quantity_used": 3.5,
  "remaining_quantity": 46.5
}
```

---

## gRPC METHODS

### GetBOM

```protobuf
message GetBOMRequest {
  string bom_id = 1;
}

message GetBOMResponse {
  BOM bom = 1;
}
```

### ValidateBOMActive

```protobuf
message ValidateBOMActiveRequest {
  string product_id = 1;
}

message ValidateBOMActiveResponse {
  bool has_active_bom = 1;
  string bom_id = 2;
  int32 version = 3;
}
```

---

## EVENTS

### Events Published

```yaml
manufacturing.bom.created:
manufacturing.bom.approved:
manufacturing.wo.created:
manufacturing.wo.released:
manufacturing.wo.started:
manufacturing.wo.completed:
  payload:
    wo_id: uuid
    wo_number: string
    product_id: uuid
    batch_number: string
    good_quantity: decimal
    
manufacturing.qc.failed:
manufacturing.ncr.created:
manufacturing.ncr.closed:
```

### Events Subscribed

```yaml
sales.order.confirmed:
  action: Create work order (if make-to-order)
  
wms.stock.low_stock_alert:
  action: Trigger production planning
```

---

## BUSINESS LOGIC

### BOM Versioning

- Product có thể có multiple BOM versions
- Chỉ 1 version ACTIVE tại 1 thời điểm
- Version history preserved
- Copy BOM tạo version mới

### Work Order Lifecycle

```
PLANNED → RELEASED → IN_PROGRESS → QC_PENDING → COMPLETED
                                               ↓
                                           CANCELLED
```

### Yield Tracking

```
Yield % = (good_quantity / planned_quantity) * 100

Example:
Planned: 1000 bottles
Actual: 980 bottles
Good: 975 bottles
Rejected: 5 bottles

Yield = 975/1000 = 97.5%
```

---

## CONFIGURATION

```bash
MFG_SERVICE_PORT=8087
MFG_GRPC_PORT=9087
MFG_DB_HOST=postgres
MFG_DB_PORT=5438
MFG_DB_NAME=manufacturing_db

# BOM Encryption
BOM_ENCRYPTION_KEY=<32-byte-hex-key>
BOM_ENCRYPTION_ALGORITHM=AES-256-GCM

# Permissions
RESTRICT_BOM_VIEW=true
BOM_FULL_ACCESS_ROLES=rd_manager,production_manager,ceo
```

---

## MONITORING METRICS

```
manufacturing_boms_total{status}
manufacturing_wo_total{status}
manufacturing_wo_yield_percentage
manufacturing_qc_pass_rate
manufacturing_ncr_total{severity}
```

---

## DEPENDENCIES

- **Master Data** (gRPC): Product, material info
- **WMS** (gRPC): Reserve & issue materials, create product lots
- **Sales** (Subscribe): Order events
- **NATS**: Events

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
