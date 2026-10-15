-- Migration: a break belongs to (at most) one shift — breaks.shift_id replaces shift_breaks
--
-- Migration 020 linked shifts and breaks through the shift_breaks table. A break is only
-- ever taken by one shift, so the link moves onto the break: breaks.shift_id is the shift
-- that takes it, NULL while no shift does (added on Shift › Breaks, assigned later).
-- Deleting a shift sets its breaks' shift_id back to NULL (service.ShiftService).
--
-- Rollback:
--   recreate shift_breaks as in 020 and fill it from breaks.shift_id, then
--   ALTER TABLE breaks DROP COLUMN shift_id;

BEGIN;

ALTER TABLE breaks ADD COLUMN IF NOT EXISTS shift_id INTEGER REFERENCES shifts(id);

CREATE INDEX IF NOT EXISTS idx_breaks_shift_id ON breaks(shift_id);

-- Each break was linked to one shift at most; keep that shift
UPDATE breaks b
SET shift_id = links.shift_id
FROM (
    SELECT sb.break_id, MIN(sb.shift_id) AS shift_id
    FROM shift_breaks sb
    JOIN shifts s ON s.id = sb.shift_id AND s.is_deleted = false
    GROUP BY sb.break_id
) links
WHERE b.id = links.break_id;

DROP TABLE IF EXISTS shift_breaks;

COMMIT;
