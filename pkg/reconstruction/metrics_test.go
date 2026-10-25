package reconstruction

import (
	"testing"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestCalculateHopVelocity(t *testing.T) {
	// 1. Stealthy path
	stealthyPath := model.AttackerPath{
		Hops: []model.ReconstructedHop{
			{HopIndex: 1},
			{HopIndex: 2},
		},
		DwellTimeSeconds: 120.0, // 2 hops in 2 minutes = 1 hop/min
	}
	v1 := CalculateHopVelocity(stealthyPath)
	if v1.VelocityCategory != "STEALTHY_SLOW" {
		t.Errorf("expected STEALTHY_SLOW, got %s", v1.VelocityCategory)
	}
	if v1.HopsPerMinute != 1.0 {
		t.Errorf("expected 1.0 HPM, got %f", v1.HopsPerMinute)
	}

	// 2. High velocity spray path
	noisyPath := model.AttackerPath{
		Hops: []model.ReconstructedHop{
			{HopIndex: 1}, {HopIndex: 2}, {HopIndex: 3}, {HopIndex: 4}, {HopIndex: 5},
		},
		DwellTimeSeconds: 6.0, // 5 hops in 6 seconds = 50 hops/min
	}
	v2 := CalculateHopVelocity(noisyPath)
	if v2.VelocityCategory != "HIGH_VELOCITY_SPRAY" {
		t.Errorf("expected HIGH_VELOCITY_SPRAY, got %s", v2.VelocityCategory)
	}
}
