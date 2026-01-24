-- Migration: Create materials table
-- Raw materials, packaging, and consumables with cosmetics-specific fields

CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    description TEXT,
    material_type VARCHAR(50) NOT NULL, -- RAW_MATERIAL, PACKAGING, CONSUMABLE, SEMI_FINISHED
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    base_unit_id UUID NOT NULL REFERENCES units_of_measure(id),
    purchase_unit_id UUID REFERENCES units_of_measure(id),
    stock_unit_id UUID REFERENCES units_of_measure(id),
    
    -- Cosmetics specific fields
    inci_name VARCHAR(500), -- International Nomenclature of Cosmetic Ingredients
    cas_number VARCHAR(50), -- Chemical Abstracts Service number
    is_allergen BOOLEAN NOT NULL DEFAULT false,
    allergen_info TEXT,
    is_organic BOOLEAN NOT NULL DEFAULT false,
    is_natural BOOLEAN NOT NULL DEFAULT false,
    is_vegan BOOLEAN NOT NULL DEFAULT false,
    origin_country VARCHAR(100),
    
    -- Storage conditions
    storage_condition VARCHAR(50) NOT NULL DEFAULT 'AMBIENT', -- AMBIENT, COLD, FROZEN
    min_temp DECIMAL(5, 2), -- Minimum storage temperature in Celsius
    max_temp DECIMAL(5, 2), -- Maximum storage temperature in Celsius
    storage_instructions TEXT,
    shelf_life_days INT NOT NULL DEFAULT 365,
    
    -- Safety
    is_hazardous BOOLEAN NOT NULL DEFAULT false,
    hazard_class VARCHAR(50),
    safety_data_sheet_url VARCHAR(500),
    
    -- Procurement
    default_supplier_id UUID, -- Will be FK when Supplier Service is integrated
    lead_time_days INT NOT NULL DEFAULT 14,
    min_order_qty DECIMAL(18, 4) NOT NULL DEFAULT 1,
    reorder_point DECIMAL(18, 4) NOT NULL DEFAULT 0,
    safety_stock DECIMAL(18, 4) NOT NULL DEFAULT 0,
    max_stock_qty DECIMAL(18, 4),
    
    -- Costing
    standard_cost DECIMAL(18, 4) NOT NULL DEFAULT 0,
    last_purchase_cost DECIMAL(18, 4),
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, discontinued, pending_approval
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Indexes
CREATE INDEX idx_materials_code ON materials(code);
CREATE INDEX idx_materials_type ON materials(material_type);
CREATE INDEX idx_materials_category ON materials(category_id);
CREATE INDEX idx_materials_inci ON materials(inci_name);
CREATE INDEX idx_materials_cas ON materials(cas_number);
CREATE INDEX idx_materials_storage ON materials(storage_condition);
CREATE INDEX idx_materials_status ON materials(status);
CREATE INDEX idx_materials_deleted_at ON materials(deleted_at);
CREATE INDEX idx_materials_supplier ON materials(default_supplier_id);

-- Full text search index
CREATE INDEX idx_materials_search ON materials USING gin(
    to_tsvector('english', coalesce(name, '') || ' ' || coalesce(name_en, '') || ' ' || coalesce(inci_name, '') || ' ' || coalesce(code, ''))
);

-- Comments
COMMENT ON TABLE materials IS 'Raw materials, packaging, and consumables for cosmetics manufacturing';
COMMENT ON COLUMN materials.inci_name IS 'International Nomenclature of Cosmetic Ingredients - standard name';
COMMENT ON COLUMN materials.cas_number IS 'Chemical Abstracts Service registry number';
COMMENT ON COLUMN materials.storage_condition IS 'AMBIENT (room temp), COLD (2-8°C), FROZEN (<-18°C)';
