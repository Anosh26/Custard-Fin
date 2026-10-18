# CustardFin Project State

## Current Milestone
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.

## Done
- Initialized `AGENTS.md` and `CLAUDE.md` with updated **Mandatory Granular Testing Policy**.
- Initialized Tier 0 and Tier 1 memory files (`INVARIANTS.md`, `STATE.md`, `PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).
- Established `libs/contracts` with comprehensive unit tests for `Manifest`, `Heartbeat`, and `Signing`.
- Scaffolded repository skeleton (`apps/edge`, `apps/api`, `apps/web`, `deploy`, `scripts`, `tests`, `docs/adr/0001`, `docs/api/openapi.yaml`, CI workflow).
- Verified local FFmpeg 8.0.1 and FFprobe with NVENC hardware acceleration.
- Implemented `apps/edge` media scanner with streaming hasher, `ffprobe` prober, and incremental change detection.
- Implemented pure-Go SQLite persistence store (`modernc.org/sqlite`) in `apps/edge/internal/adapters/sqlite`:
  - Versioned schema migrations (`media`, `watch_progress`, `schema_migrations`).
  - WAL journal mode, `busy_timeout=5000`, foreign keys enabled with cascading deletes.
  - `MediaRepository` and `WatchProgressRepository` implementations.
  - Wired into `apps/edge/cmd/edged/main.go` with startup directory scan.
  - Comprehensive unit tests and real SQLite integration test (`TestScanner_IntegrationWithRealSQLite`).
- Implemented managed subprocess FFmpeg 4s HLS transcoder (`apps/edge/internal/adapters/ffmpeg`):
  - Invariant 6: Strictly chunked into 4-second HLS segments (`.ts`) with keyframe forcing (`expr:gte(t,n_forced*4)`).
  - Hardware acceleration: Probes and uses NVIDIA NVENC (`h264_nvenc`) with automatic fallback to CPU (`libx264`).
  - Multi-rendition quality ladder (1080p, 720p, 480p) without upscaling.
  - Automated master playlist generation (`master.m3u8`).
  - Granular unit tests and real-world NVENC + CPU integration tests verified.

## In Progress
- Milestone 2: Edge local HTTP streaming server (`apps/edge/internal/adapters/httpapi`) for LAN direct playback.

## Blocked
- None.

## Next 3 Tasks
1. Implement Edge local HTTP streaming server (`apps/edge/internal/adapters/httpapi`) for LAN direct playback with CORS and byte ranges.
2. Complete Milestone 2 local streaming test (standalone mini-Jellyfin verification).
3. Transition to Milestone 3: Cloud API (`apps/api`) auth, PostgreSQL catalog, and manifest sync.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
