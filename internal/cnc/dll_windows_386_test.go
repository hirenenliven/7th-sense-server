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
