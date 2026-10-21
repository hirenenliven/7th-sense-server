-- Migration: downtime reasons and idle reasons
--
-- Replace the old mst_downtime_reason / mst_idle_reason tables. Each reason belongs to a
-- customer (user_id) and is what an operator picks when a machine stops (downtime) or
-- stands idle. Reason code is optional; when given it is unique per customer, and the
-- reason text is unique per customer too (both case-insensitive).
--
-- Rollback: DROP TABLE IF EXISTS downtime_reasons, idle_reasons;

BEGIN;

CREATE TABLE IF NOT EXISTS downtime_reasons (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    reason_code VARCHAR(10),
    reason VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- A customer cannot use the same reason code, or the same reason, twice
CREATE UNIQUE INDEX IF NOT EXISTS downtime_reasons_user_reason_code_key
    ON downtime_reasons (user_id, LOWER(reason_code)) WHERE is_deleted = false AND reason_code IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS downtime_reasons_user_reason_key
    ON downtime_reasons (user_id, LOWER(reason)) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_downtime_reasons_user_id ON downtime_reasons(user_id);
CREATE INDEX IF NOT EXISTS idx_downtime_reasons_is_deleted ON downtime_reasons(is_deleted);

CREATE TRIGGER update_downtime_reasons_updated_at BEFORE UPDATE ON downtime_reasons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE IF NOT EXISTS idle_reasons (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    reason_code VARCHAR(10),
    reason VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idle_reasons_user_reason_code_key
    ON idle_reasons (user_id, LOWER(reason_code)) WHERE is_deleted = false AND reason_code IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idle_reasons_user_reason_key
    ON idle_reasons (user_id, LOWER(reason)) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_idle_reasons_user_id ON idle_reasons(user_id);
CREATE INDEX IF NOT EXISTS idx_idle_reasons_is_deleted ON idle_reasons(is_deleted);

CREATE TRIGGER update_idle_reasons_updated_at BEFORE UPDATE ON idle_reasons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
