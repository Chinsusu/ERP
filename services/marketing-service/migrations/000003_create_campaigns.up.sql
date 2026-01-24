-- Marketing Campaigns
CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    campaign_type VARCHAR(50) NOT NULL, -- PRODUCT_LAUNCH, SEASONAL, PROMOTION, AWARENESS, INFLUENCER
    
    -- Timeline
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    
    -- Target
    target_audience TEXT,
    channels JSONB, -- Array of channels: INSTAGRAM, FACEBOOK, YOUTUBE, TIKTOK
    
    -- Budget
    budget DECIMAL(18,2) DEFAULT 0,
    spent DECIMAL(18,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'VND',
    
    -- Goals (JSON for flexibility)
    goals JSONB, -- {awareness: "100k impressions", engagement: "5k interactions", sales: "500 units"}
    
    -- Products linked to campaign
    products JSONB, -- Array of product UUIDs
    
    -- Performance
    impressions INT DEFAULT 0,
    reach INT DEFAULT 0,
    engagement INT DEFAULT 0,
    conversions INT DEFAULT 0,
    revenue_generated DECIMAL(18,2) DEFAULT 0,
    
    -- Status
    status VARCHAR(50) DEFAULT 'DRAFT', -- DRAFT, PLANNED, ACTIVE, PAUSED, COMPLETED, CANCELLED
    
    notes TEXT,
    created_by UUID,
    updated_by UUID,
    launched_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_campaigns_code ON campaigns(campaign_code);
CREATE INDEX idx_campaigns_type ON campaigns(campaign_type);
CREATE INDEX idx_campaigns_status ON campaigns(status);
CREATE INDEX idx_campaigns_dates ON campaigns(start_date, end_date);
