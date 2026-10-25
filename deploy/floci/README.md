# Running CustardFin with Floci Local Cloud Emulator

> **Component**: Floci Integration Guide  
> **Status**: Verified & Ready  
> **Target Cloud API**: AWS S3 Emulation on Port `4566`  

---

## 1. What is Floci?

[Floci](https://github.com/floci-io/floci) is a lightweight, 100% free, and open-source (MIT licensed) local cloud emulator family. Unlike proprietary or freemium cloud emulation platforms, Floci provides fully open emulations of AWS, Azure, and GCP APIs.

In CustardFin:
- **Port 4566** emulates the **AWS S3 Object Storage API**.
- Serves as the **Origin Shield pull-through chunk cache** (Invariants 1 & 7).
- CustardFin's Go S3 adapter (`apps/api/internal/adapters/s3/chunk_cache.go`) communicates with Floci via standard AWS SDK for Go v2 path-style requests.

---

## 2. Ports and Service Map

| Service | Container / Process | Port | Purpose in CustardFin |
|---|---|---|---|
| **Floci (AWS S3)** | `custardfin-floci` | `4566` | Origin Shield HLS temporary cache (auto-purging) |
| **PostgreSQL** | `custardfin-postgres` | `5432` | User accounts, registered nodes, public catalog metadata |
| **Redis** | `custardfin-redis` | `6379` | Ephemeral node heartbeats (`heartbeat:{node_id}` 15s TTL) |
| **FRP Server** | `custardfin-frps` | `7000`, `8000-8020` | Outbound reverse tunnel relay for CGNAT traversal |
| **Cloud API** | `apid` (Go process) | `8080` | Stream router, admission control, public catalog REST API |
| **Web Frontend** | `web` (Next.js) | `3000` | Community discovery catalog & HLS.js streaming player |
| **Edge Node** | `edged` (Go process) | `8090` | Home node media scanner, NVENC/CPU transcoder, LAN streamer |

---

## 3. Step-by-Step: Running the Whole System with Floci

### Step 1: Start Docker Desktop
Ensure Docker Desktop is running on your machine.

### Step 2: Spin Up the Floci Cloud Infrastructure
From the repository root:
```bash
docker compose -f deploy/floci/compose.floci.yml up -d
```
*This boots Floci on port 4566, PostgreSQL on 5432, Redis on 6379, FRP on 7000, and automatically provisions the `custardfin-origin-cache` S3 bucket.*

### Step 3: Start the Cloud API (`apps/api`)
Open **Terminal 1** and launch the API server configured for Floci:

**On Windows (PowerShell):**
```powershell
$env:API_ADDR = ":8080"
$env:DATABASE_URL = "postgres://custard:custard_secret@localhost:5432/custardfin?sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
$env:S3_ENDPOINT = "http://localhost:4566"
$env:S3_BUCKET = "custardfin-origin-cache"
$env:S3_ACCESS_KEY = "test"
$env:S3_SECRET_KEY = "test"
$env:S3_REGION = "us-east-1"
$env:TUNNEL_BASE_URL = "http://localhost:8090"

go run ./apps/api/cmd/apid
```

**On Linux/macOS (Bash):**
```bash
export API_ADDR=":8080"
export DATABASE_URL="postgres://custard:custard_secret@localhost:5432/custardfin?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export S3_ENDPOINT="http://localhost:4566"
export S3_BUCKET="custardfin-origin-cache"
export S3_ACCESS_KEY="test"
export S3_SECRET_KEY="test"
export S3_REGION="us-east-1"
export TUNNEL_BASE_URL="http://localhost:8090"

go run ./apps/api/cmd/apid
```

Verify health:
```bash
curl http://localhost:8080/healthz
# Response: {"status":"ok"}
```

---

### Step 4: Start the Next.js Web UI (`apps/web`)
Open **Terminal 2**:
```bash
cd apps/web
npm run dev
```
The community streaming catalog will be accessible at `http://localhost:3000`.

---

### Step 5: Start the Home Edge Daemon (`apps/edge`)
Open **Terminal 3**:
```powershell
# Seed demo video fixtures if not already present
powershell -ExecutionPolicy Bypass -File scripts/seed-media.ps1

# Run the Edge Daemon
go run ./apps/edge/cmd/edged
```
The edge daemon will:
1. Scan `./media` and register video files into the local pure-Go SQLite database.
2. Probe GPU hardware acceleration (NVENC on RTX 4060).
3. Send periodic heartbeats to the cloud API (`POST /v1/nodes/{id}/heartbeat`).

---

### Step 6: Verify Origin Shield Caching in Floci
Open **Terminal 4**:
Run the automated load simulation script:
```powershell
powershell -ExecutionPolicy Bypass -File scripts/load-10-viewers.ps1
```
You will observe:
- **First request**: Cache MISS $\to$ Pulled over reverse tunnel $\to$ Stored in Floci S3 bucket.
- **Subsequent requests**: Cache HIT $\to$ Served from Floci on port 4566 with `< 5ms` latency.

---

## 4. Teardown

To shut down the Floci infrastructure:
```bash
docker compose -f deploy/floci/compose.floci.yml down
```
To purge persistent emulator volumes as well:
```bash
docker compose -f deploy/floci/compose.floci.yml down -v
```
