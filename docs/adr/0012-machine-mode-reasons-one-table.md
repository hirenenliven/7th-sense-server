# 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log

- **Date:** 2026-11-06
- **Status:** Accepted

## Context

An operator must say why a machine was down or idle, by clicking its bar on the Productivity Gantt chart. The master lists exist (`downtime_reasons`, `idle_reasons`, migration 027), but nothing recorded a reason against a period. The old system had two tables, `downtime_reason` and `idle_reason`, with a machine, a date and two clock times and no link to the period. The chart's bars are `machine_mode_log` rows, clipped to the chart's window.

## Decision

- One new table, `machine_mode_reasons` (migration 033): `mode_log_id` (the period), `downtime_reason_id` or `idle_reason_id` (a CHECK allows exactly one), `started_at`, `ended_at`, soft delete. A period may have several rows, each for a part of it.
- One set of files and one endpoint, `/machine-mode-reasons`. The client names only the period; the server works out the list from its mode (Downtime 6, Idle 5) and refuses other modes and a period still running.
- The time is stored as two instants, so a reason past midnight needs nothing special.
- `MachineModeReasonService` checks: inside the period; no overlap with the period's other rows (the period's row is locked in the transaction, so two saves at once cannot both pass); an active reason of the machine's customer. A row keeps a reason later made inactive or deleted.
- Whoever may read the machine may add, edit and delete, because operators enter reasons (`MachineService` decides, others 404). No role guard on the routes.

## Alternatives considered

- **Two tables and two sets of files, like the masters.** The rules are the same for both, so every rule would be written and tested twice. Rejected.
- **One `reason_id` plus a type column.** The database could not check the reason exists. Rejected for the two foreign keys.
- **Machine, date and clock times like the old tables.** Nothing would tie a reason to its period, and a period's end could not be checked. Rejected.
- **An exclusion constraint for the overlap rule.** It needs the `btree_gist` extension, which this database does not have, for one rule the row lock already guards. Rejected.
- **Only managers may record reasons (`managersOnly`).** The operator at the machine is the one who knows why it stopped. Rejected.

## Consequences

- Reports can join a reason to its period, mode and machine through `mode_log_id`.
- Deleting a `machine_mode_log` row removes its reasons (`ON DELETE CASCADE`). Nothing deletes those rows today.
- Adding a reason to another mode means adding it to `reasonTypeOfPeriod` and, if it needs its own master list, a new column in the CHECK.
- The Gantt chart does not yet mark which bars have a reason.
