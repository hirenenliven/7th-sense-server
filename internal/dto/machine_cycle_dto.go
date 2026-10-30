package dto

import (
	"fmt"
	"net/http"
	"time"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

// MachineCycleMaxDays is the longest date range GET /machine-cycles reads at once.
const MachineCycleMaxDays = 31

// MachineCycleHistoryHours is how many hours back GET /health/machines/{id}/cycles reads (the
// mode chart page). Shorter than the mode history: a machine can run thousands of cycles a day.
const MachineCycleHistoryHours = 24

// MachineCycleListRequest is the query of GET /machine-cycles: the machine and the days its
// cycles started on, From … To (YYYY-MM-DD, both included). An empty To is today, an empty
// From is To.
type MachineCycleListRequest struct {
	MachineID int64
	From      string
	To        string
}

// Validate checks the machine and the dates: From not after To, at most MachineCycleMaxDays days.
func (listRequest *MachineCycleListRequest) Validate() error {
	if listRequest.MachineID <= 0 {
		return response.SelectMachine
	}
	if listRequest.To != "" && !validDate(listRequest.To) {
		return response.NewError(http.StatusBadRequest, "Please enter valid to date")
	}
	if listRequest.From != "" && !validDate(listRequest.From) {
		return response.NewError(http.StatusBadRequest, "Please enter valid from date")
	}
	// to is the start of the day after To, so From after To leaves no time between them
	from, to := listRequest.Range(time.Now())
	if !from.Before(to) {
		return response.NewError(http.StatusBadRequest, "From date must not be after to date")
	}
	if to.Sub(from) > MachineCycleMaxDays*24*time.Hour {
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("Please select at most %d days", MachineCycleMaxDays))
	}
	return nil
}

// Range returns the time span the cycles started in: from the start of the From day to the
// start of the day after To, in the server's time zone; now gives today for an empty To.
func (listRequest *MachineCycleListRequest) Range(now time.Time) (from, to time.Time) {
	toDay := listRequest.To
	if toDay == "" {
		toDay = now.Format(dateLayout)
	}
	fromDay := listRequest.From
	if fromDay == "" {
		fromDay = toDay
	}
	fromDate, _ := time.ParseInLocation(dateLayout, fromDay, time.Local)
	toDate, _ := time.ParseInLocation(dateLayout, toDay, time.Local)
	return fromDate, toDate.AddDate(0, 0, 1)
}

// MachineCycleSummary adds up the cycles of a GET /machine-cycles response.
// AverageCycleSec is the mean length of the completed cycles (0 when there is none).
type MachineCycleSummary struct {
	CycleCount      int `json:"cycleCount"`
	CompletedCount  int `json:"completedCount"`
	UnfinishedCount int `json:"unfinishedCount"`
	RunningCount    int `json:"runningCount"`
	PartsMade       int `json:"partsMade"`
	ProductionSec   int `json:"productionSec"`
	HoldSec         int `json:"holdSec"`
	DowntimeSec     int `json:"downtimeSec"`
	OfflineSec      int `json:"offlineSec"`
	AverageCycleSec int `json:"averageCycleSec"`
}

// MachineCyclesResponse is the response of GET /machine-cycles: the machine's cycles that
// started in From … To (To excluded), oldest first, and their totals.
type MachineCyclesResponse struct {
	MachineID   int64                `json:"machineId"`
	MachineName string               `json:"machineName"`
	From        time.Time            `json:"from"`
	To          time.Time            `json:"to"`
	Cycles      []model.MachineCycle `json:"cycles"`
	Summary     MachineCycleSummary  `json:"summary"`
}
