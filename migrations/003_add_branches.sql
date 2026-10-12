-- Add branches table, role_id, branch_id, and company_id to users
-- This migration adds branch, role, and company hierarchy support to the user system

BEGIN;

-- Create branches table
CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL, -- Branch code (e.g., 'BR_0600', 'HO_0200')
    name VARCHAR(100) NOT NULL, -- Branch name (e.g., 'Rajkot', 'Metoda')
    email VARCHAR(100),
    address TEXT,
    city VARCHAR(100),
    state_id INTEGER REFERENCES states(id) ON DELETE SET NULL,
    customer_id INTEGER REFERENCES users(id) ON DELETE SET NULL, -- Customer/company user who owns this branch
    gst_number VARCHAR(30),
    pan_number VARCHAR(30),
    color_code VARCHAR(10),
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive'
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for branches
CREATE INDEX IF NOT EXISTS idx_branches_state_id ON branches(state_id);
CREATE INDEX IF NOT EXISTS idx_branches_customer_id ON branches(customer_id);
CREATE INDEX IF NOT EXISTS idx_branches_status ON branches(status);
CREATE INDEX IF NOT EXISTS idx_branches_is_deleted ON branches(is_deleted);

-- Add role_id column to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE SET NULL;

-- Add branch_id column to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS branch_id INTEGER REFERENCES branches(id) ON DELETE SET NULL;

-- Add customer_id column to branches table
ALTER TABLE branches
    ADD COLUMN IF NOT EXISTS customer_id INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Make role_id NOT NULL (if it exists and doesn't have data)
DO $$
BEGIN
    -- Check if role_id exists and has NULL values
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'role_id') THEN
        -- If no NULL values, make it NOT NULL
        IF NOT EXISTS (SELECT 1 FROM users WHERE role_id IS NULL) THEN
            ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;
        END IF;
    END IF;
END $$;

-- Add customer_id column to branches table
ALTER TABLE branches
    ADD COLUMN IF NOT EXISTS customer_id INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Add company_id column to users table (self-referencing for company users and sub-users)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS company_id INTEGER REFERENCES users(id) ON DELETE SET NULL;

-- Create indexes for user lookups
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_branch_id ON users(branch_id);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);

-- Insert sample branches based on old database
INSERT INTO branches (code, name, email, address, city, gst_number, pan_number, color_code, status) VALUES
('BR_0600', 'Rajkot', 'shailesh.dodiya@jyoti.co.in', 'Samrat Industrial Area Main Road, B/h S.T. Work Shop,Gondal Road, Rajkot - 360 004.', 'Rajkot', '24AABCJ1947R1Z5', 'AABCJ1947R', '#f56954', 'active'),
('HO_0200', 'Metoda', 'metoda.service@jyoti.co.in', 'Plot No.,P-4, Lodhika GIDC, Vill.-Metoda, Dist.-Rajkot - 360 021.', 'Metoda', '24AABCJ1947R1Z5', 'AABCJ1947R', '#f39c12', 'active')
ON CONFLICT (code) DO NOTHING;

-- Create trigger for updated_at on branches
CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON branches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
