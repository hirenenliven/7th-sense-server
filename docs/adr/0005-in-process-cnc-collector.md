# 0005. CNC collector runs inside the API process

- **Date:** 2026-11-03
- **Status:** Accepted (Inferred from existing code)

## Context

Machine modes, cycles and variables are read from CNC controllers through a vendor DLL.

## Decision

The server starts a collector (`internal/collector`, `internal/monitor`) that writes `machine_mode_log`, `machine_cycles` and `machine_variable_values`; the API reads those tables.

## Alternatives considered

- A separate collector service
- Push from the machines

## Consequences

- One deployable; dashboards read the same database.
- The API process depends on the Windows DLL; collector load shares the API host.
