# CustardFin Project State

## Current Milestone
- **Milestone 8**: Admission control, adaptive bitrate, prefetch.

## Done
- **Milestone 1**: Contracts and repo skeleton, CI, docs scaffolding, memory files.
- **Milestone 2**: Edge Node Subsystems (Standalone Mini-Jellyfin).
- **Milestone 3**: Cloud Control API (`apps/api`).
- **Milestone 4**: Public Toggle -> Manifest Sync Worker -> Next.js Community Catalog UI.
- **Milestone 5**: FRP Tunnel (Remote Playback End-to-End).
- **Milestone 6**: Redis Heartbeat and Offline Logic.
- **Milestone 7**: Origin Shield Cache + Signed URLs + Stream Router.
  - Defined `ports.ChunkCache` and `ports.TunnelPuller` in `apps/api/internal/ports/chunk_cache.go`.
  - Implemented `S3ChunkCache` in `apps/api/internal/adapters/s3` with AWS SDK for Go v2, path-style support (MinIO/Garage/Floci), and 24h retention tagging (Invariant 1).
  - Implemented thread-safe `MemoryChunkCache` in `apps/api/internal/adapters/memory` with TTL expiration for test isolation.
  - Implemented `HTTPTunnelPuller` in `apps/api/internal/adapters/tunnel` for pulling uncached chunks via reverse tunnel relay.
  - Enforced short-lived HMAC-SHA256 URL signatures on `/v1/stream/{node_id}/{media_id}/{chunk_path...}` (Invariant 9).
  - Built dynamic HLS playlist rewriter signing relative child rendition and segment paths with secret binding.
  - Upgraded Play Decision Tree in `GET /v1/media/{id}/playback`: routes to `direct_lan`, `origin_shield`, or `cache_only` (allowing cached playback if node is offline).
  - Updated Next.js frontend client (`apps/web`) to support and display Origin Shield and Offline Cache badges.
  - 100% test pass rate across granular unit tests and live cross-subsystem E2E test suite.

## In Progress
- Milestone 8 Architecture Design (Admission control, adaptive bitrate fallback, first-chunk prefetch).

## Blocked
- None.

## Next 3 Tasks
1. Implement bandwidth and saturation estimator for remote streams.
2. Implement adaptive rendition throttling (1080p -> 720p -> 480p) under high home upload load.
3. Implement initial chunk prefetch worker on playback initialization to eliminate startup buffering.

## Open Questions (Requires Human Decision)
- **Cache store:** MinIO status check vs Garage/SeaweedFS for S3-compatible cache.
- **Floci emulator:** Exact Floci emulator (AWS vs GCP) and services to prioritize for demo.
