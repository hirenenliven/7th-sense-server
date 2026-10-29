//go:build windows && 386

package cnc

import (
	"testing"
	"unsafe"
)

// The Go struct must match the DLL's C++ struct byte for byte.
func TestSnapshotSize(t *testing.T) {
	if size := unsafe.Sizeof(cncMasterSnapshot{}); size != snapshotSize {
		t.Fatalf("cncMasterSnapshot is %d bytes, the DLL sends %d", size, snapshotSize)
	}
}

// The fields the machine mode is decided from must sit where the DLL writes them
// (MSVC #pragma pack(8), x86). A field moved but compensated by padding keeps the size right
// and reads zeros, e.g. AlarmNo 0 and no Downtime.
func TestSnapshotModeFieldOffsets(t *testing.T) {
	var snapshot cncMasterSnapshot
	offsets := []struct {
		field string
		got   uintptr
		want  uintptr
	}{
		{"IsConnected", unsafe.Offsetof(snapshot.IsConnected), 0},
		{"MachineMode", unsafe.Offsetof(snapshot.MachineMode), 16},
		{"ProcStatus", unsafe.Offsetof(snapshot.ProcStatus), 18},
		{"IsEmergency", unsafe.Offsetof(snapshot.IsEmergency), 20},
		{"MachPos", unsafe.Offsetof(snapshot.MachPos), 48},
		{"AlarmNo", unsafe.Offsetof(snapshot.AlarmNo), 34044},
		{"AlarmMsg", unsafe.Offsetof(snapshot.AlarmMsg), 34048},
		{"OpMsgNo", unsafe.Offsetof(snapshot.OpMsgNo), 34176},
		{"OpMsg", unsafe.Offsetof(snapshot.OpMsg), 34180},
	}
	for _, offset := range offsets {
		if offset.got != offset.want {
			t.Errorf("%s is at byte %d, the DLL writes it at %d", offset.field, offset.got, offset.want)
		}
	}
}
