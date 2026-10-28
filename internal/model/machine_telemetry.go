package model

import "time"

// Machine modes (machine_modes, the ids of the old machines_modes table) the collector sets.
// The user modes (2 Loading / Unloading, 7 Planned Downtime, 13 Leave Planned Downtime) and
// the cycle modes (8, 10, 11) are not set by the collector.
const (
	ModeSetup        = 1
	ModeInProduction = 3
	ModeHold         = 4
	ModeIdle         = 5
	ModeDowntime     = 6
	ModeOffline      = 9
	ModeRapid        = 12
)

// CNC engines (controllers.cnc_engine, the DLL's machineType).
const (
	EngineSimulator  = 0
	EngineFanuc      = 1
	EngineSiemens    = 2
	EngineSyntec     = 3
	EngineMitsubishi = 4
)

// Machine event types (machine_events.event_type).
const (
	EventCollectorStart = "collector_start"
	EventCollectorStop  = "collector_stop"
	EventValuesDropped  = "values_dropped"
	EventConfigError    = "config_error"
	EventConnectOK      = "connect_ok"
	EventConnectFailed  = "connect_failed"
	EventDisconnected   = "disconnected"
	EventStale          = "stale"
	EventReconnect      = "reconnect"
	EventDLLError       = "dll_error"
	EventAlarmOn        = "alarm_on"
	EventAlarmOff       = "alarm_off"
	EventEmergencyOn    = "emergency_on"
	EventEmergencyOff   = "emergency_off"
)

// MonitoredMachine is a machine a collector connects to: active and not deleted.
// Password is what the collector logs in to the controller with; it is never logged.
type MonitoredMachine struct {
	ID                    int64  `db:"id"`
	Name                  string `db:"name"`
	LiveIP                string `db:"live_ip"`
	Port                  int    `db:"port"`
	Username              string `db:"username"`
	Password              string `db:"password" json:"-"`
	ControllerVersionType int    `db:"controller_version_type"` // 1 = 828D, 2 = 840D, 0 = not set
	CNCEngine             *int   `db:"cnc_engine"`              // NULL = the controller has no engine set
}

// SameConnection reports whether other is reached the same way (a change needs a reconnect).
func (machine MonitoredMachine) SameConnection(other MonitoredMachine) bool {
	return machine.LiveIP == other.LiveIP && machine.Port == other.Port &&
		machine.Username == other.Username && machine.Password == other.Password &&
		equalEngine(machine.CNCEngine, other.CNCEngine)
}

func equalEngine(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// VariableConfig is one variable a machine reads, with its machine overrides applied
// (cnc_variables + cnc_variable_addresses + machine_variables).
type VariableConfig struct {
	ID            int      `db:"id"`
	Code          string   `db:"code"`
	DataType      string   `db:"data_type"` // number | bool | text
	StoreMode     int      `db:"store_mode"`
	Deadband      *float64 `db:"deadband"`
	MinIntervalMs int      `db:"min_interval_ms"`
	Source        string   `db:"source"` // snapshot | plc | nc
	Address       string   `db:"address"`
}

// IsText reports whether the variable holds text (machine_telementry.value_text).
func (variable VariableConfig) IsText() bool {
	return variable.DataType == "text"
}

// LiveValue is one row of machine_variable_values: the latest value for a machine-variable.
// Only numeric values are cached here (value_text is not stored in this table).
type LiveValue struct {
	MachineID  int64
	VariableID int
	Value      *float64
	UpdatedAt  time.Time
}

// TelemetryValue is one row of machine_telementry: Value for number / bool variables only.
type TelemetryValue struct {
	MachineID  int64
	VariableID int
	Value      *float64
	RecordedAt time.Time
}

// ModeChange closes the machine's open machine_mode_log period at At and opens one in ModeID.
type ModeChange struct {
	MachineID   int64
	ModeID      int
	ProgramName string
	At          time.Time
}

// MachineModePeriod is one row of machine_mode_log with its mode name. EndedAt is nil for the
// current period, whose DurationSec runs up to now.
type MachineModePeriod struct {
	ID          int64      `db:"id" json:"id"`
	ModeID      int        `db:"mode_id" json:"modeId"`
	ModeName    string     `db:"mode_name" json:"modeName"`
	ProgramName string     `db:"program_name" json:"programName"`
	StartedAt   time.Time  `db:"started_at" json:"startedAt"`
	EndedAt     *time.Time `db:"ended_at" json:"endedAt"`
	DurationSec int        `db:"duration_sec" json:"durationSec"`
}

// MachineModeStatus is a machine and the mode it is in now: its open machine_mode_log period
// (machines.mode_id when it has none). ModeID is nil for a machine never collected; ModeSince
// and ModeDurationSec are nil without an open period.
type MachineModeStatus struct {
	MachineID       int64      `db:"machine_id" json:"machineId"`
	MachineName     string     `db:"machine_name" json:"machineName"`
	CustomerName    string     `db:"customer_name" json:"customerName"`
	CompanyName     string     `db:"company_name" json:"companyName"`
	ModeID          *int       `db:"mode_id" json:"modeId"`
	ModeName        string     `db:"mode_name" json:"modeName"`
	OnlineStatus    bool       `db:"online_status" json:"onlineStatus"`
	ProgramName     string     `db:"program_name" json:"programName"`
	ModeSince       *time.Time `db:"mode_since" json:"modeSince"`
	ModeDurationSec *int       `db:"mode_duration_sec" json:"modeDurationSec"`
}

// MachineEvent is one row of machine_events. MachineID is nil for collector events.
type MachineEvent struct {
	MachineID *int64
	Shard     int
	Type      string
	Code      *int
	Message   string
	At        time.Time
}

// MachineLiveStatus is one row of machine_live_status.
type MachineLiveStatus struct {
	MachineID      int64
	Shard          int
	IsConnected    bool
	ModeID         int
	ConnectedSince *time.Time
	LastSeenAt     *time.Time
	LastChangeAt   *time.Time
	SampleAt       *time.Time
	MissingNodes   string
	CallbackCount  int64
	StoredCount    int64
	DroppedCount   int64
	ErrorCount     int64
	LastError      string
}
