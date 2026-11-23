# 0024. Telemetry rows through a RabbitMQ queue, written in batches

- **Date:** 2026-11-20
- **Status:** Proposed

## Context

Goal: thousands of machines × ~40 variables must be collected without the database writes slowing the machine workers or losing data (implementation guide: [`../rabbitmq-telemetry-queue-guide.md`](../rabbitmq-telemetry-queue-guide.md); it names this ADR 0023, a number taken meanwhile by the PLC session ADR).

Before this change one goroutine per machine (`internal/monitor/worker.go`; PLC: `plc_worker.go`) did, for **every reading**, synchronously: publish to the live hub, upsert **every** numeric value (changed or not) into `machine_variable_values`, insert the changed ones into `machine_telementry`. Each write was its own statement with a 5 s timeout on a pgx pool of the default size, max(4, CPUs) connections. At 1000 machines × 40 variables × 1 reading/s that is ~40,000 single-row-set statements a second queueing for a handful of connections, and readings back up behind the database (`cnc: reading rate` shows them as `skipped`). A failed telemetry insert was dropped (`dropped_count`).

Baseline (Step 0): `DB_MAX_CONNS` was not configurable; the per-machine rate is in the `cnc: reading rate` log line, losses in `machine_live_status.dropped_count`.

## Decision

The high-volume rows go through RabbitMQ; everything else stays a direct write.

| Topic | Decision |
|---|---|
| Library | `github.com/rabbitmq/amqp091-go` v1.17 (the official client; `streadway/amqp` is deprecated). Pure Go, builds for `GOARCH=386`. |
| Broker | RabbitMQ 3.13+ / 4.x with the management plugin. Dev: the Windows installer (tested on 4.2.2) or Docker `rabbitmq:management`. |
| Scope | In the queue: `machine_telementry` and `machine_variable_values` rows, from the CNC and the PLC worker. Not in the queue: modes, events, cycles, `machine_live_status` — low volume, order-sensitive, and the worker needs their result (a refused cycle open). The live hub is still fed in-process first. No API or client change. |
| Topology | Durable direct exchange `telemetry`, routing key `values`, durable quorum queue `telemetry.values`, dead letters to the fanout exchange `telemetry.dlx` → queue `telemetry.values.dead`. No length limit (the publish buffer and the direct-write fallback already bound memory on the server side). Declared on every connect (idempotent); a queue that exists with other arguments fails with a message saying to delete it. |
| Message | One per reading of one machine: `model.TelemetryReading {machineId, readAt, source: cnc/plc, live: [{variableId, value}], telemetry: [...]}`, JSON, persistent. |
| Publishing | `queue.Publisher`: `Publish` only puts the message in a bounded buffer (20000 messages) and returns false when it is full; one goroutine owns a confirm-mode channel, sends groups of 500 and waits for their confirms; what is nacked or not confirmed is sent again. |
| Broker down / buffer full | The worker writes that reading **directly**, exactly as before; every fallback is counted (`fallback`) and the first per minute logged. The publisher reconnects with backoff 1 s doubling to 60 s. A broker down at start does not stop the server. |
| Consuming | `RABBITMQ_CONSUMERS` (default 4 since 2026-11-23, was 2) writers (`monitor.telemetryWriter`), prefetch 1000, batch until 5000 rows or 200 ms after the first message, then **one** `InsertTelemetry` and **one** `UpsertLiveValues` (newest per machine + variable) through `MachineTelemetryService`; ack multiple only after both succeed. |
| Database error | The batch is **kept unacknowledged** and the same write retried after a wait (1 s doubling to 60 s). The guide said nack + requeue; RabbitMQ 4.x quorum queues dead-letter a message after 20 deliveries by default, so requeueing through a long database outage would push good data to the dead-letter queue. Keeping the batch costs nothing: prefetch stops new deliveries meanwhile, and a lost connection redelivers it. |
| Refused value / bad message | A value PostgreSQL refuses (SQLSTATE class 22 / 23 → `repository.ErrInvalidTelemetry` → `response.InvalidTelemetry`) makes the writer write the batch reading by reading: accepted ones are acked, refused ones go to the dead-letter queue. A message that cannot be decoded goes there too. Every dead letter is logged. |
| Delivery | At least once. `machine_telementry` gets a unique index `(machine_id, variable_id, recorded_at)` (migration 039) and `ON CONFLICT DO NOTHING`; `machine_variable_values.updated_at` is now the **read time** and the upsert has `WHERE updated_at <= EXCLUDED.updated_at`, so a repeated or reordered value never moves a value back. Both also apply to direct writes. |
| Indexes | The unique index serves every query of migration 038's `(machine_id, variable_id, recorded_at DESC)` index (a b-tree is scanned both ways), so 039 drops 038's index: one index to maintain per row instead of two. |
| Unchanged live values | The CNC worker upserts changed values on every reading and all of them on the first reading after a load and then every 30 s (`CNCConfig.LiveRefresh`, fixed; before, every value on every reading). The Virtual Screen shows `updated_at` as data freshness (Stale chip after 60 s) and the diagnostics show it as "last stored", so it must stay below 60 s. The PLC worker already stores only changes. |
| Off switch | `RABBITMQ_ENABLED=false` (default): no broker needed, the workers write directly. |
| Code | `internal/queue`: RabbitMQ only (connect, declare, publish, consume, reconnect, depth), no SQL, logs through a `Logger` the caller passes. `internal/monitor/telemetry_queue.go` (shared by both collectors, nil = direct writes) and `telemetry_writer.go` (batching). Started by `collector.Supervisor.StartTelemetryQueue` before the collectors; stopped after them: publisher flushes (what it cannot send is written directly), then writers finish and ack their batch, then the pool closes. |
| Pool | `DB_MAX_CONNS` (default 20) sizes the pgx pool (`db.Connect` via `pgxpool.ParseConfig`). |
| Logs | Always written (logType `queue`; few lines): connects, fallbacks, dead letters, and once a minute `queue: statistics` (published / confirmed / fallback / reconnects, messages / rows / batches written, average batch rows, average write ms, write errors, dead letters, queue depth and dead-letter depth). The broker password is never logged. |

## Load test (Step 8)

`tests/telemetry_queue_test.go` `TestTelemetryQueueEndToEnd` (runs with `RABBITMQ_TEST_URL`) publishes simulated readings through the real publisher → RabbitMQ → writer → PostgreSQL; **every** variable changes on every reading (the worst case: 40 telemetry + 40 live rows per reading). Dev machine, RabbitMQ 4.2.2 and PostgreSQL 18 + TimescaleDB on the same Windows PC, the dev server running beside it, migration 039 not yet applied:

| Machines × 40 variables, 1 reading/s, 15 s | Queue | Writers | Result |
|---|---|---|---|
| 50 | quorum | 2 | 20,000 rows, 9 batches, 46 ms / batch, depth 0 |
| 1000 (80,000 rows/s offered) | quorum | 2 | 1.2 M rows in 15.4 s (kept up), 5,021 rows / batch, 95 ms / batch, depth 0, fallback 0 |
| 2000 (160,000 rows/s offered) | quorum | 2 | fell behind: publish buffer full, 8,000 readings refused (they would have been written directly) |
| 2000 | quorum | 4 | no fallback; 2.4 M rows done 35 s after the 15 s burst; 382 ms / batch |
| 2000 | classic | 4 | no fallback; done 23 s after the burst; 248 ms / batch |

Reading: the database is the limit on this PC at roughly 65–80k rows/s; more writers do not raise it (each batch gets slower). The queue absorbs a burst above that and drains it without loss. Real load is far below the worst case, since only changed values reach `machine_telementry` and the 30 s live refresh removes most live upserts. Defaults stay 2 writers on a quorum queue. If production numbers show one statement per batch is the limit, the next step is `CopyFrom` into a temporary table and `INSERT … ON CONFLICT DO NOTHING` — only then.

### Load test after migration 039 (2026-11-23)

`tests/telemetry_load_test.go` `TestTelemetryLoad` (runs with `TELEMETRY_LOAD=1`, settings in its doc comment) uses the product's settings from `config.Load`, one goroutine per machine like the workers, the direct write when the queue refuses a reading, and measures RabbitMQ alone, rows offered / stored, queue depth, how old the newest stored data is, fallbacks and lost rows. Same dev PC, migration 039 applied (unique index + `ON CONFLICT`), 2 writers, quorum queue:

| Run | Sent | Result |
|---|---|---|
| RabbitMQ alone (no database) | as fast as taken | 7,400–11,900 msg/s (17.7–28.4 MB/s, 2.4 KB message of 40 + 40 values) ≈ 3–4.8 lakh values/s |
| 1000 machines × 40, 1/s, 60 s, every value changes | 80,000 rows/s | 2.4 M telemetry rows, 0 lost, 0 fallback; stored ~31k rows/s; queue peak 36,643, empty 1 min 35 s after the end |
| 1000 machines × 40, 1/s, 60 s, 25 % change | 23,400 rows/s | 673,918 rows, 0 lost; 800 readings (1.3 %) written directly while confirms were slow |
| 1500 machines × 40, 1/s, 60 s, 25 % change | 35,100 rows/s | 1,010,601 rows, 0 lost, 0 fallback; done 16 s after the end |
| 1000 machines × 40, **5/s**, 60 s, every value changes (2 lakh values/s) | 400,000 rows/s | 12,000,000 telemetry rows, **0 lost**; 245,189 readings queued, 54,811 (18 %) written directly; sending held up by the direct writes (2 min 48 s instead of 60 s, 71,462 values/s really sent); queue peak 224,134; empty 5 min 10 s after the end; 7 min 58 s in all; ~50k rows/s stored |

Reading: RabbitMQ is not the limit; PostgreSQL is, at ~30–85k rows/s on this PC (RabbitMQ and PostgreSQL share one disk, so publisher confirms slow down while the database writes hard). Nothing was lost in any run. Two things to fix before such loads, not changed by the test: (1) the two writers deadlock on the `machine_variable_values` upsert (`deadlock detected`, 40P01; 51 retries in the 5/s run): rows in a batch are not sorted, so sort live rows by `(machine_id, variable_id)` before the upsert; (2) under overload the direct-write fallback from every worker takes the pool's connections from the writers (5.6 s per batch), so the writers slow down exactly when the queue grows. Queue depth here counts ready messages; up to prefetch × writers (2,000) more are being written.

## Alternatives considered

- **Batch in process, no broker** (a channel per collector feeding a batch writer): the same write saving, but readings in memory are lost on a crash or a database outage longer than the buffer, and several collector processes cannot share one writer. Kept as the fallback path only.
- **Kafka / NATS JetStream:** fine for this volume, but RabbitMQ is already installed on the target servers, has a management UI the team knows, and per-message acks fit "ack after the database commit".
- **Requeue on database errors (the guide's first plan):** see "Database error" above: quorum delivery limits would dead-letter good data during a long outage.
- **Drop 038's index and keep only a non-unique one:** no duplicate protection for at-least-once delivery.

## Consequences

- New dependency `rabbitmq/amqp091-go`; new package `internal/queue`; migration 039 (run it while collection is paused — it removes duplicates and builds a unique index on a large hypertable).
- New settings, kept to four: `DB_MAX_CONNS`, `RABBITMQ_ENABLED`, `RABBITMQ_URL`, `RABBITMQ_CONSUMERS` (the one throughput knob that matters: the load test shows the database is the limit, so batch size, prefetch and buffers gain nothing as settings). Everything else (names, queue type, batch 5000 rows / 200 ms, prefetch 1000, publish buffer 20000, reconnect 60 s, live refresh 30 s) is fixed in `config.Load`; a first version had all of them as `RABBITMQ_*` / `CNC_LIVE_REFRESH_SEC` / `CNC_LOG_QUEUE` variables, removed as unneeded. A restart applies changes.
- Data reaches the tables a few hundred ms later when queued; no screen depends on it (they read live values from the socket).
- `machine_variable_values.updated_at` is the read time, not the database's clock (same within milliseconds).
- `machine_live_status.stored_count` / `dropped_count` mean "telemetry rows queued or written" / "neither queued nor written".
- At most 20000 readings (the publish buffer) live only in memory; a crash of the server process loses them (a clean stop writes them). Everything confirmed by the broker survives restarts of the server and the broker.
- Fixed on the way: a successful direct upsert after `GetLiveValues` failed wrote into a nil map and panicked (the "Known issue" in `docs/cnc-data-collection-flow.md`).
- Later, not in this decision: several collector processes (configurable shards, one writer); the live hub across processes (a fanout exchange); a diagnostics card for the queue.
