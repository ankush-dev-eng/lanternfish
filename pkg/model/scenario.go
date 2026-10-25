package model

import "time"

// ScenarioType designates the attack narrative script.
type ScenarioType string

const (
	ScenarioStealthy ScenarioType = "stealthy"
	ScenarioNoisy    ScenarioType = "noisy"
	ScenarioInsider  ScenarioType = "insider"
)

// ScenarioStatus indicates the runtime state of the attack simulator.
type ScenarioStatus string

const (
	StatusIdle      ScenarioStatus = "idle"
	StatusRunning   ScenarioStatus = "running"
	StatusPaused    ScenarioStatus = "paused"
	StatusCompleted ScenarioStatus = "completed"
)

// SimulationStep defines an individual attack progression stage.
type SimulationStep struct {
	StepNumber     int           `json:"step_number"`
	Delay          time.Duration `json:"delay"`
	Phase          string        `json:"phase"` // e.g. Initial Access, Discovery, Lateral Movement
	MitreTechnique string        `json:"mitre_technique"`
	SourceNode     string        `json:"source_node"`
	TargetNode     string        `json:"target_node"`
	Actor          string        `json:"actor"`
	EventType      EventType     `json:"event_type"`
	Description    string        `json:"description"`
	DecoyTripwire  string        `json:"decoy_tripwire,omitempty"`
	RawPayload     string        `json:"raw_payload,omitempty"`
}

// SimulationState holds live simulation telemetry and statistics.
type SimulationState struct {
	ScenarioID       ScenarioType   `json:"scenario_id"`
	ScenarioName     string         `json:"scenario_name"`
	Status           ScenarioStatus `json:"status"`
	SpeedMultiplier  float64        `json:"speed_multiplier"`
	CurrentStep      int            `json:"current_step"`
	TotalSteps       int            `json:"total_steps"`
	StartedAt        *time.Time     `json:"started_at,omitempty"`
	CompletedAt      *time.Time     `json:"completed_at,omitempty"`
	PatientZeroNode  string         `json:"patient_zero_node,omitempty"`
	ActiveAttackerAt string         `json:"active_attacker_at,omitempty"`
	TimeToDetectMs   int64          `json:"time_to_detect_ms"` // ms between start and first alert
	AlertsCount      int            `json:"alerts_count"`
	DecoysTriggered  int            `json:"decoys_triggered"`
	HopsCount        int            `json:"hops_count"`
}
