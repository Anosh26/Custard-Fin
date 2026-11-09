import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

svg_fixed = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 520" width="1200" height="520" style="background-color: #ffffff; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
  <defs>
    <!-- Soft Drop Shadow (#0f172a at 6% opacity) -->
    <filter id="shadow" x="-5%" y="-5%" width="110%" height="110%">
      <feDropShadow dx="0" dy="4" stdDeviation="6" flood-color="#0f172a" flood-opacity="0.06" />
    </filter>
    <filter id="shadow-sm" x="-5%" y="-5%" width="110%" height="110%">
      <feDropShadow dx="0" dy="2" stdDeviation="3" flood-color="#0f172a" flood-opacity="0.05" />
    </filter>

    <!-- Arrowhead Markers -->
    <marker id="arrow-green" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#059669" />
    </marker>
    <marker id="arrow-red" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#dc2626" />
    </marker>
    <marker id="arrow-cyan" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#0284c7" />
    </marker>
    <marker id="arrow-amber" viewBox="0 0 10 10" refX="7" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#d97706" />
    </marker>

    <!-- Subtle Gradients -->
    <linearGradient id="grad-home" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#faf5ff" />
      <stop offset="100%" stop-color="#f3e8ff" stop-opacity="0.5" />
    </linearGradient>
    <linearGradient id="grad-cgnat" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#fff1f2" />
      <stop offset="100%" stop-color="#ffe4e6" stop-opacity="0.5" />
    </linearGradient>
    <linearGradient id="grad-cloud" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#f0f9ff" />
      <stop offset="100%" stop-color="#e0f2fe" stop-opacity="0.4" />
    </linearGradient>
  </defs>

  <!-- ================= HEADER ================= -->
  <text x="600" y="34" text-anchor="middle" font-size="19" font-weight="800" fill="#0f172a" letter-spacing="-0.3">
    Carrier-Grade NAT (CGNAT) Traversal via Outbound Reverse Tunneling
  </text>
  <text x="600" y="54" text-anchor="middle" font-size="11.5" font-weight="500" fill="#64748b">
    Bypassing RFC 6598 ISP IPv4 Exhaustion Without Inbound Port Forwarding via Persistent TLS 1.3 Multiplexed Conduits
  </text>

  <!-- ================= ZONE 1: HOME EDGE ================= -->
  <g transform="translate(30, 72)">
    <rect x="0" y="0" width="330" height="425" rx="12" fill="url(#grad-home)" stroke="#7c3aed" stroke-width="1.5" filter="url(#shadow)" />
    
    <!-- Zone Header Banner -->
    <path d="M 0 12 C 0 5.37 5.37 0 12 0 L 318 0 C 324.63 0 330 5.37 330 12 L 330 38 L 0 38 Z" fill="#7c3aed" />
    <text x="165" y="24" text-anchor="middle" font-size="12" font-weight="700" fill="#ffffff" letter-spacing="0.5">
      HOME EDGE (192.168.1.0/24)
    </text>

    <!-- Server Node Box (edged) -->
    <g transform="translate(20, 54)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="290" height="106" rx="8" fill="#ffffff" stroke="#6366f1" stroke-width="1.5" />
      
      <!-- Mini Server Rack Chassis Icon -->
      <rect x="16" y="18" width="54" height="70" rx="5" fill="#1e1b4b" stroke="#4338ca" stroke-width="1.5" />
      <line x1="24" y1="32" x2="62" y2="32" stroke="#4338ca" stroke-width="2" />
      <line x1="24" y1="44" x2="62" y2="44" stroke="#4338ca" stroke-width="2" />
      <line x1="24" y1="56" x2="46" y2="56" stroke="#4338ca" stroke-width="2" />
      <!-- Power LEDs -->
      <circle cx="28" cy="74" r="3.5" fill="#10b981" />
      <circle cx="40" cy="74" r="2.5" fill="#38bdf8" />

      <!-- Server Labels -->
      <text x="84" y="30" font-size="12.5" font-weight="700" fill="#1e1b4b">Home Server (edged)</text>
      <rect x="218" y="16" width="58" height="18" rx="4" fill="#eef2ff" stroke="#6366f1" stroke-width="1" />
      <text x="247" y="28" text-anchor="middle" font-family="monospace" font-size="9.5" font-weight="700" fill="#4338ca">:8089</text>

      <text x="84" y="50" font-size="10" fill="#475569">• Standalone Mini-Jellyfin Engine</text>
      <text x="84" y="66" font-size="10" fill="#475569">• Embedded frpc Reverse Client</text>
      <text x="84" y="82" font-size="10" font-weight="600" fill="#7c3aed">• Private NAT IP: 192.168.1.105</text>
    </g>

    <!-- Bottom Row: NVMe Storage (Left) + LAN Viewers (Right) -->
    <!-- NVMe Cylinder Card -->
    <g transform="translate(20, 204)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="128" height="150" rx="8" fill="#ffffff" stroke="#7c3aed" stroke-width="1.2" />
      
      <!-- 3D Purple Storage Cylinder -->
      <g transform="translate(36, 18)">
        <path d="M 0 10 L 0 46 C 0 52 56 52 56 46 L 56 10 Z" fill="#6366f1" />
        <ellipse cx="28" cy="10" rx="28" ry="8" fill="#a855f7" stroke="#4338ca" stroke-width="1.2" />
        <path d="M 0 22 C 0 28 56 28 56 22" fill="none" stroke="#e9d5ff" stroke-width="1.2" />
        <path d="M 0 34 C 0 40 56 40 56 34" fill="none" stroke="#e9d5ff" stroke-width="1.2" />
        <path d="M 0 10 L 0 46 C 0 52 56 52 56 46 L 56 10" fill="none" stroke="#4338ca" stroke-width="1.2" />
      </g>

      <text x="64" y="88" text-anchor="middle" font-size="10.5" font-weight="700" fill="#1e1b4b">NVMe Storage</text>
      <text x="64" y="103" text-anchor="middle" font-size="9" fill="#64748b">Master HLS Chunks</text>
      <text x="64" y="118" text-anchor="middle" font-family="monospace" font-size="8.5" fill="#7c3aed">4.0s .ts + .m3u8</text>
      <text x="64" y="134" text-anchor="middle" font-size="8.5" font-weight="600" fill="#4338ca">100% Data at Home</text>
    </g>

    <!-- Internal Connector: Server <-> NVMe -->
    <line x1="84" y1="160" x2="84" y2="204" stroke="#7c3aed" stroke-width="1.8" stroke-dasharray="3,2" />

    <!-- LAN Viewers Card -->
    <g transform="translate(162, 204)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="148" height="150" rx="8" fill="#ecfdf5" stroke="#059669" stroke-width="1.5" />
      
      <!-- Desktop Monitor Icon -->
      <rect x="42" y="18" width="64" height="42" rx="4" fill="#064e3b" stroke="#059669" stroke-width="1.5" />
      <polygon points="66,40 66,30 78,35" fill="#10b981" />
      <rect x="66" y="60" width="16" height="6" fill="#047857" />
      <rect x="54" y="66" width="40" height="3" rx="1.5" fill="#047857" />

      <text x="74" y="88" text-anchor="middle" font-size="10.5" font-weight="700" fill="#064e3b">LAN Viewers</text>
      <text x="74" y="103" text-anchor="middle" font-size="9" font-weight="600" fill="#059669">Sub-4ms Direct Wi-Fi</text>
      <text x="74" y="118" text-anchor="middle" font-size="8.5" fill="#047857">Zero WAN Bandwidth</text>
      <text x="74" y="134" text-anchor="middle" font-size="8.5" font-weight="700" fill="#059669">Cost: ₹0.00 (Free)</text>
    </g>

    <!-- Direct Green Arrow: Server -> LAN Viewers (LAN Bypass) -->
    <line x1="236" y1="160" x2="236" y2="198" stroke="#059669" stroke-width="2.2" marker-end="url(#arrow-green)" />
    <rect x="188" y="171" width="96" height="18" rx="9" fill="#059669" />
    <text x="236" y="183" text-anchor="middle" font-size="9" font-weight="700" fill="#ffffff">LAN Bypass</text>

    <!-- Bottom Pill Banner -->
    <rect x="16" y="374" width="298" height="34" rx="17" fill="#059669" filter="url(#shadow-sm)" />
    <text x="165" y="395" text-anchor="middle" font-size="10" font-weight="700" fill="#ffffff" letter-spacing="0.2">
      ✓ ZERO INBOUND PORTS REQUIRED — CGNAT Immune
    </text>
  </g>

  <!-- ================= ZONE 2: ISP CGNAT BARRIER ================= -->
  <g transform="translate(415, 72)">
    <rect x="0" y="0" width="190" height="425" rx="12" fill="url(#grad-cgnat)" stroke="#dc2626" stroke-width="1.5" stroke-dasharray="6,3" filter="url(#shadow)" />
    
    <!-- Header Banner -->
    <path d="M 0 12 C 0 5.37 5.37 0 12 0 L 178 0 C 184.63 0 190 5.37 190 12 L 190 48 L 0 48 Z" fill="#dc2626" />
    <text x="95" y="21" text-anchor="middle" font-size="11.5" font-weight="700" fill="#ffffff" letter-spacing="0.4">
      ISP CGNAT BARRIER
    </text>
    <text x="95" y="38" text-anchor="middle" font-family="monospace" font-size="9" font-weight="600" fill="#fecdd3">
      RFC 6598 (100.64.0.0/10)
    </text>

    <!-- Top Sub-Box: INBOUND BLOCKED -->
    <g transform="translate(13, 60)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="164" height="142" rx="8" fill="#ffffff" stroke="#ef4444" stroke-width="1.5" />
      
      <!-- Red Shield with Cross -->
      <circle cx="82" cy="22" r="13" fill="#dc2626" />
      <text x="82" y="27" text-anchor="middle" font-size="13" font-weight="800" fill="#ffffff">✕</text>

      <text x="82" y="48" text-anchor="middle" font-size="10" font-weight="800" fill="#991b1b">INBOUND BLOCKED</text>
      <line x1="16" y1="56" x2="148" y2="56" stroke="#fecdd3" stroke-width="1" />
      
      <!-- Descriptive Pill inside the box -->
      <rect x="14" y="64" width="136" height="20" rx="4" fill="#fee2e2" stroke="#f87171" stroke-width="0.8" />
      <text x="82" y="77" text-anchor="middle" font-size="8.5" font-weight="800" fill="#b91c1c">① Inbound SYN Dropped</text>

      <text x="82" y="99" text-anchor="middle" font-size="8.5" fill="#475569">No Public Home IPv4</text>
      <text x="82" y="113" text-anchor="middle" font-size="8.5" fill="#475569">Port Forwarding Fails</text>
      <rect x="14" y="121" width="136" height="15" rx="3" fill="#fff1f2" />
      <text x="82" y="132" text-anchor="middle" font-family="monospace" font-size="7.5" font-weight="600" fill="#991b1b">ERR_CONNECTION_REFUSED</text>
    </g>

    <!-- Bottom Sub-Box: OUTBOUND PASS (Tunnel Gateway) -->
    <g transform="translate(13, 218)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="164" height="190" rx="8" fill="#ffffff" stroke="#10b981" stroke-width="1.8" />
      
      <!-- Green Checkmark Circle -->
      <circle cx="82" cy="22" r="14" fill="#059669" />
      <text x="82" y="28" text-anchor="middle" font-size="14" font-weight="800" fill="#ffffff">✓</text>

      <text x="82" y="50" text-anchor="middle" font-size="10.5" font-weight="800" fill="#047857">OUTBOUND PASS</text>
      <line x1="16" y1="58" x2="148" y2="58" stroke="#a7f3d0" stroke-width="1" />
      
      <text x="82" y="74" text-anchor="middle" font-size="9" font-weight="600" fill="#059669">Stateful NAT Mapping</text>
      <text x="82" y="90" text-anchor="middle" font-size="8.5" fill="#475569">Reverse Tunneling (FRP)</text>

      <!-- Dedicated Green Outbound Tunnel Gateway Badge (Fits neatly with 12px margin, 0 overflow) -->
      <g transform="translate(12, 102)">
        <rect x="0" y="0" width="140" height="46" rx="6" fill="#047857" stroke="#10b981" stroke-width="1.2" />
        <text x="70" y="18" text-anchor="middle" font-size="9.5" font-weight="800" fill="#ffffff">
          ② Outbound TLS Tunnel
        </text>
        <rect x="18" y="24" width="104" height="16" rx="3" fill="#065f46" />
        <text x="70" y="35" text-anchor="middle" font-family="monospace" font-size="8.5" font-weight="700" fill="#a7f3d0">
          FRP Reverse Port :7000
        </text>
      </g>

      <text x="82" y="166" text-anchor="middle" font-size="8.5" font-weight="600" fill="#059669">Multiplexed yamux Stream</text>
      <text x="82" y="180" text-anchor="middle" font-size="8" fill="#64748b">15s Keep-Alive Heartbeat</text>
    </g>
  </g>

  <!-- ================= ZONE 3: PUBLIC CLOUD TIER ================= -->
  <g transform="translate(650, 72)">
    <rect x="0" y="0" width="520" height="425" rx="12" fill="url(#grad-cloud)" stroke="#0284c7" stroke-width="1.5" filter="url(#shadow)" />
    
    <!-- Header Banner -->
    <path d="M 0 12 C 0 5.37 5.37 0 12 0 L 508 0 C 514.63 0 520 5.37 520 12 L 520 38 L 0 38 Z" fill="#0284c7" />
    <text x="260" y="24" text-anchor="middle" font-size="12" font-weight="700" fill="#ffffff" letter-spacing="0.5">
      PUBLIC CLOUD TIER (GCP VM / RELAY)
    </text>

    <!-- Top Row: 3 Service Boxes Side-by-Side -->
    <!-- Box A: FRP SERVER -->
    <g transform="translate(16, 54)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="152" height="126" rx="8" fill="#ffffff" stroke="#0284c7" stroke-width="1.5" />
      <rect x="0" y="0" width="152" height="26" rx="8" fill="#e0f2fe" />
      <text x="76" y="17" text-anchor="middle" font-size="10" font-weight="800" fill="#0369a1">FRP SERVER (frps)</text>
      
      <text x="76" y="48" text-anchor="middle" font-family="monospace" font-size="9.5" font-weight="700" fill="#0284c7">Listens on :7000 (TLS)</text>
      <text x="76" y="68" text-anchor="middle" font-size="9.5" fill="#334155">Multiplexed Tunnel Hub</text>
      <text x="76" y="86" text-anchor="middle" font-size="9" fill="#64748b">Static Public IPv4 Anchor</text>
      <rect x="20" y="98" width="112" height="18" rx="4" fill="#f0f9ff" />
      <text x="76" y="110" text-anchor="middle" font-size="8.5" font-weight="600" fill="#0369a1">yamux Stream Mux</text>
    </g>

    <!-- Box B: STREAM ROUTER -->
    <g transform="translate(184, 54)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="152" height="126" rx="8" fill="#ffffff" stroke="#0284c7" stroke-width="1.5" />
      <rect x="0" y="0" width="152" height="26" rx="8" fill="#e0f2fe" />
      <text x="76" y="17" text-anchor="middle" font-size="10" font-weight="800" fill="#0369a1">STREAM ROUTER (apid)</text>
      
      <text x="76" y="48" text-anchor="middle" font-size="9.5" font-weight="600" fill="#0f172a">Subnet Proximity Router</text>
      <text x="76" y="68" text-anchor="middle" font-size="9.5" font-weight="700" fill="#dc2626">10-Stream Limiter Cap</text>
      <text x="76" y="86" text-anchor="middle" font-size="9" fill="#64748b">Short-Lived HMAC Tokens</text>
      <rect x="20" y="98" width="112" height="18" rx="4" fill="#f0f9ff" />
      <text x="76" y="110" text-anchor="middle" font-family="monospace" font-size="8.5" font-weight="600" fill="#0369a1">Go Control Plane :8080</text>
    </g>

    <!-- Box C: ORIGIN SHIELD CACHE -->
    <g transform="translate(352, 54)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="152" height="126" rx="8" fill="#fefce8" stroke="#d97706" stroke-width="1.5" />
      <rect x="0" y="0" width="152" height="26" rx="8" fill="#fef08a" />
      <text x="76" y="17" text-anchor="middle" font-size="10" font-weight="800" fill="#b45309">ORIGIN SHIELD CACHE</text>
      
      <text x="76" y="48" text-anchor="middle" font-size="9.5" font-weight="700" fill="#d97706">S3 24h LRU Buffer</text>
      <text x="76" y="68" text-anchor="middle" font-size="9.5" font-weight="600" fill="#78350f">Single-Ingest Upload</text>
      <text x="76" y="86" text-anchor="middle" font-size="9" fill="#92400e">Zero Permanent Storage</text>
      <rect x="16" y="98" width="120" height="18" rx="4" fill="#ffffff" stroke="#f59e0b" stroke-width="0.8" />
      <text x="76" y="110" text-anchor="middle" font-size="8.5" font-weight="700" fill="#b45309">1x Pull → 10x Fan-Out</text>
    </g>

    <!-- Internal Cloud Horizontal Connectors -->
    <line x1="168" y1="117" x2="184" y2="117" stroke="#0284c7" stroke-width="2" marker-end="url(#arrow-cyan)" />
    <line x1="336" y1="117" x2="352" y2="117" stroke="#d97706" stroke-width="2" marker-end="url(#arrow-amber)" />

    <!-- Bottom Section: Public Community Viewers (10x Load) -->
    <g transform="translate(16, 204)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="488" height="136" rx="8" fill="#ffffff" stroke="#38bdf8" stroke-width="1.5" />
      <text x="244" y="21" text-anchor="middle" font-size="11" font-weight="800" fill="#0369a1">
        Public Community Viewers (10x Concurrent Load)
      </text>

      <!-- Device 1: Laptop (Viewer 1 - Cache Miss Trigger) -->
      <g transform="translate(18, 34)">
        <rect x="0" y="0" width="82" height="86" rx="6" fill="#fefce8" stroke="#f59e0b" stroke-width="1.2" />
        <rect x="21" y="12" width="40" height="26" rx="3" fill="#78350f" />
        <path d="M 15 38 L 67 38 L 71 43 L 11 43 Z" fill="#d97706" />
        <text x="41" y="58" text-anchor="middle" font-size="9.5" font-weight="700" fill="#b45309">Viewer 1</text>
        <text x="41" y="71" text-anchor="middle" font-size="8" font-weight="600" fill="#d97706">Laptop (1x Pull)</text>
        <text x="41" y="81" text-anchor="middle" font-family="monospace" font-size="7.5" fill="#92400e">MISS → FETCH</text>
      </g>

      <!-- Device 2: Smartphone 1 (Viewers 2-3) -->
      <g transform="translate(110, 34)">
        <rect x="0" y="0" width="82" height="86" rx="6" fill="#f0f9ff" stroke="#38bdf8" stroke-width="1.2" />
        <rect x="29" y="10" width="24" height="36" rx="4" fill="#0369a1" />
        <circle cx="41" cy="42" r="1.5" fill="#ffffff" />
        <text x="41" y="58" text-anchor="middle" font-size="9.5" font-weight="700" fill="#0369a1">Viewers 2–3</text>
        <text x="41" y="71" text-anchor="middle" font-size="8" fill="#0284c7">Smartphone</text>
        <text x="41" y="81" text-anchor="middle" font-family="monospace" font-size="7.5" font-weight="700" fill="#059669">CACHE HIT</text>
      </g>

      <!-- Device 3: Smartphone 2 (Viewers 4-5) -->
      <g transform="translate(202, 34)">
        <rect x="0" y="0" width="82" height="86" rx="6" fill="#f0f9ff" stroke="#38bdf8" stroke-width="1.2" />
        <rect x="29" y="10" width="24" height="36" rx="4" fill="#0369a1" />
        <circle cx="41" cy="42" r="1.5" fill="#ffffff" />
        <text x="41" y="58" text-anchor="middle" font-size="9.5" font-weight="700" fill="#0369a1">Viewers 4–5</text>
        <text x="41" y="71" text-anchor="middle" font-size="8" fill="#0284c7">Smartphone</text>
        <text x="41" y="81" text-anchor="middle" font-family="monospace" font-size="7.5" font-weight="700" fill="#059669">CACHE HIT</text>
      </g>

      <!-- Device 4: Tablet (Viewers 6-8) -->
      <g transform="translate(294, 34)">
        <rect x="0" y="0" width="82" height="86" rx="6" fill="#f0f9ff" stroke="#38bdf8" stroke-width="1.2" />
        <rect x="25" y="10" width="32" height="36" rx="4" fill="#0369a1" />
        <circle cx="41" cy="42" r="1.5" fill="#ffffff" />
        <text x="41" y="58" text-anchor="middle" font-size="9.5" font-weight="700" fill="#0369a1">Viewers 6–8</text>
        <text x="41" y="71" text-anchor="middle" font-size="8" fill="#0284c7">Tablet iPad</text>
        <text x="41" y="81" text-anchor="middle" font-family="monospace" font-size="7.5" font-weight="700" fill="#059669">CACHE HIT</text>
      </g>

      <!-- Device 5: Desktop (Viewers 9-10) -->
      <g transform="translate(386, 34)">
        <rect x="0" y="0" width="82" height="86" rx="6" fill="#f0f9ff" stroke="#38bdf8" stroke-width="1.2" />
        <rect x="21" y="10" width="40" height="26" rx="3" fill="#0369a1" />
        <rect x="35" y="36" width="12" height="5" fill="#0284c7" />
        <rect x="27" y="41" width="28" height="3" rx="1" fill="#0284c7" />
        <text x="41" y="58" text-anchor="middle" font-size="9.5" font-weight="700" fill="#0369a1">Viewers 9–10</text>
        <text x="41" y="71" text-anchor="middle" font-size="8" fill="#0284c7">Desktop PC</text>
        <text x="41" y="81" text-anchor="middle" font-family="monospace" font-size="7.5" font-weight="700" fill="#059669">CACHE HIT</text>
      </g>
    </g>

    <!-- Arrow from Origin Shield down to Viewers Box -->
    <line x1="428" y1="180" x2="428" y2="199" stroke="#d97706" stroke-width="2.2" marker-end="url(#arrow-amber)" />

    <!-- Bottom Origin Shield Bandwidth Callout Banner -->
    <g transform="translate(16, 352)" filter="url(#shadow-sm)">
      <rect x="0" y="0" width="488" height="58" rx="8" fill="#fefce8" stroke="#f59e0b" stroke-width="1.5" />
      
      <!-- Vector Lightning Icon + Title Line -->
      <polygon points="108,8 102,17 107,17 104,24 113,14 108,14" fill="#d97706" />
      <text x="252" y="19" text-anchor="middle" font-size="10" font-weight="800" fill="#b45309">
        Origin Shield Bandwidth Absorption (Invariant 7)
      </text>

      <line x1="24" y1="25" x2="464" y2="25" stroke="#fde68a" stroke-width="1" />

      <text x="244" y="38" text-anchor="middle" font-size="9" font-weight="600" fill="#78350f">
        Viewer 1 triggers single pull via reverse tunnel (1x Uplink: ~2.5 MB)
      </text>

      <text x="244" y="51" text-anchor="middle" font-size="9" font-weight="700" fill="#059669">
        Viewers 2–10 served directly from Cloud Cache (90% Home Bandwidth Saved)
      </text>
    </g>
  </g>

  <!-- ================= CONNECTING HIGHWAY ARROWS (ZERO OVERLAP GUARANTEED) ================= -->
  
  <!-- LANE 1 (TOP): Inbound SYN Dropped (From Cloud border at y=145 heading left into CGNAT Red Shield) -->
  <g>
    <path d="M 648 145 L 608 145" fill="none" stroke="#dc2626" stroke-width="2.5" stroke-dasharray="5,3" marker-end="url(#arrow-red)" />
    <!-- Bounce-back reflection arc showing rejection -->
    <path d="M 610 145 Q 624 130 642 133" fill="none" stroke="#ef4444" stroke-width="1.5" />
    <rect x="611" y="114" width="34" height="16" rx="3" fill="#fee2e2" stroke="#f87171" stroke-width="0.8" />
    <text x="628" y="125" text-anchor="middle" font-size="7.5" font-weight="800" fill="#b91c1c">DROP</text>
  </g>

  <!-- LANE 2 (BOTTOM): Outbound TLS 1.3 Tunnel -->
  <!-- Exit Home Server at (340, 107) -> drop down in left gap (x=388) -> enter CGNAT Gateway at (427, 343) -->
  <path d="M 340 107 L 388 107 L 388 343 L 427 343" fill="none" stroke="#059669" stroke-width="3" />
  
  <!-- Exit CGNAT Gateway at (567, 343) -> go into right gap (x=628) -> rise up in right gap ONLY TO y=190 -> enter Cloud FRP Server at (666, 190) -->
  <!-- NOTICE: y=190 is 45px BELOW the red arrow (y=145). THEY NEVER INTERSECT! -->
  <path d="M 567 343 L 628 343 L 628 190 L 664 190" fill="none" stroke="#059669" stroke-width="3" marker-end="url(#arrow-green)" />

  <!-- Directional Pulse Indicator Dots inside Gaps -->
  <circle cx="388" cy="225" r="4" fill="#10b981" />
  <circle cx="628" cy="265" r="4" fill="#10b981" />

</svg>"""

svg_file = os.path.join(fig_dir, "diag_01_nat_traversal.svg")
png_file = os.path.join(fig_dir, "diag_01_nat_traversal.png")
html_file = os.path.join(fig_dir, "diag_01_render_hd.html")

with open(svg_file, "w", encoding="utf-8") as f:
    f.write(svg_fixed)

# Render at 3600 x 1560 Ultra-HD
w, h = 3600, 1560
html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  * {{ margin: 0; padding: 0; box-sizing: border-box; }}
  html, body {{ width: {w}px; height: {h}px; background: #ffffff; overflow: hidden; }}
  svg {{ width: {w}px; height: {h}px; display: block; }}
</style>
</head>
<body>
  {svg_fixed}
</body>
</html>"""

with open(html_file, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    "--hide-scrollbars",
    f"--window-size={w},{h}",
    f"--screenshot={png_file}",
    html_file
]
subprocess.run(args, capture_output=True)

if os.path.exists(html_file):
    os.remove(html_file)

print(f"Rendered Ultra-HD {png_file} - Size: {os.path.getsize(png_file)} bytes")
