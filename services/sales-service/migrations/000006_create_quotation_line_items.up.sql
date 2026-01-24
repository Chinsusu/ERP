-- Quotation line items
CREATE TABLE IF NOT EXISTS quotation_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quotation_id UUID NOT NULL REFERENCES quotations(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    product_id UUID NOT NULL,
    product_code VARCHAR(50),
    product_name VARCHAR(200),
    quantity DECIMAL(18,3) NOT NULL,
    uom_id UUID,
    unit_price DECIMAL(18,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(18,2) DEFAULT 0,
    tax_percent DECIMAL(5,2) DEFAULT 10,
    tax_amount DECIMAL(18,2) DEFAULT 0,
    line_total DECIMAL(18,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(quotation_id, line_number)
);

-- Create indexes
CREATE INDEX idx_quotation_items_quotation ON quotation_line_items(quotation_id);
CREATE INDEX idx_quotation_items_product ON quotation_line_items(product_id);
