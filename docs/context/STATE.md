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

- **User Accounts & Private Video Ownership**:
  - Cloud API (`apps/api`): Added `/v1/auth/me`, enriched auth responses with `user_id` and `username`. Catalog search now accepts `allowedNodeIDs` permitting owners to discover both public and private items.
  - Strict privacy enforcement: `/v1/media/{id}/playback` and `/v1/media/{id}/thumbnail` strictly require direct LAN or authenticated node owner credentials for private items; unauthorized callers receive HTTP 403 Forbidden.
  - Enriched catalog and playback responses with `uploader_username`, `node_name`, and `is_public`.
  - Edge Node (`apps/edge`): `HTTPManifestSyncer` transmits actual `media.IsPublic` state; private toggle syncs lightweight manifest with `is_public: false` instead of un-syncing/deleting, keeping the owner aware on the Cloud.
  - Web UI (`apps/web`): Created `AuthModal` (Sign In / Register / 1-click Quick Demo Login as `@demo`), added user account button in `NavbarClientControls`, created `CatalogFeed` and `WatchClient` components showing `@username • node_name` and `🔒 Private (Only You)` badges with immediate token awareness. Added direct "Watch" links in `EdgeLibraryModal`.

## In Progress
- Live cloud deployment and evaluation of user accounts and private video visibility.

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
