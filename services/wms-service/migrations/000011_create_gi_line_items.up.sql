-- Goods Issue Line Items table
CREATE TABLE IF NOT EXISTS gi_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goods_issue_id UUID NOT NULL REFERENCES goods_issues(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    material_id UUID NOT NULL,
    requested_qty DECIMAL(15,4) NOT NULL,
    issued_qty DECIMAL(15,4) NOT NULL,
    unit_id UUID NOT NULL,
    lot_id UUID REFERENCES lots(id),
    location_id UUID REFERENCES locations(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(goods_issue_id, line_number)
);

-- Indexes
CREATE INDEX idx_gi_items_issue ON gi_line_items(goods_issue_id);
CREATE INDEX idx_gi_items_material ON gi_line_items(material_id);
CREATE INDEX idx_gi_items_lot ON gi_line_items(lot_id);
