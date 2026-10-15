-- Migration: licence key subscription duration and end date
--
-- duration: the subscription plan of the key (dto.LicenceKeyDurations: M1 … M12, Y1 … Y3,
--   LIFETIME, TRIAL, TRIAL_1Y). expires_on: the last day the licence is valid, computed by the
--   server from the day the key is saved (NULL = never ends, LIFETIME). A licence is expired
--   once expires_on < CURRENT_DATE. Keys saved before this migration become LIFETIME.
--
-- Rollback:
--   ALTER TABLE licence_keys DROP COLUMN expires_on, DROP COLUMN duration;

BEGIN;

ALTER TABLE licence_keys ADD COLUMN IF NOT EXISTS duration VARCHAR(20) NOT NULL DEFAULT 'LIFETIME';
ALTER TABLE licence_keys ALTER COLUMN duration DROP DEFAULT;
ALTER TABLE licence_keys ADD COLUMN IF NOT EXISTS expires_on DATE;

CREATE INDEX IF NOT EXISTS idx_licence_keys_expires_on ON licence_keys(expires_on);

COMMIT;
