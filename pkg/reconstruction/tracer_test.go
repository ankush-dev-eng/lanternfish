package reconstruction

import (
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
)

func TestReconstructPath(t *testing.T) {
	graph := network.GenerateEnterpriseTopology()

	t0 := time.Now()
	t1 := t0.Add(5 * time.Second)
	t2 := t0.Add(10 * time.Second)

	events := []model.Event{
		{
			ID:         "evt-01",
			Timestamp:  t0,
			Type:       model.EventProcessExec,
			Actor:      "adversary_x",
			SourceNode: "gw-edge-01",
			TargetNode: "ws-dev-02",
			Metadata: map[string]interface{}{
				"phase":           string(model.PhaseInitialAccess),
				"mitre_technique": "T1566",
				"description":     "Spearphishing foothold",
			},
		},
		{
			ID:         "evt-02",
			Timestamp:  t1,
			Type:       model.EventLateralHop,
			Actor:      "adversary_x",
			SourceNode: "ws-dev-02",
			TargetNode: "api-gw-01",
			Metadata: map[string]interface{}{
				"phase":           string(model.PhaseLateralMovement),
				"mitre_technique": "T1021",
				"description":     "SSH lateral jump",
			},
		},
		{
			ID:         "evt-03",
			Timestamp:  t2,
			Type:       model.EventCredentialUse,
			Actor:      "adversary_x",
			SourceNode: "api-gw-01",
			TargetNode: "db-legacy-backup",
			TripwireID: "TRIP-HOST-DBLEGACY",
			Metadata: map[string]interface{}{
				"phase":           string(model.PhaseDecoyInteraction),
				"mitre_technique": "T1021",
				"description":     "Probe decoy host",
			},
		},
	}

	alerts := []model.Alert{
		{
			ID:         "alt-03",
			RawEventID: "evt-03",
			Severity:   model.SeverityHigh,
			Confidence: 98,
		},
	}

	path := ReconstructPath(graph, events, alerts)

	if path.PatientZeroID != "ws-dev-02" {
		t.Fatalf("expected PatientZero ws-dev-02, got %s", path.PatientZeroID)
	}

	if len(path.Hops) != 3 {
		t.Fatalf("expected 3 hops, got %d", len(path.Hops))
	}

	if len(path.NodesVisited) != 4 { // ws-dev-02, gw-edge-01, api-gw-01, db-legacy-backup
		t.Fatalf("expected 4 nodes visited, got %d", len(path.NodesVisited))
	}

	if len(path.DecoysTouched) != 1 || path.DecoysTouched[0] != "TRIP-HOST-DBLEGACY" {
		t.Fatalf("expected TRIP-HOST-DBLEGACY decoy touched, got %+v", path.DecoysTouched)
	}

	if path.DwellTimeSeconds < 9.0 {
		t.Fatalf("expected dwell time >= 9s, got %f", path.DwellTimeSeconds)
	}

	if len(path.KillChainPhases) != 3 {
		t.Fatalf("expected 3 kill chain phases, got %d", len(path.KillChainPhases))
	}
}
