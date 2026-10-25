package model

import "time"

// ContainmentRecommendation outlines actionable steps for the incident response team.
type ContainmentRecommendation struct {
	Priority    string `json:"priority"` // IMMEDIATE, HIGH, MEDIUM
	Target      string `json:"target"`   // e.g. "ws-dev-02 (10.10.40.12)"
	Action      string `json:"action"`   // e.g. "Isolate host from VLAN"
	Rationale   string `json:"rationale"`
	MitreRemedy string `json:"mitre_remedy"`
}

// IncidentReport is the comprehensive forensic summary exported via API or Markdown.
type IncidentReport struct {
	ID                string                      `json:"id"`
	GeneratedAt       time.Time                   `json:"generated_at"`
	Scenario          ScenarioType                `json:"scenario"`
	Status            string                      `json:"status"`
	ExecutiveSummary  string                      `json:"executive_summary"`
	RootCauseAnalysis string                      `json:"root_cause_analysis"`
	IsAIGenerated     bool                        `json:"is_ai_generated"`
	NarrativeModel    string                      `json:"narrative_model"` // e.g. "Gemini 2.5 Flash" or "LANTERNFISH Deterministic Forensics Engine"
	AttackerPath      AttackerPath                `json:"attacker_path"`
	Alerts            []Alert                     `json:"alerts"`
	DecoysTriggered   []Decoy                     `json:"decoys_triggered"`
	Timeline          []TimelineItem              `json:"timeline"`
	Recommendations   []ContainmentRecommendation `json:"recommendations"`
	Metrics           map[string]interface{}      `json:"metrics"`
}
