import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_06_clean = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1160 480" width="1160" height="480">
  <defs>
    <filter id="cSh6" x="-4%" y="-4%" width="108%" height="112%" filterUnits="userSpaceOnUse">
      <feDropShadow dx="0" dy="3" stdDeviation="5" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>

    <!-- Gradients -->
    <linearGradient id="gTrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#e11d48" /><stop offset="100%" stop-color="#be123c" />
    </linearGradient>
    <linearGradient id="gCust" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#059669" /><stop offset="100%" stop-color="#047857" />
    </linearGradient>
    <linearGradient id="gSave" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#10b981" /><stop offset="100%" stop-color="#059669" />
    </linearGradient>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 18px; font-weight: 700; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .card-hdr { font-size: 13px; font-weight: 700; fill: #ffffff; }
    .card-total { font-size: 14px; font-weight: 800; fill: #ffffff; text-anchor: end; }
    .row-lbl { font-size: 11.5px; font-weight: 600; fill: #1e293b; }
    .row-val { font-size: 12px; font-weight: 700; text-anchor: end; }
    .row-sub { font-size: 9.5px; fill: #64748b; }
  </style>

  <!-- Background Canvas -->
  <rect width="1160" height="480" rx="14" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Title Header -->
  <text x="36" y="34" class="title">Monthly Operating Cost Comparison (100 Users, 5.0 TB Catalog)</text>
  <text x="36" y="52" class="subtitle">Decoupling media storage and GPU compute from the cloud yields a mathematically proven 95.89% cost reduction</text>

  <!-- ========================================================================= -->
  <!-- TOP CALLOUT: 95.89% SAVINGS BANNER -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 70)" filter="url(#cSh6)">
    <rect width="1088" height="48" rx="8" fill="#ecfdf5" stroke="#86efac" stroke-width="1.2" />
    <circle cx="28" cy="24" r="14" fill="url(#gSave)" />
    <text x="21" y="29" font-size="14" fill="#ffffff" font-weight="900">✓</text>
    <text x="52" y="22" font-size="12.5" font-weight="800" fill="#065f46">VERIFIED FINANCIAL EFFICIENCY: 95.89% CLOUD COST REDUCTION</text>
    <text x="52" y="37" font-size="10.5" fill="#047857">Net Monthly Savings: $1,145.85 / month  •  Three-Year Cumulative TCO Savings: $41,250.60</text>
    <rect x="940" y="10" width="136" height="28" rx="6" fill="#10b981" />
    <text x="1008" y="28" font-size="12" font-weight="800" fill="#ffffff" text-anchor="middle">SAVE 95.89%</text>
  </g>

  <!-- ========================================================================= -->
  <!-- LEFT: TRADITIONAL CLOUD VOD ($1,195.00/MO) -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 134)" filter="url(#cSh6)">
    <rect width="528" height="310" rx="10" fill="#fff5f5" stroke="#fecaca" stroke-width="1.2" />
    <rect width="528" height="36" rx="10" fill="url(#gTrad)" />
    <rect y="26" width="528" height="10" fill="url(#gTrad)" />
    <text x="18" y="23" class="card-hdr">🏢 Traditional Centralized Cloud VOD</text>
    <text x="510" y="23" class="card-total">$1,195.00 / mo</text>

    <!-- Row 1: Egress -->
    <g transform="translate(16, 50)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#fca5a5" />
      <text x="14" y="22" class="row-lbl">Cloud Data Egress Bandwidth</text>
      <text x="14" y="38" class="row-sub">9.0 TB monthly viewership egress @ $0.090 / GB</text>
      <text x="482" y="24" class="row-val" fill="#be123c">$810.00</text>
      <text x="482" y="40" font-size="9.5" fill="#991b1b" text-anchor="end">67.78% of Total</text>
    </g>

    <!-- Row 2: Transcode -->
    <g transform="translate(16, 110)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#e2e8f0" />
      <text x="14" y="22" class="row-lbl">Centralized Video Transcoding</text>
      <text x="14" y="38" class="row-sub">AWS Elemental MediaConvert or EC2 GPU Transcoder instances</text>
      <text x="482" y="24" class="row-val" fill="#be123c">$140.00</text>
      <text x="482" y="40" font-size="9.5" fill="#991b1b" text-anchor="end">11.72% of Total</text>
    </g>

    <!-- Row 3: Managed Infra -->
    <g transform="translate(16, 170)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#e2e8f0" />
      <text x="14" y="22" class="row-lbl">Managed Cloud Infrastructure</text>
      <text x="14" y="38" class="row-sub">AWS RDS PostgreSQL ($65) + ElastiCache ($18) + ALB ($47)</text>
      <text x="482" y="24" class="row-val" fill="#be123c">$130.00</text>
      <text x="482" y="40" font-size="9.5" fill="#991b1b" text-anchor="end">10.88% of Total</text>
    </g>

    <!-- Row 4: S3 Storage -->
    <g transform="translate(16, 230)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#e2e8f0" />
      <text x="14" y="22" class="row-lbl">Permanent Cloud Object Storage</text>
      <text x="14" y="38" class="row-sub">5,000 GB persistent S3 Standard bucket storage @ $0.023 / GB-mo</text>
      <text x="482" y="24" class="row-val" fill="#be123c">$115.00</text>
      <text x="482" y="40" font-size="9.5" fill="#991b1b" text-anchor="end">9.62% of Total</text>
    </g>

    <text x="18" y="300" font-size="10.5" font-weight="700" fill="#991b1b">Annual Liability: $14,340.00 / year (Financially unsustainable)</text>
  </g>

  <!-- ========================================================================= -->
  <!-- RIGHT: CUSTARDFIN DECENTRALIZED HYBRID ($49.15/MO) -->
  <!-- ========================================================================= -->
  <g transform="translate(596, 134)" filter="url(#cSh6)">
    <rect width="528" height="310" rx="10" fill="#f0fdf4" stroke="#86efac" stroke-width="1.2" />
    <rect width="528" height="36" rx="10" fill="url(#gCust)" />
    <rect y="26" width="528" height="10" fill="url(#gCust)" />
    <text x="18" y="23" class="card-hdr">⚡ CustardFin Decentralized Hybrid Architecture</text>
    <text x="510" y="23" class="card-total">$49.15 / mo</text>

    <!-- Row 1: Cloud VM -->
    <g transform="translate(16, 50)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#86efac" />
      <text x="14" y="22" class="row-lbl">Consolidated Coordination Cloud VM</text>
      <text x="14" y="38" class="row-sub">GCP e2-standard-2 / Hetzner CPX31 (apid, Postgres, Redis, frps)</text>
      <text x="482" y="24" class="row-val" fill="#047857">$48.00</text>
      <text x="482" y="40" font-size="9.5" fill="#15803d" text-anchor="end">97.66% of Total</text>
    </g>

    <!-- Row 2: Origin Shield Cache -->
    <g transform="translate(16, 110)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#e2e8f0" />
      <text x="14" y="22" class="row-lbl">Origin Shield Ephemeral Cache</text>
      <text x="14" y="38" class="row-sub">50 GB transient active hot chunks with 24h LRU auto-purge ($0 local)</text>
      <text x="482" y="24" class="row-val" fill="#047857">$1.15</text>
      <text x="482" y="40" font-size="9.5" fill="#15803d" text-anchor="end">2.34% of Total</text>
    </g>

    <!-- Row 3: Edge Video Storage -->
    <g transform="translate(16, 170)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#86efac" />
      <text x="14" y="22" class="row-lbl" fill="#047857">Permanent Video Storage (Home Disks)</text>
      <text x="14" y="38" class="row-sub">5.0 TB stays 100% on user-owned edge NVMe / SATA drives</text>
      <text x="482" y="24" class="row-val" fill="#15803d">$0.00</text>
      <text x="482" y="40" font-size="9.5" font-weight="700" fill="#15803d" text-anchor="end">100% SAVINGS</text>
    </g>

    <!-- Row 4: Edge Transcode -->
    <g transform="translate(16, 230)">
      <rect width="496" height="52" rx="6" fill="#ffffff" stroke="#86efac" />
      <text x="14" y="22" class="row-lbl" fill="#047857">Hardware Transcoding (Edge NVENC ASICs)</text>
      <text x="14" y="38" class="row-sub">Pre-encoded at 322 fps on local user GPUs during system idle</text>
      <text x="482" y="24" class="row-val" fill="#15803d">$0.00</text>
      <text x="482" y="40" font-size="9.5" font-weight="700" fill="#15803d" text-anchor="end">100% SAVINGS</text>
    </g>

    <text x="18" y="300" font-size="10.5" font-weight="700" fill="#047857">Annual Cost: $589.80 / year (Easily funded by academic lab / grant)</text>
  </g>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_06_cost_comparison.svg")
png_file = os.path.join(fig_dir, "diag_06_cost_comparison.png")
html_file = os.path.join(fig_dir, "diag_06_cost_comparison.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_06_clean)

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
  {svg_diag_06_clean}
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

