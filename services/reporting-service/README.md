# Reporting Service

Reporting Service cung cấp dashboards, reports, và data exports cho ERP Cosmetics.

## Features

- ✅ **Dashboards** - Configurable dashboards với widgets
- ✅ **Reports** - Pre-built reports với parameters
- ✅ **Stats API** - Real-time KPIs và metrics
- ✅ **Export** - CSV, Excel exports
- ✅ **Widgets** - KPI, Charts, Tables

## Quick Start

```bash
# Set environment
export DB_HOST=localhost DB_PORT=5432 DB_USER=erp_user \
       DB_PASSWORD=erp_password DB_NAME=erp_reporting

# Run migrations
make migrate-up

# Run service
make run
```

## API Endpoints

### Dashboards
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/dashboards` | List dashboards |
| POST | `/api/v1/dashboards` | Create dashboard |
| GET | `/api/v1/dashboards/default` | Get default dashboard |
| GET | `/api/v1/dashboards/:id` | Get dashboard with widgets |
| PUT | `/api/v1/dashboards/:id` | Update dashboard |
| DELETE | `/api/v1/dashboards/:id` | Delete dashboard |
| POST | `/api/v1/dashboards/:id/widgets` | Add widget |

### Widgets
| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/api/v1/widgets/:id` | Update widget |
| DELETE | `/api/v1/widgets/:id` | Delete widget |

### Reports
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/reports` | List report definitions |
| GET | `/api/v1/reports/:id` | Get report definition |
| POST | `/api/v1/reports/:id/execute` | Execute report |
| GET | `/api/v1/reports/:id/executions` | List executions |
| GET | `/api/v1/reports/executions/:id` | Get execution |
| GET | `/api/v1/reports/executions/:id/download` | Download export |

### Stats
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/stats/dashboard` | All dashboard stats |
| GET | `/api/v1/stats/inventory` | Inventory KPIs |
| GET | `/api/v1/stats/sales` | Sales KPIs |
| GET | `/api/v1/stats/production` | Production KPIs |
| GET | `/api/v1/stats/procurement` | Procurement KPIs |

## Pre-built Reports

| Code | Type | Description |
|------|------|-------------|
| `STOCK_SUMMARY` | INVENTORY | Stock by warehouse/material |
| `EXPIRY_REPORT` | INVENTORY | Lots expiring soon |
| `STOCK_MOVEMENT` | INVENTORY | Stock ins/outs |
| `LOW_STOCK_ITEMS` | INVENTORY | Below reorder point |
| `PO_SUMMARY` | PROCUREMENT | PO status summary |
| `SUPPLIER_PERFORMANCE` | PROCUREMENT | Supplier ratings |
| `PRODUCTION_OUTPUT` | PRODUCTION | WO completion |
| `QC_SUMMARY` | PRODUCTION | QC pass/fail rates |
| `SALES_SUMMARY` | SALES | Sales by period |
| `TOP_PRODUCTS` | SALES | Best sellers |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8092 | HTTP port |
| `EXPORT_PATH` | /tmp/exports | Export file path |
| `MAX_EXPORT_ROWS` | 100000 | Max rows per export |
| `STATS_CACHE_TTL` | 300 | Stats cache (seconds) |

## Architecture

```
reporting-service/
├── cmd/main.go
├── internal/
│   ├── config/
│   ├── domain/
│   │   ├── entity/          # ReportDefinition, Dashboard, Widget
│   │   └── repository/
│   ├── infrastructure/
│   │   ├── persistence/     # PostgreSQL repositories
│   │   ├── export/          # CSV, Excel exporters
│   │   └── aggregator/      # Stats aggregation
│   ├── usecase/
│   │   ├── dashboard/
│   │   ├── report/
│   │   └── stats/
│   └── delivery/http/
└── migrations/
```

## License

Proprietary - ERP Cosmetics System
