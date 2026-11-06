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

// OEERow is one day of GET /productivity/oee over the machines together, weighted by their
// planned / run time and closed cycles (percent). Performance,
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

// MachineSummary is one machine of GET /productivity/machine-summaries: its OEE parts and cycle
// counts of one day (the whole day, like GET /productivity/oee and /cycle-counts without shifts).
// Availability is nil when the day has no planned time yet; Performance, Quality and OEE are nil
// when no cycle closed. CycleCount counts every cycle started; running ones are neither
// completed nor unfinished.
type MachineSummary struct {
	MachineID       int64    `json:"machineId"`
	Availability    *float64 `json:"availability"`
	Performance     *float64 `json:"performance"`
	Quality         *float64 `json:"quality"`
	OEE             *float64 `json:"oee"`
	CycleCount      int      `json:"cycleCount"`
	CompletedCount  int      `json:"completedCount"`
	UnfinishedCount int      `json:"unfinishedCount"`
}

// NewMachineSummary builds machineID's row from its cycle counts; the OEE parts are set by
// WithOEE.
func NewMachineSummary(machineID int64, cycles MachineCycleSummary) MachineSummary {
	return MachineSummary{MachineID: machineID, CycleCount: cycles.CycleCount,
		CompletedCount: cycles.CompletedCount, UnfinishedCount: cycles.UnfinishedCount}
}

// WithOEE sets the day's OEE parts (one machine's row of GET /productivity/oee).
func (summary MachineSummary) WithOEE(row OEERow) MachineSummary {
	availability := row.Availability
	summary.Availability, summary.Performance, summary.Quality, summary.OEE = &availability, row.Performance, row.Quality, row.OEE
	return summary
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

// HomeTopAlarms is how many alarm codes the Home dashboard lists (the most raised first).
const HomeTopAlarms = 10

// HomeDashboardRequest is the query of GET /dashboard/home: the day (YYYY-MM-DD, empty = today)
// and, for the superadmin, the customer whose machines it covers (nil = every machine the user
// may read).
type HomeDashboardRequest struct {
	Date       string
	CustomerID *int64
}

// Validate checks the day: a valid date, not after today.
func (homeRequest *HomeDashboardRequest) Validate() error {
	if homeRequest.Date == "" {
		return nil
	}
	if !validDate(homeRequest.Date) {
		return response.NewError(http.StatusBadRequest, "Please enter valid date")
	}
	if day, _ := homeRequest.Range(time.Now()); day.After(time.Now()) {
		return response.NewError(http.StatusBadRequest, "Date must not be after today")
	}
	return nil
}

// Range returns the start of the day and the start of the next day (server time zone).
func (homeRequest *HomeDashboardRequest) Range(now time.Time) (day, nextDay time.Time) {
	listRequest := MachineCycleListRequest{From: homeRequest.Date, To: homeRequest.Date}
	return listRequest.Range(now)
}

// HomeDashboardResponse is the response of GET /dashboard/home: the machines the user may read,
// together. Fleet is now (whatever the day); the rest covers the day up to ServerTime. OEE is
// nil before the day has any time; Plan is nil when no machine has a job scheduled that day.
type HomeDashboardResponse struct {
	Date       string        `json:"date"`
	ServerTime time.Time     `json:"serverTime"`
	Fleet      HomeFleet     `json:"fleet"`
	OEE        *OEERow       `json:"oee"`
	Cycles     HomeCycles    `json:"cycles"`
	Plan       *HomePlan     `json:"plan"`
	Time       HomeTime      `json:"time"`
	Hourly     []HomeHour    `json:"hourly"`
	Machines   []HomeMachine `json:"machines"`
	Alarms     []HomeAlarm   `json:"alarms"`
}

// HomeFleet counts the active machines now: online / offline, never collected (AwaitingData),
// and how many are in each mode; Inactive counts the machines switched off in the set-up (not
// collected, left out of everything else).
type HomeFleet struct {
	Total        int             `json:"total"`
	Online       int             `json:"online"`
	Offline      int             `json:"offline"`
	AwaitingData int             `json:"awaitingData"`
	Inactive     int             `json:"inactive"`
	Modes        []HomeModeCount `json:"modes"`
}

// Add counts an active machine in mode modeID now (nil = never collected).
func (fleet *HomeFleet) Add(machine model.Machine, modeID *int) {
	fleet.Total++
	if machine.OnlineStatus {
		fleet.Online++
	} else {
		fleet.Offline++
	}
	if modeID == nil {
		fleet.AwaitingData++
	}
}

// HomeModeCount is how many machines are in mode ModeID now.
type HomeModeCount struct {
	ModeID int `json:"modeId"`
	Count  int `json:"count"`
}

// HomeCycles adds up the cycles the machines started in the day and the parts they made.
type HomeCycles struct {
	CycleCount      int `json:"cycleCount"`
	CompletedCount  int `json:"completedCount"`
	UnfinishedCount int `json:"unfinishedCount"`
	RunningCount    int `json:"runningCount"`
	PartsMade       int `json:"partsMade"`
}

// NewHomeCycles returns the Home fields of summary.
func NewHomeCycles(summary MachineCycleSummary) HomeCycles {
	return HomeCycles{CycleCount: summary.CycleCount, CompletedCount: summary.CompletedCount,
		UnfinishedCount: summary.UnfinishedCount, RunningCount: summary.RunningCount, PartsMade: summary.PartsMade}
}

// HomePlan compares the quantity scheduled for the day with the parts made by the machines that
// have a schedule that day (machines without one are not counted against the plan).
type HomePlan struct {
	Quantity  int `json:"quantity"`
	PartsMade int `json:"partsMade"`
}

// HomeTime is where the machines' time went in the day: seconds per mode group (ModeGroup*)
// and per machine mode out of SpanSec, the day so far times the active machines. Time with no
// record is SpanSec minus the groups (or the modes, which add up to the same).
type HomeTime struct {
	SpanSec int64           `json:"spanSec"`
	Groups  []HomeTimeGroup `json:"groups"`
	Modes   []HomeTimeMode  `json:"modes"`
}

// HomeTimeMode is the seconds the machines spent in machine mode ModeID (only modes with time).
type HomeTimeMode struct {
	ModeID      int   `json:"modeId"`
	DurationSec int64 `json:"durationSec"`
}

// HomeTimeGroup is the seconds the machines spent in one mode group.
type HomeTimeGroup struct {
	Group       string `json:"group"`
	DurationSec int64  `json:"durationSec"`
}

// HomeHour is the parts the machines made in one hour (0-23) of the day and of the day before
// (cycles completed in that hour).
type HomeHour struct {
	Hour              int `json:"hour"`
	PartsMade         int `json:"partsMade"`
	PreviousPartsMade int `json:"previousPartsMade"`
}

// HomeMachine is one active machine: its mode now and since when (nil without an open period),
// and its day: OEE parts (nil as in MachineSummary), cycles, parts and scheduled quantity (nil
// when it has no job that day).
type HomeMachine struct {
	MachineID       int64      `json:"machineId"`
	MachineName     string     `json:"machineName"`
	CustomerName    string     `json:"customerName"`
	ModeID          *int       `json:"modeId"`
	OnlineStatus    bool       `json:"onlineStatus"`
	ModeSince       *time.Time `json:"modeSince"`
	Availability    *float64   `json:"availability"`
	Performance     *float64   `json:"performance"`
	Quality         *float64   `json:"quality"`
	OEE             *float64   `json:"oee"`
	CycleCount      int        `json:"cycleCount"`
	CompletedCount  int        `json:"completedCount"`
	UnfinishedCount int        `json:"unfinishedCount"`
	PartsMade       int        `json:"partsMade"`
	PlannedQuantity *int       `json:"plannedQuantity"`
}

// NewHomeMachine builds machine's row from its current mode (modeID, since; nil when unknown),
// its cycles of the day and its planned quantity; the OEE parts are set by WithOEE.
func NewHomeMachine(machine model.Machine, modeID *int, since *time.Time, cycles MachineCycleSummary, planned *int) HomeMachine {
	return HomeMachine{
		MachineID: machine.ID, MachineName: machine.Name, CustomerName: machine.CustomerName,
		ModeID: modeID, OnlineStatus: machine.OnlineStatus, ModeSince: since,
		CycleCount: cycles.CycleCount, CompletedCount: cycles.CompletedCount,
		UnfinishedCount: cycles.UnfinishedCount, PartsMade: cycles.PartsMade, PlannedQuantity: planned,
	}
}

// WithOEE returns homeMachine with the OEE parts of row.
func (homeMachine HomeMachine) WithOEE(row OEERow) HomeMachine {
	availability := row.Availability
	homeMachine.Availability, homeMachine.Performance, homeMachine.Quality, homeMachine.OEE = &availability, row.Performance, row.Quality, row.OEE
	return homeMachine
}

// HomeAlarm is one alarm code raised in the day: how often, on how many machines, its message
// and when it was last raised.
type HomeAlarm struct {
	Code         int       `json:"code"`
	Message      string    `json:"message"`
	Count        int       `json:"count"`
	MachineCount int       `json:"machineCount"`
	LastAt       time.Time `json:"lastAt"`
}

// NewHomeAlarm returns the Home fields of alarm.
func NewHomeAlarm(alarm model.AlarmCount) HomeAlarm {
	return HomeAlarm{Code: alarm.Code, Message: alarm.Message, Count: alarm.Count,
		MachineCount: alarm.MachineCount, LastAt: alarm.LastAt}
}
