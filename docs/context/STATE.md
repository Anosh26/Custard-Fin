# CustardFin Project State

## Current Milestone
- **Milestone 5**: FRP Tunnel — Remote Playback End-to-End (Outbound-only home edge tunnel to relay server).

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
- **Milestone 4**: Public Toggle -> Manifest Sync Worker -> Next.js Community Catalog UI.
  - Edge manifest syncer (`apps/edge/internal/adapters/sync`): enforces <= 3 KB manifest cap, 4s segment constraint, dispatches to `PUT /v1/nodes/{id}/media/{mid}/manifest` and `DELETE ...`.
  - Edge public toggle: `PUT /api/media/{id}/public` triggers automatic manifest sync/unsync and persists `is_public` in SQLite.
  - Cloud API stream routing: `GET /v1/media/{id}/playback` endpoint with direct LAN stream resolution.
  - Next.js Web UI (`apps/web`):
    - Next.js 15.1.0 App Router with TypeScript strict mode, React Server Components.
    - HLS player client component (`apps/web/src/components/Player.tsx`) using `hls.js` with quality rendition switcher and native Safari fallback.
    - Glassmorphism dark aesthetic design system (`globals.css`) with responsive video cards, online status indicators, and duration badges.
    - Community discovery catalog (`/`) with search filter `?q=`, and watch player route (`/watch/[id]`).
    - 100% build and lint passing: `next build`, `tsc --noEmit`, `eslint`.

## In Progress
- Milestone 5 Architecture Design & Setup (FRP tunnel: outbound TLS connection from edge to cloud relay).

## Blocked
- None.

## Next 3 Tasks
1. Design Milestone 5 FRP tunnel architecture (outbound-only TLS connection, CGNAT bypass, managed subprocess vs embedded client).
2. Create `deploy/frps.toml` configuration and relay test container setup.
3. Implement `apps/edge/internal/adapters/frp` managed client adapter with auto-reconnect backoff.

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
