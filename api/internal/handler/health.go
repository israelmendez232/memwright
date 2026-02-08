package handler

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
}

type HealthHandler struct {
	environment string
}

func NewHealthHandler(environment string) *HealthHandler {
	return &HealthHandler{
		environment: environment,
	}
}

func (handler *HealthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:      "ok",
		Environment: handler.environment,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}
