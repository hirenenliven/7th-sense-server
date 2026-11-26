# 0029. Telemetry queue batches are stored with COPY

- **Date:** 2026-11-25
- **Status:** Proposed

## Context

The telemetry queue writers ([ADR 0024](0024-rabbitmq-telemetry-queue.md)) store each batch's `machine_telementry` rows with one `INSERT … SELECT FROM unnest(arrays) ON CONFLICT DO NOTHING` (`InsertTelemetry`). ADR 0024 named COPY as the next step if that statement became the limit. `ON CONFLICT DO NOTHING` matters: the queue delivers at least once, so a batch can arrive again, and the unique index of migration 039 must skip the rows already stored. Plain COPY has no `ON CONFLICT`: one duplicate fails the whole COPY.

## Decision

`MachineTelemetryRepository.CopyTelemetry`, called by `MachineTelemetryService.WriteTelemetryBatch` (and so by every queue batch, the per-reading retry of a refused batch and the unsent readings written at stop):

1. opens a savepoint in the batch's transaction (`tx.Begin`),
2. runs a binary `COPY machine_telementry (machine_id, variable_id, value, recorded_at)` (pgx `CopyFrom`),
3. on success releases the savepoint; on a unique violation (23505, a redelivered batch) rolls back to it and writes the batch with `InsertTelemetry`, which skips the stored rows; any other error is mapped as before (`mapTelemetryError`: a refused value → `ErrInvalidTelemetry`).

Direct writes (one reading of a few dozen rows, `worker.insertTelemetry`, the PLC worker) keep `InsertTelemetry`: COPY's extra round trip gains nothing there.

Measured on the dev PC (same database as ADR 0026 / 0027), 30 batches × 10,000 rows, each in its own transaction, twice:

| Write | Per batch | Rows/s |
|---|---|---|
| `InsertTelemetry` (unnest) | 77–80 ms | ~125–130k |
| `CopyTelemetry` | 38–40 ms | ~250–260k |
| `CopyTelemetry`, every batch a duplicate (fallback) | 45 ms | ~220k |

## Alternatives considered

- **COPY into a temporary table, then `INSERT … SELECT … ON CONFLICT DO NOTHING`:** works for every batch, but still pays the full INSERT into the hypertable plus a temp table write, so it gains far less; duplicates are rare (only after a lost ack or a crash).
- **Keep the INSERT:** half the telemetry write rate.

## Consequences

- The `machine_telementry` part of a batch writes about twice as fast; `machine_variable_values` (upsert) is unchanged and is now the larger part of a batch's time.
- A redelivered batch costs a failed COPY plus the INSERT (rare).
- `CopyTelemetry` needs a transaction (savepoint); it returns an error on a pool.
- The disk limit of ADR 0027 still holds for sustained load: COPY saves CPU and statement work, not WAL or index writes.
