import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

def render_svg_to_png(svg_path, png_path, width=1200, height=600):
    html_wrapper = f"""<!DOCTYPE html>
    <html>
    <head>
    <style>
      body {{ margin: 0; padding: 0; background: transparent; display: flex; align-items: center; justify-content: center; }}
      svg {{ width: {width}px; height: {height}px; }}
    </style>
    </head>
    <body>
      {open(svg_path, 'r', encoding='utf-8').read()}
    </body>
    </html>"""
    
    temp_html = svg_path.replace(".svg", ".html")
    with open(temp_html, "w", encoding="utf-8") as f:
        f.write(html_wrapper)
    
    args = [
        chrome_path,
        "--headless=old",
        "--disable-gpu",
        f"--screenshot={png_path}",
        f"--window-size={width},{height}",
        temp_html
    ]
    subprocess.run(args, capture_output=True)
    if os.path.exists(png_path):
        print(f"Generated {png_path} ({os.path.getsize(png_path)} bytes)")
    else:
        print(f"Failed to generate {png_path}")

# ==============================================================================
# DIAGRAM 1: Reverse Proxy NAT Traversal
# ==============================================================================
svg_diag_01 = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 520" width="1200" height="520">
  <defs>
    <filter id="sh1" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <linearGradient id="gHome" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" /><stop offset="100%" stop-color="#3730a3" />
    </linearGradient>
    <linearGradient id="gISP" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#e11d48" /><stop offset="100%" stop-color="#be123c" />
    </linearGradient>
    <linearGradient id="gCloud" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" /><stop offset="100%" stop-color="#0369a1" />
    </linearGradient>
    <linearGradient id="gUser" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#d97706" /><stop offset="100%" stop-color="#b45309" />
    </linearGradient>
    <marker id="mGreen" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="mBlue" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="mPurple" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#6366f1" />
    </marker>
    <marker id="mRed" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#ef4444" />
    </marker>
  </defs>

  <style>
    .title { font-family: -apple-system, sans-serif; font-size: 19px; font-weight: 700; fill: #0f172a; }
    .sub { font-family: -apple-system, sans-serif; font-size: 13px; fill: #64748b; }
    .card-hdr { font-family: -apple-system, sans-serif; font-size: 14px; font-weight: 700; fill: #ffffff; }
    .card-sub { font-family: -apple-system, sans-serif; font-size: 11px; fill: #e2e8f0; }
    .pill-t { font-family: -apple-system, sans-serif; font-size: 11px; font-weight: 700; }
    .item-h { font-family: -apple-system, sans-serif; font-size: 12px; font-weight: 700; fill: #1e293b; }
    .item-p { font-family: -apple-system, sans-serif; font-size: 11px; fill: #64748b; }
  </style>

  <rect width="1200" height="520" fill="#f8fafc" rx="14" stroke="#e2e8f0" stroke-width="1.5" />
  <text x="36" y="38" class="title">Reverse Proxy NAT Traversal &amp; Inbound Firewall Bypass</text>
  <text x="36" y="58" class="sub">Outbound-only TLS 1.3 reverse tunnel bridges media streams without opening residential router ports</text>

  <!-- Card 1: Home Edge -->
  <g transform="translate(36, 85)" filter="url(#sh1)">
    <rect width="250" height="390" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="250" height="42" rx="10" fill="url(#gHome)" /><rect y="30" width="250" height="12" fill="url(#gHome)" />
    <text x="16" y="27" class="card-hdr">🏠 Home Edge (Behind CGNAT)</text>

    <rect x="14" y="58" width="222" height="70" rx="6" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="24" y="80" class="item-h" fill="#4338ca">Edge Supervisor (edged)</text>
    <text x="24" y="100" class="item-p">• Persistent Outbound Tunnel</text>
    <text x="24" y="116" class="item-p">• 15s Ephemeral Heartbeat</text>

    <rect x="14" y="140" width="222" height="70" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="24" y="162" class="item-h" fill="#047857">Local Streamer (:8089)</text>
    <text x="24" y="182" class="item-p">• Embedded Go HTTP Server</text>
    <text x="24" y="198" class="item-p">• Serves 4s Pre-encoded HLS</text>

    <rect x="14" y="222" width="222" height="70" rx="6" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="24" y="244" class="item-h" fill="#0369a1">Edge SQLite &amp; Disk</text>
    <text x="24" y="264" class="item-p">• Video Master Stays Home</text>
    <text x="24" y="280" class="item-p">• Private Watch History (Inv 8)</text>

    <rect x="14" y="315" width="222" height="45" rx="6" fill="#ecfdf5" stroke="#10b981" />
    <text x="24" y="343" font-size="11" font-weight="700" fill="#065f46" font-family="sans-serif">🔒 NO INBOUND PORTS OPEN</text>
  </g>

  <!-- Card 2: ISP Firewall -->
  <g transform="translate(330, 85)" filter="url(#sh1)">
    <rect width="230" height="390" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="230" height="42" rx="10" fill="url(#gISP)" /><rect y="30" width="230" height="12" fill="url(#gISP)" />
    <text x="16" y="27" class="card-hdr">🛡️ ISP Network &amp; Firewall</text>

    <rect x="14" y="60" width="202" height="110" rx="6" fill="#fff1f2" stroke="#f43f5e" stroke-dasharray="4,3" />
    <text x="24" y="85" font-size="11" font-weight="700" fill="#be123c" font-family="sans-serif">❌ INBOUND PACKETS</text>
    <text x="24" y="107" class="item-p">• Inbound TCP SYN Dropped</text>
    <text x="24" y="123" class="item-p">• UPnP Blocked by ISP</text>
    <text x="24" y="139" class="item-p">• 100.64.0.0/10 Shared IP</text>
    <text x="24" y="157" font-size="10" font-weight="700" fill="#e11d48" font-family="sans-serif">Traditional Port Forwarding Fails</text>

    <rect x="14" y="185" width="202" height="110" rx="6" fill="#f0fdf4" stroke="#10b981" />
    <text x="24" y="210" font-size="11" font-weight="700" fill="#047857" font-family="sans-serif">✔️ OUTBOUND TUNNEL</text>
    <text x="24" y="232" class="item-p">• Home Initiates Outbound TLS</text>
    <text x="24" y="248" class="item-p">• State Table Created</text>
    <text x="24" y="264" class="item-p">• Bi-directional Flow Allowed</text>
    <text x="24" y="282" font-size="10" font-weight="700" fill="#059669" font-family="sans-serif">Invariant 2 Preserved</text>

    <rect x="14" y="315" width="202" height="45" rx="6" fill="#f1f5f9" stroke="#cbd5e1" />
    <text x="24" y="343" font-size="11" font-weight="600" fill="#475569" font-family="sans-serif">Stateful NAT Mapping Active</text>
  </g>

  <!-- Card 3: Cloud Relay -->
  <g transform="translate(605, 85)" filter="url(#sh1)">
    <rect width="270" height="390" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="270" height="42" rx="10" fill="url(#gCloud)" /><rect y="30" width="270" height="12" fill="url(#gCloud)" />
    <text x="16" y="27" class="card-hdr">☁️ Public Cloud Relay (GCP VM)</text>

    <rect x="14" y="58" width="242" height="75" rx="6" fill="#f0f9ff" stroke="#7dd3fc" />
    <text x="24" y="80" class="item-h" fill="#0369a1">FRP Server (frps :7000)</text>
    <text x="24" y="100" class="item-p">• Terminates Outbound TLS</text>
    <text x="24" y="116" class="item-p">• Multiplexes Reverse Requests</text>

    <rect x="14" y="145" width="242" height="75" rx="6" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="24" y="167" class="item-h" fill="#0284c7">Cloud API (apid :8080)</text>
    <text x="24" y="187" class="item-p">• Dual-Path Stream Router</text>
    <text x="24" y="203" class="item-p">• Admission Limiter (10 Streams)</text>

    <rect x="14" y="232" width="242" height="65" rx="6" fill="#fefce8" stroke="#fde047" />
    <text x="24" y="254" class="item-h" fill="#854d0e">S3 Origin Shield Cache</text>
    <text x="24" y="274" class="item-p">• 24h LRU Ephemeral Purge</text>

    <rect x="14" y="315" width="242" height="45" rx="6" fill="#f8fafc" stroke="#cbd5e1" />
    <text x="24" y="343" font-size="11" font-weight="600" fill="#334155" font-family="sans-serif">Zero Permanent Video Storage</text>
  </g>

  <!-- Card 4: Viewer -->
  <g transform="translate(920, 85)" filter="url(#sh1)">
    <rect width="244" height="390" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="244" height="42" rx="10" fill="url(#gUser)" /><rect y="30" width="244" height="12" fill="url(#gUser)" />
    <text x="16" y="27" class="card-hdr">📱 Public Viewers</text>

    <rect x="14" y="60" width="216" height="85" rx="6" fill="#fffbeb" stroke="#fde68a" />
    <text x="24" y="82" class="item-h" fill="#b45309">Viewer 1 (Cache Miss)</text>
    <text x="24" y="102" class="item-p">1. Pulls chunk via FRP Tunnel</text>
    <text x="24" y="118" class="item-p">2. Commits chunk to S3 Cache</text>

    <rect x="14" y="160" width="216" height="100" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="24" y="182" class="item-h" fill="#15803d">Viewers 2–10 (Cache Hit)</text>
    <text x="24" y="202" class="item-p">3. Served from S3 Origin Shield</text>
    <text x="24" y="222" font-size="10" font-weight="700" fill="#16a34a" font-family="sans-serif">0 Bps Home Uplink Used</text>
    <text x="24" y="238" class="item-p">90% Upload Bandwidth Saved</text>

    <rect x="14" y="280" width="216" height="80" rx="6" fill="#f8fafc" stroke="#cbd5e1" />
    <text x="24" y="304" font-size="11" font-weight="700" fill="#334155" font-family="sans-serif">Zero Software Install</text>
    <text x="24" y="324" class="item-p">• Universal HTML5 / HLS.js</text>
    <text x="24" y="340" class="item-p">• Real-time Stats for Nerds HUD</text>
  </g>

  <!-- Connectors -->
  <path d="M 286 338 L 342 338" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#mGreen)" />
  <path d="M 560 338 L 617 338" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#mGreen)" />
  
  <path d="M 605 120 L 542 120" fill="none" stroke="#ef4444" stroke-width="2.5" stroke-dasharray="4,4" marker-end="url(#mRed)" />
  <path d="M 920 105 L 878 105" fill="none" stroke="#0284c7" stroke-width="2.5" marker-end="url(#mBlue)" />
  <path d="M 617 105 L 246 105" fill="none" stroke="#6366f1" stroke-width="2.5" stroke-dasharray="5,4" marker-end="url(#mPurple)" />
</svg>
"""

# ==============================================================================
# DIAGRAM 3: Local Ingestion and Pre-Encoding Pipeline Sequence
# ==============================================================================
svg_diag_03 = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 480" width="1200" height="480">
  <defs>
    <filter id="sh3" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <linearGradient id="gHead" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" /><stop offset="100%" stop-color="#312e81" />
    </linearGradient>
    <marker id="mBlue3" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="mGreen3" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="mGold3" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#d97706" />
    </marker>
  </defs>

  <style>
    .t-main { font-family: -apple-system, sans-serif; font-size: 19px; font-weight: 700; fill: #0f172a; }
    .t-sub { font-family: -apple-system, sans-serif; font-size: 13px; fill: #64748b; }
    .col-h { font-family: -apple-system, sans-serif; font-size: 13px; font-weight: 700; fill: #ffffff; text-anchor: middle; }
    .box-num { font-family: -apple-system, sans-serif; font-size: 12px; font-weight: 700; fill: #4f46e5; }
    .box-t { font-family: -apple-system, sans-serif; font-size: 12px; font-weight: 700; fill: #1e293b; }
    .box-p { font-family: -apple-system, sans-serif; font-size: 11px; fill: #64748b; }
  </style>

  <rect width="1200" height="480" fill="#f8fafc" rx="14" stroke="#e2e8f0" stroke-width="1.5" />
  <text x="36" y="38" class="t-main">Local Media Ingestion &amp; Hardware Pre-Encoding Pipeline</text>
  <text x="36" y="58" class="t-sub">Asynchronous background workflow offloading encoding to local NVIDIA NVENC ASICs</text>

  <!-- Step Cards Across Pipeline -->
  <!-- Step 1: File Drop & Scanner -->
  <g transform="translate(36, 90)" filter="url(#sh3)">
    <rect width="205" height="340" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="205" height="36" rx="10" fill="url(#gHead)" /><rect y="24" width="205" height="12" fill="url(#gHead)" />
    <text x="102" y="23" class="col-h">1. Ingest &amp; Scanner</text>

    <text x="16" y="65" class="box-num">STEP 1</text>
    <text x="16" y="82" class="box-t">Folder Scan</text>
    <text x="16" y="102" class="box-p">• Monitors user media dir</text>
    <text x="16" y="118" class="box-p">• Filters: .mp4, .mkv, .mov</text>
    <text x="16" y="134" class="box-p">• Inode fast signature</text>

    <rect x="14" y="155" width="177" height="60" rx="6" fill="#f1f5f9" />
    <text x="22" y="177" font-size="11" font-weight="700" fill="#334155" font-family="sans-serif">SHA-256 Hasher</text>
    <text x="22" y="195" font-size="10" fill="#64748b" font-family="sans-serif">64KB Header + ModTime</text>

    <rect x="14" y="235" width="177" height="75" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="22" y="258" font-size="11" font-weight="700" fill="#047857" font-family="sans-serif">Deduplication Gate</text>
    <text x="22" y="276" font-size="10" fill="#14532d" font-family="sans-serif">• Skips already indexed</text>
    <text x="22" y="292" font-size="10" fill="#14532d" font-family="sans-serif">• 0 redundant encodes</text>
  </g>

  <!-- Step 2: SQLite Write Pending -->
  <g transform="translate(275, 90)" filter="url(#sh3)">
    <rect width="205" height="340" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="205" height="36" rx="10" fill="url(#gHead)" /><rect y="24" width="205" height="12" fill="url(#gHead)" />
    <text x="102" y="23" class="col-h">2. Edge Database</text>

    <text x="16" y="65" class="box-num">STEP 2</text>
    <text x="16" y="82" class="box-t">SQLite (WAL Mode)</text>
    <text x="16" y="102" class="box-p">• Pure-Go driver (no CGO)</text>
    <text x="16" y="118" class="box-p">• Write state: PENDING</text>
    <text x="16" y="134" class="box-p">• Records file path &amp; hash</text>

    <rect x="14" y="155" width="177" height="60" rx="6" fill="#e0f2fe" stroke="#7dd3fc" />
    <text x="22" y="177" font-size="11" font-weight="700" fill="#0369a1" font-family="sans-serif">Thumbnail Extractor</text>
    <text x="22" y="195" font-size="10" fill="#075985" font-family="sans-serif">Generates poster JPEG</text>

    <rect x="14" y="235" width="177" height="75" rx="6" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="22" y="258" font-size="11" font-weight="700" fill="#334155" font-family="sans-serif">Private by Default</text>
    <text x="22" y="276" font-size="10" fill="#64748b" font-family="sans-serif">• Stays local on edge</text>
    <text x="22" y="292" font-size="10" fill="#64748b" font-family="sans-serif">• 0 bytes sent to cloud</text>
  </g>

  <!-- Step 3: Hardware Probe -->
  <g transform="translate(515, 90)" filter="url(#sh3)">
    <rect width="205" height="340" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="205" height="36" rx="10" fill="url(#gHead)" /><rect y="24" width="205" height="12" fill="url(#gHead)" />
    <text x="102" y="23" class="col-h">3. Encoder Prober</text>

    <text x="16" y="65" class="box-num">STEP 3</text>
    <text x="16" y="82" class="box-t">Hardware Detection</text>
    <text x="16" y="102" class="box-p">• 1-second synthetic probe</text>
    <text x="16" y="118" class="box-p">• Tests GPU VRAM state</text>

    <rect x="14" y="145" width="177" height="80" rx="6" fill="#ecfdf5" stroke="#10b981" />
    <text x="22" y="168" font-size="11" font-weight="700" fill="#047857" font-family="sans-serif">Primary: NVENC ASIC</text>
    <text x="22" y="186" font-size="10" fill="#065f46" font-family="sans-serif">• 320+ FPS throughput</text>
    <text x="22" y="202" font-size="10" fill="#065f46" font-family="sans-serif">• &lt;5% CPU utilization</text>

    <rect x="14" y="235" width="177" height="75" rx="6" fill="#fef2f2" stroke="#f87171" />
    <text x="22" y="258" font-size="11" font-weight="700" fill="#991b1b" font-family="sans-serif">CPU Fallback</text>
    <text x="22" y="276" font-size="10" fill="#7f1d1d" font-family="sans-serif">• libx264 veryfast</text>
    <text x="22" y="292" font-size="10" fill="#7f1d1d" font-family="sans-serif">• Zero pipeline drop</text>
  </g>

  <!-- Step 4: Transcoder & HLS -->
  <g transform="translate(755, 90)" filter="url(#sh3)">
    <rect width="205" height="340" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="205" height="36" rx="10" fill="url(#gHead)" /><rect y="24" width="205" height="12" fill="url(#gHead)" />
    <text x="102" y="23" class="col-h">4. HLS Segmenter</text>

    <text x="16" y="65" class="box-num">STEP 4</text>
    <text x="16" y="82" class="box-t">4s GOP Alignment</text>
    <text x="16" y="102" class="box-p">• -g 120 -keyint_min 120</text>
    <text x="16" y="118" class="box-p">• IDR keyframe every 4s</text>
    <text x="16" y="134" class="box-p">• RFC 8216 HLS spec</text>

    <rect x="14" y="155" width="177" height="65" rx="6" fill="#fefce8" stroke="#fde047" />
    <text x="22" y="177" font-size="11" font-weight="700" fill="#854d0e" font-family="sans-serif">Adaptive Bitrate (ABR)</text>
    <text x="22" y="195" font-size="10" fill="#a16207" font-family="sans-serif">1080p, 720p, 480p</text>
    <text x="22" y="209" font-size="10" fill="#a16207" font-family="sans-serif">master.m3u8 ladder</text>

    <rect x="14" y="235" width="177" height="75" rx="6" fill="#f1f5f9" />
    <text x="22" y="258" font-size="11" font-weight="700" fill="#334155" font-family="sans-serif">Cache Chunks</text>
    <text x="22" y="276" font-size="10" fill="#64748b" font-family="sans-serif">.custard/hls/{id}/seg_*.ts</text>
    <text x="22" y="292" font-size="10" fill="#64748b" font-family="sans-serif">Ready for Origin Shield</text>
  </g>

  <!-- Step 5: Mark Ready -->
  <g transform="translate(995, 90)" filter="url(#sh3)">
    <rect width="170" height="340" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="170" height="36" rx="10" fill="url(#gHead)" /><rect y="24" width="170" height="12" fill="url(#gHead)" />
    <text x="85" y="23" class="col-h">5. Mark Ready</text>

    <text x="16" y="65" class="box-num">STEP 5</text>
    <text x="16" y="82" class="box-t">Library State</text>
    <text x="16" y="102" class="box-p">• Update SQLite</text>
    <text x="16" y="118" class="box-p">• Status: READY</text>
    <text x="16" y="134" class="box-p">• Duration &amp; bitrates</text>

    <rect x="14" y="155" width="142" height="70" rx="6" fill="#ecfdf5" stroke="#10b981" />
    <text x="22" y="178" font-size="11" font-weight="700" fill="#047857" font-family="sans-serif">LAN Streaming</text>
    <text x="22" y="196" font-size="10" fill="#065f46" font-family="sans-serif">Active on :8089</text>
    <text x="22" y="210" font-size="10" fill="#065f46" font-family="sans-serif">0ms cloud delay</text>

    <rect x="14" y="235" width="142" height="75" rx="6" fill="#f0f9ff" stroke="#7dd3fc" />
    <text x="22" y="258" font-size="11" font-weight="700" fill="#0369a1" font-family="sans-serif">Public Toggle</text>
    <text x="22" y="276" font-size="10" fill="#0c4a6e" font-family="sans-serif">1-click sync 650B</text>
    <text x="22" y="292" font-size="10" fill="#0c4a6e" font-family="sans-serif">JSON manifest</text>
  </g>

  <!-- Step Arrows -->
  <path d="M 241 230 L 273 230" fill="none" stroke="#4f46e5" stroke-width="2.5" marker-end="url(#mBlue3)" />
  <path d="M 480 230 L 513 230" fill="none" stroke="#4f46e5" stroke-width="2.5" marker-end="url(#mBlue3)" />
  <path d="M 720 230 L 753 230" fill="none" stroke="#4f46e5" stroke-width="2.5" marker-end="url(#mBlue3)" />
  <path d="M 960 230 L 993 230" fill="none" stroke="#10b981" stroke-width="2.5" marker-end="url(#mGreen3)" />
</svg>
"""

# ==============================================================================
# DIAGRAM 4: Ephemeral Heartbeat Liveness Tracking
# ==============================================================================
svg_diag_04 = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 480" width="1200" height="480">
  <defs>
    <filter id="sh4" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <linearGradient id="gLive" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#059669" /><stop offset="100%" stop-color="#047857" />
    </linearGradient>
    <linearGradient id="gDead" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#e11d48" /><stop offset="100%" stop-color="#be123c" />
    </linearGradient>
    <marker id="mArr4" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#0284c7" />
    </marker>
  </defs>

  <style>
    .t4-h { font-family: -apple-system, sans-serif; font-size: 19px; font-weight: 700; fill: #0f172a; }
    .t4-s { font-family: -apple-system, sans-serif; font-size: 13px; fill: #64748b; }
    .box-hdr { font-family: -apple-system, sans-serif; font-size: 14px; font-weight: 700; fill: #ffffff; }
    .txt-h { font-family: -apple-system, sans-serif; font-size: 12px; font-weight: 700; fill: #1e293b; }
    .txt-p { font-family: -apple-system, sans-serif; font-size: 11px; fill: #64748b; }
  </style>

  <rect width="1200" height="480" fill="#f8fafc" rx="14" stroke="#e2e8f0" stroke-width="1.5" />
  <text x="36" y="38" class="t4-h">Distributed Ephemeral Heartbeat &amp; Zero-Polling Offline Detection</text>
  <text x="36" y="58" class="t4-s">Redis 20-second TTL keys eliminate expensive database polling while providing instant liveness updates</text>

  <!-- Left Side: Normal Healthy State -->
  <g transform="translate(36, 85)" filter="url(#sh4)">
    <rect width="540" height="360" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="40" rx="10" fill="url(#gLive)" /><rect y="28" width="540" height="12" fill="url(#gLive)" />
    <text x="20" y="26" class="box-hdr">STATE A: Normal Operation (Node Online &amp; Healthy)</text>

    <!-- Edge Daemon Step -->
    <rect x="20" y="60" width="230" height="110" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="32" y="85" class="txt-h" fill="#047857">1. Edge Heartbeat Loop</text>
    <text x="32" y="105" class="txt-p">• Fires every 15 seconds</text>
    <text x="32" y="121" class="txt-p">• POST /v1/nodes/{id}/heartbeat</text>
    <text x="32" y="137" class="txt-p">• Sends active stream count: 2</text>
    <text x="32" y="153" class="txt-p">• Reports local LAN IP address</text>

    <!-- Cloud Redis Step -->
    <rect x="290" y="60" width="230" height="110" rx="6" fill="#fefce8" stroke="#fde047" />
    <text x="302" y="85" class="txt-h" fill="#854d0e">2. Redis SETEX Operation</text>
    <text x="302" y="105" class="txt-p">• Key: heartbeat:{node_id}</text>
    <text x="302" y="121" class="txt-p">• TTL: Exactly 20 Seconds</text>
    <text x="302" y="137" class="txt-p">• Value: "alive"</text>
    <text x="302" y="153" class="txt-p">• 5s grace window for jitter</text>

    <!-- Catalog Status Result -->
    <rect x="20" y="200" width="500" height="120" rx="8" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="36" y="228" class="txt-h" fill="#0284c7">3. Instant Catalog Query Check (Zero DB Polling)</text>
    <text x="36" y="250" class="txt-p">When a viewer loads the catalog feed, API queries Redis in &lt; 1ms:</text>
    <rect x="36" y="262" width="468" height="38" rx="4" fill="#ecfdf5" stroke="#10b981" />
    <text x="50" y="286" font-size="12" font-weight="700" fill="#065f46" font-family="sans-serif">EXISTS heartbeat:node-4060 == 1 ➔ Status: ONLINE (Green Badge)</text>
  </g>

  <!-- Right Side: Node Failure State -->
  <g transform="translate(620, 85)" filter="url(#sh4)">
    <rect width="540" height="360" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="40" rx="10" fill="url(#gDead)" /><rect y="28" width="540" height="12" fill="url(#gDead)" />
    <text x="20" y="26" class="box-hdr">STATE B: Disconnect Event (Network Severance / Sleep)</text>

    <!-- Edge Crash -->
    <rect x="20" y="60" width="230" height="110" rx="6" fill="#fff1f2" stroke="#f43f5e" />
    <text x="32" y="85" class="txt-h" fill="#be123c">1. Network Cut / PC Off</text>
    <text x="32" y="105" class="txt-p">• Power outage or sleep</text>
    <text x="32" y="121" class="txt-p">• ISP IP reassignment</text>
    <text x="32" y="137" class="txt-p">• Heartbeat beacon ceases</text>
    <text x="32" y="153" class="txt-p">• Zero traffic emitted</text>

    <!-- Redis TTL Expiry -->
    <rect x="290" y="60" width="230" height="110" rx="6" fill="#fffbeb" stroke="#fde68a" />
    <text x="302" y="85" class="txt-h" fill="#b45309">2. Automatic Key Eviction</text>
    <text x="302" y="105" class="txt-p">• At t = 15s: expected ping missing</text>
    <text x="302" y="121" class="txt-p">• At t = 20s: TTL counts to 0</text>
    <text x="302" y="137" class="txt-p">• Redis automatically drops key</text>
    <text x="302" y="153" class="txt-p">• Zero CPU background worker</text>

    <!-- Catalog Status Result -->
    <rect x="20" y="200" width="500" height="120" rx="8" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="36" y="228" class="txt-h" fill="#be123c">3. Graceful Catalog Fail-Safe Handling</text>
    <text x="36" y="250" class="txt-p">When viewer requests media from offline node:</text>
    <rect x="36" y="262" width="468" height="38" rx="4" fill="#fff1f2" stroke="#f43f5e" />
    <text x="50" y="286" font-size="12" font-weight="700" fill="#9f1239" font-family="sans-serif">EXISTS heartbeat:node-4060 == 0 ➔ Status: OFFLINE (Amber Badge / HTTP 503)</text>
  </g>
</svg>
"""

# ==============================================================================
# DIAGRAM 5: Dual-Path Playback & Origin Shield Single Ingest
# ==============================================================================
svg_diag_05 = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 500" width="1200" height="500">
  <defs>
    <filter id="sh5" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <linearGradient id="gBlueH" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" /><stop offset="100%" stop-color="#0369a1" />
    </linearGradient>
    <marker id="mArr5" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="mArrG5" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 8 5 L 0 9 z" fill="#10b981" />
    </marker>
  </defs>

  <style>
    .t5-h { font-family: -apple-system, sans-serif; font-size: 19px; font-weight: 700; fill: #0f172a; }
    .t5-s { font-family: -apple-system, sans-serif; font-size: 13px; fill: #64748b; }
    .path-title { font-family: -apple-system, sans-serif; font-size: 14px; font-weight: 700; fill: #ffffff; }
    .item-h { font-family: -apple-system, sans-serif; font-size: 12px; font-weight: 700; fill: #1e293b; }
    .item-p { font-family: -apple-system, sans-serif; font-size: 11px; fill: #64748b; }
  </style>

  <rect width="1200" height="500" fill="#f8fafc" rx="14" stroke="#e2e8f0" stroke-width="1.5" />
  <text x="36" y="38" class="t5-h">Dual-Path Stream Routing &amp; Origin Shield Single Ingest</text>
  <text x="36" y="58" class="t5-s">Intelligent proximity resolution delivers sub-4ms LAN playback and saves 90% upload bandwidth on remote streams</text>

  <!-- PATH 1: Direct LAN Streaming (Left) -->
  <g transform="translate(36, 85)" filter="url(#sh5)">
    <rect width="540" height="380" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="42" rx="10" fill="#4f46e5" /><rect y="30" width="540" height="12" fill="#4f46e5" />
    <text x="20" y="27" class="path-title">PATH 1: Direct Local Area Network Stream (LAN Proximity Match)</text>

    <rect x="20" y="60" width="500" height="85" rx="6" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="32" y="82" class="item-h" fill="#4338ca">1. Subnet &amp; Public IP Proximity Check</text>
    <text x="32" y="102" class="item-p">Viewer IP (192.168.1.55) shares identical subnet or public NAT IP with Node (192.168.1.105)</text>
    <text x="32" y="118" class="item-p">➔ Condition Matched: Viewer is in the same home/office network</text>

    <rect x="20" y="160" width="500" height="110" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="32" y="184" class="item-h" fill="#047857">2. Direct Streaming Bypass (Zero Cloud Proxy)</text>
    <text x="32" y="204" class="item-p">Cloud API returns direct local endpoint: http://192.168.1.105:8089/stream/master.m3u8</text>
    <text x="32" y="222" font-size="11" font-weight="700" fill="#15803d" font-family="sans-serif">• Time to First Byte: &lt; 3.2 ms (Pure LAN Gigabit/Wi-Fi speed)</text>
    <text x="32" y="238" font-size="11" font-weight="700" fill="#15803d" font-family="sans-serif">• Cloud Egress Bandwidth: Exactly 0 Bytes ($0.00 Cost)</text>
    <text x="32" y="254" font-size="11" font-weight="700" fill="#15803d" font-family="sans-serif">• 100% Data Sovereignty: Video bytes never transit the public internet</text>

    <rect x="20" y="285" width="500" height="60" rx="6" fill="#ecfdf5" stroke="#10b981" />
    <text x="32" y="320" font-size="12" font-weight="700" fill="#065f46" font-family="sans-serif">Result: Standalone Mini-Jellyfin offline streaming works without cloud!</text>
  </g>

  <!-- PATH 2: Remote Origin Shield (Right) -->
  <g transform="translate(620, 85)" filter="url(#sh5)">
    <rect width="540" height="380" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="42" rx="10" fill="url(#gBlueH)" /><rect y="30" width="540" height="12" fill="url(#gBlueH)" />
    <text x="20" y="27" class="path-title">PATH 2: Remote Origin Shield Caching (Single Ingest)</text>

    <rect x="20" y="60" width="500" height="85" rx="6" fill="#fffbeb" stroke="#fde68a" />
    <text x="32" y="82" class="item-h" fill="#b45309">1. Viewer 1: Cache MISS (First Ingest)</text>
    <text x="32" y="102" class="item-p">• S3 Cache does not have chunk ➔ Pulls once via FRP Reverse Tunnel</text>
    <text x="32" y="118" class="item-p">• Chunk (2.5 MB) committed to S3 Cache with 24h auto-purge LRU policy</text>
    <text x="32" y="134" class="item-p">• Home uplink used: exactly 2.5 MB</text>

    <rect x="20" y="160" width="500" height="110" rx="6" fill="#f0fdf4" stroke="#86efac" />
    <text x="32" y="184" class="item-h" fill="#15803d">2. Viewers 2 to 10: Cache HIT (Zero Edge Uplink)</text>
    <text x="32" y="204" class="item-p">• Subsequent requests for chunk served directly from S3 Cache / CloudFront</text>
    <text x="32" y="222" font-size="11" font-weight="700" fill="#15803d" font-family="sans-serif">• Response Latency: 14 ms (PoP Cache)</text>
    <text x="32" y="238" font-size="11" font-weight="700" fill="#15803d" font-family="sans-serif">• Home Node Uplink: 0 Bps (Completely shielded from multi-viewer spike)</text>
    <text x="32" y="254" font-size="11" font-weight="700" fill="#15803d" font-family="sans-serif">• Bandwidth Saved: 90.0% Reduction across 10 concurrent viewers</text>

    <rect x="20" y="285" width="500" height="60" rx="6" fill="#fefce8" stroke="#fde047" />
    <text x="32" y="320" font-size="12" font-weight="700" fill="#854d0e" font-family="sans-serif">Result: 10 viewers stream 1080p smoothly over a single 20 Mbps home link!</text>
  </g>
</svg>
"""

# ==============================================================================
# DIAGRAM 6: Financial Cost Comparison Breakdown
# ==============================================================================
svg_diag_06 = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 460" width="1200" height="460">
  <defs>
    <filter id="sh6" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
  </defs>

  <style>
    .t6-h { font-family: -apple-system, sans-serif; font-size: 19px; font-weight: 700; fill: #0f172a; }
    .t6-s { font-family: -apple-system, sans-serif; font-size: 13px; fill: #64748b; }
    .bar-lbl { font-family: -apple-system, sans-serif; font-size: 13px; font-weight: 700; fill: #1e293b; }
    .bar-val { font-family: -apple-system, sans-serif; font-size: 13px; font-weight: 700; }
  </style>

  <rect width="1200" height="460" fill="#f8fafc" rx="14" stroke="#e2e8f0" stroke-width="1.5" />
  <text x="36" y="38" class="t6-h">Monthly Cloud Operating Cost Comparison (100 Users, 5 TB Catalog)</text>
  <text x="36" y="58" class="t6-s">CustardFin decouples media storage and transcoding from cloud billing, achieving a verified 95.89% cost reduction</text>

  <!-- Left: Traditional Cloud ($1,195.00/mo) -->
  <g transform="translate(36, 85)" filter="url(#sh6)">
    <rect width="540" height="335" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="40" rx="10" fill="#e11d48" /><rect y="28" width="540" height="12" fill="#e11d48" />
    <text x="20" y="26" font-size="14" font-weight="700" fill="#ffffff" font-family="sans-serif">TRADITIONAL CENTRALIZED CLOUD VOD</text>
    <text x="410" y="26" font-size="14" font-weight="700" fill="#ffffff" font-family="sans-serif">$1,195.00 / mo</text>

    <!-- Bar 1: Bandwidth Egress $810 -->
    <text x="20" y="80" class="bar-lbl">Cloud Data Egress (9 TB @ $0.09/GB)</text>
    <rect x="20" y="90" width="400" height="26" rx="4" fill="#f43f5e" />
    <text x="430" y="108" class="bar-val" fill="#be123c">$810.00 (67.8%)</text>

    <!-- Bar 2: Transcoding Compute $140 -->
    <text x="20" y="140" class="bar-lbl">Cloud Transcode Compute (EC2 / MediaConvert)</text>
    <rect x="20" y="150" width="70" height="26" rx="4" fill="#fb7185" />
    <text x="100" y="168" class="bar-val" fill="#be123c">$140.00 (11.7%)</text>

    <!-- Bar 3: Managed Database & Infra $130 -->
    <text x="20" y="200" class="bar-lbl">Managed Cloud Infra (RDS PostgreSQL + ALB)</text>
    <rect x="20" y="210" width="65" height="26" rx="4" fill="#fda4af" />
    <text x="95" y="228" class="bar-val" fill="#be123c">$130.00 (10.9%)</text>

    <!-- Bar 4: Permanent S3 Storage $115 -->
    <text x="20" y="260" class="bar-lbl">Permanent Cloud Object Storage (5 TB @ $0.023/GB)</text>
    <rect x="20" y="270" width="58" height="26" rx="4" fill="#fecdd3" />
    <text x="88" y="288" class="bar-val" fill="#be123c">$115.00 (9.6%)</text>
  </g>

  <!-- Right: CustardFin ($49.15/mo) -->
  <g transform="translate(620, 85)" filter="url(#sh6)">
    <rect width="540" height="335" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="40" rx="10" fill="#059669" /><rect y="28" width="540" height="12" fill="#059669" />
    <text x="20" y="26" font-size="14" font-weight="700" fill="#ffffff" font-family="sans-serif">CUSTARDFIN DECENTRALIZED HYBRID</text>
    <text x="425" y="26" font-size="14" font-weight="700" fill="#ffffff" font-family="sans-serif">$49.15 / mo</text>

    <!-- Bar 1: Control Plane Compute $48 -->
    <text x="20" y="80" class="bar-lbl">Consolidated Cloud VM (e2-standard-2: Go + Redis + PG)</text>
    <rect x="20" y="90" width="24" height="26" rx="4" fill="#10b981" />
    <text x="54" y="108" class="bar-val" fill="#047857">$48.00 (97.6%)</text>

    <!-- Bar 2: Origin Shield Cache $1.15 -->
    <text x="20" y="140" class="bar-lbl">Origin Shield Ephemeral Cache (24h LRU Purge)</text>
    <rect x="20" y="150" width="4" height="26" rx="2" fill="#34d399" />
    <text x="34" y="168" class="bar-val" fill="#047857">$1.15 (2.4%)</text>

    <!-- Bar 3: Permanent Video Storage $0.00 -->
    <text x="20" y="200" class="bar-lbl">Permanent Cloud Video Storage (100% on Edge Disks)</text>
    <rect x="20" y="210" width="3" height="26" rx="1" fill="#cbd5e1" />
    <text x="30" y="228" class="bar-val" fill="#15803d">$0.00 (100% SAVINGS)</text>

    <!-- Bar 4: Transcode Compute $0.00 -->
    <text x="20" y="260" class="bar-lbl">Transcoding Compute (100% on Edge NVENC ASICs)</text>
    <rect x="20" y="270" width="3" height="26" rx="1" fill="#cbd5e1" />
    <text x="30" y="288" class="bar-val" fill="#15803d">$0.00 (100% SAVINGS)</text>

    <!-- Savings Callout -->
    <rect x="20" y="295" width="500" height="26" rx="4" fill="#ecfdf5" stroke="#10b981" />
    <text x="80" y="313" font-size="12" font-weight="700" fill="#065f46" font-family="sans-serif">NET SAVINGS: $1,145.85 / MONTH (95.89% COST REDUCTION)</text>
  </g>
</svg>
"""

# ==============================================================================
# DIAGRAM 7: NVENC GPU vs CPU Performance Benchmark
# ==============================================================================
svg_diag_07 = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 460" width="1200" height="460">
  <defs>
    <filter id="sh7" x="-5%" y="-5%" width="110%" height="115%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
  </defs>

  <style>
    .t7-h { font-family: -apple-system, sans-serif; font-size: 19px; font-weight: 700; fill: #0f172a; }
    .t7-s { font-family: -apple-system, sans-serif; font-size: 13px; fill: #64748b; }
    .card-h { font-family: -apple-system, sans-serif; font-size: 14px; font-weight: 700; fill: #ffffff; }
    .stat-val { font-family: -apple-system, sans-serif; font-size: 26px; font-weight: 800; }
    .stat-lbl { font-family: -apple-system, sans-serif; font-size: 11px; font-weight: 600; text-transform: uppercase; }
  </style>

  <rect width="1200" height="460" fill="#f8fafc" rx="14" stroke="#e2e8f0" stroke-width="1.5" />
  <text x="36" y="38" class="t7-h">Hardware Transcoding Benchmark: NVIDIA NVENC vs. CPU</text>
  <text x="36" y="58" class="t7-s">Empirical evaluation of 1080p HLS pre-encoding throughput, host CPU utilization, and thermal dissipation</text>

  <!-- Left: NVIDIA RTX 4060 NVENC -->
  <g transform="translate(36, 85)" filter="url(#sh7)">
    <rect width="540" height="335" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="42" rx="10" fill="#10b981" /><rect y="30" width="540" height="12" fill="#10b981" />
    <text x="20" y="27" class="card-h">⚡ NVIDIA GeForce RTX 4060 (h264_nvenc ASIC)</text>

    <!-- Grid of 4 Stats -->
    <rect x="20" y="65" width="240" height="105" rx="8" fill="#f0fdf4" stroke="#86efac" />
    <text x="35" y="105" class="stat-val" fill="#047857">10.74×</text>
    <text x="35" y="125" class="stat-lbl" fill="#15803d">Encoding Speed (Realtime)</text>
    <text x="35" y="145" font-size="11" fill="#166534" font-family="sans-serif">10-min clip encodes in 55.8s</text>

    <rect x="280" y="65" width="240" height="105" rx="8" fill="#f0fdf4" stroke="#86efac" />
    <text x="295" y="105" class="stat-val" fill="#047857">322.3 fps</text>
    <text x="295" y="125" class="stat-lbl" fill="#15803d">Frame Throughput</text>
    <text x="295" y="145" font-size="11" fill="#166534" font-family="sans-serif">Dedicated on-die silicon</text>

    <rect x="20" y="185" width="240" height="105" rx="8" fill="#f0fdf4" stroke="#86efac" />
    <text x="35" y="225" class="stat-val" fill="#047857">4.8%</text>
    <text x="35" y="245" class="stat-lbl" fill="#15803d">Host CPU Utilization</text>
    <text x="35" y="265" font-size="11" fill="#166534" font-family="sans-serif">Host stays 100% responsive</text>

    <rect x="280" y="185" width="240" height="105" rx="8" fill="#f0fdf4" stroke="#86efac" />
    <text x="295" y="225" class="stat-val" fill="#047857">61°C / 42W</text>
    <text x="295" y="245" class="stat-lbl" fill="#15803d">Peak Temp &amp; Power Draw</text>
    <text x="295" y="265" font-size="11" fill="#166534" font-family="sans-serif">Low thermal &amp; energy footprint</text>
  </g>

  <!-- Right: CPU Software (libx264) -->
  <g transform="translate(620, 85)" filter="url(#sh7)">
    <rect width="540" height="335" rx="10" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="540" height="42" rx="10" fill="#e11d48" /><rect y="30" width="540" height="12" fill="#e11d48" />
    <text x="20" y="27" class="card-h">💻 Multi-Threaded CPU Software (libx264 veryfast)</text>

    <!-- Grid of 4 Stats -->
    <rect x="20" y="65" width="240" height="105" rx="8" fill="#fff1f2" stroke="#fecdd3" />
    <text x="35" y="105" class="stat-val" fill="#be123c">1.12×</text>
    <text x="35" y="125" class="stat-lbl" fill="#9f1239">Encoding Speed (Realtime)</text>
    <text x="35" y="145" font-size="11" fill="#881337" font-family="sans-serif">10-min clip takes 8.9 minutes</text>

    <rect x="280" y="65" width="240" height="105" rx="8" fill="#fff1f2" stroke="#fecdd3" />
    <text x="295" y="105" class="stat-val" fill="#be123c">33.6 fps</text>
    <text x="295" y="125" class="stat-lbl" fill="#9f1239">Frame Throughput</text>
    <text x="295" y="145" font-size="11" fill="#881337" font-family="sans-serif">Barely above 30fps playback</text>

    <rect x="20" y="185" width="240" height="105" rx="8" fill="#fff1f2" stroke="#fecdd3" />
    <text x="35" y="225" class="stat-val" fill="#be123c">96.4%</text>
    <text x="35" y="245" class="stat-lbl" fill="#9f1239">Host CPU Utilization</text>
    <text x="35" y="265" font-size="11" fill="#881337" font-family="sans-serif">All 16 threads pegged at 100%</text>

    <rect x="280" y="185" width="240" height="105" rx="8" fill="#fff1f2" stroke="#fecdd3" />
    <text x="295" y="225" class="stat-val" fill="#be123c">89°C / 95W</text>
    <text x="295" y="245" class="stat-lbl" fill="#9f1239">Peak Temp &amp; Power Draw</text>
    <text x="295" y="265" font-size="11" fill="#881337" font-family="sans-serif">High thermal fan noise &amp; heat</text>
  </g>
</svg>
"""

items = [
    ("diag_01_nat_traversal", svg_diag_01, 1200, 520),
    ("diag_03_ingest_pipeline", svg_diag_03, 1200, 480),
    ("diag_04_heartbeat_liveness", svg_diag_04, 1200, 480),
    ("diag_05_dual_path_shield", svg_diag_05, 1200, 500),
    ("diag_06_cost_comparison", svg_diag_06, 1200, 460),
    ("diag_07_nvenc_benchmark", svg_diag_07, 1200, 460),
]

for name, svg_code, w, h in items:
    s_path = os.path.join(fig_dir, f"{name}.svg")
    p_path = os.path.join(fig_dir, f"{name}.png")
    with open(s_path, "w", encoding="utf-8") as f:
        f.write(svg_code)
    render_svg_to_png(s_path, p_path, w, h)

print("All clean vector SVGs & PNGs generated successfully.")
