-- Sample Requests for KOL marketing
CREATE TABLE sample_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_number VARCHAR(50) UNIQUE NOT NULL,
    kol_id UUID REFERENCES kols(id) NOT NULL,
    campaign_id UUID REFERENCES campaigns(id),
    collaboration_id UUID REFERENCES kol_collaborations(id),
    
    request_date DATE NOT NULL DEFAULT CURRENT_DATE,
    request_reason TEXT,
    
    -- Delivery Address (can override KOL address)
    delivery_address TEXT,
    recipient_name VARCHAR(200),
    recipient_phone VARCHAR(50),
    
    -- Expectations
    expected_post_date DATE,
    expected_reach INT DEFAULT 0,
    
    -- Value
    total_items INT DEFAULT 0,
    total_value DECIMAL(18,2) DEFAULT 0,
    
    -- Status Workflow: DRAFT → PENDING_APPROVAL → APPROVED → SHIPPED → DELIVERED → FEEDBACK_RECEIVED | REJECTED
    status VARCHAR(50) DEFAULT 'DRAFT',
    
    -- Approval
    approved_by UUID,
    approved_at TIMESTAMP,
    rejection_reason TEXT,
    
    notes TEXT,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sample_requests_number ON sample_requests(request_number);
CREATE INDEX idx_sample_requests_kol ON sample_requests(kol_id);
CREATE INDEX idx_sample_requests_campaign ON sample_requests(campaign_id);
CREATE INDEX idx_sample_requests_status ON sample_requests(status);
CREATE INDEX idx_sample_requests_date ON sample_requests(request_date);
