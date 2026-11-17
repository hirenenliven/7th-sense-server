-- Migration: operator tasks (server ADR 0015)
--
-- The checks an operator does on a machine (clean filters, look for leaks, …). Two tables, like
-- the old maintenance_tasks / mm_operator_tasks:
--   maintenance_tasks — the master list of task descriptions, shared by every customer; seeded with
--                       the old list; the superadmin adds more from the Add Operator Task form.
--   operator_tasks    — which tasks apply to a customer's machines of a model (Machine Type); no
--                       model = every machine of the customer. Only the superadmin writes them.
-- Both are soft-deleted. Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/036_add_operator_tasks.sql
--
-- Rollback:
--   DROP TABLE IF EXISTS operator_tasks; DROP TABLE IF EXISTS maintenance_tasks;

BEGIN;

CREATE TABLE IF NOT EXISTS maintenance_tasks (
    id SERIAL PRIMARY KEY,
    description VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT maintenance_tasks_status_check CHECK (status IN ('active', 'inactive'))
);

-- One task per description (case-insensitive) among those not deleted
CREATE UNIQUE INDEX IF NOT EXISTS maintenance_tasks_description_key
    ON maintenance_tasks (LOWER(description)) WHERE is_deleted = false;

CREATE TABLE IF NOT EXISTS operator_tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),               -- the customer
    model_id INTEGER REFERENCES models(id),                      -- Machine Type; NULL = every machine
    maintenance_task_id INTEGER NOT NULL REFERENCES maintenance_tasks(id),
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- A task once per customer and machine type (no model counts as one type) among those not deleted
CREATE UNIQUE INDEX IF NOT EXISTS operator_tasks_customer_model_task_key
    ON operator_tasks (user_id, COALESCE(model_id, 0), maintenance_task_id) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_operator_tasks_user ON operator_tasks (user_id) WHERE is_deleted = false;

-- The old system's task list (maintenance_tasks). Tasks 10 and 11 were cut short in the old
-- screens; their endings here are a best reading.
INSERT INTO maintenance_tasks (description)
SELECT seed.description
FROM (VALUES
    ('Clean Panel AC Filter'),
    ('Check for any Oil leakage'),
    ('Maintain Chilling medium tank levelon chiller'),
    ('Check for any leakage in chiller lines'),
    ('Check for any Hyd. Oil leakage'),
    ('Check input air pressure'),
    ('Check CTS filter pressure gauge'),
    ('Check Axis and Spindle lubrication tank level'),
    ('Clean machine main table and internal covers'),
    ('Check all external covers, doors, windows are properly closed'),
    ('Check spindle and axes noise in moving/running condition'),
    ('Check movement and noise in chip conveyor chain and gear box'),
    ('Check no interlock by passed'),
    ('Check for any kind of warning message'),
    ('Clean OP, KB and MCP')
) AS seed (description)
WHERE NOT EXISTS (
    SELECT 1 FROM maintenance_tasks existing
    WHERE LOWER(existing.description) = LOWER(seed.description) AND existing.is_deleted = false
);

COMMIT;
