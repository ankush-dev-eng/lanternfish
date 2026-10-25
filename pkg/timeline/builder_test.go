package timeline

import (
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestBuildTimeline(t *testing.T) {
	decoys := decoy.NewRegistry()

	t0 := time.Now()
	t1 := t0.Add(2 * time.Second)
	t2 := t0.Add(4 * time.Second)

	events := []model.Event{
		{
			ID:         "evt-02",
			Timestamp:  t1,
			Type:       model.EventFileRead,
			Actor:      "adversary_alpha",
			SourceNode: "ws-dev-02",
			TargetNode: "ws-dev-02",
			TripwireID: "TRIP-FILE-PASSWORDS",
			Metadata: map[string]interface{}{
				"phase":           string(model.PhaseCredentialAccess),
				"mitre_technique": "T1083",
				"description":     "Read passwords file",
			},
		},
		{
			ID:         "evt-01",
			Timestamp:  t0,
			Type:       model.EventProcessExec,
			Actor:      "adversary_alpha",
			SourceNode: "gw-edge-01",
			TargetNode: "ws-dev-02",
			Metadata: map[string]interface{}{
				"phase":           string(model.PhaseInitialAccess),
				"mitre_technique": "T1566",
				"description":     "Phishing execution",
			},
		},
		{
			ID:         "evt-03",
			Timestamp:  t2,
			Type:       model.EventExfilBeacon,
			Actor:      "adversary_alpha",
			SourceNode: "ws-dev-02",
			TargetNode: "gw-edge-01",
			Metadata: map[string]interface{}{
				"phase":           string(model.PhaseExfiltration),
				"mitre_technique": "T1041",
				"description":     "C2 beacon",
			},
		},
	}

	alerts := []model.Alert{
		{
			ID:             "alt-02",
			RawEventID:     "evt-02",
			Severity:       model.SeverityHigh,
			Confidence:     96,
			TripwireID:     "TRIP-FILE-PASSWORDS",
			DecoyName:      "passwords_2024_master.xlsx",
			MitreTechnique: "T1083",
		},
	}

	timeline := BuildTimeline(events, alerts, decoys)

	if len(timeline) != 3 {
		t.Fatalf("expected 3 timeline items, got %d", len(timeline))
	}

	// Verify chronological ordering (evt-01 was at t0, evt-02 at t1)
	if timeline[0].ID != "tl-002" || timeline[0].Phase != model.PhaseInitialAccess {
		t.Errorf("expected earliest event to be first, got %+v", timeline[0])
	}

	// Verify decoy detection mapping
	if !timeline[1].IsDecoyTrigger || timeline[1].DecoyName != "passwords_2024_master.xlsx" {
		t.Errorf("expected decoy trigger on item 1, got %+v", timeline[1])
	}

	// Verify grouping
	grouped := GroupByPhase(timeline)
	if len(grouped) != 3 {
		t.Errorf("expected 3 distinct phases in grouped timeline, got %d", len(grouped))
	}
}
