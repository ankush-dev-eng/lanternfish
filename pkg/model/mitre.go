package model

import "fmt"

// MitreTactic represents a high-level cyber kill-chain objective.
type MitreTactic string

const (
	TacticInitialAccess     MitreTactic = "Initial Access"
	TacticDiscovery         MitreTactic = "Discovery"
	TacticCredentialAccess  MitreTactic = "Credential Access"
	TacticLateralMovement   MitreTactic = "Lateral Movement"
	TacticCollection        MitreTactic = "Collection"
	TacticExfiltration      MitreTactic = "Exfiltration"
)

// MitreTechniqueDefinition provides authoritative metadata and documentation links.
type MitreTechniqueDefinition struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Tactic      MitreTactic `json:"tactic"`
	URL         string      `json:"url"`
	Description string      `json:"description"`
}

// GetMitreDefinition returns enriched metadata for standard ATT&CK technique IDs.
func GetMitreDefinition(techniqueID string) MitreTechniqueDefinition {
	switch techniqueID {
	case "T1566.001":
		return MitreTechniqueDefinition{
			ID:          "T1566.001",
			Name:        "Phishing: Spearphishing Attachment",
			Tactic:      TacticInitialAccess,
			URL:         "https://attack.mitre.org/techniques/T1566/001/",
			Description: "Adversaries send targeted emails with malicious attachments to gain code execution.",
		}
	case "T1083":
		return MitreTechniqueDefinition{
			ID:          "T1083",
			Name:        "File and Directory Discovery",
			Tactic:      TacticDiscovery,
			URL:         "https://attack.mitre.org/techniques/T1083/",
			Description: "Adversaries enumerate files and directories to find sensitive documents or credentials.",
		}
	case "T1078", "T1078.003":
		return MitreTechniqueDefinition{
			ID:          techniqueID,
			Name:        "Valid Accounts",
			Tactic:      TacticCredentialAccess,
			URL:         "https://attack.mitre.org/techniques/T1078/",
			Description: "Adversaries abuse legitimate or honeytoken credentials to bypass access controls.",
		}
	case "T1021", "T1021.001", "T1021.002", "T1021.004":
		return MitreTechniqueDefinition{
			ID:          techniqueID,
			Name:        "Remote Services",
			Tactic:      TacticLateralMovement,
			URL:         "https://attack.mitre.org/techniques/T1021/",
			Description: "Adversaries log into remote network services like SSH, RDP, or SMB to pivot laterally.",
		}
	case "T1041":
		return MitreTechniqueDefinition{
			ID:          "T1041",
			Name:        "Exfiltration Over C2 Channel",
			Tactic:      TacticExfiltration,
			URL:         "https://attack.mitre.org/techniques/T1041/",
			Description: "Adversaries steal sensitive data by transmitting it over existing command-and-control channels.",
		}
	default:
		return MitreTechniqueDefinition{
			ID:          techniqueID,
			Name:        "Unclassified Technique",
			Tactic:      TacticDiscovery,
			URL:         fmt.Sprintf("https://attack.mitre.org/search?q=%s", techniqueID),
			Description: "Enterprise threat actor operational procedure.",
		}
	}
}
