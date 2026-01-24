-- Supplier evaluations (quarterly/annual performance reviews)
CREATE TABLE IF NOT EXISTS supplier_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    evaluation_date DATE NOT NULL,
    evaluation_period VARCHAR(20) NOT NULL,  -- e.g., "2024-Q1", "2024-H1", "2024"
    quality_score DECIMAL(3,2) NOT NULL CHECK (quality_score >= 1 AND quality_score <= 5),
    delivery_score DECIMAL(3,2) NOT NULL CHECK (delivery_score >= 1 AND delivery_score <= 5),
    price_score DECIMAL(3,2) NOT NULL CHECK (price_score >= 1 AND price_score <= 5),
    service_score DECIMAL(3,2) NOT NULL CHECK (service_score >= 1 AND service_score <= 5),
    documentation_score DECIMAL(3,2) NOT NULL CHECK (documentation_score >= 1 AND documentation_score <= 5),
    overall_score DECIMAL(3,2) NOT NULL,
    on_time_delivery_rate DECIMAL(5,2),  -- percentage
    quality_acceptance_rate DECIMAL(5,2),  -- percentage
    lead_time_adherence DECIMAL(5,2),  -- percentage
    strengths TEXT,
    weaknesses TEXT,
    action_items TEXT,
    evaluated_by UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SUBMITTED', 'APPROVED')),
    approved_by UUID,
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_supplier_evaluations_supplier_id ON supplier_evaluations(supplier_id);
CREATE INDEX idx_supplier_evaluations_period ON supplier_evaluations(evaluation_period);
CREATE INDEX idx_supplier_evaluations_status ON supplier_evaluations(status);
