package monitor

import (
	"slices"
	"testing"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/model"
)

// ncReading is a connected DLL reading with the raw values the mode rules read.
func ncReading(opMode, procStatus int, feedOverride float64, alarmNumber int, emergency bool) cnc.Reading {
	var reading cnc.Reading
	reading.Connected = true
	reading.Values[cnc.FieldMachineMode].Number = float64(opMode)
	reading.Values[cnc.FieldProcStatus].Number = float64(procStatus)
	reading.Values[cnc.FieldFeedOverride].Number = feedOverride
	reading.Values[cnc.FieldAlarmNo].Number = float64(alarmNumber)
	if emergency {
		reading.Values[cnc.FieldIsEmergency].Number = 1
	}
	return reading
}

// The raw values each controller sends, decided the way the vendor sample (G0_V0.0.2_DLL
// main.go) does on the live machines.
func TestDecideMachineMode(t *testing.T) {
	fanuc, siemens, syntec := model.EngineFanuc, model.EngineSiemens, model.EngineSyntec
	tests := []struct {
		name       string
		engine     *int
		reading    cnc.Reading
		wantMode   int
		wantReason string
	}{
		{"fanuc AUTO START", &fanuc, ncReading(1, 3, 120, 0, false), model.ModeInProduction, "AUTO, program running, feed override > 0, no alarm"},
		{"fanuc AUTO HOLD", &fanuc, ncReading(1, 2, 120, 0, false), model.ModeHold, "program held"},
		{"fanuc AUTO Stop", &fanuc, ncReading(1, 1, 120, 0, false), model.ModeHold, "program held"},
		{"fanuc MSTR", &fanuc, ncReading(1, 4, 120, 0, false), model.ModeHold, "program held"},
		{"fanuc MDI START", &fanuc, ncReading(0, 3, 120, 0, false), model.ModeHold, "program running, not in AUTO"},
		{"fanuc START feed override 0", &fanuc, ncReading(1, 3, 0, 0, false), model.ModeHold, "program running, feed override ≤ 0"},
		{"fanuc START with alarm", &fanuc, ncReading(1, 3, 120, 100, false), model.ModeHold, "program running, alarm active"},
		{"fanuc AUTO Reset", &fanuc, ncReading(1, 0, 120, 0, false), model.ModeIdle, "program reset, no alarm (any mode)"},
		{"fanuc JOG Reset", &fanuc, ncReading(5, 0, 120, 0, false), model.ModeIdle, "program reset, no alarm (any mode)"},
		{"fanuc alarm, Reset", &fanuc, ncReading(1, 0, 120, 100, false), model.ModeDowntime, "alarm active, no program running"},
		{"fanuc E-stop, HOLD", &fanuc, ncReading(1, 2, 120, 0, true), model.ModeDowntime, "alarm active, no program running"},
		{"siemens AUTO Active", &siemens, ncReading(2, 2, 100, 0, false), model.ModeInProduction, "AUTO, program running, feed override > 0, no alarm"},
		{"siemens AUTO Waiting", &siemens, ncReading(2, 3, 100, 0, false), model.ModeInProduction, "AUTO, program running, feed override > 0, no alarm"},
		{"siemens Stop", &siemens, ncReading(2, 1, 100, 0, false), model.ModeHold, "program held"},
		{"siemens Interrupted", &siemens, ncReading(2, 4, 100, 0, false), model.ModeHold, "program held"},
		{"siemens MDA Active", &siemens, ncReading(1, 2, 100, 0, false), model.ModeHold, "program running, not in AUTO"},
		{"siemens JOG Reset", &siemens, ncReading(0, 0, 100, 0, false), model.ModeIdle, "program reset, no alarm (any mode)"},
		{"syntec read as Siemens", &syntec, ncReading(2, 2, 100, 0, false), model.ModeInProduction, "AUTO, program running, feed override > 0, no alarm"},
		{"no engine read as Siemens", nil, ncReading(2, 4, 100, 0, false), model.ModeHold, "program held"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := decideMachineMode(machineStateFromReading(test.engine, test.reading))
			if decision.mode != test.wantMode || decision.reason != test.wantReason {
				t.Errorf("decideMachineMode = %s (%q), want %s (%q)", modeName(decision.mode), decision.reason,
					modeName(test.wantMode), test.wantReason)
			}
		})
	}
	if decision := decideMachineMode(machineStateFromReading(&fanuc, cnc.Reading{})); decision.mode != model.ModeOffline {
		t.Errorf("not connected = %s, want Offline", modeName(decision.mode))
	}
}

// Every combination of the rule inputs matches a rule (the "no rule matched" fallback is never
// reached), and Idle is exactly "program reset and no alarm", whatever the operation mode.
func TestDecideMachineModeEveryCombination(t *testing.T) {
	for _, engine := range []int{model.EngineFanuc, model.EngineSiemens} {
		for opMode := 0; opMode <= 9; opMode++ {
			for procStatus := 0; procStatus <= 6; procStatus++ {
				for _, feedOverride := range []float64{0, 100} {
					for _, alarmNumber := range []int{0, 100} {
						for _, emergency := range []bool{false, true} {
							state := machineStateFromReading(&engine, ncReading(opMode, procStatus, feedOverride, alarmNumber, emergency))
							decision := decideMachineMode(state)
							if decision.reason == "no rule matched" || len(decision.conditions) == 0 {
								t.Fatalf("%+v: %s %q %v", state, modeName(decision.mode), decision.reason, decision.conditions)
							}
							wantIdle := state.program == programReset && alarmNumber == 0 && !emergency
							if (decision.mode == model.ModeIdle) != wantIdle {
								t.Errorf("%+v: %s, Idle expected %t", state, modeName(decision.mode), wantIdle)
							}
						}
					}
				}
			}
		}
	}
}

func TestMachineStateFromReading(t *testing.T) {
	reading := ncReading(1, 2, 80, 0, true)
	reading.Values[cnc.FieldProgName].Text = "JYOTI"

	fanuc := model.EngineFanuc
	state := machineStateFromReading(&fanuc, reading)
	want := machineState{engine: model.EngineFanuc, connected: true, opMode: 1, procStatus: 2, program: programHeld,
		feedOverride: 80, emergency: true, programName: "JYOTI"}
	if state != want {
		t.Errorf("machineStateFromReading = %+v, want %+v", state, want)
	}
}

func TestMachineStateFromReadingWithoutEngine(t *testing.T) {
	if state := machineStateFromReading(nil, cnc.Reading{}); state.engine != engineNone {
		t.Errorf("engine = %d, want engineNone", state.engine)
	}
}

// Only the inputs a rule checks are logged: a feed override of 0 holds a running program, so
// the Hold line shows the program and the feed override, not the alarm or the opMode.
func TestDecideMachineModeConditions(t *testing.T) {
	siemens := model.EngineSiemens
	decision := decideMachineMode(machineStateFromReading(&siemens, ncReading(2, 2, 0, 0, false)))
	want := []any{"procStatus", "2 (Active)", "programState", "running", "feedOverride", 0.0}
	if decision.mode != model.ModeHold || !slices.Equal(decision.conditions, want) {
		t.Errorf("decideMachineMode = %s %v, want Hold %v", modeName(decision.mode), decision.conditions, want)
	}
	if offline := decideMachineMode(machineState{}); !slices.Equal(offline.conditions, []any{"connected", false}) {
		t.Errorf("offline conditions = %v", offline.conditions)
	}
}
