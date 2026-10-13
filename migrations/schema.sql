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



-- Categories table (Product categories)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    description TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Brands table (Machine brands)
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    description TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table (Machine products)
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    image VARCHAR(255) DEFAULT 'no_image_available.png',
    description TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Controllers table (Machine controllers)
CREATE TABLE IF NOT EXISTS controllers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50),
    variant VARCHAR(100),
    brand_id INTEGER REFERENCES brands(id) ON DELETE SET NULL,
    description TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Controller products junction table (Many-to-many relationship)
CREATE TABLE IF NOT EXISTS controller_products (
    id SERIAL PRIMARY KEY,
    controller_id INTEGER REFERENCES controllers(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(controller_id, product_id)
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

-- Create indexes for categories
CREATE UNIQUE INDEX IF NOT EXISTS categories_name_key ON categories (name) WHERE is_deleted = false;
CREATE UNIQUE INDEX IF NOT EXISTS categories_code_key ON categories (code) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_categories_is_deleted ON categories(is_deleted);

-- Create indexes for brands
CREATE UNIQUE INDEX IF NOT EXISTS brands_name_key ON brands (name) WHERE is_deleted = false;
CREATE UNIQUE INDEX IF NOT EXISTS brands_code_key ON brands (code) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_brands_is_deleted ON brands(is_deleted);

-- Create indexes for products
CREATE UNIQUE INDEX IF NOT EXISTS products_name_key ON products (name) WHERE is_deleted = false;
CREATE UNIQUE INDEX IF NOT EXISTS products_code_key ON products (code) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_deleted ON products(is_deleted);

-- Create indexes for controllers
CREATE UNIQUE INDEX IF NOT EXISTS controllers_name_key ON controllers (name) WHERE is_deleted = false;
CREATE UNIQUE INDEX IF NOT EXISTS controllers_code_key ON controllers (code) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_controllers_brand_id ON controllers(brand_id);
CREATE INDEX IF NOT EXISTS idx_controllers_is_deleted ON controllers(is_deleted);

-- Create indexes for controller_products
CREATE INDEX IF NOT EXISTS idx_controller_products_controller_id ON controller_products(controller_id);
CREATE INDEX IF NOT EXISTS idx_controller_products_product_id ON controller_products(product_id);
CREATE INDEX IF NOT EXISTS idx_controller_products_is_deleted ON controller_products(is_deleted);

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
    company_name VARCHAR(150),
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

-- Machines table (Machine inventory and details)
CREATE TABLE IF NOT EXISTS machines (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    machine_type_id INTEGER,
    product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
    name VARCHAR(50) NOT NULL,
    serial_number VARCHAR(50) NOT NULL,
    live_ip VARCHAR(20) NOT NULL,
    port INTEGER NOT NULL,
    username VARCHAR(15) NOT NULL,
    password VARCHAR(255) NOT NULL, -- Encrypted password
    local_ip VARCHAR(20) NOT NULL,
    brand_id INTEGER REFERENCES brands(id) ON DELETE SET NULL,
    controller_id INTEGER REFERENCES controllers(id) ON DELETE SET NULL,
    controller_version VARCHAR(10) NOT NULL,
    accuracy_id INTEGER,
    line_id INTEGER,
    installation_date DATE NOT NULL,
    bill_date DATE NOT NULL,
    warranty_date DATE NOT NULL,
    amc VARCHAR(50),
    is_insured BOOLEAN NOT NULL DEFAULT false,
    online_status BOOLEAN NOT NULL DEFAULT false, -- true=ONLINE, false=OFFLINE
    maintenance_required BOOLEAN NOT NULL DEFAULT false,
    hourly_cost DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'inactive', 'deleted'
    is_monitor BOOLEAN NOT NULL DEFAULT false,
    online_date DATE,
    offline_date TIMESTAMPTZ,
    mode_id INTEGER,
    controller_version_type INTEGER, -- 1=828D, 2=840D
    health_ids VARCHAR(100),
    longitude VARCHAR(50),
    latitude VARCHAR(50),
    cycle_analysis BOOLEAN NOT NULL DEFAULT false,
    machine_system_type VARCHAR(20) NOT NULL DEFAULT 'NewPlc', -- 'NewPlc', 'OldPlc'
    security_mode VARCHAR(20) NOT NULL DEFAULT 'Basic256Sha256', -- 'None', 'Basic256Sha256'
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Category machines junction table (Many-to-many relationship)
-- Allows categories to directly reference machines (replaces comma-separated iConfigIds from old DB)
CREATE TABLE IF NOT EXISTS category_machines (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    machine_id INTEGER REFERENCES machines(id) ON DELETE CASCADE,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, machine_id)
);

-- User machines table (Machine assignments to users)
CREATE TABLE IF NOT EXISTS user_machines (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    machine_id INTEGER REFERENCES machines(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    UNIQUE(user_id, machine_id)
);

-- Create index for user machine lookups
CREATE INDEX IF NOT EXISTS idx_user_machines_user_id ON user_machines(user_id);
CREATE INDEX IF NOT EXISTS idx_user_machines_machine_id ON user_machines(machine_id);
CREATE INDEX IF NOT EXISTS idx_user_machines_is_deleted ON user_machines(is_deleted);

-- Create indexes for machines
CREATE UNIQUE INDEX IF NOT EXISTS machines_name_key ON machines (name) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_id ON machines(user_id);
CREATE INDEX IF NOT EXISTS idx_machines_product_id ON machines(product_id);
CREATE INDEX IF NOT EXISTS idx_machines_brand_id ON machines(brand_id);
CREATE INDEX IF NOT EXISTS idx_machines_controller_id ON machines(controller_id);
CREATE INDEX IF NOT EXISTS idx_machines_online_status ON machines(online_status);
CREATE INDEX IF NOT EXISTS idx_machines_status ON machines(status);
CREATE INDEX IF NOT EXISTS idx_machines_is_deleted ON machines(is_deleted);

-- Create indexes for category_machines
CREATE INDEX IF NOT EXISTS idx_category_machines_category_id ON category_machines(category_id);
CREATE INDEX IF NOT EXISTS idx_category_machines_machine_id ON category_machines(machine_id);
CREATE INDEX IF NOT EXISTS idx_category_machines_is_deleted ON category_machines(is_deleted);

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

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_brands_updated_at BEFORE UPDATE ON brands
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_controllers_updated_at BEFORE UPDATE ON controllers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_machines_updated_at BEFORE UPDATE ON machines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_company_profiles_updated_at BEFORE UPDATE ON company_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
