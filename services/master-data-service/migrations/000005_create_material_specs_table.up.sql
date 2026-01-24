-- Migration: Create material_specifications table
-- Extended specifications, certificates, and custom attributes

CREATE TABLE IF NOT EXISTS material_specifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    
    -- Specifications
    spec_name VARCHAR(255) NOT NULL,
    spec_value TEXT,
    spec_unit VARCHAR(50),
    min_value DECIMAL(18, 6),
    max_value DECIMAL(18, 6),
    
    -- Certificates
    certificate_type VARCHAR(100), -- GMP, ISO, ORGANIC, HALAL, etc.
    certificate_number VARCHAR(100),
    certificate_issuer VARCHAR(255),
    certificate_expiry_date DATE,
    certificate_file_url VARCHAR(500),
    
    -- Quality
    quality_grade VARCHAR(50), -- A, B, C, PHARMACEUTICAL, COSMETIC, FOOD
    purity_percentage DECIMAL(5, 2),
    test_method VARCHAR(255),
    
    -- Custom attributes (JSONB for flexibility)
    custom_attributes JSONB,
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Indexes
CREATE INDEX idx_material_specs_material ON material_specifications(material_id);
CREATE INDEX idx_material_specs_cert_type ON material_specifications(certificate_type);
CREATE INDEX idx_material_specs_cert_expiry ON material_specifications(certificate_expiry_date);

-- Comments
COMMENT ON TABLE material_specifications IS 'Extended specifications and certificates for materials';
COMMENT ON COLUMN material_specifications.custom_attributes IS 'Flexible JSONB field for industry-specific attributes';
