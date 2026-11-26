# Telemetry load test results

Results of `load-test.ps1` (the RabbitMQ telemetry queue, ADR 0024). How to plan test cases: [telemetry-load-test-plan.md](telemetry-load-test-plan.md). Why the settings are what they are: [ADR 0030](adr/0030-gzip-telemetry-queue-messages.md).

**Add every new test session as a new entry at the bottom of [Results log](#results-log), newest last.** Keep the same table columns, so sessions can be compared.

## How to run

Nothing to edit in `load-test.ps1`. It reads `.env` and passes the test settings itself.

```powershell
.\load-test.ps1 -Values 40000 -Seconds 60      # 1 minute, every value changing (worst case)
.\load-test.ps1 -Values 40000 -Seconds 300     # 5 minutes: the real check that nothing piles up
.\load-test.ps1 -Values 40000 -Seconds 60 -Change 50   # realistic: half the values change per reading
```

| Parameter | Meaning |
|---|---|
| `-Values` | telemetry values per second. Machines = Values / 40 (40 variables each, 1 reading/s) |
| `-Seconds` | how long to send (default 120) |
| `-Change` | % of values that change per reading (100 = worst case) |
| `-Rate` | readings per machine per second (default 1) |

Each run saves its full output to `logs\load-test-<values>[-c<change>]-<date>-<time>.log`.

**Before a session:**
- Close heavy disk users (builds, copies). This PC's SSD is the limit (see Conclusion).
- If an earlier run was killed, delete its left-over queues `telemetry.load.*` in the RabbitMQ management UI. A stopped run cannot remove them, and they hold messages on the same disk.

### `.env` settings these results used

```
DB_MAX_CONNS=20
RABBITMQ_ENABLED=true
RABBITMQ_CONSUMERS=4
RABBITMQ_PREFETCH=1000
RABBITMQ_BATCH_ROWS=20000
RABBITMQ_BATCH_MS=1000
RABBITMQ_PUBLISH_BUFFER=60000
RABBITMQ_LIVE_EVERY_N=4
# RABBITMQ_QUEUE_TYPE=quorum (default)
```

### How to read the result

| Line in the output | Good | Means |
|---|---|---|
| `RESULT` | `KEPT UP` | every value stored as fast as it was sent |
| `Stored … (lost N)` | lost 0 | values that never reached `machine_telementry` |
| `Queue … empty X after sending stopped` | under 5 s | backlog left when sending ended. More = it fell behind |
| `Path … N written directly` | 0 | readings the publish buffer could not take |
| `Data age p95` | a few seconds | how old the newest stored value was while sending |

Only `machine_telementry` values count. `machine_variable_values` upserts still run, but they are not part of the throughput. The test dates its rows 50 years ahead and drops those chunks afterwards, so the real data is untouched.

## Test PC

| | |
|---|---|
| CPU | Intel i5-8400, 6 cores / 6 threads |
| RAM | 16 GB (PostgreSQL `shared_buffers` ≈ 5 GB) |
| Disk | Kingston A400 240 GB SATA SSD, no DRAM cache, 87 % full (29 GB free) |
| Software | PostgreSQL 18 + TimescaleDB, RabbitMQ 4.2.2 (quorum queue), all on this one PC |

## Results log

### 2026-11-26: 40,000 values/s target (ADR 0030)

**Goal:** sustain 40,000 values/s (1,000 machines × 40 variables × 1 reading/s, every value changing) with no RabbitMQ backlog.

**Before:** messages were plain JSON, 2.6 KB each; publish buffer 20,000.

| Run | Values/s | Length | Stored | Lost | Queue empty after | Direct writes | Data age p95 | Result |
|---|---|---|---|---|---|---|---|---|
| Earlier sessions (2026-11-25) | 40,000 | 60 s | 29k–39k/s | 0 | 3.5–27 s | 0–1,167 | 5–28 s | ⚠️ about half fell behind |
| Baseline | 40,000 | 60 s | 37,837/s | 0 | 3.5 s | 0 | 5.1 s | ✅ |
| Clean database | 40,000 | 60 s | 35,614/s | 0 | 7 s | 900 | 23 s | ❌ SSD stall |
| 2× load | 80,000 | 60 s | 44,770/s | 0 | 47 s | 3,266 | 54 s | ❌ |
| 5 minutes | 40,000 | 300 s | 39,799/s | 0 | 1.6 s | 5,911 | 35 s | ❌ buffer overflowed in an SSD stall |

**What was measured (where the time went):**

| Part | Measurement | Bottleneck? |
|---|---|---|
| PostgreSQL COPY | ~100 ms per 16,000 rows. Writers busy ~20 % of the time | No |
| PostgreSQL alone | 40,000 rows/s for 60 s, ~8 MB/s, no slow writes | No |
| RabbitMQ alone | ~8,000–9,000 messages/s (9× the need) | No |
| CPU / DB pool | ~33 % CPU, 4 busy connections of 24 open | No |
| **SSD** | writes froze for 5–40 s (1–8 s per write instead of ~1 ms). RabbitMQ wrote every message twice: 2.8 MB/s as ~2,050 small writes/s | **Yes** |

**Changes:**
1. Queue messages are gzipped: 3,044 → 505 bytes (83 % smaller). RabbitMQ disk writes fell from 2.8 MB/s and ~2,050 writes/s to 0.8 MB/s and ~80 writes/s.
2. `RABBITMQ_PUBLISH_BUFFER` 20,000 → 60,000: a minute of readings (~30 MB) rides out an SSD stall in memory.
3. `RABBITMQ_QUEUE_TYPE` is now read from `.env`. It was ignored before.
4. Load test: counts telemetry only, uses its own future-dated chunks (no more mass DELETEs slowing the next run), and uses realistic 3-decimal values.

**After:**

| Run | Values/s | Length | Stored | Lost | Queue empty after | Direct writes | Data age p95 | Result |
|---|---|---|---|---|---|---|---|---|
| **5 minutes** | **40,000** | **300 s** | **39,776/s (12,000,000 of 12,000,000)** | **0** | **1.8 s** | **0** | 12 s | ✅ **kept up** |
| 60 s, run 1 | 40,000 | 60 s | 39,706/s | 0 | 0.5 s | 0 | 2.0 s | ✅ |
| 60 s, run 2 | 40,000 | 60 s | 39,247/s | 0 | 1.2 s | 0 | 17 s | ✅ |
| 60 s, run 3 | 40,000 | 60 s | 39,176/s | 0 | 1.4 s | 0 | 2.1 s | ✅ |
| 60 s, run 4 | 40,000 | 60 s | 39,284/s | 0 | 1.2 s | 0 | 2.1 s | ✅ |
| 60 s, runs 5–7 | 40,000 | 60 s | 28k–33k/s | 0 | 13–25 s | 0 | 11–25 s | ⚠️ SSD stall, then caught up |

**Conclusion:**

| Question | Answer |
|---|---|
| Does 40,000 values/s work? | **Yes.** 5-minute run and 4 of 7 one-minute runs kept up; every run stored every value |
| Why do some runs still fall behind? | The SSD freezes for 5–10 s (nearly full, no DRAM cache). The queue catches up within 13–25 s after |
| Anything lost? | No, 0 values lost in every run |
| Bottleneck now | The disk only. PostgreSQL writes ~11–13 MB/s, close to the minimum for 40,000 rows/s |
| To make it pass every time | Free space on C: or move the PostgreSQL + RabbitMQ data to an NVMe SSD with DRAM. Add Windows Defender exclusions for both data folders |
| Storage warning | 40,000 values/s writes ~7.5 GB per hour before compression: 29 GB free lasts ~4 hours |

### 2026-11-26: verification run, 5 minutes at 40,000 values/s

**Goal:** check the ADR 0030 changes on a fresh 5-minute run, with the settings above.
Command: `.\load-test.ps1 -Values 40000 -Seconds 300`. Log: `logs\load-test-40000-20261126-111230.log`.

| Run | Values/s | Length | Stored | Lost | Queue empty after | Direct writes | Data age p95 | Result |
|---|---|---|---|---|---|---|---|---|
| 5 minutes | 40,000 | 300 s | 39,859/s (12,000,000 of 12,000,000) | 0 | 1.2 s | 0 | 7.9 s (max 12.6 s) | ✅ kept up |

**Minute by minute** (the `queue: statistics` lines; 60,000 readings arrive every minute):

| Minute | Readings written | Avg batch write | Writers busy | What happened |
|---|---|---|---|---|
| 1 | 57,903 | 117 ms | ~12 % | normal |
| 2 | 59,945 | **368 ms** | ~42 % | short SSD slowdown: writes took 3× longer, the queue peaked at **5,360 readings** (~5 s), then caught up in the same minute |
| 3 | 59,909 | 102 ms | ~11 % | normal |
| 4 | 59,906 | 110 ms | ~11 % | normal |
| 5 | 60,496 | 223 ms | ~23 % | smaller slowdown, absorbed |

**Conclusion:**

| Question | Answer |
|---|---|
| Any issue? | **No.** Every value stored, 0 lost, 0 direct writes, queue empty 1.2 s after sending stopped |
| Peak queue 5,360? | Not a problem: a ~5 s backlog during one SSD slowdown in minute 2, drained within that minute. The buffer (60,000) and the queue absorbed it |
| Still a bottleneck? | Only the SSD, now and then. Even in the slow minute the writers were busy ~42 % of the time (normally ~12 %), so the software has more than 2× headroom |
| What would remove the rest | A faster or emptier disk for PostgreSQL + RabbitMQ (see the previous conclusion). No code or `.env` change is needed for 40,000 values/s |
| Next test | Find the real ceiling: `.\load-test.ps1 -Values 60000 -Seconds 300`. If it keeps up, 40,000 has 50 % headroom |

### 2026-11-26: ceiling check, 5 minutes at 60,000 values/s

**Goal:** find how much headroom 40,000 has. Command: `.\load-test.ps1 -Values 60000 -Seconds 300` (1,500 machines). Log: `logs\load-test-60000-20261126-112805.log`.

| Run | Values/s | Length | Stored | Lost | Queue empty after | Direct writes | Data age p95 | Result |
|---|---|---|---|---|---|---|---|---|
| 5 minutes | 60,000 | 300 s | 55,264/s (18,000,000 of 18,000,000) | 0 | 25.8 s | 0 | 20.2 s (max 26 s) | ❌ fell behind (small) |

**Minute by minute** (90,000 readings arrive every minute):

| Minute | Readings written | Values/s stored | Avg batch write | Writers busy | Backlog at end of minute |
|---|---|---|---|---|---|
| 1 | 87,000 | 58,000 | 126 ms | 21 % | 0 |
| 2 | 90,566 | 60,377 | 132 ms | 23 % | 0 |
| 3 | 72,131 | 48,087 | 174 ms | 26 % | 14,353 (SSD slowdown) |
| 4 | 101,566 | 67,711 | 163 ms | 31 % | 6,300 (catching up) |
| 5 | 89,019 | 59,346 | 254 ms | 45 % | 7,350 (second slowdown) |

**Conclusion:**

| Question | Answer |
|---|---|
| Does 60,000 values/s work? | Between SSD slowdowns, yes (minutes 1–2). Not reliably: it catches up at only ~67,700/s, ~13 % above what arrives, so one slowdown takes minutes to clear and a second one leaves a backlog |
| Anything lost? | No: 0 lost, 0 direct writes |
| Bottleneck | The SSD again. Writers ≤ 45 % busy, CPU fine; RabbitMQ confirms and PostgreSQL writes both slowed together |
| Would more writers / bigger batches help? | No: the writers are waiting on the disk, not working |
| Limits on this PC | **40,000 = reliable** (~70 % catch-up headroom). **~55,000 = average ceiling**. **60,000 = not reliable** |
| For 60,000+ | A faster or emptier disk (NVMe with DRAM) for PostgreSQL + RabbitMQ; no code change needed |

<!-- Next session: copy the "### <date>: <goal>" block above, fill in the tables, add its conclusion. -->
