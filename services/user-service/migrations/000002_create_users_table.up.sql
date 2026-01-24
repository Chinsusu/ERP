-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    employee_code VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    manager_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_employee_code ON users(employee_code);
CREATE INDEX idx_users_department_id ON users(department_id);
CREATE INDEX idx_users_manager_id ON users(manager_id);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_name ON users(first_name, last_name);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();

-- Add foreign key to departments table for manager_id
ALTER TABLE departments
    ADD CONSTRAINT fk_departments_manager
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL;
