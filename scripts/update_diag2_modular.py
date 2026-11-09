import os
import subprocess
import pathlib

fig_dir = pathlib.Path(r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures")
fig_dir.mkdir(parents=True, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

# ==============================================================================
# HIGH-PRECISION MODULAR ARCHITECTURE FOR FIGURE 3.1
# ==============================================================================
# - Clean technical dot grid background matching Figures 5.1 & 6.1
# - Pure <g> components with exact pixel geometry (no symbol viewBox scaling bugs)
# - Perfectly aligned cards and zero overlapping lines
# - Exact 3 parallel conduits across the central ISP CGNAT firewall
# ==============================================================================

svg_content = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 640" width="1200" height="640" style="background-color: #f8fafc; font-family: system-ui, -apple-system, 'SF Pro Display', 'Inter', Roboto, Helvetica, sans-serif; font-variant-numeric: tabular-nums;">
  <defs>
    <!-- Technical Background Dot Pattern -->
    <pattern id="dotGrid" x="0" y="0" width="20" height="20" patternUnits="userSpaceOnUse">
      <circle cx="2" cy="2" r="1" fill="#e2e8f0" />
    </pattern>

    <!-- Modern Drop Shadows -->
    <filter id="shadow-card" x="-8%" y="-8%" width="116%" height="116%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
      <feDropShadow dx="0" dy="1" stdDeviation="2" flood-color="#0f172a" flood-opacity="0.04" />
    </filter>
    <filter id="shadow-badge" x="-10%" y="-10%" width="120%" height="120%">
      <feDropShadow dx="0" dy="2" stdDeviation="3" flood-color="#0f172a" flood-opacity="0.08" />
    </filter>

    <!-- Gradients -->
    <linearGradient id="grad-home-tier" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#ffffff" />
      <stop offset="100%" stop-color="#fdf4ff" />
    </linearGradient>
    <linearGradient id="grad-cloud-tier" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#ffffff" />
      <stop offset="100%" stop-color="#f0f9ff" />
    </linearGradient>
    <linearGradient id="grad-nvenc" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0%" stop-color="#10b981" />
      <stop offset="100%" stop-color="#047857" />
    </linearGradient>
    <linearGradient id="grad-shield" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0%" stop-color="#f59e0b" />
      <stop offset="100%" stop-color="#d97706" />
    </linearGradient>
    <linearGradient id="grad-firewall" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#ef4444" />
      <stop offset="100%" stop-color="#b91c1c" />
    </linearGradient>

    <!-- Arrow Markers -->
    <marker id="arrow-green" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#10b981" />
    </marker>
    <marker id="arrow-indigo" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#6366f1" />
    </marker>
    <marker id="arrow-amber" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#f59e0b" />
    </marker>
    <marker id="arrow-blue" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#0284c7" />
    </marker>
  </defs>

  <!-- Background Base & Technical Dot Grid -->
  <rect width="1200" height="640" fill="#f8fafc" />
  <rect width="1200" height="640" fill="url(#dotGrid)" />

  <!-- Global Header Bar -->
  <rect x="0" y="0" width="1200" height="70" fill="#ffffff" stroke="#e2e8f0" stroke-width="1" filter="url(#shadow-card)" />
  <text x="40" y="32" font-size="18" font-weight="900" fill="#0f172a" letter-spacing="-0.3">
    CustardFin Two-Tier Hybrid Cloud-Edge Architecture
  </text>
  <text x="40" y="52" font-size="11" font-weight="500" fill="#64748b">
    Master Video Binaries Stay at Home (Left) · Outbound-Only TLS 1.3 Reverse Tunnel (Center) · Coordination &amp; Origin Shield in Cloud (Right)
  </text>
  <rect x="1010" y="20" width="150" height="28" rx="14" fill="#ede9fe" stroke="#c084fc" />
  <text x="1085" y="38" font-size="10" font-weight="800" fill="#6b21a8" text-anchor="middle">MCA Cloud Project</text>

  <!-- ========================================================================= -->
  <!-- ZONE 1: HOME EDGE NETWORK (LAN) [LEFT TIER: X=30..510] -->
  <!-- ========================================================================= -->
  <g transform="translate(30, 85)">
    <!-- Zone Container -->
    <rect width="480" height="530" rx="16" fill="url(#grad-home-tier)" stroke="#d8b4fe" stroke-width="1.8" filter="url(#shadow-card)" />
    
    <!-- Zone Header Banner -->
    <rect x="20" y="14" width="220" height="28" rx="14" fill="#7c3aed" filter="url(#shadow-badge)" />
    <text x="130" y="32" font-size="11" font-weight="800" fill="#ffffff" text-anchor="middle">🏠 HOME EDGE TIER (LAN)</text>
    <text x="255" y="32" font-size="10" font-weight="700" fill="#6b21a8">Subnet: 192.168.1.0/24</text>
    <text x="460" y="32" font-size="9" font-weight="800" fill="#059669" text-anchor="end">● ONLINE</text>

    <!-- ROW 1: Storage -> GPU -> HLS Packager (Y=55) -->
    <!-- Card 1: NVMe Storage -->
    <g transform="translate(20, 55)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#c7d2fe" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="14" width="30" height="24" rx="4" fill="#e0e7ff" stroke="#6366f1" stroke-width="1.2" />
      <circle cx="21" cy="26" r="4" fill="#6366f1" />
      <circle cx="33" cy="26" r="1.5" fill="#4338ca" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#1e1b4b">NVMe Disk</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#4338ca">Master Videos</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#e0e7ff" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#ede9fe" />
      <text x="62.5" y="72" font-size="8.5" font-weight="800" fill="#6b21a8" text-anchor="middle">🔒 100% Home Disks</text>
    </g>

    <!-- Arrow 1: Storage -> GPU -->
    <path d="M 145 102.5 L 175 102.5" stroke="#6366f1" stroke-width="2.5" marker-end="url(#arrow-indigo)" fill="none" />

    <!-- Card 2: RTX 4060 GPU Transcoder -->
    <g transform="translate(175, 55)">
      <rect width="130" height="95" rx="10" fill="#ffffff" stroke="#a7f3d0" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="12" width="32" height="30" rx="4" fill="url(#grad-nvenc)" stroke="#065f46" stroke-width="1.2" />
      <rect x="17" y="17" width="22" height="20" rx="2" fill="#064e3b" />
      <text x="28" y="30" font-size="7.5" font-weight="900" fill="#34d399" text-anchor="middle">NVENC</text>
      <text x="50" y="24" font-size="10.5" font-weight="800" fill="#064e3b">RTX 4060 GPU</text>
      <text x="50" y="36" font-size="8.5" font-weight="600" fill="#047857">10.7× Real-Time</text>
      <line x1="12" y1="50" x2="118" y2="50" stroke="#d1fae5" stroke-width="1" />
      <rect x="12" y="58" width="106" height="22" rx="4" fill="#ecfdf5" />
      <text x="65" y="73" font-size="8.5" font-weight="800" fill="#065f46" text-anchor="middle">⚡ 322 fps | &lt;5% CPU</text>
    </g>

    <!-- Arrow 2: GPU -> HLS Packager -->
    <path d="M 305 102.5 L 335 102.5" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrow-green)" fill="none" />

    <!-- Card 3: HLS Packager -->
    <g transform="translate(335, 55)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#bae6fd" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="14" width="30" height="24" rx="4" fill="#0284c7" />
      <rect x="15" y="17" width="4" height="3" rx="1" fill="#ffffff" />
      <rect x="22" y="17" width="4" height="3" rx="1" fill="#ffffff" />
      <rect x="29" y="17" width="4" height="3" rx="1" fill="#ffffff" />
      <polygon points="24,24 24,32 31,28" fill="#38bdf8" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#0c4a6e">HLS Packager</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#0284c7">4s GOP (g=120)</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#e0f2fe" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#f0f9ff" />
      <text x="62.5" y="72" font-size="8.5" font-weight="700" fill="#0369a1" text-anchor="middle">1080p · 720p · 480p</text>
    </g>

    <!-- Subtitle explaining Row 1 -->
    <text x="240" y="172" font-size="9" font-weight="600" fill="#7c3aed" text-anchor="middle">
      Automated Local Ingestion: Scan → GPU Transcode (322 fps) → 4s HLS Segments
    </text>
    <line x1="20" y1="182" x2="460" y2="182" stroke="#f3e8ff" stroke-width="1.5" stroke-dasharray="4,4" />

    <!-- ROW 2: SQLite -> Edge Daemon -> LAN Viewers (Y=195) -->
    <!-- Card 4: Edge SQLite Database -->
    <g transform="translate(20, 195)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#ddd6fe" stroke-width="1.5" filter="url(#shadow-card)" />
      <ellipse cx="27" cy="22" rx="14" ry="5" fill="#8b5cf6" />
      <path d="M 13 22 L 13 36 A 14 5 0 0 0 41 36 L 41 22 Z" fill="#ede9fe" stroke="#8b5cf6" stroke-width="1" />
      <ellipse cx="27" cy="36" rx="14" ry="5" fill="#8b5cf6" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#3b0764">Edge SQLite</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#7c3aed">Pure-Go / WAL</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#ede9fe" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#faf5ff" />
      <text x="62.5" y="72" font-size="8.5" font-weight="700" fill="#6b21a8" text-anchor="middle">Private Watch History</text>
    </g>

    <!-- Arrow 3: SQLite -> Edge Daemon -->
    <path d="M 145 242.5 L 175 242.5" stroke="#8b5cf6" stroke-width="2.5" marker-end="url(#arrow-indigo)" fill="none" />

    <!-- Card 5: Edge Daemon (:8089) -->
    <g transform="translate(175, 195)">
      <rect width="130" height="95" rx="10" fill="#ffffff" stroke="#c7d2fe" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="14" width="28" height="24" rx="3" fill="#1e1b4b" />
      <rect x="15" y="18" width="14" height="2" rx="1" fill="#818cf8" />
      <rect x="15" y="24" width="14" height="2" rx="1" fill="#818cf8" />
      <circle cx="34" cy="20" r="1.5" fill="#10b981" />
      <circle cx="34" cy="26" r="1.5" fill="#38bdf8" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#1e1b4b">Edge Daemon</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#4338ca">Port :8089 (Go)</text>
      <line x1="12" y1="48" x2="118" y2="48" stroke="#e0e7ff" stroke-width="1" />
      <rect x="12" y="56" width="106" height="24" rx="4" fill="#f5f3ff" />
      <text x="65" y="72" font-size="8.5" font-weight="800" fill="#5b21b6" text-anchor="middle">Mini-Jellyfin Server</text>
    </g>

    <!-- Arrow 4: Daemon -> LAN Viewers -->
    <path d="M 305 242.5 L 335 242.5" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrow-green)" fill="none" />

    <!-- Card 6: Direct LAN Viewers -->
    <g transform="translate(335, 195)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#a7f3d0" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="15" width="28" height="20" rx="3" fill="#047857" />
      <polygon points="23,21 23,29 29,25" fill="#a7f3d0" />
      <path d="M 22 35 L 30 35 L 32 39 L 20 39 Z" fill="#065f46" />
      <text x="46" y="24" font-size="10.5" font-weight="800" fill="#064e3b">LAN Viewers</text>
      <text x="46" y="36" font-size="8.5" font-weight="600" fill="#047857">Home Wi-Fi</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#d1fae5" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#ecfdf5" />
      <text x="62.5" y="72" font-size="8.5" font-weight="800" fill="#065f46" text-anchor="middle">&lt;4ms | ₹0.00 Egress</text>
    </g>

    <!-- Flow Label for Local Streaming -->
    <text x="320" y="310" font-size="8.5" font-weight="800" fill="#047857" text-anchor="middle">
      Direct Local Playback (Gigabit Wi-Fi)
    </text>

    <!-- ROW 3: Edge Invariants Callout Card (Y=330) -->
    <rect x="20" y="330" width="440" height="175" rx="12" fill="#ffffff" stroke="#e9d5ff" stroke-width="1.2" />
    <text x="35" y="354" font-size="11" font-weight="800" fill="#4c1d95">Edge Invariant Guarantees</text>
    
    <!-- 4 Feature Badges -->
    <rect x="35" y="368" width="190" height="30" rx="6" fill="#fdf4ff" stroke="#f0abfc" />
    <text x="45" y="387" font-size="9" font-weight="700" fill="#701a75">✔ 100% Video on Home Disk</text>

    <rect x="235" y="368" width="210" height="30" rx="6" fill="#fdf4ff" stroke="#f0abfc" />
    <text x="245" y="387" font-size="9" font-weight="700" fill="#701a75">✔ Private Watch History in SQLite</text>

    <rect x="35" y="408" width="190" height="30" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="45" y="427" font-size="9" font-weight="700" fill="#14532d">✔ Auto GPU Probe / CPU Fallback</text>

    <rect x="235" y="408" width="210" height="30" rx="6" fill="#eff6ff" stroke="#93c5fd" />
    <text x="245" y="427" font-size="9" font-weight="700" fill="#1e3a8a">✔ Standalone Mini-Jellyfin Offline</text>

    <text x="35" y="465" font-size="9" font-weight="600" fill="#64748b">
      Daemon footprint: &lt;20 MB idle RAM · Zero open router ports · Pure-Go static binary
    </text>
    <text x="35" y="485" font-size="8.5" font-weight="500" fill="#94a3b8">
      Marking a video "Public" syncs a 650-byte manifest; video files are never copied to cloud.
    </text>
  </g>

  <!-- ========================================================================= -->
  <!-- CENTER BARRIER: ISP CGNAT & 3 PARALLEL CONDUITS [X=525..665] -->
  <!-- ========================================================================= -->
  <g transform="translate(525, 85)">
    <!-- Central Firewall Tower -->
    <rect x="5" y="0" width="130" height="530" rx="14" fill="#ffffff" stroke="#fecdd3" stroke-width="1.8" filter="url(#shadow-card)" />
    
    <!-- Firewall Shield Badge -->
    <rect x="17" y="16" width="106" height="28" rx="14" fill="url(#grad-firewall)" filter="url(#shadow-badge)" />
    <text x="70" y="34" font-size="11" font-weight="800" fill="#ffffff" text-anchor="middle">🛡️ ISP CGNAT</text>
    
    <!-- Inbound Block Badge -->
    <rect x="15" y="52" width="110" height="22" rx="4" fill="#fee2e2" />
    <text x="70" y="66" font-size="8" font-weight="800" fill="#991b1b" text-anchor="middle">❌ Inbound Ports Blocked</text>
    
    <text x="70" y="88" font-size="8" font-weight="600" fill="#e11d48" text-anchor="middle">RFC 6598 (Jio / Airtel)</text>
    <text x="70" y="100" font-size="7.5" fill="#64748b" text-anchor="middle">No Public IPv4 Address</text>

    <line x1="15" y1="114" x2="125" y2="114" stroke="#fecdd3" stroke-width="1" stroke-dasharray="3,3" />

    <!-- 3 PARALLEL CONDUIT BOXES -->
    <!-- Conduit 1: Reverse TLS Tunnel -->
    <rect x="13" y="130" width="114" height="65" rx="6" fill="#ecfdf5" stroke="#10b981" stroke-width="1.5" />
    <text x="70" y="150" font-size="9" font-weight="800" fill="#065f46" text-anchor="middle">① Outbound Tunnel</text>
    <text x="70" y="165" font-size="8" font-weight="600" fill="#047857" text-anchor="middle">TLS 1.3 / Port :7000</text>
    <text x="70" y="180" font-size="7.5" font-weight="700" fill="#10b981" text-anchor="middle">Persistent Yamux TCP</text>

    <!-- Conduit 2: Manifest Sync -->
    <rect x="13" y="235" width="114" height="65" rx="6" fill="#eef2ff" stroke="#6366f1" stroke-width="1.5" />
    <text x="70" y="255" font-size="9" font-weight="800" fill="#312e81" text-anchor="middle">② Manifest Sync</text>
    <text x="70" y="270" font-size="8" font-weight="600" fill="#4338ca" text-anchor="middle">~650 Byte JSON Only</text>
    <text x="70" y="285" font-size="7.5" font-weight="700" fill="#6366f1" text-anchor="middle">Zero Video Upload</text>

    <!-- Conduit 3: 15s Ephemeral Heartbeat -->
    <rect x="13" y="340" width="114" height="65" rx="6" fill="#fffbeb" stroke="#f59e0b" stroke-width="1.5" />
    <text x="70" y="360" font-size="9" font-weight="800" fill="#78350f" text-anchor="middle">③ Node Heartbeat</text>
    <text x="70" y="375" font-size="8" font-weight="600" fill="#b45309" text-anchor="middle">Every 15s Beacon</text>
    <text x="70" y="390" font-size="7.5" font-weight="700" fill="#d97706" text-anchor="middle">20s Redis TTL Liveness</text>

    <!-- Bottom firewall notice -->
    <line x1="15" y1="430" x2="125" y2="430" stroke="#fecdd3" stroke-width="1" stroke-dasharray="3,3" />
    <text x="70" y="455" font-size="8.5" font-weight="800" fill="#be123c" text-anchor="middle">Inbound SYN Dropped</text>
    <text x="70" y="472" font-size="7.5" fill="#64748b" text-anchor="middle">Outbound Only Allowed</text>
    <text x="70" y="488" font-size="7.5" font-weight="700" fill="#047857" text-anchor="middle">CGNAT Fully Bypassed</text>
  </g>

  <!-- Connectors: From Home Daemon (center at X=270, Y=327.5) to Conduits -->
  <!-- Line 1 to Conduit 1 (input at X=538, Y=247.5) -->
  <path d="M 335 327.5 C 420 327.5, 450 247.5, 538 247.5" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrow-green)" stroke-dasharray="6,3" />
  <!-- Line 2 to Conduit 2 (input at X=538, Y=352.5) -->
  <path d="M 335 327.5 C 420 327.5, 450 352.5, 538 352.5" fill="none" stroke="#6366f1" stroke-width="2.5" marker-end="url(#arrow-indigo)" />
  <!-- Line 3 to Conduit 3 (input at X=538, Y=457.5) -->
  <path d="M 335 327.5 C 420 327.5, 450 457.5, 538 457.5" fill="none" stroke="#f59e0b" stroke-width="2.5" marker-end="url(#arrow-amber)" />

  <!-- ========================================================================= -->
  <!-- ZONE 2: PUBLIC CLOUD CONTROL PLANE [RIGHT TIER: X=690..1170] -->
  <!-- ========================================================================= -->
  <g transform="translate(690, 85)">
    <!-- Zone Container -->
    <rect width="480" height="530" rx="16" fill="url(#grad-cloud-tier)" stroke="#7dd3fc" stroke-width="1.8" filter="url(#shadow-card)" />
    
    <!-- Zone Header Banner -->
    <rect x="20" y="14" width="250" height="28" rx="14" fill="#0284c7" filter="url(#shadow-badge)" />
    <text x="145" y="32" font-size="11" font-weight="800" fill="#ffffff" text-anchor="middle">☁️ CLOUD CONTROL PLANE</text>
    <text x="285" y="32" font-size="10" font-weight="700" fill="#0369a1">₹4,129/mo VPS Stack</text>
    <text x="460" y="32" font-size="8.5" font-weight="800" fill="#0284c7" text-anchor="end">GCP VM</text>

    <!-- ROW 1: FRP Server -> Stream Router -> Origin Shield (Y=55) -->
    <!-- Card 7: FRP Server Relay (:7000) -->
    <g transform="translate(20, 55)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#bae6fd" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="14" width="30" height="24" rx="4" fill="#0284c7" />
      <polygon points="27,18 19,34 35,34" fill="#bae6fd" />
      <circle cx="27" cy="24" r="2" fill="#0369a1" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#0c4a6e">FRP Server</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#0284c7">Port :7000 (TLS)</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#e0f2fe" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#f0f9ff" />
      <text x="62.5" y="72" font-size="8.5" font-weight="700" fill="#0369a1" text-anchor="middle">Tunnel Terminator</text>
    </g>

    <!-- Arrow 5: FRP Relay -> Stream Router -->
    <path d="M 145 102.5 L 175 102.5" stroke="#0284c7" stroke-width="2.5" marker-end="url(#arrow-blue)" fill="none" />

    <!-- Card 8: Go Stream Router (apid) -->
    <g transform="translate(175, 55)">
      <rect width="130" height="95" rx="10" fill="#ffffff" stroke="#93c5fd" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="14" width="32" height="24" rx="4" fill="#1d4ed8" />
      <path d="M 19 26 L 28 19 L 37 26 M 28 19 L 28 32" stroke="#ffffff" stroke-width="1.6" fill="none" />
      <text x="50" y="24" font-size="10.5" font-weight="800" fill="#1e3a8a">Stream Router</text>
      <text x="50" y="36" font-size="8.5" font-weight="600" fill="#2563eb">Cloud API (apid)</text>
      <line x1="12" y1="48" x2="118" y2="48" stroke="#dbeafe" stroke-width="1" />
      <rect x="12" y="56" width="106" height="24" rx="4" fill="#eff6ff" />
      <text x="65" y="72" font-size="8.5" font-weight="800" fill="#1e40af" text-anchor="middle">10-Stream Limiter</text>
    </g>

    <!-- Arrow 6: Stream Router -> Origin Shield -->
    <path d="M 305 102.5 L 335 102.5" stroke="#f59e0b" stroke-width="2.5" marker-end="url(#arrow-amber)" fill="none" />

    <!-- Card 9: Origin Shield S3 Cache -->
    <g transform="translate(335, 55)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#fde68a" stroke-width="1.5" filter="url(#shadow-card)" />
      <path d="M 27 15 L 39 20 L 39 31 C 39 37 27 42 27 42 C 27 42 15 37 15 31 L 15 20 Z" fill="url(#grad-shield)" />
      <path d="M 21 28 L 25 32 L 32 23" fill="none" stroke="#ffffff" stroke-width="1.6" />
      <text x="46" y="24" font-size="10.5" font-weight="800" fill="#78350f">Origin Shield</text>
      <text x="46" y="36" font-size="8.5" font-weight="600" fill="#d97706">S3 Object Cache</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#fef3c7" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#fffbeb" />
      <text x="62.5" y="72" font-size="8.5" font-weight="800" fill="#b45309" text-anchor="middle">24h Auto-Purge LRU</text>
    </g>

    <!-- Subtitle explaining Row 1 -->
    <text x="240" y="172" font-size="9" font-weight="600" fill="#0284c7" text-anchor="middle">
      Remote Video Route: Tunnel Pull (Miss) → Origin Shield (S3 Cache) → Signed URL (Hit)
    </text>
    <line x1="20" y1="182" x2="460" y2="182" stroke="#e0f2fe" stroke-width="1.5" stroke-dasharray="4,4" />

    <!-- ROW 2: Databases -> Next.js Web -> Remote Public Viewers (Y=195) -->
    <!-- Card 10: PostgreSQL + Redis Databases -->
    <g transform="translate(20, 195)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.5" filter="url(#shadow-card)" />
      <!-- Mini Postgres -->
      <ellipse cx="27" cy="22" rx="12" ry="4.5" fill="#0284c7" />
      <path d="M 15 22 L 15 34 A 12 4.5 0 0 0 39 34 L 39 22 Z" fill="#e0f2fe" stroke="#0284c7" stroke-width="1" />
      <text x="27" y="44" font-size="7" font-weight="800" fill="#0369a1" text-anchor="middle">Postgres</text>
      <!-- Mini Redis -->
      <ellipse cx="58" cy="22" rx="12" ry="4.5" fill="#ef4444" />
      <path d="M 46 22 L 46 34 A 12 4.5 0 0 0 70 34 L 70 22 Z" fill="#fee2e2" stroke="#ef4444" stroke-width="1" />
      <text x="58" y="44" font-size="7" font-weight="800" fill="#b91c1c" text-anchor="middle">Redis</text>
      <!-- Auth Tag -->
      <rect x="76" y="18" width="40" height="24" rx="3" fill="#0f172a" />
      <text x="96" y="30" font-size="7.5" font-weight="800" fill="#ffffff" text-anchor="middle">JWT</text>
      <text x="96" y="39" font-size="6.5" font-weight="700" fill="#94a3b8" text-anchor="middle">Argon2id</text>
      <line x1="12" y1="50" x2="113" y2="50" stroke="#f1f5f9" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#f8fafc" />
      <text x="62.5" y="72" font-size="8" font-weight="700" fill="#334155" text-anchor="middle">Catalog &amp; 15s Heartbeats</text>
    </g>

    <!-- Arrow 7: DBs -> Next.js Web -->
    <path d="M 145 242.5 L 175 242.5" stroke="#0284c7" stroke-width="2.5" marker-end="url(#arrow-blue)" fill="none" />

    <!-- Card 11: Next.js 15 Presentation Tier -->
    <g transform="translate(175, 195)">
      <rect width="130" height="95" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="14" width="30" height="24" rx="3" fill="#0f172a" />
      <circle cx="17" cy="20" r="1.5" fill="#ef4444" />
      <circle cx="22" cy="20" r="1.5" fill="#f59e0b" />
      <circle cx="27" cy="20" r="1.5" fill="#10b981" />
      <polygon points="23,25 23,32 29,28.5" fill="#ffffff" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#0f172a">Next.js 15 UI</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#475569">App Router + HLS</text>
      <line x1="12" y1="48" x2="118" y2="48" stroke="#f1f5f9" stroke-width="1" />
      <rect x="12" y="56" width="106" height="24" rx="4" fill="#f8fafc" />
      <text x="65" y="72" font-size="8.5" font-weight="700" fill="#0f172a" text-anchor="middle">Stats for Nerds HUD</text>
    </g>

    <!-- Arrow 8: Web -> Public Viewers -->
    <path d="M 305 242.5 L 335 242.5" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrow-green)" fill="none" />

    <!-- Card 12: Remote Public Viewers -->
    <g transform="translate(335, 195)">
      <rect width="125" height="95" rx="10" fill="#ffffff" stroke="#86efac" stroke-width="1.5" filter="url(#shadow-card)" />
      <rect x="12" y="15" width="24" height="16" rx="2" fill="#1e293b" />
      <rect x="10" y="31" width="28" height="2" fill="#64748b" />
      <rect x="30" y="20" width="12" height="18" rx="2" fill="#0f172a" stroke="#cbd5e1" stroke-width="0.8" />
      <text x="48" y="24" font-size="10.5" font-weight="800" fill="#064e3b">Public Viewers</text>
      <text x="48" y="36" font-size="8.5" font-weight="600" fill="#15803d">Mobile &amp; PC</text>
      <line x1="12" y1="48" x2="113" y2="48" stroke="#dcfce7" stroke-width="1" />
      <rect x="12" y="56" width="101" height="24" rx="4" fill="#f0fdf4" />
      <text x="62.5" y="72" font-size="8.5" font-weight="800" fill="#166534" text-anchor="middle">90% Upload Saved</text>
    </g>

    <!-- Flow Label for Remote Streaming -->
    <text x="320" y="310" font-size="8.5" font-weight="800" fill="#15803d" text-anchor="middle">
      Origin Shield Cache Fan-Out
    </text>

    <!-- ROW 3: Cloud Economic Invariants Callout Card (Y=330) -->
    <rect x="20" y="330" width="440" height="175" rx="12" fill="#ffffff" stroke="#bae6fd" stroke-width="1.2" />
    <text x="35" y="354" font-size="11" font-weight="800" fill="#075985">Cloud Invariant Guarantees</text>
    
    <!-- 4 Feature Badges -->
    <rect x="35" y="368" width="195" height="30" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="45" y="387" font-size="9" font-weight="700" fill="#14532d">✔ 0 GB Permanent Video in Cloud</text>

    <rect x="235" y="368" width="210" height="30" rx="6" fill="#f0f9ff" stroke="#7dd3fc" />
    <text x="245" y="387" font-size="9" font-weight="700" fill="#0369a1">✔ HMAC-SHA256 Signed URLs</text>

    <rect x="35" y="408" width="195" height="30" rx="6" fill="#fffbeb" stroke="#fde68a" />
    <text x="45" y="427" font-size="9" font-weight="700" fill="#78350f">✔ Single-Ingest Cache (90% Saved)</text>

    <rect x="235" y="408" width="210" height="30" rx="6" fill="#fef2f2" stroke="#fca5a5" />
    <text x="245" y="427" font-size="9" font-weight="700" fill="#991b1b">✔ Strict 10-Stream Limit (HTTP 429)</text>

    <text x="35" y="465" font-size="9" font-weight="600" fill="#64748b">
      Verified monthly operating cost: ₹4,128.60 ($49.15) · 95.89% cloud bill reduction
    </text>
    <text x="35" y="485" font-size="8.5" font-weight="500" fill="#94a3b8">
      Saves ₹34.65 Lakhs over 3 years for 100 users with full personal data sovereignty.
    </text>
  </g>

  <!-- Connectors: From Conduits (output at X=652) to Cloud Services -->
  <!-- Line 1: Conduit 1 (Reverse Tunnel, Y=247.5) -> FRP Server Relay (X=710, Y=187.5) -->
  <path d="M 652 247.5 C 675 247.5, 680 187.5, 710 187.5" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrow-green)" stroke-dasharray="6,3" />

  <!-- Line 2: Conduit 2 (Manifest Sync, Y=352.5) -> Databases/API (X=710, Y=327.5) -->
  <path d="M 652 352.5 C 675 352.5, 680 327.5, 710 327.5" fill="none" stroke="#6366f1" stroke-width="2.5" marker-end="url(#arrow-indigo)" />

  <!-- Line 3: Conduit 3 (Heartbeat, Y=457.5) -> Redis Liveness (X=710, Y=347.5) -->
  <path d="M 652 457.5 C 685 457.5, 685 347.5, 710 347.5" fill="none" stroke="#f59e0b" stroke-width="2.5" marker-end="url(#arrow-amber)" />

</svg>
"""

# Save SVG
svg_path = fig_dir / "diag_02_architecture_overview.svg"
with open(svg_path, "w", encoding="utf-8") as f:
    f.write(svg_content)
print(f"Saved clean SVG to {svg_path}")

# Save HTML wrapper for headless Chrome screenshot (3600 x 1920)
html_path = fig_dir / "diag_02_architecture_overview.html"
html_content = f"""<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    * {{ margin: 0; padding: 0; box-sizing: border-box; }}
    html, body {{ width: 3600px; height: 1920px; background: #f8fafc; overflow: hidden; }}
    svg {{ width: 3600px; height: 1920px; display: block; }}
  </style>
</head>
<body>
  {svg_content}
</body>
</html>
"""

with open(html_path, "w", encoding="utf-8") as f:
    f.write(html_content)
print(f"Saved HTML wrapper to {html_path}")

# Render to Ultra-HD PNG using exact Chrome flags
png_path = fig_dir / "diag_02_architecture_overview.png"
cmd = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    "--hide-scrollbars",
    "--window-size=3600,1920",
    f"--screenshot={png_path}",
    str(html_path)
]

print("Rendering Ultra-HD PNG via headless Chrome...")
subprocess.run(cmd, check=True)
print(f"Successfully generated {png_path} ({os.path.getsize(png_path)} bytes)")
