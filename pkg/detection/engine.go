package detection

import (
	"sync"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// Engine manages the evaluation of security rules across telemetry events.
type Engine struct {
	mu     sync.RWMutex
	graph  *model.NetworkGraph
	decoys *decoy.Registry
	rules  []Rule
	alerts []model.Alert
}

// NewEngine creates an initialized Detection Engine armed with standard deception and anomaly rules.
func NewEngine(graph *model.NetworkGraph, decoys *decoy.Registry) *Engine {
	rules := []Rule{
		&DecoyFileAccessRule{},
		&HoneytokenCredentialRule{},
		&DecoyHostConnectionRule{},
		&AbnormalLateralMovementRule{},
		&OffHoursAccessRule{},
		&ExfiltrationBeaconRule{},
	}

	return &Engine{
		graph:  graph,
		decoys: decoys,
		rules:  rules,
		alerts: make([]model.Alert, 0, 50),
	}
}

// ProcessEvent applies all active detection rules to an incoming event, returning any triggered alerts.
func (e *Engine) ProcessEvent(evt model.Event) []model.Alert {
	e.mu.Lock()
	defer e.mu.Unlock()

	var triggered []model.Alert

	for _, rule := range e.rules {
		alert := rule.Evaluate(evt, e.graph, e.decoys)
		if alert != nil {
			// Enrich with hostnames
			if srcNode := e.graph.FindNode(alert.SourceNode); srcNode != nil {
				alert.SourceHostname = srcNode.Hostname
			}
			if dstNode := e.graph.FindNode(alert.TargetNode); dstNode != nil {
				alert.TargetHostname = dstNode.Hostname
			}

			// Enrich with MITRE ATT&CK definition
			mitreDef := model.GetMitreDefinition(alert.MitreTechnique)
			alert.MitreTactic = mitreDef.Tactic
			alert.MitreURL = mitreDef.URL

			triggered = append(triggered, *alert)
			e.alerts = append(e.alerts, *alert)
		}
	}

	return triggered
}

// GetAlerts returns a copy of all alerts recorded by the detection engine.
func (e *Engine) GetAlerts() []model.Alert {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]model.Alert, len(e.alerts))
	copy(result, e.alerts)
	return result
}

// Reset flushes all recorded alerts.
func (e *Engine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.alerts = make([]model.Alert, 0, 50)
}
