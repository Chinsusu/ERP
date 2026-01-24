-- Create departments table
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    manager_id UUID, -- Will reference users table (added later)
    level INT NOT NULL DEFAULT 0,
    path VARCHAR(500) NOT NULL DEFAULT '/',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_departments_code ON departments(code);
CREATE INDEX idx_departments_parent_id ON departments(parent_id);
CREATE INDEX idx_departments_path ON departments(path);
CREATE INDEX idx_departments_status ON departments(status);
CREATE INDEX idx_departments_deleted_at ON departments(deleted_at);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_departments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER departments_updated_at
    BEFORE UPDATE ON departments
    FOR EACH ROW
    EXECUTE FUNCTION update_departments_updated_at();
