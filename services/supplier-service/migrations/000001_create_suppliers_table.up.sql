-- Suppliers master table
CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    tax_code VARCHAR(50),
    supplier_type VARCHAR(20) NOT NULL CHECK (supplier_type IN ('MANUFACTURER', 'TRADER', 'IMPORTER')),
    business_type VARCHAR(20) NOT NULL DEFAULT 'DOMESTIC' CHECK (business_type IN ('DOMESTIC', 'INTERNATIONAL')),
    email VARCHAR(255),
    phone VARCHAR(50),
    fax VARCHAR(50),
    website VARCHAR(255),
    payment_terms VARCHAR(50) DEFAULT 'Net 30',
    currency VARCHAR(3) DEFAULT 'VND',
    credit_limit DECIMAL(15,2) DEFAULT 0,
    bank_name VARCHAR(255),
    bank_account VARCHAR(100),
    bank_branch VARCHAR(255),
    quality_rating DECIMAL(3,2) DEFAULT 0,
    delivery_rating DECIMAL(3,2) DEFAULT 0,
    service_rating DECIMAL(3,2) DEFAULT 0,
    overall_rating DECIMAL(3,2) DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'BLOCKED', 'INACTIVE')),
    blocked_reason TEXT,
    blocked_by UUID,
    blocked_at TIMESTAMP,
    approved_by UUID,
    approved_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_suppliers_code ON suppliers(code);
CREATE INDEX idx_suppliers_name ON suppliers(name);
CREATE INDEX idx_suppliers_status ON suppliers(status);
CREATE INDEX idx_suppliers_supplier_type ON suppliers(supplier_type);
CREATE INDEX idx_suppliers_deleted_at ON suppliers(deleted_at);
