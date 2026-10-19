# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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

