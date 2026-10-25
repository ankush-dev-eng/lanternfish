/**
 * LANTERNFISH — Interactive HTML5 Canvas Network Topology & Intruder Trail Renderer
 * Oceanic Field-Instrument Design Palette:
 * Abyss: #0b1517 | Kelp: #1f3b36 | Foam: #e9e4d4 | Amber: #e0a33a | Rust: #c2553a | Mist: #8fa8a0
 */

class NetworkCanvas {
  constructor(canvasId, containerId) {
    this.canvas = document.getElementById(canvasId);
    this.container = document.getElementById(containerId);
    this.ctx = this.canvas.getContext('2d');

    this.nodes = [];
    this.links = [];
    this.subnets = [];
    this.path = null;
    this.activeAttackerNode = null;
    this.patientZeroNode = null;
    this.revealDecoys = false;
    this.triggeredDecoyNodeIDs = new Set();

    // Pan & Zoom
    this.scale = 1.0;
    this.panX = 40;
    this.panY = 20;
    this.isDragging = false;
    this.dragStartX = 0;
    this.dragStartY = 0;

    // Hover & Interaction
    this.hoveredNode = null;
    this.tooltipEl = document.getElementById('nodeTooltip');
    this.onNodeClickCallback = null;

    // Animation loop timer
    this.animTime = 0;
    this.trailPackets = [];

    this.initEvents();
    this.resize();
    this.animate();
  }

  setData(networkData) {
    if (!networkData) return;
    this.nodes = networkData.nodes || [];
    this.links = networkData.links || [];
    this.subnets = networkData.subnets || [];
    this.centerGraph();
  }

  setPath(attackerPath) {
    this.path = attackerPath;
    if (attackerPath && attackerPath.patient_zero_id) {
      this.patientZeroNode = attackerPath.patient_zero_id;
    }
  }

  setActiveAttacker(nodeId) {
    this.activeAttackerNode = nodeId;
  }

  setTriggeredDecoys(decoyList) {
    this.triggeredDecoyNodeIDs.clear();
    if (decoyList) {
      for (const d of decoyList) {
        if (d.triggered && d.node_id) {
          this.triggeredDecoyNodeIDs.add(d.node_id);
        }
      }
    }
  }

  setRevealDecoys(reveal) {
    this.revealDecoys = reveal;
  }

  centerGraph() {
    if (!this.nodes.length) return;
    let minX = Infinity, maxX = -Infinity, minY = Infinity, maxY = -Infinity;
    for (const n of this.nodes) {
      if (n.position) {
        if (n.position.x < minX) minX = n.position.x;
        if (n.position.x > maxX) maxX = n.position.x;
        if (n.position.y < minY) minY = n.position.y;
        if (n.position.y > maxY) maxY = n.position.y;
      }
    }

    const rect = this.canvas.getBoundingClientRect();
    const graphWidth = maxX - minX + 220;
    const graphHeight = maxY - minY + 220;

    const scaleX = rect.width / graphWidth;
    const scaleY = rect.height / graphHeight;
    this.scale = Math.min(Math.max(Math.min(scaleX, scaleY) * 0.92, 0.55), 1.3);

    const centerX = (minX + maxX) / 2;
    const centerY = (minY + maxY) / 2;

    this.panX = rect.width / 2 - centerX * this.scale;
    this.panY = rect.height / 2 - centerY * this.scale;
  }

  initEvents() {
    window.addEventListener('resize', () => this.resize());

    // Mouse drag pan
    this.canvas.addEventListener('mousedown', (e) => {
      this.isDragging = true;
      this.canvas.style.cursor = 'grabbing';
      this.dragStartX = e.clientX - this.panX;
      this.dragStartY = e.clientY - this.panY;
    });

    window.addEventListener('mousemove', (e) => {
      if (this.isDragging) {
        this.panX = e.clientX - this.dragStartX;
        this.panY = e.clientY - this.dragStartY;
      } else {
        this.handleMouseMove(e);
      }
    });

    window.addEventListener('mouseup', () => {
      this.isDragging = false;
      this.canvas.style.cursor = this.hoveredNode ? 'pointer' : 'grab';
    });

    // Zoom
    this.canvas.addEventListener('wheel', (e) => {
      e.preventDefault();
      const zoomFactor = e.deltaY < 0 ? 1.1 : 0.9;
      const newScale = Math.min(Math.max(this.scale * zoomFactor, 0.4), 2.5);

      const rect = this.canvas.getBoundingClientRect();
      const mouseX = e.clientX - rect.left;
      const mouseY = e.clientY - rect.top;

      this.panX = mouseX - (mouseX - this.panX) * (newScale / this.scale);
      this.panY = mouseY - (mouseY - this.panY) * (newScale / this.scale);
      this.scale = newScale;
    }, { passive: false });

    // Click on node
    this.canvas.addEventListener('click', (e) => {
      const node = this.getNodeAtMouse(e);
      if (node && this.onNodeClickCallback) {
        this.onNodeClickCallback(node);
      }
    });
  }

  handleMouseMove(e) {
    const node = this.getNodeAtMouse(e);
    this.hoveredNode = node;
    this.canvas.style.cursor = node ? 'pointer' : 'grab';

    if (node && this.tooltipEl) {
      const isRevealedDecoy = this.revealDecoys && node.is_decoy;
      const isTripped = this.triggeredDecoyNodeIDs.has(node.id);

      let decoyBadge = '';
      if (isRevealedDecoy || isTripped) {
        decoyBadge = `<div style="color: #e0a33a; font-weight: bold; margin-top: 4px;">🎣 DECOY LURE HOST (${node.tripwire_id || 'TRIPWIRE'})</div>`;
      }

      this.tooltipEl.innerHTML = `
        <div style="font-weight: 700; color: #e9e4d4; font-size: 12px;">${node.hostname}</div>
        <div style="color: #8fa8a0; font-size: 10px; margin-bottom: 4px;">${node.ip} &bull; ${node.role}</div>
        <div style="color: #cfc9b8; font-size: 10px;">${node.description || ''}</div>
        <div style="color: #5c8d89; font-size: 9.5px; margin-top: 4px;">OS: ${node.os}</div>
        ${decoyBadge}
      `;

      const rect = this.canvas.getBoundingClientRect();
      this.tooltipEl.style.left = `${e.clientX - rect.left + 14}px`;
      this.tooltipEl.style.top = `${e.clientY - rect.top + 14}px`;
      this.tooltipEl.style.display = 'block';
    } else if (this.tooltipEl) {
      this.tooltipEl.style.display = 'none';
    }
  }

  getNodeAtMouse(e) {
    const rect = this.canvas.getBoundingClientRect();
    const mouseX = (e.clientX - rect.left - this.panX) / this.scale;
    const mouseY = (e.clientY - rect.top - this.panY) / this.scale;

    for (const node of this.nodes) {
      if (!node.position) continue;
      const dx = mouseX - node.position.x;
      const dy = mouseY - node.position.y;
      if (Math.sqrt(dx * dx + dy * dy) <= 24) {
        return node;
      }
    }
    return null;
  }

  resize() {
    const dpr = window.devicePixelRatio || 1;
    const rect = this.container.getBoundingClientRect();
    this.canvas.width = rect.width * dpr;
    this.canvas.height = rect.height * dpr;
    this.canvas.style.width = `${rect.width}px`;
    this.canvas.style.height = `${rect.height}px`;
    this.ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  }

  animate() {
    this.animTime += 0.02;
    this.render();
    requestAnimationFrame(() => this.animate());
  }

  render() {
    const ctx = this.ctx;
    const rect = this.canvas.getBoundingClientRect();
    ctx.clearRect(0, 0, rect.width, rect.height);

    ctx.save();
    ctx.translate(this.panX, this.panY);
    ctx.scale(this.scale, this.scale);

    // 1. Draw subtle bathymetric depth grid
    this.drawDepthRings(ctx);

    // 2. Draw normal topology links
    this.drawLinks(ctx);

    // 3. Draw animated attacker traversal trails
    this.drawAttackerTrails(ctx);

    // 4. Draw network nodes
    this.drawNodes(ctx);

    ctx.restore();
  }

  drawDepthRings(ctx) {
    ctx.save();
    ctx.strokeStyle = 'rgba(31, 59, 54, 0.45)';
    ctx.lineWidth = 1;
    ctx.setLineDash([4, 6]);

    // Concentric depth rings
    const centerX = 520;
    const centerY = 360;
    for (let r = 180; r <= 680; r += 160) {
      ctx.beginPath();
      ctx.arc(centerX, centerY, r, 0, Math.PI * 2);
      ctx.stroke();
    }

    // Bathymetric depth label
    ctx.fillStyle = 'rgba(143, 168, 160, 0.25)';
    ctx.font = '9px "IBM Plex Mono", monospace';
    ctx.fillText("ZONE 0: PERIMETER DMZ", 60, 200);
    ctx.fillText("ZONE 1: INTERNAL APP TIER", 380, 100);
    ctx.fillText("ZONE 2: RESTRICTED CORE DB", 740, 100);
    ctx.restore();
  }

  drawLinks(ctx) {
    const nodeMap = new Map();
    for (const n of this.nodes) {
      nodeMap.set(n.id, n);
    }

    for (const link of this.links) {
      const src = nodeMap.get(link.source);
      const dst = nodeMap.get(link.target);
      if (!src || !dst || !src.position || !dst.position) continue;

      ctx.save();
      ctx.beginPath();
      ctx.moveTo(src.position.x, src.position.y);
      ctx.lineTo(dst.position.x, dst.position.y);

      if (link.type === 'admin_vpn') {
        ctx.strokeStyle = 'rgba(92, 141, 137, 0.5)';
        ctx.setLineDash([3, 3]);
      } else if (link.type === 'replication') {
        ctx.strokeStyle = 'rgba(122, 92, 141, 0.45)';
        ctx.setLineDash([6, 4]);
      } else {
        ctx.strokeStyle = 'rgba(31, 59, 54, 0.85)';
      }
      ctx.lineWidth = 1.2;
      ctx.stroke();
      ctx.restore();
    }
  }

  drawAttackerTrails(ctx) {
    if (!this.path || !this.path.hops || this.path.hops.length === 0) return;

    const nodeMap = new Map();
    for (const n of this.nodes) {
      nodeMap.set(n.id, n);
    }

    const totalHops = this.path.hops.length;

    for (let i = 0; i < totalHops; i++) {
      const hop = this.path.hops[i];
      const src = nodeMap.get(hop.from_node);
      const dst = nodeMap.get(hop.to_node);
      if (!src || !dst || !src.position || !dst.position) continue;

      const isLatest = i === totalHops - 1;
      const alpha = 0.45 + (i / totalHops) * 0.55;

      ctx.save();
      ctx.beginPath();
      ctx.moveTo(src.position.x, src.position.y);
      ctx.lineTo(dst.position.x, dst.position.y);

      // Traversed edge color
      if (hop.decoy_tripped) {
        ctx.strokeStyle = `rgba(224, 163, 58, ${alpha})`;
        ctx.lineWidth = isLatest ? 3.5 : 2.5;
        // Amber glowing lure trail
        ctx.shadowColor = '#e0a33a';
        ctx.shadowBlur = 10;
      } else {
        ctx.strokeStyle = `rgba(194, 85, 58, ${alpha})`;
        ctx.lineWidth = isLatest ? 3.0 : 2.0;
        ctx.shadowColor = '#c2553a';
        ctx.shadowBlur = 6;
      }
      ctx.stroke();
      ctx.restore();

      // Traveling packet animation along recent edge
      if (isLatest) {
        const t = (this.animTime * 1.5) % 1.0;
        const px = src.position.x + (dst.position.x - src.position.x) * t;
        const py = src.position.y + (dst.position.y - src.position.y) * t;

        ctx.save();
        ctx.fillStyle = hop.decoy_tripped ? '#e0a33a' : '#c2553a';
        ctx.beginPath();
        ctx.arc(px, py, 3.5, 0, Math.PI * 2);
        ctx.fill();
        ctx.restore();
      }
    }
  }

  drawNodes(ctx) {
    for (const node of this.nodes) {
      if (!node.position) continue;
      const { x, y } = node.position;

      const isHovered = this.hoveredNode && this.hoveredNode.id === node.id;
      const isPatientZero = this.patientZeroNode === node.id;
      const isActiveAttacker = this.activeAttackerNode === node.id;
      const isTrippedDecoy = this.triggeredDecoyNodeIDs.has(node.id);
      const isRevealedDecoy = this.revealDecoys && node.is_decoy;

      ctx.save();

      // 1. Decoy Amber Lure Pulse Halo (if tripped or revealed)
      if (isTrippedDecoy || (isRevealedDecoy && isTrippedDecoy)) {
        const pulse = 18 + Math.sin(this.animTime * 4) * 6;
        ctx.beginPath();
        ctx.arc(x, y, pulse, 0, Math.PI * 2);
        ctx.fillStyle = 'rgba(224, 163, 58, 0.2)';
        ctx.fill();

        ctx.beginPath();
        ctx.arc(x, y, pulse + 6, 0, Math.PI * 2);
        ctx.strokeStyle = 'rgba(224, 163, 58, 0.4)';
        ctx.lineWidth = 1;
        ctx.stroke();
      } else if (isRevealedDecoy) {
        // Operator mode revealed lure: amber dotted halo
        ctx.beginPath();
        ctx.arc(x, y, 22, 0, Math.PI * 2);
        ctx.strokeStyle = 'rgba(224, 163, 58, 0.7)';
        ctx.setLineDash([3, 3]);
        ctx.lineWidth = 1.5;
        ctx.stroke();
      }

      // 2. Active Attacker Pulse Marker
      if (isActiveAttacker) {
        const pulse = 22 + Math.sin(this.animTime * 6) * 4;
        ctx.beginPath();
        ctx.arc(x, y, pulse, 0, Math.PI * 2);
        ctx.strokeStyle = 'rgba(194, 85, 58, 0.8)';
        ctx.lineWidth = 2;
        ctx.stroke();
      }

      // 3. Node Circle Body
      const radius = isHovered ? 16 : 14;
      ctx.beginPath();
      ctx.arc(x, y, radius, 0, Math.PI * 2);

      // Node background by role
      ctx.fillStyle = this.getNodeColor(node);
      ctx.fill();

      // Border
      ctx.strokeStyle = isHovered ? '#e9e4d4' : '#1f3b36';
      if (isRevealedDecoy || isTrippedDecoy) {
        ctx.strokeStyle = '#e0a33a';
      }
      ctx.lineWidth = isHovered ? 2 : 1.5;
      ctx.stroke();

      // 4. Role Icon glyph
      ctx.fillStyle = '#e9e4d4';
      ctx.font = '10px "IBM Plex Mono", monospace';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      const glyph = this.getNodeGlyph(node);
      ctx.fillText(glyph, x, y);

      // 5. Patient Zero Badge
      if (isPatientZero) {
        ctx.fillStyle = '#c2553a';
        ctx.font = 'bold 8.5px "IBM Plex Mono", monospace';
        ctx.fillText("P0", x, y - radius - 8);
      }

      // 6. Label under node
      ctx.fillStyle = isHovered ? '#e9e4d4' : '#8fa8a0';
      ctx.font = isHovered ? 'bold 10px "IBM Plex Mono", monospace' : '9px "IBM Plex Mono", monospace';
      ctx.fillText(node.hostname.replace('.corp', ''), x, y + radius + 12);

      ctx.restore();
    }
  }

  getNodeColor(node) {
    // If not revealed and is decoy, disguise as regular host!
    if (node.is_decoy && !this.revealDecoys && !this.triggeredDecoyNodeIDs.has(node.id)) {
      if (node.id.includes('db')) return '#1e2433';
      if (node.id.includes('finance')) return '#1e2d27';
    }

    if (this.triggeredDecoyNodeIDs.has(node.id)) {
      return '#3d2b0e'; // Amber burnt core
    }
    if (this.revealDecoys && node.is_decoy) {
      return '#332714'; // Amber lure tinted
    }

    switch (node.role) {
      case 'gateway': return '#1a3330';
      case 'web_server': return '#1b3834';
      case 'app_server': return '#1d3d37';
      case 'database': return '#282438';
      case 'domain_controller': return '#302636';
      case 'backup_vault': return '#1e2428';
      case 'file_server': return '#1f352e';
      case 'workstation': return '#1c2d38';
      case 'decoy_host': return '#332714';
      default: return '#1f3b36';
    }
  }

  getNodeGlyph(node) {
    if (this.triggeredDecoyNodeIDs.has(node.id)) {
      return 'TRP';
    }
    if (this.revealDecoys && node.is_decoy) {
      return 'LUR';
    }
    switch (node.role) {
      case 'gateway': return 'GW';
      case 'web_server': return 'WEB';
      case 'app_server': return 'APP';
      case 'database': return 'DB';
      case 'domain_controller': return 'DC';
      case 'backup_vault': return 'VLT';
      case 'file_server': return 'FS';
      case 'workstation': return 'WS';
      case 'decoy_host': return 'DH';
      default: return '::';
    }
  }
}
