-- Lots table - Critical for FEFO and traceability
CREATE TABLE IF NOT EXISTS lots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lot_number VARCHAR(30) UNIQUE NOT NULL, -- LOT-YYYYMM-XXXX
    material_id UUID NOT NULL,
    supplier_id UUID,
    supplier_lot_number VARCHAR(50),
    manufactured_date DATE,
    expiry_date DATE NOT NULL,
    received_date DATE NOT NULL,
    grn_id UUID,
    qc_status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, PASSED, FAILED, QUARANTINE
    status VARCHAR(20) DEFAULT 'AVAILABLE', -- AVAILABLE, RESERVED, BLOCKED, EXPIRED
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Critical indexes for FEFO logic
CREATE INDEX idx_lots_expiry ON lots(expiry_date);
CREATE INDEX idx_lots_material ON lots(material_id);
CREATE INDEX idx_lots_status ON lots(status);
CREATE INDEX idx_lots_qc_status ON lots(qc_status);
CREATE INDEX idx_lots_supplier ON lots(supplier_id);
CREATE INDEX idx_lots_grn ON lots(grn_id);
