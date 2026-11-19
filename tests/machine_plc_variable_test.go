package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// newPlcVariableRouter serves the PLC node id routes, acting as authenticate's user.
func newPlcVariableRouter(authenticate middleware.Middleware) http.Handler {
	return routes.New(routes.Handlers{MachinePlcVariable: handler.NewMachinePlcVariableHandler(service.NewMachinePlcVariableService(
		testDBPool, repository.NewMachinePlcVariableRepository(testDBPool), newMachineService()))}, authenticate)
}

func plcVariablesFor(t *testing.T, router http.Handler, path string) dto.MachinePlcVariablesResponse {
	t.Helper()
	result := sendRequest(t, router, http.MethodGet, path, nil)
	expectStatus(t, result, http.StatusOK, "success")
	var list dto.MachinePlcVariablesResponse
	decodeData(t, result, &list)
	return list
}

func savePlcNodeID(t *testing.T, router http.Handler, path, nodeID string) dto.MachinePlcVariableResponse {
	t.Helper()
	result := sendRequest(t, router, http.MethodPut, path, map[string]string{"nodeId": nodeID})
	expectStatus(t, result, http.StatusOK, "success")
	var row dto.MachinePlcVariableResponse
	decodeData(t, result, &row)
	return row
}

// TestMachinePlcVariables checks who may set a machine's PLC node ids, that the default is used
// until one is saved, that saving the default or resetting drops the machine's own (keeping other
// overrides), and that the collector ignores a node id it cannot read yet. The machine stays
// Inactive so a collector running beside the tests does not connect it. Needs the seeded Siemens
// PLC addresses (cnc_variable_addresses).
func TestMachinePlcVariables(t *testing.T) {
	company := newTestCompany(t, "plc"+uniqueSuffix(), 1)
	machineID := company.machineIDs[0]
	execDiagnosticsSQL(t, `UPDATE machines SET status = 'inactive', controller_version_type = 1 WHERE id = $1`, machineID)
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `DELETE FROM machine_variables WHERE machine_id = $1`, machineID)
	})

	router := newPlcVariableRouter(asTestUser)
	listPath := fmt.Sprintf("%s/machines/%d/plc-variables", apiV1, machineID)

	t.Run("only the superadmin", func(t *testing.T) {
		customer := newPlcVariableRouter(actAs(company.customerID, model.RoleCustomer))
		expectStatus(t, sendRequest(t, customer, http.MethodGet, listPath, nil), http.StatusForbidden, forbiddenMessage)
		expectStatus(t, sendRequest(t, customer, http.MethodPut, listPath+"/1", map[string]string{"nodeId": "ns=2;s=x"}), http.StatusForbidden, forbiddenMessage)
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines/999999999/plc-variables", nil), http.StatusNotFound, "Machine not found")
	})

	t.Run("a controller without PLC defaults reads the Siemens defaults", func(t *testing.T) {
		for _, engine := range []any{nil, model.EngineFanuc} {
			execDiagnosticsSQL(t, `UPDATE controllers SET cnc_engine = $2 WHERE id = $1`, company.controllerID, engine)
			list := plcVariablesFor(t, router, listPath)
			if len(list.Variables) == 0 || list.DefaultCount != len(list.Variables) || list.UpdatedCount != 0 {
				t.Fatalf("engine %v: %d variables, defaultCount %d, updatedCount %d", engine, len(list.Variables), list.DefaultCount, list.UpdatedCount)
			}
			for _, row := range list.Variables {
				if row.NodeStatus != dto.PlcNodeDefault || row.NodeID == "" || row.NodeID != row.DefaultNodeID {
					t.Errorf("engine %v: %+v", engine, row)
				}
			}
		}

		// The collector takes the same default once it can read plc addresses
		list := plcVariablesFor(t, router, listPath)
		fanuc := model.EngineFanuc
		configs, err := repository.NewMachineTelemetryRepository(testDBPool).FindVariableConfig(context.Background(),
			model.MonitoredMachine{ID: machineID, CNCEngine: &fanuc, ControllerVersionType: 1}, []string{cnc.SourceSnapshot, model.VariableSourcePLC})
		if err != nil {
			t.Fatal(err)
		}
		addresses := map[int]string{}
		for _, config := range configs {
			addresses[config.ID] = config.Address
		}
		for _, row := range list.Variables {
			address, read := addresses[row.VariableID]
			// A variable with a Fanuc snapshot field reads that (cheaper); node ids start with ns=
			if read && strings.HasPrefix(address, "ns=") && address != row.DefaultNodeID {
				t.Errorf("%s: collector reads %q, page shows default %q", row.Code, address, row.DefaultNodeID)
			}
		}
	})

	execDiagnosticsSQL(t, `UPDATE controllers SET cnc_engine = $2 WHERE id = $1`, company.controllerID, model.EngineSiemens)

	list := plcVariablesFor(t, router, listPath)
	var withDefault []dto.MachinePlcVariableResponse
	for _, row := range list.Variables {
		if row.DefaultNodeID != "" {
			withDefault = append(withDefault, row)
		}
	}
	if len(withDefault) < 2 {
		t.Fatalf("expected Siemens PLC defaults, got %d", len(withDefault))
	}
	first, second := withDefault[0], withDefault[1]
	firstPath := fmt.Sprintf("%s/%d", listPath, first.VariableID)
	secondPath := fmt.Sprintf("%s/%d", listPath, second.VariableID)

	t.Run("every variable starts on its default", func(t *testing.T) {
		if list.UpdatedCount != 0 {
			t.Errorf("updatedCount %d", list.UpdatedCount)
		}
		for _, row := range list.Variables {
			if row.NodeStatus != dto.PlcNodeDefault || row.NodeID == "" || row.NodeID != row.DefaultNodeID {
				t.Errorf("%s: %+v", row.Code, row)
			}
		}
	})

	t.Run("invalid node ids and variables are refused", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPut, firstPath, map[string]string{"nodeId": " "}), http.StatusBadRequest, response.EnterPlcNodeID.Message)
		expectStatus(t, sendRequest(t, router, http.MethodPut, firstPath, map[string]string{"nodeId": "DB3300.DBX3.0"}), http.StatusBadRequest, response.InvalidPlcNodeID.Message)
		expectStatus(t, sendRequest(t, router, http.MethodPut, listPath+"/999999", map[string]string{"nodeId": "ns=2;s=x"}), http.StatusNotFound, response.PlcVariableNotFound.Message)
	})

	customNodeID := "ns=2;s=/Random/Plc/DB9999.DBX1.1"

	t.Run("a saved node id is the machine's own", func(t *testing.T) {
		row := savePlcNodeID(t, router, firstPath, "  "+customNodeID+" ")
		if row.IsDefault || row.NodeID != customNodeID || row.DefaultNodeID != first.DefaultNodeID {
			t.Errorf("saved row %+v", row)
		}
		if again := plcVariablesFor(t, router, listPath); again.UpdatedCount != 1 {
			t.Errorf("after save: updatedCount %d", again.UpdatedCount)
		}
	})

	t.Run("the collector ignores a node id it cannot read", func(t *testing.T) {
		engine := model.EngineSiemens
		configs, err := repository.NewMachineTelemetryRepository(testDBPool).FindVariableConfig(context.Background(),
			model.MonitoredMachine{ID: machineID, CNCEngine: &engine, ControllerVersionType: 1}, []string{cnc.SourceSnapshot})
		if err != nil {
			t.Fatal(err)
		}
		for _, config := range configs {
			if config.ID == first.VariableID && config.Address == customNodeID {
				t.Errorf("a plc node id came back for a snapshot-only collector: %+v", config)
			}
		}
	})

	t.Run("saving the default resets", func(t *testing.T) {
		row := savePlcNodeID(t, router, firstPath, first.DefaultNodeID)
		if !row.IsDefault || row.NodeID != first.DefaultNodeID {
			t.Errorf("row %+v", row)
		}
		var count int
		if err := testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_variables WHERE machine_id = $1`, machineID).Scan(&count); err != nil || count != 0 {
			t.Errorf("machine_variables rows %d (%v)", count, err)
		}
	})

	t.Run("reset keeps the variable's other overrides", func(t *testing.T) {
		execDiagnosticsSQL(t, `INSERT INTO machine_variables (machine_id, id, is_enabled) VALUES ($1, $2, false)`, machineID, first.VariableID)
		savePlcNodeID(t, router, firstPath, customNodeID)
		result := sendRequest(t, router, http.MethodDelete, firstPath, nil)
		expectStatus(t, result, http.StatusOK, "success")

		var enabled bool
		var source *string
		err := testDBPool.QueryRow(context.Background(), `SELECT is_enabled, source FROM machine_variables WHERE machine_id = $1 AND id = $2`,
			machineID, first.VariableID).Scan(&enabled, &source)
		if err != nil || enabled || source != nil {
			t.Errorf("row after reset: enabled %v, source %v (%v)", enabled, source, err)
		}
	})

	t.Run("reset all", func(t *testing.T) {
		savePlcNodeID(t, router, firstPath, customNodeID)
		savePlcNodeID(t, router, secondPath, customNodeID)
		result := sendRequest(t, router, http.MethodDelete, listPath, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var reset dto.MachinePlcVariablesResponse
		decodeData(t, result, &reset)
		if reset.UpdatedCount != 0 {
			t.Errorf("updatedCount after reset all: %d", reset.UpdatedCount)
		}
		if again := plcVariablesFor(t, router, listPath); again.UpdatedCount != 0 {
			t.Errorf("stored updatedCount after reset all: %d", again.UpdatedCount)
		}
	})
}
