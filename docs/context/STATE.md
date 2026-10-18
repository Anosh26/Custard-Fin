# CustardFin Project State

## Current Milestone
- **Milestone 3**: Cloud Control API (`apps/api`): Auth, PostgreSQL catalog, manifest sync, and nodes registry.

## Done
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.
  - Initialized `AGENTS.md` and `CLAUDE.md` with updated **Mandatory Granular Testing Policy**.
  - Tier 0 and Tier 1 memory files (`INVARIANTS.md`, `STATE.md`, `PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).
  - Multi-module workspace `go.work` (`libs/contracts`, `apps/edge`, `apps/api`).
  - `libs/contracts` with comprehensive unit tests for `Manifest` (<= 3 KB cap), `Heartbeat` (10-stream cap), and `Signing` (HMAC-SHA256 URL signer).
- **Milestone 2**: Edge Node Subsystems (Standalone Mini-Jellyfin).
  - Media scanner with streaming hasher (64 KB chunk O(1) RAM) and `ffprobe` managed subprocess prober.
  - Pure-Go SQLite persistence store (`modernc.org/sqlite`) with WAL mode, single-writer pooling, and versioned schema migrations.
  - Managed subprocess FFmpeg 4-second HLS transcoder (`apps/edge/internal/adapters/ffmpeg`):
    - Invariant 6: Strict 4-second HLS segments (`.ts`) with keyframe forcing (`expr:gte(t,n_forced*4)`).
    - Hardware acceleration: Probes and uses NVIDIA NVENC (`h264_nvenc`) with automatic CPU fallback (`libx264`).
    - Multi-rendition quality ladder (1080p, 720p, 480p) without upscaling.
  - Local HTTP streaming server (`apps/edge/internal/adapters/httpapi`):
    - Direct LAN playback, CORS, byte ranges (`Accept-Ranges: bytes`), and static HLS serving.
    - Endpoints: `GET /healthz`, `GET /api/media`, `GET /api/media/{id}`, `GET /stream/{id}/{path...}`, `GET /api/watch/{id}`, `POST /api/watch/{id}`.
  - Wired into `apps/edge/cmd/edged/main.go` with graceful shutdown.
  - 100% of unit and integration tests passing.

## In Progress
- Milestone 3: Cloud Control API architecture design, PostgreSQL schema, JWT auth, and manifest/catalog endpoints.

## Blocked
- None.

## Next 3 Tasks
1. Review Cloud Control API implementation design with the user (PostgreSQL schema, pgx, JWT/Argon2id, manifest/catalog endpoints).
2. Implement PostgreSQL database migrations and repository in `apps/api/internal/adapters/postgres`.
3. Implement Auth & Node registration endpoints matching `docs/api/openapi.yaml`.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
