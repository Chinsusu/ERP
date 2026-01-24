-- Returns table
CREATE TABLE IF NOT EXISTS returns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_number VARCHAR(20) NOT NULL UNIQUE,
    sales_order_id UUID NOT NULL REFERENCES sales_orders(id),
    shipment_id UUID REFERENCES shipments(id),
    return_date DATE NOT NULL DEFAULT CURRENT_DATE,
    return_reason VARCHAR(100),
    return_type VARCHAR(20) DEFAULT 'REFUND' CHECK (return_type IN ('REFUND', 'EXCHANGE', 'CREDIT')),
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'RECEIVED', 'INSPECTED', 'COMPLETED', 'REJECTED')),
    subtotal DECIMAL(18,2) DEFAULT 0,
    refund_amount DECIMAL(18,2) DEFAULT 0,
    notes TEXT,
    approved_by UUID,
    approved_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Return line items
CREATE TABLE IF NOT EXISTS return_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_id UUID NOT NULL REFERENCES returns(id) ON DELETE CASCADE,
    so_line_item_id UUID REFERENCES so_line_items(id),
    product_id UUID NOT NULL,
    product_code VARCHAR(50),
    product_name VARCHAR(200),
    quantity DECIMAL(18,3) NOT NULL,
    unit_price DECIMAL(18,2) NOT NULL,
    reason VARCHAR(200),
    condition VARCHAR(20) DEFAULT 'GOOD' CHECK (condition IN ('GOOD', 'DAMAGED', 'DEFECTIVE', 'EXPIRED')),
    action VARCHAR(20) DEFAULT 'REFUND' CHECK (action IN ('REFUND', 'EXCHANGE', 'CREDIT', 'DISPOSE')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_returns_number ON returns(return_number);
CREATE INDEX idx_returns_order ON returns(sales_order_id);
CREATE INDEX idx_returns_status ON returns(status);
CREATE INDEX idx_returns_date ON returns(return_date);
CREATE INDEX idx_return_items_return ON return_line_items(return_id);
CREATE INDEX idx_return_items_product ON return_line_items(product_id);
