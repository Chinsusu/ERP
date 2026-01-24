-- KOL Tiers for classification
CREATE TABLE kol_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    min_followers INT DEFAULT 0,
    max_followers INT,
    auto_approve_samples BOOLEAN DEFAULT false,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kol_tiers_code ON kol_tiers(code);
CREATE INDEX idx_kol_tiers_active ON kol_tiers(is_active);
