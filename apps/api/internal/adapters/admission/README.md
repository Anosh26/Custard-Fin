# Admission Package (`apps/api/internal/adapters/admission`)

## Purpose
Evaluates home edge node saturation levels based on heartbeats recorded in Redis (`ports.HeartbeatStore`), implementing admission control and adaptive bitrate throttling levels (FR-10).

## Public Interface
- `NewDefaultEstimator(cfg Config, heartbeatStore ports.HeartbeatStore, logger *slog.Logger) *DefaultEstimator`
- `EstimateLoad(ctx context.Context, nodeID string) (domain.LoadLevel, int, error)`

## Thresholds & Load Levels
- **Normal** (0–3 active streams): Full ladder available (1080p, 720p, 480p).
- **Moderate** (4–6 active streams): 1080p throttled out, serving 720p and 480p.
- **Saturated** (7–9 active streams): 1080p and 720p throttled out, serving 480p only.
- **Exhausted** (10+ active streams): Hard limit reached, returns HTTP 429 `NODE_BUSY` (Invariant 10).
