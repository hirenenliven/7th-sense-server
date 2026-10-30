package monitor

import (
	"fmt"
	"slices"
	"time"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/model"
)

// The machine mode is decided from the NC state of every connected reading, not from the
// stored variables: a machine needs no mode variables configured, and the mode is right from
// the first reading after a (re)connect even when no value changed since the last run. The
// rules are the ones tested on live Fanuc and Siemens machines with the vendor sample
// (G0_V0.0.2_DLL main.go, decideMachineMode).
//
// The NC values come from the DLL snapshot:
//
//	opMode       MachineMode         AUTO is Fanuc 1 (MEM), Siemens 2 (autoOpModes)
//	program      ProcStatus          running / held / reset per engine (programStatesByEngine)
//	feedRateOvr  FeedRatePercetnage  feed override, %
//	alarmActive  AlarmNo ≠ 0 or IsEmergency
//
// MachineMode and ProcStatus are numbered differently per engine (nc_status.go names them):
// ProcStatus 2 is HOLD on Fanuc but Active (running) on Siemens. An engine without its own
// table (simulator, Syntec, Mitsubishi, none set) is read with the Siemens numbering.
//
// With CNC_LOG_MODES=true every mode change is logged with the final mode first, the rule that
// matched and only the inputs that rule checks (modeDecision.conditions); CNC_LOG_MODE_CHECKS
// logs the same whenever those inputs change (cnc_log.go), so a mode can be checked against the
// HMI.

// reasonNoConnection is the reason logged for Offline.
const reasonNoConnection = "no NC connection"

// programState is the NC program's state, whatever numbering the controller uses.
type programState int

const (
	programReset programState = iota // also a ProcStatus value the engine's table does not have
	programRunning
	programHeld
)

var programStateNames = map[programState]string{
	programReset: "reset", programRunning: "running", programHeld: "held",
}

// autoOpModes is the MachineMode value of AUTO per engine; every other value is a manual mode
// (Fanuc MDI, EDIT, JOG, REF; Siemens JOG, MDA).
var autoOpModes = map[int]int{
	model.EngineFanuc:   1,
	model.EngineSiemens: 2,
}

// programStatesByEngine turns ProcStatus into a programState per engine:
//
//	Fanuc    0 Reset; 3 START -> running; 1 Stop, 2 HOLD (feed hold), 4 MSTR (tool retract /
//	         manual start, not cutting the program) -> held
//	Siemens  0 Reset; 2 Active -> running (as Fanuc START); 1 Stop, 3 Waiting (as Fanuc HOLD),
//	         4 Interrupted -> held
var programStatesByEngine = map[int]map[int]programState{
	model.EngineFanuc:   {0: programReset, 1: programHeld, 2: programHeld, 3: programRunning, 4: programHeld},
	model.EngineSiemens: {0: programReset, 1: programHeld, 2: programRunning, 3: programHeld, 4: programHeld},
}

// numberingEngine is the engine whose numbering a machine's MachineMode / ProcStatus are read
// with: its own when the rules know it, otherwise Siemens.
func numberingEngine(engine int) int {
	if _, known := autoOpModes[engine]; known {
		return engine
	}
	return model.EngineSiemens
}

// machineState is the NC state the mode is decided from. engine, procStatus, alarmNumber,
// alarmMessage and emergency keep the raw values for the log.
type machineState struct {
	engine       int // controllers.cnc_engine, engineNone when not set
	connected    bool
	opMode       int
	procStatus   int
	program      programState
	feedOverride float64 // %
	alarmNumber  int
	alarmMessage string
	emergency    bool
	programName  string
}

// isRunning reports whether a program is running (Fanuc START, Siemens Active).
func (state machineState) isRunning() bool {
	return state.program == programRunning
}

// isAuto reports whether the controller is in AUTO, in the machine's engine numbering.
func (state machineState) isAuto() bool {
	return state.opMode == autoOpModes[numberingEngine(state.engine)]
}

// isAlarmActive reports whether an NC alarm or the emergency stop is active.
func (state machineState) isAlarmActive() bool {
	return state.alarmNumber != 0 || state.emergency
}

// modeDecision is the mode a reading decides: the machine_modes id, the rule that matched and
// conditions, the log attributes of only the inputs that rule checks.
type modeDecision struct {
	mode       int
	reason     string
	conditions []any
}

// offlineDecision is the decision for a machine the DLL does not reach (no reading, a failed
// connect or an IsConnected=0 reading).
var offlineDecision = modeDecision{
	mode: model.ModeOffline, reason: reasonNoConnection, conditions: []any{"connected", false},
}

// The condition attributes of the mode rules. opMode and procStatus are named the way the
// machine's engine numbers them; programState is what the rules read from procStatus.

func (state machineState) opModeCondition() []any {
	return []any{"opMode", fmt.Sprintf("%d (%s)", state.opMode, opModeName(state.engine, state.opMode))}
}

func (state machineState) programConditions() []any {
	return []any{"procStatus", state.procStatusText(), "programState", programStateNames[state.program]}
}

func (state machineState) feedOverrideCondition() []any {
	return []any{"feedOverride", state.feedOverride}
}

func (state machineState) alarmConditions() []any {
	return []any{"alarmNo", state.alarmNumber, "emergency", state.emergency}
}

// procStatusText is the raw ProcStatus with its name for the machine's engine, e.g. "3 (START)".
func (state machineState) procStatusText() string {
	return fmt.Sprintf("%d (%s)", state.procStatus, procStatusName(state.engine, state.procStatus))
}

// machineStateFromReading reads the NC state from a DLL reading of a machine with engine
// (controllers.cnc_engine, nil when not set).
func machineStateFromReading(engine *int, reading cnc.Reading) machineState {
	machineEngine := engineOf(engine)
	procStatus := int(reading.Values[cnc.FieldProcStatus].Number)
	return machineState{
		engine:       machineEngine,
		connected:    reading.Connected,
		opMode:       int(reading.Values[cnc.FieldMachineMode].Number),
		procStatus:   procStatus,
		program:      programStateFromProcStatus(machineEngine, procStatus),
		feedOverride: reading.Values[cnc.FieldFeedOverride].Number,
		alarmNumber:  int(reading.Values[cnc.FieldAlarmNo].Number),
		alarmMessage: reading.Values[cnc.FieldAlarmMsg].Text,
		emergency:    reading.Values[cnc.FieldIsEmergency].Number != 0,
		programName:  reading.Values[cnc.FieldProgName].Text,
	}
}

// programStateFromProcStatus turns a ProcStatus into a programState in engine's numbering; a
// value the table does not have counts as reset.
func programStateFromProcStatus(engine, procStatus int) programState {
	return programStatesByEngine[numberingEngine(engine)][procStatus]
}

// decideMachineMode returns the mode state decides, with the rule that decided it and the
// inputs that rule checks. First match wins:
//
//	Offline (9)        no NC connection
//	Downtime (6)       alarm active and program not running
//	In Production (3)  AUTO, program running, feed override > 0, no alarm
//	Hold (4)           program held, or running with feed override ≤ 0, an alarm or not in AUTO
//	Idle (5)           program reset and no alarm, in any operation mode
//
// Every connected state matches one of these rules: the last return is never reached.
func decideMachineMode(state machineState) modeDecision {
	running := state.isRunning()
	alarm := state.isAlarmActive()
	switch {
	case !state.connected:
		return offlineDecision
	case alarm && !running:
		return modeDecision{model.ModeDowntime, "alarm active, no program running",
			slices.Concat(state.alarmConditions(), state.programConditions())}
	case state.isAuto() && running && state.feedOverride > 0 && !alarm:
		return modeDecision{model.ModeInProduction, "AUTO, program running, feed override > 0, no alarm",
			slices.Concat(state.opModeCondition(), state.programConditions(), state.feedOverrideCondition(),
				state.alarmConditions())}
	}
	if reason, conditions := holdReason(state); reason != "" {
		return modeDecision{model.ModeHold, reason, conditions}
	}
	if state.program == programReset && !alarm {
		return modeDecision{model.ModeIdle, "program reset, no alarm (any mode)",
			slices.Concat(state.programConditions(), state.alarmConditions())}
	}
	return modeDecision{model.ModeIdle, "no rule matched",
		slices.Concat(state.opModeCondition(), state.programConditions(), state.feedOverrideCondition(),
			state.alarmConditions())}
}

// holdReason returns why a program is held (Fanuc Stop / HOLD / MSTR, Siemens Stop / Waiting /
// Interrupted, or running but not cutting: feed override 0, an alarm, or not in AUTO) with the
// inputs that show it, or "" when it is not held.
func holdReason(state machineState) (string, []any) {
	switch {
	case state.program == programHeld:
		return "program held", state.programConditions()
	case !state.isRunning():
		return "", nil
	case state.feedOverride <= 0:
		return "program running, feed override ≤ 0", slices.Concat(state.programConditions(), state.feedOverrideCondition())
	case state.isAlarmActive():
		return "program running, alarm active", slices.Concat(state.programConditions(), state.alarmConditions())
	case !state.isAuto():
		return "program running, not in AUTO", slices.Concat(state.programConditions(), state.opModeCondition())
	}
	return "", nil
}

// recordModeFromReading decides the mode of a connected reading and records it when it changed.
func (machineWorker *worker) recordModeFromReading(at time.Time, reading cnc.Reading) {
	state := machineStateFromReading(machineWorker.machine.CNCEngine, reading)
	decision := decideMachineMode(state)
	machineWorker.logModeCheck(decision)
	machineWorker.logDowntimeCheck(state, decision.mode, decision.reason)
	machineWorker.recordModeChange(at, decision)
}
