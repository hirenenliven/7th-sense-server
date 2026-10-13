-- Migration: Add filters array to user_profiles table
-- Version: 010
-- Description: Adds a TEXT array column to store user filter preferences

BEGIN;

-- Add filters column as TEXT array to user_profiles table
ALTER TABLE user_profiles
ADD COLUMN IF NOT EXISTS filters TEXT[] DEFAULT '{}';

COMMIT;
