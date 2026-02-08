package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"memwright/api/internal/handler"
)

func TestHealthHandler_Success(t *testing.T) {
	healthHandler := handler.NewHealthHandler("development")

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response handler.HealthResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response.Status)
	}

	if response.Environment != "development" {
		t.Errorf("expected environment 'development', got '%s'", response.Environment)
	}
}

func TestHealthHandler_ProductionEnvironment(t *testing.T) {
	healthHandler := handler.NewHealthHandler("production")

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler.ServeHTTP(recorder, request)

	var response handler.HealthResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Environment != "production" {
		t.Errorf("expected environment 'production', got '%s'", response.Environment)
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	healthHandler := handler.NewHealthHandler("development")

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, "/health", nil)
			recorder := httptest.NewRecorder()

			healthHandler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status code %d for %s, got %d", http.StatusMethodNotAllowed, method, recorder.Code)
			}
		})
	}
}
