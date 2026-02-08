package handler

import (
	"net/http"

	"memwright/api/pkg/logger"
)

// Dependencies holds all handler dependencies.
type Dependencies struct {
	Logger      logger.Logger
	Environment string
}

// RegisterRoutes registers all API routes on the given mux.
func RegisterRoutes(mux *http.ServeMux, deps Dependencies) {
	healthHandler := NewHealthHandler(deps.Environment)

	mux.Handle("/health", healthHandler)
}
