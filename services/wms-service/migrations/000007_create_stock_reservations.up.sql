-- Stock reservations table
CREATE TABLE IF NOT EXISTS stock_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL,
    lot_id UUID REFERENCES lots(id),
    location_id UUID REFERENCES locations(id),
    quantity DECIMAL(15,4) NOT NULL,
    unit_id UUID NOT NULL,
    reservation_type VARCHAR(30) NOT NULL, -- SALES_ORDER, WORK_ORDER, TRANSFER
    reference_id UUID NOT NULL,
    reference_number VARCHAR(30),
    status VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, RELEASED, FULFILLED, EXPIRED
    expires_at TIMESTAMP,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    released_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_reservations_material ON stock_reservations(material_id);
CREATE INDEX idx_reservations_reference ON stock_reservations(reference_id);
CREATE INDEX idx_reservations_status ON stock_reservations(status);
CREATE INDEX idx_reservations_expires ON stock_reservations(expires_at);
