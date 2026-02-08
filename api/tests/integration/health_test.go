package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"memwright/api/internal/handler"
)

func TestHealthEndpointIntegration(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/health", handler.NewHealthHandler("test"))

	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}

	var body struct {
		Status      string `json:"status"`
		Environment string `json:"environment"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body.Status)
	}

	if body.Environment != "test" {
		t.Errorf("expected environment 'test', got '%s'", body.Environment)
	}
}

func TestHealthEndpointReadiness(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/health", handler.NewHealthHandler("development"))

	server := httptest.NewServer(mux)
	defer server.Close()

	// Simulate multiple health checks as container orchestrators would do
	for i := 0; i < 5; i++ {
		response, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("health check %d failed: %v", i+1, err)
		}
		response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("health check %d: expected status %d, got %d", i+1, http.StatusOK, response.StatusCode)
		}
	}
}
