package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/detection"
	"github.com/ankush-dev-eng/lanternfish/pkg/eventbus"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
	"github.com/ankush-dev-eng/lanternfish/pkg/server"
	"github.com/ankush-dev-eng/lanternfish/pkg/simulator"
)

func TestEndToEndScenarios(t *testing.T) {
	scenarios := []struct {
		name              model.ScenarioType
		expectedTrips     int
		expectedMinAlerts int
	}{
		{model.ScenarioStealthy, 3, 3},
		{model.ScenarioNoisy, 3, 3},
		{model.ScenarioInsider, 3, 3},
	}

	for _, sc := range scenarios {
		t.Run(string(sc.name), func(t *testing.T) {
			graph := network.GenerateEnterpriseTopology()
			decoys := decoy.NewRegistry()
			bus := eventbus.NewBus(1000)
			detector := detection.NewEngine(graph, decoys)
			simEngine := simulator.NewEngine(bus, decoys, detector)

			srv := server.New(server.Config{
				Graph:     graph,
				Decoys:    decoys,
				Bus:       bus,
				Detector:  detector,
				SimEngine: simEngine,
			})
			handler := srv.Handler()

			// 1. Trigger Simulation via POST /api/simulate with 100x speed
			body, _ := json.Marshal(map[string]interface{}{
				"scenario": string(sc.name),
				"speed":    100.0,
			})
			req := httptest.NewRequest("POST", "/api/simulate", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected 200 on /api/simulate, got %d", w.Code)
			}

			// Wait for simulation to finish execution
			time.Sleep(300 * time.Millisecond)

			state := simEngine.GetState()
			if state.Status != model.StatusCompleted {
				t.Fatalf("expected status completed, got %s", state.Status)
			}

			// 2. Query /api/alerts
			req = httptest.NewRequest("GET", "/api/alerts", nil)
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected 200 on /api/alerts, got %d", w.Code)
			}
			var alerts []model.Alert
			if err := json.Unmarshal(w.Body.Bytes(), &alerts); err != nil {
				t.Fatalf("failed to decode alerts: %v", err)
			}
			if len(alerts) < sc.expectedMinAlerts {
				t.Fatalf("expected >= %d alerts, got %d", sc.expectedMinAlerts, len(alerts))
			}

			// 3. Query /api/report.md
			req = httptest.NewRequest("GET", "/api/report.md", nil)
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected 200 on /api/report.md, got %d", w.Code)
			}
			mdReport := w.Body.String()
			if !strings.Contains(mdReport, "Executive Summary") {
				t.Fatal("expected Executive Summary in Markdown report")
			}
			if !strings.Contains(mdReport, "Root Cause Analysis") {
				t.Fatal("expected Root Cause Analysis in Markdown report")
			}
			if !strings.Contains(mdReport, "Containment Recommendations") {
				t.Fatal("expected Containment Recommendations in Markdown report")
			}

			// 4. Query /api/report.json
			req = httptest.NewRequest("GET", "/api/report.json", nil)
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected 200 on /api/report.json, got %d", w.Code)
			}
			var rep model.IncidentReport
			if err := json.Unmarshal(w.Body.Bytes(), &rep); err != nil {
				t.Fatalf("failed to decode report.json: %v", err)
			}
			if len(rep.Timeline) == 0 {
				t.Fatal("expected non-empty timeline in incident report")
			}
			if len(rep.Recommendations) == 0 {
				t.Fatal("expected containment recommendations in report")
			}
		})
	}
}
