import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_02_clean = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1160 550" width="1160" height="550">
  <defs>
    <filter id="shadowBox" x="-4%" y="-4%" width="108%" height="112%" filterUnits="userSpaceOnUse">
      <feDropShadow dx="0" dy="3" stdDeviation="5" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>

    <!-- Gradients -->
    <linearGradient id="gEdgeHdr" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" /><stop offset="100%" stop-color="#3730a3" />
    </linearGradient>
    <linearGradient id="gCloudHdr" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" /><stop offset="100%" stop-color="#0369a1" />
    </linearGradient>
    <linearGradient id="gLanHdr" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#059669" /><stop offset="100%" stop-color="#047857" />
    </linearGradient>
    <linearGradient id="gRemoteHdr" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#d97706" /><stop offset="100%" stop-color="#b45309" />
    </linearGradient>

    <!-- Markers -->
    <marker id="mGreen" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="mBlue" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="mIndigo" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#6366f1" />
    </marker>
    <marker id="mAmber" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#f59e0b" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 18px; font-weight: 700; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .tier-title { font-size: 13px; font-weight: 700; fill: #ffffff; }
    .card-title { font-size: 12px; font-weight: 700; }
    .card-desc { font-size: 10.5px; fill: #64748b; }
    .tag { font-size: 9.5px; font-weight: 700; }
    .arrow-label { font-size: 10px; font-weight: 700; }
  </style>

  <!-- Canvas Background -->
  <rect width="1160" height="550" rx="14" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Header -->
  <text x="36" y="34" class="title">CustardFin Two-Tier Distributed System Architecture</text>
  <text x="36" y="52" class="subtitle">Decoupled Media Data Plane (Home Edge) and Ephemeral Coordination Control Plane (Cloud Tier)</text>

  <!-- ========================================================================= -->
  <!-- TIER 1: HOME EDGE NETWORK (TOP) -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 72)" filter="url(#shadowBox)">
    <!-- Outer Container -->
    <rect width="1088" height="195" rx="10" fill="#f8fafc" stroke="#c7d2fe" stroke-width="1.2" />
    <rect width="1088" height="32" rx="10" fill="url(#gEdgeHdr)" />
    <rect y="22" width="1088" height="10" fill="url(#gEdgeHdr)" />
    <text x="16" y="21" class="tier-title">🏠 HOME EDGE TIER — Decentralized Media Data Plane (Go Static Daemon, &lt;20MB Idle RAM)</text>

    <!-- Node Card: Storage & Ingestion -->
    <rect x="16" y="44" width="225" height="135" rx="8" fill="#ffffff" stroke="#e2e8f0" />
    <text x="28" y="66" class="card-title" fill="#3730a3">1. Media Master &amp; SQLite</text>
    <text x="28" y="85" class="card-desc">• Local NVMe Video Storage</text>
    <text x="28" y="103" class="card-desc">• Pure-Go SQLite WAL (Private Log)</text>
    <text x="28" y="121" class="card-desc">• Video Never Leaves Unshared</text>
    <rect x="28" y="138" width="82" height="16" rx="3" fill="#ede9fe" />
    <text x="34" y="150" class="tag" fill="#5b21b6">INVARIANTS 1 &amp; 8</text>

    <!-- Node Card: Transcoder -->
    <rect x="257" y="44" width="235" height="135" rx="8" fill="#ffffff" stroke="#e2e8f0" />
    <text x="269" y="66" class="card-title" fill="#3730a3">2. Ingest &amp; Transcoder</text>
    <text x="269" y="85" class="card-desc">• Recursive Scanner &amp; SHA-256</text>
    <text x="269" y="103" class="card-desc">• NVENC GPU (320+ fps, 10.7x)</text>
    <text x="269" y="121" class="card-desc">• Pre-encoded 4s fMP4 HLS</text>
    <rect x="269" y="138" width="96" height="16" rx="3" fill="#e0e7ff" />
    <text x="275" y="150" class="tag" fill="#3730a3">GPU AUTO-FALLBACK</text>

    <!-- Node Card: Local Streamer -->
    <rect x="508" y="44" width="235" height="135" rx="8" fill="#ffffff" stroke="#86efac" />
    <text x="520" y="66" class="card-title" fill="#047857">3. Local Streamer (:8089)</text>
    <text x="520" y="85" class="card-desc">• Embedded Mini-Jellyfin Server</text>
    <text x="520" y="103" class="card-desc">• Serves LAN Clients Directly</text>
    <text x="520" y="121" class="card-desc">• Zero Cloud Relay Latency (&lt;4ms)</text>
    <rect x="520" y="138" width="86" height="16" rx="3" fill="#d1fae5" />
    <text x="526" y="150" class="tag" fill="#065f46">OFFLINE ENGINE</text>

    <!-- Arrow Home -> LAN Clients -->
    <path d="M 743 111 L 815 111" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#mGreen)" />
    <rect x="748" y="90" width="60" height="17" rx="3" fill="#dcfce7" />
    <text x="752" y="102" font-size="9" font-weight="700" fill="#047857">Direct Wi-Fi</text>

    <!-- LAN Clients Card -->
    <rect x="825" y="44" width="247" height="135" rx="8" fill="#f0fdf4" stroke="#86efac" />
    <rect x="825" y="44" width="247" height="28" rx="8" fill="url(#gLanHdr)" />
    <rect y="64" x="825" width="247" height="8" fill="url(#gLanHdr)" />
    <text x="837" y="63" font-size="11.5" font-weight="700" fill="#ffffff">📶 Direct LAN Viewers</text>
    <text x="837" y="93" font-size="11" font-weight="700" fill="#0f172a">Home Network (Same Subnet)</text>
    <text x="837" y="112" class="card-desc">• Stream directly from :8089</text>
    <text x="837" y="130" class="card-desc">• 0 bytes cloud egress / 0 Internet</text>
    <rect x="837" y="143" width="115" height="18" rx="4" fill="#dcfce7" />
    <text x="843" y="156" font-size="9.5" font-weight="700" fill="#15803d">⚡ SUB-4ms LATENCY</text>
  </g>

  <!-- ========================================================================= -->
  <!-- CHANNELS BETWEEN HOME AND CLOUD (VERTICAL LANES) -->
  <!-- ========================================================================= -->
  <!-- Channel 1: Outbound Reverse Tunnel -->
  <path d="M 160 267 L 160 338" fill="none" stroke="#6366f1" stroke-width="2.5" marker-end="url(#mIndigo)" />
  <rect x="85" y="293" width="150" height="20" rx="4" fill="#eef2ff" stroke="#c7d2fe" />
  <text x="93" y="307" class="arrow-label" fill="#4338ca">① Outbound TLS Tunnel</text>

  <!-- Channel 2: Manifest Sync -->
  <path d="M 374 267 L 374 338" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#mGreen)" />
  <rect x="296" y="293" width="156" height="20" rx="4" fill="#ecfdf5" stroke="#a7f3d0" />
  <text x="304" y="307" class="arrow-label" fill="#047857">② ~3 KB Manifest Sync</text>

  <!-- Channel 3: Heartbeat Ping -->
  <path d="M 625 267 L 625 338" fill="none" stroke="#f59e0b" stroke-width="2.5" marker-end="url(#mAmber)" />
  <rect x="555" y="293" width="140" height="20" rx="4" fill="#fffbeb" stroke="#fde68a" />
  <text x="563" y="307" class="arrow-label" fill="#b45309">③ 15s Heartbeat Ping</text>

  <!-- ========================================================================= -->
  <!-- TIER 2: PUBLIC CLOUD CONTROL PLANE (BOTTOM) -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 342)" filter="url(#shadowBox)">
    <!-- Outer Container -->
    <rect width="1088" height="190" rx="10" fill="#f0f9ff" stroke="#bae6fd" stroke-width="1.2" />
    <rect width="1088" height="32" rx="10" fill="url(#gCloudHdr)" />
    <rect y="22" width="1088" height="10" fill="url(#gCloudHdr)" />
    <text x="16" y="21" class="tier-title">☁️ CLOUD CONTROL TIER — Ephemeral Coordination Plane ($49.15/mo Lightweight VPS / Lab VM)</text>

    <!-- Cloud Box 1: FRP Relay & Gateway -->
    <rect x="16" y="44" width="225" height="130" rx="8" fill="#ffffff" stroke="#e2e8f0" />
    <text x="28" y="66" class="card-title" fill="#0369a1">1. FRP Server (:7000)</text>
    <text x="28" y="85" class="card-desc">• Terminates Outbound TLS 1.3</text>
    <text x="28" y="103" class="card-desc">• Reverse Stream Multiplexing</text>
    <text x="28" y="121" class="card-desc">• Outbound-Only Bridge (Inv 2)</text>
    <rect x="28" y="135" width="84" height="16" rx="3" fill="#e0f2fe" />
    <text x="34" y="147" class="tag" fill="#0369a1">REVERSE RELAY</text>

    <!-- Cloud Box 2: Cloud API Router -->
    <rect x="257" y="44" width="235" height="130" rx="8" fill="#ffffff" stroke="#e2e8f0" />
    <text x="269" y="66" class="card-title" fill="#0f172a">2. Cloud API (apid)</text>
    <text x="269" y="85" class="card-desc">• Dual-Path Stream Router</text>
    <text x="269" y="103" class="card-desc">• 10-Stream Admission Cap</text>
    <text x="269" y="121" class="card-desc">• HMAC Signed Chunk URLs</text>
    <rect x="269" y="135" width="92" height="16" rx="3" fill="#f1f5f9" />
    <text x="275" y="147" class="tag" fill="#334155">ROUTER &amp; AUTH</text>

    <!-- Cloud Box 3: Origin Shield & Redis -->
    <rect x="508" y="44" width="235" height="130" rx="8" fill="#ffffff" stroke="#fde047" />
    <text x="520" y="66" class="card-title" fill="#b45309">3. Origin Shield &amp; State</text>
    <text x="520" y="85" class="card-desc">• S3 Cache (Single Hot Pull)</text>
    <text x="520" y="103" class="card-desc">• 24h Auto-Purge (Zero Perm)</text>
    <text x="520" y="121" class="card-desc">• Redis Heartbeat TTL (15s)</text>
    <rect x="520" y="135" width="86" height="16" rx="3" fill="#fef3c7" />
    <text x="526" y="147" class="tag" fill="#92400e">90% UPLINK SAVED</text>

    <!-- Arrow Cloud -> Remote Clients -->
    <path d="M 743 109 L 815 109" fill="none" stroke="#0284c7" stroke-width="2.5" marker-end="url(#mBlue)" />
    <rect x="748" y="88" width="60" height="17" rx="3" fill="#e0f2fe" />
    <text x="752" y="100" font-size="9" font-weight="700" fill="#0369a1">Public Web</text>

    <!-- Remote Clients Card -->
    <rect x="825" y="44" width="247" height="130" rx="8" fill="#fffbeb" stroke="#fde68a" />
    <rect x="825" y="44" width="247" height="28" rx="8" fill="url(#gRemoteHdr)" />
    <rect y="64" x="825" width="247" height="8" fill="url(#gRemoteHdr)" />
    <text x="837" y="63" font-size="11.5" font-weight="700" fill="#ffffff">🌍 Remote Public Viewers</text>
    <text x="837" y="93" font-size="11" font-weight="700" fill="#0f172a">Next.js 15 Web Player</text>
    <text x="837" y="112" class="card-desc">• HLS.js adaptive bitrate player</text>
    <text x="837" y="130" class="card-desc">• Served from S3 Origin Shield</text>
    <rect x="837" y="143" width="115" height="18" rx="4" fill="#fef3c7" />
    <text x="843" y="156" font-size="9.5" font-weight="700" fill="#b45309">ZERO CLIENT INSTALL</text>
  </g>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_02_architecture_overview.svg")
png_file = os.path.join(fig_dir, "diag_02_architecture_overview.png")
html_file = os.path.join(fig_dir, "diag_02_architecture_overview.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_02_clean)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1160px; height: 550px; }}
</style>
</head>
<body>
  {svg_diag_02_clean}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1160,550",
    html_file
]
subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

