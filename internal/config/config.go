// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"os"
)

// Config holds all application settings.
type Config struct {
	DatabaseURL string
	ServerAddr  string
}

// Load reads configuration from environment variables.
// DATABASE_URL is required, SERVER_ADDR defaults to ":8080".
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return &Config{
		DatabaseURL: dbURL,
		ServerAddr:  addr,
	}, nil
}
