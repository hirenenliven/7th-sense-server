-- Add user profiles, company profiles, and user machines tables
-- This migration adds extended user information tables

BEGIN;

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

-- Create indexes for user profiles
CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_is_deleted ON user_profiles(is_deleted);

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

-- Create indexes for company profiles
CREATE INDEX IF NOT EXISTS idx_company_profiles_user_id ON company_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_company_profiles_is_deleted ON company_profiles(is_deleted);

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

-- Create indexes for user machines
CREATE INDEX IF NOT EXISTS idx_user_machines_user_id ON user_machines(user_id);
CREATE INDEX IF NOT EXISTS idx_user_machines_machine_id ON user_machines(machine_id);
CREATE INDEX IF NOT EXISTS idx_user_machines_is_deleted ON user_machines(is_deleted);

-- Create triggers for updated_at
CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_company_profiles_updated_at BEFORE UPDATE ON company_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
