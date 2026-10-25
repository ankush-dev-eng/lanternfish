package timeline

import (
	"fmt"
	"sort"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// BuildTimeline constructs a sorted, narrative-driven chronological sequence from events and alerts.
func BuildTimeline(events []model.Event, alerts []model.Alert, decoys *decoy.Registry) []model.TimelineItem {
	// Index alerts by RawEventID
	alertByEventID := make(map[string]model.Alert)
	for _, a := range alerts {
		if a.RawEventID != "" {
			alertByEventID[a.RawEventID] = a
		}
	}

	items := make([]model.TimelineItem, 0, len(events))

	for i, evt := range events {
		matchedAlert, hasAlert := alertByEventID[evt.ID]

		phase := model.PhaseReconnaissance
		if phaseStr, ok := evt.Metadata["phase"].(string); ok && phaseStr != "" {
			phase = model.KillChainPhase(phaseStr)
		} else if hasAlert && matchedAlert.TripwireID != "" {
			phase = model.PhaseDecoyInteraction
		}

		severity := model.SeverityLow
		if hasAlert {
			severity = matchedAlert.Severity
		} else if evt.Type == model.EventLateralHop {
			severity = model.SeverityMedium
		} else if evt.Type == model.EventExfilBeacon {
			severity = model.SeverityCritical
		}

		title := string(evt.Type)
		desc := ""
		if d, ok := evt.Metadata["description"].(string); ok && d != "" {
			desc = d
			title = desc
		} else if hasAlert {
			title = matchedAlert.Title
			desc = matchedAlert.Description
		}

		mitre := ""
		if m, ok := evt.Metadata["mitre_technique"].(string); ok && m != "" {
			mitre = m
		} else if hasAlert {
			mitre = matchedAlert.MitreTechnique
		}

		isDecoy := hasAlert || evt.TripwireID != ""
		decoyName := ""
		if isDecoy {
			if matchedAlert.DecoyName != "" {
				decoyName = matchedAlert.DecoyName
			} else if evt.TripwireID != "" && decoys != nil {
				if d := decoys.FindByTripwireID(evt.TripwireID); d != nil {
					decoyName = d.Name
				}
			}
		}

		conf := 75
		if hasAlert {
			conf = matchedAlert.Confidence
		} else if isDecoy {
			conf = 95
		}

		item := model.TimelineItem{
			ID:             fmt.Sprintf("tl-%03d", i+1),
			Timestamp:      evt.Timestamp,
			Phase:          phase,
			Severity:       severity,
			Title:          title,
			Description:    desc,
			SourceNode:     evt.SourceNode,
			TargetNode:     evt.TargetNode,
			MitreTechnique: mitre,
			IsDecoyTrigger: isDecoy,
			DecoyName:      decoyName,
			Confidence:     conf,
		}
		items = append(items, item)
	}

	// Sort chronologically
	sort.Slice(items, func(i, j int) bool {
		return items[i].Timestamp.Before(items[j].Timestamp)
	})

	return items
}

// GroupByPhase organizes timeline entries by their respective MITRE kill-chain phase.
func GroupByPhase(items []model.TimelineItem) map[model.KillChainPhase][]model.TimelineItem {
	grouped := make(map[model.KillChainPhase][]model.TimelineItem)
	for _, item := range items {
		grouped[item.Phase] = append(grouped[item.Phase], item)
	}
	return grouped
}
