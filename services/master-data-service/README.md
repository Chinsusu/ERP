# Master Data Service

Master Data Service for the ERP Cosmetics System - manages materials, products, categories, and units of measure.

## Status

✅ **Complete** - All components implemented and ready for testing.

## Features

- 📦 **Materials Management** - Raw materials, packaging, consumables with INCI/CAS support
- 🏷️ **Products Management** - Finished goods with cosmetic license tracking
- 📂 **Categories** - Hierarchical categorization with tree structure
- ⚖️ **Units of Measure** - UoM with conversion support

## Cosmetics Industry Features

- **INCI Names** - International Nomenclature of Cosmetic Ingredients
- **CAS Numbers** - Chemical Abstracts Service registry
- **Allergen Tracking** - Safety information for materials
- **Storage Conditions** - Ambient, Cold (2-8°C), Frozen (<-18°C)
- **Cosmetic License** - License number and expiry tracking

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP)
- **Database**: PostgreSQL
- **Message Queue**: NATS JetStream

## Ports

- **HTTP**: 8083
- **gRPC**: 9083 (planned)

## Quick Start

### 1. Create Database

```bash
PGPASSWORD=postgres123 psql -h localhost -U postgres -c "CREATE DATABASE master_data_db;"
```

### 2. Run Migrations

```bash
cd /opt/ERP/services/master-data-service
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres123 DB_NAME=master_data_db
make migrate-up
```

### 3. Run Service

```bash
make run
# Service runs on http://localhost:8083
```

## API Endpoints

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/categories` | List categories |
| GET | `/api/v1/categories/tree` | Get hierarchical tree |
| POST | `/api/v1/categories` | Create category |
| GET | `/api/v1/categories/:id` | Get by ID |
| PUT | `/api/v1/categories/:id` | Update |
| DELETE | `/api/v1/categories/:id` | Delete |

### Units of Measure

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/units` | List units |
| POST | `/api/v1/units` | Create unit |
| GET | `/api/v1/units/:id` | Get by ID |
| POST | `/api/v1/units/convert` | Convert value |

### Materials

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/materials` | List materials |
| GET | `/api/v1/materials/search?q=vitamin` | Search |
| POST | `/api/v1/materials` | Create |
| GET | `/api/v1/materials/:id` | Get by ID |
| PUT | `/api/v1/materials/:id` | Update |
| DELETE | `/api/v1/materials/:id` | Delete |
| POST | `/api/v1/materials/:id/specifications` | Add spec |

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | List products |
| GET | `/api/v1/products/search?q=serum` | Search |
| POST | `/api/v1/products` | Create |
| GET | `/api/v1/products/:id` | Get by ID |
| PUT | `/api/v1/products/:id` | Update |
| DELETE | `/api/v1/products/:id` | Delete |
| POST | `/api/v1/products/:id/images` | Add image |

## Database Schema

7 tables:
- `categories` - Hierarchical categories (materialized path)
- `units_of_measure` - Units with conversion factors
- `unit_conversions` - Bidirectional conversions
- `materials` - Raw materials and packaging
- `material_specifications` - Extended specs and certificates
- `products` - Finished goods
- `product_images` - Product images

## Example Requests

### Create Material

```bash
curl -X POST http://localhost:8083/api/v1/materials \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Vitamin C (Ascorbic Acid)",
    "material_type": "RAW_MATERIAL",
    "inci_name": "Ascorbic Acid",
    "cas_number": "50-81-7",
    "storage_condition": "COLD",
    "base_unit_id": "a0000000-0000-0000-0000-000000000001"
  }'
```

### Convert Units

```bash
curl -X POST http://localhost:8083/api/v1/units/convert \
  -H "Content-Type: application/json" \
  -d '{
    "value": 1000,
    "from_unit_id": "a0000000-0000-0000-0000-000000000002",
    "to_unit_id": "a0000000-0000-0000-0000-000000000001"
  }'
```

## Events Published

- `master_data.material.created`
- `master_data.material.updated`
- `master_data.product.created`
- `master_data.product.updated`
- `master_data.category.created`

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8083 | HTTP port |
| DB_HOST | localhost | PostgreSQL host |
| DB_NAME | master_data_db | Database name |
| NATS_URL | nats://localhost:4222 | NATS URL |
| AUTO_GENERATE_CODES | true | Auto-generate material/product codes |

## Project Structure

```
master-data-service/
├── cmd/main.go                          # Entry point
├── internal/
│   ├── config/                          # Configuration
│   ├── domain/
│   │   ├── entity/                      # Domain entities
│   │   └── repository/                  # Repository interfaces
│   ├── usecase/                         # Business logic
│   │   ├── category/
│   │   ├── unit/
│   │   ├── material/
│   │   └── product/
│   ├── delivery/http/                   # HTTP handlers
│   │   ├── handler/
│   │   ├── dto/
│   │   └── router.go
│   └── infrastructure/                  # External implementations
│       ├── persistence/postgres/
│       └── event/
└── migrations/                          # SQL migrations
```

---

**Port**: 8083 (HTTP), 9083 (gRPC)  
**Database**: `master_data_db`  
**Status**: Ready for testing
