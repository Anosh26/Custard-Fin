# CustardFin Project State

## Current Milestone
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.

## Done
- Initialized `AGENTS.md` and `CLAUDE.md`.
- Initialized Tier 0 and Tier 1 memory files (`INVARIANTS.md`, `STATE.md`, `PROJECT_BRIEF.md`, `GLOSSARY.md`, `STACK_VERSIONS.md`, `SESSION_LOG.md`).
- Established `libs/contracts` (Manifest & Heartbeat Go types + JSON Schemas, HMAC signed URL utils, unit tests).
- Scaffolded repository skeleton (`apps/edge`, `apps/api`, `apps/web`, `deploy`, `scripts`, `tests`, `docs/adr/0001`, `docs/api/openapi.yaml`, CI workflow).
- Verified local FFmpeg 8.0.1 with NVENC support (`h264_nvenc`, `hevc_nvenc`, `av1_nvenc`).

## In Progress
- Host Go toolchain availability (winget GoLang.Go ready to install).
- Transitioning to Milestone 2: Edge local scanner, SQLite storage, FFmpeg transcoder.

## Blocked
- None.

## Next 3 Tasks
1. Confirm/install Go toolchain (`winget install --id GoLang.Go -e`) to execute and test Go code directly.
2. Milestone 2: Implement edge media scanner and SQLite database (`modernc.org/sqlite`) for local metadata and watch progress.
3. Milestone 2: Implement managed subprocess FFmpeg HLS transcoder (4s chunks, NVENC with CPU fallback).

## Open Questions (Requires Human Decision)
- **Heartbeat timing:** BRD target "offline within 15 s" vs grace window. (Currently will make TTL configurable, defaulting to 15s).
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **FRP deployment:** Running FRP client/server as managed child process vs embedded.
