-- Shipments table
CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_number VARCHAR(20) NOT NULL UNIQUE,
    sales_order_id UUID NOT NULL REFERENCES sales_orders(id),
    shipped_date DATE,
    estimated_delivery_date DATE,
    actual_delivery_date DATE,
    carrier VARCHAR(100),
    tracking_number VARCHAR(100),
    shipping_method VARCHAR(50),
    shipping_cost DECIMAL(18,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PICKED', 'PACKED', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'RETURNED')),
    recipient_name VARCHAR(200),
    recipient_phone VARCHAR(20),
    delivery_address TEXT,
    notes TEXT,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_shipments_number ON shipments(shipment_number);
CREATE INDEX idx_shipments_order ON shipments(sales_order_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_date ON shipments(shipped_date);
CREATE INDEX idx_shipments_tracking ON shipments(tracking_number);
