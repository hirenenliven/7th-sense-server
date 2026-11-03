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
		machineService, telemetryService, cycleService, repository.NewShiftRepository(testDBPool)))}, asTestUser)
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

	for _, path := range []string{
		fmt.Sprintf("/productivity/oee?machineIds=%d&from=%s&to=%s", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/oee?machineIds=%d&from=%s&to=%s&schedule=with", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/cycle-counts?machineIds=%d&from=%s&to=%s", machineID, weekAgo, today),
		fmt.Sprintf("/productivity/finish-cycles?machineId=%d&date=%s", machineID, today),
	} {
		if result := sendRequest(t, router, http.MethodGet, apiV1+path, nil); !result.Success {
			t.Errorf("GET %s: %s", path, result.Message)
		}
	}

	if result := sendRequest(t, router, http.MethodGet, apiV1+"/productivity/oee", nil); result.StatusCode != http.StatusBadRequest {
		t.Errorf("OEE without a machine: status %d, want 400", result.StatusCode)
	}
}

func mustDecode(t *testing.T, data json.RawMessage, target any) {
	t.Helper()
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
