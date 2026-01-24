-- Zones table
CREATE TABLE IF NOT EXISTS zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(30) NOT NULL, -- RECEIVING, QUARANTINE, STORAGE, COLD, PICKING, SHIPPING
    temperature_min DECIMAL(5,2), -- For cold storage
    temperature_max DECIMAL(5,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(warehouse_id, code)
);

-- Create indexes
CREATE INDEX idx_zones_warehouse ON zones(warehouse_id);
CREATE INDEX idx_zones_type ON zones(zone_type);
