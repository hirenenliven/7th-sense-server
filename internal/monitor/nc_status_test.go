package monitor

import (
	"testing"

	"7th-sense-server/internal/model"
)

// The same raw value means different things per engine (vendor sample modeNames / statusNames).
func TestNCStatusNamesPerEngine(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Fanuc opMode 1", opModeName(model.EngineFanuc, 1), "AUTO"},
		{"Siemens opMode 1", opModeName(model.EngineSiemens, 1), "MDA"},
		{"Siemens opMode 2", opModeName(model.EngineSiemens, 2), "AUTO"},
		{"Fanuc opMode 9", opModeName(model.EngineFanuc, 9), "REF"},
		{"Fanuc procStatus 2", procStatusName(model.EngineFanuc, 2), "HOLD"},
		{"Fanuc procStatus 3", procStatusName(model.EngineFanuc, 3), "START"},
		{"Siemens procStatus 2", procStatusName(model.EngineSiemens, 2), "Active"},
		{"Siemens procStatus 4", procStatusName(model.EngineSiemens, 4), "Interrupted"},
		{"unknown value", procStatusName(model.EngineFanuc, 7), "unknown(7)"},
		{"no engine", opModeName(engineNone, 2), "unknown(2)"},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("%s = %q, want %q", test.name, test.got, test.want)
		}
	}
}
