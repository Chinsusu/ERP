-- File categories table
CREATE TABLE IF NOT EXISTS file_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Validation rules
    allowed_extensions JSONB DEFAULT '[]',   -- ["pdf", "jpg", "png"]
    max_file_size BIGINT NOT NULL,           -- In bytes
    
    -- Storage configuration
    storage_bucket VARCHAR(100) NOT NULL,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_file_categories_code ON file_categories(code);

COMMENT ON TABLE file_categories IS 'File category definitions with validation rules';
