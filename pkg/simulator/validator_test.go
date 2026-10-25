package simulator

import (
	"testing"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/network"
)

func TestValidateScenarioAgainstTopology(t *testing.T) {
	graph := network.GenerateEnterpriseTopology()

	scenarios := []model.ScenarioType{
		model.ScenarioStealthy,
		model.ScenarioNoisy,
		model.ScenarioInsider,
	}

	for _, sc := range scenarios {
		if err := ValidateScenarioAgainstTopology(sc, graph); err != nil {
			t.Fatalf("scenario %s failed validation: %v", sc, err)
		}
	}

	// Verify error on unknown node reference
	emptyGraph := &model.NetworkGraph{
		Nodes: []model.Node{},
	}
	if err := ValidateScenarioAgainstTopology(model.ScenarioStealthy, emptyGraph); err == nil {
		t.Fatal("expected error validating scenario against empty graph, got nil")
	}
}
