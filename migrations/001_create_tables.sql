-- 7th Sense Server Database Schema
-- PostgreSQL Database Creation Script
-- Run this to create all tables in the database

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

-- ============================================================================
-- CORE USER TABLE
-- ============================================================================

-- Users table (Core user data)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    mobile VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for address lookups
CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);

-- ============================================================================
-- RELATIONSHIP TABLES (Many-to-Many)
-- ============================================================================

-- User roles (Many-to-Many)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

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
