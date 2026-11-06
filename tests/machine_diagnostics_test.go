package tests

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// newDiagnosticsRouter serves the diagnostics routes, acting as authenticate's user.
func newDiagnosticsRouter(authenticate middleware.Middleware, hub *live.Hub) http.Handler {
	telemetryService := service.NewMachineTelemetryService(testDBPool, repository.NewMachineTelemetryRepository(testDBPool))
	return routes.New(routes.Handlers{MachineDiagnostics: handler.NewMachineDiagnosticsHandler(service.NewMachineDiagnosticsService(
		repository.NewMachineDiagnosticsRepository(testDBPool), newMachineService(), telemetryService, hub,
		config.CNCConfig{Quiet: 5 * time.Minute}, "386"))}, authenticate)
}

func diagnosticsFor(t *testing.T, router http.Handler, machineID int64) dto.MachineDiagnosticsResponse {
	t.Helper()
	result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/machines/%d/diagnostics", apiV1, machineID), nil)
	expectStatus(t, result, http.StatusOK, "success")
	var diagnostics dto.MachineDiagnosticsResponse
	decodeData(t, result, &diagnostics)
	return diagnostics
}

func expectCheck(t *testing.T, check dto.DiagnosticCheck, status, messagePart string) {
	t.Helper()
	if check.Status != status || !strings.Contains(check.Message, messagePart) {
		t.Errorf("check %s: got %s %q; want %s containing %q", check.Key, check.Status, check.Message, status, messagePart)
	}
}

func checkOf(t *testing.T, diagnostics dto.MachineDiagnosticsResponse, key string) dto.DiagnosticCheck {
	t.Helper()
	for _, check := range diagnostics.Checks {
		if check.Key == key {
			return check
		}
	}
	t.Fatalf("no check %q in %+v", key, diagnostics.Checks)
	return dto.DiagnosticCheck{}
}

func execDiagnosticsSQL(t *testing.T, query string, args ...any) {
	t.Helper()
	if _, err := testDBPool.Exec(context.Background(), query, args...); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
}

// TestMachineDiagnostics checks who may read a machine's diagnostics, what they return from the
// database and the hub, and the network test against a local port. The machine stays Inactive: a
// collector running beside the tests would connect an active one with an engine. The reason logic
// of active machines is tested in internal/service without a database.
func TestMachineDiagnostics(t *testing.T) {
	company := newTestCompany(t, "dg"+uniqueSuffix(), 1)
	machineID := company.machineIDs[0]
	execDiagnosticsSQL(t, `UPDATE machines SET status = 'inactive' WHERE id = $1`, machineID)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = testDBPool.Exec(ctx, `DELETE FROM machine_live_status WHERE machine_id = $1`, machineID)
		_, _ = testDBPool.Exec(ctx, `DELETE FROM machine_events WHERE machine_id = $1`, machineID)
		_, _ = testDBPool.Exec(ctx, `DELETE FROM machine_mode_log WHERE machine_id = $1`, machineID)
	})

	hub := live.NewHub(200 * time.Millisecond)
	router := newDiagnosticsRouter(asTestUser, hub)
	diagnosticsPath := fmt.Sprintf("%s/machines/%d/diagnostics", apiV1, machineID)
	networkPath := diagnosticsPath + "/network"

	t.Run("only the superadmin", func(t *testing.T) {
		customer := newDiagnosticsRouter(actAs(company.customerID, model.RoleCustomer), hub)
		expectStatus(t, sendRequest(t, customer, http.MethodGet, diagnosticsPath, nil), http.StatusForbidden, forbiddenMessage)
		expectStatus(t, sendRequest(t, customer, http.MethodGet, networkPath, nil), http.StatusForbidden, forbiddenMessage)
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines/999999999/diagnostics", nil), http.StatusNotFound, "Machine not found")
	})

	t.Run("an inactive machine stops the checks after it", func(t *testing.T) {
		diagnostics := diagnosticsFor(t, router, machineID)
		expectCheck(t, checkOf(t, diagnostics, dto.DiagnosticCheckSetup), dto.DiagnosticFailed, "Inactive")
		expectCheck(t, checkOf(t, diagnostics, dto.DiagnosticCheckController), dto.DiagnosticSkipped, "not collecting")
		// The fixture's controller has no engine: every variable says so
		if len(diagnostics.Variables) == 0 || diagnostics.Variables[0].Readable || !strings.Contains(diagnostics.Variables[0].Reason, "no CNC engine") {
			t.Errorf("variables without an engine: %+v", diagnostics.Variables[:1])
		}
	})

	execDiagnosticsSQL(t, `UPDATE controllers SET cnc_engine = $2 WHERE id = $1`, company.controllerID, model.EngineFanuc)

	t.Run("variables say whether they can be read", func(t *testing.T) {
		diagnostics := diagnosticsFor(t, router, machineID)
		byCode := map[string]dto.DiagnosticVariableRow{}
		readable := 0
		for _, row := range diagnostics.Variables {
			byCode[row.Code] = row
			if row.Readable {
				readable++
			} else if row.Reason == "" {
				t.Errorf("%s is not readable and gives no reason", row.Code)
			}
		}
		if readable == 0 || !byCode["SpindleLoad"].Readable || byCode["SpindleLoad"].Address == "" {
			t.Errorf("%d readable, SpindleLoad %+v", readable, byCode["SpindleLoad"])
		}
		if readable == len(diagnostics.Variables) {
			t.Error("expected some variables without a Fanuc snapshot address")
		}
	})

	reportStatus := func(t *testing.T, connected bool) {
		t.Helper()
		execDiagnosticsSQL(t, `INSERT INTO machine_live_status (machine_id, is_connected, updated_at) VALUES ($1, $2, NOW())
			ON CONFLICT (machine_id) DO UPDATE SET is_connected = $2, updated_at = NOW()`, machineID, connected)
	}

	t.Run("values, status, mode and events are returned", func(t *testing.T) {
		hub.PublishValues(machineID, time.Now(), map[string]any{"SpindleLoad": 41.5, "progName": "O1234"})
		reportStatus(t, true)
		execDiagnosticsSQL(t, `INSERT INTO machine_events (machine_id, event_type, message) VALUES ($1, $2, 'no usable reading for 30s')`, machineID, model.EventStale)
		execDiagnosticsSQL(t, `INSERT INTO machine_mode_log (machine_id, mode_id, program_name, started_at) VALUES ($1, $2, 'O1234', NOW())`, machineID, model.ModeIdle)

		diagnostics := diagnosticsFor(t, router, machineID)
		for _, row := range diagnostics.Variables {
			if row.Code == "SpindleLoad" && row.Value != 41.5 {
				t.Errorf("SpindleLoad value %v", row.Value)
			}
			if row.Code == "progName" && row.Value != "O1234" {
				t.Errorf("progName value %v", row.Value)
			}
		}
		if diagnostics.Status == nil || !diagnostics.Status.IsConnected {
			t.Errorf("status %+v", diagnostics.Status)
		}
		if diagnostics.Mode == nil || diagnostics.Mode.ModeID != model.ModeIdle || diagnostics.Mode.ProgramName != "O1234" {
			t.Errorf("mode %+v", diagnostics.Mode)
		}
		if len(diagnostics.Events) != 1 || diagnostics.Events[0].Type != model.EventStale {
			t.Errorf("events %+v", diagnostics.Events)
		}
		expectCheck(t, checkOf(t, diagnostics, dto.DiagnosticCheckMode), dto.DiagnosticOK, "reset")
	})

	t.Run("network", func(t *testing.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()
		openPort := listener.Addr().(*net.TCPAddr).Port
		network := func(t *testing.T, ip string, port int) dto.MachineNetworkResponse {
			t.Helper()
			execDiagnosticsSQL(t, `UPDATE machines SET live_ip = $2, port = $3 WHERE id = $1`, machineID, ip, port)
			result := sendRequest(t, router, http.MethodGet, networkPath, nil)
			expectStatus(t, result, http.StatusOK, "success")
			var answer dto.MachineNetworkResponse
			decodeData(t, result, &answer)
			return answer
		}
		reportStatus(t, false)

		open := network(t, "127.0.0.1", openPort)
		expectCheck(t, open.IP, dto.DiagnosticOK, "127.0.0.1")
		expectCheck(t, open.Port, dto.DiagnosticOK, "accepted")
		if open.LatencyMs == nil {
			t.Error("no latency for an open port")
		}

		// A port nothing listens on: the host is there, the port is the problem
		_ = listener.Close()
		refused := network(t, "127.0.0.1", openPort)
		expectCheck(t, refused.IP, dto.DiagnosticOK, "127.0.0.1")
		expectCheck(t, refused.Port, dto.DiagnosticFailed, "refused")
		if refused.LatencyMs != nil {
			t.Error("latency for a closed port")
		}

		// The collector is connected: a failed test is the test's problem, so it only warns
		reportStatus(t, true)
		expectCheck(t, network(t, "127.0.0.1", openPort).Port, dto.DiagnosticWarning, "refused")
		reportStatus(t, false)

		// Text that is not an address can never reach ping as an option
		invalid := network(t, "-bad.invalid", openPort)
		expectCheck(t, invalid.IP, dto.DiagnosticFailed, "could not be resolved")
		expectCheck(t, invalid.Port, dto.DiagnosticSkipped, "Not tested")

		execDiagnosticsSQL(t, `UPDATE machines SET live_ip = '' WHERE id = $1`, machineID)
		expectStatus(t, sendRequest(t, router, http.MethodGet, networkPath, nil), http.StatusBadRequest, "This machine has no live IP and port to test")
	})
}
