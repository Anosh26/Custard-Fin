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

## In Progress
- Milestone 2: Managed subprocess FFmpeg 4-second HLS transcoder (`apps/edge/internal/adapters/ffmpeg`) with NVENC acceleration and automatic CPU fallback.

## Blocked
- None.

## Next 3 Tasks
1. Implement managed subprocess FFmpeg HLS transcoder (NVENC probe, 4s chunks, 1080p/720p/480p ladders, CPU fallback).
2. Implement Edge local HTTP streaming server (`apps/edge/internal/adapters/httpapi`) for direct LAN streaming.
3. Transition to Milestone 3: Cloud API (`apps/api`) auth, PostgreSQL catalog, and manifest sync.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
