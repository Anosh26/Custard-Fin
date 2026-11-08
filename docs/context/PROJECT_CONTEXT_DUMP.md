# CustardFin Comprehensive Project Context Dump

> **Generated:** 2026-11-09  
> **Course:** MCA Cloud Computing Lab (End-Semester Project)  
> **Master Report Title:** Cloud Computing Project Report  
> **Repository:** `https://github.com/Anosh26/Custard-Fin.git` (`origin/master`)

---

## 1. Executive Mission & Architectural Anchor
CustardFin is a free, open-source hybrid edge-cloud media streaming and community catalog platform that resolves the fundamental trade-off between private self-hosted media servers (Jellyfin, Plex) and centralized public video-on-demand platforms (YouTube, AWS CloudFront).

### The Invariant Core
> *"Video binaries remain permanently on the content owner's home machine; the cloud holds only the metadata catalog and ephemeral stream coordination."*

Marking a video **Public** syncs a lightweight ~650-byte JSON manifest, **never the video binary**. When remote viewers stream public media, an in-cloud **Origin Shield** cache ingests hot 4-second HLS segments on-demand exactly once, absorbing viral viewer fan-out while preserving residential home uplinks.

---

## 2. Academic Project Credentials
- **Course / Degree:** Master of Computer Applications (MCA) — Cloud Computing Lab
- **Official Report Title:** **`Cloud Computing Project Report`**
- **Project Candidates (Identified Strictly by Seat Numbers):**
  - **Anish Ghadi** — Seat No.: **`25P0320009`**
  - **Anosh Chodankar** — Seat No.: **`25P0320010`**
- **Academic Deliverables:**
  - Master Academic Report: [`report/main.pdf`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/report/main.pdf) (62 pages, 2.41 MB, 0 errors, full IEEEtran references).
  - Modular LaTeX Source: [`report/main.tex`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/report/main.tex) and [`report/sections/*.tex`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/report/sections/).
  - Interactive Presentation Deck: [`docs/presentation/index.html`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/presentation/index.html) and [`docs/presentation/CustardFin_Presentation.pptx`](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/presentation/CustardFin_Presentation.pptx).

---

## 3. Current Milestone & Live Deployment Status
- **Current Milestone:** **Milestone 10: Hardening, Demo Profiles, Live Cloud Deployment & Final Report** (**COMPLETE & VERIFIED**).
- **All 10 Milestones Implemented:**
  1. *Milestone 1:* Contracts, repo skeleton, CI, documentation scaffolding, memory context.
  2. *Milestone 2:* Edge node subsystem (`apps/edge`), SQLite WAL, scanner, FFmpeg subprocess, local mini-Jellyfin streamer.
  3. *Milestone 3:* Cloud Control API (`apps/api`), PostgreSQL migrations, node registration, manifest storage.
  4. *Milestone 4:* Public toggle, manifest sync worker, Next.js community catalog UI.
  5. *Milestone 5:* FRP reverse tunnel, CGNAT traversal, authenticated TLS proxying.
  6. *Milestone 6:* Ephemeral 15s Redis heartbeat, TTL key expiration, offline UI indicators.
  7. *Milestone 7:* Origin Shield S3-compatible pull-through cache, HMAC signed tokens, stream router.
  8. *Milestone 8:* Admission control (10-stream cap, HTTP 429), adaptive bitrate, first-chunk prefetch.
  9. *Milestone 9:* Hardware-accelerated NVENC encoder probe, 1-second canary test, automatic CPU fallback, benchmark suite.
  10. *Milestone 10:* Multi-stage Dockerfiles, GCP cloud VM deployment, 10-viewer load simulation scripts, academic report, and runbooks.

### Multi-Topology Deployment Status
1. **Live Cloud Tier (Google Cloud Platform):**
   - Host: Debian 12 Compute Engine VM at `34.132.188.14`.
   - Running Services: Caddy reverse proxy, Next.js web application (`:3000`), Go API daemon (`:8080`), PostgreSQL (`:5432`), Redis (`:6379`), MinIO S3 Origin Shield (`:9000`), and FRPS tunnel server (`:7000`, `:7001`).
2. **Physical Home Edge Node:**
   - Host: Fedora Workstation PC with NVIDIA GeForce RTX 4060 (8 GB GDDR6).
   - Running Services: `edged` daemon (`:8089`), `frpc` client maintaining outbound TLS connection to GCP relay VM.
3. **Live Mobile & Cross-Device Verification:**
   - Verified end-to-end 1080p streaming on Android mobile (Brave browser) over cellular 5G network.
   - Verified sub-4ms direct LAN playback on local Wi-Fi devices.
   - Verified Origin Shield single ingest and cache hits via live player "Stats for Nerds" HUD.

---

## 4. Architectural Invariants (Non-Negotiable)
1. **Zero Permanent Cloud Binaries:** Only temporary 24-hour LRU cache for hot media chunks.
2. **Outbound-Only Connections from Edge:** Never require open inbound router ports; seamlessly traverses residential Carrier-Grade NAT (CGNAT) via FRP reverse tunnel.
3. **100% Free & Open-Source Software (FOSS):** Zero proprietary SDKs, zero commercial licenses.
4. **Technology Stack:** Go for all backends (`apps/edge`, `apps/api`); Next.js App Router for frontend UI; pure Node runtime (no Bun).
5. **Strict Redis Boundary:** Redis is used exclusively for 15-second heartbeat TTL keys. Background workers are native Go goroutines with state in PostgreSQL.
6. **Pre-Encoded RFC 8216 HLS:** 4-second chunking with strict GOP keyframe alignment (`-g 120`), not wasteful on-the-fly live transcoding.
7. **Origin Shield Guarantee:** Home edge uploads each hot video segment **exactly once**; subsequent viewers are served by the cloud cache.
8. **Edge Privacy & Sovereignty:** Private media and watch history reside exclusively on edge SQLite databases.
9. **HMAC Cryptographic Security:** All chunk delivery is authenticated via short-lived HMAC signed URLs.
10. **10-Stream Concurrency Limit:** Automated admission limiter caps each edge node to 10 concurrent distinct streams, returning HTTP `429 Too Many Requests` (`Retry-After: 30`) to protect home broadband uplinks.

---

## 5. Cost Analysis & Financial Proof (in Indian Rupees ₹)
- **Baseline Exchange Benchmark:** `$1.00 USD = ₹84.00 INR`.
- **Unit Pricing Benchmarks:**
  - S3 Cloud Object Storage: `₹1.932 / GB-month` (`$0.023`).
  - Cloud Egress Bandwidth: `₹7.56 / GB` (`$0.090`).
  - Cloud Transcode Compute: `₹1.26 / minute` (`$0.015`).

### Financial Comparison Matrix (100-User Community, 5 TB Library, 9 TB Monthly Egress)

| Cost Dimension | Traditional Cloud VOD (AWS/YouTube) | CustardFin Hybrid Architecture | Monthly Savings | % Reduction |
|---|---|---|---|---|
| **Raw Media Storage** | `₹9,660.00` (5 TB @ S3) | **`₹0.00`** (100% Edge NVMe Disks) | `+₹9,660.00` | `100.0%` |
| **Transcoding Compute** | `₹11,760.00` (Cloud vCPU/EC2) | **`₹0.00`** (Edge NVENC ASICs) | `+₹11,760.00` | `100.0%` |
| **Data Egress Bandwidth** | `₹68,040.00` (9 TB @ `₹7.56/GB`) | **`₹0.00`** (LAN + Bundled VPS Bandwidth) | `+₹68,040.00` | `100.0%` |
| **Control Tier & Metadata** | `₹10,920.00` (AWS RDS Postgres) | `₹1,470.00` (Lightweight VPS Postgres) | `+₹9,450.00` | `86.5%` |
| **Heartbeat In-Memory Cache** | `₹0.00` (N/A) | `₹420.00` (Redis Cache Instance) | `-₹420.00` | N/A |
| **Origin Shield Temporary Cache**| `₹0.00` (N/A) | `₹193.20` (100 GB 24h LRU Buffer) | `-₹193.20` | N/A |
| **Reverse Tunnel Relay VM** | `₹0.00` (N/A) | `₹2,045.40` (FRP Server Instance) | `-₹2,045.40` | N/A |
| **TOTAL MONTHLY COST** | **`₹1,00,380.00 / month`** | **`₹4,128.60 / month`** | **`+₹96,251.40`** | **`95.89%`** |
| **ANNUAL EXPENDITURE** | **`₹12,04,560.00 / year`** | **`₹49,543.20 / year`** | **`+₹11,55,016.80`** | **`95.89%`** |
| **3-YEAR TOTAL TCO** | **`₹36,13,680.00`** | **`₹1,48,629.60`** | **`+₹34,65,050.40`** | **`95.89%`** |

### Asymptotic Storage Scaling Proof
$$\lim_{N \to \infty} C_{\text{custard, storage}}(N) = O(1) \le \text{₹193.20/month} \quad (\text{vs. } O(N) \to \infty \text{ on traditional cloud})$$

---

## 6. Creative Visual Graphics & Metaphors Suite
All generic flowchart boxes were redesigned into rich, domain-specific visual metaphors:

1. **Figure 2.1 — Residential NAT Traversal (`diag_01_nat_traversal.png`):**
   - 3D NVMe cylinder, server chassis with blinking LEDs, monitor on stand, physical brick firewall with padlock shield (inbound SYN dropped $\times$, outbound TLS pass $\checkmark$), cloud, laptop & mobile silhouettes.
2. **Figure 3.1 — Two-Tier Hybrid Architecture (`diag_02_architecture_overview.png`):**
   - 3D storage cylinders, green NVIDIA GPU ASIC with golden pins (NVENC 322+ fps), film strip with sprocket cuts, server chassis, monitor, firewall shield, cloud contour with database barrels, and remote devices.
3. **Figure 3.2 — Local Ingestion Pipeline (`diag_03_ingest_pipeline.png`):**
   - Video file icon, SHA-256 binary key, green NVENC chip with active probe, film strip with 4s cuts (GOP $g=120$), 3D SQLite WAL database cylinder, TV/monitor.
4. **Figure 3.3 — Heartbeat Liveness Tracking (`diag_04_heartbeat_liveness.png`):**
   - Active green pulsing ECG heartbeat waveform, 15s timer, Redis RAM cylinder, global online beacon vs. severed flatline ECG, 30s grace stopwatch, automatic memory eviction trash drop.
5. **Figure 3.4 — Dual-Path Routing & Origin Shield (`diag_05_dual_path_shield.png`):**
   - RTX 4060 PC tower, Stream router code decision, Direct LAN bypass with Smart TV and `₹0.00` egress, Origin Shield S3 LRU cache, 10-viewer fan-out with laptops/smartphones/tablets and 90% uplink saved.
6. **Figure 5.1 — Cost Comparison in Indian Rupees (`diag_06_cost_comparison.png`):**
   - 3D stacked cost comparison towers in Indian Rupees (`₹1,00,380` vs `₹4,129`), circular 95.9% savings medallion, and zero-cost checkmarks.
7. **Figure 6.1 — Hardware NVENC vs CPU Benchmark (`diag_07_nvenc_benchmark.png`):**
   - Tachometer speedometers (322.3 fps vs 33.6 fps), circular CPU donut rings (4.8% vs 96.4%), glass mercury thermometers (61°C / 42W vs 89°C / 95W).

### High-Resolution Live UI Screenshots
- `report/figures/ui_catalog_feed.png` — Public Community Catalog Feed.
- `report/figures/ui_player_stats.png` — HLS Video Player with "Stats for Nerds" HUD.
- `report/figures/ui_edge_library.png` — Edge Library & Public/Private Sync Modal.
- `report/figures/ui_mini_jellyfin.png` — Standalone Mini-Jellyfin Zero-Cloud Local Player.
- `report/figures/ui_dev_settings.png` — Developer Settings & Simulation Console.

---

## 7. Subsystem Ownership & Division of Work
- **Anish Ghadi (Seat No.: 25P0320009):**
  - Edge Node Ingestion Engine, Subprocess FFmpeg Transcoder, NVENC GPU Acceleration with 1s Canary Probe & CPU Fallback, SQLite WAL Database Architecture, Mini-Jellyfin Zero-Cloud Offline Web Hub, Hardware Benchmarks.
- **Anosh Chodankar (Seat No.: 25P0320010):**
  - Cloud Control Tier REST API (Go `net/http`), PostgreSQL Database Design & `pgx` Adapters, Redis Heartbeat Liveness State Machine, FRP Reverse Tunnel Integration, Origin Shield S3 Cache, Next.js App Router Web UI.
- **Shared Collaborative Architecture:**
  - 10 Non-Negotiable Invariants Definition, 3 KB JSON Manifest Schema, Stream Router Dual-Path Decision Logic, Admission Limiter (10-Stream Cap), 10-Viewer Automated Load Simulation Suite, Modular Academic LaTeX Report (`report/main.pdf`).

---

## 8. Version Control State
- **Git Branch:** `master`
- **Latest Commit:** `dfe6dab` — *"docs(report): complete 62-page academic report with creative graphics and INR conversion"*
- **GitHub Remote:** `https://github.com/Anosh26/Custard-Fin.git`
- **Working Tree:** Pristine and clean (`git status` reports zero uncommitted files).

