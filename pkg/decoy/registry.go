package decoy

import (
	"sync"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// Registry manages thread-safe decoy inventory state and tripwire activations.
type Registry struct {
	mu     sync.RWMutex
	decoys map[string]*model.Decoy
}

// NewRegistry initializes a registry populated with the default seeded decoys.
func NewRegistry() *Registry {
	r := &Registry{
		decoys: make(map[string]*model.Decoy),
	}
	r.Reset()
	return r
}

// Reset re-initializes all decoys to their untriggered starting state.
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.decoys = make(map[string]*model.Decoy)
	for _, d := range GenerateDefaultDecoys() {
		copyDecoy := d
		r.decoys[copyDecoy.ID] = &copyDecoy
	}
}

// GetAll returns a slice of all decoys.
func (r *Registry) GetAll() []model.Decoy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Decoy, 0, len(r.decoys))
	for _, d := range r.decoys {
		result = append(result, *d)
	}
	return result
}

// FindByID returns a decoy by its ID.
func (r *Registry) FindByID(id string) *model.Decoy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if d, ok := r.decoys[id]; ok {
		copyDecoy := *d
		return &copyDecoy
	}
	return nil
}

// FindByTripwireID locates a decoy by its tripwire identifier.
func (r *Registry) FindByTripwireID(tripwireID string) *model.Decoy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, d := range r.decoys {
		if d.TripwireID == tripwireID {
			copyDecoy := *d
			return &copyDecoy
		}
	}
	return nil
}

// FindByPathOrName looks up a decoy matching a file path, credential name, or hostname.
func (r *Registry) FindByPathOrName(targetNode, identifier string) *model.Decoy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, d := range r.decoys {
		if (d.NodeID == targetNode || targetNode == "") &&
			(d.PathOrIdentifier == identifier || d.Name == identifier) {
			copyDecoy := *d
			return &copyDecoy
		}
	}
	return nil
}

// TriggerTripwire activates the decoy lure, recording timestamps and caller notes.
func (r *Registry) TriggerTripwire(tripwireID, actor, notes string) (*model.Decoy, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, d := range r.decoys {
		if d.TripwireID == tripwireID {
			now := time.Now()
			d.Triggered = true
			d.TriggerCount++
			if d.FirstTriggeredAt == nil {
				d.FirstTriggeredAt = &now
			}
			d.LastTriggeredAt = &now
			d.AttackerNotes = notes
			copyDecoy := *d
			return &copyDecoy, true
		}
	}
	return nil, false
}

// TriggeredCount returns the number of unique decoys that have been tripped.
func (r *Registry) TriggeredCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, d := range r.decoys {
		if d.Triggered {
			count++
		}
	}
	return count
}
