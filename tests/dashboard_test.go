package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// newDashboardRouter builds the real router with the /productivity and /machine-health chart
// handlers, acting as the superadmin test user.
func newDashboardRouter() http.Handler {
	machineService := newMachineService()
	telemetryService := service.NewMachineTelemetryService(testDBPool, repository.NewMachineTelemetryRepository(testDBPool))
	cycleService := service.NewMachineCycleService(repository.NewMachineCycleRepository(testDBPool), machineService)
	return routes.New(routes.Handlers{Dashboard: handler.NewDashboardHandler(service.NewDashboardService(
		machineService, telemetryService, cycleService, repository.NewShiftRepository(testDBPool),
		repository.NewProductionScheduleRepository(testDBPool)))}, asTestUser)
}

// TestDashboardCharts reads every chart of the machine with the most recent mode periods (read
// only: the data is what the CNC collector recorded).
func TestDashboardCharts(t *testing.T) {
	var machineID int64
	err := testDBPool.QueryRow(context.Background(), `SELECT l.machine_id FROM machine_mode_log l
		JOIN machines m ON m.id = l.machine_id AND m.is_deleted = false
		ORDER BY l.started_at DESC LIMIT 1`).Scan(&machineID)
	if err != nil {
		t.Skipf("no machine with mode periods: %v", err)
	}
	router := newDashboardRouter()
	today := time.Now().Format("2006-01-02")
	weekAgo := time.Now().AddDate(0, 0, -6).Format("2006-01-02")

	timeline := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/productivity/mode-timeline?machineId=%d", apiV1, machineID), nil)
	if !timeline.Success {
		t.Fatalf("mode timeline: %s", timeline.Message)
	}
	var day dto.ModeTimelineResponse
	mustDecode(t, timeline.Data, &day)
	for _, period := range day.Periods {
		if period.StartedAt.Before(day.From) || (period.EndedAt != nil && period.EndedAt.After(day.To)) {
			t.Errorf("period %d not clipped to the day: %v – %v", period.ID, period.StartedAt, period.EndedAt)
		}
	}

	machineTime := sendRequest(t, router, http.MethodGet,
		fmt.Sprintf("%s/machine-health/machine-time-by-day?machineId=%d&from=%s&to=%s", apiV1, machineID, weekAgo, today), nil)
	if !machineTime.Success {
		t.Fatalf("machine time by day: %s", machineTime.Message)
	}
	var rows []dto.MachineTimeRow
	mustDecode(t, machineTime.Data, &rows)
	if len(rows) != 7 {
		t.Fatalf("machine time by day: %d rows, want 7", len(rows))
	}
	t.Logf("machine %d: %d periods today, first day %+v, today %+v", machineID, len(day.Periods), rows[0], rows[6])
	// A past day adds up to 24 h (rounding to 0.01 h per group)
	first := rows[0]
	if total := first.Offline + first.Production + first.Idle + first.Hold + first.Downtime; total < 23.9 || total > 24.1 {
		t.Errorf("machine time of %s adds up to %.2f h, want 24", first.Date, total)
	}

	recent := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/productivity/mode-timeline?machineId=%d&minutes=30", apiV1, machineID), nil)
	if !recent.Success {
		t.Fatalf("mode timeline of 30 minutes: %s", recent.Message)
	}
	var lastMinutes dto.ModeTimelineResponse
	mustDecode(t, recent.Data, &lastMinutes)
	if span := lastMinutes.To.Sub(lastMinutes.From); span != 30*time.Minute {
		t.Errorf("mode timeline of 30 minutes spans %v", span)
	}
	t.Logf("last 30 minutes: %d periods, program %q, op mode %q", len(lastMinutes.Periods), lastMinutes.ProgramName, lastMinutes.OpMode)

	// Mode totals: today's add up to at most the time since midnight; all time covers at least today
	todayTotals := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/productivity/mode-totals?machineId=%d&date=%s", apiV1, machineID, today), nil)
	if !todayTotals.Success {
		t.Fatalf("mode totals of today: %s", todayTotals.Message)
	}
	var todaySum dto.ModeTotalsResponse
	mustDecode(t, todayTotals.Data, &todaySum)
	todaySeconds := int64(0)
	for _, total := range todaySum.Totals {
		todaySeconds += total.DurationSec
	}
	midnight, _ := time.ParseInLocation("2006-01-02", today, time.Local)
	if limit := int64(time.Since(midnight).Seconds()) + 1; todaySeconds > limit {
		t.Errorf("mode totals of today add up to %d s, more than the %d s since midnight", todaySeconds, limit)
	}

	allTotals := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/productivity/mode-totals?machineId=%d&all=true", apiV1, machineID), nil)
	if !allTotals.Success {
		t.Fatalf("mode totals of all time: %s", allTotals.Message)
	}
	var allSum dto.ModeTotalsResponse
	mustDecode(t, allTotals.Data, &allSum)
	allSeconds := int64(0)
	for _, total := range allSum.Totals {
		allSeconds += total.DurationSec
	}
	if allSeconds < todaySeconds {
		t.Errorf("mode totals of all time (%d s) are less than today's (%d s)", allSeconds, todaySeconds)
	}
	if allSum.From == nil {
		t.Errorf("mode totals of all time have no start")
	}
	t.Logf("mode totals: today %d s, all time %d s since %v", todaySeconds, allSeconds, allSum.From)

	if result := sendRequest(t, router, http.MethodGet, apiV1+"/productivity/mode-totals?all=true", nil); result.StatusCode != http.StatusBadRequest {
		t.Errorf("mode totals without a machine: status %d, want 400", result.StatusCode)
	}

	for _, path := range []string{
		fmt.Sprintf("/productivity/mode-totals?machineId=%d&from=%s&to=%s", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/oee?machineIds=%d&from=%s&to=%s", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/oee?machineIds=%d&from=%s&to=%s&schedule=with", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/cycle-counts?machineIds=%d&from=%s&to=%s", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/finish-cycles?machineId=%d&date=%s", machineID, today),
	} {
		if result := sendRequest(t, router, http.MethodGet, apiV1+path, nil); !result.Success {
			t.Errorf("GET %s: %s", path, result.Message)
		}
	}

	summaries := sendRequest(t, router, http.MethodGet, apiV1+"/productivity/machine-summaries", nil)
	if !summaries.Success {
		t.Fatalf("machine summaries: %s", summaries.Message)
	}
	var machineRows []dto.MachineSummary
	mustDecode(t, summaries.Data, &machineRows)
	for _, row := range machineRows {
		if row.CompletedCount+row.UnfinishedCount > row.CycleCount {
			t.Errorf("machine %d: %d completed + %d unfinished > %d cycles", row.MachineID, row.CompletedCount, row.UnfinishedCount, row.CycleCount)
		}
		if row.OEE != nil && (*row.OEE < 0 || *row.OEE > 100) {
			t.Errorf("machine %d: OEE %v out of range", row.MachineID, *row.OEE)
		}
	}
	if result := sendRequest(t, router, http.MethodGet, apiV1+"/productivity/machine-summaries?date=bad", nil); result.StatusCode != http.StatusBadRequest {
		t.Errorf("machine summaries with a bad date: status %d, want 400", result.StatusCode)
	}

	if result := sendRequest(t, router, http.MethodGet, apiV1+"/productivity/oee", nil); result.StatusCode != http.StatusBadRequest {
		t.Errorf("OEE without a machine: status %d, want 400", result.StatusCode)
	}
}

// TestHomeDashboard reads the Home dashboard of every machine (superadmin) and checks that its
// parts add up: fleet counts, cycles, time within the day so far, OEE parts in range.
func TestHomeDashboard(t *testing.T) {
	router := newDashboardRouter()
	result := sendRequest(t, router, http.MethodGet, apiV1+"/dashboard/home", nil)
	if !result.Success {
		t.Fatalf("home dashboard: %s", result.Message)
	}
	var home dto.HomeDashboardResponse
	mustDecode(t, result.Data, &home)

	fleet := home.Fleet
	if fleet.Online+fleet.Offline != fleet.Total {
		t.Errorf("fleet: %d online + %d offline != %d total", fleet.Online, fleet.Offline, fleet.Total)
	}
	inModes := 0
	for _, mode := range fleet.Modes {
		inModes += mode.Count
	}
	if inModes+fleet.AwaitingData != fleet.Total {
		t.Errorf("fleet: %d in modes + %d awaiting data != %d total", inModes, fleet.AwaitingData, fleet.Total)
	}
	if len(home.Machines) != fleet.Total {
		t.Errorf("home: %d machine rows, want %d", len(home.Machines), fleet.Total)
	}
	if cycles := home.Cycles; cycles.CompletedCount+cycles.UnfinishedCount+cycles.RunningCount != cycles.CycleCount {
		t.Errorf("cycles: %+v do not add up", cycles)
	}
	if len(home.Hourly) != 24 {
		t.Errorf("hourly: %d rows, want 24", len(home.Hourly))
	}
	recorded := int64(0)
	for _, group := range home.Time.Groups {
		recorded += group.DurationSec
	}
	if recorded > home.Time.SpanSec+int64(fleet.Total) {
		t.Errorf("time: %d s recorded, more than the %d s span", recorded, home.Time.SpanSec)
	}
	byMode := int64(0)
	for _, mode := range home.Time.Modes {
		byMode += mode.DurationSec
	}
	if diff := recorded - byMode; diff < 0 || diff > int64(len(home.Time.Groups)+len(home.Time.Modes)) {
		t.Errorf("time: %d s by mode, %d s by group", byMode, recorded)
	}
	if home.OEE != nil && home.OEE.OEE != nil && (*home.OEE.OEE < 0 || *home.OEE.OEE > 100) {
		t.Errorf("fleet OEE %v out of range", *home.OEE.OEE)
	}
	t.Logf("home: fleet %+v, cycles %+v, plan %+v, %d alarms", fleet, home.Cycles, home.Plan, len(home.Alarms))

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if result := sendRequest(t, router, http.MethodGet, apiV1+"/dashboard/home?date="+yesterday, nil); !result.Success {
		t.Errorf("home dashboard of yesterday: %s", result.Message)
	}
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	for _, query := range []string{"?date=bad", "?date=" + tomorrow, "?userId=x"} {
		if result := sendRequest(t, router, http.MethodGet, apiV1+"/dashboard/home"+query, nil); result.StatusCode != http.StatusBadRequest {
			t.Errorf("home dashboard %s: status %d, want 400", query, result.StatusCode)
		}
	}
}

func mustDecode(t *testing.T, data json.RawMessage, target any) {
	t.Helper()
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
