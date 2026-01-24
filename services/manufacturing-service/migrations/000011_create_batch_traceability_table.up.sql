-- Batch Traceability - links material lots to product lots
CREATE TABLE IF NOT EXISTS batch_traceability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_order_id UUID NOT NULL REFERENCES work_orders(id),
    wo_material_issue_id UUID REFERENCES wo_material_issues(id),
    
    -- Material side (input)
    material_id UUID NOT NULL,
    material_lot_id UUID NOT NULL,
    material_lot_number VARCHAR(50) NOT NULL,
    material_quantity DECIMAL(15,4) NOT NULL,
    material_uom_id UUID NOT NULL,
    supplier_lot_number VARCHAR(100),
    
    -- Product side (output)
    product_id UUID NOT NULL,
    product_lot_id UUID, -- May be NULL until WO is completed
    product_lot_number VARCHAR(50),
    
    -- Linkage
    trace_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for forward and backward tracing
CREATE INDEX idx_batch_trace_material_lot ON batch_traceability(material_lot_id);
CREATE INDEX idx_batch_trace_product_lot ON batch_traceability(product_lot_id);
CREATE INDEX idx_batch_trace_work_order ON batch_traceability(work_order_id);
CREATE INDEX idx_batch_trace_material_id ON batch_traceability(material_id);
CREATE INDEX idx_batch_trace_product_id ON batch_traceability(product_id);
