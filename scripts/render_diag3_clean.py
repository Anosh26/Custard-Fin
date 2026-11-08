import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_diag_03_clean = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1160 390" width="1160" height="390">
  <defs>
    <filter id="pShadow" x="-5%" y="-5%" width="110%" height="114%" filterUnits="userSpaceOnUse">
      <feDropShadow dx="0" dy="3" stdDeviation="5" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>

    <!-- Header Gradients -->
    <linearGradient id="gScan" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" /><stop offset="100%" stop-color="#3730a3" />
    </linearGradient>
    <linearGradient id="gProbe" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0d9488" /><stop offset="100%" stop-color="#0f766e" />
    </linearGradient>
    <linearGradient id="gEncode" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" /><stop offset="100%" stop-color="#0369a1" />
    </linearGradient>
    <linearGradient id="gDb" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#7c3aed" /><stop offset="100%" stop-color="#6d28d9" />
    </linearGradient>
    <linearGradient id="gReady" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#10b981" /><stop offset="100%" stop-color="#059669" />
    </linearGradient>

    <!-- Markers -->
    <marker id="arrPipe" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#64748b" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .title { font-size: 18px; font-weight: 700; fill: #0f172a; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .stage-hdr { font-size: 12px; font-weight: 700; fill: #ffffff; text-anchor: middle; }
    .card-h { font-size: 12.5px; font-weight: 700; fill: #0f172a; }
    .item-p { font-size: 10.5px; fill: #475569; }
    .badge-t { font-size: 9.5px; font-weight: 700; }
  </style>

  <!-- Background Base -->
  <rect width="1160" height="390" rx="14" fill="#ffffff" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Header -->
  <text x="36" y="34" class="title">Local Media Ingestion &amp; Hardware Pre-Encoding Pipeline</text>
  <text x="36" y="52" class="subtitle">Asynchronous edge workflow: Automated GPU probe, 4-second GOP segmentation, and atomic SQLite indexing</text>

  <!-- ========================================================================= -->
  <!-- STAGE 1: SCAN & FINGERPRINT -->
  <!-- ========================================================================= -->
  <g transform="translate(36, 75)" filter="url(#pShadow)">
    <rect width="198" height="280" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="198" height="34" rx="10" fill="url(#gScan)" />
    <rect y="24" width="198" height="10" fill="url(#gScan)" />
    <text x="99" y="22" class="stage-hdr">1. Ingest &amp; Scanner</text>

    <!-- Sub-box 1: Folder Scan -->
    <rect x="12" y="46" width="174" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="66" class="card-h" fill="#4338ca">File Ingestion</text>
    <text x="20" y="84" class="item-p">• Recursive directory scan</text>
    <text x="20" y="100" class="item-p">• MP4, MKV, MOV formats</text>

    <!-- Sub-box 2: Deduplication -->
    <rect x="12" y="118" width="174" height="74" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="138" class="card-h" fill="#312e81">SHA-256 Hash</text>
    <text x="20" y="156" class="item-p">• 64KB Header + ModTime</text>
    <text x="20" y="172" class="item-p">• Skips duplicate files</text>
    <rect x="20" y="178" width="94" height="12" rx="2" fill="#ede9fe" />
    <text x="23" y="188" font-size="8.5" font-weight="700" fill="#6d28d9">ZERO RE-ENCODE</text>

    <!-- Poster extract -->
    <rect x="12" y="200" width="174" height="66" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="20" y="220" class="card-h" fill="#047857">Poster Extraction</text>
    <text x="20" y="238" class="item-p">• Extracts 1080p JPEG</text>
    <text x="20" y="254" class="item-p">• Ready for local UI</text>
  </g>

  <!-- Arrow 1 -> 2 -->
  <path d="M 234 215 L 258 215" fill="none" stroke="#64748b" stroke-width="2.5" marker-end="url(#arrPipe)" />

  <!-- ========================================================================= -->
  <!-- STAGE 2: HARDWARE ENCODER PROBE -->
  <!-- ========================================================================= -->
  <g transform="translate(262, 75)" filter="url(#pShadow)">
    <rect width="198" height="280" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="198" height="34" rx="10" fill="url(#gProbe)" />
    <rect y="24" width="198" height="10" fill="url(#gProbe)" />
    <text x="99" y="22" class="stage-hdr">2. Hardware Probe</text>

    <!-- Sub-box 1: Probe Logic -->
    <rect x="12" y="46" width="174" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="66" class="card-h" fill="#0f766e">Auto-Detection</text>
    <text x="20" y="84" class="item-p">• 1-second synthetic probe</text>
    <text x="20" y="100" class="item-p">• Probes driver &amp; session</text>

    <!-- Sub-box 2: Primary GPU -->
    <rect x="12" y="118" width="174" height="74" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="20" y="138" class="card-h" fill="#047857">NVIDIA NVENC</text>
    <text x="20" y="156" class="item-p">• 322.3 fps (10.7x realtime)</text>
    <text x="20" y="172" class="item-p">• Host CPU usage: &lt;5%</text>
    <rect x="20" y="178" width="88" height="12" rx="2" fill="#dcfce7" />
    <text x="23" y="188" font-size="8.5" font-weight="700" fill="#15803d">PRIMARY ASIC</text>

    <!-- Sub-box 3: Fallback -->
    <rect x="12" y="200" width="174" height="66" rx="6" fill="#fff1f2" stroke="#fca5a5" />
    <text x="20" y="220" class="card-h" fill="#be123c">CPU Fallback</text>
    <text x="20" y="238" class="item-p">• libx264 (veryfast)</text>
    <text x="20" y="254" class="item-p">• Seamless error fallback</text>
  </g>

  <!-- Arrow 2 -> 3 -->
  <path d="M 460 215 L 484 215" fill="none" stroke="#64748b" stroke-width="2.5" marker-end="url(#arrPipe)" />

  <!-- ========================================================================= -->
  <!-- STAGE 3: HLS PRE-ENCODING -->
  <!-- ========================================================================= -->
  <g transform="translate(488, 75)" filter="url(#pShadow)">
    <rect width="198" height="280" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="198" height="34" rx="10" fill="url(#gEncode)" />
    <rect y="24" width="198" height="10" fill="url(#gEncode)" />
    <text x="99" y="22" class="stage-hdr">3. HLS Segmentation</text>

    <!-- Sub-box 1: Keyframe alignment -->
    <rect x="12" y="46" width="174" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="66" class="card-h" fill="#0369a1">GOP Alignment</text>
    <text x="20" y="84" class="item-p">• 4-second fixed segments</text>
    <text x="20" y="100" class="item-p">• IDR frame every 120 frames</text>

    <!-- Sub-box 2: Output Chunks -->
    <rect x="12" y="118" width="174" height="74" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="138" class="card-h" fill="#0284c7">Multi-Bitrate HLS</text>
    <text x="20" y="156" class="item-p">• 1080p, 720p, 480p tiers</text>
    <text x="20" y="172" class="item-p">• Generates master.m3u8</text>
    <rect x="20" y="178" width="94" height="12" rx="2" fill="#e0f2fe" />
    <text x="23" y="188" font-size="8.5" font-weight="700" fill="#0369a1">RFC 8216 SPEC</text>

    <!-- Sub-box 3: Chunk storage -->
    <rect x="12" y="200" width="174" height="66" rx="6" fill="#fefce8" stroke="#fde047" />
    <text x="20" y="220" class="card-h" fill="#854d0e">Local Chunks</text>
    <text x="20" y="238" class="item-p">• Stored in .custard/hls/</text>
    <text x="20" y="254" class="item-p">• Ready for Origin Shield</text>
  </g>

  <!-- Arrow 3 -> 4 -->
  <path d="M 686 215 L 710 215" fill="none" stroke="#64748b" stroke-width="2.5" marker-end="url(#arrPipe)" />

  <!-- ========================================================================= -->
  <!-- STAGE 4: SQLITE WAL COMMIT -->
  <!-- ========================================================================= -->
  <g transform="translate(714, 75)" filter="url(#pShadow)">
    <rect width="198" height="280" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="198" height="34" rx="10" fill="url(#gDb)" />
    <rect y="24" width="198" height="10" fill="url(#gDb)" />
    <text x="99" y="22" class="stage-hdr">4. SQLite Commit</text>

    <!-- Sub-box 1: DB Transaction -->
    <rect x="12" y="46" width="174" height="64" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="66" class="card-h" fill="#6d28d9">Atomic Record</text>
    <text x="20" y="84" class="item-p">• Writes metadata to SQLite</text>
    <text x="20" y="100" class="item-p">• Pure-Go driver (no CGO)</text>

    <!-- Sub-box 2: WAL Engine -->
    <rect x="12" y="118" width="174" height="74" rx="6" fill="#ffffff" stroke="#e2e8f0" />
    <text x="20" y="138" class="card-h" fill="#5b21b6">WAL Mode</text>
    <text x="20" y="156" class="item-p">• Non-blocking concurrent reads</text>
    <text x="20" y="172" class="item-p">• Zero database lock contention</text>
    <rect x="20" y="178" width="94" height="12" rx="2" fill="#ede9fe" />
    <text x="23" y="188" font-size="8.5" font-weight="700" fill="#6d28d9">HIGH CONCURRENCY</text>

    <!-- Sub-box 3: Private Storage -->
    <rect x="12" y="200" width="174" height="66" rx="6" fill="#f8fafc" stroke="#cbd5e1" />
    <text x="20" y="220" class="card-h" fill="#334155">Private at Home</text>
    <text x="20" y="238" class="item-p">• Watch progress stays edge</text>
    <text x="20" y="254" class="item-p">• Preserves Invariant 8</text>
  </g>

  <!-- Arrow 4 -> 5 -->
  <path d="M 912 215 L 936 215" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#arrPipe)" />

  <!-- ========================================================================= -->
  <!-- STAGE 5: DUAL AVAILABILITY (READY) -->
  <!-- ========================================================================= -->
  <g transform="translate(940, 75)" filter="url(#pShadow)">
    <rect width="184" height="280" rx="10" fill="#f0fdf4" stroke="#86efac" stroke-width="1.2" />
    <rect width="184" height="34" rx="10" fill="url(#gReady)" />
    <rect y="24" width="184" height="10" fill="url(#gReady)" />
    <text x="92" y="22" class="stage-hdr">5. Active Playback</text>

    <!-- Option 1: LAN Instant -->
    <rect x="10" y="46" width="164" height="98" rx="6" fill="#ffffff" stroke="#86efac" />
    <text x="18" y="66" class="card-h" fill="#047857">LAN Streaming</text>
    <text x="18" y="84" class="item-p">• Embedded :8089 engine</text>
    <text x="18" y="100" class="item-p">• Instant playback (&lt;4ms)</text>
    <text x="18" y="116" class="item-p">• 0 internet data used</text>
    <rect x="18" y="124" width="86" height="12" rx="2" fill="#dcfce7" />
    <text x="21" y="134" font-size="8.5" font-weight="700" fill="#15803d">OFFLINE READY</text>

    <!-- Option 2: Public Toggle -->
    <rect x="10" y="154" width="164" height="112" rx="6" fill="#ffffff" stroke="#7dd3fc" />
    <text x="18" y="174" class="card-h" fill="#0369a1">Public Toggle</text>
    <text x="18" y="192" class="item-p">• 1-click sharing toggle</text>
    <text x="18" y="208" class="item-p">• Syncs ~3 KB manifest</text>
    <text x="18" y="224" class="item-p">• Discoverable globally</text>
    <rect x="18" y="234" width="94" height="12" rx="2" fill="#e0f2fe" />
    <text x="21" y="244" font-size="8.5" font-weight="700" fill="#0369a1">NO VIDEO UPLOAD</text>
  </g>

</svg>
"""

svg_file = os.path.join(fig_dir, "diag_03_ingest_pipeline.svg")
png_file = os.path.join(fig_dir, "diag_03_ingest_pipeline.png")
html_file = os.path.join(fig_dir, "diag_03_ingest_pipeline.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_diag_03_clean)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1160px; height: 390px; }}
</style>
</head>
<body>
  {svg_diag_03_clean}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1160,390",
    html_file
]
subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

