-- KOL Posts (track content published by KOLs)
CREATE TABLE kol_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kol_id UUID REFERENCES kols(id) NOT NULL,
    collaboration_id UUID REFERENCES kol_collaborations(id),
    campaign_id UUID REFERENCES campaigns(id),
    
    platform VARCHAR(50) NOT NULL, -- INSTAGRAM, YOUTUBE, TIKTOK, FACEBOOK
    post_type VARCHAR(50), -- POST, STORY, REEL, VIDEO, LIVE
    
    post_url TEXT,
    post_date DATE NOT NULL,
    
    -- Content
    caption TEXT,
    products_mentioned JSONB, -- Array of product UUIDs
    
    -- Metrics
    likes INT DEFAULT 0,
    comments INT DEFAULT 0,
    shares INT DEFAULT 0,
    views INT DEFAULT 0,
    reach INT DEFAULT 0,
    engagement_rate DECIMAL(5,2) DEFAULT 0,
    
    -- Analysis
    sentiment VARCHAR(50), -- POSITIVE, NEUTRAL, NEGATIVE
    summary TEXT, -- Brief summary of the review
    
    -- Screenshots / Evidence
    screenshot_urls JSONB, -- Array of URLs
    
    verified BOOLEAN DEFAULT false,
    verified_by UUID,
    verified_at TIMESTAMP,
    
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kol_posts_kol ON kol_posts(kol_id);
CREATE INDEX idx_kol_posts_collaboration ON kol_posts(collaboration_id);
CREATE INDEX idx_kol_posts_campaign ON kol_posts(campaign_id);
CREATE INDEX idx_kol_posts_platform ON kol_posts(platform);
CREATE INDEX idx_kol_posts_date ON kol_posts(post_date);
