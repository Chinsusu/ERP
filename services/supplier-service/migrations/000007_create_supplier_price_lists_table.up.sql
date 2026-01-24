-- Supplier price lists (historical pricing)
CREATE TABLE IF NOT EXISTS supplier_price_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    material_id UUID NOT NULL,  -- References master_data_db.materials
    unit_price DECIMAL(15,4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    min_order_qty DECIMAL(15,4) NOT NULL DEFAULT 1,
    price_break_qty DECIMAL(15,4),  -- quantity for price break
    price_break_price DECIMAL(15,4),  -- price at quantity break
    effective_from DATE NOT NULL,
    effective_to DATE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_supplier_price_lists_supplier_id ON supplier_price_lists(supplier_id);
CREATE INDEX idx_supplier_price_lists_material_id ON supplier_price_lists(material_id);
CREATE INDEX idx_supplier_price_lists_effective ON supplier_price_lists(effective_from, effective_to);
