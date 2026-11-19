package dto

import (
	"time"

	"7th-sense-server/internal/model"
)

// Results of one diagnostics check.
const (
	DiagnosticOK      = "ok"
	DiagnosticWarning = "warning"
	DiagnosticFailed  = "failed"
	DiagnosticSkipped = "skipped" // an earlier check decides it, or there is nothing to check
)

// Keys of the checks GET /machines/{id}/diagnostics returns, in the order the page runs them.
const (
	DiagnosticCheckSetup      = "setup"
	DiagnosticCheckCollector  = "collector"
	DiagnosticCheckController = "controller"
	DiagnosticCheckMode       = "mode"
	DiagnosticCheckVariables  = "variables"
)

// DiagnosticCheck is the result of one check: what was found, and for a problem what to do.
type DiagnosticCheck struct {
	Key     string `json:"key"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

// MachineDiagnosticsResponse is the response of GET /machines/{id}/diagnostics: what the server
// knows about one machine's data collection now. It sends nothing to the machine, so the page can
// reload it often; GET /machines/{id}/diagnostics/network tests the network.
type MachineDiagnosticsResponse struct {
	ServerTime  time.Time                   `json:"serverTime"`
	Machine     DiagnosticMachine           `json:"machine"`
	LiveEnabled bool                        `json:"liveEnabled"` // LIVE_ENABLED: the socket exists
	QuietSec    int                         `json:"quietSec"`    // CNC_QUIET_SEC
	Status      *model.DiagnosticLiveStatus `json:"status"`      // nil before the first report
	Mode        *model.DiagnosticMode       `json:"mode"`        // nil without an open period
	OpMode      string                      `json:"opMode"`
	Checks      []DiagnosticCheck           `json:"checks"`
	Variables   []DiagnosticVariableRow     `json:"variables"` // the CNC DLL's (snapshot) variables
	PLC         DiagnosticPLC               `json:"plc"`       // the OPC UA PLC reader's side
	Events      []model.DiagnosticEvent     `json:"events"`
}

// DiagnosticPLC is what the OPC UA PLC reader (ADR 0021) does for the machine: whether it reads
// it, its connection, a sentence on what is wrong with the hint to fix it, and every PLC variable.
type DiagnosticPLC struct {
	Enabled        bool                    `json:"enabled"`   // PLC_ENABLED
	Supported      bool                    `json:"supported"` // the controller runs an OPC UA server (Siemens)
	State          string                  `json:"state"`     // opcua.State*, "off" or "not_reported"
	Endpoint       string                  `json:"endpoint"`
	SecurityMode   string                  `json:"securityMode"`
	SecurityPolicy string                  `json:"securityPolicy"`
	ConnectedSince *time.Time              `json:"connectedSince"`
	LastError      string                  `json:"lastError"`
	LastErrorAt    *time.Time              `json:"lastErrorAt"`
	LastChangeAt   *time.Time              `json:"lastChangeAt"`
	ChangeCount    int64                   `json:"changeCount"`
	Status         string                  `json:"status"` // ok | warning | failed | skipped, as a DiagnosticCheck
	Message        string                  `json:"message"`
	Hint           string                  `json:"hint"`
	Variables      []DiagnosticPLCVariable `json:"variables"`
}

// DiagnosticPLCVariable is one PLC variable of the machine: its node id, what the controller says
// about it (nodeStatus "" = fine), the live value from the hub and the last stored number.
type DiagnosticPLCVariable struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Unit        string     `json:"unit"`
	Category    string     `json:"category"`
	DataType    string     `json:"dataType"`
	NodeID      string     `json:"nodeId"`
	Enabled     bool       `json:"enabled"`    // false = turned off for the machine (machine_variables)
	Subscribed  bool       `json:"subscribed"` // the controller accepted the node
	NodeStatus  string     `json:"nodeStatus"`
	Value       any        `json:"value"`
	StoredValue *float64   `json:"storedValue"`
	StoredAt    *time.Time `json:"storedAt"`
}

// DiagnosticMachine is what identifies the machine and how the server reaches it.
type DiagnosticMachine struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	CustomerName      string `json:"customerName"`
	CompanyName       string `json:"companyName"`
	Status            string `json:"status"`
	LiveIP            string `json:"liveIp"`
	LocalIP           string `json:"localIp"`
	Port              int    `json:"port"`
	Engine            string `json:"engine"` // "" when the controller has none
	BrandName         string `json:"brandName"`
	ControllerName    string `json:"controllerName"`
	ControllerVersion string `json:"controllerVersion"`
	ModelName         string `json:"modelName"`
}

// DiagnosticVariableRow is one variable of the page: whether the machine can read it (else
// Reason), the latest value the server holds in memory (nil = not read since the server started)
// and the last number stored in the database.
type DiagnosticVariableRow struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Unit        string     `json:"unit"`
	Category    string     `json:"category"`
	DataType    string     `json:"dataType"`
	Source      string     `json:"source"`
	Address     string     `json:"address"`
	Readable    bool       `json:"readable"`
	Reason      string     `json:"reason,omitempty"`
	Value       any        `json:"value"`
	StoredValue *float64   `json:"storedValue"`
	StoredAt    *time.Time `json:"storedAt"`
}

// What one IP and port test says about the machine (MachineNetworkResponse.MachineState).
const (
	MachineStateOn            = "powered_on"     // the port accepted a connection: the machine is on
	MachineStatePortClosed    = "port_closed"    // the IP answers, the port does not
	MachineStateIPUnreachable = "ip_unreachable" // no answer from the IP: off, or the network is down
	MachineStateIPInvalid     = "ip_invalid"     // the stored IP is not an address and does not resolve
)

// MachineNetworkResponse is the response of GET /machines/{id}/diagnostics/network: whether the
// machine's live IP and port answer the server, tested when asked.
type MachineNetworkResponse struct {
	CheckedAt    time.Time       `json:"checkedAt"`
	Address      string          `json:"address"` // ip:port as tested
	MachineState string          `json:"machineState" enums:"powered_on,port_closed,ip_unreachable,ip_invalid"`
	IP           DiagnosticCheck `json:"ip"`
	Port         DiagnosticCheck `json:"port"`
	LatencyMs    *int64          `json:"latencyMs"` // the TCP connect, nil when it failed
}
