-- NCRs (Non-Conformance Reports)
CREATE TABLE IF NOT EXISTS ncrs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ncr_number VARCHAR(30) NOT NULL UNIQUE, -- NCR-YYYY-XXXX
    ncr_date DATE NOT NULL,
    nc_type VARCHAR(30) NOT NULL, -- MATERIAL, PROCESS, PRODUCT, EQUIPMENT
    severity VARCHAR(20) NOT NULL DEFAULT 'MEDIUM', -- LOW, MEDIUM, HIGH, CRITICAL
    status VARCHAR(30) NOT NULL DEFAULT 'OPEN', -- OPEN, INVESTIGATION, CORRECTIVE_ACTION, CLOSED
    
    -- Reference
    reference_type VARCHAR(30), -- WORK_ORDER, QC_INSPECTION, GRN, LOT
    reference_id UUID,
    
    product_id UUID,
    material_id UUID,
    lot_id UUID,
    lot_number VARCHAR(50),
    
    -- Issue details
    description TEXT NOT NULL,
    quantity_affected DECIMAL(15,4),
    uom_id UUID,
    
    -- Investigation
    root_cause TEXT,
    investigation_date TIMESTAMP,
    investigated_by UUID,
    
    -- Actions
    immediate_action TEXT,
    corrective_action TEXT,
    preventive_action TEXT,
    
    -- Disposition
    disposition VARCHAR(30), -- USE_AS_IS, REWORK, SCRAP, RETURN_TO_SUPPLIER
    disposition_quantity DECIMAL(15,4),
    disposition_date TIMESTAMP,
    disposition_by UUID,
    
    -- Closure
    closed_at TIMESTAMP,
    closed_by UUID,
    closure_notes TEXT,
    
    -- Audit
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_nc_type CHECK (nc_type IN ('MATERIAL', 'PROCESS', 'PRODUCT', 'EQUIPMENT')),
    CONSTRAINT chk_severity CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    CONSTRAINT chk_ncr_status CHECK (status IN ('OPEN', 'INVESTIGATION', 'CORRECTIVE_ACTION', 'CLOSED')),
    CONSTRAINT chk_disposition CHECK (disposition IS NULL OR disposition IN ('USE_AS_IS', 'REWORK', 'SCRAP', 'RETURN_TO_SUPPLIER'))
);

CREATE INDEX idx_ncrs_status ON ncrs(status);
CREATE INDEX idx_ncrs_severity ON ncrs(severity);
CREATE INDEX idx_ncrs_nc_type ON ncrs(nc_type);
CREATE INDEX idx_ncrs_lot_id ON ncrs(lot_id);
CREATE INDEX idx_ncrs_reference ON ncrs(reference_type, reference_id);
