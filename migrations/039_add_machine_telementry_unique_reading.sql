-- Migration: one machine_telementry row per machine, variable and reading time (server ADR 0024)
--
-- The telemetry queue (RabbitMQ) delivers at least once: a batch the writer stored but could not
-- acknowledge is delivered again. A unique index on (machine_id, variable_id, recorded_at) lets
-- InsertTelemetry skip such rows (ON CONFLICT DO NOTHING). On a TimescaleDB hypertable a unique
-- index must hold the time column; this one does, and it is created on every chunk.
--
-- It also serves every query of migration 038's (machine_id, variable_id, recorded_at DESC) index
-- (PostgreSQL scans a b-tree both ways), so that index is dropped: one index less to write per row.
--
-- Rows stored twice before (same machine, variable and time) are removed first, keeping one.
-- On a large table this takes a while and blocks the collector's inserts: run it while collection
-- can pause (e.g. stop the server). Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/039_add_machine_telementry_unique_reading.sql
--
-- Rollback:
--   DROP INDEX IF EXISTS machine_telementry_machine_variable_reading_key;
--   CREATE INDEX IF NOT EXISTS machine_telementry_machine_variable_time_idx
--       ON machine_telementry (machine_id, variable_id, recorded_at DESC);

BEGIN;

-- machine_telementry has no id: a row is told apart by its chunk (tableoid) and position (ctid).
DELETE FROM machine_telementry stored
USING (
    SELECT table_id, row_id
    FROM (
        SELECT tableoid AS table_id, ctid AS row_id,
               row_number() OVER (PARTITION BY machine_id, variable_id, recorded_at) AS copy_number
        FROM machine_telementry
    ) numbered
    WHERE copy_number > 1
) duplicate
WHERE stored.tableoid = duplicate.table_id AND stored.ctid = duplicate.row_id;

CREATE UNIQUE INDEX IF NOT EXISTS machine_telementry_machine_variable_reading_key
    ON machine_telementry (machine_id, variable_id, recorded_at);

DROP INDEX IF EXISTS machine_telementry_machine_variable_time_idx;

COMMIT;
