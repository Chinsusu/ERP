# 05 - MASTER DATA SERVICE

## TỔNG QUAN

Master Data Service quản lý dữ liệu tham khảo chính cho toàn bộ hệ thống ERP mỹ phẩm, bao gồm nguyên liệu, sản phẩm, danh mục, và đơn vị đo lường.

### Responsibilities

✅ Material master data (raw materials, packaging, semi-finished)  
✅ Product master data (finished goods)  
✅ Categories & subcategories (hierarchical)  
✅ Units of measure (UoM) & conversions  
✅ Material specifications (INCI name, CAS number, allergens)  
✅ Material-Supplier approved list  
✅ Price lists & costing

### Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP) + gRPC
- **Database**: PostgreSQL
- **File Storage**: MinIO (SDS files, product images)

### Ports

- HTTP: `8083`
- gRPC: `9083`

---

## KEY FEATURES FOR COSMETICS INDUSTRY

### 1. INCI & CAS Tracking

**INCI (International Nomenclature of Cosmetic Ingredients)**: Tên chuẩn quốc tế cho nguyên liệu mỹ phẩm.  
**CAS Number**: Số đăng ký Chemical Abstracts Service.

Example:
- Material: Vitamin C
- INCI Name: Ascorbic Acid
- CAS Number: 50-81-7

### 2. Allergen Management

Track allergens trong nguyên liệu để:
- Compliance với regulations
- Product labeling requirements
- Customer safety

### 3. Origin & Certifications

- Country of origin
- Organic certification
- Natural certification
- Vegan/Cruelty-free status

### 4. Storage Requirements

- Cold storage (2-8°C)
- Room temperature
- Special handling requirements

---

## DATABASE SCHEMA (Summary)

Chi tiết schema đã có trong `02-SERVICE-SPECIFICATIONS.md`. Đây là các tables chính:

### Core Tables

- `categories` - Hierarchical categorization
- `units_of_measure` - UoM với conversion factors
- `materials` - Raw materials & packaging
- `products` - Finished goods
- `material_suppliers` - Approved supplier list per material

### Key Fields

**materials table**:
- `inci_name`, `cas_number`
- `is_organic`, `is_natural`
- `requires_cold_storage`, `shelf_life_days`
- `is_hazardous`, `allergen_info`
- `safety_data_sheet_url`

**products table**:
- `cosmetic_license_number`, `license_expiry_date`
- `ingredients_summary` (for label)
- `target_skin_type`
- `barcode`, `sku`

---

## API ENDPOINTS (Highlights)

### Materials

```
GET    /api/v1/master-data/materials
POST   /api/v1/master-data/materials
GET    /api/v1/master-data/materials/:id
PUT    /api/v1/master-data/materials/:id
DELETE /api/v1/master-data/materials/:id
GET    /api/v1/master-data/materials/:id/suppliers
POST   /api/v1/master-data/materials/:id/suppliers
```

#### GET /api/v1/master-data/materials

**Query Parameters**:
- `material_type`: RAW_MATERIAL, PACKAGING, SEMI_FINISHED
- `category_id`: uuid
- `is_organic`: boolean
- `requires_cold_storage`: boolean
- `is_active`: boolean
- `search`: string (name, material_code, inci_name, cas_number)
- `page`, `limit`

**Response 200**:
```json
{
  "data": [
    {
      "id": "material-uuid",
      "material_code": "RM-001",
      "name": "Vitamin C (Ascorbic Acid)",
      "inci_name": "Ascorbic Acid",
      "cas_number": "50-81-7",
      "material_type": "RAW_MATERIAL",
      "category": {
        "id": "cat-uuid",
        "name": "Active Ingredients"
      },
      "is_organic": false,
      "is_natural": true,
      "requires_cold_storage": true,
      "shelf_life_days": 365,
      "standard_cost": 150000,
      "base_uom": {
        "id": "uom-uuid",
        "code": "KG",
        "name": "Kilogram"
      },
      "is_active": true
    }
  ],
  "pagination": {...}
}
```

#### POST /api/v1/master-data/materials

**Permission**: `master_data:material:create`

**Request**:
```json
{
  "material_code": "RM-150",
  "name": "Hyaluronic Acid",
  "description": "High molecular weight HA for skin hydration",
  "material_type": "RAW_MATERIAL",
  "category_id": "cat-uuid",
  "inci_name": "Sodium Hyaluronate",
  "cas_number": "9067-32-7",
  "origin_country": "Japan",
  "is_organic": false,
  "is_natural": true,
  "storage_conditions": "Store in cool, dry place. Avoid direct sunlight.",
  "requires_cold_storage": false,
  "shelf_life_days": 730,
  "is_hazardous": false,
  "allergen_info": null,
  "standard_cost": 2500000,
  "currency": "VND",
  "base_uom_id": "kg-uuid",
  "purchase_uom_id": "kg-uuid",
  "stock_uom_id": "kg-uuid",
  "min_stock_quantity": 5,
  "max_stock_quantity": 50,
  "reorder_point": 10
}
```

**Response 201**:
```json
{
  "id": "new-material-uuid",
  "material_code": "RM-150",
  "name": "Hyaluronic Acid",
  "created_at": "2024-01-23T15:00:00Z"
}
```

---

### Products

```
GET    /api/v1/master-data/products
POST   /api/v1/master-data/products
GET    /api/v1/master-data/products/:id
PUT    /api/v1/master-data/products/:id
DELETE /api/v1/master-data/products/:id
```

#### POST /api/v1/master-data/products

**Permission**: `master_data:product:create`

**Request**:
```json
{
  "product_code": "FG-SERUM-001",
  "name": "Vitamin C Brightening Serum 30ml",
  "description": "Natural brightening serum with 20% Vitamin C",
  "category_id": "skincare-cat-uuid",
  "product_line": "Skincare",
  "volume": 30,
  "volume_unit": "ml",
  "target_skin_type": "All skin types, especially dull skin",
  "ingredients_summary": "Aqua, Ascorbic Acid (20%), Hyaluronic Acid, Vitamin E, Ferulic Acid",
  "usage_instructions": "Apply 2-3 drops to clean face daily, morning and evening",
  "standard_price": 450000,
  "recommended_retail_price": 550000,
  "currency": "VND",
  "packaging_type": "Glass bottle with dropper",
  "barcode": "8936012345678",
  "sku": "SERUM-VIT-C-30",
  "cosmetic_license_number": "123456/2024/CBMP",
  "license_expiry_date": "2029-12-31",
  "base_uom_id": "bottle-uuid",
  "sales_uom_id": "bottle-uuid",
  "shelf_life_months": 24,
  "image_urls": [
    "https://minio.../products/serum-001-front.jpg",
    "https://minio.../products/serum-001-back.jpg"
  ],
  "is_active": true,
  "launch_date": "2024-03-01"
}
```

---

### Categories

```
GET    /api/v1/master-data/categories
POST   /api/v1/master-data/categories
GET    /api/v1/master-data/categories/:id
PUT    /api/v1/master-data/categories/:id
DELETE /api/v1/master-data/categories/:id
```

#### GET /api/v1/master-data/categories

**Query Parameters**:
- `category_type`: MATERIAL, PRODUCT
- `parent_id`: uuid (null for root categories)
- `tree`: boolean (return as tree structure)

**Response with tree=true**:
```json
{
  "data": [
    {
      "id": "root-cat-uuid",
      "code": "ACTIVE-ING",
      "name": "Active Ingredients",
      "category_type": "MATERIAL",
      "children": [
        {
          "id": "child-cat-uuid",
          "code": "VITAMINS",
          "name": "Vitamins",
          "children": []
        }
      ]
    }
  ]
}
```

---

### Units of Measure

```
GET    /api/v1/master-data/uom
POST   /api/v1/master-data/uom
GET    /api/v1/master-data/uom/:id
PUT    /api/v1/master-data/uom/:id
DELETE /api/v1/master-data/uom/:id
POST   /api/v1/master-data/uom/convert
```

#### POST /api/v1/master-data/uom

**Request**:
```json
{
  "code": "G",
  "name": "Gram",
  "symbol": "g",
  "uom_type": "WEIGHT",
  "base_unit_id": "kg-uuid",
  "conversion_factor": 0.001
}
```

Example UoM conversions:
- 1 KG = 1000 G (conversion_factor of G = 0.001)
- 1 L = 1000 ML (conversion_factor of ML = 0.001)

#### POST /api/v1/master-data/uom/convert

Convert quantity giữa các UoM.

**Request**:
```json
{
  "value": 500,
  "from_uom_id": "g-uuid",
  "to_uom_id": "kg-uuid"
}
```

**Response**:
```json
{
  "original_value": 500,
  "original_uom": "G",
  "converted_value": 0.5,
  "converted_uom": "KG"
}
```

---

## gRPC METHODS

### GetMaterial

```protobuf
message GetMaterialRequest {
  string material_id = 1;
}

message GetMaterialResponse {
  Material material = 1;
}

message Material {
  string id = 1;
  string material_code = 2;
  string name = 3;
  string inci_name = 4;
  string material_type = 5;
  bool requires_cold_storage = 6;
  int32 shelf_life_days = 7;
  string base_uom_id = 8;
  bool is_active = 9;
}
```

### GetProduct

```protobuf
message GetProductRequest {
  string product_id = 1;
}

message GetProductResponse {
  Product product = 1;
}
```

### GetUoM

```protobuf
message GetUoMRequest {
  string uom_id = 1;
}

message GetUoMResponse {
  UoM uom = 1;
}

message UoM {
  string id = 1;
  string code = 2;
  string name = 3;
  string symbol = 4;
  double conversion_factor = 5;
}
```

### ConvertUoM

```protobuf
message ConvertUoMRequest {
  double value = 1;
  string from_uom_id = 2;
  string to_uom_id = 3;
}

message ConvertUoMResponse {
  double converted_value = 1;
}
```

### GetMaterialSuppliers

Lấy danh sách suppliers approved cho material.

```protobuf
message GetMaterialSuppliersRequest {
  string material_id = 1;
  bool preferred_only = 2;
}

message GetMaterialSuppliersResponse {
  repeated MaterialSupplier suppliers = 1;
}

message MaterialSupplier {
  string supplier_id = 1;
  string supplier_code = 2;
  string supplier_name = 3;
  bool is_preferred = 4;
  int32 lead_time_days = 5;
  double unit_price = 6;
}
```

---

## EVENTS

### Events Published

```yaml
master_data.material.created:
  payload:
    material_id: uuid
    material_code: string
    name: string
    material_type: string
    created_by: uuid

master_data.material.updated:
  payload:
    material_id: uuid
    updated_fields: array
    updated_by: uuid

master_data.material.deactivated:
  payload:
    material_id: uuid
    material_code: string
    reason: string

master_data.product.created:
  payload:
    product_id: uuid
    product_code: string
    name: string
    created_by: uuid

master_data.product.updated:
  payload:
    product_id: uuid
    updated_fields: array

master_data.product.deactivated:
  payload:
    product_id: uuid
    product_code: string
    reason: string

master_data.category.created:
  payload:
    category_id: uuid
    code: string
    name: string
    category_type: string

master_data.material.license_expiring:
  payload:
    product_id: uuid
    product_code: string
    license_number: string
    expiry_date: date
    days_until_expiry: int
```

### Events Subscribed

```yaml
supplier.created:
  action: Available for material_suppliers mapping

supplier.deactivated:
  action: Mark material_suppliers.is_approved = false for that supplier
```

---

## BUSINESS LOGIC

### Material Code Generation

Auto-generate nếu không được cung cấp:

```
Format: {TYPE}-{NNNN}
Examples:
  RM-0001 (Raw Material)
  PKG-0001 (Packaging)
  SF-0001 (Semi-Finished)
```

### Product Code Generation

```
Format: FG-{CATEGORY}-{NNNN}
Examples:
  FG-SERUM-0001
  FG-CREAM-0002
  FG-SHAMPOO-0003
```

### License Expiry Monitoring

- Check cosmetic license daily
- Send alert 90 days before expiry
- Send warning 30 days before expiry
- Auto-deactivate product on expiry date

### Material Approval Workflow

1. Create material (status: DRAFT)
2. Add specifications, SDS
3. Submit for approval
4. QC Manager review
5. Approve/Reject
6. If approved → is_active = true

---

## CONFIGURATION

```bash
# Server
MASTER_DATA_SERVICE_PORT=8083
MASTER_DATA_GRPC_PORT=9083

# Database
MASTER_DATA_DB_HOST=postgres
MASTER_DATA_DB_PORT=5434
MASTER_DATA_DB_NAME=master_data_db

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_BUCKET_SDS=sds-files
MINIO_BUCKET_PRODUCTS=product-images

# Business Rules
AUTO_GENERATE_CODES=true
LICENSE_EXPIRY_ALERT_DAYS=90,30,7
```

---

## MONITORING METRICS

```
master_data_materials_total{type="raw|packaging|semi"}
master_data_products_total{status="active|discontinued"}
master_data_categories_total
master_data_license_expiring_count
```

---

## DEPENDENCIES

- **Supplier Service** (gRPC): Validate supplier_id
- **File Service** (HTTP): Upload SDS files, product images
- **NATS**: Event publishing

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
