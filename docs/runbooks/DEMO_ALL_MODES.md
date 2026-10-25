# CustardFin — Complete Demo Presentation Runbook for All 4 Modes

> **Project**: CustardFin — Hybrid Personal Media Streaming & Community Catalog  
> **Authors**: Anish Ghadi (25P0320009), Anosh Chodankar (25P0320010)  
> **Course**: MCA Cloud Computing Lab (Semester-End Evaluation & Viva)  
> **Target Audience**: Evaluation Committee, External Examiners, and Lab Evaluators  

---

## 1. Executive Summary & 60-Second Opening Pitch

When your turn to present begins, deliver this exact opening statement:

> *"Good morning, Professors. We are presenting **CustardFin**, a decentralized hybrid media platform that bridges personal home media servers like Plex/Jellyfin with community discovery platforms like YouTube, at near-zero cloud operating cost.*
>
> *Our design is governed by **Invariant 1**: **Video binaries stay on the owner's home machine; the cloud holds only metadata and coordination.** When a user marks a video as Public, only a ~3 KB JSON manifest is synchronized to the cloud. Hot video chunks are pulled to an Origin Shield cache only when viewers stream them, cutting cloud storage and bandwidth costs by over 90%.*
>
> *CustardFin is built with strict Hexagonal Architecture in Go and Next.js, allowing it to run across four distinct execution topologies depending on the deployment environment."*

---

## 2. Mode 1: Pure Local (Zero Docker / Standalone Mode)

> **Best For**: Quick demonstration on an offline laptop without Docker installed or when container runtimes cannot be run.

### Architecture
All tiers run natively on the host PC. The Cloud API (`apid`) activates its built-in in-memory fallback stores (`ports.ChunkCache`, `ports.HeartbeatStore`, in-memory catalog repository).

```
┌────────────────────────────────────────────────────────┐
│                   SINGLE HOST MACHINE                  │
│                                                        │
│  [Next.js Web UI:3000] ──HTTP──▶ [Cloud API:8080]      │
│                                   (In-Memory Stores)   │
│                                           │ HTTP       │
│                                           ▼            │
│                             [Home Edge Daemon:8090]    │
│                             (SQLite WAL, ./media,      │
│                              RTX 4060 NVENC GPU)       │
└────────────────────────────────────────────────────────┘
```

### Pre-Demo Checklist
- Open **3 Terminal Tabs** in the project root.
- Ensure FFmpeg is available in system `PATH`.

---

### Step-by-Step Presentation Script (Mode 1)

#### Step 1: Synthesize Media Fixtures & Run GPU Benchmark (Tab 1)
```powershell
powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1
powershell -ExecutionPolicy Bypass -File scripts/bench-gpu.ps1
```
*(On Linux/macOS: `bash scripts/seed-media.sh && bash scripts/bench-gpu.sh`)*

**🗣️ What to Say:**
> *"First, we generate standardized test video fixtures and benchmark our hardware transcoding pipeline (Invariant 6). Our edge adapter detects the NVIDIA GeForce RTX 4060 GPU and achieves **16.5x realtime encoding speed** using `h264_nvenc`.*
>
> *Crucially, notice **4.0-second keyframe alignment**: every segment boundary is exact, enabling seamless adaptive bitrate switching without player rebuffering."*

#### Step 2: Launch Cloud API in In-Memory Mode (Tab 1)
```powershell
$env:API_ADDR = ":8080"
$env:DATABASE_URL = "memory"
$env:REDIS_ADDR = "memory"
$env:S3_BUCKET = ""
$env:TUNNEL_BASE_URL = "http://localhost:8090"

go run ./apps/api/cmd/apid
```

**🗣️ What to Say:**
> *"We start `apid`. Notice the logs: zero external databases or Redis instances are required. Thanks to hexagonal ports and adapters, in-memory repositories seamlessly satisfy the catalog, heartbeat, and Origin Shield interfaces for lightweight standalone operation."*

#### Step 3: Launch Home Edge Daemon (Tab 2)
```powershell
$env:EDGE_ADDR = ":8090"
$env:EDGE_MEDIA_DIR = "./media"
$env:EDGE_TRANSCODE_DIR = "./transcoded"
$env:EDGE_DB_PATH = "./data/edge.db"
$env:EDGE_CLOUD_URL = "http://localhost:8080"

go run ./apps/edge/cmd/edged
```

**🗣️ What to Say:**
> *"In Tab 2, `edged` starts up:*
> 1. *It initializes a pure-Go SQLite database with WAL journal mode (Invariant 8) so watch history stays private.*
> 2. *It scans `./media`, extracts metadata via FFprobe, and computes SHA-256 content hashes.*
> 3. *It publishes a **~1.4 KB JSON manifest** to the cloud API. The 100 MB video file stays on the local disk."*

#### Step 4: Launch Web UI & Playback (Tab 3)
```powershell
cd apps/web
npm run dev
```

**🌐 Browser Action:**
1. Open `http://localhost:3000`.
2. Open DevTools (`F12`), Network tab, filter by `.ts`.
3. Click the video card to stream.
4. Point out the signed URL (`?token=...&expires=...`) and 4-second chunk duration.

#### Step 5: Run Concurrency & Admission Simulation (Tab 1)
Stop `apid` briefly or open a 4th tab and run:
```powershell
powershell -ExecutionPolicy Bypass -File scripts/load-10-viewers.ps1
```

**🗣️ What to Say:**
> *"The simulation verifies our 10-viewer Origin Shield cache offload (90% home bandwidth savings), dynamic 1080p ladder pruning under load, and the 10-stream per-node cap (HTTP 429 Retry-After: 30)."*

---

## 3. Mode 2: Local Containerized (Dev Profile with MinIO)

> **Best For**: Single-laptop presentation demonstrating the full open-source cloud stack (PostgreSQL, Redis, MinIO S3, FRP relay) in Docker.

### Architecture
```
┌────────────────────────────────────────────────────────┐
│               DOCKER CONTAINERS (CLOUD)                │
│                                                        │
│  [PostgreSQL:5432]     [Redis:6379]                    │
│  [MinIO Cache:9000]    [FRP Server:7000]               │
└───────────────────────────▲────────────────────────────┘
                            │ Outbound Tunnel (:7000)
                            │ & Manifest Sync (:8080)
┌───────────────────────────┴────────────────────────────┐
│              HOST MACHINE (HOME EDGE NODE)             │
│                                                        │
│  - Media Directory (./media)                           │
│  - SQLite Edge DB (WAL mode, private watch history)    │
│  - FFmpeg Transcoder (RTX 4060 NVENC GPU acceleration) │
│  - Edge Daemon (edged) & Cloud API (apid)              │
│  - Next.js Web Frontend (apps/web on :3000)            │
└────────────────────────────────────────────────────────┘
```

### Pre-Demo Checklist
- Open **Docker Desktop** (ensure Docker Engine is healthy).
- Open **3 Terminal Tabs** in the project root.

---

### Step-by-Step Presentation Script (Mode 2)

#### Step 1: Boot Cloud Stack in Docker (Tab 1)
```powershell
make dev-up
docker compose -f deploy/compose.dev.yml ps
```

**🖥️ Checkpoint on Screen**: All 4 containers (`custardfin-cache`, `custardfin-frps`, `custardfin-postgres`, `custardfin-redis`) are `Up (healthy)`.

**🗣️ What to Say:**
> *"With one command, we spin up our complete cloud infrastructure in Docker:*
> - ***PostgreSQL 16***: *Relational catalog for registered nodes, users, and manifests.*
> - ***Redis 7***: *Used strictly for 15-second ephemeral node heartbeats (Invariant 5).*
> - ***MinIO (S3)***: *The temporary 24-hour Origin Shield pull-through cache.*
> - ***FRP Relay Server***: *Allows the edge node to bypass CGNAT with outbound-only TLS connections (Invariant 2)."*

#### Step 2: Seed Fixtures & GPU Benchmark (Tab 1)
```powershell
powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1
powershell -ExecutionPolicy Bypass -File scripts/bench-gpu.ps1
```

#### Step 3: Launch Cloud API Connected to Docker Services (Tab 1)
```powershell
$env:API_ADDR = ":8080"
$env:DATABASE_URL = "postgres://custard:custard_secret@localhost:5432/custardfin?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
$env:S3_ENDPOINT = "http://localhost:9000"
$env:S3_BUCKET = "custardfin-origin-cache"
$env:S3_ACCESS_KEY = "minioadmin"
$env:S3_SECRET_KEY = "minioadmin"

go run ./apps/api/cmd/apid
```

#### Step 4: Launch Home Edge Daemon (Tab 2)
```powershell
go run ./apps/edge/cmd/edged
```

**🖥️ Point Out in Edge Logs**:
- `sqlite database initialized (WAL mode enabled)`: Watch progress is kept locally (Invariant 8).
- `sync worker published manifest bytes=1420 is_public=true`: **Only 1.4 KB was sent to the cloud**, NOT the video file!
- `heartbeat sent to cloud active_streams=0 gpu=NVENC`: Reports liveness to Redis every 15s.

#### Step 5: Web UI Playback & DevTools Inspection (Tab 3)
```powershell
cd apps/web
npm run dev
```

1. Navigate to `http://localhost:3000`.
2. Open DevTools (`F12`), Network tab, filter by `.ts`.
3. Click the video card to stream.
4. Point out `X-Cache: MISS` on first chunk (pulled on-demand via tunnel from edge), and `X-Cache: HIT` on subsequent requests (served from MinIO cache).

#### Step 6: 10-Viewer Load & Concurrency Simulation (Tab 1 or new tab)
```powershell
powershell -ExecutionPolicy Bypass -File scripts/load-10-viewers.ps1
```

#### Step 7: Teardown
```powershell
make dev-down
```

---

## 4. Mode 3: Floci Local Cloud Emulator (Multi-PC Lab Setup)

> **Best For**: College lab presentation across 2 or 3 distinct computers (e.g. Lab PC A = Cloud, Student Laptop B = Home Node, Phone/Laptop C = Viewer) on the same Wi-Fi or lab switch.

### Architecture
```
   ┌─────────────────────────────────────────────────────────────┐
   │            MACHINE A: CLOUD CONTROL PLANE & RELAY           │
   │                  (IP: e.g. 192.168.1.100)                   │
   │                                                             │
   │  Docker Compose:                                            │
   │   - Floci S3 Emulator (Port 4566) - Origin Shield           │
   │   - PostgreSQL 16 (Port 5432) - Catalog                     │
   │   - Redis 7 (Port 6379) - Node Heartbeats                   │
   │   - FRP Server (Port 7000 & 8000-8020) - Relay Tunnel       │
   │                                                             │
   │  Host Services:                                             │
   │   - Cloud API (apid) on 0.0.0.0:8080                        │
   │   - Web Frontend (Next.js) on 0.0.0.0:3000                  │
   └──────────────────────────────▲──────────────────────────────┘
                                  │
         ┌────────────────────────┴────────────────────────┐
         │                                                 │
         │ Outbound-only FRP Tunnel (:7000)                │ HTTP Streams (:3000)
         │ & Manifest Sync (:8080)                         │
         │                                                 │
┌────────┴────────────────────────┐       ┌────────────────┴──────────────┐
│  MACHINE B: HOME EDGE NODE      │       │  MACHINE C: VIEWER CLIENT     │
│  (Student PC with RTX 4060 GPU) │       │  (Laptop, Mobile Phone, Tablet│
│  IP: e.g. 192.168.1.105         │       │  IP: Any IP on LAN or Hotspot)│
│                                 │       │                               │
│  - Media library on hard drive  │       │  - Connects to:               │
│  - NVENC GPU Transcoding (16x)  │       │    http://192.168.1.100:3000  │
│  - SQLite (private watch data)  │       │  - Zero plugins; standard HLS │
│  - Outbound-only connection     │       │  - Streams via Origin Shield  │
└─────────────────────────────────┘       └───────────────────────────────┘
```

### Pre-Demo Checklist
- Connect all machines to the same Wi-Fi router, lab LAN switch, or phone hotspot.
- On Machine A, run `ipconfig` (Windows) or `hostname -I` (Linux) to find its LAN IP (e.g., `192.168.1.100`).

---

### Step-by-Step Presentation Script (Mode 3)

#### Step 1: Start Floci Cloud Stack on Machine A (Host)
In Machine A Terminal 1:
```powershell
make floci-up
```
*(Spins up Floci MIT AWS S3 emulator on port 4566, PostgreSQL, Redis, and FRP relay).*

#### Step 2: Start Cloud API & Web UI on Machine A
In Machine A Terminal 2 (Cloud API):
```powershell
$env:API_ADDR = "0.0.0.0:8080"
$env:DATABASE_URL = "postgres://custard:custard_secret@localhost:5432/custardfin?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
$env:S3_ENDPOINT = "http://localhost:4566"
$env:S3_BUCKET = "custardfin-origin-cache"
$env:S3_ACCESS_KEY = "test"
$env:S3_SECRET_KEY = "test"
$env:S3_REGION = "us-east-1"
$env:TUNNEL_BASE_URL = "http://localhost:8000"

go run ./apps/api/cmd/apid
```

In Machine A Terminal 3 (Web UI):
```powershell
cd apps/web
npm run dev -- -H 0.0.0.0 -p 3000
```

#### Step 3: Start Home Edge Node on Machine B (RTX 4060 Laptop)
On Machine B (replace `192.168.1.100` with Machine A's IP):
```powershell
$CLOUD_IP = "192.168.1.100"

$env:EDGE_ADDR = ":8090"
$env:EDGE_MEDIA_DIR = "./media"
$env:EDGE_TRANSCODE_DIR = "./transcoded"
$env:EDGE_DB_PATH = "./data/edge.db"
$env:EDGE_CLOUD_URL = "http://${CLOUD_IP}:8080"
$env:EDGE_NODE_ID = "550e8400-e29b-41d4-a716-446655440001"
$env:EDGE_NODE_TOKEN = "custardfin_demo_node_token"

# Outbound-only FRP reverse tunnel
$env:EDGE_FRP_ENABLED = "true"
$env:EDGE_FRP_SERVER_ADDR = $CLOUD_IP
$env:EDGE_FRP_SERVER_PORT = "7000"
$env:EDGE_FRP_AUTH_TOKEN = "custardfin_demo_frp_token"

go run ./apps/edge/cmd/edged
```

**🗣️ What to Say (Crucial Viva Moment):**
> *"Professors, observe Machine B: it is acting as a homeowner's personal computer. Notice that **no inbound router ports were opened on Machine B** (Invariant 2). It established an outbound TLS tunnel to Machine A on port 7000.*
>
> *When Machine B scanned its hard drive, it synced only a **3 KB JSON manifest** to Machine A's PostgreSQL database. The video file did not move."*

#### Step 4: Stream from Machine C (Any Phone or Tablet)
1. Pick up a smartphone or open a second student's laptop connected to the Wi-Fi.
2. Open Chrome/Safari and visit: `http://192.168.1.100:3000`.
3. Tap on the video to play:
   - The video plays seamlessly on the phone.
   - Machine A pulls the hot chunk once through the tunnel from Machine B and caches it in Floci S3.
   - If another student opens the same video on their phone, Machine B does zero work—Machine A serves it directly from Floci S3 cache!

#### Step 5: Teardown Machine A
```powershell
make floci-down
```

---

## 5. Mode 4: Google Cloud Platform (GCP Live Cloud Way)

> **Best For**: Final evaluation demonstrating live public internet traversal, real CGNAT bypass, and production TLS.

### Architecture
- **Cloud Control Plane**: Google Cloud Compute Engine VM (`e2-standard-2`, Ubuntu 22.04 LTS, Static Public IP) running `deploy/compose.demo.yml` with Caddy reverse proxy on ports 80 and 443.
- **Home Edge Node**: Student laptop running on home/campus Wi-Fi behind CGNAT, connected across the public internet to the GCP VM.

---

### Step-by-Step Presentation Script (Mode 4)

#### Step 1: Pre-Deploy the GCP VM (Done in Advance)
On the GCP VM (via SSH):
```bash
git clone https://github.com/Anosh26/Custard-Fin.git
cd Custard-Fin
cp deploy/env.example .env
docker compose -f deploy/compose.demo.yml up -d --build
curl http://localhost:8080/healthz
```

#### Step 2: Connect Home Edge Node from Student Laptop (Live Demo)
On your laptop (in college lab or at home):
```powershell
$GCP_IP = "34.X.X.X"  # Replace with your GCP VM Public IP

$env:EDGE_CLOUD_URL = "http://${GCP_IP}:8080"
$env:EDGE_FRP_ENABLED = "true"
$env:EDGE_FRP_SERVER_ADDR = $GCP_IP
$env:EDGE_FRP_SERVER_PORT = "7000"
$env:EDGE_FRP_AUTH_TOKEN = "custardfin_demo_frp_token"
$env:EDGE_NODE_ID = "550e8400-e29b-41d4-a716-446655440001"
$env:EDGE_NODE_TOKEN = "custardfin_demo_node_token"

go run ./apps/edge/cmd/edged
```

#### Step 3: Live Audience Participation (Ask the Professor to Open on Their Phone)
**🗣️ What to Say:**
> *"Professors, please open your smartphone browsers and navigate to `http://<GCP_EXTERNAL_IP>`.*
>
> *You are viewing our Next.js community catalog served from our Google Cloud VM in Mumbai. When you tap play on that video:*
> 1. *GCP checks its Origin Shield S3 cache.*
> 2. *On cache miss, GCP pulls the 4-second chunk through the reverse tunnel from my laptop sitting right here on college Wi-Fi.*
> 3. *GCP caches the chunk in S3 and serves it to your phone.*
> 4. *When the next professor taps play, GCP serves it directly from cache without touching my laptop's Wi-Fi upload!*
>
> *At no point did we open a port on college Wi-Fi, and at no point did we store the entire movie in Google Cloud storage."*

---

## 6. Viva Q&A Defense Card (Memorize These 10 Answers)

| # | Anticipated Viva Question | Winning Technical Answer |
|---|---|---|
| **1** | **Why not just upload all videos to YouTube or S3?** | Storing terabytes of personal media on S3 costs $0.023/GB/month plus egress fees. CustardFin keeps files on home hard drives at $0 cloud cost. Cloud S3 only acts as a temporary 24-hour auto-purging cache for hot chunks currently being watched. |
| **2** | **How do you bypass CGNAT without opening router ports?** | The home edge daemon initiates an **outbound-only TLS reverse tunnel** (`frpc`) to the cloud relay server (`frps`). Because the connection originates from inside the home, NAT firewalls permit it without inbound port forwarding. |
| **3** | **Why Redis only for heartbeats?** | Using Redis for task queues introduces state duplication and distributed race conditions. Background tasks run as standard Go goroutines with transactional state stored in PostgreSQL. Redis is used solely for atomic 15s TTL keys (`heartbeat:{node_id}`). |
| **4** | **Why 4-second HLS segments?** | 4 seconds matches standard HLS GOP (Group of Pictures) recommendations. Shorter chunks (e.g. 1s) cause high HTTP request overhead; longer chunks (e.g. 10s) slow down adaptive quality switching and increase initial playback latency. |
| **5** | **What happens if 10 users watch 10 DIFFERENT titles?** | Origin Shield cannot cache distinct titles. Therefore, admission control estimates home bandwidth saturation and dynamically sheds 1080p down to 480p. Once 10 streams are reached, HTTP 429 `NODE_BUSY` protects home ISP bandwidth. |
| **6** | **What prevents unauthorized chunk scraping?** | Every chunk URL is cryptographically signed using HMAC-SHA256 with a short 15-minute expiration timestamp and bound to the specific viewer session and media ID. Unsigned requests receive HTTP 403 Forbidden. |
| **7** | **Why pre-encoded HLS instead of live transcoding?** | Live transcoding requires continuous heavy GPU/CPU compute during playback and fails if multiple viewers watch at once. Pre-encoding on ingest makes streaming a lightweight static file transfer, easily supporting 10 concurrent streams with < 20 MB idle RAM. |
| **8** | **How does the video player adapt to bad networks?** | The master playlist defines quality tiers (1080p, 720p, 480p) cut at exact 4.0s keyframe boundaries. HLS.js automatically switches down to a lower rendition without rebuffering or audio drift. |
| **9** | **What happens when a node turns off?** | The 15-second Redis heartbeat key expires. If a viewer requests an uncached video, the API immediately returns HTTP 503 `NODE_OFFLINE`. If chunks are already in the Origin Shield cache, they remain playable for 24 hours. |
| **10** | **Why Go for the backend and Next.js for the UI?** | Go compiles to a static binary with < 20 MB idle RAM footprint, essential for running background daemons on low-power home PCs. Next.js provides modern App Router Server-Side Rendering and sleek responsive interfaces. |
