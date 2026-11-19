# 0017. Machine accuracy as one shared list, managed by the superadmin

- **Date:** 2026-11-18
- **Status:** Accepted

## Context

A machine (`machines.accuracy_id`) and a part process (`part_processes.accuracy_id`) have a machine accuracy level. Until now the levels were a fixed list in the client (`constants/machines.js`: 1 Good, 2 Average, 3 Poor) with no table behind the ids, so nobody could add or rename a level without a code change. The old system had a `machine_accuracy` table and a Machine Accuracy screen (Accuracy Id, Accuracy, Edit / Delete). The superadmin wants that screen back in Account Center, with one list applied to every customer and user: nothing per customer.

## Decision

- **Table `machine_accuracies`** (migration 037): `id`, `name` (≤ 50, unique ignoring case among rows not deleted), soft delete. Seeded with 1 Good, 2 Average, 3 Poor (the ids machines and part processes already store), then the sequence continues after them.
- **API `/machine-accuracies`**: every logged-in user reads (`GET`, ordered by id, with `usageCount` = machines + part processes not deleted that use it); only the superadmin adds, renames and deletes (`superAdminOnly`). No `userId`, no customer column, no status.
- **An accuracy in use cannot be deleted** (400 naming the count): those machines and processes would point at nothing. Renaming changes it everywhere.
- **Machines and parts check the id**: `MachineService.checkReferences` and `PartService.checkReferences` refuse an accuracy not in the list ("Please select valid machine accuracy" / "… level"), now that the list can change.
- **Client**: Account Center › Machine Accuracy (form beside the list, like Task List); the machine and part forms load the list from the API instead of the fixed constant.

## Alternatives considered

- **Keep the fixed client list.** No way to add a level without a release. Rejected (the request).
- **A list per customer** (`user_id`, like downtime reasons). Explicitly not wanted: the levels are the same everywhere. Rejected.
- **Foreign keys from `machines` / `part_processes`.** A soft-deleted accuracy still satisfies a foreign key, and the delete is already refused while in use; a service check gives a readable 400 instead of a constraint error. Not added.
- **Status (active / inactive) like the task list.** Not asked for, and a level in use can simply be renamed. Left out; easy to add later.

## Consequences

- Existing machines and part processes keep their accuracy with no data change (seed keeps ids 1–3).
- Migration 037 must run before the server: without the table, `GET /machine-accuracies` and every machine / part save that names an accuracy fail.
- `usageCount` is two correlated counts per row; fine for a list of a handful of levels.
