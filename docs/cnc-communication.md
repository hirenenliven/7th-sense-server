# CncCommunication.dll — API & Integration Guide (Go Client)

> Text version of `docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf`, kept searchable for tools and the knowledge graph. The PDF is the source; update both together.

## 1. Overview

CncCommunication.dll is a dynamic link library implemented in C++ designed to stream real-time telemetry from industrial Computer Numerical Control (CNC) controllers directly to external software clients via asynchronous event loops.

The Go application interfaces with the DLL by loading exports, issuing a connection call (`CNCConnect`), and handling inbound updates through unmanaged callback pointers dispatching instances of `CncMasterSnapshot`.

The DLL encapsulates the entire OPC UA session — endpoint discovery, secure channel negotiation, and certificate handling — so the Go client requires no OPC UA certificate bundles, zip archives, or session configuration. It only supplies the controller address, port, and credentials.

## 2. System Architecture

```
┌──────────────┐   OPC UA    ┌───────────────────────────┐
│  SINUMERIK   │ ──────────► │   CncCommunication.dll    │
│  CNC (PLC)   │  port 4840  │  (open62541 inside, x86)  │
└──────────────┘             └─────────────┬─────────────┘
                                           │ stdcall callback
                                           │ (CncMasterSnapshot, 34440 bytes)
                                           ▼
                             ┌───────────────────────────┐
                             │         Go Client         │
                             │ callback → copy → render  │
                             └───────────────────────────┘
```

Data flow (event-driven, no polling in the client):

1. The Go client loads the DLL and registers two stdcall callbacks (`syscall.NewCallback`).
2. `CNCConnect` starts a background worker thread inside the DLL.
3. The worker polls the controller, assembles a `CncMasterSnapshot`, and invokes the data callback with a pointer to the buffer.
4. The callback copies the payload into the managed Go struct and signals the renderer via a buffered channel.
5. The main goroutine renders the live ANSI dashboard on each signal.

## 3. Package Artifacts & Dependencies

| Filename | Dependency Domain | Description |
|---|---|---|
| `open62541.dll` | OPC UA Runtime | Open-source C-based OPC UA stack. |
| `libcrypto-3.dll`, `libssl-3.dll` | OpenSSL | TLS encryption and handshake providers for OPC UA transport. |
| `Fwlib32.dll`, `fwlibe1.dll` | Fanuc FOCAS2 | Low-level proprietary protocol drivers for Fanuc controllers. |
| `CncCommunication.dll` | Core Runtime | Primary orchestration and export library. |
| `client_key.der` (`.pem`) | Security Profile | Private key corresponding to client certificate. |
| `client_cert.der` (`.pem`) | Security Profile | X.509 client certificate for Siemens OPC UA authentication. |

### 3.1 Environment & Prerequisites

| Requirement | Value |
|---|---|
| Target Architecture | 32-bit DLL (x86) — compile the Go client with `set GOARCH=386` |
| Credentials | `SIEMENS_USER` / `SIEMENS_PASS` environment variables set before `CNCConnect` |

## 4. API Reference

### 4.1 Function Signatures

```c
// Starts the background worker, opens network sockets, and begins polling
int32_t CNCConnect(
    const char* ip,          // controller IP address (ANSI, NUL-terminated)
    int32_t     port,        // controller port (e.g. 4840)
    uintptr_t   dataCb,      // stdcall callback: void OnData(const CncMasterSnapshot*)
    uintptr_t   errCb,       // stdcall callback: void OnErr(int32_t code, const char* msg)
    int32_t     machineType  // engine id, see §4.2
);

// Stops polling threads and releases underlying protocol sockets
void CNCdisconnect(int32_t machineId);
```

### 4.2 Machine Identifier Mapping

| Enum Value | Machine Controller Type | Driver Protocol |
|---|---|---|
| 0 | Simulator | Local loopback mock engine |
| 1 | Fanuc | FOCAS2 Ethernet (fwlib32) |
| 2 | Siemens | Sinumerik OPC UA (active profile in this project) |
| 3 | Syntec | Syntec Open API / TCP |
| 4 | Mitsubishi | — (not specified) |

## 5. DLL Loading & Binding (Go)

```go
package main

import (
	"fmt"
	"os"
	"syscall"
)

var (
	mod            = syscall.NewLazyDLL("CncCommunication.dll")
	procConnect    = mod.NewProc("CNCConnect")
	procDisconnect = mod.NewProc("CNCdisconnect")
)

func init() {
	// Eagerly resolve the DLL handle at application startup
	if err := mod.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: Failed loading CncCommunication.dll: %v\n", err)
		os.Exit(1)
	}
}
```

## 6. Snapshot Data Structure

### 6.1 Layout

`CncMasterSnapshot` matches the C++ definition 1:1 (`#pragma pack(8)`, x86) and totals exactly **34,440 bytes**. The Go client asserts this at startup:

```go
const expectSize = 34440
if unsafe.Sizeof(CncMasterSnapshot{}) != expectSize {
	// abort -- struct drift between client and DLL
}
```

### 6.2 Data Field Reference

Field names are as defined in the DLL, including their spelling (`SpindeOverrideActual`, `...Percetnage`); the Go struct must use the same order and types.

**Connection & state**

| Struct Member | Type | Units / Encoding | Description |
|---|---|---|---|
| `IsConnected` | `int16` | Boolean (0/1) | Network connection state |
| `ResultStatus` | `int16` | Error Flag | Controller operational status register |
| `SampleTimestampMs` | `int64` | Milliseconds | Epoch timestamp of snapshot emission |
| `MachineMode` | `int16` | Enum (0=JOG, 1=MDA, 2=AUTO) | Controller operating mode |
| `ProcStatus` | `int16` | Enum (0=Reset, 1=Stop, 2=Run) | NC channel execution state |
| `IsEmergency` | `int16` | Boolean (0/1) | Emergency stop (E-STOP) state |

**Axes & positions**

| Struct Member | Type | Units / Encoding | Description |
|---|---|---|---|
| `AxisNames` | `[5][4]byte` | ASCII Strings | Axis designations (e.g., X, Y, Z, B, C) |
| `MachPos` | `[5]float64` | Millimeters / Degrees | Machine Coordinate System (MCS) positions |
| `AbsPos` | `[5]float64` | Inches / Degrees | Work Coordinate System (WCS) positions |
| `RelPos` | `[5]float64` | Inches / Degrees | Relative incremental offsets |
| `DistToGo` | `[5]float64` | Inches / Degrees | Remaining vector to destination |

**Spindle, feed & loads**

| Struct Member | Type | Units / Encoding | Description |
|---|---|---|---|
| `SpindeOverrideActual` | `float64` | RPM | Active spindle surface velocity |
| `SpindleLoadPercetnage` | `int16` | Percent (%) | Spindle drive current load |
| `SpindeOverridePercentage` | `int16` | Percent (%) | Spindle potentiometer position |
| `FeedRateActual` | `float64` | mm/min or in/min | Active linear tool translation rate |
| `FeedRatePercetnage` | `int16` | Percent (%) | Feedrate override potentiometer percentage |
| `AxisLoadX/Y/ZPercentage` | `int16` | Percent (%) | Servo axis drive loads |

**Program**

| Struct Member | Type | Units / Encoding | Description |
|---|---|---|---|
| `ProgNumber` | `int32` | Integer ID | Program sequence identifier (O-number) |
| `CurProgLine` | `int32` | Line Count | Current block executing in interpreter |
| `TotalProgLines` | `int32` | Line Count | Total line count of active file |
| `ProgName` | `[64]byte` | ASCII String | File name of active NC program |
| `ActiveProgPath` | `[256]byte` | ASCII Path | Storage path of active NC artifact |

**Tool, G-code & alarms**

| Struct Member | Type | Units / Encoding | Description |
|---|---|---|---|
| `ToolNumber` | `int32` | Slot/ID | Spindle-mounted active tool identifier |
| `ToolLength` | `float64` | Metric / Imperial | Primary height compensation value |
| `ToolRadius` | `float64` | Metric / Imperial | Cutter radius compensation value |
| `ActiveGCode` | `[16]byte` | ASCII | Currently active G-code group |
| `GcodeGroup` | `[512]byte` | ASCII (reused) | DLL reuses this field for the `;`-separated missing-node list |
| `AlarmNo` | `int32` | Code (0 = OK) | Primary alarm or fault number |
| `AlarmMsg` | `[128]byte` | UTF-8 / ASCII | Alarm description message |

## 7. Asynchronous Callback Architecture

Because callbacks originate from an unmanaged C++ worker thread, callbacks must remain strictly non-blocking.

```go
package main

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

var (
	dataMutex  sync.RWMutex
	lastRecord CncMasterSnapshot
	hasRecord  bool

	// Declare callbacks as package-level singletons
	dataCallback = syscall.NewCallback(onDataCallback)
	errCallback  = syscall.NewCallback(onErrorCallback)
)

func onDataCallback(ptr uintptr) uintptr {
	if ptr == 0 {
		return 0
	}
	// Copy unmanaged memory buffer into managed Go struct
	dataMutex.Lock()
	lastRecord = *(*CncMasterSnapshot)(unsafe.Pointer(ptr))
	hasRecord = true
	dataMutex.Unlock()
	return 0
}

func onErrorCallback(code, messagePtr uintptr) uintptr {
	fmt.Printf("[DLL ERROR] Event Code: %d\n", code)
	return 0
}
```

## 8. Running the Go Client

```shell
cd C:\Avim\Prinsi\G0_V0.0.2_DLL
set GOARCH=386
go run sinumerik_client.go          # live dashboard
go run sinumerik_client.go check    # OPC UA node existence report
```

Press Ctrl+C in the dashboard to disconnect cleanly.

### Check-mode OPC UA node mapping

| Data | OPC UA Node Path |
|---|---|
| Operating mode | `/Bag/State/opMode[u1]` |
| Program status | `/Channel/State/acProg[u1,1]` |
| Spindle speed / override | `/Channel/Spindle/actSpeed[u1,1]`, `speedOvr[u1,1]` |
| Feed rate / override | `/Channel/State/actFeedRateIpo[u1]`, `feedRateOvr[u1,1]` |
| Axis / spindle loads | `/Channel/MachineAxis/aaLoad[u1,1..4]` |
| Positions | `/Channel/MachineAxis/actToolBasePos[u1,n]`, `/Channel/GeometricAxis/aaPcsRel[u1,n]` |
| Program info | `/Channel/ProgramInfo/progName[u1,1]`, `actLineNumber` |
| Tool | `/Channel/State/actTNumber[u1,1]` |

## 9. Controller-Specific Profiles

### 9.1 Siemens Sinumerik (OPC UA) — Engine ID: 2

- **Default Ports:** 4840 (Native NCU) or 4841 (custom industrial port profile).
- **Security Certificates:** place `client_cert.der` and `client_key.der` directly in the active binary directory.
- **Trust Requirements:** the client certificate must be signed or marked as trusted within the Sinumerik NCU Access Administration interface. Connections using unapproved certificates are rejected with status `BadRequestNotAllowed`.
- **Credentials:** passed to the DLL via `SIEMENS_USER` / `SIEMENS_PASS` environment variables, set before calling `CNCConnect`.

### 9.2 Fanuc — Engine ID: 1

- **Default Port:** 8193 (default FOCAS interface).
- **Dependencies:** requires `Fwlib32.dll` and `fwlibe1.dll` present in system lookup paths.

### 9.3 Syntec — Engine ID: 3

- **Default Port:** typically 5001 or controller-assigned port.

## 10. Diagnostics & Troubleshooting

| Observable Error / Symptom | Root Cause Analysis | Remediation Step |
|---|---|---|
| `0xC0000135` / DLL load failed: module could not be found | Missing transitive native dependencies — a dependency DLL is not resolvable. | All 6 DLLs must be in `Lib\` next to the client. Run `dumpbin /dependents CncCommunication.dll` to find and add missing dependent `.dll` files. |
| `%1 is not a valid Win32 application` | 32-bit DLL loaded by a 64-bit host process. | Compile the Go client with `set GOARCH=386`; validate the GOARCH setting against the DLL architecture. |
| Memory panic / corrupted floats | Struct size mismatch or memory alignment issues. | Verify `unsafe.Sizeof(CncMasterSnapshot{}) == 34440`. Validate GOARCH setting against DLL architecture. |
| `CNCConnect` returns `-1` | Host unreachable or routing timeout. | Run `ping <ip>` and `Test-NetConnection -ComputerName <ip> -Port <port>` to check the network path. |
| Connects but `CONNECTED NO` | Wrong IP/credentials, or OPC UA disabled on the CNC. | Verify `SIEMENS_USER` / `SIEMENS_PASS` and endpoint address; confirm the OPC UA server is enabled on the NCU. |
| `BadRequestNotAllowed` | Siemens certificate trust rejection. | Export `client_cert.der` into the NCU certificate storage under trusted certificates. |
| `BadNodeIdUnknown` | OPC UA node is missing on the target model. | Node is unavailable for this hardware model; the field will default to 0. Use check mode (`go run sinumerik_client.go check`) to enumerate available nodes. |
| Garbled console output | Old console window has a broken terminal mode. | Close and reopen the terminal window. |
