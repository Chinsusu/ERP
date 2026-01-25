-- Dashboards table
CREATE TABLE IF NOT EXISTS dashboards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Layout configuration
    layout_type VARCHAR(50) DEFAULT 'GRID', -- GRID, FREEFORM
    layout JSONB DEFAULT '{}',              -- Grid configuration
    
    -- Status
    is_default BOOLEAN DEFAULT false,       -- Default dashboard for role
    is_system BOOLEAN DEFAULT false,        -- Cannot be deleted
    is_active BOOLEAN DEFAULT true,
    
    -- Access
    visibility VARCHAR(50) DEFAULT 'PRIVATE', -- PRIVATE, SHARED, PUBLIC
    shared_with JSONB DEFAULT '[]',           -- User/role IDs
    
    -- Audit
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_dashboards_code ON dashboards(code);
CREATE INDEX idx_dashboards_user ON dashboards(created_by);
CREATE INDEX idx_dashboards_default ON dashboards(is_default) WHERE is_default = true;

COMMENT ON TABLE dashboards IS 'User dashboards with widget configurations';
