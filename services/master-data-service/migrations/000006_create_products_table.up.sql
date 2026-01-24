-- Migration: Create products table
-- Finished goods with cosmetics-specific fields

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    sku VARCHAR(100) NOT NULL UNIQUE,
    barcode VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    description TEXT,
    
    -- Classification
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    product_line VARCHAR(100), -- Skincare, Haircare, Bodycare
    brand VARCHAR(100),
    
    -- Physical properties
    volume DECIMAL(10, 2),
    volume_unit VARCHAR(20), -- ml, g, oz
    weight DECIMAL(10, 2),
    weight_unit VARCHAR(20), -- g, kg
    
    -- Cosmetics regulatory
    cosmetic_license_number VARCHAR(100),
    license_expiry_date DATE,
    registration_country VARCHAR(100),
    
    -- Product info
    ingredients_summary TEXT, -- For product label
    target_skin_type VARCHAR(255),
    usage_instructions TEXT,
    warnings TEXT,
    
    -- Packaging
    packaging_type VARCHAR(100), -- Bottle, Tube, Jar, Box
    packaging_material VARCHAR(100), -- Glass, Plastic, Aluminum
    
    -- Pricing
    standard_cost DECIMAL(18, 4) NOT NULL DEFAULT 0,
    standard_price DECIMAL(18, 4) NOT NULL DEFAULT 0,
    recommended_retail_price DECIMAL(18, 4),
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    
    -- Units
    base_unit_id UUID NOT NULL REFERENCES units_of_measure(id),
    sales_unit_id UUID REFERENCES units_of_measure(id),
    
    -- Shelf life
    shelf_life_months INT NOT NULL DEFAULT 24,
    
    -- Launch & Status
    launch_date DATE,
    discontinue_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, discontinued, pending_launch
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Indexes
CREATE INDEX idx_products_code ON products(code);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_product_line ON products(product_line);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_license_expiry ON products(license_expiry_date);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);

-- Full text search index
CREATE INDEX idx_products_search ON products USING gin(
    to_tsvector('english', coalesce(name, '') || ' ' || coalesce(name_en, '') || ' ' || coalesce(sku, '') || ' ' || coalesce(code, ''))
);

-- Comments
COMMENT ON TABLE products IS 'Finished goods (cosmetic products)';
COMMENT ON COLUMN products.cosmetic_license_number IS 'Registration number from health authority';
COMMENT ON COLUMN products.ingredients_summary IS 'INCI ingredients list for product label';
