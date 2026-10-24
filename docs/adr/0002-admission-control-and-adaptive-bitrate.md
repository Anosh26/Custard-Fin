# ADR 0002: Admission Control, Adaptive Bitrate Throttling, and First-Chunk Prefetch

## Status
Accepted

## Context
In CustardFin, home nodes host video binaries and stream uncached chunks across an FRP reverse tunnel relay to the cloud. While the Origin Shield cache eliminates upstream load when multiple viewers watch the same media, concurrent viewers watching distinct, uncached titles compete for limited home upload bandwidth (typically 10–50 Mbps). Unmitigated multi-stream 1080p demand leads to packet queue bufferbloat, tunnel stalls, and playback buffering. Furthermore, startup latency (TTFB) over tunnel relays introduces initial viewer delay.

## Decision
1. **Saturation Estimation**:
   - Define `ports.SaturationEstimator` in `apps/api/internal/ports/admission.go`.
   - Implement `admission.DefaultEstimator` assessing concurrent streams from Redis heartbeats (`heartbeat:{node_id}`):
     - `LoadNormal` (0–3 streams): Full quality ladder (1080p, 720p, 480p).
     - `LoadModerate` (4–6 streams): Throttle out 1080p, serving 720p and 480p.
     - `LoadSaturated` (7–9 streams): Strict throttle, serving 480p only.
     - `LoadExhausted` (10+ streams): Hard cap rejection with HTTP 429 `NODE_BUSY` (Invariant 10).
2. **Dynamic Playlist Rendition Throttling**:
   - In `apps/api/internal/adapters/httpapi`, enhance the master playlist rewriter (`rewritePlaylistWithLoad`) to dynamically prune variant streams (`#EXT-X-STREAM-INF`) disallowed under current node saturation.
   - HLS players (`hls.js`) automatically negotiate available variants without client-side error.
   - Safety fallback: If all variants would be pruned, preserve the lowest rendition rather than outputting an empty playlist.
3. **First-Chunk Prefetching**:
   - Implement `prefetch.BackgroundPrefetcher` in `apps/api/internal/adapters/prefetch`.
   - On `GET /v1/media/{id}/playback`, asynchronously preheat the initial segment (`segment_000.ts`) into the S3/memory Origin Shield cache (`ports.ChunkCache`) via a bounded non-blocking goroutine.
   - Subsequent client chunk requests result in immediate cache `HIT`s, eliminating TTFB playback delay.

## Consequences
- **Positive**:
  - Protects home upstream bandwidth against bufferbloat during multi-viewer sessions.
  - Decreases bandwidth consumption by 40%–65% under moderate/high load.
  - Eliminates startup buffering on playback initiation.
  - Fully transparent to client video players with zero proprietary player SDK dependencies.
- **Negative**:
  - Viewers on saturated nodes may temporarily stream at 720p or 480p until upstream capacity eases.
