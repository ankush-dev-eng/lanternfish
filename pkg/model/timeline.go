package model

import "time"

// KillChainPhase groups events into MITRE ATT&CK / cyber kill-chain progression stages.
type KillChainPhase string

const (
	PhaseInitialAccess    KillChainPhase = "Initial Access"
	PhaseReconnaissance   KillChainPhase = "Discovery & Recon"
	PhaseCredentialAccess KillChainPhase = "Credential Hunting"
	PhaseLateralMovement  KillChainPhase = "Lateral Movement"
	PhaseDecoyInteraction KillChainPhase = "Tripwire Engagement"
	PhaseCollection       KillChainPhase = "Collection"
	PhaseExfiltration     KillChainPhase = "Exfiltration Attempt"
)

// TimelineItem represents a human-readable, sorted, categorized incident narrative checkpoint.
type TimelineItem struct {
	ID             string         `json:"id"`
	Timestamp      time.Time      `json:"timestamp"`
	Phase          KillChainPhase `json:"phase"`
	Severity       AlertSeverity  `json:"severity"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	SourceNode     string         `json:"source_node"`
	TargetNode     string         `json:"target_node"`
	MitreTechnique string         `json:"mitre_technique"`
	IsDecoyTrigger bool           `json:"is_decoy_trigger"`
	DecoyName      string         `json:"decoy_name,omitempty"`
	Confidence     int            `json:"confidence"`
}
