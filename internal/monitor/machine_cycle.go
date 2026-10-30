package monitor

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"time"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

// A cycle is one run of the part program (machine_cycles, migration 030). It is followed from
// the same DLL readings that decide the mode:
//
//	reset -> running                       a cycle opens
//	ProgramPartCount steps up              completed (part_count); still running: the next opens
//	CurProgLine jumps back to the start    completed (program_loop), the next opens at once; only
//	  while running                        while the counter has not counted for this program
//	running -> reset, no alarm             completed (program_end)
//	alarm or E-stop -> reset               unfinished (alarm_reset)
//	held -> reset                          unfinished (operator_reset)
//	open past CNC_CYCLE_MAX_HOURS          unfinished (timeout)
//	after Offline / a restart:
//	  reset                                unfinished (offline_reset / collector_stop)
//	  another program                      unfinished (program_changed)
//	  running or held, same program        the same cycle goes on, the gap counts as offline
//
// A cycle closed after Offline ends when the machine went offline; the gap is not counted. The
// program name is compared only after a gap: a main program cannot change while it runs, and
// a sub program call (M98) may show its own name. A counter step within cycleCounterGrace
// after a cycle closed is that cycle's part (the counter can lag M30 / the loop by a reading).
//
// The time of each mode inside the cycle is added up in memory as the mode changes and
// written every cycleCheckpointInterval and when the cycle closes; a cycle a previous run left
// open is rebuilt from machine_mode_log when the worker starts.

const (
	// cycleCheckpointInterval is how often a running cycle's seconds are written.
	cycleCheckpointInterval = 5 * time.Minute
	// loopMinimumDropFraction: CurProgLine falling by at least this share of the program's
	// lines while running is the program starting over (machine 195: line 15 -> 3 of 15).
	// Smaller jumps back are loops inside the program.
	loopMinimumDropFraction = 0.5
	// maxPartCountStep: a bigger jump of ProgramPartCount is the counter being set on the
	// control, not parts made; it only moves the baseline.
	maxPartCountStep = 20
	// cycleCounterGrace: a counter step this soon after a cycle closed is that cycle's part.
	cycleCounterGrace = 5 * time.Second
)

// cycleReading is the part of a DLL reading a cycle is followed from.
type cycleReading struct {
	at          time.Time
	program     programState
	programName string
	partCount   int
	programLine int
	totalLines  int
	alarm       bool
}

// cycleReadingFrom takes the cycle values out of reading, whose NC state is state. ok is false
// for the DLL's first snapshot of a connection, which is all zeros (no program, no counter, no
// lines) until the real values follow a second later: it would look like a reset and a counter
// drop.
func cycleReadingFrom(state machineState, reading cnc.Reading) (cycleReading, bool) {
	current := cycleReading{
		at:          reading.ReceivedAt,
		program:     state.program,
		programName: state.programName,
		partCount:   int(reading.Values[cnc.FieldProgramPartCount].Number),
		programLine: int(reading.Values[cnc.FieldCurProgLine].Number),
		totalLines:  int(reading.Values[cnc.FieldTotalProgLines].Number),
		alarm:       state.isAlarmActive(),
	}
	emptySnapshot := state.procStatus == 0 && current.programName == "" && current.partCount == 0 && current.totalLines == 0
	return current, !emptySnapshot
}

// cycleStep is what one reading does to the machine's cycle.
type cycleStep struct {
	closeReason  string // close the open cycle for this reason ("" = keep it)
	open         bool   // open a new cycle at the reading
	openedByLoop bool   // the new cycle starts at a program loop
	parts        int    // parts the counter counted at this reading
}

// decideCycleStep returns what current does, given the last usable reading before it (nil after
// Offline or a restart, when what happened in between is unknown), the open cycle's program
// (hasOpen false when none) and whether the part counter has counted for this program. First
// match wins; a cycle opens only where its start is seen.
func decideCycleStep(previous *cycleReading, current cycleReading, hasOpen bool, openProgram string, counterCounts bool) cycleStep {
	if previous == nil {
		switch {
		case !hasOpen:
			return cycleStep{}
		case current.program == programReset:
			return cycleStep{closeReason: model.CycleEndOfflineReset}
		case current.programName != "" && openProgram != "" && current.programName != openProgram:
			return cycleStep{closeReason: model.CycleEndProgramChanged}
		}
		return cycleStep{}
	}

	running := current.program == programRunning
	var step cycleStep
	parts := current.partCount - previous.partCount
	switch {
	case parts > 0 && parts <= maxPartCountStep:
		step = cycleStep{closeReason: model.CycleEndPartCount, open: running, parts: parts}
	case running && previous.program == programRunning && !counterCounts && isProgramLoop(*previous, current):
		step = cycleStep{closeReason: model.CycleEndProgramLoop, open: true, openedByLoop: true}
	case running && previous.program == programReset:
		step = cycleStep{open: true}
	case current.program == programReset && previous.program != programReset:
		step = cycleStep{closeReason: resetReason(*previous, current)}
	}
	if !hasOpen {
		step.closeReason = ""
	} else if step.closeReason == "" {
		step.open = false // the open cycle goes on
	}
	return step
}

// isProgramLoop reports whether CurProgLine jumped from previous back to the start of the same
// program: by at least loopMinimumDropFraction of its lines, with the line count unchanged
// (a sub program reports its own lines).
func isProgramLoop(previous, current cycleReading) bool {
	if current.totalLines <= 0 || current.totalLines != previous.totalLines {
		return false
	}
	drop := previous.programLine - current.programLine
	return drop > 0 && float64(drop) >= float64(current.totalLines)*loopMinimumDropFraction
}

// resetReason is why a program that went to reset between previous and current ended.
func resetReason(previous, current cycleReading) string {
	switch {
	case previous.alarm || current.alarm:
		return model.CycleEndAlarmReset
	case previous.program == programRunning:
		return model.CycleEndProgramEnd
	}
	return model.CycleEndOperatorReset
}

// cycleTracker is a worker's cycle state; only the worker goroutine touches it.
type cycleTracker struct {
	loaded        bool                // the open cycle of a previous run was looked for
	open          *model.MachineCycle // the running cycle, nil when none
	lastClosed    *model.MachineCycle // the cycle closed last, for a part counted late
	previous      *cycleReading       // the last usable reading; nil after Offline or a restart
	restored      bool                // open was left by a previous run (no reading since)
	counterCounts bool                // ProgramPartCount counted for the current program
	checkpointAt  time.Time           // open's seconds were last written then

	mode         int                   // the machine mode time is being counted in
	countedUntil time.Time             // mode time is counted into open up to here
	modeTime     map[int]time.Duration // time per machine_modes id inside open
	// offlineSince is when the machine went Offline; it is kept until the next usable reading,
	// so the all-zero snapshot that follows a reconnect (decided Idle) does not lose it.
	offlineSince *time.Time
}

// changeMode counts the time up to at in the old mode and goes on in mode. Offline forgets the
// last reading: what happens while offline is not seen.
func (tracker *cycleTracker) changeMode(at time.Time, mode int) {
	tracker.countUntil(at)
	if mode == model.ModeOffline {
		if tracker.offlineSince == nil {
			since := at
			tracker.offlineSince = &since
		}
		tracker.previous = nil
	}
	tracker.mode = mode
}

// countUntil adds the time from countedUntil to at to the current mode of the open cycle. An
// earlier at (a reading timestamped before the last tick) counts nothing and moves nothing back.
func (tracker *cycleTracker) countUntil(at time.Time) {
	if !at.After(tracker.countedUntil) {
		return
	}
	if tracker.open != nil {
		tracker.modeTime[tracker.mode] += at.Sub(tracker.countedUntil)
	}
	tracker.countedUntil = at
}

// begin makes cycle the open cycle, counting from its start.
func (tracker *cycleTracker) begin(cycle model.MachineCycle) {
	tracker.open = &cycle
	tracker.modeTime = map[int]time.Duration{}
	tracker.countedUntil = cycle.StartedAt
	tracker.checkpointAt = cycle.StartedAt
	tracker.restored = false
}

// progress returns the open cycle with its seconds counted up to at.
func (tracker *cycleTracker) progress(at time.Time) model.MachineCycle {
	tracker.countUntil(at)
	cycle := *tracker.open
	tracker.fillSeconds(&cycle)
	return cycle
}

// close ends the open cycle at now for reason, with the counter at partCount and parts
// counted at the end, and returns it as it is to be saved. A cycle closed for the gap after
// Offline (endAtOffline) ends when the machine went offline and the gap is taken off its
// offline time.
func (tracker *cycleTracker) close(now time.Time, reason string, partCount, parts int, endAtOffline bool) model.MachineCycle {
	tracker.countUntil(now)
	end := now
	if endAtOffline && tracker.offlineSince != nil && tracker.offlineSince.Before(now) {
		end = *tracker.offlineSince
		gap := now.Sub(end)
		tracker.modeTime[model.ModeOffline] = max(0, tracker.modeTime[model.ModeOffline]-gap)
	}
	if end.Before(tracker.open.StartedAt) {
		end = tracker.open.StartedAt
	}
	cycle := *tracker.open
	tracker.fillSeconds(&cycle)
	cycle.Status, cycle.EndReason, cycle.EndedAt = model.CycleStatusFor(reason), reason, &end
	cycle.PartCountEnd = &partCount
	cycle.PartsMade += parts
	tracker.lastClosed = &cycle
	tracker.open = nil
	return cycle
}

// fillSeconds writes the counted time per mode into cycle's seconds. Idle and the user modes do
// not happen inside a cycle (a reset ends it), so they are not kept.
func (tracker *cycleTracker) fillSeconds(cycle *model.MachineCycle) {
	cycle.ProductionSec = secondsOf(tracker.modeTime[model.ModeInProduction])
	cycle.HoldSec = secondsOf(tracker.modeTime[model.ModeHold])
	cycle.DowntimeSec = secondsOf(tracker.modeTime[model.ModeDowntime])
	cycle.OfflineSec = secondsOf(tracker.modeTime[model.ModeOffline])
}

// secondsOf rounds duration to whole seconds, never below 0.
func secondsOf(duration time.Duration) int {
	return max(0, int(math.Round(duration.Seconds())))
}

// lateParts reports whether parts counted at current belong to the cycle closed last: the
// counter stepped within cycleCounterGrace of its end, while no cycle is open or the open one
// began at that end (a program loop).
func (tracker *cycleTracker) lateParts(current cycleReading) bool {
	closed := tracker.lastClosed
	if closed == nil || closed.EndedAt == nil || current.at.Sub(*closed.EndedAt) > cycleCounterGrace {
		return false
	}
	return tracker.open == nil || tracker.open.StartedAt.Equal(*closed.EndedAt)
}

// restore goes on with cycle, left open by a previous run: the time per mode since its start is
// rebuilt from the machine's mode periods, and counting goes on in the period still open (the
// Offline one the collector opened at start) from its start.
func (tracker *cycleTracker) restore(cycle model.MachineCycle, periods []model.MachineModePeriod) {
	tracker.begin(cycle)
	tracker.restored = true
	for _, period := range periods {
		start := period.StartedAt
		if start.Before(cycle.StartedAt) {
			start = cycle.StartedAt
		}
		if period.EndedAt == nil {
			tracker.mode, tracker.countedUntil = period.ModeID, start
			continue
		}
		if period.EndedAt.After(start) {
			tracker.modeTime[period.ModeID] += period.EndedAt.Sub(start)
		}
	}
	tracker.offlineSince = nil
	if tracker.mode == model.ModeOffline {
		since := tracker.countedUntil
		tracker.offlineSince = &since
	}
}

// ---------------------------------------------------------------- worker

// followCycle applies one connected reading to the machine's cycle: it opens, closes and counts
// parts as decideCycleStep says, and writes each change.
func (machineWorker *worker) followCycle(reading cnc.Reading) {
	state := machineStateFromReading(machineWorker.machine.CNCEngine, reading)
	current, ok := cycleReadingFrom(state, reading)
	if !ok || !machineWorker.loadOpenCycle() {
		return
	}
	tracker := &machineWorker.cycles
	openProgram := ""
	if tracker.open != nil {
		openProgram = tracker.open.ProgramName
	}
	step := decideCycleStep(tracker.previous, current, tracker.open != nil, openProgram, tracker.counterCounts)
	machineWorker.applyCycleStep(current, step)
	tracker.previous, tracker.restored, tracker.offlineSince = &current, false, nil
}

// applyCycleStep counts parts, closes and opens cycles as step says.
func (machineWorker *worker) applyCycleStep(current cycleReading, step cycleStep) {
	tracker := &machineWorker.cycles
	if step.parts > 0 {
		tracker.counterCounts = true
		if tracker.lateParts(current) {
			machineWorker.countLateParts(current, step.parts)
			step.closeReason = ""
			step.open = step.open && tracker.open == nil
		}
	}
	if step.closeReason != "" {
		machineWorker.closeCycle(current, step)
	}
	if step.open {
		machineWorker.openCycle(current, step.openedByLoop)
	}
}

// loadOpenCycle looks once per worker for a cycle a previous run left open and goes on with it.
// It reports whether the cycle state is known (false while the database cannot be read).
func (machineWorker *worker) loadOpenCycle() bool {
	tracker := &machineWorker.cycles
	if tracker.loaded {
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	cycle, err := machineWorker.collector.cycleService.GetOpenCycle(ctx, machineWorker.machine.ID)
	if err != nil {
		machineWorker.log(logTypeCycle, slog.LevelError, "cnc: load open cycle", "error", err)
		return false
	}
	tracker.loaded = true
	if cycle == nil {
		return true
	}
	periods, err := machineWorker.collector.service.GetModePeriods(ctx, machineWorker.machine.ID, cycle.StartedAt)
	if err != nil {
		machineWorker.log(logTypeCycle, slog.LevelError, "cnc: rebuild open cycle", "cycleId", cycle.ID, "error", err)
	}
	tracker.restore(*cycle, periods)
	restored := *tracker.open
	tracker.fillSeconds(&restored)
	machineWorker.log(logTypeCycle, slog.LevelInfo, "cnc: cycle restored", "cycleId", cycle.ID, "program", cycle.ProgramName,
		"startedAt", cycle.StartedAt, "mode", modeName(tracker.mode), "productionSec", restored.ProductionSec,
		"holdSec", restored.HoldSec, "downtimeSec", restored.DowntimeSec, "offlineSec", restored.OfflineSec)
	return true
}

// countLateParts gives parts to the cycle closed last (lateParts): a program loop it ended at
// becomes a counted part.
func (machineWorker *worker) countLateParts(current cycleReading, parts int) {
	closed := machineWorker.cycles.lastClosed
	closed.PartsMade += parts
	closed.PartCountEnd = &current.partCount
	if closed.EndReason == model.CycleEndProgramLoop {
		closed.EndReason = model.CycleEndPartCount
	}
	machineWorker.saveCycle(*closed, "cnc: count part on closed cycle")
	machineWorker.log(logTypeCycle, slog.LevelInfo, "cnc: cycle part counted after its end", "cycleId", closed.ID,
		"parts", parts, "partCount", current.partCount, "endReason", closed.EndReason)
}

// openCycle opens a cycle at current: the service matches its job and shift and stores it. The
// counter's habit is learned again for a program other than the last cycle's.
func (machineWorker *worker) openCycle(current cycleReading, openedByLoop bool) {
	tracker := &machineWorker.cycles
	if tracker.lastClosed == nil || tracker.lastClosed.ProgramName != current.programName {
		tracker.counterCounts = false
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	cycle, err := machineWorker.collector.cycleService.Create(ctx, model.MachineCycle{
		MachineID: machineWorker.machine.ID, ProgramName: current.programName,
		StartedAt: current.at, PartCountStart: &current.partCount,
	})
	if err != nil {
		machineWorker.log(logTypeCycle, slog.LevelError, "cnc: open cycle", "program", current.programName, "error", err)
		return
	}
	tracker.begin(cycle)
	machineWorker.log(logTypeCycle, slog.LevelInfo, "cnc: cycle opened", "cycleId", cycle.ID, "program", cycle.ProgramName,
		"partCount", current.partCount, "byLoop", openedByLoop, "scheduleItemId", cycle.ScheduleItemID,
		"partId", cycle.PartID, "processId", cycle.ProcessID, "shiftId", cycle.ShiftID, "shiftDate", cycle.ShiftDate,
		"standardCycleSec", cycle.StandardCycleSec, "loadingSec", cycle.LoadingSec)
}

// closeCycle ends the open cycle for step.closeReason. One closed for what happened during
// Offline ends when the machine went offline; one a previous run left ends as collector_stop.
func (machineWorker *worker) closeCycle(current cycleReading, step cycleStep) {
	reason := step.closeReason
	afterGap := reason == model.CycleEndOfflineReset || reason == model.CycleEndProgramChanged
	if reason == model.CycleEndOfflineReset && machineWorker.cycles.restored {
		reason = model.CycleEndCollectorStop
	}
	machineWorker.finishCycle(current.at, reason, current.partCount, step.parts, afterGap)
}

// finishCycle closes the open cycle at now for reason (see cycleTracker.close), saves and logs it.
func (machineWorker *worker) finishCycle(now time.Time, reason string, partCount, parts int, endAtOffline bool) {
	cycle := machineWorker.cycles.close(now, reason, partCount, parts, endAtOffline)
	machineWorker.saveCycle(cycle, "cnc: close cycle")
	machineWorker.log(logTypeCycle, slog.LevelInfo, "cnc: cycle "+cycle.Status, "cycleId", cycle.ID, "endReason", reason,
		"program", cycle.ProgramName, "startedAt", cycle.StartedAt, "endedAt", *cycle.EndedAt,
		"durationSec", secondsOf(cycle.EndedAt.Sub(cycle.StartedAt)), "productionSec", cycle.ProductionSec,
		"holdSec", cycle.HoldSec, "downtimeSec", cycle.DowntimeSec, "offlineSec", cycle.OfflineSec,
		"partsMade", cycle.PartsMade)
}

// checkCycle runs on the worker's tick: a cycle open past CNC_CYCLE_MAX_HOURS is closed as
// timeout, a running one has its seconds written every cycleCheckpointInterval.
func (machineWorker *worker) checkCycle(now time.Time) {
	tracker := &machineWorker.cycles
	if tracker.open == nil {
		return
	}
	if now.Sub(tracker.open.StartedAt) >= machineWorker.collector.config.CycleMaxDuration {
		partCount := 0
		if tracker.previous != nil {
			partCount = tracker.previous.partCount
		} else if tracker.open.PartCountEnd != nil {
			partCount = *tracker.open.PartCountEnd
		}
		machineWorker.finishCycle(now, model.CycleEndTimeout, partCount, 0, false)
		return
	}
	if now.Sub(tracker.checkpointAt) >= cycleCheckpointInterval {
		machineWorker.saveCycleProgress(now)
	}
}

// saveCycleProgress writes the open cycle's seconds counted up to now.
func (machineWorker *worker) saveCycleProgress(now time.Time) {
	tracker := &machineWorker.cycles
	if tracker.open == nil {
		return
	}
	tracker.checkpointAt = now
	machineWorker.saveCycle(tracker.progress(now), "cnc: save cycle progress")
}

// saveCycle writes cycle; an open cycle deleted with its machine is dropped from the tracker.
func (machineWorker *worker) saveCycle(cycle model.MachineCycle, operation string) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	err := machineWorker.collector.cycleService.Update(ctx, cycle)
	if err == nil {
		return
	}
	if tracker := &machineWorker.cycles; errors.Is(err, response.MachineCycleNotFound) && tracker.open != nil && tracker.open.ID == cycle.ID {
		tracker.open = nil
	}
	machineWorker.log(logTypeCycle, slog.LevelError, operation, "cycleId", cycle.ID, "error", err)
}
