-- KOL Collaborations (link KOLs to Campaigns)
CREATE TABLE kol_collaborations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collaboration_code VARCHAR(50) UNIQUE NOT NULL,
    campaign_id UUID REFERENCES campaigns(id),
    kol_id UUID REFERENCES kols(id) NOT NULL,
    
    -- Terms
    collaboration_type VARCHAR(50), -- SPONSORED_POST, PRODUCT_REVIEW, GIVEAWAY, AMBASSADOR
    agreed_fee DECIMAL(18,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'VND',
    
    -- Deliverables
    expected_posts INT DEFAULT 1,
    actual_posts INT DEFAULT 0,
    platforms JSONB, -- Which platforms will be used
    
    -- Timeline
    start_date DATE,
    end_date DATE,
    content_deadline DATE,
    
    -- Performance (aggregated from posts)
    total_impressions INT DEFAULT 0,
    total_engagement INT DEFAULT 0,
    total_reach INT DEFAULT 0,
    
    -- Payment
    payment_status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, PARTIAL, PAID
    paid_amount DECIMAL(18,2) DEFAULT 0,
    
    status VARCHAR(50) DEFAULT 'DRAFT', -- DRAFT, ACTIVE, COMPLETED, CANCELLED
    notes TEXT,
    
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_collaborations_campaign ON kol_collaborations(campaign_id);
CREATE INDEX idx_collaborations_kol ON kol_collaborations(kol_id);
CREATE INDEX idx_collaborations_status ON kol_collaborations(status);
