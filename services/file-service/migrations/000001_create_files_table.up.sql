-- Files table
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- File information
    original_name VARCHAR(500) NOT NULL,
    stored_name VARCHAR(500) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    
    -- Storage location
    bucket_name VARCHAR(100) NOT NULL,
    object_path TEXT NOT NULL,
    
    -- Categorization
    category VARCHAR(50) NOT NULL,       -- DOCUMENT, IMAGE, CERTIFICATE, CONTRACT, REPORT
    entity_type VARCHAR(50),             -- SUPPLIER, MATERIAL, PO, PR, WO, etc.
    entity_id UUID,
    
    -- Access control
    is_public BOOLEAN DEFAULT false,
    access_token VARCHAR(255),           -- For secure download links
    
    -- Metadata
    metadata JSONB DEFAULT '{}',
    checksum VARCHAR(64),                -- SHA-256 hash
    
    -- Audit
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP                 -- Optional expiry for temporary files
);

CREATE INDEX idx_files_entity ON files(entity_type, entity_id);
CREATE INDEX idx_files_category ON files(category);
CREATE INDEX idx_files_created_by ON files(created_by);
CREATE INDEX idx_files_bucket ON files(bucket_name);
CREATE INDEX idx_files_access_token ON files(access_token);

COMMENT ON TABLE files IS 'File metadata and storage references';
COMMENT ON COLUMN files.stored_name IS 'UUID-based name in object storage';
COMMENT ON COLUMN files.access_token IS 'Token for generating secure download URLs';
