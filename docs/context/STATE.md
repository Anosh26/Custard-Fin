# CustardFin Project State

## Current Milestone
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.

## Done
- Initialized `AGENTS.md` and `CLAUDE.md`.
- Initialized Tier 0 and Tier 1 memory files (`INVARIANTS.md`, `STATE.md`, `PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).
- Established `libs/contracts` (Manifest & Heartbeat Go types + JSON Schemas, HMAC signed URL utils, unit tests).
- Scaffolded repository skeleton (`apps/edge`, `apps/api`, `apps/web`, `deploy`, `scripts`, `tests`, `docs/adr/0001`, `docs/api/openapi.yaml`, CI workflow).
- Verified local FFmpeg 8.0.1 and FFprobe with NVENC hardware acceleration.
- Implemented `apps/edge` media scanner:
  - `domain.Media` entity with UUID generation.
  - `ports.Scanner`, `ports.Prober`, `ports.MediaRepository` interfaces.
  - `adapters/fs.Hasher`: streaming O(1) memory SHA-256 calculation.
  - `adapters/fs.FileScanner`: recursive traversal with extension filtering and path traversal safeguards.
  - `adapters/ffprobe.SubprocessProber`: managed `ffprobe` execution with JSON stream parsing.
  - `app.ScanService`: incremental scanning with fast-path skip for unchanged files.

## In Progress
- Milestone 2: Pure-Go SQLite persistence store (`modernc.org/sqlite`) for edge media and watch progress.

## Blocked
- None.

## Next 3 Tasks
1. Implement pure-Go SQLite store (`modernc.org/sqlite`) in `apps/edge/internal/adapters/sqlite` with WAL mode and schema migrations.
2. Wire `ScanService` with SQLite repository in `apps/edge/cmd/edged/main.go`.
3. Implement managed subprocess FFmpeg 4s HLS transcoder (`apps/edge/internal/adapters/ffmpeg`) with NVENC probe and automatic CPU fallback.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
