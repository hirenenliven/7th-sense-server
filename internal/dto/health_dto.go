package dto

import (
	"time"

	"7th-sense-server/internal/model"
)

// MachineModeHistoryDays is how many days back GET /health/machines/{id}/modes reads.
const MachineModeHistoryDays = 10

// MachineModeTotal is the time a machine spent in one mode within the history window.
type MachineModeTotal struct {
	ModeID      int    `json:"modeId"`
	ModeName    string `json:"modeName"`
	PeriodCount int    `json:"periodCount"`
	DurationSec int    `json:"durationSec"`
}

// MachineModesResponse is the response of GET /health/machines: every machine that is not
// deleted with the mode it is in now, ordered by name. At is the server time it was read.
type MachineModesResponse struct {
	At       time.Time                 `json:"at"`
	Machines []model.MachineModeStatus `json:"machines"`
}

// MachineModeHistoryResponse is the response of GET /health/machines/{id}/modes: the
// machine's mode periods of the last MachineModeHistoryDays days, oldest first, and the
// time per mode.
type MachineModeHistoryResponse struct {
	MachineID     int64                     `json:"machineId"`
	MachineName   string                    `json:"machineName"`
	CurrentModeID *int64                    `json:"currentModeId"`
	OnlineStatus  bool                      `json:"onlineStatus"`
	From          time.Time                 `json:"from"`
	To            time.Time                 `json:"to"`
	Periods       []model.MachineModePeriod `json:"periods"`
	Totals        []MachineModeTotal        `json:"totals"`
}
