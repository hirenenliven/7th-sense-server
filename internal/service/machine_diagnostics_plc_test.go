package service

import (
	"testing"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/opcua"
)

// TestPLCSection checks the PLC part of the diagnostics: which machines are read, the verdict
// per worker state, the hint per error and the rows of the variables.
func TestPLCSection(t *testing.T) {
	siemens, fanuc := engineOf(model.EngineSiemens), engineOf(model.EngineFanuc)
	machine := model.Machine{LiveIP: "10.0.0.5", Port: 4840}
	variables := []model.DiagnosticVariable{
		{ID: 39, Code: "hydPressure", Enabled: true, Source: model.VariableSourcePLC, Address: "ns=2;s=/Random/Plc/DB9004.DBD8"},
		{ID: 40, Code: "hydTemp", Enabled: true, Source: model.VariableSourcePLC, Address: "ns=2;s=/Random/Plc/DB9004.DBD12"},
	}
	now := time.Now()

	cases := []struct {
		name     string
		enabled  bool
		engine   *int
		status   opcua.MachineStatus
		reported bool
		want     string
		part     string
	}{
		{"off", false, siemens, opcua.MachineStatus{}, false, dto.DiagnosticSkipped, "PLC_ENABLED"},
		{"fanuc", true, fanuc, opcua.MachineStatus{}, false, dto.DiagnosticSkipped, "Siemens"},
		{"no worker", true, siemens, opcua.MachineStatus{}, false, dto.DiagnosticFailed, "CNC_MACHINE_IDS"},
		{"login refused", true, siemens, opcua.MachineStatus{State: opcua.StateFailed, LastError: "connect: BadUserAccessDenied"}, true,
			dto.DiagnosticFailed, "Username and Password"},
		{"certificate", true, siemens, opcua.MachineStatus{State: opcua.StateFailed, LastError: "connect: BadSecurityChecksFailed"}, true,
			dto.DiagnosticFailed, "client_cert.der"},
		{"unreachable", true, siemens, opcua.MachineStatus{State: opcua.StateFailed, LastError: "get endpoints of opc.tcp://10.0.0.5:4840: i/o timeout"}, true,
			dto.DiagnosticFailed, "4840"},
		{"node refused", true, siemens, opcua.MachineStatus{State: opcua.StateConnected, LastChangeAt: &now,
			Nodes: map[int]opcua.NodeStatus{39: {}, 40: {Status: "BadNodeIdUnknown"}}}, true, dto.DiagnosticWarning, "refused 1 of 2"},
		{"connected", true, siemens, opcua.MachineStatus{State: opcua.StateConnected, LastChangeAt: &now,
			Nodes: map[int]opcua.NodeStatus{39: {}, 40: {}}}, true, dto.DiagnosticOK, "subscribed to 2"},
	}
	for _, testCase := range cases {
		section := plcSection(testCase.enabled, testCase.engine, machine, testCase.status, testCase.reported, variables, map[string]any{"hydPressure": 12.5})
		expectDiagnostic(t, dto.DiagnosticCheck{Key: testCase.name, Status: section.Status, Message: section.Message, Hint: section.Hint},
			testCase.want, testCase.part)
	}

	section := plcSection(true, siemens, machine, cases[6].status, true, variables, map[string]any{"hydPressure": 12.5})
	if row := section.Variables[0]; !row.Subscribed || row.Value != 12.5 || row.NodeID != variables[0].Address {
		t.Errorf("hydPressure row %+v", row)
	}
	if row := section.Variables[1]; row.Subscribed || row.NodeStatus != "BadNodeIdUnknown" {
		t.Errorf("hydTemp row %+v", row)
	}
	if off := plcSection(true, fanuc, machine, opcua.MachineStatus{}, false, nil, nil); off.State != "off" || off.Endpoint != "opc.tcp://10.0.0.5:4840" {
		t.Errorf("fanuc section %+v", off)
	}
}
