package service

import (
	"strings"
	"testing"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
)

func engineOf(engine int) *int { return &engine }

func expectDiagnostic(t *testing.T, check dto.DiagnosticCheck, status, part string) {
	t.Helper()
	if check.Status != status || !strings.Contains(check.Message+" "+check.Hint, part) {
		t.Errorf("got %s %q (hint %q); want %s containing %q", check.Status, check.Message, check.Hint, status, part)
	}
}

func TestSetupCheck(t *testing.T) {
	ready := model.Machine{Status: "active", ControllerName: "0i-F", LiveIP: "10.0.0.5", Port: 8193}
	fanuc := engineOf(model.EngineFanuc)
	cases := []struct {
		name   string
		change func(machine *model.Machine)
		engine *int
		status string
		part   string
	}{
		{"ready", func(*model.Machine) {}, fanuc, dto.DiagnosticOK, "Fanuc controller, connects to 10.0.0.5:8193"},
		{"inactive", func(machine *model.Machine) { machine.Status = "inactive" }, fanuc, dto.DiagnosticFailed, "Inactive"},
		{"no controller", func(machine *model.Machine) { machine.ControllerName = "" }, nil, dto.DiagnosticFailed, "no controller"},
		{"no engine", func(*model.Machine) {}, nil, dto.DiagnosticFailed, "no CNC engine"},
		{"no ip", func(machine *model.Machine) { machine.LiveIP = " " }, fanuc, dto.DiagnosticFailed, "no Live IP"},
		{"no port", func(machine *model.Machine) { machine.Port = 0 }, fanuc, dto.DiagnosticFailed, "no port"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			machine := ready
			testCase.change(&machine)
			expectDiagnostic(t, setupCheck(machine, testCase.engine), testCase.status, testCase.part)
		})
	}
}

func TestCollectorCheck(t *testing.T) {
	now := time.Now()
	machine := model.Machine{ID: 7, Status: "active"}
	fanuc := engineOf(model.EngineFanuc)
	running := &MachineDiagnosticsService{architecture: collectorArchitecture}
	fresh := &model.DiagnosticLiveStatus{UpdatedAt: now.Add(-5 * time.Second)}
	old := &model.DiagnosticLiveStatus{UpdatedAt: now.Add(-time.Hour)}
	recentEvent := []model.DiagnosticEvent{{At: now.Add(-time.Minute)}}
	oldEvent := []model.DiagnosticEvent{{At: now.Add(-time.Hour)}}

	wrongBuild := &MachineDiagnosticsService{architecture: "amd64"}
	expectDiagnostic(t, wrongBuild.collectorCheck(machine, fanuc, fresh, nil, now), dto.DiagnosticFailed, "amd64")

	leftOut := &MachineDiagnosticsService{architecture: collectorArchitecture,
		cncConfig: config.CNCConfig{Machines: config.MachineSelection{IDs: map[int64]bool{8: true}}}}
	expectDiagnostic(t, leftOut.collectorCheck(machine, fanuc, fresh, nil, now), dto.DiagnosticFailed, "CNC_MACHINE_IDS")

	expectDiagnostic(t, running.collectorCheck(machine, nil, nil, nil, now), dto.DiagnosticSkipped, "set-up")
	expectDiagnostic(t, running.collectorCheck(machine, fanuc, fresh, nil, now), dto.DiagnosticOK, "reporting")
	// A worker between connects writes no status but records its attempts
	expectDiagnostic(t, running.collectorCheck(machine, fanuc, old, recentEvent, now), dto.DiagnosticOK, "retrying")
	expectDiagnostic(t, running.collectorCheck(machine, fanuc, old, oldEvent, now), dto.DiagnosticWarning, "No sign")
	expectDiagnostic(t, running.collectorCheck(machine, fanuc, nil, nil, now), dto.DiagnosticWarning, "No sign")
}

func TestControllerCheck(t *testing.T) {
	now := time.Now()
	service := &MachineDiagnosticsService{cncConfig: config.CNCConfig{Quiet: 5 * time.Minute}}
	passed := dto.DiagnosticCheck{Status: dto.DiagnosticOK}
	failed := dto.DiagnosticCheck{Status: dto.DiagnosticFailed}
	fanuc, siemens := engineOf(model.EngineFanuc), engineOf(model.EngineSiemens)
	seenLongAgo := now.Add(-10 * time.Minute)
	connected := &model.DiagnosticLiveStatus{IsConnected: true, UpdatedAt: now, LastSeenAt: &seenLongAgo}
	notConnected := &model.DiagnosticLiveStatus{UpdatedAt: now}
	event := func(eventType, message string) []model.DiagnosticEvent {
		return []model.DiagnosticEvent{{Type: model.EventDisconnected}, {Type: eventType, Message: message, At: now}}
	}
	check := func(setup dto.DiagnosticCheck, engine *int, status *model.DiagnosticLiveStatus, events []model.DiagnosticEvent) dto.DiagnosticCheck {
		return service.controllerCheck(setup, passed, engine, status, events, isControllerConnected(status, now), now)
	}

	expectDiagnostic(t, check(failed, fanuc, notConnected, nil), dto.DiagnosticSkipped, "not collecting")
	// Silence while connected is normal: the DLL sends a reading only on a change
	expectDiagnostic(t, check(passed, fanuc, connected, nil), dto.DiagnosticOK, "only when something changes")
	dropping := *connected
	dropping.DroppedCount = 3
	expectDiagnostic(t, check(passed, fanuc, &dropping, nil), dto.DiagnosticWarning, "3 values could not be stored")
	// A row left by a stopped server still says connected, but nobody is reporting
	stopped := *connected
	stopped.UpdatedAt = now.Add(-time.Hour)
	expectDiagnostic(t, check(passed, fanuc, &stopped, nil), dto.DiagnosticFailed, "does not answer")

	expectDiagnostic(t, check(passed, fanuc, nil, nil), dto.DiagnosticFailed, "never reported")
	expectDiagnostic(t, check(passed, fanuc, notConnected, event(model.EventConnectFailed, "cnc: CNCConnect returned -1")), dto.DiagnosticFailed, "CNCConnect returned -1")
	expectDiagnostic(t, check(passed, fanuc, notConnected, event(model.EventConfigError, "no field")), dto.DiagnosticFailed, "no field")
	expectDiagnostic(t, check(passed, fanuc, notConnected, event(model.EventStale, "no usable reading")), dto.DiagnosticFailed, "FOCAS")
	expectDiagnostic(t, check(passed, siemens, notConnected, event(model.EventStale, "no usable reading")), dto.DiagnosticFailed, "OPC UA")
	certificate := &model.DiagnosticLiveStatus{UpdatedAt: now, LastError: "12: BadRequestNotAllowed"}
	expectDiagnostic(t, check(passed, siemens, certificate, nil), dto.DiagnosticFailed, "certificate")
}

func TestModeCheck(t *testing.T) {
	expectDiagnostic(t, modeCheck(nil, false), dto.DiagnosticFailed, "No mode")
	expectDiagnostic(t, modeCheck(&model.DiagnosticMode{ModeID: model.ModeOffline}, false), dto.DiagnosticFailed, "Offline")
	expectDiagnostic(t, modeCheck(&model.DiagnosticMode{ModeID: model.ModeOffline}, true), dto.DiagnosticWarning, "next reading")
	expectDiagnostic(t, modeCheck(&model.DiagnosticMode{ModeID: model.ModeIdle, ModeName: "Idle"}, true), dto.DiagnosticOK, "Program reset")
}

func TestVariablesCheck(t *testing.T) {
	row := func(value any) dto.DiagnosticVariableRow {
		return dto.DiagnosticVariableRow{Readable: true, Value: value}
	}
	unreadable := dto.DiagnosticVariableRow{Reason: "no address"}

	expectDiagnostic(t, variablesCheck([]dto.DiagnosticVariableRow{unreadable}, true), dto.DiagnosticFailed, "no variable")
	expectDiagnostic(t, variablesCheck([]dto.DiagnosticVariableRow{row(1.0)}, false), dto.DiagnosticSkipped, "once the controller is connected")
	expectDiagnostic(t, variablesCheck([]dto.DiagnosticVariableRow{row(nil)}, true), dto.DiagnosticWarning, "no reading has arrived")
	expectDiagnostic(t, variablesCheck([]dto.DiagnosticVariableRow{row(0.0), row("")}, true), dto.DiagnosticWarning, "0 or empty")
	expectDiagnostic(t, variablesCheck([]dto.DiagnosticVariableRow{row(0.0), row(41.5), row(nil), unreadable}, true), dto.DiagnosticOK, "2 of 3 readable")
}

func TestVariableUnreadableReason(t *testing.T) {
	fanuc := engineOf(model.EngineFanuc)
	readable := model.DiagnosticVariable{Enabled: true, Source: "snapshot", Address: "AlarmNo"}
	if reason := variableUnreadableReason(fanuc, readable); reason != "" {
		t.Errorf("readable variable: %q", reason)
	}
	cases := []struct {
		name     string
		engine   *int
		variable model.DiagnosticVariable
		part     string
	}{
		{"turned off", fanuc, model.DiagnosticVariable{Source: "snapshot", Address: "AlarmNo"}, "Turned off"},
		{"no engine", nil, model.DiagnosticVariable{Enabled: true}, "no CNC engine"},
		{"plc only", fanuc, model.DiagnosticVariable{Enabled: true, KnownSources: "plc"}, "Only has a plc address"},
		{"no address", fanuc, model.DiagnosticVariable{Enabled: true}, "No address"},
		{"override to nc", fanuc, model.DiagnosticVariable{Enabled: true, Source: "nc", Address: "/x"}, "read from nc"},
		{"no such field", fanuc, model.DiagnosticVariable{Enabled: true, Source: "snapshot", Address: "Nope"}, `no field "Nope"`},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if reason := variableUnreadableReason(testCase.engine, testCase.variable); !strings.Contains(reason, testCase.part) {
				t.Errorf("got %q; want %q", reason, testCase.part)
			}
		})
	}
}
