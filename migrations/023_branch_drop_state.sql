-- Migration: branches drop state_id — a branch's state comes through its city (user_city_id)
--
-- Rollback:
--   ALTER TABLE branches ADD COLUMN state_id INTEGER REFERENCES states(id) ON DELETE SET NULL;

BEGIN;

ALTER TABLE branches DROP COLUMN IF EXISTS state_id;

COMMIT;
