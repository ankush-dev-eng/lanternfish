package simulator

import (
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/eventbus"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

type mockDetector struct {
	processed []model.Event
}

func (m *mockDetector) ProcessEvent(evt model.Event) []model.Alert {
	m.processed = append(m.processed, evt)
	if evt.TripwireID != "" {
		return []model.Alert{
			{
				ID:         "mock-alt-" + evt.ID,
				RuleName:   "MockDecoyRule",
				Severity:   model.SeverityHigh,
				Confidence: 95,
				TripwireID: evt.TripwireID,
			},
		}
	}
	return nil
}

func TestSimulatorEngine(t *testing.T) {
	bus := eventbus.NewBus(100)
	decoys := decoy.NewRegistry()
	detector := &mockDetector{}

	engine := NewEngine(bus, decoys, detector)

	initialState := engine.GetState()
	if initialState.Status != model.StatusIdle {
		t.Fatalf("expected idle status, got %s", initialState.Status)
	}

	// Run stealthy scenario at 100x speed for instant test execution
	err := engine.Start(model.ScenarioStealthy, 100.0)
	if err != nil {
		t.Fatalf("failed to start simulation: %v", err)
	}

	// Wait for completion (9 steps at ~15-20ms each = ~200ms)
	time.Sleep(400 * time.Millisecond)

	finalState := engine.GetState()
	if finalState.Status != model.StatusCompleted {
		t.Fatalf("expected completed status, got %s", finalState.Status)
	}

	if finalState.CurrentStep != finalState.TotalSteps {
		t.Fatalf("expected all steps to finish, current=%d, total=%d", finalState.CurrentStep, finalState.TotalSteps)
	}

	if finalState.AlertsCount == 0 {
		t.Fatal("expected alerts to be triggered during scenario")
	}

	if finalState.DecoysTriggered == 0 {
		t.Fatal("expected decoys to be triggered during scenario")
	}

	if len(detector.processed) == 0 {
		t.Fatal("expected detector to process events")
	}

	// Test Reset
	engine.Reset()
	resetState := engine.GetState()
	if resetState.Status != model.StatusIdle {
		t.Fatalf("expected status idle after reset, got %s", resetState.Status)
	}
	if resetState.CurrentStep != 0 {
		t.Fatalf("expected step 0 after reset, got %d", resetState.CurrentStep)
	}
}
