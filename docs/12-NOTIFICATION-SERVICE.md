# 12 - NOTIFICATION SERVICE

## TỔNG QUAN

Notification Service quản lý tất cả notifications trong hệ thống: email, in-app notifications, SMS (future), và alert rules.

### Responsibilities

✅ Email notifications  
✅ In-app notifications  
✅ Push notifications (future)  
✅ SMS alerts (future)  
✅ Alert rules configuration  
✅ Notification templates  
✅ Notification history & read status

### Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP) + gRPC
- **Database**: PostgreSQL
- **Email**: SMTP (Gmail, SendGrid, or local SMTP server)
- **Queue**: NATS (async notification processing)

### Ports

- HTTP: `8090`
- gRPC: `9090`

---

## NOTIFICATION TYPES

### 1. Transactional Emails

Triggered by specific events:
- Password reset
- Account locked
- Order confirmation
- Shipment notification
- Invoice generated

### 2. Alert Notifications

System-generated alerts:
- Stock below reorder point
- Material expiring soon (30, 7 days)
- Certificate expiring (GMP, ISO)
- Work order delayed
- Quality check failed

### 3. Approval Notifications

Workflow approvals:
- PR pending approval
- PO pending approval
- BOM pending approval
- Sample request pending approval

### 4. In-App Notifications

Real-time notifications in the UI:
- New message
- Task assigned
- Document shared
- Mention in comment

---

## DATABASE SCHEMA

```sql
-- Notification templates
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY,
    template_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    notification_type VARCHAR(50), -- EMAIL, IN_APP, SMS, PUSH
    subject_template TEXT,
    body_template TEXT, -- HTML for email, plain text for others
    variables JSONB, -- Available variables for this template
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notification queue
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    notification_type VARCHAR(50) NOT NULL,
    template_id UUID REFERENCES notification_templates(id),
    
    recipient_user_id UUID,
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(50),
    
    subject TEXT,
    body TEXT,
    
    metadata JSONB, -- Event data, links, etc.
    
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, SENT, FAILED, RETRYING
    priority VARCHAR(20) DEFAULT 'NORMAL', -- HIGH, NORMAL, LOW
    
    sent_at TIMESTAMP,
    failed_at TIMESTAMP,
    retry_count INT DEFAULT 0,
    error_message TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(recipient_user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_type ON notifications(notification_type);

-- In-app notifications
CREATE TABLE user_notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(300) NOT NULL,
    message TEXT NOT NULL,
    notification_type VARCHAR(50), -- INFO, SUCCESS, WARNING, ERROR
    category VARCHAR(50), -- APPROVAL, ALERT, MESSAGE, SYSTEM
    
    link_url TEXT, -- Deep link to relevant page
    
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_notif_user ON user_notifications(user_id);
CREATE INDEX idx_user_notif_read ON user_notifications(is_read);

-- Alert rules
CREATE TABLE alert_rules (
    id UUID PRIMARY KEY,
    rule_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    rule_type VARCHAR(50), -- STOCK_LOW, EXPIRY, CERTIFICATE_EXPIRY, etc.
    conditions JSONB, -- Rule conditions
    
    notification_template_id UUID REFERENCES notification_templates(id),
    recipients JSONB, -- Array of user_ids or emails
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## API ENDPOINTS

### Templates

```
GET    /api/v1/notifications/templates
POST   /api/v1/notifications/templates
GET    /api/v1/notifications/templates/:id
PUT    /api/v1/notifications/templates/:id
```

### Send Notification

```
POST   /api/v1/notifications/send
```

#### POST /api/v1/notifications/send

**Permission**: `notification:send`

**Request**:
```json
{
  "notification_type": "EMAIL",
  "template_code": "ORDER_CONFIRMATION",
  "recipient_user_id": "user-uuid",
  "variables": {
    "customer_name": "Nguyễn Văn A",
    "order_number": "SO-2024-001",
    "total_amount": "44,550,000 VND",
    "delivery_date": "2024-02-05"
  },
  "priority": "HIGH"
}
```

### In-App Notifications

```
GET    /api/v1/notifications/in-app
POST   /api/v1/notifications/in-app
PATCH  /api/v1/notifications/in-app/:id/read
PATCH  /api/v1/notifications/in-app/read-all
DELETE /api/v1/notifications/in-app/:id
```

#### GET /api/v1/notifications/in-app

**Response**:
```json
{
  "data": [
    {
      "id": "notif-uuid",
      "title": "Purchase Order Pending Approval",
      "message": "PO-2024-015 from ABC Chemicals awaiting your approval",
      "type": "WARNING",
      "category": "APPROVAL",
      "link_url": "/procurement/purchase-orders/uuid",
      "is_read": false,
      "created_at": "2024-01-23T15:00:00Z"
    }
  ],
  "unread_count": 5
}
```

---

## NOTIFICATION TEMPLATES

### Password Reset Email

```html
Subject: Reset Your Password

Hi {{user_name}},

You requested to reset your password for ERP Cosmetics System.

Click the link below to reset your password:
{{reset_link}}

This link will expire in 1 hour.

If you didn't request this, please ignore this email.

Thanks,
ERP Cosmetics Team
```

### Stock Low Alert

```
Subject: Low Stock Alert: {{material_name}}

Material: {{material_code}} - {{material_name}}
Current Stock: {{current_quantity}} {{uom}}
Reorder Point: {{reorder_point}} {{uom}}
Shortage: {{shortage}} {{uom}}

Action Required: Please create a Purchase Requisition.

View Stock: {{stock_link}}
```

### Certificate Expiring Alert

```
Subject: Certificate Expiring Soon: {{supplier_name}}

Supplier: {{supplier_name}}
Certificate Type: {{certificate_type}}
Certificate Number: {{certificate_number}}
Expiry Date: {{expiry_date}}
Days Until Expiry: {{days_remaining}}

Action Required: Request updated certificate from supplier.

View Certificate: {{certificate_link}}
```

---

## EVENT SUBSCRIPTIONS

Notification Service subscribes to events from all services:

```yaml
# Auth Service
auth.user.password_changed:
  → Send: Password changed confirmation email

auth.user.account_locked:
  → Send: Account locked alert email

# Procurement
procurement.pr.submitted:
  → Send: Approval notification to manager

procurement.po.created:
  → Send: PO created notification to purchasing team

# WMS
wms.stock.low_stock_alert:
  → Send: Low stock alert to procurement team

wms.lot.expiring_soon:
  → Send: Expiry alert to warehouse manager

# Supplier
supplier.certification.expiring:
  → Send: Certificate expiry alert to procurement

# Manufacturing
manufacturing.wo.completed:
  → Send: Production completion notification

manufacturing.qc.failed:
  → Send: QC failure alert to production manager

# Sales
sales.order.confirmed:
  → Send: Order confirmation email to customer
```

---

## ALERT RULES CONFIGURATION

### Low Stock Alert

```json
{
  "rule_code": "STOCK_LOW_ALERT",
  "name": "Low Stock Alert",
  "rule_type": "STOCK_LOW",
  "conditions": {
    "check_interval": "1h",
    "threshold_type": "REORDER_POINT"
  },
  "notification_template_id": "stock-low-template-uuid",
  "recipients": [
    {"role": "PROCUREMENT_MANAGER"},
    {"role": "WAREHOUSE_MANAGER"}
  ]
}
```

### Certificate Expiry Alert

```json
{
  "rule_code": "CERT_EXPIRY_90_DAYS",
  "name": "Certificate Expiring in 90 Days",
  "rule_type": "CERTIFICATE_EXPIRY",
  "conditions": {
    "days_before_expiry": 90,
    "check_interval": "daily"
  },
  "notification_template_id": "cert-expiry-template-uuid",
  "recipients": [
    {"role": "PROCUREMENT_MANAGER"},
    {"user_id": "compliance-officer-uuid"}
  ]
}
```

---

## CONFIGURATION

```bash
NOTIFICATION_SERVICE_PORT=8090
NOTIFICATION_GRPC_PORT=9090

# Email SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=noreply@company.com
SMTP_PASSWORD=<secret>
SMTP_FROM_EMAIL=noreply@company.com
SMTP_FROM_NAME=ERP Cosmetics System

# Retry Configuration
MAX_RETRY_ATTEMPTS=3
RETRY_DELAY=5m

# Rate Limiting
MAX_EMAILS_PER_MINUTE=60
```

---

## MONITORING METRICS

```
notification_sent_total{type, status}
notification_failed_total{type, error_type}
notification_processing_duration_seconds
email_delivery_rate
```

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
