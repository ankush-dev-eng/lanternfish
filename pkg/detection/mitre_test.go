package detection

import (
	"strings"
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
)

func TestMitreEnrichment(t *testing.T) {
	graph := network.GenerateEnterpriseTopology()
	decoys := decoy.NewRegistry()
	engine := NewEngine(graph, decoys)

	evt := model.Event{
		ID:         "evt-mitre-test",
		Timestamp:  time.Now(),
		Type:       model.EventFileRead,
		Actor:      "threat_actor_test",
		SourceNode: "ws-dev-02",
		TargetNode: "ws-dev-02",
		TripwireID: "TRIP-FILE-PASSWORDS",
	}

	alerts := engine.ProcessEvent(evt)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}

	alert := alerts[0]
	if alert.MitreTactic != model.TacticDiscovery {
		t.Errorf("expected TacticDiscovery, got %s", alert.MitreTactic)
	}

	if !strings.HasPrefix(alert.MitreURL, "https://attack.mitre.org/techniques/T1083") {
		t.Errorf("unexpected MitreURL: %s", alert.MitreURL)
	}
}
