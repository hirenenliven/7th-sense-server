-- Migration: a machine's place in its customer's organization (server ADR 0025)
--
-- A customer picks the organization levels its screens use (user_profiles.filters: country, state,
-- city, branch, location, shop, line, cell, supervisor). The Add / Edit Machine form shows exactly
-- those levels, and the machine is saved at the deepest one picked: one of the columns below (or
-- line_id, which machines already had) holds it, the others stay NULL. The levels above it are not
-- stored: they come from the organization itself (machine_organization_paths), so a level moved to
-- another parent takes its machines along and nothing goes stale.
-- Run it once:
--   psql -U postgres -d 7th_sense -f migrations/040_add_machine_organization.sql

ALTER TABLE machines
    ADD COLUMN IF NOT EXISTS user_country_id integer REFERENCES user_country(id),
    ADD COLUMN IF NOT EXISTS user_state_id integer REFERENCES user_state(id),
    ADD COLUMN IF NOT EXISTS user_city_id integer REFERENCES user_city(id),
    ADD COLUMN IF NOT EXISTS branch_id integer REFERENCES branches(id),
    ADD COLUMN IF NOT EXISTS user_location_id integer REFERENCES user_location(id),
    ADD COLUMN IF NOT EXISTS user_shop_id integer REFERENCES user_shop(id),
    ADD COLUMN IF NOT EXISTS user_cell_id integer REFERENCES user_cell(id),
    ADD COLUMN IF NOT EXISTS user_supervisor_id integer REFERENCES user_supervisor(id);

-- One place per machine. Machines saved before this migration only have line_id, so they pass.
ALTER TABLE machines DROP CONSTRAINT IF EXISTS machines_one_organization_level;
ALTER TABLE machines ADD CONSTRAINT machines_one_organization_level CHECK (
    num_nonnulls(user_country_id, user_state_id, user_city_id, branch_id, user_location_id,
                 user_shop_id, line_id, user_cell_id, user_supervisor_id) <= 1);

-- The delete checks of every level look machines up by these columns.
CREATE INDEX IF NOT EXISTS idx_machines_user_country_id ON machines (user_country_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_state_id ON machines (user_state_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_city_id ON machines (user_city_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_branch_id ON machines (branch_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_location_id ON machines (user_location_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_shop_id ON machines (user_shop_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_line_id ON machines (line_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_cell_id ON machines (user_cell_id) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_machines_user_supervisor_id ON machines (user_supervisor_id) WHERE is_deleted = false;

-- Every machine with its whole path, walked up from the level it is saved at: each level is the
-- machine's own column, else the parent of the level below. Levels below the saved one are NULL.
-- Soft-deleted levels still give their name (a level with machines cannot be deleted anyway).
CREATE OR REPLACE VIEW machine_organization_paths AS
SELECT m.id AS machine_id,
       co.id AS user_country_id, co.name AS user_country_name,
       st.id AS user_state_id, st.name AS user_state_name,
       ci.id AS user_city_id, ci.name AS user_city_name,
       br.id AS branch_id, br.name AS branch_name,
       lo.id AS user_location_id, lo.name AS user_location_name,
       sh.id AS user_shop_id, sh.name AS user_shop_name,
       li.id AS user_line_id, li.name AS user_line_name,
       ce.id AS user_cell_id, ce.name AS user_cell_name,
       su.id AS user_supervisor_id, su.name AS user_supervisor_name
FROM machines m
LEFT JOIN user_supervisor su ON su.id = m.user_supervisor_id
LEFT JOIN user_cell ce ON ce.id = COALESCE(m.user_cell_id, su.user_cell_id)
LEFT JOIN user_line li ON li.id = COALESCE(m.line_id, ce.user_line_id)
LEFT JOIN user_shop sh ON sh.id = COALESCE(m.user_shop_id, li.user_shop_id)
LEFT JOIN user_location lo ON lo.id = COALESCE(m.user_location_id, sh.user_location_id)
LEFT JOIN branches br ON br.id = COALESCE(m.branch_id, lo.branch_id)
LEFT JOIN user_city ci ON ci.id = COALESCE(m.user_city_id, br.user_city_id)
LEFT JOIN user_state st ON st.id = COALESCE(m.user_state_id, ci.user_state_id)
LEFT JOIN user_country co ON co.id = COALESCE(m.user_country_id, st.user_country_id);
