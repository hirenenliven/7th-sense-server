-- Migration: customer breaks (master list), shifts link breaks, holidays / week offs
--
-- 1. breaks: a customer's (user_id) breaks — name, start and end time ("Tea" 10:30–10:45).
--    Name unique per customer (case-insensitive). A break may run past midnight.
-- 2. shift_breaks becomes the link between a shift and the customer's breaks it takes
--    (no soft delete, like part_process_machines). Migration 019 kept each shift's
--    breaks inline in shift_breaks; those rows move to breaks (one per customer, name and
--    times; a name used with other times gets the times added to it) and are linked back.
-- 3. holidays: a customer's holidays (name, start and end date) and week offs (a weekday,
--    1 = Monday … 7 = Sunday; name optional). A weekday is a customer's week off once.
--
-- Rollback (the inline breaks of 019 are not restored):
--   DROP TABLE IF EXISTS holidays, shift_breaks, breaks;
--   then re-run the shift_breaks part of 019_add_shifts.sql.

BEGIN;

-- 1. Breaks ------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS breaks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL CHECK (end_time <> start_time),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS breaks_user_name_key
    ON breaks (user_id, LOWER(name)) WHERE is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_breaks_user_id ON breaks(user_id);
CREATE INDEX IF NOT EXISTS idx_breaks_is_deleted ON breaks(is_deleted);

CREATE TRIGGER update_breaks_updated_at BEFORE UPDATE ON breaks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 2. Shift ↔ break links (replacing the inline breaks of 019) ----------------

ALTER TABLE shift_breaks RENAME TO shift_breaks_019;

-- One break per customer, name and times; a name already used with other times gets
-- the times added ("Tea 10:30-10:45") so names stay unique per customer.
INSERT INTO breaks (user_id, name, start_time, end_time)
SELECT user_id,
       CASE WHEN ROW_NUMBER() OVER (PARTITION BY user_id, LOWER(name) ORDER BY start_time, end_time) = 1 THEN name
            ELSE LEFT(name, 88) || ' ' || to_char(start_time, 'HH24:MI') || '-' || to_char(end_time, 'HH24:MI')
       END,
       start_time, end_time
FROM (
    SELECT DISTINCT s.user_id, old.name, old.start_time, old.end_time
    FROM shift_breaks_019 old
    JOIN shifts s ON s.id = old.shift_id
    WHERE old.is_deleted = false AND s.is_deleted = false
) inline_breaks;

CREATE TABLE IF NOT EXISTS shift_breaks (
    id SERIAL PRIMARY KEY,
    shift_id INTEGER NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    break_id INTEGER NOT NULL REFERENCES breaks(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT shift_breaks_shift_break_key UNIQUE (shift_id, break_id)
);

-- shift_id lookups use the unique index above; this one serves "which shifts take break X"
CREATE INDEX IF NOT EXISTS idx_shift_breaks_break_id ON shift_breaks(break_id);

CREATE TRIGGER update_shift_breaks_updated_at BEFORE UPDATE ON shift_breaks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

INSERT INTO shift_breaks (shift_id, break_id)
SELECT DISTINCT old.shift_id, b.id
FROM shift_breaks_019 old
JOIN shifts s ON s.id = old.shift_id AND s.is_deleted = false
JOIN breaks b ON b.user_id = s.user_id AND b.is_deleted = false
    AND b.start_time = old.start_time AND b.end_time = old.end_time
    AND (LOWER(b.name) = LOWER(old.name) OR b.name LIKE LEFT(old.name, 88) || ' %')
WHERE old.is_deleted = false;

DROP TABLE shift_breaks_019;

-- The old table held the usual names while the new one was created: take them back
ALTER TABLE shift_breaks RENAME CONSTRAINT shift_breaks_pkey1 TO shift_breaks_pkey;
ALTER TABLE shift_breaks RENAME CONSTRAINT shift_breaks_shift_id_fkey1 TO shift_breaks_shift_id_fkey;
ALTER SEQUENCE shift_breaks_id_seq1 RENAME TO shift_breaks_id_seq;

-- 3. Holidays and week offs --------------------------------------------------

CREATE TABLE IF NOT EXISTS holidays (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    holiday_type VARCHAR(10) NOT NULL CHECK (holiday_type IN ('holiday', 'week_off')),
    name VARCHAR(100),
    start_date DATE,
    end_date DATE,
    week_day SMALLINT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- A holiday has a name and dates; a week off has a weekday
    CONSTRAINT holidays_kind_check CHECK (
        (holiday_type = 'holiday' AND name IS NOT NULL AND start_date IS NOT NULL
            AND end_date IS NOT NULL AND end_date >= start_date AND week_day IS NULL)
        OR (holiday_type = 'week_off' AND week_day BETWEEN 1 AND 7 AND start_date IS NULL AND end_date IS NULL)
    )
);

-- A weekday is a customer's week off once
CREATE UNIQUE INDEX IF NOT EXISTS holidays_user_week_day_key
    ON holidays (user_id, week_day) WHERE holiday_type = 'week_off' AND is_deleted = false;

-- The same holiday (name and start date) is not added twice
CREATE UNIQUE INDEX IF NOT EXISTS holidays_user_name_date_key
    ON holidays (user_id, LOWER(name), start_date) WHERE holiday_type = 'holiday' AND is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_holidays_user_id ON holidays(user_id);
CREATE INDEX IF NOT EXISTS idx_holidays_is_deleted ON holidays(is_deleted);

CREATE TRIGGER update_holidays_updated_at BEFORE UPDATE ON holidays
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
