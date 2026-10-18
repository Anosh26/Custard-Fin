## 2026-10-18 - Milestone 3: Cloud Control API (Auth, PostgreSQL Catalog, Node Pairing & Manifest Sync)
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 3 (Cloud API: Auth, Postgres, manifest and catalog endpoints)
- **Summary:** Built `apps/api` following Hexagonal Architecture and the Mandatory Granular Testing Policy. Implemented bcrypt password hashing, HMAC-SHA256 JWT bearer token issuance/validation, versioned PostgreSQL schema migrations (`users`, `nodes`, `catalog_items`, `catalog_renditions`), parameterized pgx query repositories, in-memory repositories for rapid isolated testing, and full OpenAPI-compliant HTTP endpoints (`/v1/auth/register`, `/v1/auth/login`, `/v1/nodes`, `/v1/nodes/{id}/media/{mid}/manifest`, `/v1/catalog`, `/healthz`). Enforced strict 3 KB manifest cap (Invariant 1) and node stream security. Wired `cmd/apid/main.go` daemon with graceful shutdown. All granular unit and integration tests passing.
- **Files Touched:**
  - `apps/api/cmd/apid/main.go`
  - `apps/api/internal/domain/{user.go, user_test.go, node.go, node_test.go, catalog.go, catalog_test.go, errors.go}`
  - `apps/api/internal/ports/{user_repository.go, node_repository.go, catalog_repository.go, auth.go}`
  - `apps/api/internal/adapters/auth/{bcrypt_hasher.go, bcrypt_hasher_test.go, jwt_service.go, jwt_service_test.go}`
  - `apps/api/internal/adapters/postgres/{db.go, migrations.go, user_repo.go, node_repo.go, catalog_repo.go, migrations/001_initial_schema.sql}`
  - `apps/api/internal/adapters/memory/{user_repo.go, node_repo.go, catalog_repo.go, memory_test.go}`
  - `apps/api/internal/adapters/httpapi/{server.go, server_test.go, middleware.go, README.md}`
  - `docs/context/STACK_VERSIONS.md`
  - `docs/context/STATE.md`
- **Follow-ups:**
  - Milestone 4: Edge public toggle -> manifest sync worker -> Next.js catalog UI.

## 2026-10-18 - Milestone 2: Edge Local HTTP Streaming Server (Standalone Mini-Jellyfin)
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Implemented `apps/edge/internal/adapters/httpapi` standard library `net/http` streaming server for direct LAN playback. Delivered CORS headers, structured logging, panic recovery, static HLS file serving with byte ranges (`Accept-Ranges: bytes`), media catalog endpoints (`GET /api/media`, `GET /api/media/{id}`), local watch progress tracking (`GET /api/watch/{id}`, `POST /api/watch/{id}`), and path traversal prevention. Wired HTTP server into `apps/edge/cmd/edged/main.go` with graceful context shutdown. Comprehensive unit and integration test suite passing. Milestone 2 is fully complete.
- **Files Touched:**
  - `apps/edge/internal/adapters/httpapi/server.go`
  - `apps/edge/internal/adapters/httpapi/server_test.go`
  - `apps/edge/internal/adapters/httpapi/README.md`
  - `apps/edge/cmd/edged/main.go`
  - `docs/context/STATE.md`
- **Follow-ups:**
  - Milestone 3: Cloud Control API (`apps/api`) setup (PostgreSQL schemas, `pgx`, `sqlc`, JWT auth, manifest sync, and catalog discovery).

## 2026-10-18 - Milestone 2: Managed Subprocess FFmpeg 4s HLS Transcoder
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Implemented `apps/edge/internal/adapters/ffmpeg` subprocess transcoder adhering to Invariant 6 (strict 4-second HLS segments). Added GPU probe utilizing NVIDIA NVENC (`h264_nvenc`) with verified minimum dimension constraints, automatic CPU fallback to `libx264`, forced keyframe multiples (`expr:gte(t,n_forced*4)`), multi-rendition quality ladder calculation without upscaling, and master playlist generation (`master.m3u8`). Verified with granular unit tests and real GPU encode integration test on RTX card (0.96s encode).
- **Files Touched:**
  - `apps/edge/internal/ports/transcoder.go`
  - `apps/edge/internal/adapters/ffmpeg/transcoder.go`
  - `apps/edge/internal/adapters/ffmpeg/transcoder_test.go`
  - `apps/edge/internal/adapters/ffmpeg/README.md`
  - `docs/context/STATE.md`

## 2026-10-18 - Milestone 2: Edge SQLite Store & Schema Migrations
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Implemented pure-Go SQLite store (`modernc.org/sqlite` v1.60.1) adhering to Invariant 8 (private data stays home). Delivered versioned schema migrations for `media` and `watch_progress` tables, WAL journal mode with `busy_timeout=5000`, single-writer connection pooling, `MediaRepository`, and `WatchProgressRepository` with cascading deletes. Wired scanner with SQLite repository in `edged/main.go`. Added comprehensive unit tests and verified with full integration test `TestScanner_IntegrationWithRealSQLite`.
- **Files Touched:**
  - `apps/edge/internal/adapters/sqlite/db.go`
  - `apps/edge/internal/adapters/sqlite/db_test.go`
  - `apps/edge/internal/adapters/sqlite/migrations.go`
  - `apps/edge/internal/adapters/sqlite/media_repo.go`
  - `apps/edge/internal/adapters/sqlite/media_repo_test.go`
  - `apps/edge/internal/adapters/sqlite/watch_repo.go`
  - `apps/edge/internal/adapters/sqlite/watch_repo_test.go`
  - `apps/edge/internal/adapters/sqlite/README.md`
  - `apps/edge/internal/domain/watch.go`
  - `apps/edge/internal/domain/watch_test.go`
  - `apps/edge/internal/ports/scanner.go`
  - `apps/edge/internal/app/scanner_integration_test.go`
  - `apps/edge/cmd/edged/main.go`
  - `docs/context/STATE.md`

## 2026-10-18 - Milestone 2: Mandatory Testing Policy & Comprehensive Test Suite
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 2 (Edge scan, transcode, SQLite, local stream)
- **Summary:** Updated `AGENTS.md` with the Mandatory Granular Testing Policy requiring comprehensive unit tests for every package, module, and function, plus integration tests. Built and ran exhaustive test suites across `libs/contracts`, `domain`, `adapters/fs`, `adapters/ffprobe`, and `app`, including real-world integration testing with generated ffmpeg clips and subprocess ffprobe. Caught and resolved edge cases in empty/whitespace root directory handling and test variable declarations. All quality gates and tests are green.
- **Files Touched:**
  - `AGENTS.md`
  - `libs/contracts/manifest_test.go`
  - `libs/contracts/heartbeat_test.go`
  - `libs/contracts/signing_test.go`
  - `apps/edge/internal/domain/media_test.go`
  - `apps/edge/internal/adapters/fs/hasher_test.go`
  - `apps/edge/internal/adapters/fs/scanner.go`
  - `apps/edge/internal/adapters/fs/scanner_test.go`
  - `apps/edge/internal/adapters/ffprobe/prober_test.go`
  - `apps/edge/internal/app/scan_service_test.go`
  - `apps/edge/internal/app/scanner_integration_test.go`
  - `docs/context/STATE.md`
  - `apps/edge/internal/domain/media.go`
  - `apps/edge/internal/domain/media_test.go`
  - `apps/edge/internal/ports/scanner.go`
  - `apps/edge/internal/adapters/fs/hasher.go`
  - `apps/edge/internal/adapters/fs/hasher_test.go`
  - `apps/edge/internal/adapters/fs/scanner.go`
  - `apps/edge/internal/adapters/fs/scanner_test.go`
  - `apps/edge/internal/adapters/ffprobe/prober.go`
  - `apps/edge/internal/adapters/ffprobe/prober_test.go`
  - `apps/edge/internal/app/scan_service.go`
  - `apps/edge/internal/app/scan_service_test.go`
  - `docs/context/STATE.md`
- **Follow-ups:**
  - Implement edge SQLite store (`modernc.org/sqlite`) for local metadata and watch progress.

## 2026-10-18 - Session Start & Milestone 1 Initialization
- **Author/Agents:** CustardFin Pair Programming Agent
- **Milestone:** Milestone 1 (Contracts and repo skeleton, CI, docs scaffolding, memory files)
- **Summary:** Initialized `AGENTS.md` master prompt, `CLAUDE.md`, and Context Tier 0/1 memory system (`INVARIANTS.md`, `STATE.md`, `PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).
- **Files Touched:**
  - `AGENTS.md`
  - `CLAUDE.md`
  - `docs/context/INVARIANTS.md`
  - `docs/context/STATE.md`
  - `docs/context/PROJECT_BRIEF.md`
  - `docs/context/GLOSSARY.md`
  - `docs/context/STACK_VERSIONS.md`
  - `docs/context/SESSION_LOG.md`
- **Follow-ups:**
  - Establish `libs/contracts` (Manifest & Heartbeat schema).
  - Create directory skeleton (`apps/edge`, `apps/api`, `apps/web`, `deploy`, `scripts`, `tests`).
  - Set up `go.work` and root `Makefile`.
