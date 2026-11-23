# RabbitMQ telemetry queue: step-by-step implementation guide

> **Implemented 2026-11-20** as [ADR 0024](adr/0024-rabbitmq-telemetry-queue.md) (0023 was taken meanwhile; wherever this guide says 0023, read 0024). Differences from this plan are listed there: a database error keeps the batch and retries instead of requeueing, unchanged live values are refreshed every 30 s, migration 039 drops 038's index, and only `RABBITMQ_ENABLED`, `RABBITMQ_URL` and `RABBITMQ_CONSUMERS` are settings (the other `RABBITMQ_*` variables and `CNC_LOG_QUEUE` named below are fixed in code). Steps 9 (diagnostics card) and 5 (later) are not done.

> Hand this file to Claude in `7th-sense-server/`. Follow the steps in order; each step ends with a check. Before starting, follow `.claude/CLAUDE.md` (Graphify first step, layered architecture, coding standards, Documentation Rules). This guide says **what** to build and **where**; the code itself follows the existing patterns of the named files.

## 1. Why, and what today looks like

Goal: 1000s of machines × ~40 variables each must be collected without the database writes slowing the machine workers or losing data.

Today (verified in code, 2026-11-19):

- One goroutine per machine (`internal/monitor/worker.go`, PLC: `plc_worker.go`). For **every reading** it, inside its own loop and synchronously:
  - publishes to the live socket first (`live.Hub.PublishValues`, in memory: fast, keep it),
  - upserts **every numeric value, changed or not** into `machine_variable_values` (`worker.upsertLiveValues` → `MachineTelemetryService.UpsertLiveValues`, `INSERT … unnest … ON CONFLICT`),
  - inserts the **changed** values into `machine_telementry` (`worker.insertTelemetry` → `InsertTelemetry`),
  - plus mode changes, events, cycles and status rows (low volume).
- Each write is its own statement with a 5 s timeout (`writeTimeout`). A failed telemetry insert is dropped (`dropped_count`).
- `db.Connect` uses `pgxpool.New` with defaults: **MaxConns = max(4, CPU count)**. With 1000 workers each writing per reading, the workers queue for a handful of connections, and readings back up behind the database.
- One process holds at most 990 machines (`CNC_SLOTS`, DLL callback slots), and the shard is fixed (`collectorShard`, `collectorShardCount` in `internal/collector/supervisor.go`). Thousands of machines therefore need several collector processes later.

Rough load at 1000 machines × 40 variables × 1 reading/s: ~40,000 live-value upserts/s, mostly of unchanged values, and up to that many telemetry inserts. That is far too many single statements for one pool.

**What RabbitMQ fixes:** workers hand their rows to a durable queue and return immediately. A small number of writer consumers take rows off the queue in **large batches** and write them in a few statements. Bursts are absorbed by the queue rather than by the workers. Data survives a database outage or restart. Later, several collector processes can publish to the same queue.

**What RabbitMQ does not fix by itself:** the write cost. The gain comes from batching and from not upserting unchanged values. Both are part of this plan (Steps 3 and 6).

## 2. Scope

**In the queue (high volume):**
- `machine_telementry` rows (changed numbers);
- `machine_variable_values` rows (latest value);
- from both the CNC worker and the PLC worker.

**Not in the queue (keep direct, as today):** mode changes (`machine_mode_log`, `machines.mode_id`), events (`machine_events`), cycles (`machine_cycles`), `machine_live_status`. They are low volume and order-sensitive. The worker also needs their result straight away: for example, the cycle open is refused with `response.MachineCycleRunning` and the worker reacts to that.

**Unchanged:**
- **Live socket:** the hub is still fed in-process before anything is queued.
- **No API change:** no OpenAPI or Postman change.
- **No client change.** I checked the client: it reads values via the socket and the REST endpoints that query these tables. Data now reaches the tables a few hundred ms later, which no screen depends on. Diagnostics' "last stored" just lags slightly.

**Off switch:** `RABBITMQ_ENABLED=false` must give **exactly today's behaviour** (direct writes). This is the default until it is proven in production.

## 3. Design decisions (put these in the ADR)

| Topic | Decision |
|---|---|
| Library | `github.com/rabbitmq/amqp091-go`: the official client maintained by the RabbitMQ team (`streadway/amqp` is deprecated). Pure Go, builds for `GOARCH=386`. |
| Broker | RabbitMQ 3.13+ (or 4.x) with the management plugin. Dev: Docker `rabbitmq:management` (AMQP 5672, UI 15672). Windows server: the official installer (needs Erlang). |
| Topology | One durable **direct exchange** `telemetry`, routing key `values`, bound to one durable **quorum queue** `telemetry.values` (survives broker restarts, replicated when clustered). Its dead-letter exchange `telemetry.dlx` routes to queue `telemetry.values.dead`. The app declares all of this at start; declaring is idempotent. |
| Message | **One message per reading of one machine**, not one per value: `{machineId, readAt, source: "cnc" or "plc", live: [{variableId, value}], telemetry: [{variableId, value}]}`. JSON encoding. Persistent delivery mode. ~40 values per message keeps the message rate equal to the reading rate. |
| Publishing | Publisher confirms on. Publishing is **non-blocking for the worker**: rows go into a bounded in-memory buffer, and one publisher goroutine drains it to RabbitMQ. |
| Broker down | The publisher reconnects with backoff (1 s doubling to 1 min). While it is down, the buffer fills. When the buffer is full, the worker **falls back to today's direct database write** for that reading. Nothing is silently lost, and every fallback is counted and logged. |
| Consuming | `RABBITMQ_CONSUMERS` (default 2) consumers with prefetch `RABBITMQ_PREFETCH` (default 1000). Each collects messages until `RABBITMQ_BATCH_ROWS` rows (5000) or `RABBITMQ_BATCH_MS` (200 ms), then writes **one telemetry insert and one live-value upsert** for the whole batch. It acks the batch (ack multiple) only after both writes commit. |
| Write failure | The database refuses or is down: nack with requeue, back off, try again; messages wait safely in the queue. A message that cannot be decoded, or is refused by a data error (not a connection error), is nacked without requeue and goes to the dead-letter queue. Every dead letter is logged. |
| Delivery | **At least once.** A redelivered batch must not create duplicates or move a value backwards (Step 3). |
| Ordering | Not guaranteed across consumers, so a live value carries its read time, and an older value never overwrites a newer one (Step 3). |
| Where code lives | `internal/queue`: the RabbitMQ side only (connect, declare, publish with confirms, consume, reconnect). No SQL and no business rules, like `internal/opcua`. The batching writer goroutine goes in `internal/monitor`, because concurrency lives there. It writes through `MachineTelemetryService`, never a repository directly. |

## 4. Step-by-step

### Step 0: Measure first, and fix the pool (no RabbitMQ yet)

1. Make the pool size configurable: `DB_MAX_CONNS` (e.g. default 20) in `config.DBConfig`, applied in `db.Connect` via `pgxpool.ParseConfig`. Add it to the Configuration block in `.claude/CLAUDE.md`.
2. Write down the baseline, so Step 8 can show the gain:
   - the `cnc: reading rate` log line (callbacks / processed / skipped / stored per 10 s per machine);
   - `dropped_count` in `machine_live_status`;
   - database CPU.

**Check:** build, vet, `go test ./tests/...` unchanged. Note the baseline numbers in the ADR.

### Step 1: ADR and configuration

1. Write `docs/adr/0023-rabbitmq-telemetry-queue.md` (Proposed) from Sections 1–3, and add its row to `docs/adr/README.md`.
2. Add `RabbitMQConfig` to `internal/config/config.go`. Every value gets a default, so nothing changes when unset:
   - `RABBITMQ_ENABLED` (false)
   - `RABBITMQ_URL` (`amqp://guest:guest@localhost:5672/`; never log the password)
   - `RABBITMQ_EXCHANGE` (telemetry)
   - `RABBITMQ_QUEUE` (telemetry.values)
   - `RABBITMQ_CONSUMERS` (2)
   - `RABBITMQ_PREFETCH` (1000)
   - `RABBITMQ_BATCH_ROWS` (5000)
   - `RABBITMQ_BATCH_MS` (200)
   - `RABBITMQ_PUBLISH_BUFFER` (number of messages; 20000)
   - `RABBITMQ_WRITER` (true: this process also runs the consumers)
3. Add a log switch `CNC_LOG_QUEUE` (logType `queue`) in `internal/monitor/cnc_log.go`, the same way `CNC_LOG_PLC` was added.
4. `go get github.com/rabbitmq/amqp091-go`.

**Check:** builds for 64-bit and `GOARCH=386`; the new env vars are listed in `.claude/CLAUDE.md` › Configuration.

### Step 2: Message model

1. Add the message type to `internal/model/machine_telemetry.go`: a batch of one reading (machine id, read time, source, live rows, telemetry rows). Reuse the existing `LiveValue` / `TelemetryValue` shapes where they fit.
2. `model.LiveValue` must carry the **read time**, not the database's `CURRENT_TIMESTAMP` (see Step 3).

**Check:** builds.

### Step 3: Make the two writes batch-safe (repository + migration)

1. **Live values:** `MachineTelemetryRepository.UpsertLiveValues`
   - Store `updated_at` = the reading's time from the row (an `unnest` column), not `CURRENT_TIMESTAMP`.
   - Add `WHERE machine_variable_values.updated_at <= EXCLUDED.updated_at` to the `DO UPDATE`, so an older, redelivered or reordered value never overwrites a newer one.
   - Inside one batch, first keep only the newest row per (machine, variable): Postgres refuses the same key twice in one `ON CONFLICT` statement.
2. **Telemetry duplicates:** a redelivered message would insert its rows again. Add a migration:
   - a unique index on `machine_telementry (machine_id, variable_id, recorded_at)`. On a hypertable a unique index must include the time column; this one does.
   - Before creating it, check for existing duplicates and remove them in the migration.
   - Then `InsertTelemetry` gets `ON CONFLICT DO NOTHING`.
   - Note: this index can replace migration 038's non-unique `(machine_id, variable_id, recorded_at DESC)` index. Decide in the ADR whether to keep both, and check the Axis Values query plan (ADR 0019) still uses an index.
3. Both methods must stay correct for the direct path (`RABBITMQ_ENABLED=false`).

**Check:** existing tests pass (`tests/dashboard_test.go` Axis Values, `tests/plc_collector_test.go`, `tests/virtual_screen_test.go`); add a test where the same rows are written twice (no duplicate) and an older live value is written after a newer one (the newer value stays).

### Step 4: `internal/queue`, the RabbitMQ adapter

One small package, the protocol side only:

1. **Connect + declare:**
   - dial `RABBITMQ_URL`;
   - declare the exchange, the quorum queue (with dead-letter arguments), the DLX and the dead queue, and the binding.
   - Watch `NotifyClose` and reconnect with backoff (1 s doubling to 1 min). Re-declare after every reconnect.
2. **Publisher:**
   - one channel in confirm mode, one goroutine owning it;
   - `Publish(batch)` only puts the message in the bounded buffer and returns false when the buffer is full (the caller then writes directly);
   - the goroutine sends persistent messages and waits for confirms;
   - a nacked or unconfirmed message is retried; on shutdown it is written directly, so nothing is lost.
3. **Consumer:** a channel with `Qos(prefetch)` that delivers messages to a handler. Ack / nack (with or without requeue) are exposed to the caller as plain functions.
4. Counters for the logs and diagnostics: buffered, published, confirmed, failed, fallback, reconnects.

**Check:** a unit test of reconnect / declare is not needed. An integration test in `internal/queue` runs against a real broker only when `RABBITMQ_TEST_URL` is set (otherwise `t.Skip`): publish N messages, consume N, all confirmed.

### Step 5: Writer: queue → database, in batches (`internal/monitor/telemetry_writer.go`)

1. `RABBITMQ_CONSUMERS` goroutines. Each one accumulates delivered messages until the batch is full (rows) or `RABBITMQ_BATCH_MS` passed since the first message, then:
   1. merges all telemetry rows into one slice and all live rows into one slice (newest per machine + variable);
   2. calls `MachineTelemetryService.InsertTelemetry` and then `UpsertLiveValues`;
   3. if both succeed, acks every message of the batch (ack multiple up to the last delivery tag);
   4. on a database error, nacks with requeue and waits (backoff) before taking more;
   5. on a decode / data error, nacks that message without requeue (dead letter) and logs it.
2. Log once per minute (logType `queue`): messages, rows written, batches, average batch size, write ms, nacks, dead letters, queue depth (from the queue declare's message count).

**Check:** a unit test of the batching rules, without RabbitMQ or a database: flush by size, flush by time, newest live value wins, ack only after success.

### Step 6: Workers publish instead of writing

1. In `worker.go` (CNC):
   - Where `upsertLiveValues` + `insertTelemetry` are called for a reading, build one batch message instead and `Publish` it.
   - If the queue is off, or `Publish` returns false (buffer full or broker down), call today's two direct writes; also count `fallback`.
   - Update the worker's in-memory `liveValues` cache when the reading is **accepted**: queued or written. Today it is updated only after a successful upsert. Change detection stays in memory, exactly as today.
   - `storedCount` / `droppedCount` now mean "queued or written" / "could not be queued nor written". Update their comments.
   - Keep `hub.PublishValues` first, unchanged.
   - Keep `recordAlarmChanges` triggered exactly as today, when telemetry rows exist.
2. **Stop upserting unchanged values** (the biggest saving):
   - Send live rows only for values that changed, plus all of them on the first reading after a connect or load.
   - Before doing this, check every reader of `machine_variable_values.updated_at`:
     - diagnostics `StoredAt` (`machine_diagnostics_repository.go`);
     - Virtual Screen `updatedAt` (`virtual_screen_repository.go`), which is shown as data freshness and drives the client's Stale chip (60 s).
   - If any screen needs "last read" rather than "last changed", keep that from `machine_live_status.last_reading_at`, or refresh unchanged values at most once per `N` seconds per machine. Write the choice in the ADR, and update the CLAUDE.md line "upserts every reading to detect changes".
3. Do the same in `plc_worker.go` (`handleChanges` → `store`); it already sends only changes.

**Check:**
- `RABBITMQ_ENABLED=false`: all tests pass unchanged, including `TestPLCCollector` (stop your own dev server first; it also reads test machines).
- With a broker: the same tests pass with the writer running in the test.

### Step 7: Wiring, start and stop order

1. `cmd/server/main.go` / `collector.Supervisor`:
   - when `RABBITMQ_ENABLED`, connect the queue before the collectors start;
   - start the writer before the collectors (`RABBITMQ_WRITER=true`);
   - pass the publisher to the CNC and PLC collectors (nil = direct writes).
   - A broker that is down at start must **not** stop the server: the publisher keeps reconnecting and workers fall back to direct writes.
2. Shutdown order:
   1. stop collectors;
   2. flush the publisher buffer (publish, or write directly if the broker is gone);
   3. let the writer finish and ack its current batch;
   4. close the queue connection;
   5. then the database pool.

   It must fit inside `CNC_STOP_TIMEOUT_SEC`.
3. `cnc collector: settings` log: add `rabbitmq: enabled / url host (no password) / writer / consumers`.

**Check:**
- Start without a broker: the server runs, logs reconnect attempts, and data is still stored (fallback).
- Start the broker: the publisher connects and writes go through the queue.
- Stop the broker while running: no data is lost. Rows fall back to direct writes, or wait in the queue.
- Restart the server: queued messages are written.

### Step 8: Load test and tuning

1. A test command or test (skipped without `RABBITMQ_TEST_URL`) that publishes simulated readings: e.g. 2000 machines × 40 variables × 1/s for 2 minutes, through the real publisher → queue → writer → database.
2. Watch:
   - queue depth in the management UI: it must stay flat, not grow;
   - rows/s written, batch size, write ms;
   - database CPU and connections;
   - fallback count (should be 0).
3. Tune `RABBITMQ_BATCH_ROWS` / `BATCH_MS` / `CONSUMERS` / `PREFETCH` and `DB_MAX_CONNS`. Put the final values and the before/after numbers (vs Step 0) in the ADR. If one statement per batch is still the limit, switch `InsertTelemetry` to pgx `CopyFrom` into a temporary table, then `INSERT … ON CONFLICT DO NOTHING`. Do that only if the numbers say so.

**Check:** the queue depth stays bounded at the target load; set the ADR to Accepted once agreed.

### Step 9: Diagnostics (small, optional)

Add a `queue` block to `GET /machines/{id}/diagnostics`:
- enabled, connected, queue depth, fallback count, last error;
- from the publisher / writer counters, via a read-only status struct passed in like `opcua.StatusBoard`.

Only if wanted. It changes the API: swag annotation, `go generate ./docs/api`, the client `docs/pages/machine-diagnostics.md`, and a small card in `MachineDiagnostics.jsx`. Otherwise skip, and the logs are the monitoring.

### Step 10: Documentation (same task, before saying done)

- `.claude/CLAUDE.md`:
  - Tech stack: add `rabbitmq/amqp091-go`.
  - Directory: add `internal/queue`.
  - CNC data collection › Flow and Tables written: queue path, fallback, at-least-once, live values only on change if Step 6.2 is done.
  - Configuration: the env vars.
- `docs/cnc-data-collection-flow.md` and `docs/live-data-flow.md`: the new write path (a Mermaid sequence: worker → hub → publisher → RabbitMQ → writer → PostgreSQL), and that the socket is unchanged.
- `README.md` (server): how to run RabbitMQ in dev (Docker command, management UI URL), and that `RABBITMQ_ENABLED=false` needs no broker.
- ADR 0023 final numbers.
- Client: no change unless Step 9.
- Reply line: `Docs: pages [..] | openapi [none, or diagnostics if Step 9] | postman [..] | adr [0023]`.

## 5. Later, not in this task

- **Thousands of machines need several collector processes:** one process holds at most 990 machines (DLL slots). Make the shard configurable (`CNC_SHARD` / `CNC_SHARDS` instead of the fixed `collectorShard` / `collectorShardCount`) and run one process per shard. All of them publish to the same queue, and one of them (or a separate writer-only process with `CNC_ENABLED=false`, `RABBITMQ_WRITER=true`) writes.
- **The live socket in several processes:** each process's hub knows only its own machines. A browser connected to process A would not see machines collected by process B. When there is more than one process:
  - publish live values to a fanout exchange `live`;
  - every API process feeds its hub from it;
  - or put the API behind one process.

  That needs its own ADR.
- Moving events, modes and cycles to the queue: only if they ever become a load; their ordering and the read-back the cycle logic needs make them a poor fit.

## 6. Rules the implementer must keep

- Layers: handler → service → repository. The writer calls `MachineTelemetryService`. `internal/queue` has no SQL.
- `RABBITMQ_ENABLED=false` = today's behaviour, byte for byte in the database.
- Never block a machine worker on RabbitMQ or the database for the high-volume rows. Never drop data silently: fall back, requeue or dead-letter, and count every case.
- Live socket values still go to the hub first, in-process.
- Every collector log line goes through `writeLog` (`cnc_log.go`), logType `queue`; never log the broker password.
- `gofmt`, `go build ./...`, `go vet ./...`, both `GOARCH`s; `go test ./tests/...`; broker-dependent tests skip without `RABBITMQ_TEST_URL`.
- No `git commit` / `push` unless asked.
