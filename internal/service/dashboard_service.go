package service

import (
	"context"
	"math"
	"slices"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
)

// dayDuration is one calendar day of the dashboard charts.
const dayDuration = 24 * time.Hour

// Machine modes that count as cutting time (OEE run time) and as planned stops (left out of the
// planned time).
var (
	runModes         = []int{model.ModeInProduction, model.ModeRapid}
	plannedStopModes = []int{7, 13} // Planned Downtime, Leave Planned Downtime (user modes)
)

// timeWindow is a span of time the charts count, From included, To excluded.
type timeWindow struct {
	From time.Time
	To   time.Time
}

// DashboardService serves the charts of the Productivity and Machine Health pages from what
// the CNC collector recorded: mode periods (machine_mode_log) and cycles (machine_cycles). Users
// read the machines they may access (MachineService decides which). Days are calendar days in
// the server's time zone, like GET /machine-cycles.
type DashboardService struct {
	machineService   *MachineService
	telemetryService *MachineTelemetryService
	cycleService     *MachineCycleService
	shiftRepository  *repository.ShiftRepository
}

// NewDashboardService creates a DashboardService.
func NewDashboardService(machineService *MachineService, telemetryService *MachineTelemetryService,
	cycleService *MachineCycleService, shiftRepository *repository.ShiftRepository) *DashboardService {
	return &DashboardService{machineService: machineService, telemetryService: telemetryService,
		cycleService: cycleService, shiftRepository: shiftRepository}
}

// GetModeTimeline returns the mode periods of dashboardRequest's machine on its To day, or in its
// last Minutes up to now, clipped to that window (the Productivity Gantt chart), with the program
// and operation mode it has now.
func (dashboardService *DashboardService) GetModeTimeline(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) (dto.ModeTimelineResponse, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, dashboardRequest.MachineIDs[:1])
	if err != nil {
		return dto.ModeTimelineResponse{}, err
	}
	machine := machines[0]
	dashboardRequest.From = dashboardRequest.To
	from, to := dashboardRequest.Range(now)
	if dashboardRequest.Minutes > 0 {
		from, to = now.Add(-time.Duration(dashboardRequest.Minutes)*time.Minute), now
	}
	periods, err := dashboardService.periodsBetween(ctx, machine.ID, from, to)
	if err != nil {
		return dto.ModeTimelineResponse{}, err
	}
	opMode, err := dashboardService.telemetryService.GetOpModeName(ctx, machine.ID)
	if err != nil {
		return dto.ModeTimelineResponse{}, err
	}
	programName := ""
	if len(periods) > 0 {
		programName = periods[len(periods)-1].ProgramName
	}

	clipped := []dto.TimelinePeriod{}
	for _, period := range periods {
		start, end := periodSpan(period, now)
		start, end = maxTime(start, from), minTime(end, to)
		if !end.After(start) {
			continue
		}
		period.StartedAt = start
		if period.EndedAt != nil || end.Before(now) {
			period.EndedAt = &end
		}
		clipped = append(clipped, dto.NewTimelinePeriod(period))
	}

	return dto.ModeTimelineResponse{
		ProgramName: programName, OpMode: opMode, From: from, To: to, ServerTime: now, Periods: clipped,
	}, nil
}

// GetOEE returns one OEE row per day with planned time, averaged over the machines:
// availability = run time / planned time, performance = standard cycle time / run time of the
// completed cycles that have a standard time, quality = completed / closed cycles. Planned time
// is the day (or the shifts picked); "with schedule" counts only the customer's shifts that
// work that weekday, and planned stops are left out.
func (dashboardService *DashboardService) GetOEE(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) (dto.OEEResponse, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, dashboardRequest.MachineIDs)
	if err != nil {
		return dto.OEEResponse{}, err
	}
	shifts, err := dashboardService.shiftsFor(ctx, machines[0], dashboardRequest.ShiftIDs, dashboardRequest.WithSchedule)
	if err != nil {
		return dto.OEEResponse{}, err
	}
	from, to := dashboardRequest.Range(now)
	programs := map[string]bool{}
	sums := map[string]*oeeSum{}

	for _, machine := range machines {
		// A shift past midnight runs into the day after To
		periods, err := dashboardService.periodsBetween(ctx, machine.ID, from, to.Add(dayDuration))
		if err != nil {
			return dto.OEEResponse{}, err
		}
		cycles, err := dashboardService.cycleService.cycleRepository.FindAll(ctx, machine.ID, from, to.Add(dayDuration))
		if err != nil {
			return dto.OEEResponse{}, err
		}
		for day := from; day.Before(to) && day.Before(now); day = day.AddDate(0, 0, 1) {
			windows := dayWindows(day, shifts, dashboardRequest.WithSchedule, now)
			for _, cycle := range cyclesIn(cycles, windows, "") {
				if cycle.ProgramName != "" {
					programs[cycle.ProgramName] = true
				}
			}
			dayOEE, ok := machineDayOEE(periods, cyclesIn(cycles, windows, dashboardRequest.Program), windows, dashboardRequest.Program, now)
			if !ok {
				continue
			}
			date := day.Format(shiftDateLayout)
			if sums[date] == nil {
				sums[date] = &oeeSum{}
			}
			sums[date].add(dayOEE)
		}
	}

	rows := []dto.OEERow{}
	for day := from; day.Before(to); day = day.AddDate(0, 0, 1) {
		date := day.Format(shiftDateLayout)
		if sum := sums[date]; sum != nil {
			rows = append(rows, sum.row(date))
		}
	}
	return dto.OEEResponse{Rows: rows, Programs: sortedKeys(programs)}, nil
}

// GetCycleCounts returns one row per day up to today: the cycles the machines started in the
// day (or the shifts picked), and how many of them completed or stayed unfinished.
func (dashboardService *DashboardService) GetCycleCounts(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) ([]dto.CycleCountRow, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, dashboardRequest.MachineIDs)
	if err != nil {
		return nil, err
	}
	shifts, err := dashboardService.shiftsFor(ctx, machines[0], dashboardRequest.ShiftIDs, false)
	if err != nil {
		return nil, err
	}
	from, to := dashboardRequest.Range(now)

	rows := []dto.CycleCountRow{}
	for day := from; day.Before(to) && day.Before(now); day = day.AddDate(0, 0, 1) {
		rows = append(rows, dto.CycleCountRow{Date: day.Format(shiftDateLayout)})
	}
	for _, machine := range machines {
		cycles, err := dashboardService.cycleService.cycleRepository.FindAll(ctx, machine.ID, from, to.Add(dayDuration))
		if err != nil {
			return nil, err
		}
		for index := range rows {
			day := from.AddDate(0, 0, index)
			summary := summarizeCycles(cyclesIn(cycles, dayWindows(day, shifts, false, now), ""))
			rows[index].Produced += summary.CycleCount
			rows[index].Finished += summary.CompletedCount
			rows[index].Unfinished += summary.UnfinishedCount
		}
	}
	return rows, nil
}

// GetFinishCycles returns the completed cycles of dashboardRequest's machine on its To day (in
// the shifts and program picked), the totals of every cycle started then and the programs run.
func (dashboardService *DashboardService) GetFinishCycles(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) (dto.FinishCyclesResponse, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, dashboardRequest.MachineIDs[:1])
	if err != nil {
		return dto.FinishCyclesResponse{}, err
	}
	machine := machines[0]
	shifts, err := dashboardService.shiftsFor(ctx, machine, dashboardRequest.ShiftIDs, false)
	if err != nil {
		return dto.FinishCyclesResponse{}, err
	}
	dashboardRequest.From = dashboardRequest.To
	day, _ := dashboardRequest.Range(now)
	cycles, err := dashboardService.cycleService.cycleRepository.FindAll(ctx, machine.ID, day, day.Add(2*dayDuration))
	if err != nil {
		return dto.FinishCyclesResponse{}, err
	}

	windows := dayWindows(day, shifts, false, now)
	programs := map[string]bool{}
	for _, cycle := range cyclesIn(cycles, windows, "") {
		if cycle.ProgramName != "" {
			programs[cycle.ProgramName] = true
		}
	}
	picked := cyclesIn(cycles, windows, dashboardRequest.Program)
	completed := []dto.FinishedCycle{}
	var lastFinishedAt *time.Time
	for _, cycle := range picked {
		if cycle.Status != model.CycleCompleted {
			continue
		}
		completed = append(completed, dto.NewFinishedCycle(cycle))
		if lastFinishedAt == nil || cycle.EndedAt.After(*lastFinishedAt) {
			lastFinishedAt = cycle.EndedAt
		}
	}
	return dto.FinishCyclesResponse{
		Cycles: completed, Summary: dto.NewFinishCyclesSummary(summarizeCycles(picked)),
		LastFinishedAt: lastFinishedAt, Programs: sortedKeys(programs),
	}, nil
}

// GetMachineTimeByDay returns, per day up to today, the hours dashboardRequest's machine spent
// in each mode group; time without a mode period counts as offline.
func (dashboardService *DashboardService) GetMachineTimeByDay(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) ([]dto.MachineTimeRow, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, dashboardRequest.MachineIDs[:1])
	if err != nil {
		return nil, err
	}
	from, to := dashboardRequest.Range(now)
	periods, err := dashboardService.periodsBetween(ctx, machines[0].ID, from, to)
	if err != nil {
		return nil, err
	}

	rows := []dto.MachineTimeRow{}
	for day := from; day.Before(to) && day.Before(now); day = day.AddDate(0, 0, 1) {
		window := timeWindow{From: day, To: minTime(day.AddDate(0, 0, 1), now)}
		seconds := map[string]float64{}
		recorded := 0.0
		for _, period := range periods {
			overlap := overlapSeconds(period, window, now)
			seconds[modeGroup(period.ModeID)] += overlap
			recorded += overlap
		}
		seconds[dto.ModeGroupOffline] += math.Max(0, window.To.Sub(window.From).Seconds()-recorded)
		rows = append(rows, dto.MachineTimeRow{
			Date:       day.Format(shiftDateLayout),
			Offline:    hours(seconds[dto.ModeGroupOffline]),
			Production: hours(seconds[dto.ModeGroupProduction]),
			Idle:       hours(seconds[dto.ModeGroupIdle]),
			Hold:       hours(seconds[dto.ModeGroupHold]),
			Downtime:   hours(seconds[dto.ModeGroupDowntime]),
		})
	}
	return rows, nil
}

// machinesFor returns the machines of ids, or "Machine not found" when actor may not read one.
func (dashboardService *DashboardService) machinesFor(ctx context.Context, actor model.AuthUser, ids []int64) ([]model.Machine, error) {
	machines := make([]model.Machine, 0, len(ids))
	for _, id := range ids {
		machine, err := dashboardService.machineService.GetByID(ctx, actor, nil, id)
		if err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

// shiftsFor returns the shifts of machine's customer that the charts count: the ones picked
// (ids of other customers are ignored), every active shift when none is picked and
// withSchedule, else none (the whole day).
func (dashboardService *DashboardService) shiftsFor(ctx context.Context, machine model.Machine, shiftIDs []int64, withSchedule bool) ([]model.Shift, error) {
	if machine.UserID == nil || (len(shiftIDs) == 0 && !withSchedule) {
		return nil, nil
	}
	shifts, err := dashboardService.shiftRepository.FindAll(ctx, machine.UserID)
	if err != nil {
		return nil, err
	}
	picked := []model.Shift{}
	for _, shift := range shifts {
		if len(shiftIDs) > 0 && !slices.Contains(shiftIDs, shift.ID) {
			continue
		}
		if len(shiftIDs) == 0 && shift.Status != "active" {
			continue
		}
		picked = append(picked, shift)
	}
	return picked, nil
}

// periodsBetween returns the machine's mode periods that overlap from … to.
func (dashboardService *DashboardService) periodsBetween(ctx context.Context, machineID int64, from, to time.Time) ([]model.MachineModePeriod, error) {
	periods, err := dashboardService.telemetryService.GetModePeriods(ctx, machineID, from)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(periods, func(period model.MachineModePeriod) bool {
		return !period.StartedAt.Before(to)
	}), nil
}

// dayWindows returns the time counted on day, up to now: the whole day without shifts, else
// each shift from its start on day (a shift past midnight ends the next day). withSchedule
// skips shifts that do not work on day's weekday; a customer without shifts counts the day.
func dayWindows(day time.Time, shifts []model.Shift, withSchedule bool, now time.Time) []timeWindow {
	if len(shifts) == 0 {
		return clipWindows([]timeWindow{{From: day, To: day.AddDate(0, 0, 1)}}, now)
	}
	weekday := int32(day.Weekday())
	if weekday == 0 {
		weekday = 7 // ISO: Sunday is 7
	}
	windows := []timeWindow{}
	for _, shift := range shifts {
		if withSchedule && len(shift.WorkingDays) > 0 && !slices.Contains(shift.WorkingDays, weekday) {
			continue
		}
		start, startErr := time.Parse("15:04", shift.StartTime)
		end, endErr := time.Parse("15:04", shift.EndTime)
		if startErr != nil || endErr != nil {
			continue
		}
		window := timeWindow{
			From: day.Add(time.Duration(start.Hour())*time.Hour + time.Duration(start.Minute())*time.Minute),
			To:   day.Add(time.Duration(end.Hour())*time.Hour + time.Duration(end.Minute())*time.Minute),
		}
		if !window.To.After(window.From) {
			window.To = window.To.Add(dayDuration)
		}
		windows = append(windows, window)
	}
	return clipWindows(windows, now)
}

// clipWindows ends windows at now and drops the ones that have not started.
func clipWindows(windows []timeWindow, now time.Time) []timeWindow {
	clipped := []timeWindow{}
	for _, window := range windows {
		window.To = minTime(window.To, now)
		if window.To.After(window.From) {
			clipped = append(clipped, window)
		}
	}
	return clipped
}

// cyclesIn returns the cycles that started inside one of windows and ran program ("" = any).
func cyclesIn(cycles []model.MachineCycle, windows []timeWindow, program string) []model.MachineCycle {
	inside := []model.MachineCycle{}
	for _, cycle := range cycles {
		if program != "" && cycle.ProgramName != program {
			continue
		}
		for _, window := range windows {
			if !cycle.StartedAt.Before(window.From) && cycle.StartedAt.Before(window.To) {
				inside = append(inside, cycle)
				break
			}
		}
	}
	return inside
}

// machineOEE is one machine's OEE parts for one day (percent); nil = no closed cycle.
type machineOEE struct {
	availability float64
	performance  *float64
	quality      *float64
}

// machineDayOEE works out one machine's day; ok is false when the day has no planned time.
func machineDayOEE(periods []model.MachineModePeriod, cycles []model.MachineCycle, windows []timeWindow, program string, now time.Time) (machineOEE, bool) {
	planned, run := 0.0, 0.0
	for _, window := range windows {
		planned += window.To.Sub(window.From).Seconds()
		for _, period := range periods {
			overlap := overlapSeconds(period, window, now)
			switch {
			case slices.Contains(plannedStopModes, period.ModeID):
				planned -= overlap
			case slices.Contains(runModes, period.ModeID) && (program == "" || period.ProgramName == program):
				run += overlap
			}
		}
	}
	if planned <= 0 {
		return machineOEE{}, false
	}
	result := machineOEE{availability: math.Min(100, run/planned*100)}

	completed, closed, standardSec, productionSec := 0, 0, 0.0, 0.0
	for _, cycle := range cycles {
		if cycle.Status == model.CycleRunning {
			continue
		}
		closed++
		if cycle.Status != model.CycleCompleted {
			continue
		}
		completed++
		if cycle.StandardCycleSec != nil && *cycle.StandardCycleSec > 0 && cycle.ProductionSec > 0 {
			standardSec += float64(*cycle.StandardCycleSec * max(cycle.PartsMade, 1))
			productionSec += float64(cycle.ProductionSec)
		}
	}
	if closed > 0 {
		quality := float64(completed) / float64(closed) * 100
		result.quality = &quality
		// Without standard cycle times the run counts at full speed
		performance := 100.0
		if productionSec > 0 {
			performance = math.Min(100, standardSec/productionSec*100)
		}
		result.performance = &performance
	}
	return result, true
}

// oeeSum adds up the machines' OEE parts of one day to average them.
type oeeSum struct {
	machines, closedMachines           int
	availability, performance, quality float64
}

func (sum *oeeSum) add(value machineOEE) {
	sum.machines++
	sum.availability += value.availability
	if value.performance != nil && value.quality != nil {
		sum.closedMachines++
		sum.performance += *value.performance
		sum.quality += *value.quality
	}
}

func (sum *oeeSum) row(date string) dto.OEERow {
	availability := sum.availability / float64(sum.machines)
	row := dto.OEERow{Date: date, Availability: round1(availability)}
	if sum.closedMachines > 0 {
		performance := sum.performance / float64(sum.closedMachines)
		quality := sum.quality / float64(sum.closedMachines)
		oee := round1(availability * performance * quality / 10_000)
		performance, quality = round1(performance), round1(quality)
		row.Performance, row.Quality, row.OEE = &performance, &quality, &oee
	}
	return row
}

// modeGroup returns the chart group of machine mode modeID. Setup and loading are not cutting,
// so they count as idle; planned stops count as downtime.
func modeGroup(modeID int) string {
	switch modeID {
	case model.ModeInProduction, model.ModeRapid:
		return dto.ModeGroupProduction
	case model.ModeHold:
		return dto.ModeGroupHold
	case model.ModeSetup, 2, model.ModeIdle: // 2 = Loading / Unloading
		return dto.ModeGroupIdle
	case model.ModeDowntime, 7, 13:
		return dto.ModeGroupDowntime
	default:
		return dto.ModeGroupOffline
	}
}

// periodSpan returns when period started and ended (now while it is open).
func periodSpan(period model.MachineModePeriod, now time.Time) (time.Time, time.Time) {
	if period.EndedAt != nil {
		return period.StartedAt, *period.EndedAt
	}
	return period.StartedAt, now
}

// overlapSeconds returns how many seconds of period fall inside window.
func overlapSeconds(period model.MachineModePeriod, window timeWindow, now time.Time) float64 {
	start, end := periodSpan(period, now)
	start, end = maxTime(start, window.From), minTime(end, window.To)
	if !end.After(start) {
		return 0
	}
	return end.Sub(start).Seconds()
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func minTime(first, second time.Time) time.Time {
	if first.Before(second) {
		return first
	}
	return second
}

func maxTime(first, second time.Time) time.Time {
	if first.After(second) {
		return first
	}
	return second
}

func hours(seconds float64) float64 { return math.Round(seconds/36) / 100 }

func round1(value float64) float64 { return math.Round(value*10) / 10 }
