//go:build windows && 386

package cnc

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
)

// cncMasterSnapshot matches the DLL's C++ CncMasterSnapshot 1:1 (#pragma pack(8), x86):
// 34,440 bytes. Field names keep the DLL's spelling. Copied from
// docs/dll-usecase/G0_V0.0.4_DLL/sinumerik_client.go.
type cncMasterSnapshot struct {
	IsConnected       int16
	ResultStatus      int16
	_                 [4]byte
	SampleTimestampMs int64

	MachineMode int16 // 0=JOG 1=MDA 2=AUTO
	ProcStatus  int16 // 0=Reset 1=Stop 2=Active 3=Wait 4=Interrupted
	IsEmergency int16
	AxisNames   [5][4]byte
	_           [6]byte

	MachPos  [AxisCount]float64 // machine position (MCS)
	AbsPos   [AxisCount]float64 // work position (WCS)
	RelPos   [AxisCount]float64
	DistToGo [AxisCount]float64

	MachDec      [5]int16
	AbsDec       [5]int16
	RelDec       [5]int16
	DistToGoDec  [5]int16
	MachUnit     [5]int16
	AbsUnit      [5]int16
	RelUnit      [5]int16
	DistToGoUnit [5]int16

	LinearDec     int16
	RotaryDec     int16
	AbsLinearDec  int16
	AbsRotaryDec  int16
	LinearUnit    int16
	RotaryUnit    int16
	AbsLinearUnit int16
	AbsRotaryUnit int16

	SpindeOverrideActual     float64 // actual spindle rpm
	SpindleDec               int16
	SpindleLoadPercetnage    int16
	SpindeOverridePercentage int16
	_                        [2]byte

	FeedRateActual      float64
	FeedDec             int16
	FeedRatePercetnage  int16
	FeedUnitCode        int16
	AxisLoadXPercetnage int16
	AxisLoadYPercetnage int16
	AxisLoadZPercetnage int16

	ProgNumber        int32
	CurProgLine       int32
	TotalProgLines    int32
	ProgProgress      int32
	ProgramPartCount  int32
	RequiredPartCount int32

	ProgName       [64]byte
	ActiveProgPath [256]byte
	ProgList       [32768]byte

	TimerMinute int32
	TimerMsec   int32

	ToolNumber    int32
	ToolLength    float64
	ToolRadius    float64
	ToolOffsetDec int16

	ActiveGCode [16]byte
	GcodeGroup  [512]byte // the DLL puts the list of missing nodes here
	CmdFeedOvr  [16]byte
	CmdSpOvr    [16]byte
	_           [2]byte

	AlarmNo  int32
	AlarmMsg [128]byte
	OpMsgNo  int32
	OpMsg    [256]byte

	_ [4]byte
}

const snapshotSize = 34440

// maxSlots stays under Go's limit of about 2,000 callbacks per process (2 per slot).
const maxSlots = 990

// slot is one pair of DLL callbacks. The callback does not say which machine sent it, so each
// connected machine owns a slot; connection is nil while the slot is free.
type slot struct {
	connection    atomic.Pointer[Connection]
	dataCallback  uintptr
	errorCallback uintptr
}

type dllDriver struct {
	connectProc    *syscall.LazyProc
	disconnectProc *syscall.LazyProc

	mu        sync.Mutex
	freeSlots []*slot // FIFO: a freed slot is reused last, so late callbacks of the old machine die out first

	siemensMu sync.Mutex // the DLL reads the Siemens login from the process environment
}

// loadDLL checks the snapshot layout, loads the DLL from libDir, resolves CNCConnect /
// CNCdisconnect and creates the callback slots.
func loadDLL(cncConfig config.CNCConfig) (*dllDriver, error) {
	if size := unsafe.Sizeof(cncMasterSnapshot{}); size != snapshotSize {
		return nil, fmt.Errorf("cnc: snapshot struct is %d bytes, the DLL sends %d", size, snapshotSize)
	}

	libDir := findLibDir(cncConfig.LibDir, cncConfig.DLLName)
	setDllDirectory(libDir) // the DLL's own dependencies (open62541, libssl, Fwlib32, ...) are found there
	dll := syscall.NewLazyDLL(filepath.Join(libDir, cncConfig.DLLName))
	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("cnc: load %s from %s: %w", cncConfig.DLLName, libDir, err)
	}
	driver := &dllDriver{connectProc: dll.NewProc("CNCConnect"), disconnectProc: dll.NewProc("CNCdisconnect")}
	if err := driver.connectProc.Find(); err != nil {
		return nil, fmt.Errorf("cnc: %w", err)
	}
	if err := driver.disconnectProc.Find(); err != nil {
		return nil, fmt.Errorf("cnc: %w", err)
	}

	// Callbacks are created once for the whole process (Go cannot free them), so a fixed set
	// of slots is made here and reused by every connection.
	slotCount := min(cncConfig.Slots, maxSlots)
	for range slotCount {
		callbackSlot := &slot{}
		callbackSlot.dataCallback = syscall.NewCallback(func(snapshot uintptr) uintptr {
			callbackSlot.onData(snapshot)
			return 0
		})
		callbackSlot.errorCallback = syscall.NewCallback(func(code, message uintptr) uintptr {
			callbackSlot.onError(code, message)
			return 0
		})
		driver.freeSlots = append(driver.freeSlots, callbackSlot)
	}
	slog.Info("cnc dll loaded", "dll", cncConfig.DLLName, "dir", libDir, "slots", slotCount, "snapshotBytes", snapshotSize)
	return driver, nil
}

// connectMachine takes a free callback slot, binds a new Connection to it and calls
// CNCConnect. On failure the slot is freed at once; on success it is freed by Disconnect.
func (driver *dllDriver) connectMachine(machine model.MonitoredMachine) (*Connection, int, error) {
	driver.mu.Lock()
	if len(driver.freeSlots) == 0 {
		driver.mu.Unlock()
		return nil, 0, ErrNoFreeSlot
	}
	callbackSlot := driver.freeSlots[0]
	driver.freeSlots = driver.freeSlots[1:]
	driver.mu.Unlock()

	// Bound before CNCConnect: the DLL thread may call back before CNCConnect returns.
	connection := newConnection()
	callbackSlot.connection.Store(connection)

	id, err := driver.callCNCConnect(machine, callbackSlot)
	if err != nil {
		callbackSlot.connection.Store(nil)
		driver.releaseSlot(callbackSlot)
		return nil, id, err
	}
	connection.disconnect = func() {
		callbackSlot.connection.Store(nil) // late callbacks are ignored from here on
		driver.disconnectProc.Call(uintptr(id))
		driver.releaseSlot(callbackSlot)
	}
	return connection, id, nil
}

// callCNCConnect calls CNCConnect, which starts the DLL's polling thread for the machine; a
// negative result is a failed connect. Siemens connects run one at a time: the DLL reads the login
// from SIEMENS_USER / SIEMENS_PASS in the process environment.
func (driver *dllDriver) callCNCConnect(machine model.MonitoredMachine, callbackSlot *slot) (int, error) {
	engine := *machine.CNCEngine
	if engine == model.EngineSiemens {
		driver.siemensMu.Lock()
		defer driver.siemensMu.Unlock()
		os.Setenv("SIEMENS_USER", machine.Username)
		os.Setenv("SIEMENS_PASS", machine.Password)
	}

	ip, err := syscall.BytePtrFromString(machine.LiveIP)
	if err != nil {
		return 0, fmt.Errorf("cnc: invalid ip %q", machine.LiveIP)
	}
	result, _, _ := driver.connectProc.Call(uintptr(unsafe.Pointer(ip)), uintptr(machine.Port),
		callbackSlot.dataCallback, callbackSlot.errorCallback, uintptr(engine))
	runtime.KeepAlive(ip)

	id := int(int32(result))
	if id < 0 {
		return id, fmt.Errorf("cnc: CNCConnect returned %d", id)
	}
	return id, nil
}

// releaseSlot puts callbackSlot back at the end of the free list.
func (driver *dllDriver) releaseSlot(callbackSlot *slot) {
	driver.mu.Lock()
	driver.freeSlots = append(driver.freeSlots, callbackSlot)
	driver.mu.Unlock()
}

// onData runs on the DLL thread for every poll. It reads only the fields the collector uses
// straight from the DLL's buffer (never the 32 KB ProgList) and must not block. A panic is
// swallowed: it would otherwise unwind through the DLL and end the process.
func (callbackSlot *slot) onData(pointer uintptr) {
	// DATA_COLLECTION_STEP_13: Receive a snapshot from the DLL thread and publish it to the worker.
	defer func() { _ = recover() }()
	connection := callbackSlot.connection.Load()
	if connection == nil || pointer == 0 {
		return
	}
	snapshot := (*cncMasterSnapshot)(unsafe.Pointer(pointer))

	reading := Reading{
		ReceivedAt:   time.Now(),
		Connected:    snapshot.IsConnected != 0,
		MissingNodes: cstr(snapshot.GcodeGroup[:]),
	}
	if snapshot.SampleTimestampMs > 0 {
		reading.SampleAt = time.UnixMilli(snapshot.SampleTimestampMs)
	}
	numbers := [FieldCount]float64{
		FieldMachineMode:       float64(snapshot.MachineMode),
		FieldProcStatus:        float64(snapshot.ProcStatus),
		FieldIsEmergency:       float64(snapshot.IsEmergency),
		FieldCurProgLine:       float64(snapshot.CurProgLine),
		FieldTotalProgLines:    float64(snapshot.TotalProgLines),
		FieldProgramPartCount:  float64(snapshot.ProgramPartCount),
		FieldRequiredPartCount: float64(snapshot.RequiredPartCount),
		FieldSpindleSpeed:      snapshot.SpindeOverrideActual,
		FieldSpindleLoad:       float64(snapshot.SpindleLoadPercetnage),
		FieldSpindleOverride:   float64(snapshot.SpindeOverridePercentage),
		FieldFeedRate:          snapshot.FeedRateActual,
		FieldFeedOverride:      float64(snapshot.FeedRatePercetnage),
		FieldAxisLoadX:         float64(snapshot.AxisLoadXPercetnage),
		FieldAxisLoadY:         float64(snapshot.AxisLoadYPercetnage),
		FieldAxisLoadZ:         float64(snapshot.AxisLoadZPercetnage),
		FieldToolNumber:        float64(snapshot.ToolNumber),
		FieldAlarmNo:           float64(snapshot.AlarmNo),
	}
	for field, number := range numbers {
		reading.Values[field].Number = number
	}
	for axis := range AxisCount {
		reading.Values[FieldMachPos+axis].Number = snapshot.MachPos[axis]
		reading.Values[FieldAbsPos+axis].Number = snapshot.AbsPos[axis]
		reading.Values[FieldRelPos+axis].Number = snapshot.RelPos[axis]
		reading.Values[FieldDistToGo+axis].Number = snapshot.DistToGo[axis]
	}
	reading.Values[FieldProgName].Text = cstr(snapshot.ProgName[:])
	reading.Values[FieldActiveProgPath].Text = cstr(snapshot.ActiveProgPath[:])
	reading.Values[FieldAlarmMsg].Text = cstr(snapshot.AlarmMsg[:])

	connection.publishReading(reading)
}

// onError runs on the DLL thread: void OnErr(int32_t code, const char* msg).
func (callbackSlot *slot) onError(code, message uintptr) {
	defer func() { _ = recover() }()
	if connection := callbackSlot.connection.Load(); connection != nil {
		connection.publishError(int(int32(code)), cstrAt(message))
	}
}

// cstr reads a NUL-terminated C string from a byte buffer.
func cstr(buffer []byte) string {
	text := string(buffer)
	if end := strings.IndexByte(text, 0); end >= 0 {
		text = text[:end]
	}
	return strings.TrimSpace(text)
}

// cstrAt reads a NUL-terminated C string (at most 512 bytes) at address.
func cstrAt(address uintptr) string {
	if address == 0 {
		return ""
	}
	var text []byte
	for offset := range 512 {
		character := *(*byte)(unsafe.Add(unsafe.Pointer(address), offset))
		if character == 0 {
			break
		}
		text = append(text, character)
	}
	return string(text)
}

// findLibDir returns the folder holding the DLL: next to the binary first, then the working
// directory. An absolute libDir is used as it is.
func findLibDir(libDir, dllName string) string {
	if filepath.IsAbs(libDir) {
		return libDir
	}
	var candidates []string
	if executable, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(executable), libDir))
	}
	if workingDir, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(workingDir, libDir))
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(filepath.Join(candidate, dllName)); err == nil && !info.IsDir() {
			return candidate
		}
	}
	absolute, _ := filepath.Abs(libDir)
	return absolute
}

func setDllDirectory(path string) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	pathPointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	kernel32.NewProc("SetDllDirectoryW").Call(uintptr(unsafe.Pointer(pathPointer)))
}
