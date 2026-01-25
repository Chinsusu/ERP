-- Report definitions table
CREATE TABLE IF NOT EXISTS report_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Report categorization
    report_type VARCHAR(50) NOT NULL, -- INVENTORY, SALES, PROCUREMENT, PRODUCTION, FINANCIAL
    category VARCHAR(50),             -- SUB-category
    
    -- Query configuration
    query_template TEXT NOT NULL,     -- SQL template with placeholders
    data_source VARCHAR(100),         -- Service to query (wms, sales, etc.)
    
    -- Parameter definitions
    parameters JSONB DEFAULT '[]',    -- [{name, type, required, default, options}]
    
    -- Column definitions
    columns JSONB DEFAULT '[]',       -- [{field, header, type, format, width}]
    
    -- Sorting and grouping
    default_sort JSONB,               -- {field, direction}
    group_by JSONB,                   -- Grouping fields
    
    -- Access control
    required_permission VARCHAR(100),
    
    -- Status
    is_system BOOLEAN DEFAULT false,  -- System reports cannot be deleted
    is_active BOOLEAN DEFAULT true,
    
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_report_definitions_type ON report_definitions(report_type);
CREATE INDEX idx_report_definitions_active ON report_definitions(is_active);

COMMENT ON TABLE report_definitions IS 'Report template definitions';
COMMENT ON COLUMN report_definitions.query_template IS 'SQL template with {{.ParamName}} placeholders';
COMMENT ON COLUMN report_definitions.parameters IS 'JSON array of parameter definitions';
