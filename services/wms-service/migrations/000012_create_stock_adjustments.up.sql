-- Stock Adjustments table
CREATE TABLE IF NOT EXISTS stock_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    adjustment_number VARCHAR(30) UNIQUE NOT NULL, -- ADJ-YYYY-XXXX
    adjustment_date DATE NOT NULL,
    adjustment_type VARCHAR(30) NOT NULL, -- CYCLE_COUNT, DAMAGE, EXPIRY, CORRECTION
    location_id UUID NOT NULL REFERENCES locations(id),
    material_id UUID NOT NULL,
    lot_id UUID REFERENCES lots(id),
    system_qty DECIMAL(15,4) NOT NULL,
    actual_qty DECIMAL(15,4) NOT NULL,
    adjustment_qty DECIMAL(15,4) NOT NULL, -- (actual - system)
    unit_id UUID NOT NULL,
    reason VARCHAR(100),
    notes TEXT,
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    created_by UUID NOT NULL,
    approved_by UUID,
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_adjustments_type ON stock_adjustments(adjustment_type);
CREATE INDEX idx_adjustments_status ON stock_adjustments(status);
CREATE INDEX idx_adjustments_location ON stock_adjustments(location_id);
CREATE INDEX idx_adjustments_material ON stock_adjustments(material_id);
CREATE INDEX idx_adjustments_date ON stock_adjustments(adjustment_date);
