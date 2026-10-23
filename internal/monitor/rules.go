package monitor

import "7th-sense-server/internal/model"

// Variables the mode rules and events read (cnc_variables ids, the old mm_variables ids).
const (
	variableRapidMode  = 5
	variableOpMode     = 14
	variableProgName   = 15
	variableAlarm      = 25
	variableProgStatus = 55
	variableEmergency  = 200
	variableAlarmMsg   = 207
)

// opMode and progStatus values from the controller.
const (
	opModeAuto            = 2
	progStatusStop        = 1
	progStatusActive      = 2
	progStatusWaiting     = 3
	progStatusInterrupted = 4
)

// modeFor turns the machine's state into a machine_modes id. First match wins:
//
//	not connected                          9 Offline
//	emergency stop or an alarm             6 Downtime
//	JOG / MDA                              1 Setup
//	AUTO, running, rapid (G0)              12 Rapid Mode
//	AUTO, running                          3 In Production
//	AUTO, stopped / waiting / interrupted  4 Hold
//	otherwise (reset)                      5 Idle
//
// values holds the number variables of the reading by variable id; a missing variable counts
// as 0, a missing opMode as AUTO (the program status alone decides).
func modeFor(connected bool, values map[int]float64) int {
	if !connected {
		return model.ModeOffline
	}
	if values[variableEmergency] != 0 || values[variableAlarm] != 0 {
		return model.ModeDowntime
	}
	if opMode, ok := values[variableOpMode]; ok && opMode != opModeAuto {
		return model.ModeSetup
	}
	switch values[variableProgStatus] {
	case progStatusActive:
		if values[variableRapidMode] == 1 {
			return model.ModeRapid
		}
		return model.ModeInProduction
	case progStatusStop, progStatusWaiting, progStatusInterrupted:
		return model.ModeHold
	}
	return model.ModeIdle
}
