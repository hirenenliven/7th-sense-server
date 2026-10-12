-- Migration: Add product-related tables (categories, brands, products, controllers)
-- Version: 005
-- Description: Creates tables for machine products, brands, controllers, and categories

BEGIN;

-- Categories table (Product categories)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE,
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
    code VARCHAR(50) UNIQUE,
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
    code VARCHAR(50) UNIQUE,
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
    code VARCHAR(50) UNIQUE,
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
-- CREATE INDEXES
-- ============================================================================

-- Indexes for categories
CREATE UNIQUE INDEX IF NOT EXISTS categories_name_key ON categories (name) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_categories_is_deleted ON categories(is_deleted);

-- Indexes for brands
CREATE UNIQUE INDEX IF NOT EXISTS brands_name_key ON brands (name) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_brands_is_deleted ON brands(is_deleted);

-- Indexes for products
CREATE UNIQUE INDEX IF NOT EXISTS products_name_key ON products (name) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_deleted ON products(is_deleted);

-- Indexes for controllers
CREATE UNIQUE INDEX IF NOT EXISTS controllers_name_key ON controllers (name) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_controllers_brand_id ON controllers(brand_id);
CREATE INDEX IF NOT EXISTS idx_controllers_is_deleted ON controllers(is_deleted);

-- Indexes for controller_products
CREATE INDEX IF NOT EXISTS idx_controller_products_controller_id ON controller_products(controller_id);
CREATE INDEX IF NOT EXISTS idx_controller_products_product_id ON controller_products(product_id);
CREATE INDEX IF NOT EXISTS idx_controller_products_is_deleted ON controller_products(is_deleted);

-- ============================================================================
-- CREATE TRIGGERS FOR UPDATED_AT
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

COMMIT;

-- Comment to verify migration ran successfully
COMMENT ON TABLE categories IS 'Product categories - Added in migration 005';
COMMENT ON TABLE brands IS 'Machine brands - Added in migration 005';
COMMENT ON TABLE products IS 'Machine products - Added in migration 005';
COMMENT ON TABLE controllers IS 'Machine controllers - Added in migration 005';
COMMENT ON TABLE controller_products IS 'Controller-Product relationships - Added in migration 005';
