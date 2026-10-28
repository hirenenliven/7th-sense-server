package monitor

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
)

func TestWriteLog(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil)))
	defer slog.SetDefault(previous)

	settings := config.CNCLogSettings{Events: false, Modes: true, ModeChecks: false, Variables: false, Other: false}
	machine := &model.MonitoredMachine{ID: 225, Name: "VM14"}

	writeLog(settings, logTypeEvent, slog.LevelInfo, machine, "event switched off")
	writeLog(settings, logTypeVariable, slog.LevelWarn, machine, "variables switched off")
	writeLog(settings, logTypeEvent, slog.LevelError, machine, "event error")
	writeLog(settings, logTypeMode, slog.LevelInfo, machine, "mode on")
	writeLog(settings, logTypeModeCheck, slog.LevelInfo, machine, "mode check switched off")
	writeLog(settings, logTypeOther, slog.LevelInfo, nil, "other switched off")
	writeLog(settings, logTypeOther, slog.LevelError, nil, "other error")

	lines := output.String()
	for _, want := range []string{
		`msg="[machine 225 VM14] event error"`,
		`msg="[machine 225 VM14] mode on" logType=mode machineId=225 machine=VM14`,
		`msg="other error" logType=other`,
	} {
		if !strings.Contains(lines, want) {
			t.Errorf("missing %s in:\n%s", want, lines)
		}
	}
	for _, unwanted := range []string{"event switched off", "variables switched off", "mode check switched off",
		"other switched off"} {
		if strings.Contains(lines, unwanted) {
			t.Errorf("unexpected %q in:\n%s", unwanted, lines)
		}
	}
}
