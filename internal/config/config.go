package config

import (
	"log/slog"
	"os"
	"strconv"
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
	// Enabled makes the API server connect the machines and collect their data. false = API only.
	Enabled bool
	// Slots is the number of DLL callback pairs: the most machines the server can hold (max 990).
	Slots int
	// LibDir and DLLName locate the CNC DLL (LibDir is relative to the server binary or the
	// working directory unless absolute).
	LibDir  string
	DLLName string
	// Stale: no usable reading for this long means the connection is lost.
	Stale time.Duration
	// MaxConnecting limits concurrent CNCConnect calls per collector.
	MaxConnecting int
	// StopTimeout bounds the collection's graceful stop.
	StopTimeout time.Duration
	// LogValues logs the changed variable values of each reading in one line (monitor/value_log.go). For
	// analysing data frequency only: at 10 readings/s it writes many lines per machine.
	LogValues bool
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
			Enabled:       getEnv("CNC_ENABLED", "true") == "true",
			Slots:         getEnvInt("CNC_SLOTS", 300),
			LibDir:        getEnv("CNC_LIB_DIR", "Lib"),
			DLLName:       getEnv("CNC_DLL_NAME", "FanucCncCommunication.dll"),
			Stale:         time.Duration(getEnvInt("CNC_STALE_SEC", 30)) * time.Second,
			MaxConnecting: getEnvInt("CNC_MAX_CONNECTING", 20),
			StopTimeout:   time.Duration(getEnvInt("CNC_STOP_TIMEOUT_SEC", 15)) * time.Second,
			LogValues:     getEnv("CNC_LOG_VALUES", "false") == "true",
			CertDir:       getEnv("CNC_CERT_DIR", ""),
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
