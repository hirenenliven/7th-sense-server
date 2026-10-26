package cnc

import "testing"

// Migration 029 stores axis addresses as "<array>[<index>]"; each must find its own field.
func TestFieldIndexAxisPositions(t *testing.T) {
	tests := []struct {
		address string
		want    int
	}{
		{"MachPos[0]", FieldMachPos},
		{"MachPos[4]", FieldMachPos + 4},
		{"AbsPos[0]", FieldAbsPos},
		{"RelPos[2]", FieldRelPos + 2},
		{"DistToGo[4]", FieldDistToGo + 4},
		{"AlarmMsg", FieldAlarmMsg},
	}
	for _, test := range tests {
		if got, ok := FieldIndex(test.address); !ok || got != test.want {
			t.Errorf("FieldIndex(%q) = %d, %v; want %d, true", test.address, got, ok, test.want)
		}
	}
	if _, ok := FieldIndex("MachPos[5]"); ok {
		t.Errorf("FieldIndex(%q) found a field; the snapshot has %d axes", "MachPos[5]", AxisCount)
	}
}
