package detection

import (
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
)

func TestDetectionEngine(t *testing.T) {
	graph := network.GenerateEnterpriseTopology()
	decoys := decoy.NewRegistry()
	engine := NewEngine(graph, decoys)

	// 1. Test Decoy File Access Rule
	fileEvt := model.Event{
		ID:         "evt-f1",
		Timestamp:  time.Now(),
		Type:       model.EventFileRead,
		Actor:      "adversary_alpha",
		SourceNode: "ws-dev-02",
		TargetNode: "ws-dev-02",
		TripwireID: "TRIP-FILE-PASSWORDS",
	}
	alerts := engine.ProcessEvent(fileEvt)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for decoy file access, got %d", len(alerts))
	}
	if alerts[0].Severity != model.SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", alerts[0].Severity)
	}
	if alerts[0].MitreTechnique != "T1083" {
		t.Errorf("expected MITRE T1083, got %s", alerts[0].MitreTechnique)
	}

	// 2. Test Honeytoken Credential Rule
	credEvt := model.Event{
		ID:         "evt-c1",
		Timestamp:  time.Now(),
		Type:       model.EventCredentialUse,
		Actor:      "adversary_alpha",
		SourceNode: "app-srv-01",
		TargetNode: "db-postgres-01",
		TripwireID: "TRIP-CRED-DBROOT",
	}
	alerts = engine.ProcessEvent(credEvt)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for honeytoken use, got %d", len(alerts))
	}
	if alerts[0].Severity != model.SeverityCritical {
		t.Errorf("expected CRITICAL severity, got %s", alerts[0].Severity)
	}
	if alerts[0].MitreTechnique != "T1078" {
		t.Errorf("expected MITRE T1078, got %s", alerts[0].MitreTechnique)
	}

	// 3. Test Decoy Host Rule
	hostEvt := model.Event{
		ID:         "evt-h1",
		Timestamp:  time.Now(),
		Type:       model.EventNetworkConnection,
		Actor:      "adversary_alpha",
		SourceNode: "app-srv-01",
		TargetNode: "db-legacy-backup",
		TripwireID: "TRIP-HOST-DBLEGACY",
	}
	alerts = engine.ProcessEvent(hostEvt)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for decoy host interaction, got %d", len(alerts))
	}
	if alerts[0].MitreTechnique != "T1021" {
		t.Errorf("expected MITRE T1021, got %s", alerts[0].MitreTechnique)
	}

	// 4. Test Lateral Movement Rule
	lateralEvt := model.Event{
		ID:         "evt-l1",
		Timestamp:  time.Now(),
		Type:       model.EventLateralHop,
		Actor:      "adversary_alpha",
		SourceNode: "ws-dev-02",
		TargetNode: "api-gw-01",
	}
	alerts = engine.ProcessEvent(lateralEvt)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for lateral hop, got %d", len(alerts))
	}

	// 5. Test Exfiltration Beacon Rule
	exfilEvt := model.Event{
		ID:         "evt-e1",
		Timestamp:  time.Now(),
		Type:       model.EventExfilBeacon,
		Actor:      "adversary_alpha",
		SourceNode: "app-srv-01",
		TargetNode: "gw-edge-01",
	}
	alerts = engine.ProcessEvent(exfilEvt)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for exfil beacon, got %d", len(alerts))
	}
	if alerts[0].Severity != model.SeverityCritical {
		t.Errorf("expected CRITICAL severity, got %s", alerts[0].Severity)
	}

	// Verify all alerts stored
	allAlerts := engine.GetAlerts()
	if len(allAlerts) != 5 {
		t.Fatalf("expected 5 total alerts stored, got %d", len(allAlerts))
	}

	// Verify reset
	engine.Reset()
	if len(engine.GetAlerts()) != 0 {
		t.Fatal("expected 0 alerts after reset")
	}
}
