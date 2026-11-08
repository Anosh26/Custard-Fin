import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(fig_dir, exist_ok=True)
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_content = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 560" width="1200" height="560">
  <defs>
    <!-- Shadow Filters -->
    <filter id="shadowBox" x="-5%" y="-5%" width="110%" height="115%" filterUnits="userSpaceOnUse">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.08" />
    </filter>
    <filter id="glowGreen" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="4" result="blur" />
      <feComposite in="SourceGraphic" in2="blur" operator="over" />
    </filter>
    <filter id="glowBlue" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="4" result="blur" />
      <feComposite in="SourceGraphic" in2="blur" operator="over" />
    </filter>

    <!-- Gradients -->
    <linearGradient id="bgGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#ffffff" />
      <stop offset="100%" stop-color="#f8fafc" />
    </linearGradient>

    <linearGradient id="edgeChassis" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#334155" />
      <stop offset="100%" stop-color="#0f172a" />
    </linearGradient>

    <linearGradient id="routerGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" />
      <stop offset="100%" stop-color="#312e81" />
    </linearGradient>

    <linearGradient id="shieldGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#0284c7" />
      <stop offset="50%" stop-color="#0369a1" />
      <stop offset="100%" stop-color="#075985" />
    </linearGradient>

    <linearGradient id="lanCard" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#ecfdf5" />
      <stop offset="100%" stop-color="#d1fae5" />
    </linearGradient>

    <linearGradient id="remoteCard" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#f0f9ff" />
      <stop offset="100%" stop-color="#e0f2fe" />
    </linearGradient>

    <!-- Markers -->
    <marker id="arrGreen" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#10b981" />
    </marker>
    <marker id="arrBlue" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#0284c7" />
    </marker>
    <marker id="arrAmber" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto">
      <path d="M 0 1 L 9 5 L 0 9 z" fill="#f59e0b" />
    </marker>
  </defs>

  <style>
    text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; }
    .title { font-size: 20px; font-weight: 800; fill: #0f172a; letter-spacing: -0.3px; }
    .subtitle { font-size: 12px; fill: #64748b; }
    .sect-title { font-size: 13px; font-weight: 700; }
    .badge-txt { font-size: 10px; font-weight: 700; }
    .body-txt { font-size: 11px; fill: #334155; }
    .code-txt { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 10px; }
  </style>

  <!-- Background Base -->
  <rect width="1200" height="560" rx="16" fill="url(#bgGrad)" stroke="#e2e8f0" stroke-width="1.5" />

  <!-- Grid Accent Pattern -->
  <g opacity="0.04" stroke="#0f172a" stroke-width="1">
    <line x1="60" y1="0" x2="60" y2="560" />
    <line x1="200" y1="0" x2="200" y2="560" />
    <line x1="380" y1="0" x2="380" y2="560" />
    <line x1="580" y1="0" x2="580" y2="560" />
    <line x1="780" y1="0" x2="780" y2="560" />
    <line x1="980" y1="0" x2="980" y2="560" />
    <line x1="0" y1="70" x2="1200" y2="70" />
    <line x1="0" y1="210" x2="1200" y2="210" />
    <line x1="0" y1="380" x2="1200" y2="380" />
  </g>

  <!-- Header -->
  <g transform="translate(40, 24)">
    <rect width="32" height="32" rx="8" fill="#4f46e5" />
    <path d="M 8 16 L 24 16 M 16 8 L 24 16 L 16 24" fill="none" stroke="#ffffff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" />
    <text x="44" y="20" class="title">Dual-Path Stream Routing &amp; Origin Shield Single Ingest</text>
    <text x="44" y="36" class="subtitle">Dynamic proximity detection routes LAN viewers directly at sub-4ms latency while cloud shield absorbs multi-viewer remote fan-out</text>
  </g>

  <!-- ========================================================================= -->
  <!-- ZONE 1: HOME EDGE SOURCE (LEFT) -->
  <!-- ========================================================================= -->
  <g transform="translate(40, 80)" filter="url(#shadowBox)">
    <!-- Edge Server Box -->
    <rect width="220" height="440" rx="12" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="220" height="34" rx="12" fill="url(#edgeChassis)" />
    <rect y="22" width="220" height="12" fill="url(#edgeChassis)" />
    <text x="16" y="22" font-size="12" font-weight="700" fill="#ffffff">🏠 Home Edge (RTX 4060)</text>

    <!-- Subnet Label -->
    <rect x="14" y="44" width="192" height="22" rx="4" fill="#f1f5f9" stroke="#e2e8f0" />
    <text x="22" y="59" class="code-txt" fill="#475569">LAN IP: 192.168.1.105</text>

    <!-- Edge PC Workstation Illustration -->
    <g transform="translate(30, 80)">
      <!-- PC Case -->
      <rect width="80" height="120" rx="6" fill="#1e293b" stroke="#475569" stroke-width="1.5" />
      <rect x="8" y="10" width="64" height="24" rx="3" fill="#0f172a" />
      <!-- Power button & LEDs -->
      <circle cx="20" cy="22" r="3" fill="#10b981" filter="url(#glowGreen)" />
      <circle cx="32" cy="22" r="2" fill="#38bdf8" />
      <circle cx="42" cy="22" r="2" fill="#38bdf8" />
      <!-- Fan grill / RGB light bar -->
      <rect x="8" y="44" width="64" height="4" rx="2" fill="#6366f1" />
      <rect x="8" y="54" width="64" height="4" rx="2" fill="#6366f1" />
      <rect x="8" y="64" width="64" height="4" rx="2" fill="#6366f1" />
      <rect x="8" y="74" width="64" height="4" rx="2" fill="#6366f1" />
      <!-- Transcoder Badge -->
      <rect x="8" y="90" width="64" height="20" rx="3" fill="#0f172a" stroke="#10b981" stroke-width="1" />
      <text x="14" y="103" font-size="8.5" font-weight="700" fill="#10b981">NVENC Ready</text>

      <!-- Local NVMe Storage Barrel -->
      <g transform="translate(95, 20)">
        <ellipse cx="25" cy="10" rx="22" ry="7" fill="#64748b" />
        <rect x="3" y="10" width="44" height="28" fill="#475569" />
        <ellipse cx="25" cy="38" rx="22" ry="7" fill="#475569" />
        <ellipse cx="25" cy="10" rx="19" ry="5.5" fill="#94a3b8" />
        <text x="12" y="27" font-size="8" font-weight="700" fill="#ffffff">HLS Disk</text>
        <text x="9" y="52" font-size="8.5" fill="#64748b">4s Chunks</text>
      </g>
    </g>

    <!-- Mini-Jellyfin Local Streamer Port -->
    <rect x="14" y="225" width="192" height="75" rx="6" fill="#f8fafc" stroke="#e2e8f0" />
    <text x="22" y="244" font-size="11" font-weight="700" fill="#0f172a">Internal Streamer (Go)</text>
    <text x="22" y="260" font-size="9.5" fill="#64748b">• Local HTTP Port: <tspan font-family="monospace" fill="#0284c7">:8089</tspan></text>
    <text x="22" y="274" font-size="9.5" fill="#64748b">• Direct Master Playlist Serving</text>
    <text x="22" y="288" font-size="9.5" fill="#10b981" font-weight="600">✓ Standalone Mini-Jellyfin</text>

    <!-- FRP Client Uplink -->
    <rect x="14" y="312" width="192" height="75" rx="6" fill="#eff6ff" stroke="#bfdbfe" />
    <text x="22" y="331" font-size="11" font-weight="700" fill="#1e40af">FRP Tunnel Daemon</text>
    <text x="22" y="347" font-size="9.5" fill="#1e3a8a">• Outbound-only TLS tunnel</text>
    <text x="22" y="361" font-size="9.5" fill="#1e3a8a">• CGNAT &amp; Zero port-forward</text>
    <text x="22" y="375" font-size="9.5" fill="#0284c7" font-weight="600">✓ Active Secure Pipeline</text>

    <!-- Physical Uplink Spec -->
    <rect x="14" y="398" width="192" height="28" rx="4" fill="#fef3c7" stroke="#fde68a" />
    <text x="22" y="416" font-size="9.5" font-weight="700" fill="#92400e">Uplink Budget: 20 Mbps Home Cap</text>
  </g>

  <!-- ========================================================================= -->
  <!-- ZONE 2: SMART STREAM ROUTER (CENTER) -->
  <!-- ========================================================================= -->
  <g transform="translate(300, 80)" filter="url(#shadowBox)">
    <!-- Router Box -->
    <rect width="250" height="210" rx="12" fill="#ffffff" stroke="#cbd5e1" stroke-width="1.2" />
    <rect width="250" height="34" rx="12" fill="url(#routerGrad)" />
    <rect y="22" width="250" height="12" fill="url(#routerGrad)" />
    <text x="16" y="22" font-size="12" font-weight="700" fill="#ffffff">🧭 Stream Router Decision Gate</text>

    <!-- Code Decision Block -->
    <g transform="translate(14, 46)">
      <rect width="222" height="88" rx="6" fill="#0f172a" />
      <text x="12" y="20" class="code-txt" fill="#94a3b8">// Proximity Decision</text>
      <text x="12" y="36" class="code-txt" fill="#38bdf8">if <tspan fill="#f1f5f9">clientIP.InSubnet(node)</tspan> {</text>
      <text x="24" y="52" class="code-txt" fill="#4ade80">return DirectLAN_URL()</text>
      <text x="12" y="68" class="code-txt" fill="#38bdf8">} else {</text>
      <text x="24" y="84" class="code-txt" fill="#60a5fa">return OriginShield_URL()</text>
      <text x="12" y="84" class="code-txt" fill="#38bdf8" dx="180">}</text>
    </g>

    <!-- Admission Control Check -->
    <rect x="14" y="146" width="222" height="50" rx="6" fill="#fffbeb" stroke="#fef3c7" />
    <text x="22" y="164" font-size="10.5" font-weight="700" fill="#92400e">Admission Limiter: ≤ 10 Streams</text>
    <text x="22" y="180" font-size="9" fill="#78350f">Reject 11th distinct stream → HTTP 429</text>
    <text x="22" y="191" font-size="8.5" font-weight="600" fill="#b45309">"Node Busy, Retry Later" (Preserves ISP)</text>
  </g>

  <!-- Connectors from Edge to Router -->
  <path d="M 260 215 L 300 185" fill="none" stroke="#6366f1" stroke-width="2.5" stroke-dasharray="4,3" />

  <!-- ========================================================================= -->
  <!-- BRANCH A: PATH A - DIRECT LAN STREAMING (BOTTOM-LEFT / CENTER) -->
  <!-- ========================================================================= -->
  <!-- Diverging Green Pipe to LAN Path -->
  <path d="M 360 290 L 360 330 L 340 330" fill="none" stroke="#10b981" stroke-width="3" marker-end="url(#arrGreen)" />
  <rect x="368" y="300" width="95" height="18" rx="3" fill="#ecfdf5" stroke="#a7f3d0" />
  <text x="374" y="313" font-size="9" font-weight="700" fill="#047857">MATCH: SAME LAN</text>

  <g transform="translate(300, 335)" filter="url(#shadowBox)">
    <rect width="250" height="185" rx="10" fill="url(#lanCard)" stroke="#6ee7b7" stroke-width="1.5" />
    <!-- Header Badge -->
    <rect x="12" y="10" width="226" height="24" rx="5" fill="#059669" />
    <text x="20" y="26" font-size="11" font-weight="700" fill="#ffffff">📶 PATH A: Direct LAN Bypass</text>

    <!-- Metrics Cards -->
    <g transform="translate(12, 42)">
      <rect width="108" height="42" rx="4" fill="#ffffff" stroke="#a7f3d0" />
      <text x="8" y="16" font-size="9" fill="#065f46">Latency (TTFB)</text>
      <text x="8" y="34" font-size="14" font-weight="800" fill="#059669">&lt; 3.2 ms</text>

      <rect x="116" width="110" height="42" rx="4" fill="#ffffff" stroke="#a7f3d0" />
      <text x="8" y="16" font-size="9" fill="#065f46">Cloud Egress Cost</text>
      <text x="8" y="34" font-size="14" font-weight="800" fill="#059669">&#x20B9;0.00 / Zero</text>
    </g>

    <!-- TV / Client Illustration -->
    <g transform="translate(16, 94)">
      <!-- Mini TV Screen -->
      <rect width="70" height="44" rx="3" fill="#0f172a" stroke="#334155" stroke-width="1.5" />
      <rect x="4" y="4" width="62" height="36" rx="2" fill="#047857" />
      <polygon points="30,16 30,28 42,22" fill="#ffffff" />
      <line x1="35" y1="44" x2="35" y2="52" stroke="#334155" stroke-width="2" />
      <line x1="25" y1="52" x2="45" y2="52" stroke="#334155" stroke-width="2" />

      <!-- Text details -->
      <text x="80" y="14" font-size="10.5" font-weight="700" fill="#065f46">Living Room Smart TV</text>
      <text x="80" y="28" font-size="9" fill="#047857">Pure local Gigabit Wi-Fi</text>
      <text x="80" y="42" font-size="9" fill="#047857">Zero WAN internet used</text>
    </g>

    <!-- Offline Resilience Badge -->
    <rect x="12" y="152" width="226" height="22" rx="4" fill="#ffffff" stroke="#10b981" />
    <text x="20" y="167" font-size="9" font-weight="700" fill="#065f46">✓ Operates offline during ISP internet cut</text>
  </g>

  <!-- Direct LAN Cable / Beam to Smart TV -->
  <path d="M 260 250 C 275 250, 285 360, 300 360" fill="none" stroke="#10b981" stroke-width="3" stroke-linecap="round" />
  <circle cx="280" cy="305" r="4" fill="#10b981" filter="url(#glowGreen)" />

  <!-- ========================================================================= -->
  <!-- BRANCH B: PATH B - ORIGIN SHIELD SINGLE INGEST (RIGHT) -->
  <!-- ========================================================================= -->
  <!-- Diverging Blue Pipe to Cloud Origin Shield -->
  <path d="M 550 140 L 600 140" fill="none" stroke="#0284c7" stroke-width="3.5" marker-end="url(#arrBlue)" />
  <rect x="535" y="112" width="95" height="18" rx="3" fill="#e0f2fe" stroke="#bae6fd" />
  <text x="541" y="125" font-size="9" font-weight="700" fill="#0369a1">NO: REMOTE</text>

  <!-- Cloud Tier Container Box -->
  <g transform="translate(615, 80)" filter="url(#shadowBox)">
    <rect width="545" height="440" rx="14" fill="url(#remoteCard)" stroke="#7dd3fc" stroke-width="1.5" />

    <!-- Cloud Header Banner -->
    <rect width="545" height="34" rx="14" fill="url(#shieldGrad)" />
    <rect y="22" width="545" height="12" fill="url(#shieldGrad)" />
    <text x="18" y="22" font-size="12.5" font-weight="700" fill="#ffffff">🛡️ PATH B: Cloud Origin Shield &amp; Multi-Viewer Fan-Out</text>

    <!-- ======================================== -->
    <!-- COMPONENT 1: ORIGIN SHIELD CACHE HUB -->
    <!-- ======================================== -->
    <g transform="translate(18, 50)">
      <!-- Origin Shield Graphical Emblem -->
      <path d="M 45 6 Q 65 0 85 6 Q 105 20 85 60 Q 65 85 45 92 Q 25 85 5 60 Q -15 20 5 6 Q 25 0 45 6 Z" fill="url(#shieldGrad)" stroke="#38bdf8" stroke-width="1.5" transform="translate(10, 0)" />
      <!-- Inner Shield Core -->
      <circle cx="55" cy="40" r="22" fill="#0f172a" />
      <text x="44" y="44" font-size="16" fill="#38bdf8">S3</text>

      <!-- Cache Specs & Architecture -->
      <g transform="translate(100, 5)">
        <text x="0" y="16" font-size="12" font-weight="700" fill="#0c4a6e">Origin Shield S3 LRU Bucket</text>
        <text x="0" y="32" class="body-txt">• Temporary 24-Hour Auto-Purge Cache</text>
        <text x="0" y="47" class="body-txt">• Pulls chunk <tspan font-weight="700" fill="#d97706">EXACTLY ONCE</tspan> from Home Edge</text>
        <text x="0" y="62" class="body-txt">• Zero permanent cloud storage liability (<tspan font-weight="700" fill="#0284c7">Invariant 1</tspan>)</text>
        <text x="0" y="78" font-size="10" font-weight="700" fill="#0369a1">&#x20B9;0.00 Permanent Storage Cost</text>
      </g>
    </g>

    <!-- Divider Line -->
    <line x1="20" y1="150" x2="525" y2="150" stroke="#bae6fd" stroke-width="1.2" stroke-dasharray="3,3" />

    <!-- ======================================== -->
    <!-- STEP 1: CACHE MISS (PULL ONCE VIA TUNNEL) -->
    <!-- ======================================== -->
    <g transform="translate(18, 160)">
      <rect width="509" height="74" rx="8" fill="#ffffff" stroke="#fcd34d" stroke-width="1.2" />
      <!-- Miss Badge -->
      <rect x="12" y="10" width="82" height="18" rx="3" fill="#fef3c7" stroke="#f59e0b" />
      <text x="18" y="23" font-size="9" font-weight="700" fill="#92400e">CACHE MISS</text>

      <text x="104" y="24" font-size="11.5" font-weight="700" fill="#92400e">Viewer 1 Requests New Segment (e.g. seg_042.ts)</text>
      <text x="14" y="46" font-size="10" fill="#475569">1. Go API pulls segment over FRP Tunnel from home edge (<tspan font-weight="700" fill="#b45309">1x Uplink: ~2.5 MB</tspan>)</text>
      <text x="14" y="62" font-size="10" fill="#475569">2. Writes chunk into S3 Cache bucket with signed access token &amp; 24h LRU timestamp</text>
    </g>

    <!-- Connective Pulse Arrow from Edge Tunnel to Cache Miss -->
    <!-- (Visual tunnel conduit) -->

    <!-- ======================================== -->
    <!-- STEP 2: CACHE HIT (FAN-OUT TO 10 VIEWERS) -->
    <!-- ======================================== -->
    <g transform="translate(18, 244)">
      <rect width="509" height="180" rx="8" fill="#ffffff" stroke="#86efac" stroke-width="1.2" />
      <!-- Hit Badge -->
      <rect x="12" y="10" width="76" height="18" rx="3" fill="#dcfce7" stroke="#16a34a" />
      <text x="18" y="23" font-size="9" font-weight="700" fill="#15803d">CACHE HIT</text>

      <text x="98" y="24" font-size="11.5" font-weight="700" fill="#15803d">Viewers 2 through 10 Request Same Segment</text>
      <text x="14" y="45" font-size="10" fill="#334155">• Served directly from High-Throughput S3 Cache / Cloud CDN</text>
      <text x="14" y="60" font-size="10" font-weight="700" fill="#059669">• Residential Home Uplink Load: EXACTLY 0 BYTES / 0 BPS</text>

      <!-- Multi-Viewer Fan-Out Devices Grid -->
      <g transform="translate(14, 72)">
        <!-- Device 1: Laptop -->
        <g transform="translate(0, 0)">
          <rect width="40" height="26" rx="2" fill="#1e293b" />
          <rect x="2" y="2" width="36" height="20" rx="1" fill="#38bdf8" />
          <polygon points="17,8 17,16 25,12" fill="#ffffff" />
          <rect x="-4" y="26" width="48" height="4" rx="1" fill="#64748b" />
          <text x="4" y="40" font-size="7.5" fill="#64748b">Viewer 2</text>
        </g>

        <!-- Device 2: Smartphone -->
        <g transform="translate(60, 2)">
          <rect width="18" height="30" rx="3" fill="#1e293b" />
          <rect x="2" y="3" width="14" height="22" rx="1" fill="#38bdf8" />
          <circle cx="9" cy="27" r="1" fill="#94a3b8" />
          <text x="0" y="40" font-size="7.5" fill="#64748b">Viewer 3</text>
        </g>

        <!-- Device 3: Smartphone -->
        <g transform="translate(98, 2)">
          <rect width="18" height="30" rx="3" fill="#1e293b" />
          <rect x="2" y="3" width="14" height="22" rx="1" fill="#38bdf8" />
          <circle cx="9" cy="27" r="1" fill="#94a3b8" />
          <text x="0" y="40" font-size="7.5" fill="#64748b">Viewer 4</text>
        </g>

        <!-- Device 4: Tablet -->
        <g transform="translate(136, 0)">
          <rect width="32" height="26" rx="2" fill="#1e293b" />
          <rect x="2" y="2" width="28" height="22" rx="1" fill="#38bdf8" />
          <polygon points="13,9 13,17 19,13" fill="#ffffff" />
          <text x="3" y="40" font-size="7.5" fill="#64748b">Viewer 5</text>
        </g>

        <!-- Dots indicating 6..10 -->
        <text x="180" y="18" font-size="16" font-weight="700" fill="#0284c7">. . .</text>

        <!-- Device 10: Desktop -->
        <g transform="translate(225, 0)">
          <rect width="38" height="24" rx="2" fill="#1e293b" />
          <rect x="2" y="2" width="34" height="18" rx="1" fill="#38bdf8" />
          <polygon points="16,8 16,14 22,11" fill="#ffffff" />
          <rect x="16" y="24" width="6" height="5" fill="#64748b" />
          <rect x="11" y="29" width="16" height="2" fill="#64748b" />
          <text x="0" y="40" font-size="7.5" fill="#64748b">Viewer 10</text>
        </g>

        <!-- Uplink Savings Calculation Callout -->
        <g transform="translate(285, -2)">
          <rect width="195" height="44" rx="5" fill="#ecfdf5" stroke="#a7f3d0" />
          <text x="8" y="16" font-size="9.5" font-weight="700" fill="#065f46">Uplink Preservation Ratio:</text>
          <text x="8" y="32" font-size="12" font-weight="800" fill="#059669">90% Bandwidth Saved</text>
          <text x="140" y="32" font-size="9" fill="#047857">(1x vs 10x)</text>
        </g>
      </g>

      <!-- Bottom Invariant Guarantee -->
      <rect x="12" y="148" width="485" height="22" rx="4" fill="#f0fdf4" stroke="#86efac" />
      <text x="20" y="163" font-size="9.5" font-weight="700" fill="#166534">Invariant 7 Upheld: Home uplink never saturates under viral community traffic</text>
    </g>

  </g>

  <!-- Connective FRP Tunnel Pipeline Arch across zones -->
  <path d="M 230 350 C 270 350, 480 340, 600 230 C 608 220, 620 210, 633 205" fill="none" stroke="#f59e0b" stroke-width="2.5" stroke-dasharray="5,4" marker-end="url(#arrAmber)" />
  <rect x="375" y="260" width="130" height="18" rx="3" fill="#fffbeb" stroke="#fcd34d" />
  <text x="381" y="273" font-size="8.5" font-weight="700" fill="#b45309">FRP Pull (1x on Miss)</text>

</svg>"""

svg_file = os.path.join(fig_dir, "diag_05_dual_path_shield.svg")
png_file = os.path.join(fig_dir, "diag_05_dual_path_shield.png")
html_file = os.path.join(fig_dir, "diag_05_dual_path_shield.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_content)

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  body {{ margin: 0; padding: 0; background: #ffffff; display: flex; align-items: center; justify-content: center; }}
  svg {{ width: 1200px; height: 560px; }}
</style>
</head>
<body>
  {svg_content}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    f"--screenshot={png_file}",
    "--window-size=1200,560",
    html_file
]
subprocess.run(args, capture_output=True)
print(f"Rendered {png_file} - Size: {os.path.getsize(png_file)} bytes")

