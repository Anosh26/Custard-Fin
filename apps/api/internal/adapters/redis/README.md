# Redis Adapter (`apps/api/internal/adapters/redis`)

## Purpose
Manages ephemeral node heartbeats in Redis. Strictly complies with **Invariant 5** (*"Redis is used for heartbeats only. One key per node `heartbeat:{node_id}` with TTL; nothing else"*).

## Public Interface
- `NewHeartbeatStore(client *redis.Client) *HeartbeatStore`: Constructs a Redis-backed heartbeat store.
  - `RecordHeartbeat(ctx, hb, ttl)`: Validates and serializes the `Heartbeat` struct to JSON with a short TTL (default: 20s).
  - `IsNodeOnline(ctx, nodeID)`: Fast existence check (`EXISTS heartbeat:{node_id}`).
  - `GetHeartbeat(ctx, nodeID)`: Retrieves active heartbeat payload (e.g. stream count and GPU capabilities).
  - `BatchIsNodeOnline(ctx, nodeIDs)`: Pipelines `EXISTS` commands across multiple node IDs in a single round-trip.

## Invariant 5 Enforcement
No background jobs, queues, or other state are permitted in Redis. Background jobs in CustardFin are Go goroutines with state in PostgreSQL.
