//go:build windows && 386

// Command partcount-probe checks when CncCommunication.dll reports a Fanuc part count change.
// Beside the DLL connection it reads the counter straight from the control with Fanuc's
// Fwlib32.dll (parameter 6711 and macro #3901) and prints, every second and at every change,
// what the control holds and what the DLL's callbacks say. It only reads; run it beside the
// server while a cycle ends (M30):
//
//	$env:GOARCH="386"; go run ./cmd/partcount-probe -ip 192.168.0.254 -port 8193 -minutes 10
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
)

const (
	// Fanuc part count: parameter 6711 (parts machined) and system macro #3901 (the same counter).
	partCountParameter = 6711
	partCountMacro     = 3901
	focasTimeoutSec    = 10
	focasOK            = 0
)

// fanucCounter reads the part counter straight from the control through Fwlib32.dll.
type fanucCounter struct {
	handle    uint16
	readParam *syscall.LazyProc
	readMacro *syscall.LazyProc
	free      *syscall.LazyProc
}

// connectFanucCounter opens a FOCAS handle to ip:port.
func connectFanucCounter(libDir, ip string, port int) (*fanucCounter, error) {
	absoluteDir, err := filepath.Abs(libDir)
	if err != nil {
		return nil, err
	}
	fwlib := syscall.NewLazyDLL(filepath.Join(absoluteDir, "Fwlib32.dll"))
	if err := fwlib.Load(); err != nil {
		return nil, fmt.Errorf("load Fwlib32.dll: %w", err)
	}
	ipText, err := syscall.BytePtrFromString(ip)
	if err != nil {
		return nil, err
	}
	counter := &fanucCounter{
		readParam: fwlib.NewProc("cnc_rdparam"),
		readMacro: fwlib.NewProc("cnc_rdmacro"),
		free:      fwlib.NewProc("cnc_freelibhndl"),
	}
	result, _, _ := fwlib.NewProc("cnc_allclibhndl3").Call(uintptr(unsafe.Pointer(ipText)), uintptr(port),
		focasTimeoutSec, uintptr(unsafe.Pointer(&counter.handle)))
	if int16(result) != focasOK {
		return nil, fmt.Errorf("cnc_allclibhndl3 returned %d", int16(result))
	}
	return counter, nil
}

// parameter returns parameter 6711, or an error text.
func (counter *fanucCounter) parameter() string {
	var buffer struct { // IODBPSD with a 2-word value
		dataNumber int16
		dataType   int16
		value      int32
	}
	result, _, _ := counter.readParam.Call(uintptr(counter.handle), partCountParameter, 0,
		unsafe.Sizeof(buffer), uintptr(unsafe.Pointer(&buffer)))
	if int16(result) != focasOK {
		return fmt.Sprintf("err%d", int16(result))
	}
	return fmt.Sprint(buffer.value)
}

// macro returns macro #3901, or an error text.
func (counter *fanucCounter) macro() string {
	var buffer struct { // ODBM
		dataNumber int16
		dummy      int16
		value      int32
		decimals   int16
		_          int16
	}
	const odbmLength = 10
	result, _, _ := counter.readMacro.Call(uintptr(counter.handle), partCountMacro, odbmLength,
		uintptr(unsafe.Pointer(&buffer)))
	if int16(result) != focasOK {
		return fmt.Sprintf("err%d", int16(result))
	}
	value := float64(buffer.value)
	for range buffer.decimals {
		value /= 10
	}
	return fmt.Sprint(value)
}

func (counter *fanucCounter) close() {
	counter.free.Call(uintptr(counter.handle))
}

func main() {
	ip := flag.String("ip", "192.168.0.254", "controller IP")
	port := flag.Int("port", 8193, "FOCAS port")
	libDir := flag.String("lib", "Lib", "folder with CncCommunication.dll and Fwlib32.dll")
	minutes := flag.Int("minutes", 10, "how long to watch")
	flag.Parse()

	if err := run(*ip, *port, *libDir, time.Duration(*minutes)*time.Minute); err != nil {
		fmt.Fprintln(os.Stderr, "partcount-probe:", err)
		os.Exit(1)
	}
}

func run(ip string, port int, libDir string, duration time.Duration) error {
	// A FOCAS handle works only on the OS thread that opened it (EW_HANDLE -8 elsewhere).
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	driver, err := cnc.LoadDriver(config.CNCConfig{Slots: 1, LibDir: libDir, DLLName: "CncCommunication.dll"})
	if err != nil {
		return err
	}
	counter, err := connectFanucCounter(libDir, ip, port)
	if err != nil {
		return err
	}
	defer counter.close()

	engine := model.EngineFanuc
	connection, _, err := driver.Connect(model.MonitoredMachine{Name: "probe", LiveIP: ip, Port: port, CNCEngine: &engine})
	if err != nil {
		return err
	}
	defer connection.Disconnect()

	fmt.Println("time          control:6711  control:#3901 | dll:ProgramPartCount  ProcStatus  callbacks/s  callbacks total  dll sample age")
	lastLine, lastCallbacks := "", int64(0)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	end := time.After(duration)
	for {
		select {
		case <-end:
			return nil
		case now := <-ticker.C:
			reading, ok := connection.Latest()
			callbacks := connection.Stats().Callbacks
			dllCount, procStatus, sampleAge := "-", "-", "-"
			if ok {
				dllCount = fmt.Sprint(reading.Values[cnc.FieldProgramPartCount].Number)
				procStatus = fmt.Sprint(reading.Values[cnc.FieldProcStatus].Number)
				if !reading.SampleAt.IsZero() {
					sampleAge = now.Sub(reading.SampleAt).Round(100 * time.Millisecond).String()
				}
			}
			counts := fmt.Sprintf("%12s  %13s | %20s  %10s", counter.parameter(), counter.macro(), dllCount, procStatus)
			marker := ""
			if counts != lastLine {
				marker = "   <== changed"
			}
			fmt.Printf("%s  %s  %11d  %15d  %14s%s\n", now.Format("15:04:05.000"), counts,
				callbacks-lastCallbacks, callbacks, sampleAge, marker)
			lastLine, lastCallbacks = counts, callbacks
		}
	}
}
