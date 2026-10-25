package simulator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/eventbus"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/reconstruction"
)

// DetectionEngineHook defines an interface for processing events synchronously as they occur.
type DetectionEngineHook interface {
	ProcessEvent(evt model.Event) []model.Alert
}

// Engine coordinates the asynchronous execution of simulated intrusion scripts.
type Engine struct {
	mu           sync.RWMutex
	bus          *eventbus.Bus
	decoys       *decoy.Registry
	detector     DetectionEngineHook
	graph        *model.NetworkGraph
	state        model.SimulationState
	cancelFn     context.CancelFunc
	pauseChan    chan struct{}
	resumeChan   chan struct{}
	isPaused     bool
	firstAlertAt *time.Time
}

// NewEngine instantiates a simulation engine attached to the telemetry bus and decoy registry.
func NewEngine(bus *eventbus.Bus, decoys *decoy.Registry, detector DetectionEngineHook) *Engine {
	return &Engine{
		bus:        bus,
		decoys:     decoys,
		detector:   detector,
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
		state: model.SimulationState{
			ScenarioID:      model.ScenarioStealthy,
			ScenarioName:    "Stealthy APT & Lure Harvest",
			Status:          model.StatusIdle,
			SpeedMultiplier: 1.0,
			CurrentStep:     0,
			TotalSteps:      len(GetScenarioSteps(model.ScenarioStealthy)),
		},
	}
}

// SetGraph binds the network graph for dynamic path reconstruction.
func (e *Engine) SetGraph(graph *model.NetworkGraph) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.graph = graph
}

// SetDetector allows late binding of the detection engine.
func (e *Engine) SetDetector(detector DetectionEngineHook) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.detector = detector
}

// GetState returns an atomic snapshot of current simulation status.
func (e *Engine) GetState() model.SimulationState {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.state
}

// Start begins execution of the selected intrusion scenario at the specified speed.
func (e *Engine) Start(scenarioType model.ScenarioType, speed float64) error {
	e.mu.Lock()

	// Stop any existing run
	if e.cancelFn != nil {
		e.cancelFn()
	}

	if speed <= 0 {
		speed = 1.0
	}

	steps := GetScenarioSteps(scenarioType)
	scenarioName := "Stealthy APT & Lure Harvest"
	if scenarioType == model.ScenarioNoisy {
		scenarioName = "Aggressive DMZ Spray & Recon"
	} else if scenarioType == model.ScenarioInsider {
		scenarioName = "Malicious Insider Exfiltration"
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.cancelFn = cancel
	e.isPaused = false
	e.firstAlertAt = nil

	now := time.Now()
	e.state = model.SimulationState{
		ScenarioID:       scenarioType,
		ScenarioName:     scenarioName,
		Status:           model.StatusRunning,
		SpeedMultiplier:  speed,
		CurrentStep:      0,
		TotalSteps:       len(steps),
		StartedAt:        &now,
		PatientZeroNode:  steps[0].TargetNode,
		ActiveAttackerAt: steps[0].TargetNode,
		TimeToDetectMs:   0,
		AlertsCount:      0,
		DecoysTriggered:  0,
		HopsCount:        0,
	}

	e.mu.Unlock()

	// Notify initial state
	e.bus.PublishState(e.GetState())

	go e.runSimulation(ctx, steps, speed)
	return nil
}

// Pause pauses active step execution.
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state.Status == model.StatusRunning && !e.isPaused {
		e.isPaused = true
		e.state.Status = model.StatusPaused
		e.bus.PublishState(e.state)
	}
}

// Resume continues a paused simulation.
func (e *Engine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.state.Status == model.StatusPaused && e.isPaused {
		e.isPaused = false
		e.state.Status = model.StatusRunning
		e.bus.PublishState(e.state)
		select {
		case e.resumeChan <- struct{}{}:
		default:
		}
	}
}

// Reset halts execution, clears simulation metrics, and resets decoys and events.
func (e *Engine) Reset() {
	e.mu.Lock()
	if e.cancelFn != nil {
		e.cancelFn()
		e.cancelFn = nil
	}
	e.isPaused = false
	e.firstAlertAt = nil

	e.state = model.SimulationState{
		ScenarioID:      model.ScenarioStealthy,
		ScenarioName:    "Stealthy APT & Lure Harvest",
		Status:          model.StatusIdle,
		SpeedMultiplier: 1.0,
		CurrentStep:     0,
		TotalSteps:      len(GetScenarioSteps(model.ScenarioStealthy)),
	}
	e.mu.Unlock()

	e.decoys.Reset()
	e.bus.Clear()
	e.bus.PublishState(e.GetState())
}

func (e *Engine) runSimulation(ctx context.Context, steps []model.SimulationStep, speed float64) {
	for i, step := range steps {
		// Calculate step delay scaled by speed
		scaledDelay := time.Duration(float64(step.Delay) / speed)

		// Wait or listen for cancel
		select {
		case <-ctx.Done():
			return
		case <-time.After(scaledDelay):
		}

		// Handle pause check
		e.mu.RLock()
		paused := e.isPaused
		e.mu.RUnlock()

		if paused {
			select {
			case <-ctx.Done():
				return
			case <-e.resumeChan:
			}
		}

		// Trigger Decoy lure if step interacts with a tripwire
		if step.DecoyTripwire != "" {
			e.decoys.TriggerTripwire(step.DecoyTripwire, step.Actor, step.Description)
		}

		// Construct telemetry event
		evt := model.Event{
			ID:         fmt.Sprintf("evt-%03d-%s", i+1, step.EventType),
			Timestamp:  time.Now(),
			Type:       step.EventType,
			Actor:      step.Actor,
			SourceNode: step.SourceNode,
			TargetNode: step.TargetNode,
			TripwireID: step.DecoyTripwire,
			Success:    true,
			RawPayload: step.RawPayload,
			Metadata: map[string]interface{}{
				"phase":           step.Phase,
				"mitre_technique": step.MitreTechnique,
				"description":     step.Description,
				"step_number":     step.StepNumber,
			},
		}

		// Publish event to telemetry bus
		e.bus.PublishEvent(evt)

		// Run detection engine hook if available
		var newAlerts []model.Alert
		e.mu.RLock()
		det := e.detector
		graph := e.graph
		e.mu.RUnlock()

		if det != nil {
			newAlerts = det.ProcessEvent(evt)
			for _, alert := range newAlerts {
				e.bus.PublishAlert(alert)
			}
		}

		// Dynamically reconstruct attacker trajectory if network topology is bound
		var currentPath *model.AttackerPath
		if graph != nil {
			var allAlerts []model.Alert
			if ap, ok := det.(interface{ GetAlerts() []model.Alert }); ok {
				allAlerts = ap.GetAlerts()
			}
			p := reconstruction.ReconstructPath(graph, e.bus.GetEvents(), allAlerts)
			currentPath = &p
			e.bus.PublishPath(p)
		}

		// Update state counters
		e.mu.Lock()
		e.state.CurrentStep = i + 1
		e.state.ActiveAttackerAt = step.TargetNode
		if currentPath != nil {
			e.state.HopsCount = len(currentPath.Hops)
			if currentPath.PatientZeroID != "" {
				e.state.PatientZeroNode = currentPath.PatientZeroID
			}
		} else if step.EventType == model.EventLateralHop {
			e.state.HopsCount++
		}
		if len(newAlerts) > 0 {
			e.state.AlertsCount += len(newAlerts)
			if e.firstAlertAt == nil && e.state.StartedAt != nil {
				now := time.Now()
				e.firstAlertAt = &now
				e.state.TimeToDetectMs = now.Sub(*e.state.StartedAt).Milliseconds()
			}
		}
		e.state.DecoysTriggered = e.decoys.TriggeredCount()
		e.mu.Unlock()

		// Publish state update to stream
		e.bus.PublishState(e.GetState())
	}

	// Simulation finished
	e.mu.Lock()
	now := time.Now()
	e.state.Status = model.StatusCompleted
	e.state.CompletedAt = &now
	e.mu.Unlock()

	e.bus.PublishState(e.GetState())
}
