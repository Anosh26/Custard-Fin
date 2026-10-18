# apps/api (apid)

CustardFin Cloud Control API server.

## Purpose
- Coordinates discovery, users, node registrations, and public media catalog.
- Evaluates stream router decisions (LAN direct vs Origin cache vs tunnel pull).
- Origin Shield pull-through cache manager: fetches chunks over tunnel once and stores temporarily in S3-compatible cache.
- Tracks node availability via Redis heartbeats (`heartbeat:{node_id}`).
- Enforces admission control (max 10 concurrent remote streams per node).

## Hexagonal Package Structure
- `cmd/apid/`: Main entrypoint and DI wiring.
- `internal/domain/`: User, Node, CatalogItem, RoutingDecision.
- `internal/app/`: Use cases (AuthService, CatalogService, HeartbeatService, PlaybackRouter, OriginShieldService).
- `internal/ports/`: Interfaces (`UserRepo`, `CatalogRepo`, `HeartbeatTracker`, `ChunkCache`, `TunnelGateway`).
- `internal/adapters/`:
  - `postgres/`: PostgreSQL database persistence using `pgx` and `sqlc`.
  - `redis/`: Redis client for node heartbeats only.
  - `s3/`: S3-compatible chunk cache adapter (AWS SDK for Go v2).
  - `httpapi/`: REST handlers conforming to `docs/api/openapi.yaml`.
  - `workers/`: Go goroutine background workers (cache purger, sync processor).

## Running Locally
```bash
go run ./cmd/apid
```
