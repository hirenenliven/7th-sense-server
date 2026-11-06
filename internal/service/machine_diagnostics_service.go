package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

const (
	diagnosticEventLimit = 30
	// The worker writes machine_live_status every 10 s while it holds a connection
	statusFreshFor = 30 * time.Second
	// Between connects a worker waits up to 5 min (+20 % jitter) and then records a reconnect event
	workerActiveWithin = 6 * time.Minute
	// The architecture the 32-bit CNC DLL needs (cnc.ErrNotSupported)
	collectorArchitecture = "386"
)

// modeMeanings says what makes the collector decide each mode (monitor.decideMachineMode).
var modeMeanings = map[int]string{
	model.ModeInProduction: "AUTO, program running, feed override above 0, no alarm.",
	model.ModeHold:         "Program held, or running with feed override 0, an alarm or not in AUTO.",
	model.ModeIdle:         "Program reset and no alarm.",
	model.ModeDowntime:     "Alarm or emergency stop with the program not running.",
	model.ModeOffline:      "No connection to the controller.",
}

// MachineDiagnosticsService answers the superadmin's question "why does this machine send no
// data?": it puts together the machine's set-up, what its collector worker reported, its variables
// and, when asked, whether its IP and port answer. Reading it changes nothing.
type MachineDiagnosticsService struct {
	diagnosticsRepository *repository.MachineDiagnosticsRepository
	machineService        *MachineService
	telemetryService      *MachineTelemetryService
	hub                   *live.Hub // nil when LIVE_ENABLED is off
	cncConfig             config.CNCConfig
	architecture          string // this build's GOARCH
}

// NewMachineDiagnosticsService creates a MachineDiagnosticsService; architecture is this build's
// GOARCH.
func NewMachineDiagnosticsService(diagnosticsRepository *repository.MachineDiagnosticsRepository, machineService *MachineService,
	telemetryService *MachineTelemetryService, hub *live.Hub, cncConfig config.CNCConfig, architecture string) *MachineDiagnosticsService {
	return &MachineDiagnosticsService{diagnosticsRepository: diagnosticsRepository, machineService: machineService,
		telemetryService: telemetryService, hub: hub, cncConfig: cncConfig, architecture: architecture}
}

// GetDiagnostics returns the machine's set-up, collector status, mode, variables and recent events,
// with the checks worked out from them.
func (diagnosticsService *MachineDiagnosticsService) GetDiagnostics(ctx context.Context, actor model.AuthUser, machineID int64) (dto.MachineDiagnosticsResponse, error) {
	machine, err := diagnosticsService.machineService.GetByID(ctx, actor, nil, machineID)
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}
	engine, err := diagnosticsService.diagnosticsRepository.FindEngine(ctx, machineID)
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}
	status, err := diagnosticsService.diagnosticsRepository.FindLiveStatus(ctx, machineID)
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}
	mode, err := diagnosticsService.diagnosticsRepository.FindCurrentMode(ctx, machineID)
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}
	events, err := diagnosticsService.diagnosticsRepository.FindRecentEvents(ctx, machineID, diagnosticEventLimit)
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}
	versionType := 0
	if machine.ControllerVersionType != nil {
		versionType = *machine.ControllerVersionType
	}
	variables, err := diagnosticsService.diagnosticsRepository.FindVariables(ctx, machineID, engine, versionType, []string{cnc.SourceSnapshot})
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}
	opMode, err := diagnosticsService.telemetryService.GetOpModeName(ctx, machineID)
	if err != nil {
		return dto.MachineDiagnosticsResponse{}, err
	}

	now := time.Now()
	rows := diagnosticsService.variableRows(machineID, engine, variables)
	connected := isControllerConnected(status, now)
	setup := setupCheck(machine, engine)
	collector := diagnosticsService.collectorCheck(machine, engine, status, events, now)
	return dto.MachineDiagnosticsResponse{
		ServerTime:  now,
		Machine:     diagnosticMachine(machine, engine),
		LiveEnabled: diagnosticsService.hub != nil,
		QuietSec:    int(diagnosticsService.cncConfig.Quiet.Seconds()),
		Status:      status,
		Mode:        mode,
		OpMode:      opMode,
		Checks: []dto.DiagnosticCheck{
			setup,
			collector,
			diagnosticsService.controllerCheck(setup, collector, engine, status, events, connected, now),
			modeCheck(mode, connected),
			variablesCheck(rows, connected),
		},
		Variables: rows,
		Events:    events,
	}, nil
}

// CheckNetwork tests the machine's live IP and port from this server. When the test fails while the
// collector holds a working connection, the test is what is in doubt, so it only warns.
func (diagnosticsService *MachineDiagnosticsService) CheckNetwork(ctx context.Context, actor model.AuthUser, machineID int64) (dto.MachineNetworkResponse, error) {
	machine, err := diagnosticsService.machineService.GetByID(ctx, actor, nil, machineID)
	if err != nil {
		return dto.MachineNetworkResponse{}, err
	}
	if strings.TrimSpace(machine.LiveIP) == "" || machine.Port <= 0 {
		return dto.MachineNetworkResponse{}, response.NewError(http.StatusBadRequest, "This machine has no live IP and port to test")
	}
	result := probeNetwork(ctx, machine.LiveIP, machine.Port)

	status, err := diagnosticsService.diagnosticsRepository.FindLiveStatus(ctx, machineID)
	if err != nil {
		return dto.MachineNetworkResponse{}, err
	}
	if isControllerConnected(status, result.CheckedAt) {
		for _, check := range []*dto.DiagnosticCheck{&result.IP, &result.Port} {
			if check.Status == dto.DiagnosticFailed {
				check.Status = dto.DiagnosticWarning
				check.Hint = "The collector is connected to this machine right now, so its data is not affected. This test may be blocked by a firewall."
			}
		}
	}
	return result, nil
}

// isControllerConnected reports whether the worker said connected and is still reporting: after a
// server stop the row keeps its last value until the next start.
func isControllerConnected(status *model.DiagnosticLiveStatus, now time.Time) bool {
	return status != nil && status.IsConnected && now.Sub(status.UpdatedAt) <= statusFreshFor
}

func diagnosticMachine(machine model.Machine, engine *int) dto.DiagnosticMachine {
	engineName := ""
	if engine != nil {
		engineName = model.EngineNames[*engine]
		if engineName == "" {
			engineName = fmt.Sprintf("Engine %d", *engine)
		}
	}
	return dto.DiagnosticMachine{
		ID: machine.ID, Name: machine.Name, CustomerName: machine.CustomerName, CompanyName: machine.CompanyName,
		Status: machine.Status, LiveIP: machine.LiveIP, LocalIP: machine.LocalIP, Port: machine.Port, Engine: engineName,
		BrandName: machine.BrandName, ControllerName: machine.ControllerName, ControllerVersion: machine.ControllerVersion,
		ModelName: machine.ModelName,
	}
}

// setupCheck: the collector connects only an active machine with an engine, IP and port.
func setupCheck(machine model.Machine, engine *int) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: dto.DiagnosticCheckSetup, Status: dto.DiagnosticFailed}
	switch {
	case machine.Status != "active":
		check.Message = "The machine is Inactive, so the server does not connect to it."
		check.Hint = "Set Status to Active on the machine form."
	case machine.ControllerName == "":
		check.Message = "The machine has no controller, so the server does not know how to talk to it."
		check.Hint = "Pick the controller on the machine form."
	case engine == nil:
		check.Message = fmt.Sprintf("Controller %q has no CNC engine (Fanuc, Siemens, …), so the server cannot connect to it.", machine.ControllerName)
		check.Hint = "Set controllers.cnc_engine for this controller, then save the machine."
	case strings.TrimSpace(machine.LiveIP) == "":
		check.Message = "The machine has no Live IP."
		check.Hint = "Enter the Live IP on the machine form."
	case machine.Port <= 0:
		check.Message = "The machine has no port."
		check.Hint = "Enter the port on the machine form."
	default:
		check.Status = dto.DiagnosticOK
		check.Message = fmt.Sprintf("Active, %s controller, connects to %s:%d.", model.EngineNames[*engine], machine.LiveIP, machine.Port)
	}
	return check
}

// collectorCheck: is this server collecting this machine at all?
func (diagnosticsService *MachineDiagnosticsService) collectorCheck(machine model.Machine, engine *int, status *model.DiagnosticLiveStatus, events []model.DiagnosticEvent, now time.Time) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: dto.DiagnosticCheckCollector, Status: dto.DiagnosticFailed}
	switch {
	case diagnosticsService.architecture != collectorArchitecture:
		check.Message = fmt.Sprintf("This server is a %s build, and the CNC DLL is 32-bit, so nothing is collected.", diagnosticsService.architecture)
		check.Hint = "Build with GOARCH=386: .\\run.ps1, .\\build.ps1 or go tool air."
	case !diagnosticsService.cncConfig.Machines.Includes(machine.ID):
		check.Message = "CNC_MACHINE_IDS leaves this machine out, so the server does not connect to it."
		check.Hint = "Add its id to CNC_MACHINE_IDS, or set it to ALL, and restart the server."
	case machine.Status != "active" || engine == nil:
		check.Status = dto.DiagnosticSkipped
		check.Message = "No worker runs until the set-up above is fixed."
	case status != nil && now.Sub(status.UpdatedAt) <= statusFreshFor:
		check.Status = dto.DiagnosticOK
		check.Message = "A collector worker is running for this machine and reporting."
	case len(events) > 0 && now.Sub(events[0].At) <= workerActiveWithin:
		check.Status = dto.DiagnosticOK
		check.Message = "A collector worker is running for this machine and retrying the connection."
	default:
		check.Status = dto.DiagnosticWarning
		check.Message = "No sign of a collector worker for this machine in the last few minutes."
		check.Hint = "The CNC DLL may have failed to load, or the collector stopped. Look for \"cnc\" errors in logs/app.log."
	}
	return check
}

// controllerCheck: does the DLL reach the controller, and if not, what did it last say?
func (diagnosticsService *MachineDiagnosticsService) controllerCheck(setup, collector dto.DiagnosticCheck, engine *int, status *model.DiagnosticLiveStatus, events []model.DiagnosticEvent, connected bool, now time.Time) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: dto.DiagnosticCheckController}
	if !connected && (setup.Status == dto.DiagnosticFailed || collector.Status == dto.DiagnosticFailed || collector.Status == dto.DiagnosticSkipped) {
		check.Status = dto.DiagnosticSkipped
		check.Message = "Not tried: the server is not collecting this machine (see above)."
		return check
	}
	if !connected {
		check.Status = dto.DiagnosticFailed
		check.Message, check.Hint = explainNotConnected(engine, status, events)
		return check
	}

	check.Status = dto.DiagnosticOK
	check.Message = "The DLL is connected to the controller."
	if status.ConnectedSince != nil {
		check.Message = fmt.Sprintf("The DLL has been connected to the controller for %s.", shortDuration(now.Sub(*status.ConnectedSince)))
	}
	// Silence is normal: the DLL calls back only when something changes
	if status.LastSeenAt != nil && now.Sub(*status.LastSeenAt) > statusFreshFor {
		check.Message += fmt.Sprintf(" Last reading %s ago: the DLL sends one only when something changes, so a held or idle machine stays quiet (the connection is reopened after %s).",
			shortDuration(now.Sub(*status.LastSeenAt)), diagnosticsService.cncConfig.Quiet)
	}
	if status.DroppedCount > 0 {
		check.Status = dto.DiagnosticWarning
		check.Hint = fmt.Sprintf("%d values could not be stored in the database: see logs/app.log.", status.DroppedCount)
	}
	return check
}

// explainNotConnected gives the reason from the worker's newest connection event and the DLL's last
// error, in the words of the DLL guide (docs/cnc-communication.md, section 10).
func explainNotConnected(engine *int, status *model.DiagnosticLiveStatus, events []model.DiagnosticEvent) (message, hint string) {
	lastError := ""
	if status != nil {
		lastError = status.LastError
	}
	var cause *model.DiagnosticEvent
	for index := range events {
		switch events[index].Type {
		case model.EventConnectFailed, model.EventConfigError, model.EventDLLError, model.EventStale:
			cause = &events[index]
		}
		if cause != nil {
			break
		}
	}
	if lastError == "" && cause != nil && cause.Type == model.EventDLLError {
		lastError = cause.Message
	}

	switch {
	case strings.Contains(lastError, "BadRequestNotAllowed"):
		return "The Siemens controller rejected this server's certificate: " + lastError,
			"Add client_cert.der to the trusted certificates on the NCU."
	case cause == nil && status == nil:
		return "The server has never reported a connection to this controller.",
			"Check the IP and port results above, and logs/app.log."
	case cause != nil && cause.Type == model.EventConnectFailed:
		return "The DLL could not open a connection: " + cause.Message,
			"Usually the host is unreachable or times out: see the IP and port results above."
	case cause != nil && cause.Type == model.EventConfigError:
		return "The collector reports a set-up problem: " + cause.Message, ""
	}

	message = "The DLL started a session, but the controller does not answer (it reports not connected)."
	if lastError != "" {
		message += " Last DLL error: " + lastError
	}
	switch {
	case engine != nil && *engine == model.EngineSiemens:
		hint = "Check SIEMENS_USER / SIEMENS_PASS, that the OPC UA server is enabled on the NCU and that the address is right."
	case engine != nil && *engine == model.EngineFanuc:
		hint = "Check the controller is on, FOCAS Ethernet is enabled on it, the port is right and Fwlib32.dll + fwlibe1.dll are in Lib\\."
	default:
		hint = "Check the controller is on and that the machine's controller has the right CNC engine."
	}
	return message, hint
}

func modeCheck(mode *model.DiagnosticMode, connected bool) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: dto.DiagnosticCheckMode, Status: dto.DiagnosticOK}
	switch {
	case mode == nil:
		check.Status = dto.DiagnosticFailed
		check.Message = "No mode has been recorded for this machine yet."
	case mode.ModeID == model.ModeOffline:
		check.Status = dto.DiagnosticFailed
		check.Message = "The machine is Offline: no connection to the controller."
		if connected {
			check.Status = dto.DiagnosticWarning
			check.Message = "The machine shows Offline although the DLL is connected: its next reading decides the mode."
		}
	default:
		check.Message = fmt.Sprintf("%s. %s", mode.ModeName, modeMeanings[mode.ModeID])
	}
	return check
}

// variablesCheck: are values arriving? Unchanged values are not a fault: the DLL and the socket send
// changes only.
func variablesCheck(rows []dto.DiagnosticVariableRow, connected bool) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: dto.DiagnosticCheckVariables}
	readable, withValue, nonZero := 0, 0, 0
	for _, row := range rows {
		if !row.Readable {
			continue
		}
		readable++
		if row.Value == nil {
			continue
		}
		withValue++
		if row.Value != 0.0 && row.Value != "" {
			nonZero++
		}
	}
	notReadable := len(rows) - readable
	switch {
	case readable == 0:
		check.Status = dto.DiagnosticFailed
		check.Message = "This machine has no variable the server can read."
		check.Hint = "Each variable lists its reason below; most often the controller has no CNC engine."
	case !connected:
		check.Status = dto.DiagnosticSkipped
		check.Message = fmt.Sprintf("%d variables can be read once the controller is connected. Values below are the last ones known.", readable)
	case withValue == 0:
		check.Status = dto.DiagnosticWarning
		check.Message = fmt.Sprintf("Connected, but no reading has arrived since this server started (%d variables waiting).", readable)
		check.Hint = "The DLL sends a reading when something changes. Change a value on the machine; if none arrives, restart the connection by saving the machine."
	case nonZero == 0:
		check.Status = dto.DiagnosticWarning
		check.Message = fmt.Sprintf("All %d values read are 0 or empty.", withValue)
		check.Hint = "The controller is answering with empty data: check the machine type (engine) and that the controller exposes these data points."
	default:
		check.Status = dto.DiagnosticOK
		check.Message = fmt.Sprintf("%d of %d readable variables have a value.", withValue, readable)
		if notReadable > 0 {
			check.Message += fmt.Sprintf(" %d more cannot be read on this machine; each lists why.", notReadable)
		}
	}
	return check
}

// variableRows adds to each variable whether it is readable (else why) and the latest value in the
// hub, which also holds text values the database does not store.
func (diagnosticsService *MachineDiagnosticsService) variableRows(machineID int64, engine *int, variables []model.DiagnosticVariable) []dto.DiagnosticVariableRow {
	values := map[string]any{}
	if diagnosticsService.hub != nil {
		codes := make([]string, len(variables))
		for index, variable := range variables {
			codes[index] = variable.Code
		}
		values = diagnosticsService.hub.Snapshot(machineID, codes).Values
	}
	rows := make([]dto.DiagnosticVariableRow, len(variables))
	for index, variable := range variables {
		reason := variableUnreadableReason(engine, variable)
		rows[index] = dto.DiagnosticVariableRow{
			ID: variable.ID, Code: variable.Code, Name: variable.Name, Unit: variable.Unit, Category: variable.Category,
			DataType: variable.DataType, Source: variable.Source, Address: variable.Address,
			Readable: reason == "", Reason: reason, Value: values[variable.Code],
			StoredValue: variable.StoredValue, StoredAt: variable.StoredAt,
		}
	}
	return rows
}

// variableUnreadableReason says why the collector cannot read variable on this machine, "" when it
// can (the rules of worker.loadVariableConfig).
func variableUnreadableReason(engine *int, variable model.DiagnosticVariable) string {
	_, isSnapshotField := cnc.FieldIndex(variable.Address)
	switch {
	case !variable.Enabled:
		return "Turned off for this machine (machine_variables)."
	case engine == nil:
		return "The machine's controller has no CNC engine."
	case variable.Address == "" && variable.KnownSources != "":
		return fmt.Sprintf("Only has a %s address; the server reads DLL snapshot fields only today.", variable.KnownSources)
	case variable.Address == "":
		return "No address for this controller type."
	case variable.Source != cnc.SourceSnapshot:
		return fmt.Sprintf("Set to read from %s; the server reads DLL snapshot fields only today.", variable.Source)
	case !isSnapshotField:
		return fmt.Sprintf("The DLL snapshot has no field %q.", variable.Address)
	}
	return ""
}

// shortDuration writes a duration for a sentence: "45s", "12m", "3h 05m", "2d 4h".
func shortDuration(duration time.Duration) string {
	seconds := int64(duration.Seconds())
	switch {
	case seconds < 60:
		return fmt.Sprintf("%ds", max(seconds, 0))
	case seconds < 3600:
		return fmt.Sprintf("%dm", seconds/60)
	case seconds < 86400:
		return fmt.Sprintf("%dh %02dm", seconds/3600, seconds%3600/60)
	default:
		return fmt.Sprintf("%dd %dh", seconds/86400, seconds%86400/3600)
	}
}
