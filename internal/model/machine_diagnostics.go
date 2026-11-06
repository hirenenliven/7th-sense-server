package model

import "time"

// EngineNames names the CNC engines (controllers.cnc_engine).
var EngineNames = map[int]string{
	EngineSimulator:  "Simulator",
	EngineFanuc:      "Fanuc",
	EngineSiemens:    "Siemens",
	EngineSyntec:     "Syntec",
	EngineMitsubishi: "Mitsubishi",
}

// DiagnosticLiveStatus is a machine's machine_live_status row: what its collector worker last
// reported. UpdatedAt moves every 10 s while the worker holds a connection.
type DiagnosticLiveStatus struct {
	IsConnected    bool       `db:"is_connected" json:"isConnected"`
	ConnectedSince *time.Time `db:"connected_since" json:"connectedSince"`
	LastSeenAt     *time.Time `db:"last_seen_at" json:"lastSeenAt"`
	LastChangeAt   *time.Time `db:"last_change_at" json:"lastChangeAt"`
	CallbackCount  int64      `db:"callback_count" json:"callbackCount"`
	StoredCount    int64      `db:"stored_count" json:"storedCount"`
	DroppedCount   int64      `db:"dropped_count" json:"droppedCount"`
	ErrorCount     int64      `db:"error_count" json:"errorCount"`
	LastError      string     `db:"last_error" json:"lastError"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}

// DiagnosticMode is the mode a machine is in now: its open machine_mode_log period.
type DiagnosticMode struct {
	ModeID      int       `db:"mode_id" json:"modeId"`
	ModeName    string    `db:"mode_name" json:"modeName"`
	ProgramName string    `db:"program_name" json:"programName"`
	Since       time.Time `db:"started_at" json:"since"`
}

// DiagnosticVariable is one active cnc_variables row as this machine would read it, with the last
// number stored for it. KnownSources lists every source that has an address for the machine's
// engine, to say why a variable without a readable Address cannot be read.
type DiagnosticVariable struct {
	ID           int        `db:"id"`
	Code         string     `db:"code"`
	Name         string     `db:"name"`
	Unit         string     `db:"unit"`
	Category     string     `db:"category"`
	DataType     string     `db:"data_type"`
	Enabled      bool       `db:"enabled"`
	Source       string     `db:"source"`
	Address      string     `db:"address"`
	KnownSources string     `db:"known_sources"`
	StoredValue  *float64   `db:"stored_value"`
	StoredAt     *time.Time `db:"stored_at"`
}

// DiagnosticEvent is one machine_events row of a machine.
type DiagnosticEvent struct {
	ID      int64     `db:"id" json:"id"`
	Type    string    `db:"event_type" json:"type"`
	Code    *int      `db:"code" json:"code"`
	Message string    `db:"message" json:"message"`
	At      time.Time `db:"created_at" json:"at"`
}
