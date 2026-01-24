-- Approved Supplier List (ASL) - which suppliers can supply which materials
CREATE TABLE IF NOT EXISTS approved_supplier_list (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    material_id UUID NOT NULL,  -- References master_data_db.materials
    priority INT NOT NULL DEFAULT 1 CHECK (priority >= 1 AND priority <= 3),  -- 1=primary, 2=secondary, 3=backup
    lead_time_days INT NOT NULL DEFAULT 14,
    min_order_qty DECIMAL(15,4) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,4) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    valid_from DATE NOT NULL,
    valid_to DATE,
    approved_by UUID NOT NULL,
    approved_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'EXPIRED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(supplier_id, material_id, valid_from)
);

-- Indexes
CREATE INDEX idx_asl_supplier_id ON approved_supplier_list(supplier_id);
CREATE INDEX idx_asl_material_id ON approved_supplier_list(material_id);
CREATE INDEX idx_asl_status ON approved_supplier_list(status);
CREATE INDEX idx_asl_priority ON approved_supplier_list(priority);
