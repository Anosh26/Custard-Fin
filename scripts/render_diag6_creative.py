import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_06_creative = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1180 500" width="1180" height="500">
  <defs>
    <filter id="costShadow" x="-10%" y="-10%" width="120%" height="120%">
      <feDropShadow dx="0" dy="6" stdDeviation="8" flood-color="#0f172a" flood-opacity="0.08" />
    </filter>

    <!-- Pillar Gradients -->
    <linearGradient id="gEgress" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#f43f5e" /><stop offset="100%" stop-color="#e11d48" />
    </linearGradient>
    <linearGradient id="gTranscode" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#fb7185" /><stop offset="100%" stop-color="#f43f5e" />
    </linearGradient>
    <linearGradient id="gInfra" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#fda4af" /><stop offset="100%" stop-color="#fb7185" />
    </linearGradient>
    <linearGradient id="gStorage" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#fecdd3" /><stop offset="100%" stop-color="#fda4af" />
    </linearGradient>

    <!-- CustardFin Gradients -->
    <linearGradient id="gCustVM" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#10b981" /><stop offset="100%" stop-color="#059669" />
    </linearGradient>
    <linearGradient id="gCustShield" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#6ee7b7" /><stop offset="100%" stop-color="#34d399" />
    </linearGradient>

    <!-- Badge Gradient -->
    <linearGradient id="gBadgeSave" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#10b981" /><stop offset="100%" stop-color="#047857" />
    </linearGradient>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 19px; font-weight: 800; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .pillar-hdr { font-size: 14px; font-weight: 800; }
    .pillar-price { font-size: 20px; font-weight: 900; }
    .seg-lbl { font-size: 11px; font-weight: 700; fill: #ffffff; }
    .seg-val { font-size: 10px; fill: #ffffff; opacity: 0.95; }
    .note-h { font-size: 11.5px; font-weight: 700; fill: #0f172a; }
    .note-p { font-size: 10px; fill: #64748b; }
  </style>

  <!-- Canvas Base -->
  <rect width="1180" height="500" rx="16" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Header -->
  <text x="40" y="36" class="title">Monthly Operating Cost Comparison (100 Users, 5.0 TB Catalog)</text>
  <text x="40" y="54" class="subtitle">Decoupling video binaries and GPU encoding from the cloud yields a mathematically proven 95.89% recurring expenditure reduction</text>

  <!-- ========================================================================= -->
  <!-- TOWER 1: TRADITIONAL CLOUD VOD ($1,195.00/MO) -->
  <!-- ========================================================================= -->
  <g transform="translate(60, 90)" filter="url(#costShadow)">
    <!-- Base Platform Card -->
    <rect width="360" height="375" rx="12" fill="#fff5f5" stroke="#fecaca" stroke-width="1.2" />
    <text x="24" y="32" class="pillar-hdr" fill="#991b1b">TRADITIONAL CLOUD VOD</text>
    <text x="24" y="58" class="pillar-price" fill="#be123c">$1,195.00 <tspan font-size="12" font-weight="600" fill="#991b1b">/ month</tspan></text>
    <text x="24" y="74" font-size="10" font-weight="700" fill="#e11d48">$14,340.00 / year recurring liability</text>

    <!-- 3D Stacked Tower Graphic -->
    <!-- Segment 1: Cloud Egress ($810 - Height: 140px) -->
    <rect x="24" y="90" width="312" height="135" rx="6" fill="url(#gEgress)" />
    <text x="36" y="112" class="seg-lbl">Cloud Data Egress (9 TB @ $0.09/GB)</text>
    <text x="36" y="128" class="seg-val">$810.00 / mo • 67.78% of Total Cost</text>
    <text x="36" y="142" font-size="9" fill="#fee2e2">Prohibitive public internet bandwidth toll</text>

    <!-- Segment 2: Transcoding ($140 - Height: 42px) -->
    <rect x="24" y="230" width="312" height="38" rx="5" fill="url(#gTranscode)" />
    <text x="36" y="248" class="seg-lbl">Transcode Compute (EC2 / MediaConvert)</text>
    <text x="36" y="260" class="seg-val">$140.00 / mo • 11.72% of Total</text>

    <!-- Segment 3: Managed Infra ($130 - Height: 38px) -->
    <rect x="24" y="272" width="312" height="36" rx="5" fill="url(#gInfra)" />
    <text x="36" y="289" class="seg-lbl">Managed DBs (RDS + Redis + ALB)</text>
    <text x="36" y="301" class="seg-val">$130.00 / mo • 10.88% of Total</text>

    <!-- Segment 4: S3 Storage ($115 - Height: 35px) -->
    <rect x="24" y="312" width="312" height="35" rx="5" fill="url(#gStorage)" />
    <text x="36" y="328" class="seg-lbl" fill="#881337">S3 Object Storage (5 TB @ $0.023/GB)</text>
    <text x="36" y="340" class="seg-val" fill="#881337">$115.00 / mo • 9.62% of Total</text>

    <text x="24" y="365" font-size="9" font-weight="700" fill="#991b1b">Result: Financially unfeasible for schools, labs, or students</text>
  </g>

  <!-- ========================================================================= -->
  <!-- CENTER METRIC: 95.89% SAVINGS BADGE & ARROW -->
  <!-- ========================================================================= -->
  <g transform="translate(450, 150)">
    <!-- Large Savings Circle -->
    <circle cx="140" cy="80" r="70" fill="url(#gBadgeSave)" filter="url(#costShadow)" />
    <circle cx="140" cy="80" r="62" fill="none" stroke="#ffffff" stroke-width="2" stroke-dasharray="4,3" />
    <text x="140" y="68" font-size="28" font-weight="900" fill="#ffffff" text-anchor="middle">95.9%</text>
    <text x="140" y="86" font-size="11" font-weight="800" fill="#d1fae5" text-anchor="middle">SAVINGS</text>
    <text x="140" y="100" font-size="9" fill="#ffffff" opacity="0.9" text-anchor="middle">VERIFIED</text>

    <!-- Savings Callouts -->
    <rect x="15" y="180" width="250" height="70" rx="8" fill="#f8fafc" stroke="#e2e8f0" filter="url(#costShadow)" />
    <text x="140" y="202" font-size="12" font-weight="800" fill="#047857" text-anchor="middle">-$1,145.85 / MONTH</text>
    <text x="140" y="218" font-size="10" font-weight="700" fill="#0f172a" text-anchor="middle">-$13,750.20 / year saved</text>
    <text x="140" y="234" font-size="9.5" fill="#64748b" text-anchor="middle">3-Year TCO Savings: $41,250.60</text>

    <!-- Comparative Flow Arrow (Red -> Green) -->
    <path d="M -15 80 L 45 80" fill="none" stroke="#e11d48" stroke-width="3" stroke-dasharray="4,4" />
    <path d="M 235 80 L 295 80" fill="none" stroke="#10b981" stroke-width="3" />
    <polygon points="295,75 305,80 295,85" fill="#10b981" />
  </g>

  <!-- ========================================================================= -->
  <!-- TOWER 2: CUSTARDFIN DECENTRALIZED HYBRID ($49.15/MO) -->
  <!-- ========================================================================= -->
  <g transform="translate(760, 90)" filter="url(#costShadow)">
    <!-- Base Platform Card -->
    <rect width="360" height="375" rx="12" fill="#f0fdf4" stroke="#86efac" stroke-width="1.2" />
    <text x="24" y="32" class="pillar-hdr" fill="#065f46">CUSTARDFIN HYBRID ARCHITECTURE</text>
    <text x="24" y="58" class="pillar-price" fill="#047857">$49.15 <tspan font-size="12" font-weight="600" fill="#065f46">/ month</tspan></text>
    <text x="24" y="74" font-size="10" font-weight="700" fill="#059669">$589.80 / year total operating expenditure</text>

    <!-- CustardFin Segment 1: Consolidated Cloud VM ($48 - Height: 55px) -->
    <rect x="24" y="90" width="312" height="52" rx="6" fill="url(#gCustVM)" />
    <text x="36" y="112" class="seg-lbl">Consolidated Coordination Cloud VM</text>
    <text x="36" y="128" class="seg-val">$48.00 / mo • GCP e2-standard-2 / Hetzner CPX31</text>

    <!-- CustardFin Segment 2: Origin Shield Cache ($1.15 - Height: 24px) -->
    <rect x="24" y="146" width="312" height="26" rx="4" fill="url(#gCustShield)" />
    <text x="36" y="163" class="seg-lbl" fill="#064e3b">Origin Shield Ephemeral S3 Cache: $1.15 / mo</text>

    <!-- Zero Cost Callout 1: Video Storage -->
    <g transform="translate(24, 185)">
      <rect width="312" height="48" rx="6" fill="#ffffff" stroke="#86efac" />
      <circle cx="20" cy="24" r="9" fill="#10b981" />
      <text x="16" y="28" font-size="11" fill="#ffffff" font-weight="900">✓</text>
      <text x="38" y="21" class="note-h" fill="#047857">Permanent Video Storage: $0.00</text>
      <text x="38" y="36" class="note-p">100% on user-owned edge NVMe / SATA drives</text>
      <text x="295" y="28" font-size="11" font-weight="800" fill="#15803d" text-anchor="end">FREE</text>
    </g>

    <!-- Zero Cost Callout 2: Video Transcoding -->
    <g transform="translate(24, 240)">
      <rect width="312" height="48" rx="6" fill="#ffffff" stroke="#86efac" />
      <circle cx="20" cy="24" r="9" fill="#10b981" />
      <text x="16" y="28" font-size="11" fill="#ffffff" font-weight="900">✓</text>
      <text x="38" y="21" class="note-h" fill="#047857">Hardware Transcoding: $0.00</text>
      <text x="38" y="36" class="note-p">Offloaded to local RTX GPUs (322 fps NVENC)</text>
      <text x="295" y="28" font-size="11" font-weight="800" fill="#15803d" text-anchor="end">FREE</text>
    </g>

    <!-- Zero Cost Callout 3: Egress Bandwidth -->
    <g transform="translate(24, 295)">
      <rect width="312" height="48" rx="6" fill="#ffffff" stroke="#86efac" />
      <circle cx="20" cy="24" r="9" fill="#10b981" />
      <text x="16" y="28" font-size="11" fill="#ffffff" font-weight="900">✓</text>
      <text x="38" y="21" class="note-h" fill="#047857">Egress Bandwidth: $0.00 Overages</text>
      <text x="38" y="36" class="note-p">Direct LAN bypass + Bundled VPS traffic</text>
      <text x="295" y="28" font-size="11" font-weight="800" fill="#15803d" text-anchor="end">FREE</text>
    </g>

    <text x="24" y="365" font-size="9" font-weight="700" fill="#047857">Result: Easily hosted on student budget or lab server</text>
  </g>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_06_cost_comparison.svg")
png_file = os.path.join(fig_dir, "diag_06_cost_comparison.png")
html_file = os.path.join(fig_dir, "diag_06_cost_comparison.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_06_creative)

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
  {svg_diag_06_creative}
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

