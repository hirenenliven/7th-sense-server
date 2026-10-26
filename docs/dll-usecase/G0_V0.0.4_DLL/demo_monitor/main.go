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

// =============================================================================
//
//	Machine CONFIG
//
// =============================================================================
type MachineDef struct {
	Title string
	Type  int // 1 = Fanuc, 2 = Siemens
	IP    string
	Port  int
	User  string // Siemens only
	Pass  string // Siemens only
}

var (
	// active machine target
	TargetMachine = "Siemens"

	Presets = map[string]MachineDef{
		"Fanuc": {
			Title: "FANUC",
			Type:  1,
			IP:    "192.168.0.254",
			Port:  8193,
		},
		"Siemens": {
			Title: "SINUMERIK",
			Type:  2,
			IP:    "192.168.0.254",
			Port:  8193,
			User:  "JYOTI",
			Pass:  "12345678",
		},
	}
)

// =============================================================================
//
//	SNAPSHOT STRUCT (C++ CncMasterSnapshot layout)
//
// =============================================================================
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

	MachPos  [5]float64
	AbsPos   [5]float64
	RelPos   [5]float64
	DistToGo [5]float64

	MachDec      [5]int16
	AbsDec       [5]int16
	RelDec       [5]int16
	DistToGoDec  [5]int16
	MachUnit     [5]int16 // 0=mm, 1=inch
	AbsUnit      [5]int16
	RelUnit      [5]int16
	DistToGoUnit [5]int16

	LinearDec, RotaryDec, AbsLinearDec, AbsRotaryDec     int16
	LinearUnit, RotaryUnit, AbsLinearUnit, AbsRotaryUnit int16

	SpindeOverrideActual     float64
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

	ProgNumber, CurProgLine, TotalProgLines           int32
	ProgProgress, ProgramPartCount, RequiredPartCount int32

	ProgName       [64]byte
	ActiveProgPath [256]byte
	ProgList       [32768]byte

	TimerMinute, TimerMsec int32
	ToolNumber             int32
	ToolLength, ToolRadius float64
	ToolOffsetDec          int16

	ActiveGCode [16]byte
	GcodeGroup  [512]byte
	CmdFeedOvr  [16]byte
	CmdSpOvr    [16]byte
	_           [2]byte

	AlarmNo  int32
	AlarmMsg [128]byte
	OpMsgNo  int32
	OpMsg    [256]byte
	_        [4]byte
}

// =============================================================================
//
//	GLOBAL STATE & DLL
//
// =============================================================================
func initDLL() {
	libDir, _ := filepath.Abs("Lib")
	ptr, _ := syscall.UTF16PtrFromString(libDir)
	syscall.NewLazyDLL("kernel32.dll").NewProc("SetDllDirectoryW").Call(uintptr(unsafe.Pointer(ptr)))

	mod = syscall.NewLazyDLL(filepath.Join(libDir, "CncCommunication.dll"))
	procConnect = mod.NewProc("CNCConnect")
	procDisconnect = mod.NewProc("CNCdisconnect")
}

var (
	mu             sync.Mutex
	latest         CncMasterSnapshot
	have           bool
	newData        = make(chan struct{}, 1)
	mod            *syscall.LazyDLL
	procConnect    *syscall.LazyProc
	procDisconnect *syscall.LazyProc
	cfg            MachineDef

	modes = map[int16]string{0: "JOG", 1: "MDI", 2: "AUTO"}
	progs = map[int16]string{0: "Reset", 1: "Stop", 2: "Active", 3: "Waiting", 4: "Interrupted"}
)

func cncDataCallback(s uintptr) uintptr {
	if s == 0 {
		return 0
	}
	mu.Lock()
	latest = *(*CncMasterSnapshot)(unsafe.Pointer(s))
	have = true
	mu.Unlock()

	select {
	case newData <- struct{}{}:
	default:
	}
	return 0
}

func onErr(code uintptr, msg uintptr) uintptr {
	fmt.Printf("DLL error: %d\n", code)
	return 0
}

func cstr(b []byte) string {
	s := string(b)
	if i := strings.IndexByte(s, 0); i >= 0 {
		return s[:i]
	}
	return s
}

func enableVT() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	var m uint32
	h := os.Stdout.Fd()
	k32.NewProc("GetConsoleMode").Call(h, uintptr(unsafe.Pointer(&m)))
	k32.NewProc("SetConsoleMode").Call(h, uintptr(m|0x0004))
}

// =============================================================================
//
//	DASHBOARD RENDERER
//
// =============================================================================
func render(d *CncMasterSnapshot) string {
	var b strings.Builder
	line := func(format string, a ...interface{}) {
		fmt.Fprintf(&b, format+"\033[K\n", a...)
	}

	ts := time.Now()
	if d.SampleTimestampMs > 0 {
		ts = time.UnixMilli(d.SampleTimestampMs)
	}

	title := fmt.Sprintf(" %s LIVE DATA   %s:%d   %s", cfg.Title, cfg.IP, cfg.Port, ts.Format("15:04:05"))
	line("\033[1;97;44m%-78s\033[0m", title)

	modeStr := modes[d.MachineMode]
	if modeStr == "" {
		modeStr = fmt.Sprintf("%d", d.MachineMode)
	}
	progStr := progs[d.ProcStatus]
	if progStr == "" {
		progStr = fmt.Sprintf("%d", d.ProcStatus)
	}

	conn := "\033[91mNO\033[0m"
	if d.IsConnected == 1 {
		conn = "\033[92mYES\033[0m"
	}

	line("\033[90mMODE:\033[0m \033[1m%-6s\033[0m  \033[90mSTATUS:\033[0m \033[1m%-8s\033[0m  \033[90mCONNECTED:\033[0m %s", modeStr, progStr, conn)
	line("\033[90m%s\033[0m", strings.Repeat("-", 78))

	line("\033[1m%-5s %13s %13s %13s %13s %13s\033[0m", "AXIS", "MCS (mm)", "WCS (inch)", "WCS (mm)", "DIST-TO-GO", "RELATIVE")

	for i := 0; i < 5; i++ {
		name := cstr(d.AxisNames[i][:])
		if name == "" {
			continue
		}

		// Axis color according to image: X = Cyan, Y = Green, Z = Yellow
		axisCol := "\033[1;97m" // default white
		switch strings.ToUpper(name) {
		case "X":
			axisCol = "\033[1;96m" // Cyan
		case "Y":
			axisCol = "\033[1;92m" // Green
		case "Z":
			axisCol = "\033[1;93m" // Yellow
		}

		if name == "A" || name == "B" || name == "C" {
			line("%s%-5s %12.3f° %12.3f° %13s %12.3f° %12.3f°\033[0m", axisCol, name, d.MachPos[i], d.AbsPos[i], "-", d.DistToGo[i], d.RelPos[i])
		} else {
			mcs := d.MachPos[i]
			if d.MachUnit[i] == 1 {
				mcs *= 25.4
			}
			wInch, wMm := d.AbsPos[i], d.AbsPos[i]
			if d.AbsUnit[i] == 1 {
				wMm = d.AbsPos[i] * 25.4
			} else {
				wInch = d.AbsPos[i] / 25.4
			}
			line("%s%-5s %13.3f %13.4f %13.3f %13.4f %13.4f\033[0m", axisCol, name, mcs, wInch, wMm, d.DistToGo[i], d.RelPos[i])
		}
	}
	line("\033[90m%s\033[0m", strings.Repeat("-", 78))

	line("\033[1;34mSPINDLE \033[0m\033[90mact:\033[0m %-5.0f rpm  \033[90mload:\033[0m %-3d%%  \033[90movr:\033[0m %-3d%%  \033[90mcmd:\033[0m %s",
		d.SpindeOverrideActual, d.SpindleLoadPercetnage, d.SpindeOverridePercentage, cstr(d.CmdSpOvr[:]))
	line("\033[1;34mFEED    \033[0m\033[90mact:\033[0m %-7.1f  \033[90movr:\033[0m %-3d%%  \033[90mcmd:\033[0m %s",
		d.FeedRateActual, d.FeedRatePercetnage, cstr(d.CmdFeedOvr[:]))
	line("\033[1;34mSERVO   \033[0m\033[90mload X:\033[0m %-3d%%  \033[90mY:\033[0m %-3d%%  \033[90mZ:\033[0m %-3d%%  |  \033[1;34mTOOL:\033[0m %d",
		d.AxisLoadXPercetnage, d.AxisLoadYPercetnage, d.AxisLoadZPercetnage, d.ToolNumber)
	line("\033[90m%s\033[0m", strings.Repeat("-", 78))

	line("\033[1;34mPROGRAM \033[0m\033[90mname:\033[0m %s  \033[90mline:\033[0m %d  \033[90mpath:\033[0m %s",
		cstr(d.ProgName[:]), d.CurProgLine, cstr(d.ActiveProgPath[:]))
	line("\033[1;34mG-CODE  \033[0m\033[90mact:\033[0m %s", cstr(d.ActiveGCode[:]))
	line("\033[90m%s\033[0m\n\033[90mPress Ctrl+C to stop\033[0m", strings.Repeat("-", 78))

	return b.String()
}

// =============================================================================
//
//	MAIN
//
// =============================================================================
func main() {
	//1. Check : Terminal Argument  (e.g. `go run main.go Fanuc` or `2`)
	for _, arg := range os.Args[1:] {
		low := strings.ToLower(arg)
		if low == "1" {
			low = "Fanuc"
		} else if low == "2" {
			low = "Siemens"
		}
		if _, ok := Presets[low]; ok {
			TargetMachine = low
		}
	}

	cfg = Presets[TargetMachine]

	//2. Siemens needs env credentials set
	if cfg.Type == 2 {
		os.Setenv("SIEMENS_USER", cfg.User)
		os.Setenv("SIEMENS_PASS", cfg.Pass)
	}

	//3. Load DLL and connect
	initDLL()
	if err := mod.Load(); err != nil {
		fmt.Printf("Failed to load DLL: %v\n", err)
		return
	}

	//4. Connect to CNC
	ipPtr, _ := syscall.BytePtrFromString(cfg.IP)
	ret, _, _ := procConnect.Call(
		uintptr(unsafe.Pointer(ipPtr)),
		uintptr(cfg.Port),
		syscall.NewCallback(cncDataCallback),
		syscall.NewCallback(onErr),
		uintptr(cfg.Type),
	)

	id := int32(ret)
	if id < 0 {
		fmt.Printf("Connection error: %d\n", id)
		return
	}

	//5. Enable VT100 and start dashboard loop
	enableVT()
	fmt.Print("\033[2J\033[?25l") // Clear screen, hide cursor
	defer fmt.Print("\033[?25h")  // Show cursor on exit

	//6. Handle Ctrl+C and new data
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	//7. Main loop
	for {
		select {
		case <-sig:
			procDisconnect.Call(uintptr(id))
			fmt.Print("\033[2J\033[HDisconnected.\n")
			return
		case <-newData:
			mu.Lock()
			snap := latest
			ok := have
			mu.Unlock()
			if ok {
				fmt.Print("\033[H" + render(&snap) + "\033[J")
			}
		}
	}
}
