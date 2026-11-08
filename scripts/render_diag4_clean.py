import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_04_clean = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1160 440" width="1160" height="440">
  <defs>
    <filter id="boxShadow" x="-4%" y="-4%" width="108%" height="112%" filterUnits="userSpaceOnUse">
      <feDropShadow dx="0" dy="3" stdDeviation="5" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>

    <!-- Gradients -->
    <linearGradient id="gOnline" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#059669" /><stop offset="100%" stop-color="#047857" />
    </linearGradient>
    <linearGradient id="gOffline" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#e11d48" /><stop offset="100%" stop-color="#be123c" />
    </linearGradient>

    <!-- Markers -->
    <marker id="arrGreen4" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="arrRed4" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#ef4444" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 18px; font-weight: 700; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .panel-hdr { font-size: 12.5px; font-weight: 700; fill: #ffffff; }
    .step-t { font-size: 12px; font-weight: 700; fill: #0f172a; }
    .step-d { font-size: 10.5px; fill: #475569; }
    .badge { font-size: 9px; font-weight: 700; }
  </style>

  <!-- Background Base -->
  <rect width="1160" height="440" rx="14" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Title Header -->
  <text x="36" y="34" class="title">Distributed Ephemeral Heartbeat &amp; Zero-Polling Liveness State Machine</text>
  <text x="36" y="52" class="subtitle">Redis TTL keys automate distributed node health tracking without periodic database scans or background worker CPU</text>

  <!-- ========================================================================= -->
  <!-- PANEL 1: ONLINE / HEALTHY FLOW -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 72)" filter="url(#boxShadow)">
    <rect width="1088" height="160" rx="10" fill="#f0fdf4" stroke="#86efac" stroke-width="1.2" />
    <rect width="1088" height="30" rx="10" fill="url(#gOnline)" />
    <rect y="20" width="1088" height="10" fill="url(#gOnline)" />
    <text x="16" y="20" class="panel-hdr">🟢 STATE 1: Normal Operation — 15s Heartbeat Beacon &amp; Key Renewal</text>

    <!-- Box 1: Edge Ping -->
    <rect x="16" y="42" width="220" height="102" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="26" y="62" class="step-t" fill="#047857">1. Edge Daemon (:8089)</text>
    <text x="26" y="80" class="step-d">• 15-second background timer</text>
    <text x="26" y="96" class="step-d">• Sends active streams &amp; LAN IP</text>
    <rect x="26" y="106" width="92" height="16" rx="3" fill="#dcfce7" />
    <text x="31" y="118" class="badge" fill="#15803d">15s TIMER TICK</text>

    <!-- Arrow 1 -> 2 -->
    <path d="M 236 93 L 268 93" fill="none" stroke="#10b981" stroke-width="2.2" marker-end="url(#arrGreen4)" />

    <!-- Box 2: Cloud API -->
    <rect x="278" y="42" width="230" height="102" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="288" y="62" class="step-t" fill="#0f172a">2. Cloud API (apid)</text>
    <text x="288" y="80" class="step-d">• Authenticates node JWT</text>
    <text x="288" y="96" class="step-d">• Writes ephemeral key to Redis</text>
    <rect x="288" y="106" width="112" height="16" rx="3" fill="#e0f2fe" />
    <text x="293" y="118" class="badge" fill="#0369a1">POST /v1/heartbeat</text>

    <!-- Arrow 2 -> 3 -->
    <path d="M 508 93 L 540 93" fill="none" stroke="#10b981" stroke-width="2.2" marker-end="url(#arrGreen4)" />

    <!-- Box 3: Redis SETEX -->
    <rect x="550" y="42" width="230" height="102" rx="6" fill="#ffffff" stroke="#fde047" />
    <text x="560" y="62" class="step-t" fill="#b45309">3. Redis Ephemeral Store</text>
    <text x="560" y="80" class="step-d">• Key: heartbeat:{node_id}</text>
    <text x="560" y="96" class="step-d">• Exact TTL = 45 Seconds</text>
    <rect x="560" y="106" width="102" height="16" rx="3" fill="#fef3c7" />
    <text x="565" y="118" class="badge" fill="#92400e">SETEX (TTL: 45s)</text>

    <!-- Arrow 3 -> 4 -->
    <path d="M 780 93 L 812 93" fill="none" stroke="#10b981" stroke-width="2.2" marker-end="url(#arrGreen4)" />

    <!-- Box 4: Global Catalog -->
    <rect x="822" y="42" width="250" height="102" rx="6" fill="#ffffff" stroke="#86efac" />
    <text x="832" y="62" class="step-t" fill="#047857">4. Public Web Catalog</text>
    <text x="832" y="80" class="step-d">• Check EXISTS in &lt;1ms</text>
    <text x="832" y="96" class="step-d">• Full video discoverability active</text>
    <rect x="832" y="106" width="125" height="16" rx="3" fill="#d1fae5" />
    <text x="837" y="118" class="badge" fill="#065f46">STATUS: ONLINE (1ms)</text>
  </g>

  <!-- ========================================================================= -->
  <!-- PANEL 2: OFFLINE / DISCONNECT FLOW -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 252)" filter="url(#boxShadow)">
    <rect width="1088" height="160" rx="10" fill="#fff5f5" stroke="#fecaca" stroke-width="1.2" />
    <rect width="1088" height="30" rx="10" fill="url(#gOffline)" />
    <rect y="20" width="1088" height="10" fill="url(#gOffline)" />
    <text x="16" y="20" class="panel-hdr">🔴 STATE 2: Disconnect Event — Automatic Memory Eviction &amp; Catalog Protection</text>

    <!-- Box 1: Edge Severed -->
    <rect x="16" y="42" width="220" height="102" rx="6" fill="#ffffff" stroke="#fca5a5" />
    <text x="26" y="62" class="step-t" fill="#991b1b">1. Node Disconnection</text>
    <text x="26" y="80" class="step-d">• Power loss, sleep, or ISP drop</text>
    <text x="26" y="96" class="step-d">• Heartbeat beacon ceases</text>
    <rect x="26" y="106" width="94" height="16" rx="3" fill="#fee2e2" />
    <text x="31" y="118" class="badge" fill="#991b1b">0 UPLINK TRAFFIC</text>

    <!-- Arrow 1 -> 2 -->
    <path d="M 236 93 L 268 93" fill="none" stroke="#ef4444" stroke-width="2.2" marker-end="url(#arrRed4)" />

    <!-- Box 2: Grace Window -->
    <rect x="278" y="42" width="230" height="102" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="288" y="62" class="step-t" fill="#0f172a">2. Grace Window (30s)</text>
    <text x="288" y="80" class="step-d">• Absorbs residential network jitter</text>
    <text x="288" y="96" class="step-d">• Node allowed 2 missed pings</text>
    <rect x="288" y="106" width="108" height="16" rx="3" fill="#f1f5f9" />
    <text x="293" y="118" class="badge" fill="#475569">JITTER TOLERANCE</text>

    <!-- Arrow 2 -> 3 -->
    <path d="M 508 93 L 540 93" fill="none" stroke="#ef4444" stroke-width="2.2" marker-end="url(#arrRed4)" />

    <!-- Box 3: Key Expired -->
    <rect x="550" y="42" width="230" height="102" rx="6" fill="#ffffff" stroke="#fca5a5" />
    <text x="560" y="62" class="step-t" fill="#991b1b">3. Key Auto-Eviction</text>
    <text x="560" y="80" class="step-d">• Redis TTL reaches 0</text>
    <text x="560" y="96" class="step-d">• Key automatically purged</text>
    <rect x="560" y="106" width="112" height="16" rx="3" fill="#fee2e2" />
    <text x="565" y="118" class="badge" fill="#991b1b">0 CPU PURGE TASK</text>

    <!-- Arrow 3 -> 4 -->
    <path d="M 780 93 L 812 93" fill="none" stroke="#ef4444" stroke-width="2.2" marker-end="url(#arrRed4)" />

    <!-- Box 4: Safe Routing -->
    <rect x="822" y="42" width="250" height="102" rx="6" fill="#ffffff" stroke="#fca5a5" />
    <text x="832" y="62" class="step-t" fill="#991b1b">4. Playback Routing Guard</text>
    <text x="832" y="80" class="step-d">• Serve hot cached chunks from S3</text>
    <text x="832" y="96" class="step-d">• If uncached ➔ Friendly HTTP 503</text>
    <rect x="832" y="106" width="138" height="16" rx="3" fill="#fee2e2" />
    <text x="837" y="118" class="badge" fill="#991b1b">STATUS: OFFLINE (Safe)</text>
  </g>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_04_heartbeat_liveness.svg")
png_file = os.path.join(fig_dir, "diag_04_heartbeat_liveness.png")
html_file = os.path.join(fig_dir, "diag_04_heartbeat_liveness.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_04_clean)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1160px; height: 440px; }}
</style>
</head>
<body>
  {svg_diag_04_clean}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1160,440",
    html_file
]
subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

