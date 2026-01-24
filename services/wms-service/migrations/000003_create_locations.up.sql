-- Locations table
CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id UUID NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    code VARCHAR(30) NOT NULL, -- A01-R02-S03-B01 (Aisle-Rack-Shelf-Bin)
    aisle VARCHAR(10),
    rack VARCHAR(10),
    shelf VARCHAR(10),
    bin VARCHAR(10),
    capacity DECIMAL(10,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(zone_id, code)
);

-- Create indexes
CREATE INDEX idx_locations_zone ON locations(zone_id);
CREATE INDEX idx_locations_code ON locations(code);
