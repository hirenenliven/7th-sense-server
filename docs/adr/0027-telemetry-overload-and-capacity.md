# 0027. Telemetry queue under overload: larger confirm groups, the writers' own database pool; the 2 lakh values/s ceiling

- **Date:** 2026-11-24
- **Status:** Proposed

## Context

The target is 2 lakh (200,000) values/s, e.g. 1000 machines × 40 variables × 5 readings/s. Each value is two row writes: one `machine_telementry` insert and one `machine_variable_values` upsert, so 400,000 row writes/s. The 5/s run of [ADR 0024](0024-rabbitmq-telemetry-queue.md) stored about 50k rows/s. Sending was held up as well: 18 % of the readings were written directly, and the whole run took 2 min 48 s instead of 60 s. [ADR 0026](0026-telemetry-postgres-write-tuning.md) raised PostgreSQL alone to 25–34k values/s on this PC (one DRAM-less SATA SSD shared with RabbitMQ, 16 GB RAM).

Two faults in our own code made overload worse:

1. **The publisher waited for each group of 500 confirms before it sent more.** A quorum queue confirms after its fsync, which takes longer when PostgreSQL keeps the shared disk busy. The publisher's rate then fell to 500 messages per confirm wait. The 20,000-message buffer filled up, and the workers fell back to direct writes.
2. **The direct writes and the queue writers shared one pool** (`DB_MAX_CONNS` 20). Up to 1000 workers were waiting for a connection, and each sent two small synchronous transactions. The 4 writers waited in the same line, so a batch took 5.6 s. The writers slowed down exactly when the queue grew.

## Decision

| Topic | Decision |
|---|---|
| Confirm groups | `queue.publishGroupSize` 500 → 5000. A group still takes only what is buffered, so latency at low load does not change. Under load, one confirm wait covers up to 5000 messages. Retries are the same: a message that is not confirmed is sent again. |
| Writers' pool | `collector.Supervisor.openWriterPool` opens a pool of `RABBITMQ_CONSUMERS` connections that only the queue writers use (and the unsent readings written at stop). `DB_MAX_CONNS` stays for the API and the direct writes. If the extra pool cannot be opened, the writers share the main pool as before and an error is logged. The pool is closed after the queue stops. |

The total number of PostgreSQL connections goes up by `RABBITMQ_CONSUMERS` (4).

## What no setting or code change can reach on this disk

As long as every value is stored as its own row, 2 lakh values/s is beyond this hardware:

| | Per value | At 200,000 values/s |
|---|---|---|
| Rows written | 1 insert + 1 upsert | 400,000 row writes/s; PostgreSQL alone managed 25–34k values/s (ADR 0026) |
| Table + index on disk | ~60 B row + ~45 B unique index entry | ≈ 1.8 TB per day |
| WAL written | 409–455 B (ADR 0026) | ≈ 86 MB/s sustained, ≈ 7 TB per day. A DRAM-less A400 slows down when its SLC cache is full, and is rated for 80–300 TB written in total (by size), so it would wear out in weeks. |

RabbitMQ absorbs bursts: the queue grows on disk and the writers drain it afterwards. A load that stays above what PostgreSQL stores makes the queue grow without end.

## Alternatives considered (each needs a decision first)

- **Store fewer values** (a deadband per variable or `store_mode`, already in the schema but not used): the biggest gain, since load and override values change on every reading. It changes what is recorded.
- **One `machine_telementry` row per reading** (machine, time, variable ids[], values[]) instead of one per value: about 10× fewer bytes of heap, index and WAL with 40 values per reading. It needs a migration, and the Machine Health queries (`MachineHealthRepository`) must be rewritten.
- **TimescaleDB compression policy** on chunks older than a day: about 10–20× less disk. Writes are not faster, but at this rate it is needed for disk space with any of the options.
- **Coalesce live values** across batches (one `machine_variable_values` upsert per machine-variable per second): up to 2× fewer row writes at 5 readings/s, none at 1 reading/s. Live values written in the last second could be lost on a crash (the worker's 30 s refresh writes them again).
- **More writers / bigger batches:** the disk is the limit, not the connection count or the statement count.

## Consequences

- Under overload RabbitMQ takes the readings instead of the workers writing directly, and the writers keep their connections. The database writes at its own top speed while the queue holds the rest.
- The publisher holds up to 5000 messages (~12 MB) between confirms, instead of 500.
- Rollback: `publishGroupSize = 500`, and in `StartTelemetryQueue` build the service on `supervisor.dbPool`.
