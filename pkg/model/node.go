package model

// NodeRole represents the operational function of a simulated network host.
type NodeRole string

const (
	RoleGateway          NodeRole = "gateway"
	RoleWebServer        NodeRole = "web_server"
	RoleAppServer        NodeRole = "app_server"
	RoleDatabase         NodeRole = "database"
	RoleFileServer       NodeRole = "file_server"
	RoleDomainController NodeRole = "domain_controller"
	RoleBackupVault      NodeRole = "backup_vault"
	RoleWorkstation      NodeRole = "workstation"
	RoleDecoyHost        NodeRole = "decoy_host"
)

// NodePosition holds layout coordinates for 2D/Canvas rendering.
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Node represents a computing asset in the simulated enterprise topology.
type Node struct {
	ID          string       `json:"id"`
	Hostname    string       `json:"hostname"`
	IP          string       `json:"ip"`
	Subnet      string       `json:"subnet"`
	OS          string       `json:"os"`
	Role        NodeRole     `json:"role"`
	Services    []string     `json:"services"`
	IsDecoy     bool         `json:"is_decoy"`
	TripwireID  string       `json:"tripwire_id,omitempty"`
	Position    NodePosition `json:"position"`
	Description string       `json:"description"`
	DecoyAssets []string     `json:"decoy_assets,omitempty"`
}
