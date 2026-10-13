-- Migration: Dark theme by default
-- Version: 012
-- Description: The UI is dark by design and now also ships a light theme, chosen per user
-- (user_profiles.theme: 'dark_version' | 'light_version'). New profiles default to dark.
-- Existing rows keep their stored value.
-- Rollback: ALTER TABLE user_profiles ALTER COLUMN theme SET DEFAULT 'light_version';

BEGIN;

ALTER TABLE user_profiles ALTER COLUMN theme SET DEFAULT 'dark_version';

COMMIT;
