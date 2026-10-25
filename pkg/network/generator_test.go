package network

import (
	"testing"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestGenerateEnterpriseTopology(t *testing.T) {
	graph := GenerateEnterpriseTopology()

	if len(graph.Nodes) < 14 {
		t.Fatalf("expected at least 14 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Subnets) == 0 {
		t.Fatal("expected defined subnets, got 0")
	}

	if len(graph.Links) == 0 {
		t.Fatal("expected network links, got 0")
	}

	// Verify FindNode
	gw := graph.FindNode("gw-edge-01")
	if gw == nil {
		t.Fatal("expected gateway node gw-edge-01 to exist")
	}
	if gw.Role != model.RoleGateway {
		t.Fatalf("expected gateway role, got %s", gw.Role)
	}

	// Verify FindNodeByHostname
	dc := graph.FindNodeByHostname("dc-corp-01.corp")
	if dc == nil {
		t.Fatal("expected domain controller dc-corp-01.corp to exist")
	}

	// Verify decoy hosts are flagged
	decoyCount := 0
	for _, n := range graph.Nodes {
		if n.IsDecoy {
			decoyCount++
			if n.TripwireID == "" {
				t.Errorf("decoy node %s is missing tripwire ID", n.ID)
			}
		}
	}

	if decoyCount < 2 {
		t.Fatalf("expected at least 2 decoy hosts, found %d", decoyCount)
	}
}
