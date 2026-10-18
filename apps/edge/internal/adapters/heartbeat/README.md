# Edge Heartbeat Adapter (`apps/edge/internal/adapters/heartbeat`)

## Purpose
Manages periodic outbound liveness, load, and hardware capability reporting from the home edge daemon to the Cloud Control tier (`POST /v1/nodes/{id}/heartbeat`).

## Public Interface
- `Config`:
  - `CloudBaseURL`: Base URL of Cloud Control API (e.g. `http://localhost:8080`).
  - `NodeID`, `NodeToken`: Node identity and authenticated secret credentials.
  - `Version`: Running edge daemon version.
  - `Interval`: Heartbeat frequency (default: 15s).
  - `GPUCaps`: Detected hardware encoding acceleration flags (e.g. NVENC).
  - `StreamTracker`: `ports.StreamTracker` provider for active concurrent streams.
- `NewWorker(cfg Config) (*Worker, error)`: Instantiates and validates the heartbeat worker.
- `Worker.SendHeartbeat(ctx, hb) error`: Dispatches a single heartbeat immediately.
- `Worker.Start(ctx) error`: Initiates periodic background ticker routine.
- `Worker.Stop() error`: Gracefully cancels ticker and terminates worker.

## Invariant Compliance
- **Invariant 2 (Outbound-only connections):** Heartbeats are strictly outbound HTTP calls from the edge to the cloud.
- **Invariant 10 (10-stream cap):** Active streams are tracked and capped at 10 (`contracts.MaxConcurrentStreamsPerNode`).

## How to Test
```bash
go test -v ./internal/adapters/heartbeat
```
