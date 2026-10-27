# CNC Data Collection Flow

This document explains, step by step, how the server finds the active CNC machines, connects to
them and stores their data. Every step names the **file** and the **function** that does the work.

> **Tip:** every step has a marker comment in the code. Search the project for
> `DATA_COLLECTION_STEP_` (or for one step, e.g. `DATA_COLLECTION_STEP_8`) to jump straight to it.

---

## 1. The big picture

- Data collection runs **inside the API server process**. There is no separate service.
  Starting the server also starts the collection.
- The server talks to the machines through a vendor DLL: `Lib\FanucCncCommunication.dll`.
  One DLL handles every controller type (Fanuc, Siemens, …).
- The DLL is **32-bit**, so the server must be built with `GOARCH=386`
  (`.\run.ps1`, `.\build.ps1` and `go tool air` do this).
- Each machine gets its own **worker** (a goroutine). One machine failing never affects the others.

```
cmd/server/main.go                 Server start
        │
        ▼
internal/collector/supervisor.go   Supervisor  ── loads the DLL, keeps the collector running
        │
        ▼
internal/monitor/collector.go      Collector   ── finds active machines, starts one worker each,
        │                                         follows database changes
        ▼
internal/monitor/worker.go         Worker (1 per machine) ── connects, reads, stores, reconnects
        │
        ▼
internal/cnc/cnc.go                Driver      ── Go wrapper around the DLL
internal/cnc/dll_windows_386.go
        │
        ▼
FanucCncCommunication.dll          Vendor DLL  ── talks to the real CNC controller

Database access goes through:
internal/service/machine_telemetry_service.go  →  internal/repository/machine_telemetry_repository.go
```

### Main parts

| Part | File | Job |
|---|---|---|
| Supervisor | `internal/collector/supervisor.go` | Starts collection when the server starts, restarts the collector if it fails to start, stops everything on shutdown |
| Collector | `internal/monitor/collector.go` | Loads active machines, keeps exactly one worker per active machine, listens for database changes |
| Worker | `internal/monitor/worker.go` | Handles one machine: connect → read → store → reconnect |
| Mode decision | `internal/monitor/machine_mode.go` | Decides the machine mode (Production, Hold, …) from the NC state of each reading |
| Driver | `internal/cnc/cnc.go`, `internal/cnc/dll_windows_386.go` | Loads the DLL, connects machines, receives the DLL callbacks |
| Service | `internal/service/machine_telemetry_service.go` | Business rules and transactions for collection |
| Repository | `internal/repository/machine_telemetry_repository.go` | All SQL for collection |

---

## 2. All steps at a glance

| Step | File | Function | What it does |
|---|---|---|---|
| 1 | `cmd/server/main.go` | `run` | Starts CNC collection before the HTTP server starts |
| 2 | `internal/collector/supervisor.go` | `StartCNCDataCollection` | Checks `CNC_ENABLED` |
| 3 | `internal/cnc/cnc.go` → `dll_windows_386.go` | `LoadDriver` → `loadDLL` | Loads the DLL and creates callback slots |
| 4 | `internal/collector/supervisor.go` | `StartCNCDataCollection` → `runCollectorWithRetry` | Runs the collector in the background, retries if it cannot start |
| 5 | `internal/monitor/collector.go` | `RunCollection` | Prepares the database: every machine starts offline |
| 6 | `internal/monitor/collector.go` | `RunCollection` | Records `collector_start`, starts the sync loop |
| 7 | `internal/monitor/collector.go` | `runDatabaseSynchronizationLoop` | Listens for database changes (LISTEN) |
| 8 | `internal/monitor/collector.go` | `syncWorkersWithActiveMachines` | **Loads all active machines from the database** |
| 9 | `internal/monitor/collector.go` | `startOrUpdateWorker` | Starts one worker per active machine |
| 10 | `internal/monitor/worker.go` | `runConnectionLoop` → `loadVariableConfig` | Loads which values to read for this machine |
| 11 | `internal/monitor/worker.go` | `connectWithinLimit` | **Connects the machine** (limited number at a time) |
| 12 | `internal/cnc/cnc.go` → `dll_windows_386.go` | `Driver.Connect` → `connectMachine` → `callCNCConnect` | Takes a callback slot and calls the DLL's `CNCConnect` |
| 13 | `internal/cnc/dll_windows_386.go` | `slot.onData` | The DLL sends a reading; it is handed to the worker |
| 14 | `internal/monitor/worker.go` | `collectReadings` | **Data collection loop**: waits for readings, reloads, ticks |
| 15 | `internal/monitor/worker.go` | `processReading` | Stores the reading, works out alarms and mode |
| 16 | `internal/monitor/worker.go` | `collectReadings` (ticker) | Detects a dead connection, saves the live status |
| 17 | `internal/monitor/worker.go` | `runConnectionLoop` → `recordDisconnection` | Disconnects, marks the machine offline, retries |
| 18 | `internal/monitor/collector.go` | `handleChangeNotification` | Reacts when a machine or variable is edited |
| 19 | `cmd/server/main.go` | `run` (defer) → `Supervisor.Stop` | Stops collection on server shutdown |
| 20 | `internal/monitor/collector.go` | `RunCollection` (end) | Stops workers, marks machines offline, records `collector_stop` |

Steps 1–4 run once at server start. Steps 5–9 run each time the collector starts. Steps 10–17 run
for **every machine**, again and again while the server is up. Step 18 runs whenever someone edits
data. Steps 19–20 run at shutdown.

---

## 3. Startup: steps 1 – 4

### Step 1 — Start collection with the server
- **File:** `cmd/server/main.go`
- **Function:** `run`
- **What it does:**
  1. Loads the config, connects to PostgreSQL and builds the HTTP routes (normal API startup).
  2. Creates the supervisor: `collector.NewSupervisor(appConfig.CNC, dbPool)`.
  3. Calls `collectors.StartCNCDataCollection()`.
  4. Starts the HTTP server.
- **Good to know:** if `StartCNCDataCollection` returns an error (wrong build), the server does not start.

### Step 2 — Check that collection is enabled
- **File:** `internal/collector/supervisor.go`
- **Function:** `Supervisor.StartCNCDataCollection`
- **What it does:** reads `CNC_ENABLED`. If it is `false`, it logs
  `cnc data collection disabled` and returns. The API keeps working without collection.

### Step 3 — Load the DLL
- **Files / functions:**
  - `internal/cnc/cnc.go` → `LoadDriver`
  - `internal/cnc/dll_windows_386.go` → `loadDLL`, `findLibDir`, `setDllDirectory`
  - `internal/cnc/dll_other.go` → `loadDLL` (non-386 builds: always returns `ErrNotSupported`)
- **What it does:**
  1. Checks that the Go snapshot struct is exactly **34,440 bytes**, the same size as the DLL's C++ struct.
  2. Finds the `Lib` folder: first next to the `.exe`, then in the working directory (`findLibDir`).
  3. Tells Windows to look for the DLL's own dependencies in that folder (`setDllDirectory`).
  4. Loads `FanucCncCommunication.dll` and finds two functions in it: `CNCConnect` and `CNCdisconnect`.
  5. Creates **callback slots**: `CNC_SLOTS` of them (default 300, max 990). Each slot is a pair of
     Go functions the DLL calls: one for data (`onData`), one for errors (`onError`).
- **What can happen:**

  | Result | What the server does |
  |---|---|
  | Build is not 32-bit (`cnc.ErrNotSupported`) | Returns `ErrWrongBuild` → **server does not start** |
  | DLL file missing or broken | Logs `cnc dll unavailable` → API runs, **no collection** |
  | OK | Continues to step 4 |

- **Why slots?** The DLL's callback does not say which machine sent the data. So every connected
  machine gets its own slot, and the slot knows its machine. Go allows only about 2,000 callbacks per
  process (2 per slot), so the maximum is 990 slots = 990 machines.

### Step 4 — Run the collector in the background
- **File:** `internal/collector/supervisor.go`
- **Functions:** `Supervisor.StartCNCDataCollection`, `Supervisor.runCollectorWithRetry`
- **What it does:**
  1. Creates `MachineTelemetryService` (all database work for collection).
  2. Starts a goroutine that runs `runCollectorWithRetry` and returns at once, so the HTTP server
     does not wait for the database or the machines.
  3. `runCollectorWithRetry` creates a `monitor.Collector` (shard 0 of 1 = all machines) and calls
     `collector.RunCollection(ctx)`.
  4. If `RunCollection` fails to start (for example the database is not ready), it waits and tries
     again: **1 s, 2 s, 4 s … up to 1 minute** between attempts.

---

## 4. Find the machines: steps 5 – 9

### Step 5 — Prepare the database
- **File:** `internal/monitor/collector.go`
- **Function:** `Collector.RunCollection`
- **Calls:** `MachineTelemetryService.IsHypertable`, `MachineTelemetryService.InitializeShardForCollection`
  → `MachineTelemetryRepository.MarkShardOffline`
- **What it does:**
  1. Checks that `machine_telementry` is a TimescaleDB hypertable. If not, it only logs an error
     (fine for development, too slow for production).
  2. Marks **every machine offline** in one transaction:
     - a mode period left open by the last run is closed at the machine's `last_seen_at`
       (so the time the server was down counts as offline),
     - an Offline period is opened in `machine_mode_log`,
     - `machine_live_status` and `machines.online_status` say offline, and the counters restart.
- **Why:** a machine only shows as online after its worker really connects to it.

### Step 6 — Record the start and begin following the database
- **File:** `internal/monitor/collector.go`
- **Functions:** `Collector.RunCollection`, `Collector.recordCollectorEvent`
- **What it does:**
  1. Stores a `collector_start` event in `machine_events`.
  2. Starts a goroutine that runs `runDatabaseSynchronizationLoop` (step 7).
  3. Then `RunCollection` just **waits** until the server shuts down (step 20).

### Step 7 — Listen for database changes
- **File:** `internal/monitor/collector.go`
- **Function:** `Collector.runDatabaseSynchronizationLoop`
- **Calls:** `MachineTelemetryService.ListenForConfigurationChanges` → `db.ListenForDatabaseNotifications`
  (`internal/db/postgres.go`)
- **What it does:**
  1. Opens one database connection and runs `LISTEN machines_changed` and `LISTEN cnc_config_changed`.
  2. **Once listening works**, it loads all machines (step 8).
  3. For every notification after that, it calls `handleChangeNotification` (step 18).
  4. If the connection is lost, it waits (1 s, 2 s … up to 1 min) and starts again from 1.
- **Why load the machines only after LISTEN works?** A change made between "load machines" and
  "start listening" would otherwise be missed. This also covers reconnects: after every reconnect the
  full list is loaded again.

### Step 8 — Load all active machines
- **File:** `internal/monitor/collector.go`
- **Function:** `Collector.syncWorkersWithActiveMachines`
- **Calls:** `MachineTelemetryService.GetActiveMachines` → `MachineTelemetryRepository.FindActiveMachines`
- **What it does:** reads the machines to collect:

  ```sql
  SELECT m.id, m.name, m.live_ip, m.port, username, password, controller_version_type, c.cnc_engine
  FROM machines m
  LEFT JOIN controllers c ON c.id = m.controller_id AND c.is_deleted = false
  WHERE m.is_deleted = false AND m.status = 'active'
    AND m.id % shards = shard          -- today: shard 0 of 1 = every machine
  ORDER BY m.id
  ```

  The protocol to use (`cnc_engine`) comes from the machine's **controller**.

### Step 9 — Start one worker per machine
- **File:** `internal/monitor/collector.go`
- **Functions:** `Collector.syncWorkersWithActiveMachines`, `Collector.startOrUpdateWorker`,
  `Collector.stopWorker`, `isMachineUnchanged`
- **What it does:** for each active machine, `startOrUpdateWorker` decides:

  | Situation | Action |
  |---|---|
  | No worker yet | Create one (`newMachineWorker`) and start `runConnectionLoop` in a goroutine |
  | Worker exists, nothing changed | Do nothing |
  | Only name or controller version changed | Keep the connection; give the worker the new details and ask it to reload its variables (`queueMachineUpdate`) |
  | IP, port, login or engine changed | Stop the old worker, wait until it has disconnected, start a new one |

  Then every worker whose machine is **no longer active** is stopped (`stopWorker`).

---

## 5. Connect and collect, per machine: steps 10 – 17

Each worker runs these steps in its own goroutine. Only that goroutine changes the worker's state,
so the worker needs no locks for its own data.

### Step 10 — Get ready to connect
- **File:** `internal/monitor/worker.go`
- **Functions:** `worker.runConnectionLoop`, `worker.loadVariableConfig`, `worker.applyPendingMachineUpdate`
- **Calls:** `MachineTelemetryService.GetVariableConfig` → `MachineTelemetryRepository.FindVariableConfig`
- **What it does:**
  1. If the machine's controller has **no `cnc_engine`**, it stores a `config_error` event and waits
     without connecting (until the machine is changed or the server stops).
  2. Before every connect attempt it loads the machine's **variable list**: which values to read and
     where they are in the DLL snapshot (for example variable 14 = `MachineMode`).
     - The address for the machine's controller version is preferred over the general one.
     - Per-machine overrides in `machine_variables` are applied.
     - Today only `snapshot` addresses can be read; `plc` / `nc` need a new DLL function.
  3. Each address is matched to a snapshot field (`cnc.FieldIndex` in `internal/cnc/fields.go`).
     An unknown address gives a `config_error` event and that variable is skipped.
  4. Logs `cnc: variables loaded` with the map it built (`"14:opMode<-MachineMode"`, …). Active
     variables with no address the server can read (today: everything that is not a DLL snapshot
     field, such as the Siemens node ids) are listed once in `cnc: active variables not readable`.

### Step 11 — Connect, a few machines at a time
- **File:** `internal/monitor/worker.go`
- **Function:** `worker.connectWithinLimit`
- **What it does:** waits for a free place in `connectLimiter` (at most `CNC_MAX_CONNECTING`,
  default 20, connect calls at the same time), then calls `driver.Connect`.
- **Why:** `CNCConnect` can block while the DLL connects. Without a limit, a server start with many
  machines would try to open every connection at once.

### Step 12 — Call the DLL's CNCConnect
- **Files / functions:**
  - `internal/cnc/cnc.go` → `Driver.Connect` (returns `ErrNoEngine` if no engine)
  - `internal/cnc/dll_windows_386.go` → `dllDriver.connectMachine`, `dllDriver.callCNCConnect`
- **What it does:**
  1. Takes a free callback slot. If none is left → `ErrNoFreeSlot` (raise `CNC_SLOTS`).
  2. Links a new `Connection` to the slot **before** calling the DLL, because the DLL may send data
     before `CNCConnect` even returns.
  3. For **Siemens** only: sets `SIEMENS_USER` / `SIEMENS_PASS` in the process environment (the DLL
     reads the login from there). A lock makes Siemens connects run one at a time so two machines do
     not overwrite each other's login.
  4. Calls `CNCConnect(ip, port, dataCallback, errorCallback, engine)`. The DLL starts its own
     thread that polls the machine.
  5. A negative result means the connect failed: the slot is freed again.
- **Result in the database:**

  | Result | Event in `machine_events` |
  |---|---|
  | Connected | `connect_ok` ("connection 5 (192.168.1.10:8193, engine 1)") |
  | Connect failed | `connect_failed` → retry later (step 17) |
  | No free slot | `config_error` → retry later |

  After `connect_ok` the log shows `cnc: machine subscribed` (machine, IP, port, engine,
  connection id, number of variables).

  Engine values (`controllers.cnc_engine`): 0 Simulator, 1 Fanuc, 2 Siemens, 3 Syntec, 4 Mitsubishi.

### Step 13 — The DLL sends data
- **File:** `internal/cnc/dll_windows_386.go`
- **Functions:** `slot.onData`, `slot.onError`
- **Then:** `internal/cnc/cnc.go` → `Connection.publishReading` / `publishError` → `signalWorker`
- **What it does:**
  1. On every poll the DLL calls `onData` on **its own thread** with a pointer to the snapshot.
  2. `onData` copies only the fields the server uses (never the 32 KB program list) into a `Reading`:
     mode, status, program, spindle, feed, axis loads, tool, alarm, and the **axis positions** of
     5 axes (`MachPos`, `AbsPos`, `RelPos`, `DistToGo`, named `MachPos[0]` … in `internal/cnc/fields.go`).
     Migration 029 stores X, Y and Z (`[0]`, `[1]`, `[2]`) as variables 208–219.
  3. `publishReading` saves it as the connection's newest reading and wakes the worker.
  4. Errors from the DLL go through `onError` → `publishError` (counted in the connection stats).
- **Rules for this code:** it must never block and never panic (a panic would go through the DLL and
  kill the whole server), so panics are caught. If the worker is busy, several wake-ups collapse into
  one; the worker simply takes the newest reading.

### Step 14 — The data collection loop
- **File:** `internal/monitor/worker.go`
- **Function:** `worker.collectReadings`
- **What it does:** waits for one of four things, again and again:

  | Event | Action |
  |---|---|
  | Server stops / worker stopped | Leave the loop (`stopRequested`) |
  | Variable reload requested | Load the variable list again (`loadVariableConfig`), no reconnect |
  | DLL sent something | Store any new DLL error (`recordDLLError`), take the newest reading, skip it if it was already processed, then `processReading` (step 15) |
  | 1-second tick | Stale check and live status (step 16) |

### Step 15 — Store the reading
- **File:** `internal/monitor/worker.go`
- **Functions:** `worker.processReading`, `upsertLiveValues`, `insertTelemetry`,
  `recordAlarmChanges`, `recordModeChange`, and `decideMachineMode` in `internal/monitor/machine_mode.go`
- **What it does:**
  1. If the DLL says the machine is **not connected**: set mode Offline and stop here.
  2. The first time (or after a variable reload) load the current values from
     `machine_variable_values` into memory (`GetLiveValues`).
  3. For every variable:
     - **Text** variables (program name, alarm message) are **not stored**; they are only used for
       the mode log and alarm events.
     - **Numeric** variables are compared with the cached value.
  4. Every numeric value is **upserted** into `machine_variable_values` (one row per machine and variable).
  5. Only values that **changed** are inserted into `machine_telementry` (time series).
     With `CNC_LOG_VALUES=true` each changed value is also logged as `cnc: variable value changed`
     (machine, variable id and code, address, previous and new value, read time).
  6. If anything changed, it also writes `alarm_on` / `alarm_off` and `emergency_on` /
     `emergency_off` events when they change.
  7. On **every** connected reading it decides the mode from the NC state of the reading
     (`machineStateFromReading` → `decideMachineMode`) and, if it is different, closes the old
     period and opens a new one in `machine_mode_log` (and sets `machines.mode_id`), in one
     transaction. The mode does not depend on the machine's variables, and it is right from the
     first reading after a start or reconnect even when no value changed.
- **NC state read from the snapshot:** opMode = `MachineMode` (0 JOG, 1 MDA, 2 AUTO); program =
  `ProcStatus`, which the DLL reads from `/Channel/State/acProg` (0 reset, 1 stopped, 2 active,
  3 waiting, 4 interrupted; `programStateFromProcStatus`); feedRateOvr = `FeedRatePercetnage`;
  alarmActive = `AlarmNo` ≠ 0 or `IsEmergency`. chanStatus is not in the snapshot, so it is
  unknown and decides nothing until the DLL sends it.
- **How the mode is chosen** (`decideMachineMode`, first match wins; running = program running
  or waiting):

  | Condition | Mode (id) |
  |---|---|
  | No NC connection | Offline (9) |
  | Alarm active and not running | Downtime (6) |
  | AUTO, running, chanStatus active (or unknown), feed override > 0, no alarm | In Production (3) |
  | Program interrupted / stopped, chanStatus interrupted, or running with feed override ≤ 0, an alarm or not in AUTO | Hold (4) |
  | Anything else (reset, not running, JOG) | Idle (5) |

  Setup (1) and Rapid Mode (12) are no longer decided.

- **Log:** each change is logged once as `cnc: mode changed` with `from`, `to`, `reason` (the rule
  that matched, e.g. `program running, feed override ≤ 0`) and the NC inputs `opMode`,
  `procStatus` (raw value and state), `chanStatus`, `feedOverride`, `alarmNo`, `emergency`,
  `ncProgram`, so a mode can be checked against the HMI. It is always on; `CNC_LOG_VALUES` only
  controls the per-reading `cnc: variable values changed` line.

- **If a write fails:** the error is logged and saved in `machine_live_status.last_error`. Failed
  telemetry rows are lost and counted in `dropped_count`; the next reading is stored normally.

### Step 16 — Every second: dead connection check and live status
- **File:** `internal/monitor/worker.go`
- **Functions:** `worker.collectReadings` (ticker), `worker.saveLiveStatus`
- **What it does:**
  1. If there was **no connected reading for `CNC_STALE_SEC`** (default 30 s), the connection is
     treated as dead: leave the loop (`stopConnectionStale`) → step 17.
  2. Writes the `machine_live_status` row when something changed, or at least every 10 seconds
     (connected since, last seen, counters, last error, missing nodes).
  3. Every 10 seconds logs `cnc: reading rate` (`logReadingRate`): readings the DLL sent
     (`callbacks`, `callbacksPerSecond`), readings the worker handled (`processed`), readings it
     had to skip because it was still busy with the previous one (`skipped`), and values stored
     (`storedValues`, `storedValuesPerSecond`). A growing `skipped` means the worker, mostly its
     database writes, is slower than the machine.

### Step 17 — Disconnect and try again
- **File:** `internal/monitor/worker.go`
- **Functions:** `worker.runConnectionLoop`, `worker.recordDisconnection`
- **Calls:** `Connection.Disconnect` → DLL `CNCdisconnect`, then the slot is freed (`releaseSlot`)
- **What it does:**
  1. Disconnects. From this moment late callbacks from the DLL are ignored.
  2. Stores a `stale` event (if it went stale) and a `disconnected` event.
  3. Sets mode Offline and saves the live status.
  4. Waits, then goes back to step 10 and stores a `reconnect` event ("attempt 2", …).
- **Retry wait:** 1 s, 2 s, 4 s … up to **5 minutes**, with ±20 % random jitter so machines that
  failed together do not all retry at the same moment. The wait goes back to 1 s once a connection
  delivered real data.

---

## 6. When data is edited: step 18

### Step 18 — React to machine and variable changes
- **File:** `internal/monitor/collector.go`
- **Function:** `Collector.handleChangeNotification`
- **Where notifications come from:** database triggers in `migrations/028_cnc_machine_telemetry.sql`

| Channel | Sent when | Payload | What the collector does |
|---|---|---|---|
| `machines_changed` | A machine is added or deleted, or one of `is_deleted, status, live_ip, port, username, password, brand_id, controller_id, controller_version_type` changes | machine id | Reads the machine again (`GetActiveMachine`). Not active any more → `stopWorker`. Otherwise → `startOrUpdateWorker` (step 9 rules) |
| `cnc_config_changed` | `cnc_variables`, `cnc_variable_addresses` or `machine_variables` change | `machine:<id>` or `variables` | Asks that one worker, or all workers, to reload their variable list (`requestVariableReload`). No reconnect |

So when a user edits a machine on the website, collection follows within moments, without polling.

---

## 7. Shutdown: steps 19 – 20

### Step 19 — Stop collection with the server
- **Files / functions:** `cmd/server/main.go` → `run` (deferred call), `internal/collector/supervisor.go` → `Supervisor.Stop`
- **Order on Ctrl+C / SIGTERM:**
  1. The HTTP server shuts down (`server.Shutdown`).
  2. `Supervisor.Stop` cancels the collection and waits at most `CNC_STOP_TIMEOUT_SEC` + 5 s
     (so a machine stuck inside the DLL cannot hang the shutdown).
  3. The database pool closes last, so the collector can still write "offline".

### Step 20 — Stop the workers and mark machines offline
- **File:** `internal/monitor/collector.go`
- **Functions:** `Collector.RunCollection` (end), `waitForWorkersToStop`, `recordCollectorEvent`
- **Calls:** `MachineTelemetryService.FinalizeShardCollection` → `MarkShardOffline`
- **What it does:**
  1. Stops the sync loop and every worker. Each worker disconnects and records its machine offline (step 17).
  2. Waits for the workers, at most `CNC_STOP_TIMEOUT_SEC` in total.
  3. Marks every machine offline (counters are kept, to help find problems).
  4. Stores a `collector_stop` event.

---

## 8. Reference

### Log lines (in `logs/app.log`, search for the message)

| Message | Level | When | Step |
|---|---|---|---|
| `cnc dll loaded` / `cnc data collection started` | Info | Server start | 3, 4 |
| `cnc collector started` | Info | Collector run begins | 6 |
| `cnc collector: machines loaded` | Info | Machine list loaded (count + `machineIds`) | 8 |
| `cnc collector: worker started` | Info | A worker starts for a machine | 9 |
| `cnc: variables loaded` | Info | Variable list and its variable → field map | 10 |
| `cnc: active variables not readable` | Warn | Active variables the server cannot read | 10 |
| `cnc machine event` (`connect_ok`, `connect_failed`, `reconnect`, …) | Info / Warn | Every machine event | 12, 17 |
| `cnc: machine subscribed` | Info | Connected (IP, port, engine, connection id) | 12 |
| `cnc: controller data started` / `stopped` | Info / Warn | The DLL starts / stops sending connected data | 15 |
| `cnc: variable value changed` | Info | Each changed value (`CNC_LOG_VALUES=true` only) | 15 |
| `cnc: mode changed` | Info | Mode changes (from, to, program) | 15 |
| `cnc: reading rate` | Info | Every 10 s per machine | 16 |
| `cnc collector: processing machine change notification` | Info | A machine was edited | 18 |
| `cnc collector: variable config changed, reloading variables` | Info | Variables were edited | 18 |
| `cnc collector: machine no longer collected` | Info | A worker stops | 9, 18 |
| `cnc collector stopping` / `stopped`, `cnc data collection stopped` | Info | Shutdown | 19, 20 |

### Settings (`.env`, defaults in `internal/config/config.go`)

| Setting | Default | Meaning |
|---|---|---|
| `CNC_ENABLED` | `true` | `false` = API only, no collection |
| `CNC_SLOTS` | `300` | Most machines the server can connect (max 990) |
| `CNC_LIB_DIR` | `Lib` | Folder with the DLL |
| `CNC_DLL_NAME` | `FanucCncCommunication.dll` | DLL file name |
| `CNC_STALE_SEC` | `30` | Seconds without data before a connection counts as dead |
| `CNC_MAX_CONNECTING` | `20` | Connect calls allowed at the same time |
| `CNC_STOP_TIMEOUT_SEC` | `15` | Longest wait for a clean stop |
| `CNC_LOG_VALUES` | `false` | `true` logs every changed value (step 15). For analysis only: it writes many lines per second |
| `SIEMENS_USER`, `SIEMENS_PASS` | – | Set per machine by the server from the machine's login (step 12) |

### Tables

| Table | Read / written | Content |
|---|---|---|
| `machines`, `controllers` | Read | Which machines to connect, and how (step 8) |
| `cnc_variables`, `cnc_variable_addresses`, `machine_variables` | Read | Which values to read and where (step 10) |
| `machine_variable_values` | Written every reading | Latest numeric value per machine and variable (step 15) |
| `machine_telementry` | Written on change | Time series of changed numeric values (TimescaleDB, compressed after 1 day, kept 365 days) |
| `machine_mode_log` | Written on mode change | One row per mode period; `ended_at` NULL = current period |
| `machine_live_status` | Written on change / every 10 s | One row per machine: connected, last seen, counters, last error |
| `machine_events` | Written on events | Connect, disconnect, alarms, errors, collector start / stop |

### Events in `machine_events`

| Event | When | Step |
|---|---|---|
| `collector_start` / `collector_stop` | Collector starts / stops (`machine_id` is NULL) | 6 / 20 |
| `config_error` | No `cnc_engine`, unknown variable address, or no free callback slot | 10, 12 |
| `connect_ok` / `connect_failed` | Result of `CNCConnect` | 12 |
| `dll_error` | The DLL reports a new error message | 14 |
| `alarm_on` / `alarm_off` | Alarm number changes | 15 |
| `emergency_on` / `emergency_off` | Emergency stop changes | 15 |
| `stale` | No data for `CNC_STALE_SEC` | 17 |
| `disconnected` | A connection ended | 17 |
| `reconnect` | A new connect attempt starts | 17 |

`values_dropped` exists in `internal/model/machine_telemetry.go`, but nothing writes it. Lost rows
are counted in `machine_live_status.dropped_count` instead.

### Retry and timing summary

| What | Where | Wait |
|---|---|---|
| Collector cannot start | `runCollectorWithRetry` | 1 s doubling to 1 min |
| LISTEN connection lost | `runDatabaseSynchronizationLoop` | 1 s doubling to 1 min, ±20 % |
| Machine connect failed / connection lost | `runConnectionLoop` | 1 s doubling to 5 min, ±20 %; back to 1 s after real data |
| Dead connection check | `collectReadings` | every 1 s, limit `CNC_STALE_SEC` |
| Live status save | `collectReadings` | on change, at least every 10 s |
| One database write | `worker.go` (`writeTimeout`) | 5 s |

### Troubleshooting

| Symptom | Look at |
|---|---|
| Server does not start: "needs a 32-bit build" | Start with `.\run.ps1` / `go tool air`, or set `CNC_ENABLED=false` (step 3) |
| Log: `cnc dll unavailable` | DLL or its dependencies missing in `Lib\` (step 3) |
| Machine never connects, `config_error` "no cnc_engine" | Set `controllers.cnc_engine` and save the machine (step 10) |
| `config_error` "no free callback slot" | Raise `CNC_SLOTS` (max 990) (step 12) |
| Many `connect_failed` | Machine IP / port / network; check `machine_events.message` (step 12) |
| Machine keeps going `stale` | Machine is reachable but sends no connected data (step 16) |
| A value is never stored | Its address does not match a snapshot field (see `config_error`), it is a text variable, or it has no variable in `cnc_variables` (steps 10, 15) |
| Fewer rows than the machine changes | Compare `callbacksPerSecond` and `storedValuesPerSecond` in the `cnc: reading rate` log: only changed values of configured variables are stored, and `skipped` shows readings the worker could not keep up with (step 16) |
| Whole server stops suddenly | A crash inside the DLL ends the process; check `logs/app.log` |

### Known issue
In `worker.processReading`, if loading the cached values (`GetLiveValues`) fails, the cache stays
`nil`; a following successful `upsertLiveValues` then writes into a `nil` map and panics, which ends
that worker. The fix is to start with an empty map when the load fails.

### Related documents
- `docs/cnc-communication.md` — the DLL API (`CNCConnect`, callbacks, snapshot layout)

---

## 9. Data Reading Flow (Function Names Only)

**From server start to data storage:**

```
run (main.go)
  → NewSupervisor (supervisor.go)
  → StartCNCDataCollection (supervisor.go)
  → runCollectorWithRetry (supervisor.go)
  → RunCollection (collector.go)
  → runDatabaseSynchronizationLoop (collector.go)
  → syncWorkersWithActiveMachines (collector.go)
  → startOrUpdateWorker (collector.go)
  → runConnectionLoop (worker.go)
  → loadVariableConfig (worker.go)
  → connectWithinLimit (worker.go)
  → driver.Connect (cnc.go)
  → dllDriver.connectMachine (dll_windows_386.go)
  → callCNCConnect (dll_windows_386.go)
  → CNCConnect (DLL function)
  → onData (dll_windows_386.go - DLL callback)
  → publishReading (dll_windows_386.go)
  → collectReadings (worker.go)
  → Latest (worker.go)
  → processReading (worker.go)
  → upsertLiveValues (worker.go)
  → insertTelemetry (worker.go)
```

**Key functions in the data collection loop (per machine):**

```
runConnectionLoop (worker.go)
  → collectReadings (worker.go)
    → connection.Signal() (waits for DLL data)
    → Latest (worker.go)
    → processReading (worker.go)
      → upsertLiveValues (worker.go)
      → insertTelemetry (worker.go)
      → recordAlarmChanges (worker.go)
      → recordModeChange (worker.go)
```
