package narrative

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestGenerateNarrativeDeterministic(t *testing.T) {
	// Ensure GEMINI_API_KEY is not interfering with unit test
	os.Unsetenv("GEMINI_API_KEY")

	path := model.AttackerPath{
		PatientZeroID:       "ws-dev-02",
		PatientZeroHostname: "ws-dev-02.corp",
		CompromiseTime:      time.Now(),
		Hops: []model.ReconstructedHop{
			{HopIndex: 1, FromNode: "ws-dev-02", ToNode: "api-gw-01"},
		},
		NodesVisited:     []string{"ws-dev-02", "api-gw-01"},
		DecoysTouched:    []string{"TRIP-FILE-PASSWORDS"},
		DwellTimeSeconds: 15.4,
	}

	scenarios := []model.ScenarioType{
		model.ScenarioStealthy,
		model.ScenarioNoisy,
		model.ScenarioInsider,
	}

	for _, sc := range scenarios {
		out := GenerateNarrative(sc, path, nil, nil)

		if out.IsAI {
			t.Errorf("expected deterministic engine, got IsAI=true for %s", sc)
		}

		if !strings.Contains(out.ModelUsed, "Deterministic") {
			t.Errorf("expected Deterministic model name, got %s", out.ModelUsed)
		}

		if len(out.ExecutiveSummary) < 50 {
			t.Errorf("executive summary too short for %s", sc)
		}

		if len(out.RootCauseAnalysis) < 30 {
			t.Errorf("root cause too short for %s", sc)
		}

		if len(out.Recommendations) == 0 {
			t.Errorf("expected containment recommendations for %s", sc)
		}

		for _, rec := range out.Recommendations {
			if rec.Priority == "" || rec.Action == "" || rec.Target == "" {
				t.Errorf("incomplete recommendation for %s: %+v", sc, rec)
			}
		}
	}
}
