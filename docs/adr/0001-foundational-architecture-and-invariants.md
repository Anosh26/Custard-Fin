# ADR 0001: Foundational Architecture and Invariants

## Status
Accepted

## Context
CustardFin requires an architecture that combines personal media streaming (Jellyfin/Plex style) with community sharing (YouTube style) for a self-hosted cloud demo at the end of the MCA Cloud Computing semester. The home nodes operate behind ISP CGNAT without public IP addresses or open inbound ports. Centralized cloud hosting of video binaries is economically and architecturally prohibitive.

## Decision
1. **Zero Cloud Video Retention**: Home edge nodes keep all media binaries. Marking media public syncs only a <= 3 KB JSON manifest.
2. **Outbound Reverse Tunnels**: Home edge nodes use FRP (Fast Reverse Proxy) to establish outbound TLS connections to a central relay.
3. **Pre-encoded 4-second HLS**: Transcoding is done ahead of time on the edge node (NVENC with automatic CPU fallback), chunked into strict 4-second segments.
4. **Origin Shield**: The cloud acts as a temporary pull-through cache (S3-compatible, 24h LRU purge) so hot segments are uploaded from the home node at most once.
5. **Strict Tech Stack**:
   - Backend: Go (stdlib `net/http` + `log/slog`).
   - Frontend: Next.js App Router (TypeScript strict, SSR default, `hls.js` client player).
   - Redis: Node heartbeat tracking only (`heartbeat:{node_id}` TTL).
   - Cloud DB: PostgreSQL via `pgx` and `sqlc`.
   - Edge DB: Pure-Go SQLite (`modernc.org/sqlite`).
6. **Admission Control**: Cap of 10 concurrent streams per home node returning HTTP 429 Retry-After when saturated.

## Consequences
- **Positive**: Near-zero cloud storage costs; bypasses home ISP CGNAT; resilient offline detection; adheres 100% to FOSS dependencies.
- **Negative**: Origin Shield provides cache benefits primarily when multiple users watch the same content; watching 10 distinct titles simultaneously stresses the home upstream bandwidth, requiring admission control and adaptive bitrate fallbacks.
