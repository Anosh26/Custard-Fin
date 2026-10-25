# CustardFin Project State

## Current Milestone
- **Milestone 10**: Hardening, demo profile, rehearsal, runbooks (COMPLETE).
- **Status**: All 10 Milestones successfully implemented, tested, and documented. Ready for final evaluation.

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
- **Post-Milestone 10**: Interactive Video Management, Dev Console & Multi-Terminal Launcher.
  - Edge ingestion API: `POST /api/media/upload`, `POST /api/media/scan`, `POST /api/media/generate-sample`, and `GET /api/status`.
  - Next.js UI Modals: `AddVideoModal` (drag-and-drop / scan / sample generator), `EdgeLibraryModal` (local videos & 1-click Public/Private sync toggle), `DevSettingsModal` (endpoint switcher, latency ping, Invariant checklist, stream limit simulator), and Player "Stats for Nerds" live HUD.
  - Multi-terminal demo launcher (`scripts/start-all-terminals.ps1`, `scripts/start-all-terminals.bat`, `scripts/tab1`..`tab5`) providing 1-click terminal orchestration across Docker, Cloud API, Edge NVENC Daemon, Web UI, and Simulation.
  - Demo auto-pairing & host PostgreSQL port 5433 conflict resolution for rock-solid heartbeats and playback.

## In Progress
- Final demonstration rehearsal and viva preparation.

## Blocked
- None.

## Next 3 Tasks
1. Select demo topology (Mode 3 Floci Multi-PC Lab vs Mode 4 GCP Cloud) based on evaluation venue Wi-Fi.
2. Conduct mock viva presentation rehearsing Q&A from `docs/runbooks/demo.md`.
3. Stand up stack and demonstrate live hardware-accelerated streaming.

## Decisions Resolved
- **Cache store:** MinIO retained as baseline S3-compatible Origin Shield; Floci emulator fully integrated as MIT open-source option (`deploy/floci/`).
- **Deployment:** Four complete topologies documented and tested (Pure Local, Dev Containers, Floci Multi-PC Lab, and GCP Live Cloud).
