-- Migration: Add company name to company profiles
-- Version: 007
-- Description: The customer form has a required Company Name, which no table stored.

BEGIN;

ALTER TABLE company_profiles ADD COLUMN IF NOT EXISTS company_name VARCHAR(150);

COMMIT;
