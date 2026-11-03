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
