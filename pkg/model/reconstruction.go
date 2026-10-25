package model

import "time"

// ReconstructedHop represents an inferred or observed traversal between network nodes.
type ReconstructedHop struct {
	HopIndex       int       `json:"hop_index"`
	Timestamp      time.Time `json:"timestamp"`
	FromNode       string    `json:"from_node"`
	FromHostname   string    `json:"from_hostname"`
	ToNode         string    `json:"to_node"`
	ToHostname     string    `json:"to_hostname"`
	MitreTechnique string    `json:"mitre_technique"`
	Method         string    `json:"method"`
	Evidence       string    `json:"evidence"`
	AlertID        string    `json:"alert_id,omitempty"`
	Confidence     int       `json:"confidence"`
	DecoyTripped   string    `json:"decoy_tripped,omitempty"`
}

// AttackerPath encapsulates the end-to-end reconstructed intrusion trajectory.
type AttackerPath struct {
	PatientZeroID       string             `json:"patient_zero_id"`
	PatientZeroHostname string             `json:"patient_zero_hostname"`
	CompromiseTime      time.Time          `json:"compromise_time"`
	CurrentLocation     string             `json:"current_location"`
	Hops                []ReconstructedHop `json:"hops"`
	NodesVisited        []string           `json:"nodes_visited"`
	DecoysTouched       []string           `json:"decoys_touched"`
	DwellTimeSeconds    float64            `json:"dwell_time_seconds"`
	KillChainPhases     []string           `json:"kill_chain_phases"`
}
