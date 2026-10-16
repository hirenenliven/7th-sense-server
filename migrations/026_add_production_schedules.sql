-- Migration: production schedules and their items
--
-- Replaces the old production_schedule / schedule_items tables. A schedule belongs to a
-- customer (user_id) and plans one of its machines for one day (schedule_type 'single',
-- end_date = start_date) or a range of days ('multiple'). Each item is one job on that
-- machine: a part's process, the quantity to make, the shift, its supervisor / programmer
-- and operator (members of the customer) and the time it runs on schedule_date (a day of
-- the schedule). Times are "HH:MM"; an end not after the start runs past midnight.
-- total_minutes is worked out by the server from the times, never sent.
--
-- Rollback: DROP TABLE IF EXISTS production_schedule_items, production_schedules;

BEGIN;

CREATE TABLE IF NOT EXISTS production_schedules (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    machine_id INTEGER NOT NULL REFERENCES machines(id),
    schedule_type VARCHAR(10) NOT NULL CHECK (schedule_type IN ('single', 'multiple')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL CHECK (end_date >= start_date),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_production_schedules_user_id ON production_schedules(user_id);
CREATE INDEX IF NOT EXISTS idx_production_schedules_machine_dates
    ON production_schedules(machine_id, start_date, end_date) WHERE is_deleted = false;

CREATE TRIGGER update_production_schedules_updated_at BEFORE UPDATE ON production_schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS production_schedule_items (
    id SERIAL PRIMARY KEY,
    schedule_id INTEGER NOT NULL REFERENCES production_schedules(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    schedule_date DATE NOT NULL,
    job_id VARCHAR(50) NOT NULL,
    supervisor_id INTEGER NOT NULL REFERENCES users(id),
    operator_id INTEGER NOT NULL REFERENCES users(id),
    shift_id INTEGER NOT NULL REFERENCES shifts(id),
    part_id INTEGER NOT NULL REFERENCES parts(id),
    process_id INTEGER NOT NULL REFERENCES part_processes(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL CHECK (end_time <> start_time),
    total_minutes INTEGER NOT NULL DEFAULT 0 CHECK (total_minutes >= 0),
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_production_schedule_items_schedule_id ON production_schedule_items(schedule_id);
CREATE INDEX IF NOT EXISTS idx_production_schedule_items_is_deleted ON production_schedule_items(is_deleted);

CREATE TRIGGER update_production_schedule_items_updated_at BEFORE UPDATE ON production_schedule_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
