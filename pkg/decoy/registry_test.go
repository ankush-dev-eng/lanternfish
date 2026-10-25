package decoy

import (
	"testing"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestDecoyRegistry(t *testing.T) {
	reg := NewRegistry()
	decoys := reg.GetAll()

	if len(decoys) < 8 {
		t.Fatalf("expected at least 8 seeded decoys, got %d", len(decoys))
	}

	// Verify all decoys have tripwires and are flagged as decoys
	fileCount := 0
	credCount := 0
	hostCount := 0
	for _, d := range decoys {
		if !d.IsDecoy {
			t.Errorf("decoy %s has IsDecoy false", d.ID)
		}
		if d.TripwireID == "" {
			t.Errorf("decoy %s is missing TripwireID", d.ID)
		}
		switch d.Type {
		case model.DecoyTypeFile:
			fileCount++
		case model.DecoyTypeCredential:
			credCount++
		case model.DecoyTypeHost:
			hostCount++
		}
	}

	if fileCount == 0 || credCount == 0 || hostCount == 0 {
		t.Fatalf("expected diverse decoy types: files=%d, creds=%d, hosts=%d", fileCount, credCount, hostCount)
	}

	// Test Lookup by Tripwire ID
	found := reg.FindByTripwireID("TRIP-FILE-PASSWORDS")
	if found == nil {
		t.Fatal("expected to find TRIP-FILE-PASSWORDS")
	}
	if found.Name != "passwords_2024_master.xlsx" {
		t.Errorf("unexpected name: %s", found.Name)
	}

	// Test TriggerTripwire
	tripped, ok := reg.TriggerTripwire("TRIP-FILE-PASSWORDS", "adversary_actor", "Read via PowerShell cat")
	if !ok || tripped == nil {
		t.Fatal("failed to trigger tripwire")
	}
	if !tripped.Triggered || tripped.TriggerCount != 1 {
		t.Errorf("expected Triggered=true, Count=1, got %+v", tripped)
	}
	if reg.TriggeredCount() != 1 {
		t.Errorf("expected TriggeredCount=1, got %d", reg.TriggeredCount())
	}

	// Test Reset
	reg.Reset()
	if reg.TriggeredCount() != 0 {
		t.Fatalf("expected 0 triggered after reset, got %d", reg.TriggeredCount())
	}
}
