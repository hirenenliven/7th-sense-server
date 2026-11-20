# 0022. Test PLC node IDs in a short OPC UA session of their own

- **Date:** 2026-11-19
- **Status:** Proposed

## Context

On the Machine PLC node IDs page (ADR 0020) the superadmin sets the OPC UA node id each PLC variable of a machine is read from. Machines are wired differently, so a node id may not exist on a controller or may bring no value. The superadmin needs to see, with one click, whether a node id works and what it reads, also for a node id only typed and not yet saved, so a wrong one can be changed and tried again on the spot. The PLC reader (ADR 0021) already holds a subscription per Siemens machine, but it subscribes only to saved node ids and sends only changes.

## Decision

- **`POST /machines/{id}/plc-variables/test {nodeIds}`** (superadmin, 1–100 node ids, each checked like a save) reads every node id once and answers per node `ok` + `value` + `valueType`, or the controller's `status` (e.g. BadNodeIdUnknown). Nothing is stored, published on the socket or changed in `machine_variables`.
- **A session of its own**, opened for the request and closed after it (`MachinePlcVariableService.TestNodeIDs`, `service/machine_plc_node_probe.go`), connected exactly as the PLC worker connects (`opc.tcp://live_ip:port`, the machine's security mode, login and the client certificate), with one OPC UA Read of the Value attribute (`opcua.Connection.ReadValues`). Connect and read are each bounded by `PLC_CONNECT_TIMEOUT_SEC`.
- **Only machines the PLC reader reads** (Siemens, `model.ReadsPLCOverOPCUA`) with a Live IP and Port are tried; otherwise `tested: false` with the reason. A connect or read failure is part of the answer (`connected: false`, `message`, `hint` from the diagnostics' `plcErrorHint`), not an HTTP error, because it is the result the superadmin asked for.
- The PLC node IDs list also carries what the reader does (`reader` and per row `readerNodeId` / `readerStatus` from `opcua.StatusBoard`), and the client shows each variable's live value from the socket beside it.

## Alternatives considered

- **Read through the running PLC worker's session** (send the request to the worker goroutine and read on its connection): no extra session, but it only works while the worker is connected, so it cannot help find out why the reader fails, and it needs a request channel into every worker. Rejected for now.
- **Save, then watch the live value:** no new endpoint, but a wrong node id would be saved and read by the reader first, and a node whose value never changes sends nothing, so "no value" would not tell working from broken.
- **Subscribe instead of Read:** the first notification of a monitored item would also give the value, but a Read is one request and answers a refused node at once.

## Consequences

- While a test runs the NCU holds a third OPC UA session for the machine (DLL, PLC reader, test) for a moment; its session limit must allow it.
- A test takes up to twice `PLC_CONNECT_TIMEOUT_SEC` (default 10 s) on an unreachable machine; the client waits up to 30 s.
- Tests are not limited per machine: the page disables the button while a test runs, but two superadmins could test at once (more sessions).
- Basic256Sha256 is exercised by the same connect code as the PLC reader, which is still untested against a real NCU (ADR 0021).
