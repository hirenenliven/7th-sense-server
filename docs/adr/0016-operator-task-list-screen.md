# 0016. A Task List screen for the operator tasks' master list

- **Date:** 2026-11-17
- **Status:** Accepted (supersedes the task-list part of 0015)

## Context

ADR 0015 kept the maintenance task master list to a quick add from the Add Operator Task form: no screen to see the list, no edit, no delete. The superadmin wants to manage the list itself — add a task such as "Coolant not working and need to fix" and see it in the Task dropdown, correct a description, and retire tasks no longer used. Everything else in 0015 (two tables, machine type = catalog model, superadmin-only writes, company-scoped reads) stays.

## Decision

- **Task List screen** (client: Account Center › Operator Task › Task List) beside Operator Tasks, with add, edit (description, status) and delete. The form's quick add stays.
- **API**: `PUT /maintenance-tasks/{id} {description, status}` and `DELETE /maintenance-tasks/{id}`, both `superAdminOnly`. `GET /maintenance-tasks` returns active and inactive tasks with `usageCount` (operator tasks not deleted that use each), worked out in the query.
- **A task in use cannot be deleted** (400 naming the count): its operator tasks would lose their task. **Inactive** is the way to retire it: it is no longer accepted for a new operator task (server) or offered in the dropdown (client); the operator tasks that already use it keep it, and the dropdown still shows it on such a row when edited.
- **Editing a description renames it everywhere**: operator tasks point at the task, so they show the new text. Uniqueness (ignoring case) still holds.

## Alternatives considered

- **Delete a task in use and soft-delete its operator tasks too.** Silent removal of customers' checks from one click on a shared list. Rejected.
- **Edit creates a new task** (keeping old operator tasks on the old text). The usual need is a typo fix that should apply everywhere. Rejected.
- **Keep the quick add only.** No way to fix or retire a task. Rejected (the request).

## Consequences

- The seeded descriptions (including tasks 10 and 11 of 0015) can now be corrected on the screen; no migration needed.
- A task used by operator tasks stays in the list until they are removed; the screen offers Make inactive instead of Delete.
- `usageCount` is a correlated count per task; fine for a list of tens of tasks, worth revisiting only if the list grows to thousands.
