package monitor

import (
	"testing"

	"7th-sense-server/internal/model"
)

func TestModeFor(t *testing.T) {
	tests := []struct {
		name      string
		connected bool
		values    map[int]float64
		want      int
	}{
		{"offline", false, map[int]float64{variableOpMode: 2, variableProgStatus: 2}, model.ModeOffline},
		{"emergency", true, map[int]float64{variableEmergency: 1, variableOpMode: 2, variableProgStatus: 2}, model.ModeDowntime},
		{"alarm", true, map[int]float64{variableAlarm: 700115, variableOpMode: 2, variableProgStatus: 1}, model.ModeDowntime},
		{"jog", true, map[int]float64{variableOpMode: 0, variableProgStatus: 0}, model.ModeSetup},
		{"mda", true, map[int]float64{variableOpMode: 1, variableProgStatus: 2}, model.ModeSetup},
		{"auto running", true, map[int]float64{variableOpMode: 2, variableProgStatus: 2}, model.ModeInProduction},
		{"auto running rapid", true, map[int]float64{variableOpMode: 2, variableProgStatus: 2, variableRapidMode: 1}, model.ModeRapid},
		{"auto stopped", true, map[int]float64{variableOpMode: 2, variableProgStatus: 1}, model.ModeHold},
		{"auto waiting", true, map[int]float64{variableOpMode: 2, variableProgStatus: 3}, model.ModeHold},
		{"auto interrupted", true, map[int]float64{variableOpMode: 2, variableProgStatus: 4}, model.ModeHold},
		{"auto reset", true, map[int]float64{variableOpMode: 2, variableProgStatus: 0}, model.ModeIdle},
		{"no opMode, running", true, map[int]float64{variableProgStatus: 2}, model.ModeInProduction},
		{"nothing known", true, map[int]float64{}, model.ModeIdle},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := modeFor(test.connected, test.values); got != test.want {
				t.Errorf("modeFor = %d, want %d", got, test.want)
			}
		})
	}
}
