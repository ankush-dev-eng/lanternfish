package network

import (
	"testing"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

func TestValidateTopologyIntegrity(t *testing.T) {
	graph := GenerateEnterpriseTopology()

	// Default generated topology must be valid
	if err := ValidateTopologyIntegrity(graph); err != nil {
		t.Fatalf("expected valid default topology, got: %v", err)
	}

	// Test duplicate node ID detection
	corruptedGraph := *graph
	corruptedGraph.Nodes = append(corruptedGraph.Nodes, model.Node{
		ID:       "gw-edge-01", // duplicate
		Hostname: "duplicate-gw.corp",
		IP:       "10.10.99.99",
	})
	if err := ValidateTopologyIntegrity(&corruptedGraph); err == nil {
		t.Fatal("expected error on duplicate node ID, got nil")
	}

	// Test broken link reference detection
	corruptedLinkGraph := *graph
	corruptedLinkGraph.Links = append(corruptedLinkGraph.Links, model.Link{
		ID:     "broken-link",
		Source: "gw-edge-01",
		Target: "non-existent-node-404",
	})
	if err := ValidateTopologyIntegrity(&corruptedLinkGraph); err == nil {
		t.Fatal("expected error on broken link target reference, got nil")
	}
}
