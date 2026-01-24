-- QC Inspections - actual QC records
CREATE TABLE IF NOT EXISTS qc_inspections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inspection_number VARCHAR(30) NOT NULL UNIQUE, -- QC-YYYY-XXXX
    inspection_date TIMESTAMP NOT NULL,
    inspection_type VARCHAR(20) NOT NULL, -- IQC, IPQC, FQC
    checkpoint_id UUID REFERENCES qc_checkpoints(id),
    
    -- Reference (what is being inspected)
    reference_type VARCHAR(30) NOT NULL, -- WORK_ORDER, GRN, LOT
    reference_id UUID NOT NULL,
    
    product_id UUID,
    material_id UUID,
    lot_id UUID,
    lot_number VARCHAR(50),
    
    -- Quantities
    inspected_quantity DECIMAL(15,4) NOT NULL,
    accepted_quantity DECIMAL(15,4),
    rejected_quantity DECIMAL(15,4),
    sample_size INT,
    
    -- Result
    result VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PASSED, FAILED, CONDITIONAL
    overall_score DECIMAL(5,2),
    
    -- Inspector
    inspector_id UUID NOT NULL,
    inspector_name VARCHAR(100),
    
    -- Approval
    approved_by UUID,
    approved_at TIMESTAMP,
    
    test_results JSONB, -- Detailed test results
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_inspection_type CHECK (inspection_type IN ('IQC', 'IPQC', 'FQC')),
    CONSTRAINT chk_reference_type CHECK (reference_type IN ('WORK_ORDER', 'GRN', 'LOT')),
    CONSTRAINT chk_result CHECK (result IN ('PENDING', 'PASSED', 'FAILED', 'CONDITIONAL'))
);

CREATE INDEX idx_qc_inspections_type ON qc_inspections(inspection_type);
CREATE INDEX idx_qc_inspections_result ON qc_inspections(result);
CREATE INDEX idx_qc_inspections_reference ON qc_inspections(reference_type, reference_id);
CREATE INDEX idx_qc_inspections_lot_id ON qc_inspections(lot_id);
CREATE INDEX idx_qc_inspections_date ON qc_inspections(inspection_date);
