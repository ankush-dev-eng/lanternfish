package eventbus

import (
	"encoding/json"
	"sync"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// SSEMessage represents a categorized message broadcast across server-sent event streams.
type SSEMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// Bus provides thread-safe in-memory publish-subscribe telemetry routing.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[chan SSEMessage]struct{}
	events      []model.Event
	alerts      []model.Alert
	maxHistory  int
}

// NewBus creates an initialized EventBus with the specified event retention cap.
func NewBus(maxHistory int) *Bus {
	if maxHistory <= 0 {
		maxHistory = 2000
	}
	return &Bus{
		subscribers: make(map[chan SSEMessage]struct{}),
		events:      make([]model.Event, 0, 100),
		alerts:      make([]model.Alert, 0, 50),
		maxHistory:  maxHistory,
	}
}

// Subscribe returns an SSE message channel and a cleanup unregister function.
func (b *Bus) Subscribe() (<-chan SSEMessage, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan SSEMessage, 100)
	b.subscribers[ch] = struct{}{}

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if _, exists := b.subscribers[ch]; exists {
			delete(b.subscribers, ch)
			close(ch)
		}
	}

	return ch, unsubscribe
}

// broadcast dispatches an SSE message to all connected clients.
func (b *Bus) broadcast(msg SSEMessage) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.subscribers {
		select {
		case ch <- msg:
		default:
			// Non-blocking drop if consumer is congested
		}
	}
}

// PublishEvent records a telemetry event and pushes it to subscribers.
func (b *Bus) PublishEvent(evt model.Event) {
	b.mu.Lock()
	b.events = append(b.events, evt)
	if len(b.events) > b.maxHistory {
		b.events = b.events[1:]
	}
	b.mu.Unlock()

	data, err := json.Marshal(evt)
	if err == nil {
		b.broadcast(SSEMessage{
			Event: "event",
			Data:  string(data),
		})
	}
}

// PublishAlert records a detection alert and pushes it to subscribers.
func (b *Bus) PublishAlert(alt model.Alert) {
	b.mu.Lock()
	b.alerts = append(b.alerts, alt)
	b.mu.Unlock()

	data, err := json.Marshal(alt)
	if err == nil {
		b.broadcast(SSEMessage{
			Event: "alert",
			Data:  string(data),
		})
	}
}

// PublishState broadcasts the current simulation state without storing it in event log.
func (b *Bus) PublishState(state model.SimulationState) {
	data, err := json.Marshal(state)
	if err == nil {
		b.broadcast(SSEMessage{
			Event: "state",
			Data:  string(data),
		})
	}
}

// PublishPath broadcasts an updated attacker path reconstruction.
func (b *Bus) PublishPath(path model.AttackerPath) {
	data, err := json.Marshal(path)
	if err == nil {
		b.broadcast(SSEMessage{
			Event: "path",
			Data:  string(data),
		})
	}
}

// GetEvents returns a copy of all recorded telemetry events.
func (b *Bus) GetEvents() []model.Event {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]model.Event, len(b.events))
	copy(result, b.events)
	return result
}

// GetAlerts returns a copy of all triggered alerts.
func (b *Bus) GetAlerts() []model.Alert {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]model.Alert, len(b.alerts))
	copy(result, b.alerts)
	return result
}

// Clear flushes all historical events and alerts.
func (b *Bus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.events = make([]model.Event, 0, 100)
	b.alerts = make([]model.Alert, 0, 50)
}
