# 30. Telemetry queue messages are gzip-compressed

Date: 2026-11-26

## Status

Proposed

## Context

The goal was 40,000 telemetry values/s (1,000 machines × 40 variables × 1 reading/s, every value changing), not counting `machine_variable_values` upserts. Load runs passed sometimes and fell behind other times. Measurements on the development PC (i5-8400, 6 cores; Kingston A400 240 GB SATA SSD, DRAM-less, 87 % full):

- PostgreSQL was not the limit. The writers' COPY took ~100 ms per 16,000-row batch (writers busy ~20 % of the time), and PostgreSQL alone took 40,000 rows/s for 60 s at ~8 MB/s without one slow write (WAL ~171 bytes/row).
- RabbitMQ on its own carries ~8,000–9,000 messages/s.
- The backlog came from **SSD stalls**: write latency jumped from ~1 ms to 1–8 s per write for 5–40 s, and both the broker's confirms and PostgreSQL waited. RabbitMQ alone triggered them too. The quorum queue writes every message twice (Raft log + segment file), 2.6 KB of JSON per reading, as ~2,000 small writes/s.
- The load test itself made it worse: each run DELETEd millions of rows in place, so the next run wrote into reused space and paid for vacuum and checkpoint writes that production (no deletes) never has.

## Decision

- `monitor.EncodeReading` writes each queue message as gzip (BestSpeed) JSON, with pooled gzip writers and readers. A reading's JSON repeats the same keys for every value, so it shrinks 83 % (3,044 → 505 bytes with realistic 3-decimal values). `decodeReading` detects the gzip magic number and still reads plain JSON, so messages queued before the upgrade are stored as before.
- `RABBITMQ_PUBLISH_BUFFER` defaults to 60,000 (was 20,000): a minute of 1,000 readings/s, ~30 MB of gzipped messages. A disk stall of the broker is then ridden out in memory instead of switching readings to synchronous direct writes, which stall as well.
- `RABBITMQ_QUEUE_TYPE` is read from the environment, as ADR 0028 says (`config.Load` had it fixed to quorum).
- The load test counts telemetry only (watermark, drain and rate from `machine_telementry`; live rows are checked after, untimed). Its readings are dated 50 years ahead and their chunks are dropped after the run, and its values are realistic 3-decimal numbers so compression is not flattered.

## Alternatives considered

- **A message format without the repeated keys, or without live values that repeat the telemetry ones.** Smaller still, but it changes the message contract on both sides; gzip gets most of the gain with two functions and stays readable as JSON after decompression.
- **More writers or bigger batches.** The writers were mostly idle, and bigger batches did not help in earlier runs (ADR 0027).
- **Classic queue.** Possible through `RABBITMQ_QUEUE_TYPE`, but it loses the quorum queue's guarantees, and gzip alone removed the broker's share of the stalls.

## Consequences

- RabbitMQ writes fell from 2.8 MB/s and ~2,050 write calls/s to ~0.8 MB/s and ~80/s at 40,000 values/s. PostgreSQL (~11–13 MB/s) is now the main writer, close to what storing 40,000 rows/s with the unique index and WAL needs.
- Results at 40,000 values/s, every value changing: 5-minute run kept up (12,000,000 of 12,000,000 stored, 39,776 values/s, no fallback, queue empty 1.8 s after sending, p95 data age 12 s across one SSD stall). In 60 s runs, 4 of 7 kept up, and every one stored every value. The 3 that fell behind each hit an SSD stall of 1–8 s per write and caught up within 13–25 s.
- What remains is the disk: on a nearly full DRAM-less SSD, stalls still happen. A drive with DRAM (or NVMe) and free space for PostgreSQL and RabbitMQ is what removes them; 40,000 values/s also stores ~7.5 GB of heap an hour before compression.
- Messages in `telemetry.values.dead` are gzip and need gunzip to read. CPU cost ~25 µs per message.
