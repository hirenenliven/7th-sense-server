package monitor

import (
	"testing"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/model"
)

func TestDecideMachineMode(t *testing.T) {
	cutting := machineState{connected: true, opMode: opModeAuto, program: programRunning, feedOverride: 100}
	with := func(change func(state *machineState)) machineState {
		state := cutting
		change(&state)
		return state
	}
	tests := []struct {
		name  string
		state machineState
		want  int
	}{
		{"offline", with(func(state *machineState) { state.connected = false }), model.ModeOffline},
		{"alarm, reset", with(func(state *machineState) { state.alarmNumber, state.program = 3000, programReset }), model.ModeDowntime},
		{"alarm, stopped", with(func(state *machineState) { state.emergency, state.program = true, programStopped }), model.ModeDowntime},
		{"auto running", cutting, model.ModeInProduction},
		{"auto waiting", with(func(state *machineState) { state.program = programWaiting }), model.ModeInProduction},
		{"auto running, channel active", with(func(state *machineState) { state.channel = channelActive }), model.ModeInProduction},
		{"auto running, channel reset", with(func(state *machineState) { state.channel = channelReset }), model.ModeIdle},
		{"feed override 0", with(func(state *machineState) { state.feedOverride = 0 }), model.ModeHold},
		{"alarm while running", with(func(state *machineState) { state.alarmNumber = 700115 }), model.ModeHold},
		{"running in MDA", with(func(state *machineState) { state.opMode = 1 }), model.ModeHold},
		{"stopped", with(func(state *machineState) { state.program = programStopped }), model.ModeHold},
		{"interrupted", with(func(state *machineState) { state.program = programInterrupted }), model.ModeHold},
		{"channel interrupted", with(func(state *machineState) { state.channel = channelInterrupted }), model.ModeHold},
		{"auto reset", with(func(state *machineState) { state.program = programReset }), model.ModeIdle},
		{"jog", with(func(state *machineState) { state.opMode, state.program = 0, programReset }), model.ModeIdle},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, reason := decideMachineMode(test.state); got != test.want || reason == "" {
				t.Errorf("decideMachineMode = %s (%q), want %s", modeName(got), reason, modeName(test.want))
			}
		})
	}
}

func TestMachineStateFromReading(t *testing.T) {
	var reading cnc.Reading
	reading.Connected = true
	reading.Values[cnc.FieldMachineMode].Number = opModeAuto
	reading.Values[cnc.FieldProcStatus].Number = procStatusWaiting
	reading.Values[cnc.FieldFeedOverride].Number = 80
	reading.Values[cnc.FieldIsEmergency].Number = 1
	reading.Values[cnc.FieldProgName].Text = "PART.MPF"

	siemens := model.EngineSiemens
	state := machineStateFromReading(&siemens, reading)
	want := machineState{engine: model.EngineSiemens, connected: true, opMode: opModeAuto, procStatus: procStatusWaiting, program: programWaiting,
		channel: channelUnknown, feedOverride: 80, emergency: true, programName: "PART.MPF"}
	if state != want {
		t.Errorf("machineStateFromReading = %+v, want %+v", state, want)
	}
}

func TestMachineStateFromReadingWithoutEngine(t *testing.T) {
	if state := machineStateFromReading(nil, cnc.Reading{}); state.engine != engineNone {
		t.Errorf("engine = %d, want engineNone", state.engine)
	}
}
