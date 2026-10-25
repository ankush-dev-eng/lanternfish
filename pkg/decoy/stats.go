package decoy

import (
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// CategoryStats summarizes active and triggered counts by decoy artifact class.
type CategoryStats struct {
	TotalFiles        int `json:"total_files"`
	TrippedFiles      int `json:"tripped_files"`
	TotalCredentials  int `json:"total_credentials"`
	TrippedCredentials int `json:"tripped_credentials"`
	TotalHosts        int `json:"total_hosts"`
	TrippedHosts      int `json:"tripped_hosts"`
}

// GetCategoryStats calculates inventory distributions across decoy categories.
func (r *Registry) GetCategoryStats() CategoryStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var stats CategoryStats
	for _, d := range r.decoys {
		switch d.Type {
		case model.DecoyTypeFile:
			stats.TotalFiles++
			if d.Triggered {
				stats.TrippedFiles++
			}
		case model.DecoyTypeCredential:
			stats.TotalCredentials++
			if d.Triggered {
				stats.TrippedCredentials++
			}
		case model.DecoyTypeHost:
			stats.TotalHosts++
			if d.Triggered {
				stats.TrippedHosts++
			}
		}
	}
	return stats
}
