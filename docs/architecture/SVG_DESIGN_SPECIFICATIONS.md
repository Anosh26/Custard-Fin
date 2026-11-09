# CustardFin — Master SVG Diagram Production Prompts

> **Audience:** Professional Vector Graphic Designers / Technical Illustrators  
> **Target Document:** CustardFin Academic Project Report (`report/main.pdf`)  
> **Design Theme:** Clean, modern, technical vector illustration (Stripe/Linear/AWS Architecture style). White canvas (`#ffffff`), subtle gradients, soft drop shadows (`#0f172a` at 6% opacity), crisp modern sans-serif typography (Inter / SF Pro / Segoe UI), dedicated color zoning per subsystem, and **zero overlapping arrows or text elements**.

---

## Global Design Guidelines & Color Palette

- **Home Edge Zone:** Lavender / Violet / Indigo (`#7c3aed`, `#6366f1`, `#4338ca`, bg `#faf5ff`)
- **Residential Firewall / CGNAT:** Crimson / Coral (`#dc2626`, `#ef4444`, `#991b1b`, bg `#fff1f2`)
- **Cloud Control Tier (GCP VM):** Cyan / Sky Blue (`#0284c7`, `#0369a1`, `#38bdf8`, bg `#f0f9ff`)
- **Direct LAN Streaming:** Emerald Green (`#059669`, `#10b981`, `#047857`, bg `#ecfdf5`)
- **Origin Shield & Cache:** Warm Amber / Gold (`#d97706`, `#f59e0b`, `#b45309`, bg `#fefce8`)
- **Typography:** Sans-serif, bold clear titles (18-20px), section headers (11-13px), body labels (9-10px), monospace code snippets (9px).
- **Invariants to respect:** Indian Rupee symbol `₹` for all monetary costs; strict separation of layout lanes to guarantee zero element overlap.

---

## 1. Diagram 1: Residential Carrier-Grade NAT (CGNAT) Traversal
- **File Name:** `diag_01_nat_traversal.svg` (and `.png`)
- **Report Location:** Chapter 2 (Literature Survey & Foundations), **Figure 2.1**
- **Dimensions:** `1200 x 520 px` (Aspect ratio ~21:9 landscape)
- **Concept:** Illustrate how home media servers hosted behind ISP Carrier-Grade NAT (RFC 6598 `100.64.0.0/10`) cannot accept inbound connections, and how CustardFin overcomes this via an outbound TLS 1.3 reverse tunnel to a public cloud relay VM, eliminating port-forwarding.

### Design Prompt / Specification
```text
Create a clean, widescreen 3-zone technical network architecture diagram titled:
"Carrier-Grade NAT (CGNAT) Traversal via Outbound Reverse Tunneling"

LAYOUT (Left to Right across 3 clearly separated vertical containers):
1. LEFT CONTAINER (Width 330px, Light Violet #faf5ff): "HOME EDGE (192.168.1.0/24)"
   - Top: Mini rack chassis / desktop server tower labeled "Home Server (edged)", with power LED and port ":8089".
   - Bottom Left: 3D purple storage cylinder labeled "NVMe Storage (Master HLS at Home)".
   - Bottom Right: Living room desktop monitor labeled "LAN Viewers (Sub-4ms Direct Wi-Fi)".
   - Direct green arrow between server and monitor labeled "LAN Bypass".
   - Bottom Banner: Green pill "ZERO INBOUND PORTS REQUIRED — CGNAT & Firewall Immune".

2. CENTER CONTAINER (Width 180px, Light Red #fff1f2): "ISP CGNAT BARRIER"
   - Subtitle: "RFC 6598 (100.64.0.0/10)".
   - Top Section: Red padlock/shield with white cross "✕".
     Box labeled "INBOUND BLOCKED — Direct TCP SYN Dropped, No Public Home IPv4, Port Forwarding Fails".
   - Bottom Section: Green circular checkmark "✓".
     Box labeled "OUTBOUND PASS — Persistent TLS 1.3 Conduit, Reverse Tunneling (FRP), Stateful NAT Mapping, Outbound Port :7000".

3. RIGHT CONTAINER (Width 530px, Light Cyan #f0f9ff): "PUBLIC CLOUD TIER (GCP VM / RELAY)"
   - Top Row (3 service boxes side-by-side):
     a. "FRP SERVER (frps) — Listens on :7000 (TLS), Multiplexed Tunnel Hub".
     b. "STREAM ROUTER (apid) — Subnet Proximity Router, 10-Stream Limiter Cap, Short-Lived HMAC Tokens".
     c. "ORIGIN SHIELD CACHE — S3 24h LRU Buffer, Single-Ingest Upload, Zero Permanent Storage".
   - Bottom Section: "Public Community Viewers (10x Load)" showing 5 device icons (Laptop, Smartphone 1, Smartphone 2, Tablet, Desktop) labeled Viewer 1 through Viewer 10.
   - Banner below devices: "⚡ Origin Shield Bandwidth Absorption: Viewer 1 triggers single pull via reverse tunnel (1x Uplink: ~2.5 MB); Viewers 2–10 served directly from Cloud Cache (90% Home Bandwidth Saved)".

CONNECTING HIGHWAY ARROWS (Clear separate horizontal lanes — DO NOT OVERLAP):
- Lane 1 (Top, Red Dashed Arrow): From Cloud toward Home, hitting the center red shield and bouncing back: "① Inbound SYN Dropped".
- Lane 2 (Bottom, Thick Solid Green Arrow with Dashes through Firewall): Originating from Home Edge, piercing cleanly through the "OUTBOUND PASS" section of the firewall, curving up into the cloud FRP Server: "② Outbound TLS 1.3 Tunnel (FRP)".
```

---

## 2. Diagram 2: Two-Tier Hybrid Edge-Cloud Architecture Overview
- **File Name:** `diag_02_architecture_overview.svg` (and `.png`)
- **Report Location:** Chapter 3 (System Architecture), **Figure 3.1**
- **Dimensions:** `1200 x 580 px`
- **Concept:** Master architectural blueprint showing the physical and logical decoupling of CustardFin into two cooperating tiers: Home Edge Tier (`apps/edge`) and Cloud Control Tier (`apps/api` + `apps/web`), connected strictly via outbound TLS and 650-byte metadata manifests.

### Design Prompt / Specification
```text
Create a high-impact, professional dual-tier cloud architecture diagram titled:
"CustardFin Distributed Hybrid Edge-Cloud System Architecture"
Subtitle: "Decoupled Edge Storage, Subprocess NVENC Pre-transcoding, Ephemeral Redis Liveness, and Origin Shield Playback Routing"

LAYOUT (Split horizontally into TWO distinct tiers):

TIER 1 (LEFT HALF, Width 540px, Dark Slate / Cyber Aesthetic #0f172a / #1e293b): "HOME EDGE TIER (apps/edge)"
- Subtitle: "Customer-Owned PC / NAS — Zero Cloud Storage Liability"
- Core Components inside Edge:
  1. "Hardware Ingest Engine":
     - Green NVIDIA GPU ASIC microchip with visible pins labeled "NVENC Hardware Encoder (322.3 fps / 10.7x realtime)".
     - Fallback icon: "CPU libx264 Fallback Engine".
     - Filmstrip graphic sliced into 4-second blocks labeled "RFC 8216 HLS (GOP g=120, 1080p/720p/480p)".
  2. "Data Sovereignty Store":
     - 3D Database cylinder labeled "Pure-Go SQLite (WAL Mode)" with badges "Private Watch History" and "Local Media Hashes".
  3. "Local & Remote Gateways":
     - "Standalone Mini-Jellyfin (:8089) — Sub-4ms Direct LAN Playback".
     - "FRP Tunnel Client (frpc) — Outbound TLS 1.3 Connection to Cloud Relay".
- Edge Bottom Callout: "Hardware Spec: RTX 4060 PC | Idle RAM < 20 MB | Cost: ₹0.00 / Free".

TIER 2 (RIGHT HALF, Width 580px, Crisp White / Cloud Blue #f8fafc / #0284c7): "CLOUD CONTROL TIER (apps/api & apps/web)"
- Subtitle: "Lightweight Stateless Coordination — ₹4,128.60 / month total cost"
- Core Components inside Cloud:
  1. "Web Presentation Layer":
     - Web browser frame labeled "Next.js 14 App Router (:3000)" showing Community Catalog and HLS.js Player with "Stats for Nerds" HUD.
  2. "Go Control Plane (apid :8080)":
     - REST API endpoints: /v1/catalog, /v1/playback, /v1/heartbeat, /v1/nodes.
     - "Stream Router & Admission Limiter (10-Stream Concurrency Cap)".
  3. "Ephemeral State Infrastructure":
     - Redis Red Cylinder: "Redis In-Memory (15s Heartbeat TTL Keys ONLY)".
     - PostgreSQL Blue Cylinder: "PostgreSQL 16 (Users, Nodes, 650-Byte Public Manifests)".
     - S3 Amber Shield: "Origin Shield (MinIO/S3 24h LRU Ephemeral Cache Bucket)".
  4. "Relay Infrastructure":
     - "FRP Server (frps :7000) — Reverse Tunnel Hub".

INTER-TIER COUPLING PIPELINES (Spanning across the middle gap):
1. Top Pipeline (Violet Arrow): "Lightweight JSON Manifest (~650 Bytes)" from Edge to PostgreSQL.
2. Middle Pipeline (Green Pulsing Arrow): "15-Second Heartbeat Beacon" from Edge to Redis.
3. Bottom Pipeline (Amber Thick Pipe): "Outbound FRP TLS Tunnel" carrying on-demand chunk pulls.
- Bottom Guarantee Callout: "Invariant 1: Zero video binaries stored permanently in cloud (100% on edge disks)".
```

---

## 3. Diagram 3: Edge Media Ingestion & HLS Transcoding Pipeline
- **File Name:** `diag_03_ingest_pipeline.svg` (and `.png`)
- **Report Location:** Chapter 3 (System Architecture), **Figure 3.2**
- **Dimensions:** `1200 x 480 px`
- **Concept:** Step-by-step sequential manufacturing line showing how raw media files on disk are detected, hashed via SHA-256, probed for GPU hardware capabilities, transcoded into 4-second multi-bitrate HLS segments, indexed in SQLite, and synced as JSON manifests.

### Design Prompt / Specification
```text
Create a linear, step-by-step pipeline infographic titled:
"Edge Media Ingestion, HLS Segment Slicing, and SQLite Metadata Indexing"
Subtitle: "Automated scanner discovers raw video, probes NVENC GPU, slices 4-second GOP-aligned HLS chunks, and syncs lightweight manifest"

LAYOUT (Horizontal processing sequence across 5 circular/hexagonal stations):

STATION 1: "Raw Media Discovery & Hashing"
- Graphic: Video file icon (MKV/MP4) with magnifying glass.
- Details:
  • Local directory watch: /media/*.mp4
  • Cryptographic SHA-256 binary hash generation
  • Duplicate elimination & identity binding

STATION 2: "Hardware Capabilities Probe"
- Graphic: NVIDIA green chip with active probe needle and CPU checkmark.
- Details:
  • Encoder probe: h264_nvenc → QSV/VAAPI → libx264
  • 1-second canary transcode test
  • Automatic CPU fallback on driver/VRAM fault

STATION 3: "Managed Subprocess Transcoder (FFmpeg)"
- Graphic: Filmstrip reel with razor blade cutting 4.0s pieces.
- Details:
  • Standardized 4-second segment cuts (-hls_time 4)
  • Strict GOP keyframe alignment: -g 120 -keyint_min 120
  • Multi-bitrate ladder: 1080p (5.0 Mbps), 720p (2.5 Mbps), 480p (1.0 Mbps)
  • Master playlist generation: master.m3u8

STATION 4: "Local SQLite Metadata Persistence"
- Graphic: 3D purple database barrel with WAL badge.
- Details:
  • WAL journal mode enabled for high-concurrency read/write
  • Segments, bitrates, runtimes, and local paths indexed
  • Private watch history kept strictly on local disk (Invariant 8)

STATION 5: "Public Manifest Generation & Cloud Sync"
- Graphic: JSON document glowing with cloud upload arrow.
- Details:
  • User toggles "Public" in Edge UI
  • Compact ~650-byte JSON manifest generated
  • Synced to Cloud API via PUT /v1/nodes/{id}/media/{id}/manifest
  • Video binaries NEVER uploaded (Invariant 1)

CONNECTORS: Smooth gradient arrows connecting Station 1 → 2 → 3 → 4 → 5 with step numbers.
```

---

## 4. Diagram 4: Ephemeral Heartbeat & Node Liveness Detection
- **File Name:** `diag_04_heartbeat_liveness.svg` (and `.png`)
- **Report Location:** Chapter 3 (System Architecture), **Figure 3.3**
- **Dimensions:** `1200 x 480 px`
- **Concept:** Illustrate the zero-polling, ephemeral Redis TTL heartbeat mechanism. Shows the dual state machine: Active Online Heartbeat (every 15s) vs. Sudden Disconnect Failover (Redis TTL key expiry after 20s, graceful UI amber badge).

### Design Prompt / Specification
```text
Create a dual-scenario comparison diagram titled:
"Ephemeral Heartbeat and Automated Node Liveness Detection"
Subtitle: "Edge transmits a 15-second heartbeat beacon refreshing a Redis TTL key; missed beacons cause automatic catalog de-listing"

LAYOUT (Two distinct horizontal scenario tracks, top and bottom):

SCENARIO A (TOP TRACK, Emerald Green Theme #ecfdf5 / #059669): "SCENARIO A: NORMAL ONLINE HEARTBEAT LOOP"
- Left: Edge daemon with green pulsing ECG cardiac waveform icon.
  Text: "Edge Daemon: Emits POST /v1/nodes/{id}/heartbeat every 15.0s with active stream count".
- Middle: Redis in-memory storage barrel with glowing green LED.
  Text: "Redis TTL Key: SET heartbeat:node-123 'online' EX 20". Key refreshed at t=15s (Never expires).
- Right: Web browser catalog showing crisp Green "ONLINE" badge next to video cards.
  Text: "Catalog Status: Instant playback authorized. Edge ready to stream".

SCENARIO B (BOTTOM TRACK, Crimson / Amber Theme #fff1f2 / #dc2626): "SCENARIO B: SUDDEN EDGE POWER LOSS / OFFLINE FAILOVER"
- Left: Edge PC with severed power plug and red flatline ECG wave.
  Text: "t = 0s: Home edge loses power or internet disconnects. Beacons cease immediately".
- Middle: Redis cylinder with countdown timer:
  • t = 15s: Expected beacon missed.
  • t = 20s: 20-second Redis TTL expires! Key automatically evicted from memory.
  Graphic: Trash can / eviction puff.
- Right: Web browser catalog showing Amber "NODE OFFLINE" warning badge on video cards.
  Text: "Graceful UI Degradation: Cached chunks continue playing; uncached media alerts viewer with 'Node Offline' notice instead of frozen player".

KEY CALLOUT: "Invariant 5 Upheld: Redis is used strictly for 15s heartbeat keys; no heavy background queue overhead".
```

---

## 5. Diagram 5: Dual-Path Stream Router & Origin Shield Single Ingest
- **File Name:** `diag_05_dual_path_shield.svg` (and `.png`)
- **Report Location:** Chapter 3 (System Architecture), **Figure 3.4**
- **Dimensions:** `1200 x 560 px`
- **Concept:** Visualizing the core traffic routing engine. A client request arrives at the Go stream router: if on the same LAN subnet, traffic routes directly to the local edge at <4ms latency (₹0.00 cost). If remote, traffic routes through the Origin Shield, which pulls each hot chunk exactly once over the reverse tunnel and serves up to 10 viewers from cache.

### Design Prompt / Specification
```text
Create a dynamic routing decision and traffic fan-out diagram titled:
"Dual-Path Stream Routing & Origin Shield Single Ingest"
Subtitle: "Subnet-aware proximity routing bypasses WAN for local clients; cloud cache shields home uplink from remote viewer fan-out"

LAYOUT:
1. LEFT: "HOME EDGE NODE (192.168.1.105)"
   - Edge PC Workstation (RTX 4060, NVMe HLS disk, Mini-Jellyfin on :8089, FRP client).
   - Residential Uplink Budget: 20 Mbps limit.

2. CENTER: "SMART STREAM ROUTER DECISION GATE (apid)"
   - Illuminated decision compass / code block:
     "if clientIP.InSubnet(edgeNode) { return DirectLAN_URL() } else { return OriginShield_URL() }"
   - Concurrency Check: "Admission Limiter: If active_streams >= 10, return HTTP 429 Too Many Requests (Node Busy)".

3. PATH A (Diverging Down/Left, Green #059669): "PATH A: DIRECT LAN BYPASS"
   - Green conduit straight from Edge to Living Room Smart TV / Tablet.
   - Metrics:
     • Latency (TTFB): < 3.2 ms (Local Gigabit Wi-Fi)
     • Cloud Egress Cost: ₹0.00 / Zero
     • WAN Bandwidth Used: 0 Bps
     • Offline Resilience: Works 100% even if ISP internet is completely cut.

4. PATH B (Diverging Right, Sky Blue & Amber): "PATH B: CLOUD ORIGIN SHIELD & MULTI-VIEWER FAN-OUT"
   - Origin Shield Cache (S3 24h LRU Bucket).
   - Step 1 (Cache Miss, Amber Conduit): Viewer 1 requests new chunk seg_042.ts.
     Go API pulls chunk ONCE over FRP tunnel from home edge (~2.5 MB). Writes to S3 cache.
   - Step 2 (Cache Hit, Blue Fan-Out): Viewers 2 through 10 request same chunk.
     Served directly from Cloud S3 Cache. Home edge uplink load: EXACTLY 0 BYTES.
   - Multi-device fan-out graphic (Laptops, phones, desktops).
   - Metrics Callout: "90% Upload Bandwidth Saved (1x Upload vs 10 Viewers) — Invariant 7 Upheld".
```

---

## 6. Diagram 6: Economic Cost Comparison (in Indian Rupees ₹)
- **File Name:** `diag_06_cost_comparison.svg` (and `.png`)
- **Report Location:** Chapter 5 (Cost Analysis & Feasibility), **Figure 5.1**
- **Dimensions:** `1200 x 520 px`
- **Concept:** Clear, authoritative financial comparison for a 100-user community library (5 TB storage, 9 TB monthly egress). Demonstrates the 95.89% cloud cost reduction in Indian Rupees (`₹`) by comparing Traditional Cloud VOD against CustardFin's hybrid architecture.

### Design Prompt / Specification
```text
Create an executive financial infographic and stacked bar chart titled:
"Monthly Operational Expenditure Comparison: Traditional Cloud VOD vs. CustardFin (INR ₹)"
Subtitle: "Baseline Benchmark: $1.00 USD = ₹84.00 INR | Workload: 100-User Community, 5 TB Video Library, 9 TB Monthly Egress"

LAYOUT (Two 3D stacked comparison towers side-by-side with a central savings callout):

TOWER 1 (LEFT, Tall Red/Coral Column): "TRADITIONAL CLOUD VOD (AWS S3 + CloudFront + EC2)"
- Total Height: High (~380px)
- Stacked segments with individual Rupee costs:
  • S3 Storage (5 TB @ ₹1.932/GB-mo): ₹9,660.00 / mo
  • Transcoding Compute (EC2 / MediaConvert): ₹11,760.00 / mo
  • Data Egress Bandwidth (9 TB @ ₹7.56/GB): ₹68,040.00 / mo (Dominant 68% segment)
  • Managed RDS Database & Control Tier: ₹10,920.00 / mo
- TOWER 1 TOTAL BADGE: "₹1,00,380.00 / month" (₹12.05 Lakhs / year | 3-Year TCO: ₹36.14 Lakhs).

CENTER MEDALLION (Glowing Emerald Green Badge):
- Circular seal with ribbon: "95.89% NET SAVINGS"
- Monthly Savings: "Save ₹96,251.40 / month"
- Annual Savings: "Save ₹11,55,016.80 / year"
- 3-Year TCO Savings: "Save ₹34.65 Lakhs over 3 Years"

TOWER 2 (RIGHT, Short Emerald Green Column): "CUSTARDFIN HYBRID ARCHITECTURE"
- Total Height: Very Short (~35px, 4.1% of Tower 1)
- Stacked segments with individual Rupee costs:
  • Video Storage (100% Edge Disks): ₹0.00 (Zero)
  • Transcoding Compute (Edge NVENC ASICs): ₹0.00 (Zero)
  • Data Egress Bandwidth (LAN + Bundled VPS): ₹0.00 (Zero)
  • Lightweight PostgreSQL VPS: ₹1,470.00 / mo
  • Redis Heartbeat Instance: ₹420.00 / mo
  • Origin Shield Buffer (100 GB S3): ₹193.20 / mo
  • FRP Reverse Tunnel Relay VM: ₹2,045.40 / mo
- TOWER 2 TOTAL BADGE: "₹4,128.60 / month" (₹49,543 / year | 3-Year TCO: ₹1.49 Lakhs).

BOTTOM EQUATION CALLOUT:
"Asymptotic Storage Proof: lim (N -> inf) C_custard(N) = O(1) <= ₹193.20/mo  vs.  Traditional O(N) -> inf"
```

---

## 7. Diagram 7: Hardware Transcoding Benchmark (NVENC vs. CPU)
- **File Name:** `diag_07_nvenc_benchmark.svg` (and `.png`)
- **Report Location:** Chapter 6 (Evaluation, Testing & Results), **Figure 6.1**
- **Dimensions:** `1200 x 480 px`
- **Concept:** High-impact hardware performance benchmarking dashboard comparing dedicated hardware-accelerated NVENC GPU encoding against traditional multi-threaded CPU software encoding (`libx264`) across 3 critical metrics: Transcoding Speed, Host CPU Utilization, and Thermal/Power Draw.

### Design Prompt / Specification
```text
Create a high-tech hardware benchmarking dashboard titled:
"Hardware Transcoding Benchmarks: NVENC GPU Acceleration vs. CPU libx264"
Subtitle: "Evaluated on Intel Core i7-10700K + NVIDIA GeForce RTX 4060 (8 GB) on 4K Master Asset Transcoded to 1080p60 HLS"

LAYOUT (3 distinct comparison metric pods side-by-side):

POD 1: "TRANSCODING THROUGHPUT (SPEED MULTIPLIER)"
- Graphic: Dual automotive-style tachometers / speedometers.
  • Speedometer A (Emerald Green #10b981): Needle pegged at high speed.
    Label: "NVENC Hardware ASIC: 322.3 fps (10.74x Realtime Speed)".
    Badge: "1-hour film encoded in 5.5 minutes".
  • Speedometer B (Slate/Gray #64748b): Needle near low end.
    Label: "CPU libx264 (Software): 33.6 fps (1.12x Realtime Speed)".
    Badge: "1-hour film takes 53.6 minutes".
- Pod Summary: "9.59x Transcoding Speedup with NVENC".

POD 2: "HOST CPU UTILIZATION DURING INGEST"
- Graphic: Dual circular donut percentage rings.
  • Donut A (Emerald Green #10b981): 4.8% ring filled.
    Label: "NVENC Mode: 4.8% Host CPU Load".
    Badge: "PC remains completely responsive for gaming/work".
  • Donut B (Red #ef4444): 96.4% ring filled (Near 100% pegged).
    Label: "CPU libx264 Mode: 96.4% CPU Load (All 16 Threads Pegged)".
    Badge: "System experiences severe stutter & thermal throttling".
- Pod Summary: "91.6% Host CPU Workload Reduction".

POD 3: "THERMAL EMISSIONS & PACKAGE POWER"
- Graphic: Dual glass mercury thermometers.
  • Thermometer A (Cool Blue/Green #0284c7): Mercury level low at 61°C.
    Label: "GPU Package: 61°C | 42W Power Draw".
  • Thermometer B (Boiling Crimson #dc2626): Mercury level near boiling at 89°C.
    Label: "CPU Package: 89°C | 95W Power Draw".
- Pod Summary: "28°C Cooler & 53W Less Power Consumption".

BOTTOM FAILOVER BADGE:
"Automated Graceful Fallback: If NVENC VRAM limits are reached, edged falls back to CPU encoding within 420 ms without failing the ingestion job".
```

---

## 8. Summary of Figures in Report

| Diagram # | File Name | Placement | Title |
|---|---|---|---|
| **Diag 1** | `diag_01_nat_traversal.svg` | Chapter 2, Fig 2.1 | Carrier-Grade NAT Traversal via Outbound Reverse Tunneling |
| **Diag 2** | `diag_02_architecture_overview.svg` | Chapter 3, Fig 3.1 | CustardFin Distributed Hybrid Edge-Cloud System Architecture |
| **Diag 3** | `diag_03_ingest_pipeline.svg` | Chapter 3, Fig 3.2 | Edge Media Ingestion, HLS Segment Slicing, and SQLite Indexing |
| **Diag 4** | `diag_04_heartbeat_liveness.svg` | Chapter 3, Fig 3.3 | Ephemeral Heartbeat and Automated Node Liveness Detection |
| **Diag 5** | `diag_05_dual_path_shield.svg` | Chapter 3, Fig 3.4 | Dual-Path Stream Routing & Origin Shield Single Ingest |
| **Diag 6** | `diag_06_cost_comparison.svg` | Chapter 5, Fig 5.1 | Monthly Operational Expenditure Comparison (INR ₹) |
| **Diag 7** | `diag_07_nvenc_benchmark.svg` | Chapter 6, Fig 6.1 | Hardware Transcoding Benchmarks (NVENC vs CPU) |

