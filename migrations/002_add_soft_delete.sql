-- Soft delete for users and addresses.
-- Deleted rows stay in the table with is_deleted = true and deleted_at set;
-- the API only reads rows where is_deleted = false.

BEGIN;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

ALTER TABLE addresses
    ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- A deleted user must not block a new user from taking the same email or
-- username, so uniqueness now applies only to rows that are not deleted.
-- The index names are kept the same as the old constraints (the API relies on them).
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;
CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email) WHERE is_deleted = false;
CREATE UNIQUE INDEX IF NOT EXISTS users_username_key ON users (username) WHERE is_deleted = false;

COMMIT;
