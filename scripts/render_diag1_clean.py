import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

# Diagram 1: Clean, Minimalist, Visual-First Reverse Proxy NAT Traversal
svg_diag_01_clean = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1100 450" width="1100" height="450">
  <defs>
    <filter id="cardShadow" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <linearGradient id="hdrHome" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" /><stop offset="100%" stop-color="#4338ca" />
    </linearGradient>
    <linearGradient id="hdrISP" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#f43f5e" /><stop offset="100%" stop-color="#e11d48" />
    </linearGradient>
    <linearGradient id="hdrCloud" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" /><stop offset="100%" stop-color="#0369a1" />
    </linearGradient>
    <linearGradient id="hdrViewer" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#d97706" /><stop offset="100%" stop-color="#b45309" />
    </linearGradient>

    <!-- Markers -->
    <marker id="arrGreen" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="arrRed" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#ef4444" />
    </marker>
    <marker id="arrBlue" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="arrAmber" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#f59e0b" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .main-t { font-size: 18px; font-weight: 700; fill: #0f172a; }
    .main-s { font-size: 12px; fill: #64748b; }
    .card-t { font-size: 13px; font-weight: 700; fill: #ffffff; }
    .box-t { font-size: 12px; font-weight: 700; }
    .box-d { font-size: 11px; fill: #64748b; }
    .badge { font-size: 10px; font-weight: 700; }
    .step-t { font-size: 10.5px; font-weight: 700; }
  </style>

  <!-- Background Base -->
  <rect width="1100" height="450" rx="14" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />
  
  <!-- Header Title -->
  <text x="36" y="34" class="main-t">Reverse Proxy NAT Traversal &amp; Inbound Firewall Bypass</text>
  <text x="36" y="52" class="main-s">Outbound TLS 1.3 reverse channel eliminates residential port forwarding under ISP CGNAT (RFC 6598)</text>

  <!-- ==================== 1. HOME NODE ==================== -->
  <g transform="translate(36, 75)" filter="url(#cardShadow)">
    <rect width="210" height="330" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="210" height="36" rx="10" fill="url(#hdrHome)" />
    <rect y="26" width="210" height="10" fill="url(#hdrHome)" />
    <text x="16" y="23" class="card-t">🏠 Home Edge Node</text>

    <!-- Component 1 -->
    <rect x="12" y="50" width="186" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="22" y="72" class="box-t" fill="#3730a3">Edge Daemon (Go)</text>
    <text x="22" y="90" class="box-d">Managed static binary</text>
    <rect x="22" y="96" width="76" height="14" rx="3" fill="#ede9fe" />
    <text x="26" y="106" class="badge" fill="#5b21b6">SUPERVISOR</text>

    <!-- Component 2 -->
    <rect x="12" y="124" width="186" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="22" y="146" class="box-t" fill="#047857">Local Streamer (:8089)</text>
    <text x="22" y="164" class="box-d">4s fMP4 HLS Chunks</text>
    <rect x="22" y="170" width="70" height="14" rx="3" fill="#d1fae5" />
    <text x="26" y="180" class="badge" fill="#065f46">DATA PLANE</text>

    <!-- Component 3 -->
    <rect x="12" y="198" width="186" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="22" y="220" class="box-t" fill="#0f172a">Local NVMe &amp; SQLite</text>
    <text x="22" y="238" class="box-d">Private storage &amp; watch log</text>
    <rect x="22" y="244" width="72" height="14" rx="3" fill="#f1f5f9" />
    <text x="26" y="254" class="badge" fill="#475569">INVARIANT 1</text>

    <!-- Bottom Badge -->
    <rect x="12" y="274" width="186" height="42" rx="6" fill="#ecfdf5" stroke="#a7f3d0" />
    <text x="26" y="295" font-size="11" font-weight="700" fill="#065f46">🔒 0 Open Ports Needed</text>
    <text x="26" y="308" font-size="9.5" fill="#047857">Outbound Connections Only</text>
  </g>

  <!-- ==================== 2. ISP CGNAT & ROUTER ==================== -->
  <g transform="translate(286, 75)" filter="url(#cardShadow)">
    <rect width="180" height="330" rx="10" fill="#fff5f5" stroke="#fecaca" stroke-width="1.2" stroke-dasharray="6,4" />
    <rect width="180" height="36" rx="10" fill="url(#hdrISP)" />
    <rect y="26" width="180" height="10" fill="url(#hdrISP)" />
    <text x="16" y="23" class="card-t">🛡️ ISP CGNAT &amp; Router</text>

    <!-- Inbound Blocked -->
    <g transform="translate(12, 52)">
      <rect width="156" height="100" rx="6" fill="#ffffff" stroke="#fca5a5" />
      <circle cx="26" cy="26" r="12" fill="#fee2e2" />
      <text x="21" y="31" font-size="13" fill="#dc2626" font-weight="900">✕</text>
      <text x="46" y="24" font-size="11.5" font-weight="700" fill="#991b1b">Inbound Drops</text>
      <text x="46" y="38" class="box-d">TCP SYN blocked</text>
      <line x1="10" y1="52" x2="146" y2="52" stroke="#f1f5f9" />
      <text x="12" y="70" font-size="10" fill="#64748b">• No Public IPv4</text>
      <text x="12" y="86" font-size="10" fill="#64748b">• CGNAT: 100.64.0.0/10</text>
    </g>

    <!-- Outbound Allowed -->
    <g transform="translate(12, 166)">
      <rect width="156" height="100" rx="6" fill="#ffffff" stroke="#86efac" />
      <circle cx="26" cy="26" r="12" fill="#dcfce7" />
      <text x="20" y="31" font-size="13" fill="#16a34a" font-weight="900">✓</text>
      <text x="46" y="24" font-size="11.5" font-weight="700" fill="#166534">Outbound Pass</text>
      <text x="46" y="38" class="box-d">Stateful channel</text>
      <line x1="10" y1="52" x2="146" y2="52" stroke="#f1f5f9" />
      <text x="12" y="70" font-size="10" fill="#047857">• TLS 1.3 Outbound</text>
      <text x="12" y="86" font-size="10" fill="#047857">• Reverse Stream Active</text>
    </g>

    <rect x="12" y="278" width="156" height="38" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="296" font-size="10.5" font-weight="600" fill="#475569">Stateful NAT Mapping</text>
    <text x="20" y="309" font-size="9" fill="#94a3b8">Preserves Invariant 2</text>
  </g>

  <!-- ==================== 3. CLOUD RELAY & ORIGIN SHIELD ==================== -->
  <g transform="translate(506, 75)" filter="url(#cardShadow)">
    <rect width="268" height="330" rx="10" fill="#f0f9ff" stroke="#bae6fd" stroke-width="1.2" />
    <rect width="268" height="36" rx="10" fill="url(#hdrCloud)" />
    <rect y="26" width="268" height="10" fill="url(#hdrCloud)" />
    <text x="16" y="23" class="card-t">☁️ Public Cloud Control &amp; Relay</text>

    <!-- FRP Relay Server -->
    <rect x="12" y="50" width="244" height="66" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="22" y="72" class="box-t" fill="#0369a1">FRP Relay Server (:7000)</text>
    <text x="22" y="90" class="box-d">Terminates reverse TLS connection</text>
    <rect x="22" y="96" width="86" height="14" rx="3" fill="#e0f2fe" />
    <text x="26" y="106" class="badge" fill="#0369a1">REVERSE TUNNEL</text>

    <!-- Cloud API Router -->
    <rect x="12" y="124" width="244" height="66" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="22" y="146" class="box-t" fill="#0f172a">Cloud API Router (Go)</text>
    <text x="22" y="164" class="box-d">Evaluates subnet &amp; stream cap (max 10)</text>
    <rect x="22" y="170" width="104" height="14" rx="3" fill="#f1f5f9" />
    <text x="26" y="180" class="badge" fill="#334155">STREAM DECISION</text>

    <!-- Origin Shield Cache -->
    <rect x="12" y="198" width="244" height="66" rx="6" fill="#ffffff" stroke="#fde047" />
    <text x="22" y="220" class="box-t" fill="#b45309">Origin Shield Cache (S3)</text>
    <text x="22" y="238" class="box-d">Single-ingest hot chunks (24h LRU)</text>
    <rect x="22" y="244" width="82" height="14" rx="3" fill="#fef3c7" />
    <text x="26" y="254" class="badge" fill="#92400e">ORIGIN SHIELD</text>

    <!-- Zero Storage Reminder -->
    <rect x="12" y="274" width="244" height="42" rx="6" fill="#ffffff" stroke="#cbd5e1" />
    <text x="22" y="294" font-size="10.5" font-weight="700" fill="#0369a1">Invariant 1: Ephemeral Cache Only</text>
    <text x="22" y="307" font-size="9.5" fill="#64748b">No permanent video stored in cloud</text>
  </g>

  <!-- ==================== 4. REMOTE VIEWERS ==================== -->
  <g transform="translate(814, 75)" filter="url(#cardShadow)">
    <rect width="250" height="330" rx="10" fill="#fffbeb" stroke="#fde68a" stroke-width="1.2" />
    <rect width="250" height="36" rx="10" fill="url(#hdrViewer)" />
    <rect y="26" width="250" height="10" fill="url(#hdrViewer)" />
    <text x="16" y="23" class="card-t">👥 Remote Audience</text>

    <!-- Viewer 1 (Cache Miss) -->
    <rect x="12" y="50" width="226" height="100" rx="6" fill="#ffffff" stroke="#fcd34d" />
    <rect x="20" y="60" width="62" height="16" rx="3" fill="#fef3c7" />
    <text x="24" y="72" class="badge" fill="#92400e">VIEWER 1</text>
    <text x="88" y="73" class="box-t" fill="#92400e">First Pull (Miss)</text>
    <text x="20" y="94" font-size="10.5" fill="#475569">1. Requests segment from Cloud</text>
    <text x="20" y="110" font-size="10.5" fill="#475569">2. Tunnel streams chunk from Home</text>
    <text x="20" y="128" font-size="10.5" font-weight="700" fill="#15803d">✓ Ingested to Origin Shield</text>

    <!-- Viewers 2..10 (Cache Hit) -->
    <rect x="12" y="158" width="226" height="100" rx="6" fill="#ffffff" stroke="#86efac" />
    <rect x="20" y="168" width="82" height="16" rx="3" fill="#dcfce7" />
    <text x="24" y="180" class="badge" fill="#166534">VIEWERS 2..10</text>
    <text x="108" y="181" class="box-t" fill="#166534">Shield Hit</text>
    <text x="20" y="202" font-size="10.5" fill="#475569">3. Served directly from S3 cache</text>
    <text x="20" y="220" font-size="10.5" font-weight="700" fill="#16a34a">⚡ 0 Home Upload Used</text>
    <text x="20" y="238" font-size="10.5" font-weight="700" fill="#047857">⚡ 90% Upload Bandwidth Saved</text>

    <!-- Web UI -->
    <rect x="12" y="266" width="226" height="50" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="286" font-size="10.5" font-weight="700" fill="#0f172a">Next.js Web Player (HLS.js)</text>
    <text x="20" y="302" font-size="9.5" fill="#64748b">Zero client installation needed</text>
  </g>

  <!-- ==================== FLOW PATHS & LABELS ==================== -->
  <!-- 1. Outbound Tunnel Establishment (Green) -->
  <path d="M 246 360 L 506 360" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrGreen)" />
  <rect x="296" y="348" width="160" height="20" rx="4" fill="#ecfdf5" stroke="#10b981" stroke-width="1" />
  <text x="306" y="362" class="step-t" fill="#065f46">① Outbound TLS Tunnel</text>

  <!-- 2. Blocked Inbound Probe (Red) -->
  <path d="M 506 102 L 472 102" fill="none" stroke="#ef4444" stroke-width="2.2" stroke-dasharray="4,4" marker-end="url(#arrRed)" />
  <rect x="424" y="91" width="82" height="18" rx="3" fill="#fee2e2" stroke="#fca5a5" stroke-width="1" />
  <text x="430" y="103" font-size="9.5" font-weight="700" fill="#991b1b">SYN Dropped</text>

  <!-- 3. Viewer 1 Pull (Blue) -->
  <path d="M 814 105 L 778 105" fill="none" stroke="#0284c7" stroke-width="2.2" marker-end="url(#arrBlue)" />
  <rect x="760" y="94" width="50" height="18" rx="3" fill="#e0f2fe" />
  <text x="766" y="106" font-size="9.5" font-weight="700" fill="#0369a1">② Get</text>

  <!-- 4. Viewer 2..10 Hit (Amber) -->
  <path d="M 778 225 L 814 225" fill="none" stroke="#f59e0b" stroke-width="2.2" marker-end="url(#arrAmber)" />
  <rect x="752" y="214" width="60" height="18" rx="3" fill="#fef3c7" />
  <text x="758" y="226" font-size="9.5" font-weight="700" fill="#b45309">③ Cached</text>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_01_nat_traversal.svg")
png_file = os.path.join(fig_dir, "diag_01_nat_traversal.png")
html_file = os.path.join(fig_dir, "diag_01_nat_traversal.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_01_clean)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1100px; height: 450px; }}
</style>
</head>
<body>
  {svg_diag_01_clean}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1100,450",
    html_file
]
res = subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

