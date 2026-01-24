-- PR Approvals history table
CREATE TABLE IF NOT EXISTS pr_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pr_id UUID NOT NULL REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    approver_name VARCHAR(255),
    approval_level VARCHAR(30) NOT NULL,
    action VARCHAR(20) NOT NULL CHECK (action IN ('APPROVED', 'REJECTED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_pr_approvals_pr_id ON pr_approvals(pr_id);
CREATE INDEX idx_pr_approvals_approver ON pr_approvals(approver_id);
