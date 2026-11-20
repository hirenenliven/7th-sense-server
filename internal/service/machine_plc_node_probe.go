package service

import (
	"context"
	"fmt"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/opcua"
)

// plcProbeCloseTimeout bounds closing the test session, so a dead controller cannot hold the request.
const plcProbeCloseTimeout = 2 * time.Second

// TestNodeIDs reads each node id once from the machine's OPC UA server, in a session of its own
// opened and closed for this request, the way the PLC reader connects (live IP and port, security
// mode, login, client certificate). Nothing is stored or published. A machine the PLC reader does
// not read (not Siemens) or without an address is not tried (tested = false); a connect or read
// failure is part of the answer (connected = false, message, hint), not an error.
func (plcVariableService *MachinePlcVariableService) TestNodeIDs(ctx context.Context, actor model.AuthUser, machineID int64, request dto.TestPlcNodeIDsRequest) (dto.PlcNodeTestResponse, error) {
	machine, err := plcVariableService.machineService.GetByID(ctx, actor, nil, machineID)
	if err != nil {
		return dto.PlcNodeTestResponse{}, err
	}
	engine, err := plcVariableService.diagnosticsRepository.FindEngine(ctx, machineID)
	if err != nil {
		return dto.PlcNodeTestResponse{}, err
	}

	test := dto.PlcNodeTestResponse{
		Endpoint:     opcua.Endpoint(machine.LiveIP, machine.Port),
		SecurityMode: machine.SecurityMode,
		TestedAt:     time.Now(),
		Results:      []dto.PlcNodeTestResult{},
	}
	switch {
	case !model.ReadsPLCOverOPCUA(engine):
		test.Message = "PLC node IDs are tested over OPC UA on Siemens controllers only."
		test.Hint = "This controller is not read over OPC UA, so its PLC node IDs cannot be tested."
		return test, nil
	case machine.LiveIP == "" || machine.Port <= 0:
		test.Message = "The machine has no Live IP or Port."
		test.Hint = "Set Live IP and Port on the machine form (Siemens OPC UA is usually 4840)."
		return test, nil
	}

	test.Tested = true
	plcVariableService.readNodeIDs(ctx, machine, request.NodeIDs, &test)
	test.DurationMs = time.Since(test.TestedAt).Milliseconds()
	return test, nil
}

// readNodeIDs connects to machine, reads nodeIDs in one request and fills test with the outcome.
func (plcVariableService *MachinePlcVariableService) readNodeIDs(ctx context.Context, machine model.Machine, nodeIDs []string, test *dto.PlcNodeTestResponse) {
	connectTimeout := plcVariableService.plcConfig.ConnectTimeout
	connection, err := opcua.Connect(ctx, opcua.Settings{
		Endpoint:       test.Endpoint,
		SecurityMode:   machine.SecurityMode,
		Username:       machine.Username,
		Password:       machine.Password,
		Certificate:    plcVariableService.certificate,
		ConnectTimeout: connectTimeout,
		SessionTimeout: plcVariableService.plcConfig.SessionTimeout,
	})
	if err != nil {
		test.Message = "Cannot connect: " + err.Error()
		test.Hint = plcErrorHint(err.Error(), machine.SecurityMode)
		return
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), plcProbeCloseTimeout)
		defer cancel()
		connection.Close(closeCtx)
	}()
	test.Connected, test.SecurityPolicy = true, connection.SecurityPolicy

	readCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	results, err := connection.ReadValues(readCtx, nodeIDs)
	if err != nil {
		test.Message = "Connected, but the read failed: " + err.Error()
		test.Hint = plcErrorHint(err.Error(), machine.SecurityMode)
		return
	}

	working := 0
	for _, result := range results {
		row := dto.PlcNodeTestResult{NodeID: result.NodeID, Status: result.Status, ValueType: result.ValueType, SourceTime: result.SourceTime}
		if result.Status == "" {
			row.OK, row.Value = true, plcValueOf(result.Value)
			working++
		}
		test.Results = append(test.Results, row)
	}
	test.Message = fmt.Sprintf("Connected; %d of %d node IDs read a value.", working, len(results))
	if working < len(results) {
		test.Hint = "BadNodeIdUnknown: the node does not exist on this controller; BadUserAccessDenied / BadNotReadable: this login may not read it."
	}
}

// plcValueOf is a read value as the API sends it: text, or a number (booleans 1 / 0).
func plcValueOf(value opcua.Value) any {
	if value.IsText {
		return value.Text
	}
	return value.Number
}
