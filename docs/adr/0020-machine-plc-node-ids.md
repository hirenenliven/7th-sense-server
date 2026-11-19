# 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise

- **Date:** 2026-11-19
- **Status:** Proposed

## Context

Each PLC variable has a default OPC UA node ID per CNC engine and controller version in `cnc_variable_addresses` (source `plc`, today only Siemens: 45 rows, 828D / 840D / every version). Machines are wired differently, so some need another node ID for the same variable. The collector must read the machine's node ID when one is set and the default otherwise, and the superadmin needs a screen to see and change them per machine quickly. For now this is for PLC variables only.

## Decision

- **No new table.** The machine's own node ID is the existing per-machine override row `machine_variables (machine_id, id, source, address)` with `source = 'plc'`. `MachineTelemetryRepository.FindVariableConfig` already reads `COALESCE(mv.address, default)`, and the table's trigger already sends `cnc_config_changed`, so the machine's worker reloads its variables at once.
- **Default stays live.** Saving the default node ID, or resetting, removes the machine's node ID instead of storing a copy, so a later change of the default reaches every machine on it. A row that also holds other overrides (`is_enabled`, `deadband`, `min_interval_ms`) keeps them; one that only held the node ID is deleted (the table has no soft delete).
- **Which variables and defaults:** every variable with a `plc` address in `cnc_variable_addresses` for any engine. Without a node ID of its own a variable is **Default**: the default is the best `plc` address in `cnc_variable_addresses`: the machine's own CNC engine and controller version first, else another engine's (PLC defaults exist for Siemens only, so a Fanuc machine reads the Siemens node IDs by default), then its version before every version before another version. `FindVariableConfig` (collector) and the diagnostics rank plc addresses the same way, so the screen and the connection agree; snapshot and NC addresses still come only from the machine's engine. Each row has `nodeStatus` default / updated.
- **API:** `GET /machines/{id}/plc-variables`, `PUT` / `DELETE /machines/{id}/plc-variables/{variableId}`, `DELETE /machines/{id}/plc-variables` (reset all); `superAdminOnly`, like the diagnostics (installation details of a customer's controller). Node ID validated as OPC UA (`ns=<n>;s=…`, also `i=` / `g=` / `b=`, ≤ 255).
- **Readable sources only.** The DLL reads snapshot fields only today. `FindVariableConfig` now uses a machine's own address only when its source is one the collector can read, else the readable default — so a stored PLC node ID waits for the DLL instead of raising a `config_error` event on every reload.
- **Client:** All Machines › a machine › **PLC node IDs** (`/machines/:id/plc-variables`): a table with the node ID edited in the row (Enter saves, Esc undoes), Default / Updated / Unsaved status, Save N changes, reset per row and Reset all.

## Alternatives considered

- **A new `machine_plc_node_ids` table.** Duplicates `machine_variables`, needs its own trigger and a second COALESCE in the collector query. Rejected.
- **List only the variables with a default for the machine's engine.** First version; a Fanuc machine (no PLC defaults) showed no variables at all, so its node IDs could not be entered. Rejected.
- **"Not set" for a variable with no default for the machine's engine** (second version): the machine read nothing until each node ID was typed, and the user expects every machine to start on the default node IDs. Rejected: another engine's plc address is the default.
- **Copy every default into the machine when it is created.** Defaults could no longer be corrected in one place, and "is it default?" would need a comparison. Rejected.
- **Managers (customer / admin) edit their own machines.** A wrong node ID silently stops data; kept to the superadmin like the diagnostics. Can be widened later.
- **Snapshot and NC addresses on the same screen.** Out of scope for now (PLC only); the same rows and endpoints could take a `source` later.

## Consequences

- No migration. Writes go through the existing trigger, so the collector follows edits without polling.
- PLC node IDs are read over OPC UA on Siemens machines (ADR 0021); a saved node ID is picked up at once (the PLC worker subscribes again).
- A PLC node ID override on a variable whose default is snapshot or NC would not be shown here; this screen only writes `plc` rows for variables with a `plc` default.
- Test: `tests/machine_plc_variable_test.go`.
