# CustardFin Project State

## Current Milestone
- **Milestone 10**: Hardening, demo profile, rehearsal, runbooks.

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
  - Implemented multi-encoder hardware capability probe (`ProbeCapabilities`) checking NVENC, QSV, and VAAPI.
  - Implemented automatic graceful CPU fallback (`libx264`) on any hardware encoder failure during HLS chunking (Invariant 6).
  - Built standalone benchmark suite in `scripts/bench-gpu.sh` and `scripts/bench-gpu.ps1` measuring execution wall-clock time, speed multiplier (up to 10.7x realtime on GeForce GTX), and speedup ratio.
  - Built automated media seeder `scripts/seed-media.sh` generating synthetic 1080p, 720p, and 480p fixtures into `./media`.
  - Added Makefile targets `make bench` and `make seed`.

## In Progress
- Milestone 10: Demo profile, load simulation scripts, and viva presentation runbook.

## Blocked
- None.

## Next 3 Tasks
1. Implement 10-viewer load simulation script (`scripts/load-10-viewers.sh`).
2. Create self-hosted demo docker-compose profile (`deploy/compose.demo.yml`).
3. Author step-by-step presentation runbook (`docs/runbooks/demo.md`).

## Open Questions (Requires Human Decision)
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
