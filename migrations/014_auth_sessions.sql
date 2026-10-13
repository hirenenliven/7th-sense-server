-- Migration: Login sessions in auth_tokens
-- Version: 014
-- Description: Login stores one auth_tokens row per session (token_type 'session').
-- The token column holds the SHA-256 hash of the cookie value, never the value itself.
-- Logout soft-deletes the row, so auth_tokens gets is_deleted / deleted_at like every table.

BEGIN;

ALTER TABLE auth_tokens ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE auth_tokens ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Session lookup: by token among live rows.
CREATE INDEX IF NOT EXISTS idx_auth_tokens_live_token ON auth_tokens(token) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_auth_tokens_is_deleted ON auth_tokens(is_deleted);

COMMIT;
