-- GRN (Goods Receipt Notes) table
CREATE TABLE IF NOT EXISTS grns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grn_number VARCHAR(30) UNIQUE NOT NULL, -- GRN-YYYY-XXXX
    grn_date DATE NOT NULL,
    po_id UUID,
    po_number VARCHAR(30),
    supplier_id UUID,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    delivery_note_number VARCHAR(50),
    vehicle_number VARCHAR(20),
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, IN_PROGRESS, COMPLETED, CANCELLED
    qc_status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, PASSED, FAILED
    qc_notes TEXT,
    notes TEXT,
    received_by UUID,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_grns_po ON grns(po_id);
CREATE INDEX idx_grns_status ON grns(status);
CREATE INDEX idx_grns_supplier ON grns(supplier_id);
CREATE INDEX idx_grns_warehouse ON grns(warehouse_id);
CREATE INDEX idx_grns_date ON grns(grn_date);
