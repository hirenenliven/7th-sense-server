# 0015. Operator tasks: a shared task list, assigned per customer and machine model

- **Date:** 2026-11-17
- **Status:** Accepted — the task-list part (no task screen, no edit or delete) superseded by [0016](0016-operator-task-list-screen.md)

## Context

The old system let the superadmin say which checks an operator does on a customer's machines ("Check for any Oil leakage", "Clean Panel AC Filter"). It had a master list, `maintenance_tasks` (id, description, status), and `mm_operator_tasks` (a free-text `iMachineType`, the task id, the customer, a date). The machine type was the old product ("NX Series High-End", "VMC-Entry Level"). The new system has no product or machine type list in use: machines are described by the catalog Brand → Controller → Model (`machines.product_id` and `machine_type_id` point at nothing).

## Decision

- **Two tables** (migration 036): `maintenance_tasks` (description ≤ 100, unique ignoring case among rows not deleted, status active / inactive, soft delete), seeded with the 15 old tasks; `operator_tasks` (`user_id` customer, `model_id` nullable, `maintenance_task_id`, soft delete) with a unique index on `(user_id, COALESCE(model_id, 0), maintenance_task_id)` among rows not deleted.
- **Machine Type is a catalog model** (`models.id`), optional. No model means every machine of the customer and counts as its own type for the duplicate rule. No new master list to keep in step with the catalog.
- **The task list is shared** by every customer. The superadmin adds a task from the Add Operator Task form (`POST /maintenance-tasks`); there is no separate task screen yet, no edit or delete of a task.
- **Writes are superadmin only** (`superAdminOnly`): `POST /maintenance-tasks`, `POST/PUT/DELETE /operator-tasks`. **Reads** are for every logged-in user, company-scoped like parts (`companyOwner`), so an operator checklist can read them later without an API change.
- The service checks what a row names before saving: a customer (`UserRepository.IsCustomer`), a catalog model (`MachineModelService.GetByID`), an active task. The database constraints back these up.

## Alternatives considered

- **A new `machine_types` list** (and `machines.machine_type_id`). It would match the old screen's words, but needs its own screen and a way to set it on every machine; the catalog model already says what a machine is. Rejected for now.
- **Free text for the machine type**, as the old varchar. Typos split one type into several and nothing ties a task to machines. Rejected.
- **A separate task list screen** (add / edit / delete tasks). More to build for a list that rarely changes; quick add covers adding. Deferred.
- **Task text stored on each operator task** (no master list). Every customer would retype the same checks and they could not be counted together. Rejected.

## Consequences

- Tasks 10 and 11 of the old list were cut short in the old screens; the seed completes them ("…properly closed", "…in moving/running condition"). Changing a seeded text needs a migration until a task screen exists.
- A model or customer that is soft-deleted keeps its operator tasks (the foreign keys only guard hard deletes); the list still shows the saved names.
- An operator checklist (marking tasks done per machine and day, the old `oper_mainten_tasks`) can be built on these tables: a machine's tasks are its customer's rows with its `model_id` or none.
