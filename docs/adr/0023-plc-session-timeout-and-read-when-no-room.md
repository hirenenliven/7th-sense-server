# 0023. PLC reader: short session timeout, and read the nodes the controller has no room to monitor

- **Date:** 2026-11-20
- **Status:** Proposed

## Context

The PLC reader (ADR 0021) subscribes to every PLC variable of a Siemens machine: one monitored item per variable. A Sinumerik OPC UA server monitors only a limited number of items at once, shared by every client connected to it (the CNC DLL's session, the PLC reader, node tests, HMI / SCADA tools). On a live machine most PLC variables came back `StatusBadTooManyMonitoredItems` and were never read, while a script that *read* the same node ids got every value: a Read reserves no monitored item.

gopcua asks the controller for a 20-minute session timeout by default. A server that is killed instead of shut down (a stop from the IDE, a crash, a hot reload) cannot close its session, so the controller keeps that session and all its monitored items for 20 minutes. Each restart within that time adds a full set of items, which is how the limit was reached.

## Decision

- **Short session timeout:** `opcua.Settings.SessionTimeout`, from `PLC_SESSION_TIMEOUT_SEC` (default 120), is asked for by the PLC worker and the node test. A killed server's items are freed after about two minutes. A live session talks to the controller at least every `PLC_HEALTH_SEC` (30), so it never times out.
- **Read the nodes the controller has no room for:** a monitored item refused with `BadTooManyMonitoredItems` only (`opcua.ItemResult.NoRoom`) is read instead, in one OPC UA Read for all of them every `PLC_POLL_INTERVAL_MS` (default 1000), on the same session (`plcWorker.poll`). Only a value or status that changed since the previous read is handed on (`polledChanges`), as a publish would, so the hub, `machine_variable_values` and `machine_telementry` see the same thing as for a subscribed node. A failed read ends the session like a failed health read. Such a node shows no refusal on the diagnostics; the log says which nodes are read instead (`plc: controller has no room to monitor these nodes, reading them instead`).
- Every other refusal (BadNodeIdUnknown, BadUserAccessDenied, an invalid node id) is not read: reading would fail the same way.

## Alternatives considered

- **Poll every PLC variable, no subscription:** no dependence on the item limit at all, but one Read per interval for every variable of every machine, and changes faster than the interval are missed. Subscriptions stay the default; reading is only for what does not fit.
- **Retry the refused items every few minutes:** they come back once the controller frees room, but get nothing until then. Reading gives them values at once, so a retry adds little.
- **Raise the limit on the controller:** right when the limit really is too low for the machine's variables, but it is a Siemens setting / licence outside the server, and does not help against leftover sessions.

## Consequences

- New settings: `PLC_SESSION_TIMEOUT_SEC` (120), `PLC_POLL_INTERVAL_MS` (1000).
- A node that is read instead gets its values up to `PLC_POLL_INTERVAL_MS` late, and a value that changes and changes back between two reads is not seen.
- The controller answers one extra Read per interval per machine while any node is read instead; nothing when every node fits.
- A server that is killed holds the controller's items for up to `PLC_SESSION_TIMEOUT_SEC` instead of 20 minutes.
- Amends ADR 0021's "Subscriptions, not polling": subscriptions first, reads only for nodes the controller has no room for.
