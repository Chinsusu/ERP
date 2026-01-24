-- Sample Items (products in a sample request)
CREATE TABLE sample_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sample_request_id UUID REFERENCES sample_requests(id) ON DELETE CASCADE NOT NULL,
    line_number INT NOT NULL,
    
    product_id UUID NOT NULL,
    product_code VARCHAR(50),
    product_name VARCHAR(255),
    
    quantity INT NOT NULL DEFAULT 1,
    unit_value DECIMAL(18,2) DEFAULT 0, -- Value per unit for tracking cost
    total_value DECIMAL(18,2) DEFAULT 0,
    
    lot_id UUID, -- If tracked from WMS
    lot_number VARCHAR(100),
    
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sample_items_request ON sample_items(sample_request_id);
CREATE INDEX idx_sample_items_product ON sample_items(product_id);
