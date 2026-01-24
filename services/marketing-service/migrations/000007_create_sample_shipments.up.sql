-- Sample Shipments (tracking delivery to KOLs)
CREATE TABLE sample_shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_number VARCHAR(50) UNIQUE NOT NULL,
    sample_request_id UUID REFERENCES sample_requests(id) NOT NULL,
    
    shipment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    courier VARCHAR(100), -- Giao Hang Nhanh, VNPost, etc.
    tracking_number VARCHAR(100),
    
    -- Recipient (copied from sample request or KOL)
    recipient_name VARCHAR(200),
    recipient_phone VARCHAR(50),
    delivery_address TEXT,
    
    -- Tracking
    estimated_delivery DATE,
    actual_delivery DATE,
    
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, SHIPPED, IN_TRANSIT, DELIVERED, RETURNED
    
    delivery_notes TEXT,
    proof_of_delivery TEXT, -- URL to image/signature
    
    created_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sample_shipments_request ON sample_shipments(sample_request_id);
CREATE INDEX idx_sample_shipments_tracking ON sample_shipments(tracking_number);
CREATE INDEX idx_sample_shipments_status ON sample_shipments(status);
