package dto

import (
	"fmt"
	"net/http"
	"time"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

// ModeTimelineMaxMinutes is the longest "last N minutes" window of GET /productivity/mode-timeline.
const ModeTimelineMaxMinutes = 24 * 60

// Machine mode groups of the dashboard charts (client constants/machines.js MACHINE_MODES).
const (
	ModeGroupOffline    = "offline"
	ModeGroupProduction = "production"
	ModeGroupIdle       = "idle"
	ModeGroupHold       = "hold"
	ModeGroupDowntime   = "downtime"
)

// DashboardRequest is the query of the /productivity and /machine-health charts: the machines,
// the days From … To (YYYY-MM-DD, both included; an empty To is today, an empty From is To),
// the shifts (none = the whole day), the program (empty = every program), for OEE whether
// availability counts the scheduled shift time only, for the mode timeline the last Minutes
// up to now (0 = the To day), and for the mode totals AllTime (every recorded period, no dates).
type DashboardRequest struct {
	MachineIDs   []int64
	From         string
	To           string
	ShiftIDs     []int64
	Program      string
	WithSchedule bool
	Minutes      int
	AllTime      bool
}

// Validate checks the machines and the dates (at most MachineCycleMaxDays days).
func (dashboardRequest *DashboardRequest) Validate() error {
	if dashboardRequest.Minutes < 0 || dashboardRequest.Minutes > ModeTimelineMaxMinutes {
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("Please select at most %d minutes", ModeTimelineMaxMinutes))
	}
	// All time has no dates, so only the machine is checked (no day limit)
	if dashboardRequest.AllTime {
		if len(dashboardRequest.MachineIDs) == 0 || dashboardRequest.MachineIDs[0] <= 0 {
			return response.SelectMachine
		}
		return nil
	}
	listRequest := dashboardRequest.cycleRequest()
	return listRequest.Validate()
}

// Range returns the start of the From day and the start of the day after To (server time zone).
func (dashboardRequest *DashboardRequest) Range(now time.Time) (from, to time.Time) {
	listRequest := dashboardRequest.cycleRequest()
	return listRequest.Range(now)
}

func (dashboardRequest *DashboardRequest) cycleRequest() MachineCycleListRequest {
	listRequest := MachineCycleListRequest{From: dashboardRequest.From, To: dashboardRequest.To}
	if len(dashboardRequest.MachineIDs) > 0 {
		listRequest.MachineID = dashboardRequest.MachineIDs[0]
	}
	return listRequest
}

// The Productivity and Machine Health charts return only what their page draws. Each row type
// below is filled by one constructor: to add or drop a field, change the struct and that
// constructor (and the client chart that reads it).

// ModeTimelineResponse is the response of GET /productivity/mode-timeline: the machine's mode
// periods of one day or of the last minutes, clipped to From … To (the open period has no EndedAt
// and runs up to ServerTime), the program it runs now (its latest period's) and its operation mode.
type ModeTimelineResponse struct {
	ProgramName string           `json:"programName"`
	OpMode      string           `json:"opMode"`
	From        time.Time        `json:"from"`
	To          time.Time        `json:"to"`
	ServerTime  time.Time        `json:"serverTime"`
	Periods     []TimelinePeriod `json:"periods"`
}

// ModeTotalsResponse is the response of GET /productivity/mode-totals: the time the machine spent
// in each mode within From … To (periods clipped to it; the open period runs up to ServerTime).
// For all time From is the start of the first recorded period (nil when nothing is recorded).
type ModeTotalsResponse struct {
	From       *time.Time           `json:"from"`
	To         time.Time            `json:"to"`
	ServerTime time.Time            `json:"serverTime"`
	Totals     []model.ModeDuration `json:"totals"`
}

// TimelinePeriod is one bar of the mode Gantt chart.
type TimelinePeriod struct {
	ID        int64      `json:"id"`
	ModeID    int        `json:"modeId"`
	ModeName  string     `json:"modeName"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt"`
}

// NewTimelinePeriod returns the chart fields of period.
func NewTimelinePeriod(period model.MachineModePeriod) TimelinePeriod {
	return TimelinePeriod{
		ID:        period.ID,
		ModeID:    period.ModeID,
		ModeName:  period.ModeName,
		StartedAt: period.StartedAt,
		EndedAt:   period.EndedAt,
	}
}

// OEERow is one day of GET /productivity/oee, averaged over the machines (percent). Performance,
// Quality and OEE are nil when no cycle closed that day.
type OEERow struct {
	Date         string   `json:"date"`
	Availability float64  `json:"availability"`
	Performance  *float64 `json:"performance"`
	Quality      *float64 `json:"quality"`
	OEE          *float64 `json:"oee"`
}

// OEEResponse is the response of GET /productivity/oee: one row per day with planned time, and
// the programs the machines ran in the range (for the program filter).
type OEEResponse struct {
	Rows     []OEERow `json:"rows"`
	Programs []string `json:"programs"`
}

// CycleCountRow is one day of GET /productivity/cycle-counts, summed over the machines.
// Produced counts every cycle started; running cycles are neither finished nor unfinished.
type CycleCountRow struct {
	Date       string `json:"date"`
	Produced   int    `json:"produced"`
	Finished   int    `json:"finished"`
	Unfinished int    `json:"unfinished"`
}

// FinishCyclesResponse is the response of GET /productivity/finish-cycles: the machine's
// completed cycles of the day (in the shifts and program picked), the totals of every cycle
// started then, and the programs it ran that day.
type FinishCyclesResponse struct {
	Cycles         []FinishedCycle     `json:"cycles"`
	Summary        FinishCyclesSummary `json:"summary"`
	LastFinishedAt *time.Time          `json:"lastFinishedAt"`
	Programs       []string            `json:"programs"`
}

// FinishedCycle is one bar of the Finish Cycle Summary: when it finished, its program, its length
// and its seconds per machine mode.
type FinishedCycle struct {
	EndedAt       *time.Time `json:"endedAt"`
	ProgramName   string     `json:"programName"`
	DurationSec   int        `json:"durationSec"`
	ProductionSec int        `json:"productionSec"`
	HoldSec       int        `json:"holdSec"`
	DowntimeSec   int        `json:"downtimeSec"`
	OfflineSec    int        `json:"offlineSec"`
}

// NewFinishedCycle returns the chart fields of cycle.
func NewFinishedCycle(cycle model.MachineCycle) FinishedCycle {
	return FinishedCycle{
		EndedAt:       cycle.EndedAt,
		ProgramName:   cycle.ProgramName,
		DurationSec:   cycle.DurationSec,
		ProductionSec: cycle.ProductionSec,
		HoldSec:       cycle.HoldSec,
		DowntimeSec:   cycle.DowntimeSec,
		OfflineSec:    cycle.OfflineSec,
	}
}

// FinishCyclesSummary is the day's totals behind the Productivity KPI tiles.
type FinishCyclesSummary struct {
	CycleCount           int `json:"cycleCount"`
	CompletedCount       int `json:"completedCount"`
	UnfinishedCount      int `json:"unfinishedCount"`
	AverageCycleSec      int `json:"averageCycleSec"`
	AverageProductionSec int `json:"averageProductionSec"`
}

// NewFinishCyclesSummary returns the KPI fields of summary.
func NewFinishCyclesSummary(summary MachineCycleSummary) FinishCyclesSummary {
	return FinishCyclesSummary{
		CycleCount:           summary.CycleCount,
		CompletedCount:       summary.CompletedCount,
		UnfinishedCount:      summary.UnfinishedCount,
		AverageCycleSec:      summary.AverageCycleSec,
		AverageProductionSec: summary.AverageProductionSec,
	}
}

// MachineTimeRow is one day of GET /machine-health/machine-time-by-day: hours per mode group
// (time without a mode period counts as offline; today runs up to now).
type MachineTimeRow struct {
	Date       string  `json:"date"`
	Offline    float64 `json:"offline"`
	Production float64 `json:"production"`
	Idle       float64 `json:"idle"`
	Hold       float64 `json:"hold"`
	Downtime   float64 `json:"downtime"`
}
