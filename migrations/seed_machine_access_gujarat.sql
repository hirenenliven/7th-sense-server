-- Sample data for the user form's machine access (dev only, ADR 0010). Run after migration 032:
--   psql -U postgres -d 7th_sense -f migrations/seed_machine_access_gujarat.sql
--
-- Owner: customer cust_15 (users.id 49, "Bogan - Nikolaus"), whose organization already has India.
--   India → Gujarat → Rajkot    → Rajkot Main (GJRJK01)    → Rajkot Unit 1    → Machining Shop → GJ Line A: GJ-VMC-01, GJ-VMC-02
--                                                                                               → GJ Line B: GJ-HMC-01
--                   → Ahmedabad → Ahmedabad Main (GJAMD01) → Ahmedabad Unit 1 → Turning Shop   → GJ Line C: GJ-LATHE-01
-- The customer's machine Poway (Andhra Pradesh) stays outside Gujarat.
-- Machines are added inactive, so the CNC collector does not try to connect to their sample IPs.
-- Safe to re-run: rows that already exist (same name under the same parent) are skipped.

BEGIN;

INSERT INTO user_country (user_id, name)
SELECT 49, 'India'
WHERE NOT EXISTS (SELECT 1 FROM user_country WHERE user_id = 49 AND name = 'India' AND is_deleted = false);

INSERT INTO user_state (user_id, user_country_id, name)
SELECT 49, c.id, 'Gujarat' FROM user_country c
WHERE c.user_id = 49 AND c.name = 'India' AND c.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_state s WHERE s.user_country_id = c.id AND s.name = 'Gujarat' AND s.is_deleted = false);

INSERT INTO user_city (user_id, user_state_id, name)
SELECT 49, s.id, city_name FROM user_state s, (VALUES ('Rajkot'), ('Ahmedabad')) AS cities(city_name)
WHERE s.user_id = 49 AND s.name = 'Gujarat' AND s.is_deleted = false
  AND NOT EXISTS (SELECT 1 FROM user_city x WHERE x.user_state_id = s.id AND x.name = city_name AND x.is_deleted = false);

INSERT INTO branches (code, name, customer_id, user_city_id, status)
SELECT b.code, b.name, 49, c.id, 'active'
FROM (VALUES ('GJRJK01', 'Rajkot Main', 'Rajkot'), ('GJAMD01', 'Ahmedabad Main', 'Ahmedabad')) AS b(code, name, city_name)
JOIN user_city c ON c.user_id = 49 AND c.name = b.city_name AND c.is_deleted = false
JOIN user_state s ON s.id = c.user_state_id AND s.name = 'Gujarat'
WHERE NOT EXISTS (SELECT 1 FROM branches x WHERE x.code = b.code AND x.is_deleted = false);

INSERT INTO user_location (user_id, branch_id, name)
SELECT 49, br.id, l.name
FROM (VALUES ('GJRJK01', 'Rajkot Unit 1'), ('GJAMD01', 'Ahmedabad Unit 1')) AS l(code, name)
JOIN branches br ON br.code = l.code AND br.is_deleted = false
WHERE NOT EXISTS (SELECT 1 FROM user_location x WHERE x.branch_id = br.id AND x.name = l.name AND x.is_deleted = false);

INSERT INTO user_shop (user_id, user_location_id, name)
SELECT 49, lo.id, s.name
FROM (VALUES ('Rajkot Unit 1', 'Machining Shop'), ('Ahmedabad Unit 1', 'Turning Shop')) AS s(location_name, name)
JOIN user_location lo ON lo.user_id = 49 AND lo.name = s.location_name AND lo.is_deleted = false
WHERE NOT EXISTS (SELECT 1 FROM user_shop x WHERE x.user_location_id = lo.id AND x.name = s.name AND x.is_deleted = false);

INSERT INTO user_line (user_id, user_shop_id, name)
SELECT 49, sh.id, l.name
FROM (VALUES ('Machining Shop', 'GJ Line A'), ('Machining Shop', 'GJ Line B'), ('Turning Shop', 'GJ Line C')) AS l(shop_name, name)
JOIN user_shop sh ON sh.user_id = 49 AND sh.name = l.shop_name AND sh.is_deleted = false
WHERE NOT EXISTS (SELECT 1 FROM user_line x WHERE x.user_shop_id = sh.id AND x.name = l.name AND x.is_deleted = false);

-- Brand / controller / model of the customer's machine Poway (63 / 151 / 41)
INSERT INTO machines (user_id, name, serial_number, live_ip, port, username, password, local_ip, brand_id,
    controller_id, model_id, controller_version, accuracy_id, line_id, installation_date, hourly_cost, status)
SELECT 49, m.name, m.serial, '10.0.0.' || m.host, 4840, 'operator', 'secret123', '192.168.10.' || m.host, 63,
    151, 41, 'V1.0', 1, li.id, DATE '2026-01-15', 750, 'inactive'
FROM (VALUES ('GJ-VMC-01', 'SN-GJ-VMC-01', 'GJ Line A', 11), ('GJ-VMC-02', 'SN-GJ-VMC-02', 'GJ Line A', 12),
             ('GJ-HMC-01', 'SN-GJ-HMC-01', 'GJ Line B', 13), ('GJ-LATHE-01', 'SN-GJ-LATHE-01', 'GJ Line C', 14))
     AS m(name, serial, line_name, host)
JOIN user_line li ON li.user_id = 49 AND li.name = m.line_name AND li.is_deleted = false
WHERE NOT EXISTS (SELECT 1 FROM machines x WHERE x.user_id = 49 AND x.name = m.name AND x.is_deleted = false);

COMMIT;

-- Check: the Gujarat machines and their path
SELECT m.id, m.name, p.user_country_id, p.user_state_id, p.user_city_id, p.branch_id,
       p.user_location_id, p.user_shop_id, p.user_line_id, li.name AS line
FROM machines m
JOIN user_line li ON li.id = m.line_id
JOIN organization_paths p ON p.user_line_id = m.line_id
WHERE m.user_id = 49 AND m.is_deleted = false
ORDER BY m.id;
