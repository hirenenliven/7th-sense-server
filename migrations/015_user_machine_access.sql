-- Migration: user_machine_access, the machines each user may access
--
-- One row per user-machine assignment. Users with the supervisor, operator, qc or
-- shop_head role see only the machines assigned here; a customer and its admins see
-- every machine of the customer without rows here (see service.MachineService).
--
-- Unlike most tables there is no is_deleted: an assignment is either there or not, so
-- removing one deletes the row and assigning it again inserts a new one.
-- Future per-machine permissions can hang off this row (user_machine_access.id).
--
-- The older user_machines table (migration 004) was never used and is left as it is.
--
-- Rollback: DROP TABLE IF EXISTS user_machine_access;

BEGIN;

CREATE TABLE IF NOT EXISTS user_machine_access (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_machine_access_user_machine_key UNIQUE (user_id, machine_id)
);

-- user_id lookups use the unique index above; this one serves "who can access machine X"
CREATE INDEX IF NOT EXISTS idx_user_machine_access_machine_id ON user_machine_access(machine_id);

CREATE TRIGGER update_user_machine_access_updated_at BEFORE UPDATE ON user_machine_access
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
