-- Migration: machine cycles (replaces the old mm_cycletime / mm_cycletime_history)
--
-- A cycle is one run of the part program on a machine, one row per cycle, written by the CNC
-- collector (internal/monitor) from the same DLL readings that decide the machine mode; the
-- collector code comes in a later change. Run it once, on any database (it creates only what
-- is missing):
--   psql -U postgres -d 7th_sense -f migrations/030_add_machine_cycles.sql
--
-- Old system -> this table
--   mm_cycletime.iCycleId            id
--   iMachineId, vProgramName         machine_id, program_name
--   iItemId (the schedule job)       schedule_item_id (+ part_id, process_id, shift_id)
--   iInProduction / iInHold          production_sec / hold_sec
--   iInLdUld                         loading_sec (gap since the previous cycle ended)
--   iInRapid                         not kept: Rapid Mode is no longer decided
--   iPartCount                       part_count_start / part_count_end / parts_made
--   iCstatus 0 running, 1 unfinished, 2 ended
--                                    status running / unfinished / completed
--   dCreatedDate / dCendTime         started_at / ended_at
--   mm_cycletime_history             not needed: one table, indexed by machine and start
--   mm_tool_life, batch_tbl          not ported yet (the DLL reads no tool data)
--
-- Signals in the DLL snapshot (program state from ProcStatus, machine_mode.go: running / held
-- / reset; Fanuc 3 START = running, Siemens 2 Active = running):
--   ProgramPartCount   steps up when the control counts a part (M30 or a counting M code).
--                      Not every program counts: machine 195 ran an M99 loop for days with
--                      the counter fixed at 262
--   CurProgLine        drops from the end of the program back to its start while running
--                      (M99 loop, a new part without leaving AUTO); machine 195: line 15 -> 3
--                      every 26.9 s
--   ProcStatus         reset -> running opens a cycle; running / held -> reset ends one
--   ProgName           a different program closes the cycle
--   The first snapshot after a connect is all zeros (ProcStatus, ProgramPartCount,
--   TotalProgLines 0, ProgName empty) and the real values follow within a second: such a
--   reading is skipped for cycles, so a reconnect is never a reset or a counter drop.
--   A counter that goes down (reset on the control, or 262 -> 0 -> 262 around that zero
--   snapshot) only moves the baseline; it never ends a cycle and never counts parts.
--
-- How a cycle opens and closes (first match wins; end_reason in brackets):
--   reset -> running                     a new cycle opens (status running)
--   running <-> held                     same cycle; held time counts as hold
--   ProgramPartCount up                  completed (part_count); if the program is still
--                                        running the next cycle opens at the same moment
--   CurProgLine back to the start while  completed (program_loop), next cycle opens at once;
--     running, counter not moved         the counter, when it counts, wins over the loop
--   running -> reset, no alarm           completed (program_end): the program reached M30 / M02
--   alarm or E-stop active -> reset      unfinished (alarm_reset)
--   held -> reset, no alarm, no part     unfinished (operator_reset): the operator aborted it
--   another ProgName                     unfinished (program_changed), a new cycle opens
--   open longer than the collector's     unfinished (timeout), so a machine left in Hold for
--     maximum cycle length               days does not distort the reports
--   Offline / server restart             the cycle stays open. The first connected reading
--                                        after it: running or held with the same program ->
--                                        same cycle, the gap counts as offline_sec; reset ->
--                                        unfinished (offline_reset, collector_stop when the
--                                        server stopped) and ended_at = when the machine went
--                                        offline or the last reading before the stop; that
--                                        last gap is not counted, so a night switched off is
--                                        not a 14-hour cycle
--
-- Writes: one INSERT when the cycle opens; the worker adds up the seconds per mode in memory
-- and writes them every 5 minutes (checkpoint) and once when the cycle closes. After a server
-- restart the seconds of a cycle that is still open are rebuilt from machine_mode_log between
-- started_at and the restart. Inside a cycle the mode is In Production, Hold, Downtime or
-- Offline; Idle (program reset) is always between cycles and becomes the next loading_sec.
-- Cycle length = ended_at - started_at (now - started_at while running), not stored.
--
-- Data table like machine_mode_log: written only by the collector, so no soft delete
-- (is_deleted / deleted_at) and no API writes.
--
-- Rollback:
--   DROP TABLE IF EXISTS machine_cycles;

BEGIN;

CREATE TABLE IF NOT EXISTS machine_cycles (
    id BIGSERIAL PRIMARY KEY,
    machine_id INTEGER NOT NULL REFERENCES machines(id) ON DELETE CASCADE,

    -- What ran. program_name is the DLL ProgName at cycle start (same as machine_mode_log);
    -- the job is matched at cycle start: the machine's production_schedule_items row whose
    -- day, shift and time window hold started_at, preferring the one whose part process has
    -- this program_name. NULL when nothing was planned or nothing matched.
    program_name VARCHAR(64),
    schedule_item_id INTEGER REFERENCES production_schedule_items(id),
    part_id INTEGER REFERENCES parts(id),
    process_id INTEGER REFERENCES part_processes(id),
    shift_id INTEGER REFERENCES shifts(id),
    -- Day the shift started (a 22:00-06:00 shift run at 02:00 belongs to the previous day);
    -- the calendar day of started_at when no shift covers it. A cycle is never split at a
    -- shift change: it belongs to the shift it started in.
    shift_date DATE,

    -- Planned cycle time (part_processes.cycle_time + its sub program cycle times, seconds)
    -- copied at cycle start, so later edits to the part do not change old cycles
    standard_cycle_sec INTEGER,

    status VARCHAR(20) NOT NULL DEFAULT 'running',
    -- Why the cycle closed (see the header); NULL while running
    end_reason VARCHAR(30),

    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,

    -- Seconds per mode inside the cycle; they add up to ended_at - started_at once closed
    production_sec INTEGER NOT NULL DEFAULT 0,
    hold_sec INTEGER NOT NULL DEFAULT 0,
    downtime_sec INTEGER NOT NULL DEFAULT 0,
    offline_sec INTEGER NOT NULL DEFAULT 0,
    -- Loading / unloading: from the end of the machine's previous cycle to this start
    -- (NULL for the first cycle, or when the gap holds Offline or a shift change)
    loading_sec INTEGER,

    -- DLL ProgramPartCount (cnc_variables 201) at start and end; parts_made = the steps up in
    -- between (a counter reset on the control is not subtracted), 0 when it did not count
    part_count_start INTEGER,
    part_count_end INTEGER,
    parts_made INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT machine_cycles_status_check
        CHECK (status IN ('running', 'completed', 'unfinished')),
    -- completed: part_count, program_loop, program_end; unfinished: the others
    CONSTRAINT machine_cycles_end_reason_check
        CHECK (end_reason IS NULL
               OR (status = 'completed' AND end_reason IN ('part_count', 'program_loop', 'program_end'))
               OR (status = 'unfinished' AND end_reason IN ('alarm_reset', 'operator_reset', 'program_changed',
                                                            'offline_reset', 'collector_stop', 'timeout'))),
    -- A running cycle has no end; a closed one has its end and reason
    CONSTRAINT machine_cycles_open_check
        CHECK ((status = 'running') = (ended_at IS NULL)
               AND (status = 'running') = (end_reason IS NULL)),
    CONSTRAINT machine_cycles_period_check CHECK (ended_at IS NULL OR ended_at >= started_at),
    CONSTRAINT machine_cycles_seconds_check
        CHECK (production_sec >= 0 AND hold_sec >= 0 AND downtime_sec >= 0 AND offline_sec >= 0
               AND (loading_sec IS NULL OR loading_sec >= 0)
               AND (standard_cycle_sec IS NULL OR standard_cycle_sec >= 0)),
    CONSTRAINT machine_cycles_parts_made_check CHECK (parts_made >= 0)
);

-- Upgrade a table made by the first version of this file: total_sec was derivable, the
-- end_reason list lacked program_loop / operator_reset / timeout, and the shift_date index
-- had no query yet
ALTER TABLE machine_cycles DROP COLUMN IF EXISTS total_sec;
DROP INDEX IF EXISTS idx_machine_cycles_shift_date;
ALTER TABLE machine_cycles DROP CONSTRAINT IF EXISTS machine_cycles_seconds_check,
    ADD CONSTRAINT machine_cycles_seconds_check
        CHECK (production_sec >= 0 AND hold_sec >= 0 AND downtime_sec >= 0 AND offline_sec >= 0
               AND (loading_sec IS NULL OR loading_sec >= 0)
               AND (standard_cycle_sec IS NULL OR standard_cycle_sec >= 0));
ALTER TABLE machine_cycles DROP CONSTRAINT IF EXISTS machine_cycles_end_reason_check,
    ADD CONSTRAINT machine_cycles_end_reason_check
        CHECK (end_reason IS NULL
               OR (status = 'completed' AND end_reason IN ('part_count', 'program_loop', 'program_end'))
               OR (status = 'unfinished' AND end_reason IN ('alarm_reset', 'operator_reset', 'program_changed',
                                                            'offline_reset', 'collector_stop', 'timeout')));

-- At most one running cycle per machine (like machine_mode_log_open_key)
CREATE UNIQUE INDEX IF NOT EXISTS machine_cycles_open_key
    ON machine_cycles (machine_id) WHERE status = 'running';

-- Cycle history of a machine, newest first (chart page, reports by date range)
CREATE INDEX IF NOT EXISTS idx_machine_cycles_machine_started
    ON machine_cycles (machine_id, started_at DESC);

-- Production per job
CREATE INDEX IF NOT EXISTS idx_machine_cycles_schedule_item
    ON machine_cycles (schedule_item_id) WHERE schedule_item_id IS NOT NULL;

COMMIT;
