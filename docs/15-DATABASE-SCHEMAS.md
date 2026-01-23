# 15 - DATABASE SCHEMAS

## TỔNG QUAN

Tài liệu tổng hợp database schemas cho tất cả 15 services, indexes, constraints, và best practices.

---

## DATABASE STRATEGY

### Database Per Service

Mỗi service có PostgreSQL database riêng:

| Service | Database Name | Port |
|---------|---------------|------|
| Auth | `auth_db` | 5432 |
| User | `user_db` | 5433 |
| Master Data | `master_data_db` | 5434 |
| Supplier | `supplier_db` | 5435 |
| Procurement | `procurement_db` | 5436 |
| WMS | `wms_db` | 5437 |
| Manufacturing | `manufacturing_db` | 5438 |
| Sales | `sales_db` | 5439 |
| Marketing | `marketing_db` | 5440 |
| Finance | `finance_db` | 5441 |
| Reporting | `reporting_db` | 5442 |
| Notification | `notification_db` | 5443 |
| AI | `ai_db` | 5444 |

---

## COMMON PATTERNS

### UUID Primary Keys

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE example (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ...
);
```

### Timestamps

```sql
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
deleted_at TIMESTAMP  -- Soft delete
```

### Auto-update Trigger

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_example_updated_at
BEFORE UPDATE ON example
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
```

---

## INDEX STRATEGIES

### Essential Indexes

```sql
-- Foreign keys
CREATE INDEX idx_table_foreign_key ON table_name(foreign_key_id);

-- Query filters
CREATE INDEX idx_table_status ON table_name(status) WHERE deleted_at IS NULL;

-- Composite indexes (order matters!)
CREATE INDEX idx_table_composite ON table_name(user_id, created_at DESC);

-- Unique constraints
CREATE UNIQUE INDEX idx_table_unique ON table_name(code) WHERE deleted_at IS NULL;
```

### Partial Indexes

```sql
-- Only index active records
CREATE INDEX idx_active_suppliers ON suppliers(name) 
WHERE status = 'APPROVED' AND deleted_at IS NULL;
```

---

## DATA INTEGRITY

### Check Constraints

```sql
-- Ensure positive quantities
CHECK (quantity >= 0)

-- Ensure end date after start date
CHECK (end_date >= start_date)

-- Enum-like values
CHECK (status IN ('DRAFT', 'APPROVED', 'REJECTED'))
```

### Foreign Key Constraints

```sql
-- CASCADE delete
FOREIGN KEY (parent_id) REFERENCES parent_table(id) ON DELETE CASCADE

-- RESTRICT delete (default)
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT

-- SET NULL on delete
FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL
```

---

## COSMETICS-SPECIFIC SCHEMAS

### Material with INCI/CAS

```sql
CREATE TABLE materials (
    inci_name VARCHAR(300),  -- International nomenclature
    cas_number VARCHAR(50),  -- Chemical registry number
    is_organic BOOLEAN DEFAULT false,
    requires_cold_storage BOOLEAN DEFAULT false,
    shelf_life_days INT,
    allergen_info TEXT,
    safety_data_sheet_url TEXT
);
```

### Lot/Batch Tracking

```sql
CREATE TABLE lots (
    lot_number VARCHAR(100) NOT NULL UNIQUE,
    material_id UUID,
    supplier_lot_number VARCHAR(100), -- Traceability to supplier
    manufactured_date DATE,
    expiry_date DATE,
    qc_status VARCHAR(50),
    CONSTRAINT chk_expiry CHECK (expiry_date > manufactured_date)
);
```

### BOM with Encryption

```sql
CREATE TABLE boms (
    formula_details JSONB,  -- ENCRYPTED! AES-256-GCM
    confidentiality_level VARCHAR(50) DEFAULT 'RESTRICTED'
);

-- Encryption handled at application layer
-- Only users with 'manufacturing:bom:view_formula' permission can decrypt
```

---

---

## COMPREHENSIVE ER DIAGRAMS BY SERVICE GROUP

### Auth & User Services

```mermaid
erDiagram
    USERS ||--o{ USER_CREDENTIALS : has
    USERS ||--o{ USER_ROLES : assigned
    USERS }o--|| DEPARTMENTS : belongs_to
    USERS }o--o| USERS : manager
    DEPARTMENTS ||--o{ DEPARTMENTS : parent
    
    ROLES ||--o{ USER_ROLES : assigned_to
    ROLES ||--o{ ROLE_PERMISSIONS : has
    PERMISSIONS ||--o{ ROLE_PERMISSIONS : granted_in
    
    USER_CREDENTIALS ||--o{ REFRESH_TOKENS : generates
    USER_CREDENTIALS ||--o{ SESSIONS : maintains
    USER_CREDENTIALS ||--o{ PASSWORD_RESET_TOKENS : requests

    USERS {
        uuid id PK
        varchar employee_code UK
        varchar email UK
        varchar first_name
        varchar last_name
        uuid department_id FK
        uuid manager_id FK
        date hire_date
        varchar employment_status
    }
    
    USER_CREDENTIALS {
        uuid id PK
        uuid user_id FK,UK
        varchar email UK
        varchar password_hash
        boolean is_active
        int failed_login_attempts
        timestamp locked_until
    }
    
    ROLES {
        uuid id PK
        varchar name UK
        boolean is_system
    }
    
    PERMISSIONS {
        uuid id PK
        varchar code UK
        varchar service
        varchar resource
        varchar action
    }
    
    DEPARTMENTS {
        uuid id PK
        varchar code UK
        varchar name
        uuid parent_id FK
        uuid manager_id FK
    }
```

### Master Data & Supplier Services

```mermaid
erDiagram
    MATERIALS ||--o{ MATERIAL_SUPPLIERS : sourced_from
    SUPPLIERS ||--o{ MATERIAL_SUPPLIERS : supplies
    SUPPLIERS ||--o{ SUPPLIER_CERTIFICATIONS : has
    SUPPLIERS ||--o{ SUPPLIER_EVALUATIONS : evaluated_by
    
    MATERIALS }o--|| CATEGORIES : categorized_in
    PRODUCTS }o--|| CATEGORIES : categorized_in
    MATERIALS }o--|| UNITS_OF_MEASURE : measured_in
    PRODUCTS }o--|| UNITS_OF_MEASURE : measured_in
    
    CATEGORIES ||--o{ CATEGORIES : parent

    MATERIALS {
        uuid id PK
        varchar material_code UK
        varchar name
        varchar inci_name
        varchar cas_number
        varchar material_type
        uuid category_id FK
        uuid base_uom_id FK
        boolean is_organic
        boolean requires_cold_storage
        int shelf_life_days
        text allergen_info
    }
    
    SUPPLIERS {
        uuid id PK
        varchar supplier_code UK
        varchar name
        varchar tax_code
        varchar supplier_type
        varchar status
        decimal overall_rating
    }
    
    MATERIAL_SUPPLIERS {
        uuid id PK
        uuid material_id FK
        uuid supplier_id FK
        boolean is_preferred
        int lead_time_days
        decimal unit_price
    }
    
    SUPPLIER_CERTIFICATIONS {
        uuid id PK
        uuid supplier_id FK
        varchar certification_type
        varchar certificate_number
        date expiry_date
        varchar status
    }
    
    PRODUCTS {
        uuid id PK
        varchar product_code UK
        varchar name
        uuid category_id FK
        varchar cosmetic_license_number
        date license_expiry_date
        decimal standard_price
    }
```

### Procurement & WMS Services

```mermaid
erDiagram
    PURCHASE_REQUISITIONS ||--o{ PR_LINE_ITEMS : contains
    PURCHASE_ORDERS ||--o{ PO_LINE_ITEMS : contains
    PO_LINE_ITEMS ||--o{ GRN_LINE_ITEMS : received_via
    
    PR_LINE_ITEMS }o--|| MATERIALS : requests
    PO_LINE_ITEMS }o--|| MATERIALS : orders
    GRN_LINE_ITEMS }o--|| MATERIALS : receives
    
    PURCHASE_ORDERS }o--|| SUPPLIERS : ordered_from
    GRNS }o--|| PURCHASE_ORDERS : fulfills
    GRNS }o--|| WAREHOUSES : received_at
    
    GRNS ||--o{ GRN_LINE_ITEMS : contains
    GRN_LINE_ITEMS ||--|| LOTS : creates
    
    LOTS ||--o{ STOCK : stored_as
    STOCK }o--|| LOCATIONS : stored_at
    LOCATIONS }o--|| ZONES : in
    ZONES }o--|| WAREHOUSES : in
    
    STOCK ||--o{ STOCK_MOVEMENTS : tracked_by
    STOCK ||--o{ STOCK_RESERVATIONS : reserved_for

    PURCHASE_REQUISITIONS {
        uuid id PK
        varchar pr_number UK
        date pr_date
        uuid requested_by FK
        varchar status
        decimal total_amount
    }
    
    PURCHASE_ORDERS {
        uuid id PK
        varchar po_number UK
        uuid supplier_id FK
        date po_date
        date expected_delivery_date
        varchar status
        decimal total_amount
    }
    
    GRNS {
        uuid id PK
        varchar grn_number UK
        uuid po_id FK
        uuid warehouse_id FK
        date grn_date
        varchar status
    }
    
    LOTS {
        uuid id PK
        varchar lot_number UK
        uuid material_id FK
        date manufactured_date
        date expiry_date
        varchar qc_status
    }
    
    WAREHOUSES {
        uuid id PK
        varchar warehouse_code UK
        varchar name
        varchar warehouse_type
    }
    
    ZONES {
        uuid id PK
        varchar zone_code UK
        uuid warehouse_id FK
        varchar zone_type
    }
    
    LOCATIONS {
        uuid id PK
        varchar location_code UK
        uuid zone_id FK
    }
    
    STOCK {
        uuid id PK
        uuid location_id FK
        uuid lot_id FK
        decimal quantity
        decimal reserved_quantity
        decimal available_quantity
    }
```

### Manufacturing & Sales Services

```mermaid
erDiagram
    PRODUCTS ||--o{ BOMS : has
    BOMS ||--o{ BOM_LINE_ITEMS : contains
    BOM_LINE_ITEMS }o--|| MATERIALS : uses
    
    BOMS ||--o{ WORK_ORDERS : produces_via
    WORK_ORDERS ||--o{ WO_MATERIAL_ISSUES : consumes
    WO_MATERIAL_ISSUES }o--|| LOTS : from
    
    WORK_ORDERS ||--o{ QC_INSPECTIONS : inspected_by
    WORK_ORDERS ||--o{ BATCH_TRACEABILITY : tracked_by
    
    CUSTOMERS ||--o{ QUOTATIONS : requests
    CUSTOMERS ||--o{ SALES_ORDERS : places
    QUOTATIONS ||--o{ QUOTATION_LINE_ITEMS : contains
    SALES_ORDERS ||--o{ SO_LINE_ITEMS : contains
    SO_LINE_ITEMS }o--|| PRODUCTS : orders
    
    SALES_ORDERS ||--o{ STOCK_RESERVATIONS : reserves

    BOMS {
        uuid id PK
        varchar bom_number UK
        uuid product_id FK
        int version
        decimal batch_size
        jsonb formula_details
        varchar status
        decimal material_cost
    }
    
    WORK_ORDERS {
        uuid id PK
        varchar wo_number UK
        uuid product_id FK
        uuid bom_id FK
        varchar batch_number
        decimal planned_quantity
        decimal actual_quantity
        varchar status
    }
    
    QC_INSPECTIONS {
        uuid id PK
        varchar inspection_type
        uuid reference_id FK
        date inspection_date
        decimal inspected_quantity
        decimal accepted_quantity
        varchar result
        jsonb test_results
    }
    
    CUSTOMERS {
        uuid id PK
        varchar customer_code UK
        varchar name
        varchar customer_type
        decimal credit_limit
    }
    
    SALES_ORDERS {
        uuid id PK
        varchar so_number UK
        uuid customer_id FK
        date so_date
        date delivery_date
        varchar status
        decimal total_amount
    }
    
    BATCH_TRACEABILITY {
        uuid id PK
        uuid work_order_id FK
        uuid material_lot_id FK
        uuid product_lot_id FK
        decimal quantity_used
    }
```

---

## PERFORMANCE OPTIMIZATION

### Materialized Views

For reporting queries:

```sql
CREATE MATERIALIZED VIEW mv_stock_summary AS
SELECT 
    m.material_code,
    m.name,
    w.warehouse_name,
    SUM(s.quantity) as total_quantity,
    SUM(s.reserved_quantity) as total_reserved
FROM stock s
JOIN materials m ON s.material_id = m.id
JOIN locations l ON s.location_id = l.id
JOIN warehouses w ON l.warehouse_id = w.id
GROUP BY m.material_code, m.name, w.warehouse_name;

-- Refresh schedule
REFRESH MATERIALIZED VIEW CONCURRENTLY mv_stock_summary;
```

### Partitioning

For large transaction tables:

```sql
-- Partition by month
CREATE TABLE stock_movements (
    id UUID,
    movement_date DATE NOT NULL,
    ...
) PARTITION BY RANGE (movement_date);

CREATE TABLE stock_movements_2024_01 PARTITION OF stock_movements
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
```

---

## BACKUP & RECOVERY

### Daily Backup Script

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/postgres"

for DB in auth_db user_db master_data_db procurement_db wms_db manufacturing_db sales_db
do
    pg_dump -h postgres -U postgres $DB | gzip > $BACKUP_DIR/${DB}_${DATE}.sql.gz
done

# Retention: Keep 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete
```

---

## MIGRATION STRATEGY

### Tools

- **golang-migrate**: Database migration tool
- **Flyway** (alternative): Database migration tool

### Migration Files

```
migrations/
  000001_create_users_table.up.sql
  000001_create_users_table.down.sql
  000002_add_users_email_index.up.sql
  000002_add_users_email_index.down.sql
```

### Example Migration

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 000001_create_users_table.down.sql
DROP TABLE users;
```

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-23  
**Author**: ERP Development Team
