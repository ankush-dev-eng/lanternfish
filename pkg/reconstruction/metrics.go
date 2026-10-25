package reconstruction

import (
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// TraversalVelocityMetrics measures the operational pace of an intruder's lateral movement.
type TraversalVelocityMetrics struct {
	TotalHops        int     `json:"total_hops"`
	DwellSeconds     float64 `json:"dwell_seconds"`
	HopsPerMinute    float64 `json:"hops_per_minute"`
	VelocityCategory string  `json:"velocity_category"` // "STEALTHY_SLOW", "STANDARD_PACED", "HIGH_VELOCITY_SPRAY"
}

// CalculateHopVelocity derives velocity metrics from the reconstructed attacker path.
func CalculateHopVelocity(path model.AttackerPath) TraversalVelocityMetrics {
	totalHops := len(path.Hops)
	dwellSecs := path.DwellTimeSeconds

	var hpm float64
	if dwellSecs > 0 && totalHops > 0 {
		hpm = float64(totalHops) / (dwellSecs / 60.0)
	}

	category := "STEALTHY_SLOW"
	if hpm > 25.0 {
		category = "HIGH_VELOCITY_SPRAY"
	} else if hpm > 6.0 {
		category = "STANDARD_PACED"
	}

	return TraversalVelocityMetrics{
		TotalHops:        totalHops,
		DwellSeconds:     dwellSecs,
		HopsPerMinute:    hpm,
		VelocityCategory: category,
	}
}
