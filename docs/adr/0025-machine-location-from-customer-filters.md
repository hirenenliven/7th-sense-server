# 0025. A machine is saved at one level of its customer's organization; the path above is read, not stored

- **Date:** 2026-11-23
- **Status:** Proposed

## Context

A customer picks, on the customer form, the organization levels its screens use (`user_profiles.filters`: country, state, city, branch, location, shop, line, cell, supervisor). Nothing used them yet. A machine only had `machines.line_id`, so it could not be placed at a level above Line, or at a cell or supervisor, and could not be filtered by country or state without a line.

What we need: the Add / Edit Machine form offers exactly the customer's levels (dependent selects, each optional), and the machine lists filter by any of them, e.g. country-wise, state-wise or location-wise for filters Country, State, Location.

Facts that shape the choice:

- The organization is one strict chain: Country → State → City → Branch → Location → Shop → Line → Cell → Supervisor, each row pointing at its parent. The lowest level picked therefore fixes every level above it.
- Levels can be moved to another parent (`PUT /user-line` takes `userShopId`, the same for every level), so ancestors copied onto a machine would go stale.
- The user machine filter (ADR 0010) already walks the organization (`organization_paths`), but only from a machine's line.

## Decision

| Topic | Decision |
|---|---|
| Storage (migration 040) | One nullable column per level on `machines` (`user_country_id` … `user_supervisor_id`, plus the existing `line_id`), with foreign keys and `CHECK (num_nonnulls(...) <= 1)`: the machine is saved at the **lowest level picked**, the other columns stay NULL. Existing machines (only `line_id`) pass unchanged; no data migration. |
| Path | View `machine_organization_paths` walks up from the saved level (`COALESCE(own column, parent of the level below)`) and gives every level's id and name. `machineSelect` joins it: responses carry the whole path (`userCountryId` … `userSupervisorId` with names; `lineId` / `lineName` are the path's line) and `customerFilters` (the customer's filters, hierarchy order). |
| API | Body `location` (`model.OrganizationPath`, the nine ids; `lineId` still accepted as the line). `MachineService.checkLocation`: every level is optional; the levels named must be one chain of the customer's organization (`MachineRepository.IsPathInOrganization`, walked up from the lowest level, "Please select valid location"). Levels the customer does not use may be named too, so a machine keeps a lower place it was saved at. An empty location is allowed. New `GET /machines/location-filters?userId=` gives the form the customer's levels (the superadmin names the customer; others get their own company's). |
| Access rule | `accessibleMachine` (ADR 0010's `allMachines`) matches the machine's own path from `machine_organization_paths` instead of `organization_paths` by `line_id`, so a machine saved at a cell, a supervisor or above its line counts under the levels above it. |
| Deletes | Every level refuses delete while machines are saved at it ("This cell has machines. Move or delete them first"; Line already did). Descendants are already covered: a level with children cannot be deleted. |
| Filtering | The machine lists filter in the client, from the path every machine already carries (the lists load every machine the user may read). No query parameters were added. |

## Alternatives considered

- **Store every level of the path on the machine** (nine columns, all filled): simplest filtering, but a level moved to another parent leaves its machines with the old ancestors unless every level update also rewrites machines.
- **One polymorphic pair `(organization_level, organization_id)`**: no foreign key, and every read would need a CASE per level.
- **Keep only `line_id` and require a line**: does not allow customers whose filters end above Line (e.g. Country, State, Location), or placing a machine at a cell / supervisor.
- **Server-side list filters (`?userCountryId=` …)**: not needed while the lists load the whole set; can be added on the same view later.

## Consequences

- The machine response has more fields (path ids and names, `customerFilters`); clients that read `lineId` still work, now the path's line (also for a machine saved at a cell or supervisor).
- Moving a level moves its machines with it, by design.
- A customer's filters can change after machines were saved: a machine keeps its place (levels no longer used are still accepted on save); lists and labels follow the current filters.
- The user machine filter still ends at Line; picking cells / supervisors there is a separate change.
- Tests: `tests/machine_location_test.go` (required levels, chain, saved level, path, access through the line, delete refusal), `TestMachineFilterMatrix` passes on the new view.
