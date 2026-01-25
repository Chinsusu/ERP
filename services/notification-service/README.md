# Notification Service

Notification Service quản lý tất cả notifications trong hệ thống ERP Cosmetics: email, in-app notifications, và alert rules.

## Features

- ✅ **Email Notifications** - SMTP integration với HTML templates
- ✅ **In-App Notifications** - Real-time notifications trong UI
- ✅ **Notification Templates** - Customizable templates với variables
- ✅ **Alert Rules** - Configurable automated alerts
- ✅ **Event Subscriptions** - Subscribe to events từ các services khác
- ✅ **Retry Logic** - Automatic retry với exponential backoff

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 16+
- NATS Server

### Environment Variables

```bash
# Service
SERVICE_NAME=notification-service
PORT=8090
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=erp_user
DB_PASSWORD=erp_password
DB_NAME=erp_notification

# NATS
NATS_URL=nats://localhost:4222

# SMTP (optional, uses mock sender if not configured)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=noreply@company.com
SMTP_FROM_NAME=ERP Cosmetics System
SMTP_USE_TLS=true
```

### Run Migrations

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=erp_user DB_PASSWORD=erp_password DB_NAME=erp_notification
make migrate-up
```

### Run Service

```bash
make run
# or
go run cmd/main.go
```

## API Endpoints

### Health Check

```
GET /health
GET /ready
```

### Notifications

```
POST /api/v1/notifications/send           # Send notification
```

### Templates

```
GET    /api/v1/notifications/templates       # List templates
POST   /api/v1/notifications/templates       # Create template
GET    /api/v1/notifications/templates/:id   # Get template
PUT    /api/v1/notifications/templates/:id   # Update template
DELETE /api/v1/notifications/templates/:id   # Delete template
```

### In-App Notifications

```
GET    /api/v1/notifications/in-app              # List user's notifications
POST   /api/v1/notifications/in-app              # Create notification
GET    /api/v1/notifications/in-app/unread-count # Get unread count
PATCH  /api/v1/notifications/in-app/read-all     # Mark all as read
PATCH  /api/v1/notifications/in-app/:id/read     # Mark as read
DELETE /api/v1/notifications/in-app/:id          # Delete notification
```

### Alert Rules

```
GET    /api/v1/alert-rules              # List rules
POST   /api/v1/alert-rules              # Create rule
GET    /api/v1/alert-rules/:id          # Get rule
PUT    /api/v1/alert-rules/:id          # Update rule
DELETE /api/v1/alert-rules/:id          # Delete rule
POST   /api/v1/alert-rules/:id/activate    # Activate rule
POST   /api/v1/alert-rules/:id/deactivate  # Deactivate rule
```

## Notification Types

| Type | Description |
|------|-------------|
| EMAIL | Email notification via SMTP |
| IN_APP | In-app notification displayed in UI |
| BOTH | Send both email and in-app notification |

## Alert Rule Types

| Type | Description |
|------|-------------|
| STOCK_LOW | Stock below reorder point |
| LOT_EXPIRY | Lot expiring within X days |
| CERT_EXPIRY | Certificate expiring within X days |
| APPROVAL_PENDING | Approval pending > X hours |
| TEMP_OUT_OF_RANGE | Temperature out of range (cold storage) |

## Event Subscriptions

Service subscribes to these events from other services:

| Event | Action |
|-------|--------|
| `wms.stock.low_stock_alert` | Send low stock notification |
| `wms.lot.expiring_soon` | Send lot expiry notification |
| `supplier.certification.expiring` | Send certificate expiry notification |
| `procurement.pr.submitted` | Notify approvers |
| `procurement.po.created` | Notify purchasing team |
| `manufacturing.qc.failed` | Notify production manager |
| `sales.order.confirmed` | Send order confirmation |

## Default Templates

Seed data includes these templates:

- `PASSWORD_RESET` - Password reset email
- `ACCOUNT_LOCKED` - Account locked alert
- `LOW_STOCK_ALERT` - Low stock notification
- `LOT_EXPIRING_ALERT` - Lot expiry notification
- `CERT_EXPIRING_ALERT` - Certificate expiry notification
- `PR_PENDING_APPROVAL` - PR approval notification
- `PO_CREATED` - PO created notification
- `QC_FAILED` - QC failed notification
- `ORDER_CONFIRMATION` - Order confirmation email

## Architecture

```
notification-service/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── domain/
│   │   ├── entity/            # Domain entities
│   │   └── repository/        # Repository interfaces
│   ├── infrastructure/
│   │   ├── persistence/       # PostgreSQL repositories
│   │   ├── email/            # SMTP sender
│   │   └── event/            # NATS pub/sub
│   ├── usecase/              # Business logic
│   └── delivery/
│       └── http/             # HTTP handlers
└── migrations/               # Database migrations
```

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## License

Proprietary - ERP Cosmetics System
