-- Migration: fewer writes per telemetry row, smaller indexes to keep in memory (server ADR 0026)
--
-- 1. Drops machine_telementry_recorded_at_idx, the (recorded_at DESC) index TimescaleDB adds to a
--    new hypertable. Every query of the table names machine_id and variable_id (MachineHealthRepository),
--    which the unique (machine_id, variable_id, recorded_at) index of migration 039 serves, so this
--    index was only written, never needed: one index less per inserted row.
-- 2. New chunks hold 1 day instead of 7: the indexes of the chunk being written stay small enough
--    to stay in memory, so inserts do not read index pages from disk. Existing chunks keep 7 days.
-- 3. machine_variable_values pages keep 20 % free (fillfactor 80): its rows are updated all the
--    time and only value / updated_at change (not indexed), so with room on the page PostgreSQL
--    updates in place (HOT) without touching the primary key index.
--
-- Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/041_tune_telemetry_writes.sql
--
-- Rollback:
--   CREATE INDEX IF NOT EXISTS machine_telementry_recorded_at_idx ON machine_telementry (recorded_at DESC);
--   SELECT set_chunk_time_interval('machine_telementry', INTERVAL '7 days');
--   ALTER TABLE machine_variable_values RESET (fillfactor);

BEGIN;

DROP INDEX IF EXISTS machine_telementry_recorded_at_idx;

-- Only on a hypertable: development databases may keep machine_telementry a plain table.
DO $$
BEGIN
    -- Nested: timescaledb_information does not exist without the extension.
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb') THEN
        IF EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'machine_telementry') THEN
            PERFORM set_chunk_time_interval('machine_telementry', INTERVAL '1 day');
        END IF;
    END IF;
END
$$;

ALTER TABLE machine_variable_values SET (fillfactor = 80);

COMMIT;

-- Rewrites the existing pages with the free room (fillfactor applies to new pages only). The table
-- holds one row per machine-variable, so this is quick; it locks the table meanwhile.
VACUUM FULL machine_variable_values;
