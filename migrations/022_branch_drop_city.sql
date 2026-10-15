-- Migration: branches keep their city only as user_city_id — the free-text city column goes
--
-- The branch form picks the city from the customer's cities (user_city); its name is read
-- through the join (cityName), so the copied text in branches.city is no longer written or read.
--
-- Rollback:
--   ALTER TABLE branches ADD COLUMN city VARCHAR(100);

BEGIN;

ALTER TABLE branches DROP COLUMN IF EXISTS city;

COMMIT;
