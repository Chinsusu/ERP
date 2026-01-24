-- PO Receipts (linked to WMS GRN)
CREATE TABLE IF NOT EXISTS po_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    po_line_item_id UUID NOT NULL REFERENCES po_line_items(id),
    grn_id UUID,
    grn_number VARCHAR(50),
    received_qty DECIMAL(15,4) NOT NULL,
    received_date DATE NOT NULL DEFAULT CURRENT_DATE,
    received_by UUID,
    qc_status VARCHAR(20) DEFAULT 'PENDING' CHECK (qc_status IN ('PENDING', 'PASSED', 'FAILED', 'PARTIAL')),
    qc_notes TEXT,
    batch_number VARCHAR(100),
    lot_number VARCHAR(100),
    expiry_date DATE,
    storage_location VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_po_receipts_po_id ON po_receipts(po_id);
CREATE INDEX idx_po_receipts_line_item ON po_receipts(po_line_item_id);
CREATE INDEX idx_po_receipts_grn ON po_receipts(grn_id);
CREATE INDEX idx_po_receipts_date ON po_receipts(received_date);
