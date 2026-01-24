-- QC Checkpoints - templates for QC inspections
CREATE TABLE IF NOT EXISTS qc_checkpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    checkpoint_type VARCHAR(20) NOT NULL, -- IQC, IPQC, FQC
    applies_to VARCHAR(20) NOT NULL DEFAULT 'ALL', -- ALL, MATERIAL, PRODUCT
    test_items JSONB NOT NULL, -- Array of test items with parameters
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_checkpoint_type CHECK (checkpoint_type IN ('IQC', 'IPQC', 'FQC')),
    CONSTRAINT chk_applies_to CHECK (applies_to IN ('ALL', 'MATERIAL', 'PRODUCT'))
);

CREATE INDEX idx_qc_checkpoints_type ON qc_checkpoints(checkpoint_type);
CREATE INDEX idx_qc_checkpoints_active ON qc_checkpoints(is_active);
