# CustardFin Universal Running & Deployment Guide

> **Project**: CustardFin — Hybrid Personal Media Streaming & Community Catalog  
> **Authors**: Anish Ghadi (25P0320009), Anosh Chodankar (25P0320010)  
> **Target Audience**: Developers, Lab Evaluators, and Cloud Deployment Engineers  

---

## 1. Execution Modes at a Glance

CustardFin is architected with strict decoupling and hexagonal ports/adapters. It can run in four distinct operational topologies:

| Mode | External Dependencies | S3 Object Cache | Best For | Typical Topology |
|---|---|---|---|---|
| **Mode 1: Pure Local** | None (Zero Docker required) | Built-in In-Memory LRU | Quick tests, offline laptops, code inspection | Single machine |
| **Mode 2: Local Containers** | Docker (PostgreSQL, Redis, MinIO, FRP) | MinIO (:9000/:9001) | Full local integration testing | Single machine |
| **Mode 3: Floci (Multi-PC Lab)** | Docker (Floci S3 emulator, PostgreSQL, Redis, FRP) | Floci S3 (:4566) | College lab demo across 2+ computers | Multi-PC LAN / Wi-Fi |
| **Mode 4: Google Cloud (GCP)** | GCP VM (Compute Engine e2-standard-2) | MinIO / S3 on Cloud VM | Final live cloud evaluation & viva | Cloud VM + Home Edge Node |

---

## 2. Prerequisites Across All Modes

- **Go**: Version 1.22+ (Multi-module workspace: `apps/edge`, `apps/api`, `libs/contracts`).
- **Node.js**: Version 20+ LTS with `npm` (for `apps/web` Next.js frontend).
- **FFmpeg & FFprobe**: Version 6.0+ (Must be available in system `PATH`).
  - Supports NVIDIA GPU hardware acceleration (`h264_nvenc`) with automatic CPU fallback (`libx264`).
- **Docker & Docker Compose**: (Required for Modes 2, 3, and 4; **NOT** required for Mode 1).

---

## 3. Mode 1: Pure Local (Zero Docker / Standalone Mode)

In this mode, all subsystems run natively on a single PC. `apid` activates its built-in in-memory fallback repositories (in-memory SQLite/RAM catalog, in-memory Redis heartbeat tracker, in-memory Origin Shield chunk cache). No database or container setup is needed!

```
                    ┌─────────────────────────┐
                    │    Web UI (:3000)       │
                    │  Next.js App Router     │
                    └────────────┬────────────┘
                                 │ HTTP
                                 ▼
┌────────────────────────┐  HTTP   ┌─────────────────────────┐
│   Home Edge (:8090)    │◄───────┤   Cloud API (:8080)     │
│  edged (SQLite WAL,    │         │  apid (In-Memory Stores)│
│  NVENC, ./media scan)  │         │  Origin Shield In-Memory│
└────────────────────────┘         └─────────────────────────┘
```

### Step 1: Synthesize Demo Media Fixtures
```powershell
# Windows (PowerShell)
powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1

# Linux / macOS (Bash)
bash scripts/seed-media.sh
```

### Step 2: Start Cloud Control API (`apps/api`)
Open **Terminal 1**:

**Windows PowerShell:**
```powershell
$env:API_ADDR = ":8080"
$env:DATABASE_URL = "memory"
$env:REDIS_ADDR = "memory"
$env:S3_BUCKET = ""
$env:TUNNEL_BASE_URL = "http://localhost:8090"

go run ./apps/api/cmd/apid
```

**Linux / macOS Bash:**
```bash
export API_ADDR=":8080"
export DATABASE_URL="memory"
export REDIS_ADDR="memory"
export S3_BUCKET=""
export TUNNEL_BASE_URL="http://localhost:8090"

go run ./apps/api/cmd/apid
```

### Step 3: Start Home Edge Daemon (`apps/edge`)
Open **Terminal 2**:

**Windows PowerShell:**
```powershell
$env:EDGE_ADDR = ":8090"
$env:EDGE_MEDIA_DIR = "./media"
$env:EDGE_TRANSCODE_DIR = "./transcoded"
$env:EDGE_DB_PATH = "./data/edge.db"
$env:EDGE_CLOUD_URL = "http://localhost:8080"

go run ./apps/edge/cmd/edged
```

### Step 4: Start Web Interface (`apps/web`)
Open **Terminal 3**:
```bash
cd apps/web
npm install
npm run dev
```

### Step 5: Test Playback
- Open `http://localhost:3000` in your web browser.
- Run the GPU benchmark: `powershell -ExecutionPolicy Bypass -File scripts/bench-gpu.ps1` (or `bash scripts/bench-gpu.sh`).
- Run the 10-viewer load simulation: `powershell -ExecutionPolicy Bypass -File scripts/load-10-viewers.ps1`.

---

## 4. Mode 2: Local Containerized (Dev Profile with MinIO)

Spins up standard open-source infrastructure (PostgreSQL, Redis, MinIO S3, FRP server) via Docker Compose.

### Step 1: Start Supporting Containers
```bash
make dev-up
# Equivalent to: docker compose -f deploy/compose.dev.yml up -d
```

### Step 2: Seed Media & Run Services
```powershell
# 1. Seed demo media
powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1

# 2. Terminal 1: Cloud API
$env:API_ADDR = ":8080"
$env:DATABASE_URL = "postgres://custard:custard_secret@localhost:5432/custardfin?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
$env:S3_ENDPOINT = "http://localhost:9000"
$env:S3_BUCKET = "custardfin-origin-cache"
$env:S3_ACCESS_KEY = "minioadmin"
$env:S3_SECRET_KEY = "minioadmin"
go run ./apps/api/cmd/apid

# 3. Terminal 2: Edge Daemon
go run ./apps/edge/cmd/edged

# 4. Terminal 3: Web UI
cd apps/web && npm run dev
```

### Teardown
```bash
make dev-down
```

---

## 5. Mode 3: Floci Local Cloud Emulator (Multi-PC Lab Setup)

This mode demonstrates the **true decentralized power of CustardFin** in a computer lab setting using multiple machines connected to the same Wi-Fi or lab LAN switch.

### Architecture & Machine Allocation

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
   │  Native Services:                                           │
   │   - Cloud API (apid) on Port 8080                           │
   │   - Web Frontend (Next.js) on Port 3000                     │
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

### Setup on Machine A (Cloud Server / Lab Host)

1. Find Machine A's local IP address:
   - Windows: `ipconfig` (e.g. `192.168.1.100`)
   - Linux: `ip a` or `hostname -I`
2. Start the Floci cloud infrastructure:
   ```bash
   make floci-up
   # Or: docker compose -f deploy/floci/compose.floci.yml up -d
   ```
3. Start the Cloud API (`apid`) binding to all network interfaces:
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
4. Start the Web UI (`apps/web`):
   ```bash
   cd apps/web
   npm run dev -- -H 0.0.0.0 -p 3000
   ```

---

### Setup on Machine B (Home Edge Node - RTX 4060 PC)

Machine B acts as the user's home computer. It requires **ZERO inbound open ports** (Invariant 2).

1. Seed test videos:
   ```powershell
   powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1
   ```
2. Set configuration pointing outbound to Machine A (`192.168.1.100`):
   ```powershell
   $CLOUD_IP = "192.168.1.100"  # Replace with Machine A's IP

   $env:EDGE_ADDR = ":8090"
   $env:EDGE_MEDIA_DIR = "./media"
   $env:EDGE_TRANSCODE_DIR = "./transcoded"
   $env:EDGE_DB_PATH = "./data/edge.db"
   $env:EDGE_CLOUD_URL = "http://${CLOUD_IP}:8080"
   $env:EDGE_NODE_ID = "550e8400-e29b-41d4-a716-446655440001"
   $env:EDGE_NODE_TOKEN = "custardfin_demo_node_token"

   # Enable outbound FRP reverse tunnel
   $env:EDGE_FRP_ENABLED = "true"
   $env:EDGE_FRP_SERVER_ADDR = $CLOUD_IP
   $env:EDGE_FRP_SERVER_PORT = "7000"
   $env:EDGE_FRP_AUTH_TOKEN = "custardfin_demo_frp_token"

   go run ./apps/edge/cmd/edged
   ```
3. Watch Machine B's terminal:
   - It probes NVENC GPU acceleration.
   - It hashes and catalogs `./media/` videos.
   - It establishes the outbound reverse tunnel to Machine A.
   - It publishes the **~3 KB JSON metadata manifest** to Machine A (video files remain on Machine B!).

---

### Setup on Machine C (Viewer Client)

1. Open any browser on any phone, laptop, or tablet connected to the same Wi-Fi.
2. Navigate to: `http://192.168.1.100:3000`
3. Click any video:
   - Notice the video plays instantly!
   - Open Developer Tools (`F12 > Network`): observe the `.m3u8` and `.ts` chunk requests.
   - Verify the `X-Cache: HIT` response from Floci S3 cache on port 4566.

---

## 6. Mode 4: Google Cloud Platform (GCP Live Cloud Way)

For the final live evaluation, the cloud control tier is deployed to a Google Cloud Platform VM with a static public IP and Caddy HTTPS.

### VM Sizing & Budget ($40 Semester Credit)
- **Machine Type**: `e2-standard-2` (2 vCPU, 8 GB RAM) — costs ~$0.067/hour (~$1.60/day).
- **Disk**: 30 GB Balanced Persistent Disk.
- **Region**: Closest to your college (e.g. `asia-south1` Mumbai or `asia-south2` Delhi).

### Step 1: Create GCP Firewall Rules
In Google Cloud Console (`VPC network > Firewall`) or via `gcloud`:
```bash
gcloud compute firewall-rules create allow-custardfin \
    --direction=INGRESS \
    --priority=1000 \
    --network=default \
    --action=ALLOW \
    --rules=tcp:80,tcp:443,tcp:7000,tcp:8000-8020 \
    --source-ranges=0.0.0.0/0 \
    --target-tags=custardfin-server
```

### Step 2: Deploy Cloud Stack on the GCP VM
SSH into the GCP VM:
```bash
# 1. Clone repository
git clone https://github.com/Anosh26/Custard-Fin.git
cd Custard-Fin

# 2. Configure environment
cp deploy/env.example .env

# 3. Launch demo stack (Caddy, Postgres, Redis, MinIO, FRP, API, Web)
docker compose -f deploy/compose.demo.yml up -d --build

# 4. Verify health
curl http://localhost:8080/healthz
```

### Step 3: Connect Home Edge Node Across the Public Internet
On your home laptop / RTX 4060 PC at home or college:
```powershell
$GCP_IP = "34.X.X.X"  # Replace with GCP VM's Public IP

$env:EDGE_CLOUD_URL = "http://${GCP_IP}/api"
$env:EDGE_FRP_ENABLED = "true"
$env:EDGE_FRP_SERVER_ADDR = $GCP_IP
$env:EDGE_FRP_SERVER_PORT = "7000"
$env:EDGE_FRP_AUTH_TOKEN = "custardfin_demo_frp_token"
$env:EDGE_NODE_ID = "550e8400-e29b-41d4-a716-446655440001"
$env:EDGE_NODE_TOKEN = "custardfin_demo_node_token"

go run ./apps/edge/cmd/edged
```
*Because the home node makes only outbound connections, it seamlessly punches through college Wi-Fi CGNAT and firewalls.*

---

## 7. Operational Validation & Benchmark Commands

Run these automated verification scripts to validate performance during viva evaluations:

| Action | PowerShell (Windows) | Bash (Linux/macOS) | What It Verifies |
|---|---|---|---|
| **Synthesize Media** | `powershell -File scripts/seed-media.ps1` | `bash scripts/seed-media.sh` | Generates 1080p, 720p, 480p fixtures via FFmpeg |
| **GPU Benchmark** | `powershell -File scripts/bench-gpu.ps1` | `bash scripts/bench-gpu.sh` | Validates NVENC 16x speedup & strict 4.0s keyframes |
| **10-Viewer Simulation** | `powershell -File scripts/load-10-viewers.ps1` | `bash scripts/load-10-viewers.sh` | Origin Shield cache hits (90%), ladder pruning, & 10-stream cap |
| **Unit & E2E Tests** | `go test ./...` | `go test -v ./...` | Validates all Go packages and end-to-end integration |
| **Build Web App** | `cd apps/web && npm run build` | `cd apps/web && npm run build` | Compiles Next.js production SSR bundle |

---

## 8. Summary Checklist for Demo Day

- [ ] Video files remain on local edge PC hard drive (**Invariant 1**).
- [ ] Only ~3 KB JSON manifests are synced to the cloud catalog (**FR-5**).
- [ ] Node connections are outbound-only through FRP tunnel (**Invariant 2**).
- [ ] Transcoder probes NVENC and falls back to CPU automatically (**Invariant 6**).
- [ ] Origin Shield serves cached chunks on cache HIT with `< 5ms` latency (**Invariant 7**).
- [ ] Concurrency limiter enforces 10-stream remote cap (`HTTP 429`) (**Invariant 10**).
- [ ] Watch history stays private in edge SQLite database (**Invariant 8**).
