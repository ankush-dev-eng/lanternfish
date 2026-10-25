package model

import "time"

// DecoyType defines the class of deceptive artifact deployed in the network.
type DecoyType string

const (
	DecoyTypeFile        DecoyType = "file"
	DecoyTypeCredential  DecoyType = "credential"
	DecoyTypeHost        DecoyType = "host"
	DecoyTypeServicePort DecoyType = "service_port"
)

// Decoy represents an instrumented decoy digital asset acting as an anglerfish lure.
type Decoy struct {
	ID               string     `json:"id"`
	TripwireID       string     `json:"tripwire_id"`
	Type             DecoyType  `json:"type"`
	Name             string     `json:"name"`
	PathOrIdentifier string     `json:"path_or_identifier"`
	NodeID           string     `json:"node_id"`
	NodeHostname     string     `json:"node_hostname"`
	PayloadSummary   string     `json:"payload_summary"`
	IsDecoy          bool       `json:"is_decoy"` // Always true; operator toggle reveals in UI
	Triggered        bool       `json:"triggered"`
	TriggerCount     int        `json:"trigger_count"`
	FirstTriggeredAt *time.Time `json:"first_triggered_at,omitempty"`
	LastTriggeredAt  *time.Time `json:"last_triggered_at,omitempty"`
	AttackerNotes    string     `json:"attacker_notes,omitempty"`
	BaitEnticement   string     `json:"bait_enticement"` // Why attacker finds it irresistible
}
