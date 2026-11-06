-- Migration: machine access by organization (user form › Access › Machines)
--
-- A supervisor, operator, QC or shop head gets machines either one by one (user_machine_access)
-- or through a filter on the company's organization: Country → State → City → Branch → Location
-- → Shop → Line. With all_machines on, the user may access every machine of the company whose
-- line lies under every level the filter names (no level = every company machine), including
-- machines added later. With it off the filter only remembers the selection of the form.
-- One row per user; like user_machine_access there is no soft delete (saving replaces the row).
-- Run it once:
--   psql -U postgres -d 7th_sense -f migrations/032_add_user_machine_filters.sql

CREATE TABLE IF NOT EXISTS user_machine_filters (
    id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id),
    user_country_id integer REFERENCES user_country(id),
    user_state_id integer REFERENCES user_state(id),
    user_city_id integer REFERENCES user_city(id),
    branch_id integer REFERENCES branches(id),
    user_location_id integer REFERENCES user_location(id),
    user_shop_id integer REFERENCES user_shop(id),
    user_line_id integer REFERENCES user_line(id),
    all_machines boolean NOT NULL DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS user_machine_filters_user_key ON user_machine_filters (user_id);

-- Every path of a company's organization, from each country down to each line. A level without
-- children ends its path with NULLs below it, so a chain of levels is valid when some row has
-- them all, and a machine's region is the row of its line (machines.line_id).
CREATE OR REPLACE VIEW organization_paths AS
SELECT co.user_id,
       co.id AS user_country_id,
       st.id AS user_state_id,
       ci.id AS user_city_id,
       br.id AS branch_id,
       lo.id AS user_location_id,
       sh.id AS user_shop_id,
       li.id AS user_line_id
FROM user_country co
LEFT JOIN user_state st ON st.user_country_id = co.id AND st.user_id = co.user_id AND st.is_deleted = false
LEFT JOIN user_city ci ON ci.user_state_id = st.id AND ci.user_id = co.user_id AND ci.is_deleted = false
LEFT JOIN branches br ON br.user_city_id = ci.id AND br.is_deleted = false
LEFT JOIN user_location lo ON lo.branch_id = br.id AND lo.user_id = co.user_id AND lo.is_deleted = false
LEFT JOIN user_shop sh ON sh.user_location_id = lo.id AND sh.user_id = co.user_id AND sh.is_deleted = false
LEFT JOIN user_line li ON li.user_shop_id = sh.id AND li.user_id = co.user_id AND li.is_deleted = false
WHERE co.is_deleted = false;
