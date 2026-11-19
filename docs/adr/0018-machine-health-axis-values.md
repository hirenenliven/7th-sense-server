# 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes

- **Date:** 2026-11-18
- **Status:** Accepted (the "always downsample" part superseded by [0019](0019-axis-values-raw-and-zoom.md))

## Context

The Machine Health tab needs X, Y and Z of an axis (machine position by default) as three charts, read from `machine_telementry`. The table holds every changed value of every machine (13.6 million rows locally, about 1.1 million per position variable of one machine; production will hold crores), and the team wants to see what a chart over such a table costs before building more charts. The chart is temporary and more Machine Health charts will follow, so it should not get a resource of its own.

`machine_telementry` had only an index on `recorded_at`. Values are stored only when they change.

## Decision

- **Endpoint on the existing machine-health routes**: `GET /machine-health/axis-values?machineId&axisSet&minutes` (or `date` / `from`–`to`), on `DashboardHandler` / `DashboardService` beside `machine-time-by-day` (`dashboard_routes.go`). Queries in `repository.MachineHealthRepository`, a home for the next Machine Health charts. Every logged-in user, for machines they may read.
- **Axis sets are a fixed server list** (`dto.AxisSets`: machine / work / relative position, distance to go, axis load), each three `cnc_variables` codes. The client sends a key, never codes.
- **Downsample in SQL to at most 600 points per axis** (`dto.AxisValuesMaxPoints`, `time_bucket` from the range start). Each point is the bucket's **last** value (with min, max, row count), plus the last value before the range (`startValue`); the client draws a step line. The last value comes from the bucket's `MAX(recorded_at)` looked up by index: TimescaleDB `last(value, recorded_at)` gave the same result about 3.5× slower (30 days, 3.3 million rows: 2.75 s vs 0.79 s).
- **The response says what it cost**: `stats {rowsRead, pointCount, queryMs}` (database time of the telemetry queries); the client adds the total round trip and shows both.
- **Index `(machine_id, variable_id, recorded_at DESC)`** (migration 038), so a chart reads only its machine's rows of its variables, and "last value before the range" is one lookup.

## Alternatives considered

- **Own route / handler / service** (`/axis-values`). Rejected by the team: temporary chart, more to come under Machine Health.
- **Send raw rows and downsample in the browser.** Millions of rows per request. Rejected.
- **Average per bucket.** Fast, but a position drawn as an average is a value the axis never had; last value + min / max keeps the truth at almost the same cost.
- **TimescaleDB continuous aggregates** (pre-computed buckets). Fastest for long ranges, but fixed bucket sizes, refresh policies and storage; worth it only if 7- / 30-day ranges prove too slow with the index. Left for later.

## Consequences

- Measured locally (machine 195, 13.6 million rows in the table, before migration 038): 1 hour ≈ 25 ms, 1 day ≈ 110 ms, 7 days ≈ 310 ms, 30 days ≈ 790 ms of database time. Today one machine owns almost every row, so the old `recorded_at` index happens to work; with many machines it would read every machine's rows in the range, which migration 038 prevents.
- Migration 038 builds an index on a large hypertable: run it while collection can pause.
- Query time grows with the rows in the range (30 days reads every row of those 30 days); points sent stay ≤ 600 per axis.
