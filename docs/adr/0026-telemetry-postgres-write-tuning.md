# 0026. Telemetry writes: one asynchronous commit per batch, one index less, PostgreSQL settings

- **Date:** 2026-11-24
- **Status:** Proposed

## Context

The 60 s end-to-end test of [ADR 0024](0024-rabbitmq-telemetry-queue.md) (1000 machines × 40 variables × 1 reading/s = 40,000 values/s, 4 writers) stored ~23,000 values/s while sending (316 ms per 5,000-row batch); readings piled up after ~30 s and the last was stored 12 s after the end. The writers were the limit, and the hardware stays as it is (one SATA SSD, a Kingston A400 without DRAM cache, shared by PostgreSQL, RabbitMQ and Windows; 16 GB RAM).

What the database did per batch: two statements, each its own transaction whose COMMIT waits for the WAL flush (`synchronous_commit = on`); every `machine_telementry` row written to two indexes, the unique `(machine_id, variable_id, recorded_at)` key (migration 039) and TimescaleDB's default `(recorded_at DESC)`, which no query uses (every query names machine and variable); 7-day chunks, whose indexes outgrow memory at this rate; `max_wal_size` 1 GB, so 55 checkpoints were forced by WAL volume.

## Decision

| Topic | Decision |
|---|---|
| Batch write | The queue writer stores a batch through `MachineTelemetryService.WriteTelemetryBatch`: the `machine_telementry` insert and the `machine_variable_values` upsert in **one transaction** with `SET LOCAL synchronous_commit = off` (`MachineTelemetryRepository.UseAsynchronousCommit`). One commit instead of two, and it does not wait for the disk. Direct writes (no queue, fallback) and everything else (modes, cycles, events, statuses) keep synchronous commits. |
| Index | Migration 041 drops `machine_telementry_recorded_at_idx`. |
| Chunks | Migration 041: new `machine_telementry` chunks hold 1 day (`set_chunk_time_interval`, only on a hypertable). |
| Live values | Migration 041: `machine_variable_values` fillfactor 80, so the constant updates (value / updated_at, not indexed) stay on their page (HOT). |
| PostgreSQL settings (not in a migration: `ALTER SYSTEM` + `SELECT pg_reload_conf()`, no restart) | `max_wal_size = '4GB'` (8 GB with more free disk), `checkpoint_timeout = '15min'`, `wal_compression = 'lz4'`. Set on every server that runs the collector. |

## Measurements (this PC, PostgreSQL alone, 200 batches × 5,000 values, insert + upsert)

| Setup | Values/s | WAL bytes / value |
|---|---|---|
| Before (two synchronous commits, `recorded_at` index, WAL uncompressed) | 18,856 | 460 |
| After (one asynchronous commit, no `recorded_at` index, lz4) | 25,440 – 34,401 | 409 – 455 |

Batch write in the end-to-end test: 316 ms → 42 ms while RabbitMQ was the slower side. With both writing hard, the shared disk is the limit for both: whichever writes more slows the other (a classic queue let RabbitMQ take 56,000 readings in 60 s, and PostgreSQL batches went back to ~350–540 ms). Runs back to back on this drive vary by ±30 %.

## Alternatives considered

- **`synchronous_commit = off` for the whole database:** would also apply to modes, cycles, sessions and user data. Only the queue batches accept the risk below.
- **`COPY` instead of `INSERT … unnest`:** saves CPU and parsing; the limit here is disk writes, which it does not reduce. Next step only if the CPU becomes the limit.
- **`fsync = off` / unlogged table:** a crash can corrupt or empty the table. Rejected.
- **Compression of old chunks:** saves disk space, not write time. Separate task.
- **Store fewer values (deadband per variable, `store_mode`):** the biggest possible gain (noisy analog values like loads change on every reading), but it changes what is recorded: a business decision, separate task.

## Consequences

- Up to ~1.8× more values per second on the same disk; one index less to keep per row; WAL-forced checkpoints mostly gone.
- **Risk:** if PostgreSQL or Windows crashes (not the Go server), batches committed in the last ~0.6 s (3 × `wal_writer_delay`) may be lost although RabbitMQ was already told they were stored. The database stays consistent. Accepted for telemetry samples; not used for any other data.
- Existing 7-day chunks keep their size until they age out; the 1-day chunks start with the next new chunk.
- On a database without TimescaleDB the chunk step is skipped.
- Rollback: migration 041's header; `ALTER SYSTEM RESET` the three settings; call `InsertTelemetry` + `UpsertLiveValues` again in `TelemetryQueue.writeRows`.
