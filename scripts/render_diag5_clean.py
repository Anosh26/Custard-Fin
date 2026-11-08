import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_05_clean = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1160 480" width="1160" height="480">
  <defs>
    <filter id="cShadow" x="-4%" y="-4%" width="108%" height="112%" filterUnits="userSpaceOnUse">
      <feDropShadow dx="0" dy="3" stdDeviation="5" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>

    <!-- Gradients -->
    <linearGradient id="gReq" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" /><stop offset="100%" stop-color="#3730a3" />
    </linearGradient>
    <linearGradient id="gLan5" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#059669" /><stop offset="100%" stop-color="#047857" />
    </linearGradient>
    <linearGradient id="gRemote5" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" /><stop offset="100%" stop-color="#0369a1" />
    </linearGradient>

    <!-- Markers -->
    <marker id="arrG5" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="arrB5" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="arrGold5" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#d97706" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 18px; font-weight: 700; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .path-hdr { font-size: 12.5px; font-weight: 700; fill: #ffffff; }
    .card-t { font-size: 12px; font-weight: 700; fill: #0f172a; }
    .card-d { font-size: 10.5px; fill: #475569; }
    .tag { font-size: 9px; font-weight: 700; }
  </style>

  <!-- Background Canvas -->
  <rect width="1160" height="480" rx="14" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Header -->
  <text x="36" y="34" class="title">Dual-Path Stream Routing &amp; Origin Shield Single Ingest</text>
  <text x="36" y="52" class="subtitle">Dynamic proximity routing delivers sub-4ms LAN playback and shields the residential uplink from multi-viewer load</text>

  <!-- ========================================================================= -->
  <!-- TOP: CLIENT REQUEST & DECISION GATE -->
  <!-- ========================================================================= -->
  <g transform="translate(360, 72)" filter="url(#cShadow)">
    <rect width="440" height="66" rx="8" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="440" height="26" rx="8" fill="url(#gReq)" />
    <rect y="18" width="440" height="8" fill="url(#gReq)" />
    <text x="16" y="18" font-size="11.5" font-weight="700" fill="#ffffff">🧭 Cloud Stream Router Decision Gate (apid)</text>

    <text x="16" y="44" font-size="11" font-weight="700" fill="#0f172a">Check: Client IP == Edge LAN Subnet / Public NAT IP?</text>
    <text x="16" y="58" font-size="10" fill="#64748b">Also verifies 10-stream concurrency admission cap per node</text>
  </g>

  <!-- Connectors from Router to Paths -->
  <!-- Left Branch (LAN) -->
  <path d="M 430 138 L 430 170 L 260 170 L 260 195" fill="none" stroke="#10b981" stroke-width="2.2" marker-end="url(#arrG5)" />
  <rect x="290" y="158" width="105" height="18" rx="3" fill="#ecfdf5" stroke="#a7f3d0" />
  <text x="296" y="171" font-size="9" font-weight="700" fill="#065f46">MATCH: SAME LAN</text>

  <!-- Right Branch (Remote) -->
  <path d="M 730 138 L 730 170 L 890 170 L 890 195" fill="none" stroke="#0284c7" stroke-width="2.2" marker-end="url(#arrB5)" />
  <rect x="765" y="158" width="105" height="18" rx="3" fill="#e0f2fe" stroke="#bae6fd" />
  <text x="771" y="171" font-size="9" font-weight="700" fill="#0369a1">NO MATCH: REMOTE</text>

  <!-- ========================================================================= -->
  <!-- LEFT COLUMN: PATH A - DIRECT LAN STREAMING -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 202)" filter="url(#cShadow)">
    <rect width="520" height="250" rx="10" fill="#f0fdf4" stroke="#86efac" stroke-width="1.2" />
    <rect width="520" height="30" rx="10" fill="url(#gLan5)" />
    <rect y="20" width="520" height="10" fill="url(#gLan5)" />
    <text x="16" y="20" class="path-hdr">📶 PATH A: Direct LAN Stream (Local Subnet Proximity)</text>

    <!-- Box 1: Direct URL Return -->
    <rect x="14" y="42" width="492" height="65" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="24" y="62" class="card-t" fill="#047857">1. Direct Local Endpoint Returned</text>
    <text x="24" y="80" class="card-d">• API returns LAN URL: http://192.168.1.105:8089/stream/master.m3u8</text>
    <text x="24" y="96" class="card-d">• Client bypasses Cloud Relay, Internet, and CDN completely</text>

    <!-- Box 2: Performance Metrics -->
    <rect x="14" y="115" width="492" height="75" rx="6" fill="#ffffff" stroke="#86efac" />
    <text x="24" y="135" class="card-t" fill="#047857">2. Maximum Performance &amp; Zero Cost</text>
    <text x="24" y="153" class="card-d">• Playback Latency: &lt; 3.2 ms (Local Gigabit Ethernet / Wi-Fi)</text>
    <text x="24" y="169" class="card-d">• Cloud Egress: Exactly 0 Bytes ($0.00 / 100% Free)</text>
    <text x="24" y="185" class="card-d">• Absolute Privacy: Video traffic never leaves the home router</text>

    <!-- Result Banner -->
    <rect x="14" y="198" width="492" height="38" rx="6" fill="#ecfdf5" stroke="#a7f3d0" />
    <text x="24" y="217" font-size="10.5" font-weight="700" fill="#065f46">Offline Survival: Operates flawlessly even if ISP Internet is down</text>
    <text x="24" y="230" font-size="9" fill="#047857">Mini-Jellyfin standalone UI on :8089 delivers instant local playback</text>
  </g>

  <!-- ========================================================================= -->
  <!-- RIGHT COLUMN: PATH B - ORIGIN SHIELD CACHING -->
  <!-- ========================================================================= -->
  <g transform="translate(604, 202)" filter="url(#cShadow)">
    <rect width="520" height="250" rx="10" fill="#f0f9ff" stroke="#bae6fd" stroke-width="1.2" />
    <rect width="520" height="30" rx="10" fill="url(#gRemote5)" />
    <rect y="20" width="520" height="10" fill="url(#gRemote5)" />
    <text x="16" y="20" class="path-hdr">🛡️ PATH B: Origin Shield Single Ingest (Public Internet)</text>

    <!-- Box 1: First Pull (Miss) -->
    <rect x="14" y="42" width="240" height="110" rx="6" fill="#ffffff" stroke="#fcd34d" />
    <rect x="20" y="50" width="70" height="15" rx="3" fill="#fef3c7" />
    <text x="25" y="61" class="tag" fill="#92400e">CACHE MISS</text>
    <text x="96" y="62" class="card-t" fill="#92400e">Viewer 1 (Pull Once)</text>
    <text x="20" y="82" class="card-d">1. Pulls chunk via FRP Tunnel</text>
    <text x="20" y="98" class="card-d">2. Writes chunk to S3 Cache</text>
    <text x="20" y="114" class="card-d">3. 24h auto-purge LRU policy</text>
    <text x="20" y="132" font-size="9.5" font-weight="700" fill="#b45309">Uplink: Exactly 1x Chunk Size</text>

    <!-- Box 2: Subsequent Hits -->
    <rect x="266" y="42" width="240" height="110" rx="6" fill="#ffffff" stroke="#86efac" />
    <rect x="272" y="50" width="70" height="15" rx="3" fill="#dcfce7" />
    <text x="277" y="61" class="tag" fill="#166534">CACHE HIT</text>
    <text x="348" y="62" class="card-t" fill="#166534">Viewers 2–10</text>
    <text x="272" y="82" class="card-d">1. Served directly from S3</text>
    <text x="272" y="98" class="card-d">2. Response time: &lt; 14 ms</text>
    <text x="272" y="114" class="card-d">3. Edge uplink: 0 Bps used</text>
    <text x="272" y="132" font-size="9.5" font-weight="700" fill="#15803d">90% Upload Bandwidth Saved</text>

    <!-- Result Banner -->
    <rect x="14" y="160" width="492" height="76" rx="6" fill="#ffffff" stroke="#cbd5e1" />
    <text x="24" y="180" font-size="10.5" font-weight="700" fill="#0f172a">Origin Shield Guarantee (Single Home Upload per Hot Chunk)</text>
    <text x="24" y="198" class="card-d">• 10 simultaneous viewers stream 1080p smoothly over a standard 20 Mbps home uplink</text>
    <text x="24" y="214" class="card-d">• Temporary cache automatically self-cleans — Zero permanent video stored in cloud</text>
    <text x="24" y="228" font-size="9" font-weight="700" fill="#0369a1">Invariant 1 &amp; Invariant 7 Preserved</text>
  </g>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_05_dual_path_shield.svg")
png_file = os.path.join(fig_dir, "diag_05_dual_path_shield.png")
html_file = os.path.join(fig_dir, "diag_05_dual_path_shield.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_05_clean)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1160px; height: 480px; }}
</style>
</head>
<body>
  {svg_diag_05_clean}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1160,480",
    html_file
]
subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

