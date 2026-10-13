-- Migration: Use branches in the organization hierarchy instead of user_branch
-- Version: 013
-- Description: The hierarchy is now user_country → user_state → user_city → branches → user_location.
-- branches.user_city_id (added in 011) links a branch to its city, and user_location.branch_id
-- replaces user_location.user_branch_id. Existing user_branch rows are copied into branches
-- (code 'UB' + id, owner as customer_id, soft-delete state kept) so their locations keep a parent,
-- then user_branch is dropped.

BEGIN;

-- 1. Copy user_branch rows into branches and remember the new ids.
CREATE TEMP TABLE branch_id_map (user_branch_id INTEGER PRIMARY KEY, branch_id INTEGER NOT NULL) ON COMMIT DROP;

WITH copied AS (
    INSERT INTO branches (code, name, customer_id, user_city_id, status, is_deleted, deleted_at, created_at, updated_at)
    SELECT 'UB' || lpad(id::text, 6, '0'), name, user_id, user_city_id, 'active', is_deleted, deleted_at, created_at, updated_at
    FROM user_branch
    RETURNING id, code
)
INSERT INTO branch_id_map (user_branch_id, branch_id)
SELECT substring(code FROM 3)::integer, id FROM copied;

-- 2. Point locations at branches.
ALTER TABLE user_location ADD COLUMN branch_id INTEGER;

UPDATE user_location
SET branch_id = branch_id_map.branch_id
FROM branch_id_map
WHERE user_location.user_branch_id = branch_id_map.user_branch_id;

ALTER TABLE user_location ALTER COLUMN branch_id SET NOT NULL;
ALTER TABLE user_location ADD CONSTRAINT user_location_branch_id_fkey
    FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE RESTRICT;

DROP INDEX IF EXISTS user_location_name_key;
DROP INDEX IF EXISTS idx_user_location_user_branch_id;
CREATE UNIQUE INDEX IF NOT EXISTS user_location_name_key ON user_location (user_id, branch_id, lower(name)) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_location_branch_id ON user_location(branch_id);

-- 3. Remove user_branch.
ALTER TABLE user_location DROP COLUMN user_branch_id;
DROP TABLE user_branch;

COMMIT;
