package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAndOptionsPreflight(t *testing.T) {
	srv := setupTestServer()
	handler := srv.Handler()

	// Test OPTIONS preflight
	req := httptest.NewRequest("OPTIONS", "/api/simulate", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for OPTIONS, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected Access-Control-Allow-Origin '*', got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("expected Access-Control-Allow-Methods header to be set")
	}
}
