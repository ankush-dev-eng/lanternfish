package eventbus

import (
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestEventBus(t *testing.T) {
	bus := NewBus(10)

	subCh, unsubscribe := bus.Subscribe()
	defer unsubscribe()

	// Test PublishEvent
	testEvt := model.Event{
		ID:        "evt-001",
		Timestamp: time.Now(),
		Type:      model.EventFileRead,
		Actor:     "test_adversary",
	}
	bus.PublishEvent(testEvt)

	select {
	case msg := <-subCh:
		if msg.Event != "event" {
			t.Fatalf("expected event type, got %s", msg.Event)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for event on subscriber channel")
	}

	events := bus.GetEvents()
	if len(events) != 1 || events[0].ID != "evt-001" {
		t.Fatalf("unexpected events slice: %+v", events)
	}

	// Test PublishAlert
	testAlt := model.Alert{
		ID:       "alt-001",
		RuleName: "TestRule",
		Severity: model.SeverityHigh,
	}
	bus.PublishAlert(testAlt)

	select {
	case msg := <-subCh:
		if msg.Event != "alert" {
			t.Fatalf("expected alert type, got %s", msg.Event)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for alert on subscriber channel")
	}

	alerts := bus.GetAlerts()
	if len(alerts) != 1 || alerts[0].ID != "alt-001" {
		t.Fatalf("unexpected alerts slice: %+v", alerts)
	}

	// Test Clear
	bus.Clear()
	if len(bus.GetEvents()) != 0 || len(bus.GetAlerts()) != 0 {
		t.Fatal("expected empty slices after clear")
	}
}
