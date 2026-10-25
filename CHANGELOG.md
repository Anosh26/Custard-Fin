# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Post-Milestone 10 (Interactive Video Upload, Edge Library Manager & Dev Console)**:
  - Edge Daemon Video Ingestion API: added `POST /api/media/upload` (multipart video upload with auto-transcode and manifest sync), `POST /api/media/scan` (on-demand directory re-scan), `POST /api/media/generate-sample` (1-click 15s synthetic FFmpeg test card generator), and `GET /api/status` (live node stream count and NVENC hardware probe).
  - Web UI `+ Add Video` Modal (`apps/web/src/components/AddVideoModal.tsx`) featuring drag-and-drop file upload, folder scan trigger, and 1-click test clip synthesizer.
  - Edge Library Manager Modal (`apps/web/src/components/EdgeLibraryModal.tsx`) showing all local edge videos with resolution, duration, storage path, and 1-click "Public / Private" community sharing switch.
  - Dev Settings Console (`apps/web/src/components/DevSettingsModal.tsx`) with dynamic Cloud API and Edge Node endpoint switcher, live latency ping tester, Invariant verification checklist (Invariants 1, 2, 5, 6, 7, 9, 10), and live admission control concurrency simulation button.
  - Player "Stats for Nerds" HUD in `Player.tsx` displaying stream protocol (HLS RFC 8216), 4s segment enforcement, active rendition/bitrate, Origin Shield cache status, and HMAC token validity.
  - Development Demo Auto-Seeding: `apid` automatically pre-seeds demo user and edge node credentials, resolving host PostgreSQL port 5432 conflict by mapping dev container to 5433.
  - Multi-Terminal Demo Launcher: added `scripts/start-all-terminals.ps1`, `scripts/start-all-terminals.bat`, and individual scripts `scripts/tab1-docker.ps1` through `scripts/tab5-bench.ps1` to launch synchronized terminal tabs in Windows Terminal (or separate PowerShell windows) for the 5-layer presentation stack.
- **Milestone 10 (Demo Profile, Load Simulation Scripts, GCP Deployment & Viva Runbooks)**:
  - 10-viewer load simulation scripts (`scripts/load-10-viewers.sh` and `scripts/load-10-viewers.ps1`) demonstrating Origin Shield cache hits, adaptive ladder throttling, and 10-stream cap.
  - Multi-stage Dockerfiles (`apps/api/Dockerfile`, `apps/web/Dockerfile`) and unified demo orchestration (`deploy/compose.demo.yml`, `deploy/Caddyfile`).
  - Comprehensive runbooks: GCP deployment guide (`docs/runbooks/deploy.md`), viva presentation script (`docs/runbooks/demo.md`), troubleshooting guide (`docs/runbooks/troubleshooting.md`), and universal running guide (`RUNNING.md`).
  - Open-source Floci cloud emulator deployment profile (`deploy/floci/compose.floci.yml`).
- **Milestone 9 (NVENC Path with Probe, Fallback and Benchmark Script)**:
  - Hardware accelerated NVENC encoding adapter (`apps/edge/internal/adapters/ffmpeg`) with probe order NVENC -> QSV/VAAPI -> libx264 CPU fallback.
  - Test-encode verification and strict 4-second keyframe segment alignment.
  - Cross-platform GPU benchmarking scripts (`scripts/bench-gpu.sh` and `scripts/bench-gpu.ps1`).
- **Milestone 8 (Admission Control, Adaptive Bitrate & Prefetch)**:
  - In-memory and Redis-backed admission control limiter enforcing the 10-stream per-node cap (Invariant 10) with HTTP 429 Retry-After.
  - Adaptive bitrate ladder pruner dynamically filtering high-bitrate master playlist renditions under load.
  - Background first-chunk prefetch worker (`apps/api/internal/adapters/prefetch`) pre-warming Origin Shield cache upon playback resolution.
- **Milestone 7 (Origin Shield Cache, Signed URLs & Stream Router)**:
  - AWS SDK for Go v2 S3 chunk cache adapter (`apps/api/internal/adapters/s3`) with 24h retention tagging.
  - In-memory chunk cache adapter (`apps/api/internal/adapters/memory`) for testing.
  - HTTP tunnel puller adapter (`apps/api/internal/adapters/tunnel`) for pulling chunks over FRP reverse tunnel.
  - HMAC-SHA256 signed URL verification on `/v1/stream/{node_id}/{media_id}/{chunk_path...}`.
  - Dynamic HLS playlist rewriter signing relative child rendition and segment paths.
  - Play Decision Tree upgraded to support `origin_shield`, `direct_lan`, and `cache_only` routes.
  - Next.js player UI updated with dynamic route badges.
- **Milestone 6 (Redis Heartbeat & Offline Logic)**:
  - Redis heartbeat adapter (`apps/api/internal/adapters/redis`) adhering to Invariant 5.
  - Edge heartbeat worker (`apps/edge/internal/adapters/heartbeat`) reporting NVENC capabilities and stream count.
  - Node offline detection and 10-stream cap admission enforcement (Invariant 10).
- **Milestone 5 (FRP Reverse Tunnel)**:
  - Managed subprocess FRP client (`apps/edge/internal/adapters/frp`) with TLS and CGNAT traversal.
- **Milestone 4 (Public Toggle & Next.js Catalog UI)**:
  - Edge manifest sync worker (`apps/edge/internal/adapters/sync`) enforcing <= 3 KB manifest cap.
  - Next.js App Router community catalog and HLS player component (`apps/web`).
- **Milestone 3 (Cloud API Core & PostgreSQL Catalog)**:
  - JWT auth, bcrypt password hashing, pgx repository migrations, and OpenAPI 3.0 endpoints.
- **Milestone 2 (Edge Node Subsystems)**:
  - Local directory scanner, SQLite persistence, and FFmpeg HLS pre-transcoder (NVENC/CPU).
- **Milestone 1 (Foundational Skeleton)**:
  - Contracts, docs scaffolding, memory files, CI pipeline, and project structure.

