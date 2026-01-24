-- Warehouses table
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    warehouse_type VARCHAR(20) NOT NULL, -- MAIN, COLD_STORAGE, FINISHED_GOODS, QUARANTINE
    address TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index
CREATE INDEX idx_warehouses_code ON warehouses(code);
CREATE INDEX idx_warehouses_type ON warehouses(warehouse_type);
