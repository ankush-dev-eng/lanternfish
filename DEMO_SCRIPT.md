# LANTERNFISH — 90-Second Demo Video Script

**Target Duration:** 90 seconds  
**Tone:** Authoritative, technical, focused, nautical field-instrument theme  
**Visual Setting:** Browser open at `http://localhost:8080` in fullscreen  

---

### [00:00 - 00:15] Hook & Problem Statement
- **Visual:** Camera starts on the LANTERNFISH dashboard. The animated logo gently glows; the 16-node topology is centered on the canvas with bathymetric depth rings.
- **Voiceover:**  
  *"In the deepest ocean trenches, the anglerfish uses a glowing lure to expose predators in the dark. In enterprise networks, stealthy intruders evade traditional firewalls for weeks using living-off-the-land techniques. This is LANTERNFISH — a cyber deception and intruder tracing platform built with zero third-party dependencies in Go."*

---

### [00:15 - 00:30] Decoys & The Operator Reveal Mode
- **Visual:** Mouse clicks the "DECOYS" tab showing files and credentials that look like regular production assets. Then the operator flips the **"REVEAL DECOYS"** switch.
- **Voiceover:**  
  *"Adversaries see convincing digital assets: spreadsheets, SSH keys, and legacy database hosts. But to our security operators, a single toggle spotlights the hidden tripwires. Every decoy is seeded with canary tokens with zero legitimate operational use — guaranteeing 99% confidence and zero false positives."*

---

### [00:30 - 00:50] The Intrusion Simulation & Tripwire Bite
- **Visual:** Select "APT Stealthy Intrusion", click **"SIMULATE"**. The topology animates. A red pulse appears on `ws-dev-02` (Patient Zero), hops to `api-gw-01`, and then an intruder touches `passwords.xlsx`. Instantly, an amber glow flashes on the canvas!
- **Voiceover:**  
  *"Let's launch an APT simulation. Patient Zero is compromised via phishing. The adversary pivots laterally across the internal API gateway. Now watch: the moment the attacker reads `passwords_2024_master.xlsx`, the lure trips! The alert fires within 1.2 seconds, tagged with MITRE ATT&CK technique T1083."*

---

### [00:50 - 01:10] Trajectory Reconstruction & Live Canvas Animation
- **Visual:** The canvas draws fading directional trails from `ws-dev-02` through `api-gw-01` into `app-srv-01` and `db-legacy-backup`. The HUD updates TTD (1.20s), Hops Traced (4), and Decoys Triggered (3).
- **Voiceover:**  
  *"Rather than alerting in isolation, LANTERNFISH reconstructs the full kill-chain trajectory. Our pathfinding engine traces back to Patient Zero, maps every intermediate hop, and calculates exact adversary dwell time — animated live over our Server-Sent Events stream."*

---

### [01:10 - 01:25] AI Forensics Dossier & Containment Playbook
- **Visual:** Click over to the **"DOSSIER"** tab. Show the Executive Summary, Root Cause Analysis, and the Actionable Containment Playbook cards. Click **"REPORT"** to show the print-ready HTML dossier.
- **Voiceover:**  
  *"Immediately upon containment, LANTERNFISH synthesizes a complete forensic incident report — powered by Google Gemini or our offline deterministic engine. Security teams get prioritized actions: isolate ws-dev-02, revoke the honeytoken database role, and block external C2 egress."*

---

### [01:25 - 01:30] Closing & Call to Action
- **Visual:** Quick pan back to the live glowing canvas and GitHub repository page.
- **Voiceover:**  
  *"Single binary. Zero external runtime dependencies. Instant deception. Expose the intruder with LANTERNFISH."*
