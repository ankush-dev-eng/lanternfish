package model

import "time"

// AlertSeverity indicates the threat level of an anomaly or deception tripwire.
type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "LOW"
	SeverityMedium   AlertSeverity = "MEDIUM"
	SeverityHigh     AlertSeverity = "HIGH"
	SeverityCritical AlertSeverity = "CRITICAL"
)

// Alert represents a high-fidelity security detection triggered by decoy tripwires or behavioral anomalies.
type Alert struct {
	ID             string        `json:"id"`
	RuleID         string        `json:"rule_id"`
	RuleName       string        `json:"rule_name"`
	Severity       AlertSeverity `json:"severity"`
	Confidence     int           `json:"confidence"` // 0-100 score
	MitreTechnique string        `json:"mitre_technique"` // e.g. T1078, T1083, T1021, T1041
	MitreName      string        `json:"mitre_name"`
	MitreTactic    MitreTactic   `json:"mitre_tactic,omitempty"`
	MitreURL       string        `json:"mitre_url,omitempty"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Timestamp      time.Time     `json:"timestamp"`
	SourceNode     string        `json:"source_node"`
	SourceHostname string        `json:"source_hostname"`
	TargetNode     string        `json:"target_node"`
	TargetHostname string        `json:"target_hostname"`
	TripwireID     string        `json:"tripwire_id,omitempty"`
	DecoyName      string        `json:"decoy_name,omitempty"`
	DecoyType      DecoyType     `json:"decoy_type,omitempty"`
	Actor          string        `json:"actor"`
	RawEventID     string        `json:"raw_event_id,omitempty"`
}
