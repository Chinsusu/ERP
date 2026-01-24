-- PR Line Items table
CREATE TABLE IF NOT EXISTS pr_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pr_id UUID NOT NULL REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    material_id UUID NOT NULL,
    material_code VARCHAR(50),
    material_name VARCHAR(255),
    quantity DECIMAL(15,4) NOT NULL,
    uom_id UUID,
    uom_code VARCHAR(20),
    unit_price DECIMAL(15,4) NOT NULL DEFAULT 0,
    line_total DECIMAL(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    required_date DATE,
    specifications TEXT,
    suggested_supplier_id UUID,
    suggested_supplier_name VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pr_id, line_number)
);

-- Indexes
CREATE INDEX idx_pr_line_items_pr_id ON pr_line_items(pr_id);
CREATE INDEX idx_pr_line_items_material ON pr_line_items(material_id);
