-- Supplier addresses (multiple per supplier)
CREATE TABLE IF NOT EXISTS supplier_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL CHECK (address_type IN ('BILLING', 'SHIPPING', 'FACTORY', 'OFFICE', 'WAREHOUSE')),
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    ward VARCHAR(100),
    district VARCHAR(100),
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100),
    country VARCHAR(100) NOT NULL DEFAULT 'Vietnam',
    postal_code VARCHAR(20),
    is_primary BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_supplier_addresses_supplier_id ON supplier_addresses(supplier_id);
CREATE INDEX idx_supplier_addresses_type ON supplier_addresses(address_type);
