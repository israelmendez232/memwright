package unit

import (
	"os"
	"path/filepath"
	"testing"

	"memwright/api/pkg/env"
)

func TestLoadEnvFile(t *testing.T) {
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	content := `
# Comment line
PORT=9000
ENVIRONMENT=testing
DATABASE_HOST="quoted-host"
DATABASE_PASSWORD='single-quoted'
EMPTY_VALUE=
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp env file: %v", err)
	}

	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("DATABASE_HOST")
	os.Unsetenv("DATABASE_PASSWORD")
	os.Unsetenv("EMPTY_VALUE")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("DATABASE_PASSWORD")
		os.Unsetenv("EMPTY_VALUE")
	}()

	if err := env.Load(envFile); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}

	if os.Getenv("PORT") != "9000" {
		t.Errorf("expected PORT to be '9000', got '%s'", os.Getenv("PORT"))
	}

	if os.Getenv("ENVIRONMENT") != "testing" {
		t.Errorf("expected ENVIRONMENT to be 'testing', got '%s'", os.Getenv("ENVIRONMENT"))
	}

	if os.Getenv("DATABASE_HOST") != "quoted-host" {
		t.Errorf("expected DATABASE_HOST to be 'quoted-host', got '%s'", os.Getenv("DATABASE_HOST"))
	}

	if os.Getenv("DATABASE_PASSWORD") != "single-quoted" {
		t.Errorf("expected DATABASE_PASSWORD to be 'single-quoted', got '%s'", os.Getenv("DATABASE_PASSWORD"))
	}
}

func TestLoadEnvFileDoesNotOverrideExisting(t *testing.T) {
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	content := `PORT=9000`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp env file: %v", err)
	}

	os.Setenv("PORT", "existing_value")
	defer os.Unsetenv("PORT")

	if err := env.Load(envFile); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}

	if os.Getenv("PORT") != "existing_value" {
		t.Errorf("expected PORT to remain 'existing_value', got '%s'", os.Getenv("PORT"))
	}
}

func TestLoadNonExistentFileDoesNotFail(t *testing.T) {
	if err := env.Load("/nonexistent/path/.env"); err != nil {
		t.Errorf("expected no error for nonexistent file, got: %v", err)
	}
}
