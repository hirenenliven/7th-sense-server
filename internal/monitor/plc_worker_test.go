package monitor

import (
	"testing"

	"7th-sense-server/internal/opcua"
)

// TestPolledChanges checks that reading the nodes the controller had no room to monitor hands on
// only what changed since the previous read, like a subscription.
func TestPolledChanges(t *testing.T) {
	handles := []uint32{10, 11}
	last := map[uint32]opcua.Change{}
	read := func(pressure float64, temperatureStatus string) []opcua.Change {
		return polledChanges(handles, []opcua.ReadResult{
			{Value: opcua.Value{Number: pressure}},
			{Status: temperatureStatus, Value: opcua.Value{Number: 40}},
		}, last)
	}

	first := read(21.5, "")
	if len(first) != 2 || !first[0].Good || first[0].Value.Number != 21.5 || first[1].Value.Number != 40 {
		t.Fatalf("first read: got %+v, want both values", first)
	}
	if again := read(21.5, ""); len(again) != 0 {
		t.Errorf("unchanged read: got %+v, want none", again)
	}
	changed := read(22, "")
	if len(changed) != 1 || changed[0].Handle != 10 || changed[0].Value.Number != 22 {
		t.Errorf("pressure changed: got %+v, want handle 10 = 22", changed)
	}
	bad := read(22, "StatusBadNotReadable")
	if len(bad) != 1 || bad[0].Handle != 11 || bad[0].Good || bad[0].Status != "StatusBadNotReadable" {
		t.Errorf("temperature bad: got %+v, want handle 11 not good", bad)
	}
	if back := read(22, ""); len(back) != 1 || back[0].Handle != 11 || !back[0].Good {
		t.Errorf("temperature good again: got %+v", back)
	}
	if extra := polledChanges(handles[:1], []opcua.ReadResult{{}, {}}, map[uint32]opcua.Change{}); len(extra) != 1 {
		t.Errorf("more results than handles: got %+v, want 1", extra)
	}
}
