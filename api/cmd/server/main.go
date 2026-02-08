package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"memwright/api/internal/config"
	"memwright/api/internal/handler"
	"memwright/api/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	minLevel := logger.LevelInfo
	if cfg.IsDevelopment() {
		minLevel = logger.LevelDebug
	}
	appLogger := logger.New(nil, minLevel)

	appLogger.Info("starting memwright API server environment=%s port=%d", cfg.Environment, cfg.Port)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, handler.Dependencies{
		Logger:      appLogger,
		Environment: cfg.Environment,
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		appLogger.Info("server listening addr=%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("server error: %v", err)
		}
	}()

	// Block until SIGINT or SIGTERM is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down server...")

	// Allow 30 seconds for active connections to finish
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("server forced to shutdown: %v", err)
	}

	appLogger.Info("server exited gracefully")
}
