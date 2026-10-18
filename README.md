# CustardFin

> Free, open-source hybrid personal media streaming and community catalog sharing platform. MCA Cloud Computing Lab project by Anish Ghadi (25P0320009) and Anosh Chodankar (25P0320010).

---

## Architecture Overview

Video files remain on the owner's local machine; the cloud holds only catalog discovery and coordination.

```
HOME EDGE (Go daemon, static binary, <20 MB idle RAM)        CLOUD CONTROL
 scanner → transcoder (FFmpeg, CPU/NVENC) → HLS store         Next.js (UI, SSR, hls.js)
 SQLite (paths, hashes, watch progress)                       Go API + goroutine workers
 sync worker ── 3 KB manifest ───────────────────────────────▶ PostgreSQL (users, nodes, catalog)
 heartbeat ── every 15 s ────────────────────────────────────▶ Redis (TTL key = online)
 FRP client ── outbound TLS tunnel ──────────────────────────▶ FRP server (relay VM)
 local streamer ──▶ LAN viewers (direct)                      S3-compatible cache (Origin Shield) [+ optional CDN]
```

## Quick Start (Dev)

```bash
# 1. Start cloud dependencies
docker compose -f deploy/compose.dev.yml up -d

# 2. Run Cloud API
cd apps/api && go run ./cmd/apid

# 3. Run Home Edge Daemon
cd apps/edge && go run ./cmd/edged

# 4. Run Web UI
cd apps/web && npm install && npm run dev
```

## Documentation
- Master Agent Instructions: [AGENTS.md](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/AGENTS.md)
- Invariants: [INVARIANTS.md](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/context/INVARIANTS.md)
- Project State: [STATE.md](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/context/STATE.md)
- OpenAPI Specification: [openapi.yaml](file:///c:/Users/Anosh/Documents/Notes/Projects/Cloud_Computing/sem_end/Custard%20Fin/docs/api/openapi.yaml)
