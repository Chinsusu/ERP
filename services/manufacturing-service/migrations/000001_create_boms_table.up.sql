-- BOMs table with encrypted formula_details
CREATE TABLE IF NOT EXISTS boms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_number VARCHAR(50) NOT NULL UNIQUE,
    product_id UUID NOT NULL,
    version INT NOT NULL DEFAULT 1,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT', -- DRAFT, PENDING_APPROVAL, APPROVED, OBSOLETE
    batch_size DECIMAL(15,4) NOT NULL,
    batch_unit_id UUID NOT NULL,
    
    -- ENCRYPTED field - stores AES-256-GCM encrypted JSON
    formula_details BYTEA,
    confidentiality_level VARCHAR(30) NOT NULL DEFAULT 'RESTRICTED', -- PUBLIC, INTERNAL, CONFIDENTIAL, RESTRICTED
    
    -- Costing
    material_cost DECIMAL(18,2) DEFAULT 0,
    labor_cost DECIMAL(18,2) DEFAULT 0,
    overhead_cost DECIMAL(18,2) DEFAULT 0,
    total_cost DECIMAL(18,2) DEFAULT 0,
    
    -- Validity
    effective_from DATE,
    effective_to DATE,
    
    -- Approval
    approved_by UUID,
    approved_at TIMESTAMP,
    
    -- Audit
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_bom_status CHECK (status IN ('DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'OBSOLETE')),
    CONSTRAINT chk_confidentiality CHECK (confidentiality_level IN ('PUBLIC', 'INTERNAL', 'CONFIDENTIAL', 'RESTRICTED'))
);

-- Unique constraint: product can have only one active version
CREATE UNIQUE INDEX idx_boms_product_version ON boms(product_id, version);
CREATE INDEX idx_boms_product_id ON boms(product_id);
CREATE INDEX idx_boms_status ON boms(status);
