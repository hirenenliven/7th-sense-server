-- Migration: Add organization tables
-- Version: 011
-- Description: Per-user organization hierarchy
--   user_country → user_state → user_city → user_branch → user_location
--   → user_shop → user_line → user_cell → user_supervisor
-- Every row belongs to one user (user_id). Names are unique per user within the same
-- parent, case-insensitively, ignoring soft-deleted rows (so a deleted name can be reused).
-- Parent FKs are ON DELETE RESTRICT; the service also refuses to soft-delete a row that
-- still has active children.
-- Rollback: 011_add_organization_tables_down.sql

BEGIN;

-- Country (no parent)
CREATE TABLE IF NOT EXISTS user_country (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_country_name_key ON user_country (user_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_country_user_id ON user_country(user_id);
CREATE INDEX IF NOT EXISTS idx_user_country_is_deleted ON user_country(is_deleted);
DROP TRIGGER IF EXISTS update_user_country_updated_at ON user_country;
CREATE TRIGGER update_user_country_updated_at BEFORE UPDATE ON user_country
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- State → Country
CREATE TABLE IF NOT EXISTS user_state (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_country_id INTEGER NOT NULL REFERENCES user_country(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_state_name_key ON user_state (user_id, user_country_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_state_user_id ON user_state(user_id);
CREATE INDEX IF NOT EXISTS idx_user_state_user_country_id ON user_state(user_country_id);
CREATE INDEX IF NOT EXISTS idx_user_state_is_deleted ON user_state(is_deleted);
DROP TRIGGER IF EXISTS update_user_state_updated_at ON user_state;
CREATE TRIGGER update_user_state_updated_at BEFORE UPDATE ON user_state
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- City → State
CREATE TABLE IF NOT EXISTS user_city (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_state_id INTEGER NOT NULL REFERENCES user_state(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_city_name_key ON user_city (user_id, user_state_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_city_user_id ON user_city(user_id);
CREATE INDEX IF NOT EXISTS idx_user_city_user_state_id ON user_city(user_state_id);
CREATE INDEX IF NOT EXISTS idx_user_city_is_deleted ON user_city(is_deleted);
DROP TRIGGER IF EXISTS update_user_city_updated_at ON user_city;
CREATE TRIGGER update_user_city_updated_at BEFORE UPDATE ON user_city
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Branch → City
CREATE TABLE IF NOT EXISTS user_branch (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_city_id INTEGER NOT NULL REFERENCES user_city(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_branch_name_key ON user_branch (user_id, user_city_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_branch_user_id ON user_branch(user_id);
CREATE INDEX IF NOT EXISTS idx_user_branch_user_city_id ON user_branch(user_city_id);
CREATE INDEX IF NOT EXISTS idx_user_branch_is_deleted ON user_branch(is_deleted);
DROP TRIGGER IF EXISTS update_user_branch_updated_at ON user_branch;
CREATE TRIGGER update_user_branch_updated_at BEFORE UPDATE ON user_branch
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Location → Branch
CREATE TABLE IF NOT EXISTS user_location (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_branch_id INTEGER NOT NULL REFERENCES user_branch(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_location_name_key ON user_location (user_id, user_branch_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_location_user_id ON user_location(user_id);
CREATE INDEX IF NOT EXISTS idx_user_location_user_branch_id ON user_location(user_branch_id);
CREATE INDEX IF NOT EXISTS idx_user_location_is_deleted ON user_location(is_deleted);
DROP TRIGGER IF EXISTS update_user_location_updated_at ON user_location;
CREATE TRIGGER update_user_location_updated_at BEFORE UPDATE ON user_location
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Shop → Location
CREATE TABLE IF NOT EXISTS user_shop (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_location_id INTEGER NOT NULL REFERENCES user_location(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_shop_name_key ON user_shop (user_id, user_location_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_shop_user_id ON user_shop(user_id);
CREATE INDEX IF NOT EXISTS idx_user_shop_user_location_id ON user_shop(user_location_id);
CREATE INDEX IF NOT EXISTS idx_user_shop_is_deleted ON user_shop(is_deleted);
DROP TRIGGER IF EXISTS update_user_shop_updated_at ON user_shop;
CREATE TRIGGER update_user_shop_updated_at BEFORE UPDATE ON user_shop
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Line → Shop
CREATE TABLE IF NOT EXISTS user_line (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_shop_id INTEGER NOT NULL REFERENCES user_shop(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_line_name_key ON user_line (user_id, user_shop_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_line_user_id ON user_line(user_id);
CREATE INDEX IF NOT EXISTS idx_user_line_user_shop_id ON user_line(user_shop_id);
CREATE INDEX IF NOT EXISTS idx_user_line_is_deleted ON user_line(is_deleted);
DROP TRIGGER IF EXISTS update_user_line_updated_at ON user_line;
CREATE TRIGGER update_user_line_updated_at BEFORE UPDATE ON user_line
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Cell → Line
CREATE TABLE IF NOT EXISTS user_cell (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_line_id INTEGER NOT NULL REFERENCES user_line(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_cell_name_key ON user_cell (user_id, user_line_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_cell_user_id ON user_cell(user_id);
CREATE INDEX IF NOT EXISTS idx_user_cell_user_line_id ON user_cell(user_line_id);
CREATE INDEX IF NOT EXISTS idx_user_cell_is_deleted ON user_cell(is_deleted);
DROP TRIGGER IF EXISTS update_user_cell_updated_at ON user_cell;
CREATE TRIGGER update_user_cell_updated_at BEFORE UPDATE ON user_cell
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Supervisor → Cell
CREATE TABLE IF NOT EXISTS user_supervisor (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_cell_id INTEGER NOT NULL REFERENCES user_cell(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS user_supervisor_name_key ON user_supervisor (user_id, user_cell_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_supervisor_user_id ON user_supervisor(user_id);
CREATE INDEX IF NOT EXISTS idx_user_supervisor_user_cell_id ON user_supervisor(user_cell_id);
CREATE INDEX IF NOT EXISTS idx_user_supervisor_is_deleted ON user_supervisor(is_deleted);
DROP TRIGGER IF EXISTS update_user_supervisor_updated_at ON user_supervisor;
CREATE TRIGGER update_user_supervisor_updated_at BEFORE UPDATE ON user_supervisor
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add user_city_id to branches table (nullable for global branches)
ALTER TABLE branches
    ADD COLUMN IF NOT EXISTS user_city_id INTEGER REFERENCES user_city(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_branches_user_city_id ON branches(user_city_id);

COMMIT;
