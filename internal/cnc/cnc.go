// Package cnc talks to CNC controllers through the vendor DLL (docs/cnc-communication.md).
// It only connects, receives readings and disconnects: no database code, no business rules.
//
// The DLL is 32-bit: it is loaded only in windows/386 builds (dll_windows_386.go). Every other
// build returns ErrNotSupported from Open.
package cnc

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
)

var (
	// ErrNotSupported is returned by Open when the DLL cannot be used in this build.
	ErrNotSupported = errors.New("cnc: the CNC DLL is 32-bit; build the server with GOARCH=386 on Windows")
	// ErrNoFreeSlot is returned by Connect when every callback slot is in use (CNC_SLOTS).
	ErrNoFreeSlot = errors.New("cnc: no free callback slot")
	// ErrNoEngine is returned by Connect when the machine's controller has no cnc_engine.
	ErrNoEngine = errors.New("cnc: the machine's controller has no cnc_engine")
)

// SourceSnapshot is the only variable source the current DLL can read: a field of its
// CncMasterSnapshot. 'plc' and 'nc' addresses need a DLL read function (plan §12).
const SourceSnapshot = "snapshot"

// Value is one snapshot field: Text for text fields, Number for the rest.
type Value struct {
	Number float64
	Text   string
}

// Reading is one DLL callback.
type Reading struct {
	ReceivedAt   time.Time // server clock
	SampleAt     time.Time // DLL clock (zero when the DLL sent none)
	Connected    bool      // the DLL reaches the controller
	MissingNodes string    // nodes the controller does not have
	Values       [FieldCount]Value
}

// Stats are a connection's counters since it was opened.
type Stats struct {
	Callbacks int64
	Errors    int64
	LastError string
}

// Connection is one connected machine. The DLL thread publishes readings into it;
// the machine's worker waits on Signal and takes the newest reading with Latest.
type Connection struct {
	mu         sync.Mutex
	latest     Reading
	hasReading bool
	stats      Stats
	signal     chan struct{}
	disconnect func()
	once       sync.Once
}

func newConnection() *Connection {
	return &Connection{signal: make(chan struct{}, 1)}
}

// Signal receives a value when a new reading or error arrived. Several arrivals before the
// worker wakes collapse into one; Latest always returns the newest reading.
func (connection *Connection) Signal() <-chan struct{} {
	return connection.signal
}

// Latest returns the newest reading and whether there is one yet.
func (connection *Connection) Latest() (Reading, bool) {
	connection.mu.Lock()
	defer connection.mu.Unlock()
	return connection.latest, connection.hasReading
}

// Stats returns the connection's counters.
func (connection *Connection) Stats() Stats {
	connection.mu.Lock()
	defer connection.mu.Unlock()
	return connection.stats
}

// Disconnect closes the connection. Later calls do nothing.
func (connection *Connection) Disconnect() {
	connection.once.Do(func() {
		if connection.disconnect != nil {
			connection.disconnect()
		}
	})
}

// publish stores a reading and wakes the worker. It never blocks (it runs on the DLL thread).
func (connection *Connection) publish(reading Reading) {
	connection.mu.Lock()
	connection.latest, connection.hasReading = reading, true
	connection.stats.Callbacks++
	connection.mu.Unlock()
	connection.wake()
}

// publishError stores a DLL error and wakes the worker. It never blocks.
func (connection *Connection) publishError(code int, message string) {
	connection.mu.Lock()
	connection.stats.Errors++
	connection.stats.LastError = fmt.Sprintf("%d: %s", code, message)
	connection.mu.Unlock()
	connection.wake()
}

func (connection *Connection) wake() {
	select {
	case connection.signal <- struct{}{}:
	default:
	}
}

// Driver connects machines through the DLL.
type Driver struct {
	dll  *dllDriver
	name string
}

// Open loads the DLL. A collector cannot work without it.
func Open(cncConfig config.CNCConfig) (*Driver, error) {
	dll, err := openDLL(cncConfig)
	if err != nil {
		return nil, err
	}
	return &Driver{dll: dll, name: cncConfig.DLLName}, nil
}

// Name describes the driver for logs and events.
func (driver *Driver) Name() string {
	return driver.name
}

// Sources returns the variable sources the driver can read.
func (driver *Driver) Sources() []string {
	return []string{SourceSnapshot}
}

// Connect starts receiving readings from machine. It returns the DLL connection id
// (CNCConnect's result, also on failure) for events. It may block while the DLL connects.
func (driver *Driver) Connect(machine model.MonitoredMachine) (*Connection, int, error) {
	if machine.CNCEngine == nil {
		return nil, 0, ErrNoEngine
	}
	return driver.dll.connect(machine)
}
