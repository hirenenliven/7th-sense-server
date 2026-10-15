-- Migration: licence keys
--
-- licence_keys: the licence keys the superadmin generates for a customer (user_id), e.g.
-- "K7M2Q-XR9TB-4HWZP-N8CDF" (4 groups of 5 characters, no I, O, 0 or 1). Keys are generated
-- by the server (POST /licence-keys/generate) and saved in a batch (POST /licence-keys).
-- A key is unique among the keys not deleted.
--
-- Rollback:
--   DROP TABLE IF EXISTS licence_keys;

BEGIN;

CREATE TABLE IF NOT EXISTS licence_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    licence_key VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS licence_keys_licence_key_key
    ON licence_keys (licence_key) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_licence_keys_user_id ON licence_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_licence_keys_is_deleted ON licence_keys(is_deleted);

CREATE TRIGGER update_licence_keys_updated_at BEFORE UPDATE ON licence_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
