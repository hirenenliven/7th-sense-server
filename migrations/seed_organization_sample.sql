-- Sample organization data for testing (dev only). Run after migration 011.
-- India → Gujarat → Rajkot → Main Branch → U1 → Shop1 → L1, L2 → Cell 1 → Supervisor A
-- Owner: user 6 (the current dev user, see internal/middleware/auth.go). Safe to re-run:
-- rows that already exist (same name under the same parent) are skipped.

BEGIN;

INSERT INTO user_country (user_id, name)
SELECT 6, 'India'
WHERE NOT EXISTS (SELECT 1 FROM user_country WHERE user_id = 6 AND lower(name) = 'india' AND is_deleted = false);

INSERT INTO user_state (user_id, user_country_id, name)
SELECT 6, c.id, 'Gujarat' FROM user_country c
WHERE c.user_id = 6 AND c.name = 'India' AND c.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_state s WHERE s.user_country_id = c.id AND s.name = 'Gujarat' AND s.is_deleted = false);

INSERT INTO user_city (user_id, user_state_id, name)
SELECT 6, s.id, 'Rajkot' FROM user_state s
WHERE s.user_id = 6 AND s.name = 'Gujarat' AND s.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_city x WHERE x.user_state_id = s.id AND x.name = 'Rajkot' AND x.is_deleted = false);

INSERT INTO branches (code, name, customer_id, user_city_id, status)
SELECT 'MAINRJK', 'Main Branch', 6, p.id, 'active' FROM user_city p
WHERE p.user_id = 6 AND p.name = 'Rajkot' AND p.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM branches x WHERE x.code = 'MAINRJK');

INSERT INTO user_location (user_id, branch_id, name)
SELECT 6, p.id, 'U1' FROM branches p
WHERE p.code = 'MAINRJK' AND p.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_location x WHERE x.branch_id = p.id AND x.name = 'U1' AND x.is_deleted = false);

INSERT INTO user_shop (user_id, user_location_id, name)
SELECT 6, p.id, 'Shop1' FROM user_location p
WHERE p.user_id = 6 AND p.name = 'U1' AND p.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_shop x WHERE x.user_location_id = p.id AND x.name = 'Shop1' AND x.is_deleted = false);

INSERT INTO user_line (user_id, user_shop_id, name)
SELECT 6, p.id, line_name FROM user_shop p, (VALUES ('L1'), ('L2')) AS lines(line_name)
WHERE p.user_id = 6 AND p.name = 'Shop1' AND p.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_line x WHERE x.user_shop_id = p.id AND x.name = line_name AND x.is_deleted = false);

INSERT INTO user_cell (user_id, user_line_id, name)
SELECT 6, p.id, 'Cell 1' FROM user_line p
WHERE p.user_id = 6 AND p.name = 'L1' AND p.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_cell x WHERE x.user_line_id = p.id AND x.name = 'Cell 1' AND x.is_deleted = false);

INSERT INTO user_supervisor (user_id, user_cell_id, name)
SELECT 6, p.id, 'Supervisor A' FROM user_cell p
WHERE p.user_id = 6 AND p.name = 'Cell 1' AND p.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_supervisor x WHERE x.user_cell_id = p.id AND x.name = 'Supervisor A' AND x.is_deleted = false);

COMMIT;
