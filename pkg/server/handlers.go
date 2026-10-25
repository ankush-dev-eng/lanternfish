package server

import (
	"encoding/json"
	"net/http"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/reconstruction"
	"github.com/ankush-dev-eng/lanternfish/pkg/report"
	"github.com/ankush-dev-eng/lanternfish/pkg/timeline"
)

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","app":"lanternfish","version":"1.0.0"}`))
}

func (s *Server) handleGetNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.graph)
}

func (s *Server) handleGetDecoys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoys := s.decoys.GetAll()
	json.NewEncoder(w).Encode(decoys)
}

func (s *Server) handleGetAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	alerts := s.detector.GetAlerts()
	json.NewEncoder(w).Encode(alerts)
}

func (s *Server) handleGetTimeline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events := s.bus.GetEvents()
	alerts := s.detector.GetAlerts()
	tl := timeline.BuildTimeline(events, alerts, s.decoys)
	json.NewEncoder(w).Encode(tl)
}

func (s *Server) handleGetReportJSON(w http.ResponseWriter, r *http.Request) {
	rep := s.compileIncidentReport()
	data, err := report.ExportJSON(rep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=\"lanternfish-report.json\"")
	w.Write(data)
}

func (s *Server) handleGetReportMD(w http.ResponseWriter, r *http.Request) {
	rep := s.compileIncidentReport()
	md := report.ExportMarkdown(rep)
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"lanternfish-report.md\"")
	w.Write([]byte(md))
}

func (s *Server) handleGetReportHTML(w http.ResponseWriter, r *http.Request) {
	rep := s.compileIncidentReport()
	html := report.ExportPrintableHTML(rep)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

type simulateRequest struct {
	Scenario string  `json:"scenario"`
	Speed    float64 `json:"speed"`
}

func (s *Server) handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req simulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Default to stealthy if empty
		req.Scenario = string(model.ScenarioStealthy)
		req.Speed = 1.0
	}

	scenarioType := model.ScenarioType(req.Scenario)
	if scenarioType == "" {
		scenarioType = model.ScenarioStealthy
	}

	speed := req.Speed
	if speed <= 0 {
		speed = 1.0
	}

	if err := s.simEngine.Start(scenarioType, speed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.simEngine.GetState())
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.simEngine.Reset()
	s.detector.Reset()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"reset","message":"Simulation, alerts, and decoys reinitialized"}`))
}

func (s *Server) handleGetTrace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events := s.bus.GetEvents()
	alerts := s.detector.GetAlerts()
	path := reconstruction.ReconstructPath(s.graph, events, alerts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(path)
}

func (s *Server) compileIncidentReport() model.IncidentReport {
	events := s.bus.GetEvents()
	alerts := s.detector.GetAlerts()
	path := reconstruction.ReconstructPath(s.graph, events, alerts)
	tl := timeline.BuildTimeline(events, alerts, s.decoys)
	allDecoys := s.decoys.GetAll()
	state := s.simEngine.GetState()

	metrics := map[string]interface{}{
		"time_to_detect_ms": state.TimeToDetectMs,
		"hops_count":        state.HopsCount,
		"alerts_count":      len(alerts),
		"decoys_triggered":  state.DecoysTriggered,
		"status":            state.Status,
	}

	return report.BuildReport(state.ScenarioID, path, alerts, allDecoys, tl, metrics)
}
