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

## How to Run (4 Operational Modes)

CustardFin supports 4 flexible execution topologies depending on your environment. See [**`RUNNING.md`**](RUNNING.md) and [**`docs/runbooks/RUNNING_GUIDE.md`**](docs/runbooks/RUNNING_GUIDE.md) for full instructions:

1. **[Mode 1: Pure Local (Zero Docker)](docs/runbooks/RUNNING_GUIDE.md#3-mode-1-pure-local-zero-docker--standalone-mode)**: In-memory fallback stores; zero external services needed.
2. **[Mode 2: Local Containerized](docs/runbooks/RUNNING_GUIDE.md#4-mode-2-local-containerized-dev-profile-with-minio)**: Full local Docker stack (`make dev-up`).
3. **[Mode 3: Floci Local Cloud Emulator (Multi-PC Lab)](docs/runbooks/RUNNING_GUIDE.md#5-mode-3-floci-local-cloud-emulator-multi-pc-lab-setup)**: Multi-computer lab demo with Floci S3 emulator (:4566), RTX 4060 edge node, and remote viewer clients.
4. **[Mode 4: Google Cloud Platform (GCP Live Cloud)](docs/runbooks/RUNNING_GUIDE.md#6-mode-4-google-cloud-platform-gcp-live-cloud-way)**: Real cloud deployment on GCP VM with Caddy HTTPS and outbound reverse tunnel punching through CGNAT.

## Documentation & Runbooks
- **Deployment & Running Guide**: [docs/runbooks/RUNNING_GUIDE.md](docs/runbooks/RUNNING_GUIDE.md)
- **Viva Demo Script & Defense Guide**: [docs/runbooks/demo.md](docs/runbooks/demo.md)
- **GCP Cloud Deployment Guide**: [docs/runbooks/deploy.md](docs/runbooks/deploy.md)
- **Troubleshooting & Emergency Recovery**: [docs/runbooks/troubleshooting.md](docs/runbooks/troubleshooting.md)
- **Master Agent Prompt**: [AGENTS.md](AGENTS.md)
- **Ten Architectural Invariants**: [docs/context/INVARIANTS.md](docs/context/INVARIANTS.md)
- **Project State**: [docs/context/STATE.md](docs/context/STATE.md)
- **OpenAPI Specification**: [docs/api/openapi.yaml](docs/api/openapi.yaml)
