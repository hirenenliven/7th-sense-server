# Telemetry load test plan: step by step to 1000 machines

Results of every session are recorded in [telemetry-load-test-results.md](telemetry-load-test-results.md).

How to prove the telemetry pipeline (CNC/PLC worker → RabbitMQ → writers → PostgreSQL, [ADR 0024](adr/0024-rabbitmq-telemetry-queue.md), [0027](adr/0027-telemetry-overload-and-capacity.md), [0028](adr/0028-telemetry-queue-settings-from-env.md), [0029](adr/0029-telemetry-batch-copy.md)) holds 1000 machines, increasing the load one step at a time with `load-test.ps1`.

## 1. What the script does

`.\load-test.ps1 -Values <n> -Seconds <s> -Rate <r>` runs `TestTelemetryLoad` (`tests/telemetry_load_test.go`):

- Simulates **machines = Values ÷ (40 × Rate)** (rounded up), each with **40 variables**, sending **Rate readings per second**.
- **Every value changes on every reading** (the worst case). Each value = **2 rows**: `machine_telementry` (COPY) + `machine_variable_values` (upsert).
- Uses the real publisher, queue, writers and database code, with the settings from `config.Load` (`RABBITMQ_*`, `DB_MAX_CONNS`).
- Uses test machine ids and its own queue; deletes both afterwards. **Never stop it with Ctrl+C** (the cleanup would not run).
- Prints a `conclusion` block and saves the whole output to `logs\load-test-<values>-<time>.log`.

| Machines | Rate 1/s | Rate 2/s | Rate 5/s |
|---|---|---|---|
| 10 | 400 | 800 | 2 000 |
| 100 | 4 000 | 8 000 | 20 000 |
| 250 | 10 000 | 20 000 | 50 000 |
| 500 | 20 000 | 40 000 | 100 000 |
| 1000 | 40 000 | 80 000 | 200 000 |

(values/s = `-Values`; rows/s written = 2 × that)

## 2. Before every run

1. PostgreSQL and RabbitMQ are running (`.env`: `DB_*`, `RABBITMQ_URL`). Migrations 039 and 041 are applied.
2. **Stop the dev server** (`go tool air` / `run.ps1`): it writes to the same database and disk and would distort the numbers.
3. Close heavy programs; note the hardware (CPU, RAM, disk type, whether PostgreSQL and RabbitMQ share the disk).
4. Settings print in the first line of the log (`settings: N writers …, batch … rows / …, prefetch …`). Check they are what you expect. The test reads `.env`, but a variable already set in the shell wins. To try other values, set them in the same PowerShell window before the run, e.g. `$env:RABBITMQ_BATCH_ROWS = '20000'`, and remove them after (`Remove-Item Env:RABBITMQ_BATCH_ROWS`).
5. Leave about 1 minute between runs (checkpoints and autovacuum settle).

## 3. Pass criteria (every test case)

The `RESULT` line of the conclusion, and these values from it:

| Check | Pass | Where |
|---|---|---|
| Result | `KEPT UP` | `RESULT` line |
| Lost | 0 | `Stored … (lost N)`: any loss is a **bug**, stop and investigate |
| Fallback | 0 | `Path … N written directly` |
| Sending time | ≈ `-Seconds` (within 10 %) | `Sent … sending took` |
| Queue empty | < 5 s after sending stopped | `Queue … empty after` |
| Data age p95 | < 3 s (BATCH_MS 1 s + write time) | `Data age` |
| Peak queue | stays small (a few batches), does not grow with time | `Queue peak` |

`FELL BEHIND` with lost 0 is not a bug. It means this hardware's limit is below the step: the stored values/s shown is roughly the capacity. Record it, stop stepping up, and go to section 6.

## 4. Test cases: step up (worst case, every value changes)

Run in order. Go to the next step only when the current one passes.

| ID | Machines | Command | Values/s | Rows/s | Purpose |
|---|---|---|---|---|---|
| LT-01 | 10 | `.\load-test.ps1 -Values 400 -Seconds 60` | 400 | 800 | Smoke test: setup, cleanup, log file |
| LT-02 | 50 | `.\load-test.ps1 -Values 2000 -Seconds 60` | 2 000 | 4 000 | Small plant |
| LT-03 | 100 | `.\load-test.ps1 -Values 4000 -Seconds 60` | 4 000 | 8 000 | |
| LT-04 | 250 | `.\load-test.ps1 -Values 10000 -Seconds 60` | 10 000 | 20 000 | Script's default load |
| LT-05 | 500 | `.\load-test.ps1 -Values 20000 -Seconds 60` | 20 000 | 40 000 | Half target |
| LT-06 | 750 | `.\load-test.ps1 -Values 30000 -Seconds 60` | 30 000 | 60 000 | Around this PC's measured limit (ADR 0026: 25–34k values/s before COPY) |
| LT-07 | 1000 | `.\load-test.ps1 -Values 40000 -Seconds 60` | 40 000 | 80 000 | **Target: 1000 machines at 1 reading/s** |
| LT-08 | 1000 | `.\load-test.ps1 -Values 40000 -Seconds 600` | 40 000 | 80 000 | Soak: 10 min at the target; checkpoints, autovacuum and chunk growth show up only over time |

## 5. Test cases: faster readings (only after LT-08 passes)

The same 1000 machines sending more often (a Siemens PLC publishes every 500 ms, ADR 0021).

| ID | Machines | Command | Values/s | Rows/s | Purpose |
|---|---|---|---|---|---|
| LT-09 | 1000 | `.\load-test.ps1 -Values 80000 -Seconds 60 -Rate 2` | 80 000 | 160 000 | 1000 machines × 2 readings/s |
| LT-10 | 1000 | `.\load-test.ps1 -Values 200000 -Seconds 60 -Rate 5` | 200 000 | 400 000 | The 2 lakh values/s goal of ADR 0027; expected to fall behind on the dev PC |

## 6. Test cases: tuning (at the first step that fell behind, or LT-07)

Change **one** setting at a time, rerun the same command, and compare `stored values/s` and `Data age` with the baseline. Keep a value only if it clearly helps. Put the winners in `.env`.

| ID | Setting (PowerShell, before the run) | Try | What to look at |
|---|---|---|---|
| LT-T1 | baseline | defaults | the reference numbers |
| LT-T2 | `$env:RABBITMQ_CONSUMERS` | 2, 4, 8 | more writers help only until the disk is busy (ADR 0027) |
| LT-T3 | `$env:RABBITMQ_BATCH_ROWS` | 5000, 10000, 20000 | fewer commits vs larger transactions |
| LT-T4 | `$env:RABBITMQ_BATCH_MS` | 500, 1000, 2000 | data age vs batch size under lighter load |
| LT-T5 | `$env:RABBITMQ_PREFETCH` | 500, 1000, 2000 | keep ≥ BATCH_ROWS ÷ 80 (40 values × 2 rows) |
| LT-T6 | `$env:RABBITMQ_PUBLISH_BUFFER` | 20000, 50000 | fallback count in short bursts (costs RAM) |
| LT-T7 | `$env:DB_MAX_CONNS` | 20, 40 | direct writes only; the writers have their own pool |

## 7. Optional: realistic load (not every value changes)

On a real machine most variables stay the same between readings. `load-test.ps1` always runs the worst case. For a realistic run, call the test directly (25 % of values change, 1000 machines):

```powershell
$env:TELEMETRY_LOAD='1'; $env:TELEMETRY_LOAD_PHASES='realistic'; $env:TELEMETRY_LOAD_MACHINES='1000'
$env:TELEMETRY_LOAD_CHANGE='0.25'; $env:TELEMETRY_LOAD_SECONDS='60'
cmd /c "go test ./tests -run TestTelemetryLoad$ -count=1 -v -timeout 120m 2>&1"
Remove-Item Env:TELEMETRY_LOAD, Env:TELEMETRY_LOAD_PHASES, Env:TELEMETRY_LOAD_MACHINES, Env:TELEMETRY_LOAD_CHANGE, Env:TELEMETRY_LOAD_SECONDS
```

This phase prints no `conclusion` block: read the summary table at the end of the output (`Kept up`, `Lost`, `Fallback`, `Data age` columns).

| ID | Change share | Purpose |
|---|---|---|
| LT-R1 | 0.25 | Typical plant: expected to pass well below the worst-case limit |
| LT-R2 | 0.50 | Busy plant |

## 8. Manual resilience checks (with the real server, not the script)

Run with the server (`.\run.ps1`, `RABBITMQ_ENABLED=true`) and a few connected machines, or during a long LT-08 run where noted. Check `logs/app.log` (logType `queue`) and the database.

| ID | Action | Expected |
|---|---|---|
| LT-F1 | Stop RabbitMQ for 2 min, start it again | `fallback` count rises (direct writes), no gaps in `machine_telementry`; publisher reconnects (1 s doubling up to `RABBITMQ_RECONNECT_MAX_SEC`) and the queue is used again |
| LT-F2 | Stop PostgreSQL for 1 min, start it again | Writers log `writing a batch failed, keeping it`; queue depth grows, then drains; no loss, no dead letters |
| LT-F3 | Stop the server normally while load runs | Log `unsent readings written directly at stop`; after restart nothing is missing |
| LT-F4 | Kill the server process (Task Manager) during load | Unacked messages are delivered again after restart; duplicates are skipped (COPY falls back to INSERT … ON CONFLICT, ADR 0029); row counts match |

## 9. Record the results

Copy this table into the run notes (or the ADR when a limit or setting changes). Fill in one row per run from the `conclusion` block.

| ID | Date | Settings changed | Machines | Values/s asked | Values/s stored | Fallback | Lost | Peak queue | Drain | Data age p95 | Result |
|---|---|---|---|---|---|---|---|---|---|---|---|
| LT-01 | | defaults | 10 | 400 | | | | | | | |
| LT-02 | | defaults | 50 | 2 000 | | | | | | | |
| LT-03 | | defaults | 100 | 4 000 | | | | | | | |
| LT-04 | | defaults | 250 | 10 000 | | | | | | | |
| LT-05 | | defaults | 500 | 20 000 | | | | | | | |
| LT-06 | | defaults | 750 | 30 000 | | | | | | | |
| LT-07 | | defaults | 1000 | 40 000 | | | | | | | |
| LT-08 | | defaults, 600 s | 1000 | 40 000 | | | | | | | |

**If LT-07 falls behind**: the limit is the database's disk, not the queue (ADR 0027). In order of effect: move PostgreSQL to a fast NVMe disk of its own (RabbitMQ on another disk or host), store fewer rows (realistic change share, deadband), then the tuning cases in section 6.
