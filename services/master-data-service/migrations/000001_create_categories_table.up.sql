-- Migration: Create categories table
-- Hierarchical categories using materialized path pattern

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    name_en VARCHAR(255),
    description TEXT,
    category_type VARCHAR(50) NOT NULL DEFAULT 'MATERIAL', -- MATERIAL, PRODUCT
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    path VARCHAR(1000) NOT NULL DEFAULT '/', -- Materialized path: /ROOT/PARENT/CHILD/
    level INT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    updated_by UUID
);

-- Indexes
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_path ON categories(path);
CREATE INDEX idx_categories_type ON categories(category_type);
CREATE INDEX idx_categories_status ON categories(status);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Comments
COMMENT ON TABLE categories IS 'Hierarchical categories for materials and products';
COMMENT ON COLUMN categories.path IS 'Materialized path for efficient tree queries';
COMMENT ON COLUMN categories.category_type IS 'MATERIAL for raw materials/packaging, PRODUCT for finished goods';
