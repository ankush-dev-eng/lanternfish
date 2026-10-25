package report

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestReportExporter(t *testing.T) {
	path := model.AttackerPath{
		PatientZeroID:       "ws-dev-02",
		PatientZeroHostname: "ws-dev-02.corp",
		Hops: []model.ReconstructedHop{
			{HopIndex: 1, FromNode: "ws-dev-02", ToNode: "api-gw-01", Method: "SSH", MitreTechnique: "T1021"},
		},
		NodesVisited:     []string{"ws-dev-02", "api-gw-01"},
		DecoysTouched:    []string{"TRIP-FILE-PASSWORDS"},
		DwellTimeSeconds: 12.3,
	}

	decoys := []model.Decoy{
		{
			ID:             "decoy-1",
			Name:           "passwords.xlsx",
			Type:           model.DecoyTypeFile,
			NodeHostname:   "ws-dev-02.corp",
			TripwireID:     "TRIP-FILE-PASSWORDS",
			BaitEnticement: "Cleartext credentials",
			Triggered:      true,
		},
	}

	alerts := []model.Alert{
		{
			ID:             "alt-1",
			RuleName:       "DecoyFileRule",
			Severity:       model.SeverityHigh,
			MitreTechnique: "T1083",
		},
	}

	timeline := []model.TimelineItem{
		{
			ID:             "tl-1",
			Timestamp:      time.Now(),
			Phase:          model.PhaseInitialAccess,
			Severity:       model.SeverityHigh,
			Title:          "Phishing Foothold",
			MitreTechnique: "T1566",
		},
	}

	rep := BuildReport(model.ScenarioStealthy, path, alerts, decoys, timeline, map[string]interface{}{"ttd_ms": 1200})

	// 1. Test JSON export
	jsonBytes, err := ExportJSON(rep)
	if err != nil {
		t.Fatalf("failed to export JSON: %v", err)
	}
	var unmarshaled model.IncidentReport
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal exported JSON: %v", err)
	}
	if unmarshaled.Scenario != model.ScenarioStealthy {
		t.Errorf("expected ScenarioStealthy, got %s", unmarshaled.Scenario)
	}

	// 2. Test Markdown export
	md := ExportMarkdown(rep)
	if !strings.Contains(md, "LANTERNFISH Incident Forensics Report") {
		t.Fatal("expected report title in Markdown output")
	}
	if !strings.Contains(md, "ws-dev-02") {
		t.Fatal("expected patient zero in Markdown output")
	}
	if !strings.Contains(md, "Executive Summary") {
		t.Fatal("expected Executive Summary section in Markdown output")
	}

	// 3. Test HTML export
	html := ExportPrintableHTML(rep)
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Fatal("expected DOCTYPE in HTML export")
	}
	if !strings.Contains(html, "window.print()") {
		t.Fatal("expected print handler in HTML export")
	}
}
