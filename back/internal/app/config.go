// Package app is the composition root: it reads the configuration, wires the bounded
// contexts together and runs the HTTP server.
package app

import (
	"errors"
	"os"
	"strings"
)

// Config is read once from the environment at startup.
type Config struct {
	HTTPAddr           string
	RedisAddr          string
	CommonsUserAgent   string
	CORSAllowedOrigins []string
}

// LoadConfig reads the configuration from the environment variables.
func LoadConfig() (Config, error) {
	config := Config{
		HTTPAddr:         envOr("HTTP_ADDR", ":8080"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		CommonsUserAgent: os.Getenv("COMMONS_USER_AGENT"),
	}
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		config.CORSAllowedOrigins = strings.Split(origins, ",")
	}

	var missing []string
	if config.RedisAddr == "" {
		missing = append(missing, "REDIS_ADDR")
	}
	if config.CommonsUserAgent == "" {
		missing = append(missing, "COMMONS_USER_AGENT")
	}
	if len(missing) > 0 {
		return Config{}, errors.New("missing environment variables: " + strings.Join(missing, ", "))
	}
	return config, nil
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
