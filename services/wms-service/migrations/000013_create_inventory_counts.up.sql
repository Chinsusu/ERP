-- Inventory Counts table
CREATE TABLE IF NOT EXISTS inventory_counts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    count_number VARCHAR(30) UNIQUE NOT NULL, -- CNT-YYYY-XXXX
    count_date DATE NOT NULL,
    count_type VARCHAR(20) NOT NULL, -- FULL, CYCLE, SPOT
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    zone_id UUID REFERENCES zones(id),
    status VARCHAR(20) DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, COMPLETED, CANCELLED
    notes TEXT,
    started_by UUID NOT NULL,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_inventory_counts_warehouse ON inventory_counts(warehouse_id);
CREATE INDEX idx_inventory_counts_status ON inventory_counts(status);
CREATE INDEX idx_inventory_counts_date ON inventory_counts(count_date);
