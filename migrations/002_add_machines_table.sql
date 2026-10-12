-- Migration: Add machines table and related updates
-- Version: 006
-- Description: Creates machines table, updates user_machines foreign key, adds indexes and triggers

-- ============================================================================
-- CREATE MACHINES TABLE
-- ============================================================================

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

-- ============================================================================
-- UPDATE USER_MACHINES FOREIGN KEY
-- ============================================================================

-- Drop existing foreign key constraint if it exists (for safety)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'user_machines_machine_id_fkey'
    ) THEN
        ALTER TABLE user_machines DROP CONSTRAINT user_machines_machine_id_fkey;
    END IF;
END $$;

-- Add proper foreign key constraint to machines table
ALTER TABLE user_machines
    ADD CONSTRAINT user_machines_machine_id_fkey
    FOREIGN KEY (machine_id) REFERENCES machines(id) ON DELETE CASCADE;

-- ============================================================================
-- CREATE CATEGORY-MACHINES JUNCTION TABLE
-- ============================================================================

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

-- ============================================================================
-- CREATE INDEXES FOR MACHINES
-- ============================================================================

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
-- CREATE MISSING TRIGGERS FOR UPDATED_AT
-- ============================================================================

-- Function to update updated_at timestamp (skip if already exists)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for categories
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create triggers for brands
DROP TRIGGER IF EXISTS update_brands_updated_at ON brands;
CREATE TRIGGER update_brands_updated_at BEFORE UPDATE ON brands
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create triggers for products
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create triggers for controllers
DROP TRIGGER IF EXISTS update_controllers_updated_at ON controllers;
CREATE TRIGGER update_controllers_updated_at BEFORE UPDATE ON controllers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create triggers for machines
DROP TRIGGER IF EXISTS update_machines_updated_at ON machines;
CREATE TRIGGER update_machines_updated_at BEFORE UPDATE ON machines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

-- Comment to verify migration ran successfully
COMMENT ON TABLE machines IS 'Machine inventory and details - Added in migration 002';
