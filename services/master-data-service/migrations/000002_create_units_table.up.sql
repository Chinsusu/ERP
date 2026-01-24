-- Migration: Create units_of_measure table

CREATE TABLE IF NOT EXISTS units_of_measure (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    symbol VARCHAR(10) NOT NULL,
    uom_type VARCHAR(50) NOT NULL, -- WEIGHT, VOLUME, QUANTITY, LENGTH, AREA
    is_base_unit BOOLEAN NOT NULL DEFAULT false,
    base_unit_id UUID REFERENCES units_of_measure(id) ON DELETE SET NULL,
    conversion_factor DECIMAL(18, 8) NOT NULL DEFAULT 1,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_uom_type ON units_of_measure(uom_type);
CREATE INDEX idx_uom_base_unit ON units_of_measure(base_unit_id);
CREATE INDEX idx_uom_status ON units_of_measure(status);
CREATE INDEX idx_uom_deleted_at ON units_of_measure(deleted_at);

-- Comments
COMMENT ON TABLE units_of_measure IS 'Units of measure with conversion support';
COMMENT ON COLUMN units_of_measure.conversion_factor IS 'Factor to convert to base unit (e.g., G to KG = 0.001)';
COMMENT ON COLUMN units_of_measure.is_base_unit IS 'True if this is a base unit (KG, L, PCS)';
