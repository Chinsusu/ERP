-- Stock table with available qty calculation
CREATE TABLE IF NOT EXISTS stock (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    zone_id UUID NOT NULL REFERENCES zones(id),
    location_id UUID NOT NULL REFERENCES locations(id),
    material_id UUID NOT NULL,
    lot_id UUID REFERENCES lots(id),
    quantity DECIMAL(15,4) NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_qty DECIMAL(15,4) NOT NULL DEFAULT 0 CHECK (reserved_qty >= 0),
    available_qty DECIMAL(15,4) GENERATED ALWAYS AS (quantity - reserved_qty) STORED,
    unit_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(location_id, material_id, lot_id),
    CHECK (reserved_qty <= quantity)
);

-- Indexes for stock queries
CREATE INDEX idx_stock_material ON stock(material_id);
CREATE INDEX idx_stock_lot ON stock(lot_id);
CREATE INDEX idx_stock_location ON stock(location_id);
CREATE INDEX idx_stock_warehouse ON stock(warehouse_id);
CREATE INDEX idx_stock_zone ON stock(zone_id);
