# 0010. Machine access by an organization filter, worked out on every request

- **Date:** 2026-11-05
- **Status:** Accepted

## Context

A customer (or its admin, or the superadmin) gives a supervisor, operator, QC or shop head access to machines on the user form. Until now that was a list of ticked machines (`user_machine_access`). Customers asked to pick machines by their organization — Country → State → City → Branch → Location → Shop → Line — and, when they choose "all machines" of such a selection, to have machines added there later given to the user automatically, without editing the user again. A machine sits on a line (`machines.line_id`); cells and supervisors are below the line and no machine points at them. The customer form's Filters checklist (`user_filters`) is a different feature (which filters a customer's screens offer) and stays for customers only.

## Decision

- New table `user_machine_filters` (migration 032), one row per user: `user_country_id`, `user_state_id`, `user_city_id`, `branch_id`, `user_location_id`, `user_shop_id`, `user_line_id` (NULL = any level) and `all_machines`. Like `user_machine_access` it has no soft delete; saving the user replaces the row, deleting the user removes it.
- New view `organization_paths`: every path of a company's organization from each country down to each line (LEFT JOINs, so a level without children ends its path with NULLs).
- With `all_machines`, `MachineRepository` (`accessibleMachine`) lets the user access every company machine whose line's path matches every level the filter names (no level = every company machine), on every request. Nothing is copied into `user_machine_access`, so a new machine, or a line moved to another shop, is followed at once; the live sockets already reload a user's machines on `machines_changed`.
- Without `all_machines` the row only keeps the form's selection; access is the ticked machines as before.
- `UserService.checkMachineFilter` refuses levels that are not one chain of the company's organization ("Please select valid machine filter").
- The client loads the company's organization lists and narrows the machine list itself (`pages/Customers/machineFilter.js`, same rule as the view).

## Alternatives considered

- Expand "all machines" into `user_machine_access` rows when saving, and add rows for new machines in `MachineService.Create`: simple reads, but every place that changes a machine's line or the organization would have to keep the rows in step.
- Store the chosen level only (e.g. one state id): a lower level implies the ones above, but the form also keeps the levels above it, and a multi-level row costs nothing more to match.
- Add `country_id`, `state_id`, … columns to `machines`: duplicates the organization and goes stale when it changes.

## Consequences

- Reading machines for a filtered user joins `organization_paths`; the organization tables are small (per company), so this stays cheap.
- Machines without a line match only a filter without levels.
- Users created before this change keep their ticked machines and have no filter.
- One filter per user; several selections (e.g. two states) would need more rows per user.
