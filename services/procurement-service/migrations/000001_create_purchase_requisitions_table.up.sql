-- Purchase Requisitions (PR) table
CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pr_number VARCHAR(20) NOT NULL UNIQUE,
    pr_date DATE NOT NULL DEFAULT CURRENT_DATE,
    required_date DATE NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL' CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT')),
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SUBMITTED', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED', 'CONVERTED_TO_PO', 'CANCELLED')),
    requester_id UUID NOT NULL,
    department_id UUID,
    justification TEXT,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    approval_level VARCHAR(30),
    notes TEXT,
    submitted_at TIMESTAMP,
    approved_at TIMESTAMP,
    approved_by UUID,
    rejected_at TIMESTAMP,
    rejected_by UUID,
    rejection_reason TEXT,
    po_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_pr_number ON purchase_requisitions(pr_number);
CREATE INDEX idx_pr_status ON purchase_requisitions(status);
CREATE INDEX idx_pr_requester ON purchase_requisitions(requester_id);
CREATE INDEX idx_pr_required_date ON purchase_requisitions(required_date);
CREATE INDEX idx_pr_deleted_at ON purchase_requisitions(deleted_at);
