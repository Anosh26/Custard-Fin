import xml.etree.ElementTree as ET
import os
import subprocess
import pathlib

svg_raw = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 480" width="1200" height="480" style="background-color: #ffffff; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
  <defs>
    <!-- Soft Drop Shadows -->
    <filter id="shadow" x="-5%" y="-5%" width="110%" height="110%">
      <feDropShadow dx="0" dy="4" stdDeviation="5" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <filter id="shadow-sm" x="-5%" y="-5%" width="110%" height="110%">
      <feDropShadow dx="0" dy="2" stdDeviation="3" flood-color="#0f172a" flood-opacity="0.05" />
    </filter>
    <filter id="glow-green" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="3" result="blur" />
      <feComposite in="SourceGraphic" in2="blur" operator="over" />
    </filter>

    <!-- Transition Gradients for Connectors -->
    <linearGradient id="grad-1-2" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%" stop-color="#64748b" />
      <stop offset="100%" stop-color="#10b981" />
    </linearGradient>
    <linearGradient id="grad-2-3" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%" stop-color="#10b981" />
      <stop offset="100%" stop-color="#ef4444" />
    </linearGradient>
    <linearGradient id="grad-3-4" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%" stop-color="#ef4444" />
      <stop offset="100%" stop-color="#8b5cf6" />
    </linearGradient>
    <linearGradient id="grad-4-5" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%" stop-color="#8b5cf6" />
      <stop offset="100%" stop-color="#f59e0b" />
    </linearGradient>

    <!-- Arrowhead Markers for Gradients -->
    <marker id="arrow-green" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#10b981" />
    </marker>
    <marker id="arrow-red" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#ef4444" />
    </marker>
    <marker id="arrow-purple" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#8b5cf6" />
    </marker>
    <marker id="arrow-amber" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#f59e0b" />
    </marker>
  </defs>

  <!-- ================= GLOBAL HEADER ================= -->
  <text x="600" y="38" text-anchor="middle" font-size="20" font-weight="800" fill="#0f172a" letter-spacing="-0.3">
    Edge Media Ingestion, HLS Segment Slicing, and SQLite Metadata Indexing
  </text>
  <text x="600" y="58" text-anchor="middle" font-size="11.5" font-weight="500" fill="#64748b">
    Automated scanner discovers raw video, probes NVENC GPU, slices 4-second GOP-aligned HLS chunks, and syncs lightweight manifest
  </text>

  <!-- ================= HORIZONTAL PIPELINE CONNECTORS ================= -->
  <!-- Drawn first so they sit under the node circles -->
  <!-- 1 -> 2 -->
  <line x1="160" y1="140" x2="335" y2="140" stroke="url(#grad-1-2)" stroke-width="4" marker-end="url(#arrow-green)" />
  <!-- 2 -> 3 -->
  <line x1="400" y1="140" x2="575" y2="140" stroke="url(#grad-2-3)" stroke-width="4" marker-end="url(#arrow-red)" />
  <!-- 3 -> 4 -->
  <line x1="640" y1="140" x2="815" y2="140" stroke="url(#grad-3-4)" stroke-width="4" marker-end="url(#arrow-purple)" />
  <!-- 4 -> 5 -->
  <line x1="880" y1="140" x2="1055" y2="140" stroke="url(#grad-4-5)" stroke-width="4" marker-end="url(#arrow-amber)" />

  <!-- ================= STATION 1: RAW MEDIA DISCOVERY (Slate / Blue) ================= -->
  <!-- Center: X = 120 -->
  <g transform="translate(12, 0)">
    <!-- Node Icon -->
    <circle cx="108" cy="140" r="36" fill="#f8fafc" stroke="#475569" stroke-width="2.5" filter="url(#shadow-sm)" />
    <!-- Step Badge -->
    <circle cx="108" cy="104" r="12" fill="#475569" />
    <text x="108" y="108.5" text-anchor="middle" font-size="11" font-weight="800" fill="#ffffff">1</text>
    <!-- Video + Mag Glass Icon -->
    <rect x="94" y="124" width="22" height="28" rx="2" fill="#94a3b8" />
    <polygon points="98,128 112,128 112,148 98,148" fill="#e2e8f0" />
    <circle cx="118" cy="144" r="7" fill="none" stroke="#0f172a" stroke-width="2" />
    <line x1="123" y1="149" x2="128" y2="154" stroke="#0f172a" stroke-width="2.5" stroke-linecap="round" />

    <!-- Detail Box -->
    <rect x="0" y="210" width="216" height="180" rx="8" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.5" filter="url(#shadow)" />
    <path d="M 0 218 C 0 213.58 3.58 210 8 210 L 208 210 C 212.42 210 216 213.58 216 218 L 216 254 L 0 254 Z" fill="#f1f5f9" />
    
    <text x="108" y="231" text-anchor="middle" font-size="11" font-weight="800" fill="#0f172a">Raw Media Discovery</text>
    <text x="108" y="246" text-anchor="middle" font-size="11" font-weight="800" fill="#0f172a">&amp; Hashing</text>
    
    <circle cx="16" cy="276" r="2.5" fill="#475569" />
    <text x="26" y="280" font-size="9.5" fill="#334155">Local watch: /media/*.mp4</text>
    
    <circle cx="16" cy="302" r="2.5" fill="#475569" />
    <text x="26" y="306" font-size="9.5" fill="#334155">Cryptographic SHA-256</text>
    <text x="26" y="318" font-size="9.5" fill="#334155">binary hash generation</text>
    
    <circle cx="16" cy="344" r="2.5" fill="#475569" />
    <text x="26" y="348" font-size="9.5" fill="#334155">Duplicate elimination &amp;</text>
    <text x="26" y="360" font-size="9.5" fill="#334155">strict identity binding</text>
  </g>

  <!-- ================= STATION 2: HARDWARE PROBE (Emerald) ================= -->
  <!-- Center: X = 360 -->
  <g transform="translate(252, 0)">
    <!-- Node Icon -->
    <circle cx="108" cy="140" r="36" fill="#ecfdf5" stroke="#10b981" stroke-width="2.5" filter="url(#shadow-sm)" />
    <!-- Step Badge -->
    <circle cx="108" cy="104" r="12" fill="#10b981" />
    <text x="108" y="108.5" text-anchor="middle" font-size="11" font-weight="800" fill="#ffffff">2</text>
    <!-- NVENC / CPU Icon -->
    <rect x="94" y="128" width="28" height="28" rx="3" fill="#047857" stroke="#10b981" stroke-width="1.5" />
    <!-- CPU check fallback -->
    <rect x="114" y="122" width="12" height="12" rx="2" fill="#334155" />
    <path d="M 116 128 L 118 130 L 123 124" fill="none" stroke="#10b981" stroke-width="1.5" />
    <!-- Probe needle -->
    <line x1="84" y1="142" x2="94" y2="142" stroke="#0f172a" stroke-width="1.5" stroke-dasharray="2,1" />

    <!-- Detail Box -->
    <rect x="0" y="210" width="216" height="180" rx="8" fill="#ffffff" stroke="#a7f3d0" stroke-width="1.5" filter="url(#shadow)" />
    <path d="M 0 218 C 0 213.58 3.58 210 8 210 L 208 210 C 212.42 210 216 213.58 216 218 L 216 254 L 0 254 Z" fill="#ecfdf5" />
    
    <text x="108" y="231" text-anchor="middle" font-size="11" font-weight="800" fill="#064e3b">Hardware Capabilities</text>
    <text x="108" y="246" text-anchor="middle" font-size="11" font-weight="800" fill="#064e3b">Probe &amp; Fallback</text>
    
    <circle cx="16" cy="276" r="2.5" fill="#10b981" />
    <text x="26" y="280" font-family="monospace" font-size="9" fill="#047857">h264_nvenc → QSV → libx264</text>
    
    <circle cx="16" cy="302" r="2.5" fill="#10b981" />
    <text x="26" y="306" font-size="9.5" fill="#334155">1-second canary transcode test</text>
    <text x="26" y="318" font-size="9.5" fill="#334155">for driver / VRAM verification</text>
    
    <circle cx="16" cy="344" r="2.5" fill="#10b981" />
    <text x="26" y="348" font-size="9.5" fill="#334155">Automatic CPU fallback on</text>
    <text x="26" y="360" font-size="9.5" fill="#334155">fault without failing job</text>
  </g>

  <!-- ================= STATION 3: HLS TRANSCODING (Crimson) ================= -->
  <!-- Center: X = 600 -->
  <g transform="translate(492, 0)">
    <!-- Node Icon -->
    <circle cx="108" cy="140" r="36" fill="#fef2f2" stroke="#ef4444" stroke-width="2.5" filter="url(#shadow-sm)" />
    <!-- Step Badge -->
    <circle cx="108" cy="104" r="12" fill="#ef4444" />
    <text x="108" y="108.5" text-anchor="middle" font-size="11" font-weight="800" fill="#ffffff">3</text>
    <!-- Filmstrip & Razor Icon -->
    <rect x="94" y="126" width="28" height="28" rx="2" fill="#7f1d1d" />
    <!-- Film holes -->
    <circle cx="98" cy="130" r="1.5" fill="#fef2f2" /> <circle cx="98" cy="136" r="1.5" fill="#fef2f2" /> <circle cx="98" cy="142" r="1.5" fill="#fef2f2" /> <circle cx="98" cy="148" r="1.5" fill="#fef2f2" />
    <circle cx="118" cy="130" r="1.5" fill="#fef2f2" /> <circle cx="118" cy="136" r="1.5" fill="#fef2f2" /> <circle cx="118" cy="142" r="1.5" fill="#fef2f2" /> <circle cx="118" cy="148" r="1.5" fill="#fef2f2" />
    <!-- Razor Line -->
    <line x1="88" y1="140" x2="128" y2="140" stroke="#fca5a5" stroke-width="2" stroke-dasharray="4,2" />

    <!-- Detail Box -->
    <rect x="0" y="210" width="216" height="180" rx="8" fill="#ffffff" stroke="#fca5a5" stroke-width="1.5" filter="url(#shadow)" />
    <path d="M 0 218 C 0 213.58 3.58 210 8 210 L 208 210 C 212.42 210 216 213.58 216 218 L 216 254 L 0 254 Z" fill="#fef2f2" />
    
    <text x="108" y="231" text-anchor="middle" font-size="11" font-weight="800" fill="#7f1d1d">Subprocess Transcoder</text>
    <text x="108" y="246" text-anchor="middle" font-size="11" font-weight="800" fill="#7f1d1d">(FFmpeg / HLS Slicer)</text>
    
    <circle cx="16" cy="276" r="2.5" fill="#ef4444" />
    <text x="26" y="280" font-family="monospace" font-size="9" fill="#991b1b">-hls_time 4 (4.0s Segments)</text>
    
    <circle cx="16" cy="302" r="2.5" fill="#ef4444" />
    <text x="26" y="306" font-family="monospace" font-size="9" fill="#991b1b">-g 120 -keyint_min 120</text>
    <text x="26" y="318" font-size="9" fill="#334155">(Strict GOP Keyframe align)</text>
    
    <circle cx="16" cy="344" r="2.5" fill="#ef4444" />
    <text x="26" y="348" font-size="9.5" fill="#334155">Multi-bitrate ladder profile:</text>
    <text x="26" y="360" font-size="9.5" font-weight="600" fill="#7f1d1d">1080p, 720p, 480p &amp; master.m3u8</text>
  </g>

  <!-- ================= STATION 4: SQLITE INDEXING (Purple) ================= -->
  <!-- Center: X = 840 -->
  <g transform="translate(732, 0)">
    <!-- Node Icon -->
    <circle cx="108" cy="140" r="36" fill="#f5f3ff" stroke="#8b5cf6" stroke-width="2.5" filter="url(#shadow-sm)" />
    <!-- Step Badge -->
    <circle cx="108" cy="104" r="12" fill="#8b5cf6" />
    <text x="108" y="108.5" text-anchor="middle" font-size="11" font-weight="800" fill="#ffffff">4</text>
    <!-- DB Barrel Icon -->
    <path d="M 94 130 L 94 150 C 94 154 122 154 122 150 L 122 130 Z" fill="#7c3aed" />
    <ellipse cx="108" cy="130" rx="14" ry="4" fill="#a855f7" stroke="#6d28d9" stroke-width="1.2" />
    <rect x="110" y="142" width="22" height="10" rx="2" fill="#f3e8ff" stroke="#6d28d9" stroke-width="1" />
    <text x="121" y="149" text-anchor="middle" font-size="6.5" font-weight="800" fill="#6d28d9">WAL</text>

    <!-- Detail Box -->
    <rect x="0" y="210" width="216" height="180" rx="8" fill="#ffffff" stroke="#c4b5fd" stroke-width="1.5" filter="url(#shadow)" />
    <path d="M 0 218 C 0 213.58 3.58 210 8 210 L 208 210 C 212.42 210 216 213.58 216 218 L 216 254 L 0 254 Z" fill="#f5f3ff" />
    
    <text x="108" y="231" text-anchor="middle" font-size="11" font-weight="800" fill="#4c1d95">Local SQLite Metadata</text>
    <text x="108" y="246" text-anchor="middle" font-size="11" font-weight="800" fill="#4c1d95">Persistence Engine</text>
    
    <circle cx="16" cy="276" r="2.5" fill="#8b5cf6" />
    <text x="26" y="280" font-size="9.5" fill="#334155">WAL journal mode enabled</text>
    
    <circle cx="16" cy="302" r="2.5" fill="#8b5cf6" />
    <text x="26" y="306" font-size="9.5" fill="#334155">Segments, runtimes, bitrates</text>
    <text x="26" y="318" font-size="9.5" fill="#334155">and local paths indexed</text>
    
    <circle cx="16" cy="344" r="2.5" fill="#8b5cf6" />
    <text x="26" y="348" font-size="9.5" fill="#334155">Private watch history kept</text>
    <text x="26" y="360" font-size="9.5" font-weight="700" fill="#6d28d9">strictly on disk (Invariant 8)</text>
  </g>

  <!-- ================= STATION 5: MANIFEST SYNC (Amber) ================= -->
  <!-- Center: X = 1080 -->
  <g transform="translate(972, 0)">
    <!-- Node Icon -->
    <circle cx="108" cy="140" r="36" fill="#fffbeb" stroke="#f59e0b" stroke-width="2.5" filter="url(#shadow-sm)" />
    <!-- Step Badge -->
    <circle cx="108" cy="104" r="12" fill="#f59e0b" />
    <text x="108" y="108.5" text-anchor="middle" font-size="11" font-weight="800" fill="#ffffff">5</text>
    <!-- JSON Doc + Upload Arrow -->
    <rect x="96" y="124" width="24" height="30" rx="2" fill="#fef3c7" stroke="#d97706" stroke-width="1.5" />
    <text x="108" y="136" text-anchor="middle" font-family="monospace" font-size="7" font-weight="700" fill="#b45309">JSON</text>
    <path d="M 108 152 L 108 142 M 104 146 L 108 142 L 112 146" fill="none" stroke="#d97706" stroke-width="1.5" />

    <!-- Detail Box -->
    <rect x="0" y="210" width="216" height="180" rx="8" fill="#ffffff" stroke="#fde68a" stroke-width="1.5" filter="url(#shadow)" />
    <path d="M 0 218 C 0 213.58 3.58 210 8 210 L 208 210 C 212.42 210 216 213.58 216 218 L 216 254 L 0 254 Z" fill="#fffbeb" />
    
    <text x="108" y="231" text-anchor="middle" font-size="11" font-weight="800" fill="#78350f">Public Manifest Gen</text>
    <text x="108" y="246" text-anchor="middle" font-size="11" font-weight="800" fill="#78350f">&amp; Cloud Synchronization</text>
    
    <circle cx="16" cy="276" r="2.5" fill="#f59e0b" />
    <text x="26" y="280" font-size="9.5" fill="#334155">User toggles "Public" in Edge UI</text>
    
    <circle cx="16" cy="302" r="2.5" fill="#f59e0b" />
    <text x="26" y="306" font-size="9.5" fill="#334155">Compact ~650-byte manifest</text>
    <text x="26" y="318" font-family="monospace" font-size="8.5" fill="#b45309">PUT /v1/nodes/media/manifest</text>
    
    <circle cx="16" cy="344" r="2.5" fill="#f59e0b" />
    <text x="26" y="348" font-size="9.5" font-weight="700" fill="#b45309">Video binaries NEVER uploaded</text>
    <text x="26" y="360" font-size="9.5" fill="#78350f">(Invariant 1 strict enforcement)</text>
  </g>

  <!-- ================= BOTTOM DECORATIVE / CLOSING LINE ================= -->
  <line x1="12" y1="440" x2="1188" y2="440" stroke="#e2e8f0" stroke-width="1.5" stroke-dasharray="8,4" />
  <rect x="500" y="430" width="200" height="20" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1" />
  <text x="600" y="444" text-anchor="middle" font-size="9.5" font-weight="600" fill="#64748b">
    Fully Automated Edge Workload
  </text>

</svg>"""

# Parse to verify XML validity
try:
    ET.fromstring(svg_raw)
    print("XML Validation Passed: Valid SVG!")
except Exception as e:
    print("XML Validation Failed:", e)
    raise

svg_file = pathlib.Path("report/figures/diag_03_ingest_pipeline.svg")
html_file = pathlib.Path("report/figures/diag_03_ingest_pipeline.html")
png_file = pathlib.Path("report/figures/diag_03_ingest_pipeline.png")

svg_file.write_text(svg_raw, encoding="utf-8")
print(f"Saved {svg_file}")

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  * {{ margin: 0; padding: 0; box-sizing: border-box; }}
  html, body {{ width: 3600px; height: 1440px; background: #ffffff; overflow: hidden; }}
  svg {{ width: 3600px; height: 1440px; display: block; }}
</style>
</head>
<body>
{svg_raw}
</body>
</html>"""

html_file.write_text(html_content, encoding="utf-8")
print(f"Saved {html_file}")

chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    "--hide-scrollbars",
    "--window-size=3600,1440",
    f"--screenshot={png_file.resolve()}",
    str(html_file.resolve())
]

subprocess.run(args, check=True)
print(f"Rendered {png_file} Ultra-HD -> {png_file.stat().st_size} bytes")
