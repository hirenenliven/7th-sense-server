# CNC DLL Integration Plan

Plan to read live machine data through `FanucCncCommunication.dll` (reference client:
`docs/dll-usecase/G0_V0.0.4_DLL/sinumerik_client.go`) and store it for the API.
Status: plan only, nothing implemented yet.

## DLL limits that shape the design

1. **The DLL is 32-bit.** It is built for x86 (`GOARCH=386`, `Fwlib32.dll`, snapshot struct
   34,440 bytes with `pack(8)` x86 layout). The API server builds as 64-bit and cannot load it,
   so the DLL runs in a **separate 386 collector program**. A native crash then cannot take down
   the API.
2. **The callback does not say which machine sent the data.** `onData(snapshotPtr)` receives only
   the snapshot, and the snapshot has no machine ID. With many connections in one process, each
   machine needs its own callback pair. Go creates at most about 2,000 callbacks per process and
   never frees them (about 1,000 machines at 2 per machine). So: a **fixed set of reusable slots**
   per process, and **machines split across several collector processes** (for example 100–200
   machines each).
3. **Login details are shared by the whole process.** The DLL reads them from the `SIEMENS_USER`
   and `SIEMENS_PASS` environment variables, and one process has only one set. Connects must run
   one at a time under a lock: set the environment, then connect. **Ask the DLL vendor to
   confirm** the variables are read only during `CNCConnect` and never on the DLL's own internal
   reconnects.
   - Best fix: a vendor function such as `CNCConnectEx(ip, port, user, pass, type, userData)`
     that passes `userData` back in the callbacks. That removes limits 2 and 3, and the rest of
     this design stays the same.

## How data flows

```
cmd/collector (windows/386, one process per shard)
  load monitored machines → one worker per machine → DLL callback → latest snapshot in its slot
  → flusher every N seconds → service → repository → PostgreSQL / TimescaleDB
API server (64-bit, unchanged) → reads the live tables through the existing machine access rules
```

Both programs live in the same Go module and share `config`, `db`, `model` and `repository`.
The layer order stays: collector → service → repository.

## Steps

### 0. Check the DLL behaviour on the test machine (no code)

- How the DLL reports a dropped connection: `IsConnected`, the error callback, or the callbacks
  simply stopping.
- Whether `CNCConnect` blocks, and how often callbacks arrive.
- How limit 3 behaves with two connections open at once.
- The `machineType` code for Fanuc (`2` is Siemens).

### 1. Migration `027_add_machine_telemetry.sql`

- `machine_live_status`: one row per machine, updated in place — connected/online state, mode,
  program status, program name, part count, alarm, spindle/feed values, `last_seen_at`. Keeps the
  frequent writes off the `machines` table.
- `machine_telemetry`: TimescaleDB hypertable keyed by `(machine_id, ts)` with the numeric fields.
  Rows are written when values change or every N seconds, not on every callback. Add retention
  and compression policies.
- `machine_events`: connect, disconnect and alarm events, for uptime / downtime later.

### 2. `internal/cnc/` — the DLL wrapper only

Files use `//go:build windows && 386`, with a stub file for other builds so `go build ./...` and
the tests still pass for the server.

- `dll.go`: copy `initDLL`, `findLibDir`, `setDllDirectory` as they are. Wrap
  `Connect(ip, port, user, pass, type)` (serial, under the lock from limit 3) and `Disconnect(id)`.
- `snapshot.go`: copy `CncMasterSnapshot`, the size check, `cstr`, and the mode/status name maps.
  Add `ToLiveStatus()` to convert a snapshot to the model.
- `slots.go`: a fixed set of reusable callback pairs, created once at startup. A slot is assigned
  to a machine and freed when it disconnects.
  - The callback copies the snapshot into the slot (mutex or atomic pointer) and never blocks,
    like `onData` today.
  - It recovers from panics.
  - It ignores callbacks that arrive after the slot is freed (generation counter).

### 3. `internal/collector/` — managing the machines

- `manager.go`: every 60 s, loads machines with `is_monitor = true AND is_deleted = false` in this
  process's shard (`id % shards == shard`). Starts new workers, stops removed ones, restarts a
  worker when the machine's IP, port or login changed.
- `worker.go`: per-machine state (Disconnected → Connecting → Connected → Stale).
  - **Stale** = no callback for T seconds, or `IsConnected = 0`.
  - On failure: disconnect, then reconnect with backoff and jitter (1 s up to 5 min).
  - A limit on concurrent connects avoids a connection flood after a restart.
  - A failing machine affects only itself.
- `flusher.go`: one goroutine reads every slot's latest snapshot on a timer. Writes
  `machine_telemetry` in one batch with pgx `CopyFrom` and updates `machine_live_status` in one
  `INSERT … ON CONFLICT`. DB write cost stays roughly the same for 1 or 1,000 machines.

### 4. Data layer (user CRUD pattern)

- `model/machine_telemetry.go`
- `repository/machine_telemetry_repository.go`: `ListMonitored`, `CopyTelemetry`,
  `UpsertLiveStatus`, `InsertEvents`
- `service/machine_telemetry_service.go`

### 5. `cmd/collector/main.go`

- Loads config. New environment variables: `COLLECTOR_SHARD`, `COLLECTOR_SHARDS`,
  `COLLECTOR_FLUSH_SEC`, `COLLECTOR_STALE_SEC`, `COLLECTOR_MAX_CONNECTING`, `CNC_LIB_DIR`.
- Opens the DB pool and starts the manager.
- On SIGTERM: disconnects every machine and flushes once more.
- The old `runCheck` becomes `collector -check <machineId>` for setting up a new machine.
- Build: `GOARCH=386 go build -o collector.exe ./cmd/collector`; ship with the `Lib/` folder and
  the certificates.
- Run each shard as a Windows service (NSSM or `sc`) with automatic restart, so a DLL crash
  restarts only that shard.

### 6. API (after data is flowing)

- `GET /machines/{id}/live` and `GET /machines/live`, reading `machine_live_status`, access
  through the existing `MachineService.machineScope`.
- Update `.claude/CLAUDE.md` and run `graphify update .`.

## Scaling path

- **1 test machine:** shard 1 of 1.
- **1,000+ machines:** about 5–10 shard processes on one Windows host, or across hosts. No code
  changes, only `COLLECTOR_SHARD` / `COLLECTOR_SHARDS`.
- **Memory:** about 34 KB per snapshot; 200 machines per 32-bit process is about 7 MB.

## Unchanged

- The API server, its 64-bit build, existing handlers, services and routes, and `machines` CRUD.
- Runtime columns are still not changed through the API.

## Security issues in the sample code

- `sinumerik_client.go` has the PLC username and password hard-coded.
- `client_key.pem` and `client_key.der` are committed under `docs/`.

The collector reads each machine's login from the `machines` table instead. Remove the keys from
git, and change them if this repository is pushed anywhere.
