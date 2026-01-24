-- QC Inspection Items - individual test results
CREATE TABLE IF NOT EXISTS qc_inspection_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inspection_id UUID NOT NULL REFERENCES qc_inspections(id) ON DELETE CASCADE,
    item_number INT NOT NULL,
    test_name VARCHAR(100) NOT NULL,
    test_method VARCHAR(100),
    specification VARCHAR(200),
    target_value VARCHAR(100),
    min_value VARCHAR(100),
    max_value VARCHAR(100),
    actual_value VARCHAR(100),
    uom VARCHAR(20),
    result VARCHAR(20) NOT NULL, -- PASS, FAIL, N/A
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_item_result CHECK (result IN ('PASS', 'FAIL', 'N/A'))
);

CREATE INDEX idx_qc_inspection_items_inspection_id ON qc_inspection_items(inspection_id);
