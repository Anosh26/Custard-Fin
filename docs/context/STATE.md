# CustardFin Project State

## Current Milestone
- **Milestone 7**: Origin Shield cache + signed URLs + stream router.

## Done
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.
- **Milestone 2**: Edge Node Subsystems (Standalone Mini-Jellyfin).
- **Milestone 3**: Cloud Control API (`apps/api`).
- **Milestone 4**: Public Toggle -> Manifest Sync Worker -> Next.js Community Catalog UI.
- **Milestone 5**: FRP Tunnel (Remote Playback End-to-End).
- **Milestone 6**: Redis Heartbeat and Offline Logic.
  - Defined `ports.HeartbeatStore` in `apps/api/internal/ports/heartbeat_store.go`.
  - Implemented `apps/api/internal/adapters/redis/heartbeat_store.go` adhering to Invariant 5 (Redis used strictly for `heartbeat:{node_id}` TTL keys).
  - Implemented `apps/api/internal/adapters/memory/heartbeat_store.go` with simulated error injection for tests.
  - Registered `POST /v1/nodes/{id}/heartbeat` with `RequireNodeAuth` token verification and Invariant 10 validation.
  - Enriched `GET /v1/catalog` with batch node liveness checks (`is_node_online`).
  - Implemented Play Decision Tree in `GET /v1/media/{id}/playback`: returns HTTP 503 `NODE_OFFLINE` when key is absent/expired, and HTTP 429 `NODE_BUSY` with `Retry-After: 30` when active streams reach cap ($\ge 10$).
  - Built edge heartbeat adapter (`apps/edge/internal/adapters/heartbeat`) running periodic background reporter (15s ticker) with GPU caps and stream load reporting.
  - Implemented `ports.StreamTracker` on edge HTTP server with automatic session pruning.
  - Wired into both `cmd/apid` and `cmd/edged` with 12-factor configuration and graceful shutdown.
  - 100% test pass rate across unit tests and live cross-subsystem E2E test suite.

## In Progress
- Milestone 7 Architecture Design & Setup (S3-compatible ChunkCache adapter, Origin Shield pull-through cache, HMAC-SHA256 signed chunk URLs).

## Blocked
- None.

## Next 3 Tasks
1. Design `ports.ChunkCache` interface and S3-compatible adapter (`apps/api/internal/adapters/s3`).
2. Implement short-lived signed URL generator and validator for hot chunks (HMAC with node/media/expiry binding).
3. Implement stream router cache-hit vs miss pull-through mechanism via tunnel.

## Open Questions (Requires Human Decision)
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.

