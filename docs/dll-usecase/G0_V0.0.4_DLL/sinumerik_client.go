package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// ===================== SETTINGS =====================

const (
	plcIP       = "157.10.57.69"
	plcPort     = 4840
	siemensUser = "JYOTI"
	siemensPass = "12345678"
	endpoint    = "opc.tcp://157.10.57.69:4840"

	verbose            = false
	emergencyAvailable = false
)

// ===================== DLL BINDINGS =====================

const (
	dllDir  = "Lib"
	dllName = "FanucCncCommunication.dll"
)

var (
	libPath        string
	mod            *syscall.LazyDLL
	procConnect    *syscall.LazyProc
	procDisconnect *syscall.LazyProc
)

func setDllDirectory(path string) {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	p, _ := syscall.UTF16PtrFromString(path)
	k32.NewProc("SetDllDirectoryW").Call(uintptr(unsafe.Pointer(p)))
}

// findLibDir returns the absolute Lib folder: exe dir first (built exe),
func findLibDir() string {
	var cands []string
	if exe, err := os.Executable(); err == nil {
		cands = append(cands, filepath.Join(filepath.Dir(exe), dllDir))
	}
	if wd, err := os.Getwd(); err == nil {
		cands = append(cands, filepath.Join(wd, dllDir))
	}
	for _, c := range cands {
		if fi, err := os.Stat(filepath.Join(c, dllName)); err == nil && !fi.IsDir() {
			return c
		}
	}
	abs, _ := filepath.Abs(dllDir)
	return abs
}

// initDLL must be called once at the start of main().
func initDLL() {
	libPath = findLibDir()
	setDllDirectory(libPath)
	mod = syscall.NewLazyDLL(filepath.Join(libPath, dllName))
	procConnect = mod.NewProc("CNCConnect")
	procDisconnect = mod.NewProc("CNCdisconnect")
}

// ===================== SNAPSHOT STRUCT =====================
// Must match C++ CncMasterSnapshot 1:1 (#pragma pack(8), x86) = 34440 bytes.

type CncMasterSnapshot struct {
	IsConnected       int16
	ResultStatus      int16
	_                 [4]byte
	SampleTimestampMs int64

	MachineMode int16 // 0=JOG 1=MDA 2=AUTO
	ProcStatus  int16 // 0=Reset 1=Stop 2=Active 3=Wait 4=Interrupted
	IsEmergency int16
	AxisNames   [5][4]byte
	_           [6]byte

	MachPos  [5]float64 // MCS (mm)
	AbsPos   [5]float64 // WCS (inch)
	RelPos   [5]float64
	DistToGo [5]float64

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
	GcodeGroup  [512]byte // DLL reuses this field for missing-node list
	CmdFeedOvr  [16]byte
	CmdSpOvr    [16]byte
	_           [2]byte

	AlarmNo  int32
	AlarmMsg [128]byte
	OpMsgNo  int32
	OpMsg    [256]byte

	_ [4]byte
}

const expectSize = 34440

// ===================== HELPERS =====================

var modeNames = map[int16]string{0: "JOG", 1: "MDI", 2: "AUTO"}
var progNames = map[int16]string{0: "Reset", 1: "Stop", 2: "Active", 3: "Waiting", 4: "Interrupted"}

// logf prints only when verbose = true.
func logf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

// cstr reads a NUL-terminated C string from a byte buffer.
func cstr(b []byte) string {
	s := string(b)
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return s
}

func nameOr(m map[int16]string, n int16) string {
	if s, ok := m[n]; ok {
		return s
	}
	return fmt.Sprintf("%d", n)
}

func isRotary(name string) bool {
	switch strings.ToUpper(name) {
	case "A", "B", "C":
		return true
	}
	return false
}

// ===================== ANSI COLOURS =====================

const (
	cReset  = "\033[0m"
	cBold   = "\033[1m"
	cDim    = "\033[90m"
	cRed    = "\033[91m"
	cGreen  = "\033[92m"
	cYellow = "\033[93m"
	cBlue   = "\033[94m"
	cCyan   = "\033[96m"
	cWhite  = "\033[97m"
)

func paint(color, s string) string { return color + s + cReset }
func lbl(s string) string          { return paint(cDim, s) }

// pv pads to width; grey for "n/a"/"-", white otherwise.
func pv(width int, s string) string {
	p := fmt.Sprintf("%-*s", width, s)
	if s == "n/a" || s == "-" {
		return paint(cDim, p)
	}
	return paint(cBold+cWhite, p)
}

func axisColor(name string) string {
	switch strings.ToUpper(name) {
	case "X":
		return cCyan
	case "Y":
		return cGreen
	case "Z":
		return cYellow
	}
	return cWhite
}

// enableVT turns on ANSI escape support in the Windows console.
func enableVT() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	getMode := k32.NewProc("GetConsoleMode")
	setMode := k32.NewProc("SetConsoleMode")
	h := os.Stdout.Fd()
	var m uint32
	if r, _, _ := getMode.Call(h, uintptr(unsafe.Pointer(&m))); r == 0 {
		return
	}
	setMode.Call(h, uintptr(m|0x0004))
}

// ===================== DASHBOARD RENDER =====================

func render(d *CncMasterSnapshot) string {
	var b strings.Builder
	line := func(format string, a ...interface{}) {
		fmt.Fprintf(&b, format+"\033[K\n", a...)
	}
	bar := paint(cDim, strings.Repeat("-", 78))
	sec := func(name string) string { return paint(cBold+cBlue, fmt.Sprintf("%-8s", name)) }

	// header
	ts := time.Now()
	if d.SampleTimestampMs > 0 {
		ts = time.UnixMilli(d.SampleTimestampMs)
	}
	title := fmt.Sprintf(" SINUMERIK LIVE DATA   %s   %s", endpoint, ts.Format("15:04:05"))
	line("%s", "\033[1;97;44m"+fmt.Sprintf("%-78s", title)+cReset)

	// mode / status / emergency / connected
	modeCol := cWhite
	switch d.MachineMode {
	case 0:
		modeCol = cCyan
	case 1:
		modeCol = cYellow
	case 2:
		modeCol = cGreen
	}
	emerg := paint(cDim, "n/a")
	if emergencyAvailable {
		if d.IsEmergency != 0 {
			emerg = paint(cBold+cRed, "ACTIVE")
		} else {
			emerg = paint(cGreen, "OK")
		}
	}
	conn := paint(cBold+cRed, "NO")
	if d.IsConnected == 1 {
		conn = paint(cBold+cGreen, "YES")
	}
	line("%s %s   %s %s   %s %s   %s %s",
		lbl("MODE"), paint(cBold+modeCol, fmt.Sprintf("%-6s", nameOr(modeNames, d.MachineMode))),
		lbl("STATUS"), pv(8, nameOr(progNames, d.ProcStatus)),
		lbl("EMERGENCY"), emerg,
		lbl("CONNECTED"), conn)
	line("%s", bar)

	// axis table
	line("%s", paint(cBold+cWhite, fmt.Sprintf("%-5s %13s %13s %13s %13s %13s",
		"AXIS", "MCS (mm)", "WCS (inch)", "WCS (mm)", "DIST-TO-GO", "RELATIVE")))
	for i := 0; i < 5; i++ {
		name := cstr(d.AxisNames[i][:])
		if name == "" {
			continue
		}
		var row string
		if isRotary(name) {
			row = fmt.Sprintf("%-5s %13s %13s %13s %13s %13s", name,
				fmt.Sprintf("%.3f°", d.MachPos[i]),
				fmt.Sprintf("%.3f°", d.AbsPos[i]),
				"-",
				fmt.Sprintf("%.3f°", d.DistToGo[i]),
				fmt.Sprintf("%.3f°", d.RelPos[i]))
		} else {
			row = fmt.Sprintf("%-5s %13s %13s %13s %13s %13s", name,
				fmt.Sprintf("%.3f", d.MachPos[i]),
				fmt.Sprintf("%.4f", d.AbsPos[i]),
				fmt.Sprintf("%.3f", d.AbsPos[i]*25.4), // inch -> mm
				fmt.Sprintf("%.4f", d.DistToGo[i]),
				fmt.Sprintf("%.4f", d.RelPos[i]))
		}
		line("%s", cBold+axisColor(name)+row+cReset)
	}
	line("%s", lbl("(MCS: mm | WCS: Inch & converted mm | Rotary: degrees)"))
	line("%s", bar)

	// spindle / feed / servo / tool
	line("%s %s %s %s   %s %s %s   %s %s %s   %s %s", sec("SPINDLE"),
		lbl("act"), pv(5, fmt.Sprintf("%.0f", d.SpindeOverrideActual)), lbl("rpm"),
		lbl("load"), pv(3, fmt.Sprintf("%d", d.SpindleLoadPercetnage)), lbl("%"),
		lbl("ovr"), pv(3, fmt.Sprintf("%d", d.SpindeOverridePercentage)), lbl("%"),
		lbl("cmd"), pv(0, cstr(d.CmdSpOvr[:])))
	line("%s %s %s     %s %s %s   %s %s", sec("FEED"),
		lbl("act"), pv(7, fmt.Sprintf("%.1f", d.FeedRateActual)),
		lbl("ovr"), pv(3, fmt.Sprintf("%d", d.FeedRatePercetnage)), lbl("%"),
		lbl("cmd"), pv(0, cstr(d.CmdFeedOvr[:])))
	line("%s %s X %s %s   Y %s %s   Z %s %s", sec("SERVO"), lbl("load"),
		pv(3, fmt.Sprintf("%d", d.AxisLoadXPercetnage)), lbl("%"),
		pv(3, fmt.Sprintf("%d", d.AxisLoadYPercetnage)), lbl("%"),
		pv(3, fmt.Sprintf("%d", d.AxisLoadZPercetnage)), lbl("%"))
	line("%s %s %s", sec("TOOL"), lbl("number"), pv(0, fmt.Sprintf("%d", d.ToolNumber)))
	line("%s", bar)

	// program info
	line("%s %s %s", sec("PROGRAM"), lbl("name"), pv(0, cstr(d.ProgName[:])))
	line("%s %s %s", sec(""), lbl("path"), pv(0, cstr(d.ActiveProgPath[:])))
	line("%s %s %s", sec(""), lbl("line"), pv(0, fmt.Sprintf("%d", d.CurProgLine)))
	line("%s %s %s", sec("G-CODE"), lbl("active"), pv(0, cstr(d.ActiveGCode[:])))
	line("%s", bar)
	line("%s", lbl("Press Ctrl+C to stop"))
	return b.String()
}

// ===================== CHECK MODE =====================
// go run sinumerik_client.go check  -> lists which nodes exist on this machine.

type checkItem struct {
	label string
	path  string
	state string // CURRENT, ZERO, EMPTY, MISSING
	val   string
}

func runCheck() {
	// wait up to 10 s for first data
	var snap CncMasterSnapshot
	got := false
	for i := 0; i < 100; i++ {
		mu.Lock()
		snap, got = latest, have
		mu.Unlock()
		if got && snap.IsConnected == 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !got || snap.IsConnected != 1 {
		fmt.Println("No data received from the DLL (not connected).")
		return
	}
	time.Sleep(500 * time.Millisecond)
	mu.Lock()
	snap = latest
	mu.Unlock()
	d := &snap

	// missing node paths sent by the DLL
	missing := map[string]bool{}
	for _, p := range strings.Split(cstr(d.GcodeGroup[:]), ";") {
		if p != "" {
			missing[p] = true
		}
	}

	num := func(label, path string, v float64) checkItem {
		if v == 0 {
			return checkItem{label, path, "ZERO", "0"}
		}
		return checkItem{label, path, "CURRENT", fmt.Sprintf("%v", v)}
	}
	str := func(label, path string, b []byte) checkItem {
		s := cstr(b)
		if s == "" {
			return checkItem{label, path, "EMPTY", ""}
		}
		return checkItem{label, path, "CURRENT", s}
	}

	items := []checkItem{
		{"MachineMode", "/Bag/State/opMode[u1]", "CURRENT", nameOr(modeNames, d.MachineMode)},
		{"ProcStatus", "/Channel/State/acProg[u1,1]", "CURRENT", nameOr(progNames, d.ProcStatus)},
		num("IsEmergency", "/Plc/DB2700.DBX0.1", float64(d.IsEmergency)),
		num("SpindleActual (rpm)", "/Channel/Spindle/actSpeed[u1,1]", d.SpindeOverrideActual),
		num("SpindleLoad %", "/Channel/MachineAxis/aaLoad[u1,4]", float64(d.SpindleLoadPercetnage)),
		num("SpindleOverride %", "/Channel/Spindle/speedOvr[u1,1]", float64(d.SpindeOverridePercentage)),
		str("CmdSpindle", "/Channel/Spindle/cmdSpeed[u1,1]", d.CmdSpOvr[:]),
		num("FeedActual", "/Channel/State/actFeedRateIpo[u1]", d.FeedRateActual),
		num("FeedOverride %", "/Channel/MachineAxis/feedRateOvr[u1,1]", float64(d.FeedRatePercetnage)),
		str("CmdFeed", "/Channel/State/cmdFeedRateIpo[u1,1]", d.CmdFeedOvr[:]),
		num("ServoLoadX", "/Channel/MachineAxis/aaLoad[u1,1]", float64(d.AxisLoadXPercetnage)),
		num("ServoLoadY", "/Channel/MachineAxis/aaLoad[u1,2]", float64(d.AxisLoadYPercetnage)),
		num("ServoLoadZ", "/Channel/MachineAxis/aaLoad[u1,3]", float64(d.AxisLoadZPercetnage)),
		str("ProgName", "/Channel/ProgramInfo/progName[u1,1]", d.ProgName[:]),
		str("ActiveProgPath", "/Channel/ProgramInfo/workPandProgName[u1,1]", d.ActiveProgPath[:]),
		num("CurProgLine", "/Channel/ProgramInfo/actLineNumber", float64(d.CurProgLine)),
		num("PartCount", "/Channel/State/actParts", float64(d.ProgramPartCount)),
		num("ToolNumber", "/Channel/State/actTNumber[u1,1]", float64(d.ToolNumber)),
		str("ActiveGCode", "/Channel/SelectedFunctions/ncFktAct[u1,8]", d.ActiveGCode[:]),
	}

	// axis nodes (indices match the C++ node table)
	axisIdx := []int{1, 2, 3, 5, 6}
	for i := 0; i < 5; i++ {
		n := cstr(d.AxisNames[i][:])
		ix := fmt.Sprintf("[u1,%d]", axisIdx[i])
		items = append(items,
			num(fmt.Sprintf("MachPos[%d] %s", i, n), "/Channel/MachineAxis/actToolBasePos"+ix, d.MachPos[i]),
			num(fmt.Sprintf("AbsPos[%d] %s", i, n), "/Channel/GeometricAxis/progREPOS"+ix, d.AbsPos[i]),
			num(fmt.Sprintf("RelPos[%d] %s", i, n), "/Channel/GeometricAxis/aaPcsRel"+ix, d.RelPos[i]),
			num(fmt.Sprintf("DistToGo[%d] %s", i, n), "/Channel/GeometricAxis/progDistToGo"+ix, d.DistToGo[i]),
		)
	}

	// DLL-reported missing nodes override state
	for i := range items {
		if missing[items[i].path] {
			items[i].state = "MISSING"
			items[i].val = ""
		}
	}

	section := func(title, state string) int {
		fmt.Printf("\n===== %s =====\n", title)
		count := 0
		for _, it := range items {
			if it.state != state {
				continue
			}
			count++
			if it.val != "" {
				fmt.Printf("  %-26s = %s\n", it.label, it.val)
			} else {
				fmt.Printf("  %-26s %s\n", it.label, it.path)
			}
		}
		if count == 0 {
			fmt.Println("  (none)")
		}
		return count
	}

	m := section("MISSING (node not on this machine)", "MISSING")
	c := section("CURRENT (node available, data is coming)", "CURRENT")
	z := section("ZERO (node available, value is 0 right now)", "ZERO")
	e := section("EMPTY (node available, text is empty)", "EMPTY")

	fmt.Printf("\nTotal %d | missing %d | current %d | zero %d | empty %d\n", len(items), m, c, z, e)
}

// ===================== CALLBACKS (stdcall) =====================

var (
	mu     sync.Mutex
	latest CncMasterSnapshot // newest snapshot
	have   bool
)

// newData: callback -> main goroutine signal.
// Buffer 1 + non-blocking send = the DLL thread never blocks.
var newData = make(chan struct{}, 1)

// onData stores the snapshot copy and pings the renderer. Do NOT render here.
func onData(s uintptr) uintptr {
	if s == 0 {
		return 0
	}
	mu.Lock()
	latest = *(*CncMasterSnapshot)(unsafe.Pointer(s))
	have = true
	mu.Unlock()

	select {
	case newData <- struct{}{}:
	default: // renderer busy -> skip ping, no block
	}
	return 0
}

func onErr(code uintptr, msg uintptr) uintptr {
	fmt.Printf("DLL error %d\n", code)
	return 0
}

var (
	dataCb = syscall.NewCallback(onData)
	errCb  = syscall.NewCallback(onErr)
)

// ===================== MAIN =====================

func main() {
	// DLL reads credentials from env — set before CNCConnect
	os.Setenv("SIEMENS_USER", siemensUser)
	os.Setenv("SIEMENS_PASS", siemensPass)

	// struct layout guard
	if unsafe.Sizeof(CncMasterSnapshot{}) != expectSize {
		fmt.Printf("Struct layout mismatch: Go=%d, C++=%d. Aborting.\n",
			unsafe.Sizeof(CncMasterSnapshot{}), expectSize)
		return
	}
	logf("sizeof(CncMasterSnapshot) = %d (expected %d)\n",
		unsafe.Sizeof(CncMasterSnapshot{}), expectSize)

	// load DLL (from Lib folder + register it so its dependencies resolve)
	initDLL()
	logf("DLL folder: %s\n", libPath)
	if err := mod.Load(); err != nil {
		fmt.Println("DLL load failed:", err)
		fmt.Println("Expected", dllName, "in:", libPath)
		return
	}

	// connect (machineType 2 = Siemens)
	ipPtr, _ := syscall.BytePtrFromString(plcIP)
	logf("Connecting via DLL (Siemens, %s) ...\n", endpoint)
	ret, _, _ := procConnect.Call(
		uintptr(unsafe.Pointer(ipPtr)),
		plcPort,
		dataCb,
		errCb,
		2,
	)
	id := int32(ret)
	if id < 0 {
		fmt.Println("CNCConnect failed:", id)
		return
	}
	logf("Connected, machineId = %d\n", id)

	// one-shot check mode
	if len(os.Args) > 1 && os.Args[1] == "check" {
		runCheck()
		procDisconnect.Call(uintptr(id))
		return
	}

	// live dashboard: event-driven, no timer
	enableVT()
	fmt.Print("\033[2J\033[?25l") // clear screen, hide cursor
	defer fmt.Print("\033[?25h")  // restore cursor

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sig: // Ctrl+C -> disconnect
			procDisconnect.Call(uintptr(id))
			fmt.Print("\033[2J\033[H")
			fmt.Println("Disconnected.")
			return

		case <-newData: // new data -> render instantly
			mu.Lock()
			snap := latest
			ok := have
			mu.Unlock()
			if !ok {
				continue
			}
			fmt.Print("\033[H" + render(&snap) + "\033[J")
		}
	}
}
