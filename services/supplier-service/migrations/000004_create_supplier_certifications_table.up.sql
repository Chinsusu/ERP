-- Supplier certifications (GMP, ISO, Organic, etc.)
CREATE TABLE IF NOT EXISTS supplier_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    certification_type VARCHAR(20) NOT NULL CHECK (certification_type IN ('GMP', 'ISO9001', 'ISO22716', 'ORGANIC', 'ECOCERT', 'HALAL', 'COSMOS', 'OTHER')),
    certificate_number VARCHAR(100) NOT NULL,
    issuing_body VARCHAR(255) NOT NULL,
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    document_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'VALID' CHECK (status IN ('VALID', 'EXPIRING_SOON', 'EXPIRED')),
    verified_by UUID,
    verified_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_supplier_certifications_supplier_id ON supplier_certifications(supplier_id);
CREATE INDEX idx_supplier_certifications_type ON supplier_certifications(certification_type);
CREATE INDEX idx_supplier_certifications_status ON supplier_certifications(status);
CREATE INDEX idx_supplier_certifications_expiry_date ON supplier_certifications(expiry_date);
