-- BOM Line Items - component materials
CREATE TABLE IF NOT EXISTS bom_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_id UUID NOT NULL REFERENCES boms(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    material_id UUID NOT NULL,
    item_type VARCHAR(30) NOT NULL DEFAULT 'MATERIAL', -- MATERIAL, PACKAGING, CONSUMABLE
    quantity DECIMAL(15,4) NOT NULL,
    uom_id UUID NOT NULL,
    quantity_min DECIMAL(15,4),
    quantity_max DECIMAL(15,4),
    is_critical BOOLEAN DEFAULT FALSE,
    scrap_percentage DECIMAL(5,2) DEFAULT 0,
    unit_cost DECIMAL(18,4) DEFAULT 0,
    total_cost DECIMAL(18,2) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_item_type CHECK (item_type IN ('MATERIAL', 'PACKAGING', 'CONSUMABLE'))
);

CREATE INDEX idx_bom_line_items_bom_id ON bom_line_items(bom_id);
CREATE INDEX idx_bom_line_items_material_id ON bom_line_items(material_id);
CREATE UNIQUE INDEX idx_bom_line_items_bom_line ON bom_line_items(bom_id, line_number);
