-- Key Opinion Leaders (KOLs/Influencers)
CREATE TABLE kols (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kol_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    tier_id UUID REFERENCES kol_tiers(id),
    category VARCHAR(100), -- BEAUTY_BLOGGER, INFLUENCER, CELEBRITY, EXPERT
    
    -- Social Media Platforms
    instagram_handle VARCHAR(100),
    instagram_followers INT DEFAULT 0,
    youtube_channel VARCHAR(200),
    youtube_subscribers INT DEFAULT 0,
    tiktok_handle VARCHAR(100),
    tiktok_followers INT DEFAULT 0,
    facebook_page VARCHAR(200),
    facebook_followers INT DEFAULT 0,
    
    -- Engagement Metrics
    avg_engagement_rate DECIMAL(5,2) DEFAULT 0, -- percentage
    niche VARCHAR(100), -- Skincare, Makeup, Haircare, etc.
    
    -- Business Terms
    collaboration_rate DECIMAL(18,2) DEFAULT 0, -- Fee per post
    currency VARCHAR(10) DEFAULT 'VND',
    preferred_products JSONB, -- Array of product categories they prefer
    
    -- Address
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Vietnam',
    
    -- Performance Tracking
    total_posts INT DEFAULT 0,
    total_samples_received INT DEFAULT 0,
    total_collaborations INT DEFAULT 0,
    last_collaboration_date DATE,
    
    notes TEXT,
    status VARCHAR(50) DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE, BLACKLISTED
    
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kols_code ON kols(kol_code);
CREATE INDEX idx_kols_tier ON kols(tier_id);
CREATE INDEX idx_kols_category ON kols(category);
CREATE INDEX idx_kols_niche ON kols(niche);
CREATE INDEX idx_kols_status ON kols(status);
CREATE INDEX idx_kols_instagram ON kols(instagram_handle);
