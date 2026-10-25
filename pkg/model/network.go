package model

// SubnetDef describes a logical network segment.
type SubnetDef struct {
	Name        string `json:"name"`
	CIDR        string `json:"cidr"`
	Zone        string `json:"zone"`
	Description string `json:"description"`
}

// NetworkGraph represents the complete simulated enterprise infrastructure.
type NetworkGraph struct {
	Nodes   []Node      `json:"nodes"`
	Links   []Link      `json:"links"`
	Subnets []SubnetDef `json:"subnets"`
}

// FindNode searches the graph for a node by ID.
func (g *NetworkGraph) FindNode(id string) *Node {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return &g.Nodes[i]
		}
	}
	return nil
}

// FindNodeByHostname searches the graph for a node by hostname.
func (g *NetworkGraph) FindNodeByHostname(hostname string) *Node {
	for i := range g.Nodes {
		if g.Nodes[i].Hostname == hostname {
			return &g.Nodes[i]
		}
	}
	return nil
}
