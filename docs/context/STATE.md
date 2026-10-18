# CustardFin Project State

## Current Milestone
- **Milestone 4**: Public toggle -> manifest sync worker -> Next.js community catalog UI.

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
- **Milestone 3**: Cloud Control API (`apps/api`).
  - Hexagonal architecture: domain entities (`User`, `Node`, `CatalogItem`, `Rendition`), ports, adapters (`auth`, `postgres`, `memory`, `httpapi`).
  - Password hashing with bcrypt, JWT bearer token issuance and validation.
  - PostgreSQL versioned schema migrations (`001_initial_schema.sql` for users, nodes, catalog_items, catalog_renditions) with parameterized `pgx` queries.
  - In-memory repository adapters for rapid, isolated unit testing.
  - REST endpoints adhering to `docs/api/openapi.yaml`: `/v1/auth/register`, `/v1/auth/login`, `/v1/nodes`, `/v1/nodes/{id}/media/{mid}/manifest`, `/v1/catalog`, `/healthz`.
  - Invariant 1 enforced: manifest payloads strictly capped at <= 3 KB.
  - Main daemon `cmd/apid/main.go` wired with 12-factor config and graceful shutdown.
  - 100% test coverage across all modules and full HTTP integration suite passing.

## In Progress
- Milestone 4: Edge public toggle -> manifest sync worker -> Next.js community catalog UI.

## Blocked
- None.

## Next 3 Tasks
1. Explain Milestone 4 implementation design to user (Edge manifest sync worker on public toggle, Next.js App Router catalog UI).
2. Implement Edge sync worker in `apps/edge` (`PUT /v1/nodes/{id}/media/{mid}/manifest` when marked public).
3. Scaffold and build Next.js App Router UI (`apps/web`) with server components and community catalog discovery.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
