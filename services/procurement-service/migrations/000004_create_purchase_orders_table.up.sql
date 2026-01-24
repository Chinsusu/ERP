-- Purchase Orders (PO) table
CREATE TABLE IF NOT EXISTS purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    po_number VARCHAR(20) NOT NULL UNIQUE,
    po_date DATE NOT NULL DEFAULT CURRENT_DATE,
    pr_id UUID REFERENCES purchase_requisitions(id),
    supplier_id UUID NOT NULL,
    supplier_code VARCHAR(50),
    supplier_name VARCHAR(255),
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SUBMITTED', 'CONFIRMED', 'PARTIALLY_RECEIVED', 'FULLY_RECEIVED', 'CLOSED', 'CANCELLED')),
    delivery_address TEXT,
    delivery_terms VARCHAR(50) DEFAULT 'EXW',
    payment_terms VARCHAR(50) DEFAULT 'Net 30',
    expected_delivery_date DATE,
    actual_delivery_date DATE,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    grand_total DECIMAL(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    notes TEXT,
    created_by UUID NOT NULL,
    submitted_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    confirmed_by UUID,
    cancelled_at TIMESTAMP,
    cancelled_by UUID,
    cancellation_reason TEXT,
    closed_at TIMESTAMP,
    amendment_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_po_number ON purchase_orders(po_number);
CREATE INDEX idx_po_status ON purchase_orders(status);
CREATE INDEX idx_po_supplier ON purchase_orders(supplier_id);
CREATE INDEX idx_po_pr ON purchase_orders(pr_id);
CREATE INDEX idx_po_expected_date ON purchase_orders(expected_delivery_date);
CREATE INDEX idx_po_deleted_at ON purchase_orders(deleted_at);
