import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_01_creative = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1180 500" width="1180" height="500">
  <defs>
    <filter id="cShadow1" x="-10%" y="-10%" width="120%" height="120%">
      <feDropShadow dx="0" dy="5" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.08" />
    </filter>
    <filter id="glowGreen" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="3" result="blur" />
      <feComposite in="SourceGraphic" in2="blur" operator="over" />
    </filter>

    <linearGradient id="gHomeBg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#faf5ff" /><stop offset="100%" stop-color="#f3e8ff" />
    </linearGradient>
    <linearGradient id="gCloudBg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#f0f9ff" /><stop offset="100%" stop-color="#e0f2fe" />
    </linearGradient>
    <linearGradient id="gWall" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#ef4444" /><stop offset="50%" stop-color="#dc2626" /><stop offset="100%" stop-color="#991b1b" />
    </linearGradient>

    <!-- Markers -->
    <marker id="mGreen1" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="mRed1" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#ef4444" />
    </marker>
    <marker id="mBlue1" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#0284c7" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 19px; font-weight: 800; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .label { font-size: 11px; font-weight: 700; fill: #0f172a; text-anchor: middle; }
    .sub { font-size: 9.5px; fill: #64748b; text-anchor: middle; }
    .flow-badge { font-size: 9.5px; font-weight: 800; }
  </style>

  <!-- Background Base -->
  <rect width="1180" height="500" rx="16" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Header -->
  <text x="40" y="36" class="title">Reverse Proxy NAT Traversal &amp; Inbound Firewall Bypass</text>
  <text x="40" y="54" class="subtitle">Home edge node establishes an outbound TLS 1.3 tunnel through ISP CGNAT, avoiding open residential router ports</text>

  <!-- ========================================================================= -->
  <!-- 1. HOME RESIDENTIAL SUBNET (LEFT) -->
  <!-- ========================================================================= -->
  <g transform="translate(40, 75)">
    <!-- Boundary Pill -->
    <rect width="320" height="390" rx="14" fill="url(#gHomeBg)" stroke="#d8b4fe" stroke-width="1.5" stroke-dasharray="6,4" />
    <rect x="20" y="14" width="200" height="24" rx="12" fill="#7c3aed" />
    <text x="120" y="30" font-size="10.5" font-weight="800" fill="#ffffff" text-anchor="middle">🏠 HOME EDGE (LAN)</text>

    <!-- Edge Server Host Chassis -->
    <g transform="translate(95, 60)" filter="url(#cShadow1)">
      <rect width="130" height="85" rx="8" fill="#1e293b" stroke="#475569" stroke-width="1.5" />
      <rect x="14" y="16" width="60" height="14" rx="2" fill="#0f172a" stroke="#64748b" />
      <rect x="14" y="36" width="60" height="14" rx="2" fill="#0f172a" stroke="#64748b" />
      <rect x="14" y="56" width="60" height="14" rx="2" fill="#0f172a" stroke="#64748b" />
      <circle cx="90" cy="24" r="3.5" fill="#10b981" />
      <circle cx="104" cy="24" r="3.5" fill="#38bdf8" />
      <text x="65" y="108" class="label">Home Server (edged)</text>
      <text x="65" y="122" class="sub">Local Streamer (:8089)</text>
    </g>

    <!-- Local Media Storage Cylinder -->
    <g transform="translate(30, 210)" filter="url(#cShadow1)">
      <path d="M 5 18 L 5 55 A 30 10 0 0 0 65 55 L 65 18 Z" fill="#ffffff" stroke="#6366f1" stroke-width="1.5" />
      <ellipse cx="35" cy="18" rx="30" ry="10" fill="#818cf8" stroke="#4f46e5" stroke-width="1.5" />
      <text x="35" y="80" class="label">NVMe Storage</text>
      <text x="35" y="94" class="sub">4s HLS Segments</text>
      <text x="35" y="105" font-size="8.5" font-weight="700" fill="#4f46e5" text-anchor="middle">Master at Home</text>
    </g>

    <!-- Local LAN Client Monitor Screen -->
    <g transform="translate(180, 210)" filter="url(#cShadow1)">
      <rect x="8" y="8" width="80" height="52" rx="4" fill="#1e293b" stroke="#64748b" stroke-width="1.2" />
      <rect x="12" y="12" width="72" height="44" rx="2" fill="#064e3b" />
      <polygon points="44,28 44,42 56,35" fill="#34d399" />
      <path d="M 42 60 L 54 60 L 58 72 L 38 72 Z" fill="#475569" />
      <text x="48" y="92" class="label">LAN Viewers</text>
      <text x="48" y="105" font-size="8.5" font-weight="700" fill="#047857" text-anchor="middle">Sub-4ms Wi-Fi</text>
    </g>

    <!-- Wi-Fi Path from Server to LAN Client -->
    <path d="M 160 145 L 160 185 L 220 185 L 220 205" fill="none" stroke="#10b981" stroke-width="2" marker-end="url(#mGreen1)" />
    <text x="195" y="178" font-size="8" font-weight="700" fill="#047857">Direct LAN</text>

    <!-- Port Forwarding Badge -->
    <rect x="25" y="340" width="270" height="32" rx="6" fill="#ecfdf5" stroke="#86efac" />
    <text x="160" y="360" font-size="10.5" font-weight="800" fill="#065f46" text-anchor="middle">🔒 ZERO INBOUND PORTS OPEN (CGNAT SAFE)</text>
  </g>

  <!-- ========================================================================= -->
  <!-- 2. ISP CGNAT & FIREWALL (CENTER BARRIER) -->
  <!-- ========================================================================= -->
  <g transform="translate(400, 75)">
    <!-- Firewall Wall Graphic -->
    <rect width="130" height="390" rx="12" fill="#fff1f2" stroke="#fca5a5" stroke-width="1.5" stroke-dasharray="6,4" />
    
    <!-- Shield Top -->
    <g transform="translate(35, 20)">
      <path d="M 10 10 Q 30 0 50 10 Q 50 50 30 70 Q 10 50 10 10 Z" fill="url(#gWall)" stroke="#7f1d1d" stroke-width="1.5" filter="url(#cShadow1)" />
      <text x="30" y="42" font-size="18" fill="#ffffff" font-weight="900" text-anchor="middle">✕</text>
    </g>

    <!-- Barrier Text -->
    <text x="65" y="115" font-size="11" font-weight="800" fill="#991b1b" text-anchor="middle">ISP CGNAT</text>
    <text x="65" y="130" font-size="9" font-weight="600" fill="#be123c" text-anchor="middle">RFC 6598 Subnet</text>
    <text x="65" y="142" font-size="8.5" fill="#e11d48" text-anchor="middle">100.64.0.0/10</text>

    <!-- Blocked Inbound Probe Illustration -->
    <g transform="translate(15, 160)">
      <rect width="100" height="70" rx="6" fill="#ffffff" stroke="#fca5a5" />
      <text x="50" y="24" font-size="9" font-weight="800" fill="#991b1b" text-anchor="middle">INBOUND BLOCKED</text>
      <text x="50" y="40" font-size="8.5" fill="#64748b" text-anchor="middle">TCP SYN Dropped</text>
      <text x="50" y="54" font-size="8" fill="#64748b" text-anchor="middle">No Public IPv4</text>
      <text x="50" y="64" font-size="7.5" font-weight="700" fill="#be123c" text-anchor="middle">UPnP Disabled</text>
    </g>

    <!-- Allowed Outbound Tunnel Illustration -->
    <g transform="translate(15, 250)">
      <rect width="100" height="75" rx="6" fill="#f0fdf4" stroke="#86efac" />
      <text x="50" y="24" font-size="9" font-weight="800" fill="#166534" text-anchor="middle">OUTBOUND PASS</text>
      <text x="50" y="40" font-size="8.5" fill="#047857" text-anchor="middle">Persistent TLS 1.3</text>
      <text x="50" y="54" font-size="8" fill="#047857" text-anchor="middle">Reverse Channel</text>
      <text x="50" y="66" font-size="7.5" font-weight="700" fill="#15803d" text-anchor="middle">Stateful NAT Active</text>
    </g>

    <text x="65" y="360" font-size="8.5" font-weight="700" fill="#64748b" text-anchor="middle">Traditional Port</text>
    <text x="65" y="372" font-size="8.5" font-weight="700" fill="#64748b" text-anchor="middle">Forwarding Fails</text>
  </g>

  <!-- ========================================================================= -->
  <!-- 3. PUBLIC CLOUD RELAY & ORIGIN SHIELD (RIGHT) -->
  <!-- ========================================================================= -->
  <g transform="translate(570, 75)">
    <!-- Actual Vector Cloud Background Shape -->
    <path d="M 60 70 
             A 45 45 0 0 1 140 45 
             A 55 55 0 0 1 240 40 
             A 50 50 0 0 1 340 60 
             A 45 45 0 0 1 400 115 
             A 55 55 0 0 1 410 220 
             A 60 60 0 0 1 380 330 
             A 50 50 0 0 1 300 380 
             A 60 60 0 0 1 150 385 
             A 55 55 0 0 1 50 330 
             A 45 45 0 0 1 25 220 
             A 45 45 0 0 1 30 130 
             A 45 45 0 0 1 60 70 Z" 
          fill="url(#gCloudBg)" stroke="#7dd3fc" stroke-width="1.8" filter="url(#cShadow1)" />

    <!-- Cloud Title Pill -->
    <rect x="90" y="14" width="240" height="24" rx="12" fill="#0284c7" />
    <text x="210" y="30" font-size="10.5" font-weight="800" fill="#ffffff" text-anchor="middle">☁️ PUBLIC CLOUD RELAY (GCP VM)</text>

    <!-- FRP Relay Server -->
    <g transform="translate(50, 85)" filter="url(#cShadow1)">
      <rect width="110" height="65" rx="6" fill="#ffffff" stroke="#0284c7" stroke-width="1.5" />
      <polygon points="55,10 38,32 72,32" fill="#38bdf8" />
      <circle cx="55" cy="40" r="4" fill="#0284c7" />
      <text x="55" y="85" class="label">FRP Server (:7000)</text>
      <text x="55" y="98" class="sub">Reverse TLS Relay</text>
    </g>

    <!-- Stream Router API -->
    <g transform="translate(180, 85)" filter="url(#cShadow1)">
      <rect width="105" height="65" rx="6" fill="#ffffff" stroke="#0284c7" stroke-width="1.5" />
      <rect x="15" y="16" width="75" height="32" rx="3" fill="#e0f2fe" />
      <text x="52.5" y="32" font-size="9" font-weight="800" fill="#0369a1" text-anchor="middle">STREAM ROUTER</text>
      <text x="52.5" y="44" font-size="8" fill="#0284c7" text-anchor="middle">10-Cap Limiter</text>
      <text x="52.5" y="85" class="label">Cloud API (apid)</text>
      <text x="52.5" y="98" class="sub">Dual-Path Decision</text>
    </g>

    <!-- Origin Shield S3 Cache -->
    <g transform="translate(305, 85)" filter="url(#cShadow1)">
      <rect width="105" height="65" rx="6" fill="#fefce8" stroke="#f59e0b" stroke-width="1.5" />
      <path d="M 52.5 12 Q 70 2 82 12 Q 82 35 70 48 Q 52.5 54 40 48 Q 25 35 25 12 Z" transform="translate(0, 0) scale(0.65) translate(28, 6)" fill="#f59e0b" />
      <text x="52.5" y="85" class="label">Origin Shield</text>
      <text x="52.5" y="98" class="sub">S3 Ephemeral Cache</text>
      <text x="52.5" y="110" font-size="8.5" font-weight="700" fill="#b45309" text-anchor="middle">24h Auto-Purge</text>
    </g>

    <!-- Remote Viewers (Laptop & Phone) -->
    <g transform="translate(170, 235)" filter="url(#cShadow1)">
      <rect x="20" y="12" width="80" height="50" rx="4" fill="#1e293b" stroke="#475569" />
      <rect x="24" y="16" width="72" height="42" rx="2" fill="#0f172a" />
      <polygon points="56,30 56,44 68,37" fill="#f59e0b" />
      <path d="M 12 62 L 108 62 L 104 68 L 16 68 Z" fill="#64748b" />
      <!-- Phone silhouette -->
      <rect x="110" y="20" width="24" height="44" rx="4" fill="#0f172a" stroke="#cbd5e1" stroke-width="1.2" />
      <text x="70" y="88" class="label">👥 Public Viewers</text>
      <text x="70" y="102" class="sub">HTML5 + HLS.js (Zero Install)</text>
      <text x="70" y="114" font-size="8.5" font-weight="700" fill="#15803d" text-anchor="middle">⚡ 90% Uplink Saved (Hit Shield)</text>
    </g>

    <!-- Connectors inside Cloud -->
    <path d="M 155 117 L 180 117" fill="none" stroke="#0284c7" stroke-width="2" marker-end="url(#mBlue1)" />
    <path d="M 285 117 L 305 117" fill="none" stroke="#0284c7" stroke-width="2" marker-end="url(#mBlue1)" />
    <path d="M 232 150 L 232 230" fill="none" stroke="#0284c7" stroke-width="2.2" marker-end="url(#mBlue1)" />
  </g>

  <!-- ========================================================================= -->
  <!-- CONDUIT ARROWS THROUGH THE FIREWALL -->
  <!-- ========================================================================= -->
  <!-- 1. Outbound TLS Reverse Tunnel (Green conduit passing through firewall) -->
  <path d="M 225 155 C 320 155, 360 365, 415 365 L 530 365 C 570 365, 590 190, 620 190 L 620 155" 
        fill="none" stroke="#10b981" stroke-width="3" stroke-dasharray="6,4" marker-end="url(#mGreen1)" />
  <rect x="420" y="352" width="140" height="22" rx="4" fill="#ecfdf5" stroke="#10b981" />
  <text x="490" y="367" class="flow-badge" fill="#065f46" text-anchor="middle">① Outbound TLS Tunnel</text>

  <!-- 2. Blocked Inbound SYN Packet (Red dashed ray bouncing off shield) -->
  <path d="M 620 115 L 470 115" fill="none" stroke="#ef4444" stroke-width="2.5" stroke-dasharray="4,4" marker-end="url(#mRed1)" />
  <rect x="475" y="104" width="80" height="18" rx="3" fill="#fee2e2" stroke="#fca5a5" />
  <text x="515" y="117" font-size="8.5" font-weight="800" fill="#991b1b" text-anchor="middle">SYN DROPPED</text>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_01_nat_traversal.svg")
png_file = os.path.join(fig_dir, "diag_01_nat_traversal.png")
html_file = os.path.join(fig_dir, "diag_01_nat_traversal.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_01_creative)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1180px; height: 500px; }}
</style>
</head>
<body>
  {svg_diag_01_creative}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1180,500",
    html_file
]
subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

