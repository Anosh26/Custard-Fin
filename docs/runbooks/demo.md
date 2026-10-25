# CustardFin Demo & Viva Presentation Runbook

> **Project**: CustardFin — Hybrid Personal Media Streaming & Community Catalog  
> **Team**: Anish Ghadi (25P0320009), Anosh Chodankar (25P0320010)  
> **Target Audience**: MCA Cloud Computing Lab Viva / Evaluation Committee  

---

## 1. Executive Summary & Pitch (First 60 Seconds)

**What is CustardFin?**
CustardFin is an open-source, hybrid media platform that combines **personal home media hosting** (like Plex or Jellyfin) with a **global community discovery catalog** (like YouTube), with near-zero cloud operating costs.

**The Golden Architectural Rule (Invariant 1):**
> *The video binaries stay on the owner's home machine. The cloud holds only the metadata catalog and coordination.*

When a user toggles a video to **Public**, only a **~3 KB JSON manifest** is synced to the cloud PostgreSQL database. Video chunks are uploaded to the cloud **on-demand** only when remote viewers request them, and stored in a temporary 24-hour LRU Origin Shield cache.

---

## 2. The 10 Invariants (Viva Defense Card)

| # | Invariant | How It Is Implemented in Code |
|---|---|---|
| **1** | No video binaries permanently in cloud | Edge stores video; S3 is a temporary 24h LRU Origin Shield (`ports.ChunkCache`) |
| **2** | Outbound-only home connections | Outbound FRP reverse tunnel client (`ports.TunnelManager`); 0 open router ports |
| **3** | 100% Free & Open Source | PostgreSQL, Redis, MinIO, Go stdlib, Next.js, Caddy, FFmpeg |
| **4** | Go Backend + Next.js UI | `apps/edge` and `apps/api` in Go; `apps/web` in Next.js (Node runtime) |
| **5** | Redis for heartbeats only | TTL keys `heartbeat:{node_id}` (15s); Goroutines + Postgres for jobs |
| **6** | Pre-encoded HLS (4s chunks) | Subprocess FFmpeg with NVENC GPU probe and automatic CPU fallback |
| **7** | Origin Shield caching | Hot chunks uploaded once over tunnel; cache serves all viewers |
| **8** | Private data stays home | Watch progress and local history kept in Edge pure-Go SQLite |
| **9** | Signed short-lived URLs | HMAC-SHA256 tokens with 15-minute expiry bound to node and media |
| **10** | 10-stream remote cap | Admission limiter returns HTTP 429 (`Retry-After: 30`) when node is busy |

---

## 3. Step-by-Step Live Demo Script

### Phase 1: Boot Cloud Infrastructure & Health Check (1 min)

1. Open Terminal 1: Start the cloud control plane:
   ```bash
   make demo-up
   ```
   *Explanation to Professor*: "This spins up PostgreSQL, Redis for node heartbeats, MinIO for the Origin Shield chunk cache, the FRP tunnel relay server, and Caddy reverse proxy."

2. Verify all services are healthy:
   ```bash
   make demo-check
   ```
   *Expected Output*: `HTTP 200 OK` from `http://localhost/healthz`.

---

### Phase 2: Transcoding & GPU Acceleration Benchmark (2 mins)

1. Open Terminal 2: Seed test media and run the GPU benchmark:
   ```bash
   make seed
   make bench
   ```
2. **What to point out on screen**:
   - The benchmark automatically probes the local GPU (`h264_nvenc` on NVIDIA hardware) vs CPU (`libx264`).
   - Highlight the **Transcoding Speed Multiplier** (e.g. **10.7x realtime** on RTX/GTX GPU).
   - Point out **Keyframe Alignment**: Every chunk is cut exactly at 4.0-second boundaries (`-g 120 -keyint_min 120` at 30 fps) to enable seamless adaptive bitrate switching without player rebuffering.
   - Point out **Graceful Fallback**: If GPU memory or session limits are reached, the edge automatically falls back to CPU encoding without failing the job.

---

### Phase 3: Start Edge Daemon & Sync Catalog (2 mins)

1. Start the home edge daemon:
   ```bash
   go run ./apps/edge/cmd/edged
   ```
2. **What to point out in edge logs**:
   - `sqlite database initialized (WAL mode enabled)`: Shows watch history stays private at home (Invariant 8).
   - `scanned directory ./media`: Found new videos and generated 4-second HLS renditions (1080p, 720p, 480p).
   - `sync worker published manifest`: Highlight that only a **~3 KB JSON manifest** was sent to the cloud, **NOT** the video file.
   - `heartbeat worker started`: Outbound heartbeat sent every 15s to keep the node active in Redis.

---

### Phase 4: Web Browser Exploration (2 mins)

1. Open your browser to `http://localhost` (or `http://<GCP_EXTERNAL_IP>` if running on cloud):
   - **Community Catalog**: Shows all public videos from all active edge nodes.
   - **Click on a Video**: Starts streaming via HLS.js.
   - Open Browser Developer Tools (`F12 > Network tab`):
     - Filter by `.m3u8` and `.ts`.
     - Point out the signed URL with HMAC token parameters (`?token=...&expires=...`).
     - Point out chunk size and 4-second duration.

---

### Phase 5: Concurrency, Origin Shield & 10-Stream Cap (3 mins)

1. Open Terminal 3: Run the automated 10-viewer load simulation:
   ```bash
   make load-test
   ```
2. **What this verifies live for the evaluation committee**:
   - **Scenario 1 (Origin Shield)**:
     - 10 viewers request the same video simultaneously.
     - Viewer 1 triggers an on-demand pull from home edge through the reverse tunnel (`X-Cache: MISS`).
     - Viewers 2 through 10 are served directly from the cloud S3 cache (`X-Cache: HIT`).
     - Home upload bandwidth saved: **90%**!
   - **Scenario 2 (Adaptive Rendition Throttling)**:
     - Shows dynamic master playlist rewriting. Under heavy load (4+ streams), 1080p is pruned, leaving 720p/480p to prevent buffer stalls on home internet.
   - **Scenario 3 (Per-Node Hard Cap - Invariant 10)**:
     - An 11th distinct stream receives `HTTP 429 Too Many Requests (NODE_BUSY)` with a `Retry-After: 30` header.
   - **Scenario 4 (Node Offline Resilience)**:
     - Cached chunks remain playable even if the home PC loses power. Uncached media cleanly reports `HTTP 503 NODE_OFFLINE`.

---

## 4. Anticipated Viva Questions & Winning Answers

### Q1: "Why don't you just upload all the videos directly to YouTube or S3?"
> **Answer**:  
> "A typical personal media library is hundreds of gigabytes or terabytes. Uploading all of that to cloud object storage costs substantial monthly storage fees ($0.023/GB/month) and egress bandwidth fees. With CustardFin, 100% of the media library remains on the user's existing hard drive at zero cloud storage cost. The cloud only holds a temporary, 24-hour auto-purging cache for hot chunks currently being watched."

### Q2: "How does remote streaming work if the home node is behind CGNAT or a home Wi-Fi router without public IP?"
> **Answer**:  
> "Home edge nodes establish an outbound-only TLS reverse tunnel to our cloud FRP relay server (`frps`) on port 7000. Inbound port forwarding is never required on the user's home router. When a remote viewer requests a chunk that is not in the cloud cache, the cloud API pulls it through this established reverse tunnel, caches it in S3, and serves the signed URL."

### Q3: "What prevents a malicious user from downloading all chunks directly from the cache without permission?"
> **Answer**:  
> "Every chunk URL is cryptographically signed using HMAC-SHA256 with a short 15-minute expiration timestamp and bound to the specific viewer session and media ID. Direct unauthorized access is rejected with HTTP 403 Forbidden."

### Q4: "Why did you use pre-encoded HLS rather than live on-the-fly transcoding?"
> **Answer**:  
> "Live transcoding requires heavy continuous GPU/CPU compute during playback and fails if multiple viewers watch at once. By pre-encoding into standardized 4-second HLS segments upon ingest using NVENC hardware acceleration, serving video becomes a simple, lightweight static file transfer. An edge machine can effortlessly serve 10 concurrent streams with virtually zero CPU overhead during playback."

### Q5: "What role does Redis play, and why not use it for background queues?"
> **Answer**:  
> "By Invariant 5, Redis is strictly used as an ephemeral key-value store with 15-second TTLs to track node online heartbeats. All background job processing (media scanning, transcode pipelines, manifest synchronization) is handled directly by native Go goroutines with transactional state recorded in PostgreSQL or SQLite. This keeps our architecture clean, lean, and avoids unnecessary queue dependencies."
