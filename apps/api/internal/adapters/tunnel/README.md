# Tunnel Puller Adapter

`apps/api/internal/adapters/tunnel` implements `ports.TunnelPuller`.

## Purpose
Provides the cloud-side pull mechanism to fetch uncached HLS playlists and media segments from edge nodes through the FRP reverse tunnel relay (Invariant 2, Invariant 7).

## Public Interface
- `Config`: Configuration (`BaseURL`, `URLPattern`, `Timeout`).
- `NewHTTPTunnelPuller(cfg Config) (*HTTPTunnelPuller, error)`: Standard constructor.
- `PullChunk(ctx context.Context, nodeID, mediaID, chunkPath string) (io.ReadCloser, string, int64, error)`: Fetches a chunk via HTTP.

## Testing
Run unit tests with mock HTTP test servers:
```bash
go test -v ./apps/api/internal/adapters/tunnel/...
```
