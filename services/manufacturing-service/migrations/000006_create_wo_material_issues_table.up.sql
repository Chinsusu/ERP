-- WO Material Issues - actual materials issued from WMS
CREATE TABLE IF NOT EXISTS wo_material_issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_order_id UUID NOT NULL REFERENCES work_orders(id) ON DELETE CASCADE,
    wo_line_item_id UUID REFERENCES wo_line_items(id),
    issue_number VARCHAR(30) NOT NULL, -- ISS-YYYY-XXXX
    issue_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    material_id UUID NOT NULL,
    lot_id UUID NOT NULL, -- Reference to WMS lot
    lot_number VARCHAR(50) NOT NULL,
    quantity DECIMAL(15,4) NOT NULL,
    uom_id UUID NOT NULL,
    wms_movement_id UUID, -- Reference to WMS stock movement
    issued_by UUID,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wo_material_issues_work_order_id ON wo_material_issues(work_order_id);
CREATE INDEX idx_wo_material_issues_lot_id ON wo_material_issues(lot_id);
CREATE INDEX idx_wo_material_issues_material_id ON wo_material_issues(material_id);
