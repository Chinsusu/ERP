-- Work Orders table
CREATE TABLE IF NOT EXISTS work_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wo_number VARCHAR(30) NOT NULL UNIQUE, -- WO-YYYY-XXXX
    wo_date DATE NOT NULL,
    product_id UUID NOT NULL,
    bom_id UUID NOT NULL REFERENCES boms(id),
    status VARCHAR(30) NOT NULL DEFAULT 'PLANNED', -- PLANNED, RELEASED, IN_PROGRESS, QC_PENDING, COMPLETED, CANCELLED
    priority VARCHAR(20) DEFAULT 'NORMAL', -- LOW, NORMAL, HIGH, URGENT
    
    -- Planning
    planned_quantity DECIMAL(15,4) NOT NULL,
    uom_id UUID NOT NULL,
    planned_start_date TIMESTAMP,
    planned_end_date TIMESTAMP,
    
    -- Execution
    actual_start_date TIMESTAMP,
    actual_end_date TIMESTAMP,
    actual_quantity DECIMAL(15,4),
    good_quantity DECIMAL(15,4),
    rejected_quantity DECIMAL(15,4),
    yield_percentage DECIMAL(5,2),
    
    -- Batch/Lot info
    batch_number VARCHAR(50),
    output_lot_id UUID, -- Reference to WMS lot after completion
    
    -- Reference
    sales_order_id UUID,
    production_line VARCHAR(50),
    shift VARCHAR(20),
    supervisor_id UUID,
    
    -- Notes
    notes TEXT,
    
    -- Audit
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_wo_status CHECK (status IN ('PLANNED', 'RELEASED', 'IN_PROGRESS', 'QC_PENDING', 'COMPLETED', 'CANCELLED')),
    CONSTRAINT chk_wo_priority CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT'))
);

CREATE INDEX idx_work_orders_status ON work_orders(status);
CREATE INDEX idx_work_orders_product_id ON work_orders(product_id);
CREATE INDEX idx_work_orders_bom_id ON work_orders(bom_id);
CREATE INDEX idx_work_orders_batch_number ON work_orders(batch_number);
CREATE INDEX idx_work_orders_dates ON work_orders(planned_start_date, planned_end_date);
