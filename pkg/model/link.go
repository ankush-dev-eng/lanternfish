package model

// LinkType categorizes the communication channel between network nodes.
type LinkType string

const (
	LinkTypeLAN       LinkType = "lan"
	LinkTypeDMZ       LinkType = "dmz"
	LinkTypeAdminVPN  LinkType = "admin_vpn"
	LinkTypeReplication LinkType = "replication"
)

// Link represents an active or permissible network route between two nodes.
type Link struct {
	ID       string   `json:"id"`
	Source   string   `json:"source"`
	Target   string   `json:"target"`
	Type     LinkType `json:"type"`
	Protocol string   `json:"protocol"`
	Port     int      `json:"port"`
}
