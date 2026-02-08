package unit

import (
	"os"
	"testing"

	"memwright/api/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	clearEnvVars()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Errorf("expected default environment 'development', got '%s'", cfg.Environment)
	}

	if cfg.DatabaseHost != "localhost" {
		t.Errorf("expected default database host 'localhost', got '%s'", cfg.DatabaseHost)
	}

	if cfg.DatabasePort != 5432 {
		t.Errorf("expected default database port 5432, got %d", cfg.DatabasePort)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	clearEnvVars()

	os.Setenv("PORT", "3000")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("DATABASE_HOST", "db.example.com")
	os.Setenv("DATABASE_PORT", "5433")
	defer clearEnvVars()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Port)
	}

	if cfg.Environment != "production" {
		t.Errorf("expected environment 'production', got '%s'", cfg.Environment)
	}

	if cfg.DatabaseHost != "db.example.com" {
		t.Errorf("expected database host 'db.example.com', got '%s'", cfg.DatabaseHost)
	}

	if cfg.DatabasePort != 5433 {
		t.Errorf("expected database port 5433, got %d", cfg.DatabasePort)
	}
}

func TestIsDevelopment(t *testing.T) {
	clearEnvVars()
	os.Setenv("ENVIRONMENT", "development")
	defer clearEnvVars()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.IsDevelopment() {
		t.Error("expected IsDevelopment() to return true for development environment")
	}

	if cfg.IsProduction() {
		t.Error("expected IsProduction() to return false for development environment")
	}
}

func TestIsProduction(t *testing.T) {
	clearEnvVars()
	os.Setenv("ENVIRONMENT", "production")
	defer clearEnvVars()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.IsProduction() {
		t.Error("expected IsProduction() to return true for production environment")
	}

	if cfg.IsDevelopment() {
		t.Error("expected IsDevelopment() to return false for production environment")
	}
}

func TestDatabaseURL(t *testing.T) {
	clearEnvVars()
	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("DATABASE_USER", "testuser")
	os.Setenv("DATABASE_PASSWORD", "testpass")
	os.Setenv("DATABASE_NAME", "testdb")
	os.Setenv("DATABASE_SSLMODE", "require")
	defer clearEnvVars()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require"
	if cfg.DatabaseURL() != expected {
		t.Errorf("expected database URL '%s', got '%s'", expected, cfg.DatabaseURL())
	}
}

func TestInvalidPortFallsBackToDefault(t *testing.T) {
	clearEnvVars()
	os.Setenv("PORT", "invalid")
	defer clearEnvVars()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port to fall back to 8080 for invalid value, got %d", cfg.Port)
	}
}

func clearEnvVars() {
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("DATABASE_HOST")
	os.Unsetenv("DATABASE_PORT")
	os.Unsetenv("DATABASE_USER")
	os.Unsetenv("DATABASE_PASSWORD")
	os.Unsetenv("DATABASE_NAME")
	os.Unsetenv("DATABASE_SSLMODE")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRATION_HOURS")
}
