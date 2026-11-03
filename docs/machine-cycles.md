# Machine cycles

A **cycle** is one run of the part program on a machine: it starts when the program starts
running and ends when the part is finished or the run is stopped. The CNC collector records
one row per cycle in `machine_cycles` (migration `030_add_machine_cycles.sql`), decided in
`internal/monitor/machine_cycle.go` from the same DLL readings that decide the machine mode.

## One running cycle per machine

A machine has **at most one running cycle at a time**. Any number of completed and unfinished
cycles can exist beside it.

| Where | How it is enforced |
|---|---|
| Database | Unique index `machine_cycles_open_key` on `(machine_id) WHERE status = 'running'`: a second running row is refused (`repository.ErrDuplicateCycle` → `response.MachineCycleRunning`). |
| Collector | `decideCycleStep` opens a new cycle only when none is open, or in the same reading that closes the open one (close first, then open). While a cycle is open, every other reading keeps it. |
| Recovery | If saving a close failed, the old row is still `running` and the next open is refused. The worker then loads that row again (`tracker.loaded = false`) and closes it by the rules below, so a machine is never left with a stale running cycle and no new ones. |
| Restart | When a worker starts, it loads the machine's running cycle (`GetOpenCycle`) and goes on with it instead of opening a second one. Its seconds are rebuilt from `machine_mode_log`. |

## How a cycle ends

The program state comes from `ProcStatus` per engine (`machine_mode.go`): **running** (Fanuc 3
START, Siemens 2 Active), **held** (Fanuc 1/2/4, Siemens 1/3/4) or **reset** (0). An alarm is
`AlarmNo` ≠ 0 or `IsEmergency`. "Previous" is the last usable reading before the current one.
The first match wins, in the order of the rows below (`decideCycleStep`).

Names used in the **Conditions** column:

| Name | Meaning |
|---|---|
| `prev` / `now` | The previous usable reading / the current reading. `prev = none` after Offline, a lost connection or a server restart (what happened in between is unknown). |
| `program` | `running` / `held` / `reset`, from `ProcStatus` as above. |
| `alarm` | `AlarmNo ≠ 0` or `IsEmergency = 1`. |
| `count` | `ProgramPartCount`; `step = now.count − prev.count`. |
| `line` / `lines` | `CurProgLine` / `TotalProgLines`. |
| `name` | `ProgName`; `cycle.name` is the program the open cycle started with. |

Feed override (`FeedRatePercetnage`), spindle override and the operation mode (AUTO / MDI / JOG)
are **not** checked for cycles. They only decide the mode inside the cycle: running with feed
override 0 is Hold, so the cycle goes on and the time counts as `hold_sec`.

### Opens

| Conditions | Notes |
|---|---|
| `prev.program = reset` AND `now.program = running` AND no cycle open (or one closes at this reading) | A cycle opens only where its start is seen. Held → running does not open one (the held cycle goes on). A machine that connects while already running gets no cycle until its next reset → running. |
| `step` 1 … 20 or a program loop (below) AND `now.program = running` | The next cycle starts at the moment the previous one is completed, but its row is written only when the program has not reset within 5 s (`cycleStartGrace`). The counter can count a reading before M30 resets the program; that reset ends the same part, so no 0 s cycle is written. |

### Completed

| End reason | Shown as | How it is identified | Conditions |
|---|---|---|---|
| `part_count` | Part counted | `ProgramPartCount` went up by 1 … 20 since the previous reading. A bigger jump is the counter being set on the control and is ignored. If the program is still running, the next cycle opens at once. | `1 ≤ step ≤ 20` (any `program`, any `alarm`; checked first, so it wins over a reset in the same reading) |
| `program_loop` | Program started over | Running in both readings, `CurProgLine` fell by at least half the program's lines with `TotalProgLines` and `ProgName` unchanged, and the name is the cycle's own program (an M99 loop, e.g. line 15 → 3). A call into a sub program (M98) also jumps back to low lines, shows the sub program's name and can keep `TotalProgLines`. It is not a loop (machine 195: `O1002` line 6 → `O9001` line 2). Not used once the part counter has counted for this program. The next cycle opens as above. | `prev.program = running` AND `now.program = running` AND counter has not counted for this program AND `now.lines > 0` AND `now.lines = prev.lines` AND `now.name = prev.name` AND (`cycle.name = ''` OR `now.name = cycle.name`) AND `prev.line − now.line ≥ now.lines × 0.5` |
| `program_end` | Program finished | Previous running, now reset, no alarm in either reading (M30 / M02). | `prev.program = running` AND `now.program = reset` AND `prev.alarm = no` AND `now.alarm = no` AND no `step` |

A counter step up to 5 s after a cycle closed (the counter can lag M30 or the loop by a
reading) is counted on that cycle. A `program_loop` cycle then becomes `part_count`.

### Unfinished

| End reason | Shown as | How it is identified | Conditions | When the cycle ends |
|---|---|---|---|---|
| `alarm_reset` | Stopped by an alarm | Program goes to reset while an alarm or E-stop is active, in the previous or the current reading. | `prev.program ≠ reset` (running or held) AND `now.program = reset` AND (`prev.alarm = yes` OR `now.alarm = yes`) AND no `step` | At the reset reading |
| `operator_reset` | Stopped by the operator | Previous held (feed hold, stop), now reset, no alarm: the operator aborted the part. | `prev.program = held` AND `now.program = reset` AND `prev.alarm = no` AND `now.alarm = no` AND no `step` | At the reset reading |
| `offline_reset` | Machine went offline | The machine was Offline (or the connection was lost) with the cycle open, and the first usable reading after it is reset. | `prev = none` AND cycle open AND cycle not left by a previous server run AND `now.program = reset` | When the machine went offline; the offline gap is not counted |
| `collector_stop` | Server restarted | Same as `offline_reset`, but the cycle was left open by a previous server run (loaded at worker start). | `prev = none` AND cycle restored at worker start (no reading since) AND `now.program = reset` | When the machine went offline at the restart |
| `program_changed` | Program changed | After Offline or a restart, the first usable reading runs a different `ProgName`. While connected the name is not compared: a main program cannot change while it runs, and a sub program call (M98) can show its own name. | `prev = none` AND cycle open AND `now.program ≠ reset` (running or held) AND `now.name ≠ ''` AND `cycle.name ≠ ''` AND `now.name ≠ cycle.name` | When the machine went offline |
| `timeout` | Open too long, closed | The cycle has been open for `CNC_CYCLE_MAX_HOURS` (default 24), e.g. a machine left in Hold for days. Checked every second on the worker's tick. | `now − cycle.started_at ≥ CNC_CYCLE_MAX_HOURS` (any `program`, `alarm`, connection) | At that moment |

After Offline or a restart with the program still **running or held** and the same program,
the same cycle goes on and the gap counts as `offline_sec`.

### Readings that are skipped

- The DLL's first snapshot after a connect is all zeros (ProcStatus 0, no program, no counter,
  no lines). It would look like a reset and a counter drop, so it is ignored for cycles.
- A counter that goes down (reset on the control) only moves the baseline. It never ends a
  cycle.

## Where to see it

- Mode chart page `/mode-chart/`: the Cycles row on the timeline and the Cycles section (last 24 hours).
- `GET /api/v1/health/machines/{id}/cycles` (public, last 24 h), `GET /api/v1/machine-cycles`,
  `GET /api/v1/machine-cycles/{id}` (logged in, machines the user may read).
- `logs/app.log` with `CNC_LOG_CYCLES=true`: `cnc: cycle opened`, `cnc: cycle completed` /
  `unfinished` (with `endReason` and seconds per mode), `cnc: cycle restored`.

## How to verify the cycles

A cycle row is only as right as the signals it came from, so check it in three layers:
the row itself, the row against the other data the collector wrote from the same readings,
and the row against what really happened at the machine. The queries use machine `225` and
the last 24 hours; change them as needed. Mode ids: 3 In Production, 4 Hold, 5 Idle,
6 Downtime, 9 Offline.

### 1. The rows are consistent (SQL, no machine needed)

Each query should return **no rows**.

```sql
-- a) Seconds per mode add up to the cycle length (±3 s for rounding each mode)
SELECT id, end_reason, started_at, ended_at,
       EXTRACT(EPOCH FROM ended_at - started_at)::int AS length_sec,
       production_sec + hold_sec + downtime_sec + offline_sec AS mode_sec
FROM machine_cycles
WHERE machine_id = 225 AND started_at >= now() - interval '24 hours' AND status <> 'running'
  AND abs(EXTRACT(EPOCH FROM ended_at - started_at)
          - (production_sec + hold_sec + downtime_sec + offline_sec)) > 3;

-- b) No two cycles of a machine overlap (back-to-back cycles share one instant, that is fine)
SELECT a.id, b.id
FROM machine_cycles a
JOIN machine_cycles b ON b.machine_id = a.machine_id AND b.id > a.id
 AND b.started_at < COALESCE(a.ended_at, now()) AND a.started_at < COALESCE(b.ended_at, now())
WHERE a.machine_id = 225 AND a.started_at >= now() - interval '24 hours';

-- c) Parts agree with the end reason and the counter (a counter reset on the control inside a
--    cycle also shows here; check it against the control)
SELECT id, end_reason, part_count_start, part_count_end, parts_made
FROM machine_cycles
WHERE machine_id = 225 AND started_at >= now() - interval '24 hours' AND status <> 'running'
  AND ((end_reason = 'part_count' AND parts_made = 0)
       OR (end_reason = 'program_loop' AND parts_made > 0)
       OR part_count_end - part_count_start <> parts_made);

-- d) Loading time is the gap since the previous cycle
SELECT id, loading_sec,
       EXTRACT(EPOCH FROM started_at - LAG(ended_at) OVER w)::int AS gap_sec
FROM machine_cycles
WHERE machine_id = 225 AND started_at >= now() - interval '24 hours'
WINDOW w AS (PARTITION BY machine_id ORDER BY started_at)
ORDER BY started_at;  -- loading_sec should be gap_sec, or NULL across a shift / day / Offline
```

### 2. The rows agree with the mode log (SQL)

The mode and the cycle are decided from the same reading, so each end reason leaves a known
pattern in `machine_mode_log`:

```sql
SELECT c.id, c.end_reason, c.started_at, c.ended_at,
  (SELECT l.mode_id FROM machine_mode_log l WHERE l.machine_id = c.machine_id
     AND l.started_at <= c.started_at - interval '1 second'
     AND (l.ended_at IS NULL OR l.ended_at > c.started_at - interval '1 second')
   ORDER BY l.started_at DESC LIMIT 1) AS mode_before_start,
  (SELECT l.mode_id FROM machine_mode_log l WHERE l.machine_id = c.machine_id
     AND l.started_at <= c.started_at + interval '1 second'
     AND (l.ended_at IS NULL OR l.ended_at > c.started_at + interval '1 second')
   ORDER BY l.started_at DESC LIMIT 1) AS mode_after_start,
  (SELECT l.mode_id FROM machine_mode_log l WHERE l.machine_id = c.machine_id
     AND l.started_at <= c.ended_at - interval '1 second'
     AND (l.ended_at IS NULL OR l.ended_at > c.ended_at - interval '1 second')
   ORDER BY l.started_at DESC LIMIT 1) AS mode_before_end,
  (SELECT l.mode_id FROM machine_mode_log l WHERE l.machine_id = c.machine_id
     AND l.started_at <= c.ended_at + interval '1 second'
     AND (l.ended_at IS NULL OR l.ended_at > c.ended_at + interval '1 second')
   ORDER BY l.started_at DESC LIMIT 1) AS mode_after_end
FROM machine_cycles c
WHERE c.machine_id = 225 AND c.started_at >= now() - interval '24 hours'
ORDER BY c.started_at;
```

| End reason | Mode before end | Mode after end |
|---|---|---|
| `program_end` | 3 (4 if feed override was 0 at M30) | 5 Idle |
| `operator_reset` | 4 Hold | 5 Idle |
| `alarm_reset` | 6, 4 or 3 | 6 Downtime (5 if the alarm cleared at the reset reading) |
| `part_count`, `program_loop` | 3 or 4 | 3 / 4 (still running, next cycle) or 5 |
| `offline_reset`, `collector_stop`, `program_changed` | 3, 4 or 6 | 9 Offline (the cycle ends when the machine went offline) |
| `timeout` | any, usually 4 | any |

At the start: `mode_after_start` is 3 or 4 (program running), `mode_before_start` is 5 Idle,
or 3 / 4 when the cycle follows a `part_count` / `program_loop` back to back. A cycle whose
pattern does not match is the one to look at in step 3.

### 3. The raw signals around one cycle (telemetry and logs)

`machine_telementry` stores every change of a numeric variable, so when the machine has the
variables `55` (ProcStatus), `201` (ProgramPartCount) and `204` (CurProgLine) in
`machine_variables`, the transitions that decided a cycle can be read back:

```sql
-- Which of the cycle signals the machine stores
SELECT id FROM machine_variables WHERE machine_id = 225 AND id IN (55, 201, 204);

-- Every change of those signals from 10 s before the cycle to 10 s after it
SELECT t.recorded_at, t.variable_id, t.value
FROM machine_telementry t, machine_cycles c
WHERE c.id = :cycle_id AND t.machine_id = c.machine_id AND t.variable_id IN (55, 201, 204)
  AND t.recorded_at BETWEEN c.started_at - interval '10 seconds'
                        AND COALESCE(c.ended_at, now()) + interval '10 seconds'
ORDER BY t.recorded_at;
```

Read the rows with the **Conditions** tables above: e.g. a `program_end` must show ProcStatus
going from the running value (Fanuc 3, Siemens 2) to 0 at `ended_at`; a `part_count` must show
201 stepping up by 1 … 20 at `ended_at`; an `operator_reset` must show a held value (Fanuc 1/2/4,
Siemens 1/3/4) just before the 0.

Alarms are not a stored variable, so for the alarm side (and when the variables are not
assigned) follow the machine in the log: set `CNC_MACHINE_IDS=<id>`, `CNC_LOG_SNAPSHOT=true`
(ProcStatus, AlarmNo, IsEmergency, FeedRatePercetnage, ProgName by name), `CNC_LOG_VARIABLES=true`
(`201:partCount 262 -> 263`, CurProgLine), `CNC_LOG_CYCLES=true` and `CNC_LOG_MODES=true`, then
read `logs/app.log` lines starting `[machine <id>`: every `cnc: cycle completed` / `unfinished`
should follow the snapshot / variable change its conditions name, at the same time.

### 4. At the machine (floor test)

Run the tests in [Floor tests on the HUMA panel](#floor-tests-on-the-huma-panel-fanuc) below.
Over a full shift, also compare the count of `completed` cycles and the sum of `parts_made` with
the control's part counter and the parts in the bin, and a few cycle lengths with a stopwatch. A
difference points to a missed transition: find the cycle around it and read its signals as in
step 3.

## Floor tests on the HUMA panel (Fanuc)

Step-by-step tests for a tester at the Jyoti **HUMA** panel of a Fanuc machine. Each test says
what to press, what the panel shows, and which `machine_cycles` row should appear. Together they
cover every row of the *Opens*, *Completed* and *Unfinished* tables above.

### HUMA buttons used

| HUMA button / area | What it does in these tests | Program state |
|---|---|---|
| **AUTO** | Selects memory (auto) mode. Run every test in AUTO. | — |
| **FEED START** (yellow, under AUTO) | Cycle start: starts or continues the program. | running (status `START`) |
| **CYCLE STOP** (red, top left) | Feed hold: pauses the program. | held (status `HOLD` / `STOP`) |
| **RESET** (top left) | Ends the program and clears the alarm. | reset |
| **M01** | Optional stop ON: the program stops at `M01` (status bar shows `Stop: M0/M1 active`). | held |
| **SBK** | Single block: the program stops after every line. | held between lines |
| **Feed Rate** dial (left) | Feed override, 0 … 120 %. | unchanged (the mode becomes Hold at 0 %) |
| Status bar (top left, next to the green icon) | Shows `START` while running, `HOLD` / `STOP` while held, and the active stop or alarm. | — |
| **Time / Count / Progress** box | `Cycle Time` and `Part Count`, e.g. `262 / 200`: the first number is `ProgramPartCount`. | — |
| **Program** / **File Browser** tabs | Create and select the test programs. | — |

### Before the tests (once)

1. In `.env` set the following and restart the server (`.\run.ps1`):
   ```
   CNC_MACHINE_IDS=<machine_id>
   CNC_LOG_CYCLES=true
   CNC_LOG_MODES=true
   CNC_LOG_SNAPSHOT=true
   ```
2. Open `/mode-chart/?machineId=<machine_id>`. The machine must be online (not Offline).
3. On the HUMA, create these programs (**Program** tab). They only dwell (`G4`): **the axes do
   not move**, so they are safe to run without a part.

   ```
   <CYCLE_END>        <CYCLE_LOOP>       <CYCLE_LONG>       <CYCLE_B>          <CYCLE_ALARM>
   G4X10.             G4X10.             G4X300.            G4X300.            G4X10.
   G4X10.             G4X10.             M30                M30                #3000=1(CYCLE TEST ALARM)
   M01                G4X10.                                                   M30
   G4X10.             M99
   M30
   ```

   - `CYCLE_END`: about 30 s, ends with M30 (stops at `M01` only when the **M01** button is ON).
   - `CYCLE_LOOP`: about 30 s per loop, starts over for ever (M99), like the `JYOTI` program.
   - `CYCLE_LONG` / `CYCLE_B`: 5 minutes, for the offline and restart tests.
   - `CYCLE_ALARM`: raises macro alarm 3001 after 10 s. If custom macros are not enabled on the
     control, press the machine's **E-stop** button instead (TEST-11).
4. Before **every** test: press **RESET**, select **AUTO**, set **Feed Rate** to 100 %, set
   **M01** and **SBK** to OFF, and select the program the test names.

### Check query

After each test, run this query (or look at the Cycles section on the mode chart page). The
newest row is at the top.

```sql
SELECT id, program_name, status, end_reason,
       to_char(started_at, 'HH24:MI:SS') AS started,
       to_char(ended_at, 'HH24:MI:SS') AS ended,
       EXTRACT(EPOCH FROM COALESCE(ended_at, now()) - started_at)::int AS length_sec,
       production_sec, hold_sec, downtime_sec, offline_sec,
       loading_sec, part_count_start, part_count_end, parts_made
FROM machine_cycles
WHERE machine_id = <machine_id>
ORDER BY started_at DESC
LIMIT 5;
```

Check the seconds only on **closed** rows. A running row writes its seconds every 5 minutes, so
until then they show 0. Seconds may be off by 1–2 s. Note the clock time (`HH:MM:SS`) when you
press a button, so you can compare it with `started` / `ended`.

### Test list

| Test | Name | Covers | Expected row(s) |
|---|---|---|---|
| TEST-01 | Start a cycle | Opens: reset → running | `running`, no end reason |
| TEST-02 | Normal end (M30) | `program_end` / `part_count` | `completed` |
| TEST-03 | Three parts one after another | `part_count`, `loading_sec` | 3 × `completed` |
| TEST-04 | Program loop (M99) | `program_loop`, next opens at once | 3 × `completed` + 1 × `unfinished` |
| TEST-05 | Feed hold, then continue | held → running keeps the cycle | 1 × `completed`, `hold_sec` ≈ 30 |
| TEST-06 | Feed override 0 % | feed is not a cycle condition | 1 × `completed`, `hold_sec` ≈ 20 |
| TEST-07 | Optional stop M01 | held at M01 keeps the cycle | 1 × `completed`, `hold_sec` ≈ 20 |
| TEST-08 | Single block | held between lines keeps the cycle | 1 × `completed` |
| TEST-09 | Feed hold, then RESET | `operator_reset` | `unfinished` |
| TEST-10 | RESET while running | running → reset, no alarm | `completed` / `program_end` |
| TEST-11 | Alarm while running | `alarm_reset` | `unfinished` |
| TEST-12 | Alarm while held | `alarm_reset` (held) | `unfinished` |
| TEST-13 | Offline, program still running | same cycle goes on | 1 × `completed`, `offline_sec` > 0 |
| TEST-14 | Offline, RESET while offline | `offline_reset` | `unfinished` |
| TEST-15 | Offline, other program started | `program_changed` | `unfinished` |
| TEST-16 | Server restart, program still running | restored cycle goes on | 1 × `completed`, `offline_sec` > 0 |
| TEST-17 | Server stopped, RESET, server started | `collector_stop` | `unfinished` |
| TEST-18 | Server started while program runs | a cycle opens only where its start is seen | no row until the next start |
| TEST-19 | Cycle open too long (optional, 1 h) | `timeout` | `unfinished` |
| TEST-20 | Part counter changed by hand (optional) | a counter jump > 20 is ignored | no extra row |

---

#### TEST-01 · Start a cycle

**Do on HUMA**
1. Select `CYCLE_END`, press **FEED START**.

**You see on HUMA:** status `START`, `Cycle Time` counting.

**Expected:** within 1–2 s a new row: `program_name = CYCLE_END`, `status = running`,
`end_reason` empty, `ended` empty, `part_count_start` = the panel's Part Count. Mode chart: In
Production. Log: `cnc: cycle opened`.

**Pass if** exactly one new `running` row appears. Let the program finish and go on with TEST-02.

#### TEST-02 · Normal end (M30)

**Do on HUMA**
1. Select `CYCLE_END`, press **FEED START**, wait about 30 s until the program ends by itself.

**You see on HUMA:** status back to reset. `Part Count` goes up by 1 if the control counts at M30.

**Expected:** the row is closed: `status = completed`, `length_sec` ≈ 30, `production_sec` ≈ 30,
`hold_sec = downtime_sec = offline_sec = 0`.
- Part Count went up: `end_reason = part_count` (or `program_end`), `parts_made = 1`,
  `part_count_end = part_count_start + 1`.
- Part Count did not move: `end_reason = program_end`, `parts_made = 0`.

Log: `cnc: cycle completed`. Mode chart: Idle after the end.

#### TEST-03 · Three parts one after another

**Do on HUMA**
1. Run `CYCLE_END` to the end (FEED START, wait for the end).
2. Wait about 10 s (this is the "loading" time), press **FEED START** again.
3. Repeat until it has run 3 times.

**Expected:** 3 `completed` rows of about 30 s each. The 2nd and 3rd rows have `loading_sec` ≈ 10
(the wait before you pressed FEED START). If the control counts at M30, `parts_made = 1` on each
row and the panel's Part Count went up by 3.

#### TEST-04 · Program loop (M99)

**Do on HUMA**
1. Select `CYCLE_LOOP`, press **FEED START**.
2. Let it loop 3 times (about 90 s). Each loop starts over at the first line.
3. Press **CYCLE STOP**, then **RESET**.

**Expected:** 4 rows, all `CYCLE_LOOP`:
- 3 × `completed` / `program_loop`, each about 30 s long. Each one's `ended` is the next one's
  `started` (no gap, `loading_sec` = 0).
- 1 × `unfinished` / `operator_reset`: the 4th loop, which you stopped.

If this control counts parts at M99, you see `part_count` instead of `program_loop`; that is also
a pass.

#### TEST-05 · Feed hold, then continue

**Do on HUMA**
1. Select `CYCLE_END`, press **FEED START**.
2. After about 5 s press **CYCLE STOP** (status `HOLD`), wait 30 s.
3. Press **FEED START** again and let it finish.

**Expected:** **one** row only (the hold does not end or start a cycle): `completed`,
`hold_sec` ≈ 30, `production_sec` ≈ 30, `length_sec` ≈ 60. Mode chart: Hold for 30 s inside the
cycle.

#### TEST-06 · Feed override 0 %

**Do on HUMA**
1. Select `CYCLE_END`, press **FEED START**.
2. Turn the **Feed Rate** dial to 0 %, wait 20 s, turn it back to 100 %. Let it finish.

**Expected:** one `completed` row with `hold_sec` ≈ 20 (running with 0 % feed counts as Hold, but
the cycle goes on). The dwell does not wait for the feed, so `length_sec` stays about 30 and
`production_sec` ≈ 10.

#### TEST-07 · Optional stop M01

**Do on HUMA**
1. Turn **M01** ON. Select `CYCLE_END`, press **FEED START**.
2. After about 20 s the program stops at `M01`: the status bar shows `Stop: M0/M1 active`.
3. Wait 20 s, press **FEED START**, let it finish. Turn **M01** OFF.

**Expected:** one `completed` row, `hold_sec` ≈ 20, `length_sec` ≈ 50.

#### TEST-08 · Single block

**Do on HUMA**
1. Turn **SBK** ON. Select `CYCLE_END`, press **FEED START**.
2. Each time the program stops after a line, wait about 5 s and press **FEED START**, until it
   ends. Turn **SBK** OFF.

**Expected:** **one** `completed` row for the whole program (not one per line), with the waits
in `hold_sec`.

#### TEST-09 · Feed hold, then RESET (operator aborts the part)

**Do on HUMA**
1. Select `CYCLE_END`, press **FEED START**.
2. After about 10 s press **CYCLE STOP**, wait 5 s, press **RESET**.

**Expected:** `unfinished` / `operator_reset`, `ended` = when you pressed RESET,
`production_sec` ≈ 10, `hold_sec` ≈ 5, `parts_made = 0`. Log: `cnc: cycle unfinished`.

#### TEST-10 · RESET while running (no hold)

**Do on HUMA**
1. Select `CYCLE_END`, press **FEED START**.
2. After about 10 s press **RESET** directly, **without** CYCLE STOP.

**Expected:** `completed` / `program_end`, `length_sec` ≈ 10. This is the current rule: running →
reset without an alarm looks the same as M30 in the DLL data, so it counts as a completed cycle.
To abort a part as unfinished, the operator must press CYCLE STOP first (TEST-09).

#### TEST-11 · Alarm while running

**Do on HUMA**
1. Select `CYCLE_ALARM`, press **FEED START**.
2. After 10 s alarm 3001 `CYCLE TEST ALARM` appears (bell / alarm bar) and the program stops.
   (Without custom macros: run `CYCLE_LONG` and press the machine's **E-stop** after 10 s.)
3. Wait 10 s, press **RESET** (release the E-stop first, if you used it).

**Expected:** `unfinished` / `alarm_reset`, `production_sec` ≈ 10. On Fanuc the macro alarm and
the E-stop both reset the program at once, so `ended` = when the alarm appeared / the E-stop was
pressed (not when you pressed RESET) and `downtime_sec` ≈ 0. Mode chart: Downtime after the
cycle, until RESET. If `ended` is the RESET time instead, the control held the program in alarm
first; the time in between is then in `downtime_sec`. Both are a pass.

#### TEST-12 · Alarm while held

**Do on HUMA**
1. Select `CYCLE_LONG`, press **FEED START**.
2. After 10 s press **CYCLE STOP** (status `HOLD`).
3. Press the machine's **E-stop**, wait 10 s, release it, press **RESET**.

**Expected:** `unfinished` / `alarm_reset` (not `operator_reset`: the alarm wins),
`production_sec` ≈ 10, `hold_sec` = the time from CYCLE STOP to the E-stop. `ended` = when the
E-stop was pressed (it resets the program).

#### TEST-13 · Offline, program still running

**Do on HUMA**
1. Select `CYCLE_LONG`, press **FEED START**.
2. After 20 s disconnect the machine's network (LAN cable or switch port). Wait until the mode
   chart shows the machine **Offline** (up to about 30 s), then wait 30 s more.
3. Reconnect the network **while the program is still running** (status `START`). Let it finish.

**Expected:** **one** row (the same cycle goes on): `completed`, `offline_sec` ≈ the time it showed
Offline, `length_sec` ≈ 300.

#### TEST-14 · Offline, RESET while offline

**Do on HUMA**
1. Select `CYCLE_LONG`, press **FEED START**.
2. After 20 s disconnect the network and wait until the mode chart shows **Offline**.
3. On HUMA press **RESET**. Wait 30 s, reconnect the network.

**Expected:** `unfinished` / `offline_reset`. `ended` = when the machine **went Offline** (not when
you reconnected). The offline gap is not counted: `offline_sec` ≈ 0, `length_sec` ≈ 20 + the
detection time.

#### TEST-15 · Offline, other program started

**Do on HUMA**
1. Select `CYCLE_LONG`, press **FEED START**.
2. After 20 s disconnect the network and wait until the mode chart shows **Offline**.
3. On HUMA press **RESET**, select `CYCLE_B`, press **FEED START**.
4. Reconnect the network **while `CYCLE_B` is running**.

**Expected:** the `CYCLE_LONG` row is `unfinished` / `program_changed`, `ended` = when it went
Offline. No row for `CYCLE_B` appears: its start was not seen (the machine was offline). A
`CYCLE_B` row appears only the next time it is started from reset.

#### TEST-16 · Server restart, program still running

**Do on HUMA**
1. Select `CYCLE_LONG`, press **FEED START**.
2. After 20 s stop the server (Ctrl+C in the `run.ps1` window). Wait 30 s, start it again.
3. Let the program finish.

**Expected:** log `cnc: cycle restored` after the start, then **one** row: `completed`,
`offline_sec` ≈ the time the server was down, `length_sec` ≈ 300.

#### TEST-17 · Server stopped, RESET, server started

**Do on HUMA**
1. Select `CYCLE_LONG`, press **FEED START**.
2. After 20 s stop the server. On HUMA press **RESET**. Start the server again.

**Expected:** `unfinished` / `collector_stop`. `ended` = when the server was stopped (the last
reading before it), not when it started again.

#### TEST-18 · Server started while the program runs

**Do on HUMA**
1. Stop the server. Make sure no row is `running` for the machine (check query).
2. Select `CYCLE_END`, press **FEED START**, start the server while it runs.

**Expected:** **no** row for this run (its start was not seen). Run `CYCLE_END` once more from
reset: that run gets a normal row.

#### TEST-19 · Cycle open too long (optional, takes 1 hour)

**Do on HUMA**
1. Set `CNC_CYCLE_MAX_HOURS=1` in `.env`, restart the server.
2. Turn **M01** ON, run `CYCLE_END`. It stops at `M01` (held). Leave it there for more than 1 hour.
3. Press **RESET**, turn **M01** OFF, set `CNC_CYCLE_MAX_HOURS=24` again and restart the server.

**Expected:** `unfinished` / `timeout`, `ended` ≈ `started` + 1 h, `hold_sec` ≈ 3580. Pressing RESET
afterwards does not add or change a row.

#### TEST-20 · Part counter changed by hand (optional)

**Do on HUMA**
1. Run `CYCLE_LONG`. While it runs, change the part count on the control's setting screen by
   more than 20 (e.g. 262 → 0, or 262 → 300).
2. Let it finish.

**Expected:** the counter change does not end the cycle: one `completed` row for the whole run.
`parts_made` does not include the jump.

### Result sheet

Copy this table, fill it in during the tests and attach the check query output.

| Test | Time started | Row id(s) | status / end_reason seen | Seconds OK? | Pass / Fail | Notes |
|---|---|---|---|---|---|---|
| TEST-01 | | | | | | |
| TEST-02 | | | | | | |
| TEST-03 | | | | | | |
| TEST-04 | | | | | | |
| TEST-05 | | | | | | |
| TEST-06 | | | | | | |
| TEST-07 | | | | | | |
| TEST-08 | | | | | | |
| TEST-09 | | | | | | |
| TEST-10 | | | | | | |
| TEST-11 | | | | | | |
| TEST-12 | | | | | | |
| TEST-13 | | | | | | |
| TEST-14 | | | | | | |
| TEST-15 | | | | | | |
| TEST-16 | | | | | | |
| TEST-17 | | | | | | |
| TEST-18 | | | | | | |
| TEST-19 | | | | | | |
| TEST-20 | | | | | | |

When a test fails, keep the row id and the time, and read that cycle's signals as in step 3 of
*How to verify the cycles* (`logs/app.log` lines `[machine <machine_id> …` around that time).
