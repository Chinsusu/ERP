# Supplier Service

Supplier Management Service for ERP Cosmetics System - Phase 2.1 Supply Chain.

## Status

✅ **Complete** - All components implemented and ready for testing.

## Features

- 📦 **Supplier Management** - CRUD with auto-generated codes (SUP-0001)
- 📍 **Addresses** - Multiple types (billing, shipping, factory)
- 👥 **Contacts** - Multiple contacts per supplier (sales, technical, QC)
- 📜 **Certifications** - GMP, ISO9001, ISO22716, Organic, Ecocert, Halal
- ⚠️ **Expiry Tracking** - Auto status update (VALID→EXPIRING_SOON→EXPIRED)
- ⭐ **Evaluations** - Quarterly/annual performance reviews
- 📊 **Rating** - Auto-calculated from evaluations

## Cosmetics Industry Features

| Feature | Description |
|---------|-------------|
| **GMP Compliance** | Track GMP certificates, block on expiry |
| **ISO 22716** | Cosmetics manufacturing standard |
| **Organic Certs** | Ecocert, COSMOS, USDA Organic |
| **Expiry Alerts** | 90/30/7 days before certificate expiry |
| **Performance Rating** | Quality, delivery, service scores |

## Ports

- **HTTP**: 8084
- **gRPC**: 9084 (planned)

## Quick Start

### 1. Create Database

```bash
PGPASSWORD=postgres123 psql -h localhost -U postgres -c "CREATE DATABASE supplier_db;"
```

### 2. Run Migrations

```bash
cd /opt/ERP/services/supplier-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres123 DB_NAME=supplier_db
make migrate-up
```

### 3. Run Service

```bash
make run
# Service runs on http://localhost:8084
```

## API Endpoints

### Suppliers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/suppliers` | List suppliers (filter by type, status) |
| POST | `/api/v1/suppliers` | Create supplier (status=PENDING) |
| GET | `/api/v1/suppliers/:id` | Get with addresses, contacts, certs |
| PUT | `/api/v1/suppliers/:id` | Update supplier info |
| PATCH | `/api/v1/suppliers/:id/approve` | Approve (status=APPROVED) |
| PATCH | `/api/v1/suppliers/:id/block` | Block with reason |

### Addresses & Contacts

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/suppliers/:id/addresses` | List addresses |
| POST | `/api/v1/suppliers/:id/addresses` | Add address |
| GET | `/api/v1/suppliers/:id/contacts` | List contacts |
| POST | `/api/v1/suppliers/:id/contacts` | Add contact |

### Certifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/suppliers/:id/certifications` | List certificates |
| POST | `/api/v1/suppliers/:id/certifications` | Add certificate |
| GET | `/api/v1/certifications/expiring?days=90` | Expiring certs |

### Evaluations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/suppliers/:id/evaluations` | List evaluations |
| POST | `/api/v1/suppliers/:id/evaluations` | Create evaluation |

## Database Schema

7 tables:
- `suppliers` - Master data with rating
- `supplier_addresses` - Multiple per supplier
- `supplier_contacts` - Multiple per supplier
- `supplier_certifications` - GMP, ISO, Organic with expiry
- `supplier_evaluations` - Performance reviews
- `approved_supplier_list` - Material-supplier mapping
- `supplier_price_lists` - Historical pricing

## Example Requests

### Create Supplier

```bash
curl -X POST http://localhost:8084/api/v1/suppliers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ABC Chemicals Vietnam",
    "legal_name": "Công ty TNHH ABC Chemicals",
    "tax_code": "0123456789",
    "supplier_type": "MANUFACTURER",
    "business_type": "DOMESTIC",
    "email": "sales@abc-chem.vn",
    "currency": "VND"
  }'
```

### Add Certification

```bash
curl -X POST http://localhost:8084/api/v1/suppliers/{id}/certifications \
  -H "Content-Type: application/json" \
  -d '{
    "certification_type": "GMP",
    "certificate_number": "GMP-VN-2024-001",
    "issuing_body": "Vietnam FDA",
    "issue_date": "2024-01-01",
    "expiry_date": "2027-01-01"
  }'
```

## Events Published (NATS)

- `supplier.created`
- `supplier.approved`
- `supplier.blocked`
- `supplier.certification.added`
- `supplier.certification.expiring`
- `supplier.certification.expired`
- `supplier.evaluation.completed`

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8084 | HTTP port |
| `DB_NAME` | supplier_db | Database name |
| `CERT_EXPIRY_ALERT_DAYS` | 90 | Alert threshold |
| `AUTO_BLOCK_ON_GMP_EXPIRY` | true | Block on GMP expiry |

---

**Port**: 8084 (HTTP), 9084 (gRPC)  
**Database**: `supplier_db`  
**Status**: Ready for testing
