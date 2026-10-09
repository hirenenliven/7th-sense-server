package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds settings read from environment variables.
type Config struct {
	Port string
}

// Load reads the config from the environment, falling back to defaults.
// It loads variables from .env file if present.
func Load() Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env, using system environment variables")
	}

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
