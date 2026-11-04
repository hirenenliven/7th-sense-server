package monitor

import (
	"testing"
	"time"

	"7th-sense-server/internal/model"
)

// readingAt is a cycle reading second seconds after a fixed start.
func readingAt(second float64, program programState, name string, line, partCount int) cycleReading {
	start := time.Date(2026, 11, 2, 17, 40, 46, 0, time.Local)
	return cycleReading{
		at: start.Add(time.Duration(second * float64(time.Second))), program: program, programName: name,
		programLine: line, totalLines: 40, partCount: partCount,
	}
}

// The floor test on machine 195 (cycles 232-234): O1002 starts, is held, runs on, calls sub
// program O9001 (line 6 -> 2, ProgName O9001), the counter counts and M30 resets the program a
// reading later. That is one cycle, completed by the part count.
func TestDecideCycleStepSubProgramCallIsNotALoop(t *testing.T) {
	beforeCall := readingAt(32, programRunning, "O1002", 6, 269)
	inSubProgram := readingAt(33, programRunning, "O9001", 2, 269)
	if step := decideCycleStep(&beforeCall, inSubProgram, true, "O1002", false); step != (cycleStep{}) {
		t.Fatalf("sub program call: got %+v, want the cycle to go on", step)
	}

	counted := readingAt(43, programRunning, "O9001", 7, 270)
	step := decideCycleStep(&inSubProgram, counted, true, "O1002", false)
	if step.closeReason != model.CycleEndPartCount || !step.open || step.parts != 1 {
		t.Fatalf("part count: got %+v", step)
	}
}

func TestIsProgramLoop(t *testing.T) {
	tests := []struct {
		name        string
		previous    cycleReading
		current     cycleReading
		openProgram string
		want        bool
	}{
		{"main program starts over (M99)", readingAt(0, programRunning, "JYOTI", 30, 0), readingAt(1, programRunning, "JYOTI", 2, 0), "JYOTI", true},
		{"call into a sub program", readingAt(0, programRunning, "testavim", 22, 0), readingAt(1, programRunning, "O9001", 2, 0), "testavim", false},
		{"sub program repeats inside the cycle", readingAt(0, programRunning, "O9001", 30, 0), readingAt(1, programRunning, "O9001", 2, 0), "testavim", false},
		{"program name unknown", readingAt(0, programRunning, "", 30, 0), readingAt(1, programRunning, "", 2, 0), "", true},
		{"small jump back", readingAt(0, programRunning, "JYOTI", 10, 0), readingAt(1, programRunning, "JYOTI", 8, 0), "JYOTI", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isProgramLoop(test.previous, test.current, test.openProgram); got != test.want {
				t.Errorf("isProgramLoop = %v, want %v", got, test.want)
			}
		})
	}
}

// Machine 195 (cycles 520, 523): RESET while running shows IsEmergency 1 for 0.2 - 0.6 s at the
// reset reading. That is not an E-stop until it stays on for emergencyResetGrace.
func TestIsResetKeyEmergency(t *testing.T) {
	running := readingAt(0, programRunning, "JYOTI", 4, 274)
	resetWithEmergency := readingAt(0.1, programReset, "JYOTI", 4, 274)
	resetWithEmergency.emergency = true

	tracker := cycleTracker{previous: &running}
	if !tracker.isResetKeyEmergency(resetWithEmergency) {
		t.Fatal("a fresh IsEmergency at the reset must be taken as the RESET key")
	}

	withAlarm := resetWithEmergency
	withAlarm.alarm = true
	if tracker.isResetKeyEmergency(withAlarm) {
		t.Fatal("an AlarmNo is a real alarm")
	}

	pressedBefore := running.at.Add(-emergencyResetGrace)
	tracker.emergencySince = &pressedBefore
	if tracker.isResetKeyEmergency(resetWithEmergency) {
		t.Fatal("an IsEmergency on for emergencyResetGrace before the reset is an E-stop")
	}
}

func TestFollowEmergency(t *testing.T) {
	on := readingAt(0, programReset, "JYOTI", 0, 274)
	on.emergency = true
	tracker := cycleTracker{checkEmergency: true}
	tracker.followEmergency(on)
	if tracker.emergencySince == nil || !tracker.emergencySince.Equal(on.at) || !tracker.checkEmergency {
		t.Fatalf("IsEmergency on: since %v, check %v", tracker.emergencySince, tracker.checkEmergency)
	}

	off := readingAt(0.4, programReset, "JYOTI", 0, 274)
	tracker.followEmergency(off)
	if tracker.emergencySince != nil || tracker.checkEmergency {
		t.Fatal("IsEmergency off within the grace: the cycle keeps its end reason")
	}
}

// The next cycle's time is counted from its start while it waits, so it keeps it when it opens.
func TestAwaitNextCountsFromItsStart(t *testing.T) {
	start := readingAt(0, programRunning, "O1002", 2, 270)
	tracker := cycleTracker{mode: model.ModeInProduction}
	tracker.awaitNext(start, false)
	tracker.countUntil(start.at.Add(cycleStartGrace))
	if got := tracker.modeTime[model.ModeInProduction]; got != cycleStartGrace {
		t.Fatalf("counted %v, want %v", got, cycleStartGrace)
	}

	tracker.changeMode(start.at.Add(cycleStartGrace), model.ModeOffline)
	if tracker.next != nil {
		t.Fatal("Offline must drop the waiting cycle: its start is no longer known to go on")
	}
}

// The controller counts at most one part per M30: a bigger step is the counter set by hand.
func TestDecideCycleStepCountsOnePartPerStep(t *testing.T) {
	running := readingAt(0, programRunning, "O1002", 10, 287)
	if step := decideCycleStep(&running, readingAt(1, programRunning, "O1002", 11, 288), true, "O1002", true); step.closeReason != model.CycleEndPartCount || step.parts != 1 {
		t.Fatalf("one part: got %+v", step)
	}
	if step := decideCycleStep(&running, readingAt(1, programRunning, "O1002", 11, 289), true, "O1002", true); step.closeReason != "" || step.parts != 0 {
		t.Fatalf("counter set by hand (+2): got %+v, want the cycle to go on", step)
	}
}

// Machine 195, 2026-11-03: cycle 568 ended (M30 reset) at 17:09:40 with the counter at 287; the
// DLL sent 288 only with the next start at 17:10:01. That part is cycle 568's.
func TestLatePartsBelongToTheCycleClosedLast(t *testing.T) {
	ended := readingAt(0, programReset, "O1002", 0, 287).at
	closed := model.MachineCycle{EndedAt: &ended}

	nextStart := readingAt(21, programRunning, "O1002", 1, 288)
	tracker := cycleTracker{lastClosed: &closed}
	if !tracker.lateParts(nextStart) {
		t.Fatal("a step with no cycle open is the closed cycle's part")
	}

	tracker.open = &model.MachineCycle{StartedAt: nextStart.at}
	if !tracker.lateParts(readingAt(23, programRunning, "O1002", 2, 288)) {
		t.Fatal("a step within cycleCounterGrace of the next start is the closed cycle's part")
	}
	if tracker.lateParts(readingAt(40, programRunning, "O1002", 30, 288)) {
		t.Fatal("a step later in the next cycle is that cycle's own part")
	}

	closed.PartsMade = 1
	tracker.open = nil
	if tracker.lateParts(nextStart) {
		t.Fatal("a cycle has at most one part")
	}
}
