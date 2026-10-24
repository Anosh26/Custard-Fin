# Prefetch Package (`apps/api/internal/adapters/prefetch`)

## Purpose
Proactively preheats the Origin Shield cache with initial video chunks (`segment_000.ts`) when playback is initiated, eliminating startup time-to-first-byte (TTFB) delay and buffering over reverse tunnels.

## Public Interface
- `NewBackgroundPrefetcher(cfg Config, chunkCache ports.ChunkCache, tunnelPuller ports.TunnelPuller, logger *slog.Logger) (*BackgroundPrefetcher, error)`
- `Prefetch(ctx context.Context, nodeID, mediaID, chunkPath string) error`
- `TriggerPrefetch(ctx context.Context, nodeID, mediaID, chunkPath string)`
- `Wait()`
