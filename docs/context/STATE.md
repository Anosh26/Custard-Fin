# CustardFin Project State

## Current Milestone
- **Milestone 9**: NVENC path with probe, fallback and benchmark script.

## Done
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.
- **Milestone 2**: Edge Node Subsystems (Standalone Mini-Jellyfin).
- **Milestone 3**: Cloud Control API (`apps/api`).
- **Milestone 4**: Public Toggle -> Manifest Sync Worker -> Next.js Community Catalog UI.
- **Milestone 5**: FRP Tunnel (Remote Playback End-to-End).
- **Milestone 6**: Redis Heartbeat and Offline Logic.
- **Milestone 7**: Origin Shield Cache + Signed URLs + Stream Router.
- **Milestone 8**: Admission Control, Adaptive Bitrate & Prefetch.
  - Defined `ports.SaturationEstimator` and `ports.Prefetcher` in `apps/api/internal/ports/admission.go`.
  - Implemented `domain.LoadLevel` (`LoadNormal`, `LoadModerate`, `LoadSaturated`, `LoadExhausted`) and resolution allowance rules.
  - Built `DefaultEstimator` in `apps/api/internal/adapters/admission` mapping heartbeat active stream counts to saturation tiers.
  - Implemented `BackgroundPrefetcher` in `apps/api/internal/adapters/prefetch` to warm Origin Shield cache ahead of time on playback start.
  - Implemented dynamic master playlist rendition throttling in `httpapi.rewritePlaylistWithLoad`, pruning 1080p under moderate load and 1080p+720p under saturated load, with single-rendition safety fallback.
  - Wired prefetcher into `handleGetPlayback` and saturation estimator into `handleStreamChunk`.
  - Added full test coverage for admission, prefetch, and adaptive playlist rewriting.

## In Progress
- Milestone 9: Hardware acceleration probe refinement and benchmark scripts.

## Blocked
- None.

## Next 3 Tasks
1. Implement NVENC hardware probe validation and benchmark script (`scripts/bench-gpu`).
2. Verify CPU fallback transition under GPU saturation or unsupported codec formats.
3. Advance to Milestone 10 (Self-hosted cloud demo profile, rehearsal scripts, demo runbooks).

## Open Questions (Requires Human Decision)
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
