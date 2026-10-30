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
The first match wins.

### Completed

| End reason | Shown as | How it is identified |
|---|---|---|
| `part_count` | Part counted | `ProgramPartCount` went up by 1 … 20 since the previous reading. A bigger jump is the counter being set on the control and is ignored. If the program is still running, the next cycle opens at once. |
| `program_loop` | Program started over | Running in both readings, `CurProgLine` fell by at least half the program's lines with `TotalProgLines` unchanged (an M99 loop, e.g. line 15 → 3). Not used once the part counter has counted for this program. The next cycle opens at once. |
| `program_end` | Program finished | Previous running, now reset, no alarm in either reading (M30 / M02). |

A counter step up to 5 s after a cycle closed (the counter can lag M30 or the loop by a
reading) is counted on that cycle. A `program_loop` cycle then becomes `part_count`.

### Unfinished

| End reason | Shown as | How it is identified | When the cycle ends |
|---|---|---|---|
| `alarm_reset` | Stopped by an alarm | Program goes to reset while an alarm or E-stop is active, in the previous or the current reading. | At the reset reading |
| `operator_reset` | Stopped by the operator | Previous held (feed hold, stop), now reset, no alarm: the operator aborted the part. | At the reset reading |
| `offline_reset` | Machine went offline | The machine was Offline (or the connection was lost) with the cycle open, and the first usable reading after it is reset. | When the machine went offline; the offline gap is not counted |
| `collector_stop` | Server restarted | Same as `offline_reset`, but the cycle was left open by a previous server run (loaded at worker start). | When the machine went offline at the restart |
| `program_changed` | Program changed | After Offline or a restart, the first usable reading runs a different `ProgName`. While connected the name is not compared: a main program cannot change while it runs, and a sub program call (M98) can show its own name. | When the machine went offline |
| `timeout` | Open too long, closed | The cycle has been open for `CNC_CYCLE_MAX_HOURS` (default 24), e.g. a machine left in Hold for days. Checked every second on the worker's tick. | At that moment |

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
