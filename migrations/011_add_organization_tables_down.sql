-- Rollback of migration 011 (organization tables). Drops children before parents.
-- WARNING: deletes every organization record of every user.

BEGIN;

DROP TABLE IF EXISTS user_supervisor;
DROP TABLE IF EXISTS user_cell;
DROP TABLE IF EXISTS user_line;
DROP TABLE IF EXISTS user_shop;
DROP TABLE IF EXISTS user_location;
DROP TABLE IF EXISTS user_branch;
DROP TABLE IF EXISTS user_city;
DROP TABLE IF EXISTS user_state;
DROP TABLE IF EXISTS user_country;

-- Remove user_city_id from branches table
ALTER TABLE branches
    DROP COLUMN IF EXISTS user_city_id;

COMMIT;
