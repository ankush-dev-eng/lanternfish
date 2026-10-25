package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/detection"
	"github.com/ankush-dev-eng/lanternfish/pkg/eventbus"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
	"github.com/ankush-dev-eng/lanternfish/pkg/simulator"
)

func setupTestServer() *Server {
	graph := network.GenerateEnterpriseTopology()
	decoys := decoy.NewRegistry()
	bus := eventbus.NewBus(100)
	detector := detection.NewEngine(graph, decoys)
	simEngine := simulator.NewEngine(bus, decoys, detector)

	return New(Config{
		Graph:     graph,
		Decoys:    decoys,
		Bus:       bus,
		Detector:  detector,
		SimEngine: simEngine,
	})
}

func TestServerEndpoints(t *testing.T) {
	srv := setupTestServer()
	handler := srv.Handler()

	// 1. Healthz
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /healthz, got %d", w.Code)
	}

	// 2. /api/network
	req = httptest.NewRequest("GET", "/api/network", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/network, got %d", w.Code)
	}

	// 3. /api/decoys
	req = httptest.NewRequest("GET", "/api/decoys", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/decoys, got %d", w.Code)
	}

	// 4. /api/alerts
	req = httptest.NewRequest("GET", "/api/alerts", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/alerts, got %d", w.Code)
	}

	// 5. /api/timeline
	req = httptest.NewRequest("GET", "/api/timeline", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/timeline, got %d", w.Code)
	}

	// 6. /api/report.json
	req = httptest.NewRequest("GET", "/api/report.json", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/report.json, got %d", w.Code)
	}

	// 7. /api/report.md
	req = httptest.NewRequest("GET", "/api/report.md", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/report.md, got %d", w.Code)
	}

	// 8. /report
	req = httptest.NewRequest("GET", "/report", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /report, got %d", w.Code)
	}

	// 9. /api/simulate
	simBody, _ := json.Marshal(map[string]interface{}{
		"scenario": "stealthy",
		"speed":    10.0,
	})
	req = httptest.NewRequest("POST", "/api/simulate", bytes.NewBuffer(simBody))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/simulate, got %d", w.Code)
	}

	// 10. /api/reset
	req = httptest.NewRequest("POST", "/api/reset", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /api/reset, got %d", w.Code)
	}
}
