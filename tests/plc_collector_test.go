package tests

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/ua"

	"7th-sense-server/internal/cert"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/monitor"
	"7th-sense-server/internal/opcua"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// plcTestVariableID is hydPressure: an active number variable with only a plc address (Siemens),
// so the PLC worker reads it.
const plcTestVariableID = 39

// startPLCTestServer runs a gopcua OPC UA server (security None, anonymous) on a free localhost
// port with a HydPressure node, until the test ends.
func startPLCTestServer(t *testing.T) (*server.MapNamespace, int) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	opcServer := server.New(
		server.EnableSecurity("None", ua.MessageSecurityModeNone),
		server.EnableAuthMode(ua.UserTokenTypeAnonymous),
		server.EndPoint("127.0.0.1", port),
		server.SetLogger(slog.New(slog.DiscardHandler)),
	)
	namespace := server.NewMapNamespace(opcServer, "PlcTest")
	namespace.Data["HydPressure"] = 41.5 // SetValue needs a started server
	if err := opcServer.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = opcServer.Close() })
	return namespace, port
}

// waitFor polls check every 200 ms until it returns true or 20 s pass.
func waitFor(t *testing.T, what string, check func() bool) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if check() {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// TestPLCCollector runs the PLC collector against a real OPC UA server: an active Siemens machine
// whose hydPressure node id points at the server is connected, its value reaches the live hub,
// machine_variable_values and machine_telementry, a change is stored, and the diagnostics show
// the connection and the value. A server-side collector running beside the tests (PLC_ENABLED,
// every machine) could connect the test machine too while it is active.
func TestPLCCollector(t *testing.T) {
	namespace, port := startPLCTestServer(t)
	company := newTestCompany(t, "plcc"+uniqueSuffix(), 1)
	machineID := company.machineIDs[0]
	nodeID := fmt.Sprintf("ns=%d;s=HydPressure", namespace.ID())

	execDiagnosticsSQL(t, `UPDATE controllers SET cnc_engine = $2 WHERE id = $1`, company.controllerID, model.EngineSiemens)
	execDiagnosticsSQL(t, `UPDATE machines SET status = 'active', live_ip = '127.0.0.1', port = $2, security_mode = 'None',
		username = '', password = '' WHERE id = $1`, machineID, port)
	execDiagnosticsSQL(t, `INSERT INTO machine_variables (machine_id, id, source, address) VALUES ($1, $2, 'plc', $3)`,
		machineID, plcTestVariableID, nodeID)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = testDBPool.Exec(ctx, `UPDATE machines SET status = 'inactive' WHERE id = $1`, machineID)
		_, _ = testDBPool.Exec(ctx, `DELETE FROM machine_variables WHERE machine_id = $1`, machineID)
		_, _ = testDBPool.Exec(ctx, `DELETE FROM machine_variable_values WHERE machine_id = $1`, machineID)
		_, _ = testDBPool.Exec(ctx, `DELETE FROM machine_telementry WHERE machine_id = $1`, machineID)
	})

	hub := live.NewHub(200 * time.Millisecond)
	board := opcua.NewStatusBoard()
	cncConfig := config.CNCConfig{
		StopTimeout: 5 * time.Second,
		Machines:    config.MachineSelection{IDs: map[int64]bool{machineID: true}},
		Logs:        config.CNCLogSettings{PLC: true},
		PLC: config.PLCConfig{Enabled: true, PublishInterval: 100 * time.Millisecond, SamplingInterval: 50 * time.Millisecond,
			ConnectTimeout: 5 * time.Second, HealthInterval: 5 * time.Second},
	}
	telemetryService := service.NewMachineTelemetryService(testDBPool, repository.NewMachineTelemetryRepository(testDBPool))
	collector := monitor.NewPlcCollector(cncConfig, telemetryService, hub, board, cert.Paths{})

	ctx, stop := context.WithCancel(context.Background())
	collected := make(chan struct{})
	go func() {
		defer close(collected)
		collector.RunCollection(ctx)
	}()
	t.Cleanup(func() {
		stop()
		<-collected
	})

	storedValue := func() *float64 {
		var value *float64
		_ = testDBPool.QueryRow(context.Background(), `SELECT value FROM machine_variable_values WHERE machine_id = $1 AND variable_id = $2`,
			machineID, plcTestVariableID).Scan(&value)
		return value
	}
	telemetryCount := func() int {
		var count int
		_ = testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_telementry WHERE machine_id = $1 AND variable_id = $2`,
			machineID, plcTestVariableID).Scan(&count)
		return count
	}

	t.Run("the first value is published and stored", func(t *testing.T) {
		waitFor(t, "machine_variable_values = 41.5", func() bool { value := storedValue(); return value != nil && *value == 41.5 })
		if count := telemetryCount(); count != 1 {
			t.Errorf("machine_telementry rows %d, want 1", count)
		}
		if value := hub.Snapshot(machineID, []string{"hydPressure"}).Values["hydPressure"]; value != 41.5 {
			t.Errorf("hub value %v", value)
		}
		status, found := board.Get(machineID)
		if !found || status.State != opcua.StateConnected || status.Nodes[plcTestVariableID].Status != "" {
			t.Errorf("board status %+v", status)
		}
	})

	t.Run("a change is stored once", func(t *testing.T) {
		namespace.SetValue("HydPressure", 43.25)
		waitFor(t, "machine_variable_values = 43.25", func() bool { value := storedValue(); return value != nil && *value == 43.25 })
		if count := telemetryCount(); count != 2 {
			t.Errorf("machine_telementry rows %d, want 2", count)
		}
	})

	t.Run("the diagnostics show the PLC section", func(t *testing.T) {
		router := routes.New(routes.Handlers{MachineDiagnostics: handler.NewMachineDiagnosticsHandler(service.NewMachineDiagnosticsService(
			repository.NewMachineDiagnosticsRepository(testDBPool), newMachineService(), telemetryService, hub,
			config.CNCConfig{Quiet: 5 * time.Minute, PLC: config.PLCConfig{Enabled: true}}, "386", board))}, asTestUser)
		result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/machines/%d/diagnostics", apiV1, machineID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var diagnostics dto.MachineDiagnosticsResponse
		decodeData(t, result, &diagnostics)

		plc := diagnostics.PLC
		if !plc.Supported || plc.State != opcua.StateConnected || plc.Status == dto.DiagnosticFailed {
			t.Errorf("plc section %+v", plc)
		}
		found := false
		for _, variable := range plc.Variables {
			if variable.ID == plcTestVariableID {
				found = true
				if variable.NodeID != nodeID || !variable.Subscribed || variable.Value != 43.25 {
					t.Errorf("hydPressure row %+v", variable)
				}
			}
		}
		if !found {
			t.Errorf("hydPressure not in the PLC section: %+v", plc.Variables)
		}
		for _, variable := range diagnostics.Variables {
			if variable.ID == plcTestVariableID {
				t.Errorf("hydPressure also listed with the snapshot variables: %+v", variable)
			}
		}
	})
	t.Run("the PLC node IDs page shows the reader and tests node ids", func(t *testing.T) {
		router := routes.New(routes.Handlers{MachinePlcVariable: handler.NewMachinePlcVariableHandler(service.NewMachinePlcVariableService(
			testDBPool, repository.NewMachinePlcVariableRepository(testDBPool), newMachineService(),
			repository.NewMachineDiagnosticsRepository(testDBPool), cncConfig.PLC, board, cert.Paths{}))}, asTestUser)
		path := fmt.Sprintf("%s/machines/%d/plc-variables", apiV1, machineID)

		result := sendRequest(t, router, http.MethodGet, path, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var list dto.MachinePlcVariablesResponse
		decodeData(t, result, &list)
		if !list.Reader.Enabled || !list.Reader.Supported || list.Reader.State != opcua.StateConnected {
			t.Errorf("reader %+v", list.Reader)
		}
		for _, row := range list.Variables {
			if row.VariableID == plcTestVariableID && (row.ReaderNodeID != nodeID || row.ReaderStatus != "") {
				t.Errorf("hydPressure reader node %q status %q", row.ReaderNodeID, row.ReaderStatus)
			}
		}

		missing := fmt.Sprintf("ns=%d;s=NoSuchNode", namespace.ID())
		result = sendRequest(t, router, http.MethodPost, path+"/test", map[string]any{"nodeIds": []string{" " + nodeID, missing, nodeID}})
		expectStatus(t, result, http.StatusOK, "success")
		var test dto.PlcNodeTestResponse
		decodeData(t, result, &test)
		if !test.Tested || !test.Connected || len(test.Results) != 2 {
			t.Fatalf("test %+v", test)
		}
		if good := test.Results[0]; good.NodeID != nodeID || !good.OK || good.Value != 43.25 || good.ValueType != "Double" {
			t.Errorf("good node %+v", good)
		}
		if bad := test.Results[1]; bad.NodeID != missing || bad.OK || bad.Status == "" {
			t.Errorf("missing node %+v", bad)
		}

		result = sendRequest(t, router, http.MethodPost, path+"/test", map[string]any{"nodeIds": []string{"not a node id"}})
		expectStatus(t, result, http.StatusBadRequest, "Please enter a valid OPC UA node ID, e.g. ns=2;s=/Random/Plc/DB3300.DBX3.0")
	})
}
