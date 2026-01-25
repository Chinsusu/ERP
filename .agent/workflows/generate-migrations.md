---
description: Generate database migrations for a service
---

# Generate Database Migrations

## Usage
```
/generate-migrations [SERVICE_NAME]
```

## Steps

1. Read the service documentation from `/opt/ERP/docs/` to identify required tables
2. Create migration files in `services/[service-name]/migrations/`
3. Use golang-migrate format with naming: `YYYYMMDDHHMMSS_description.up.sql` / `.down.sql`

## Migration Requirements

- Include all indexes for commonly queried fields
- Include foreign keys with appropriate ON DELETE actions (CASCADE, SET NULL, RESTRICT)
- Add standard audit columns: `created_at`, `updated_at`, `deleted_at` for soft delete
- Use `gen_random_uuid()` for UUID primary keys
- Use `TIMESTAMP` with `DEFAULT CURRENT_TIMESTAMP` for timestamps

## Example Migration Structure

```sql
-- UP Migration
CREATE TABLE table_name (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- fields...
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_table_name_field ON table_name(field);

-- DOWN Migration  
DROP TABLE IF EXISTS table_name;
```

## Seed Data
If the service requires seed data (default roles, permissions, etc.), create a separate migration file with `_seed` suffix.
