// Package live keeps the latest CNC values of every machine in memory and decides what to push
// to which Socket.IO room (ADR 0009, docs/live-socket-events.md). The collector's workers publish
// into the Hub without waiting for anything; the hub's own goroutine (Run) builds the messages
// and hands them to the transport (handler.LiveHandler) through the EmitFunc. No socket library
// and no SQL here.
package live

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"7th-sense-server/internal/model"
)

// Events the server sends to the browser.
const (
	EventValues  = "machine:values" // ValuesMessage: the changed codes of one room's variable set
	EventStatus  = "machine:status" // StatusMessage: online / mode / program changed
	EventMachine = "machine:event"  // EventMessage: alarm, emergency stop, cycle opened / closed
)

// Channels a subscription may ask for besides variable codes.
const (
	ChannelStatus = "status"
	ChannelEvent  = "event"
)

// Machine events the hub forwards on top of the machine_events types: the worker opens and
// closes cycles in machine_cycles, not as machine_events rows.
const (
	EventCycleOpened = "cycle_opened"
	EventCycleClosed = "cycle_closed"
)

// DefaultInterval is a subscription's interval when the browser names none.
const DefaultInterval = time.Second

// roomHashLength is how many hex characters of the SHA-256 of a variable set name its room.
const roomHashLength = 16

// forwardedEvents are the machine events screens act on (they reload their REST data); the
// other machine_events types (connect_ok, dll_error, …) are not pushed.
var forwardedEvents = map[string]bool{
	model.EventAlarmOn: true, model.EventAlarmOff: true,
	model.EventEmergencyOn: true, model.EventEmergencyOff: true,
	EventCycleOpened: true, EventCycleClosed: true,
}

// EmitFunc sends event with payload to every socket in room.
type EmitFunc func(room, event string, payload any)

// Status is what a machine does now: online is mode ≠ Offline, the same rule as
// machines.online_status.
type Status struct {
	Online      bool      `json:"online"`
	ModeID      int       `json:"modeId"`
	ModeSince   time.Time `json:"modeSince"`
	ProgramName string    `json:"programName"`
}

// ValuesMessage is machine:values: only the codes of the room's set that changed since its last
// message, newest value of each. Seq is the machine's sequence number of the newest change in it.
type ValuesMessage struct {
	MachineID int64          `json:"machineId"`
	Seq       int64          `json:"seq"`
	At        time.Time      `json:"at"`
	Values    map[string]any `json:"values"`
}

// StatusMessage is machine:status.
type StatusMessage struct {
	MachineID int64 `json:"machineId"`
	Seq       int64 `json:"seq"`
	Status
}

// EventMessage is machine:event: Type is a machine_events type (alarm_on, emergency_on, …) or
// cycle_opened / cycle_closed.
type EventMessage struct {
	MachineID int64     `json:"machineId"`
	Seq       int64     `json:"seq"`
	Type      string    `json:"type"`
	Code      *int      `json:"code,omitempty"`
	Message   string    `json:"message,omitempty"`
	At        time.Time `json:"at"`
}

// Snapshot is the current value of the requested codes of one machine, sent once when a screen
// starts watching. The status fields are left out until the collector reported the machine, and
// a code never read is left out of Values. Messages with a Seq not above it are older than it.
type Snapshot struct {
	MachineID int64          `json:"machineId"`
	Seq       int64          `json:"seq"`
	At        *time.Time     `json:"at"`
	Values    map[string]any `json:"values"`
	*Status
}

// Hub holds the latest values per machine and the rooms that watch them.
type Hub struct {
	minInterval time.Duration
	wake        chan struct{} // holds one request to send the outbox; more collapse into it

	mu       sync.Mutex
	emit     EmitFunc
	machines map[int64]*machineState
	rooms    map[string]*valuesRoom // by room name
	outbox   []outgoing             // status and event messages waiting for Run
}

// machineState is one machine's latest values and its value rooms.
type machineState struct {
	seq    int64
	at     time.Time      // when a value last changed
	values map[string]any // code → float64 or string
	codes  map[string]bool
	status *Status
	rooms  map[string]*valuesRoom
}

// valuesRoom is a machine + variable set + interval that sockets share: its message is built once.
type valuesRoom struct {
	name      string
	machineID int64
	codes     map[string]bool
	interval  time.Duration
	members   int
	pending   map[string]any // changed codes not sent yet, newest value wins
	seq       int64          // machine seq of the newest pending change
	at        time.Time      // time of the newest pending change
	flushedAt time.Time
}

// outgoing is one message for one room.
type outgoing struct {
	room    string
	event   string
	payload any
}

// NewHub creates a Hub whose rooms send at most once per minInterval.
func NewHub(minInterval time.Duration) *Hub {
	return &Hub{
		minInterval: minInterval,
		wake:        make(chan struct{}, 1),
		machines:    map[int64]*machineState{},
		rooms:       map[string]*valuesRoom{},
	}
}

// SetEmitter makes the hub send its messages through emit (the Socket.IO server).
func (hub *Hub) SetEmitter(emit EmitFunc) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.emit = emit
}

// Run sends the waiting messages until ctx ends: status and events as soon as they arrive,
// each room's changed values at most once per its interval.
func (hub *Hub) Run(ctx context.Context) {
	ticker := time.NewTicker(hub.minInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-hub.wake:
			hub.Flush(time.Now())
		case now := <-ticker.C:
			hub.Flush(now)
		}
	}
}

// Flush sends the outbox and the rooms whose interval passed at now and that have changes.
// Run calls it; tests call it directly.
func (hub *Hub) Flush(now time.Time) {
	hub.mu.Lock()
	messages := hub.outbox
	hub.outbox = nil
	for _, room := range hub.rooms {
		if len(room.pending) == 0 || now.Sub(room.flushedAt) < room.interval {
			continue
		}
		messages = append(messages, outgoing{room: room.name, event: EventValues, payload: ValuesMessage{
			MachineID: room.machineID, Seq: room.seq, At: room.at, Values: room.pending,
		}})
		room.pending = map[string]any{}
		room.flushedAt = now
	}
	emit := hub.emit
	hub.mu.Unlock()

	if emit == nil {
		return
	}
	for _, message := range messages {
		emit(message.room, message.event, message.payload)
	}
}

// SetMachineCodes records the variable codes the collector reads for machineID: subscriptions
// to other codes of it are answered as unknown.
func (hub *Hub) SetMachineCodes(machineID int64, codes []string) {
	if hub == nil {
		return
	}
	known := make(map[string]bool, len(codes))
	for _, code := range codes {
		known[code] = true
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.machine(machineID).codes = known
}

// PublishValues takes one reading of machineID (code → float64 or string, every variable read):
// the values that differ from the latest become pending in the rooms whose set has them. It only
// merges under the lock and never waits for a socket.
func (hub *Hub) PublishValues(machineID int64, at time.Time, values map[string]any) {
	if hub == nil {
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	machine := hub.machine(machineID)
	changed := false
	for code, value := range values {
		if previous, found := machine.values[code]; found && previous == value {
			continue
		}
		if !changed {
			changed = true
			machine.seq++
			machine.at = at
		}
		machine.values[code] = value
		for _, room := range machine.rooms {
			if room.codes[code] {
				room.pending[code] = value
				room.seq, room.at = machine.seq, at
			}
		}
	}
}

// PublishStatus records the machine's status and sends it to its status room when it changed.
func (hub *Hub) PublishStatus(machineID int64, status Status) {
	if hub == nil {
		return
	}
	hub.mu.Lock()
	machine := hub.machine(machineID)
	if machine.status != nil && *machine.status == status {
		hub.mu.Unlock()
		return
	}
	machine.seq++
	machine.status = &status
	hub.outbox = append(hub.outbox, outgoing{room: StatusRoom(machineID), event: EventStatus,
		payload: StatusMessage{MachineID: machineID, Seq: machine.seq, Status: status}})
	hub.mu.Unlock()
	hub.requestFlush()
}

// PublishEvent sends a machine event to the machine's event room when screens act on its type.
func (hub *Hub) PublishEvent(machineID int64, eventType string, code *int, message string, at time.Time) {
	if hub == nil || !forwardedEvents[eventType] {
		return
	}
	hub.mu.Lock()
	machine := hub.machine(machineID)
	machine.seq++
	hub.outbox = append(hub.outbox, outgoing{room: EventRoom(machineID), event: EventMachine,
		payload: EventMessage{MachineID: machineID, Seq: machine.seq, Type: eventType, Code: code, Message: message, At: at}})
	hub.mu.Unlock()
	hub.requestFlush()
}

// JoinValues adds a member to the room of machineID + codes + interval and returns its name, the
// codes kept (sorted, duplicates removed) and the codes the machine does not read. While the
// collector has not reported the machine's codes yet every code is kept. room is "" when no code
// is left.
func (hub *Hub) JoinValues(machineID int64, codes []string, interval time.Duration) (room string, kept, unknown []string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	machine := hub.machine(machineID)
	for _, code := range codes {
		switch {
		case slices.Contains(kept, code) || slices.Contains(unknown, code):
		case machine.codes != nil && !machine.codes[code]:
			unknown = append(unknown, code)
		default:
			kept = append(kept, code)
		}
	}
	if len(kept) == 0 {
		return "", nil, unknown
	}
	slices.Sort(kept)
	interval = hub.Interval(interval)
	room = valuesRoomName(machineID, kept, interval)
	existing := hub.rooms[room]
	if existing == nil {
		existing = &valuesRoom{name: room, machineID: machineID, codes: map[string]bool{}, interval: interval,
			pending: map[string]any{}}
		for _, code := range kept {
			existing.codes[code] = true
		}
		hub.rooms[room] = existing
		machine.rooms[room] = existing
	}
	existing.members++
	return room, kept, unknown
}

// LeaveValues removes a member from room; a room without members is dropped and no longer built.
func (hub *Hub) LeaveValues(room string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	existing := hub.rooms[room]
	if existing == nil {
		return
	}
	existing.members--
	if existing.members > 0 {
		return
	}
	delete(hub.rooms, room)
	if machine := hub.machines[existing.machineID]; machine != nil {
		delete(machine.rooms, room)
	}
}

// Snapshot returns the latest values of codes of machineID and its status.
func (hub *Hub) Snapshot(machineID int64, codes []string) Snapshot {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	machine := hub.machine(machineID)
	snapshot := Snapshot{MachineID: machineID, Seq: machine.seq, Values: map[string]any{}}
	for _, code := range codes {
		if value, found := machine.values[code]; found {
			snapshot.Values[code] = value
		}
	}
	if !machine.at.IsZero() {
		at := machine.at
		snapshot.At = &at
	}
	if machine.status != nil {
		status := *machine.status
		snapshot.Status = &status
	}
	return snapshot
}

// Interval returns requested raised to the minimum interval, or DefaultInterval when zero.
func (hub *Hub) Interval(requested time.Duration) time.Duration {
	if requested <= 0 {
		return max(DefaultInterval, hub.minInterval)
	}
	return max(requested, hub.minInterval)
}

// machine returns the state of machineID, creating it. The caller holds hub.mu.
func (hub *Hub) machine(machineID int64) *machineState {
	machine := hub.machines[machineID]
	if machine == nil {
		machine = &machineState{values: map[string]any{}, rooms: map[string]*valuesRoom{}}
		hub.machines[machineID] = machine
	}
	return machine
}

// requestFlush wakes Run without waiting.
func (hub *Hub) requestFlush() {
	select {
	case hub.wake <- struct{}{}:
	default:
	}
}

// StatusRoom is the room of machineID's status messages.
func StatusRoom(machineID int64) string {
	return "m:" + strconv.FormatInt(machineID, 10) + ":status"
}

// EventRoom is the room of machineID's event messages.
func EventRoom(machineID int64) string {
	return "m:" + strconv.FormatInt(machineID, 10) + ":event"
}

// valuesRoomName names the room of machineID + sorted codes + interval:
// m:<id>:v:<hash of the codes>@<interval ms>. The browser never names rooms.
func valuesRoomName(machineID int64, sortedCodes []string, interval time.Duration) string {
	sum := sha256.Sum256([]byte(strings.Join(sortedCodes, ",")))
	return "m:" + strconv.FormatInt(machineID, 10) + ":v:" + hex.EncodeToString(sum[:])[:roomHashLength] +
		"@" + strconv.FormatInt(interval.Milliseconds(), 10)
}
