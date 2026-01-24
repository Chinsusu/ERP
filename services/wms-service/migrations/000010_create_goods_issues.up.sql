-- Goods Issues table
CREATE TABLE IF NOT EXISTS goods_issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_number VARCHAR(30) UNIQUE NOT NULL, -- GI-YYYY-XXXX
    issue_date DATE NOT NULL,
    issue_type VARCHAR(30) NOT NULL, -- PRODUCTION, SALES, SAMPLE, SCRAP, ADJUSTMENT
    reference_type VARCHAR(30), -- WORK_ORDER, SALES_ORDER, SAMPLE_REQUEST
    reference_id UUID,
    reference_number VARCHAR(30),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, CONFIRMED, COMPLETED, CANCELLED
    notes TEXT,
    issued_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_goods_issues_type ON goods_issues(issue_type);
CREATE INDEX idx_goods_issues_status ON goods_issues(status);
CREATE INDEX idx_goods_issues_warehouse ON goods_issues(warehouse_id);
CREATE INDEX idx_goods_issues_reference ON goods_issues(reference_type, reference_id);
CREATE INDEX idx_goods_issues_date ON goods_issues(issue_date);
