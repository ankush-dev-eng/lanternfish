package narrative

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
)

// NarrativeOutput holds the generated incident summary, root cause, and recommendations.
type NarrativeOutput struct {
	ExecutiveSummary  string
	RootCauseAnalysis string
	Recommendations   []model.ContainmentRecommendation
	ModelUsed         string
	IsAI              bool
}

// GenerateNarrative creates an incident analysis using Gemini if configured, or deterministic DFIR rules.
func GenerateNarrative(scenario model.ScenarioType, path model.AttackerPath, alerts []model.Alert, decoys []model.Decoy) NarrativeOutput {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey != "" {
		if out, err := callGemini(apiKey, scenario, path, alerts, decoys); err == nil && out.ExecutiveSummary != "" {
			return *out
		}
	}

	return generateDeterministic(scenario, path, alerts, decoys)
}

func generateDeterministic(scenario model.ScenarioType, path model.AttackerPath, alerts []model.Alert, decoys []model.Decoy) NarrativeOutput {
	var execSummary strings.Builder
	var rootCause strings.Builder
	var recs []model.ContainmentRecommendation

	switch scenario {
	case model.ScenarioStealthy:
		execSummary.WriteString(fmt.Sprintf(
			"LANTERNFISH Deception Engine detected an advanced persistent threat (APT) intrusion originating at patient-zero endpoint '%s'. "+
				"The adversary executed low-and-slow reconnaissance before interacting with decoy asset 'passwords_2024_master.xlsx'. "+
				"The intrusion traversed %d lateral hops across internal segments before triggering honeytoken credential 'db_root_backup_svc' and probe trap 'db-legacy-backup'. "+
				"High-confidence tripwire alerts exposed the adversary's dwell footprint across %d hosts within %.1f seconds.",
			path.PatientZeroHostname, len(path.Hops), len(path.NodesVisited), path.DwellTimeSeconds,
		))

		rootCause.WriteString(fmt.Sprintf(
			"Initial compromise occurred via Trojanized developer payload execution on %s (T1566.001). "+
				"The adversary escalated privileges by harvesting seeded decoy credentials, leading to lateral movement via SSH into the internal API gateway "+
				"and subsequent interaction with monitored decoy database instances.",
			path.PatientZeroHostname,
		))

		recs = []model.ContainmentRecommendation{
			{
				Priority:    "IMMEDIATE",
				Target:      path.PatientZeroHostname,
				Action:      "Isolate host from Corporate VLAN via 802.1X / EDR quarantine",
				Rationale:   "Patient zero execution point where initial Trojanized payload established persistence.",
				MitreRemedy: "M1038 - Execution Prevention",
			},
			{
				Priority:    "IMMEDIATE",
				Target:      "db_root_backup_svc",
				Action:      "Revoke and audit all Active Directory and Postgres sessions for honeytoken role",
				Rationale:   "Credential tripwire was actively submitted in database authentication attempts.",
				MitreRemedy: "M1027 - Password Policies",
			},
			{
				Priority:    "HIGH",
				Target:      "Perimeter Boundary (198.51.100.44)",
				Action:      "Block external destination IP on Perimeter Gateway gw-edge-01 firewall",
				Rationale:   "Identified as external Command-and-Control egress beacon destination.",
				MitreRemedy: "M1031 - Network Intrusion Prevention",
			},
			{
				Priority:    "MEDIUM",
				Target:      "api-gw-01 & app-srv-01",
				Action:      "Rotate internal SSH host keys and review transient authorized_keys entries",
				Rationale:   "Intermediate lateral hop nodes utilized by intruder to reach core database segments.",
				MitreRemedy: "M1041 - Account Use Policies",
			},
		}

	case model.ScenarioNoisy:
		execSummary.WriteString(fmt.Sprintf(
			"A high-velocity external cyber intrusion was detected penetrating the DMZ Web perimeter at '%s'. "+
				"The threat actor conducted aggressive port scanning and automated lateral sweeps, tripping multiple deception tripwires including 'id_ed25519_bastion_deploy' "+
				"and decoy server 'finance-share-archive.corp'. Threat containment is urgent due to concurrent data exfiltration attempts.",
			path.PatientZeroHostname,
		))

		rootCause.WriteString(
			"Unpatched buffer vulnerability on public reverse proxy web-prod-01 allowed initial arbitrary command execution (T1190). " +
				"The attacker leveraged cached SSH keys to spray internal microservice tiers before probing unauthenticated decoy shares.",
		)

		recs = []model.ContainmentRecommendation{
			{
				Priority:    "IMMEDIATE",
				Target:      "web-prod-01.corp",
				Action:      "Sever public reverse proxy ingress and patch Nginx reverse proxy buffer vulnerability",
				Rationale:   "External ingress patient zero point exploited via HTTP payload buffer overflow.",
				MitreRemedy: "M1048 - Application Vulnerability Mitigation",
			},
			{
				Priority:    "IMMEDIATE",
				Target:      "id_ed25519_bastion_deploy",
				Action:      "Invalidate deploy key and prune authorized_keys across entire internal DMZ cluster",
				Rationale:   "Decoy SSH deploy key read from disk and attempted against internal bastions.",
				MitreRemedy: "M1027 - Revoke Compromised Keys",
			},
			{
				Priority:    "HIGH",
				Target:      "DNS Gateway & Perimeter Egress",
				Action:      "Enforce DNS sinkhole and block anomalous TXT chunk tunneling queries",
				Rationale:   "Attacker attempted high-bandwidth exfiltration over DNS protocol channels.",
				MitreRemedy: "M1030 - Network Segmentation",
			},
		}

	case model.ScenarioInsider:
		execSummary.WriteString(fmt.Sprintf(
			"LANTERNFISH identified anomalous insider threat activity originating from endpoint '%s' during off-hours (02:14 AM). "+
				"The rogue actor accessed confidential network shares and opened seeded decoy document 'payroll_Q3_exec_compensation.csv', "+
				"followed by unauthorized attempts to authenticate as fake domain administrator 'svc_global_admin'.",
			path.PatientZeroHostname,
		))

		rootCause.WriteString(
			"Privileged insider misuse on workstation ws-finance-01. The operator abused local access to harvest sensitive compensation data, " +
				"staged a password-protected 7z archive, and attempted unauthorized transmission to external storage.",
		)

		recs = []model.ContainmentRecommendation{
			{
				Priority:    "IMMEDIATE",
				Target:      "ws-finance-01 (10.10.60.33)",
				Action:      "Revoke Kerberos TGT and suspend user account session immediately",
				Rationale:   "Off-hours access and deliberate interaction with lure executive compensation spreadsheets.",
				MitreRemedy: "M1018 - User Account Management",
			},
			{
				Priority:    "IMMEDIATE",
				Target:      "svc_global_admin",
				Action:      "Verify honeytoken domain account lockout status in Active Directory",
				Rationale:   "Target of unauthorized Kerberoasting / ticket request by internal endpoint.",
				MitreRemedy: "M1026 - Privileged Account Management",
			},
			{
				Priority:    "HIGH",
				Target:      "External IP 203.0.113.88",
				Action:      "Block destination endpoint at gateway proxy and preserve egress upload telemetry for legal review",
				Rationale:   "Identified as personal cloud storage drop utilized for intellectual property exfiltration.",
				MitreRemedy: "M1057 - Data Loss Prevention",
			},
		}
	}

	return NarrativeOutput{
		ExecutiveSummary:  execSummary.String(),
		RootCauseAnalysis: rootCause.String(),
		Recommendations:   recs,
		ModelUsed:         "LANTERNFISH Deterministic Forensics Engine",
		IsAI:              false,
	}
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func callGemini(apiKey string, scenario model.ScenarioType, path model.AttackerPath, alerts []model.Alert, decoys []model.Decoy) (*NarrativeOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey)

	prompt := fmt.Sprintf(`You are a top-tier Lead Cyber Incident Response Investigator and Threat Intelligence Officer.
Analyze this cyber deception detection report from the LANTERNFISH defense platform:

SCENARIO: %s
PATIENT ZERO: %s (%s)
TOTAL LATERAL HOPS: %d
DWELL TIME: %.1f seconds
NODES VISITED: %s
DECOYS TRIGGERED: %s
TOTAL ALERTS: %d

Generate a concise, authoritative DFIR incident summary in strict JSON format with this exact schema:
{
  "executive_summary": "1-2 paragraphs summarizing the intrusion, patient zero, decoys tripped, and threat containment status.",
  "root_cause_analysis": "1 paragraph detailing initial compromise vector and privilege escalation path.",
  "recommendations": [
    {
      "priority": "IMMEDIATE | HIGH | MEDIUM",
      "target": "affected asset name or hostname",
      "action": "clear operational containment action",
      "rationale": "why this action must be taken immediately",
      "mitre_remedy": "e.g. M1038 - Execution Prevention"
    }
  ]
}
Respond with raw JSON only. Do not wrap in markdown quotes.`,
		scenario,
		path.PatientZeroHostname, path.PatientZeroID,
		len(path.Hops),
		path.DwellTimeSeconds,
		strings.Join(path.NodesVisited, ", "),
		strings.Join(path.DecoysTouched, ", "),
		len(alerts),
	)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API returned status %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gResp geminiResponse
	if err := json.Unmarshal(respBytes, &gResp); err != nil {
		return nil, err
	}

	if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	rawText := gResp.Candidates[0].Content.Parts[0].Text
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimPrefix(rawText, "```")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	var parsed struct {
		ExecutiveSummary  string                            `json:"executive_summary"`
		RootCauseAnalysis string                            `json:"root_cause_analysis"`
		Recommendations   []model.ContainmentRecommendation `json:"recommendations"`
	}

	if err := json.Unmarshal([]byte(rawText), &parsed); err != nil {
		return nil, err
	}

	return &NarrativeOutput{
		ExecutiveSummary:  parsed.ExecutiveSummary,
		RootCauseAnalysis: parsed.RootCauseAnalysis,
		Recommendations:   parsed.Recommendations,
		ModelUsed:         "Gemini 2.5 Flash",
		IsAI:              true,
	}, nil
}
