-- Migration: clienteles keep country and state names instead of ISO codes
--
-- country_code / state_code (migration 017) become country_name / state_name
-- VARCHAR(100), holding names such as "India" and "Gujarat" (like the customer
-- details). Safe to run on a database where the columns were already renamed.
--
-- Rollback:
--   ALTER TABLE clienteles RENAME COLUMN country_name TO country_code;
--   ALTER TABLE clienteles RENAME COLUMN state_name TO state_code;
--   (and shrink them back to VARCHAR(2) / VARCHAR(10) once they hold codes again)

BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'clienteles' AND column_name = 'country_code') THEN
        ALTER TABLE clienteles RENAME COLUMN country_code TO country_name;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'clienteles' AND column_name = 'state_code') THEN
        ALTER TABLE clienteles RENAME COLUMN state_code TO state_name;
    END IF;
END $$;

ALTER TABLE clienteles ALTER COLUMN country_name TYPE VARCHAR(100);
ALTER TABLE clienteles ALTER COLUMN state_name TYPE VARCHAR(100);

COMMIT;
