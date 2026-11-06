-- Migration: why a machine was down or idle (Productivity › Machine Mode Timeline popup)
--
-- A user clicks a finished Downtime or Idle bar of the Gantt chart and records one or more reasons
-- for parts of that period. Replaces the old downtime_reason / idle_reason tables. One table for
-- both: exactly one of the two reason columns is set, by the mode of the period (server ADR 0012).
-- The rows are what users entered, so they are soft-deleted, unlike machine_mode_log.
-- Run it once:
--   psql -U postgres -d 7th_sense -f migrations/033_add_machine_mode_reasons.sql
--
-- Rollback:
--   DROP TABLE IF EXISTS machine_mode_reasons;

BEGIN;

CREATE TABLE IF NOT EXISTS machine_mode_reasons (
    id BIGSERIAL PRIMARY KEY,
    mode_log_id BIGINT NOT NULL REFERENCES machine_mode_log(id) ON DELETE CASCADE,
    downtime_reason_id INTEGER REFERENCES downtime_reasons(id),
    idle_reason_id INTEGER REFERENCES idle_reasons(id),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ended_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT machine_mode_reasons_one_reason_check
        CHECK (num_nonnulls(downtime_reason_id, idle_reason_id) = 1),
    CONSTRAINT machine_mode_reasons_period_check CHECK (ended_at > started_at)
);

CREATE INDEX IF NOT EXISTS idx_machine_mode_reasons_mode_log
    ON machine_mode_reasons (mode_log_id, started_at) WHERE is_deleted = false;

DROP TRIGGER IF EXISTS update_machine_mode_reasons_updated_at ON machine_mode_reasons;
CREATE TRIGGER update_machine_mode_reasons_updated_at
    BEFORE UPDATE ON machine_mode_reasons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
