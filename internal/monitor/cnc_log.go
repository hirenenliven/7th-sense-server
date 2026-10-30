package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
)

// Every log line of the CNC collector goes through writeLog in this file (through
// Collector.log and worker.log), with a log type and a level. These .env settings decide what is
// written (restart to apply):
//
//	CNC_MACHINE_IDS=225,247   only these machines are connected, so only they are logged;
//	                          ALL or empty = every active machine (collector.go)
//	CNC_LOG_EVENTS=true        logTypeEvent: connect, disconnect, alarms, DLL errors
//	CNC_LOG_MODES=true         logTypeMode: mode changes
//	CNC_LOG_MODE_CHECKS=true   logTypeModeCheck: the mode the NC inputs decide, whenever they change
//	CNC_LOG_VARIABLES=false    logTypeVariable: the changed values of each reading, one line
//	CNC_LOG_OTHER=true         logTypeOther: start, stop, connections, reading rate, …
//	CNC_LOG_SNAPSHOT=false     logTypeSnapshot: the raw DLL status fields, whenever one changes
//	CNC_LOG_DOWNTIME=true      logTypeDowntime: the Downtime inputs (machine type, alarm, E-stop,
//	                          program status) and whether the machine is Downtime, whenever one changes
//
// Every error is written, whatever the switches say.
// A line about a machine starts with "[machine 225 VM14]" and has machineId and machine; every
// line has logType, so app.log can be searched by machine or by type.

// logType is the kind of a CNC log line; CNCLogSettings switches each kind on or off.
type logType string

const (
	logTypeEvent     logType = "event"
	logTypeVariable  logType = "variable"
	logTypeMode      logType = "mode"
	logTypeModeCheck logType = "modeCheck"
	logTypeSnapshot  logType = "snapshot"
	logTypeDowntime  logType = "downtime"
	logTypeOther     logType = "other"
)

// isLogTypeOn reports whether a line of kind at level is written: errors always, the rest
// when their setting is on.
func isLogTypeOn(settings config.CNCLogSettings, kind logType, level slog.Level) bool {
	if level >= slog.LevelError {
		return true
	}
	switch kind {
	case logTypeEvent:
		return settings.Events
	case logTypeMode:
		return settings.Modes
	case logTypeModeCheck:
		return settings.ModeChecks
	case logTypeVariable:
		return settings.Variables
	case logTypeSnapshot:
		return settings.Snapshot
	case logTypeDowntime:
		return settings.Downtime
	}
	return settings.Other
}

// writeLog writes one CNC log line when its kind is on. machine is the machine the line is
// about (nil for none): its id and name start the message and are added as attributes.
func writeLog(settings config.CNCLogSettings, kind logType, level slog.Level, machine *model.MonitoredMachine,
	message string, attributes ...any) {
	if !isLogTypeOn(settings, kind, level) {
		return
	}
	lineAttributes := []any{"logType", string(kind)}
	if machine != nil {
		message = fmt.Sprintf("[machine %d %s] %s", machine.ID, machine.Name, message)
		lineAttributes = append(lineAttributes, "machineId", machine.ID, "machine", machine.Name)
	}
	slog.Log(context.Background(), level, message, append(lineAttributes, attributes...)...)
}

// log writes a collector line that is about no single machine (pass "machineId" as an
// attribute when only the id is known).
func (collector *Collector) log(kind logType, level slog.Level, message string, attributes ...any) {
	writeLog(collector.config.Logs, kind, level, nil, message, attributes...)
}

// logMachine writes a collector line about machine.
func (collector *Collector) logMachine(kind logType, level slog.Level, machine model.MonitoredMachine, message string,
	attributes ...any) {
	writeLog(collector.config.Logs, kind, level, &machine, message, attributes...)
}

// log writes a line about the worker's machine.
func (machineWorker *worker) log(kind logType, level slog.Level, message string, attributes ...any) {
	writeLog(machineWorker.collector.config.Logs, kind, level, &machineWorker.machine, message, attributes...)
}

// logCollectionSettings logs once, when the collector starts, which machines it connects and
// which CNC logs are on.
func (collector *Collector) logCollectionSettings() {
	machines := any(config.AllMachines)
	if !collector.config.Machines.IsAll() {
		ids := make([]int64, 0, len(collector.config.Machines.IDs))
		for id := range collector.config.Machines.IDs {
			ids = append(ids, id)
		}
		slices.Sort(ids)
		machines = ids
	}
	logs := collector.config.Logs
	// Written whatever CNC_LOG_OTHER says, so app.log always shows what is collected and logged.
	slog.Info("cnc collector: settings", "logType", string(logTypeOther), "machines", machines,
		"logEvents", logs.Events, "logModes", logs.Modes, "logModeChecks", logs.ModeChecks,
		"logVariables", logs.Variables, "logOther", logs.Other, "logSnapshot", logs.Snapshot,
		"logDowntime", logs.Downtime)
}

// ---------------------------------------------------------------- snapshot (CNC_LOG_SNAPSHOT)

// snapshotStatus is the part of a DLL reading the snapshot log compares: every status and
// message field, raw as the DLL sends it. Positions, loads and speeds are left out: they change
// on almost every reading.
type snapshotStatus struct {
	connected       bool
	missingNodes    string
	machineMode     float64
	procStatus      float64
	isEmergency     float64
	feedOverride    float64
	spindleOverride float64
	alarmNo         float64
	alarmMessage    string
	programName     string
	status          cnc.Status
}

func snapshotStatusOf(reading cnc.Reading) snapshotStatus {
	return snapshotStatus{
		connected:       reading.Connected,
		missingNodes:    reading.MissingNodes,
		machineMode:     reading.Values[cnc.FieldMachineMode].Number,
		procStatus:      reading.Values[cnc.FieldProcStatus].Number,
		isEmergency:     reading.Values[cnc.FieldIsEmergency].Number,
		feedOverride:    reading.Values[cnc.FieldFeedOverride].Number,
		spindleOverride: reading.Values[cnc.FieldSpindleOverride].Number,
		alarmNo:         reading.Values[cnc.FieldAlarmNo].Number,
		alarmMessage:    reading.Values[cnc.FieldAlarmMsg].Text,
		programName:     reading.Values[cnc.FieldProgName].Text,
		status:          reading.Status,
	}
}

// logSnapshot logs a reading's raw status fields, with the DLL's field names, when they differ
// from the last ones logged for the machine (the first reading of a connection is always logged).
// MachineMode and ProcStatus are also named for the machine's engine. The NC timer
// (TimerMinute, TimerMsec) is logged with them but not compared: it counts on every reading
// while a program runs.
func (machineWorker *worker) logSnapshot(reading cnc.Reading) {
	if !machineWorker.collector.config.Logs.Snapshot {
		return
	}
	current := snapshotStatusOf(reading)
	if machineWorker.lastSnapshot != nil && *machineWorker.lastSnapshot == current {
		return
	}
	machineWorker.lastSnapshot = &current
	engine := engineOf(machineWorker.machine.CNCEngine)
	machineWorker.log(logTypeSnapshot, slog.LevelInfo, "cnc: snapshot status",
		"IsConnected", current.connected, "ResultStatus", current.status.ResultStatus,
		"MachineMode", current.machineMode, "machineModeName", opModeName(engine, int(current.machineMode)),
		"ProcStatus", current.procStatus, "procStatusName", procStatusName(engine, int(current.procStatus)),
		"IsEmergency", current.isEmergency,
		"TimerMinute", reading.Values[cnc.FieldTimerMinute].Number, "TimerMsec", reading.Values[cnc.FieldTimerMsec].Number,
		"FeedRatePercetnage", current.feedOverride, "SpindeOverridePercentage", current.spindleOverride,
		"CmdFeedOvr", current.status.CmdFeedOverride, "CmdSpOvr", current.status.CmdSpindleOverride,
		"AlarmNo", current.alarmNo, "AlarmMsg", current.alarmMessage,
		"OpMsgNo", current.status.OperatorMessageNo, "OpMsg", current.status.OperatorMessage,
		"ActiveGCode", current.status.ActiveGCode, "ProgName", current.programName,
		"missingNodes", current.missingNodes)
}

// ---------------------------------------------------------------- events (CNC_LOG_EVENTS)

// logMachineEvent logs a machine event; failures (connect failed, configuration and DLL errors,
// stale connection) as warnings.
func (machineWorker *worker) logMachineEvent(eventType string, code *int, message string) {
	attributes := []any{"event", eventType}
	if code != nil {
		attributes = append(attributes, "code", *code)
	}
	if message != "" {
		attributes = append(attributes, "message", message)
	}
	level := slog.LevelInfo
	switch eventType {
	case model.EventConnectFailed, model.EventConfigError, model.EventDLLError, model.EventStale:
		level = slog.LevelWarn
	}
	machineWorker.log(logTypeEvent, level, "cnc machine event "+eventType, attributes...)
}

// ---------------------------------------------------------------- modes (CNC_LOG_MODES, CNC_LOG_MODE_CHECKS)

// logModeChange logs that the machine's mode changes to mode: the rule that decided it and
// the NC inputs it was decided from (none for Offline).
func (machineWorker *worker) logModeChange(mode int, reason string, inputs []any) {
	from, to := modeName(machineWorker.mode), modeName(mode)
	attributes := []any{"from", from, "to", to, "reason", reason, "program", machineWorker.program}
	machineWorker.log(logTypeMode, slog.LevelInfo, fmt.Sprintf("cnc: mode changed %s -> %s", from, to),
		append(attributes, inputs...)...)
}

// logModeCheck logs the NC inputs of a reading with the mode they decide, when the inputs
// differ from the last ones logged (not on every reading), so a mode can be compared with the
// HMI even while it does not change.
func (machineWorker *worker) logModeCheck(state machineState, mode int, reason string) {
	if !machineWorker.collector.config.Logs.ModeChecks ||
		(machineWorker.lastModeCheck != nil && *machineWorker.lastModeCheck == state) {
		return
	}
	machineWorker.lastModeCheck = &state
	attributes := []any{"decided", modeName(mode), "current", modeName(machineWorker.mode),
		"modeChanges", mode != machineWorker.mode, "reason", reason}
	machineWorker.log(logTypeModeCheck, slog.LevelInfo, "cnc: mode check "+modeName(mode),
		append(attributes, state.logAttributes()...)...)
}

// ---------------------------------------------------------------- downtime (CNC_LOG_DOWNTIME)

// engineNames name the CNC engines (controllers.cnc_engine) for the log.
var engineNames = map[int]string{
	model.EngineSimulator: "Simulator", model.EngineFanuc: "Fanuc", model.EngineSiemens: "Siemens",
	model.EngineSyntec: "Syntec", model.EngineMitsubishi: "Mitsubishi",
}

// engineName names a machine's CNC engine; nil is "none".
func engineName(engine *int) string {
	if engine == nil {
		return "none"
	}
	if name, found := engineNames[*engine]; found {
		return name
	}
	return fmt.Sprintf("unknown(%d)", *engine)
}

// downtimeCheck is what the Downtime rule (alarm active and no program running) is decided
// from, with the decision; logDowntimeCheck compares it with the last one logged.
type downtimeCheck struct {
	alarmNumber  int
	alarmMessage string
	emergency    bool
	procStatus   int
	program      programState
	mode         int
}

// logDowntimeCheck logs, when it differs from the last one logged for the machine, whether the
// machine is Downtime and why: its machine type, alarm, E-stop and program status. The decided
// mode is added because a machine that is not Downtime may still be Hold or Idle.
func (machineWorker *worker) logDowntimeCheck(state machineState, mode int, reason string) {
	if !machineWorker.collector.config.Logs.Downtime {
		return
	}
	current := downtimeCheck{
		alarmNumber: state.alarmNumber, alarmMessage: state.alarmMessage, emergency: state.emergency,
		procStatus: state.procStatus, program: state.program, mode: mode,
	}
	if machineWorker.lastDowntimeCheck != nil && *machineWorker.lastDowntimeCheck == current {
		return
	}
	machineWorker.lastDowntimeCheck = &current
	isDowntime := mode == model.ModeDowntime
	machineWorker.log(logTypeDowntime, slog.LevelInfo, fmt.Sprintf("cnc: downtime check downtime=%t", isDowntime),
		"machineType", engineName(machineWorker.machine.CNCEngine),
		"downtime", isDowntime,
		"alarmActive", state.isAlarmActive(),
		"alarmNo", state.alarmNumber,
		"alarmMsg", state.alarmMessage,
		"emergency", state.emergency,
		"programRunning", state.isRunning(),
		"procStatus", state.procStatusText(),
		"programState", programStateNames[state.program],
		"mode", modeName(mode),
		"reason", reason)
}

// ---------------------------------------------------------------- variables (CNC_LOG_VARIABLES)

// valueChange is one variable whose value changed in a reading.
type valueChange struct {
	variable model.VariableConfig
	previous *float64 // nil for the first value the machine sends for the variable
	value    float64
}

// valueChangeLog collects the changed values of one reading; processReading adds them with
// add and logs them together with flush. When logging is off it collects nothing, so a reading
// costs no extra work.
type valueChangeLog struct {
	enabled bool
	changes []valueChange
}

// newValueChangeLog starts the change log of one reading; enabled is CNC_LOG_VARIABLES.
func newValueChangeLog(enabled bool) *valueChangeLog {
	return &valueChangeLog{enabled: enabled}
}

// add records that variable changed from previous to value.
func (changeLog *valueChangeLog) add(variable model.VariableConfig, previous *float64, value float64) {
	if !changeLog.enabled {
		return
	}
	changeLog.changes = append(changeLog.changes, valueChange{variable: variable, previous: previous, value: value})
}

// flush logs every changed value of the reading in one line for machineWorker's machine: how
// many changed and, per variable, "id:code previous -> value" (previous is "none" for a first
// value).
func (changeLog *valueChangeLog) flush(machineWorker *worker, readAt time.Time) {
	if !changeLog.enabled || len(changeLog.changes) == 0 {
		return
	}
	changes := make([]string, 0, len(changeLog.changes))
	for _, change := range changeLog.changes {
		changes = append(changes, fmt.Sprintf("%d:%s %s -> %s", change.variable.ID, change.variable.Code,
			formatPreviousValue(change.previous), formatValue(change.value)))
	}
	machineWorker.log(logTypeVariable, slog.LevelInfo, "cnc: variable values changed",
		"count", len(changes), "changes", changes, "readAt", readAt.Format(time.RFC3339Nano))
}

// formatPreviousValue formats a previous value for the log; nil is "none".
func formatPreviousValue(previous *float64) string {
	if previous == nil {
		return "none"
	}
	return formatValue(*previous)
}

// formatValue formats a value without trailing zeros (12.5, not 12.500000).
func formatValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
