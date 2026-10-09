package config

import "os"

// Config holds settings read from environment variables.
type Config struct {
	Port string
}

// Load reads the config from the environment, falling back to defaults.
func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
