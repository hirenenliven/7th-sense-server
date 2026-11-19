package opcua

import (
	"maps"
	"sync"
	"time"
)

// Connection states of a machine's PLC reader (MachineStatus.State).
const (
	StateConnecting  = "connecting"   // connecting or waiting to retry
	StateConnected   = "connected"    // session open, subscribed
	StateFailed      = "failed"       // the last connect or the connection failed; retrying
	StateNoVariables = "no_variables" // nothing to read: no active PLC variable with a node id
)

// NodeStatus is what the controller said about one subscribed node: Status "" while it is fine,
// else why not (refused at subscribe, or its last value was bad).
type NodeStatus struct {
	NodeID string    `json:"nodeId"`
	Status string    `json:"status"`
	At     time.Time `json:"at"`
}

// MachineStatus is what a machine's PLC reader is doing now, for the diagnostics.
type MachineStatus struct {
	State          string             `json:"state"`
	Endpoint       string             `json:"endpoint"`
	SecurityMode   string             `json:"securityMode"`
	SecurityPolicy string             `json:"securityPolicy"`
	ConnectedSince *time.Time         `json:"connectedSince"`
	LastError      string             `json:"lastError"`
	LastErrorAt    *time.Time         `json:"lastErrorAt"`
	LastChangeAt   *time.Time         `json:"lastChangeAt"`
	ChangeCount    int64              `json:"changeCount"`
	Nodes          map[int]NodeStatus `json:"-"` // by variable id
	UpdatedAt      time.Time          `json:"updatedAt"`
}

// StatusBoard keeps every machine's PLC reader status in memory: the PLC workers write it, the
// diagnostics read it. Safe for concurrent use; a nil board ignores writes and reads nothing.
type StatusBoard struct {
	mu       sync.Mutex
	machines map[int64]*MachineStatus
}

// NewStatusBoard creates an empty StatusBoard.
func NewStatusBoard() *StatusBoard {
	return &StatusBoard{machines: map[int64]*MachineStatus{}}
}

// Get returns a copy of the machine's status; false when no PLC worker reported it.
func (board *StatusBoard) Get(machineID int64) (MachineStatus, bool) {
	if board == nil {
		return MachineStatus{}, false
	}
	board.mu.Lock()
	defer board.mu.Unlock()
	status, found := board.machines[machineID]
	if !found {
		return MachineStatus{}, false
	}
	copied := *status
	copied.Nodes = maps.Clone(status.Nodes)
	return copied, true
}

// Update changes the machine's status under the lock (creating it first).
func (board *StatusBoard) Update(machineID int64, change func(status *MachineStatus)) {
	if board == nil {
		return
	}
	board.mu.Lock()
	defer board.mu.Unlock()
	status, found := board.machines[machineID]
	if !found {
		status = &MachineStatus{State: StateConnecting, Nodes: map[int]NodeStatus{}}
		board.machines[machineID] = status
	}
	change(status)
	status.UpdatedAt = time.Now()
}

// Remove forgets the machine (its worker stopped).
func (board *StatusBoard) Remove(machineID int64) {
	if board == nil {
		return
	}
	board.mu.Lock()
	defer board.mu.Unlock()
	delete(board.machines, machineID)
}
