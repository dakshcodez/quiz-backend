package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration loaded from environment.
// All production config should flow through here; no scattered os.Getenv in handlers.
type Config struct {
	// Server
	Port string

	// Database (Supabase PostgreSQL connection string)
	DatabaseURL string

	// JWT
	JWTSecret   string
	JWTExpiry   time.Duration // e.g. 24 * time.Hour

	// Query timeout for DB operations
	DBTimeout time.Duration
}

// Load reads configuration from environment variables.
// Call during startup; log.Fatal on missing required values.
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, ErrMissingDatabaseURL
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, ErrMissingJWTSecret
	}

	jwtExpiryHours := 24
	if h := os.Getenv("JWT_EXPIRY_HOURS"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 {
			jwtExpiryHours = parsed
		}
	}

	dbTimeoutSec := 10
	if t := os.Getenv("DB_TIMEOUT_SEC"); t != "" {
		if parsed, err := strconv.Atoi(t); err == nil && parsed > 0 {
			dbTimeoutSec = parsed
		}
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		JWTExpiry:   time.Duration(jwtExpiryHours) * time.Hour,
		DBTimeout:   time.Duration(dbTimeoutSec) * time.Second,
	}, nil
}
