# 06 - SUPPLIER SERVICE

## TỔNG QUAN

Supplier Service quản lý toàn bộ thông tin nhà cung cấp, từ thông tin cơ bản, liên hệ, chứng nhận (GMP, ISO, Organic) đến đánh giá định kỳ và Approved Supplier List (ASL).

### Responsibilities

✅ Supplier master data management  
✅ Supplier contacts & addresses  
✅ Certifications tracking (GMP, ISO22716, Organic, Ecocert, Halal)  
✅ Certificate expiry alerts  
✅ Supplier evaluation & rating  
✅ Approved Supplier List (ASL)  
✅ Supplier performance metrics

### Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP) + gRPC
- **Database**: PostgreSQL
- **File Storage**: MinIO (certificates)

### Ports

- HTTP: `8084`
- gRPC: `9084`

---

## WHY SUPPLIER MANAGEMENT IS CRITICAL FOR COSMETICS

### 1. GMP Compliance

**GMP (Good Manufacturing Practice)** là yêu cầu bắt buộc cho ngành mỹ phẩm. Nhà cung cấp nguyên liệu phải có chứng nhận GMP hợp lệ.

**ISO 22716**: Cosmetics - Good Manufacturing Practices (GMP)

### 2. Material Traceability

Theo dõi nguồn gốc nguyên liệu từ nhà cung cấp → sản xuất → khách hàng.

### 3. Quality Assurance

- Only approved suppliers in ASL
- Certificate validity monitoring
- Supplier rating system
- Periodic evaluation

### 4. Organic & Natural Certifications

Nếu sản phẩm claim "Organic" hoặc "Natural", nguyên liệu phải có chứng nhận:
- **Ecocert**: European organic standard
- **USDA Organic**: US organic standard
- **Cosmos**: COSMetic Organic Standard

---

## DATABASE SCHEMA (Key Tables)

Chi tiết schema đã có trong `02-SERVICE-SPECIFICATIONS.md`. Summary:

### Core Tables

- `suppliers` - Supplier master data
- `supplier_addresses` - Multiple addresses (billing, shipping, factory)
- `supplier_contacts` - Key contacts (sales, technical, QC, accounting)
- `supplier_certifications` - GMP, ISO, Organic certificates với expiry tracking
- `supplier_evaluations` - Periodic performance reviews

### Key Features

**suppliers table**:
- `supplier_code`, `name`, `tax_code`
- `supplier_type`: MANUFACTURER, TRADER, IMPORTER
- `quality_rating`, `delivery_rating`, `service_rating`
- `status`: PENDING, APPROVED, BLOCKED, INACTIVE
- `payment_terms`, `credit_limit`

**supplier_certifications table**:
- `certification_type`: GMP, ISO9001, ISO22716, ORGANIC, ECOCERT, HALAL
- `certificate_number`, `issuing_body`
- `issue_date`, `expiry_date`
- `status`: VALID, EXPIRING_SOON, EXPIRED
- `certificate_file_url` (MinIO)

**supplier_evaluations table**:
- Quarterly/annual reviews
- Scores: quality, delivery, price, service, documentation
- On-time delivery rate, quality acceptance rate
- Strengths, weaknesses, action items

---

## API ENDPOINTS

### Suppliers

#### GET /api/v1/suppliers

**Permission**: `supplier:supplier:read`

**Query Parameters**:
- `supplier_type`: MANUFACTURER, TRADER, IMPORTER
- `status`: PENDING, APPROVED, BLOCKED, INACTIVE
- `country`: string
- `has_gmp`: boolean (có GMP certificate không)
- `min_rating`: decimal (1-5)
- `search`: string (name, code, tax_code)
- `page`, `limit`

**Response 200**:
```json
{
  "data": [
    {
      "id": "supplier-uuid",
      "supplier_code": "SUP-001",
      "name": "ABC Chemicals Vietnam",
      "legal_name": "Công ty TNHH ABC Chemicals Vietnam",
      "tax_code": "0123456789",
      "supplier_type": "MANUFACTURER",
      "business_type": "DOMESTIC",
      "email": "sales@abc-chemicals.vn",
      "phone": "+84 28 1234 5678",
      "website": "https://abc-chemicals.vn",
      "status": "APPROVED",
      "overall_rating": 4.5,
      "has_valid_gmp": true,
      "created_at": "2023-01-15T00:00:00Z"
    }
  ],
  "pagination": {...}
}
```

---

#### POST /api/v1/suppliers

**Permission**: `supplier:supplier:create`

**Request**:
```json
{
  "supplier_code": "SUP-150",
  "name": "Natural Ingredients Ltd",
  "legal_name": "Natural Ingredients Company Limited",
  "tax_code": "9876543210",
  "supplier_type": "MANUFACTURER",
  "business_type": "INTERNATIONAL",
  "website": "https://naturalingredients.com",
  "email": "contact@naturalingredients.com",
  "phone": "+66 2 1234 5678",
  "payment_terms": "Net 30",
  "currency": "USD",
  "bank_name": "Bangkok Bank",
  "bank_account": "123-456-789",
  "credit_limit": 50000,
  "notes": "Premium organic ingredient supplier"
}
```

**Response 201**:
```json
{
  "id": "new-supplier-uuid",
  "supplier_code": "SUP-150",
  "name": "Natural Ingredients Ltd",
  "status": "PENDING",
  "created_at": "2024-01-23T15:00:00Z"
}
```

**Note**: New suppliers start with status=PENDING, cần approval.

---

#### GET /api/v1/suppliers/:id

**Response 200**:
```json
{
  "id": "supplier-uuid",
  "supplier_code": "SUP-001",
  "name": "ABC Chemicals Vietnam",
  "legal_name": "Công ty TNHH ABC Chemicals Vietnam",
  "tax_code": "0123456789",
  "supplier_type": "MANUFACTURER",
  "business_type": "DOMESTIC",
  "website": "https://abc-chemicals.vn",
  "email": "sales@abc-chemicals.vn",
  "phone": "+84 28 1234 5678",
  "fax": "+84 28 1234 5679",
  "payment_terms": "Net 30",
  "currency": "VND",
  "credit_limit": 100000000,
  "status": "APPROVED",
  "approved_by": "admin-user-uuid",
  "approved_at": "2023-02-01T00:00:00Z",
  "quality_rating": 5,
  "delivery_rating": 4,
  "service_rating": 5,
  "overall_rating": 4.67,
  "certifications_count": 3,
  "valid_certifications": ["GMP", "ISO9001"],
  "expiring_certifications": [],
  "notes": "Reliable supplier for active ingredients",
  "created_at": "2023-01-15T00:00:00Z",
  "updated_at": "2024-01-20T00:00:00Z"
}
```

---

#### PATCH /api/v1/suppliers/:id/approve

Approve supplier vào ASL.

**Permission**: `supplier:supplier:approve`

**Request**:
```json
{
  "notes": "All documents verified. GMP certificate valid until 2026."
}
```

**Response 200**:
```json
{
  "id": "supplier-uuid",
  "status": "APPROVED",
  "approved_by": "user-uuid",
  "approved_at": "2024-01-23T15:00:00Z"
}
```

---

#### PATCH /api/v1/suppliers/:id/block

Block supplier (tạm ngưng giao dịch).

**Permission**: `supplier:supplier:approve`

**Request**:
```json
{
  "reason": "Quality issues - 3 consecutive batches rejected"
}
```

**Response 200**:
```json
{
  "id": "supplier-uuid",
  "status": "BLOCKED",
  "notes": "Quality issues - 3 consecutive batches rejected"
}
```

---

### Addresses

#### GET /api/v1/suppliers/:id/addresses

**Response 200**:
```json
{
  "data": [
    {
      "id": "addr-uuid",
      "address_type": "FACTORY",
      "address_line1": "123 Industrial Park",
      "ward": "Ward 5",
      "district": "Binh Thanh",
      "city": "Ho Chi Minh City",
      "country": "Vietnam",
      "postal_code": "70000",
      "is_primary": true
    },
    {
      "id": "addr-uuid-2",
      "address_type": "BILLING",
      "address_line1": "456 Main Street",
      "city": "Ho Chi Minh City",
      "is_primary": false
    }
  ]
}
```

---

#### POST /api/v1/suppliers/:id/addresses

**Request**:
```json
{
  "address_type": "SHIPPING",
  "address_line1": "789 Port Road",
  "district": "District 4",
  "city": "Ho Chi Minh City",
  "country": "Vietnam",
  "postal_code": "70000",
  "is_primary": false
}
```

---

### Contacts

#### GET /api/v1/suppliers/:id/contacts

**Response 200**:
```json
{
  "data": [
    {
      "id": "contact-uuid",
      "contact_type": "PRIMARY",
      "full_name": "Nguyễn Văn A",
      "position": "Sales Manager",
      "department": "Sales",
      "email": "nguyen.a@abc-chemicals.vn",
      "phone": "+84 28 1234 5678",
      "mobile": "+84 901 234 567",
      "is_primary": true
    },
    {
      "id": "contact-uuid-2",
      "contact_type": "TECHNICAL",
      "full_name": "Trần Thị B",
      "position": "Technical Support",
      "email": "tran.b@abc-chemicals.vn",
      "is_primary": false
    }
  ]
}
```

---

#### POST /api/v1/suppliers/:id/contacts

**Request**:
```json
{
  "contact_type": "QUALITY",
  "full_name": "Lê Văn C",
  "position": "QC Manager",
  "department": "Quality Control",
  "email": "le.c@abc-chemicals.vn",
  "phone": "+84 28 1234 5680",
  "mobile": "+84 902 345 678",
  "is_primary": false
}
```

---

### Certifications

#### GET /api/v1/suppliers/:id/certifications

**Response 200**:
```json
{
  "data": [
    {
      "id": "cert-uuid",
      "certification_type": "GMP",
      "certificate_number": "GMP-VN-2023-001",
      "issuing_body": "Vietnam Food Administration",
      "issue_date": "2023-06-01",
      "expiry_date": "2026-05-31",
      "status": "VALID",
      "days_until_expiry": 850,
      "certificate_file_url": "https://minio.../certificates/gmp-cert.pdf",
      "verified_by": "user-uuid",
      "verification_date": "2023-06-15"
    },
    {
      "id": "cert-uuid-2",
      "certification_type": "ISO9001",
      "certificate_number": "ISO-2024-456",
      "issuing_body": "TUV Nord",
      "issue_date": "2024-01-01",
      "expiry_date": "2025-12-31",
      "status": "EXPIRING_SOON",
      "days_until_expiry": 60,
      "certificate_file_url": "https://minio.../certificates/iso9001.pdf"
    }
  ]
}
```

---

#### POST /api/v1/suppliers/:id/certifications

Upload certificate.

**Permission**: `supplier:certification:manage`

**Request (multipart/form-data)**:
```
certification_type: GMP
certificate_number: GMP-VN-2024-002
issuing_body: Vietnam Food Administration
issue_date: 2024-01-20
expiry_date: 2027-01-19
file: certificate.pdf (max 10MB)
notes: New GMP certificate after audit
```

**Response 201**:
```json
{
  "id": "new-cert-uuid",
  "certification_type": "GMP",
  "certificate_number": "GMP-VN-2024-002",
  "expiry_date": "2027-01-19",
  "status": "VALID",
  "certificate_file_url": "https://minio.../certificates/cert-uuid.pdf"
}
```

---

#### GET /api/v1/suppliers/certifications/expiring

Danh sách certificates sắp hết hạn (trong 90 ngày).

**Permission**: `supplier:certification:manage`

**Query Parameters**:
- `days`: int (default: 90) - số ngày trước expiry

**Response 200**:
```json
{
  "data": [
    {
      "supplier_id": "supplier-uuid",
      "supplier_name": "ABC Chemicals",
      "certification_type": "ISO9001",
      "certificate_number": "ISO-2024-456",
      "expiry_date": "2025-12-31",
      "days_until_expiry": 60,
      "status": "EXPIRING_SOON"
    }
  ],
  "total": 5
}
```

---

### Evaluations

#### GET /api/v1/suppliers/:id/evaluations

**Response 200**:
```json
{
  "data": [
    {
      "id": "eval-uuid",
      "evaluation_date": "2024-01-15",
      "evaluation_period": "Q1-2024",
      "quality_score": 5,
      "delivery_score": 4,
      "price_score": 4,
      "service_score": 5,
      "documentation_score": 5,
      "overall_score": 4.6,
      "on_time_delivery_rate": 95.5,
      "quality_acceptance_rate": 98.2,
      "strengths": "Excellent quality, responsive technical support",
      "weaknesses": "Occasional delivery delays",
      "action_items": "Improve delivery scheduling",
      "evaluated_by": "user-uuid",
      "status": "APPROVED"
    }
  ]
}
```

---

#### POST /api/v1/suppliers/:id/evaluations

**Permission**: `supplier:evaluation:create`

**Request**:
```json
{
  "evaluation_date": "2024-01-23",
  "evaluation_period": "2023-Q4",
  "quality_score": 5,
  "delivery_score": 4,
  "price_score": 4,
  "service_score": 5,
  "documentation_score": 5,
  "on_time_delivery_rate": 92.3,
  "quality_acceptance_rate": 97.8,
  "lead_time_adherence": 85.0,
  "strengths": "Consistent quality, good communication",
  "weaknesses": "Price increases in Q4",
  "action_items": "Negotiate fixed price contract for 2024"
}
```

---

## gRPC METHODS

### GetSupplier

```protobuf
message GetSupplierRequest {
  string supplier_id = 1;
}

message GetSupplierResponse {
  Supplier supplier = 1;
}

message Supplier {
  string id = 1;
  string supplier_code = 2;
  string name = 3;
  string supplier_type = 4;
  string status = 5;
  double overall_rating = 6;
  bool has_valid_gmp = 7;
}
```

---

### ValidateSupplierActive

```protobuf
message ValidateSupplierActiveRequest {
  string supplier_id = 1;
}

message ValidateSupplierActiveResponse {
  bool is_active = 1;
  bool is_approved = 2;
  string status = 3;
}
```

---

### GetSupplierCertifications

```protobuf
message GetSupplierCertificationsRequest {
  string supplier_id = 1;
  bool valid_only = 2; // Only return valid certificates
}

message GetSupplierCertificationsResponse {
  repeated Certification certifications = 1;
}

message Certification {
  string certification_type = 1;
  string certificate_number = 2;
  string expiry_date = 3;
  string status = 4;
}
```

---

## EVENTS

### Events Published

```yaml
supplier.created:
  payload:
    supplier_id: uuid
    supplier_code: string
    name: string
    created_by: uuid

supplier.approved:
  payload:
    supplier_id: uuid
    supplier_code: string
    name: string
    approved_by: uuid
    approved_at: datetime

supplier.blocked:
  payload:
    supplier_id: uuid
    supplier_code: string
    reason: string
    blocked_by: uuid

supplier.certification.added:
  payload:
    supplier_id: uuid
    certification_type: string
    certificate_number: string
    expiry_date: date

supplier.certification.expiring:
  payload:
    supplier_id: uuid
    supplier_name: string
    certification_type: string
    certificate_number: string
    expiry_date: date
    days_until_expiry: int
  # Sent 90, 30, 7 days before expiry

supplier.certification.expired:
  payload:
    supplier_id: uuid
    supplier_name: string
    certification_type: string
    certificate_number: string
    expired_date: date

supplier.evaluation.completed:
  payload:
    supplier_id: uuid
    evaluation_id: uuid
    evaluation_period: string
    overall_score: decimal
    evaluated_by: uuid
```

---

## BUSINESS LOGIC

### Supplier Code Generation

```
Format: SUP-{NNNN}
Example: SUP-0001, SUP-0002
```

### Approval Workflow

1. Create supplier (status = PENDING)
2. Add addresses, contacts
3. Upload certifications (GMP required)
4. Submit for approval
5. Procurement Manager reviews
6. Approve → status = APPROVED (added to ASL)

### Certificate Expiry Monitoring

**Daily Job**:
- Check all certificates
- Update status:
  - VALID: > 90 days to expiry
  - EXPIRING_SOON: ≤ 90 days to expiry
  - EXPIRED: past expiry date

**Alerts**:
- 90 days before: Email to Procurement team
- 30 days before: Email + in-app notification
- 7 days before: Daily email
- On expiry: Block supplier if GMP expired

### Supplier Rating Calculation

```
overall_rating = (
    quality_score + 
    delivery_score + 
    service_score
) / 3

Updated after each evaluation
```

### Evaluation Schedule

- **Quarterly**: For critical suppliers (Class A)
- **Bi-annually**: For regular suppliers (Class B)
- **Annually**: For occasional suppliers (Class C)

---

## CONFIGURATION

```bash
# Server
SUPPLIER_SERVICE_PORT=8084
SUPPLIER_GRPC_PORT=9084

# Database
SUPPLIER_DB_HOST=postgres
SUPPLIER_DB_PORT=5435
SUPPLIER_DB_NAME=supplier_db

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_BUCKET_CERTIFICATES=supplier-certificates

# Certificate Settings
CERTIFICATE_EXPIRY_ALERT_DAYS=90,30,7
AUTO_BLOCK_ON_GMP_EXPIRY=true

# Evaluation
REQUIRE_ANNUAL_EVALUATION=true
```

---

## MONITORING METRICS

```
supplier_service_suppliers_total{status="approved|pending|blocked"}
supplier_service_certifications_expiring_count
supplier_service_evaluations_total
supplier_service_avg_supplier_rating
```

---

## DEPENDENCIES

- **File Service** (HTTP): Upload certificates
- **Notification Service** (Event): Send expiry alerts
- **Master Data Service** (Subscribe): Material-supplier mappings
- **NATS**: Event publishing

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
