# 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector

- **Date:** 2026-11-19
- **Status:** Proposed (amended by [0023](0023-plc-session-timeout-and-read-when-no-room.md))

## Context

PLC variables (source `plc` in `cnc_variable_addresses`, OPC UA node ids such as `ns=2;s=/Random/Plc/DB9006.DBD8`: hydraulic pressure, temperatures, chiller flow …) were configured per machine (ADR 0020) but never read: `CncCommunication.dll` only hands over its fixed snapshot fields and has no read function for other nodes. These node ids belong to the Siemens Sinumerik OPC UA server (port 4840), which the DLL already talks to for the snapshot. The values must be read, stored like the other variables, pushed to the screens over the live socket the moment they change, and shown on the diagnostics page.

## Decision

- **Library: `github.com/gopcua/opcua` v0.9.1** (MIT). Chosen because it is the most used pure-Go OPC UA stack, maintained (v0.9.1 released 2026-08-09, regular releases since 2019), needs Go ≥ 1.23, builds for `GOARCH=386` (the server's DLL build) and 64-bit alike, and supports what we need: endpoint discovery, security policies None / Basic256Sha256 with our client certificate, username / anonymous login, subscriptions with monitored items (data-change notifications), and a test server package. Wrapped in `internal/opcua` (connect, subscribe, value conversion, `StatusBoard`), the protocol side only, like `internal/cnc` for the DLL.
- **Subscriptions, not polling.** One subscription per machine (publish `PLC_PUBLISH_INTERVAL_MS` 500, sampling `PLC_SAMPLING_INTERVAL_MS` 250), one monitored item per PLC variable; the controller sends only changes. A read of the server time every `PLC_HEALTH_SEC` (30) notices a dead connection (a switched-off machine sends no TCP reset, and a subscription is silent while nothing changes); three publish errors in a row also reconnect. gopcua's own auto-reconnect is off: the worker reconnects with backoff (1 s doubling to 5 min) and subscribes again, so its state stays simple.
- **In-process, beside the CNC collector.** `monitor.PlcCollector` (started by `collector.Supervisor.StartPLCDataCollection`, `PLC_ENABLED`, default on) runs one `plcWorker` per active Siemens machine (`model.ReadsPLCOverOPCUA`; `CNC_MACHINE_IDS` applies) and follows `machines_changed` / `cnc_config_changed` like the CNC collector. Endpoint `opc.tcp://live_ip:port`, the machine's `security_mode`, `username` / `password` (empty = anonymous), certificate `client_cert.der` / `client_key.der` with application URI `urn:go-opcua-client`.
- **Same storage and live path as the snapshot variables.** Each publish goes to `live.Hub.PublishValues` first, then numbers are upserted in `machine_variable_values` and a changed number is inserted in `machine_telementry` (text is published, not stored). The hub keeps the known codes per reader (`SetMachineCodes(machineID, live.ReaderCNC | live.ReaderPLC, codes)`), so the two workers do not overwrite each other.
- **One reader per variable.** `service.VariableSources(engine, plcEnabled)` = snapshot (+ plc on a Siemens machine with PLC on); the CNC worker, the PLC worker and the diagnostics all pick each variable's address from it (`FindVariableConfig`, cheapest source first), and the CNC worker leaves variables whose address is plc to the PLC worker.
- **The client certificate is kept.** `cert.EnsureCertificates` (was `GenerateCertificates`, which made a new certificate at every start) keeps an existing, unexpired certificate and creates one valid 10 years only when missing or expired: the NCU trusts that exact certificate, for the DLL and the PLC reader alike.
- **Diagnostics.** `GET /machines/{id}/diagnostics` gains `plc`: whether the machine is read, the worker's state from the in-memory `opcua.StatusBoard`, endpoint, security policy, last error, a verdict with message and hint, and every PLC variable (node id, the controller's answer, live value, last stored number). The client shows it as a "PLC variables (OPC UA)" section with its own live-socket subscription.

## Alternatives considered

- **`github.com/awcullen/opcua`.** Capable, but last release 2024-12 and far fewer users. Rejected.
- **open62541 through cgo, or ask the DLL vendor for a node read function.** cgo complicates the 32-bit build and the vendor function does not exist yet; gopcua needs neither. Rejected for now; the DLL keeps the snapshot.
- **Poll each node (Read every N seconds).** More traffic, slower than a subscription for a change, and still needs the same reconnect logic. Rejected.
- **A separate OPC UA service.** Another process to deploy and watch, and it would need the hub over the network. The server already runs the collector in-process. Rejected.
- **gopcua's auto-reconnect.** Restores subscriptions itself, but its state changes are harder to report and test; our own backoff loop matches the CNC worker. Rejected.

## Consequences

- PLC variables of Siemens machines are read, stored and live; Fanuc and other engines are not read over OPC UA (their PLC node ids stay configuration only; the diagnostics say so).
- Each Siemens machine now holds two OPC UA sessions (the DLL's and the PLC reader's); Sinumerik limits sessions per NCU, so this must fit the controller's licence.
- Basic256Sha256 against a real NCU is not covered by tests (gopcua's test server implements security None only): try it on a machine first. The client certificate must be trusted on the NCU (`BadSecurityChecksFailed` / certificate errors otherwise; the diagnostics hint says so).
- New settings: `PLC_ENABLED`, `PLC_PUBLISH_INTERVAL_MS`, `PLC_SAMPLING_INTERVAL_MS`, `PLC_CONNECT_TIMEOUT_SEC`, `PLC_HEALTH_SEC`, `CNC_LOG_PLC`.
- Tests: `internal/opcua/connection_test.go` (gopcua test server: subscribe, refused / invalid nodes, change, lost connection), `tests/plc_collector_test.go` (end to end: OPC UA server → hub, `machine_variable_values`, `machine_telementry`, diagnostics), `internal/service/machine_diagnostics_plc_test.go`, `internal/live/hub_test.go` (codes per reader).
