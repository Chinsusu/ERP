-- Migration: Create unit_conversions table
-- Explicit bidirectional conversions between units

CREATE TABLE IF NOT EXISTS unit_conversions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_unit_id UUID NOT NULL REFERENCES units_of_measure(id) ON DELETE CASCADE,
    to_unit_id UUID NOT NULL REFERENCES units_of_measure(id) ON DELETE CASCADE,
    conversion_factor DECIMAL(18, 8) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_unit_conversion UNIQUE(from_unit_id, to_unit_id),
    CONSTRAINT ck_different_units CHECK (from_unit_id != to_unit_id)
);

-- Indexes
CREATE INDEX idx_unit_conv_from ON unit_conversions(from_unit_id);
CREATE INDEX idx_unit_conv_to ON unit_conversions(to_unit_id);

-- Comments
COMMENT ON TABLE unit_conversions IS 'Explicit conversion factors between units';
COMMENT ON COLUMN unit_conversions.conversion_factor IS 'Multiply from_unit value by this to get to_unit value';
