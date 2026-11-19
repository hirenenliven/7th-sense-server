-- Migration: machine accuracies (server ADR 0017)
--
-- The Machine Accuracy list of the machine form (machines.accuracy_id) and the part process form
-- (part_processes.accuracy_id). Until now a fixed list in the client (1 Good, 2 Average, 3 Poor);
-- now one table shared by every customer, written by the superadmin only (Account Center ›
-- Machine Accuracy). The seed keeps those ids, so saved machines and processes keep their accuracy.
-- Soft-deleted. Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/037_add_machine_accuracies.sql
--
-- Rollback:
--   DROP TABLE IF EXISTS machine_accuracies;

BEGIN;

CREATE TABLE IF NOT EXISTS machine_accuracies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- One accuracy per name (case-insensitive) among those not deleted
CREATE UNIQUE INDEX IF NOT EXISTS machine_accuracies_name_key
    ON machine_accuracies (LOWER(name)) WHERE is_deleted = false;

-- The client's old fixed list, with the ids machines and part processes already store
INSERT INTO machine_accuracies (id, name)
VALUES (1, 'Good'), (2, 'Average'), (3, 'Poor')
ON CONFLICT (id) DO NOTHING;

-- New rows continue after the seeded ids
SELECT setval(pg_get_serial_sequence('machine_accuracies', 'id'),
    GREATEST((SELECT MAX(id) FROM machine_accuracies), 1));

COMMIT;
