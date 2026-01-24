-- Sales Order line items
CREATE TABLE IF NOT EXISTS so_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sales_order_id UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    product_id UUID NOT NULL,
    product_code VARCHAR(50),
    product_name VARCHAR(200),
    quantity DECIMAL(18,3) NOT NULL,
    shipped_quantity DECIMAL(18,3) DEFAULT 0,
    uom_id UUID,
    unit_price DECIMAL(18,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(18,2) DEFAULT 0,
    tax_percent DECIMAL(5,2) DEFAULT 10,
    tax_amount DECIMAL(18,2) DEFAULT 0,
    line_total DECIMAL(18,2) DEFAULT 0,
    reservation_id UUID,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sales_order_id, line_number)
);

-- Create indexes
CREATE INDEX idx_so_items_order ON so_line_items(sales_order_id);
CREATE INDEX idx_so_items_product ON so_line_items(product_id);
CREATE INDEX idx_so_items_reservation ON so_line_items(reservation_id);
