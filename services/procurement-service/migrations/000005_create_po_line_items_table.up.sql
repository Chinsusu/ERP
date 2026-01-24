-- PO Line Items table
CREATE TABLE IF NOT EXISTS po_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    pr_line_item_id UUID,
    line_number INT NOT NULL,
    material_id UUID NOT NULL,
    material_code VARCHAR(50),
    material_name VARCHAR(255),
    quantity DECIMAL(15,4) NOT NULL,
    received_qty DECIMAL(15,4) NOT NULL DEFAULT 0,
    pending_qty DECIMAL(15,4) NOT NULL,
    uom_id UUID,
    uom_code VARCHAR(20),
    unit_price DECIMAL(15,4) NOT NULL,
    line_total DECIMAL(15,2) NOT NULL,
    tax_rate DECIMAL(5,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    expected_date DATE,
    specifications TEXT,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PARTIAL', 'COMPLETE', 'CANCELLED')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(po_id, line_number)
);

-- Indexes
CREATE INDEX idx_po_line_items_po_id ON po_line_items(po_id);
CREATE INDEX idx_po_line_items_material ON po_line_items(material_id);
CREATE INDEX idx_po_line_items_pr_line ON po_line_items(pr_line_item_id);
