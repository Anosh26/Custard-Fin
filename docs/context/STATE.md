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
  - OpenXML PowerPoint presentation (`docs/presentation/CustardFin_Presentation.pptx`), OpenXML dump (`presentation.xml`), and interactive web slides (`docs/presentation/index.html`) with pastel palette and animated custard mascot.

- **User Accounts, Edge Node Claiming & Playback Resilience**:
  - Cloud API (`apps/api`): Added `POST /v1/nodes/{id}/claim` allowing authenticated users (e.g. `@anish`) to claim / bind their edge node, instantly re-attributing catalog uploads and granting private video visibility. Added `owner_username` and `owner_user_id` to `GET /v1/nodes/{id}/status`. Fixed prefetch segment name typo to `1080p/segment_0000.ts`.
  - Edge Node Daemon (`apps/edge`): Modified `edged/main.go` and `server.go` to generate crisp poster thumbnails immediately and ONLY sync media manifests to the Cloud catalog once `master.m3u8` is confirmed on disk, preventing catalog publishing during in-flight encoding.
  - Web UI (`apps/web`): Added `claimNode` and `fetchNodeStatus` in `api-client.ts`; added node ownership banner and 1-click `Link to @username` button in `EdgeLibraryModal.tsx`. Enhanced `Player.tsx` with automatic re-polling during stream warm-up and resilient retry handling.

## In Progress
- Final verification of node claiming, uploader attribution, and large video playback on GCP VM.

## Blocked
- None.

## Next 3 Tasks
1. Conduct mock viva rehearsal covering architecture diagrams, Play Decision logic, and Origin Shield proofs.
2. Demonstrate live hardware-accelerated streaming on GCP VM (`34.132.188.14`) and local edge node during final evaluation.
3. Review Q&A cheat sheet in `docs/runbooks/demo.md` for team members Anish Ghadi & Anosh Chodankar.

## Decisions Resolved
- **Cache store:** MinIO retained as baseline S3-compatible Origin Shield; Floci emulator fully integrated as MIT open-source option (`deploy/floci/`).
- **Deployment:** Four complete topologies documented and tested (Pure Local, Dev Containers, Floci Multi-PC Lab, and GCP Live Cloud).
- **Reverse Tunnel Routing:** FRPS configured with `customDomains` and cloud puller injects `req.Host = nodeID` for zero-open-port home streaming.
