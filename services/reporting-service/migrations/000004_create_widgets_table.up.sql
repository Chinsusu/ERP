-- Dashboard widgets table
CREATE TABLE IF NOT EXISTS widgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dashboard_id UUID NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    
    -- Widget type
    widget_type VARCHAR(50) NOT NULL, -- KPI, LINE_CHART, BAR_CHART, PIE_CHART, TABLE, GAUGE, MAP
    
    -- Display
    title VARCHAR(200) NOT NULL,
    subtitle VARCHAR(300),
    icon VARCHAR(50),
    
    -- Data configuration
    data_source VARCHAR(100),         -- API endpoint or report code
    query_params JSONB DEFAULT '{}',  -- Query parameters
    refresh_interval INT DEFAULT 300, -- Seconds (0 = no auto refresh)
    
    -- Chart configuration
    config JSONB DEFAULT '{}',        -- Chart-specific config (colors, labels, etc.)
    
    -- Position in grid
    position_x INT DEFAULT 0,
    position_y INT DEFAULT 0,
    width INT DEFAULT 4,              -- Grid columns (1-12)
    height INT DEFAULT 2,             -- Grid rows
    
    -- Status
    is_visible BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_widgets_dashboard ON widgets(dashboard_id);
CREATE INDEX idx_widgets_type ON widgets(widget_type);

COMMENT ON TABLE widgets IS 'Dashboard widget configurations';
COMMENT ON COLUMN widgets.config IS 'Chart-specific configuration (colors, legends, etc.)';
