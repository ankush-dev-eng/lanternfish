package simulator

import (
	"fmt"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// ValidateScenarioAgainstTopology checks that every step references known nodes in the graph.
func ValidateScenarioAgainstTopology(scenarioType model.ScenarioType, graph *model.NetworkGraph) error {
	steps := GetScenarioSteps(scenarioType)
	if len(steps) == 0 {
		return fmt.Errorf("scenario %s has no steps defined", scenarioType)
	}

	for _, step := range steps {
		if step.SourceNode != "" && graph.FindNode(step.SourceNode) == nil {
			return fmt.Errorf("step %d references unknown source node: %s", step.StepNumber, step.SourceNode)
		}
		if step.TargetNode != "" && graph.FindNode(step.TargetNode) == nil {
			return fmt.Errorf("step %d references unknown target node: %s", step.StepNumber, step.TargetNode)
		}
	}

	return nil
}
