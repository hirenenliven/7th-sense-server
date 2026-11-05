package service

import (
	"context"
	"maps"
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

// DashboardService serves the charts of the Home, Productivity and Machine Health pages from what
// the CNC collector recorded: mode periods (machine_mode_log) and cycles (machine_cycles). Users
// read the machines they may access (MachineService decides which). Days are calendar days in
// the server's time zone, like GET /machine-cycles.
type DashboardService struct {
	machineService     *MachineService
	telemetryService   *MachineTelemetryService
	cycleService       *MachineCycleService
	shiftRepository    *repository.ShiftRepository
	scheduleRepository *repository.ProductionScheduleRepository
}

// NewDashboardService creates a DashboardService.
func NewDashboardService(machineService *MachineService, telemetryService *MachineTelemetryService,
	cycleService *MachineCycleService, shiftRepository *repository.ShiftRepository,
	scheduleRepository *repository.ProductionScheduleRepository) *DashboardService {
	return &DashboardService{machineService: machineService, telemetryService: telemetryService,
		cycleService: cycleService, shiftRepository: shiftRepository, scheduleRepository: scheduleRepository}
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

// GetModeTotals returns the time the machine spent in each mode over the days From … To (up to
// now), or over every recorded period when AllTime.
func (dashboardService *DashboardService) GetModeTotals(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) (dto.ModeTotalsResponse, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, dashboardRequest.MachineIDs[:1])
	if err != nil {
		return dto.ModeTotalsResponse{}, err
	}
	from, to := time.Time{}, now
	if !dashboardRequest.AllTime {
		from, to = dashboardRequest.Range(now)
		to = minTime(to, now)
	}
	totals, err := dashboardService.telemetryService.GetModeDurations(ctx, machines[0].ID, from, to)
	if err != nil {
		return dto.ModeTotalsResponse{}, err
	}

	// A range starts at its From day; all time starts at the first recorded period
	var first *time.Time
	if !dashboardRequest.AllTime {
		first = &from
	}
	for index := range totals {
		if first == nil || totals[index].FirstAt.Before(*first) {
			first = &totals[index].FirstAt
		}
	}
	return dto.ModeTotalsResponse{From: first, To: to, ServerTime: now, Totals: totals}, nil
}

// GetOEE returns one OEE row per day with planned time, over all the machines together
// (weighted by time and cycles, see oeeSum): availability = run time / planned time, performance = standard cycle time × parts counted
// (ProgramPartCount) / run time of the completed cycles that have a standard time, quality =
// completed / closed cycles. Planned time
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

// GetMachineSummaries returns, for every machine actor may read, its OEE parts and cycle counts
// of dashboardRequest's To day (default today; the whole day, no shifts or program): the
// Monitoring machine list. Same rules as GetOEE and GetCycleCounts for one machine.
func (dashboardService *DashboardService) GetMachineSummaries(ctx context.Context, actor model.AuthUser, dashboardRequest dto.DashboardRequest) ([]dto.MachineSummary, error) {
	now := time.Now()
	machines, err := dashboardService.machineService.GetAll(ctx, actor, nil)
	if err != nil {
		return nil, err
	}
	dashboardRequest.From = dashboardRequest.To
	day, nextDay := dashboardRequest.Range(now)
	windows := dayWindows(day, nil, false, now)

	summaries := make([]dto.MachineSummary, 0, len(machines))
	for _, machine := range machines {
		cycles, err := dashboardService.cycleService.cycleRepository.FindAll(ctx, machine.ID, day, nextDay)
		if err != nil {
			return nil, err
		}
		dayCycles := cyclesIn(cycles, windows, "")
		summary := dto.NewMachineSummary(machine.ID, summarizeCycles(dayCycles))
		if day.Before(now) {
			periods, err := dashboardService.periodsBetween(ctx, machine.ID, day, nextDay)
			if err != nil {
				return nil, err
			}
			if dayOEE, ok := machineDayOEE(periods, dayCycles, windows, "", now); ok {
				sum := oeeSum{}
				sum.add(dayOEE)
				summary = summary.WithOEE(sum.row(day.Format(shiftDateLayout)))
			}
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

// GetHomeDashboard returns the Home page: the active machines actor may read (of
// homeRequest's customer for the superadmin) together. The fleet is counted now; OEE, cycles,
// parts vs plan, time per mode group, parts per hour and alarms cover the day up to now, with
// the same rules as GetOEE / GetMachineSummaries (whole day, no shifts), machines weighted by
// their time (oeeSum).
func (dashboardService *DashboardService) GetHomeDashboard(ctx context.Context, actor model.AuthUser, homeRequest dto.HomeDashboardRequest) (dto.HomeDashboardResponse, error) {
	now := time.Now()
	machines, err := dashboardService.machineService.GetAll(ctx, actor, homeRequest.CustomerID)
	if err != nil {
		return dto.HomeDashboardResponse{}, err
	}
	day, nextDay := homeRequest.Range(now)
	home := dto.HomeDashboardResponse{Date: day.Format(shiftDateLayout), ServerTime: now,
		Hourly: make([]dto.HomeHour, 24), Machines: []dto.HomeMachine{}, Alarms: []dto.HomeAlarm{}}
	for hour := range home.Hourly {
		home.Hourly[hour].Hour = hour
	}

	active, machineIDs := []model.Machine{}, []int64{}
	for _, machine := range machines {
		if machine.Status != dto.MachineStatusActive {
			home.Fleet.Inactive++
			continue
		}
		active = append(active, machine)
		machineIDs = append(machineIDs, machine.ID)
	}
	openPeriods, planned, err := dashboardService.homeLookups(ctx, machineIDs, home.Date)
	if err != nil {
		return dto.HomeDashboardResponse{}, err
	}
	alarms, err := dashboardService.telemetryService.GetAlarmCounts(ctx, machineIDs, day, nextDay, dto.HomeTopAlarms)
	if err != nil {
		return dto.HomeDashboardResponse{}, err
	}
	for _, alarm := range alarms {
		home.Alarms = append(home.Alarms, dto.NewHomeAlarm(alarm))
	}

	windows := dayWindows(day, nil, false, now)
	fleetOEE, fleetCycles := oeeSum{}, dto.MachineCycleSummary{}
	groupSeconds := map[string]float64{}
	modeCounts := map[int]int{}
	for _, machine := range active {
		modeID, since := currentMode(machine, openPeriods)
		home.Fleet.Add(machine, modeID)
		if modeID != nil {
			modeCounts[*modeID]++
		}

		// The day before is read too, for the parts per hour it is compared with
		cycles, err := dashboardService.cycleService.cycleRepository.FindAll(ctx, machine.ID, day.AddDate(0, 0, -1), nextDay)
		if err != nil {
			return dto.HomeDashboardResponse{}, err
		}
		addHourlyParts(home.Hourly, cycles, day)
		dayCycles := cyclesIn(cycles, windows, "")
		machineCycles := summarizeCycles(dayCycles)
		addCycleSummary(&fleetCycles, machineCycles)
		plannedQuantity := planned[machine.ID]
		if plannedQuantity != nil {
			if home.Plan == nil {
				home.Plan = &dto.HomePlan{}
			}
			home.Plan.Quantity += *plannedQuantity
			home.Plan.PartsMade += machineCycles.PartsMade
		}
		homeMachine := dto.NewHomeMachine(machine, modeID, since, machineCycles, plannedQuantity)

		if day.Before(now) {
			periods, err := dashboardService.periodsBetween(ctx, machine.ID, day, nextDay)
			if err != nil {
				return dto.HomeDashboardResponse{}, err
			}
			for _, window := range windows {
				home.Time.SpanSec += int64(window.To.Sub(window.From).Seconds())
				for _, period := range periods {
					groupSeconds[modeGroup(period.ModeID)] += overlapSeconds(period, window, now)
				}
			}
			if dayOEE, ok := machineDayOEE(periods, dayCycles, windows, "", now); ok {
				machineOEE := oeeSum{}
				machineOEE.add(dayOEE)
				homeMachine = homeMachine.WithOEE(machineOEE.row(home.Date))
				fleetOEE.add(dayOEE)
			}
		}
		home.Machines = append(home.Machines, homeMachine)
	}

	home.Fleet.Modes = homeModeCounts(modeCounts)
	home.Cycles = dto.NewHomeCycles(fleetCycles)
	home.Time.Groups = homeTimeGroups(groupSeconds)
	if fleetOEE.plannedSec > 0 {
		row := fleetOEE.row(home.Date)
		home.OEE = &row
	}
	return home, nil
}

// homeLookups returns the open mode period and the quantity scheduled on day of each of
// machineIDs (machines without one are not in the maps).
func (dashboardService *DashboardService) homeLookups(ctx context.Context, machineIDs []int64, day string) (map[int64]model.OpenModePeriod, map[int64]*int, error) {
	periods, err := dashboardService.telemetryService.GetOpenModePeriods(ctx, machineIDs)
	if err != nil {
		return nil, nil, err
	}
	quantities, err := dashboardService.scheduleRepository.SumMachineQuantities(ctx, machineIDs, day)
	if err != nil {
		return nil, nil, err
	}
	openPeriods := make(map[int64]model.OpenModePeriod, len(periods))
	for _, period := range periods {
		openPeriods[period.MachineID] = period
	}
	planned := make(map[int64]*int, len(quantities))
	for _, quantity := range quantities {
		planned[quantity.MachineID] = &quantity.Quantity
	}
	return openPeriods, planned, nil
}

// currentMode returns the mode machine is in now and since when: its open mode period, else
// machines.mode_id with no start; nil for a machine never collected.
func currentMode(machine model.Machine, openPeriods map[int64]model.OpenModePeriod) (*int, *time.Time) {
	if period, found := openPeriods[machine.ID]; found {
		return &period.ModeID, &period.StartedAt
	}
	if machine.ModeID == nil {
		return nil, nil
	}
	modeID := int(*machine.ModeID)
	return &modeID, nil
}

// addHourlyParts adds the parts of the cycles completed on day, and on the day before, to the
// hour (0-23) they completed in.
func addHourlyParts(hours []dto.HomeHour, cycles []model.MachineCycle, day time.Time) {
	previousDay, nextDay := day.AddDate(0, 0, -1), day.AddDate(0, 0, 1)
	for _, cycle := range cycles {
		if cycle.Status != model.CycleCompleted || cycle.EndedAt == nil || cycle.PartsMade == 0 {
			continue
		}
		endedAt := cycle.EndedAt.In(day.Location())
		switch {
		case !endedAt.Before(day) && endedAt.Before(nextDay):
			hours[endedAt.Hour()].PartsMade += cycle.PartsMade
		case !endedAt.Before(previousDay) && endedAt.Before(day):
			hours[endedAt.Hour()].PreviousPartsMade += cycle.PartsMade
		}
	}
}

// addCycleSummary adds machine's cycle counts and parts to fleet.
func addCycleSummary(fleet *dto.MachineCycleSummary, machine dto.MachineCycleSummary) {
	fleet.CycleCount += machine.CycleCount
	fleet.CompletedCount += machine.CompletedCount
	fleet.UnfinishedCount += machine.UnfinishedCount
	fleet.RunningCount += machine.RunningCount
	fleet.PartsMade += machine.PartsMade
}

// homeModeCounts returns how many machines are in each mode, by mode id.
func homeModeCounts(counts map[int]int) []dto.HomeModeCount {
	modes := make([]dto.HomeModeCount, 0, len(counts))
	for _, modeID := range slices.Sorted(maps.Keys(counts)) {
		modes = append(modes, dto.HomeModeCount{ModeID: modeID, Count: counts[modeID]})
	}
	return modes
}

// homeTimeGroups returns the seconds of every mode group, in the order the chart stacks them.
func homeTimeGroups(seconds map[string]float64) []dto.HomeTimeGroup {
	order := []string{dto.ModeGroupProduction, dto.ModeGroupHold, dto.ModeGroupIdle, dto.ModeGroupDowntime, dto.ModeGroupOffline}
	groups := make([]dto.HomeTimeGroup, 0, len(order))
	for _, group := range order {
		groups = append(groups, dto.HomeTimeGroup{Group: group, DurationSec: int64(seconds[group])})
	}
	return groups
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

// machineOEE is one machine's day as the seconds and counts its OEE parts come from, so several
// machines add up weighted by their time and cycles (oeeSum), not by averaging percentages.
type machineOEE struct {
	plannedSec, runSec         float64 // availability = run / planned
	standardSec, productionSec float64 // performance = standard / production (completed cycles with a standard time)
	completed, closed          int     // quality = completed / closed cycles
}

// machineDayOEE works out one machine's day; ok is false when the day has no planned time.
func machineDayOEE(periods []model.MachineModePeriod, cycles []model.MachineCycle, windows []timeWindow, program string, now time.Time) (machineOEE, bool) {
	result := machineOEE{}
	for _, window := range windows {
		result.plannedSec += window.To.Sub(window.From).Seconds()
		for _, period := range periods {
			overlap := overlapSeconds(period, window, now)
			switch {
			case slices.Contains(plannedStopModes, period.ModeID):
				result.plannedSec -= overlap
			case slices.Contains(runModes, period.ModeID) && (program == "" || period.ProgramName == program):
				result.runSec += overlap
			}
		}
	}
	if result.plannedSec <= 0 {
		return machineOEE{}, false
	}

	for _, cycle := range cycles {
		if cycle.Status == model.CycleRunning {
			continue
		}
		result.closed++
		if cycle.Status != model.CycleCompleted {
			continue
		}
		result.completed++
		// Parts come only from the controller's counter (M30). One part can take several
		// cycles, so a cycle without a part adds its run time but no standard time.
		if cycle.StandardCycleSec != nil && *cycle.StandardCycleSec > 0 {
			result.standardSec += float64(*cycle.StandardCycleSec * cycle.PartsMade)
			result.productionSec += float64(cycle.ProductionSec)
		}
	}
	return result, true
}

// oeeSum adds up the machines' days (seconds and cycles) so a machine counts by how long it was
// planned and ran and by how many cycles it closed: a machine planned 10 minutes weighs less
// than one planned 20 hours. For one machine it is that machine's OEE.
type oeeSum struct {
	machineOEE
}

func (sum *oeeSum) add(value machineOEE) {
	sum.plannedSec += value.plannedSec
	sum.runSec += value.runSec
	sum.standardSec += value.standardSec
	sum.productionSec += value.productionSec
	sum.completed += value.completed
	sum.closed += value.closed
}

func (sum *oeeSum) row(date string) dto.OEERow {
	availability := math.Min(100, sum.runSec/sum.plannedSec*100)
	row := dto.OEERow{Date: date, Availability: round1(availability)}
	if sum.closed > 0 {
		quality := float64(sum.completed) / float64(sum.closed) * 100
		// Without standard cycle times the run counts at full speed
		performance := 100.0
		if sum.productionSec > 0 {
			performance = math.Min(100, sum.standardSec/sum.productionSec*100)
		}
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
