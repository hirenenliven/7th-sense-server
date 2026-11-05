package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"7th-sense-server/internal/dto"
)

// TestVirtualScreen reads the virtual screen of the machine with the most live values (read
// only), and refuses a request without a machine.
func TestVirtualScreen(t *testing.T) {
	var machineID int64
	err := testDBPool.QueryRow(context.Background(), `SELECT v.machine_id FROM machine_variable_values v
		JOIN machines m ON m.id = v.machine_id AND m.is_deleted = false
		GROUP BY v.machine_id ORDER BY count(*) DESC LIMIT 1`).Scan(&machineID)
	if err != nil {
		t.Skipf("no machine with live values: %v", err)
	}
	router := newDashboardRouter()

	result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/productivity/virtual-screen?machineId=%d", apiV1, machineID), nil)
	if !result.Success {
		t.Fatalf("virtual screen: %s", result.Message)
	}
	var screen dto.VirtualScreenResponse
	mustDecode(t, result.Data, &screen)
	if screen.MachineID != machineID {
		t.Errorf("virtual screen of machine %d, want %d", screen.MachineID, machineID)
	}
	if screen.Live.UpdatedAt == nil {
		t.Error("virtual screen has no live values")
	}
	if screen.Live.FeedUnit == "" {
		t.Error("virtual screen has no feed unit")
	}
	if screen.Job != nil && screen.Job.Quantity <= 0 {
		t.Errorf("job %d has quantity %d", screen.Job.ScheduleItemID, screen.Job.Quantity)
	}
	t.Logf("machine %d: mode %v %q, program %q, live %+v, alarm %+v, cycle %+v, job %+v",
		machineID, screen.ModeID, screen.ModeName, screen.ProgramName, screen.Live, screen.Alarm, screen.Cycle, screen.Job)

	missing := sendRequest(t, router, http.MethodGet, apiV1+"/productivity/virtual-screen", nil)
	if missing.Success || missing.StatusCode != http.StatusBadRequest {
		t.Errorf("virtual screen without machineId: status %d, want 400", missing.StatusCode)
	}
}
