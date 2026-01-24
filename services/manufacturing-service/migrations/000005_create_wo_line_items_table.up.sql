-- WO Line Items - planned materials for work order
CREATE TABLE IF NOT EXISTS wo_line_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_order_id UUID NOT NULL REFERENCES work_orders(id) ON DELETE CASCADE,
    bom_line_item_id UUID REFERENCES bom_line_items(id),
    line_number INT NOT NULL,
    material_id UUID NOT NULL,
    planned_quantity DECIMAL(15,4) NOT NULL,
    issued_quantity DECIMAL(15,4) DEFAULT 0,
    uom_id UUID NOT NULL,
    is_critical BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wo_line_items_work_order_id ON wo_line_items(work_order_id);
CREATE INDEX idx_wo_line_items_material_id ON wo_line_items(material_id);
