package monitor

import (
	"fmt"
	"log/slog"
	"slices"
)

// Debug logging of chosen machines lives in this file: CNC_DEBUG_MACHINE_IDS="225,247" turns it
// on for those machines only; every machine is still collected and stored. For a debug machine:
//
//   - "cnc debug: mode check [machine 225 VM14] ..." is logged whenever the NC inputs of the mode
//     change (not on every reading), with the mode they decide, the rule, and whether the
//     recorded mode changes;
//   - its changed variable values are logged as with CNC_LOG_VALUES=true (value_log.go).

// logDebugMachines logs the debug machines once when the collector starts.
func logDebugMachines(machineIDs map[int64]bool) {
	if len(machineIDs) == 0 {
		return
	}
	ids := make([]int64, 0, len(machineIDs))
	for id := range machineIDs {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	slog.Info("cnc debug: logging machines in detail (CNC_DEBUG_MACHINE_IDS)", "machineIds", ids)
}

// logModeCheck logs the NC inputs of a debug machine's reading with the mode they decide, when
// the inputs differ from the last ones logged.
func (machineWorker *worker) logModeCheck(state machineState, mode int, reason string) {
	if !machineWorker.debug || (machineWorker.lastModeCheck != nil && *machineWorker.lastModeCheck == state) {
		return
	}
	machineWorker.lastModeCheck = &state
	message := fmt.Sprintf("cnc debug: mode check %s %s", machineWorker.logLabel(), modeName(mode))
	attributes := []any{"machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name,
		"decided", modeName(mode), "current", modeName(machineWorker.mode), "modeChanges", mode != machineWorker.mode,
		"reason", reason}
	slog.Info(message, append(attributes, state.logAttributes()...)...)
}
