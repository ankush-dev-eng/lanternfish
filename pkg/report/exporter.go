package report

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ankush-dev-eng/lanternfish/pkg/model"
	"github.com/ankush-dev-eng/lanternfish/pkg/narrative"
)

// BuildReport constructs an end-to-end incident dossier from session telemetry.
func BuildReport(scenario model.ScenarioType, path model.AttackerPath, alerts []model.Alert, decoys []model.Decoy, timeline []model.TimelineItem, metrics map[string]interface{}) model.IncidentReport {
	narr := narrative.GenerateNarrative(scenario, path, alerts, decoys)

	triggeredDecoys := make([]model.Decoy, 0)
	for _, d := range decoys {
		if d.Triggered {
			triggeredDecoys = append(triggeredDecoys, d)
		}
	}

	return model.IncidentReport{
		ID:                fmt.Sprintf("IR-%d", time.Now().Unix()),
		GeneratedAt:       time.Now(),
		Scenario:          scenario,
		Status:            "CONTAINED_BY_DECEPTION",
		ExecutiveSummary:  narr.ExecutiveSummary,
		RootCauseAnalysis: narr.RootCauseAnalysis,
		IsAIGenerated:     narr.IsAI,
		NarrativeModel:    narr.ModelUsed,
		AttackerPath:      path,
		Alerts:            alerts,
		DecoysTriggered:   triggeredDecoys,
		Timeline:          timeline,
		Recommendations:   narr.Recommendations,
		Metrics:           metrics,
	}
}

// ExportJSON serializes the incident report to indented JSON bytes.
func ExportJSON(rep model.IncidentReport) ([]byte, error) {
	return json.MarshalIndent(rep, "", "  ")
}

// ExportMarkdown renders a clean, professional DFIR report in GitHub-flavored Markdown.
func ExportMarkdown(rep model.IncidentReport) string {
	var sb strings.Builder

	sb.WriteString("# 🐟 LANTERNFISH Incident Forensics Report\n\n")
	sb.WriteString(fmt.Sprintf("**Incident ID:** `%s`  \n", rep.ID))
	sb.WriteString(fmt.Sprintf("**Generated At:** `%s`  \n", rep.GeneratedAt.UTC().Format("2006-01-02 15:04:05 UTC")))
	sb.WriteString(fmt.Sprintf("**Attack Scenario:** `%s`  \n", rep.Scenario))
	sb.WriteString(fmt.Sprintf("**Defense Status:** `%s`  \n", rep.Status))
	sb.WriteString(fmt.Sprintf("**Intelligence Engine:** `%s` (AI: `%t`)  \n\n", rep.NarrativeModel, rep.IsAIGenerated))

	sb.WriteString("## 1. Executive Summary\n\n")
	sb.WriteString(rep.ExecutiveSummary)
	sb.WriteString("\n\n")

	sb.WriteString("## 2. Root Cause Analysis\n\n")
	sb.WriteString(rep.RootCauseAnalysis)
	sb.WriteString("\n\n")

	sb.WriteString("## 3. Key Incident Metrics\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|---|---|\n")
	sb.WriteString(fmt.Sprintf("| Patient Zero | `%s` (%s) |\n", rep.AttackerPath.PatientZeroHostname, rep.AttackerPath.PatientZeroID))
	sb.WriteString(fmt.Sprintf("| Traversal Hops | %d |\n", len(rep.AttackerPath.Hops)))
	sb.WriteString(fmt.Sprintf("| Compromised / Visited Nodes | %d |\n", len(rep.AttackerPath.NodesVisited)))
	sb.WriteString(fmt.Sprintf("| Decoy Assets Triggered | %d |\n", len(rep.DecoysTriggered)))
	sb.WriteString(fmt.Sprintf("| High-Fidelity Alerts | %d |\n", len(rep.Alerts)))
	sb.WriteString(fmt.Sprintf("| Adversary Dwell Time | %.1fs |\n\n", rep.AttackerPath.DwellTimeSeconds))

	sb.WriteString("## 4. Reconstructed Attacker Trajectory\n\n")
	if len(rep.AttackerPath.Hops) == 0 {
		sb.WriteString("*No lateral hops recorded.*\n\n")
	} else {
		sb.WriteString("| Hop | Source | Target | MITRE Technique | Decoy Tripwire | Method |\n")
		sb.WriteString("|---|---|---|---|---|---|\n")
		for _, h := range rep.AttackerPath.Hops {
			decoyTag := "None"
			if h.DecoyTripped != "" {
				decoyTag = fmt.Sprintf("`%s`", h.DecoyTripped)
			}
			sb.WriteString(fmt.Sprintf("| %d | `%s` | `%s` | `%s` | %s | %s |\n",
				h.HopIndex, h.FromHostname, h.ToHostname, h.MitreTechnique, decoyTag, h.Method))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 5. Decoys & Honeypots Triggered\n\n")
	if len(rep.DecoysTriggered) == 0 {
		sb.WriteString("*No decoys were engaged.*\n\n")
	} else {
		sb.WriteString("| Asset Name | Type | Host | Tripwire ID | Bait Enticement |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, d := range rep.DecoysTriggered {
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | `%s` | %s |\n",
				d.Name, d.Type, d.NodeHostname, d.TripwireID, d.BaitEnticement))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 6. Actionable Containment Recommendations\n\n")
	if len(rep.Recommendations) == 0 {
		sb.WriteString("*No immediate containment actions generated.*\n\n")
	} else {
		for i, rec := range rep.Recommendations {
			sb.WriteString(fmt.Sprintf("### %d. [%s] %s\n\n", i+1, rec.Priority, rec.Action))
			sb.WriteString(fmt.Sprintf("- **Target:** `%s`\n", rec.Target))
			sb.WriteString(fmt.Sprintf("- **Rationale:** %s\n", rec.Rationale))
			sb.WriteString(fmt.Sprintf("- **MITRE Mitigation:** `%s`\n\n", rec.MitreRemedy))
		}
	}

	sb.WriteString("## 7. Chronological Incident Timeline\n\n")
	if len(rep.Timeline) == 0 {
		sb.WriteString("*Timeline is empty.*\n\n")
	} else {
		sb.WriteString("| Time | Phase | Severity | Title | MITRE | Decoy Bait |\n")
		sb.WriteString("|---|---|---|---|---|---|\n")
		for _, tl := range rep.Timeline {
			decoyStr := "-"
			if tl.IsDecoyTrigger {
				decoyStr = fmt.Sprintf("🎣 `%s`", tl.DecoyName)
			}
			sb.WriteString(fmt.Sprintf("| `%s` | %s | **%s** | %s | `%s` | %s |\n",
				tl.Timestamp.Format("15:04:05"), tl.Phase, tl.Severity, tl.Title, tl.MitreTechnique, decoyStr))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n*Report generated autonomously by LANTERNFISH Cyber Deception Platform.*\n")
	return sb.String()
}

// ExportPrintableHTML formats the report into an elegant, print-ready HTML page.
func ExportPrintableHTML(rep model.IncidentReport) string {
	mdContent := ExportMarkdown(rep)

	// Simple styled standalone HTML wrapper using the deep-sea field-instrument palette
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>LANTERNFISH Incident Report - %s</title>
  <style>
    @media print {
      body { background: #fff !important; color: #111 !important; font-size: 11pt; }
      .no-print { display: none !important; }
      table { page-break-inside: avoid; }
      h1, h2, h3 { page-break-after: avoid; }
    }
    body {
      background: #0b1517;
      color: #e9e4d4;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      line-height: 1.6;
      max-width: 900px;
      margin: 40px auto;
      padding: 0 24px;
    }
    h1, h2, h3 {
      font-family: "Fraunces", Georgia, serif;
      color: #e0a33a;
      border-bottom: 1px solid #1f3b36;
      padding-bottom: 8px;
    }
    h1 { font-size: 2rem; }
    h2 { font-size: 1.4rem; margin-top: 32px; }
    h3 { font-size: 1.1rem; color: #c2553a; }
    code, pre {
      font-family: "IBM Plex Mono", Consolas, monospace;
      background: #152724;
      color: #e0a33a;
      padding: 2px 6px;
      border-radius: 3px;
    }
    table {
      width: 100%%;
      border-collapse: collapse;
      margin: 16px 0;
      background: #11201d;
    }
    th, td {
      border: 1px solid #1f3b36;
      padding: 10px 14px;
      text-align: left;
    }
    th {
      background: #18302c;
      color: #8fa8a0;
      font-size: 0.85rem;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    .print-btn {
      background: #e0a33a;
      color: #0b1517;
      border: none;
      padding: 10px 20px;
      font-weight: 700;
      font-family: "IBM Plex Mono", monospace;
      cursor: pointer;
      border-radius: 4px;
      margin-bottom: 24px;
    }
    .print-btn:hover { background: #f3b74f; }
  </style>
</head>
<body>
  <div class="no-print">
    <button class="print-btn" onclick="window.print()">PRINT / SAVE AS PDF</button>
  </div>
  <pre style="white-space: pre-wrap; background: none; color: inherit; padding: 0;">%s</pre>
</body>
</html>`, rep.ID, mdContent)
}
