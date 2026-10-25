# CustardFin — Step-by-Step Docker Demo Presentation Guide

> **Project**: CustardFin (MCA Cloud Computing Lab)  
> **Team**: Anish Ghadi (25P0320009), Anosh Chodankar (25P0320010)  
> **Host Environment**: Windows 11 with NVIDIA GeForce RTX 4060 GPU  
> **Presentation Duration**: ~8–10 Minutes

---

## 1. Architecture Setup for the Presentation

In CustardFin's hybrid architecture:
- **Cloud Control Plane (Docker)**: PostgreSQL (metadata catalog), Redis (ephemeral heartbeats), MinIO (Origin Shield cache), FRP Server (NAT traversal relay), and optionally Cloud API + Web UI.
- **Home Edge Node (Host PC)**: Runs natively on the host machine with direct access to local media files and the NVIDIA RTX 4060 GPU (`h264_nvenc`). It initiates an **outbound-only tunnel** to Docker.

```
┌────────────────────────────────────────────────────────┐
│               DOCKER CONTAINERS (CLOUD)                │
│                                                        │
│  [PostgreSQL:5432]     [Redis:6379]                    │
│  [MinIO Cache:9000]    [FRP Server:7000]               │
│  [Cloud API:8080]      [Next.js Web:3000]              │
└───────────────────────────▲────────────────────────────┘
                            │ Outbound Tunnel (:7000)
                            │ & Manifest Sync (:8080)
┌───────────────────────────┴────────────────────────────┐
│              HOST MACHINE (HOME EDGE NODE)             │
│                                                        │
│  - Media Directory (./media)                           │
│  - SQLite Edge DB (WAL mode, private watch history)    │
│  - FFmpeg Transcoder (RTX 4060 NVENC GPU acceleration) │
│  - Edge Daemon (edged)                                 │
└────────────────────────────────────────────────────────┘
```

---

## 2. Pre-Demo Checklist (2 Minutes Before Presentation)

Before the evaluators arrive, prepare your environment:

1. **Docker Desktop**: Open Docker Desktop and ensure the Docker Engine is running (green indicator).
2. **Open Three Terminal Tabs** (in VS Code or Windows Terminal):
   - **Tab 1 — Cloud & Simulation** (`CWD: Custard Fin root`)
   - **Tab 2 — Home Edge Daemon** (`CWD: Custard Fin root`)
   - **Tab 3 — Web Client / Logs** (`CWD: apps/web` or root)
3. **Browser**: Have Google Chrome or Edge open to `http://localhost:3000` (or `http://localhost` if using Caddy).

---

## 3. Step-by-Step Presentation Script

---

### Step 1: Pitch & The Core Architectural Invariant (1 Minute)

#### 🗣️ Speaker Notes (What to Say)
> *"Good morning, Professors. Today we are presenting **CustardFin**, a hybrid media streaming platform that bridges personal media servers like Plex/Jellyfin with community discovery platforms like YouTube.*
>
> *Our core design principle is **Invariant 1**: **Video binaries stay on the owner's home machine; the cloud holds only metadata and coordination.** When a user marks a video as Public, only a tiny 3 KB JSON manifest is synchronized to the cloud. Hot video chunks are pulled to an Origin Shield cache only when viewers stream them, cutting cloud storage and bandwidth costs to near-zero."*

---

### Step 2: Spin Up Cloud Infrastructure in Docker (1.5 Minutes)

#### 💻 Action: Tab 1 (PowerShell)
Run:
```powershell
make dev-up
```
*(Or `docker compose -f deploy/compose.dev.yml up -d`)*

Verify all containers are healthy:
```powershell
docker compose -f deploy/compose.dev.yml ps
```

#### 🖥️ What Appears on Screen
```text
NAME                  IMAGE                       STATUS                   PORTS
custardfin-cache      pgsty/minio:latest          Up (healthy)             0.0.0.0:9000-9001
custardfin-frps       snowdreamtech/frps:latest   Up                       0.0.0.0:7000, 7500, 8000-8020
custardfin-postgres   postgres:16-alpine          Up (healthy)             0.0.0.0:5432
custardfin-redis      redis:7-alpine              Up (healthy)             0.0.0.0:6379
```

#### 🗣️ Speaker Notes
> *"With one command, we spin up our cloud dependencies in Docker:*
> 1. ***PostgreSQL 16***: *Relational catalog for registered nodes, users, and manifests.*
> 2. ***Redis 7***: *Used strictly for ephemeral 15-second node heartbeats (Invariant 5).*
> 3. ***MinIO (S3)***: *Acts as the temporary 24-hour Origin Shield pull-through cache.*
> 4. ***FRP Relay Server***: *Allows home nodes behind CGNAT to be reached without opening incoming router ports (Invariant 2)."*

---

### Step 3: Hardware NVENC GPU Transcoding Benchmark (2 Minutes)

#### 💻 Action: Tab 1 (PowerShell)
Seed test media fixtures and run the GPU benchmark:
```powershell
powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1
powershell -ExecutionPolicy Bypass -File scripts/bench-gpu.ps1
```

#### 🖥️ What Appears on Screen
```text
[BENCH] Probing hardware acceleration capabilities...
[BENCH] Detected NVENC: NVIDIA GeForce RTX 4060 Laptop GPU
[BENCH] Running NVENC benchmark (h264_nvenc, 1080p, 4s GOP)...
[BENCH] NVENC Transcoding Time: 2.73s (16.48x Realtime)
[BENCH] Running CPU baseline (libx264, 1080p, 4s GOP)...
[BENCH] CPU Transcoding Time:   3.63s (12.39x Realtime)
[BENCH] Speedup Factor: 1.33x faster (4.0s keyframes verified)
```

#### 🗣️ Speaker Notes
> *"Here we demonstrate **Invariant 6**: Pre-encoded 4-second HLS chunks with hardware acceleration.*
>
> *Our edge adapter automatically probes the system for hardware encoders (`h264_nvenc` on NVIDIA RTX, `h264_qsv` on Intel, `h264_vaapi` on AMD). On this RTX 4060, NVENC encodes at **16.5x realtime**.*
>
> *Crucially, notice **keyframe alignment**: every chunk is cut at exactly 4.0-second boundaries. This guarantees that when a viewer's network fluctuates, the video player can switch between 1080p, 720p, and 480p without any rebuffering or audio drift.*
>
> *If GPU memory is exhausted, the system automatically and silently falls back to CPU encoding without crashing."*

---

### Step 4: Launch Cloud API & Home Edge Daemon (2 Minutes)

#### 💻 Action: Tab 1 — Start Cloud API
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

#### 💻 Action: Tab 2 — Start Home Edge Daemon
```powershell
go run ./apps/edge/cmd/edged
```

#### 🖥️ What Appears in Edge Logs (Point this out!)
```text
INFO sqlite database initialized (WAL mode enabled) db=./data/edge.db
INFO media scanner started dir=./media
INFO scanned media file path=media/sample_1080p.mp4 hash=3b8a1f...
INFO transcoding complete renditions=3 duration=45s
INFO sync worker published manifest media_id=... bytes=1420 is_public=true
INFO heartbeat sent to cloud active_streams=0 gpu=NVENC
```

#### 🗣️ Speaker Notes
> *"Watch the edge daemon logs:*
> 1. ***SQLite with WAL mode***: *Watch progress and file locations stay strictly on the local machine (Invariant 8).*
> 2. ***Manifest Sync***: *Notice the size of the sync payload: **1,420 bytes (~1.4 KB)**! We did NOT upload the 100 MB MP4 file. Only the playlist structure and content hashes were registered in the cloud.*
> 3. ***Heartbeat Worker***: *Sends an outbound heartbeat every 15 seconds to Redis to announce that the node is online."*

---

### Step 5: Web UI Playback & HMAC Signed URLs (1.5 Minutes)

#### 💻 Action: Tab 3 — Start Web Client
```powershell
cd apps/web
npm run dev
```

#### 🌐 Browser Action:
1. Open `http://localhost:3000`.
2. The homepage displays the **Community Catalog** showing the public video just synced from the edge node.
3. Open **Chrome DevTools** (`F12`), switch to the **Network** tab, and type `m3u8` or `ts` in the filter box.
4. Click on the video card to start watching.

#### 🖥️ What to Point Out in DevTools
1. Show the signed playback URL:
   ```text
   http://localhost:8080/v1/stream/.../segment_000.ts?token=a8f9b2...&expires=1740000000
   ```
2. Click the response headers of any `.ts` segment:
   - Notice `X-Cache: MISS` on first play (pulled on-demand via tunnel from edge).
   - Refresh the page and notice `X-Cache: HIT` (served from MinIO cache).

#### 🗣️ Speaker Notes
> *"Here in the web player, HLS.js smoothly streams the video. Looking at Developer Tools:*
> - *Every chunk URL has an **HMAC-SHA256 signature** with a 15-minute expiration timestamp (Invariant 9). Direct hotlinking without an authenticated session is blocked.*
> - *The response header confirms **Origin Shield**: the first viewer triggered an on-demand pull through the tunnel; all subsequent chunk requests hit the cloud cache directly."*

---

### Step 6: 10-Viewer Concurrency & Admission Control Simulation (2 Minutes)

#### 💻 Action: Open a new terminal and run:
```powershell
powershell -ExecutionPolicy Bypass -File scripts/load-10-viewers.ps1
```

#### 🖥️ What Appears on Screen
```text
================================================================================
CUSTARDFIN 10-VIEWER CONCURRENCY & ADMISSION SIMULATION
================================================================================

[SCENARIO 1] Origin Shield Cache Offload (10 Viewers, Same Title)
Viewer 01: [MISS] Pulled through reverse tunnel from home edge (bytes=245812)
Viewer 02: [HIT ] Served from Cloud Origin Shield Cache (bytes=245812)
...
Viewer 10: [HIT ] Served from Cloud Origin Shield Cache (bytes=245812)
=> Bandwidth Saved on Home Upload: 90.0% (Hot chunk served from Cloud S3)

[SCENARIO 2] Adaptive Bitrate & Rendition Pruning Under Node Load
Active Streams: 1  => Available Renditions: 1080p, 720p, 480p (Load: Normal)
Active Streams: 5  => Available Renditions: 720p, 480p (Load: Moderate - 1080p pruned)
Active Streams: 8  => Available Renditions: 480p (Load: Saturated - 1080p & 720p pruned)

[SCENARIO 3] Per-Node Concurrency Hard Cap (Invariant 10)
Stream #10: HTTP 200 OK (Within allowable limit)
Stream #11: HTTP 429 Too Many Requests | Retry-After: 30 | Error: NODE_BUSY

[SCENARIO 4] Node Offline Resilience (Invariant 1 & 7)
Cached Title: HTTP 200 OK (Served from 24h Origin Shield Cache while node offline)
Uncached Title: HTTP 503 Service Unavailable (NODE_OFFLINE)
```

#### 🗣️ Speaker Notes
> *"This load test validates our four key cloud computing mechanisms:*
> 1. ***Origin Shield (Invariant 7)***: *10 viewers stream the same viral video. Viewer 1 uploads the chunk once; Viewers 2 through 10 are served from cache. Home upload bandwidth is reduced by **90%**.*
> 2. ***Adaptive Master Playlist Rewriting***: *As home node upload gets busy, the cloud API dynamically prunes the 1080p stream from the master playlist, guiding players to 720p or 480p before buffering occurs.*
> 3. ***10-Stream Concurrency Limit (Invariant 10)***: *The 11th distinct stream receives HTTP 429 `NODE_BUSY` with `Retry-After: 30`, preventing domestic ISP connections from collapsing.*
> 4. ***Offline Resilience***: *If the home PC turns off, cached media remains watchable globally."*

---

### Step 7: Clean Teardown

#### 💻 Action:
```powershell
make dev-down
```
*(Or `docker compose -f deploy/compose.dev.yml down`)*

#### 🗣️ Speaker Notes
> *"Thank you! We are now open for your questions."*

---

## 4. Quick Reference Card for Tough Viva Questions

| Question | Winning Technical Answer |
|---|---|
| **Why not upload directly to S3 or YouTube?** | Storing terabytes of personal video permanently on cloud object storage incurs unsustainable monthly storage fees. CustardFin costs ~$0 for storage because files stay on local hard drives. Cloud cache only retains active chunks for 24 hours. |
| **How do you bypass CGNAT without port forwarding?** | The edge node establishes an outbound TLS reverse tunnel (`frpc`) to the cloud relay server (`frps`). Because the connection originates from inside the home, NAT firewalls permit it without inbound port forwarding. |
| **Why Redis only for heartbeats?** | Using Redis for task queues introduces state duplication and race conditions. Background tasks run as standard Go goroutines with transactional state stored in PostgreSQL. Redis is used solely for atomic 15s TTL keys (`heartbeat:{node_id}`). |
| **Why 4-second HLS segments?** | 4 seconds matches standard HLS GOP (Group of Pictures) recommendations. Shorter chunks (e.g. 1s) cause high HTTP request overhead; longer chunks (e.g. 10s) slow down adaptive quality switching and increase initial playback latency. |
| **What happens if 10 users watch 10 DIFFERENT titles?** | Origin Shield cannot cache distinct titles. Therefore, admission control estimates home bandwidth saturation and sheds 1080p down to 480p. Once 10 streams are reached, HTTP 429 `NODE_BUSY` protects home ISP bandwidth. |
