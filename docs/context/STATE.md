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
  - Added Makefile targets `make bench`, `make seed`, `make load-test`.
  - Authored comprehensive GCP deployment guide (`docs/runbooks/deploy.md`) with VPC firewall rules, ports reference, and edge client configuration.
  - Authored viva presentation script and defense guide (`docs/runbooks/demo.md`).
  - Authored emergency recovery and troubleshooting guide (`docs/runbooks/troubleshooting.md`).

## In Progress
- Rehearsal and final project demonstration preparation.

## Blocked
- None.

## Next 3 Tasks
1. Review GCP deployment guide (`docs/runbooks/deploy.md`) with professor if live cloud deployment is required.
2. Rehearse live demonstration flow following `docs/runbooks/demo.md`.
3. Merge pull request #1 into `master`.

## Decisions Resolved
- **Cache store:** MinIO retained as stable S3-compatible Origin Shield with visual console.
- **Deployment:** Self-hosted / local container stack as baseline; full GCP deployment runbook ready for live deployment with $40 credits.
