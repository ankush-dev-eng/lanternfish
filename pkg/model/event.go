package model

import "time"

// EventType enumerates simulated operating system and network telemetry.
type EventType string

const (
	EventPortScan          EventType = "port_scan"
	EventAuthAttempt       EventType = "auth_attempt"
	EventFileRead          EventType = "file_read"
	EventFileWrite         EventType = "file_write"
	EventProcessExec       EventType = "process_exec"
	EventLateralHop        EventType = "lateral_hop"
	EventDecoyInteraction  EventType = "decoy_interaction"
	EventCredentialUse     EventType = "credential_use"
	EventExfilBeacon       EventType = "exfil_beacon"
	EventNetworkConnection EventType = "network_connection"
)

// Event represents a discrete simulated security telemetry item.
type Event struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Type       EventType              `json:"type"`
	Actor      string                 `json:"actor"`
	SourceNode string                 `json:"source_node"`
	TargetNode string                 `json:"target_node"`
	SourceIP   string                 `json:"source_ip"`
	TargetIP   string                 `json:"target_ip"`
	Process    string                 `json:"process,omitempty"`
	Resource   string                 `json:"resource,omitempty"`
	DecoyID    string                 `json:"decoy_id,omitempty"`
	TripwireID string                 `json:"tripwire_id,omitempty"`
	Success    bool                   `json:"success"`
	RawPayload string                 `json:"raw_payload,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
