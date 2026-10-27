-- Imports the old machine_monitor machines (active, not offline) into machines.
-- Set customer_id below to the customer that owns them. Brand, controller and model are looked
-- up by brand name (the first controller of the brand, its first model), so no ids are hard-coded.
-- Machines the customer already has (same name, case-insensitive) are skipped, so it can be re-run.
-- Old defaults kept: username JYOTI, password 12345678, version V04.07, hourly cost 800.00.
-- status 'active': the server's CNC collector connects every inserted machine right away.
BEGIN;

WITH params AS (
    SELECT 49 AS customer_id   -- <-- the customer (users.id, role customer) that owns the machines
),
old_machines (name, ip, port, brand_name) AS (
    VALUES
    ('HM11', '192.168.173.47', 8147, 'Fanuc'),
    ('HM17', '192.168.170.13', 8123, 'Fanuc'),
    ('TC29', '192.168.162.22', 8111, 'Fanuc'),
    ('TC30', '192.168.162.23', 8158, 'Fanuc'),
    ('TC31', '192.168.162.24', 8159, 'Fanuc'),
    ('TC33', '192.168.162.149', 8297, 'Fanuc'),
    ('TC34', '192.168.162.36', 8171, 'Fanuc'),
    ('TC35', '192.168.162.37', 8172, 'Fanuc'),
    ('TC36', '192.168.162.38', 8173, 'Fanuc'),
    ('TC37', '192.168.162.39', 8174, 'Fanuc'),
    ('TM12', '192.168.173.61', 8141, 'Fanuc'),
    ('UM28', '192.168.170.42', 8127, 'Fanuc'),
    ('VM40', '192.168.162.21', 8110, 'Fanuc'),
    ('ZV03', '192.168.170.17', 8134, 'Fanuc'),
    ('GC10', '192.168.8.112', 4804, 'Siemens'),
    ('GC21', '192.168.8.115', 4871, 'Siemens'),
    ('GC22', '192.168.8.113', 4870, 'Siemens'),
    ('GC24', '192.168.8.64', 4801, 'Siemens'),
    ('GC25', '192.168.8.114', 8103, 'Siemens'),
    ('GC28', '192.168.8.151', 8303, 'Siemens'),
    ('GC33', '192.168.8.106', 4809, 'Siemens'),
    ('GC34', '192.168.8.108', 8101, 'Siemens'),
    ('GS05', '192.168.8.13', 4876, 'Siemens'),
    ('GS09', '192.168.8.62', 8216, 'Siemens'),
    ('GS12', '192.168.8.109', 8196, 'Siemens'),
    ('RD05', '192.168.162.114', 8262, 'Siemens'),
    ('RD11', '192.168.162.108', 8256, 'Siemens'),
    ('RD12', '192.168.162.107', 8255, 'Siemens'),
    ('RD13', '192.168.162.121', 8269, 'Siemens'),
    ('RD17', '192.168.162.109', 8257, 'Siemens'),
    ('RD19', '192.168.162.119', 8267, 'Siemens'),
    ('RD20', '192.168.162.118', 8266, 'Siemens'),
    ('RD22', '192.168.162.116', 8264, 'Siemens'),
    ('TC07', '192.168.8.66', 4839, 'Siemens'),
    ('TC11', '192.168.8.90', 4881, 'Siemens'),
    ('TC22', '192.168.8.65', 4803, 'Siemens'),
    ('TC39', '192.168.8.69', 8271, 'Siemens'),
    ('TC41', '192.168.8.81', 8272, 'Siemens'),
    ('TM14', '192.168.8.70', 8294, 'Siemens'),
    ('TM16', '192.168.8.67', 4873, 'Siemens'),
    ('TM18', '192.168.8.159', 8160, 'Siemens'),
    ('TM19', '192.168.8.142', 8220, 'Siemens'),
    ('TM21', '192.168.8.143', 8192, 'Siemens'),
    ('TM22', '192.168.8.103', 4821, 'Siemens'),
    ('TM23', '192.168.8.92', 8275, 'Siemens'),
    ('TM26', '192.168.8.93', 8274, 'Siemens'),
    ('UM18', '192.168.172.191', 4891, 'Siemens'),
    ('UM19', '192.168.170.54', 4811, 'Siemens'),
    ('UM21', '192.168.8.32', 4812, 'Siemens'),
    ('UM24', '192.168.170.56', 8226, 'Siemens'),
    ('UM27', '192.168.8.63', 4822, 'Siemens'),
    ('VM14', '192.168.8.86', 4892, 'Siemens'),
    ('VM33', '192.168.173.56', 4857, 'Siemens'),
    ('VM36', '192.168.173.50', 4856, 'Siemens'),
    ('VM43', '192.168.8.14', 8215, 'Siemens'),
    ('VM45', '192.168.8.101', 8273, 'Siemens'),
    ('VM46', '192.168.8.75', 8276, 'Siemens'),
    ('VM47', '192.168.8.169', 8109, 'Siemens'),
    ('VM50', '192.168.8.154', 8162, 'Siemens'),
    ('VM51', '192.168.8.160', 8163, 'Siemens'),
    ('VM53', '192.168.8.175', 8165, 'Siemens'),
    ('VM55', '192.168.8.49', 8214, 'Siemens'),
    ('VM57', '192.168.8.149', 8167, 'Siemens'),
    ('VM58', '192.168.8.119', 8168, 'Siemens'),
    ('VM74', '192.168.8.16', 8217, 'Siemens'),
    ('ZV01', '192.168.8.170', 4861, 'Siemens'),
    ('ZV02', '192.168.8.87', 4862, 'Siemens'),
    ('ZV05', '192.168.8.118', 4863, 'Siemens'),
    ('ZV06', '192.168.8.132', 4864, 'Siemens'),
    ('ZV07', '192.168.8.97', 4865, 'Siemens'),
    ('ZV08', '192.168.8.55', 4867, 'Siemens'),
    ('ZV09', '192.168.8.52', 4866, 'Siemens'),
    ('ZV10', '192.168.8.98', 4846, 'Siemens'),
    ('ZV11', '192.168.8.47', 8212, 'Siemens'),
    ('ZV12', '192.168.8.94', 4848, 'Siemens'),
    ('ZV13', '192.168.8.45', 4851, 'Siemens'),
    ('ZV15', '192.168.8.88', 4860, 'Siemens'),
    ('ZV16', '192.168.8.99', 4845, 'Siemens'),
    ('ZV18', '192.168.8.71', 4838, 'Siemens'),
    ('ZV19', '192.168.8.85', 4853, 'Siemens')
),
catalog AS (
    -- The first controller of each brand and that controller's first model.
    SELECT DISTINCT ON (b.id)
           b.name AS brand_name, b.id AS brand_id, c.id AS controller_id, mo.id AS model_id
    FROM brands b
    JOIN controllers c ON c.brand_id = b.id AND c.is_deleted = false
    JOIN models mo     ON mo.controller_id = c.id AND mo.is_deleted = false
    WHERE b.is_deleted = false
    ORDER BY b.id, c.id, mo.id
)
INSERT INTO machines (
    user_id, name, serial_number, live_ip, local_ip, port, username, password,
    brand_id, controller_id, model_id, controller_version, controller_version_type,
    installation_date, hourly_cost, health_ids, machine_system_type, security_mode,
    status, mode_id, online_status
)
SELECT params.customer_id, old.name, old.name, old.ip, old.ip, old.port, 'JYOTI', '12345678',
       catalog.brand_id, catalog.controller_id, catalog.model_id, 'V04.07', 1,  -- 1 = 828D
       CURRENT_DATE, 800.00, 'All', 'NewPlc',
       CASE WHEN old.brand_name = 'Siemens' THEN 'Basic256Sha256' ELSE 'None' END,
       'active', 9, false                                                     -- 9 = Offline until it connects
FROM old_machines old
CROSS JOIN params
JOIN catalog ON lower(catalog.brand_name) = lower(old.brand_name)
WHERE NOT EXISTS (
    SELECT 1 FROM machines existing
    WHERE existing.user_id = params.customer_id
      AND lower(existing.name) = lower(old.name)
      AND existing.is_deleted = false
)
RETURNING id, name, live_ip, port, brand_id, controller_id, model_id;

COMMIT;
