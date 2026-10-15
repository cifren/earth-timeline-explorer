package app

import (
	"slices"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("REDIS_ADDR", "redis:6379")
	t.Setenv("COMMONS_USER_AGENT", "FriseChrono/test")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,https://example.org")

	config, err := LoadConfig()

	if err != nil {
		t.Fatal(err)
	}
	if config.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want the default :8080", config.HTTPAddr)
	}
	if !slices.Equal(config.CORSAllowedOrigins, []string{"http://localhost:5173", "https://example.org"}) {
		t.Errorf("CORSAllowedOrigins = %v", config.CORSAllowedOrigins)
	}
}

func TestLoadConfigNamesMissingVariables(t *testing.T) {
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("COMMONS_USER_AGENT", "")

	_, err := LoadConfig()

	if err == nil || !strings.Contains(err.Error(), "REDIS_ADDR") || !strings.Contains(err.Error(), "COMMONS_USER_AGENT") {
		t.Fatalf("error = %v, want both missing variables named", err)
	}
}
