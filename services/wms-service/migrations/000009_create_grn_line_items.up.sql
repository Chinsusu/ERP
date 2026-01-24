-- GRN Line Items table
CREATE TABLE IF NOT EXISTS grn_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grn_id UUID NOT NULL REFERENCES grns(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    po_line_item_id UUID,
    material_id UUID NOT NULL,
    expected_qty DECIMAL(15,4),
    received_qty DECIMAL(15,4) NOT NULL,
    accepted_qty DECIMAL(15,4),
    rejected_qty DECIMAL(15,4) DEFAULT 0,
    unit_id UUID NOT NULL,
    lot_id UUID REFERENCES lots(id),
    supplier_lot_number VARCHAR(50),
    manufactured_date DATE,
    expiry_date DATE NOT NULL,
    location_id UUID REFERENCES locations(id),
    qc_status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, PASSED, FAILED
    qc_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(grn_id, line_number)
);

-- Indexes
CREATE INDEX idx_grn_items_grn ON grn_line_items(grn_id);
CREATE INDEX idx_grn_items_material ON grn_line_items(material_id);
CREATE INDEX idx_grn_items_lot ON grn_line_items(lot_id);
