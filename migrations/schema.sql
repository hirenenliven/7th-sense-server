-- 7th Sense Server Database Schema
-- PostgreSQL Database Creation Script
-- Run this to create all tables in the database (fresh setup)

-- ============================================================================
-- LOOKUP TABLES (Reference Data)
-- ============================================================================

-- Countries table
CREATE TABLE IF NOT EXISTS countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(3) UNIQUE NOT NULL, -- ISO 3-letter code (e.g., IND, USA)
    phone_code VARCHAR(10) NOT NULL, -- Dialing code (e.g., +91, +1)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- States table
CREATE TABLE IF NOT EXISTS states (
    id SERIAL PRIMARY KEY,
    country_id INTEGER REFERENCES countries(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- 'admin', 'supervisor', 'operator', 'customer'
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Branches table
CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL, -- Branch code (e.g., 'BR_0600', 'HO_0200')
    name VARCHAR(100) NOT NULL, -- Branch name (e.g., 'Rajkot', 'Metoda')
    email VARCHAR(100),
    address TEXT,
    customer_id INTEGER REFERENCES users(id) ON DELETE SET NULL, -- Customer/company user who owns this branch
    city VARCHAR(100),
    state_id INTEGER REFERENCES states(id) ON DELETE SET NULL,
    gst_number VARCHAR(30),
    pan_number VARCHAR(30),
    color_code VARCHAR(10),
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive'
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- CORE USER TABLE
-- ============================================================================

-- Users table (Core user data)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    mobile VARCHAR(20),
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE SET NULL,
    branch_id INTEGER REFERENCES branches(id) ON DELETE SET NULL,
    company_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
-- Unique indexes only apply to non-deleted records
CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email) WHERE is_deleted = false;
CREATE UNIQUE INDEX IF NOT EXISTS users_username_key ON users (username) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_branch_id ON users(branch_id);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);

-- Create indexes for branches
CREATE INDEX IF NOT EXISTS idx_branches_state_id ON branches(state_id);
CREATE INDEX IF NOT EXISTS idx_branches_customer_id ON branches(customer_id);
CREATE INDEX IF NOT EXISTS idx_branches_status ON branches(status);
CREATE INDEX IF NOT EXISTS idx_branches_is_deleted ON branches(is_deleted);

-- ============================================================================
-- ADDRESSES TABLE
-- ============================================================================

-- Addresses table (Multiple addresses per user)
CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    address_type VARCHAR(20) DEFAULT 'home', -- 'home', 'work', 'shipping', 'billing'
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state_id INTEGER REFERENCES states(id) ON DELETE SET NULL,
    country_id INTEGER REFERENCES countries(id) ON DELETE SET NULL,
    pin_code VARCHAR(20),
    contact_person VARCHAR(100),
    contact_phone VARCHAR(20),
    is_default BOOLEAN DEFAULT false,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for address lookups
CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);
CREATE INDEX IF NOT EXISTS idx_addresses_is_deleted ON addresses(is_deleted);

-- ============================================================================
-- USER PROFILES TABLE
-- ============================================================================

-- User profiles table (Extended user information)
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    designation VARCHAR(100), -- Job title/position
    grade VARCHAR(50), -- Employee grade
    birth_date DATE,
    joining_date DATE,
    blood_group VARCHAR(5), -- A+, B+, O+, etc.
    profile_picture VARCHAR(255), -- Path or URL to profile image
    contact_person VARCHAR(100), -- Emergency contact
    contact_phone VARCHAR(20), -- Emergency contact phone
    theme VARCHAR(20) DEFAULT 'light_version', -- UI theme preference
    timezone VARCHAR(50) DEFAULT 'Asia/Kolkata',
    language VARCHAR(10) DEFAULT 'en',
    start_day_time TIME, -- Day start time for operations
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for user profile lookups
CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_is_deleted ON user_profiles(is_deleted);

-- ============================================================================
-- COMPANY PROFILES TABLE
-- ============================================================================

-- Company profiles table (Company-specific details for customer users)
CREATE TABLE IF NOT EXISTS company_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    gst_number VARCHAR(30),
    pan_number VARCHAR(30),
    website VARCHAR(255),
    finix_id INTEGER, -- External system ID
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for company profile lookups
CREATE INDEX IF NOT EXISTS idx_company_profiles_user_id ON company_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_company_profiles_is_deleted ON company_profiles(is_deleted);

-- ============================================================================
-- USER MACHINES TABLE
-- ============================================================================

-- User machines table (Machine assignments to users)
CREATE TABLE IF NOT EXISTS user_machines (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    machine_id INTEGER NOT NULL, -- Reference to machines table (to be created)
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, machine_id)
);

-- Create index for user machine lookups
CREATE INDEX IF NOT EXISTS idx_user_machines_user_id ON user_machines(user_id);
CREATE INDEX IF NOT EXISTS idx_user_machines_machine_id ON user_machines(machine_id);
CREATE INDEX IF NOT EXISTS idx_user_machines_is_deleted ON user_machines(is_deleted);

-- ============================================================================
-- AUTHENTICATION TABLES
-- ============================================================================

-- Auth tokens table (JWT, refresh tokens, password reset, email verification)
CREATE TABLE IF NOT EXISTS auth_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    token_type VARCHAR(20) NOT NULL, -- 'api', 'refresh', 'reset_password', 'email_verification'
    device_type VARCHAR(50), -- 'web', 'mobile', 'api'
    device_info VARCHAR(255), -- User agent, IP, etc.
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for token lookups
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON auth_tokens(token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens(expires_at);

-- ============================================================================
-- SAMPLE DATA (For Testing)
-- ============================================================================

-- Insert default roles
INSERT INTO roles (name, description) VALUES
('admin', 'System administrator with full access'),
('supervisor', 'Supervisor with management permissions'),
('operator', 'Machine operator'),
('customer', 'Customer with limited access'),
('qc', 'Quality control staff'),
('shop_head', 'Shop head/manager')
ON CONFLICT (name) DO NOTHING;

-- Insert sample branches
INSERT INTO branches (code, name, email, address, customer_id, city, gst_number, pan_number, color_code, status) VALUES
('BR_0600', 'Rajkot', 'shailesh.dodiya@jyoti.co.in', 'Samrat Industrial Area Main Road, B/h S.T. Work Shop,Gondal Road, Rajkot - 360 004.', 6, 'Rajkot', '24AABCJ1947R1Z5', 'AABCJ1947R', '#f56954', 'active'),
('HO_0200', 'Metoda', 'metoda.service@jyoti.co.in', 'Plot No.,P-4, Lodhika GIDC, Vill.-Metoda, Dist.-Rajkot - 360 021.', 6, 'Metoda', '24AABCJ1947R1Z5', 'AABCJ1947R', '#f39c12', 'active')
ON CONFLICT (code) DO NOTHING;

-- Insert default country
INSERT INTO countries (name, code, phone_code) VALUES
('India', 'IND', '+91'),
('United States', 'USA', '+1')
ON CONFLICT (code) DO NOTHING;

-- Insert default state
INSERT INTO states (country_id, name, code) VALUES
(1, 'Gujarat', 'GJ'),
(1, 'Maharashtra', 'MH'),
(1, 'Karnataka', 'KA'),
(2, 'California', 'CA')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- FUNCTIONS FOR AUTOMATIC UPDATED_AT
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON branches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_company_profiles_updated_at BEFORE UPDATE ON company_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
