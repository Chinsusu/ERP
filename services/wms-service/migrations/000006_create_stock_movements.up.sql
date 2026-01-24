-- Stock movements table (audit trail)
CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movement_number VARCHAR(30) NOT NULL,
    movement_type VARCHAR(20) NOT NULL, -- IN, OUT, TRANSFER, ADJUSTMENT
    reference_type VARCHAR(30), -- GRN, GI, WO, TRANSFER, ADJUSTMENT, RESERVATION
    reference_id UUID,
    material_id UUID NOT NULL,
    lot_id UUID REFERENCES lots(id),
    from_location_id UUID REFERENCES locations(id),
    to_location_id UUID REFERENCES locations(id),
    quantity DECIMAL(15,4) NOT NULL,
    unit_id UUID NOT NULL,
    notes TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for movement queries
CREATE INDEX idx_movements_material ON stock_movements(material_id);
CREATE INDEX idx_movements_lot ON stock_movements(lot_id);
CREATE INDEX idx_movements_type ON stock_movements(movement_type);
CREATE INDEX idx_movements_date ON stock_movements(created_at);
CREATE INDEX idx_movements_reference ON stock_movements(reference_type, reference_id);
