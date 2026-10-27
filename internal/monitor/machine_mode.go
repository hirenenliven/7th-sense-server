package monitor

import (
	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/model"
)

// The machine mode is decided from the NC state of every connected reading, the way the
// Sinumerik HMI shows it, not from the stored variables: a machine needs no mode variables
// configured, and the mode is right from the first reading after a (re)connect even when no
// value changed since the last run.
//
// The NC values come from the DLL snapshot, which has no read function for single nodes:
//
//	opMode       MachineMode         /Bag/State/opMode (0 JOG, 1 MDA, 2 AUTO)
//	program      ProcStatus          /Channel/State/acProg (see programStateFromProcStatus)
//	chanStatus   (not in the snapshot; channelUnknown until the DLL sends it)
//	feedRateOvr  FeedRatePercetnage  feed override, %
//	alarmActive  AlarmNo ≠ 0 or IsEmergency (E-stop raises alarm 3000 as well)

// opModeAuto is opMode AUTO; JOG (0) and MDA (1) are manual modes.
const opModeAuto = 2

// ProcStatus values the DLL sends: /Channel/State/acProg, which numbers the program states
// differently from /Channel/State/progStatus (1 interrupted, 2 stopped, 3 running, 4 waiting).
const (
	procStatusReset       = 0
	procStatusStopped     = 1
	procStatusActive      = 2
	procStatusWaiting     = 3
	procStatusInterrupted = 4
)

// programState is the NC program's state, whatever numbering the controller uses.
type programState int

const (
	programReset programState = iota
	programStopped
	programRunning
	programWaiting // tool change, M function, dwell: counts as running
	programInterrupted
)

// channelState is the NC channel's state (/Channel/State/chanStatus: 0 reset, 1 interrupted,
// 2 active).
type channelState int

const (
	channelUnknown channelState = iota
	channelReset
	channelInterrupted
	channelActive
)

// machineState is the NC state the mode is decided from.
type machineState struct {
	connected    bool
	opMode       int
	program      programState
	channel      channelState
	feedOverride float64 // %
	alarmActive  bool
}

// isRunning reports whether a program is running; waiting counts as running.
func (state machineState) isRunning() bool {
	return state.program == programRunning || state.program == programWaiting
}

// machineStateFromReading reads the NC state from a DLL reading.
func machineStateFromReading(reading cnc.Reading) machineState {
	return machineState{
		connected:    reading.Connected,
		opMode:       int(reading.Values[cnc.FieldMachineMode].Number),
		program:      programStateFromProcStatus(int(reading.Values[cnc.FieldProcStatus].Number)),
		channel:      channelUnknown,
		feedOverride: reading.Values[cnc.FieldFeedOverride].Number,
		alarmActive:  reading.Values[cnc.FieldAlarmNo].Number != 0 || reading.Values[cnc.FieldIsEmergency].Number != 0,
	}
}

// programStateFromProcStatus turns the DLL's ProcStatus (acProg) into a programState; an
// unknown value counts as reset.
func programStateFromProcStatus(procStatus int) programState {
	switch procStatus {
	case procStatusStopped:
		return programStopped
	case procStatusActive:
		return programRunning
	case procStatusWaiting:
		return programWaiting
	case procStatusInterrupted:
		return programInterrupted
	}
	return programReset
}

// decideMachineMode returns the machine_modes id of state. First match wins:
//
//	Offline (9)        no NC connection
//	Downtime (6)       alarm active and no program running (an alarm stopped machining)
//	In Production (3)  AUTO, program running or waiting, channel active, feed override > 0, no alarm
//	Hold (4)           program interrupted / stopped, channel interrupted, or running with
//	                   feed override ≤ 0, an alarm or a manual mode
//	Idle (5)           anything else: reset, program not running, manual JOG
//
// An unknown channel state does not decide anything: In Production does not wait for it and
// Hold is not caused by it.
func decideMachineMode(state machineState) int {
	running := state.isRunning()
	switch {
	case !state.connected:
		return model.ModeOffline
	case state.alarmActive && !running:
		return model.ModeDowntime
	case isInProduction(state):
		return model.ModeInProduction
	case isOnHold(state):
		return model.ModeHold
	}
	return model.ModeIdle
}

// isInProduction reports whether the machine is cutting in AUTO without anything holding it.
func isInProduction(state machineState) bool {
	channelActive := state.channel == channelActive || state.channel == channelUnknown
	return state.opMode == opModeAuto && state.isRunning() && channelActive &&
		state.feedOverride > 0 && !state.alarmActive
}

// isOnHold reports whether a program is held: stopped or interrupted, or running but not
// moving (feed override 0, an alarm, or not in AUTO).
func isOnHold(state machineState) bool {
	if state.program == programInterrupted || state.program == programStopped || state.channel == channelInterrupted {
		return true
	}
	return state.isRunning() && (state.feedOverride <= 0 || state.alarmActive || state.opMode != opModeAuto)
}
