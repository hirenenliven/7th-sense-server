package monitor

import (
	"strconv"

	"7th-sense-server/internal/model"
)

// Variables the alarm events and the mode log read (cnc_variables ids, the old mm_variables
// ids). The mode itself is decided from the reading (machine_mode.go).
const (
	variableProgName  = 15
	variableAlarm     = 25
	variableEmergency = 200
	variableAlarmMsg  = 207
)

// modeName names a machine_modes id for logs.
func modeName(mode int) string {
	switch mode {
	case model.ModeSetup:
		return "Setup"
	case model.ModeInProduction:
		return "In Production"
	case model.ModeHold:
		return "Hold"
	case model.ModeIdle:
		return "Idle"
	case model.ModeDowntime:
		return "Downtime"
	case model.ModeOffline:
		return "Offline"
	case model.ModeRapid:
		return "Rapid Mode"
	}
	return strconv.Itoa(mode)
}
