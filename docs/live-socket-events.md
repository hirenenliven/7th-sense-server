# Live data — Socket.IO event contract

Decision: [ADR 0009](adr/0009-socketio-live-data.md). Code: [internal/live/hub.go](../internal/live/hub.go) (latest values, rooms, intervals), [internal/handler/live_handler.go](../internal/handler/live_handler.go) (handshake, socket registry, access), [internal/dto/live_dto.go](../internal/dto/live_dto.go) (messages and limits). Client: `../7th-sense-client/src/services/liveSocket.js`, `hooks/useMachineLive.js`, `constants/live.js`. Plain-words walkthrough (machine → server → screen, how subscribing works): [live-data-flow.md](live-data-flow.md).

## Connecting

| | |
|---|---|
| URL | `/api/v1/socket.io/` (Socket.IO default path under the API base); a "log in as" tab uses `/api/as/<key>/v1/socket.io/` |
| Transport | WebSocket only (`transports: ['websocket']`); long-polling is not served |
| Auth | The session cookie of the handshake (`ts_session`, or `ts_login_as` under `/api/as/<key>/`). No session → HTTP 401 before the upgrade (the browser sees `connect_error`) |
| Origin | When the browser sends `Origin`, its host must equal the request's `Host`, else 403 |
| Heartbeat | ping every 25 s, dropped after 20 s more without a pong |
| Server close | Logout of this session, an ended session, or more than `LIVE_MAX_SOCKETS_PER_USER` sockets of the user: the server sends a Socket.IO disconnect (`io server disconnect` in the browser, which then does not reconnect by itself) |

## Browser → server

| Event | Payload | Acknowledgement |
|---|---|---|
| `subscribe` | `{ machineIds: [225], codes: ["SpindleLoad"], channels: ["status", "event"], intervalMs: 1000 }` | `{ ok, denied: [ids], unknown: [codes], snapshots: [Snapshot] }`, or `{ ok: false, message }` when the message is invalid |
| `unsubscribe` | `{ machineIds: [225] }` | `{ ok: true }` |

- `codes` are exact `cnc_variables.code` values; `channels` are `status` and / or `event`; at least one code or channel.
- `intervalMs` is optional: 0 = 1 000 ms; below `LIVE_MIN_INTERVAL_MS` (200) it is raised to it.
- A `subscribe` for a machine **replaces** this socket's earlier subscription of that machine (the client merges its components' codes per machine first).
- Limits: 500 machines and 100 codes per message, a code at most 64 characters, a message at most 64 KB.
- `denied`: machines the user may not read (`MachineService` rule) — nothing is joined for them. `unknown`: codes the collector does not read for those machines (only once the collector loaded the machine's variables; before that every code is accepted).
- `userId`, `role` or room names in a message are ignored: the user comes from the session and the server builds every room.

**Snapshot** (one per machine not denied) — the current value of the codes asked for:

```json
{ "machineId": 225, "seq": 1840, "at": "2026-11-05T10:42:07.41+05:30",
  "values": { "SpindleLoad": 41.5, "progName": "O1234_FLANGE" },
  "online": true, "modeId": 3, "modeSince": "2026-11-05T10:31:55+05:30", "programName": "O1234_FLANGE" }
```

A code never read is left out of `values`; `at` is null and the status fields are left out until the collector reported the machine (e.g. right after a server restart).

## Server → browser

| Event | When | Payload |
|---|---|---|
| `machine:values` | Codes of the room's set changed, at most once per its interval | `{ machineId, seq, at, values: { code: value } }` — only the changed codes; numbers or text |
| `machine:status` | Mode / online changed (status channel) | `{ machineId, seq, online, modeId, modeSince, programName }`; `online` = mode ≠ Offline |
| `machine:event` | Alarm, E-stop or cycle (event channel) | `{ machineId, seq, type, code?, message?, at }`; `type` = `alarm_on`, `alarm_off`, `emergency_on`, `emergency_off`, `cycle_opened`, `cycle_closed` |
| `access:changed` | The user's machines changed | `{ removed: [ids], added: [ids] }` — removed machines' rooms were already left |
| `session:expired` | The session ended (logout, expiry, user deactivated / deleted) | `{}`, then the socket is closed |

## Ordering

`seq` counts every change of a machine (values, status, events). Ignore any message whose `seq` is not above the `seq` of the machine's latest snapshot; deltas after it are always newer. Nothing is replayed after a reconnect: the client subscribes again and the new snapshot replaces its state.

## Rooms (server side)

| Room | Members | Carries |
|---|---|---|
| `m:<machineId>:v:<hash of sorted codes>@<interval ms>` | Sockets that asked for exactly these codes of the machine at this interval | `machine:values` |
| `m:<machineId>:status` / `m:<machineId>:event` | Sockets that asked for the channel | `machine:status` / `machine:event` |
| `u:<userId>` | All sockets of a user | `access:changed` |
| `s:<session hash>` | All sockets of one session | `session:expired` |

## Settings

`LIVE_ENABLED` (true), `LIVE_MIN_INTERVAL_MS` (200), `LIVE_REVALIDATE_SEC` (60), `LIVE_MAX_SOCKETS_PER_USER` (10) — `internal/config/config.go`.

## Adding a variable to a screen

A variable the collector already reads (a `snapshot` address in `cnc_variable_addresses` for the machine's engine): add its code to the screen's list in the client's `constants/live.js` and draw it — no server change. A new `cnc_variables` row the DLL already sends is picked up by the collector through `cnc_config_changed`; PLC / NC addresses need a DLL read function first. Details: plan § 5.5 (`../../docs/socketio-live-data-plan.md`).
