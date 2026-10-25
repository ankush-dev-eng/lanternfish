package reconstruction

import (
	"fmt"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// ReconstructPath analyzes events and alerts to map patient-zero and subsequent lateral hops.
func ReconstructPath(graph *model.NetworkGraph, events []model.Event, alerts []model.Alert) model.AttackerPath {
	if len(events) == 0 {
		return model.AttackerPath{
			Hops:            []model.ReconstructedHop{},
			NodesVisited:    []string{},
			DecoysTouched:   []string{},
			KillChainPhases: []string{},
		}
	}

	// 1. Identify Patient Zero (first target or source node of the earliest event)
	earliestEvt := events[0]
	patientZeroID := earliestEvt.TargetNode
	if patientZeroID == "" {
		patientZeroID = earliestEvt.SourceNode
	}
	patientZeroHostname := patientZeroID
	if n := graph.FindNode(patientZeroID); n != nil {
		patientZeroHostname = n.Hostname
	}

	// Index alerts by RawEventID or TripwireID
	alertByEventID := make(map[string]model.Alert)
	for _, a := range alerts {
		if a.RawEventID != "" {
			alertByEventID[a.RawEventID] = a
		}
	}

	hops := make([]model.ReconstructedHop, 0, 10)
	nodesVisitedMap := make(map[string]struct{})
	nodesVisitedList := make([]string, 0, 10)
	decoysTouchedMap := make(map[string]struct{})
	decoysTouchedList := make([]string, 0, 10)
	phasesMap := make(map[string]struct{})
	phasesList := make([]string, 0, 8)

	// Record Patient Zero
	nodesVisitedMap[patientZeroID] = struct{}{}
	nodesVisitedList = append(nodesVisitedList, patientZeroID)

	currentLoc := patientZeroID
	startTime := earliestEvt.Timestamp
	latestTime := earliestEvt.Timestamp

	for _, evt := range events {
		if evt.Timestamp.After(latestTime) {
			latestTime = evt.Timestamp
		}

		// Extract Kill Chain Phase from metadata
		if phaseVal, ok := evt.Metadata["phase"].(string); ok && phaseVal != "" {
			if _, seen := phasesMap[phaseVal]; !seen {
				phasesMap[phaseVal] = struct{}{}
				phasesList = append(phasesList, phaseVal)
			}
		}

		// Record decoy interaction
		if evt.TripwireID != "" {
			if _, seen := decoysTouchedMap[evt.TripwireID]; !seen {
				decoysTouchedMap[evt.TripwireID] = struct{}{}
				decoysTouchedList = append(decoysTouchedList, evt.TripwireID)
			}
		}

		// Check for lateral hop or significant node traversal
		isHop := evt.Type == model.EventLateralHop ||
			(evt.SourceNode != "" && evt.TargetNode != "" && evt.SourceNode != evt.TargetNode)

		if isHop {
			// Source and target hostnames
			fromHost := evt.SourceNode
			if n := graph.FindNode(evt.SourceNode); n != nil {
				fromHost = n.Hostname
			}
			toHost := evt.TargetNode
			if n := graph.FindNode(evt.TargetNode); n != nil {
				toHost = n.Hostname
			}

			// Matched alert
			alertID := ""
			conf := 85
			if matchedAlert, hasAlert := alertByEventID[evt.ID]; hasAlert {
				alertID = matchedAlert.ID
				conf = matchedAlert.Confidence
			}

			mitre := "T1021"
			if m, ok := evt.Metadata["mitre_technique"].(string); ok && m != "" {
				mitre = m
			}

			desc := fmt.Sprintf("Traversal from %s to %s", fromHost, toHost)
			if d, ok := evt.Metadata["description"].(string); ok && d != "" {
				desc = d
			}

			hop := model.ReconstructedHop{
				HopIndex:       len(hops) + 1,
				Timestamp:      evt.Timestamp,
				FromNode:       evt.SourceNode,
				FromHostname:   fromHost,
				ToNode:         evt.TargetNode,
				ToHostname:     toHost,
				MitreTechnique: mitre,
				Method:         string(evt.Type),
				Evidence:       desc,
				AlertID:        alertID,
				Confidence:     conf,
				DecoyTripped:   evt.TripwireID,
			}
			hops = append(hops, hop)

			// Record node visits
			if _, seen := nodesVisitedMap[evt.TargetNode]; !seen {
				nodesVisitedMap[evt.TargetNode] = struct{}{}
				nodesVisitedList = append(nodesVisitedList, evt.TargetNode)
			}
			currentLoc = evt.TargetNode
		}
	}

	dwellTime := latestTime.Sub(startTime).Seconds()
	if dwellTime < 0 {
		dwellTime = 0
	}

	return model.AttackerPath{
		PatientZeroID:       patientZeroID,
		PatientZeroHostname: patientZeroHostname,
		CompromiseTime:      startTime,
		CurrentLocation:     currentLoc,
		Hops:                hops,
		NodesVisited:        nodesVisitedList,
		DecoysTouched:       decoysTouchedList,
		DwellTimeSeconds:    dwellTime,
		KillChainPhases:     phasesList,
	}
}
