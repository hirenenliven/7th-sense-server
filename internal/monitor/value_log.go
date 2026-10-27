package monitor

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"7th-sense-server/internal/model"
)

// Variable value change logging lives in this file only: CNC_LOG_VALUES=true turns it on,
// false (the default) turns it off. processReading collects the changed values of one reading
// with valueChangeLog.add and logs them together with valueChangeLog.flush.

// valueChange is one variable whose value changed in a reading.
type valueChange struct {
	variable model.VariableConfig
	previous *float64 // nil for the first value the machine sends for the variable
	value    float64
}

// valueChangeLog collects the changed values of one reading. When logging is off it collects
// nothing, so a reading costs no extra work.
type valueChangeLog struct {
	enabled bool
	changes []valueChange
}

// newValueChangeLog starts the change log of one reading; enabled is CNC_LOG_VALUES.
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

// flush logs every changed value of the reading in one line: how many changed and, per
// variable, "id:code previous -> value" (previous is "none" for a first value).
func (changeLog *valueChangeLog) flush(machine model.MonitoredMachine, readAt time.Time) {
	if !changeLog.enabled || len(changeLog.changes) == 0 {
		return
	}
	changes := make([]string, 0, len(changeLog.changes))
	for _, change := range changeLog.changes {
		changes = append(changes, fmt.Sprintf("%d:%s %s -> %s", change.variable.ID, change.variable.Code,
			formatPreviousValue(change.previous), formatValue(change.value)))
	}
	slog.Info("cnc: variable values changed", "machineId", machine.ID, "machine", machine.Name,
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
