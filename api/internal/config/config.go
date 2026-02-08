package config

import (
	"fmt"
	"os"
	"strconv"

	"memwright/api/pkg/env"
)

type Config struct {
	Port        int
	Environment string
	LogLevel    string

	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string

	JWTSecret          string
	JWTExpirationHours int
}

func Load() (*Config, error) {
	if err := env.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	return &Config{
		Port:        getEnvInt("PORT", 8080),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		DatabaseHost:     getEnv("POSTGRES_HOST", "localhost"),
		DatabasePort:     getEnvInt("POSTGRES_PORT", 5432),
		DatabaseUser:     getEnv("POSTGRES_USER", "memwright"),
		DatabasePassword: getEnv("POSTGRES_PASSWORD", ""),
		DatabaseName:     getEnv("POSTGRES_DB", "memwright"),
		DatabaseSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),

		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
	}, nil
}

func (config *Config) IsDevelopment() bool {
	return config.Environment == "development"
}

func (config *Config) IsProduction() bool {
	return config.Environment == "production"
}

func (config *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		config.DatabaseSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
