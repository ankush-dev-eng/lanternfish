/**
 * LANTERNFISH — Frontend Application Controller & Real-Time SSE Consumer
 */

(function () {
  'use strict';

  // Application State
  const state = {
    network: null,
    decoys: [],
    alerts: [],
    timeline: [],
    simState: {
      status: 'idle',
      current_step: 0,
      total_steps: 0,
      time_to_detect_ms: 0,
      alerts_count: 0,
      decoys_triggered: 0,
      hops_count: 0,
      scenario_id: 'stealthy',
    },
    attackerPath: null,
    revealDecoys: false,
    eventSource: null,
  };

  let canvas = null;

  // DOM Elements
  const btnStart = document.getElementById('btnStart');
  const btnPause = document.getElementById('btnPause');
  const btnReset = document.getElementById('btnReset');
  const scenarioSelect = document.getElementById('scenarioSelect');
  const speedSelect = document.getElementById('speedSelect');
  const toggleRevealDecoys = document.getElementById('toggleRevealDecoys');
  const btnCenterGraph = document.getElementById('btnCenterGraph');

  const hudTTD = document.getElementById('hudTTD');
  const hudHops = document.getElementById('hudHops');
  const hudDecoys = document.getElementById('hudDecoys');
  const hudAlerts = document.getElementById('hudAlerts');
  const hudActor = document.getElementById('hudActor');
  const hudLocation = document.getElementById('hudLocation');
  const hudPatientZero = document.getElementById('hudPatientZero');
  const hudSeverityMax = document.getElementById('hudSeverityMax');

  const tabAlertCount = document.getElementById('tabAlertCount');
  const tabDecoyCount = document.getElementById('tabDecoyCount');
  const alertFeedContainer = document.getElementById('alertFeedContainer');
  const decoyListContainer = document.getElementById('decoyListContainer');
  const timelineListContainer = document.getElementById('timelineListContainer');
  const decoyMaskNotice = document.getElementById('decoyMaskNotice');

  const dossierExecSummary = document.getElementById('dossierExecSummary');
  const dossierRootCause = document.getElementById('dossierRootCause');
  const dossierRecommendations = document.getElementById('dossierRecommendations');
  const dossierModelBadge = document.getElementById('dossierModelBadge');

  const nodeModal = document.getElementById('nodeModal');
  const modalNodeTitle = document.getElementById('modalNodeTitle');
  const modalNodeBody = document.getElementById('modalNodeBody');
  const modalCloseBtn = document.getElementById('modalCloseBtn');

  // Initialize
  window.addEventListener('DOMContentLoaded', async () => {
    canvas = new NetworkCanvas('networkCanvas', 'canvasContainer');
    canvas.onNodeClickCallback = openNodeModal;

    setupTabs();
    setupControls();
    await loadInitialData();
    connectSSE();
  });

  // Setup Tabs
  function setupTabs() {
    const tabBtns = document.querySelectorAll('.tab-btn');
    tabBtns.forEach((btn) => {
      btn.addEventListener('click', () => {
        tabBtns.forEach((b) => b.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach((c) => c.classList.remove('active'));

        btn.classList.add('active');
        const targetId = btn.getAttribute('data-tab');
        const targetContent = document.getElementById(targetId);
        if (targetContent) {
          targetContent.classList.add('active');
        }
        if (targetId === 'tabDossier') {
          fetchDossier();
        }
      });
    });
  }

  // Setup User Interaction Controls
  function setupControls() {
    btnStart.addEventListener('click', startSimulation);
    btnPause.addEventListener('click', togglePause);
    btnReset.addEventListener('click', resetSimulation);

    btnCenterGraph.addEventListener('click', () => {
      if (canvas) canvas.centerGraph();
    });

    toggleRevealDecoys.addEventListener('change', (e) => {
      state.revealDecoys = e.target.checked;
      if (canvas) {
        canvas.setRevealDecoys(state.revealDecoys);
      }
      renderDecoys();
      if (decoyMaskNotice) {
        if (state.revealDecoys) {
          decoyMaskNotice.innerHTML = `<strong>OPERATOR MODE: REVEAL ACTIVE.</strong> Hidden tripwire tokens and bait assets are spotlighted with amber halos.`;
          decoyMaskNotice.style.color = '#e0a33a';
        } else {
          decoyMaskNotice.innerHTML = `<strong>OPERATOR MODE:</strong> Showing all digital assets. Toggle <em>"REVEAL DECOYS"</em> above to spotlight hidden lures.`;
          decoyMaskNotice.style.color = '#8fa8a0';
        }
      }
    });

    modalCloseBtn.addEventListener('click', () => {
      nodeModal.close();
    });
    nodeModal.addEventListener('click', (e) => {
      if (e.target === nodeModal) nodeModal.close();
    });
  }

  // API Requests
  async function loadInitialData() {
    try {
      const [netRes, decoyRes, alertRes] = await Promise.all([
        fetch('/api/network'),
        fetch('/api/decoys'),
        fetch('/api/alerts'),
      ]);

      if (netRes.ok) {
        state.network = await netRes.json();
        canvas.setData(state.network);
      }
      if (decoyRes.ok) {
        state.decoys = await decoyRes.json();
        renderDecoys();
        tabDecoyCount.textContent = state.decoys.length;
      }
      if (alertRes.ok) {
        state.alerts = await alertRes.json();
        renderAlerts();
      }
    } catch (err) {
      console.error('Failed to load initial data:', err);
    }
  }

  // Server-Sent Events Telemetry Stream
  function connectSSE() {
    if (state.eventSource) {
      state.eventSource.close();
    }

    state.eventSource = new EventSource('/api/stream');

    state.eventSource.addEventListener('connected', () => {
      document.getElementById('gridStatusDot').className = 'status-dot pulse-amber';
      document.getElementById('gridStatusText').textContent = 'GRID ARMED \u2022 DECOYS ACTIVE';
    });

    state.eventSource.addEventListener('alert', (e) => {
      try {
        const alert = JSON.parse(e.data);
        state.alerts.unshift(alert);
        renderAlerts();
        updateHUD();
      } catch (err) {
        console.error('Error parsing SSE alert:', err);
      }
    });

    state.eventSource.addEventListener('event', () => {
      // Refresh timeline when a new event arrives
      fetchTimeline();
    });

    state.eventSource.addEventListener('state', (e) => {
      try {
        const simState = JSON.parse(e.data);
        state.simState = simState;
        updateSimulationUI(simState);
        updateHUD();
      } catch (err) {
        console.error('Error parsing SSE state:', err);
      }
    });

    state.eventSource.addEventListener('path', (e) => {
      try {
        const path = JSON.parse(e.data);
        state.attackerPath = path;
        canvas.setPath(path);
      } catch (err) {
        console.error('Error parsing SSE path:', err);
      }
    });

    state.eventSource.onerror = () => {
      document.getElementById('gridStatusDot').className = 'status-dot';
      document.getElementById('gridStatusDot').style.background = '#c2553a';
      document.getElementById('gridStatusText').textContent = 'STREAM RECONNECTING...';
    };
  }

  // Actions
  async function startSimulation() {
    const scenario = scenarioSelect.value;
    const speed = parseFloat(speedSelect.value) || 1.0;

    btnStart.disabled = true;
    btnPause.disabled = false;
    btnPause.textContent = 'PAUSE';

    try {
      const res = await fetch('/api/simulate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ scenario, speed }),
      });
      if (!res.ok) throw new Error('Simulation start failed');
    } catch (err) {
      console.error(err);
      btnStart.disabled = false;
    }
  }

  async function togglePause() {
    if (state.simState.status === 'running') {
      btnPause.textContent = 'RESUME';
      // Call pause if needed or wait
    } else if (state.simState.status === 'paused') {
      btnPause.textContent = 'PAUSE';
    }
  }

  async function resetSimulation() {
    try {
      await fetch('/api/reset', { method: 'POST' });
      state.alerts = [];
      state.timeline = [];
      state.attackerPath = null;

      // Reload decoys
      const res = await fetch('/api/decoys');
      if (res.ok) {
        state.decoys = await res.json();
      }

      canvas.setPath(null);
      canvas.setActiveAttacker(null);
      canvas.setTriggeredDecoys([]);

      renderAlerts();
      renderDecoys();
      renderTimeline();
      updateHUD();

      btnStart.disabled = false;
      btnPause.disabled = true;
      btnPause.textContent = 'PAUSE';

      hudTTD.textContent = '--';
      hudHops.textContent = '0';
      hudPatientZero.textContent = 'Patient Zero: None';
      hudActor.textContent = 'STANDBY';
      hudLocation.textContent = 'Location: Awaiting Ingress';
      hudSeverityMax.textContent = 'Max Severity: None';
    } catch (err) {
      console.error('Reset failed:', err);
    }
  }

  function updateSimulationUI(sim) {
    if (sim.status === 'running') {
      btnStart.disabled = true;
      btnPause.disabled = false;
      btnPause.textContent = 'PAUSE';
      canvas.setActiveAttacker(sim.active_attacker_at);
      hudActor.textContent = sim.scenario_id.toUpperCase();
      hudLocation.textContent = `At: ${sim.active_attacker_at || 'Perimeter'}`;
      hudPatientZero.textContent = `Patient Zero: ${sim.patient_zero_node || 'None'}`;
    } else if (sim.status === 'completed') {
      btnStart.disabled = false;
      btnPause.disabled = true;
      hudLocation.textContent = `Completed (${sim.total_steps} steps)`;
      fetchDossier();
    }
  }

  function updateHUD() {
    // TTD
    if (state.simState.time_to_detect_ms > 0) {
      const sec = (state.simState.time_to_detect_ms / 1000).toFixed(2);
      hudTTD.textContent = `${sec}s`;
    } else {
      hudTTD.textContent = '--';
    }

    // Hops
    hudHops.textContent = state.simState.hops_count || 0;

    // Decoys
    const tripped = state.decoys.filter((d) => d.triggered).length;
    hudDecoys.textContent = `${tripped} / ${state.decoys.length}`;

    // Alerts
    hudAlerts.textContent = state.alerts.length;
    tabAlertCount.textContent = state.alerts.length;

    // Max Severity
    let maxSev = 'None';
    for (const a of state.alerts) {
      if (a.severity === 'CRITICAL') {
        maxSev = 'CRITICAL';
        break;
      }
      if (a.severity === 'HIGH' && maxSev !== 'CRITICAL') {
        maxSev = 'HIGH';
      }
      if (a.severity === 'MEDIUM' && maxSev === 'None') {
        maxSev = 'MEDIUM';
      }
    }
    hudSeverityMax.textContent = `Max Severity: ${maxSev}`;
    if (maxSev === 'CRITICAL') {
      hudSeverityMax.className = 'hud-sub text-rust';
    } else if (maxSev === 'HIGH') {
      hudSeverityMax.className = 'hud-sub text-amber';
    } else {
      hudSeverityMax.className = 'hud-sub';
    }

    // Update canvas decoys
    if (canvas) {
      canvas.setTriggeredDecoys(state.decoys);
    }
  }

  // Render Alert Feed
  function renderAlerts() {
    if (state.alerts.length === 0) {
      alertFeedContainer.innerHTML = `
        <div class="empty-state">
          <div class="empty-icon">&#10059;</div>
          <div class="empty-title">ALL TRIPWIRES ARMED</div>
          <div class="empty-desc">No deceptive tripwires or anomalous hops detected yet. Launch a simulation to initiate telemetry ingestion.</div>
        </div>`;
      return;
    }

    let html = '';
    for (const alt of state.alerts) {
      const timeStr = new Date(alt.timestamp).toLocaleTimeString();
      const tripwireBadge = alt.tripwire_id
        ? `<span class="meta-pill meta-tripwire">🎣 ${alt.tripwire_id}</span>`
        : '';
      const mitreBadge = alt.mitre_technique
        ? `<span class="meta-pill">MITRE ${alt.mitre_technique}</span>`
        : '';
      const confidenceBadge = alt.confidence
        ? `<span class="meta-pill">${alt.confidence}% Conf</span>`
        : '';

      html += `
        <div class="alert-card severity-${alt.severity}">
          <div class="alert-header">
            <span class="alert-badge badge-${alt.severity}">${alt.severity}</span>
            <span class="alert-time">${timeStr}</span>
          </div>
          <div class="alert-title">${escapeHTML(alt.title)}</div>
          <div class="alert-desc">${escapeHTML(alt.description)}</div>
          <div class="alert-meta">
            ${tripwireBadge}
            ${mitreBadge}
            ${confidenceBadge}
            <span class="meta-pill">Src: ${escapeHTML(alt.source_node || '-')}</span>
            <span class="meta-pill">Dst: ${escapeHTML(alt.target_node || '-')}</span>
          </div>
        </div>`;
    }
    alertFeedContainer.innerHTML = html;
  }

  // Render Decoy Inventory
  function renderDecoys() {
    if (state.decoys.length === 0) {
      decoyListContainer.innerHTML = `<div class="empty-state"><div class="empty-desc">No decoys loaded.</div></div>`;
      return;
    }

    let html = '';
    for (const d of state.decoys) {
      const isTripped = d.triggered;
      const isRevealed = state.revealDecoys;

      let typePill = isRevealed
        ? `<span class="decoy-pill revealed">🎣 ${d.type.toUpperCase()} LURE</span>`
        : `<span class="decoy-pill">DIGITAL ASSET</span>`;

      let tripStatus = isTripped
        ? `<span class="status-tripped">&#9888; TRIPPED (${d.trigger_count}x)</span>`
        : `<span class="status-armed">&#10003; Armed &amp; Monitored</span>`;

      let tripwireIdDisplay = isRevealed
        ? `<div style="font-size: 9.5px; color: #e0a33a; margin-top: 2px;">Token ID: <code>${d.tripwire_id}</code></div>`
        : '';

      html += `
        <div class="decoy-card ${isTripped ? 'tripped' : ''}">
          <div class="decoy-header">
            <span class="decoy-name">${escapeHTML(d.name)}</span>
            ${typePill}
          </div>
          <div class="decoy-path">${escapeHTML(d.path_or_identifier)} &bull; Host: ${escapeHTML(d.node_hostname)}</div>
          <div class="decoy-enticement">Bait Rationale: "${escapeHTML(d.bait_enticement || d.payload_summary)}"</div>
          <div class="decoy-trip-status">
            ${tripStatus}
          </div>
          ${tripwireIdDisplay}
        </div>`;
    }
    decoyListContainer.innerHTML = html;
  }

  // Fetch & Render Timeline
  async function fetchTimeline() {
    try {
      const res = await fetch('/api/timeline');
      if (res.ok) {
        state.timeline = await res.json();
        renderTimeline();
      }
    } catch (err) {
      console.error('Failed to fetch timeline:', err);
    }
  }

  function renderTimeline() {
    if (!state.timeline || state.timeline.length === 0) {
      timelineListContainer.innerHTML = `
        <div class="empty-state">
          <div class="empty-title">TIMELINE IDLE</div>
          <div class="empty-desc">Incident checkpoints will stream here categorized by cyber kill-chain progression stages.</div>
        </div>`;
      return;
    }

    let html = '';
    for (const item of state.timeline) {
      const timeStr = new Date(item.timestamp).toLocaleTimeString();
      const decoyBadge = item.is_decoy_trigger
        ? `<span class="meta-pill meta-tripwire">🎣 Bait: ${escapeHTML(item.decoy_name || 'Lure')}</span>`
        : '';

      html += `
        <div class="timeline-card ${item.is_decoy_trigger ? 'decoy-engaged' : ''}">
          <div class="timeline-phase">${escapeHTML(item.phase)} &bull; ${timeStr}</div>
          <div class="timeline-title">${escapeHTML(item.title)}</div>
          <div class="timeline-detail">${escapeHTML(item.description)}</div>
          <div class="alert-meta" style="margin-top: 4px;">
            <span class="alert-badge badge-${item.severity}">${item.severity}</span>
            <span class="meta-pill">MITRE ${escapeHTML(item.mitre_technique || 'T1021')}</span>
            ${decoyBadge}
          </div>
        </div>`;
    }
    timelineListContainer.innerHTML = html;
  }

  // Fetch & Render Dossier
  async function fetchDossier() {
    try {
      const res = await fetch('/api/report.json');
      if (!res.ok) return;
      const rep = await res.json();

      if (rep.executive_summary) {
        dossierExecSummary.textContent = rep.executive_summary;
      }
      if (rep.root_cause_analysis) {
        dossierRootCause.textContent = rep.root_cause_analysis;
      }
      if (rep.narrative_model) {
        dossierModelBadge.textContent = `MODEL: ${rep.narrative_model.toUpperCase()}`;
      }

      if (rep.recommendations && rep.recommendations.length > 0) {
        let recHtml = '';
        for (const rec of rep.recommendations) {
          recHtml += `
            <div class="rec-card">
              <div class="rec-header">
                <span class="rec-priority rec-priority-${rec.priority}">${escapeHTML(rec.priority)}</span>
                <span style="font-size: 10px; color: #8fa8a0;">${escapeHTML(rec.target)}</span>
              </div>
              <div class="rec-action">${escapeHTML(rec.action)}</div>
              <div class="rec-desc">${escapeHTML(rec.rationale)}</div>
              <div style="font-size: 9.5px; color: #5c8d89; margin-top: 4px;">Mitigation: ${escapeHTML(rec.mitre_remedy)}</div>
            </div>`;
        }
        dossierRecommendations.innerHTML = recHtml;
      }
    } catch (err) {
      console.error('Failed to load dossier:', err);
    }
  }

  // Node Inspection Modal
  function openNodeModal(node) {
    modalNodeTitle.textContent = `${node.hostname} (${node.ip})`;

    const isTripped = state.decoys.some((d) => d.node_id === node.id && d.triggered);
    const decoysOnNode = state.decoys.filter((d) => d.node_id === node.id);

    let decoySection = '';
    if (decoysOnNode.length > 0) {
      let decoyListHtml = '';
      for (const d of decoysOnNode) {
        decoyListHtml += `<li><strong>${escapeHTML(d.name)}</strong> (${d.type}) — ${d.triggered ? '<span style="color: #e0a33a;">TRIPPED</span>' : 'ARMED'}</li>`;
      }
      decoySection = `
        <div class="modal-field">
          <div class="modal-field-label">INSTRUMENTED DECOYS ON HOST</div>
          <ul style="padding-left: 18px; margin-top: 4px; font-size: 11px; color: #e9e4d4;">
            ${decoyListHtml}
          </ul>
        </div>`;
    }

    modalNodeBody.innerHTML = `
      <div class="modal-field">
        <div class="modal-field-label">SUBNET &amp; SECURITY ZONE</div>
        <div class="modal-field-value">${escapeHTML(node.subnet)}</div>
      </div>
      <div class="modal-field">
        <div class="modal-field-label">ROLE &amp; FUNCTION</div>
        <div class="modal-field-value">${escapeHTML(node.role)} &bull; ${escapeHTML(node.description || '')}</div>
      </div>
      <div class="modal-field">
        <div class="modal-field-label">OPERATING SYSTEM</div>
        <div class="modal-field-value">${escapeHTML(node.os)}</div>
      </div>
      <div class="modal-field">
        <div class="modal-field-label">ACTIVE NETWORK SERVICES</div>
        <div class="modal-field-value">${(node.services || []).map((s) => `<code>${escapeHTML(s)}</code>`).join(' ')}</div>
      </div>
      ${decoySection}
    `;

    nodeModal.showModal();
  }

  function escapeHTML(str) {
    if (!str) return '';
    return str
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  }
})();
