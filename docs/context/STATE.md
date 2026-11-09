# CustardFin Project State

## Current Milestone
- **Milestone 10**: Hardening, demo profile, rehearsal, runbooks (COMPLETE & VERIFIED LIVE).
- **Status**: All 10 Milestones successfully implemented, tested, documented, and verified live on Google Cloud Platform with remote mobile streaming. Ready for final evaluation.

## Done
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.
- **Milestone 2**: Edge Node Subsystems (Standalone Mini-Jellyfin).
- **Milestone 3**: Cloud Control API (`apps/api`).
- **Milestone 4**: Public Toggle -> Manifest Sync Worker -> Next.js Community Catalog UI.
- **Milestone 5**: FRP Tunnel (Remote Playback End-to-End).
- **Milestone 6**: Redis Heartbeat and Offline Logic.
- **Milestone 7**: Origin Shield Cache + Signed URLs + Stream Router.
- **Milestone 8**: Admission Control, Adaptive Bitrate & Prefetch.
- **Milestone 9**: NVENC Path with Probe, Fallback and Benchmark Script.
- **Milestone 10**: Demo Profile, Load Simulation Scripts, GCP Deployment & Viva Runbooks.
  - Built 10-viewer load simulation scripts (`scripts/load-10-viewers.sh` and `scripts/load-10-viewers.ps1`) demonstrating Origin Shield cache hits, adaptive ladder throttling, and 10-stream cap.
  - Created multi-stage Dockerfiles (`apps/api/Dockerfile`, `apps/web/Dockerfile`) and unified demo orchestration (`deploy/compose.demo.yml`, `deploy/Caddyfile`).
  - Authored comprehensive GCP deployment guide (`docs/runbooks/deploy.md`), viva presentation script (`docs/runbooks/demo.md`), troubleshooting guide (`docs/runbooks/troubleshooting.md`), and universal running guide (`RUNNING.md`).
  - Created Floci local cloud emulator deployment profile (`deploy/floci/compose.floci.yml`, `deploy/floci/README.md`) with `make floci-up`/`make floci-down`.
- **Live Cloud Deployment & Mobile Verification (GCP)**:
  - Deployed full cloud stack on GCP VM (`34.132.188.14`) with Caddy, Next.js, Go API, PostgreSQL, Redis, MinIO (Origin Shield), and FRPS.
  - Connected Fedora edge daemon (`edged`) via outbound FRP reverse tunnel (`frpc`) with NVIDIA NVENC acceleration.
  - Verified end-to-end HLS playback, buffering recovery, seeking, video uploading, and 1-click Public/Private catalog sync.
  - Verified remote streaming on mobile (Brave browser on Android) with Origin Shield cache hits and live HUD stats.
- **Interactive Video Management, Dev Console & Presentation Deck**:
  - Edge ingestion API (`/api/media/upload`, `/api/media/scan`, `/api/media/generate-sample`).
  - Next.js UI Modals: `AddVideoModal`, `EdgeLibraryModal`, `DevSettingsModal`, and Player "Stats for Nerds" HUD.
  - Modal Viewport Portaling (`ModalPortal.tsx`) to avoid CSS `backdrop-filter` clipping.
  - Complete 8-slide presentation suite across OpenXML PowerPoint (`docs/presentation/CustardFin_Presentation.pptx`), OpenXML dump (`docs/presentation/presentation.xml`), generator script (`scripts/generate_presentation.py`), and interactive HTML5 web deck (`docs/presentation/index.html`) with hardware-agnostic universal GPU/CPU wording, GCP setup, and per-slide custard mascot animations.

- **User Accounts, Edge Node Claiming & Playback Resilience**:
  - Cloud API (`apps/api`): Added `POST /v1/nodes/{id}/claim` allowing authenticated users (e.g. `@anish`) to claim / bind their edge node, instantly re-attributing catalog uploads and granting private video visibility. Added `owner_username` and `owner_user_id` to `GET /v1/nodes/{id}/status`. Fixed prefetch segment name typo to `1080p/segment_0000.ts`.
  - Edge Node Daemon (`apps/edge`): Modified `edged/main.go` and `server.go` to generate crisp poster thumbnails immediately and ONLY sync media manifests to the Cloud catalog once `master.m3u8` is confirmed on disk, preventing catalog publishing during in-flight encoding.
  - Web UI (`apps/web`): Added `claimNode` and `fetchNodeStatus` in `api-client.ts`; added node ownership banner and 1-click `Link to @username` button in `EdgeLibraryModal.tsx`. Enhanced `Player.tsx` with automatic re-polling during stream warm-up and resilient retry handling.

- **Direct Local LAN Streaming & Standalone Mini-Jellyfin (FR-3, FR-6)**:
  - Embedded Standalone Web UI (`apps/edge`): Added responsive, zero-cloud single-page web application served on `GET /` of `edged` (:8090). Any Wi-Fi device can stream directly with < 5ms delay, browse local library, and save watch history without contacting the cloud tier.
  - Dynamic LAN IP Detection & Reporting: `edged` detects primary outbound LAN IPv4 address (or `EDGE_LAN_ADDR`) and transmits `lan_addr` in 15s heartbeats.
  - Cloud API Intelligent Stream Router: Automatically routes viewers on the same home network (matching public NAT IP or private network) to `http://<node-lan-ip>:8090/stream/...` instead of hardcoded `localhost:8090`, eliminating cloud hops and WAN bandwidth for local viewers.
  - Web UI Integration: Local Edge Library links directly with `?lan=true` and surfaces 1-click shortcut to `:8090` standalone player.

- **Comprehensive Academic Project Report (`report/`)**:
  - Authored a 62-page publication-grade LaTeX project report explicitly titled "Cloud Computing Project Report" formatted with IEEEtran bibliography and standard academic front matter.
  - Strictly specified student identification by Seat Numbers: Anish Ghadi (Seat No.: 25P0320009) and Anosh Chodankar (Seat No.: 25P0320010).
  - Fully converted all economic figures, mathematical proofs, and comparison tables to Indian Rupees (₹) with baseline benchmark $1.00 USD = ₹84.00 INR (e.g. Traditional Cloud: ₹1,00,380.00/mo vs CustardFin: ₹4,128.60/mo; 95.89% savings; ₹11.55 Lakhs/yr saved).
  - Replaced generic nested boxes across all diagrams with custom creative visual metaphors:
    - Diag 1: NAT Traversal with 3D NVMe cylinder, server chassis, and brick firewall with padlock shield (resolved green reverse tunnel conduit intersection on Inbound SYN Dropped badge).
    - Diag 2: Two-tier architecture with NVIDIA GPU chip (NVENC), film strips, and cloud database barrels.
    - Diag 3: Automated ingest with SHA-256 binary key, GPU probe, 4s GOP film strip, and SQLite WAL cylinder.
    - Diag 4: Heartbeat liveness with active pulsing green ECG waveform vs severed flatline and 30s grace stopwatch.
    - Diag 5: Dual-path proximity routing with Smart TV LAN bypass (₹0.00 egress) and Origin Shield S3 cache fan-out (90% uplink saved).
    - Diag 6: Economic cost comparison with 3D stacked cost towers in INR (₹1,00,380 vs ₹4,129) and 95.9% savings medallion.
    - Diag 7: Hardware acceleration benchmarks with tachometer speedometers (322 fps vs 33 fps), donut rings, and mercury thermometers.
  - Rendered all diagrams at 3x Ultra-HD resolution (3600px width, ~400 DPI Retina quality) via headless Chrome, producing pristine, razor-sharp vector graphics in `report/main.pdf` (63 pages, 6.31 MB).

## In Progress
- Final rehearsal and submission of `report/main.pdf`.

## Blocked
- None.

## Next 3 Tasks
1. Deliver `report/main.pdf` and LaTeX source files to evaluation committee.
2. Conduct mock viva rehearsal covering architecture diagrams, Play Decision logic, and Origin Shield proofs.
3. Demonstrate live hardware-accelerated streaming on GCP VM (`34.132.188.14`) and local edge node during final evaluation.

## Decisions Resolved
- **Cache store:** MinIO retained as baseline S3-compatible Origin Shield; Floci emulator fully integrated as MIT open-source option (`deploy/floci/`).
- **Deployment:** Four complete topologies documented and tested (Pure Local, Dev Containers, Floci Multi-PC Lab, and GCP Live Cloud).
- **Reverse Tunnel Routing:** FRPS configured with `customDomains` and cloud puller injects `req.Host = nodeID` for zero-open-port home streaming.
