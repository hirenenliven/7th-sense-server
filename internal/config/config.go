package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds settings read from environment variables.
type Config struct {
	Port string
	DB   DBConfig
	// CookieSecure sends the session cookie over HTTPS only. Set
	// COOKIE_SECURE=true whenever the API is served over HTTPS.
	CookieSecure bool
	// UploadDir is the directory uploaded files are kept in (part drawings under
	// UploadDir/part-drawings). It must survive restarts and deploys.
	UploadDir string
	// CNC holds the settings of the CNC data collectors (docs/cnc-telemetry-collector-plan.md).
	CNC CNCConfig
}

// CNCConfig holds the CNC_* settings of the data collection the API server runs in-process.
type CNCConfig struct {
	// Slots is the number of DLL callback pairs: the most machines the server can hold (max 990).
	Slots int
	// LibDir and DLLName locate the CNC DLL (LibDir is relative to the server binary or the
	// working directory unless absolute).
	LibDir  string
	DLLName string
	// Stale: no usable reading for this long, while the DLL's last word was "not connected" (or it
	// sent nothing yet), means the connection is lost: the machine goes Offline and is reconnected.
	Stale time.Duration
	// Quiet: no callback at all for this long after a connected reading. The DLL stops calling
	// back while nothing changes (a held or idle Siemens machine sends nothing), so quiet is not
	// lost; after Quiet the connection is reopened in case the DLL thread hung, keeping the mode.
	Quiet time.Duration
	// MaxConnecting limits concurrent CNCConnect calls per collector.
	MaxConnecting int
	// StopTimeout bounds the collection's graceful stop.
	StopTimeout time.Duration
	// CycleMaxDuration: a cycle open this long is closed as unfinished (timeout), so a machine
	// left in Hold for days does not distort the reports.
	CycleMaxDuration time.Duration
	// Machines (CNC_MACHINE_IDS) are the machines the server connects: "225,247" only these, the
	// others are not connected (no data, mode or status) — for analysing a few machines; "ALL" or
	// empty every active machine.
	Machines MachineSelection
	// Logs picks which CNC logs are written (monitor/cnc_log.go).
	Logs CNCLogSettings
	// CertDir is the directory where client certificates are generated for OPC UA authentication.
	// If empty, certificates are generated next to the server binary.
	CertDir string
}

// DBConfig holds the PostgreSQL connection settings.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Load reads the config from the environment, falling back to defaults.
// Variables from a .env file are loaded first if the file exists.
func Load() Config {
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using system environment variables")
	}

	return Config{
		Port: getEnv("PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_DATABASE", "7th_sense"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		CookieSecure: getEnv("COOKIE_SECURE", "false") == "true",
		UploadDir:    getEnv("UPLOAD_DIR", "uploads"),
		CNC: CNCConfig{
			Slots:            getEnvInt("CNC_SLOTS", 300),
			LibDir:           getEnv("CNC_LIB_DIR", "Lib"),
			DLLName:          getEnv("CNC_DLL_NAME", "CncCommunication.dll"),
			Stale:            time.Duration(getEnvInt("CNC_STALE_SEC", 30)) * time.Second,
			Quiet:            time.Duration(getEnvInt("CNC_QUIET_SEC", 300)) * time.Second,
			MaxConnecting:    getEnvInt("CNC_MAX_CONNECTING", 20),
			StopTimeout:      time.Duration(getEnvInt("CNC_STOP_TIMEOUT_SEC", 15)) * time.Second,
			CycleMaxDuration: time.Duration(getEnvInt("CNC_CYCLE_MAX_HOURS", 24)) * time.Hour,
			Machines:         getEnvMachineSelection("CNC_MACHINE_IDS"),
			Logs: CNCLogSettings{
				Events:     getEnvBool("CNC_LOG_EVENTS", true),
				Modes:      getEnvBool("CNC_LOG_MODES", true),
				ModeChecks: getEnvBool("CNC_LOG_MODE_CHECKS", true),
				Variables:  getEnvBool("CNC_LOG_VARIABLES", false),
				Other:      getEnvBool("CNC_LOG_OTHER", true),
				Snapshot:   getEnvBool("CNC_LOG_SNAPSHOT", false),
				Downtime:   getEnvBool("CNC_LOG_DOWNTIME", true),
				Cycles:     getEnvBool("CNC_LOG_CYCLES", true),
			},
			CertDir: getEnv("CNC_CERT_DIR", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

// getEnvInt reads a positive whole number; a missing or invalid value falls back.
func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(getEnv(key, ""))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

// getEnvBool reads "true" or "false" (any case); a missing or other value falls back.
func getEnvBool(key string, fallback bool) bool {
	switch strings.ToLower(getEnv(key, "")) {
	case "true":
		return true
	case "false":
		return false
	}
	return fallback
}

// CNCLogSettings switch the CNC logs on and off (CNC_LOG_EVENTS, CNC_LOG_MODES, CNC_LOG_SNAPSHOT,
// CNC_LOG_MODE_CHECKS, CNC_LOG_VARIABLES, CNC_LOG_OTHER, CNC_LOG_DOWNTIME); monitor/cnc_log.go
// writes them.
// Errors are written whatever these say.
type CNCLogSettings struct {
	// Events logs machine events: connect, disconnect, alarms.
	Events bool
	// Modes logs mode changes.
	Modes bool
	// ModeChecks logs, whenever the NC inputs change, the mode they decide (even when the mode
	// stays the same).
	ModeChecks bool
	// Variables logs the changed variable values of each reading in one line; at 10 readings/s
	// it writes many lines per machine.
	Variables bool
	// Other logs everything else: collector start / stop, connections, variables loaded,
	// reading rate every 10 s per machine.
	Other bool
	// Snapshot logs the raw DLL snapshot status fields (mode, program, override, alarm, E-stop,
	// operator message, …) whenever one changes; for analysing what a controller sends.
	Snapshot bool
	// Downtime logs the inputs of the Downtime rule (machine type, alarm, E-stop, program
	// status) and whether they make the machine Downtime, whenever one of them changes.
	Downtime bool
	// Cycles logs each cycle opened, closed and restored, with its seconds per mode and parts.
	Cycles bool
}

// AllMachines is the CNC_MACHINE_IDS value that selects every machine (as does empty).
const AllMachines = "ALL"

// MachineSelection is a set of machines picked in a setting; no ids means every machine.
type MachineSelection struct {
	IDs map[int64]bool
}

// Includes reports whether machine id is selected.
func (selection MachineSelection) Includes(id int64) bool {
	return selection.IsAll() || selection.IDs[id]
}

// IsAll reports whether every machine is selected ("ALL" or empty).
func (selection MachineSelection) IsAll() bool {
	return len(selection.IDs) == 0
}

// getEnvMachineSelection reads "ALL" (any case), empty, or positive machine ids separated by
// commas or spaces ("225, 247"); an invalid id is logged and left out.
func getEnvMachineSelection(key string) MachineSelection {
	selection := MachineSelection{IDs: map[int64]bool{}}
	value := strings.TrimSpace(getEnv(key, ""))
	if strings.EqualFold(value, AllMachines) {
		return selection
	}
	for _, entry := range strings.FieldsFunc(value, func(character rune) bool {
		return character == ',' || character == ' '
	}) {
		id, err := strconv.ParseInt(entry, 10, 64)
		if err != nil || id <= 0 {
			slog.Warn("invalid machine id ignored", "variable", key, "value", entry)
			continue
		}
		selection.IDs[id] = true
	}
	return selection
}
