package decoy

import (
	"testing"
)

func TestCategoryStats(t *testing.T) {
	reg := NewRegistry()
	stats := reg.GetCategoryStats()

	if stats.TotalFiles == 0 || stats.TotalCredentials == 0 || stats.TotalHosts == 0 {
		t.Fatalf("expected non-zero totals for all categories, got %+v", stats)
	}

	if stats.TrippedFiles != 0 || stats.TrippedCredentials != 0 || stats.TrippedHosts != 0 {
		t.Fatalf("expected 0 tripped initially, got %+v", stats)
	}

	// Trigger a file decoy
	reg.TriggerTripwire("TRIP-FILE-PASSWORDS", "actor_1", "Read file")
	updatedStats := reg.GetCategoryStats()
	if updatedStats.TrippedFiles != 1 {
		t.Fatalf("expected 1 tripped file, got %d", updatedStats.TrippedFiles)
	}
}
