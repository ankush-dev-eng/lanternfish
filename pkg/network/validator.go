package network

import (
	"fmt"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// ValidateTopologyIntegrity ensures all nodes have unique identifiers and all links connect valid hosts.
func ValidateTopologyIntegrity(graph *model.NetworkGraph) error {
	if graph == nil {
		return fmt.Errorf("network graph is nil")
	}

	nodeIDs := make(map[string]struct{})
	hostnames := make(map[string]struct{})
	ipAddrs := make(map[string]struct{})

	for _, n := range graph.Nodes {
		if n.ID == "" {
			return fmt.Errorf("node with empty ID detected")
		}
		if _, exists := nodeIDs[n.ID]; exists {
			return fmt.Errorf("duplicate node ID: %s", n.ID)
		}
		nodeIDs[n.ID] = struct{}{}

		if _, exists := hostnames[n.Hostname]; exists {
			return fmt.Errorf("duplicate hostname: %s", n.Hostname)
		}
		hostnames[n.Hostname] = struct{}{}

		if _, exists := ipAddrs[n.IP]; exists {
			return fmt.Errorf("duplicate IP address: %s", n.IP)
		}
		ipAddrs[n.IP] = struct{}{}
	}

	for _, link := range graph.Links {
		if _, ok := nodeIDs[link.Source]; !ok {
			return fmt.Errorf("link %s references non-existent source node: %s", link.ID, link.Source)
		}
		if _, ok := nodeIDs[link.Target]; !ok {
			return fmt.Errorf("link %s references non-existent target node: %s", link.ID, link.Target)
		}
	}

	return nil
}
