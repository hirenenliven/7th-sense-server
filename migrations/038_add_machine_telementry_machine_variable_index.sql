-- Migration: index machine_telementry by machine, variable and time (server ADR 0018)
--
-- machine_telementry had only (recorded_at DESC). A chart of one machine's variables over a time
-- range (GET /machine-health/axis-values) then reads every machine's rows in that range and
-- throws most away; with many machines and crores of rows that grows with the whole fleet. This
-- index lets it read only that machine's rows of those variables, and makes "last value before the
-- range" one index lookup. On a hypertable the index is created on every chunk.
--
-- Creating it on a large table takes a while and blocks the collector's inserts meanwhile: run it
-- when collection can pause (e.g. CNC_ENABLED=false). Safe to run more than once:
--   psql -U postgres -d 7th_sense -f migrations/038_add_machine_telementry_machine_variable_index.sql
--
-- Rollback:
--   DROP INDEX IF EXISTS machine_telementry_machine_variable_time_idx;

CREATE INDEX IF NOT EXISTS machine_telementry_machine_variable_time_idx
    ON machine_telementry (machine_id, variable_id, recorded_at DESC);
