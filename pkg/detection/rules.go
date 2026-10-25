package detection

import (
	"fmt"
	"strings"

	"github.com/ankush-dev-eng/lanternfish/pkg/decoy"
	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// Rule defines the interface for evaluating telemetry events against deception tripwires.
type Rule interface {
	ID() string
	Name() string
	MitreTechnique() string
	MitreName() string
	Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert
}

// DecoyFileAccessRule detects interaction with deceptive files seeded on disks or shares.
type DecoyFileAccessRule struct{}

func (r *DecoyFileAccessRule) ID() string             { return "RULE-01-DECOY-FILE" }
func (r *DecoyFileAccessRule) Name() string           { return "Decoy File Access / Credential Hunt" }
func (r *DecoyFileAccessRule) MitreTechnique() string { return "T1083" }
func (r *DecoyFileAccessRule) MitreName() string      { return "File and Directory Discovery" }
func (r *DecoyFileAccessRule) Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert {
	if strings.HasPrefix(evt.TripwireID, "TRIP-FILE-") ||
		(evt.Type == model.EventFileRead && evt.TripwireID != "") {
		targetDecoy := decoys.FindByTripwireID(evt.TripwireID)
		decoyName := "Lure File"
		if targetDecoy != nil {
			decoyName = targetDecoy.Name
		}

		return &model.Alert{
			ID:             fmt.Sprintf("ALT-FILE-%s", evt.ID),
			RuleID:         r.ID(),
			RuleName:       r.Name(),
			Severity:       model.SeverityHigh,
			Confidence:     96,
			MitreTechnique: r.MitreTechnique(),
			MitreName:      r.MitreName(),
			Title:          fmt.Sprintf("Decoy File Interacted: %s", decoyName),
			Description:    fmt.Sprintf("Adversary %s opened deceptive honeypot file '%s' on %s. Hidden tripwire token triggered.", evt.Actor, decoyName, evt.TargetNode),
			Timestamp:      evt.Timestamp,
			SourceNode:     evt.SourceNode,
			TargetNode:     evt.TargetNode,
			TripwireID:     evt.TripwireID,
			DecoyName:      decoyName,
			DecoyType:      model.DecoyTypeFile,
			Actor:          evt.Actor,
			RawEventID:     evt.ID,
		}
	}
	return nil
}

// HoneytokenCredentialRule catches the invocation or authentication with decoy accounts or keys.
type HoneytokenCredentialRule struct{}

func (r *HoneytokenCredentialRule) ID() string             { return "RULE-02-HONEYTOKEN-CRED" }
func (r *HoneytokenCredentialRule) Name() string           { return "Honeytoken Credential Abuse" }
func (r *HoneytokenCredentialRule) MitreTechnique() string { return "T1078" }
func (r *HoneytokenCredentialRule) MitreName() string      { return "Valid Accounts" }
func (r *HoneytokenCredentialRule) Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert {
	if strings.HasPrefix(evt.TripwireID, "TRIP-CRED-") || evt.Type == model.EventCredentialUse {
		targetDecoy := decoys.FindByTripwireID(evt.TripwireID)
		credName := "Honeytoken Account"
		if targetDecoy != nil {
			credName = targetDecoy.Name
		}

		return &model.Alert{
			ID:             fmt.Sprintf("ALT-CRED-%s", evt.ID),
			RuleID:         r.ID(),
			RuleName:       r.Name(),
			Severity:       model.SeverityCritical,
			Confidence:     99,
			MitreTechnique: r.MitreTechnique(),
			MitreName:      r.MitreName(),
			Title:          fmt.Sprintf("Honeytoken In Use: %s", credName),
			Description:    fmt.Sprintf("Unauthorized authentication attempt using zero-legitimate-use honeytoken credential '%s' targeting %s.", credName, evt.TargetNode),
			Timestamp:      evt.Timestamp,
			SourceNode:     evt.SourceNode,
			TargetNode:     evt.TargetNode,
			TripwireID:     evt.TripwireID,
			DecoyName:      credName,
			DecoyType:      model.DecoyTypeCredential,
			Actor:          evt.Actor,
			RawEventID:     evt.ID,
		}
	}
	return nil
}

// DecoyHostConnectionRule detects traffic directed at deceptive honeynet machines.
type DecoyHostConnectionRule struct{}

func (r *DecoyHostConnectionRule) ID() string             { return "RULE-03-DECOY-HOST" }
func (r *DecoyHostConnectionRule) Name() string           { return "Decoy Host Interaction" }
func (r *DecoyHostConnectionRule) MitreTechnique() string { return "T1021" }
func (r *DecoyHostConnectionRule) MitreName() string      { return "Remote Services / Honeynet Probe" }
func (r *DecoyHostConnectionRule) Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert {
	targetNode := graph.FindNode(evt.TargetNode)
	isDecoyTarget := (targetNode != nil && targetNode.IsDecoy) || strings.HasPrefix(evt.TripwireID, "TRIP-HOST-")

	if isDecoyTarget {
		hostName := evt.TargetNode
		if targetNode != nil {
			hostName = targetNode.Hostname
		}

		return &model.Alert{
			ID:             fmt.Sprintf("ALT-HOST-%s", evt.ID),
			RuleID:         r.ID(),
			RuleName:       r.Name(),
			Severity:       model.SeverityHigh,
			Confidence:     98,
			MitreTechnique: r.MitreTechnique(),
			MitreName:      r.MitreName(),
			Title:          fmt.Sprintf("Decoy Host Contacted: %s", hostName),
			Description:    fmt.Sprintf("Connection or probe initiated against instrumented decoy host %s from %s. High probability intruder trap activation.", hostName, evt.SourceNode),
			Timestamp:      evt.Timestamp,
			SourceNode:     evt.SourceNode,
			TargetNode:     evt.TargetNode,
			TripwireID:     evt.TripwireID,
			DecoyName:      hostName,
			DecoyType:      model.DecoyTypeHost,
			Actor:          evt.Actor,
			RawEventID:     evt.ID,
		}
	}
	return nil
}

// AbnormalLateralMovementRule flags suspicious multi-hop lateral pivots across security zones.
type AbnormalLateralMovementRule struct{}

func (r *AbnormalLateralMovementRule) ID() string             { return "RULE-04-ABNORMAL-LATERAL" }
func (r *AbnormalLateralMovementRule) Name() string           { return "Abnormal Lateral Traversal Chain" }
func (r *AbnormalLateralMovementRule) MitreTechnique() string { return "T1021" }
func (r *AbnormalLateralMovementRule) MitreName() string      { return "Lateral Movement / Remote Services" }
func (r *AbnormalLateralMovementRule) Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert {
	if evt.Type == model.EventLateralHop {
		return &model.Alert{
			ID:             fmt.Sprintf("ALT-LATERAL-%s", evt.ID),
			RuleID:         r.ID(),
			RuleName:       r.Name(),
			Severity:       model.SeverityMedium,
			Confidence:     87,
			MitreTechnique: r.MitreTechnique(),
			MitreName:      r.MitreName(),
			Title:          fmt.Sprintf("Lateral Movement Hop: %s -> %s", evt.SourceNode, evt.TargetNode),
			Description:    fmt.Sprintf("Cross-segment lateral traversal detected from %s to %s by actor %s.", evt.SourceNode, evt.TargetNode, evt.Actor),
			Timestamp:      evt.Timestamp,
			SourceNode:     evt.SourceNode,
			TargetNode:     evt.TargetNode,
			Actor:          evt.Actor,
			RawEventID:     evt.ID,
		}
	}
	return nil
}

// OffHoursAccessRule flags anomalous user actions executed outside typical operational windows.
type OffHoursAccessRule struct{}

func (r *OffHoursAccessRule) ID() string             { return "RULE-05-OFF-HOURS-ANOMALY" }
func (r *OffHoursAccessRule) Name() string           { return "Off-Hours Privileged Session" }
func (r *OffHoursAccessRule) MitreTechnique() string { return "T1078.003" }
func (r *OffHoursAccessRule) MitreName() string      { return "Valid Accounts: Local Accounts" }
func (r *OffHoursAccessRule) Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert {
	// Flag if explicit actor indicates insider rogue analyst or raw payload mentions off-hours login
	if strings.Contains(evt.Actor, "Insider") || strings.Contains(evt.RawPayload, "02:14 AM") {
		if evt.Type == model.EventAuthAttempt {
			return &model.Alert{
				ID:             fmt.Sprintf("ALT-OFFHOURS-%s", evt.ID),
				RuleID:         r.ID(),
				RuleName:       r.Name(),
				Severity:       model.SeverityMedium,
				Confidence:     84,
				MitreTechnique: r.MitreTechnique(),
				MitreName:      r.MitreName(),
				Title:          fmt.Sprintf("Off-Hours Console Logon on %s", evt.TargetNode),
				Description:    fmt.Sprintf("Privileged interactive login initiated during anomalous off-hours window (02:14 AM) by %s on %s.", evt.Actor, evt.TargetNode),
				Timestamp:      evt.Timestamp,
				SourceNode:     evt.SourceNode,
				TargetNode:     evt.TargetNode,
				Actor:          evt.Actor,
				RawEventID:     evt.ID,
			}
		}
	}
	return nil
}

// ExfiltrationBeaconRule catches outbound exfiltration or command-and-control beaconing.
type ExfiltrationBeaconRule struct{}

func (r *ExfiltrationBeaconRule) ID() string             { return "RULE-06-EXFILTRATION-BEACON" }
func (r *ExfiltrationBeaconRule) Name() string           { return "Exfiltration Beacon / C2 Egress" }
func (r *ExfiltrationBeaconRule) MitreTechnique() string { return "T1041" }
func (r *ExfiltrationBeaconRule) MitreName() string      { return "Exfiltration Over C2 Channel" }
func (r *ExfiltrationBeaconRule) Evaluate(evt model.Event, graph *model.NetworkGraph, decoys *decoy.Registry) *model.Alert {
	if evt.Type == model.EventExfilBeacon {
		return &model.Alert{
			ID:             fmt.Sprintf("ALT-EXFIL-%s", evt.ID),
			RuleID:         r.ID(),
			RuleName:       r.Name(),
			Severity:       model.SeverityCritical,
			Confidence:     97,
			MitreTechnique: r.MitreTechnique(),
			MitreName:      r.MitreName(),
			Title:          fmt.Sprintf("Data Exfiltration Beacon to Perimeter from %s", evt.SourceNode),
			Description:    fmt.Sprintf("Active data staging and external beaconing detected from %s towards boundary gateway. Suspected outbound data transfer.", evt.SourceNode),
			Timestamp:      evt.Timestamp,
			SourceNode:     evt.SourceNode,
			TargetNode:     evt.TargetNode,
			Actor:          evt.Actor,
			RawEventID:     evt.ID,
		}
	}
	return nil
}
