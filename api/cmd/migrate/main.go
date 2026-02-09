package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"memwright/api/internal/config"
	"memwright/api/pkg/env"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if err := env.Load(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	migrationsPath := getMigrationsPath()

	switch os.Args[1] {
	case "up":
		runUp(cfg.DatabaseURL(), migrationsPath)
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			steps, _ = strconv.Atoi(os.Args[2])
		}
		runDown(cfg.DatabaseURL(), migrationsPath, steps)
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Error: migration name required")
			fmt.Println("Usage: migrate create <name>")
			os.Exit(1)
		}
		createMigration(migrationsPath, os.Args[2])
	case "version":
		showVersion(cfg.DatabaseURL(), migrationsPath)
	case "force":
		if len(os.Args) < 3 {
			fmt.Println("Error: version required")
			fmt.Println("Usage: migrate force <version>")
			os.Exit(1)
		}
		version, _ := strconv.Atoi(os.Args[2])
		forceVersion(cfg.DatabaseURL(), migrationsPath, version)
	default:
		printUsage()
		os.Exit(1)
	}
}

func getMigrationsPath() string {
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	candidates := []string{
		"migrations",
		"api/migrations",
		filepath.Join(execDir, "migrations"),
		filepath.Join(execDir, "..", "migrations"),
		filepath.Join(execDir, "..", "..", "migrations"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	absPath, _ := filepath.Abs("migrations")
	return absPath
}

func runUp(dbURL, migrationsPath string) {
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
			return
		}
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migrations applied successfully")
}

func runDown(dbURL, migrationsPath string, steps int) {
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Steps(-steps); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
			return
		}
		fmt.Printf("Error rolling back migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Rolled back %d migration(s)\n", steps)
}

func createMigration(migrationsPath, name string) {
	timestamp := time.Now().Format("20060102150405")
	baseName := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join(migrationsPath, baseName+".up.sql")
	downFile := filepath.Join(migrationsPath, baseName+".down.sql")

	if err := os.WriteFile(upFile, []byte("-- Migration: "+name+"\n\n"), 0644); err != nil {
		fmt.Printf("Error creating up migration: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(downFile, []byte("-- Rollback: "+name+"\n\n"), 0644); err != nil {
		fmt.Printf("Error creating down migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created migrations:\n  %s\n  %s\n", upFile, downFile)
}

func showVersion(dbURL, migrationsPath string) {
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations applied")
			return
		}
		fmt.Printf("Error getting version: %v\n", err)
		os.Exit(1)
	}

	dirtyStr := ""
	if dirty {
		dirtyStr = " (dirty)"
	}
	fmt.Printf("Current version: %d%s\n", version, dirtyStr)
}

func forceVersion(dbURL, migrationsPath string, version int) {
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Force(version); err != nil {
		fmt.Printf("Error forcing version: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Forced version to %d\n", version)
}

func printUsage() {
	fmt.Println("Usage: migrate <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up              Apply all pending migrations")
	fmt.Println("  down [n]        Rollback n migrations (default: 1)")
	fmt.Println("  create <name>   Create a new migration")
	fmt.Println("  version         Show current migration version")
	fmt.Println("  force <version> Force set the migration version")
}
