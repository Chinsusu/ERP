-- Report executions table
CREATE TABLE IF NOT EXISTS report_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES report_definitions(id),
    
    -- Execution parameters
    parameters JSONB DEFAULT '{}',    -- Parameter values used
    
    -- Status tracking
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, RUNNING, COMPLETED, FAILED, CANCELLED
    progress INT DEFAULT 0,               -- 0-100 percentage
    
    -- Timing
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    -- Results
    row_count INT DEFAULT 0,
    result_preview JSONB,             -- First N rows for preview
    
    -- Export file
    file_path TEXT,                   -- Path to exported file
    file_format VARCHAR(20),          -- CSV, XLSX, PDF
    file_size BIGINT,
    
    -- Error handling
    error_message TEXT,
    
    -- Audit
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_report_executions_report ON report_executions(report_id);
CREATE INDEX idx_report_executions_status ON report_executions(status);
CREATE INDEX idx_report_executions_user ON report_executions(created_by);
CREATE INDEX idx_report_executions_created ON report_executions(created_at DESC);

COMMENT ON TABLE report_executions IS 'Report execution history and results';
