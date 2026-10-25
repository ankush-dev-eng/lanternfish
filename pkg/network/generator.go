package network

import (
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// GenerateEnterpriseTopology builds a realistic ~16 node simulated enterprise network graph.
func GenerateEnterpriseTopology() *model.NetworkGraph {
	subnets := []model.SubnetDef{
		{
			Name:        "Perimeter DMZ",
			CIDR:        "10.10.0.0/24",
			Zone:        "DMZ",
			Description: "Ingress security filters, WAF, and boundary gateways",
		},
		{
			Name:        "DMZ Web Tier",
			CIDR:        "10.10.10.0/24",
			Zone:        "DMZ",
			Description: "Publicly accessible load-balanced application reverse proxies",
		},
		{
			Name:        "Internal App Tier",
			CIDR:        "10.10.20.0/24",
			Zone:        "Internal",
			Description: "Core microservice business logic and API orchestration",
		},
		{
			Name:        "Secure Core DB",
			CIDR:        "10.10.30.0/24",
			Zone:        "Restricted",
			Description: "PCI-DSS database cluster and stateful caching memory tiers",
		},
		{
			Name:        "Internal Services",
			CIDR:        "10.10.40.0/24",
			Zone:        "Internal",
			Description: "Domain Controllers, Kerberos KDC, SMB file clusters",
		},
		{
			Name:        "Isolated Vault",
			CIDR:        "10.10.50.0/24",
			Zone:        "Airgapped",
			Description: "Immutable zero-trust offline backup archives",
		},
		{
			Name:        "Workstation LAN",
			CIDR:        "10.10.60.0/24",
			Zone:        "Corporate",
			Description: "End-user developer, executive, and financial endpoints",
		},
	}

	nodes := []model.Node{
		{
			ID:          "gw-edge-01",
			Hostname:    "gw-edge-01.corp",
			IP:          "10.10.0.1",
			Subnet:      "Perimeter DMZ",
			OS:          "pfSense 2.7 / FreeBSD",
			Role:        model.RoleGateway,
			Services:    []string{"HTTPS:443", "BGP", "WireGuard:51820"},
			Position:    model.NodePosition{X: 100, Y: 280},
			Description: "Perimeter ingress firewall & boundary routing proxy",
		},
		{
			ID:          "web-prod-01",
			Hostname:    "web-prod-01.corp",
			IP:          "10.10.10.11",
			Subnet:      "DMZ Web Tier",
			OS:          "Ubuntu 22.04 LTS",
			Role:        model.RoleWebServer,
			Services:    []string{"Nginx 1.24", "TLS 1.3", "HTTP:80"},
			Position:    model.NodePosition{X: 240, Y: 160},
			Description: "Public web frontend reverse proxy node A",
		},
		{
			ID:          "web-prod-02",
			Hostname:    "web-prod-02.corp",
			IP:          "10.10.10.12",
			Subnet:      "DMZ Web Tier",
			OS:          "Ubuntu 22.04 LTS",
			Role:        model.RoleWebServer,
			Services:    []string{"Nginx 1.24", "TLS 1.3", "HTTP:80"},
			Position:    model.NodePosition{X: 240, Y: 400},
			Description: "Public web frontend reverse proxy node B",
		},
		{
			ID:          "api-gw-01",
			Hostname:    "api-gw-01.corp",
			IP:          "10.10.20.5",
			Subnet:      "Internal App Tier",
			OS:          "Alpine Linux 3.19",
			Role:        model.RoleAppServer,
			Services:    []string{"Envoy Proxy:8443", "gRPC:50051"},
			Position:    model.NodePosition{X: 380, Y: 280},
			Description: "Central internal API router & identity token verifier",
		},
		{
			ID:          "app-srv-01",
			Hostname:    "app-srv-01.corp",
			IP:          "10.10.20.21",
			Subnet:      "Internal App Tier",
			OS:          "Debian 12 Bookworm",
			Role:        model.RoleAppServer,
			Services:    []string{"Go Microservices:8080", "Docker Daemon"},
			Position:    model.NodePosition{X: 520, Y: 180},
			Description: "Payment processing & checkout orchestration cluster node",
		},
		{
			ID:          "app-srv-02",
			Hostname:    "app-srv-02.corp",
			IP:          "10.10.20.22",
			Subnet:      "Internal App Tier",
			OS:          "Debian 12 Bookworm",
			Role:        model.RoleAppServer,
			Services:    []string{"Customer Portal:3000", "Node.js 20"},
			Position:    model.NodePosition{X: 520, Y: 380},
			Description: "Customer account portal backend services",
		},
		{
			ID:          "db-postgres-01",
			Hostname:    "db-postgres-01.corp",
			IP:          "10.10.30.15",
			Subnet:      "Secure Core DB",
			OS:          "Rocky Linux 9.3",
			Role:        model.RoleDatabase,
			Services:    []string{"PostgreSQL 16:5432", "pgBouncer:6432"},
			Position:    model.NodePosition{X: 680, Y: 180},
			Description: "Primary transactional customer & transaction database",
		},
		{
			ID:          "db-redis-cache",
			Hostname:    "db-redis-cache.corp",
			IP:          "10.10.30.18",
			Subnet:      "Secure Core DB",
			OS:          "Alpine Linux 3.19",
			Role:        model.RoleDatabase,
			Services:    []string{"Redis Sentinel:26379", "Redis Server:6379"},
			Position:    model.NodePosition{X: 680, Y: 380},
			Description: "In-memory authentication token & session cache",
		},
		{
			ID:          "file-share-corp",
			Hostname:    "file-share-corp.corp",
			IP:          "10.10.40.5",
			Subnet:      "Internal Services",
			OS:          "Windows Server 2022",
			Role:        model.RoleFileServer,
			Services:    []string{"SMBv3:445", "DFS Replication", "RPC:135"},
			Position:    model.NodePosition{X: 580, Y: 520},
			Description: "Enterprise corporate document repository & DFS share",
		},
		{
			ID:          "dc-corp-01",
			Hostname:    "dc-corp-01.corp",
			IP:          "10.10.40.10",
			Subnet:      "Internal Services",
			OS:          "Windows Server 2022",
			Role:        model.RoleDomainController,
			Services:    []string{"Active Directory:389", "LDAPS:636", "Kerberos:88", "DNS:53"},
			Position:    model.NodePosition{X: 820, Y: 280},
			Description: "Primary Active Directory Domain Controller (Forest Root)",
		},
		{
			ID:          "backup-vault-01",
			Hostname:    "backup-vault-01.corp",
			IP:          "10.10.50.99",
			Subnet:      "Isolated Vault",
			OS:          "FreeBSD 14 / OpenZFS",
			Role:        model.RoleBackupVault,
			Services:    []string{"SSH Bastion:2222", "ZFS Send/Receive", "Rsync:873"},
			Position:    model.NodePosition{X: 960, Y: 280},
			Description: "Immutable worm-storage backup vault with strict ACLs",
		},
		{
			ID:          "ws-dev-02",
			Hostname:    "ws-dev-02.corp",
			IP:          "10.10.60.12",
			Subnet:      "Workstation LAN",
			OS:          "macOS Sonoma 14.4",
			Role:        model.RoleWorkstation,
			Services:    []string{"SSH Client", "Docker Desktop", "Git CLI"},
			Position:    model.NodePosition{X: 240, Y: 560},
			Description: "Senior Platform Engineer development laptop",
		},
		{
			ID:          "ws-finance-01",
			Hostname:    "ws-finance-01.corp",
			IP:          "10.10.60.33",
			Subnet:      "Workstation LAN",
			OS:          "Windows 11 Enterprise",
			Role:        model.RoleWorkstation,
			Services:    []string{"Excel 365", "SAP GUI:3200", "Outlook:993"},
			Position:    model.NodePosition{X: 380, Y: 560},
			Description: "Chief Financial Officer workstation endpoint",
		},
		{
			ID:          "ws-sec-ops",
			Hostname:    "ws-sec-ops.corp",
			IP:          "10.10.60.88",
			Subnet:      "Workstation LAN",
			OS:          "Fedora 40 Workstation",
			Role:        model.RoleWorkstation,
			Services:    []string{"OpenSSH:22", "Wireshark", "Auditd"},
			Position:    model.NodePosition{X: 100, Y: 560},
			Description: "Security Analyst operational jumpbox",
		},
		// Decoy / Honeytoken hosts placed naturally into topology
		{
			ID:          "db-legacy-backup",
			Hostname:    "db-legacy-backup.corp",
			IP:          "10.10.30.91",
			Subnet:      "Secure Core DB",
			OS:          "CentOS 7.9 (Unpatched)",
			Role:        model.RoleDecoyHost,
			Services:    []string{"MySQL 5.6:3306", "Telnet:23", "FTP:21"},
			IsDecoy:     true,
			TripwireID:  "TRIP-HOST-DBLEGACY",
			Position:    model.NodePosition{X: 740, Y: 60},
			Description: "Apparent legacy unmaintained database backup (Bait Lure)",
			DecoyAssets: []string{"honeytoken:mysql_root_legacy", "file:dump_2022_customer_pii.sql"},
		},
		{
			ID:          "finance-share-archive",
			Hostname:    "finance-share-archive.corp",
			IP:          "10.10.40.89",
			Subnet:      "Internal Services",
			OS:          "Windows Server 2016",
			Role:        model.RoleDecoyHost,
			Services:    []string{"SMBv1:445", "HTTP:8080 (Directory Listing)"},
			IsDecoy:     true,
			TripwireID:  "TRIP-HOST-FINARCHIVE",
			Position:    model.NodePosition{X: 740, Y: 520},
			Description: "Unrestricted internal confidential financial archive (Bait Lure)",
			DecoyAssets: []string{"file:payroll_q3_unredacted.xlsx", "honeytoken:cfo_vpn_access.key"},
		},
	}

	links := []model.Link{
		// Perimeter ingress
		{ID: "link-01", Source: "gw-edge-01", Target: "web-prod-01", Type: model.LinkTypeDMZ, Protocol: "TCP", Port: 443},
		{ID: "link-02", Source: "gw-edge-01", Target: "web-prod-02", Type: model.LinkTypeDMZ, Protocol: "TCP", Port: 443},
		{ID: "link-03", Source: "gw-edge-01", Target: "ws-sec-ops", Type: model.LinkTypeAdminVPN, Protocol: "WireGuard", Port: 51820},

		// Web to App tier
		{ID: "link-04", Source: "web-prod-01", Target: "api-gw-01", Type: model.LinkTypeLAN, Protocol: "HTTPS", Port: 8443},
		{ID: "link-05", Source: "web-prod-02", Target: "api-gw-01", Type: model.LinkTypeLAN, Protocol: "HTTPS", Port: 8443},

		// API gateway to App services
		{ID: "link-06", Source: "api-gw-01", Target: "app-srv-01", Type: model.LinkTypeLAN, Protocol: "gRPC", Port: 50051},
		{ID: "link-07", Source: "api-gw-01", Target: "app-srv-02", Type: model.LinkTypeLAN, Protocol: "gRPC", Port: 50051},

		// App services to DB Core
		{ID: "link-08", Source: "app-srv-01", Target: "db-postgres-01", Type: model.LinkTypeLAN, Protocol: "PostgreSQL", Port: 5432},
		{ID: "link-09", Source: "app-srv-01", Target: "db-redis-cache", Type: model.LinkTypeLAN, Protocol: "Redis", Port: 6379},
		{ID: "link-10", Source: "app-srv-02", Target: "db-postgres-01", Type: model.LinkTypeLAN, Protocol: "PostgreSQL", Port: 5432},
		{ID: "link-11", Source: "app-srv-02", Target: "db-redis-cache", Type: model.LinkTypeLAN, Protocol: "Redis", Port: 6379},

		// App & Corporate Shares
		{ID: "link-12", Source: "app-srv-01", Target: "file-share-corp", Type: model.LinkTypeLAN, Protocol: "SMB", Port: 445},

		// Active Directory links
		{ID: "link-13", Source: "dc-corp-01", Target: "db-postgres-01", Type: model.LinkTypeLAN, Protocol: "Kerberos", Port: 88},
		{ID: "link-14", Source: "dc-corp-01", Target: "file-share-corp", Type: model.LinkTypeLAN, Protocol: "LDAP", Port: 389},
		{ID: "link-15", Source: "dc-corp-01", Target: "app-srv-01", Type: model.LinkTypeLAN, Protocol: "LDAP", Port: 389},

		// Backup Vault synchronization
		{ID: "link-16", Source: "backup-vault-01", Target: "db-postgres-01", Type: model.LinkTypeReplication, Protocol: "SSH", Port: 2222},
		{ID: "link-17", Source: "backup-vault-01", Target: "file-share-corp", Type: model.LinkTypeReplication, Protocol: "Rsync", Port: 873},

		// Workstation connections
		{ID: "link-18", Source: "ws-dev-02", Target: "api-gw-01", Type: model.LinkTypeLAN, Protocol: "HTTPS", Port: 8443},
		{ID: "link-19", Source: "ws-dev-02", Target: "file-share-corp", Type: model.LinkTypeLAN, Protocol: "SMB", Port: 445},
		{ID: "link-20", Source: "ws-finance-01", Target: "file-share-corp", Type: model.LinkTypeLAN, Protocol: "SMB", Port: 445},
		{ID: "link-21", Source: "ws-finance-01", Target: "dc-corp-01", Type: model.LinkTypeLAN, Protocol: "Kerberos", Port: 88},
		{ID: "link-22", Source: "ws-sec-ops", Target: "dc-corp-01", Type: model.LinkTypeAdminVPN, Protocol: "RDP", Port: 3389},

		// Decoy Host Connections (visible in topology, bait paths)
		{ID: "link-23", Source: "app-srv-01", Target: "db-legacy-backup", Type: model.LinkTypeLAN, Protocol: "MySQL", Port: 3306},
		{ID: "link-24", Source: "file-share-corp", Target: "finance-share-archive", Type: model.LinkTypeLAN, Protocol: "SMB", Port: 445},
		{ID: "link-25", Source: "ws-finance-01", Target: "finance-share-archive", Type: model.LinkTypeLAN, Protocol: "SMB", Port: 445},
		{ID: "link-26", Source: "ws-dev-02", Target: "db-legacy-backup", Type: model.LinkTypeLAN, Protocol: "TCP", Port: 3306},
	}

	return &model.NetworkGraph{
		Nodes:   nodes,
		Links:   links,
		Subnets: subnets,
	}
}
