package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config holds settings read from environment variables.
type Config struct {
	Port string
	DB   DBConfig
	// CookieSecure sends the session cookie over HTTPS only. Set
	// COOKIE_SECURE=true whenever the API is served over HTTPS.
	CookieSecure bool
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
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
