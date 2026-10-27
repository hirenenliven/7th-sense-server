package monitor

import (
	"fmt"
	"time"

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
//
// Every mode change is logged as "cnc: mode changed" with these inputs and the rule that
// matched (recordModeChange), so a mode can be checked against the HMI.

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

// reasonNoConnection is the reason logged for Offline.
const reasonNoConnection = "no NC connection"

// programState is the NC program's state, whatever numbering the controller uses.
type programState int

const (
	programReset programState = iota
	programStopped
	programRunning
	programWaiting // tool change, M function, dwell: counts as running
	programInterrupted
)

var programStateNames = map[programState]string{
	programReset: "reset", programStopped: "stopped", programRunning: "running",
	programWaiting: "waiting", programInterrupted: "interrupted",
}

// channelState is the NC channel's state (/Channel/State/chanStatus: 0 reset, 1 interrupted,
// 2 active).
type channelState int

const (
	channelUnknown channelState = iota
	channelReset
	channelInterrupted
	channelActive
)

var channelStateNames = map[channelState]string{
	channelUnknown: "unknown", channelReset: "reset", channelInterrupted: "interrupted", channelActive: "active",
}

var opModeNames = map[int]string{0: "JOG", 1: "MDA", opModeAuto: "AUTO"}

// machineState is the NC state the mode is decided from. procStatus, alarmNumber and
// emergency keep the raw values for the log.
type machineState struct {
	connected    bool
	opMode       int
	procStatus   int
	program      programState
	channel      channelState
	feedOverride float64 // %
	alarmNumber  int
	emergency    bool
	programName  string
}

// isRunning reports whether a program is running; waiting counts as running.
func (state machineState) isRunning() bool {
	return state.program == programRunning || state.program == programWaiting
}

// isAlarmActive reports whether an NC alarm or the emergency stop is active.
func (state machineState) isAlarmActive() bool {
	return state.alarmNumber != 0 || state.emergency
}

// logAttributes are the NC inputs of the decision, for "cnc: mode changed".
func (state machineState) logAttributes() []any {
	opMode := opModeNames[state.opMode]
	if opMode == "" {
		opMode = fmt.Sprintf("unknown(%d)", state.opMode)
	}
	return []any{
		"opMode", opMode,
		"procStatus", fmt.Sprintf("%d (%s)", state.procStatus, programStateNames[state.program]),
		"chanStatus", channelStateNames[state.channel],
		"feedOverride", state.feedOverride,
		"alarmNo", state.alarmNumber,
		"emergency", state.emergency,
		"ncProgram", state.programName,
	}
}

// machineStateFromReading reads the NC state from a DLL reading.
func machineStateFromReading(reading cnc.Reading) machineState {
	procStatus := int(reading.Values[cnc.FieldProcStatus].Number)
	return machineState{
		connected:    reading.Connected,
		opMode:       int(reading.Values[cnc.FieldMachineMode].Number),
		procStatus:   procStatus,
		program:      programStateFromProcStatus(procStatus),
		channel:      channelUnknown,
		feedOverride: reading.Values[cnc.FieldFeedOverride].Number,
		alarmNumber:  int(reading.Values[cnc.FieldAlarmNo].Number),
		emergency:    reading.Values[cnc.FieldIsEmergency].Number != 0,
		programName:  reading.Values[cnc.FieldProgName].Text,
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

// decideMachineMode returns the machine_modes id of state and the rule that decided it.
// First match wins:
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
func decideMachineMode(state machineState) (int, string) {
	running := state.isRunning()
	switch {
	case !state.connected:
		return model.ModeOffline, reasonNoConnection
	case state.isAlarmActive() && !running:
		return model.ModeDowntime, "alarm active, no program running"
	case isInProduction(state):
		return model.ModeInProduction, "AUTO, program running, feed override > 0, no alarm"
	}
	if reason := holdReason(state); reason != "" {
		return model.ModeHold, reason
	}
	return model.ModeIdle, "no program running"
}

// isInProduction reports whether the machine is cutting in AUTO without anything holding it.
func isInProduction(state machineState) bool {
	channelActive := state.channel == channelActive || state.channel == channelUnknown
	return state.opMode == opModeAuto && state.isRunning() && channelActive &&
		state.feedOverride > 0 && !state.isAlarmActive()
}

// holdReason returns why a program is held (stopped or interrupted, or running but not moving:
// feed override 0, an alarm, or not in AUTO), or "" when it is not held.
func holdReason(state machineState) string {
	switch {
	case state.program == programInterrupted:
		return "program interrupted"
	case state.program == programStopped:
		return "program stopped"
	case state.channel == channelInterrupted:
		return "channel interrupted"
	case !state.isRunning():
		return ""
	case state.feedOverride <= 0:
		return "program running, feed override ≤ 0"
	case state.isAlarmActive():
		return "program running, alarm active"
	case state.opMode != opModeAuto:
		return "program running, not in AUTO"
	}
	return ""
}

// recordModeFromReading decides the mode of a connected reading and records it when it changed.
func (machineWorker *worker) recordModeFromReading(at time.Time, reading cnc.Reading) {
	state := machineStateFromReading(reading)
	mode, reason := decideMachineMode(state)
	machineWorker.recordModeChange(at, mode, reason, state.logAttributes())
}
