-- Migration: shifts of a customer and their breaks
--
-- A shift belongs to a customer (user_id) and runs from start_time to end_time on its
-- working days (ISO weekdays, 1 = Monday … 7 = Sunday). A shift whose end is not after
-- its start runs past midnight (night shift). Its breaks lie inside the shift and do not
-- overlap. Shift and break hours are worked out from the times, never stored.
--
-- Rollback: DROP TABLE IF EXISTS shift_breaks, shifts;

BEGIN;

CREATE TABLE IF NOT EXISTS shifts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    shift_type VARCHAR(10) NOT NULL CHECK (shift_type IN ('day', 'night')),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL CHECK (end_time <> start_time),
    working_days INTEGER[] NOT NULL DEFAULT '{1,2,3,4,5,6,7}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- A customer cannot have two shifts with the same name (case-insensitive)
CREATE UNIQUE INDEX IF NOT EXISTS shifts_user_name_key
    ON shifts (user_id, LOWER(name)) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_shifts_user_id ON shifts(user_id);
CREATE INDEX IF NOT EXISTS idx_shifts_is_deleted ON shifts(is_deleted);

CREATE TRIGGER update_shifts_updated_at BEFORE UPDATE ON shifts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE TABLE IF NOT EXISTS shift_breaks (
    id SERIAL PRIMARY KEY,
    shift_id INTEGER NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    name VARCHAR(100) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL CHECK (end_time <> start_time),
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shift_breaks_shift_id ON shift_breaks(shift_id);
CREATE INDEX IF NOT EXISTS idx_shift_breaks_is_deleted ON shift_breaks(is_deleted);

CREATE TRIGGER update_shift_breaks_updated_at BEFORE UPDATE ON shift_breaks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
