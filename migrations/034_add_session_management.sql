-- Migration: server-side session management (server ADR 0014)
--
-- Every login session stays a row of auth_tokens (token = SHA-256 of the cookie's token). This adds
-- what a user or the superadmin needs to see and end sessions: a status, the client IP and who
-- revoked a session and why. Logout and revoke still set is_deleted / deleted_at as well, so every
-- query that filters on is_deleted keeps working. Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/034_add_session_management.sql
--
-- Rollback:
--   DROP INDEX IF EXISTS idx_auth_tokens_user_status;
--   ALTER TABLE auth_tokens DROP CONSTRAINT IF EXISTS auth_tokens_status_check,
--     DROP COLUMN IF EXISTS status, DROP COLUMN IF EXISTS ip_address, DROP COLUMN IF EXISTS revoked_at,
--     DROP COLUMN IF EXISTS revoked_by, DROP COLUMN IF EXISTS revoke_reason;

BEGIN;

ALTER TABLE auth_tokens
    ADD COLUMN IF NOT EXISTS status varchar(20) NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS ip_address varchar(45),
    ADD COLUMN IF NOT EXISTS revoked_at timestamptz,
    ADD COLUMN IF NOT EXISTS revoked_by integer REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS revoke_reason varchar(50);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'auth_tokens_status_check') THEN
        ALTER TABLE auth_tokens ADD CONSTRAINT auth_tokens_status_check
            CHECK (status IN ('active', 'logged_out', 'revoked', 'expired'));
    END IF;
END $$;

-- Rows from before this migration: a deleted session was logged out, an overdue one expired.
UPDATE auth_tokens SET status = 'logged_out'
WHERE status = 'active' AND is_deleted = true;

UPDATE auth_tokens SET status = 'expired'
WHERE status = 'active' AND expires_at <= NOW();

-- Counting and listing a user's active sessions (login conflict check, My sessions).
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_status
    ON auth_tokens (user_id, status) WHERE token_type = 'session';

COMMIT;
